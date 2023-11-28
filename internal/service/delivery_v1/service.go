package service

import (
	"context"
	"errors"

	"github.com/Shemistan/uzum_delivery/internal/convert"
	"github.com/Shemistan/uzum_delivery/internal/models"
	"github.com/Shemistan/uzum_delivery/internal/storage"
	pb_login "github.com/Shemistan/uzum_delivery/pkg/login_v1"
	"google.golang.org/grpc/metadata"
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
	Login(ctx context.Context, login string, password string) (*models.Token, error)
}

func NewService(storage storage.IStorage, loginClient pb_login.LoginV1Client) IService {
	return &service{
		storage:     storage,
		loginClient: loginClient,
	}
}

type service struct {
	storage     storage.IStorage
	loginClient pb_login.LoginV1Client
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
	user_id, err := s.GetUserFromToken(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.storage.GiveOrderToCourier(ctx, order_id, user_id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *service) CloseOrder(ctx context.Context, order_id int64) error {
	err := s.storage.UpdateOrderStatus(ctx, order_id)

	return err
}

func (s *service) Login(ctx context.Context, login string, password string) (*models.Token, error) {
	req := &pb_login.Login_Request{Login: login, Password: password}
	auth, err := s.loginClient.Login(ctx, req)
	if err != nil {
		return nil, err
	}

	return convert.GetToken(auth), nil
}

func (s *service) GetUserFromToken(ctx context.Context) (int64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, errors.New("Can't get context")
	}
	ctx = metadata.NewOutgoingContext(ctx, md)
	check, err := s.loginClient.Check(ctx, &pb_login.Check_Request{EndpointAddress: ""})
	if err != nil {
		return 0, err
	}

	return check.UserId, nil
}
