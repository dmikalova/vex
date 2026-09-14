package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Stilt-Kin
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Goblin
//
//	Skirmish.
//	After a Giant Creature is played adjacent to Stilt-Kin, ready and fight with Stilt-Kin.
var StiltKin = set.New(
	"Stilt-Kin",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "14"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Goblin),
	card.WithKeywords(card.Keyword.Skirmish),
	card.WithAbility(
		card.Trigger.AfterCreaturePlayedAdjacent, card.Conditional{
			Cond: card.ItIsOfTrait{Trait: card.Traits.Giant},
			Then: card.OnChooseCreature{
				Target: card.Target.This,
				Verbs: []card.CreatureVerb{
					card.ReadyVerb{},
					card.FightVerb{},
				},
			},
		}),
)
