package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/JohnCGriffin/overflow"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Mur466/distribcalc2/internal/utils"
	pb "github.com/Mur466/distribcalc2/proto"
)

type Config struct {
	server_host   string
	server_port   string
	max_workers   int
	poll_interval int
}

func NewConfig() Config {
	pserver_host := flag.String("host", "127.0.0.1", "gRPC server host to get job from")
	pserver_port := flag.String("port", "9090", "Port of the host")
	pmax_workers := flag.Int("workers", 3, "Maximum number of workers")
	ppoll_interval := flag.Int("pollint", 15, "Poll interval (seconds)")
	flag.Parse()

	return Config{
		server_host:   *pserver_host,
		server_port:   *pserver_port,
		max_workers:   *pmax_workers,
		poll_interval: *ppoll_interval,
	}
}

func GetAgentId() string {
	return utils.Pseudo_uuid()
}

type Monitor struct {
	mx         sync.Mutex
	inProgress map[string]*pb.OperResponse
}

func NewMonitor() *Monitor {
	return &Monitor{inProgress: map[string]*pb.OperResponse{}}
}
func (m *Monitor) Register(id string, oper *pb.OperResponse) {
	m.mx.Lock()
	defer m.mx.Unlock()
	m.inProgress[id] = oper
}
func (m *Monitor) Unregister(id string) {
	m.mx.Lock()
	defer m.mx.Unlock()
	delete(m.inProgress, id)
}

func (m *Monitor) Count() int {
	m.mx.Lock()
	defer m.mx.Unlock()
	return len(m.inProgress)
}

func (m *Monitor) List() string {
	m.mx.Lock()
	defer m.mx.Unlock()
	res := ""
	for _, oper := range m.inProgress {
		if res != "" {
			res += "; "
		}
		res += fmt.Sprintf("%v %v %v = ?", oper.Operand1, oper.Operator, oper.Operand2)
	}
	return res
}

func Worker(ch_oper <-chan *pb.OperResponse, ch_idle chan<- interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	for oper := range ch_oper {
		func() {
			// регистрируем в мониторе
			thisid := utils.Pseudo_uuid()
			monitor.Register(thisid, oper)
			// отправляем статус
			SendHeartbeat()
			defer func() {
				// завершающие действия
				// убираем из мониторе
				monitor.Unregister(thisid)
				// отправляем статус
				SendHeartbeat()
				// посылаем сигнал, что готовы к новой операции, не дожидаясь таймера
				// в горутине, чтобы избежать deadlock
				go func() { ch_idle <- struct{}{} }()
			}()

			// вычисляем
			log.Printf("Operation calc start %+v", oper)
			// заготовка под результат
			operRes := &pb.OperResultRequest{
				AgentId: agent_id,
				TaskId:  oper.TaskId,
				NodeId:  oper.NodeId,
			}
			var no_overfl bool = true
			switch {
			case oper.Operator == "+":
				operRes.Result, no_overfl = overflow.Add64(int64(oper.Operand1), int64(oper.Operand2))
			case oper.Operator == "-":
				operRes.Result, no_overfl = overflow.Sub64(int64(oper.Operand1), int64(oper.Operand2))
			case oper.Operator == "*":
				operRes.Result, no_overfl = overflow.Mul64(int64(oper.Operand1), int64(oper.Operand2))
			case oper.Operator == "/":
				if oper.Operand2 == 0 {
					operRes.Status = "error"
					operRes.Message = "Division by zero"
				} else {
					operRes.Result = int64(oper.Operand1) / int64(oper.Operand2)
					operRes.Result, no_overfl = overflow.Div64(int64(oper.Operand1), int64(oper.Operand2))
				}
			default:
				operRes.Status = "error"
				operRes.Message = "Incorrect operator [" + oper.Operator + "]"
				log.Printf("Incorrect operator [%v] for operation %+v", oper.Operator, oper)
			}
			if !no_overfl {
				operRes.Status = "error"
				operRes.Message = "Overflow"
			}
			if operRes.Status != "error" {
				// изображаем бурную деятельность
				time.Sleep(time.Duration(oper.OperatorDelay) * time.Second)
				operRes.Status = "done"
			}
			// отправляем результат
			SendResult(operRes)

		}()
	}
}

