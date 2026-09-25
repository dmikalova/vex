package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Round Table
//
//	House:  Sanctum
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Location
//
//	Each friendly Knight creature gains +1 power and taunt.
var RoundTable = set.New(
	"Round Table",
	card.House.Sanctum,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "235"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Location),
	card.WithConstant(card.ConstantAbility{
		PowerBonus: 1,
		Keywords:   card.Keywords(card.Keyword.Taunt),
		Target:     card.Target.EachFriendlyCreature.WithTrait(card.Traits.Knight),
	}),
)
