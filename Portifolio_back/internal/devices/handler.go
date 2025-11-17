package devices

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// Handler holds dependencies for HTTP handlers.
type Handler struct {
	Repo *Repo
	// You can add logger or other deps here
}

// NewHandler builds a new Handler instance.
func NewHandler(db *sql.DB) *Handler {
	return &Handler{Repo: NewRepo(db)}
}

// Helper: extract user ID from request context
// TODO: Replace with real auth extraction used in your project (e.g., from middleware)
func getUserIDFromContext(r *http.Request) (int64, error) {
	// Example: if auth middleware sets r.Context() value "user_id" (int64)
	v := r.Context().Value("user_id")
	if v == nil {
		return 0, errors.New("unauthenticated")
	}
	switch t := v.(type) {
	case int64:
		return t, nil
	case int:
		return int64(t), nil
	case string:
		// attempt parse
		if id, err := strconv.ParseInt(t, 10, 64); err == nil {
			return id, nil
		}
	}
	return 0, errors.New("invalid user in context")
}

// registerRoutes demonstrates how to register handlers on a mux.Router or http.ServeMux
// Example usage (pseudo-code from main): http.HandleFunc("/api/devices", handler.ListDevices) etc.

// ListDevices handles GET /api/devices
func (h *Handler) ListDevices(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uid, err := getUserIDFromContext(r)
	if err != nil {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)
		return
	}
	devs, err := h.Repo.ListByUser(ctx, uid)
	if err != nil {
		http.Error(w, "failed to list devices: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(devs)
}

// CreateDeviceRequest represents the payload for creating a device.
type CreateDeviceRequest struct {
	Name     string `json:"name"`
	Room     string `json:"room,omitempty"`
	Metadata string `json:"metadata,omitempty"`
}

// CreateDevice handles POST /api/devices
func (h *Handler) CreateDevice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uid, err := getUserIDFromContext(r)
	if err != nil {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)
		return
	}
	var req CreateDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	device := &Device{
		UserID:   uid,
		Name:     req.Name,
		Room:     req.Room,
		Metadata: req.Metadata,
	}
	id, err := h.Repo.Create(ctx, device)
	if err != nil {
		http.Error(w, "failed to create device: "+err.Error(), http.StatusInternalServerError)
		return
	}
	device.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(device)
}

// DeleteDevice handles DELETE /api/devices/{id}
// This handler expects the router to pass the id as a URL param; here we parse from query 'id' if not using a router lib.
func (h *Handler) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uid, err := getUserIDFromContext(r)
	if err != nil {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)
		return
	}
	// Try to get id from URL path (if router sets), otherwise from query param 'id'
	var idStr string
	// If you use gorilla/mux or chi, replace this with mux.Vars(r)["id"]
	if rid := r.URL.Query().Get("id"); rid != "" {
		idStr = rid
	} else {
		// fallback: last path segment
		segments := splitPathSegments(r.URL.Path)
		if len(segments) > 0 {
			idStr = segments[len(segments)-1]
		}
	}
	if idStr == "" {
		http.Error(w, "missing device id", http.StatusBadRequest)
		return
	}
	deviceID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.Repo.Delete(ctx, uid, deviceID); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to delete: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// splitPathSegments is a tiny helper to get path segments (no external deps).
func splitPathSegments(path string) []string {
	// remove trailing slash
	if len(path) > 1 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	segments := []string{}
	start := 0
	for i := 0; i < len(path); i++ {
		if path[i] == '/' {
			if i-start > 0 {
				segments = append(segments, path[start:i])
			}
			start = i + 1
		}
	}
	if start < len(path) {
		segments = append(segments, path[start:])
	}
	return segments
}
