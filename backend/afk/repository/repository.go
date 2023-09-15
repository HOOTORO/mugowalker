package repository

import (
	"errors"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/glebarez/sqlite"
	// "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	appdata   = "afkarena.db"
	locations = "locations.db"
	arenaDB   = "arena.db"
)

// var udb, lcdb *gorm.DB
var db *gorm.DB
var log *logrus.Logger

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
	Mode  string
	Level uint
}

type Daily struct {
	UserID uint
	gorm.Model
	Quests uint `gorm:"default:0"`
}

type Button struct {
	gorm.Model
	Name                   string
	Text                   string
	X, Y, Xwmsize, Ywmsize int
}

func init() {
	f, _ := os.OpenFile("db.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)

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

	db, err := gorm.Open(sqlite.Open(arenaDB), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	e := db.AutoMigrate(&User{}, &Daily{}, &Progress{}, &RawLocation{}, &Button{})
	if e != nil {
		log.Errorf("\nerr:%v\nduring run:%v", e, "udb connect")
	}
}

// func DbInit(fn func(string) string) {
// 	f, _ := os.OpenFile(fn("db.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)

// 	log = &logrus.Logger{
// 		Out: f,
// 		Formatter: &logrus.TextFormatter{
// 			ForceColors:               true,
// 			EnvironmentOverrideColors: true,
// 			PadLevelText:              true,
// 			TimestampFormat:           time.Stamp,
// 		},
// 		Level: logrus.TraceLevel,
// 	}

// 	udb = CreateDBConnection(fn(appdata))
// 	migrateScheme(udb, &User{}, &Daily{}, &Progress{})
// 	lcdb = CreateDBConnection(fn(locations))
// 	migrateScheme(lcdb, &RawLocation{}, &Button{})
// }

// func CreateDBConnection(dbname string) *gorm.DB {
// 	db, err := gorm.Open(sqlite.Open(dbname), &gorm.Config{})
// 	if err != nil {
// 		panic("failed to connect database")
// 	}
// 	return db
// }

// func migrateScheme(g *gorm.DB, datatypes ...interface{}) {
// 	// Migrate the schema
// 	e := g.AutoMigrate(datatypes...)
// 	if e != nil {
// 		log.Errorf("\nerr:%v\nduring run:%v", e, "udb connect")
// 	}
// }

func GetUser(user string) *User {
	var usr User
	r := db.Where("username = ?", user).First(&usr)
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

func (u *User) Quests() uint {
	var td *Daily
	r := db.Where("user_id = ? and created_at > ?", u.ID, StartOfDay(time.Now().UTC())).First(&td)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		td = &Daily{Quests: 0, UserID: u.ID}
		db.Save(td)
	}
	return td.Quests
}
func (u *User) SetQuests(q uint) {
	var usrDaily *Daily
	r := db.Where("user_id = ? and created_at > ?", u.ID, StartOfDay(time.Now().UTC())).First(&usrDaily)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		usrDaily = &Daily{Quests: q, UserID: u.ID}
	} else {
		usrDaily.Quests = q
	}
	db.Save(usrDaily)
}

func (u *User) save() {
	r := db.Save(u)
	log.Errorf("\nudb: %v user updated", r.RowsAffected)
	if r.Error != nil {
		panic("DB ERROR : " + r.Error.Error())
	}
}

func (u *User) AfterUpdate(tx *gorm.DB) (err error) {
	u.save()
	return
}

func (p *Progress) Update(level uint) {
	p.Level = level
	r := db.Save(p)
	log.Errorf("\nudb: %v user updated", r.RowsAffected)
	if r.Error != nil {
		panic("DB ERROR : " + r.Error.Error())
	}
}

func (u *User) GetProgress(loc string) *Progress {
	var p *Progress
	r := db.Where("user_id = ? and location_type = ?", u.ID, loc).Order("created_at DESC").First(&p)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		log.Error("No entries rn")
		p = &Progress{Mode: loc, Level: 0, UserID: u.ID}
		db.Save(p)
	}
	return p
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
	db.Create(&RawLocation{Name: loc, Text: txt})
}

func (b Button) String() string {
	return b.Name
}

func (b Button) Offset() (x int, y int) {
	return 0, 0
}

func (b Button) Position() (x int, y int) {
	return b.X, b.Y
}

func GetButtons(xResolution, yResolution int) []*Button {
	var bts []*Button
	r := db.Where("xwmsize = ? and ywmsize = ?", xResolution, yResolution).Find(&bts)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		log.Error("No entries rn")
	}
	return bts
}
func NewBtn(name, text string, x, y, xResolution, yResolution int) *Button {
	b := &Button{
		Name:    name,
		Text:    text,
		X:       x,
		Y:       y,
		Xwmsize: xResolution,
		Ywmsize: yResolution,
	}
	db.Save(b)
	return b
}
