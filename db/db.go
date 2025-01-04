package db

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5" // Import the PostgreSQL driver
)

func OpenDBConnection() *pgx.Conn {

	// Replace with your actual PostgreSQL connection string
	connStr := os.Getenv("POSTGRES_URI")

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the PostgreSQL database: %v", err)
	}

	fmt.Println("Connected to PostgreSQL database")
	return conn
}
