package signout


import (

	l "github.com/Mur466/distribcalc/v2/internal/logger"

)


type Service struct {
	repo Repo
}

func NewSvc(tokenbl Repo) *Service {
	return &Service{repo: tokenbl}
}


func (s *Service) Do(token string) {
	s.repo.AddTokenToBL(token)
	l.SLogger.Infof("tocken blacklisted %v",token)
}
