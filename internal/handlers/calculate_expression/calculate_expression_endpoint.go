package calculateexpression

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	
	l "github.com/Mur466/distribcalc/v2/internal/logger"
)

func MakeHandlerPost(s *Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		var extexpr ExtExpr

		frombrowser := c.GetHeader("Content-Type") != "application/json"

		if frombrowser {
			// вызвали из html-формы
			extexpr.Expr = c.PostForm("expr")
			extexpr.Ext_id = c.PostForm("ext_id")
		} else {
			// пытаемся через json
			if err := c.BindJSON(&extexpr); err != nil {
				l.Logger.Info("Error JSON",
					zap.String("JSON", err.Error()))
				return
			}
		}
		extexpr.User = s.repo.GetUser(c.GetString("username"))

		res, err := s.Do(&extexpr)

		if frombrowser {
			var errorMsg string
			if err != nil {
				errorMsg = "?error="+err.Error()
			}
			http.Redirect(c.Writer, c.Request, "/tasks"+errorMsg, http.StatusSeeOther)
		} else {
			// ответим на json
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			} else {
				c.IndentedJSON(http.StatusOK, res)
			}
		}

	}
}