func SendResult(operRes *pb.OperResultRequest) {
	log.Printf("Sending operation result: TaskId: %v, NodeId: %v, Status: %v, Result: %v", operRes.TaskId, operRes.NodeId, operRes.Status, operRes.Result)
	_, err := grpcClient.OperResult(context.Background(), operRes)
	if err != nil {
		log.Printf("failed invoking OperResult: %v", err.Error())
	}
}

func GetOperation() (*pb.OperResponse, bool) {
	operation, err := grpcClient.Oper(context.Background(), &pb.OperRequest{AgentId: agent_id})
	if err != nil {
		log.Printf("failed invoking Oper: %v", err.Error())
		return &pb.OperResponse{}, false
	}
	if operation.Operator == "" {
		// пустой оператор, значит пустая операция, нечего делать
		return &pb.OperResponse{}, false
	}
	log.Printf("Got new operation %+v", operation)
	return operation, true
}

func SendHeartbeat() {
	free_workers := config.max_workers - monitor.Count()
	verbose := monitor.List()
	var status string
	switch {
	case free_workers == 0:
		status = "load 100%"
	case free_workers == config.max_workers:
		status = "idle"
	default:
		status = "busy"
	}
	//log.Printf("Heartbeat, invoke AgentStatus, idle workers %v", free_workers)
	_, err := grpcClient.AgentStatus(context.Background(),
		&pb.AgentStatusRequest{
			AgentId:    agent_id,
			Status:     status,
			TotalProcs: int32(config.max_workers),
			IdleProcs:  int32(free_workers),
			Verbose:    verbose,
		})
	if err != nil {
		log.Printf("failed invoking AgentStatus: %v", err.Error())
	}

}

func DinDon(ch_oper chan<- *pb.OperResponse) {
	// набираем до упора
	for config.max_workers > monitor.Count() {
		if oper, ok := GetOperation(); ok {
			// тут отправить операцию воркеру
			ch_oper <- oper
		} else {
			// не дали операцию, выходим из цикла
			break
		}
	}
	SendHeartbeat()
}

func TaskChecker(ch_oper chan<- *pb.OperResponse, ch_idle <-chan interface{}, ch_stop <-chan interface{}) {
	tick := time.NewTicker(time.Duration(config.poll_interval) * time.Second)
	go func() {
		for {
			select {
			case <-tick.C:
				// таймер прозвенел
				DinDon(ch_oper)
			case <-ch_idle:
				// какой-то воркер освободился
				DinDon(ch_oper)
			case <-ch_stop:
				tick.Stop()
				close(ch_oper)
				return
			}
		}

	}()
	// первый раз не ждем таймера
	DinDon(ch_oper)
}

var config = NewConfig()
var agent_id = GetAgentId()
var free_workers = int32(config.max_workers)
var grpcClient pb.MathOperServiceClient
var monitor *Monitor = NewMonitor()

func main() {

	log.Printf("Agent started with agent_id=%v", agent_id)
	log.Printf("Config %+v", config)

	ch_oper := make(chan *pb.OperResponse)
	ch_stop := make(chan interface{})
	ch_idle := make(chan interface{})
	wg := new(sync.WaitGroup)

	// Создаем воркеров
	for i := 0; i < config.max_workers; i++ {
		wg.Add(1)
		go Worker(ch_oper, ch_idle, wg)
	}

	addr := fmt.Sprintf("%s:%s", config.server_host, config.server_port)
	grpcConn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("could not connect to grpc server %s %s", addr, err.Error())
		os.Exit(1)
	}
	// закроем соединение, когда выйдем из функции
	defer grpcConn.Close()

	grpcClient = pb.NewMathOperServiceClient(grpcConn)

	// поллер заданий, он же heartbeat
	TaskChecker(ch_oper, ch_idle, ch_stop)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	log.Print("Exit on ctrl-c signal")

	log.Printf("Waiting for workers running workers to finish. Busy workers %v of total %v", config.max_workers-int(free_workers), config.max_workers)
	close(ch_stop)
	wg.Wait()
	log.Print("All workers finished")

}
