package repositories

import (
	"context"

	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LeadRepository struct {
	db *gorm.DB
}

func NewLeadRepository(db *gorm.DB) *LeadRepository {
	return &LeadRepository{db: db}
}

func (r *LeadRepository) Create(ctx context.Context, lead *models.Lead) error {
	return r.db.WithContext(ctx).Create(lead).Error
}

func (r *LeadRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Lead, error) {
	var lead models.Lead
	if err := r.db.WithContext(ctx).First(&lead, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &lead, nil
}

func (r *LeadRepository) GetAll(ctx context.Context, status *models.LeadStatus, email *string) ([]models.Lead, error) {
	var leads []models.Lead

	query := r.db.WithContext(ctx).Model(&models.Lead{})

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if email != nil {
		query = query.Where("email = ?", *email)
	}

	err := query.Find(&leads).Error
	if err != nil {
		return nil, err
	}

	return leads, nil
}

func (r *LeadRepository) Update(ctx context.Context, lead *models.Lead) error {
	return r.db.WithContext(ctx).Save(lead).Error
}

func (r *LeadRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.Lead{}, "id = ?", id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
