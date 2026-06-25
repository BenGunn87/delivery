package model

import (
	"delivery/internal/pkg/errs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_LocationBeCorrectWhenParamsAreCorrectOnCreated(t *testing.T) {
	expectedX := uint8(1)
	expectedY := uint8(2)

	location, err := NewLocation(expectedX, expectedY)

	assert.NoError(t, err)
	assert.NotEmpty(t, location)
	assert.Equal(t, expectedX, location.X())
	assert.Equal(t, expectedY, location.Y())
}

func Test_LocationReturnErrorWhenParamsAreInCorrectOnCreated(t *testing.T) {
	// Arrange
	tests := map[string]struct {
		x        uint8
		y        uint8
		expected error
	}{
		"x 0": {
			x:        0,
			y:        1,
			expected: errs.NewMustBeBetween("x", 0, minCoordinateValue, maxCoordinateValue),
		},
		"x 250": {
			x:        250,
			y:        1,
			expected: errs.NewMustBeBetween("x", 250, minCoordinateValue, maxCoordinateValue),
		},
		"y 0": {
			x:        1,
			y:        0,
			expected: errs.NewMustBeBetween("y", 0, minCoordinateValue, maxCoordinateValue),
		},
		"y 250": {
			x:        1,
			y:        250,
			expected: errs.NewMustBeBetween("y", 250, minCoordinateValue, maxCoordinateValue),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Act
			_, err := NewLocation(test.x, test.y)

			// Assert
			if err.Error() != test.expected.Error() {
				t.Errorf("expected %v, got %v", test.expected, err)
			}
		})
	}
}
