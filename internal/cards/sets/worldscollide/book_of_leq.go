package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Book of leQ
//
//	House:  Star Alliance
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Action: Reveal the top card of your deck. If it is a non-Star Alliance card, its house becomes your active house. Otherwise, end your turn.
var BookOfLeQ = set.New(
	"Book of leQ",
	card.House.StarAlliance,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, "325"),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.Action, card.Sequence{
			Effects: []card.Effect{
				card.RevealTopOfDeck{Amount: 1},
				card.Conditional{
					Cond: card.ItIs{House: card.Houses.Except(card.House.Self)},
					Then: card.ChangeActiveHouse{To: card.TheContextualHouse},
					Else: card.EndTurn{},
				},
			},
		}),
)
