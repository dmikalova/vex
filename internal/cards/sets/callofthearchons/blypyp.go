package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Blypyp
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Martian • Scientist
//
//	Reap: The next Mars creature you play this turn enters play ready.
var Blypyp = card.New(
	"Blypyp",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 190),
	card.WithPower(2),
	card.WithTraits(card.Traits.Martian, card.Traits.Scientist),
	card.WithAbility(
		card.Trigger.Reap, card.NextPlayed{
			Of:         card.House.Mars,
			Type:       card.Type.Creature,
			EntersPlay: card.Ready{Target: card.Target.Triggering},
		}),
)
