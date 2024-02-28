package main

import (
	"os"

	routes "github.com/SHUBHAM91285/votingApp_go/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.CandidateRoutes(router)
	routes.UserRoutes(router)
	router.Run(":" + port)
}
