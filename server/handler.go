package server

import (
	"NOTESBE/repository"
	"NOTESBE/utility"
	"encoding/json"
	"net/http"
	"strconv"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ping\n"))
}

func (s *server) Signup(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var req UserReq

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if req.UserName == "" || req.PassWord == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
		return
	}

	userInfo := &repository.User{
		Username: req.UserName,
		Password: req.PassWord,
	}

	err = s.db.CreateUser(userInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User Account created"))
}

func (s *server) Login(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var req UserReq

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if req.UserName == "" || req.PassWord == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
		return
	}

	userInfo := &repository.User{
		Username: req.UserName,
		Password: req.PassWord,
	}

	userId, err := s.db.GetUser(userInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	tokenReq := &utility.TokenReq{
		Id: userId,
	}

	token, err := tokenReq.CreateJwtToken()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
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

	w.Header().Set("Content-Type", "application/json")

	var req NoteReq

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	userId, err := utility.ParseUserId(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	noteInfo := &repository.Note{
		Note:   req.Note,
		Userid: userId,
	}

	err = s.db.CreateNote(noteInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *server) GetNotes(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	userId, err := utility.ParseUserId(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	notes, err := s.db.GetNotesOfUser(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(notes)

}

func (s *server) GetNotesById(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	noteId, err := utility.ParseNoteId(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	userId, err := utility.ParseUserId(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	notes, err := s.db.GetNoteById(noteId, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(notes)

}

func (s *server) UpdateNoteById(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var req NoteReq

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	userId, err := utility.ParseUserId(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	noteId, err := utility.ParseNoteId(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	err = s.db.UpdateNoteById(noteId, userId, req.Note)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *server) DeleteNoteById(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	noteId, err := utility.ParseNoteId(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	userId, err := utility.ParseUserId(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	err = s.db.DeleteNoteById(noteId, userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *server) ShareNoteById(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var req ShareNoteReq

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	noteId, err := utility.ParseNoteId(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if req.RecieverId == 0 || req.SenderId == 0 {
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.db.ShareNoteToUser(noteId, req.SenderId, req.RecieverId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *server) GetNoteByKey(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query()
	user := query.Get("userid")
	keyword := query.Get("query")

	userId, _ := strconv.ParseUint(user, 10, 64)

	if userId == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "User id is not valid"})
		return
	}

	records, err := s.db.GetNotesByKey(userId, keyword)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(records)

}
