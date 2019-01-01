package root

import (
	"github.com/sirupsen/logrus"
)

// Config the application's configuration
type Config struct {
	Port      string `default:"5000"`
	JWTSecret string `envconfig:"jwt_secret"`
	LogLevel  string `envconfig:"log_level"`
	MongoURI  string `envconfig:"mongo_uri" default:"mongodb://foo:bar@localhost:27017"`
	MongoDB   string `envconfig:"mongo_db" default:"konveyer"`
	Logger    *logrus.Logger
}
