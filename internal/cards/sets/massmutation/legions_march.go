package massmutation

import "github.com/dmikalova/vex/internal/card"

// Legion's March
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: For the remainder of the turn, after you use a Dinosaur creature, deal 1 damage to each non-Dinosaur creature.
var LegionsMarch = set.New(
	"Legion's March",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.MM, "224"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(card.Trigger.Play, card.DamageOthersAfterUsingTrait{
		Trait:  card.Traits.Dinosaur,
		Amount: 1,
	}),
)
