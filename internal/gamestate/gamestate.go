package gamestate

import (
	"fmt"
	"encoding/json"
)

type Room struct {
	Id string
	Name string
}

type GameSave struct {
	X int `json:"x"`
	Y int `json:"y"`
	Room string `json:"room"`
	State int `json:"state"`
}

type GameState struct {
	Save GameSave `json:"save"`
	CurrGrid []string `json:"curr_grid,omitempty"`
}

// json serialization as methods
func (save *GameSave) ToJson() ([]byte, error) {
	j, err := json.Marshal(save)
	if err != nil {
		return []byte{}, err
	}
	return j, nil
}

func GameSaveFromJson(b []byte) (*GameSave, error) {
	gs := &GameSave{}
	err := json.Unmarshal(b, &gs)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling gamesave: %v, %s", err, b)
	}
	return gs, nil
}

func (gs *GameState) ToJson() ([]byte, error) {
	j, err := json.Marshal(gs)
	if err != nil {
		return []byte{}, err
	}
	return j, nil
}

func FromJson(b []byte) (*GameState, error) {
	gs := &GameState{}
	err := json.Unmarshal(b, &gs)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling gamestate: %v, %s", err, b)
	}
	return gs, nil
}

// Stringer
func (gs *GameState) String() string {
	j, err := json.Marshal(gs)
	if err != nil {
		// TODO be better - this is lazy
		panic(err)
	}
	return string(j)
}

func (gs *GameState) GetStatusBlurb() string {
	return "Status: All good!"
}

func (gs *GameState) CanMove(dx int, dy int) bool {
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
	if gs.CanMove(gs.Save.X, gs.Save.Y-1) {
		gs.Save.Y--;
	}
}

func (gs *GameState) MoveDown() {
	if gs.CanMove(gs.Save.X, gs.Save.Y+1) {
		gs.Save.Y++;
	}
}

func (gs *GameState) MoveLeft() {
	if gs.CanMove(gs.Save.X-1, gs.Save.Y) {
		gs.Save.X--;
	}
}

func (gs *GameState) MoveRight() {
	if gs.CanMove(gs.Save.X+1, gs.Save.Y) {
		gs.Save.X++;
	}
}

func (gs *GameState) GetGridLoc(x int, y int) string {
	adjustedX := max(0, len(gs.CurrGrid[0]) / 2 + x)
	adjustedY := max(0, len(gs.CurrGrid) / 2 + y)
	// Dumb exit to send weird letter we can map to some style
	if ( adjustedX < 0 || adjustedX >= len(gs.CurrGrid[0]) ) || (adjustedY < 0 || adjustedY >= len(gs.CurrGrid) ) {
		return "Q"
	}
	loc := gs.CurrGrid[adjustedY][adjustedX]
	return string(loc)
}

func (gs *GameState) GetGridLocColor(x int, y int) string {
	var colorMap = map[string]string{
		".": "#FFE4B5",
		"^": "#444444",
		"~": "#4682B4",
		"#": "#9ACD32",
		"Q": "#000000",
	}
	return colorMap[gs.GetGridLoc(x, y)]
}

func (gs *GameState) GetGridLocClass(x int, y int) string {
	var classMap = map[string]string{
		".": "desert",
		"^": "mountain",
		"~": "water",
		"#": "grass",
		"Q": "abyss",
	}
	return classMap[gs.GetGridLoc(x, y)]
}

func (gs *GameState) GetRoomHash() map[string]Room {
	var roomMap = map[string]Room {
		"mh04i224": Room{Id: "mh04i224", Name: "Continent of Euniciar"},
		"mh04dw5i": Room{Id: "mh04dw5i", Name: "Land of Patricolia"},
	}
	return roomMap
}

func (gs *GameState) GetRoom(i string) Room {
	return gs.GetRoomHash()[i]
}

func (gs *GameState) GetCurrRoom() Room {
	return gs.GetRoomHash()[gs.Save.Room]
}

func (gs *GameState) GetCurrRoomName() string {
	return gs.GetRoomHash()[gs.Save.Room].Name
}

func (gs *GameState) GetMapRange(coord int) []int {
	var coordRange []int
	for i := coord - 2; i < coord + 3; i++ {
		coordRange = append(coordRange, i)
	}
	return coordRange
}
