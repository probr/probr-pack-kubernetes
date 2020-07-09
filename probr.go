package main

import (
	"fmt"	 
	"flag"
	"os"	

	server "citihub.com/probr/api"
	api "citihub.com/probr/api/probrapi"
	"github.com/labstack/echo/v4"
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	var port = flag.Int("port", 1234, "Port for test HTTP server")
	flag.Parse()

	swagger, err := api.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}

	probrAPI := server.NewProbrAPI()

	e := echo.New()
	e.Use(echomiddleware.Logger())

	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	e.Use(middleware.OapiRequestValidator(swagger))

	// Register ProbrAPI as the handler for the interface
	api.RegisterHandlers(e, probrAPI)

	// And we serve HTTP until the world ends.
	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%d", *port)))

	

}
