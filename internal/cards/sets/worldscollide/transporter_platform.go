package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Transporter Platform
//
//	House:  Star Alliance
//	Type:   Artifact
//	Rarity: Uncommon
//	Bonus:  Æmber
//	Traits: Location
//
//	Action: Put a friendly creature and each upgrade attached to it into its owner's hand.
var TransporterPlatform = set.New(
	"Transporter Platform",
	card.House.StarAlliance,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "321"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Location),
	card.WithAbility(
		card.Trigger.Action, card.PutFromPlay{
			Target:       card.Target.FriendlyCreature,
			Destination:  card.To.Hand,
			WithUpgrades: true,
		}),
)
