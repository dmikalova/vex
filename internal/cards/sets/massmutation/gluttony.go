package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Gluttony
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Special
//	Power:  6
//	Traits: Demon • Sin
//
//	Play: For each friendly Sin creature, exalt Gluttony.
//	Reap: Move all Æmber from each friendly creature to your pool.
var Gluttony = set.New(
	"Gluttony",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Special,
	card.Provenance(card.MM, "057"),
	card.InCluster(sinsCluster),
	card.OneCopyPerDeck(),
	card.WithPower(6),
	card.WithTraits(card.Traits.Demon, card.Traits.Sin),
	card.WithAbility(card.Trigger.Play, card.ForEach{
		Times: card.CardsInPlay{
			Player: card.Controller,
			Type:   card.Type.Creature,
			Trait:  card.Traits.Sin,
		},
		Do: card.Exalt{
			Target: card.Target.This,
			Amount: 1,
		},
	}),
	card.WithAbility(card.Trigger.Reap, card.MoveAember{
		All:  true,
		From: card.Target.EachFriendlyCreature,
		To:   card.Controller,
	}),
)
