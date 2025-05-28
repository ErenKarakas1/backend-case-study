package models

type TeamAttributes struct {
	Attack    float64 `json:"attack"`
	Defense   float64 `json:"defense"`
	Midfield  float64 `json:"midfield"`
	HomeBoost float64 `json:"home_boost"`
}

type PlayStyle string

const (
	PlayStyleAttacking  PlayStyle = "attacking"
	PlayStyleDefensive  PlayStyle = "defensive"
	PlayStylePossession PlayStyle = "possession"
	PlayStyleBalanced   PlayStyle = "balanced"
)

type Team struct {
	ID         int            `json:"id"`
	Name       string         `json:"name"`
	Attributes TeamAttributes `json:"-"`
	PlayStyle  PlayStyle      `json:"-"`
}

func (t *Team) GetOverallRating() float64 {
	return (t.Attributes.Attack + t.Attributes.Defense + t.Attributes.Midfield) / 3.0
}
