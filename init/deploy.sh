#!/bin/bash
# Easy deploy script..

echo 'Fetching latest git commit'
cd /var/www/goscrobble-api
git pull

echo 'Building backend'
go build -o goscrobble cmd/goscrobble/*.go


echo 'Fetching lastest frontend commit'
cd /var/www/goscrobble-web
git pull

echo 'Installing frontend packages'
npm install --production

echo 'Building frontend'
npm run build --env production

cd ..
echo 'Restarting Go service'
systemctl restart goscrobble.service