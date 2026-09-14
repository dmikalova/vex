package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Captain Val Jericho
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Armor:  1
//	Traits: Human • Leader
//
//	During your turn, if Captain Val Jericho is in the center of your battleline, you may play one card that is not of the active house.
var CaptainValJericho = set.New(
	"Captain Val Jericho",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "326"),
	card.WithPower(5),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Human, card.Traits.Leader),
	card.WithPlayPermission(card.PlayPermission{
		NonActive: true,
		Condition: card.SourceInCenterOfBattleline{},
		Amount:    1,
	}),
)
