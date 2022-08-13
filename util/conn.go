package util

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v9"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB = dbConn()

func dbConn() *sql.DB {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	dbinfo := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s sslmode=disable", os.Getenv("USER"), os.Getenv("PWD"), os.Getenv("NAME"))

	db, err := sql.Open("postgres", dbinfo)

	if err != nil {
		panic(err)
	}

	return db
}

var Rdb = redisConn()

func redisConn() *redis.Client {
	c := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := c.Ping(context.Background()).Result()
	if err != nil {
		log.Println(err)
	}
	return c
}
