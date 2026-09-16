package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Amphora Captura
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	When resolving a bonus icon, you may resolve it as a Capture bonus icon instead.
//	Enhance Æmber Æmber Damage Damage Draw Draw.
var AmphoraCaptura = set.New(
	"Amphora Captura",
	card.House.Saurian,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.MM, "215"),
	card.WithEnhance(
		card.Bonus.Aember, card.Bonus.Aember,
		card.Bonus.Damage, card.Bonus.Damage,
		card.Bonus.Draw, card.Bonus.Draw,
	),
	card.WithTraits(card.Traits.Item),
	card.WithBonusInstead(card.BonusInstead{
		May: true,
		As:  card.Bonus.Capture,
	}),
)
