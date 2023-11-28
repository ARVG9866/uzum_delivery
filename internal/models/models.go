package models

import "time"

type Order struct {
	ID                 int64       `json:"id"`
	OrderName          string      `json:"order_name"`
	UserName           string      `json:"user_name"`
	Phone              string      `json:"phone"`
	Address            string      `json:"address"`
	Coordinate_address *Coordinate `json:"coordinate_address"`
	Coordinate_opp     *Coordinate `json:"coordinate_opp"`
	Meta               string      `json:"meta"`
	Status             string      `json:"status"`
	DeliveryAt         time.Time   `json:"delivery_at"`
	Courier_id         int64       `json:"courier_id"`
}

type Coordinate struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type Product struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	Count       int64  `json:"count"`
}

type OrderList struct {
	ID       int64   `json:"id"`
	Distance float32 `json:"distance"`
}

type Token struct {
	Access  string
	Refresh string
}
