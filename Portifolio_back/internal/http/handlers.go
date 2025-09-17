package http

import (
	"encoding/json"
	"net/http"
)

func (s *Server) listTodos(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, s.store.List())
}

func (s *Server) createTodo(w http.ResponseWriter, r *http.Request) {
	var in todo.todo
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	if in.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[strin]string{"error": "title is required"})
		return
	}
}
