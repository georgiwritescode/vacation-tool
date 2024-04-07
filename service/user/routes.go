package user

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/georgiwritescode/vacation-tool/types"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc(http.MethodGet+" /users/{id}", h.HandleGetByID)
}

func (h *Handler) HandleGetByID(w http.ResponseWriter, r *http.Request) {

	pathValue := r.PathValue("id")
	id, err := strconv.Atoi(pathValue)
	if err != nil {
		log.Fatal(err)
	}

	user, err := h.store.FindById(id)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(user)
}
