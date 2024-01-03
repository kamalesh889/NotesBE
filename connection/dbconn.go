package connection

import (
	"NotesBe/repository"
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

	isMigrate := true // change it
	if isMigrate {
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

	db.Exec("CREATE INDEX idx ON note USING GIN (note);")

}
