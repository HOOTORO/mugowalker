package repository

import (
    "errors"
    "time"
    // "worker/game"

    "github.com/fatih/color"
    log "github.com/sirupsen/logrus"
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
	Dailies   []Daily `gorm:"foreignKey:UserID"`
}

type Daily struct {
	gorm.Model
    Quests     uint8 `gorm:"default:0"`
	UserID     uint
}

func init() {
	db = CreateDBConnection(appdata)
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

func GetUser(user string) *User {
	var usr User
	r := db.Where("username = ?", user).First(&usr)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		usr = User{Username: user}
		usr.save()
	}
	return &usr
}
func (u *User) save() {
	r := db.Save(u)
	color.HiWhite("\ndb: %v user updated", r.RowsAffected)
	if r.Error != nil {
		panic("DB ERROR : " + r.Error.Error())
	}
}

func (u *User) AfterUpdate(tx *gorm.DB) (err error) {

	u.save()

	return
}

func (u *User) ActiveQuests() *Daily {
	var td *Daily
	r := db.Where("user_id = ? and created_at > ?", u.ID, Bod(time.Now())).First(&td)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		td = &Daily{Quests: 0, UserID: u.ID}
		db.Save(td)
	}

	return td
}

func (u *Daily) Update(quest uint8) {
	u.Quests = quest
	r := db.Save(u)
	color.HiWhite("\ndb: %v user updated", r.RowsAffected)
	if r.Error != nil {
		panic("DB ERROR : " + r.Error.Error())
	}
}


func Bod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func Truncate(t time.Time) time.Time {
	return t.Truncate(24 * time.Hour)
}

func NowInMoscow() time.Time {
	moscow, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Warn(err)
		return time.Now()
	}
	return time.Now().In(moscow)

}
