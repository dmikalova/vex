package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Speed Sigil
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Power
//
//	After a creature enters play, if it is the first creature played this turn, ready it.
var SpeedSigil = card.New(
	"Speed Sigil",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 293),
	card.WithAemberBonus(1),
	card.WithTraits("Power"),
	card.WithAbility(
		card.Trigger.AfterCreatureEnters, card.Conditional{
			Cond: card.FirstCreaturePlayedThisTurn{},
			Then: card.Ready{Target: card.Target.Triggering},
		}),
)
