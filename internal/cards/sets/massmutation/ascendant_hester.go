package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Ascendant Hester
//
//	House:  Sanctum
//	Type:   Gigantic Creature
//	Rarity: Special
//	Power:  8
//	Traits: Knight • Spirit
//
//	Each other friendly creature gains +2 armor for each Æmber on it.
//	Play/Fight: Each friendly creature captures 1 Æmber from your opponent.
var AscendantHester = set.Gigantic(
	"Ascendant Hester",
	card.House.Sanctum,
	card.Rarity.Special,
	card.Provenance(card.MoMu, "134"),
	card.WithPower(8),
	card.WithTraits(card.Traits.Knight, card.Traits.Spirit),
	card.WithConstant(card.ConstantAbility{
		Target:     card.Target.EachOtherFriendlyCreature,
		ArmorBonus: 2,
		PerTarget:  card.AemberOnIt,
	}),
	card.WithAbility(
		card.Trigger.PlayFight, card.CaptureAember{
			Amount: 1,
			Target: card.Target.EachFriendlyCreature,
			Source: card.Opponent,
		}),
)
