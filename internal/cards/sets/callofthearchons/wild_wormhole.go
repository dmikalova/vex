package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Wild Wormhole
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Play the top card of your deck.
var WildWormhole = card.New(
	"Wild Wormhole",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 125),
	card.WithAemberBonus(1),
	card.WithAbility(card.Trigger.Play, card.PlayTopOfDeck{}),
)
