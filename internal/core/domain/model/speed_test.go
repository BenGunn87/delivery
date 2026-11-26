package model

import (
	"delivery/internal/pkg/errs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_SpeedBeCorrectWhenParamIsCorrectOnCreated(t *testing.T) {
	expectedS := 10

	speed, err := NewSpeed(expectedS)

	assert.NoError(t, err)
	assert.NotEmpty(t, speed)
	assert.Equal(t, expectedS, speed.Speed())
}

func Test_SpeedReturnErrorWhenParamIsInCorrectOnCreated(t *testing.T) {
	// Arrange
	tests := map[string]struct {
		s        int
		expected error
	}{
		"Speed 0": {
			s:        0,
			expected: errs.NewValueMustBeGreaterThan("speed", 0, 0),
		},
		"Speed -10": {
			s:        -10,
			expected: errs.NewValueMustBeGreaterThan("speed", -10, 0),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Act
			_, err := NewSpeed(test.s)

			// Assert
			if err.Error() != test.expected.Error() {
				t.Errorf("expected %v, got %v", test.expected, err)
			}
		})
	}
}
