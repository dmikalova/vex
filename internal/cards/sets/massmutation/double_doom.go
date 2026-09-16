package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Double Doom
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Put an enemy creature into its owner's hand, and your opponent discards a random card from their hand.
var DoubleDoom = set.New(
	"Double Doom",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "020"),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.PutFromPlay{
				Target:      card.Target.EnemyCreature,
				Destination: card.To.Hand,
			},
			card.DiscardCard{
				Player:    card.Opponent,
				Zone:      card.Hand,
				Selection: card.Random{},
				Amount:    1,
			},
		}}),
)
