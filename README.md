# go-scrobble

Golang based music scrobbler. MySQL 8.0+

Currently building on Node V15.X & Go V1.16.X

With a prebuilt binary - you will still need the migrations folder + web/build folder on prod.

Copy .env.example to .env and set variables. You can use https://www.grc.com/passwords.htm to generate a JWT_SECRET.


## Setup MySQL
    create user 'goscrobble'@'%' identified by 'supersecurepass';
    create database goscrobble;
    grant all privileges on goscrobble.* to 'goscrobble'@'%';

## Local build/run
    cp .env.example .env # Fill in the blanks
    cd web && npm install && REACT_APP_API_URL=http://127.0.0.1:42069 npm start
    # In another terminal
    go mod tidy
    CGO_ENABLED=0 go run cmd/go-scrobble/*.go

Access dev frontend @ http://127.0.0.1:3000 + API @ http://127.0.0.1:42069/api/v1

## Run local tests
    go test test/*

## Prod deployment
We need to build NPM package, and then ship web/build with the binary.
    cp .env.example .env # Fill in the blanks
    cd web npm install --production && REACT_APP_API_URL=https://goscrobble.com npm run build
    go build -o goscrobble cmd/go-scrobble/*.go
    ./goscrobble