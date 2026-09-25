package massmutation

import "github.com/dmikalova/vex/internal/card"

// Maleficorn
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Traits: Mutant
//
//	After you resolve a Damage bonus icon, deal 1 damage to the same creature.
//	Enhance Damage Damage Damage Damage.
var Maleficorn = set.New(
	"Maleficorn",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "040"),
	card.WithEnhance(card.Bonus.Damage, card.Bonus.Damage, card.Bonus.Damage, card.Bonus.Damage),
	card.WithPower(5),
	card.WithTraits(card.Traits.Mutant),
	card.WithAbility(card.Trigger.AfterBonusDamage, card.DealDamage{
		Amount: 1,
		Target: card.Target.TheSameCreature,
	}),
)
