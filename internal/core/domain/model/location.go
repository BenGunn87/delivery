package model

import "delivery/internal/pkg/errs"

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
	distanceX := l.x - other.x
	distanceY := l.y - other.y

	if distanceX < 0 {
		distanceX = -distanceX
	}
	if distanceY < 0 {
		distanceY = -distanceY
	}

	return distanceX + distanceY
}
