package courier

import (
	"delivery/internal/core/domain/model"
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
	orderId     uuid.UUID
}

func NewStoragePlace(name string, totalVolume model.Volume) (*StoragePlace, error) {
	if name == "" {
		return nil, errs.NewValueIsRequired("name")
	}

	return &StoragePlace{
		baseEntity:  ddd.NewBaseEntity(uuid.New()),
		name:        name,
		totalVolume: totalVolume,
		orderId:     uuid.Nil,
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

func (storage *StoragePlace) OrderId() uuid.UUID {
	return storage.orderId
}

func (storage *StoragePlace) Equal(other *StoragePlace) bool {
	if other == nil {
		return false
	}

	return storage.baseEntity.Equal(other.baseEntity)
}

func (storage *StoragePlace) CanOrderBePlace(orderVolume model.Volume) bool {
	err := storage.checkCanBePlace(orderVolume)
	return err == nil
}

func (storage *StoragePlace) placeOrder(orderId uuid.UUID, orderVolume model.Volume) error {
	if orderId == uuid.Nil {
		return errs.NewValueIsRequired("orderId")
	}

	err := storage.checkCanBePlace(orderVolume)
	if err != nil {
		return err
	}

	storage.orderId = orderId

	return nil
}

func (storage *StoragePlace) retrieveOrder() (uuid.UUID, error) {
	if storage.IsEmpty() {
		return uuid.Nil, ErrStorageIsEmpty
	}

	orderId := storage.orderId
	storage.orderId = uuid.Nil

	return orderId, nil
}

func (storage *StoragePlace) IsEmpty() bool {
	return storage.orderId == uuid.Nil
}

func (storage *StoragePlace) checkCanBePlace(orderVolume model.Volume) error {
	if storage.orderId != uuid.Nil {
		return ErrStorageIsAlreadyTaken
	}

	if storage.totalVolume.LessThen(orderVolume) {
		return ErrOrderVolumeIsTooLarge
	}

	return nil
}
