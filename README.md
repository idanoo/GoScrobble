# go-scrobble

Golang based music scrobbler. MySQL backend.


## Setup MySQL

    create user 'goscrobble'@'%' identified by 'supersecurepass';
    create database goscrobble;
    grant all privileges on goscrobble.* to 'goscrobble'@'%';

## Local build/run
    go mod tidy
    CGO_ENABLED=0 MYSQL_HOST=127.0.0.1 MYSQL_USER=goscrobble MYSQL_PASS=supersecurepass MYSQL_DB=goscrobble go run cmd/go-scrobble/*.go


## Run local tests
    go test test/*