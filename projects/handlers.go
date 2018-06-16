package projects

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/ghodss/yaml"

	"github.com/gin-gonic/gin"
	"github.com/undeadops/konveyer/deployments"
	"github.com/undeadops/konveyer/repos/ecr"
)

// Payload - Temp Store
type Payload struct {
	Data json.RawMessage `json:"data"`
}

// GetProjects - GET /projects
func GetProjects(c *gin.Context) {
	message, err := ListProjects()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Error Listing Projects"})
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "projects": message})
}

// CreateProjects - POST /projects
func CreateProjects(c *gin.Context) {

	var p Project
	c.BindJSON(&p)

	repo, err := ecr.CreateContainerRepo(p.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Unable to create ECR Repo"})
	}
	p.Repo = repo
	if err := CreateProject(p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Unable to create Project"})
	}

	c.JSON(http.StatusCreated, gin.H{"status": "Created Project"})
}

// CreateDeployment - PUT /:project/:deployment
func CreateDeployment(c *gin.Context) {
	project := c.Param("project")
	deployment := c.Param("deployment")

	var rawPayload Payload
	//c.Bind(&rawPayload)
	c.BindJSON(&rawPayload)

	// Test Context given...
	payload, err := yaml.YAMLToJSON(rawPayload.Data)
	if err != nil {
		log.Printf("Error Converting YAML to JSON, %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Error", "message": "There was an error Handling Kubernetes Client Config"})
	}

	log.Printf("%s\n\n", payload)

	client, err := deployments.NewClient(payload)
	if err != nil {
		log.Printf("Create Deployment Failed Client Syntax Check, error: %s\n", err.Error())
		c.JSON(http.StatusNotAcceptable, gin.H{"status": "Failed", "message": "Kubernetes Config Failed to Produce Proper Client"})
	}
	var updateProject Project
	updateProject.Deployments = make(map[string]*restclient.Config)
	updateProject.Deployments[deployment] = client
	if err := UpdateProject(project, updateProject); err != nil {
		c.JSON(http.StatusNotModified, gin.H{"status": "Failed", "message": "Project Failed to Update"})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Project Deployment Added"})
	}
}

// DeployHandler - Initiate Deployment of new Software for project env
func DeployHandler(c *gin.Context) {
	project := c.Param("project")
	deployment := c.Param("deployment")
	image := c.Param("image")

	p, err := GetProject(project)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": 404, "message": "Project Not Found"})
	}

	if _, ok := p.Deployments[deployment]; ok {
		// Caching could make this step faster...
		if k := ecr.VerifyContainerImageVersion(p.Repo, image); k {
			log.Printf("DeployHandler: I match a deployment: %s -> %s\n", deployment, image)
			c.JSON(http.StatusOK, gin.H{"status": 200, "message": "I will Deploy that image"})
		}
	}
}

// DescribeProjectHandler - GET /project/{projectName}
func DescribeProjectHandler(c *gin.Context) {

	project := c.Param("project")

	message, err := GetProject(project)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Invalid Project"})
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "project": message})
}

// DescribeProjectImages - List Containers built for project
func DescribeProjectImages(c *gin.Context) {
	project := c.Param("project")

	message, err := ecr.FetchContainerVersions(project)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Unable to fetch Container Images"})
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "images": message})
}

// UpdateProjectHandler - Update an existing project
func UpdateProjectHandler(c *gin.Context) {
	project := c.Param("project")

	var updateProject Project
	c.BindJSON(&updateProject)

	if err := UpdateProject(project, updateProject); err != nil {
		c.JSON(http.StatusNotModified, gin.H{"status": "Failed", "message": "Project Failed to Update"})
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Project Updated"})
}

// // DescribeProjectDeployments - Test out or kubernetes client...
func DescribeProjectDeployments(c *gin.Context) {
	ctx := context.Background()
	project := c.Param("project")

	proj, err := GetProject(project)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "Not Found", "message": "Project Not Found"})
	}
	projDeployments := make(map[string]interface{})
	for k, v := range proj.Deployments {
		d, err := deployments.GetCurrentDeployments(ctx, v)
		if err != nil {
			// This will need to handle kube connect errors...
			c.JSON(http.StatusNotFound, gin.H{"status": 404, "message": "Problems accessing deployments for projects"})
		}
		projDeployments[k] = d
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK", "deployments": projDeployments})
}

// DescribeDeployment - GET /:project/:deployment
func DescribeDeployment(c *gin.Context) {
	project := c.Param("project")
	deployment := c.Param("deployment")

	proj, err := GetProject(project)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": 404, "message": "Project Not Found"})
	}
	if _, ok := proj.Deployments[deployment]; ok {
		c.JSON(http.StatusOK, gin.H{"status": "OK", "deployment": proj.Deployments[deployment]})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"status": 404, "message": "Deployment Not Found"})
	}
}

// // SyncProject - Syncronize project
// func SyncProject(c *gin.Context) {
// 	project := c.Param("project")

// 	proj, err := GetProject(project)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"status": "Not Found", "message": "Project Not Found"})
// 	}

// 	_, err = ecr.CreateContainerRepo(project)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"status": "Error Creating Container Repo"})
// 	}

// 	// // need to handle multiples here...
// 	// if err := deployments.CreateNameSpace(proj.Deployments["dev"]); err != nil {
// 	// 	c.JSON(http.StatusInternalServerError, gin.H{"status": "Error Creating Project Deployment Namespace"})
// 	// }

// }

// DeleteDeployment - DELETE /:project/:deployment
func DeleteDeployment(c *gin.Context) {

}
