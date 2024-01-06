package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"NOTESBE/repository"
	repomock "NOTESBE/repository/mocks"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

var (
	testServer *server
	mockrepo   *repomock.MockRepository
)

func TestServer(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockrepo = repomock.NewMockRepository(ctrl)
	testServer = &server{}

	testServer.router = mux.NewRouter()
	testServer.db = mockrepo

}

func init() {

	t := &testing.T{}
	TestServer(t)
}

func TestSignup(t *testing.T) {

	mockreqbody1 := UserReq{
		UserName: "testuser",
		PassWord: "testpass",
	}
	mockbyt, _ := json.Marshal(mockreqbody1)

	mockreqbody2 := UserReq{}
	mockbyt2, _ := json.Marshal(mockreqbody2)

	t.Run("Success case", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewBuffer(mockbyt))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mockrepo.EXPECT().CreateUser(gomock.Any()).Return(nil)

		testServer.Signup(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code)

	})
	t.Run("Invalid request Body", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewBuffer(mockbyt2))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		testServer.Signup(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

	})

	t.Run("Error from Database", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewBuffer(mockbyt))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mockrepo.EXPECT().CreateUser(gomock.Any()).Return(errors.New("Database error"))

		testServer.Signup(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

	})
}

func TestLogin(t *testing.T) {

	mockreqbody1 := UserReq{
		UserName: "testuser",
		PassWord: "testpass",
	}
	mockbyt, _ := json.Marshal(mockreqbody1)

	mockreqbody2 := UserReq{}
	mockbyt2, _ := json.Marshal(mockreqbody2)

	t.Run("Success case", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(mockbyt))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mockrepo.EXPECT().GetUser(gomock.Any()).Return(uint64(1), nil)

		testServer.Login(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

	})
	t.Run("Invalid request Body", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(mockbyt2))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		testServer.Login(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

	})

	t.Run("Error from Database", func(t *testing.T) {

		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(mockbyt))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mockrepo.EXPECT().GetUser(gomock.Any()).Return(uint64(0), errors.New("Error from Database"))

		testServer.Login(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

	})
}

func TestCreateNotes(t *testing.T) {

	mockNoteReq := NoteReq{
		Note: "TestNote",
	}
	mockNoteReqBytes, _ := json.Marshal(mockNoteReq)

	mockUserID := uint64(1)

	t.Run("Success case", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodPost, "/api/notes?userid="+strconv.FormatUint(mockUserID, 10), bytes.NewBuffer(mockNoteReqBytes))
		if err != nil {
			t.Fatal("Error creating request:", err)
		}
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mockrepo.EXPECT().CreateNote(gomock.Any()).Return(nil)

		testServer.CreateNotes(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code)

	})

	t.Run("Failure case - Userid missing", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodPost, "/api/notes"+strconv.FormatUint(mockUserID, 10), bytes.NewBuffer(mockNoteReqBytes))
		if err != nil {
			t.Fatal("Error creating request:", err)
		}
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		testServer.CreateNotes(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

	})

	t.Run("Failure case - error from database", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodPost, "/api/notes?userid="+strconv.FormatUint(mockUserID, 10), bytes.NewBuffer(mockNoteReqBytes))
		if err != nil {
			t.Fatal("Error creating request:", err)
		}
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mockrepo.EXPECT().CreateNote(gomock.Any()).Return(errors.New("error from database"))

		testServer.CreateNotes(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

	})
}

