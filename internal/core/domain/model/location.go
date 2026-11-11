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
	x int
	y int
}

func NewLocation(x int, y int) (Location, error) {
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

func (l Location) X() int {
	return l.x
}

func (l Location) Y() int {
	return l.y
}

func (l Location) Equal(other Location) bool {
	return l == other
}

func (l Location) Distance(other Location) int {
	distanceX := math.Abs(float64(l.x) - float64(other.x))
	distanceY := math.Abs(float64(l.y) - float64(other.y))

	return int(distanceX + distanceY)
}
