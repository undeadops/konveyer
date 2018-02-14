package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/undeadops/konveyer"
	"github.com/undeadops/konveyer/gen/models"
	"github.com/undeadops/konveyer/gen/restapi/operations/projects"
)

func (a *konveyer.Application) ListProjects(params projects.FindProjectsParams) middleware.Responder {
	names := []string{"puppeteer", "speculator", "bender"}
	var p models.FindProjectsOKBody
	for n := range names {
		name := &models.Project{
			Name: &names[n],
		}
		p = append(p, name)
	}
	// Print out log of something here...
	a.Log.Info("This is a message: %s", a.Port)
	return projects.NewFindProjectsOK().WithPayload(p)
}

func (a *konveyer.Application) DescribeProject(params projects.DescribeProjectsParams) middleware.Responder {
	name := params.ProjectName
	p := &models.Project{
		Name: &name,
	}
	return projects.NewDescribeProjectsOK().WithPayload(p)
}
