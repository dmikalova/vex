package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Nurse Soto
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Human
//
//	Deploy.
//	Play/Fight/Reap: Heal 3 damage from each neighboring creature.
var NurseSoto = set.New(
	"Nurse Soto",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "315"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human),
	card.WithKeywords(card.Keyword.Deploy),
	card.WithAbility(card.Trigger.PlayFightReap, card.Heal{
		Amount: 3,
		Target: card.Target.EachCreature.Neighboring(),
	}),
)
