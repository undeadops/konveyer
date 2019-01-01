package mongo_test

import (
	"log"
	"testing"

	"github.com/undeadops/konveyer/pkg"
	"github.com/undeadops/konveyer/pkg/mongo"
)

const (
	mongoUrl                 = "localhost:27017"
	dbName                   = "konveyer"
	deploymentCollectionName = "deployment"
)

func Test_DeploymentService(t *testing.T) {
	t.Run("CreateDeployment", createDeployment_should_insert_deployment_into_mongo)
}

func createDeployment_should_insert_deployment_into_mongo(t *testing.T) {
	//Arrange
	//   mongoConfig := root.MongoConfig {
	// 	Ip: "127.0.0.1:27017",
	// 	DbName: "myDb" }
	config := root.Config{
		MongoURI: "mongodb://foo:bar@localhost:27017",
	}
	session, err := mongo.NewSession(&config)
	if err != nil {
		log.Fatalf("Unable to connect to mongo: %s", err)
	}

	defer func() {
		session.DropDatabase("konveyer")
		session.Close()
	}()

	DeploymentService := mongo.NewDeploymentService(session.Copy(), &config)

	testApp := "integration_test_app"
	testNamespace := "integration_test_namespace"
	testImage := "undeadops/webby:latest"
	deploy := root.Deployment{
		App:       testApp,
		Namespace: testNamespace,
		Image:     testImage}

	//Act
	err = DeploymentService.CreateDeployment(&deploy)

	//Assert
	if err != nil {
		t.Errorf("Unable to create deployment: %s", err)
	}

	resultDeployment, _ := DeploymentService.GetDeployment(testApp, testNamespace)

	if resultDeployment.App != deploy.App {
		t.Errorf("Incorrect App. Expected `%s`, Got: `%s`", testApp, resultDeployment.App)
	}
}
