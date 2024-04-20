package register

import (
	"regexp"

	"github.com/Mur466/distribcalc/v2/internal/entities"
	"github.com/Mur466/distribcalc/v2/internal/errors"
	l "github.com/Mur466/distribcalc/v2/internal/logger"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo Repo
	}

func NewSvc(repo Repo) *Service {
	return &Service{
		repo: repo,

	}
}

func (s *Service) Do(request *request) error {

	if !usernameIsValid(request.Username) {
		return errors.ErrIvalidUsername
	}
	if !PasswordIsValid(request.Password1) {
		return errors.ErrIvalidPassword
	}
	if request.Password1 != request.Password2 {
		return errors.ErrPasswordMismatch
	}

	u := entities.User{
		Username: request.Username,
	}
	u.PasswordHash, _ = bcrypt.GenerateFromPassword([]byte(request.Password1), bcrypt.DefaultCost)
	
	id, err := s.repo.AddUser(&u)
	if err != nil {
		l.SLogger.Errorf("Register user %v failed: %v", u.Username, err)
		return err
	}
	l.SLogger.Infof("Registered user %v %v",id, u.Username)
	return nil
}

func usernameIsValid(username string) bool {
	
	// начало с буквы, далее только буквы или цифры
	ok, _ := regexp.MatchString("[a-z]([a-z0-9])+", username)
	return ok &&
		len(username) > 0 &&
	    len(username) < 30	   

}
func PasswordIsValid(password string) bool {
	ok, _ := regexp.MatchString(`[a-z0-9!#$%&? "]*`, password)
	return ok &&
		len(password) > 3 &&
	    len(password) < 30		
}