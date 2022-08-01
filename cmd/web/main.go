package main

import (
	_ "context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"deeptown.com/deepsearch/pkg/models/postgresql"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	products *postgresql.ProductModel
}

func main() {
	addr := flag.String("addr", ":4000", "Сетевой адрес веб-сервера")
	dsn := "postgres://artem:artem1906@localhost:5432/deeptown?sslmode=disable"
	flag.Parse()

	dbpool, err := openDB(dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		products: &postgresql.ProductModel{DB: dbpool},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Запуск сервера на %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	dbpool, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = dbpool.Ping(); err != nil {
		return nil, err
	}
	return dbpool, nil
}
