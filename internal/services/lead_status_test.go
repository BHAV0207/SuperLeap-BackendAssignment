package services

import (
	"testing"

	"github.com/BHAV0207/SuperLeap-BackendAssignment/internal/models"
)

func TestIsValidTransition(t *testing.T) {
	tests := []struct {
		from, to models.LeadStatus
		want     bool
	}{
		{models.StatusNew, models.StatusContacted, true},
		{models.StatusNew, models.StatusLost, true},
		{models.StatusNew, models.StatusQualified, false},
		{models.StatusNew, models.StatusConverted, false},
		{models.StatusContacted, models.StatusQualified, true},
		{models.StatusContacted, models.StatusNew, false},
		{models.StatusQualified, models.StatusConverted, true},
		{models.StatusQualified, models.StatusLost, true},
		{models.StatusConverted, models.StatusLost, false},
		{models.StatusLost, models.StatusNew, false},
		{models.LeadStatus("UNKNOWN"), models.StatusNew, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.from)+"->"+string(tt.to), func(t *testing.T) {
			if got := isValidTransition(tt.from, tt.to); got != tt.want {
				t.Errorf("isValidTransition(%s, %s) = %v, want %v", tt.from, tt.to, got, tt.want)
			}
		})
	}
}
