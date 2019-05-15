FROM alpine:latest
WORKDIR /app
COPY bin/twchd /app
ENTRYPOINT [ "/app/twchd", "-config", "/app/config.yml" ]