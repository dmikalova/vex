package massmutation

import "github.com/dmikalova/vex/internal/card"

// Blast Shielding
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Common
//	Bonus:  Æmber
//
//	This creature gains +2 armor.
//	This creature gains, "After this creature is used, you may attach Blast Shielding to a friendly neighboring creature."
var BlastShielding = set.New(
	"Blast Shielding",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Common,
	card.Provenance(card.MM, "303"),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{
		ArmorBonus: 2,
		Granted: []card.Ability{{
			Trigger: card.Trigger.UsedSelf,
			Effect: card.May{Do: card.AttachSelfTo{
				Target: card.Target.FriendlyCreature.Neighboring(),
			}},
		}},
	}),
)
