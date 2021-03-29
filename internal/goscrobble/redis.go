package goscrobble

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

var redisDb *redis.Client
var redisPrefix string

var ctx = context.Background()

// InitRedis - Boot redis connection!
func InitRedis() {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisDatabase := os.Getenv("REDIS_DB")
	redisAuth := os.Getenv("REDIS_AUTH")
	redisPrefix = os.Getenv("REDIS_PREFIX")

	redisDbNum := 0
	if redisDatabase != "" {
		redisDbNum, _ = strconv.Atoi(redisDatabase)
	}

	// Create new connection
	redisDb = redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisAuth,
		DB:       redisDbNum,
	})

	// Lets just check it's active..
	err := redisDb.Set(ctx, "testSetKey", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	redisDb.Del(ctx, "testSetKey")
	fmt.Println("Redis connected")
}

func CloseRedisConn() {
	redisDb.Close()
}

// setRedis - Uses default 24 hour TTL
func setRedisVal(key string, val string) error {
	ttl := time.Hour * time.Duration(24)
	return setRedisValTtl(key, val, ttl)
}

// setRedisTtl - Allows custom TTL
func setRedisValTtl(key string, val string, ttl time.Duration) error {
	return redisDb.Set(ctx, redisPrefix+key, val, 0).Err()
}

// getRedisVal - Returns value if exists
func getRedisVal(key string) string {
	val, err := redisDb.Get(ctx, redisPrefix+key).Result()
	if err != nil {
		if err == redis.Nil {
			return ""
		}
		log.Printf("Failed to fetch redis key (%+v) Error: %+v", key, err)
	}

	return val
}
