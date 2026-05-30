package database

import (
	"log"

	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/models"
	"gorm.io/gorm"
)

func strPtr(s string) *string {
	return &s
}

func SeedData(db *gorm.DB) error {

	leads := []models.Lead{
		{
			Name:   "Bhavya Jain",
			Email:  "bhavya@example.com",
			Phone:  strPtr("+91-9999999999"),
			Source: strPtr("website"),
			Status: models.StatusNew,
		},
		{
			Name:   "Alice Johnson",
			Email:  "alice@example.com",
			Phone:  strPtr("+1-555-123456"),
			Source: strPtr("linkedin"),
			Status: models.StatusContacted,
		},
		{
			Name:   "Bob Smith",
			Email:  "bob@example.com",
			Phone:  strPtr("+44-777777777"),
			Source: strPtr("referral"),
			Status: models.StatusQualified,
		},
		{
			Name:   "Charlie Brown",
			Email:  "charlie@example.com",
			Phone:  strPtr("+91-8888888888"),
			Source: strPtr("ads"),
			Status: models.StatusConverted,
		},
		{
			Name:   "David Wilson",
			Email:  "david@example.com",
			Phone:  strPtr("+91-7777777777"),
			Source: strPtr("cold-email"),
			Status: models.StatusLost,
		},
	}

	for _, lead := range leads {

		var existing models.Lead

		err := db.Where(
			"email = ?",
			lead.Email,
		).First(&existing).Error

		if err == nil {
			log.Printf(
				"lead already exists: %s",
				lead.Email,
			)

			continue
		}

		err = db.Create(&lead).Error
		if err != nil {
			return err
		}

		log.Printf(
			"seeded lead: %s",
			lead.Email,
		)
	}

	log.Println("database seeding completed")

	return nil
}