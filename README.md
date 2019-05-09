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
- docker run --network=host -e TWITCH_ACCOUNT="XXX" -e TWITCH_OAUTH="oauth:YYY" $USER/twchd:2.0
