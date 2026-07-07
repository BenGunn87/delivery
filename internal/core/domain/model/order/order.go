package order

import (
	"delivery/internal/core/domain/model"
	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/errs"
	"github.com/google/uuid"
)

type Order struct {
	baseAggregate *ddd.BaseAggregate[uuid.UUID]
	volume        model.Volume
}

func NewOrder(id uuid.UUID, volume model.Volume) (*Order, error) {
	if id == uuid.Nil {
		return nil, errs.NewValueIsRequired("id")
	}

	return &Order{
		baseAggregate: ddd.NewBaseAggregate(id),
		volume:        volume,
	}, nil
}

func (order *Order) Id() uuid.UUID {
	return order.baseAggregate.ID()
}

func (order *Order) Volume() model.Volume {
	return order.volume
}

func (order *Order) Equal(other *Order) bool {
	if other == nil {
		return false
	}

	return order.baseAggregate.Equal(other.baseAggregate)
}
