# Administration
<aside class="warning">
An admin JWT token is required for the below endpoints
</aside>

## GET v1/config

> Fetch configuration values:

```shell
curl "https://goscrobble.com/api/v1/config" \
  -H "Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6ZmFsc2UsImVtYWlsIjoidGVzdEB0ZXN0LmNvbSIsImV4cCI6MTY0MTAwMjUxOSwiaWF0IjoxNjQwMzk3NzE5LCJtb2QiOnRydWUsI"
```

> Response:

```json
{
  "SPOTIFY_API_ID": "abc",
  "SPOTIFY_API_SECRET": "def",
  "REGISTRATION_ENABLED": true
}
```

Fetch instance configuration values.

<aside class="notice">
Standard rate-limiting applies
</aside>

### HTTP Request

`GET https://goscrobble.com/api/v1/config`


## POST v1/config

> Update a configuration parameter:

```shell
curl "https://goscrobble.com/api/v1/config" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6ZmFsc2UsImVtYWlsIjoidGVzdEB0ZXN0LmNvbSIsImV4cCI6MTY0MTAwMjUxOSwiaWF0IjoxNjQwMzk3NzE5LCJtb2QiOnRydWUsI" \
  -X POST \
  --data '{"SPOTIFY_API_ID":"notarealapikey","SPOTIFY_API_SECRET":"notarealsecret"}'
```

> Response:

```json
{
  "message": "Config updated successfully"
}
```

Updates instance configuration values.

<aside class="notice">
Standard rate-limiting applies
</aside>

### HTTP Request

`POST https://goscrobble.com/api/v1/config`

### Query Parameters

Parameter | Required | Description
--------- | ------- | -----------
key | true | Array of key/value pairs to update. See example.