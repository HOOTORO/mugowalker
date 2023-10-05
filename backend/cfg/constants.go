package cfg

// ///////////////////////////
// Global afk activities ///
// /////////////////////////
type Mission int

const (
	PushCampain Mission = iota + 1
	ClimbKings
	ClimbWild
	ClimbGrave
	ClimbInferno
	ClimbMaul
	ClimbLight
	ClimbCelestial
	GuildBosses
	Daily
)
const (
	afkapp     = "com.lilithgames.hgame.gp"
	afktestapp = "com.lilithgames.hgame.gp.id"
)
