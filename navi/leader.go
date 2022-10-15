package navi

import "worker/adb"

type Leaderboard struct {
	Region  int
	Players []Player
	*adb.Device
}

type Player struct {
	place   int
	name    string
	damage  string
	details []FightStat
}

type FightStat struct {
	tree    string
	herodps map[string]string
}

func (board *Leaderboard) SetRegion(n int) {
	//rectangle x29, y12,  x72, y15 (rel)
	//need recognition
	board.GoForward(28, 13)
	//magic

}

// func (board *Leaderboard) PlayersArea() []Player {
// 	// rectangle x1, y23, x99, y77 (rel)
// 	//till y33 P1, y 43 .... 5 playes per area
// 	//  top 100 players per region  , 31 region
// 	// all 200

// }

func (BD *Leaderboard) NextFive() []Player
