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
	Name: "Strength",
	Description: "This skill makes you stronger",
	Cap: 9999,
}

var StrengthLevel = SkillLevel{
	Level: 5,
	name: "Strength",
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
