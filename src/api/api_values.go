package api

import (
	"github.com/gorilla/mux"
	"models"
	"net/http"
	// "github.com/golang/glog"
	"strings"
	"time"
	"utils"
)

func changeHeader(handler func (http.ResponseWriter, *http.Request)) func (http.ResponseWriter, *http.Request) {
	return func (w http.ResponseWriter, r *http.Request) {
		// glog.Infof("[%v] %v", r.Method, r.URL.String())	

		time.Sleep(time.Millisecond * utils.Cfg.ServiceSettings.TimeoutRequest)

		w.Header().Add("Content-Type", "application/json")

		if origin := r.Header.Get("Origin"); origin != "" {
		    w.Header().Set("Access-Control-Allow-Origin", origin)
		    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		    w.Header().Set("Access-Control-Allow-Headers",
		      "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		  }

		  // Stop here for a Preflighted OPTIONS request.
		  if r.Method == "OPTIONS" {
		    return
		  }

		handler(w, r)
	}
}

func InitValues(r *mux.Router) {
	sr := r.PathPrefix("/values").Subrouter()	

	sr.HandleFunc("/", changeHeader(NewValueHandler)).Methods("POST", "OPTIONS")
	sr.HandleFunc("/{value_id}", changeHeader(GetValueHandler)).Methods("GET")
	sr.HandleFunc("/{value_id}", changeHeader(UpdateValueHandler)).Methods("PUT", "OPTIONS")
	sr.HandleFunc("/{value_id}", changeHeader(DeleteValueHandler)).Methods("DELETE", "OPTIONS")

	sr.HandleFunc("/search/eq/{keys}", changeHeader(BuildSearchByMode("="))).Methods("GET")
	sr.HandleFunc("/search/contains/{keys}", changeHeader(BuildSearchByMode("@>"))).Methods("GET")
	sr.HandleFunc("/search/any/{keys}", changeHeader(BuildSearchByMode("&&"))).Methods("GET")
}

func BuildSearchByMode(mode string) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
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
	dto := models.NewValueDTO()
	dto.ValueId = mux.Vars(r)["value_id"]
	dto.UpdateFields = strings.Split(r.URL.Query().Get("fields"), ",")

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
