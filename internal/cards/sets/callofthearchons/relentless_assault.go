package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Relentless Assault
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Ready and fight with up to 3 different friendly creatures, one at a time.
var RelentlessAssault = card.New(
	"Relentless Assault",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 13),
	card.WithAbility(
		card.Trigger.Play, card.OneAtATime{
			Times:  card.Fixed(3),
			Target: card.Target.FriendlyCreature,
			Verbs:  []card.CreatureVerb{card.ReadyVerb{}, card.FightVerb{}},
		}),
)
