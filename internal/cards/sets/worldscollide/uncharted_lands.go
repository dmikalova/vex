package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Uncharted Lands
//
//	House:  Star Alliance
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Each Star Alliance creature gains, "Reap: Move 1 Æmber from Uncharted Lands to your pool."
//	Play: Place 6 Æmber from the common supply on Uncharted Lands.
var UnchartedLands = set.New(
	"Uncharted Lands",
	card.House.StarAlliance,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, "342"),
	card.WithTraits(card.Traits.Location),
	card.WithConstant(card.ConstantAbility{
		Target: card.Target.EachCreature.House(card.Houses.Named(card.House.Self)),
		Granted: []card.Ability{{
			Trigger: card.Trigger.Reap,
			Effect: card.MoveAember{
				Amount: 1,
				From:   card.Target.GrantingCard,
				To:     card.Controller,
			},
		}},
	}),
	card.WithAbility(
		card.Trigger.Play, card.PlaceAemberOnThis{Amount: 6}),
)
