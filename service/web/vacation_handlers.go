package web

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/georgiwritescode/vacation-tool/types"
	"github.com/georgiwritescode/vacation-tool/utils"
)

type VacationFormData struct {
	Vacation *types.Vacation
	Users    []*types.User
}

// HandleVacationNew shows the create vacation form
func (h *Handler) HandleVacationNew(w http.ResponseWriter, r *http.Request) {
	users, err := h.userStore.FetchAllUsers()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	tmpl, err := parseTemplate("vacation_form.html")
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	tmpl.Execute(w, VacationFormData{Vacation: &types.Vacation{}, Users: users})
}

// HandleVacationCreate processes the create vacation form
func (h *Handler) HandleVacationCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	personId, _ := strconv.Atoi(r.FormValue("person_id"))
	daysUsed, _ := strconv.Atoi(r.FormValue("days_used"))

	vacation := &types.Vacation{
		Label:     r.FormValue("label"),
		FromDate:  r.FormValue("from_date"),
		ToDate:    r.FormValue("to_date"),
		PersonId:  personId,
		DaysUsed:  daysUsed,
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}

	id, err := h.vacationStore.CreateVacation(vacation)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/vacations/%d", id), http.StatusSeeOther)
}

// HandleVacationDetail shows vacation details
func (h *Handler) HandleVacationDetail(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid vacation id"))
		return
	}

	vacation, err := h.vacationStore.FindById(id)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}

	tmpl, err := parseTemplate("vacation_detail.html")
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	tmpl.Execute(w, vacation)
}

// HandleVacationEdit shows the edit vacation form
func (h *Handler) HandleVacationEdit(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid vacation id"))
		return
	}

	vacation, err := h.vacationStore.FindById(id)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}

	users, err := h.userStore.FetchAllUsers()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	tmpl, err := parseTemplate("vacation_form.html")
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	tmpl.Execute(w, VacationFormData{Vacation: vacation, Users: users})
}

// HandleVacationUpdate processes the edit vacation form
func (h *Handler) HandleVacationUpdate(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid vacation id"))
		return
	}

	if err := r.ParseForm(); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	personId, _ := strconv.Atoi(r.FormValue("person_id"))
	daysUsed, _ := strconv.Atoi(r.FormValue("days_used"))

	vacation := &types.Vacation{
		ID:        id,
		Label:     r.FormValue("label"),
		FromDate:  r.FormValue("from_date"),
		ToDate:    r.FormValue("to_date"),
		PersonId:  personId,
		DaysUsed:  daysUsed,
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}

	if err := h.vacationStore.UpdateVacation(vacation); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/vacations/%d", id), http.StatusSeeOther)
}

// HandleVacationDelete processes vacation deletion
func (h *Handler) HandleVacationDelete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid vacation id"))
		return
	}

	if err := h.vacationStore.DeleteVacation(id); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, "/vacations", http.StatusSeeOther)
}
