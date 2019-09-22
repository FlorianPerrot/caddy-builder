package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

var (
	envVersion          = os.Getenv("VERSION")
	envDisableTelemetry = os.Getenv("DISABLE_TELEMETRY")
	envPlugins          = os.Getenv("PLUGINS")
)

var (
	output      = flag.String("o", "./caddy", "Output caddy file")
	caddyLib    = GoLib{"github.com/caddyserver/caddy", envVersion}
	buildDir, _ = ioutil.TempDir("", "template")
	tpl, _      = template.New("caddy.template").Parse(`
		package main
		
		import (
			"{{.CaddyLib.Repository}}/caddy/caddymain"
		
			{{- range .Plugins}}
				_ "{{.Repository}}"
			{{end}}
		)
		
		func main() {
			caddymain.EnableTelemetry = {{if .Telemetry}}true{{else}}false{{end}}
			caddymain.Run()
		}
	`)
)

func main() {
	flag.Parse()

	if !caddyLib.isValid() {
		log.Fatal("invalid caddy version")
	}

	plugins, err := searchPlugins(envPlugins)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Generate caddymain with plugins imports
	caddyGo, _ := os.Create(filepath.Join(buildDir, "caddy.go"))
	_ = tpl.Execute(caddyGo, TemplateParameters{plugins, envDisableTelemetry == "true", caddyLib})

	runGoCommand("mod", "init", "caddy")

	// Set Caddy and Plugins version
	if !caddyLib.isLatest() {
		runGoCommand("mod", "edit", "-require", fmt.Sprintf("%s@%s", caddyLib.Repository, caddyLib.Version))
	}
	for _, plugin := range plugins {
		if plugin.Version != "" {
			runGoCommand("mod", "edit", "-require", fmt.Sprintf("%s@%s", plugin.Repository, plugin.Version))
		}
	}

	runGoCommand("get")
	runGoCommand("build", "-ldflags", `-extldflags "-static"`, "-o", "caddy")

	// Copy caddy build on current directory
	caddyBuild, _ := os.Open(filepath.Join(buildDir, "caddy"))
	caddyFinal := outputCaddyFile()
	_, err = io.Copy(caddyFinal, caddyBuild)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Caddy build with success")
}

func searchPlugins(rawPlugins string) ([]GoLib, error) {
	var plugins []GoLib

	if rawPlugins != "" {
		for _, rawPlugin := range strings.Split(rawPlugins, ",") {
			repositoryAndVersion := strings.Split(strings.TrimSpace(rawPlugin), "@")

			version := ""
			if len(repositoryAndVersion) == 2 {
				version = repositoryAndVersion[1]
			}

			plugin := GoLib{repositoryAndVersion[0], version}

			if !plugin.isValid() {
				return nil, fmt.Errorf("invalid plugin given : %s with version %s", plugin.Repository, plugin.Version)
			}

			plugins = append(plugins, plugin)
		}
	}

	return plugins, nil
}

func runGoCommand(args ...string) {
	fmt.Printf("Run: go %s\n", strings.Join(args, " "))
	cmd := exec.Command("go", args...)
	cmd.Dir = buildDir

	out, err := cmd.CombinedOutput()
	fmt.Printf("%s\n", out)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func outputCaddyFile() *os.File {
	if dir := filepath.Dir(*output); dir != "." {
		_, err := os.Open(dir)

		if os.IsNotExist(err) {
			_ = os.MkdirAll(dir, 0755)
		}
	}

	caddy, err := os.Create(*output)
	if err != nil {
		log.Fatal(err.Error())
	}

	_ = caddy.Chmod(0777)

	return caddy
}

type TemplateParameters struct {
	Plugins []GoLib
	Telemetry bool
	CaddyLib GoLib
}

type GoLib struct {
	Repository string
	Version string
}

func (p *GoLib) isValid() bool {
	validVersion, _ := regexp.Match(`^v\d\.\d\.\d$`, []byte(p.Version))
	return p.isLatest() || validVersion
}

func (p *GoLib) isLatest() bool {
	return p.Version == ""
}
