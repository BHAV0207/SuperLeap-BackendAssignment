package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/config"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

func ConnectRedis(cfg *config.AppConfig) *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf(
			"%s:%s",
			cfg.Redis.Host,
			cfg.Redis.Port,
		),
		// redis is only a cache. If it is down, a request should fall back
		// to postgres quickly instead of waiting on the default 5s timeouts.
		DialTimeout:  500 * time.Millisecond,
		ReadTimeout:  200 * time.Millisecond,
		WriteTimeout: 200 * time.Millisecond,
		MaxRetries:   -1,
		// go-redis retries a failed dial 5 times by default, that alone added
		// ~3.5s to every request while redis was down. Fail fast instead.
		DialerRetries: 1,
	})

	if err := redisotel.InstrumentTracing(client); err != nil {
		slog.Warn("could not enable redis tracing", "error", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	defer cancel()

	// Earlier this returned nil when redis was down, which turned the cache
	// off until the next restart. go-redis reconnects by itself, so we keep the
	// client and cache errors are treated as misses in the service layer.
	err := client.Ping(ctx).Err()
	if err != nil {
		slog.Warn("redis unavailable, continuing without cache until it comes back", "error", err)
	} else {
		slog.Info("redis connected successfully")
	}

	return client
}
