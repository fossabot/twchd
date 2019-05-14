FROM alpine:latest
WORKDIR /app
COPY bin /app/
ENTRYPOINT [ "/app/twchd", "-config", "/app/config.yml", "-debug" ]