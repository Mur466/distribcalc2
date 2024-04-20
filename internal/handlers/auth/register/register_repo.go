package register

import "github.com/Mur466/distribcalc/v2/internal/entities"

type Repo interface {
	GetUser(username string) *entities.User
	AddUser(u *entities.User) (id int, err error)
}
