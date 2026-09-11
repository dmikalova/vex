package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Zyzzix the Many
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Martian • Soldier
//
//	Fight/Reap: You may reveal a creature from your hand and archive it -> give Zyzzix the Many three +1 power counters.
var ZyzzixTheMany = card.New(
	"Zyzzix the Many",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, "207"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Martian, card.Traits.Soldier),
	card.WithAbility(card.Trigger.FightReap, card.May{
		Do: card.Then{
			First: card.ArchiveCard{
				Zone:      card.Hand,
				Selection: card.Chosen{Type: card.Type.Creature, Optional: true},
				Revealed:  true,
			},
			Result: card.AddPowerCounter{
				Target: card.Target.This,
				Amount: 3,
			},
		},
	}),
)
