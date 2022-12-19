package main

import (
	"worker/adb"
	"worker/afk"
	"worker/bot"

	log "github.com/sirupsen/logrus"
)

func main() {
	const (
		host = "127.0.0.1"
		port = "5555" // MAIN + TEST
	)
	log.SetLevel(log.WarnLevel)

	dev, e := adb.Connect(host, port)
	if e != nil {
		log.Panicf("\ndev:%v\nerr:%v", dev, e)
	}
    game := "afkarena"
    player := "Devitool"

	gm := afk.New(game, player)

	bt := bot.New(dev, gm)

	// bt.gridTap(5, 13)
//	err := bt.Push(gm)
    b ,err := bt.Battle(gm)
	// err := gm.Daily()
	// err := gm.Tower()
    bt.CurrentLoc = gm.GetLocation(afk.BATTLE)
	if err != nil || !b {
		log.Fatalf("MISSION GOES ERRRRRRRRRRRRRRRRRRRRRRRRR%v", err.Error())
	}
}

// =IF(E6=MAX(E10:E79), "Fill some data first!",
// 	IF(AND(E6=0,E5=0),
// 		"KISEKI!!!",
// 		IF(AND(E6=0,E5<>0),
// 			"!UNCLAIMED IS SUM OF HIDDEN PARTS(K)!",
// 			IF(AND(E6>0,(E5/E6)>AVERAGE(D10:D79)),
// 			"!!UNCLAIMED HAS TOP 10 DROP!!",
// 			"UNCLAIMED LEFTOVERS IS BELLOW AVG(<top10)"
// 			)
// 		)

// 	)
// )

// =IF(COUNT(C10:C79)<10, "Need at least 10 records!",
// 	TEXT(
// 		SUMIF(
// 			D10:D79, ">="&D$81, C10:C79
// 			), "0.00%") & " - holds top 10 winners"
// )

// =IF(AND(ISNUMBER(E5), E5>0),
// 	TEXT((E5/C5), "0.00%") & " - unclaimed in " &E6&" chests(Top qty: "&10 - COUNTIF(D10:D79, ">"&D$81)&")",
// 	"Missing chest data"
// )
