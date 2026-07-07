package courier

import (
	"delivery/internal/core/domain/model"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func standardOrderWithVolume10() *order.Order {
	volume, _ := model.NewVolume(10)
	standardOrder, _ := order.NewOrder(
		uuid.MustParse("67dda004-f324-4868-af9a-88f77d1d28fc"),
		volume,
	)

	return standardOrder
}

func Test_StoragePlaceBeCorrectWhenParamsAreCorrectOnCreated(t *testing.T) {
	expectedOrder := standardOrderWithVolume10()
	expectedOrderId := expectedOrder.Id()
	volume, _ := model.NewVolume(uint16(20))

	name := "BackPack"
	storage, err := NewStoragePlace(name, volume, expectedOrder)

	assert.NoError(t, err)
	assert.NotEmpty(t, storage)
	assert.Equal(t, name, storage.Name())
	assert.Equal(t, volume, storage.TotalVolume())
	assert.Equal(t, &(expectedOrderId), storage.OrderId())
}

func Test_StoragePlaceBeCorrectEmptyWhenOrderNotSpecifiedOnCreated(t *testing.T) {
	volume, _ := model.NewVolume(uint16(20))

	name := "BackPack"
	storage, err := NewStoragePlace(name, volume, nil)

	assert.NoError(t, err)
	assert.NotEmpty(t, storage)
	assert.Equal(t, name, storage.Name())
	assert.Equal(t, volume, storage.TotalVolume())
	assert.Equal(t, true, storage.IsEmpty())
}

func Test_StoragePlaceReturnErrorWhenParamsAreNotCorrectOnCreated(t *testing.T) {
	expectedOrder := standardOrderWithVolume10()

	tests := map[string]struct {
		name        string
		volume      uint16
		order       *order.Order
		expectedErr error
	}{
		"wrong_name": {
			name:        "",
			order:       expectedOrder,
			volume:      20,
			expectedErr: errs.NewValueIsRequired("name"),
		},
		"order_volume_is_too_large": {
			name:        "BackPack",
			order:       expectedOrder,
			volume:      5,
			expectedErr: ErrOrderVolumeIsTooLarge,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			volume, _ := model.NewVolume(test.volume)

			_, err := NewStoragePlace(test.name, volume, test.order)
			if err.Error() != test.expectedErr.Error() {
				t.Errorf("expected %v, got %v", test.expectedErr, err)
			}
		})
	}
}

func Test_StoragePlaceShouldBeEqualWhenIdIsSame(t *testing.T) {
	volume, _ := model.NewVolume(10)
	storage, _ := NewStoragePlace("BackPack", volume, nil)

	isEqual := storage.Equal(storage)

	assert.True(t, isEqual)
}

func Test_StoragePlaceShouldBeNotEqualWhenIdIsNotSame(t *testing.T) {
	volume, _ := model.NewVolume(10)
	name := "BackPack"
	storage1, _ := NewStoragePlace(name, volume, nil)
	storage2, _ := NewStoragePlace(name, volume, nil)

	isEqual := storage1.Equal(storage2)

	assert.False(t, isEqual)
}

func Test_StoragePlaceCorrectlyDeterminesPossibilityPlacingOrder(t *testing.T) {
	testOrder := standardOrderWithVolume10()

	tests := map[string]struct {
		name           string
		volume         uint16
		initOrder      *order.Order
		placedOrder    *order.Order
		expectedResult bool
	}{
		"empty_storage_place": {
			name:           "BackPack",
			initOrder:      nil,
			placedOrder:    testOrder,
			volume:         20,
			expectedResult: true,
		},
		"used_storage_place": {
			name:           "BackPack",
			initOrder:      testOrder,
			placedOrder:    testOrder,
			volume:         20,
			expectedResult: false,
		},
		"small_storage_place": {
			name:           "BackPack",
			initOrder:      nil,
			placedOrder:    testOrder,
			volume:         5,
			expectedResult: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			volume, _ := model.NewVolume(test.volume)

			storagePlace, _ := NewStoragePlace(test.name, volume, test.initOrder)
			res := storagePlace.CanOrderBePlace(test.placedOrder)

			assert.Equal(t, test.expectedResult, res)
		})
	}
}

func Test_StoragePlaceCanBeUseWhenEmptyAndVolumeSufficient(t *testing.T) {
	volume, _ := model.NewVolume(10)
	name := "BackPack"
	storage, _ := NewStoragePlace(name, volume, nil)
	testOrder := standardOrderWithVolume10()

	err := storage.PlaceOrder(testOrder)

	assert.NoError(t, err)
}

func Test_StoragePlaceReturnErrorWhenParamsAreNotCorrectOnPlaceOrder(t *testing.T) {
	testOrder := standardOrderWithVolume10()

	tests := map[string]struct {
		name        string
		volume      uint16
		initOrder   *order.Order
		placedOrder *order.Order
		expectedErr error
	}{
		"used_storage_place": {
			name:        "BackPack",
			initOrder:   testOrder,
			placedOrder: testOrder,
			volume:      20,
			expectedErr: ErrStorageIsAlreadyTaken,
		},
		"small_storage_place": {
			name:        "BackPack",
			initOrder:   nil,
			placedOrder: testOrder,
			volume:      5,
			expectedErr: ErrOrderVolumeIsTooLarge,
		},
		"order not specified": {
			name:        "BackPack",
			initOrder:   nil,
			placedOrder: nil,
			volume:      5,
			expectedErr: errs.NewValueIsRequired("order"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			volume, _ := model.NewVolume(test.volume)

			storagePlace, _ := NewStoragePlace(test.name, volume, test.initOrder)
			err := storagePlace.PlaceOrder(test.placedOrder)

			if err.Error() != test.expectedErr.Error() {
				t.Errorf("expected %v, got %v", test.expectedErr, err)
			}
		})
	}
}

func Test_StoragePlaceShouldReturnOrderWhenNotEmpty(t *testing.T) {
	volume, _ := model.NewVolume(10)
	storage, _ := NewStoragePlace("BackPack", volume, nil)
	testOrder := standardOrderWithVolume10()

	_ = storage.PlaceOrder(testOrder)

	result, _ := storage.RetrieveOrder()
	assert.Equal(t, testOrder.Id(), *result)
	assert.Equal(t, true, storage.IsEmpty())
}

func Test_StoragePlaceReturnErrorWhenItEmptyOnRetrieveOrder(t *testing.T) {
	volume, _ := model.NewVolume(10)
	storage, _ := NewStoragePlace("BackPack", volume, nil)

	_, err := storage.RetrieveOrder()
	assert.ErrorIs(t, ErrStorageIsEmpty, err)
}
