package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Qyxxlyx Plague Master
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Martian • Scientist
//
//	Fight/Reap: Deal 3 damage to each Human trait creature, ignoring armor.
var QyxxlyxPlagueMaster = card.New(
	"Qyxxlyx Plague Master",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 198),
	card.WithPower(3),
	card.WithTraits("Martian", "Scientist"),
	card.WithFightOrReap(card.DealDamage{
		Amount:      3,
		Target:      card.Target.EachCreature.WithTrait("Human"),
		IgnoreArmor: true,
	}),
)
