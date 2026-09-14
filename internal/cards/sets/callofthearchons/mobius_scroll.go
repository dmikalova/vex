package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Mobius Scroll
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Action: Archive Mobius Scroll from play, and archive up to 2 cards from your hand.
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
					Amount:    2,
				},
			},
		}),
)
