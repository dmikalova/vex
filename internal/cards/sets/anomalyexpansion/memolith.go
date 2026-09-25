package anomalyexpansion

import "github.com/dmikalova/vex/internal/card"

// Memolith
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Special
//	Traits: Location
//
//	Action: Choose one:
//	- Put a tactic card from your hand faceup under Memolith
//	- Trigger the play effect of a Tactic grafted onto Memolith.
var Memolith = set.New(
	"Memolith",
	card.House.Brobnar,
	card.Type.Artifact,
	card.Rarity.Special,
	card.Provenance(card.WC, "A04"),
	card.WithTraits(card.Traits.Location),
	card.WithAbility(
		card.Trigger.Action, card.ChooseOne{Options: []card.Effect{
			card.PutUnderFromHand{Type: card.Type.Tactic},
			card.TriggerGraftedPlayEffect{},
		}}),
)
