package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Screechbomb
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Weapon
//
//	Versatile.
//	Action: Destroy Screechbomb. Your opponent loses 2 Æmber.
var Screechbomb = card.New(
	"Screechbomb",
	card.House.Brobnar,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 26),
	card.WithTraits("Weapon"),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(
		card.Trigger.Action, card.Sequence{
			Effects: []card.Effect{
				card.Sentence{Effect: card.Destroy{Target: card.Target.This}},
				card.Sentence{Effect: card.LoseAember{
					Player: card.Opponent,
					Amount: 2,
				}},
			},
		}),
)
