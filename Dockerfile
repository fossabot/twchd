FROM 127.0.0.1:1900/alpine:latest
WORKDIR /app
COPY bin/twchd /app
ENTRYPOINT [ "/app/twchd", "-config", "/app/config.yml" ]
