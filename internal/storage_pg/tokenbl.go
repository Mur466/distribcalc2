package storage_pg

import (
	"context"

	stderrors "errors"

	l "github.com/Mur466/distribcalc2/internal/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

func (s *StoragePg) AddTokenToBL(token string) {
	/*
		s.mx.Lock()
		defer s.mx.Unlock()
		s.tokenbl[token] = true
	*/

	_, err := s.conn.Exec(context.Background(), `
	INSERT INTO tokenbl (token) 
	VALUES ($1);
	`, token)
	if err != nil {
		var pgerr *pgconn.PgError
		if stderrors.As(err, &pgerr) && pgerr.ConstraintName == "tokenbl_pkey" {
			//no problem
		} else {
			l.Logger.Error("Error on insert to TOKENBL",
				zap.String("error", err.Error()),
				zap.String("username", string(token)),
			)
		}
	}
}

func (s *StoragePg) IsBlacklisted(token string) bool {
	/* это можно использовать как задел под кеширование, чтобы не лазить в БД на каждый запрос
	s.mx.Lock()
	defer s.mx.Unlock()
	_, found := s.tokenbl[token]
	return found
	*/
	var dummy int
	err := s.conn.QueryRow(context.Background(), `
	SELECT 1 
	  FROM tokenbl 
	 WHERE token=$1
	`, token).Scan(&dummy)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false
		} else {
			l.Logger.Error("Error on select from TOKENBL",
				zap.String("error", err.Error()),
				zap.String("token", string(token)),
			)
			return true // в случае сбоя лучше не дадим  доступ хорошему, чем пропустим плохого
		}
	}
	return true

}
