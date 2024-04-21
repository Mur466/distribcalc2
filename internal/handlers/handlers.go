package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Mur466/distribcalc2/internal/agent"
	"github.com/Mur466/distribcalc2/internal/cfg"
	"github.com/Mur466/distribcalc2/internal/task"
	"github.com/Mur466/distribcalc2/internal/utils"
)

func GetAgents(c *gin.Context) {
	c.HTML(
		200,
		"agents.html",
		gin.H{
			"AuthUser": c.GetString("username"),
			"title":    "Agents",
			"Agents":   agent.Agents,
		},
	)
}

func GetTasks(c *gin.Context) {
	// Могли бы показывать сразу task.Task, но хочется порядок новые вверху и ограничение на странице
	// поэтому берем из БД

	tasks := task.ListTasks(c.GetInt("user_id"), cfg.Cfg.RowsOnPage, 0)
	errorMessage := c.Query("error")
	c.HTML(
		200,
		"tasks.html",
		gin.H{
			"AuthUser":       c.GetString("username"),
			"title":          "Tasks",
			"Tasks":          tasks,
			"NewRandomValue": utils.Pseudo_uuid(),
			"Message":        errorMessage,
		},
	)
}

func GetConfig(c *gin.Context) {
	c.HTML(
		200,
		"config.html",
		gin.H{
			"AuthUser": c.GetString("Username"),
			"title":    "Config",
			"Config":   cfg.Cfg,
		},
	)
}
func ValidateDelay(v string, dflt int) int {
	i, err := strconv.Atoi(v)
	if v != "" && err == nil && i >= 0 {
		return i
	}
	return dflt
}
func SetConfig(c *gin.Context) {
	cfg.Cfg.DelayForAdd = ValidateDelay(c.PostForm("DelayForAdd"), cfg.Cfg.DelayForAdd)
	cfg.Cfg.DelayForSub = ValidateDelay(c.PostForm("DelayForSub"), cfg.Cfg.DelayForSub)
	cfg.Cfg.DelayForMul = ValidateDelay(c.PostForm("DelayForMul"), cfg.Cfg.DelayForMul)
	cfg.Cfg.DelayForDiv = ValidateDelay(c.PostForm("DelayForDiv"), cfg.Cfg.DelayForDiv)
	cfg.RecalcAgentTimeout()
	cfg.Cfg.RowsOnPage = ValidateDelay(c.PostForm("RowsOnPage"), cfg.Cfg.RowsOnPage)
	http.Redirect(c.Writer, c.Request, "/config", http.StatusSeeOther)
}
