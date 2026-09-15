package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Wild Wormhole
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Play the top card of your deck.
var WildWormhole = set.New(
	"Wild Wormhole",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, "125"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(card.Trigger.Play, card.PlayTopOfDeck{}),
)
