package router

import (
	"fmt"
	"net/http"

	pm "github.com/albertogviana/prometheus-middleware"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Create a new router and then return mux.Router
func NewShortifyRouter(routes Routes) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	// Add in the middleware for prometheus
	middleware := pm.NewPrometheusMiddleware(pm.Opts{})
	router.Handle("/metrics", promhttp.Handler())
	router.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})
	router.Use(middleware.InstrumentHandlerDuration)

	// Feed the router all of the information from the routes struct
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}
