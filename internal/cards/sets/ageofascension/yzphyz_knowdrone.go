package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Yzphyz Knowdrone
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Armor:  1
//	Traits: Martian • Scientist
//
//	Play: Archive a card from your hand. You may purge a card from your archives to stun a Creature.
var YzphyzKnowdrone = set.New(
	"Yzphyz Knowdrone",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "210"),
	card.WithPower(3),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Martian, card.Traits.Scientist),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{Effects: []card.Effect{
			card.ArchiveCard{Zone: card.Hand, Selection: card.Chosen{}},
			card.PurgeArchivedCardThen{
				Then: card.Stun{Target: card.Target.Creature},
			},
		}}),
)
