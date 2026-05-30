package services

import "github.com/BHAV0207/SuperLeap-BackendAssignment/internal/models"


var validTransitions = map[models.LeadStatus][]models.LeadStatus{
	models.StatusNew: {
		models.StatusContacted,
		models.StatusLost,
	},

	models.StatusContacted: {
		models.StatusQualified,
		models.StatusLost,
	},

	models.StatusQualified: {
		models.StatusConverted,
		models.StatusLost,
	},

	models.StatusConverted: {},

	models.StatusLost: {},
}

func isValidTransition(current, next models.LeadStatus) bool {
	allowedStatuses, exists := validTransitions[current]

	if !exists {
		return false
	}

	for _, status := range allowedStatuses {
		if status == next {
			return true
		}
	}

	return false
}