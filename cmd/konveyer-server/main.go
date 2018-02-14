package main

import (
	"github.com/undeadops/konveyer"

	"github.com/go-openapi/loads"

	"github.com/undeadops/konveyer/gen/restapi"
	"github.com/undeadops/konveyer/gen/restapi/operations"
	"github.com/undeadops/konveyer/gen/restapi/operations/projects"
)

func main() {

	app := konveyer.NewApp()
	defer app.Close()

	// load embedded swagger file
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		app.Log.Panic(err)
	}

	// create new service API
	api := operations.NewKonveyerAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()

	server.Port = app.Port

	api.ProjectsFindProjectsHandler = projects.FindProjectsHandlerFunc(app.handlers.ListProjects)
	api.ProjectsDescribeProjectsHandler = projects.DescribeProjectsHandlerFunc(app.handlers.DescribeProject)

	// serve API
	if err := server.Serve(); err != nil {
		app.Log.Fatalln(err)
	}
}
