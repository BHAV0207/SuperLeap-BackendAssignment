package validators

import "slices"

func IsValidStatus(status string) bool {
	validStatuses := []string{"NEW", "CONTACTED", "QUALIFIED", "LOST"}
	return slices.Contains(validStatuses, status)
}
