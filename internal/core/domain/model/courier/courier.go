package courier

import (
	"delivery/internal/core/domain/model"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/errs"
	"errors"
	"github.com/google/uuid"
	"math"
)

var (
	ErrOrderCannotBeTaken      = errors.New("order cannot be taken")
	ErrOrderNoFoundWithCourier = errors.New("order not fount with the courier")
)

type Courier struct {
	baseAggregate *ddd.BaseAggregate[uuid.UUID]
	name          string
	speed         uint8
	location      model.Location
	storagePlaces []*StoragePlace
}

func NewDefaultBag() *StoragePlace {
	volume, _ := model.NewVolume(10)
	defaultBag, _ := NewStoragePlace("Bag", volume, nil)
	return defaultBag
}

func NewCourier(name string, speed uint8, location model.Location) (*Courier, error) {
	if name == "" {
		return nil, errs.NewValueIsRequired("name")
	}
	if speed == 0 {
		return nil, errs.NewValueMustBeGreaterThan("speed", speed, 0)
	}

	return &Courier{
		baseAggregate: ddd.NewBaseAggregate(uuid.New()),
		name:          name,
		speed:         speed,
		location:      location,
		storagePlaces: []*StoragePlace{NewDefaultBag()},
	}, nil
}

func (c *Courier) Id() uuid.UUID { return c.baseAggregate.ID() }

func (c *Courier) Name() string { return c.name }

func (c *Courier) Speed() uint8 { return c.speed }

func (c *Courier) Location() model.Location { return c.location }

func (c *Courier) StoragePlaces() []*StoragePlace { return c.storagePlaces }

func (c *Courier) Equal(other *Courier) bool {
	if other == nil {
		return false
	}

	return c.baseAggregate.Equal(other.baseAggregate)
}

func (c *Courier) AddStoragePlace(name string, value uint16) error {
	val, err := model.NewVolume(value)
	if err != nil {
		return err
	}
	storagePlace, err := NewStoragePlace(name, val, nil)
	if err != nil {
		return err
	}

	c.storagePlaces = append(c.storagePlaces, storagePlace)
	return nil
}

func (c *Courier) CanTakeOrder(o order.Order) bool {
	orderVolume := o.Volume()
	for _, storagePlace := range c.storagePlaces {
		if storagePlace.CanOrderBePlace(orderVolume) {
			return true
		}
	}

	return false
}

func (c *Courier) TakeOrder(o *order.Order) error {
	orderVolume := o.Volume()

	for _, storagePlace := range c.storagePlaces {
		if storagePlace.CanOrderBePlace(orderVolume) {
			err := o.Assigned(c.Id())
			if err != nil {
				return err
			}

			err = storagePlace.PlaceOrder(o.Id(), orderVolume)
			if err != nil {
				return err
			}

			return nil
		}
	}

	return ErrOrderCannotBeTaken
}

func (c *Courier) CompleteOrder(o *order.Order) error {
	for _, storagePlace := range c.storagePlaces {
		if storagePlace.OrderId() != nil && *storagePlace.OrderId() == o.Id() {
			err := o.Completed()
			if err != nil {
				return err
			}

			_, err = storagePlace.RetrieveOrder()
			if err != nil {
				return err
			}

			return nil
		}
	}

	return ErrOrderNoFoundWithCourier
}

func (c *Courier) GetStepCount(target model.Location) uint8 {
	steps := c.location.Distance(target) / c.Speed()
	if c.location.Distance(target)%c.Speed() != 0 {
		steps++
	}

	return steps
}

func (c *Courier) Move(target model.Location) error {
	dx := float64(target.X()) - float64(c.location.X())
	dy := float64(target.Y()) - float64(c.location.Y())
	remainingRange := float64(c.speed)

	if math.Abs(dx) > remainingRange {
		dx = math.Copysign(remainingRange, dx)
	}
	remainingRange -= math.Abs(dx)

	if math.Abs(dy) > remainingRange {
		dy = math.Copysign(remainingRange, dy)
	}

	newX := c.location.X() + uint8(dx)
	newY := c.location.Y() + uint8(dy)

	newLocation, err := model.NewLocation(newX, newY)
	if err != nil {
		return err
	}
	c.location = newLocation

	return nil
}
