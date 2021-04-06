#!/bin/bash
# Easy deploy script..

echo 'Fetching latest git commit'
git pull

echo 'Building backend'
go build -o goscrobble cmd/go-scrobble/*.go

cd web
echo 'Installing frontend packages'
npm install --production

echo 'Building frontend'
npm run build --env production

cd ..
echo 'Restarting Go service'
systemctl restart goscrobble.service