package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

// SwaggerHandler handles swagger documentation routes
type SwaggerHandler struct{}

// NewSwaggerHandler creates a new swagger handler
func NewSwaggerHandler() *SwaggerHandler {
	return &SwaggerHandler{}
}

// Routes returns the router for swagger endpoints
func (h *SwaggerHandler) Routes() chi.Router {
	r := chi.NewRouter()

	// Serve swagger spec first to ensure it's available
	r.Get("/doc.json", func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers specifically for the swagger spec
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		http.ServeFile(w, r, "docs/swagger.json")
	})

	// Serve swagger docs
	r.Get("/*", httpSwagger.Handler(
		httpSwagger.URL("doc.json"), // Use relative URL
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))

	return r
}
