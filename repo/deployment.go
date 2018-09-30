package repo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"

	appsv1 "k8s.io/api/apps/v1"
)

// GetDeploymentImage - Read deployment file and return container names and image used
func (repo *Repo) GetDeploymentImage(namespace, deployment string) (map[string]string, error) {
	deployFile := fmt.Sprintf("%s/manifests/%s/%s.deploy.yaml", repo.Path, namespace, deployment)

	repo.Mutex.RLock()
	defer repo.Mutex.RUnlock()

	ybytes, err := ioutil.ReadFile(deployFile)
	if err != nil {
		return make(map[string]string), err
	}

	jsonbytes, err := yaml.YAMLToJSON(ybytes)
	if err != nil {
		return make(map[string]string), err
	}

	var d = appsv1.Deployment{}
	err = json.Unmarshal(jsonbytes, &d)
	if err != nil {
		return make(map[string]string), err
	}

	containers := make(map[string]string)

	for _, c := range d.Spec.Template.Spec.Containers {
		containers[c.Name] = c.Image
	}

	return containers, nil
}
