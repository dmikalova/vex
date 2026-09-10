package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Power of Fire
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Destroy a friendly creature -> each player loses Æmber equal to half its power, rounded down. Gain 1 chain.
var PowerOfFire = card.New(
	"Power of Fire",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "26"),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{Effects: []card.Effect{
			card.Then{
				First: card.Destroy{Target: card.Target.FriendlyCreature},
				Result: card.LoseAemberEqualTo{
					Player: card.EachPlayer,
					Count:  card.PowerOfChosen{Of: card.HalfRoundedDown},
				},
			},
			card.GainChains{Amount: 1},
		}}),
)
