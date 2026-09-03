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
//	Fight/Reap: You may reveal a creature from your hand and archive it -> give Zyzzix the Many 3 +1 power counters.
var ZyzzixTheMany = card.New(
	"Zyzzix the Many",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 207),
	card.WithPower(3),
	card.WithTraits("Martian", "Soldier"),
	card.WithFightOrReap(card.May{
		Do: card.Then{
			First: card.ArchiveFromHand{
				Count:    1,
				Type:     card.Type.Creature,
				Revealed: true,
			},
			Result: card.AddPowerCounter{
				Target: card.Target.This,
				Amount: 3,
			},
		},
	}),
)
