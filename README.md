# Caddy builder

`caddy-builder` generate caddy executable with your configuration.

Chose caddy version with `VERSION=v1.0.3`, chose to disable telemetry with `DISABLE_TELEMETRY=false` and chose your plugins with `PLUGINS=github.com/xuqingfeng/caddy-rate-limit@v1.6.3`. If caddy or plugin version is omit `latest` is selected.

## Example

Simple run `caddy-builder`, but useless without any configurations just download [here](https://caddyserver.com/download).

`VERSION=v1.0.3 DISABLE_TELEMETRY=true PLUGINS=github.com/nicolasazrak/caddy-cache,github.com/xuqingfeng/caddy-rate-limit@v1.6.3 caddy-builder`