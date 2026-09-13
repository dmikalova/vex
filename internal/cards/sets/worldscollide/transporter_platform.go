package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Transporter Platform
//
//	House:  Star Alliance
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Location
//
//	Action: Put a friendly Creature and each Upgrade attached to it into its owner's hand.
var TransporterPlatform = card.New(
	"Transporter Platform",
	card.House.StarAlliance,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "321"),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Location),
	card.WithAbility(
		card.Trigger.Action, card.PutFromPlay{
			Target:       card.Target.FriendlyCreature,
			Destination:  card.To.Hand,
			WithUpgrades: true,
		}),
)
