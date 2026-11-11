package order

import (
	"delivery/internal/core/domain/model"
	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/errs"
	"errors"
	"github.com/google/uuid"
)

var (
	ErrStatusIsNotCreated  = errors.New("order status must be Created")
	ErrStatusIsNotAssigned = errors.New("order status must be Assigned")
)

type Order struct {
	baseAggregate *ddd.BaseAggregate[uuid.UUID]
	volume        model.Volume
	status        Status
	courierId     *uuid.UUID
}

func NewOrder(id uuid.UUID, volume model.Volume) (*Order, error) {
	if id == uuid.Nil {
		return nil, errs.NewValueIsRequired("id")
	}

	return &Order{
		baseAggregate: ddd.NewBaseAggregate(id),
		volume:        volume,
		status:        StatusCreated,
		courierId:     nil,
	}, nil
}

func (order *Order) ID() uuid.UUID {
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

func (order *Order) Status() Status {
	return order.status
}

func (order *Order) CourierId() *uuid.UUID {
	if order.courierId == nil {
		return nil
	}

	courierId := *order.courierId
	return &courierId
}

func (order *Order) MarkAsAssigned(courierId uuid.UUID) error {
	if courierId == uuid.Nil {
		return errs.NewValueIsRequired("courierId")
	}

	if order.status != StatusCreated {
		return ErrStatusIsNotCreated
	}
	order.status = StatusAssigned
	order.courierId = &courierId
	return nil
}

func (order *Order) MarkAsCompleted() error {
	if order.status != StatusAssigned {
		return ErrStatusIsNotAssigned
	}
	order.status = StatusCompleted
	order.courierId = nil
	return nil
}
