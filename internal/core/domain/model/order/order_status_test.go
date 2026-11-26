package order

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_StatusShouldReturnCorrectName(t *testing.T) {
	// Arrange

	// Act

	// Assert
	assert.Equal(t, "created", StatusCreated.String())
	assert.Equal(t, "assigned", StatusAssigned.String())
	assert.Equal(t, "completed", StatusCompleted.String())
}

func Test_StatusShouldBeEqualWhenAllPropertiesEqual(t *testing.T) {
	// Arrange

	// Act

	// Assert
	assert.True(t, StatusCreated.Equal(StatusCreated))
	assert.True(t, StatusAssigned.Equal(StatusAssigned))
	assert.True(t, StatusCompleted.Equal(StatusCompleted))
}

func Test_StatusShouldBeNotEqualWhenAllPropertiesEqual(t *testing.T) {
	// Arrange

	// Act

	// Assert
	assert.False(t, StatusCompleted.Equal(StatusCreated))
	assert.False(t, StatusCreated.Equal(StatusAssigned))
	assert.False(t, StatusAssigned.Equal(StatusCompleted))
}
