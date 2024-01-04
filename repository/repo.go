package repository

import (
	"errors"
	"fmt"
	"log"
	"time"
)

type Repository interface {
	CreateUser(req *User) error
	GetUser(req *User) (uint64, error)
	CreateNote(req *Note) error
	GetNotesOfUser(userid uint64) ([]Note, error)
	GetNoteById(noteId, userid uint64) (*Note, error)
	UpdateNoteById(noteId, userid uint64, note string) error
	DeleteNoteById(noteId, userid uint64) error
	ShareNoteToUser(noteId, senderuserid, recieveruserid uint64) error
	GetNotesByKey(userid uint64, key string) ([]Note, error)
}

func (r *Database) CreateUser(req *User) error {

	result := r.DbConn.Create(req)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *Database) GetUser(req *User) (uint64, error) {

	var id uint64

	query := "select id from users where username = ? and password = ? ;"

	err := r.DbConn.Raw(query, req.Username, req.Password).Scan(&id).Error
	if err != nil {
		log.Println("Error in Fetching User Order details", err)
		return 0, err
	}

	if id == 0 {
		return 0, errors.New("User does not exist in records")
	}

	return id, nil

}

func (r *Database) CreateNote(req *Note) error {

	req.Createdat = time.Now()
	req.Updatedat = time.Now()

	result := r.DbConn.Create(req)

	if result.Error != nil {
		return result.Error
	}

	return nil

}

func (r *Database) GetNotesOfUser(userid uint64) ([]Note, error) {

	usernotes := []Note{}

	query := "select * from notes where userid = ? ; "

	err := r.DbConn.Raw(query, userid).Scan(&usernotes).Error
	if err != nil {
		log.Println("Error in Fetching Notes details of User", err)
		return nil, err
	}

	sharednotes := []Note{}

	query = `SELECT *
    FROM notes
    JOIN sharerecords ON notes.id = sharerecords.noteid
    WHERE sharerecords.reciveruserid = ? ;`

	err = r.DbConn.Raw(query, userid).Scan(&sharednotes).Error
	if err != nil {
		log.Println("Error in Fetching shared Notes details of User", err)
		return nil, err
	}

	usernotes = append(usernotes, sharednotes...)

	return usernotes, nil

}

func (r *Database) GetNoteById(noteId, userid uint64) (*Note, error) {

	noteInfo := &Note{}

	query := "select * from notes where id = ? and userid = ? ;"

	err := r.DbConn.Raw(query, noteId, userid).Scan(noteInfo).Error
	if err != nil {
		log.Println("Error in Fetching Notes detail", err)
		return nil, err
	}

	if noteInfo.Id == 0 {
		return nil, errors.New("Note does not exist in records")
	}

	return noteInfo, nil

}

func (r *Database) UpdateNoteById(noteId, userid uint64, note string) error {

	query := fmt.Sprintf("update notes set note = '%s' , updatedat = '%v' where id = %d and userid = %d  ;", note, time.Now().Format("2006-01-02 15:04:05.999999-07:00"), noteId, userid)

	result := r.DbConn.Exec(query)

	if result.Error != nil {
		log.Println("Error in Updating Note", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("this note doesn't exist in records")
	}

	return nil

}

func (r *Database) DeleteNoteById(noteId, userid uint64) error {

	query := fmt.Sprintf("delete from notes where id = %d and userid = %d ;", noteId, userid)

	result := r.DbConn.Exec(query)

	if result.Error != nil {
		log.Println("Error in Deleting Note", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("this note doesn't exist in records")
	}

	return nil

}

func (r *Database) ShareNoteToUser(noteId, senderuserid, recieveruserid uint64) error {

	user := User{}

	err := r.DbConn.First(&user, recieveruserid).Error
	if err != nil {
		log.Println("Error in Getting the Reciever User", err)
		return err
	}

	_, err = r.GetNoteById(noteId, senderuserid)
	if err != nil {
		log.Println("Error in Getting the Note ", err)
		return err
	}

	// check with noteid and recieveruserid that record exists or not in  Sharerecords

	shareInfo := &Sharerecords{
		Noteid:        noteId,
		Reciveruserid: recieveruserid,
		Senderuserid:  senderuserid,
	}

	err = r.DbConn.Create(shareInfo).Error
	if err != nil {
		return err
	}

	return nil

}

func (r *Database) GetNotesByKey(userid uint64, key string) ([]Note, error) {

	noteRecords := []Note{}

	query := `select * from notes WHERE userid = ? and note @@ to_tsquery('english', ?);`

	err := r.DbConn.Raw(query, userid, key).Scan(&noteRecords).Error
	if err != nil {
		log.Println("Error in Fetching Notes", err)
		return nil, err
	}

	sharednoteids := []uint64{}
	sharedNotes := []Note{}

	query = `SELECT id
    FROM notes
    JOIN sharerecords ON notes.id = sharerecords.noteid
    WHERE sharerecords.reciveruserid = ? ;`

	err = r.DbConn.Raw(query, userid).Scan(&sharednoteids).Error
	if err != nil {
		log.Println("Error in Fetching shared Notes details of User", err)
		return nil, err
	}

	query = `select * from notes where id in(?) and note @@ to_tsquery('english', ?);`
	err = r.DbConn.Raw(query, sharednoteids, key).Scan(&sharedNotes).Error
	if err != nil {
		log.Println("Error in Fetching shared Notes details of User", err)
		return nil, err
	}

	noteRecords = append(noteRecords, sharedNotes...)

	return noteRecords, nil

}
