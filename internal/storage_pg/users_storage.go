package storage_pg

import (
	"context"

	"github.com/Mur466/distribcalc2/internal/entities"
	l "github.com/Mur466/distribcalc2/internal/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	stderrors "errors"

	"github.com/Mur466/distribcalc2/internal/errors"
	"go.uber.org/zap"
)

type Users map[int]entities.User

func (s *StoragePg) GetUser(Username string) *entities.User {
	u := entities.User{}
	err := s.conn.QueryRow(context.Background(), `
	SELECT user_id, username, password_hash 
	  FROM users 
	 WHERE username=$1
	`, Username).Scan(&u.Id, &u.Username, &u.PasswordHash)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil
		} else {
			l.Logger.Error("Error on select from USERS",
				zap.String("error", err.Error()),
				zap.String("username", string(u.Username)),
			)
			return nil
		}
	}
	return &u
}

func (s *StoragePg) AddUser(u *entities.User) (int, error) {
	var newid int
	err := s.conn.QueryRow(context.Background(), `
	INSERT INTO users (username, password_hash) 
	VALUES ($1, $2) 
	RETURNING user_id;
	`, u.Username, u.PasswordHash).Scan(&u.Id)
	if err != nil {
		var pgerr *pgconn.PgError
		if stderrors.As(err, &pgerr) && pgerr.ConstraintName == "users_username_uk" {
			l.Logger.Info("Attempt to insert duplicate user into USERS",
				zap.String("username", string(u.Username)),
			)
			return 0, errors.ErrDuplicateUsername
		}
		l.Logger.Error("Error on insert to USERS",
			zap.String("error", err.Error()),
			zap.Int("user_id", newid),
			zap.String("username", string(u.Username)),
		)
		return 0, errors.ErrDatabaseInternalError
	}
	return newid, nil
}
