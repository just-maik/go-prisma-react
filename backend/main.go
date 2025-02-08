package main

import (
	"backend/handlers"
	"backend/middleware"
	"backend/prisma/db"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// @title           Calculation API
// @version         1.0
// @description     API for managing calculations, formulars, and nodes
// @host            localhost:8081
// @BasePath        /api

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	// Initialize Prisma client
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		return err
	}

	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()

	// Initialize router
	r := chi.NewRouter()

	// Setup middleware
	middleware.Setup(r)

	// Initialize handlers
	swaggerHandler := handlers.NewSwaggerHandler()
	calculationHandler := handlers.NewCalculationHandler(client)
	formularHandler := handlers.NewFormularHandler(client)
	nodeHandler := handlers.NewNodeHandler(client)

	// Mount routes
	r.Mount("/swagger", swaggerHandler.Routes())
	r.Mount("/api/calculations", calculationHandler.Routes())
	r.Mount("/api/formulars", formularHandler.Routes())
	r.Mount("/api/nodes", nodeHandler.Routes())

	// Start server
	fmt.Println("Server running on http://localhost:8080")
	return http.ListenAndServe(":8080", r)
}
