package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Imperial Forge
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Forge a key at +8 Æmber current cost, reduced by 1 Æmber for each Æmber on friendly Creatures -> purge Imperial Forge.
var ImperialForge = set.New(
	"Imperial Forge",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, "222"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.ForgeKey{
			Extra:     8,
			ReducedBy: card.AemberOnFriendlyCreatures{},
		}),
)
