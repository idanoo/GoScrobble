## FRONTEND VARS
These are stored in `web/.env.production` and `web/.env.development`

    REACT_APP_REGISTRATION_DISABLED=true // Disables registration
    REACT_APP_API_URL=https://goscrobble.com // Sets API URL


## BACKEND VARS
    MYSQL_HOST= // MySQL Server
    MYSQL_USER= // MySQL User
    MYSQL_PASS= // MySQL Password
    MYSQL_DB= // MySQL Database

    REDIS_URL= // Redis host
    REDIS_DB=4 // Redis DB
    REDIS_PREFIX="gs:" // Redis key prefix
    REDIS_AUTH="" // Redis password

    TIMEZONE= // Used for MySQL connection

    JWT_SECRET= // 32+ Char JWT secret
    JWT_EXPIRY=86400 // JWT expiry

    REVERSE_PROXIES=127.0.0.1 // Comma separated list of servers to ignore for IP logs
    PORT=42069 // Server port
