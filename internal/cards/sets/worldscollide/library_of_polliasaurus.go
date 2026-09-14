package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Library of Polliasaurus
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Location
//
//	Action: Move 1 Æmber from a friendly Creature to your pool.
var LibraryOfPolliasaurus = set.New(
	"Library of Polliasaurus",
	card.House.Saurian,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "204"),
	card.WithTraits(card.Traits.Location),
	card.WithAbility(
		card.Trigger.Action, card.MoveAember{
			Amount: 1,
			From:   card.Target.FriendlyCreature,
			To:     card.Controller,
		}),
)
