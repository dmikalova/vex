package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Guji Dinosaur Hunter
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Giant • Hunter
//
//	Elusive.
//	Action: Choose a creature - if it is a Dinosaur creature or it has Æmber on it, deal 6 damage to it. Otherwise, deal 2 damage to it.
var GujiDinosaurHunter = set.New(
	"Guji Dinosaur Hunter",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "38"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Giant, card.Traits.Hunter),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Action, card.ChooseCreatureThen{
			Target: card.Target.Creature,
			Then: card.Conditional{
				Cond: card.Or{Conditions: []card.Condition{
					card.ItIsOfTrait{Trait: card.Traits.Dinosaur},
					card.ItHasAember{},
				}},
				Then: card.DealDamage{Amount: 6, Target: card.Target.Triggering},
				Else: card.DealDamage{Amount: 2, Target: card.Target.Triggering},
			},
		}),
)
