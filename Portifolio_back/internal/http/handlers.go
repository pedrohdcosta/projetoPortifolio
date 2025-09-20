package apihttp

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/pedrohdcosta/projetoPortifolio/Portifolio_back/internal/todo"
)

func (s *Server) listTodos(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, s.store.List())
}

func (s *Server) createTodo(w http.ResponseWriter, r *http.Request) {
	var in todo.Todo
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	if in.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "title is required"})
		return
	}
	out := s.store.Create(in.Title)
	writeJSON(w, http.StatusCreated, out)
}

func (s *Server) getTodo(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	t, err := s.store.Get(id)
	if errors.Is(err, todo.ErrNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{"erros": "not found"})
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func (s *Server) updateTodo(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	var in todo.Todo
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	t, err := s.store.Update(id, in.Title, in.Done)
	if errors.Is(err, todo.ErrNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func (s *Server) deleteTodo(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	if err := s.store.Delete(id); errors.Is(err, todo.ErrNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func parseID(v string) (int, error) { return strconv.Atoi(v) }
