package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Sacrificial Altar
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Location
//
//	Action: Purge a friendly Human creature -> play a creature from your discard pile.
var SacrificialAltar = set.New(
	"Sacrificial Altar",
	card.House.Dis,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "78"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Location),
	card.WithAbility(
		card.Trigger.Action, card.Then{
			First: card.PurgeCreature{
				Target: card.Target.FriendlyCreature.WithTrait(card.Traits.Human),
			},
			Result: card.PlayFrom{
				From:  card.Discard,
				Types: card.Types.Of(card.Type.Creature),
			},
		}),
)
