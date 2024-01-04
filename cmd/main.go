package main

import (
	"NOTESBE/connection"
	"NOTESBE/server"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Notes Backend ")

	db, err := connection.InitializeDB()
	if err != nil {
		log.Panicln("Error in Connecting to Database:", err)
	}

	srv := server.NewServer(db)
	if err != nil {
		log.Panicln("Error in creating server:", err)
	}

	mux := server.Router(srv)
	http.Handle("/", mux)
	http.ListenAndServe(fmt.Sprintf(":%s", "8081"), mux)
}
