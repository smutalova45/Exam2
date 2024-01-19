package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"main.go/api/models"
)

func (h Handler) Car(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateCar(w, r)
	case http.MethodGet:
		values := r.URL.Query()
		if _, ok := values["id"]; !ok {
			h.GetCarList(w, r)
		} else {
			h.GetCarByID(w, r)
		}
	case http.MethodPut:
		values := r.URL.Query()
		if _, ok := values["route"]; ok {
			h.UpdateCarRoute(w, r)
		} else if _, ok := values["status"]; ok {
			h.UpdateCarStatus(w, r)
		} else {
			h.UpdateCar(w, r)
		}
	case http.MethodDelete:
		h.DeleteCar(w, r)
	case http.MethodPatch:
		values := r.URL.Query()
		if _, ok := values["status"]; ok {
			h.UpdateCarStatus(w, r)
		} else {
			h.UpdateCarRoute(w, r)
		}
	}
}

func (h Handler) CreateCar(w http.ResponseWriter, r *http.Request) {
	createcar := models.CreateCar{}
	if err := json.NewDecoder(r.Body).Decode(&createcar); err != nil {
		fmt.Println(err.Error())
		handleResponse(w, http.StatusBadRequest, err)
		return
	}
	pk, err := h.storage.Car().Create(createcar)
	if err != nil {
		handleResponse(w, 500, err)
		return
	}
	car, err := h.storage.Car().Get(models.PrimaryKey{ID: pk})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}
	handleResponse(w, http.StatusCreated, car)
}

func (h Handler) GetCarByID(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	v := values.Get("id")
	if len(v) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}
	car, err := h.storage.Car().Get(models.PrimaryKey{ID: v})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, car)
}

func (h Handler) GetCarList(w http.ResponseWriter, r *http.Request) {
	var (
		page, limit = 1, 10
		search      string
		err         error
	)
	values := r.URL.Query()
	v := values.Get("page")
	if len(v) > 0 {
		page, err = strconv.Atoi(v)
		if err != nil {
			page = 1
		}
	}
	v1 := values.Get("limit")
	if len(v1) > 0 {
		limit, err = strconv.Atoi(v1)
		if err != nil {
			fmt.Println("limit", v1)
			limit=10
		}
	}
	v2 := values.Get("search")
	if len(v2) > 0 {
		search = v2
	}

	resp, err := h.storage.Car().GetList(models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})
	if err != nil {
		handleResponse(w, 500, err)
		return
	}
	handleResponse(w, 200, resp)

}

func (h Handler) UpdateCar(w http.ResponseWriter, r *http.Request) {
	updatecar := models.Car{}
	if err := json.NewDecoder(r.Body).Decode(&updatecar); err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	pk, err := h.storage.Car().Update(updatecar)
	if err != nil {
		handleResponse(w, 500, err)
		return
	}
	car, err := h.storage.Car().Get(models.PrimaryKey{
		ID: pk,
	})
	if err != nil {
		handleResponse(w, 500, err)
		return
	}

	handleResponse(w, 200, car)
}

func (h Handler) DeleteCar(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	id := values.Get("id")
	if len(id) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}
	if err := h.storage.Car().Delete(models.PrimaryKey{
		ID: id,
	}); err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	handleResponse(w, 200, "deleted data")
}

func (h Handler) UpdateCarRoute(w http.ResponseWriter, r *http.Request) {
	updatecar := models.UpdateCarRoute{}
	if err := json.NewDecoder(r.Body).Decode(&updatecar); err != nil {
		return
	}
	err := h.storage.Car().UpdateCarRoute(updatecar)
	if err != nil {
		fmt.Println("error in route", err.Error())
		return
	}
	handleResponse(w, 200, "updated route")
}

func (h Handler) UpdateCarStatus(w http.ResponseWriter, r *http.Request) {
	updatestatus := models.UpdateCarStatus{}
	if err := json.NewDecoder(r.Body).Decode(&updatestatus); err != nil {
		handleResponse(w, http.StatusBadRequest, err)
		return
	}
	err := h.storage.Car().UpdateCarStatus(updatestatus)
	if err != nil {
		fmt.Println("error in status", err.Error())
		return
	}
	handleResponse(w, 200, "updated status")
}
