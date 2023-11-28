package app

import (
	"context"
	"fmt"
	"log"

	gateway_runtime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	"github.com/mvrilo/go-redoc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Shemistan/uzum_delivery/dev"
	"github.com/Shemistan/uzum_delivery/internal/models"
	delivery_v1 "github.com/Shemistan/uzum_delivery/internal/service/delivery_v1"
	rep "github.com/Shemistan/uzum_delivery/internal/storage"
	pb_login "github.com/Shemistan/uzum_delivery/pkg/login_v1"
)

type App struct {
	appConfig   *models.Config
	muxDelivery *gateway_runtime.ServeMux

	grpcDeliveryServer *grpc.Server
	deliveryService    delivery_v1.IService
	db                 *sqlx.DB
	reDoc              redoc.Redoc
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}
	a.setConfig()
	a.initDB()
	a.initReDoc()
	a.initGRPCServer()

	if err := a.initHTTPServer(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) setConfig() {
	if dev.DEBUG {
		err := dev.SetConfig()
		if err != nil {
			log.Fatal("failed to get config", err.Error())
		}
	}
	conf := models.Config{}

	envconfig.MustProcess("", &conf)

	a.appConfig = &conf
}

func (a *App) getService() delivery_v1.IService {
	storage := rep.NewStorage(a.db)

	conn, err := grpc.Dial(a.appConfig.App.AuthClient, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}

	loginClient := pb_login.NewLoginV1Client(conn)

	if a.deliveryService == nil {
		a.deliveryService = delivery_v1.NewService(storage, loginClient)

	}
	return a.deliveryService
}

func (a *App) getSqlConnectionString() string {
	sqlConnectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%v",
		a.appConfig.DB.User,
		a.appConfig.DB.Password,
		a.appConfig.DB.Host,
		a.appConfig.DB.Port,
		a.appConfig.DB.Database,
		a.appConfig.DB.SSLMode,
	)

	return sqlConnectionString
}
