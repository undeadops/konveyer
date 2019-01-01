package root

// Deployment - Manage Kubernetes Deployment Manifests
type Deployment struct {
	Id        string `json:"id"`
	App       string `json:"app"`
	Namespace string `json:"namespace"`
	Image     string `json:"image"`
}

// DeploymentService - Abstract Deployments away from Database
type DeploymentService interface {
	CreateDeployment(d *Deployment) error
	GetDeployment(app, namespace string) (error, Deployment)
}
