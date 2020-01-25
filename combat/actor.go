package combat

import (
	"github.com/faiface/pixel"
	"github.com/steelx/go-rpg-cgm/utilz"
)

var DefaultStats = BaseStats{
	HpNow:    300,
	HpMax:    300,
	MpNow:    300,
	MpMax:    300,
	Strength: 10, Speed: 10, Intelligence: 10,
}

// Actor is any creature or character that participates in combat
// and therefore requires stats, equipment, etc
type Actor struct {
	Id, Name   string
	Stats      Stats
	StatGrowth map[string]func() int

	PortraitTexture pixel.Picture
	Portrait        *pixel.Sprite
	Level           int
	XP, NextLevelXP float64
	Actions         []string
	Equipment       Equipment
}

/* example: ActorCreate(HeroDef)
var HeroDef = combat.ActorDef{
		Stats: combat.DefaultStats,
		StatGrowth: map[string]func() int{
			"HpMax":        dice.Create("4d50+100"),
			"MpMax":        dice.Create("2d50+100"),
			"Strength":     combat.StatsGrowth.Fast,
			"Speed":        combat.StatsGrowth.Fast,
			"Intelligence": combat.StatsGrowth.Med,
		},
	}
*/
// ActorCreate
func ActorCreate(def ActorDef) Actor {
	actorAvatar, err := utilz.LoadPicture(def.Portrait)
	utilz.PanicIfErr(err)

	a := Actor{
		Id:              def.Id,
		Name:            def.Name,
		StatGrowth:      def.StatGrowth,
		Stats:           StatsCreate(def.Stats),
		XP:              0,
		Level:           1,
		PortraitTexture: actorAvatar,
		Portrait:        pixel.NewSprite(actorAvatar, actorAvatar.Bounds()),
		Actions:         def.Actions,
	}

	a.NextLevelXP = NextLevel(a.Level)
	return a
}

func (a Actor) ReadyToLevelUp() bool {
	return a.XP >= a.NextLevelXP
}

func (a *Actor) AddXP(xp float64) bool {
	a.XP += xp
	return a.ReadyToLevelUp()
}

func (a Actor) CreateLevelUp() LevelUp {
	levelUp := LevelUp{
		XP:        -a.NextLevelXP,
		Level:     1,
		BaseStats: make(map[string]float64),
	}

	for id, diceRoll := range a.StatGrowth {
		levelUp.BaseStats[id] = float64(diceRoll())
	}

	//Pending feature
	// Additional level up code
	// e.g. if you want to apply
	// a bonus every 4 levels
	// or heal the players MP/HP

	return levelUp
}

func (a *Actor) ApplyLevel(levelUp LevelUp) {
	a.XP += levelUp.XP
	a.Level += levelUp.Level
	a.NextLevelXP = NextLevel(a.Level)

	for k, v := range levelUp.BaseStats {
		a.Stats.Base[k] += v
	}

	//Pending feature
	// Unlock any special abilities etc.
}

type ActorDef struct {
	Id         string //must match entityDef
	Stats      BaseStats
	StatGrowth map[string]func() int
	Portrait   string
	Name       string
	Actions    []string
}

type LevelUp struct {
	XP        float64
	Level     int
	BaseStats map[string]float64
}

type Equipment struct {
	Weapon  string
	Armor   string
	Access1 string
	Access2 string
}