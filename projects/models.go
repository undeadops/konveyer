package projects

import (
	"github.com/undeadops/konveyer/common"
	"github.com/undeadops/konveyer/repos/ecr"
	"github.com/worg/merger"
	"gopkg.in/mgo.v2/bson"
	restclient "k8s.io/client-go/rest"
)

// Project - Object for Storing projects in mongodb
// TODO: Change Repo to be interface,
//       To allow for multiple repos
// TODO: Deployments as interface for different
//       types of deployments
type Project struct {
	ID          bson.ObjectId                 `bson:"_id,omitempty" json:"id"`
	Name        string                        `bson:"name" json:"name"`
	Repo        ecr.Repository                `bson:"repo" json:"repo"`
	Deployments map[string]*restclient.Config `bson:"deployments" json:"deployments"`
}

// ListProjects - Return List of Project
func ListProjects() ([]Project, error) {
	r := common.GetRuntime()

	var projects []Project

	session := r.Session.Copy()
	defer session.Close()
	db := session.DB(r.Mongo.Database)
	err := db.C("projects").Find(bson.M{}).All(&projects)
	return projects, err
}

// CreateProject - Create Project in database
func CreateProject(p Project) error {
	r := common.GetRuntime()
	session := r.Session.Copy()
	defer session.Close()
	db := session.DB(r.Mongo.Database)
	err := db.C("projects").Insert(&p)
	if err != nil {
		return err
	}
	return nil
}

// GetProject - Describe Project and its settings
func GetProject(name string) (Project, error) {
	r := common.GetRuntime()

	session := r.Session.Copy()
	defer session.Close()
	db := session.DB(r.Mongo.Database)
	var project Project
	err := db.C("projects").Find(bson.M{"name": name}).One(&project)
	if err != nil {
		return Project{}, err
	}
	return project, nil
}

// UpdateProject - Update Project with new data...
func UpdateProject(name string, updates Project) error {
	r := common.GetRuntime()
	session := r.Session.Copy()
	defer session.Close()
	db := session.DB(r.Mongo.Database)
	filter := bson.M{"name": name}

	var current Project
	err := db.C("projects").Find(filter).One(&current)

	var target Project
	target = current

	if err := merger.Merge(&target, updates); err != nil {
		return err
	}
	change := bson.M{"$set": &target}
	_, err = db.C("projects").Upsert(filter, change)
	if err != nil {
		return err
	}
	return nil
}
