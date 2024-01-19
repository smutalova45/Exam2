package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"main.go/api/models"
)

func (h Handler) City(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateCity(w, r)
	case http.MethodGet:
		values := r.URL.Query()
		if _, ok := values["id"]; !ok {
			h.GetCityList(w, r)
		} else {
			h.GetCityByID(w, r)
		}
	case http.MethodPut:
		h.UpdateCity(w, r)
	case http.MethodDelete:
		h.DeleteCity(w, r)
	}
}

func (h Handler) CreateCity(w http.ResponseWriter, r *http.Request) {
	createcity := models.CreateCity{}
	if err := json.NewDecoder(r.Body).Decode(&createcity); err != nil {
		handleResponse(w, http.StatusBadRequest, err)
		return
	}

	pk, err := h.storage.City().Create(createcity)
	if err != nil {
		handleResponse(w, 500, err)
		return
	}

	city, err := h.storage.City().Get(models.PrimaryKey{ID: pk})
	if err != nil {
		handleResponse(w, 500, err)
		return
	}
	
	handleResponse(w, 200, city)

}

func (h Handler) GetCityByID(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	id := values.Get("id")
	if len(id) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}
	var err error
	city, err := h.storage.City().Get(models.PrimaryKey{ID: id})
	if err != nil {
		handleResponse(w, 500, err)
		return
	}
	handleResponse(w, 200, city)
}

func (h Handler) GetCityList(w http.ResponseWriter, r *http.Request) {
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
		}
	}
	v2 := values.Get("search")
	if len(v2) > 0 {
		search = v2
	}
	resp, err := h.storage.City().GetList(models.GetListRequest{
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

func (h Handler) UpdateCity(w http.ResponseWriter, r *http.Request) {
	updatecity := models.UpdateCity{}
	if err := json.NewDecoder(r.Body).Decode(&updatecity); err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	pk, err := h.storage.City().Update(updatecity)
	if err != nil {
		handleResponse(w, 500, err)
		return
	}
	city, err := h.storage.City().Get(models.PrimaryKey{
		ID: pk,
	})
	if err != nil {
		handleResponse(w, 500, err)
		return
	}

	handleResponse(w, 200, city)
}

func (h Handler) DeleteCity(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	id := values.Get("id")
	if len(id) <= 0 {
		handleResponse(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}
	if err := h.storage.City().Delete(models.PrimaryKey{
		ID: id,
	}); err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	handleResponse(w, 200, "deleted data")
}
