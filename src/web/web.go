package web

import (
	"api"
	"utils"
	"net/http"
)

func InitWeb() {
	examplesPath := utils.FindDir(utils.Cfg.WebSettings.ExamplesPath)

	api.Srv.Router.PathPrefix("/examples/").Handler(http.StripPrefix("/examples/", http.FileServer(http.Dir(examplesPath))))
}
