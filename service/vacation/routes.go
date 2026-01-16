package vacation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/georgiwritescode/vacation-tool/types"
	"github.com/georgiwritescode/vacation-tool/utils"
)

type Handler struct {
	store types.VacationStore
}

func NewHandler(store types.VacationStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /api/v1/vacations/{id}", h.HandleGetByID)
	router.HandleFunc("GET /api/v1/vacations/list", h.HandleListVacations)
	router.HandleFunc("POST /api/v1/vacations/create", h.HandleCreateVacation)
	router.HandleFunc("PUT /api/v1/vacations/update", h.HandleUpdateVacation)
	router.HandleFunc("DELETE /api/v1/vacations/delete/{id}", h.HandleDeleteVacation)
}

func (h *Handler) HandleGetByID(w http.ResponseWriter, r *http.Request) {

	pathValue := r.PathValue("id")
	id, err := strconv.Atoi(pathValue)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid vacation id: %v", err))
		return
	}

	vacation, err := h.store.FindById(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, vacation)
}

func (h *Handler) HandleListVacations(w http.ResponseWriter, r *http.Request) {
	vacations, err := h.store.FindAll()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, vacations)
}

func (h *Handler) HandleCreateVacation(w http.ResponseWriter, r *http.Request) {
	var vacation types.Vacation
	if err := json.NewDecoder(r.Body).Decode(&vacation); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	id, err := h.store.CreateVacation(&vacation)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]int{"id": id})
}

func (h *Handler) HandleUpdateVacation(w http.ResponseWriter, r *http.Request) {
	var vacation types.Vacation
	if err := json.NewDecoder(r.Body).Decode(&vacation); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.store.UpdateVacation(&vacation); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (h *Handler) HandleDeleteVacation(w http.ResponseWriter, r *http.Request) {
	pathValue := r.PathValue("id")
	id, err := strconv.Atoi(pathValue)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid vacation id: %v", err))
		return
	}

	if err := h.store.DeleteVacation(id); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}
