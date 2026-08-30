package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// EMP Blast
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Stun each Mars creature and each Robot trait creature, and destroy each artifact.
var EMPBlast = card.New(
	"EMP Blast",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 163),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.Stun{Target: card.Target.EachCreature.OfHouse(card.House.Mars)},
			card.Stun{Target: card.Target.EachCreature.WithTrait("Robot")},
			card.Destroy{Target: card.Target.EachArtifact},
		}}),
)
