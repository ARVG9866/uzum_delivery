package api

import (
	"context"

	"github.com/Shemistan/uzum_delivery/internal/convert"
	delivery_v1 "github.com/Shemistan/uzum_delivery/internal/service/delivery_v1"
	pb "github.com/Shemistan/uzum_delivery/pkg/delivery_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Delivery struct {
	pb.UnimplementedDeliveryV1Server

	DeliveryService delivery_v1.IService
}

func (d *Delivery) AddOrderForDelivery(ctx context.Context, req *pb.AddOrderForDelivery_Request) (*emptypb.Empty, error) {
	err := req.Validate()
	if err != nil {
		return &emptypb.Empty{}, err
	}

	err = d.DeliveryService.CreateDeliveryOrder(ctx, convert.PbToModelOrder(req.Order))

	return &emptypb.Empty{}, err
}

func (d *Delivery) GetAllOrdersForDelivery(ctx context.Context, req *pb.GetAllOrdersForDelivery_Request) (*pb.GetAllOrdersForDelivery_Response, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	res, err := d.DeliveryService.GetOrderListForDelivery(ctx, convert.PbToModelCoordinate(req.CourierCoordinate))
	if err != nil {
		return nil, err
	}

	orders := make([]*pb.GetOrder, 0, len(res))

	for _, v := range res {
		orders = append(orders, convert.ModelToPbGetOrder(v))
	}

	rtn := &pb.GetAllOrdersForDelivery_Response{
		Orders: orders,
	}

	return rtn, nil
}
func (d *Delivery) GiveOrderForDelivery(ctx context.Context, req *pb.GiveOrderForDelivery_Request) (*pb.GiveOrderForDelivery_Response, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	res, err := d.DeliveryService.GiveOrderForDelivery(ctx, req.OrderId)
	if err != nil {
		return nil, err
	}

	rtn := &pb.GiveOrderForDelivery_Response{
		Order: convert.ModelToPbOrder(res),
	}

	return rtn, nil

}
func (d *Delivery) CloseOrderForDelivery(ctx context.Context, req *pb.CloseOrderForDelivery_Request) (*emptypb.Empty, error) {
	err := req.Validate()
	if err != nil {
		return &emptypb.Empty{}, err
	}

	err = d.DeliveryService.CloseOrder(ctx, req.OrderId)

	return &emptypb.Empty{}, err
}
