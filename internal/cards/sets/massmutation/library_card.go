package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Library Card
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Action: For the remainder of the turn, each time you play another card, draw a card. Purge Library Card.
var LibraryCard = set.New(
	"Library Card",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.MM, "105"),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.Action, card.Sequence{Effects: []card.Effect{
			card.ForRemainderOfTurn{
				On: card.Event.CardPlayed,
				Do: card.Draw{Amount: 1},
			},
			card.PurgeSource{},
		}}),
)
