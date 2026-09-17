package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Auto-Vac 5150
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Action: Discard a card from your archives -> keys cost +3 Æmber during your opponent's next turn. Otherwise, archive a card from your hand.
var AutoVac5150 = set.New(
	"Auto-Vac 5150",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.MM, "101"),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.Action, card.Then{
			First: card.DiscardCard{
				Player:    card.Controller,
				Zones:     []card.Zone{card.Archives},
				Selection: card.Chosen{Optional: true},
			},
			Result: card.RaiseKeyCost{
				Player:   card.Opponent,
				Amount:   3,
				Duration: card.Duration.OpponentNextTurn,
			},
			Else: card.ArchiveCard{Zone: card.Hand, Selection: card.Chosen{}},
		}),
)
