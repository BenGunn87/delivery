package model

import "delivery/internal/pkg/errs"

type Speed struct {
	speed int
}

func NewSpeed(speed int) (Speed, error) {
	if speed <= 0 {
		return Speed{}, errs.NewValueMustBeGreaterThan(
			"speed",
			speed,
			0,
		)
	}
	return Speed{speed: speed}, nil
}
func (s Speed) Speed() int {
	return s.speed
}

func (s Speed) Equal(other Speed) bool {
	return s.speed == other.speed
}
