package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Professor Terato
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Mutant • Scientist
//
//	Each Mutant creature gains, "Reap: Draw a card."
var ProfessorTerato = set.New(
	"Professor Terato",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "095"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Mutant, card.Traits.Scientist),
	card.WithConstant(card.ConstantAbility{
		Target: card.Target.EachCreature.WithTrait(card.Traits.Mutant),
		Granted: []card.Ability{{
			Trigger: card.Trigger.Reap,
			Effect:  card.Draw{Amount: 1},
		}},
	}),
)
