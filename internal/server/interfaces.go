package server

import "github.com/gorilla/mux"

type Server interface {
	NewRouter() *mux.Router
}
