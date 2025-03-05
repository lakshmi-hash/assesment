package main

import (
	"database/sql"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/yourusername/user-management-app/backend/controllers"
	"github.com/yourusername/user-management-app/backend/routes"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORS()) // Enable CORS for Angular frontend

	// SQLite database setup
	db, err := sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatal("Failed to open DB: ", err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL
		)`)
	if err != nil {
		log.Fatal("Failed to create table: ", err)
	}

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)
	userCtrl := controllers.NewUserController(db, sq)

	// Routes
	routes.RegisterRoutes(e, userCtrl)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
