package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Shemistan/uzum_delivery/cmd/api/handlers"
	service "github.com/Shemistan/uzum_delivery/internal/service/delivery_v1"
	"github.com/gorilla/mux"
)

const (
	httpPort = ":8080"
)

func main() {
	deliveryService := service.NewService()
	handler := handlers.NewHandler(deliveryService)

	router := mux.NewRouter()

	router.HandleFunc("/healthz", handler.Healthz).Methods(http.MethodGet)
	router.HandleFunc("/deliver/v1/give-delivery/{id:[0-9]+}", handler.GiveOrder).Methods(http.MethodGet)
	router.HandleFunc("/deliver/v1/add-order/", handler.AddOrder).Methods(http.MethodPost)

	srv := &http.Server{
		Addr:        httpPort,
		ReadTimeout: time.Second * 10,
		Handler:     router,
	}

	log.Println("delivery server is running at port:", httpPort)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
