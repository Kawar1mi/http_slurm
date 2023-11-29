package api

import (
	"slurm/go-on-practice-2/http_06/internals/app/handlers"

	"github.com/gorilla/mux"
)

func CreateRoutes(usersHandler *handlers.UsersHandler, carsHandler *handlers.CarsHandler) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/users/create", usersHandler.Create).Methods("POST")
	r.HandleFunc("/users/list", usersHandler.List).Methods("GET")
	r.HandleFunc("/users/find/{id:[0-9]+}", usersHandler.Find).Methods("GET")

	r.HandleFunc("/cars/create", carsHandler.Create).Methods("POST")
	r.HandleFunc("/cars/list", carsHandler.List).Methods("GET")
	r.HandleFunc("/cars/find/{id:[0-9]+}", carsHandler.Find).Methods("GET")

	r.NotFoundHandler = r.NewRoute().HandlerFunc(handlers.NotFound).GetHandler()

	return r
}
