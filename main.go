package main

import (
	"codelit/internal/api"
	"codelit/internal/httpserver"
	"codelit/internal/repositories"
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	dbHost := os.Getenv("DOCKER_INTERNAL")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	server := httpserver.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Database setup
	db, err := sql.Open("postgres", "host="+dbHost+" port="+dbPort+" user="+dbUser+" password="+dbPassword+" dbname="+dbName+" sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	dbRepo := repositories.NewDBRepository(db)

	api.RegisterRoutes(e, dbRepo)

	log.Fatal(e.Start(":8080"))
}
