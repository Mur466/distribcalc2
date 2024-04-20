package main

import (
	"strconv"

	"github.com/Mur466/distribcalc/v2/internal/agent"
	"github.com/Mur466/distribcalc/v2/internal/cfg"
	"github.com/Mur466/distribcalc/v2/internal/db"
	"github.com/Mur466/distribcalc/v2/internal/logger"
	"github.com/Mur466/distribcalc/v2/internal/routers"
	"github.com/Mur466/distribcalc/v2/internal/storage_pg"
	"github.com/Mur466/distribcalc/v2/internal/task"
)

func main() {

	cfg.InitConfig()
	myconfig := &cfg.Cfg // todo: пробуем обойтись без глобального конфига

	logger.InitLogger()
	defer logger.Logger.Sync()

	db.InitDb()
	defer db.ShutdownDb()
	task.InitTasks()

	agent.InitAgents()

//	repo := storage_fake.New()
	repo := storage_pg.New(myconfig)
	defer repo.Stop()

	router := routers.InitRouters(repo, myconfig)
	router.Run(":" + strconv.Itoa(cfg.Cfg.HttpPort))
}

/*
curl http://localhost:8080/nodes --include --header "Content-Type: application/json" --request "POST" --data "{\"Astnode_id\": 5, \"task_id\": 1, \"Operand1\": 5, \"Operand2\": 5, \"Operator\": \"*\", \"Operator_delay\" : 20}"
*/
