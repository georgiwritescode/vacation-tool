package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/georgiwritescode/vacation-tool/db"
	"github.com/georgiwritescode/vacation-tool/person"
	"github.com/georgiwritescode/vacation-tool/vacation"
)

func main() {

	db.InitDB()

	mux := http.NewServeMux()
	personHandler := &person.Handler{}
	vacationHandler := &vacation.Handler{}

	mux.HandleFunc("GET /{$}", SendGreeting)

	//person router
	mux.HandleFunc("GET /person/{id}", personHandler.FindById)
	mux.HandleFunc("GET /person/list", personHandler.FindAll)
	mux.HandleFunc("POST /person/create", personHandler.Create)
	mux.HandleFunc("PUT /person/update/{id}", personHandler.Update)
	mux.HandleFunc("DELETE /person/delete/{id}", personHandler.Delete)

	//vacation router
	mux.HandleFunc("GET /vacation/{id}", vacationHandler.FindById)
	mux.HandleFunc("GET /vacation/list", vacationHandler.FindAll)
	mux.HandleFunc("POST /vacation/create", vacationHandler.Create)
	mux.HandleFunc("PUT /vacation/update/{id}", vacationHandler.Update)
	mux.HandleFunc("DELETE /vacation/delete/{id}", vacationHandler.Delete)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Server is listening on port", server.Addr)
	server.ListenAndServe()
}

func SendGreeting(w http.ResponseWriter, r *http.Request) {
	res := "Hello go v.1.22"
	json.NewEncoder(w).Encode(res)
}
