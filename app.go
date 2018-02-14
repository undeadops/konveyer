package konveyer

import (
	"github.com/globalsign/mgo"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

// Application - context for running app
type Application struct {
	Db    *mgo.Session
	Log   *logrus.Logger
	Debug bool   `default:"false"`
	Port  int    `default:"3000"`
	URI   string `envconfig:"MONGO_URI"`
}

// NewApp - Initialize context of app parameters
func NewApp() *Application {
	var a Application
	log := *logrus.New()

	err := envconfig.Process("", &a)
	if err != nil {
		log.Fatalf("Error Processing Environment Config: %s\n", err)
	}
	a.Log = &log

	m, err := mgo.ParseURL(a.URI)
	if err != nil {
		log.Fatalf("Error Connecting to Mongo URI: %s\n", err)
	}
	muri, err := mgo.DialWithInfo(m)
	if err != nil {
		log.Fatalf("Error Connecting to Mongo URI: %s\n", err)
	}

	a.Db = muri
	return &a
}

// Close - Clean up DB Connection
func (a *Application) Close() {
	// Close Database Connection
	a.Db.Close()
}
