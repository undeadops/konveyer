package main

import (
	"fmt"
	"log"

	"github.com/undeadops/konveyer/pkg"
	"github.com/undeadops/konveyer/pkg/api"
	"github.com/undeadops/konveyer/pkg/config"
	"github.com/undeadops/konveyer/pkg/mongo"
)

// App - Base struct for managing Application
type App struct {
	server  *api.API
	session *mongo.Session
	config  *root.Config
}

// Initialize - Setup App for use
func (a *App) Initialize() {
	a.config = config.LoadConfig()
	var err error
	a.session, err = mongo.NewSession(a.config)
	if err != nil {
		log.Fatalln("unable to connect to mongodb")
	}

	u := mongo.NewDeploymentService(a.session.Copy(), a.config)
	a.server = api.NewAPI(u, a.config)
}

// Run - Run Webserver
func (a *App) Run() {
	fmt.Println("Run")
	defer a.session.Close()
	a.server.Start()
}
