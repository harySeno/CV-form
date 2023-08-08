package main

import (
	"CV-form/pkg/restapi"
	"github.com/emicklei/go-restful/v3"
	"log"
	"net/http"
)

func main() {
	// create new web service
	webService := new(restful.WebService)
	webService.Path("/api").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	// initiate mock data
	restapi.InitializeMockData()

	// create a new route
	webService.
		Route(webService.GET("/profile/{code}").
			To(restapi.GetByCode))

	webService.
		Route(webService.POST("/profile").
			To(restapi.AddProfile))

	// create a new container
	container := restful.NewContainer()
	container.Add(webService)

	// initiate server
	log.Println(" Initiating server on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", container))
}