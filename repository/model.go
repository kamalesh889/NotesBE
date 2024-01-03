package repository

import (
	"time"
)

type Note struct {
	Id        uint64 `gorm:"primaryKey;autoIncrement"`
	Note      string
	Userid    uint64
	Createdat time.Time
	Updatedat time.Time
}

type User struct {
	Id       uint64 `gorm:"primaryKey;autoIncrement"`
	Username string `gorm:"unique"`
	Password string
}

type Sharerecords struct {
	Noteid       uint64 `gorm:"not null"`
	Senderuserid  uint64 `gorm:"not null"`
	Reciveruserid uint64 `gorm:"not null"`
}
