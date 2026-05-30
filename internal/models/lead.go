package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LeadStatus string

const (
	StatusNew       LeadStatus = "NEW"
	StatusContacted LeadStatus = "CONTACTED"
	StatusQualified LeadStatus = "QUALIFIED"
	StatusConverted LeadStatus = "CONVERTED"
	StatusLost      LeadStatus = "LOST"
)

type Lead struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string    `gorm:"not null"`
	Email     string    `gorm:"not null"`
	Phone     *string
	Source    *string
	Status    LeadStatus `gorm:"type:varchar(20);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (l *Lead) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	if l.Status == "" {
		l.Status = StatusNew
	}
	return nil
}
