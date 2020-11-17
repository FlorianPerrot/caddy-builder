# Caddy builder

`caddy-builder` generate caddy executable with your configuration.

Chose caddy version with `VERSION=v1.0.3`, chose to disable telemetry with `DISABLE_TELEMETRY=true` and chose your plugins with `PLUGINS=github.com/xuqingfeng/caddy-rate-limit@v1.6.3`.

If caddy or plugin version is omit `latest` is selected.

## Example

Simple run `caddy-builder`, but useless without any configurations just download [here](https://caddyserver.com/download).

### First case in Command line
```bash
$ go install github.com/FlorianPerrot/caddy-builder
$ VERSION=v1.0.3 DISABLE_TELEMETRY=true PLUGINS=github.com/nicolasazrak/caddy-cache,github.com/xuqingfeng/caddy-rate-limit@v1.6.3 caddy-builder
```

### Second case on Dockerfile

cf: [Dockerhub](https://cloud.docker.com/repository/docker/florianperrot/caddy-builder)

```Dockerfile
FROM florianperrot/caddy-builder as builder

ENV VERSION="v1.0.3"
ENV PLUGINS="github.com/nicolasazrak/caddy-cache,github.com/xuqingfeng/caddy-rate-limit"
ENV DISABLE_TELEMETRY="false"

RUN CGO_ENABLED=0 caddy-builder -o /install/caddy

#
# Final stage
#
FROM alpine:3.10

# install caddy
COPY --from=builder /install/caddy /usr/bin/caddy

RUN caddy -version
RUN caddy -plugins

CMD ["caddy"]
```