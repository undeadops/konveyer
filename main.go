package main

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/undeadops/konveyer/common"
	"github.com/undeadops/konveyer/projects"
)

func main() {
	runtime := common.Init()
	defer runtime.Session.Close()

	// Env Vars or cli here? or... make it part of project directly?
	runtime.Region = "us-east-2"

	r := gin.Default()
	v1 := r.Group("/api/v1")

	projects.Register(v1.Group("/projects"))
	projects.ProjectRegister(v1.Group("/project"))

	testAuth := r.Group("/api/ping")

	testAuth.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run(":" + runtime.Port) // listen and serve on 0.0.0.0:8080

}

// Index - list index
func Index(w http.ResponseWriter, r *http.Request) {
	payload := "{ status: OK }"
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
