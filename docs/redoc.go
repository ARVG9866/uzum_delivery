package docs

import "github.com/mvrilo/go-redoc"

func Initialize() redoc.Redoc {
	return redoc.Redoc{
		Title:       "Documentation of DeliverySystem",
		Description: "Documentation describes working procedures of DeliverySystem like structs, handlers, etc.",
		SpecFile:    "./docs/api_delivery_v1.swagger.json",
		SpecPath:    "/docs/api_delivery_v1.swagger.json",
		DocsPath:    "/docs",
	}

}
