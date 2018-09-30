package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func (s *Configuration) updateImage(c *gin.Context) {
	namespace := c.Params.ByName("namespace")
	deployment := c.Params.ByName("deployment")
	image := c.Params.ByName("newImage")

	u, err := url.Parse(s.Source)
	if err != nil {

	}

	files := parseDeploys(u.Path)

	modifiedRoot := u.Path + "/manifests/" + namespace

	input, err := ioutil.ReadFile(modifiedRoot + "/" + deployment + ".yaml")
	if err != nil {
		fmt.Println(err)
	}

	output := bytes.Replace(input, []byte("replaceme"), []byte("ok"), -1)

	if err = ioutil.WriteFile("modified.txt", output, 0666); err != nil {
		fmt.Println(err)
	}
}

func (s *Configuration) setAnnotations(c *gin.Context) {

}
