package massmutation

import "github.com/dmikalova/vex/internal/card"

// Sirs Colossus
//
//	House:  Sanctum
//	Type:   Gigantic Creature
//	Rarity: Rare
//	Power:  10
//	Armor:  3
//	Traits: Knight • Spirit
//
//	Taunt.
//	Play: Capture all your opponent's Æmber, distributed among any number of friendly creatures.
//	Fight: Move each Æmber on a friendly creature to the common supply.
var SirsColossus = set.Gigantic(
	"Sirs Colossus",
	card.House.Sanctum,
	card.Rarity.Rare,
	card.Provenance(card.MoMu, "152"),
	card.WithPower(10),
	card.WithArmor(3),
	card.WithTraits(card.Traits.Knight, card.Traits.Spirit),
	card.WithKeywords(card.Keyword.Taunt),
	card.WithAbility(
		card.Trigger.Play, card.DistributeCapture{
			All:    true,
			Source: card.Opponent,
		}),
	card.WithAbility(
		card.Trigger.Fight, card.MoveAemberToSupply{
			All:    true,
			Target: card.Target.FriendlyCreature,
		}),
)
