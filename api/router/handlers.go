package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	db "api/dbmodel"
	"api/util/stringmaker"
	check "api/util/validators"

	"github.com/gorilla/mux"
)

const API_VERSION = "0.5"

type ShortifyAPI struct {
	connection *db.Connection
	version    string
}

type URLMap struct {
	LongURL string `json:longurl`
}

type APIResponse struct {
	StatusMessage string `json:statusmessage`
	ShortURL      string `json:shorturl`
}

func NewShortifyAPI() *ShortifyAPI {
	s := &ShortifyAPI{
		connection: db.NewConnection(),
		version:    API_VERSION,
	}
	return s
}

func (s *ShortifyAPI) Root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Shortify API \n"+
		"Version: %v \n"+
		"GET request with the short link to get the long link redirect \n"+
		"POST request with long link to get a short link back! \n", s.version)
}

func (s *ShortifyAPI) Create(w http.ResponseWriter, r *http.Request) {
	reqBody := new(URLMap)
	responseEncoder := json.NewEncoder(w)
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := responseEncoder.Encode(&APIResponse{StatusMessage: err.Error()}); err != nil {
			fmt.Fprintf(w, "Error occured while processing post request %v \n", err.Error())
		}
		return
	}
	genShortURL := stringmaker.NewString(8)
	// Validate the passed in longUrl
	isUrl, err := check.IsUrl(reqBody.LongURL)
	if !isUrl {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, err.Error(), reqBody.LongURL)
		return
	}
	err = s.connection.AddURL(reqBody.LongURL, genShortURL)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		if err := responseEncoder.Encode(&APIResponse{StatusMessage: err.Error()}); err != nil {
			fmt.Fprintf(w, "Error %s occured while trying to add the url \n", err.Error())
		}
		return
	}
	responseEncoder.Encode(&APIResponse{StatusMessage: "OK", ShortURL: genShortURL})
}

func (s *ShortifyAPI) Show(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortURL := vars["shorturl"]
	if len(shortURL) > 0 {
		// Validate the passed in shortUrl is 8 chars
		isShort8, err := check.IsShort8(shortURL)
		if !isShort8 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, err.Error(), shortURL)
			return
		}
		// Validate the passed in shortURL has valid chars
		isShortValidChars, err := check.IsShortValidChars(shortURL)
		if !isShortValidChars {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, err.Error(), shortURL)
			return
		}
		longURL, err := s.connection.FindLongURL(shortURL)
		if err != nil {
			fmt.Fprintf(w, "Could not find a long url that corresponds to this short url %s \n", shortURL)
			return
		}
		http.Redirect(w, r, longURL, http.StatusFound)
	}
}
