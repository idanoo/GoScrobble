# Server Info
## v1/serverinfo
> Check server version and registration status:

```shell
curl "https://goscrobble.com/api/v1/serverinfo"
```

> Response:

```json
{
  "version":"0.1.1",
  "registration_enabled":"1"
}
```

This endpoint is used to get server API version and registration status.

### HTTP Request

`GET https://goscrobble.com/api/v1/serverstatus`