package api

import (
	"github.com/gorilla/mux"
	"models"
	"net/http"
	// "github.com/golang/glog"
	"strings"
)


func InitValues(r *mux.Router) {
	sr := r.PathPrefix("/values").Subrouter()

	sr.HandleFunc("/", NewValueHandler).Methods("POST")
	sr.HandleFunc("/{value_id}", GetValueHandler).Methods("GET")
	sr.HandleFunc("/{value_id}", UpdateValueHandler).Methods("PUT")
	sr.HandleFunc("/{value_id}", DeleteValueHandler).Methods("DELETE")

	sr.HandleFunc("/search/eq/{keys}", BuildSearchByMode("=")).Methods("GET")
	sr.HandleFunc("/search/contains/{keys}", BuildSearchByMode("@>")).Methods("GET")
	sr.HandleFunc("/search/any/{keys}", BuildSearchByMode("&&")).Methods("GET")
}

func BuildSearchByMode(mode string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		dto := models.NewValueDTO()
		res := models.NewResponse()
		dto.Keys.AddAsArray(strings.Split(mux.Vars(r)["keys"], ","))

		values, err := Srv.Store.ValueStore().FindByKeys(dto, mode)

		if err != nil {
			res.Message = "Ошибка поиска"
			res.DevMessage = err.Error()
			w.Write(res.ToJson())
			return	
		}

		res.StatusCode = http.StatusOK
		res.Data = values

		w.Write(res.ToJson())
	}
} 

func NewValueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	dto := models.NewValueDTO()
	res := models.NewResponse()

	if err := dto.FromJson(r.Body); err != nil {
		res.Message = "unknown"
		res.DevMessage = err.Error()
		w.Write(res.ToJson())
		return
	}
	value, err := Srv.Store.ValueStore().Create(dto)

	if err != nil {
		res.Message = "Ошибка создания"
		res.DevMessage = err.Error()
		w.Write(res.ToJson())
		return	
	}

	res.StatusCode = http.StatusOK
	res.Data = value

	w.Write(res.ToJson())
}

func GetValueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	dto := models.NewValueDTO()
	dto.ValueId = mux.Vars(r)["value_id"]
	res := models.NewResponse()

	value, err := Srv.Store.ValueStore().GetOne(dto)

	if err != nil {
		res.Message = "Ошибка загрузки"
		res.DevMessage = err.Error()
		w.Write(res.ToJson())
		return	
	}

	res.StatusCode = http.StatusOK
	res.Data = value

	w.Write(res.ToJson())
}

func UpdateValueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	dto := models.NewValueDTO()
	dto.ValueId = mux.Vars(r)["value_id"]
	res := models.NewResponse()

	if err := dto.FromJson(r.Body); err != nil {
		res.Message = "unknown"
		res.DevMessage = err.Error()
		w.Write(res.ToJson())
		return
	}
	value, err := Srv.Store.ValueStore().Update(dto)

	if err != nil {
		res.Message = "Ошибка обновления"
		res.DevMessage = err.Error()
		w.Write(res.ToJson())
		return	
	}

	res.StatusCode = http.StatusOK
	res.Data = value

	w.Write(res.ToJson())
}

func DeleteValueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	dto := models.NewValueDTO()
	dto.ValueId = mux.Vars(r)["value_id"]
	res := models.NewResponse()

	value, err := Srv.Store.ValueStore().Delete(dto)

	if err != nil {
		res.Message = "Ошибка удаления"
		res.DevMessage = err.Error()
		w.Write(res.ToJson())
		return	
	}

	res.StatusCode = http.StatusOK
	res.Data = value

	w.Write(res.ToJson())
}
