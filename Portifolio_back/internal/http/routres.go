package http

import (
	"encoding/json"
	"net/http"

	"github.com/pedrohdcosta/projetoPortifolio/Portifolio_back/internal/todo"
)

type Server struct {
	store *todo.MemoryStore
	mux   *http.ServeMux
}

func NewRouter() *Server {
	store := todo.MemoryStore()
	mux := http.NewServeMux()

	s := &Server{store: store, mux: mux}

	mux.HandleFunc("GET /helth", s.health)
	mux.HandleFunc("Get /todos", s.listTodos)
	mux.HandleFunc("POST /todos", s.createTodo)
	mux.HandleFunc("GET /todos/{id}", s.getTodo)
	mux.HandleFunc("PUT /todos/{id}", s.updateTodo)
	mux.HandleFunc("DELETE /todos/{id}", s.deleteTodo)

	return s
}

func (s *Server) ServeHttp(w http.ResponseWriter, r *http.Request) { s.mux.ServeHTTP(w, r) }

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
