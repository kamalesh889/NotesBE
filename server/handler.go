package server

import (
	"NotesBe/repository"
	"NotesBe/utility"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *server) Signup(w http.ResponseWriter, r *http.Request) {

	var req UserReq

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.UserName == "" || req.PassWord == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userInfo := &repository.User{
		Username: req.UserName,
		Password: req.PassWord,
	}

	err = s.db.CreateUser(userInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *server) Login(w http.ResponseWriter, r *http.Request) {

	var req UserReq

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.UserName == "" || req.PassWord == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userInfo := &repository.User{
		Username: req.UserName,
		Password: req.PassWord,
	}

	userId, err := s.db.GetUser(userInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tokenReq := &utility.TokenReq{
		Id: userId,
	}

	token, err := tokenReq.CreateJwtToken()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := &LoginResp{
		Token:  token.Token,
		UserId: userId,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (s *server) CreateNotes(w http.ResponseWriter, r *http.Request) {

	var req NoteReq

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	query := r.URL.Query()
	user := query.Get("userid")

	userId, _ := strconv.ParseUint(user, 10, 64)

	if userId == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	noteInfo := &repository.Note{
		Note:   req.Note,
		Userid: userId,
	}

	err = s.db.CreateNote(noteInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *server) GetNotes(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	user := query.Get("userid")

	userId, _ := strconv.ParseUint(user, 10, 64)

	if userId == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	notes, err := s.db.GetNotesOfUser(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(notes)

}

func (s *server) GetNotesById(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	note, found := vars["id"]
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	noteId, _ := strconv.ParseUint(note, 10, 64)

	if noteId == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	notes, err := s.db.GetNoteById(noteId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(notes)

}

func (s *server) UpdateNoteById(w http.ResponseWriter, r *http.Request) {

	var req NoteReq

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)

	note, found := vars["id"]
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	noteId, _ := strconv.ParseUint(note, 10, 64)

	if noteId == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.db.UpdateNoteById(noteId, req.Note)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *server) DeleteNoteById(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	note, found := vars["id"]
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	noteId, _ := strconv.ParseUint(note, 10, 64)

	if noteId == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := s.db.DeleteNoteById(noteId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *server) ShareNoteById(w http.ResponseWriter, r *http.Request) {

	var req ShareNoteReq

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)

	note, found := vars["id"]
	if !found {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	noteId, _ := strconv.ParseUint(note, 10, 64)

	if noteId == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.RecieverId == 0 || req.SenderId == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.db.ShareNoteToUser(noteId, req.SenderId, req.RecieverId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *server) GetNoteByKey(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	user := query.Get("userid")
	keyword := query.Get("query")

	userId, _ := strconv.ParseUint(user, 10, 64)

	if userId == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	records, err := s.db.GetNotesByKey(userId, keyword)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(records)

}
