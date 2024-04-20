package signout

import (
	"net/http"
	"strings"
	l "github.com/Mur466/distribcalc/v2/internal/logger"
	"github.com/gin-gonic/gin"
)

func MakeHandlerGet(s *Service) func(c *gin.Context) {
	return func(c *gin.Context) {

		var tokenString string
		frombrowser := c.GetHeader("Content-Type") != "application/json"
		if frombrowser {
			tokenString, _ = c.Cookie("token")
		} else {
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				// Split the header to get the token part
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					tokenString = parts[1]
				}
			}
		}

		l.SLogger.Infof("Log out user %v",c.GetString("Username"))
		s.Do(tokenString)
		if frombrowser {
			c.HTML(http.StatusOK, "signin.html", gin.H{})
		}
	}
}
