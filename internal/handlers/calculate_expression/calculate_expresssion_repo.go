package calculateexpression

import "github.com/Mur466/distribcalc2/internal/entities"

type Repo interface {
	GetUser(username string) *entities.User
}
