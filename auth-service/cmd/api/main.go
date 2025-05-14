package main

import (
	"auth/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const webPort = "5003"


type Config struct {
	DB *sql.DB
	Models data.Models
}


func main() {
	//connect to database
	dsn := os.Getenv("DSN")
	conn, err := openDB(dsn)
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()
	
	
	app := Config{
		DB: conn,
		Models: data.New(conn),
	}


	server := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	log.Printf("Starting auth service on port %s\n", webPort)
	err = server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}


func openDB(dsn string) (*sql.DB, error) {
	conn, err := sql.Open("pgx", dsn)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, err
	}

	// Verify the connection works
    if err := conn.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

    fmt.Println("Connected to database")

	return conn, nil
}