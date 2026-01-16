package user

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/georgiwritescode/vacation-tool/types"
	"github.com/georgiwritescode/vacation-tool/utils"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /users/{id}", h.HandleGetByID)
	router.HandleFunc("GET /users/list", h.HandleListAllUsers)
	router.HandleFunc("POST /users/create", h.HandleCreateUser)
	router.HandleFunc("PUT /users/update", h.HandleUpdateUser)
	router.HandleFunc("DELETE /users/delete/{id}", h.HandleDeleteUser)
}

func (h *Handler) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	var user types.User
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.store.UpdateUser(&user); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (h *Handler) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	pathValue := r.PathValue("id")
	id, err := strconv.Atoi(pathValue)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user id: %v", err))
		return
	}

	if err := h.store.DeleteUser(id); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *Handler) HandleGetByID(w http.ResponseWriter, r *http.Request) {

	pathValue := r.PathValue("id")
	id, err := strconv.Atoi(pathValue)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user id: %v", err))
		return
	}

	user, err := h.store.FindById(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, user)
}

func (h *Handler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {

	var user types.User

	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	res, err := h.store.CreateUser(&types.User{
		ID:           user.ID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Age:          user.Age,
		Email:        user.Email,
		VacationDays: user.VacationDays,
		NonPaidLeave: user.NonPaidLeave,
		Timestamp:    user.Timestamp,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, fmt.Sprintf("user with id: %d created", res))
}

func (h *Handler) HandleListAllUsers(w http.ResponseWriter, r *http.Request) {

	res, err := h.store.FetchAllUsers()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}

	utils.WriteJSON(w, http.StatusOK, res)
}
