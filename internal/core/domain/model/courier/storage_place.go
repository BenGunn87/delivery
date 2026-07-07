package courier

import (
	"delivery/internal/core/domain/model"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/errs"
	"errors"
	"github.com/google/uuid"
)

var (
	ErrOrderVolumeIsTooLarge = errors.New("order volume is too large")
	ErrStorageIsAlreadyTaken = errors.New("storage place is already taken")
	ErrStorageIsEmpty        = errors.New("storage place is empty")
)

type StoragePlace struct {
	baseEntity  *ddd.BaseEntity[uuid.UUID]
	name        string
	totalVolume model.Volume
	orderId     *uuid.UUID
}

func NewStoragePlace(name string, totalVolume model.Volume, order *order.Order) (*StoragePlace, error) {
	if name == "" {
		return nil, errs.NewValueIsRequired("name")
	}
	var orderId *uuid.UUID

	if order != nil {
		if totalVolume.LessThen(order.Volume()) {
			return nil, ErrOrderVolumeIsTooLarge
		}
		id := order.Id()
		orderId = &(id)
	} else {
		orderId = nil
	}

	return &StoragePlace{
		baseEntity:  ddd.NewBaseEntity(uuid.New()),
		name:        name,
		totalVolume: totalVolume,
		orderId:     orderId,
	}, nil
}

func (storage *StoragePlace) ID() uuid.UUID {
	return storage.baseEntity.ID()
}

func (storage *StoragePlace) Name() string {
	return storage.name
}

func (storage *StoragePlace) TotalVolume() model.Volume {
	return storage.totalVolume
}

func (storage *StoragePlace) OrderId() *uuid.UUID {
	if storage.orderId == nil {
		return nil
	}

	orderId := *storage.orderId
	return &orderId

}

func (storage *StoragePlace) Equal(other *StoragePlace) bool {
	if other == nil {
		return false
	}

	return storage.baseEntity.Equal(other.baseEntity)
}

func (storage *StoragePlace) CanOrderBePlace(order *order.Order) bool {
	err := storage.checkCanBePlace(order)
	if err != nil {
		return false
	}

	return true
}

func (storage *StoragePlace) PlaceOrder(order *order.Order) error {
	err := storage.checkCanBePlace(order)
	if err != nil {
		return err
	}

	orderId := order.Id()
	storage.orderId = &orderId

	return nil
}

func (storage *StoragePlace) RetrieveOrder() (*uuid.UUID, error) {
	if storage.IsEmpty() {
		return nil, ErrStorageIsEmpty
	}

	orderId := storage.orderId
	storage.orderId = nil

	return orderId, nil
}

func (storage *StoragePlace) IsEmpty() bool {
	return storage.orderId == nil
}

func (storage *StoragePlace) checkCanBePlace(order *order.Order) error {
	if order == nil {
		return errs.NewValueIsRequired("order")
	}
	if storage.orderId != nil {
		return ErrStorageIsAlreadyTaken
	}

	if storage.totalVolume.LessThen(order.Volume()) {
		return ErrOrderVolumeIsTooLarge
	}

	return nil
}
