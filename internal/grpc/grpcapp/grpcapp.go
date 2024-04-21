package grpcapp

import (
	"fmt"
	"net"

	"github.com/Mur466/distribcalc2/internal/cfg"
	"github.com/Mur466/distribcalc2/internal/grpc/mathoper_grpc"
	l "github.com/Mur466/distribcalc2/internal/logger"
	pb "github.com/Mur466/distribcalc2/proto"
	"go.uber.org/zap"
	"go.uber.org/zap/zapgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

type grpcApp struct {
	gRPCServer *grpc.Server
	addr       string
}

func New(
	cfg *cfg.Config,
) *grpcApp {

	grpclog.SetLoggerV2(zapgrpc.NewLogger(l.Logger))
	host := "" //"localhost"
	addr := fmt.Sprintf("%s:%v", host, cfg.GrpcPort)
	// создадим сервер grpc
	grpcServer := grpc.NewServer()
	// объект структуры, которая содержит реализацию
	// серверной части MathOper
	mathOperServiceServer := mathoper_grpc.NewServer()
	// зарегистрируем нашу реализацию сервера
	pb.RegisterMathOperServiceServer(grpcServer, mathOperServiceServer)

	return &grpcApp{addr: addr, gRPCServer: grpcServer}

}

func (a *grpcApp) Run() error {

	lis, err := net.Listen("tcp", a.addr) // будем ждать запросы по этому адресу

	if err != nil {
		l.Logger.Error("error starting tcp listener", zap.String("address", a.addr), zap.Error(err))
		return err
	}
	l.Logger.Info("tcp listener started at address ", zap.String("address", a.addr))

	// запустим grpc сервер
	l.Logger.Info("starting grpc serve... ")
	if err := a.gRPCServer.Serve(lis); err != nil {
		l.Logger.Error("error serving grpc", zap.Error(err))
		return err
	}
	return nil
}

// MustRun runs gRPC server and panics if any error occurs.
func (a *grpcApp) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Stop stops gRPC server.
func (a *grpcApp) Stop() {
	l.Logger.Info("stopping gRPC server", zap.String("address", a.addr))
	a.gRPCServer.GracefulStop()
}
