package main

import (
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/smallfish/simpleyaml"
)

// DesiredState - State Stored on local FileSystem or in Git
type DesiredState struct {
	Filepath  string `json:"filepath"`
	Namespace string `json:"namespace"`
	Image     string `json:"image"`
}

func (s *Configuration) status(c *gin.Context) {
	u, err := url.Parse(s.Source)
	if err != nil {
		panic(err)
	}

	switch u.Scheme {
	case "file":
		files := parseDeploys(u.Path)
		c.JSON(200, gin.H{"ok": &files})
	case "git":
		c.JSON(200, gin.H{"ok": "GET git PATH"})
	default:
		c.JSON(404, gin.H{"Not Found": "GET unknown PATH"})
	}
}

func scanDir(root string) []string {
	var files []string

	rootModified := root + "/manifests/"

	err := filepath.Walk(rootModified, func(path string, info os.FileInfo, err error) error {
		fi, err := os.Lstat(path)
		if err != nil {
			log.Fatal(err)
		}
		switch mode := fi.Mode(); {
		case mode.IsRegular():

			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}

func loadDeployFile(filename string) *simpleyaml.Yaml {
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	yaml, err := simpleyaml.NewYaml(source)
	if err != nil {
		panic(err)
	}
	return yaml
}

func parseDeploys(root string) []*DesiredState {
	files := scanDir(root)
	rootModified := root + "/manifests/"

	var ds []*DesiredState

	for _, f := range files {
		x := strings.Replace(f, rootModified, "", -1)
		s := strings.Split(x, "/")
		namespace, _ := s[0], s[1]
		deploy := &DesiredState{}

		deploy.Filepath = f
		deploy.Namespace = namespace

		y := loadDeployFile(f)
		image, err := y.Get("deployments").GetIndex(0).Get("containers").GetIndex(0).Get("image").String()
		if err != nil {
			panic(err)
		}
		deploy.Image = image

		ds = append(ds, deploy)
	}
	return ds
}
