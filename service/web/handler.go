package web

import (
	"html/template"
	"net/http"
	"path/filepath"

	"time"

	"github.com/georgiwritescode/vacation-tool/types"
	"github.com/georgiwritescode/vacation-tool/utils"
)

type DashboardData struct {
	ActiveVacations []ActiveVacationView
}

type ActiveVacationView struct {
	UserName string
	Label    string
	Dates    string
}

type Handler struct {
	userStore     types.UserStore
	vacationStore types.VacationStore
}

func NewHandler(userStore types.UserStore, vacationStore types.VacationStore) *Handler {
	return &Handler{
		userStore:     userStore,
		vacationStore: vacationStore,
	}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /", h.HandleIndex)
	router.HandleFunc("GET /users", h.HandleUsers)
	router.HandleFunc("GET /users/new", h.HandleUserNew)
	router.HandleFunc("POST /users/new", h.HandleUserCreate)
	router.HandleFunc("GET /users/{id}", h.HandleUserDetail)
	router.HandleFunc("GET /users/{id}/edit", h.HandleUserEdit)
	router.HandleFunc("POST /users/{id}/edit", h.HandleUserUpdate)
	router.HandleFunc("POST /users/{id}/delete", h.HandleUserDelete)
	router.HandleFunc("GET /vacations", h.HandleVacations)
	router.HandleFunc("GET /vacations/new", h.HandleVacationNew)
	router.HandleFunc("POST /vacations/new", h.HandleVacationCreate)
	router.HandleFunc("GET /vacations/{id}", h.HandleVacationDetail)
	router.HandleFunc("GET /vacations/{id}/edit", h.HandleVacationEdit)
	router.HandleFunc("POST /vacations/{id}/edit", h.HandleVacationUpdate)
	router.HandleFunc("POST /vacations/{id}/delete", h.HandleVacationDelete)
}

func (h *Handler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	today := time.Now().Format("2006-01-02")
	activeVacations, err := h.vacationStore.GetActiveVacations(today)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	var viewData DashboardData
	for _, v := range activeVacations {
		user, err := h.userStore.FindById(v.PersonId)
		userName := "Unknown"
		if err == nil {
			userName = user.FirstName + " " + user.LastName
		}

		viewData.ActiveVacations = append(viewData.ActiveVacations, ActiveVacationView{
			UserName: userName,
			Label:    v.Label,
			Dates:    v.FromDate + " to " + v.ToDate,
		})
	}

	tmpl, err := parseTemplate("index.html")
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	tmpl.Execute(w, viewData)
}

func (h *Handler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userStore.FetchAllUsers()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	tmpl, err := parseTemplate("users.html")
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	
	// Pass data to template
	tmpl.Execute(w, users)
}

func (h *Handler) HandleVacations(w http.ResponseWriter, r *http.Request) {
	vacations, err := h.vacationStore.FindAll()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	tmpl, err := parseTemplate("vacations.html")
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	tmpl.Execute(w, vacations)
}

func parseTemplate(name string) (*template.Template, error) {
	// Parse base layout and the specific page template
	return template.ParseFiles(
		filepath.Join("templates", "base.html"),
		filepath.Join("templates", name),
	)
}
