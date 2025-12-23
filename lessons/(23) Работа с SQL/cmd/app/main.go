package main

import (
	"context"
	"fmt"
	"log"

	"OtusGolang/23-sql/internal/app"
	"OtusGolang/23-sql/internal/config"
	"OtusGolang/23-sql/internal/repository/psql"
)

func main() {
	if err := mainImpl(); err != nil {
		log.Fatal(err)
	}
}

func mainImpl() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c, err := config.Read("configs/local.toml")
	if err != nil {
		return fmt.Errorf("cannot read config: %v", err)
	}

	r := new(psql.Repo)
	if err := r.Connect(ctx, c.PSQL.DSN); err != nil {
		return fmt.Errorf("cannot connect to psql: %v", err)
	}
	defer func() {
		if err := r.Close(); err != nil {
			log.Println("cannot close psql connection", err)
		}
	}()

	if err := r.Migrate(ctx, c.PSQL.Migration); err != nil {
		return fmt.Errorf("cannot migrate: %v", err)
	}

	a, err := app.New(r)
	if err != nil {
		return fmt.Errorf("cannot create app: %v", err)
	}

	if err := a.Run(ctx); err != nil {
		return fmt.Errorf("cannot run app: %v", err)
	}

	return nil
}
