package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Obsidian Forge
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Item
//
//	Action: Destroy any number of friendly Creatures, then forge a key at +6 Æmber current cost, reduced by 1 Æmber for each Creature destroyed this way -> purge Obsidian Forge.
var ObsidianForge = card.New(
	"Obsidian Forge",
	card.House.Dis,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "093"),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.Action, card.SacrificeToForge{
			Target: card.Target.EachFriendlyCreature,
			Extra:  6,
		}),
)
