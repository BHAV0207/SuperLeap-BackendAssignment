package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/database"
	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/models"
)

func cacheKey(id string) string {
	return fmt.Sprintf("lead:%s", id)
}

func (s *LeadService) getCachedLead(id string) (*models.Lead, error) {

	if s.redis == nil {
		return nil, nil
	}

	val, err := s.redis.Get(database.Ctx, cacheKey(id)).Result()

	if err != nil {
		return nil, err
	}

	var lead models.Lead

	err = json.Unmarshal([]byte(val), &lead)

	if err != nil {
		return nil, err
	}

	return &lead, nil
}

func (s *LeadService) setCachedLead(lead *models.Lead) {

	if s.redis == nil {
		return
	}

	data, err := json.Marshal(lead)
	if err != nil {
		return
	}

	s.redis.Set(
		database.Ctx,
		cacheKey(lead.ID.String()),
		data,
		10*time.Minute,
	)
}

func (s *LeadService) deleteCachedLead(id string) {

	if s.redis == nil {
		return
	}

	s.redis.Del(
		database.Ctx,
		cacheKey(id),
	)
}
