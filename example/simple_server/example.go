package main

import (
	"log"

	"github.com/go-openapi/loads"

	"github.com/ilyakaznacheev/go-plugger"
	"github.com/ilyakaznacheev/go-plugger/example/simple_server/restapi"
	"github.com/ilyakaznacheev/go-plugger/example/simple_server/restapi/operations"
)

func main() {
	// generate spec for swagger API
	// this step is required by go-swagger
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatal(err)
	}
	// generate swagger-API handler
	swagger := operations.NewGreetingServerAPI(swaggerSpec)

	// generate go-swagger server
	srv := restapi.NewServer(swagger)

	// wrap it to make plugable
	//
	// provide some server parameters as plugger options
	plug := plugger.NewPlug(srv,
		plugger.WithConfiguredAPI(),
		plugger.WithHost("localhost"),
		plugger.WithPort(8888))
	defer plug.Shutdown()

	// run server
	err = plug.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
