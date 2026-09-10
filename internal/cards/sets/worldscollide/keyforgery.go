package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Keyforgery
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Item
//
//	When your opponent would forge a key, they name a house. Reveal a random card from your hand. If that card is not of the named house, destroy Keyforgery and they do not forge that key.
var Keyforgery = card.New(
	"Keyforgery",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, "271"),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Item),
	card.WithGuardsOpponentForge(),
)
