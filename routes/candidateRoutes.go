package routes

import (
	controller "github.com/SHUBHAM91285/votingApp_go/controllers"

	"github.com/gin-gonic/gin"
)

func CandidateRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/candidate/:AadharCardNumber", controller.AddCandidate())
	incomingRoutes.PATCH("/candidate/:AadharCardNumber", controller.UpdateCandidate())
	incomingRoutes.DELETE("/candidate/:AadharCardNumber", controller.DeleteCandidate())
	incomingRoutes.POST("/vote/:AadharCardNumber", controller.VoteCandidate())
	incomingRoutes.GET("/vote/count", controller.GetVoteCount())
	incomingRoutes.GET("/candidates", controller.GetCandidate())
}
