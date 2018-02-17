package konveyer

import (
	"log"
	"testing"

	mgo "gopkg.in/mgo.v2"
)

func TestInitialize(t *testing.T) {
	// The tempdir is created so MongoDB has a location to store its files.
	// Contents are wiped once the server stops
	// tempDir, _ := ioutil.TempDir("", "testing")
	// Server.SetPath(tempDir)

	var r Runtime
	r.MongoUri = "mongodb://localhost:27017/konveyer-test"
	mongoinfo, err := mgo.ParseURL(r.MongoUri)
	if err != nil {
		log.Fatalln("Fatal: Unable to Parse Mongo_URI mgo=fail")
		panic(err)
	}
	r.Mongo = mongoinfo
	d, err := mgo.DialWithInfo(mongoinfo)
	if err != nil {
		log.Fatalln("Fatal: Unabled to connect to MongoDB mgo=fail")
		// TODO: Add something to wait until success but retry...
	}
	r.Session = d

	// Make sure we DropDatabase so we make absolutely sure nothing is left or locked while wiping the data and
	// close session
	r.Session.DB(r.Mongo.Database).DropDatabase()
	defer r.Session.Close()
}
