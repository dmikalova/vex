package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// EMP Blast
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Stun each Mars or Robot creature. Destroy each artifact.
var EMPBlast = set.New(
	"EMP Blast",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "163"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.Stun{
				Target: card.Target.EachCreature.With(
					card.Filter{
						House:    card.Houses.Named(card.House.Self),
						Trait:    card.Traits.Robot,
						MatchAny: true,
					},
				),
			},
			card.Destroy{Target: card.Target.EachArtifact},
		}}),
)
