FROM alpine:3.10
WORKDIR /app
COPY bin/twchd /app
ENTRYPOINT [ "/app/twchd", "-config", "/app/config.yml" ]
