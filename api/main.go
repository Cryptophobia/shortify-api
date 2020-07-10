package main

import (
	r "api/router"
	"net/http"
)

func main() {
	NewShortify := r.NewShortifyAPI()
	routes := r.InitRoutes(NewShortify)
	router := r.NewShortifyRouter(routes)
	http.ListenAndServe(":5000", router)
}
