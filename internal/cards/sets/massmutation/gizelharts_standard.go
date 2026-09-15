package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Gizelhart's Standard
//
//	House:  Sanctum
//	Type:   Artifact
//	Rarity: Uncommon
//	Bonus:  Æmber
//	Traits: Item
//
//	Each friendly creature with Æmber on it gains +1 armor.
//	Play: Exalt a friendly creature.
var GizelhartsStandard = set.New(
	"Gizelhart's Standard",
	card.House.Sanctum,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "150"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Item),
	card.WithConstant(card.ConstantAbility{
		ArmorBonus: 1,
		Target:     card.Target.EachFriendlyCreature.WithAember(),
	}),
	card.WithAbility(
		card.Trigger.Play, card.Exalt{
			Target: card.Target.FriendlyCreature,
			Amount: 1,
		}),
)
