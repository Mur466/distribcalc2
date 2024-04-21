package main

import (
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/Mur466/distribcalc/v2/internal/agent"
	"github.com/Mur466/distribcalc/v2/internal/cfg"
	"github.com/Mur466/distribcalc/v2/internal/db"
	"github.com/Mur466/distribcalc/v2/internal/grpc/grpcapp"
	"github.com/Mur466/distribcalc/v2/internal/logger"
	l "github.com/Mur466/distribcalc/v2/internal/logger"
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
	httpaddr :=":" + strconv.Itoa(cfg.Cfg.HttpPort)
	go func() {
		router.Run(httpaddr)
	}()

	// gRPC сервер
	grpcApp := grpcapp.New(myconfig)
	go func() {
		grpcApp.MustRun()
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	// остановим http
	routers.GracefulShutdown(&http.Server{
        Addr:    httpaddr, // Use the appropriate port
        Handler: router,
    })
	// остановим grpc
	grpcApp.Stop()
	l.Logger.Info("Gracefully stopped")


}

/*
curl http://localhost:8080/nodes --include --header "Content-Type: application/json" --request "POST" --data "{\"Astnode_id\": 5, \"task_id\": 1, \"Operand1\": 5, \"Operand2\": 5, \"Operator\": \"*\", \"Operator_delay\" : 20}"
*/
