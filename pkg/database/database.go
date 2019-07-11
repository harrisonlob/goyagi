package database

import (
	"crypto/tls"
	"fmt"

	"github.com/go-pg/pg"
	"github.com/harrisonlob/goyagi/pkg/config"
)

func New(cfg config.Config) (*pg.DB, error) {
	addr := fmt.Sprintf("%s:%d", cfg.DatabaseHost, cfg.DatabasePort)

	opts := &pg.Options{
		Addr:     addr,
		User:     cfg.DatabaseUser,
		Password: cfg.DatabasePassword,
		Database: cfg.DatabaseName,
	}

	if cfg.DatabaseTLS {
		opts.TLSConfig = &tls.Config{ServerName: cfg.DatabaseHost}
	}

	db := pg.Connect(opts)

	_, err := db.Exec("SELECT 1")
	if err != nil {
		return nil, err
	}

	return db, nil
}
