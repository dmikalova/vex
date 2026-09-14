package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Orator Hissaro
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Dinosaur • Politician
//
//	Deploy.
//	Play: Ready and exalt each neighboring Creature. For the remainder of the turn, those Creatures belong to house Saurian.
var OratorHissaro = set.New(
	"Orator Hissaro",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "205"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Politician),
	card.WithKeywords(card.Keyword.Deploy),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{Effects: []card.Effect{
			card.Sequence{Effects: []card.Effect{
				card.Ready{Target: card.Target.EachCreature.Neighboring()},
				card.Exalt{
					Target: card.Target.EachCreature.Neighboring(),
					Amount: 1,
				},
			}},
			card.BelongToHouse{
				Target:   card.Target.EachCreature.Neighboring(),
				House:    card.House.Self,
				Duration: card.Duration.RemainderOfPlayerTurn,
				Pronoun:  true,
			},
		}}),
)
