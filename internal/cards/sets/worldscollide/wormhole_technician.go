package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Wormhole Technician
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Cyborg • Scientist
//
//	Reap: Reveal the top card of your deck. If it is a Logos card, play it. Otherwise, archive the top card of your deck.
var WormholeTechnician = card.New(
	"Wormhole Technician",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "144"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Cyborg, card.Traits.Scientist),
	card.WithAbility(
		card.Trigger.Reap, card.Sentences{
			Effects: []card.Effect{
				card.RevealTopOfDeck{Amount: 1},
				card.Conditional{
					Cond: card.ItIs{House: card.House.Self},
					Then: card.PlayRevealedCard{},
					Else: card.ArchiveCard{Zone: card.Deck, Selection: card.Top{}},
				},
			},
		}),
)
