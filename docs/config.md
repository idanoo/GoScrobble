## Timezones
GoScrobble runs as UTC and connects to MySQL as UTC. All timezone handling is done in the frontend.

## FRONTEND VARS
These are stored in `web/.env.production` and `web/.env.development`

    REACT_APP_API_URL=https://goscrobble.com // Sets API URL


## BACKEND VARS
    MYSQL_HOST= // MySQL Server
    MYSQL_USER= // MySQL User
    MYSQL_PASS= // MySQL Password
    MYSQL_DB= // MySQL Database

    REDIS_HOST=127.0.0.1 // Redis host
    REDIS_PORT= // Redis port (defaults 6379)
    REDIS_DB=4 // Redis DB
    REDIS_PREFIX="gs:" // Redis key prefix
    REDIS_AUTH="" // Redis password

    JWT_SECRET= // 32+ Char JWT secret
    JWT_EXPIRY=1800 // JWT expiry in seconds
    REFRESH_EXPIRY=604800 // Refresh token expiry

    REVERSE_PROXIES=127.0.0.1 // Comma separated list of servers to ignore for IP logs
    PORT=42069 // Server port

    SENDGRID_API_KEY= // API KEY
    MAIL_FROM_ADDRESS= // FROM email
    MAIL_FROM_NAME= // FROM name

    DEV_MODE=false // true|false - Defaults false
    GOSCROBBLE_DOMAIN="" // Full domain for email links (https://goscrobble.com))
    
    DATA_DIRECTORY="/var/www/goscrobble-data"
    FRONTEND_DIRECTORY="/var/www/goscrobble-web"
    API_DOCS_DIRECTORY="/var/www/goscrobble-api/docs/api/build"
