package routes

import (
	"github.com/Iqbal-Maulana67/ayo-technical-test/controllers"
	"github.com/Iqbal-Maulana67/ayo-technical-test/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.POST("/auth/login", controllers.LoginAdmin)
	r.POST("/auth/register", controllers.RegisterAdmin)
	admin := r.Group("/", middlewares.Auth())
	{
		admin.POST("/teams", controllers.CreateTeam)
		admin.PUT("/teams/:id", controllers.UpdateTeam)
		admin.DELETE("/teams/:id", controllers.DeleteTeam)
		admin.POST("/players", controllers.CreatePlayer)
		admin.PUT("/players/:id", controllers.UpdatePlayer)
		admin.DELETE("/players/:id", controllers.DeletePlayer)
		admin.POST("/matches", controllers.CreateMatch)
		admin.PUT("/matches/:id", controllers.UpdateMatch)
		admin.DELETE("/matches/:id", controllers.DeleteMatch)
		admin.POST("/matches/:id/result", controllers.SubmitMatchResult)
	}
	// public
	r.GET("/teams", controllers.ListTeams)
	r.GET("/teams/:id", controllers.GetTeam)
	r.GET("/players", controllers.ListPlayers)
	r.GET("/matches", controllers.ListMatches)
	r.GET("/match_result/:id", controllers.GetMatchResult)
}
