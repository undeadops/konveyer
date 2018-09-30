package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"

	appsv1 "k8s.io/api/apps/v1"
)

var deployFile = flag.String("file", "", "Kubernetes Deployment File to Load and Modify")
var image = flag.String("image_name", "", "Container Image to update deploy with")
var argoflags = flag.Bool("argo", false, "Add Argo annotations to Deployment file")

func main() {
	flag.Parse()

	f := *deployFile
	deployment, err := getDeploy(f)
	if err != nil {
		panic(err)
	}

	if *image != "" {
		i := *image
		for _, c := range deployment.Spec.Template.Spec.Containers {
			if c.Name == "spidey-stage" {
				deployment.Spec.Template.Spec.Containers[0].Image = i
			}
		}
	}

	if *argoflags {
		labels := make(map[string]string)
		labels["applications.argoproj.io/app-name"] = "spidey-stage"
		deployment.Labels = labels
	}

	y, err := yaml.Marshal(deployment)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(y))
}

func getDeploy(deployFile string) (appsv1.Deployment, error) {
	yamlbytes, err := ioutil.ReadFile(deployFile)
	
	if err != nil {
		panic(err)
	}
	jsonbytes, err := yaml.YAMLToJSON(yamlbytes)
	if err != nil {
		return appsv1.Deployment{}, err
	}

	var d = appsv1.Deployment{}
	err = json.Unmarshal(jsonbytes, &d)
	if err != nil {
		return appsv1.Deployment{}, err
	}

	return d, nil
}
