package combat

import (
	"encoding/json"
	"fmt"
)

type CombatStats struct {
	Level     int
	HealthMax int
	Strength  int
	Defence   int
	Avoid     int
	HitRate   int
}

type EnemyData struct {
	Enemies []CombatStats
}

// Load enemies, nto direct combat stats...
func (ed *EnemyData) ToJson() ([]byte, error) {
	j, err := json.Marshal(ed)
	if err != nil {
		return []byte{}, err
	}
	return j, nil
}

func FromJson(b []byte) (*EnemyData, error) {
	ed := &EnemyData{}
	err := json.Unmarshal(b, &ed)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling enemy data: %v, %s", err, b)
	}
	return ed, nil
}

func DeriveCombatStats(level int) (CombatStats, error) {
	if level <= 0 {
		return CombatStats{}, fmt.Errorf("invalid level <= 0 : %d", level)
	}
	return CombatStats{
		Level:     level,
		HealthMax: level + 10,
		Strength:  level,
		Defence:   level,
		Avoid:     level,
		HitRate:   level,
	}, nil
}