func TestGetNotesById(t *testing.T) {

	mockNoteID := uint64(1)
	mockUserID := uint64(2)

	mockTime := time.Date(2024, time.January, 5, 18, 42, 48, 0, time.Local)

	t.Run("Success case", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/notes/{%d}?userid=%d", mockNoteID, mockUserID), nil)
		if err != nil {
			t.Fatal("Error creating request:", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatUint(mockNoteID, 10)})

		rec := httptest.NewRecorder()

		testNote := &repository.Note{
			Id: mockNoteID, Note: "Test Note 1", Userid: mockUserID, Createdat: mockTime, Updatedat: mockTime,
		}

		mockrepo.EXPECT().GetNoteById(mockNoteID, mockUserID).Return(testNote, nil)

		testServer.GetNotesById(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		var actualNotes repository.Note
		json.NewDecoder(rec.Body).Decode(&actualNotes)

		expectedNotes := *testNote
		assert.Equal(t, expectedNotes, actualNotes, "Unexpected response")
	})

	t.Run("Failure case - NoteId is missing", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/notes?userid=%d", mockUserID), nil)
		if err != nil {
			t.Fatal("Error creating request:", err)
		}

		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		testServer.GetNotesById(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

	})

	t.Run("Failure case - Error from Database", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/notes/{%d}?userid=%d", mockNoteID, mockUserID), nil)
		if err != nil {
			t.Fatal("Error creating request:", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatUint(mockNoteID, 10)})

		rec := httptest.NewRecorder()

		mockrepo.EXPECT().GetNoteById(mockNoteID, mockUserID).Return(nil, errors.New("Error from database"))

		testServer.GetNotesById(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestGetNotes(t *testing.T) {

	mockUserID := uint64(1)
	mockNoteID := uint64(5)
	mockTime := time.Date(2024, time.January, 5, 18, 42, 48, 0, time.Local)

	t.Run("Success case", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodGet, "/api/notes?userid="+strconv.FormatUint(mockUserID, 10), nil)
		if err != nil {
			t.Fatal("Error creating request:", err)
		}
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		testNotes := []repository.Note{{
			Id: mockNoteID, Note: "Test Note 1", Userid: mockUserID, Createdat: mockTime, Updatedat: mockTime},
		}

		mockrepo.EXPECT().GetNotesOfUser(gomock.Any()).Return(testNotes, nil)

		testServer.GetNotes(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		var actualNotes []repository.Note
		json.NewDecoder(rec.Body).Decode(&actualNotes)

		expectedNotes := testNotes
		assert.Equal(t, expectedNotes, actualNotes, "Unexpected response")

	})

	t.Run("Failure case - Userid missing", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodGet, "/api/notes", nil)
		if err != nil {
			t.Fatal("Error creating request:", err)
		}
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		testServer.GetNotes(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

	})

	t.Run("Failure case - error from database", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodGet, "/api/notes?userid="+strconv.FormatUint(mockUserID, 10), nil)
		if err != nil {
			t.Fatal("Error creating request:", err)
		}
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mockrepo.EXPECT().GetNotesOfUser(gomock.Any()).Return(nil, errors.New("Error from database"))

		testServer.GetNotes(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

	})
}

