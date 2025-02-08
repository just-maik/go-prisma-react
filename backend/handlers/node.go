// Package handlers provides HTTP request handlers.
package handlers

import (
	"backend/prisma/db"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// NodeHandler handles HTTP requests for nodes
type NodeHandler struct {
	db *db.PrismaClient
}

// NewNodeHandler creates a new node handler
func NewNodeHandler(db *db.PrismaClient) *NodeHandler {
	return &NodeHandler{db: db}
}

// Routes returns the router for node endpoints
func (h *NodeHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Get("/{id}", h.Get)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)

	return r
}

// List godoc
// @Summary List nodes
// @Description Get all nodes
// @Tags nodes
// @Accept json
// @Produce json
// @Success 200 {array} db.NodeModel
// @Router /nodes [get]
func (h *NodeHandler) List(w http.ResponseWriter, r *http.Request) {
	nodes, err := h.db.Node.FindMany().Exec(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(nodes)
}

// Get godoc
// @Summary Get a node
// @Description Get node by ID
// @Tags nodes
// @Accept json
// @Produce json
// @Param id path string true "Node ID"
// @Success 200 {object} db.NodeModel
// @Failure 404 {string} string "Node not found"
// @Router /nodes/{id} [get]
func (h *NodeHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	node, err := h.db.Node.FindUnique(
		db.Node.ID.Equals(id),
	).Exec(r.Context())

	if err != nil {
		http.Error(w, "Node not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(node)
}

// CreateNodeInput represents the input for creating a node
type CreateNodeInput struct {
	Name     string `json:"name" example:"My Node"`      // The name of the node
	NodeData string `json:"nodeData" example:"raw data"` // The data associated with the node
}

// Create godoc
// @Summary Create a node
// @Description Create a new node
// @Tags nodes
// @Accept json
// @Produce json
// @Param node body CreateNodeInput true "Node to create"
// @Success 201 {object} db.NodeModel
// @Router /nodes [post]
func (h *NodeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input CreateNodeInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	node, err := h.db.Node.CreateOne(
		db.Node.Name.Set(input.Name),
		db.Node.NodeData.Set(input.NodeData),
	).Exec(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(node)
}

// UpdateNodeInput represents the input for updating a node
type UpdateNodeInput struct {
	Name     *string `json:"name,omitempty" example:"Updated Node"`     // The new name of the node
	NodeData *string `json:"nodeData,omitempty" example:"updated data"` // The new data for the node
}

// Update godoc
// @Summary Update a node
// @Description Update a node by ID
// @Tags nodes
// @Accept json
// @Produce json
// @Param id path string true "Node ID"
// @Param node body UpdateNodeInput true "Node updates"
// @Success 200 {object} db.NodeModel
// @Failure 404 {string} string "Node not found"
// @Router /nodes/{id} [put]
func (h *NodeHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var input UpdateNodeInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := []db.NodeSetParam{}

	if input.Name != nil {
		params = append(params, db.Node.Name.Set(*input.Name))
	}
	if input.NodeData != nil {
		params = append(params, db.Node.NodeData.Set(*input.NodeData))
	}

	node, err := h.db.Node.FindUnique(
		db.Node.ID.Equals(id),
	).Update(
		params...,
	).Exec(r.Context())

	if err != nil {
		http.Error(w, "Node not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(node)
}

// Delete godoc
// @Summary Delete a node
// @Description Delete a node by ID
// @Tags nodes
// @Accept json
// @Produce json
// @Param id path string true "Node ID"
// @Success 204 "No Content"
// @Failure 404 {string} string "Node not found"
// @Router /nodes/{id} [delete]
func (h *NodeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, err := h.db.Node.FindUnique(
		db.Node.ID.Equals(id),
	).Delete().Exec(r.Context())

	if err != nil {
		http.Error(w, "Node not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
