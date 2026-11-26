package courier

import (
	"delivery/internal/core/domain/model"
	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/errs"
	"errors"
	"github.com/google/uuid"
	"math"
)

var (
	ErrOrderCannotBeTaken       = errors.New("order cannot be taken")
	ErrOrderNotFoundWithCourier = errors.New("order not found with the courier")
)

type Courier struct {
	baseAggregate *ddd.BaseAggregate[uuid.UUID]
	name          string
	speed         model.Speed
	location      model.Location
	storagePlaces []*StoragePlace
}

func newDefaultBag() *StoragePlace {
	volume, _ := model.NewVolume(10)
	defaultBag, _ := NewStoragePlace("Bag", volume)
	return defaultBag
}

func NewCourier(name string, speed int, location model.Location) (*Courier, error) {
	if name == "" {
		return nil, errs.NewValueIsRequired("name")
	}

	s, err := model.NewSpeed(speed)
	if err != nil {
		return nil, err
	}

	return &Courier{
		baseAggregate: ddd.NewBaseAggregate(uuid.New()),
		name:          name,
		speed:         s,
		location:      location,
		storagePlaces: []*StoragePlace{newDefaultBag()},
	}, nil
}

func (c *Courier) ID() uuid.UUID { return c.baseAggregate.ID() }

func (c *Courier) Name() string { return c.name }

func (c *Courier) Speed() model.Speed { return c.speed }

func (c *Courier) Location() model.Location { return c.location }

func (c *Courier) StoragePlaces() []StoragePlace {
	res := make([]StoragePlace, len(c.storagePlaces))
	for i, storagePlace := range c.storagePlaces {
		res[i] = *storagePlace
	}
	return res
}

func (c *Courier) Equal(other *Courier) bool {
	if other == nil {
		return false
	}

	return c.baseAggregate.Equal(other.baseAggregate)
}

func (c *Courier) AddStoragePlace(name string, value model.Volume) error {

	storagePlace, err := NewStoragePlace(name, value)
	if err != nil {
		return err
	}

	c.storagePlaces = append(c.storagePlaces, storagePlace)
	return nil
}

func (c *Courier) CanTakeOrderWithVolume(volume model.Volume) bool {
	for _, storagePlace := range c.storagePlaces {
		if storagePlace.CanOrderBePlace(volume) {
			return true
		}
	}

	return false
}

func (c *Courier) TakeOrder(orderId uuid.UUID, orderVolume model.Volume) error {

	for _, storagePlace := range c.storagePlaces {
		if storagePlace.CanOrderBePlace(orderVolume) {
			err := storagePlace.placeOrder(orderId, orderVolume)
			if err != nil {
				return err
			}

			return nil
		}
	}

	return ErrOrderCannotBeTaken
}

func (c *Courier) CompleteOrder(orderId uuid.UUID) error {
	if orderId == uuid.Nil {
		return errs.NewValueIsRequired("orderId")
	}

	for _, storagePlace := range c.storagePlaces {
		if storagePlace.OrderId() == orderId {
			_, err := storagePlace.retrieveOrder()
			if err != nil {
				return err
			}

			return nil
		}
	}

	return ErrOrderNotFoundWithCourier
}

func (c *Courier) GetStepCount(target model.Location) int {
	speed := c.speed.Speed()
	distance := c.location.Distance(target)
	steps := distance / speed
	if distance%speed != 0 {
		steps++
	}

	return steps
}

func (c *Courier) Move(target model.Location) error {
	dx := float64(target.X()) - float64(c.location.X())
	dy := float64(target.Y()) - float64(c.location.Y())

	remainingRange := float64(c.speed.Speed())

	if math.Abs(dx) > remainingRange {
		dx = math.Copysign(remainingRange, dx)
	}
	remainingRange -= math.Abs(dx)

	if math.Abs(dy) > remainingRange {
		dy = math.Copysign(remainingRange, dy)
	}

	newX := c.location.X() + int(dx)
	newY := c.location.Y() + int(dy)

	newLocation, err := model.NewLocation(newX, newY)
	if err != nil {
		return err
	}
	c.location = newLocation

	return nil
}
