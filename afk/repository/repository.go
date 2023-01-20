package repository

import (
	"errors"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	appdata   = "afkarena.db"
	locations = "locations.db"
)

var udb, lcdb *gorm.DB
var log *logrus.Logger

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
	AccountId uint
	Vip       uint
	Diamonds  uint
	Gold      uint
	Locations []Progress `gorm:"foreignKey:UserID"`
	Dailies   []Daily    `gorm:"foreignKey:UserID"`
}

type Progress struct {
	UserID uint
	gorm.Model
	Location uint
	Level    uint
}

type Daily struct {
	UserID uint
	gorm.Model
	Quests uint8 `gorm:"default:0"`
}

func DbInit(fn func(string) string) {
	f, _ := os.OpenFile(fn("db.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)

	log = &logrus.Logger{
		Out: f,
		Formatter: &logrus.TextFormatter{
			ForceColors:               true,
			EnvironmentOverrideColors: true,
			PadLevelText:              true,
			TimestampFormat:           time.Stamp,
		},
		Level: logrus.TraceLevel,
	}

	udb = CreateDBConnection(fn(appdata))
	migrateScheme(udb, &User{}, &Daily{}, &Progress{})
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
		log.Errorf("\nerr:%v\nduring run:%v", e, "udb connect")
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
func (u *User) Id() uint {
	return u.ID
}
func (u *User) Name() string {
	return u.Username
}

func (u *User) Quests() uint8 {
	var td *Daily
	r := udb.Where("user_id = ? and created_at > ?", u.ID, StartOfDay(time.Now().UTC())).First(&td)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		td = &Daily{Quests: 0, UserID: u.ID}
		udb.Save(td)
	}
	return td.Quests
}
func (u *User) SetQuests(q uint8) {
	var usrDaily *Daily
	r := udb.Where("user_id = ? and created_at > ?", u.ID, StartOfDay(time.Now().UTC())).First(&usrDaily)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		usrDaily = &Daily{Quests: q, UserID: u.ID}
	} else {
		usrDaily.Quests = q
	}
	udb.Save(usrDaily)
}

func (u *User) save() {
	r := udb.Save(u)
	log.Errorf("\nudb: %v user updated", r.RowsAffected)
	if r.Error != nil {
		panic("DB ERROR : " + r.Error.Error())
	}
}

func (u *User) AfterUpdate(tx *gorm.DB) (err error) {
	u.save()
	return
}

// func (u *Daily) Update(quest uint8) {
// 	u.Quests = quest
// 	r := udb.Save(u)
// 	log.Errorf("\nudb: %v user updated", r.RowsAffected)
// 	if r.Error != nil {
// 		panic("DB ERROR : " + r.Error.Error())
// 	}
// }

func (u *Progress) Update(level uint) {
	u.Level = level
	r := udb.Save(u)
	log.Errorf("\nudb: %v user updated", r.RowsAffected)
	if r.Error != nil {
		panic("DB ERROR : " + r.Error.Error())
	}
}

func (u *User) GetProgress(loc uint) *Progress {
	var p *Progress
	r := udb.Where("user_id = ? and location = ?", u.ID, loc).Order("created_at DESC").First(&p)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		log.Error("No entries rn")
		p = &Progress{Location: loc, Level: 0, UserID: u.ID}
		udb.Save(p)
	}
	return p
}

func (u *User) LocLevel(loc, level uint) {
	u.Locations = append(u.Locations, Progress{Location: loc, Level: level})
	udb.Save(u)
}

func StartOfDay(t time.Time) time.Time {
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
