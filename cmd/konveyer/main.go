package main

// import (
// 	"log"
// 	"os"

// 	"github.com/sirupsen/logrus"

// 	"github.com/undeadops/konveyer/pkg/api"
// 	"github.com/undeadops/konveyer/pkg/config"
// 	"github.com/undeadops/konveyer/pkg/db"
// )

var Version = "0.0.1"

// func main() {
// 	logger := logrus.New()
// 	logger.Formatter = &logrus.TextFormatter{
// 		// disable, as we set our own
// 		FullTimestamp: true,
// 	}

// 	logger.Info("Starting Up....")

// 	config, err := config.LoadConfig()
// 	if err != nil {
// 		log.Fatal("Failed to load config: " + err.Error())
// 	}

// 	// logger, err := conf.ConfigureLogging(&config.LogConfig)
// 	// if err != nil {
// 	// 	log.Fatal("Failed to configure logging: " + err.Error())
// 	// }

// 	logger.Infof("Connecting to DB")
// 	db, err := db.Connect(config.DBURI)
// 	if err != nil {
// 		logger.Fatal("Failed to connect to db: " + err.Error())
// 	}

// 	logger.Infof("Starting API on port %s", config.Port)
// 	a := api.NewAPI(config, db, Version)
// 	err = a.Serve()
// 	if err != nil {
// 		logger.WithError(err).Error("Error while running API: %v", err)
// 		os.Exit(1)
// 	}
// 	logger.Info("API Shutdown")
// }

func main() {
	a := App{}
	a.Initialize()
	a.Run()
}
