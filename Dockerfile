FROM alpine:latest
WORKDIR /app
COPY bin settings.yml /app/
ENTRYPOINT [ "/app/twchd", "-config", "/app/settings.yml", "-debug" ]