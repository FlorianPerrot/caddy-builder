package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestSearchPlugins(t *testing.T) {
	plugins, err := searchPlugins("github.com/nicolasazrak/caddy-cache,github.com/xuqingfeng/caddy-rate-limit@v1.6.3")
	if err != nil {
		t.Error(err.Error())
	}

	if len(plugins) != 2 {
		t.Error("Invalid plugins number found")
	}

	if plugins[0].Repository != "github.com/nicolasazrak/caddy-cache" || plugins[1].Repository != "github.com/xuqingfeng/caddy-rate-limit" {
		t.Error("Invalid plugin repository find")
	}

	if plugins[0].Version != "" || plugins[1].Version != "v1.6.3" {
		t.Error("Invalid plugin version find")
	}
}

func TestTemplate(t *testing.T) {
	buff := bytes.NewBuffer([]byte{})
	err := tpl.Execute(buff, TemplateParameters{[]GoLib{{"github.com/test", ""}}, true, GoLib{"github.com/caddy", ""}})
	if err != nil {
		t.Error(err.Error())
	}

	tplExecuted := buff.String()

	if !strings.Contains(tplExecuted, `"github.com/caddy/caddy/caddymain"`) {
		t.Error("Invalid caddy import")
	}
	if !strings.Contains(tplExecuted, `_ "github.com/test"`) {
		t.Error("Invalid plugin import")
	}
	if !strings.Contains(tplExecuted, "EnableTelemetry = true") {
		t.Error("Invalid telemetry configuration")
	}
}

func TestGoLib(t *testing.T) {
	invalidLib := GoLib{"github.com/test", "invalid"}
	if invalidLib.isValid() {
		t.Error("Golib should not be valid")
	}

	invalidLib.Version = "v1.2.3"
	if !invalidLib.isValid() {
		t.Error("Golib should be valid")
	}

	if invalidLib.isLatest() {
		t.Error("Golib should not be latest")
	}

	invalidLib.Version = ""
	if !invalidLib.isLatest() {
		t.Error("Golib should be latest")
	}
}