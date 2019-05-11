FROM alpine:latest
WORKDIR /app
COPY bin config/vanya83.yml /app/
ENTRYPOINT [ "/app/twchd", "-config", "/app/vanya83.yml", "-debug" ]