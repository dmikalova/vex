package massmutation

import (
	"slices"

	"github.com/dmikalova/vactrol/internal/card"
)

// isNiffleCreature reports whether a definition is a Niffle creature, for Niffle
// Kong's deck-wide pull.
func isNiffleCreature(d card.Definition) bool {
	if d.Type != card.Type.Creature {
		return false
	}
	return slices.Contains(d.Traits, card.Traits.Niffle)
}

// Niffle Kong
//
//	House:  Untamed
//	Type:   Gigantic Creature
//	Rarity: Rare
//	Power:  12
//	Armor:  2
//	Traits: Mutant • Niffle
//
//	Play: Search your deck and discard pile for any number of Niffle creatures, reveal them, and put them into your hand. Shuffle your deck.
//	Fight/Reap: You may destroy a friendly Niffle creature -> deal 3 damage to a creature. Steal 1 Æmber. Destroy an enemy artifact.
var NiffleKong = set.Gigantic(
	"Niffle Kong",
	card.House.Untamed,
	card.Rarity.Rare,
	card.Provenance(card.MM, "422"),
	card.PullsMatching("Niffle Kong Niffles", 2, 3, isNiffleCreature),
	card.WithPower(12),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Mutant, card.Traits.Niffle),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{
			Effects: []card.Effect{
				card.Search{
					Sources: []card.Zone{card.Deck, card.Discard},
					Filter: card.Filter{
						Type:  card.Type.Creature,
						Trait: card.Traits.Niffle,
					},
					Any:    true,
					Reveal: true,
					Dest:   card.To.Hand,
				},
				card.Shuffle{},
			},
		}),
	card.WithAbility(
		card.Trigger.FightReap, card.May{Do: card.Then{
			First: card.Destroy{
				Target: card.Target.FriendlyCreature.WithTrait(card.Traits.Niffle),
			},
			Result: card.Sequence{Effects: []card.Effect{
				card.DealDamage{Amount: 3, Target: card.Target.Creature},
				card.StealAember{Amount: 1},
				card.Destroy{Target: card.Target.EnemyArtifact},
			}},
		}}),
)
