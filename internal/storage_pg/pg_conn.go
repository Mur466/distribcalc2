package storage_pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/Mur466/distribcalc/v2/internal/cfg"
	l "github.com/Mur466/distribcalc/v2/internal/logger"
)

type StoragePg struct {
	conn *pgxpool.Pool
}

func New(cfg *cfg.Config)  (*StoragePg) {

	var dbURL string = fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", 
		cfg.Dbuser,
		cfg.Dbpassword,
		cfg.Dbhost,
		cfg.Dbport,
		cfg.Dbname,
	)

	//fmt.Println(dbURL)
	
	var err error

	conn, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		l.Logger.Fatal(err.Error(),
			zap.String("Dbuser",cfg.Dbuser),
			zap.String("Dbhost",cfg.Dbhost),
			zap.Int("Dbport",cfg.Dbport),
			zap.String("Dbname",cfg.Dbname),
		)
	}

	return &StoragePg{conn: conn}

}

func (s *StoragePg) Stop(){
	s.conn.Close()
}
