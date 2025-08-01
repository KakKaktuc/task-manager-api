package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/KakKaktuc/task-manager-api/internal/repository"
	"github.com/KakKaktuc/task-manager-api/pkg/models"
)

type UserHandler struct {
	Repo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{Repo: repo}
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/users")
	path = strings.Trim(path, "/")

	switch r.Method {
	case http.MethodGet:
		if path == "" {
			h.getAll(w)
		} else {
			h.getByID(w, path)
		}
	case http.MethodPost:
		h.create(w, r)
	case http.MethodPut:
		h.update(w, r, path)
	case http.MethodDelete:
		h.delete(w, path)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *UserHandler) getAll(w http.ResponseWriter) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	done := make(chan struct{})

	go func() {
		users := h.Repo.GetAll()
		json.NewEncoder(w).Encode(users)
		close(done)
	}()

	select {
	case <-ctx.Done():
		http.Error(w, "Request timed out", http.StatusGatewayTimeout)
	case <-done:
		// ok
	}
}

func (h *UserHandler) getByID(w http.ResponseWriter, idStr string) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}

	done := make(chan struct{})

	go func() {
		user, err := h.Repo.GetByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			json.NewEncoder(w).Encode(user)
		}
		close(done)
	}()

	select {
	case <-ctx.Done():
		http.Error(w, "Request timed out", http.StatusGatewayTimeout)
	case <-done:
	}
}

func (h *UserHandler) create(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	created := h.Repo.Create(user)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *UserHandler) update(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	updated, err := h.Repo.Update(id, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(updated)
}

func (h *UserHandler) delete(w http.ResponseWriter, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.Repo.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
