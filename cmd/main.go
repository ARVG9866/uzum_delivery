package main

import (
	"context"
	"log"

	"github.com/Shemistan/uzum_delivery/internal/app"
)

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatal("failed to create app")
	}

	err = a.Run()
	if err != nil {
		log.Fatal("failed to run app")
	}
}
