package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"

	"github.com/bazani/go-horse/go-rest-api/pkg/swagger/server/restapi"
	"github.com/bazani/go-horse/go-rest-api/pkg/swagger/server/restapi/operations"
)

func main() {

	// Initialize Swagger
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewHelloAPIAPI(swaggerSpec)
	server := restapi.NewServer(api)

	defer func() {
		if err := server.Shutdown(); err != nil {
			// error handle ???
			log.Fatalln(err)
		}
	}()

	server.Port = 8080

	api.CheckHealthHandler = operations.CheckHealthHandlerFunc(Health)

	api.GetHelloUserHandler = operations.GetHelloUserHandlerFunc(GetHelloUser)

	api.GetGopherNameHandler = operations.GetGopherNameHandlerFunc(GetGopherByName)

	// Start listening
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}

// Health route returns OK
func Health(operations.CheckHealthParams) middleware.Responder {
	return operations.NewCheckHealthOK().WithPayload("OK")
}

// GetHelloUser returns Hello + the parameter
func GetHelloUser(user operations.GetHelloUserParams) middleware.Responder {
	return operations.NewGetHelloUserOK().WithPayload("Hello there " + user.User + " =)")
}

// GetGopherByName returns a gopher image in png
func GetGopherByName(gopher operations.GetGopherNameParams) middleware.Responder {
	var URL string

	if gopher.Name != "" {
		URL = "https://raw.githubusercontent.com/scraly/gophers/main/" + gopher.Name + ".png"
	} else {
		// returns dr who gopher by default
		URL = "https://raw.githubusercontent.com/scraly/gophers/main/dr-who.png"
	}

	response, err := http.Get(URL)
	if err != nil {
		fmt.Println("error trying to get the gopher image")
	}

	return operations.NewGetGopherNameOK().WithPayload(response.Body)
}
