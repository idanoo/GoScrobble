# go-scrobble

Golang based music scrobbler.

Stack: Go 1.16+, Node 15+, React 17+, Postgresql 14.0+, Redis

There are prebuilt binaries/packages available.

Copy .env.example to .env and set variables. You can use https://www.grc.com/passwords.htm to generate a JWT_SECRET.

## More documentation
[Changelog](docs/changelog.md)

[Environment Variables](docs/config.md)

## Local development with docker
This assumes you have goscrobble-api and goscrobble-web cloned in the same folder.

    cp .env.development .env
    docker-compose up -d

Access API @ http://127.0.0.1:42069/api/v1
Access frontend @ http://127.0.0.1:3000
pgAdmin @ http://127.0.0.1:5050 (admin@admin.com / root)

## Prod deployment
    cp .env.production .env # Fill in the blanks
    go build -o goscrobble cmd/go-scrobble/*.go
    ./goscrobble

## Build API docs
    cd docs/api && docker run --rm --name slate -v $(pwd)/build:/srv/slate/build -v $(pwd)/source:/srv/slate/source slatedocs/slate build

## Test API Docs
   cd docs/api && docker run --rm --name slate -p 4567:4567 -v $(pwd)/source:/srv/slate/source slatedocs/slate serve

## Support development!
Feel free to support hosting and my coffee addiction https://liberapay.com/idanoo