package storage_fake

import (
	"sync"

	"github.com/Mur466/distribcalc/v2/internal/entities"
	"golang.org/x/crypto/bcrypt"
)

type storage_fake struct {
	mx *sync.Mutex
	Users
	tokenbl
}

func New() *storage_fake {
	return &storage_fake{
		Users: map[int]entities.User{
			1: {Username: "user1", PasswordHash: func() []byte { r, _ := bcrypt.GenerateFromPassword([]byte("1111"), bcrypt.DefaultCost); return r }()},
			2: {Username: "user2", PasswordHash: func() []byte { r, _ := bcrypt.GenerateFromPassword([]byte("2222"), bcrypt.DefaultCost); return r }()},
		},
		tokenbl: map[string]bool{},
		mx:      &sync.Mutex{},
	}
}

func (s *storage_fake) Stop() {

}
