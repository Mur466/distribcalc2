package register

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func MakeHandlerPost(s *Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		req := request{
			Username: c.PostForm("Username"),
			Password1: c.PostForm("Password1"),
			Password2: c.PostForm("Password2"),
		}

		err := s.Do(&req)
		if err != nil {
			// показываем страницу с сообщением об ошибке
			c.HTML(
				http.StatusOK,
				"register.html",
				gin.H{"Message": err.Error(),
					"Username": req.Username,
				},
			)
		} else {
			http.Redirect(c.Writer, c.Request, "/signin", http.StatusSeeOther)
		}

	}
}

func MakeHandlerGet() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{})
	}

}