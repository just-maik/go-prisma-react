// Package handlers provides HTTP request handlers.
package handlers

import (
	"backend/prisma/db"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// FormularHandler handles HTTP requests for formulars
type FormularHandler struct {
	db *db.PrismaClient
}

// NewFormularHandler creates a new formular handler
func NewFormularHandler(db *db.PrismaClient) *FormularHandler {
	return &FormularHandler{db: db}
}

// Routes returns the router for formular endpoints
func (h *FormularHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Get("/{id}", h.Get)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)

	// Node relationship endpoints
	r.Post("/{id}/nodes", h.AddNode)
	r.Delete("/{id}/nodes/{nodeId}", h.RemoveNode)
	r.Get("/{id}/nodes", h.ListNodes)
	r.Put("/{id}/nodes/reorder", h.ReorderNodes)

	return r
}

// List godoc
// @Summary List formulars
// @Description Get all formulars
// @Tags formulars
// @Accept json
// @Produce json
// @Success 200 {array} db.FormularModel
// @Router /formulars [get]
func (h *FormularHandler) List(w http.ResponseWriter, r *http.Request) {
	formulars, err := h.db.Formular.FindMany().Exec(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(formulars)
}

// Get godoc
// @Summary Get a formular
// @Description Get formular by ID
// @Tags formulars
// @Accept json
// @Produce json
// @Param id path string true "Formular ID"
// @Success 200 {object} db.FormularModel
// @Failure 404 {string} string "Formular not found"
// @Router /formulars/{id} [get]
func (h *FormularHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	formular, err := h.db.Formular.FindUnique(
		db.Formular.ID.Equals(id),
	).Exec(r.Context())

	if err != nil {
		http.Error(w, "Formular not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(formular)
}

// CreateFormularInput represents the input for creating a formular
type CreateFormularInput struct {
	Name string `json:"name" example:"My Formular"` // The name of the formular
}

// Create godoc
// @Summary Create a formular
// @Description Create a new formular
// @Tags formulars
// @Accept json
// @Produce json
// @Param formular body CreateFormularInput true "Formular to create"
// @Success 201 {object} db.FormularModel
// @Router /formulars [post]
func (h *FormularHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input CreateFormularInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	formular, err := h.db.Formular.CreateOne(
		db.Formular.Name.Set(input.Name),
	).Exec(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(formular)
}

// UpdateFormularInput represents the input for updating a formular
type UpdateFormularInput struct {
	Name *string `json:"name,omitempty" example:"Updated Formular"` // The new name of the formular
}

// Update godoc
// @Summary Update a formular
// @Description Update a formular by ID
// @Tags formulars
// @Accept json
// @Produce json
// @Param id path string true "Formular ID"
// @Param formular body UpdateFormularInput true "Formular updates"
// @Success 200 {object} db.FormularModel
// @Failure 404 {string} string "Formular not found"
// @Router /formulars/{id} [put]
func (h *FormularHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var input UpdateFormularInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := []db.FormularSetParam{}

	if input.Name != nil {
		params = append(params, db.Formular.Name.Set(*input.Name))
	}

	formular, err := h.db.Formular.FindUnique(
		db.Formular.ID.Equals(id),
	).Update(
		params...,
	).Exec(r.Context())

	if err != nil {
		http.Error(w, "Formular not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(formular)
}

// Delete godoc
// @Summary Delete a formular
// @Description Delete a formular by ID
// @Tags formulars
// @Accept json
// @Produce json
// @Param id path string true "Formular ID"
// @Success 204 "No Content"
// @Failure 404 {string} string "Formular not found"
// @Router /formulars/{id} [delete]
func (h *FormularHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, err := h.db.Formular.FindUnique(
		db.Formular.ID.Equals(id),
	).Delete().Exec(r.Context())

	if err != nil {
		http.Error(w, "Formular not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddNodeInput represents the input for adding a node to a formular
type AddNodeInput struct {
	NodeID string  `json:"nodeId" example:"123e4567-e89b-12d3-a456-426614174000"`           // The ID of the node to add
	NextID *string `json:"nextId,omitempty" example:"123e4567-e89b-12d3-a456-426614174001"` // Optional ID of the next node in sequence
}

// AddNode godoc
// @Summary Add a node to a formular
// @Description Add a node to a formular's sequence
// @Tags formulars
// @Accept json
// @Produce json
// @Param id path string true "Formular ID"
// @Param node body AddNodeInput true "Node to add"
// @Success 201 {object} db.FormularNodeModel
// @Router /formulars/{id}/nodes [post]
func (h *FormularHandler) AddNode(w http.ResponseWriter, r *http.Request) {
	formularID := chi.URLParam(r, "id")

	var input AddNodeInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	formularNode, err := h.db.FormularNode.CreateOne(
		db.FormularNode.Formular.Link(db.Formular.ID.Equals(formularID)),
		db.FormularNode.Node.Link(db.Node.ID.Equals(input.NodeID)),
	).Exec(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if input.NextID != nil {
		_, err = h.db.FormularNode.FindUnique(
			db.FormularNode.ID.Equals(formularNode.ID),
		).Update(
			db.FormularNode.NextID.Set(*input.NextID),
		).Exec(r.Context())

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(formularNode)
}

// RemoveNode godoc
// @Summary Remove a node from a formular
// @Description Remove a node from a formular's sequence
// @Tags formulars
// @Accept json
// @Produce json
// @Param id path string true "Formular ID"
// @Param nodeId path string true "Node ID"
// @Success 204 "No Content"
// @Failure 404 {string} string "FormularNode not found"
// @Router /formulars/{id}/nodes/{nodeId} [delete]
func (h *FormularHandler) RemoveNode(w http.ResponseWriter, r *http.Request) {
	formularID := chi.URLParam(r, "id")
	nodeID := chi.URLParam(r, "nodeId")

	// Find the formular node first
	formularNode, err := h.db.FormularNode.FindFirst(
		db.FormularNode.FormularID.Equals(formularID),
		db.FormularNode.NodeID.Equals(nodeID),
	).Exec(r.Context())

	if err != nil {
		http.Error(w, "FormularNode not found", http.StatusNotFound)
		return
	}

	// Then delete it
	_, err = h.db.FormularNode.FindUnique(
		db.FormularNode.ID.Equals(formularNode.ID),
	).Delete().Exec(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListNodes godoc
// @Summary List nodes in a formular
// @Description Get all nodes in a formular's sequence
// @Tags formulars
// @Accept json
// @Produce json
// @Param id path string true "Formular ID"
// @Success 200 {array} db.FormularNodeModel
// @Router /formulars/{id}/nodes [get]
func (h *FormularHandler) ListNodes(w http.ResponseWriter, r *http.Request) {
	formularID := chi.URLParam(r, "id")

	nodes, err := h.db.FormularNode.FindMany(
		db.FormularNode.FormularID.Equals(formularID),
	).With(
		db.FormularNode.Node.Fetch(),
	).Exec(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(nodes)
}

// ReorderNodesInput represents the input for reordering nodes in a formular
type ReorderNodesInput struct {
	NodeOrder []string `json:"nodeOrder" example:"['123e4567-e89b-12d3-a456-426614174000', '123e4567-e89b-12d3-a456-426614174001']"` // The ordered list of node IDs
}

// ReorderNodes godoc
// @Summary Reorder nodes in a formular
// @Description Update the sequence of nodes in a formular
// @Tags formulars
// @Accept json
// @Produce json
// @Param id path string true "Formular ID"
// @Param order body ReorderNodesInput true "New node order"
// @Success 200 "OK"
// @Router /formulars/{id}/nodes/reorder [put]
func (h *FormularHandler) ReorderNodes(w http.ResponseWriter, r *http.Request) {
	formularID := chi.URLParam(r, "id")

	var input ReorderNodesInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update the next pointers for each node in the order
	for i := 0; i < len(input.NodeOrder)-1; i++ {
		currentNodeID := input.NodeOrder[i]
		nextNodeID := input.NodeOrder[i+1]

		// Find the formular node first
		formularNode, err := h.db.FormularNode.FindFirst(
			db.FormularNode.FormularID.Equals(formularID),
			db.FormularNode.NodeID.Equals(currentNodeID),
		).Exec(r.Context())

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Then update it
		_, err = h.db.FormularNode.FindUnique(
			db.FormularNode.ID.Equals(formularNode.ID),
		).Update(
			db.FormularNode.NextID.Set(nextNodeID),
		).Exec(r.Context())

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Clear the next pointer for the last node
	// Find the last formular node first
	lastNode, err := h.db.FormularNode.FindFirst(
		db.FormularNode.FormularID.Equals(formularID),
		db.FormularNode.NodeID.Equals(input.NodeOrder[len(input.NodeOrder)-1]),
	).Exec(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Then update it
	_, err = h.db.FormularNode.FindUnique(
		db.FormularNode.ID.Equals(lastNode.ID),
	).Update(
		db.FormularNode.NextID.SetOptional(nil),
	).Exec(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
