package config

import (
	"log"
	"net/url"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"

	"github.com/undeadops/konveyer/pkg"
)

// LoadConfig loads the config from a file if specified, otherwise from the environment
func LoadConfig() *root.Config {

	var c root.Config

	err := envconfig.Process("", &c)
	if err != nil {
		log.Fatal(err.Error())
	}

	//config, err = populateConfig(config)

	return validateConfig(&c)
}

func validateConfig(c *root.Config) *root.Config {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		// disable, as we set our own
		FullTimestamp: true,
	}
	c.Logger = logger

	_, err := url.Parse(c.MongoURI)
	// Maybe some error checking on mongodb URI string later
	if err != nil {
		logger.Fatal("Invalid formating of Database URI (eg: mongodb://foo:bar@localhost:27017")
	}

	// Fix Port to follow :<port> formating for Go Port binding
	if !strings.HasPrefix(c.Port, ":") {
		c.Port = ":" + c.Port
	}

	return c
}
