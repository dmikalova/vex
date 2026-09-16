package worldscollide

import (
	"github.com/dmikalova/vactrol/internal/card"
	"github.com/dmikalova/vactrol/internal/cards/clusters"
)

// Praefectus Ludo
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Dinosaur • Politician
//
//	Each other friendly creature gains, "Destroyed: Move each Æmber on this creature to the common supply."
var PraefectusLudo = set.New(
	"Praefectus Ludo",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "190"),
	card.InCluster(card.Pulled(clusters.Ludo, 1, 1)),
	card.WithPower(5),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Politician),
	card.WithConstant(card.ConstantAbility{
		Target: card.Target.EachOtherFriendlyCreature,
		Granted: []card.Ability{{
			Trigger: card.Trigger.Destroyed,
			Effect:  card.MoveAemberToSupply{All: true, Target: card.Target.This},
		}},
	}),
)
