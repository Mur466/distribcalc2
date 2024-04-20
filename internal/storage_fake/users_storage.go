package storage_fake

import (
	"github.com/Mur466/distribcalc/v2/internal/entities"
	"github.com/Mur466/distribcalc/v2/internal/errors"

)

type Users map[int]entities.User

func (s *storage_fake) GetUser(Username string) *entities.User {
	s.mx.Lock()
	defer s.mx.Unlock()
	for _, v := range s.Users {
		if Username == v.Username {
			return &v
		}
	}
	return nil
}

func (s *storage_fake) AddUser(u *entities.User) (id int, err error) {
	s.mx.Lock()
	defer s.mx.Unlock()
	for _, v := range s.Users {
		if v.Username == u.Username {
			return 0, errors.ErrDuplicateUsername
		}
	}
	newid := len(s.Users) + 1
	u.Id = newid
	s.Users[newid] = *u
	return newid, nil
}
