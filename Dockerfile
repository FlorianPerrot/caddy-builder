FROM golang:1.15

COPY builder.go .

RUN go build -o /usr/bin/caddy-builder
RUN chmod +x /usr/bin/caddy-builder

CMD [ "/usr/bin/caddy-builder" ]
