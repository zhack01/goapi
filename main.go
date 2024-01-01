package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// db := mysql.ConnectDB()
	// txndb := mysql.SecondConnectDB()
}
