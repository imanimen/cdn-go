package utils

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"time"
)

var RedisClient *redis.Client

type Config struct {
	Addr     string
	Password string
	DB       int
}

func loadConfig() *Config {
	return &Config{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}
}

func connect() {
	config := loadConfig()

	// create new redis connection
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	// test the connection
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RedisClient.Ping().Err(); err != nil {
		log.Fatalf("Error connecting to redis: %v", err)
	}

	fmt.Println("Successfully connected to redis")
}

func disconnect() {
	if err := RedisClient.Close(); err != nil {
		log.Fatalf("Error disconnecting from redis: %v", err)
	} else {
		fmt.Println("Successfully disconnected from redis")
	}
}
