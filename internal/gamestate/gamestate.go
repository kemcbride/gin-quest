package gamestate

import (
	"fmt"
	"encoding/json"

	"github.com/kemcbride/gin-quest/internal/room"
)

type GameSave struct {
	X int `json:"x"`
	Y int `json:"y"`
	RoomKey string `json:"roomkey"`
	State int `json:"state"`
}

type GameState struct {
	Save GameSave `json:"save"`
	Room room.Room
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
	loc := gs.Room.GetGridLoc(dx, dy)
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

func (gs *GameState) GetRoomHash() map[string]room.Room {
	var roomMap = map[string]room.Room {
		"mh04i224": room.Room{Id: "mh04i224", Name: "Continent of Euniciar"},
		"mh04dw5i": room.Room{Id: "mh04dw5i", Name: "Land of Patricolia"},
	}
	return roomMap
}

func (gs *GameState) GetRoom(key string) room.Room {
	return gs.GetRoomHash()[key]
}

func (gs *GameState) GetCurrRoom() room.Room {
	return gs.GetRoomHash()[gs.Save.RoomKey]
}

func (gs *GameState) GetCurrRoomName() string {
	return gs.GetRoomHash()[gs.Save.RoomKey].Name
}

func (gs *GameState) GetMapRange(coord int) []int {
	var coordRange []int
	for i := coord - 2; i < coord + 3; i++ {
		coordRange = append(coordRange, i)
	}
	return coordRange
}

func (gs *GameState) GetGridLocClass(x int, y int) string {
	return gs.Room.GetGridLocClass(x, y)
}

func (gs *GameState) GetGridLocImg(x int, y int) string {
	return gs.Room.GetGridLocImg(x, y)
}

func (gs *GameState) NpcHere(x int, y int) bool {
	return gs.Room.NpcHere(x, y)
}

