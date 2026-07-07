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
	assert.Equal(t, order.Id(), id)
	assert.Equal(t, order.Volume(), volume)
}

func Test_OrderReturnErrorWhenOrderIdIsNil(t *testing.T) {
	volume, _ := model.NewVolume(uint16(10))
	expected := errs.NewValueIsRequired("id")

	_, err := NewOrder(uuid.Nil, volume)

	if err.Error() != expected.Error() {
		t.Errorf("expected %v, got %v", expected, err)
	}
}
