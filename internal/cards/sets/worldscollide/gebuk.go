package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Gebuk
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Traits: Beast
//
//	Destroyed: Discard the top card of your deck. If it is a creature, swap it with Gebuk.
var Gebuk = set.New(
	"Gebuk",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "373"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Beast),
	card.WithAbility(
		card.Trigger.Destroyed, card.Sequence{Effects: []card.Effect{
			card.DiscardTop{
				Amount: 1,
				Player: card.Controller,
			},
			card.Conditional{
				Cond: card.ItIs{Type: card.Type.Creature},
				Then: card.Swap{
					With:        card.Target.TheOtherCreature,
					FromContext: true,
				},
			},
		}},
	),
)
