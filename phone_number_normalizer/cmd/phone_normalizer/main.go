package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"

	"github.com/jkbresee/gophercises/phone_number_normalizer/internal/phone"
)

func main() {
	var (
		host     = flag.String("host", "localhost", "db host")
		port     = flag.Int("port", 5432, "db port")
		user     = flag.String("user", "jkbresee", "db user")
		password = flag.String("pass", "password", "db password")
		dbname   = flag.String("db", "gophercises_phone", "database name")
	)
	flag.Parse()

	adminDSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable",
		*host, *port, *user, *password)
	adminDB, err := sql.Open("postgres", adminDSN)
	if err != nil {
		log.Fatalf("open admin db: %v", err)
	}
	defer adminDB.Close()

	// (Re)create the target database
	if err := phone.ResetAndCreateDB(adminDB, *dbname); err != nil {
		log.Fatalf("prepare db: %v", err)
	}

	appDSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		*host, *port, *user, *password, *dbname)
	db, err := sql.Open("postgres", appDSN)
	if err != nil {
		log.Fatalf("open app db: %v", err)
	}
	defer db.Close()

	db.SetConnMaxLifetime(3 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	ctx := context.Background()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	if err := phone.CreatePhoneTable(ctx, db); err != nil {
		log.Fatalf("create table: %v", err)
	}

	// seed
	if err := phone.SeedIfEmpty(ctx, db); err != nil {
		log.Fatalf("seed: %v", err)
	}

	if err := phone.ProcessAll(ctx, db); err != nil {
		log.Fatalf("process phones: %v", err)
	}

	fmt.Println("done")
}
