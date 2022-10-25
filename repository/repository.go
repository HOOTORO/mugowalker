package repository

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const appdata = "data.db"

type User struct {
	gorm.Model
	User     string
	Filename string
}

func CreateDBConnection(dbname string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dbname), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	db.AutoMigrate(&User{})
	return db
}

func DbSave(user, fname string) {
	db := CreateDBConnection(appdata)

	var data User
	usr := User{User: user, Filename: fname}
	result := db.FirstOrCreate(&data, usr)
	if result.Error != nil {
		panic("DB ERROR : " + result.Error.Error())
	}
}

func PrevousRun(p interface{}) (string, bool) {
	db := CreateDBConnection(appdata)

	var cnf User
	res := db.Last(&cnf)
	if res.Error != nil {
		panic("DB ERROR ON PREV RUN")
	}
	// if res.RowsAffected > 0 {
	// 	f, ok := Last(cnf.User, p)
	// 	return f, ok
	// }
	return "", false
}
