package main

import (
	"log"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"

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
	api := operations.NewGreetingServerAPI(swaggerSpec)

	// simple handler based on code-generated types
	api.GetGreetingHandler = operations.GetGreetingHandlerFunc(
		func(param operations.GetGreetingParams) middleware.Responder {
			return operations.NewGetGreetingOK().
				WithPayload("Hi from simple example")
		})

	// create go-swagger server
	// no need to set api to server yet
	// because the Plugger will set API to the server
	// just create it to initialize properly
	srv := restapi.NewServer(nil)

	// wrap it to make plugable
	//
	// provide some server parameters as plugger options
	plug := plugger.NewPlug(srv, api,
		plugger.WithHost("localhost"),
		plugger.WithPort(8000))
	defer plug.Shutdown()

	// run server
	err = plug.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
