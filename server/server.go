package server

import (
	"NOTESBE/repository"

	"github.com/gorilla/mux"
)

type server struct {
	router *mux.Router
	db     repository.Repository
}

func NewServer(db *repository.Database) *server {

	s := &server{}

	s.router = mux.NewRouter()
	s.db = db

	return s
}
