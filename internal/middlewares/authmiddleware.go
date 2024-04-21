package middlewares

import (

	"fmt"
	"net/http"
	"strings"

	"github.com/Mur466/distribcalc/v2/internal/entities"
	l "github.com/Mur466/distribcalc/v2/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

)

type Repo interface {
	IsBlacklisted(token string) bool
	GetUser(Username string) *entities.User 
}



func AuthMiddleware(secret string, repo Repo) gin.HandlerFunc {
	return func(c *gin.Context) {
		l.SLogger.Debug("Authmiddleware")
		var tokenString, username string
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

		if tokenString != "" && !repo.IsBlacklisted(tokenString) {
			tokenFromString, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return "", fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}

				return []byte(secret), nil
			})

			if err != nil {
				l.SLogger.Error(err)
			} else {
				claims, ok := tokenFromString.Claims.(jwt.MapClaims)
				if ok {
					username = claims["username"].(string)
				}
			}
		}
		if user_id, err := entities.GetUserId(username, repo); err == nil {
			c.Set("username", username)
			c.Set("user_id", user_id)
			c.Next()
		} else {
			// не прошли авторизацию
			if frombrowser {
				// очищаем куку токена (отрицательный maxAge)
				c.SetCookie("token", "", -1, "", "", false, true)
				// показываем страницу логина с сообщением об ошибке
				c.HTML(
					http.StatusUnauthorized,
					"signin.html",
					gin.H{"Message": "",
						"Username": "",
					},
				)
				c.Abort()
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			}
		}
	}
}

func CheckUsername(username string, repo Repo) bool {
	// todo: надо проверить, что такой пользователь существует в repo
	if username != "" {
		u := repo.GetUser(username)
		if u != nil {
			return true
		}
	}
	return false
}
