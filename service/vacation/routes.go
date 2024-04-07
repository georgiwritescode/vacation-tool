package vacation

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/georgiwritescode/vacation-tool/types"
)

type Handler struct {
	store types.VacationStore
}

func NewHandler(store types.VacationStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /vacations/{id}", h.HandleGetByID)
}

func (h *Handler) HandleGetByID(w http.ResponseWriter, r *http.Request) {

	pathValue := r.PathValue("id")
	id, err := strconv.Atoi(pathValue)
	if err != nil {
		log.Fatal(err)
	}

	vacation, err := h.store.FindById(id)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(vacation)
}
