package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Mushroom Man
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Fungus • Human
//
//	Mushroom Man gains +3 power for each unforged key you have.
var MushroomMan = card.New(
	"Mushroom Man",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 362),
	card.WithPower(2),
	card.WithTraits("Fungus", "Human"),
	card.WithConstant(card.ConstantAbility{
		Target:     card.Target.This,
		PowerBonus: 3,
		Per:        card.UnforgedKeys{Player: card.Controller},
	}),
)
