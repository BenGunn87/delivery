package courier

import (
	"delivery/internal/core/domain/model"
	"delivery/internal/pkg/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_CourierBeCorrectWhenParamsAreCorrectOnCreate(t *testing.T) {
	// Arrange
	name := "Test"
	speed := 2
	location, _ := model.NewLocation(1, 1)

	// Act
	c, err := NewCourier(name, speed, location)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, name, c.Name())
	assert.Equal(t, speed, c.Speed().Speed())
	assert.Equal(t, location, c.Location())
	assert.Equal(t, uint16(10), c.StoragePlaces()[0].totalVolume.Volume())
}

func Test_CourierReturnErrorWhenNameIsEmptyOnCreated(t *testing.T) {
	// Arrange
	location, _ := model.NewLocation(1, 1)
	expected := errs.NewValueIsRequired("name")
	// Act
	_, err := NewCourier("", 2, location)

	// Assert

	if err.Error() != expected.Error() {
		t.Errorf("expected %v, got %v", expected, err)
	}
}

func Test_CourierReturnErrorWhenSpeedIsZeroOnCreated(t *testing.T) {
	// Arrange
	location, _ := model.NewLocation(1, 1)
	expected := errs.NewValueMustBeGreaterThan("speed", 0, 0)
	// Act
	_, err := NewCourier("Test", 0, location)

	// Assert
	if err.Error() != expected.Error() {
		t.Errorf("expected %v, got %v", expected, err)
	}
}

func Test_StoragePlaceAddedWhenParamsAreCorrectOnAddStoragePlace(t *testing.T) {
	// Arrange
	location, _ := model.NewLocation(1, 1)
	c, _ := NewCourier("Test", 1, location)
	spName := "backpack"
	spValue, _ := model.NewVolume(30)

	//Act
	err := c.AddStoragePlace(spName, spValue)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, spValue.Volume(), c.StoragePlaces()[1].totalVolume.Volume())
	assert.Equal(t, spName, c.StoragePlaces()[1].Name())
}

func Test_CourierCanTakeOrderWhenThereIsEmptyStoragePlaceWithSufficientVolume(t *testing.T) {
	// Arrange
	oId := uuid.MustParse("ec85ceee-f186-4e9c-a4dd-2929e69e586c")
	volume, _ := model.NewVolume(uint16(2))
	location, _ := model.NewLocation(1, 1)
	c, _ := NewCourier("Test", 1, location)

	// Act
	err := c.TakeOrder(oId, volume)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, c.StoragePlaces()[0].orderId, oId)
}

func Test_CourierCannotTakeOrderWhenNoEmptyStoragePlace(t *testing.T) {
	// Arrange
	id1 := uuid.MustParse("ec85ceee-f186-4e9c-a4dd-2929e69e586c")
	id2 := uuid.MustParse("ec85ceee-f186-4e9c-a4dd-2929e69e500c")
	volume, _ := model.NewVolume(uint16(2))
	location, _ := model.NewLocation(1, 1)
	c, _ := NewCourier("Test", 1, location)

	// Act
	_ = c.TakeOrder(id1, volume)
	err := c.TakeOrder(id2, volume)

	// Assert
	if err.Error() != ErrOrderCannotBeTaken.Error() {
		t.Errorf("expected %v, got %v", ErrOrderCannotBeTaken, err)
	}
}

func Test_CourierCannotTakeOrderWhenThereIsEmptyStoragePlaceWithSmallerVolume(t *testing.T) {
	// Arrange
	id := uuid.MustParse("ec85ceee-f186-4e9c-a4dd-2929e69e586c")
	volume, _ := model.NewVolume(uint16(15))
	location, _ := model.NewLocation(1, 1)
	c, _ := NewCourier("Test", 1, location)

	// Act
	err := c.TakeOrder(id, volume)

	// Assert
	if err.Error() != ErrOrderCannotBeTaken.Error() {
		t.Errorf("expected %v, got %v", ErrOrderCannotBeTaken, err)
	}
}

func Test_CourierCompleteOrderWhenOrderAssignedCourier(t *testing.T) {
	// Arrange
	id := uuid.MustParse("ec85ceee-f186-4e9c-a4dd-2929e69e586c")
	volume, _ := model.NewVolume(uint16(2))
	location, _ := model.NewLocation(1, 1)
	c, _ := NewCourier("Test", 1, location)

	// Act
	_ = c.TakeOrder(id, volume)
	err := c.CompleteOrder(id)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, c.StoragePlaces()[0].IsEmpty(), true)
}

func Test_NumberOfStepsCalculatedCorrectlyOnGetStepCount(t *testing.T) {
	// Arrange
	location, _ := model.NewLocation(1, 1)
	c, _ := NewCourier("Test", 2, location)

	tests := map[string]struct {
		x        int
		y        int
		expected int
	}{
		"target 1 1": {
			x:        1,
			y:        1,
			expected: 0,
		},
		"target 1 5": {
			x:        1,
			y:        5,
			expected: 2,
		},
		"target 5 4": {
			x:        5,
			y:        4,
			expected: 4,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			target, _ := model.NewLocation(test.x, test.y)
			// Act
			steps := c.GetStepCount(target)
			// Assert
			assert.Equal(t, test.expected, steps)
		})
	}
}

func Test_CourierMoveCorrectlyOnMove(t *testing.T) {
	// Arrange
	location, _ := model.NewLocation(5, 5)

	tests := map[string]struct {
		x    int
		y    int
		expX int
		expY int
	}{
		"target 5 5": {
			x:    5,
			y:    5,
			expX: 5,
			expY: 5,
		},
		"target 5 10": {
			x:    5,
			y:    10,
			expX: 5,
			expY: 7,
		},
		"target 10 5": {
			x:    10,
			y:    5,
			expX: 7,
			expY: 5,
		},
		"target 6 6": {
			x:    6,
			y:    6,
			expX: 6,
			expY: 6,
		},
		"target 5 4": {
			x:    5,
			y:    4,
			expX: 5,
			expY: 4,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			c, _ := NewCourier("Test", 2, location)
			target, _ := model.NewLocation(test.x, test.y)
			expLocation, _ := model.NewLocation(test.expX, test.expY)
			// Act
			err := c.Move(target)
			// Assert
			assert.NoError(t, err)
			assert.Equal(t, c.Location().Equal(expLocation), true)
		})
	}
}
