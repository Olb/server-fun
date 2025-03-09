package main

import (
	"log"

	"olbcloud.com/webapi/internal/config"
	"olbcloud.com/webapi/internal/database"
	"olbcloud.com/webapi/internal/database/mongodb"
	"olbcloud.com/webapi/internal/database/postgresql"
	"olbcloud.com/webapi/internal/handlers"
	"olbcloud.com/webapi/internal/server"
	"olbcloud.com/webapi/internal/services"
)

func main() {
	cfg := config.LoadConfig()

	var db database.DB
	var err error

	switch cfg.DBType {
	case "postgresql":
		db, err = postgresql.NewPostgreSQL(cfg.PostgresURL)
	case "mongodb":
		db, err = mongodb.NewMongoDB(cfg.MongoDBURL)
	default:
		log.Fatal("Invalid DB_TYPE. Must be 'postgresql' or 'mongodb'")
	}

	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	postService := services.NewPostService(db)
	hd := handlers.NewHandlers(postService)

	mux := server.NewServer(hd)
	server.StartServer(mux)
}
