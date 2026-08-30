package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Round Table
//
//	House:  Sanctum
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Location
//
//	Each friendly Knight trait creature gains +1 power and taunt.
var RoundTable = card.New(
	"Round Table",
	card.House.Sanctum,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 235),
	card.WithAemberBonus(1),
	card.WithTraits("Location"),
	card.WithConstantAbility(card.ConstantAbility{
		PowerBonus: 1,
		Keywords:   card.Keywords(card.Keyword.Taunt),
		Target:     card.Target.EachFriendlyCreature.WithTrait("Knight"),
	}),
)
