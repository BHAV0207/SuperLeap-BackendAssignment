package validators

import "slices"

func IsValidStatus(status string) bool {
	validStatuses := []string{"NEW", "CONTACTED", "QUALIFIED", "CONVERTED", "LOST"}
	return slices.Contains(validStatuses, status)
}
