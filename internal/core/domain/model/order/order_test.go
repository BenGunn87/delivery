package order

import (
	"delivery/internal/core/domain/model"
	"delivery/internal/pkg/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_OrderBeCorrectWhenParamsAreCorrectOnCreate(t *testing.T) {
	id := uuid.MustParse("ec85ceee-f186-4e9c-a4dd-2929e69e586c")
	volume, _ := model.NewVolume(uint16(10))

	order, err := NewOrder(id, volume)
	assert.NoError(t, err)
	assert.NotEmpty(t, order)
	assert.Equal(t, order.ID(), id)
	assert.Equal(t, order.Volume(), volume)
	assert.Equal(t, order.Status(), StatusCreated)
	assert.Empty(t, order.CourierId())
}

func Test_OrderReturnErrorWhenOrderIdIsNil(t *testing.T) {
	volume, _ := model.NewVolume(uint16(10))
	expected := errs.NewValueIsRequired("id")

	_, err := NewOrder(uuid.Nil, volume)

	if err.Error() != expected.Error() {
		t.Errorf("expected %v, got %v", expected, err)
	}
}

func Test_OrderCanBeAssignedWhenItCreated(t *testing.T) {
	// Arrange
	id := uuid.MustParse("ec85ceee-f186-4e9c-a4dd-2929e69e586c")
	courierId := uuid.MustParse("ec58ceee-f186-4e9c-a4dd-2929e69e586c")
	volume, _ := model.NewVolume(uint16(10))
	order, _ := NewOrder(id, volume)

	// Act
	err := order.MarkAsAssigned(courierId)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, *order.CourierId(), courierId)
	assert.Equal(t, order.Status(), StatusAssigned)
}

func Test_OrderAssignedReturnErrorWhenOrderStatusIsWrong(t *testing.T) {
	// Arrange
	id := uuid.MustParse("ec85ceee-f186-4e9c-a4dd-2929e69e586c")
	courierId := uuid.MustParse("ec58ceee-f186-4e9c-a4dd-2929e69e586c")
	volume, _ := model.NewVolume(uint16(10))
	order, _ := NewOrder(id, volume)
	_ = order.MarkAsAssigned(courierId)

	// Act
	err := order.MarkAsAssigned(courierId)

	// Assert
	if err.Error() != ErrStatusIsNotCreated.Error() {
		t.Errorf("expected %v, got %v", ErrStatusIsNotCreated, err)
	}
}

func Test_OrderCanBeCompletedWhenItAssigned(t *testing.T) {
	// Arrange
	id := uuid.MustParse("ec85ceee-f186-4e9c-a4dd-2929e69e586c")
	courierId := uuid.MustParse("ec58ceee-f186-4e9c-a4dd-2929e69e586c")
	volume, _ := model.NewVolume(uint16(10))
	order, _ := NewOrder(id, volume)
	_ = order.MarkAsAssigned(courierId)

	// Act
	err := order.MarkAsCompleted()

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, order.CourierId())
	assert.Equal(t, order.Status(), StatusCompleted)
}

func Test_OrderCompletedReturnErrorWhenOrderStatusIsWrong(t *testing.T) {
	// Arrange
	id := uuid.MustParse("ec85ceee-f186-4e9c-a4dd-2929e69e586c")
	volume, _ := model.NewVolume(uint16(10))
	order, _ := NewOrder(id, volume)

	// Act
	err := order.MarkAsCompleted()

	// Assert
	if err.Error() != ErrStatusIsNotAssigned.Error() {
		t.Errorf("expected %v, got %v", ErrStatusIsNotAssigned, err)
	}
}
