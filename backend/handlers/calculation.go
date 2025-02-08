// Package handlers provides HTTP request handlers.
package handlers

import (
	"backend/prisma/db"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// CalculationHandler handles HTTP requests for calculations
type CalculationHandler struct {
	db *db.PrismaClient
}

// NewCalculationHandler creates a new calculation handler
func NewCalculationHandler(db *db.PrismaClient) *CalculationHandler {
	return &CalculationHandler{db: db}
}

// Routes returns the router for calculation endpoints
func (h *CalculationHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Get("/{id}", h.Get)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)

	// Formular relationship endpoints
	r.Post("/{id}/formulars", h.AddFormular)
	r.Delete("/{id}/formulars/{formularId}", h.RemoveFormular)
	r.Get("/{id}/formulars", h.ListFormulars)
	r.Put("/{id}/formulars/reorder", h.ReorderFormulars)

	return r
}

// List godoc
// @Summary List calculations
// @Description Get all calculations
// @Tags calculations
// @Accept json
// @Produce json
// @Success 200 {array} db.CalculationModel
// @Router /calculations [get]
func (h *CalculationHandler) List(w http.ResponseWriter, r *http.Request) {
	calculations, err := h.db.Calculation.FindMany().Exec(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(calculations)
}

// Get godoc
// @Summary Get a calculation
// @Description Get calculation by ID
// @Tags calculations
// @Accept json
// @Produce json
// @Param id path string true "Calculation ID"
// @Success 200 {object} db.CalculationModel
// @Failure 404 {string} string "Calculation not found"
// @Router /calculations/{id} [get]
func (h *CalculationHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	calculation, err := h.db.Calculation.FindUnique(
		db.Calculation.ID.Equals(id),
	).Exec(r.Context())

	if err != nil {
		http.Error(w, "Calculation not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(calculation)
}

// CreateCalculationInput represents the input for creating a calculation
type CreateCalculationInput struct {
	Name string `json:"name" example:"My Calculation"` // The name of the calculation
}

// Create godoc
// @Summary Create a calculation
// @Description Create a new calculation
// @Tags calculations
// @Accept json
// @Produce json
// @Param calculation body CreateCalculationInput true "Calculation to create"
// @Success 201 {object} db.CalculationModel
// @Router /calculations [post]
func (h *CalculationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input CreateCalculationInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	calculation, err := h.db.Calculation.CreateOne(
		db.Calculation.Name.Set(input.Name),
	).Exec(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(calculation)
}

// UpdateCalculationInput represents the input for updating a calculation
type UpdateCalculationInput struct {
	Name *string `json:"name,omitempty" example:"Updated Calculation"` // The new name of the calculation
}

// Update godoc
// @Summary Update a calculation
// @Description Update a calculation by ID
// @Tags calculations
// @Accept json
// @Produce json
// @Param id path string true "Calculation ID"
// @Param calculation body UpdateCalculationInput true "Calculation updates"
// @Success 200 {object} db.CalculationModel
// @Failure 404 {string} string "Calculation not found"
// @Router /calculations/{id} [put]
func (h *CalculationHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var input UpdateCalculationInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := []db.CalculationSetParam{}

	if input.Name != nil {
		params = append(params, db.Calculation.Name.Set(*input.Name))
	}

	calculation, err := h.db.Calculation.FindUnique(
		db.Calculation.ID.Equals(id),
	).Update(
		params...,
	).Exec(r.Context())

	if err != nil {
		http.Error(w, "Calculation not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(calculation)
}

// Delete godoc
// @Summary Delete a calculation
// @Description Delete a calculation by ID
// @Tags calculations
// @Accept json
// @Produce json
// @Param id path string true "Calculation ID"
// @Success 204 "No Content"
// @Failure 404 {string} string "Calculation not found"
// @Router /calculations/{id} [delete]
func (h *CalculationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, err := h.db.Calculation.FindUnique(
		db.Calculation.ID.Equals(id),
	).Delete().Exec(r.Context())

	if err != nil {
		http.Error(w, "Calculation not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddFormularInput represents the input for adding a formular to a calculation
type AddFormularInput struct {
	FormularID string  `json:"formularId" example:"123e4567-e89b-12d3-a456-426614174000"`       // The ID of the formular to add
	NextID     *string `json:"nextId,omitempty" example:"123e4567-e89b-12d3-a456-426614174001"` // Optional ID of the next formular in sequence
}

// AddFormular godoc
// @Summary Add a formular to a calculation
// @Description Add a formular to a calculation's sequence
// @Tags calculations
// @Accept json
// @Produce json
// @Param id path string true "Calculation ID"
// @Param formular body AddFormularInput true "Formular to add"
// @Success 201 {object} db.CalculationFormularModel
// @Router /calculations/{id}/formulars [post]
func (h *CalculationHandler) AddFormular(w http.ResponseWriter, r *http.Request) {
	calculationID := chi.URLParam(r, "id")

	var input AddFormularInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	calculationFormular, err := h.db.CalculationFormular.CreateOne(
		db.CalculationFormular.Calculation.Link(db.Calculation.ID.Equals(calculationID)),
		db.CalculationFormular.Formular.Link(db.Formular.ID.Equals(input.FormularID)),
	).Exec(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if input.NextID != nil {
		_, err = h.db.CalculationFormular.FindUnique(
			db.CalculationFormular.ID.Equals(calculationFormular.ID),
		).Update(
			db.CalculationFormular.NextID.Set(*input.NextID),
		).Exec(r.Context())

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(calculationFormular)
}

// RemoveFormular godoc
// @Summary Remove a formular from a calculation
// @Description Remove a formular from a calculation's sequence
// @Tags calculations
// @Accept json
// @Produce json
// @Param id path string true "Calculation ID"
// @Param formularId path string true "Formular ID"
// @Success 204 "No Content"
// @Failure 404 {string} string "CalculationFormular not found"
// @Router /calculations/{id}/formulars/{formularId} [delete]
func (h *CalculationHandler) RemoveFormular(w http.ResponseWriter, r *http.Request) {
	calculationID := chi.URLParam(r, "id")
	formularID := chi.URLParam(r, "formularId")

	// Find the calculation formular first
	calculationFormular, err := h.db.CalculationFormular.FindFirst(
		db.CalculationFormular.CalculationID.Equals(calculationID),
		db.CalculationFormular.FormularID.Equals(formularID),
	).Exec(r.Context())

	if err != nil {
		http.Error(w, "CalculationFormular not found", http.StatusNotFound)
		return
	}

	// Then delete it
	_, err = h.db.CalculationFormular.FindUnique(
		db.CalculationFormular.ID.Equals(calculationFormular.ID),
	).Delete().Exec(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListFormulars godoc
// @Summary List formulars in a calculation
// @Description Get all formulars in a calculation's sequence
// @Tags calculations
// @Accept json
// @Produce json
// @Param id path string true "Calculation ID"
// @Success 200 {array} db.CalculationFormularModel
// @Router /calculations/{id}/formulars [get]
func (h *CalculationHandler) ListFormulars(w http.ResponseWriter, r *http.Request) {
	calculationID := chi.URLParam(r, "id")

	formulars, err := h.db.CalculationFormular.FindMany(
		db.CalculationFormular.CalculationID.Equals(calculationID),
	).With(
		db.CalculationFormular.Formular.Fetch(),
	).Exec(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(formulars)
}

// ReorderFormularsInput represents the input for reordering formulars in a calculation
type ReorderFormularsInput struct {
	FormularOrder []string `json:"formularOrder" example:"['123e4567-e89b-12d3-a456-426614174000', '123e4567-e89b-12d3-a456-426614174001']"` // The ordered list of formular IDs
}

// ReorderFormulars godoc
// @Summary Reorder formulars in a calculation
// @Description Update the sequence of formulars in a calculation
// @Tags calculations
// @Accept json
// @Produce json
// @Param id path string true "Calculation ID"
// @Param order body ReorderFormularsInput true "New formular order"
// @Success 200 "OK"
// @Router /calculations/{id}/formulars/reorder [put]
func (h *CalculationHandler) ReorderFormulars(w http.ResponseWriter, r *http.Request) {
	calculationID := chi.URLParam(r, "id")

	var input ReorderFormularsInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update the next pointers for each formular in the order
	for i := 0; i < len(input.FormularOrder)-1; i++ {
		currentFormularID := input.FormularOrder[i]
		nextFormularID := input.FormularOrder[i+1]

		// Find the calculation formular first
		calculationFormular, err := h.db.CalculationFormular.FindFirst(
			db.CalculationFormular.CalculationID.Equals(calculationID),
			db.CalculationFormular.FormularID.Equals(currentFormularID),
		).Exec(r.Context())

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Then update it
		_, err = h.db.CalculationFormular.FindUnique(
			db.CalculationFormular.ID.Equals(calculationFormular.ID),
		).Update(
			db.CalculationFormular.NextID.Set(nextFormularID),
		).Exec(r.Context())

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Clear the next pointer for the last formular
	// Find the last calculation formular first
	lastFormular, err := h.db.CalculationFormular.FindFirst(
		db.CalculationFormular.CalculationID.Equals(calculationID),
		db.CalculationFormular.FormularID.Equals(input.FormularOrder[len(input.FormularOrder)-1]),
	).Exec(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Then update it
	_, err = h.db.CalculationFormular.FindUnique(
		db.CalculationFormular.ID.Equals(lastFormular.ID),
	).Update(
		db.CalculationFormular.NextID.SetOptional(nil),
	).Exec(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
