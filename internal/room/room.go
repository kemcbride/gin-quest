package room
// Room - aka, not a map since that's a keyword.

import (
	"embed"
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
