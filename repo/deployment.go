package repo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"github.com/sirupsen/logrus"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	appsv1 "k8s.io/api/apps/v1"
)

// GetDeploymentImage - Read deployment file and return container names and image used
func (repo *Repo) GetDeploymentImage(namespace, appname string) (map[string]string, error) {
	deployFile := fmt.Sprintf("%s/manifests/%s/%s-deployment.yaml", repo.Path, namespace, appname)
	fmt.Printf("looking for deploy file at: %s", deployFile)
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
func (repo *Repo) SetDeploymentImage(namespace, deployment, image string) error {
	deployFile := fmt.Sprintf("%s/manifests/%s/%s-deployment.yaml", repo.Path, namespace, deployment)
	fmt.Println(deployFile)
	repo.Mutex.Lock()
	defer repo.Mutex.Unlock()

	// err := repo.PullRepo()
	// if err != nil {
	// 	return errors.New("Error Pulling latest repo update before modifying")
	// }

	w, err := repo.Clone.Worktree()
	if err != nil {
		return errors.New("Error Inititating Worktree in git repo")
	}

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

	err = ioutil.WriteFile(deployFile, dbytes, 0644)
	if err != nil {
		return errors.New("Error Writing modified manifest")
	}

	fmt.Println(deployFile)

	_, err = w.Add(repo.makeRelative(deployFile))
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return errors.New("Error adding updated deployment file to git")
	}

	// status, err := w.Status()
	// if err != nil {
	// 	return errors.New("Error gathering status for updated deployment file")
	// }

	commit, err := w.Commit("Konveyer Updated Image: "+repo.makeRelative(deployFile), &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Konveyer.sh",
			Email: "svc@konveyer.sh",
			When:  time.Now(),
		},
	})
	if err != nil {
		return errors.New("Error Commiting updated deployment file")
	}

	repo.Logger.WithFields(
		logrus.Fields{"git_path": repo.Path, "git_hash": commit},
	).Info("Git Commited, Head Updated")

	if err := repo.PushRepo(); err != nil {
		return errors.New("Error pushing Repo up to git-repo")
	}

	return nil
}

func (repo *Repo) makeRelative(f string) string {
	file := strings.Replace(f, repo.Path, "", 1)

	if strings.HasPrefix(file, "/") {
		file = strings.Replace(file, "/", "", 1)
	}
	return file
}
