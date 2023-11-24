package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Shemistan/uzum_delivery/internal/models"
	"github.com/jmoiron/sqlx"
)

const (
	status1 = "Open"
	status2 = "InProgress"
	status3 = "Close"
)

const orderTable = "order_for_delivery"

type IStorage interface {
	CreateOrder(context.Context, *models.Order) error
	GetOrdersForDelivery(context.Context, *models.Coordinate) ([]*models.OrderList, error)
	GiveOrderToCourier(context.Context, int64, int64) (*models.Order, error)
	UpdateOrderStatus(context.Context, int64) error
}

type storage struct {
	db *sqlx.DB
}

func NewStorage(db *sqlx.DB) IStorage {
	return &storage{db: db}
}

func (s *storage) CreateOrder(ctx context.Context, order *models.Order) error {
	str_query := fmt.Sprintf(`INSERT INTO %s (
		id, 
		order_name, 
		user_name, 
		phone, 
		address, 
		coordinate_address_x, 
		coordinate_address_y,
		coordinate_oop_x, 
		coordinate_oop_y, 
		meta, 
		status
	) 
	VALUES(%d, %s, %s, %s, %s, %f, %f, %f, %f, %s, %s)`,
		orderTable,
		order.ID,
		order.OrderName,
		order.UserName,
		order.Phone,
		order.Address,
		order.Coordinate_address.X,
		order.Coordinate_address.Y,
		order.Coordinate_opp.X,
		order.Coordinate_opp.Y,
		order.Meta,
		order.Status)

	_, err := s.db.ExecContext(ctx, str_query)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) GetOrdersForDelivery(ctx context.Context, courier_coordinate *models.Coordinate) ([]*models.OrderList, error) {
	// (x1, y1) - opp, (x2, y2) - courier,  (x3, y3) - client
	//SQRT((x2 - x1)^2 + (y2 - y1)^2) + SQRT((x3 - x2)^2 + (y3 - y2)^ 2)
	orders := make([]*models.OrderList, 0, 5)

	str_query := fmt.Sprintf(
		`SELECT
			order_id,  
			ROUND(SQRT(POWER(%f - coordinate_oop_x, 2) + POWER(%f - coordinate_oop_y, 2)) +
				SQRT(POWER(coordinate_address_x - %f, 2) + POWER(coordinate_address_y - %f, 2)), 2) AS distance
		FROM %s 
		WHERE status = %s 
		ORDER BY 1 
		LIMIT 5;`,
		courier_coordinate.X, courier_coordinate.Y, courier_coordinate.X, courier_coordinate.Y, orderTable, status1)

	rows, err := s.db.QueryContext(ctx, str_query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var row models.OrderList

		err = rows.Scan(&row.ID, &row.Distance)
		if err != nil {
			return nil, err
		}

		orders = append(orders, &row)
	}

	return orders, nil
}

func (s *storage) GiveOrderToCourier(ctx context.Context, order_id int64, courier_id int64) (*models.Order, error) {
	var order models.Order

	err := s.updateOrderStatus2(ctx, order_id, courier_id)
	if err != nil {
		return nil, err
	}

	str_query := fmt.Sprintf(`SELECT * FROM %s WHERE id=%d`, orderTable, order_id)

	err = s.db.QueryRowxContext(ctx, str_query).StructScan(&order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (s *storage) UpdateOrderStatus(ctx context.Context, order_id int64) error {
	today := time.Now()

	var rtn int64
	str_query := fmt.Sprintf(`SELECT id FROM %s WHERE id=%d AND status=%s`, orderTable, order_id, status2)
	err := s.db.QueryRowContext(ctx, str_query).Scan(&rtn)
	if err != nil {
		return err
	}
	if rtn == 0 {
		return errors.New("Order not found or status is not correct")
	}

	str_query = fmt.Sprintf(`UPDATE %s SET status=%s, delivery_at=%s WHERE id=%d`, orderTable, status3, today, order_id)
	_, err = s.db.ExecContext(ctx, str_query)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) updateOrderStatus2(ctx context.Context, order_id int64, courier_id int64) error {
	str_query := fmt.Sprintf("UPDATE %s SET status=%s, courier_id=%d WHERE id=%d", orderTable, status2, courier_id, order_id)

	_, err := s.db.ExecContext(ctx, str_query)
	if err != nil {
		return err
	}

	return nil
}
