package mathoper_grpc

import (
	"context"

	"github.com/Mur466/distribcalc2/internal/agent"
	"github.com/Mur466/distribcalc2/internal/errors"
	l "github.com/Mur466/distribcalc2/internal/logger"
	"github.com/Mur466/distribcalc2/internal/task"
	pb "github.com/Mur466/distribcalc2/proto"
	"go.uber.org/zap"
)

type Server struct {
	pb.MathOperServiceServer // сервис из сгенерированного пакета

}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Oper(ctx context.Context, in *pb.OperRequest) (*pb.OperResponse, error) {
	l.Logger.Debug("RPC call Oper", zap.String("AgentId", in.AgentId))

	if in.AgentId == "" {
		return nil, errors.ErrInvalidInput
	}
	thisagent := agent.AgentSeen(in.AgentId)

	for _, t := range task.Tasks {
		if n, ok := t.GetWaitingNodeAndSetProcess(thisagent.AgentId); ok {
			res := pb.OperResponse{
				NodeId:        int32(n.Node_id),
				TaskId:        int32(n.Task_id),
				Operand1:      int32(n.Operand1),
				Operand2:      int32(n.Operand2),
				Operator:      n.Operator,
				OperatorDelay: int32(n.Operator_delay)}
			l.SLogger.Infof("Responce to agent %v with operation %+v", thisagent.AgentId, &res)
			return &res, nil
		}
	}
	l.SLogger.Debugf("Responce to agent %v with empty operation ", thisagent.AgentId)
	return &pb.OperResponse{}, nil
}

func (s *Server) OperResult(ctx context.Context, in *pb.OperResultRequest) (*pb.Empty, error) {
	l.Logger.Info("RPC call OperResult",
		zap.String("AgentId", in.AgentId),
		zap.Int32("TaskId", in.TaskId),
		zap.Int32("NodeId", in.NodeId),
		zap.String("Status", in.Status),
		zap.Int64("Result", in.Result),
		zap.String("Message", in.Message),
	)
	if in.AgentId == "" {
		return nil, errors.ErrInvalidInput
	}
	_ = agent.AgentSeen(in.AgentId)

	for _, t := range task.Tasks {
		if t.Task_id == int(in.TaskId) {
			if len(t.TreeSlice) > int(in.NodeId) {
				if t.TreeSlice[in.NodeId].Agent_id == in.AgentId {
					func() {
						t.Mx.Lock()
						defer t.Mx.Unlock()
						// тут мы 100% одни
						t.SetNodeStatus(int(in.NodeId),
							in.Status,
							task.NodeStatusInfo{Result: in.Result, Message: in.Message},
						)
					}()
				} else {
					// получили результат не от того агента, который забрал операцию
					// просто проигнорируем, вдруг получим еще от кого надо
					// если не получим, то потом по таймауту повторно подадим
					l.Logger.Error("Expected result from one agent, got from another",
						zap.Int32("TaskId", in.TaskId),
						zap.Int32("NodeId", in.NodeId),
						zap.String("Expected agent_id", t.TreeSlice[in.NodeId].Agent_id),
						zap.String("Actual agent_id", in.AgentId),
						zap.String("Status", in.Status),
					)
				}
			} else {
				l.Logger.Error("No such NodeId in Task",
					zap.Int32("TaskId", in.TaskId),
					zap.Int32("NodeId", in.NodeId),
				)

			}
		}
	}
	return &pb.Empty{}, nil
}

func (s *Server) AgentStatus(ctx context.Context, in *pb.AgentStatusRequest) (*pb.Empty, error) {
	// тут уровень DEBUG чтобы не мусорить в логах
	l.Logger.Debug("RPC call AgentStatus",
		zap.String("AgentId", in.AgentId),
		zap.String("Status", in.Status),
		zap.Int32("TotalProcs", in.TotalProcs),
		zap.Int32("IdleProcs", in.IdleProcs),
	)
	if in.AgentId == "" {
		return nil, errors.ErrInvalidInput
	}
	a := agent.AgentSeen(in.AgentId)
	a.Status = in.Status
	a.TotalProcs = int(in.TotalProcs)
	a.IdleProcs = int(in.IdleProcs)
	a.Verbose = in.Verbose
	// заменяем значение в мапе
	agent.AgentUpdate(a)

	return &pb.Empty{}, nil
}
