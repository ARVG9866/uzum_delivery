package service

import (
	"context"

	"github.com/Shemistan/uzum_delivery/internal/models"
	"github.com/Shemistan/uzum_delivery/internal/storage"
)

const (
	status1 = "Open"
	status2 = "InProgress"
	status3 = "Close"
)

type IService interface {
	CreateDeliveryOrder(context.Context, *models.Order) error
	GetOrderListForDelivery(context.Context, *models.Coordinate) ([]*models.OrderList, error)
	GiveOrderForDelivery(context.Context, int64) (*models.Order, error)
	CloseOrder(context.Context, int64) error
}

func NewService(storage storage.IStorage) IService {
	return &service{
		storage: storage,
	}
}

type service struct {
	storage storage.IStorage
}

func (s *service) CreateDeliveryOrder(ctx context.Context, order *models.Order) error {
	order.Status = status1

	err := s.storage.CreateOrder(ctx, order)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetOrderListForDelivery(ctx context.Context, coordinate *models.Coordinate) ([]*models.OrderList, error) {
	res, err := s.storage.GetOrdersForDelivery(ctx, coordinate)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *service) GiveOrderForDelivery(ctx context.Context, order_id int64) (*models.Order, error) {
	res, err := s.storage.GiveOrderToCourier(ctx, order_id, getCourier_id())
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *service) CloseOrder(ctx context.Context, order_id int64) error {
	err := s.storage.UpdateOrderStatus(ctx, order_id)

	return err
}

func getCourier_id() int64 {
	return 1
}
