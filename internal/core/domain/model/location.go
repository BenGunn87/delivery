package model

import (
	"delivery/internal/pkg/errs"
	"math"
)

const (
	minCoordinateValue = 1
	maxCoordinateValue = 10
)

type Location struct {
	x uint8
	y uint8
}

func NewLocation(x uint8, y uint8) (Location, error) {
	if x < minCoordinateValue || x > maxCoordinateValue {
		return Location{}, errs.NewMustBeBetween(
			"x",
			x,
			minCoordinateValue,
			maxCoordinateValue,
		)
	}

	if y < minCoordinateValue || y > maxCoordinateValue {
		return Location{}, errs.NewMustBeBetween(
			"y",
			y,
			minCoordinateValue,
			maxCoordinateValue,
		)
	}

	return Location{
		x: x,
		y: y,
	}, nil
}

func (l Location) X() uint8 {
	return l.x
}

func (l Location) Y() uint8 {
	return l.y
}

func (l Location) Equal(other Location) bool {
	return l == other
}

func (l Location) Distance(other Location) uint8 {
	distanceX := math.Abs(float64(l.x) - float64(other.x))
	distanceY := math.Abs(float64(l.y) - float64(other.y))

	return uint8(distanceX + distanceY)
}
