package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Auto-Encoder
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Common
//	Traits: Item
//
//	After you discard a card from your hand, archive the top card of your deck.
var AutoEncoder = set.New(
	"Auto-Encoder",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Common,
	card.Provenance(card.MM, "066"),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.AfterDiscardFromHand, card.ArchiveCard{
			Zone:      card.Deck,
			Selection: card.Top{},
		}),
)
