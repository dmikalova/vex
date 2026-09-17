package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Bawretchadontius
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Special
//	Power:  14
//	Traits: Beast
//
//	Each friendly creature with Æmber on it gains, "Reap: Deal 4 damage to a creature."
//	Play/Fight/Reap: Exalt a friendly creature and 2 enemy creatures.
var Bawretchadontius = set.Gigantic(
	"Bawretchadontius",
	card.House.Saurian,
	card.Rarity.Special,
	card.Provenance(card.MoMu, "194"),
	card.WithPower(14),
	card.WithTraits(card.Traits.Beast),
	card.WithConstant(card.ConstantAbility{
		Target: card.Target.EachFriendlyCreature.WithAember(),
		Granted: []card.Ability{{
			Trigger: card.Trigger.Reap,
			Effect:  card.DealDamage{Amount: 4, Target: card.Target.Creature},
		}},
	}),
	card.WithAbility(
		card.Trigger.PlayFightReap, card.Sequence{Effects: []card.Effect{
			card.Exalt{Target: card.Target.FriendlyCreature, Amount: 1},
			card.Exalt{
				Target:   card.Target.EnemyCreature,
				Amount:   1,
				Times:    card.Fixed(2),
				Distinct: true,
			},
		}}),
)
