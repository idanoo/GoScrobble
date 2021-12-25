# Authentication

## v1/login

> To obtain a JWT token (Replace default credentials):

```shell
curl "https://goscrobble.com/api/v1/login" \
  -H "Content-Type: application/json" \
  --data '{"username":"abc","password":"def"}'
```

> Response:

```json
{
  "token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6ZmFsc2UsImVtYWlsIjoidGVzdEB0ZXN0LmNvbSIsImV4cCI6MTY0MTAwMjUxOSwiaWF0IjoxNjQwMzk3NzE5LCJtb2QiOnRydWUsInJlZnJlc2hfZXhwIjoxNjQxMDAyNTE5LCJyZWZyZXNoX3Rva2VuIjoiYWJjZGVmZ2hqaWtsbW5vcHF3ZXJ0dXZ3eHl6Iiwic3ViIjoiNDE4ZmRkZmItOGIxYi01MWFiLTliZHMtNGZnMDhjYTYzY2ZmIiwidXNlcm5hbWUiOiJ0ZXN0In0.fuPXjQ7IzNyttgIKpdS4-KBQ-QeHTl-BfgYkSnMCmpVrBunzMrSwr1RzxI7Xg2WWF-FHtW3Bnv9RpSqLDN4F2g"
}
```

This endpoint is used to authenticate and retrieve a JWT token.

### HTTP Request

`POST https://goscrobble.com/api/v1/login`

### Query Parameters

Parameter | Required | Description
--------- | ------- | -----------
username | true | Account username
password | true | Account password

## v1/register

> Create a new account (Replace default credentials):

```shell
curl "https://goscrobble.com/api/v1/register" \
  -H "Content-Type: application/json" \
  --data '{"email": "test@test.com", "username":"abc","password":"def"}'
```

> Response:

```json
{
  "message": "User created succesfully. You may now login"
}
```

If the server has REGISTRATION_ENABLED=true set, this endpoint will allow you to create a new account. Password must be at least 8 characters long.

### HTTP Request

`POST https://goscrobble.com/api/v1/register`

### Query Parameters

Parameter | Required | Description
--------- | ------- | -----------
email | true | Account email
username | true | Account username
password | true | Account password

## v1/sendreset

> Trigger a password reset email:

```shell
curl "https://goscrobble.com/api/v1/sendreset" \
  -H "Content-Type: application/json" \
  --data '{"email":"test@test.com"}'
```

> Response:

```json
{
  "message": "Password reset email sent"
}
```

This endpoint triggers a password reset email to be sent to the email on an account.

### HTTP Request

`POST https://goscrobble.com/api/v1/sendreset`

### Query Parameters

Parameter | Required | Description
--------- | ------- | -----------
email | true | Account username

## v1/resetpassword

> Trigger a password reset email:

```shell
curl "https://goscrobble.com/api/v1/resetpassword" \
  -H "Content-Type: application/json" \
  --data '{"token":"abcdefghijklmnopqrstuvwxyz", "password": "Hunter1"}'
```

> Response:

```json
{
  "message": "Password updated successfully!"
}
```

This endpoint confirms a password reset with a valid hash from the reset email. You must call v1/sendreset first and obtain the hash.

### HTTP Request

`POST https://goscrobble.com/api/v1/resetpassword`

### Query Parameters

Parameter | Required | Description
--------- | ------- | -----------
token | true | Reset token from the password reset email
password | true | New account password

