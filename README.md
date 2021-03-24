# go-scrobble

Golang based music scrobbler. MySQL backend.


## Setup MySQL

    create user 'goscrobble'@'%' identified by 'supersecurepass';
    create database goscrobble;
    grant all privileges on goscrobble.* to 'goscrobble'@'%';

## Local build/run
    cd web && npm install && npm start
    # In another terminal
    go mod tidy
    CGO_ENABLED=0 MYSQL_HOST=127.0.0.1 MYSQL_USER=goscrobble MYSQL_PASS=supersecurepass MYSQL_DB=goscrobble go run cmd/go-scrobble/*.go

Access test frontend through http://127.0.0.1:3000 API http://127.0.0.1:42069/api/v1

## Run local tests
    go test test/*

## Prod deployment
We need to build NPM package, and then ship web/build with the binary.

    cd web npm install --production && npm run build
    go build -o goscrobble cmd/go-scrobble/*.go
    ./goscrobble