package connection

import (
	"NOTESBE/repository"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitializeDB() (*repository.Database, error) {

	dsn := "host=" + viper.GetString("database.host") + " user=" + viper.GetString("database.user") + " password=" + viper.GetString("database.password") + " dbname=" + viper.GetString("database.name") + " port=" + viper.GetString("database.port") + " sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	if viper.GetBool("database.migrate") {
		Migrate(db)
	}

	return &repository.Database{DbConn: db}, nil
}

func Migrate(db *gorm.DB) {

	err := db.AutoMigrate(
		&repository.User{}, &repository.Note{}, &repository.Sharerecords{},
	)
	if err != nil {
		log.Fatalln(err)
		return
	}

	db.Exec("CREATE EXTENSION btree_gin;")
	db.Exec("CREATE INDEX idx ON notes USING GIN (userid, to_tsvector('english', note));")

}