func TestUpdateNoteById(t *testing.T) {

	mockNoteReq := NoteReq{
		Note: "TestNote",
	}
	mockNoteReqBytes, _ := json.Marshal(mockNoteReq)

	mockUserID := uint64(1)
	mockNoteID := uint64(5)

	t.Run("Success case", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/notes/{%d}?userid=%d", mockNoteID, mockUserID), bytes.NewBuffer(mockNoteReqBytes))
		if err != nil {
			t.Fatal("Error creating request:", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatUint(mockNoteID, 10)})
		rec := httptest.NewRecorder()

		mockrepo.EXPECT().UpdateNoteById(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		testServer.UpdateNoteById(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

	})

	t.Run("Failure case - NoteId is missing", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/notes?userid=%d", mockUserID), bytes.NewBuffer(mockNoteReqBytes))
		if err != nil {
			t.Fatal("Error creating request:", err)
		}

		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		testServer.UpdateNoteById(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

	})

	t.Run("Failure case - error from database", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/notes/{%d}?userid=%d", mockNoteID, mockUserID), bytes.NewBuffer(mockNoteReqBytes))
		if err != nil {
			t.Fatal("Error creating request:", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatUint(mockNoteID, 10)})
		rec := httptest.NewRecorder()

		mockrepo.EXPECT().UpdateNoteById(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error from database"))

		testServer.UpdateNoteById(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

	})
}

func TestDeleteNoteById(t *testing.T) {

	mockNoteID := uint64(1)
	mockUserID := uint64(2)

	t.Run("Success case", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/notes/{%d}?userid=%d", mockNoteID, mockUserID), nil)
		if err != nil {
			t.Fatal("Error creating request:", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatUint(mockNoteID, 10)})

		rec := httptest.NewRecorder()

		mockrepo.EXPECT().DeleteNoteById(mockNoteID, mockUserID).Return(nil)

		testServer.DeleteNoteById(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

	})

	t.Run("Failure case - NoteId is missing", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/notes?userid=%d", mockUserID), nil)
		if err != nil {
			t.Fatal("Error creating request:", err)
		}

		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		testServer.DeleteNoteById(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

	})

	t.Run("Failure case - Error from Database", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/notes/{%d}?userid=%d", mockNoteID, mockUserID), nil)
		if err != nil {
			t.Fatal("Error creating request:", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatUint(mockNoteID, 10)})

		rec := httptest.NewRecorder()

		mockrepo.EXPECT().DeleteNoteById(mockNoteID, mockUserID).Return(errors.New("Error from database"))

		testServer.DeleteNoteById(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestShareNoteById(t *testing.T) {

	mockNoteReq := ShareNoteReq{
		SenderId:   1,
		RecieverId: 2,
	}
	mockNoteReqBytes, _ := json.Marshal(mockNoteReq)

	mockNoteReq1 := ShareNoteReq{
		SenderId:   0,
		RecieverId: 2,
	}

	mockNoteReqBytes1, _ := json.Marshal(mockNoteReq1)

	mockNoteID := uint64(1)

	t.Run("Success case", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/notes/{%d}/share", mockNoteID), bytes.NewBuffer(mockNoteReqBytes))
		if err != nil {
			t.Fatal("Error creating request:", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatUint(mockNoteID, 10)})
		rec := httptest.NewRecorder()

		mockrepo.EXPECT().ShareNoteToUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		testServer.ShareNoteById(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

	})

	t.Run("Failure case - Invalid body", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/notes/{%d}/share", mockNoteID), bytes.NewBuffer(mockNoteReqBytes1))
		if err != nil {
			t.Fatal("Error creating request:", err)
		}
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		testServer.ShareNoteById(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

	})

	t.Run("Failure case - error from database", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/notes/{%d}/share", mockNoteID), bytes.NewBuffer(mockNoteReqBytes))
		if err != nil {
			t.Fatal("Error creating request:", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatUint(mockNoteID, 10)})

		rec := httptest.NewRecorder()

		mockrepo.EXPECT().ShareNoteToUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error from database"))

		testServer.ShareNoteById(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

	})
}

func TestGetNoteByKey(t *testing.T) {

	mockUserID := uint64(1)
	mockNoteID := uint64(5)
	mockTime := time.Date(2024, time.January, 5, 18, 42, 48, 0, time.Local)

	t.Run("Success case", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/search?userid=%d&query=%s", mockNoteID, "mocktest"), nil)
		if err != nil {
			t.Fatal("Error creating request:", err)
		}
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		testNotes := []repository.Note{{
			Id: mockNoteID, Note: "Test Note 1 mocktest", Userid: mockUserID, Createdat: mockTime, Updatedat: mockTime},
		}

		mockrepo.EXPECT().GetNotesByKey(gomock.Any(), gomock.Any()).Return(testNotes, nil)

		testServer.GetNoteByKey(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		var actualNotes []repository.Note
		json.NewDecoder(rec.Body).Decode(&actualNotes)

		expectedNotes := testNotes
		assert.Equal(t, expectedNotes, actualNotes, "Unexpected response")

	})

	t.Run("Failure case - error from database", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/search?userid=%d&query=%s", mockNoteID, "mocktest"), nil)
		if err != nil {
			t.Fatal("Error creating request:", err)
		}
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mockrepo.EXPECT().GetNotesByKey(gomock.Any(), gomock.Any()).Return(nil, errors.New("Error from database"))

		testServer.GetNoteByKey(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

	})
}
