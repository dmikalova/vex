package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Niffle Queen
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  6
//	Traits: Beast • Niffle
//
//	Each other friendly Beast creature gains +1 power.
//	Each other friendly Niffle creature gains +1 power.
var NiffleQueen = set.New(
	"Niffle Queen",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "364"),
	card.InCluster(card.Pulled(troopCallCluster, 0, 0.85)),
	card.WithPower(6),
	card.WithTraits(card.Traits.Beast, card.Traits.Niffle),
	card.WithConstant(card.ConstantAbility{
		Target: card.Target.EachFriendlyCreature.With(card.Filter{
			Trait:  card.Traits.Beast,
			Except: card.Except.Focus,
		}),
		PowerBonus: 1,
	}),
	card.WithConstant(card.ConstantAbility{
		Target: card.Target.EachFriendlyCreature.With(card.Filter{
			Trait:  card.Traits.Niffle,
			Except: card.Except.Focus,
		}),
		PowerBonus: 1,
	}),
)
