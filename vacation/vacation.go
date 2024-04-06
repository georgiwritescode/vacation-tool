package vacation

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Vacation struct {
	ID       int    `json:"id"`
	Label    string `json:"label"`
	FromDate string `json:"fromDate"`
	ToDate   string `json:"toDate"`
	PersonId int    `json:"personId"`
}

type Handler struct{}

var vacations []Vacation = []Vacation{
	{ID: 1, Label: "Vacation 1", FromDate: "2024-04-06", ToDate: "2024-04-10", PersonId: 1},
	{ID: 2, Label: "Vacation 2", FromDate: "2024-05-15", ToDate: "2024-05-20", PersonId: 2},
	{ID: 3, Label: "Vacation 3", FromDate: "2024-06-01", ToDate: "2024-06-05", PersonId: 3},
	{ID: 4, Label: "Vacation 4", FromDate: "2024-07-10", ToDate: "2024-07-15", PersonId: 4},
	{ID: 5, Label: "Vacation 5", FromDate: "2024-08-20", ToDate: "2024-08-25", PersonId: 5},
	{ID: 6, Label: "Vacation 6", FromDate: "2024-09-06", ToDate: "2024-09-10", PersonId: 1},
	{ID: 7, Label: "Vacation 7", FromDate: "2024-10-15", ToDate: "2024-10-20", PersonId: 2},
	{ID: 8, Label: "Vacation 8", FromDate: "2024-11-01", ToDate: "2024-11-05", PersonId: 3},
	{ID: 9, Label: "Vacation 9", FromDate: "2024-12-10", ToDate: "2024-12-15", PersonId: 4},
	{ID: 10, Label: "Vacation 10", FromDate: "2025-01-20", ToDate: "2025-01-25", PersonId: 5},
	{ID: 11, Label: "Vacation 11", FromDate: "2025-02-06", ToDate: "2025-02-10", PersonId: 1},
	{ID: 12, Label: "Vacation 12", FromDate: "2025-03-15", ToDate: "2025-03-20", PersonId: 2},
	{ID: 13, Label: "Vacation 13", FromDate: "2025-04-01", ToDate: "2025-04-05", PersonId: 3},
	{ID: 14, Label: "Vacation 14", FromDate: "2025-05-10", ToDate: "2025-05-15", PersonId: 4},
	{ID: 15, Label: "Vacation 15", FromDate: "2025-06-20", ToDate: "2025-06-25", PersonId: 5},
	{ID: 16, Label: "Vacation 16", FromDate: "2025-07-06", ToDate: "2025-07-10", PersonId: 1},
	{ID: 17, Label: "Vacation 17", FromDate: "2025-08-15", ToDate: "2025-08-20", PersonId: 2},
	{ID: 18, Label: "Vacation 18", FromDate: "2025-09-01", ToDate: "2025-09-05", PersonId: 3},
	{ID: 19, Label: "Vacation 19", FromDate: "2025-10-10", ToDate: "2025-10-15", PersonId: 4},
	{ID: 20, Label: "Vacation 20", FromDate: "2025-11-20", ToDate: "2025-11-25", PersonId: 5},
}

func (h *Handler) FindById(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	vacation, exists := loadData(id)
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(vacation)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var vacation Vacation
	json.NewDecoder(r.Body).Decode(&vacation)
	vacations = append(vacations, vacation)
}

func (h *Handler) FindAll(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(vacations)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(r.PathValue("id"))

	var req Vacation
	json.NewDecoder(r.Body).Decode(&req)

	for i, vacation := range vacations {
		if vacation.ID == id {
			vacations[i] = req
		}
	}

	json.NewEncoder(w).Encode(req)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(r.PathValue("id"))

	indexToRemove := -1

	for i, vacation := range vacations {
		if vacation.ID == id {
			indexToRemove = i
			break
		}
	}

	if indexToRemove != -1 {
		vacations = append(vacations[:indexToRemove], vacations[indexToRemove+1:]...)
	}

}

func loadData(x string) (Vacation, bool) {

	id, _ := strconv.Atoi(x)

	res := Vacation{}

	for _, x := range vacations {
		if x.ID == id {
			return x, true
		}
	}

	return res, false
}
