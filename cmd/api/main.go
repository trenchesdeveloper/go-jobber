package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/trenchesdeveloper/jobber/config"
	db "github.com/trenchesdeveloper/jobber/internal/db/sqlc"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.LoadConfig(".")

	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	app := &server{
		config: cfg,
		logger: logger,
	}

	// connect to the database
	conn, err := sql.Open(cfg.DBdriver, cfg.DBSource)
	if err != nil {
		log.Fatal(err)
	}
	conn.SetMaxOpenConns(30)
	conn.SetMaxIdleConns(30)
	err = conn.PingContext(ctx)
	if err != nil {
		logger.Fatal(err)
	}
	defer conn.Close()
	logger.Info("database connected")

	// create a new store
	store := db.NewStore(conn)

	// create a new server
	app.store = store

	mux := app.mount()

	if err := app.start(mux); err != nil {
		logger.Fatal(err)
	}
}
