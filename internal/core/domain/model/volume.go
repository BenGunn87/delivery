package model

import "delivery/internal/pkg/errs"

type Volume struct {
	volume uint16
}

func NewVolume(volume uint16) (Volume, error) {
	if volume == 0 {
		return Volume{}, errs.NewValueMustBeGreaterThan(
			"volume",
			volume,
			0,
		)
	}
	return Volume{volume}, nil
}

func (v Volume) Volume() uint16 {
	return v.volume
}

func (v Volume) Equal(other Volume) bool { return v.volume == other.volume }

func (v Volume) LessThen(other Volume) bool { return v.volume < other.volume }

func (v Volume) GreaterThan(other Volume) bool { return v.volume > other.volume }
