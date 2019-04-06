FROM alpine:latest
WORKDIR /app
COPY bin settings.yml mapping.json /app/
CMD [ "/app/botbot.com", "-config", "settings.yml" ]