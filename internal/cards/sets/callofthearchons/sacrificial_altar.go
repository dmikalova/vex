package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Sacrificial Altar
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Location
//
//	Action: Purge a friendly Human trait creature -> play a creature from your discard pile.
var SacrificialAltar = card.New(
	"Sacrificial Altar",
	card.House.Dis,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 78),
	card.WithAemberBonus(1),
	card.WithTraits("Location"),
	card.WithAbility(
		card.Trigger.Action, card.Then{
			First: card.PurgeCreature{
				Target: card.Target.FriendlyCreature.WithTrait("Human"),
			},
			Result: card.PlayFrom{
				From: card.Discard,
				Type: card.Type.Creature,
			},
		}),
)
