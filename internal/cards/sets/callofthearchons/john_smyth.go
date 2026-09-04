package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// "John Smyth"
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Agent • Martian
//
//	Elusive.
//	Fight/Reap: Ready a non-Agent trait Mars creature.
var JohnSmyth = card.New(
	"\"John Smyth\"",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 195),
	card.WithPower(2),
	card.WithTraits(card.Traits.Agent, card.Traits.Martian),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithFightOrReap(card.OnChooseCreature{
		Target: card.Target.Creature.OfHouse(card.House.Mars).ExceptTrait(card.Traits.Agent),
		Verbs:  []card.CreatureVerb{card.ReadyVerb{}},
	}),
)
