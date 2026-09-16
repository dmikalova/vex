package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// isMutantCreature reports whether a definition is a Mutant creature, for Dark
// Æmber Vault's deck-wide pull.
func isMutantCreature(d card.Definition) bool {
	if d.Type != card.Type.Creature {
		return false
	}
	for _, t := range d.Traits {
		if t == card.Traits.Mutant {
			return true
		}
	}
	return false
}

// Dark Æmber Vault
//
//	House:  None
//	Type:   Artifact
//	Rarity: Special
//	Traits: Location
//
//	Each friendly Mutant creature gains +2 power.
//	After a creature is played, if it is a friendly creature and it is a Mutant creature, draw a card.
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
		card.Trigger.AfterCreaturePlayed, card.Conditional{
			Cond: card.And{Conditions: []card.Condition{
				card.ItIsFriendly{},
				card.ItIsOfTrait{Trait: card.Traits.Mutant},
			}},
			Then: card.Draw{Amount: 1},
		}),
)
