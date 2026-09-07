package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Centurion Stenopius
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Armor:  2
//	Traits: Dinosaur • Soldier
//
//	Centurion Stenopius gains +3 power for each Æmber on it.
//	Play/Fight/Reap: You may exalt Centurion Stenopius.
var CenturionStenopius = card.New(
	"Centurion Stenopius",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 199),
	card.WithPower(3),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Soldier),
	card.WithConstant(card.ConstantAbility{
		Target:     card.Target.This,
		PowerBonus: 3,
		Per:        card.AemberOnThis{},
	}),
	card.WithPlayFightReap(card.May{Do: card.Exalt{
		Target: card.Target.This,
		Amount: 1,
	}}),
)
