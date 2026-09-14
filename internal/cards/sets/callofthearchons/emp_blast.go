package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// EMP Blast
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Stun each Mars Creature and each Robot Creature, and destroy each Artifact.
var EMPBlast = set.New(
	"EMP Blast",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "163"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.Stun{Target: card.Target.EachCreature.House(card.Houses.Named(card.House.Self))},
			card.Stun{Target: card.Target.EachCreature.WithTrait(card.Traits.Robot)},
			card.Destroy{Target: card.Target.EachArtifact},
		}}),
)
