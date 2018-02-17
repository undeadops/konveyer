package models

import (
	"github.com/undeadops/konveyer"
	"gopkg.in/mgo.v2/bson"
)

// const (
// 	COLLECTION = "projects"
// )

// Project - Object for Storing projects in mongodb
type Project struct {
	ID   bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Name string        `bson:"name" json:"name"`
}

// func (k *konveyer.Runtime) CopySession() *mgo.Session {

// }

// ListProjects - Return List of Project
func ListProjects(k *konveyer.Runtime) ([]Project, error) {
	var projects []Project

	session := k.Session.Copy()
	defer session.Close()
	db := session.DB(k.Mongo.Database)
	err := db.C("projects").Find(bson.M{}).All(&projects)
	return projects, err
}

// CreateProject - Create Project in database
func CreateProject(k *konveyer.Runtime, p Project) error {
	session := k.Session.Copy()
	defer session.Close()
	db := session.DB(k.Mongo.Database)
	err := db.C("projects").Insert(&p)
	if err != nil {
		return err
	}
	return nil
}

// DescribeProject - Describe Project and its settings
func DescribeProject(k *konveyer.Runtime, name string) (Project, error) {
	session := k.Session.Copy()
	defer session.Close()
	db := session.DB(k.Mongo.Database)
	var project Project
	err := db.C("projects").Find(bson.M{"name": name}).One(&project)
	if err != nil {
		return Project{}, err
	}
	return project, nil
}
