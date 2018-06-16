package projects

import (
	"github.com/gin-gonic/gin"
)

// Register - Register URL endpoints prefixed by /projects
func Register(router *gin.RouterGroup) {
	router.GET("", GetProjects)
}

// ProjectRegister - Register URL endpoints prefixed by /project
func ProjectRegister(router *gin.RouterGroup) {
	router.PUT("", CreateProjects)
	router.GET("/:project", DescribeProjectHandler)
	router.PUT("/:project", UpdateProjectHandler)
	router.GET("/:project/images", DescribeProjectImages)
	//router.GET("/:project/deployments", DescribeProjectDeployments)
	router.POST("/:project/env/:deployment", CreateDeployment)
	router.DELETE("/:project/env/:deployment", DeleteDeployment)
	router.GET("/:project/env/:deployment", DescribeDeployment)
	router.PUT("/:project/deploy/:deployment/:image", DeployHandler)
}
