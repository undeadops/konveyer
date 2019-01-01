package mongo

import (
	"github.com/undeadops/konveyer/pkg"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DeploymentService struct {
	collection *mgo.Collection
}

func NewDeploymentService(session *mgo.Session, config *root.Config) *DeploymentService {
	db := config.MongoDB
	collection := session.DB(db).C("deployment")
	collection.EnsureIndex(deploymentModelIndex())
	return &DeploymentService{collection: collection}
}

func (p *DeploymentService) CreateDeployment(u *root.Deployment) error {
	deploy := newDeploymentModel(u)
	return p.collection.Insert(&deploy)
}

func (p *DeploymentService) GetDeployment(app, namespace string) (root.Deployment, error) {
	model := DeploymentModel{}
	err := p.collection.Find(bson.M{"app": app, "namespace": namespace}).One(&model)
	return root.Deployment{
		Id:        model.Id.Hex(),
		App:       model.App,
		Namespace: model.Namespace,
		Image:     model.Image,
	}, err
}
