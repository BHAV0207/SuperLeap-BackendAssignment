package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/config"
	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func ConnectRedis(cfg *config.AppConfig) *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf(
			"%s:%s",
			cfg.Redis.Host,
			cfg.Redis.Port,
		),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err := client.Ping(ctx).Err()
	if err != nil {

		log.Printf("redis unavailable, continuing without cache: %v", err)

		return nil // this here is intentional, we want to continue without Redis if it's not available
	}

	log.Println("redis connected successfully")

	return client
}
