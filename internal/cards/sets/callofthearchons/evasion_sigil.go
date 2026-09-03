package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Evasion Sigil
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Power
//
//	Each creature gains, "Before Fight: Discard the top card of its controller's deck. If it is of the active house, the fight does not occur."
var EvasionSigil = card.New(
	"Evasion Sigil",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 286),
	card.WithAemberBonus(1),
	card.WithTraits("Power"),
	card.WithConstant(card.ConstantAbility{
		Target: card.Target.EachCreature,
		Granted: []card.Ability{{
			Trigger: card.Trigger.BeforeFight,
			Effect: card.Sentences{
				Effects: []card.Effect{
					card.DiscardTopOfDeck{},
					card.Conditional{
						Cond: card.ItIsOfHouse{House: card.TheActiveHouse},
						Then: card.CancelFight{},
					},
				},
			},
		}},
	}),
)
