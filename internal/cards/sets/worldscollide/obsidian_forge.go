package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// ObsidianForge
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Item
//
//	Action: Sacrifice any number of friendly creatures. Then, you may forge a key at +6A current cost, reduced by 1A for each creature sacrificed this way. If you do, destroy Obsidian Forge.
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
