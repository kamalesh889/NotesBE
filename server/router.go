package server

import (
	"NOTESBE/utility"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gorilla/mux"
	"github.com/juju/ratelimit"
)

func Router(s *server) *mux.Router {

	r := s.router

	limiter := tollbooth.NewLimiter(100, &limiter.ExpirableOptions{})
	r.Use(utility.RateLimitMiddleware(limiter))

	throttleLimiter := ratelimit.NewBucketWithQuantum(time.Second, 100, 100)
	r.Use(utility.RequestThrottleMiddleware(throttleLimiter))

	r.HandleFunc("/ping", Ping).Methods("GET")

	// Authentication routes
	authRouter := r.PathPrefix("/api/auth").Subrouter()
	authRouter.HandleFunc("/signup", s.Signup).Methods("POST")
	authRouter.HandleFunc("/login", s.Login).Methods("POST")

	// Notes routes
	notesRouter := r.PathPrefix("/api/notes").Subrouter()
	notesRouter.HandleFunc("", utility.VerifyToken(s.CreateNotes)).Methods("POST")
	notesRouter.HandleFunc("", utility.VerifyToken(s.GetNotes)).Methods("GET")
	notesRouter.HandleFunc("/{id}", utility.VerifyToken(s.GetNotesById)).Methods("GET")
	notesRouter.HandleFunc("/{id}", utility.VerifyToken(s.UpdateNoteById)).Methods("PUT")
	notesRouter.HandleFunc("/{id}", utility.VerifyToken(s.DeleteNoteById)).Methods("DELETE")
	notesRouter.HandleFunc("/{id}/share", utility.VerifyToken(s.ShareNoteById)).Methods("POST")

	r.HandleFunc("/api/search", utility.VerifyToken(s.GetNoteByKey)).Methods("GET")

	return r
}
