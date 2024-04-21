package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/Mur466/distribcalc2/internal/cfg"
	l "github.com/Mur466/distribcalc2/internal/logger"
)

var Conn *pgxpool.Pool

func InitDb() {

	var dbURL string = fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		cfg.Cfg.Dbuser,
		cfg.Cfg.Dbpassword,
		cfg.Cfg.Dbhost,
		cfg.Cfg.Dbport,
		cfg.Cfg.Dbname,
	)

	//fmt.Println(dbURL)

	var err error

	Conn, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		l.Logger.Fatal(err.Error(),
			zap.String("Dbuser", cfg.Cfg.Dbuser),
			zap.String("Dbhost", cfg.Cfg.Dbhost),
			zap.Int("Dbport", cfg.Cfg.Dbport),
			zap.String("Dbname", cfg.Cfg.Dbname),
		)
	}

}

func ShutdownDb() {
	Conn.Close()
}
