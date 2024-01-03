package server

import (
	"NotesBe/repository"

	"github.com/gorilla/mux"
)

type server struct {
	router *mux.Router
	db     *repository.Database
}

func NewServer(db *repository.Database) *server {

	s := &server{}

	s.router = mux.NewRouter()
	s.db = db

	return s
}
