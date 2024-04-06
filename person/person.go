package person

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Handler struct{}

var persons []Person = []Person{
	{ID: 1, FirstName: "John", LastName: "Doe", Age: 30, Email: "john.doe@example.com"},
	{ID: 2, FirstName: "Alice", LastName: "Smith", Age: 25, Email: "alice.smith@example.com"},
	{ID: 3, FirstName: "Bob", LastName: "Johnson", Age: 40, Email: "bob.johnson@example.com"},
	{ID: 4, FirstName: "Emily", LastName: "Brown", Age: 35, Email: "emily.brown@example.com"},
	{ID: 5, FirstName: "Michael", LastName: "Davis", Age: 28, Email: "michael.davis@example.com"},
}

func (h *Handler) FindById(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	person, exists := loadData(id)
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(person)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var person Person
	json.NewDecoder(r.Body).Decode(&person)
	persons = append(persons, person)
}

func (h *Handler) FindAll(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(persons)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(r.PathValue("id"))

	var req Person
	json.NewDecoder(r.Body).Decode(&req)

	for i, person := range persons {
		if person.ID == id {
			persons[i] = req
		}
	}

	json.NewEncoder(w).Encode(req)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(r.PathValue("id"))

	indexToRemove := -1

	for i, person := range persons {
		if person.ID == id {
			indexToRemove = i
			break
		}
	}

	if indexToRemove != -1 {
		persons = append(persons[:indexToRemove], persons[indexToRemove+1:]...)
	}

}

func loadData(x string) (Person, bool) {

	id, _ := strconv.Atoi(x)

	res := Person{}

	for _, x := range persons {
		if x.ID == id {
			return x, true
		}
	}

	return res, false
}
