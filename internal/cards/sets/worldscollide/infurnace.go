package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Infurnace
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Demon
//
//	Play: Purge 2 cards from a discard pile. For each Æmber bonus icon on the purged cards, your opponent loses 1 Æmber.
var Infurnace = set.New(
	"Infurnace",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "78"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Demon),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.PurgeCard{
				Zones:     []card.Zone{card.Discard},
				Player:    card.ChosenPlayer,
				Selection: card.Chosen{},
				Quantity:  card.Takes{N: card.Fixed(2)},
			},
			card.LoseAember{
				Player: card.Opponent,
				Amount: 1,
				Per: card.BonusIconsOf{
					Over: card.ThePurgedCards{},
					Kind: card.Bonus.Aember,
				},
			},
		}},
	),
)
