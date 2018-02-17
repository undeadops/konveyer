package konveyer

import (
	"log"

	"github.com/kelseyhightower/envconfig"
	mgo "gopkg.in/mgo.v2"
)

const appName = "konveyer-server"

// Runtime - Application Context instead of globals
type Runtime struct {
	Session  *mgo.Session
	Mongo    *mgo.DialInfo
	MongoUri string `default:"mongodb://localhost:27017/konveyer"`
	Port     string `default:"3000"`
	LogLevel string `default:"INFO"`
	Debug    bool   `default:"false"`
}

// Initialize Application Runtime
func Initialize() *Runtime {
	var r Runtime
	err := envconfig.Process("", &r)

	mongoinfo, err := mgo.ParseURL(r.MongoUri)
	if err != nil {
		log.Fatalln("Fatal: Unable to Parse Mongo_URI mgo=fail")
		panic(err)
	}

	d, err := mgo.DialWithInfo(mongoinfo)
	if err != nil {
		log.Fatalln("Fatal: Unabled to connect to MongoDB mgo=fail")
		// TODO: Add something to wait until success but retry...
	}
	r.Session = d
	r.Mongo = mongoinfo
	return &r
}

// Stop - method for closing mongodb connections
func (r *Runtime) Stop() {
	// Close MongoDB Connection
	r.Session.Close()
}
