package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Plague Rat
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Beast • Rat
//
//	Elusive.
//	Play: For each Rat trait creature in play, deal 1 damage to each non-Rat trait creature.
var PlagueRat = card.New(
	"Plague Rat",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 308),
	card.WithPower(1),
	card.WithTraits(card.Traits.Beast, card.Traits.Rat),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Amount: 1,
			Per: card.InPlay{
				Player: card.EachPlayer,
				Type:   card.Type.Creature,
				Trait:  card.Traits.Rat,
			},
			Target: card.Target.EachCreature.ExceptTrait(card.Traits.Rat),
		}),
)
