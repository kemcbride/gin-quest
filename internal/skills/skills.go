package skills

import ()

// Note - the struct OUGHT to have something to do with
// how the skill actually works but we can cross that bridge when we get to it.
type Skill struct {
	Name        string
	Description string
	Cap         int
}

type SkillLevel struct {
	Level int

	name string
}

var StrengthSkill = Skill{
	Name:        "Strength",
	Description: "makes you stronger",
	Cap:         9999,
}

var Skills = map[string]Skill{
	"Vitality": Skill{
		Name:        "Vitality",
		Description: "raises your max health",
		Cap:         100,
	},
	"Strength": StrengthSkill,
	"Defence": Skill{
		Name:        "Defence",
		Description: "raises your defence",
		Cap:         100,
	},
	"Dodge": Skill{
		Name:        "Dodge",
		Description: "raises your avoidance",
		Cap:         100,
	},
	"Accuracy": Skill{
		Name:        "Accuracy",
		Description: "raises your hit rate",
		Cap:         100,
	},
}

var StrengthLevel = SkillLevel{
	Level: 5,
	name:  "Strength",
}

func NewSkillLevels() map[string]int {
	m := make(map[string]int)
	for name := range Skills {
		m[name] = 0
	}
	return m
}

func (s *Skill) GetName() string {
	return s.Name
}

func (s *Skill) GetDescription() string {
	return s.Description
}

func (s *Skill) GetCap() int {
	return s.Cap
}

func (sl *SkillLevel) GetName() string {
	return sl.name
}
