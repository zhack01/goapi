package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/zhack01/goapi/services/api"
)

func main() {
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	serviceAPI := api.ServiceAPI{}

	r := gin.Default()

	Api := r.Group("/api")
	{
		//Login
		Api.POST("/login", serviceAPI.LogIn())
	}

	if os.Getenv("ENV") == "local" {
		r.Run("127.0.0.1:1111")
	} else {
		r.Run(":1111")
	}
}
