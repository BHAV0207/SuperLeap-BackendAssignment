package repositories

import (
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

func (r *LeadRepository) Create(lead *models.Lead) error {
	return r.db.Create(lead).Error
}

func (r *LeadRepository) GetByID(id uuid.UUID) (*models.Lead, error) {
	var lead models.Lead
	if err := r.db.First(&lead, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &lead, nil
}

func (r *LeadRepository) GetAll(status *models.LeadStatus) ([]models.Lead, error) {
	var leads []models.Lead

	query := r.db.Model(&models.Lead{})

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	err := query.Find(&leads).Error
	if err != nil {
		return nil, err
	}

	return leads, nil
}

func (r *LeadRepository) Update(lead *models.Lead) error {
	return r.db.Save(lead).Error
}

func (r *LeadRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Lead{}, "id = ?", id).Error
}
