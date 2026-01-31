package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"math"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	todo, nil := h.svc.CreateTODO(ctx, req.Subject, req.Description)

	return &model.CreateTODOResponse{
		TODO: *todo,
	}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	_, _ = h.svc.ReadTODO(ctx, 0, 0)
	return &model.ReadTODOResponse{}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	_, _ = h.svc.UpdateTODO(ctx, 0, "", "")
	return &model.UpdateTODOResponse{}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		var req model.CreateTODORequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		if req.Subject == "" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		todo, _ := h.svc.CreateTODO(r.Context(), req.Subject, req.Description)

		w.Header().Set("Content-type", "application/json")

		resp := model.CreateTODOResponse{
			TODO: *todo,
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Println(err)
		}

	case http.MethodPut:
		var req model.UpdateTODORequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		if req.ID == 0 || req.Subject == "" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		todo, err := h.svc.UpdateTODO(r.Context(), req.ID, req.Subject, req.Description)

		if err != nil {
			if _, ok := err.(*model.ErrNotFound); ok {
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			}
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-type", "application/json")

		resp := model.UpdateTODOResponse{
			TODO: *todo,
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Println(err)
		}

	case http.MethodGet:
		var prevID int64 = 0
		var size int64 = math.MaxInt64
		var err error

		q := r.URL.Query()
		prevIDStr := q.Get("prev_id")
		sizeStr := q.Get("size")

		if prevIDStr != "" {
			prevID, err = strconv.ParseInt(prevIDStr, 10, 64)
			if err != nil {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}
		}

		if sizeStr != "" {
			size, err = strconv.ParseInt(sizeStr, 10, 64)
			if err != nil {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}
		}

		todos, err := h.svc.ReadTODO(r.Context(), prevID,size)

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		resp := model.ReadTODOResponse{
			TODOs: todos,
		}

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusOK)

		_ = json.NewEncoder(w).Encode(resp)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

}
