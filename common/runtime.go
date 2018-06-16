package common

import (
	"log"

	"github.com/kelseyhightower/envconfig"
	mgo "gopkg.in/mgo.v2"
)

// Runtime - Application Context instead of globals
type Runtime struct {
	Session  *mgo.Session
	Mongo    *mgo.DialInfo
	MongoUri string `default:"mongodb://localhost:27017/konveyer"`
	Port     string `default:"3000"`
	Region   string `default:"us-east-1"`
	LogLevel string `default:"INFO"`
	Debug    bool   `default:"false"`
}

// Config - Runtime Object for current state
var Config *Runtime

// Initialize Application Runtime
func Init() *Runtime {
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
	Config = &r
	return Config
}

// GetRuntime - Return Runtime object
func GetRuntime() *Runtime {
	return Config
}
