package repository

import (
	"database/sql"
	"errors"

	// "worker/game"

	"github.com/fatih/color"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const appdata = "afkarena.db"

var db *gorm.DB

type Client interface {
	GetDaily() map[string]bool
	UpdateDaily() error
}

type User struct {
	gorm.Model
	Username  string
	AccountId int
	Vip       int
	Chapter   int
	Stage     int
	Diamonds  int
	Gold      int
	Daily     []Daily
}

type Daily struct {
	gorm.Model
	Loot          sql.NullBool `gorm:"default:false"`
	FastRewards   sql.NullBool `gorm:"default:false"`
	FRqty         uint8        `gorm:"default:1"`
	Likes         sql.NullBool `gorm:"default:false"`
	GuildBoss     sql.NullBool `gorm:"default:false"`
	Arena         sql.NullBool `gorm:"default:false"`
	ArenaTopEnemy sql.NullBool `gorm:"default:false"`
	UserID        uint
}

func init() {
	db = CreateDBConnection(appdata)
}

func GetUser(user string) *User {
	var usr User
	r := db.Where("username = ?", user).First(&usr)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		usr = User{Username: user}
	}
	return &usr
}

func CreateDBConnection(dbname string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dbname), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	e := db.AutoMigrate(&User{}, &Daily{})
	if e != nil {
		color.HiWhite("\nerr:%v\nduring run:%v", e, "db connect")
	}
	return db
}

func (u *User) SaveUserInfo() {
	r := db.Save(u)
	color.HiWhite("\ndb: %v user updated", r.RowsAffected)
	if r.Error != nil {
		panic("DB ERROR : " + r.Error.Error())
	}
}

// func PrevousRun(p interface{}) (string, bool) {
// 	db := CreateDBConnection(appdata)

// 	var cnf User
// 	res := db.Last(&cnf)
// 	if res.Error != nil {
// 		panic("DB ERROR ON PREV RUN")
// 	}
// 	// if res.RowsAffected > 0 {
// 	// 	f, ok := Last(cnf.User, p)
// 	// 	return f, ok
// 	// }
// 	return "", false
// }
