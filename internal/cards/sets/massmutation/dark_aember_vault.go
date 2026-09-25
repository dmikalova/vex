package massmutation

import (
	"slices"

	"github.com/dmikalova/vex/internal/card"
)

// isMutantCreature reports whether a definition is a Mutant creature, for Dark
// Æmber Vault's deck-wide pull.
func isMutantCreature(d card.Definition) bool {
	if d.Type != card.Type.Creature {
		return false
	}
	return slices.Contains(d.Traits, card.Traits.Mutant)
}

// Dark Æmber Vault
//
//	House:  None
//	Type:   Artifact
//	Rarity: Special
//	Traits: Location
//
//	Each friendly Mutant creature gains +2 power.
//	After you play a Mutant creature, draw a card.
var DarkAemberVault = set.New(
	"Dark Æmber Vault",
	card.House.None,
	card.Type.Artifact,
	card.Rarity.Special,
	card.Provenance(card.MM, "001"),
	card.Houseless(),
	card.PullsMatching("Dark Æmber Vault Mutants", 4, 6, isMutantCreature),
	card.WithTraits(card.Traits.Location),
	card.WithConstant(card.ConstantAbility{
		PowerBonus: 2,
		Target:     card.Target.EachFriendlyCreature.WithTrait(card.Traits.Mutant),
	}),
	card.WithAbility(
		card.Trigger.AfterCardPlayed, card.Conditional{
			Cond: card.ItIsOfTrait{Trait: card.Traits.Mutant},
			Then: card.Draw{Amount: 1},
		}),
)
