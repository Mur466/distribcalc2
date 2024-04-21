package entities

import "github.com/Mur466/distribcalc2/internal/errors"

type User struct {
	Id       int
	Username string
	//	Password string
	PasswordHash []byte
}

type Repo interface {
	GetUser(username string) *User
}

func GetUserId(username string, r Repo) (int, error) {
	u := r.GetUser(username)
	if u == nil {
		return 0, errors.ErrInvalidCreds
	}
	return u.Id, nil
}
