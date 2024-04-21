package agent

import (
	"sync"
	"time"

	"github.com/Mur466/distribcalc2/internal/cfg"
	l "github.com/Mur466/distribcalc2/internal/logger"
	"github.com/Mur466/distribcalc2/internal/task"
	"go.uber.org/zap"
)

type Agent struct {
	AgentId    string
	Status     string
	TotalProcs int
	IdleProcs  int
	FirstSeen  time.Time
	LastSeen   time.Time
	Verbose    string
}

var Agents map[string]*Agent = make(map[string]*Agent)
var mx sync.Mutex

// Удаляем пропавших агентов
func CleanLostAgents() {
	timeout := time.Second * time.Duration(cfg.Cfg.AgentLostTimeout)
	for _, a := range Agents {
		if time.Since(a.LastSeen) > timeout {
			// давно не видели, забудем про него
			l.Logger.Info("Agent lost",
				zap.String("agent_id", a.AgentId),
				zap.Time("Last seen", a.LastSeen),
				zap.Int("timeout sec", cfg.Cfg.AgentLostTimeout),
			)
			// но вначале передадим его задание другим
			for _, t := range task.Tasks {
				if t.Status == "in progress" {
					for _, n := range t.TreeSlice {
						if n.Status == "in progress" &&
							n.Agent_id == a.AgentId {
							t.SetNodeStatus(n.Node_id, "ready", task.NodeStatusInfo{})
						}
					}
				}
			}
			// нет больше такого агента
			mx.Lock()
			delete(Agents, a.AgentId)
			mx.Unlock()
		}
	}
}

func InitAgents() {
	tick := time.NewTicker(time.Second * time.Duration(cfg.Cfg.AgentLostTimeout))
	go func() {
		for range tick.C {
			// таймер прозвенел
			CleanLostAgents()
		}

	}()

}

func NewAgent(AgentId string) *Agent {
	return &Agent{
		AgentId:    AgentId,
		TotalProcs: 0,
		IdleProcs:  0,
		FirstSeen:  time.Now(),
		LastSeen:   time.Now(),
	}
}

func AgentSeen(AgentId string) *Agent {

	mx.Lock()
	defer mx.Unlock()
	thisagent, found := Agents[AgentId]
	if !found {
		// инициализиуем
		thisagent = NewAgent(AgentId)
		l.SLogger.Infof("Registered new agent id: %v", thisagent.AgentId)
	}
	thisagent.LastSeen = time.Now()
	Agents[AgentId] = thisagent
	return thisagent

}

func AgentUpdate(a *Agent) {

	mx.Lock()
	defer mx.Unlock()
	Agents[a.AgentId] = a

}

func (a *Agent) FirstSeenFmt() string {
	if a.FirstSeen.IsZero() {
		return "no info"
	}
	return a.FirstSeen.Format("2006-01-02 15:04:05")
}

func (a *Agent) LastSeenFmt() string {
	if a.LastSeen.IsZero() {
		return "no info"
	}
	return a.LastSeen.Format("2006-01-02 15:04:05")
}
