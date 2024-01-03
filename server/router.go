package server

import "github.com/gorilla/mux"

func Router(s *server) *mux.Router {

	r := s.router

	// Authentication routes
	authRouter := r.PathPrefix("/api/auth").Subrouter()
	authRouter.HandleFunc("/signup", s.Signup).Methods("POST")
	authRouter.HandleFunc("/login", s.Login).Methods("POST")

	// Notes routes
	notesRouter := r.PathPrefix("/api/notes").Subrouter()
	notesRouter.HandleFunc("", s.CreateNotes).Methods("POST")
	notesRouter.HandleFunc("", s.GetNotes).Methods("GET")
	notesRouter.HandleFunc("/{id}", s.GetNotesById).Methods("GET")
	notesRouter.HandleFunc("/{id}", s.UpdateNoteById).Methods("PUT")
	notesRouter.HandleFunc("/{id}", s.DeleteNoteById).Methods("DELETE")
	notesRouter.HandleFunc("/{id}/share", s.ShareNoteById).Methods("POST")

	r.HandleFunc("/api/search", s.GetNoteByKey).Methods("GET")

	return r
}
