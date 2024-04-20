package signin

import (
	"net/http"

	"github.com/Mur466/distribcalc/v2/internal/errors"
	l "github.com/Mur466/distribcalc/v2/internal/logger"
	"github.com/gin-gonic/gin"
)

func MakeHandlerPost(s *Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		var frombrowser, ok1, ok2 bool
		req := request{}
		// вначале пробуем из html-формы
		req.Username, ok1 = c.GetPostForm("Username")
		req.Password, ok2 = c.GetPostForm("Password")
		if ok1 && ok2 {
			frombrowser = true
		} else {
			frombrowser = false
			// проверим JSON
			if err := c.BindJSON(&req); err != nil {
				l.SLogger.Errorf("Error JSON %+v", err)
			}
		}

		tokenString, err := s.Do(&req)
		if err != nil {
			if frombrowser {
				// очищаем куку токена (отрицательный maxAge)
				c.SetCookie("token", "", -1, "", "", false, true)
				// показываем страницу логина с сообщением об ошибке
				c.HTML(
					http.StatusOK,
					"signin.html",
					gin.H{"Message": err.Error(),
						"Username": req.Username,
					},
				)
			} else {
				if err == errors.ErrInvalidCreds {
					c.Status(http.StatusUnauthorized)
				} else {
					c.Status(http.StatusInternalServerError)
				}
			}
		} else {
			l.SLogger.Infof("Login successful for user %v", req.Username)
			if frombrowser {
				// запоминаем в куку токен и переходим к выражениям
				c.SetCookie("token", tokenString, s.cfg.AuthTTL, "", "", false, true)
				http.Redirect(c.Writer, c.Request, "/tasks", http.StatusSeeOther)
			} else {
				c.IndentedJSON(http.StatusOK, LoginResponse{AccessToken: tokenString})
			}

		}

	}
}

func MakeHandlerGet() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.HTML(200, "signin.html", gin.H{})
	}

}
