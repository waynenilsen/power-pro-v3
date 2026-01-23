// Package main provides the entry point for the PowerPro server.
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/waynenilsen/power-pro-v3/internal/database"
	"github.com/waynenilsen/power-pro-v3/internal/server"
)

func main() {
	// Parse flags
	port := flag.Int("port", 8080, "Server port")
	dbPath := flag.String("db", "powerpro.db", "Database file path")
	migrationsPath := flag.String("migrations", "migrations", "Migrations directory path")
	flag.Parse()

	// Open database
	db, err := database.Open(database.Config{
		Path:           *dbPath,
		MigrationsPath: *migrationsPath,
	})
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create and start server
	srv := server.New(server.Config{
		Port: *port,
		DB:   db,
	})

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down server...")
		_ = srv.Stop(context.Background())
	}()

	log.Printf("Starting server on port %d", *port)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
