package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Noname
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Demon
//
//	Noname gains +1 power for each purged card.
//	Play/Fight/Reap: Purge a card from a discard pile.
var Noname = card.New(
	"Noname",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "112"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Demon),
	card.WithConstant(card.ConstantAbility{
		Target:     card.Target.This,
		PowerBonus: 1,
		Per:        card.PurgedCards{},
	}),
	card.WithAbility(card.Trigger.Play, card.PurgeCard{Zone: card.Discard}),
	card.WithAbility(card.Trigger.Fight, card.PurgeCard{Zone: card.Discard}),
	card.WithAbility(card.Trigger.Reap, card.PurgeCard{Zone: card.Discard}),
)
