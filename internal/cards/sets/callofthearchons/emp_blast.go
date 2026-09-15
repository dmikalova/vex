package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// EMP Blast
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Stun each Mars creature and each Robot creature, and destroy each artifact.
var EMPBlast = set.New(
	"EMP Blast",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "163"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.Stun{Target: card.Target.EachCreature.House(card.Houses.Named(card.House.Self))},
			card.Stun{Target: card.Target.EachCreature.WithTrait(card.Traits.Robot)},
			card.Destroy{Target: card.Target.EachArtifact},
		}}),
)
