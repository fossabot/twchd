# twchd
Twitch chat message grabber
# How to run
- cd $GOPATH/src
- git clone https://github.com/Aded175/twchd.git
- cd twchd
- dep ensure
- make generate
- make image
- docker-compose up -d
- docker run --rm --network=host -e TWITCH_ACCOUNT="XXX" -e TWITCH_OAUTH="oauth:YYY" -v $PWD/config/ZZZ.yml:/app/config.yml $USER/twchd:2.1-1-g90487dd -debug
