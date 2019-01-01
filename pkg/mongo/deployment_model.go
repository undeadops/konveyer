package mongo

import (
	"github.com/undeadops/konveyer/pkg"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DeploymentModel struct {
	Id        bson.ObjectId `bson:"_id,omitempty"`
	App       string
	Namespace string
	Image     string
}

func deploymentModelIndex() mgo.Index {
	return mgo.Index{
		Key:        []string{"app"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
}

func newDeploymentModel(d *root.Deployment) *DeploymentModel {
	deploy := DeploymentModel{App: d.App, Namespace: d.Namespace, Image: d.Image}
	return &deploy
}
