package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Gargantodon
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  16
//	Traits: Beast
//
//	Gargantodon deals 4 Damage when fighting.
//	Each Æmber that would be stolen is captured by a Creature controlled by the active player instead.
//	Gargantodon enters play stunned.
var Gargantodon = card.New(
	"Gargantodon",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "203"),
	card.WithPower(16),
	card.WithTraits(card.Traits.Beast),
	card.WithEntersPlay(card.Stun{Target: card.Target.This}),
	card.WithAttackDamage(card.AttackDamage{
		Amount: 4,
		Fixed:  true,
	}),
	card.WithReplaces(card.Instead{
		Of:   card.Event.AemberStolen,
		With: card.Capture,
	}),
)
