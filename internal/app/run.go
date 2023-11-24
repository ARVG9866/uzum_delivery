package app

import (
	"log"
	"net"
	"net/http"
	"sync"
)

func (a *App) Run() error {
	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()

		log.Fatal(a.runGRPC())
	}()

	go func() {
		defer wg.Done()

		log.Fatal(a.RunHTTP())
	}()

	go func() {
		defer wg.Done()

		log.Fatal(a.RunDocumentation())
	}()

	wg.Wait()
	return nil
}
func (a *App) runGRPC() error {
	listener, err := net.Listen("tcp", a.appConfig.App.PortGRPC)
	if err != nil {
		return err
	}
	log.Println("Delivery GRPC server running on port:", a.appConfig.App.PortGRPC)

	return a.grpcDeliveryServer.Serve(listener)
}

func (a *App) RunHTTP() error {
	log.Println("HTTP server is running on port:", a.appConfig.App.PortHTTP)

	return http.ListenAndServe(a.appConfig.App.PortHTTP, a.muxDelivery)

}

func (a *App) RunDocumentation() error {
	log.Println("Swagger documentation running on port:", a.appConfig.App.PortDocs)

	return http.ListenAndServe(a.appConfig.App.PortDocs, a.reDoc.Handler())

}
