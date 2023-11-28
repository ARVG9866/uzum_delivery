package convert

import (
	"github.com/Shemistan/uzum_delivery/internal/models"
	pb "github.com/Shemistan/uzum_delivery/pkg/delivery_v1"
	pb_login "github.com/Shemistan/uzum_delivery/pkg/login_v1"
)

func PbToModelOrder(order *pb.Order) *models.Order {
	return &models.Order{
		ID:        order.Id,
		OrderName: order.OrderName,
		UserName:  order.UserName,
		Phone:     order.UserPhone,
		Address:   order.Address,
		Coordinate_address: &models.Coordinate{
			X: order.CoordinateAddress.X,
			Y: order.CoordinateAddress.Y,
		},
		Coordinate_opp: &models.Coordinate{
			X: order.CoordinateOpp.X,
			Y: order.CoordinateOpp.Y,
		},
		Meta: order.Meta,
	}
}

func PbToModelCoordinate(coordinate *pb.Coordinate) *models.Coordinate {
	return &models.Coordinate{
		X: coordinate.X,
		Y: coordinate.Y,
	}
}

func ModelToPbGetOrder(order *models.OrderList) *pb.GetOrder {
	return &pb.GetOrder{
		OrderId:  order.ID,
		Distance: order.Distance,
	}
}

func ModelToPbOrder(order *models.Order) *pb.Order {
	return &pb.Order{
		Id:        order.ID,
		OrderName: order.OrderName,
		UserName:  order.UserName,
		UserPhone: order.Phone,
		Address:   order.Address,
		CoordinateAddress: &pb.Coordinate{
			X: order.Coordinate_address.X,
			Y: order.Coordinate_address.Y,
		},
		CoordinateOpp: &pb.Coordinate{
			X: order.Coordinate_opp.X,
			Y: order.Coordinate_opp.Y,
		},
		Meta:   order.Meta,
		Status: order.Status,
	}
}

func GetToken(auth *pb_login.Login_Response) *models.Token {
	return &models.Token{
		Refresh: auth.RefreshToken,
		Access:  auth.AccessToken,
	}
}
