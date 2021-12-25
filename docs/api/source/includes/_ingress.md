# Ingress Endpoints

Spotify and Navidrome both work on a poll based system and do not require incoming webhooks/endpoints.

## POST v1/ingress/jellyfin

> Submit a scrobble:

```shell
curl "https://goscrobble.com/api/v1/ingress/jellyfin" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6ZmFsc2UsImVtYWlsIjoidGVzdEB0ZXN0LmNvbSIsImV4cCI6MTY0MTAwMjUxOSwiaWF0IjoxNjQwMzk3NzE5LCJtb2QiOnRydWUsI" \
  -X POST \
  --data '{"Name":"<Full JSON structure located in ingress_jellyfin.go>", "Album":"Best Album","Artist":"Best Artist"}'
```

> Response:

```json
{
  "message": "success"
}
```

Adds a scrobble from Jellyfin into the database. Please view `type JellyfinRequest struct` in ingress_jellyfin.go for the full request format. You need to install the webhook plugin for this to work.

<aside class="notice">
Light rate-limiting applies
</aside>

### HTTP Request

`POST https://goscrobble.com/api/v1/ingress/jellyfin`

### Query Parameters

Parameter | Required | Type | Description
--------- | ------- | ----------- | -----------
Album | true | string | Album title
Artist | true | string | Artist name
Name | true | string | Song title
ItemType | true | string | Will only scrobble type 'Audio'
ClientName | false | string | TBD
DeviceId | false | string | TBD
DeviceName | false | string | TBD
IsAutomated | false | bool | TBD
IsPaused | false | bool | TBD
ItemId | false | string | TBD
MediaSourceId | false | string | TBD
NotificationType | false | string | TBD
Overview | false | string | TBD
PlaybackPosition | false | string | TBD
PlaybackPositionTicks | false | int | TBD
Provider_musicbrainzalbum | false | string | TBD
Provider_musicbrainzalbumartist | false | string | TBD
Provider_musicbrainzartist | false | string | TBD
Provider_musicbrainzreleasegroup | false | string | TBD
Provider_musicbrainztrack | false | string | TBD
RunTime | false | string | TBD
RunTimeTicks | false | int | TBD
ServerId | false | string | TBD
ServerName | false | string | TBD
ServerUrl | false | string | TBD
ServerVersion | false | string | TBD
Timestamp | false | string | TBD
UserId | false | TBD
Username | false | TBD
UtcTimestamp | false | TBD
Year | false | int | TBD

## POST v1/ingress/multiscrobbler

> Submit a scrobble:

```shell
curl "https://goscrobble.com/api/v1/ingress/multiscrobbler" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6ZmFsc2UsImVtYWlsIjoidGVzdEB0ZXN0LmNvbSIsImV4cCI6MTY0MTAwMjUxOSwiaWF0IjoxNjQwMzk3NzE5LCJtb2QiOnRydWUsI" \
  -X POST \
  --data '{"Track":"<Full JSON structure located in ingress_multiscrobbler.go>", "Album":"Best Album","Artist":["Best Artist1", "Best Artist2"]}'
```

> Response:

```json
{
  "message": "success"
}
```

Adds a scrobble from MultiScrobbler into the database. Please view `type MultiScrobblerRequest struct` in ingress_multiscrobbler.go for the full request format.

<aside class="notice">
Light rate-limiting applies
</aside>

### HTTP Request

`POST https://goscrobble.com/api/v1/ingress/multiscrobbler`

### Query Parameters

Parameter | Required | Type | Description
--------- | ------- | ----------- | -----------
Album | true | string | Album title
Artist | true | array | Array of Artist names
Track | true | string | Song title
PlayedAt | true | int | Unix Timestamp of scrobble time
Duration | true | int | Song length in seconds