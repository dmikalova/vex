package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Evasion Sigil
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Power
//
//	Each creature gains, "Before Fight: Discard the top card of its controller's deck. If it is of the active house, the fight does not occur."
var EvasionSigil = set.New(
	"Evasion Sigil",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "286"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Power),
	card.WithConstant(card.ConstantAbility{
		Target: card.Target.EachCreature,
		Granted: []card.Ability{{
			Trigger: card.Trigger.BeforeFight,
			Effect: card.Sentences{
				Effects: []card.Effect{
					card.DiscardTop{},
					card.Conditional{
						Cond: card.ItIs{House: card.Houses.Active},
						Then: card.CancelFight{},
					},
				},
			},
		}},
	}),
)
