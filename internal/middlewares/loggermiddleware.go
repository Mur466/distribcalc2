package middlewares

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
	  // Log the request
	  logger.Info("Incoming request",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
	  )
	  c.Next()
	}
  }