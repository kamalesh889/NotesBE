package repository

import (
	"errors"
	"fmt"
	"log"
	"time"
)

func (r *Database) CreateUser(req *User) error {

	// check for username exists or not

	result := r.DbConn.Create(req)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *Database) GetUser(req *User) (uint64, error) {

	var id uint64

	query := "select id from user where username = ? and password = ? ;"

	err := r.DbConn.Raw(query, req.Username, req.Password).Scan(id).Error
	if err != nil {
		log.Println("Error in Fetching User Order details", err)
		return 0, err
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

	query := "select * from note where userid = ? ; "

	err := r.DbConn.Raw(query, userid).Scan(usernotes).Error
	if err != nil {
		log.Println("Error in Fetching Notes details of User", err)
		return nil, err
	}

	sharednotes := []Note{}

	query = `SELECT *
    FROM note
    JOIN sharerecords ON notes.id = sharerecords.noteid
    WHERE sharerecords.reciveruserid = ? ;`

	err = r.DbConn.Raw(query, userid).Scan(sharednotes).Error
	if err != nil {
		log.Println("Error in Fetching shared Notes details of User", err)
		return nil, err
	}

	usernotes = append(usernotes, sharednotes...)

	return usernotes, nil

}

func (r *Database) GetNoteById(noteId uint64) (string, error) {

	var note string

	query := "select note from note where id = ? ;"

	err := r.DbConn.Raw(query, noteId).Scan(note).Error
	if err != nil {
		log.Println("Error in Fetching Notes detail", err)
		return "", err
	}

	return note, nil

}

func (r *Database) UpdateNoteById(noteId uint64, note string) error {

	query := fmt.Sprintf("update note set note = '%s' and updatedat = '%v' where id = %d ;", note, time.Now(), noteId)

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

func (r *Database) DeleteNoteById(noteId uint64) error {

	query := fmt.Sprintf("delete from note where id = %d ;", noteId)

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

	query := `select * from note WHERE userid = ? and note @@ to_tsquery('english', ?);`

	err := r.DbConn.Raw(query, userid, key).Scan(noteRecords).Error
	if err != nil {
		log.Println("Error in Fetching Notes", err)
		return nil, err
	}

	return noteRecords, nil

}
