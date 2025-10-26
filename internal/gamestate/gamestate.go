package gamestate

import (
	"fmt"
	"strings"
	"strconv"
)

type Room struct {
	Id int
	Name string
	Path string
}

type GameState struct {
	x int
	y int
	room int
	currGrid []string
}

func SerializeGameState(g GameState) string {
	return fmt.Sprintf("%d,%d,%d", g.x, g.y, g.room)
}

func DeserializeGameState(s string) GameState {
	segments := strings.Split(s, ",")
	var values [3]int
	for i, seg := range segments {
		val, err := strconv.Atoi(seg)
		if err != nil {
			// Lazy lazy bad
			panic(err)
		}
		values[i] = val
	}
	return GameState{
		x: values[0],
		y: values[1],
		room: values[2],
	}
}

func (gs GameState) GetStatusBlurb() string {
	return "Status: All good!"
}

func (gs GameState) GetX() int {
	return gs.x
}

func (gs GameState) GetY() int {
	return gs.y
}

func (gs *GameState) AddX(a int) int {
	return gs.x + a
}

func (gs *GameState) AddY(a int) int {
	return gs.y + a
}

func (gs GameState) CanMove(dx int, dy int) bool {
	// For now, we just want to check if it's water or mountain
	// Or, the edge of the map.
	// Eventually could handle a case where we fly or swim.
	loc := gs.GetGridLoc(dx, dy)
	switch loc {
	case	".":
		return true
	case "^":
		return false
	case "~":
		return false
	case "Q": // placeholder for edge of map?
		return false
	default:
		return true
	}
}

func (gs *GameState) MoveUp() {
	if gs.CanMove(gs.x, gs.y-1) {
		gs.y--;
	}
}

func (gs *GameState) MoveDown() {
	if gs.CanMove(gs.x, gs.y+1) {
		gs.y++;
	}
}

func (gs *GameState) MoveLeft() {
	if gs.CanMove(gs.x-1, gs.y) {
		gs.x--;
	}
}

func (gs *GameState) MoveRight() {
	if gs.CanMove(gs.x+1, gs.y) {
		gs.x++;
	}
}

func (gs GameState) GetGrid() []string {
	return gs.currGrid
}

func (gs *GameState) SetGrid(grid []string) {
	gs.currGrid = grid
}

func (gs GameState) GetGridLoc(x int, y int) string {
	adjustedX := max(0, len(gs.currGrid[0]) / 2 + x)
	adjustedY := max(0, len(gs.currGrid) / 2 + y)
	// Dumb exit to send weird letter we can map to some style
	if ( adjustedX < 0 || adjustedX >= len(gs.currGrid[0]) ) || (adjustedY < 0 || adjustedY >= len(gs.currGrid) ) {
		return "Q"
	}
	loc := gs.GetGrid()[adjustedY][adjustedX]
	return string(loc)
}

func (gs GameState) GetGridLocColor(x int, y int) string {
	var colorMap = map[string]string{
		".": "#FFE4B5",
		"^": "#444444",
		"~": "#4682B4",
		"#": "#9ACD32",
		"Q": "#000000",
	}
	return colorMap[gs.GetGridLoc(x, y)]
}

func (gs GameState) GetGridLocClass(x int, y int) string {
	var classMap = map[string]string{
		".": "desert",
		"^": "mountain",
		"~": "water",
		"#": "grass",
		"Q": "abyss",
	}
	return classMap[gs.GetGridLoc(x, y)]
}

func (gs GameState) GetRoomHash() map[int]Room {
	var roomMap = map[int]Room {
		0: Room{Id: 0, Name: "Continent of Euniciar", Path: "map-mh04i224.txt"},
		1: Room{Id: 1, Name: "Land of Patricolia", Path: "map-mh04dw5i.txt"},
	}
	return roomMap
}

func (gs GameState) GetRoom(i int) Room {
	return gs.GetRoomHash()[i]
}

func (gs GameState) GetCurrRoom() Room {
	return gs.GetRoomHash()[gs.room]
}

func (gs GameState) GetCurrRoomName() string {
	return gs.GetRoomHash()[gs.room].Name
}

func (gs GameState) GetMapRange(coord int) []int {
	var coordRange []int
	for i := coord - 2; i < coord + 3; i++ {
		coordRange = append(coordRange, i)
	}
	return coordRange
}
