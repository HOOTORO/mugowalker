package repository

type UserField int

var strs = [...]string{"", "name", "account_id", "vip", "chapter", "stage", "diamonds", "gold"}

const (
	USERNAME UserField = iota + 1
	ACCOUNTID
	VIP
	CHAPTER
	STAGE
	DIAMONDS
	GOLD
)

func (uf UserField) String() string {
	return strs[uf]
}
