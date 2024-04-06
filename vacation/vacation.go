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
}

type Handler struct{}

var vacations []Vacation = []Vacation{
	{ID: 1, Label: "Vacation 1", FromDate: "2024-04-06", ToDate: "2024-04-10"},
	{ID: 2, Label: "Vacation 2", FromDate: "2024-05-15", ToDate: "2024-05-20"},
	{ID: 3, Label: "Vacation 3", FromDate: "2024-06-01", ToDate: "2024-06-05"},
	{ID: 4, Label: "Vacation 4", FromDate: "2024-07-10", ToDate: "2024-07-15"},
	{ID: 5, Label: "Vacation 5", FromDate: "2024-08-20", ToDate: "2024-08-25"},
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
