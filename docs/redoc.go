package docs

import "github.com/mvrilo/go-redoc"

func Initialize() redoc.Redoc {
	return redoc.Redoc{
		Title:       "Documentation of ShopSystem",
		Description: "Documentation describes working procedures of ShopSystem like structs, handlers, etc.",
		SpecFile:    "./docs/api_shop_v1.swagger.json",
		SpecPath:    "/docs/api_shop_v1.swagger.json",
		DocsPath:    "/docs",
	}

}
