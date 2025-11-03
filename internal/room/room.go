package room
// Room - aka, not a map since that's a keyword.

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"
)

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
	Rate float32
	Name string
}

type Portal struct {
	Name string
	Loc Location
	Map string // Map Name / Id for destination map
	DestLoc Location // Entry location in destination map
	Img string
}

type Area struct {
	Name string
	Rectanges [][]Location
	Encounters []Encounter
}

type Metadata struct {
	Npcs []Npc
	Portals []Portal
	Areas []Area
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

func LoadMap(server embed.FS, roomKey string) MapGrid {
	file, err := server.ReadFile(fmt.Sprintf("static/map/%s/map.txt", roomKey))
	if err != nil {
		panic(err)
	}

	var grid []string
	lines := strings.Split(string(file), "\n")
	grid = append(grid, lines...)
	return grid
}

func (r *Room) LoadMap(server embed.FS) {
	r.Grid = LoadMap(server, r.Id)
}

func (r *Room) LoadMeta(server embed.FS) {
	// metadata := LoadMeta(server, r.Id)
	file, err := server.ReadFile(fmt.Sprintf("static/map/%s/meta.json", r.Id))
	if err != nil {
		panic(err)
	}

	var metajson Metadata
	err = json.Unmarshal(file, &metajson)
	if err != nil {
		panic(err)
	}

	r.Npcs = metajson.Npcs
	r.Portals = metajson.Portals
	r.Areas = metajson.Areas
}

func (r *Room) GetGridLoc(x int, y int) string {
	adjustedX := max(0, len(r.Grid[0]) / 2 + x)
	adjustedY := max(0, len(r.Grid) / 2 + y)
	// Dumb exit to send weird letter we can map to some style
	if ( adjustedX < 0 || adjustedX >= len(r.Grid[0]) ) || (adjustedY < 0 || adjustedY >= len(r.Grid) ) {
		return "Q"
	}
	loc := r.Grid[adjustedY][adjustedX]
	return string(loc)
}

func (r *Room) GetGridLocClass(x int, y int) string {
	var classMap = map[string]string{
		".": "desert",
		"^": "mountain",
		"~": "water",
		"#": "grass",
		"Q": "abyss",
	}
	return classMap[r.GetGridLoc(x, y)]
}

func (r *Room) GetGridLocImg(x int, y int) string {
	s := "" // default do return empty string
	// Use a map lookup instead? Meh.
	for _, npc := range r.Npcs {
		if (npc.Loc.X == x && npc.Loc.Y == y) {
			return npc.Img
		}
	}
	return s
}

func (r *Room) NpcHere(x int, y int) bool {
	for _, npc := range r.Npcs {
		if (npc.Loc.X == x && npc.Loc.Y == y) {
			return true
		}
	}
	return false
}

