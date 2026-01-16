package web

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/georgiwritescode/vacation-tool/types"
	"github.com/georgiwritescode/vacation-tool/utils"
)

type UserFormData struct {
	User *types.User
}

// HandleUserNew shows the create user form
func (h *Handler) HandleUserNew(w http.ResponseWriter, r *http.Request) {
	tmpl, err := parseTemplate("user_form.html")
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	tmpl.Execute(w, UserFormData{User: &types.User{}})
}

// HandleUserCreate processes the create user form
func (h *Handler) HandleUserCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	age, _ := strconv.Atoi(r.FormValue("age"))
	vacationDays, _ := strconv.Atoi(r.FormValue("vacation_days"))
	nonPaidLeave, _ := strconv.Atoi(r.FormValue("non_paid_leave"))

	user := &types.User{
		FirstName:    r.FormValue("first_name"),
		LastName:     r.FormValue("last_name"),
		Age:          age,
		Email:        r.FormValue("email"),
		VacationDays: vacationDays,
		NonPaidLeave: nonPaidLeave,
		Timestamp:    time.Now().Format("2006-01-02 15:04:05"),
	}

	id, err := h.userStore.CreateUser(user)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/users/%d", id), http.StatusSeeOther)
}

// HandleUserDetail shows user details
func (h *Handler) HandleUserDetail(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user id"))
		return
	}

	user, err := h.userStore.FindById(id)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}

	tmpl, err := parseTemplate("user_detail.html")
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	tmpl.Execute(w, user)
}

// HandleUserEdit shows the edit user form
func (h *Handler) HandleUserEdit(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user id"))
		return
	}

	user, err := h.userStore.FindById(id)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}

	tmpl, err := parseTemplate("user_form.html")
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	tmpl.Execute(w, UserFormData{User: user})
}

// HandleUserUpdate processes the edit user form
func (h *Handler) HandleUserUpdate(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user id"))
		return
	}

	if err := r.ParseForm(); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	age, _ := strconv.Atoi(r.FormValue("age"))
	vacationDays, _ := strconv.Atoi(r.FormValue("vacation_days"))
	nonPaidLeave, _ := strconv.Atoi(r.FormValue("non_paid_leave"))

	user := &types.User{
		ID:           id,
		FirstName:    r.FormValue("first_name"),
		LastName:     r.FormValue("last_name"),
		Age:          age,
		Email:        r.FormValue("email"),
		VacationDays: vacationDays,
		NonPaidLeave: nonPaidLeave,
	}

	if err := h.userStore.UpdateUser(user); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/users/%d", id), http.StatusSeeOther)
}

// HandleUserDelete processes user deletion
func (h *Handler) HandleUserDelete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user id"))
		return
	}

	if err := h.userStore.DeleteUser(id); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}
