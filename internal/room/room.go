package room
// Room - aka, not a map since that's a keyword.

type MapGrid = []string


type Location struct {
	X int
	Y int
}

type Npc struct {
	Name string
	Loc Location
	Img string
	Interaction string
}

type Encounter struct {
	Rate float
	Name string
}

type Portal {
	Name string
	Loc Location
	Map string // Map Name / Id for destination map
	DestLoc Location Location // Entry location in destination map
	Img string
}

type Area {
	Name string
	Rectanges [][]Location
	Encounters []Encounter
}

// So, a Room is the main object here.
type Room struct {
	Id string
	Name string
	Grid MapGrid
	Npcs []Npc
	Portals []Portal
	Areas []Area
}

