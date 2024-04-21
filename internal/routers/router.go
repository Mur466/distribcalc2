package routers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/Mur466/distribcalc/v2/internal/cfg"
	"github.com/Mur466/distribcalc/v2/internal/handlers"
	"github.com/Mur466/distribcalc/v2/internal/handlers/auth/register"
	"github.com/Mur466/distribcalc/v2/internal/handlers/auth/signin"
	"github.com/Mur466/distribcalc/v2/internal/handlers/auth/signout"
	calculateexpression "github.com/Mur466/distribcalc/v2/internal/handlers/calculate_expression"
	"github.com/Mur466/distribcalc/v2/internal/logger"
	"github.com/Mur466/distribcalc/v2/internal/middlewares"
	l "github.com/Mur466/distribcalc/v2/internal/logger"
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

func GracefulShutdown(server *http.Server) {
    // Create a timeout context for the shutdown process
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Shutdown the server gracefully
    if err := server.Shutdown(ctx); err != nil {
		l.Logger.Error("http server forced to shutdown",zap.Error(err))
    }

	l.Logger.Info("http server exiting")
}

