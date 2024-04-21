package signin

import "github.com/Mur466/distribcalc2/internal/entities"

type Repo interface {
	GetUser(Username string) *entities.User
}
