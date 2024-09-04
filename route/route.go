package route

import (
	"net/http"
)

type Route struct {
	mux *http.ServeMux
}

func (route *Route) ResgisterRoutes() {
	route.mux.Handle("/", http.HandlerFunc(route.HomeHanlder))
}

func (route *Route) HomeHanlder(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Home"))
}
