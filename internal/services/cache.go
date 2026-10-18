package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/models"
	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/observability"
	"github.com/redis/go-redis/v9"
)

const cacheTTL = 10 * time.Minute

func cacheKey(id string) string {
	return fmt.Sprintf("lead:%s", id)
}

func (s *LeadService) getCachedLead(ctx context.Context, id string) (*models.Lead, error) {

	if s.redis == nil {
		return nil, nil
	}

	val, err := s.redis.Get(ctx, cacheKey(id)).Result()

	if errors.Is(err, redis.Nil) {
		observability.CacheRequestsTotal.WithLabelValues("miss").Inc()
		return nil, err
	}

	if err != nil {
		// redis down or slow, caller falls back to postgres
		observability.CacheRequestsTotal.WithLabelValues("error").Inc()
		slog.WarnContext(ctx, "cache read failed", "key", cacheKey(id), "error", err)
		return nil, err
	}

	var lead models.Lead

	err = json.Unmarshal([]byte(val), &lead)

	if err != nil {
		observability.CacheRequestsTotal.WithLabelValues("error").Inc()
		return nil, err
	}

	observability.CacheRequestsTotal.WithLabelValues("hit").Inc()
	return &lead, nil
}

func (s *LeadService) setCachedLead(ctx context.Context, lead *models.Lead) {

	if s.redis == nil {
		return
	}

	data, err := json.Marshal(lead)
	if err != nil {
		return
	}

	s.redis.Set(
		ctx,
		cacheKey(lead.ID.String()),
		data,
		cacheTTL,
	)
}

func (s *LeadService) deleteCachedLead(ctx context.Context, id string) {

	if s.redis == nil {
		return
	}

	s.redis.Del(
		ctx,
		cacheKey(id),
	)
}
