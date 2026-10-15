package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB() {
	connStr := "postgres://chamba_user:chamba_password@localhost:5432/chamba_db?sslmode=disable"
	var err error
	DB, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("No se pudo conectar a la base de datos local: %v\n", err)
	}

	if err := DB.Ping(context.Background()); err != nil {
		log.Fatalf("Error al hacer ping a PostgreSQL: %v\n", err)
	}

	fmt.Println("Conexión exitosa a PostgreSQL 16 (chamba_db)")
}
