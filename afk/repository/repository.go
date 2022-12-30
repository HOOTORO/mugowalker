package repository

import (
	"errors"
	"time"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const appdata = "afkarena.db"
const locations = "locations.db"

var udb, lcdb *gorm.DB

type Client interface {
	GetDaily() map[string]bool
	UpdateDaily() error
}
type RawLocation struct {
	gorm.Model
	Name string
	Text string
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
	Quests uint8 `gorm:"default:0"`
	UserID uint
}

func DbInit(fn func(string) string) {
	udb = CreateDBConnection(fn(appdata))
	migrateScheme(udb, &User{}, &Daily{})
    lcdb = CreateDBConnection(fn(locations))
	migrateScheme(lcdb, &RawLocation{})
}

func CreateDBConnection(dbname string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dbname), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func migrateScheme(g *gorm.DB, datatypes ...interface{}) {
	// Migrate the schema
	e := g.AutoMigrate(datatypes...)
	if e != nil {
		color.HiWhite("\nerr:%v\nduring run:%v", e, "udb connect")
	}
}

func GetUser(user string) *User {
	var usr User
	r := udb.Where("username = ?", user).First(&usr)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		usr = User{Username: user}
		usr.save()
	}
	return &usr
}
func (u *User) save() {
	r := udb.Save(u)
	color.HiWhite("\nudb: %v user updated", r.RowsAffected)
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
	r := udb.Where("user_id = ? and created_at > ?", u.ID, Bod(time.Now().UTC())).First(&td)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		td = &Daily{Quests: 0, UserID: u.ID}
		udb.Save(td)
	}

	return td
}

func (u *Daily) Update(quest uint8) {
	u.Quests = quest
	r := udb.Save(u)
	color.HiWhite("\nudb: %v user updated", r.RowsAffected)
	if r.Error != nil {
		panic("DB ERROR : " + r.Error.Error())
	}
}

func Bod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
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

func RawLocData(loc, txt string) {
	lcdb.Create(&RawLocation{Name: loc, Text: txt})

}
