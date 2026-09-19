package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Subject Kirby
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Mutant
//
//	Play/Fight/Reap: Play a non-Star Alliance creature.
var SubjectKirby = set.New(
	"Subject Kirby",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "315"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Mutant),
	card.WithAbility(card.Trigger.PlayFightReap, card.PlayFrom{
		From:  card.Hand,
		House: card.Houses.Except(card.House.Self),
		Types: card.Types.Of(card.Type.Creature),
	}),
)
