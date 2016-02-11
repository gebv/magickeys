package api

import (
	"github.com/gorilla/mux"
	// "models"
	"net/http"
	"github.com/golang/glog"
)


func InitValues(r *mux.Router) {
	sr := r.PathPrefix("/values").Subrouter()

	sr.HandleFunc("/", NewValueHandler).Methods("POST")
	sr.HandleFunc("/", GetValueHandler).Methods("GET")
	// sr.HandleFunc("/", UpdateValueHandler).Methods("PUT")
	// sr.HandleFunc("/", DeleteValueHandler).Methods("DELETE")
}

func NewValueHandler(w http.ResponseWriter, r *http.Request) {
	glog.Infof("[POST] new value")
}

func GetValueHandler(w http.ResponseWriter, r *http.Request) {
	glog.Infof("[GET] %v", r.URL.Query().Get("value_id"))
}

func UpdateValueHandler(w http.ResponseWriter, r *http.Request) {
	glog.Infof("[GET] %v", r.URL.Query().Get("value_id"))
}

func DeleteValueHandler(w http.ResponseWriter, r *http.Request) {
	glog.Infof("[GET] %v", r.URL.Query().Get("value_id"))
}
