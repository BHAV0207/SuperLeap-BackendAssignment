package services

import "errors"

var (
	ErrInvalidStatusTransition = errors.New(
		"invalid status transition",
	)
)