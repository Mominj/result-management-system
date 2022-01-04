package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

func main() {
	if err := Migrate(); err != nil {
		log.Fatalln("error occur when migrate database", err)
	}
}

func Migrate() error {
	db, DBErr := sql.Open("postgres", "user=postgres password=momin1234 dbname=glogin sslmode=disable")
	if DBErr != nil {
		log.Fatalln("error while open to database", DBErr)
	}
	defer db.Close()
	flag.Parse()
	args := flag.Args()
	driver := "postgres"
	if err := goose.SetDialect(driver); err != nil {
		return fmt.Errorf("failed to set goose dialect: %v", err)
	}
	if len(args) == 0 {
		return errors.New("expected at least one arg")
	}
	command := args[0]

	migrationDir := "migrations/sql"
	if err := goose.Run(command, db, migrationDir, args[1:]...); err != nil {
		return fmt.Errorf("goose run: %v", err)
	}
	return nil
}
