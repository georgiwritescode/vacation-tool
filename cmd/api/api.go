package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/georgiwritescode/vacation-tool/service/user"
	"github.com/georgiwritescode/vacation-tool/service/vacation"
)

type ApiServer struct {
	addr string
	db   *sql.DB
}

func NewApiServer(addr string, db *sql.DB) *ApiServer {
	return &ApiServer{
		addr: addr,
		db:   db,
	}
}

func (s *ApiServer) Run() error {
	router := http.NewServeMux()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(router)

	vacationStore := vacation.NewStore(s.db)
	vacationHandler := vacation.NewHandler(vacationStore)
	vacationHandler.RegisterRoutes(router)

	log.Println("Server listening on port", s.addr)
	return http.ListenAndServe(s.addr, router)
}
