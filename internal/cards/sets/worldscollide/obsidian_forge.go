package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Obsidian Forge
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Uncommon
//	Bonus:  Æmber
//	Traits: Item
//
//	Action: Destroy any number of friendly creatures, and forge a key at +6 Æmber current cost, reduced by 1 Æmber for each creature destroyed this way -> purge Obsidian Forge.
var ObsidianForge = set.New(
	"Obsidian Forge",
	card.House.Dis,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "093"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.Action, card.Sequence{Effects: []card.Effect{
			card.DestroyChosen{Target: card.Target.EachFriendlyCreature},
			card.ForgeKey{
				Extra:     6,
				ReducedBy: card.CreaturesDestroyed{},
			},
		}}),
)
