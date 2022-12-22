package web

import (
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter() (*mux.Router, error) {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/auth/register", socketCtl.Handle).Methods(http.MethodPost)
}
