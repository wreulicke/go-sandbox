package main

import (
	"log"

	_ "statik-migrate/statik"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" //
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	_ "github.com/lib/pq" //
	"github.com/rakyll/statik/fs"
)

//go:generate statik -src=./migrations
func mainInternal() error {
	f, err := fs.New()
	if err != nil {
		return err
	}
	hfs, err := httpfs.New(f, "/")
	if err != nil {
		return err
	}
	m, err := migrate.NewWithSourceInstance("httpfs", hfs, "postgresql://postgres:postgres@localhost:15432/test?sslmode=disable")
	if err != nil {
		return err
	}
	defer m.Close()
	return m.Up()
}

func main() {
	if err := mainInternal(); err != nil {
		log.Fatal(err)
	}
}
