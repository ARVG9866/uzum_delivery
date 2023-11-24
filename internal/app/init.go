package app

import (
	"context"
	"log"

	"github.com/Shemistan/uzum_delivery/docs"
	"github.com/Shemistan/uzum_delivery/internal/api"
	_ "github.com/lib/pq"

	gateway_runtime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/Shemistan/uzum_delivery/pkg/delivery_v1"
)

func (a *App) initDB() {
	sqlConnectionString := a.getSqlConnectionString()

	var err error
	a.db, err = sqlx.Open("postgres", sqlConnectionString)
	if err != nil {
		log.Fatal("failed to opening connection to db: ", err.Error())
	}

	if err = a.db.Ping(); err != nil {
		log.Fatal("failed to connect to the database: ", err.Error())
	}
}

func (a *App) initReDoc() {
	a.reDoc = docs.Initialize()
}

func (a *App) initHTTPServer(ctx context.Context) error {
	a.muxDelivery = gateway_runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pb.RegisterDeliveryV1HandlerFromEndpoint(ctx, a.muxDelivery, a.appConfig.App.PortGRPC, opts)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initGRPCServer() {
	a.grpcDeliveryServer = grpc.NewServer()
	pb.RegisterDeliveryV1Server(
		a.grpcDeliveryServer,
		&api.Delivery{
			DeliveryService: a.getService(),
		},
	)
}
