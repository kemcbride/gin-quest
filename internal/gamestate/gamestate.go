package gamestate

import (
	"embed"
	"encoding/json"
	"fmt"

	"github.com/kemcbride/gin-quest/internal/room"
	"github.com/kemcbride/gin-quest/internal/skills"
)

type State int // Kind of like "mode" for the UI/potential actions

const (
	StateExplore State = iota
	StateTalk
	StateBattle
	StateSkills
	StateInventory
	StateOptions
)

type GameSave struct {
	X           int            `json:"x"`
	Y           int            `json:"y"`
	RoomKey     string         `json:"roomkey"`
	State       State          `json:"state"`
	Level       int            `json:"level"`
	Name        string         `json:"name"`
	SkillLevels map[string]int `json:"skilllevels"`
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

	if gs.Level <= 0 {
		return nil, fmt.Errorf("error loading gamesave from json: level cannot be <= 0, %d", gs.Level)
	}
	return gs, nil
}

func NewSkillLevels() map[string]int {
	return skills.NewSkillLevels()
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

func (gs *GameState) GetName() string {
	return gs.Save.Name
}

func (gs *GameState) GetLevel() int {
	return gs.Save.Level
}

func (gs *GameState) CanMove(dx int, dy int) bool {
	if gs.Save.State != StateExplore {
		return false
	}
	// For now, we just want to check if it's water or mountain
	// Or, the edge of the map.
	// Eventually could handle a case where we fly or swim.
	loc := gs.Room.GetGridLoc(dx, dy)
	switch loc {
	case ".":
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
		gs.Save.Y--
	}
}

func (gs *GameState) MoveDown() {
	if gs.CanMove(gs.Save.X, gs.Save.Y+1) {
		gs.Save.Y++
	}
}

func (gs *GameState) MoveLeft() {
	if gs.CanMove(gs.Save.X-1, gs.Save.Y) {
		gs.Save.X--
	}
}

func (gs *GameState) MoveRight() {
	if gs.CanMove(gs.Save.X+1, gs.Save.Y) {
		gs.Save.X++
	}
}

func (gs *GameState) Portal(server embed.FS) error {
	if gs.Room.PortalHere(gs.Save.X, gs.Save.Y) {
		// We know there's a portal where the Protagonist is standing.
		portal, found := gs.Room.GetPortal(gs.Save.X, gs.Save.Y)
		if !found {
			return nil // return without modifying game state, eg. going thru portal
		}
		// Update GameState to be in new room at destloc from portal
		gs.Save.State = StateExplore
		gs.Save.RoomKey = portal.Map
		gs.Save.X = portal.DestLoc.X
		gs.Save.Y = portal.DestLoc.Y
		gs.Room = gs.GetCurrRoom()
		gs.Room.LoadMap(server)
		gs.Room.LoadMeta(server)
	}
	return nil
}

func (gs *GameState) Talk(server embed.FS) error {
	_, found := gs.Room.GetNpc(gs.Save.X, gs.Save.Y)
	if !found {
		return nil // return without modifying game state, eg. going thru portal
	}
	// Update GameState to be in new room at destloc from portal
	gs.Save.State = StateTalk
	// This breaks everything LOL and also how do we know that it's happening??
	return nil
}

func (gs *GameState) GetRoomHash() map[string]room.Room {
	var roomMap = map[string]room.Room{
		"mh04i224":  room.Room{Id: "mh04i224", Name: "Continent of Euniciar"},
		"mh04dw5i":  room.Room{Id: "mh04dw5i", Name: "Land of Patricolia"},
		"chia-town": room.Room{Id: "chia-town", Name: "Chia Town"},
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

func (gs *GameState) GetCurrRoomDescription() string {
	// TODO: Have this use the description fields from meta.json
	return fmt.Sprintf("You're somewhere in %s.", gs.Room.Name)
}

func (gs *GameState) GetMapRange(coord int, size int) []int {
	var coordRange []int
	for i := coord - size; i < coord+size+1; i++ {
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

func (gs *GameState) ProtagHere(x int, y int) bool {
	return gs.Save.X == x && gs.Save.Y == y
}

func (gs *GameState) NpcHere(x int, y int) bool {
	return gs.Room.NpcHere(x, y)
}

func (gs *GameState) PortalHere(x int, y int) bool {
	return gs.Room.PortalHere(x, y)
}

func (gs *GameState) GetSkills() map[string]skills.Skill {
	return skills.Skills
}

func (save *GameSave) GetSkillLevels() map[string]int {
	return save.SkillLevels
}

func (save *GameSave) GetUnusedSkillPoints() int {
	// Let's say every 5 levels (incl. level 1) you get 3 points
	// and each level you get one otherwise
	// so if you're level 15, you have 3 + 3 + 3 + 3 (1, 5, 10, 15)
	// plus one for each other level (2, 3, 4, 6, 7, 8, 9, 11, 12, 13, 14)
	// aka: 3x(lvl/5 +1) + (lvl - (lvl/5 - 1))
	landmarkPoints := 3 * (save.Level/5 + 1)
	normalPoints := save.Level - max((save.Level/5-1), 0) - 1
	earnedPoints := landmarkPoints + normalPoints
	usedPoints := 0
	for _, skillLevel := range save.SkillLevels {
		usedPoints += skillLevel
	}
	return max(earnedPoints-usedPoints, 0)
}
