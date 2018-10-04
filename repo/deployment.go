package repo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"

	appsv1 "k8s.io/api/apps/v1"
)

// GetDeploymentImage - Read deployment file and return container names and image used
func (repo *Repo) GetDeploymentImage(namespace, deployment string) (map[string]string, error) {
	deployFile := fmt.Sprintf("%s/manifests/%s/%s-deployment.yaml", repo.Path, namespace, deployment)

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

// SetDeploymentImage - Update Deployment file with new image, and commit it back to master
func (repo *Repo) SetDeploymentImage(namespace, deployment string, image string) error {
	deployFile := fmt.Sprintf("%s/manifests/%s/%s-deployment.yaml", repo.Path, namespace, deployment)

	repo.Mutex.Lock()
	defer repo.Mutex.Unlock()

	ybytes, err := ioutil.ReadFile(deployFile)
	if err != nil {
		return err
	}

	jsonbytes, err := yaml.YAMLToJSON(ybytes)
	if err != nil {
		return errors.New("Error Converting Yaml to JSON from deploymentfile")
	}

	var d = appsv1.Deployment{}
	err = json.Unmarshal(jsonbytes, &d)
	if err != nil {
		return errors.New("Error Unmarshaling file into Kubernetes Deployment")
	}

	// Run a split/if on image to see if I need to handle multiple containers images
	d.Spec.Template.Spec.Containers[0].Image = image

	dbytes, err := yaml.Marshal(d)
	if err != nil {
		return errors.New("Error Converting to byte string")
	}

	return nil
}
