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
//	Play: Ready each neighboring creature, exalt each neighboring creature, and for the remainder of the turn, each neighboring creature belongs to house Saurian.
var OratorHissaro = card.New(
	"Orator Hissaro",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "205"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Politician),
	card.WithKeywords(card.Keyword.Deploy),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.Ready{Target: card.Target.EachCreature.Neighboring()},
			card.Exalt{
				Target: card.Target.EachCreature.Neighboring(),
				Amount: 1,
			},
			card.BelongToHouse{
				Target:   card.Target.EachCreature.Neighboring(),
				House:    card.House.Self,
				Duration: card.Duration.EndOfTurn,
			},
		}}),
)
