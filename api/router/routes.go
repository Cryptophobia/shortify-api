package router

import "net/http"

// Route struct for routing
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Slice of routes struct
type Routes []Route

// Function to initialize the routes and return them all
func InitRoutes(s *ShortifyAPI) Routes {
	return Routes{
		Route{
			"Root",
			"GET",
			"/",
			s.Root,
		},
		Route{
			"Show",
			"GET",
			"/{shorturl}",
			s.Show,
		},
		Route{
			"Create",
			"POST",
			"/Create",
			s.Create,
		},
	}
}
