package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Chuff Ape
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Rare
//	Power:  11
//	Traits: Beast
//
//	Taunt.
//	Chuff Ape enters play stunned.
//	Fight/Reap: You may destroy another friendly creature -> fully heal Chuff Ape.
var ChuffApe = card.New(
	"Chuff Ape",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 191),
	card.WithPower(11),
	card.WithTraits(card.Traits.Beast),
	card.WithKeywords(card.Keyword.Taunt),
	card.WithEntersPlay(card.Stun{Target: card.Target.This}),
	card.WithFightOrReap(card.May{
		Do: card.Then{
			First: card.Destroy{Target: card.Target.OtherFriendlyCreature},
			Result: card.Heal{
				Fully:  true,
				Target: card.Target.This,
			},
		},
	}),
)
