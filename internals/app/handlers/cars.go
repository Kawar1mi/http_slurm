package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"slurm/go-on-practice-2/http_06/internals/app/models"
	"slurm/go-on-practice-2/http_06/internals/app/processors"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type CarsHandler struct {
	processor *processors.CarsProcessor
}

func NewCarsHandler(processor *processors.CarsProcessor) *CarsHandler {
	handler := new(CarsHandler)
	handler.processor = processor
	return handler
}

func (handler *CarsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var newUser models.Car

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		WrapError(w, err)
		return
	}

	err = handler.processor.CreateCar(newUser)
	if err != nil {
		WrapError(w, err)
		return
	}

	var m = map[string]any{
		"result": "OK",
		"data":   "",
	}

	WrapOK(w, m)
}

func (handler *CarsHandler) List(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	var userIdFilter int64
	if vars.Get("userid") != "" {
		var err error
		userIdFilter, err = strconv.ParseInt(vars.Get("userid"), 10, 64)
		if err != nil {
			WrapError(w, err)
			return
		}
	}

	list, err := handler.processor.ListCars(
		userIdFilter,
		strings.Trim(vars.Get("brand"), "\""),
		strings.Trim(vars.Get("colour"), "\""),
		strings.Trim(vars.Get("license"), "\""),
	)
	if err != nil {
		WrapError(w, err)
		return
	}

	var m = map[string]any{
		"result": "OK",
		"data":   list,
	}

	WrapOK(w, m)

}

func (handler *CarsHandler) Find(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars["id"] == "" {
		WrapError(w, errors.New("missing id"))
		return
	}

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		WrapError(w, err)
		return
	}

	car, err := handler.processor.FindCar(id)
	if err != nil {
		WrapError(w, err)
		return
	}

	var m = map[string]any{
		"result": "OK",
		"data":   car,
	}

	WrapOK(w, m)
}
