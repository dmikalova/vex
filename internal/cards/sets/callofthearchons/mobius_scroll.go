package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Mobius Scroll
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Action: Archive Mobius Scroll from play. Archive up to 2 cards from your hand.
var MobiusScroll = set.New(
	"Mobius Scroll",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "130"),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.Action, card.Sequence{
			Effects: []card.Effect{
				card.ArchiveFromPlay{Target: card.Target.This},
				card.ArchiveCard{
					Zone:      card.Hand,
					Selection: card.Chosen{Optional: true},
					Quantity:  card.UpTo{N: card.Fixed(2)},
				},
			},
		}),
)
