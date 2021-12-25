# go-scrobble

Golang based music scrobbler.

Stack: Go 1.16+, Node 15+, React 17+, MySQL 8.0+, Redis

There are prebuilt binaries/packages available.

Copy .env.example to .env and set variables. You can use https://www.grc.com/passwords.htm to generate a JWT_SECRET.

## More documentation
[Changelog](docs/changelog.md)

[Environment Variables](docs/config.md)

## Setup MySQL
    create user 'goscrobble'@'%' identified by 'supersecurepass';
    create database goscrobble;
    grant all privileges on goscrobble.* to 'goscrobble'@'%';

## Local Development
    cp .env.example .env # Fill in the blanks
    go mod tidy
    CGO_ENABLED=0 go run cmd/go-scrobble/*.go


Access API @ http://127.0.0.1:42069/api/v1

## Prod deployment
    cp .env.example .env # Fill in the blanks
    go build -o goscrobble cmd/go-scrobble/*.go
    ./goscrobble

## Build API Docs
    cd docs/api && docker run --rm --name slate -v $(pwd)/build:/srv/slate/build -v $(pwd)/source:/srv/slate/source slatedocs/slate build


## Test API Docs
   cd docs/api && docker run --rm --name slate -p 4567:4567 -v $(pwd)/source:/srv/slate/source slatedocs/slate serve

## Support Development!
Feel free to support hosting and my coffee addiction https://liberapay.com/idanoo