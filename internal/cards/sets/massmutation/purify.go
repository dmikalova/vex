package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Purify
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Purge a Mutant creature -> discard cards from the top of its controller's deck until you discard a non-Mutant creature or run out of cards -> put the discarded creature into play under its owner's control.
var Purify = set.New(
	"Purify",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "153"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Then{
			First: card.PurgeCreature{
				Target: card.Target.Creature.WithTrait(card.Traits.Mutant),
			},
			Result: card.Then{
				First: card.DiscardUntil{
					Player: card.ItsController,
					Filter: card.Filter{
						Type:        card.Type.Creature,
						ExceptTrait: card.Traits.Mutant,
					},
				},
				Result: card.PutDiscardedIntoPlay{Type: card.Type.Creature},
			},
		}),
)
