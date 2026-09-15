package anomalyexpansion

import "github.com/dmikalova/vactrol/internal/card"

// Infomancer
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Special
//	Power:  3
//	Traits: Human • Cyborg
//
//	Elusive.
//	Play: Put a tactic card from your hand faceup under Infomancer.
//	Reap: Trigger the play effect of a Tactic grafted onto Infomancer.
var Infomancer = set.New(
	"Infomancer",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Special,
	card.Provenance(card.WC, "A02"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human, card.Traits.Cyborg),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Play, card.PutUnderFromHand{Type: card.Type.Tactic}),
	card.WithAbility(
		card.Trigger.Reap, card.TriggerGraftedPlayEffect{}),
)
