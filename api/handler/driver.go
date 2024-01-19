package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"main.go/api/models"
)

func (h Handler) Driver(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateDriver(w, r)
	case http.MethodGet:
		values := r.URL.Query()
		_, ok := values["id"]
		if !ok {
			h.GetDriverList(w, r)
		} else {
			h.GetDriverByID(w, r)
		}
	case http.MethodPut:
		h.UpdateDriver(w, r)
	case http.MethodDelete:
		h.DeleteDriver(w, r)
	}
}

func (h Handler) CreateDriver(w http.ResponseWriter, r *http.Request) {
	createdriver := models.CreateDriver{}
	if err := json.NewDecoder(r.Body).Decode(&createdriver); err != nil {
		fmt.Println(err.Error())
		handleResponse(w, 400, err)
		return
	}
	pk, err := h.storage.Driver().Create(createdriver)
	if err != nil {
		fmt.Println(err.Error())
		handleResponse(w, 500, err)
		return
	}
	driver, err := h.storage.Driver().Get(models.PrimaryKey{ID: pk})
	if err != nil {
		handleResponse(w, 500, err)
		return
	}
	handleResponse(w, 200, driver)
}

func (h Handler) GetDriverByID(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	v := values.Get("id")
	if len(v) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}
	driver, err := h.storage.Driver().Get(models.PrimaryKey{ID: v})
	if err != nil {
		handleResponse(w, 500, err)
		return
	}
	handleResponse(w, 200, driver)
}

func (h Handler) GetDriverList(w http.ResponseWriter, r *http.Request) {
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
	res, err := h.storage.Driver().GetList(models.GetListRequest{
		Page:   page,
		Limit:  limit,
		Search: search,
	})
	handleResponse(w, 200, res)
}

func (h Handler) UpdateDriver(w http.ResponseWriter, r *http.Request) {
	updatedriver := models.UpdateDriver{}
	if err := json.NewDecoder(r.Body).Decode(&updatedriver); err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	pk, err := h.storage.Driver().Update(updatedriver)
	if err != nil {
		handleResponse(w, 500, err)
		return
	}
	driver, err := h.storage.Driver().Get(models.PrimaryKey{
		ID: pk,
	})
	if err != nil {
		handleResponse(w, 500, err)
		return
	}
	handleResponse(w, 200, driver)
}

func (h Handler) DeleteDriver(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	id := values.Get("id")
	if len(id) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}
	if err := h.storage.Driver().Delete(models.PrimaryKey{
		ID: id,
	}); err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	handleResponse(w, 200, "deleted data")
}
