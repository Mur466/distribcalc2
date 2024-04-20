package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/Mur466/distribcalc/v2/internal/cfg"
	"github.com/Mur466/distribcalc/v2/internal/handlers"
	"github.com/Mur466/distribcalc/v2/internal/handlers/auth/register"
	"github.com/Mur466/distribcalc/v2/internal/handlers/auth/signin"
	"github.com/Mur466/distribcalc/v2/internal/handlers/auth/signout"
	calculateexpression "github.com/Mur466/distribcalc/v2/internal/handlers/calculate_expression"
	"github.com/Mur466/distribcalc/v2/internal/logger"
	"github.com/Mur466/distribcalc/v2/internal/middlewares"
)

// общий интерфейс для всех repo, которые мы передаем в хендлеры
type repo interface {
	signin.Repo
	signout.Repo
	register.Repo
	middlewares.Repo	
}


func InitRouters(repo repo, cfg *cfg.Config) *gin.Engine {
	router := gin.New() // gin.Default()
	router.Use(middlewares.LoggerMiddleware(logger.Logger), gin.Recovery())
	router.LoadHTMLGlob("templates/*")


	router.POST("/give-me-operation", handlers.GiveMeOperation)
	router.POST("/take-operation-result", handlers.TakeOperationResult)
	
	
	router.POST("/set-config", handlers.SetConfig)
	
	router.GET("/signin", signin.MakeHandlerGet())
	router.POST("/signin", signin.MakeHandlerPost(signin.NewSvc(repo, cfg)))
	
	router.GET("/register", register.MakeHandlerGet())
	router.POST("/register", register.MakeHandlerPost(register.NewSvc(repo)))

	// ниже хендлы с авторизацией
	router.Use(middlewares.AuthMiddleware(cfg.Secret, repo))
	router.GET("/agents", handlers.GetAgents)
	router.GET("/tasks", handlers.GetTasks)
	router.GET("/config", handlers.GetConfig)
	router.POST("/calculate-expression", calculateexpression.MakeHandlerPost(calculateexpression.NewSvc(repo)))
	router.GET("/signout", signout.MakeHandlerGet(signout.NewSvc(repo)))
	router.GET("/", handlers.GetTasks)

	return router

}

