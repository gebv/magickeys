package api

func InitApi() {
	r := Srv.Router.PathPrefix("/api/v1").Subrouter()

	// Components
	InitValues(r)
}