// Package card is a small, ergonomic facade over the game engine for *authoring*
// cards. Card-set packages import this package (as `card`) so definitions read
// like card.New(...), card.House.Brobnar, card.DealDamage{...} without referring to the
// engine's package name directly.
//
// Almost everything here is a re-export — type aliases (which are fully
// interchangeable with the engine's types) and constant/function values — of
// identifiers in package engine, shaping the authoring API. The one bit of logic
// is the card registry (see register.go): New builds a card and enrolls it in
// the database, so a set package declares each card as `var X = card.New(...)`
// — one line, no init, no central list.
package card

import (
	"github.com/dmikalova/vactrol/internal/cards/provenance"
	"github.com/dmikalova/vactrol/internal/engine"
)

// The core card-structure types you name directly when authoring. The effect
// nodes are in effects.go; the enum-like categories (house, type, rarity,
// keyword, trigger, player, target, duration) are grouped namespaces in their
// own files (types.go, target.go, duration.go).
type (
	// Definition is a card blueprint.
	Definition = engine.CardDefinition
	// Ability pairs a trigger with an effect.
	Ability = engine.Ability
	// StaticModifier is a continuous stat change or replacement applied by an upgrade.
	StaticModifier = engine.StaticModifier
	// ConstantAbility is a continuous stat change a card applies to creatures in
	// play. Its Target says which creatures it reaches (each creature when unset);
	// e.g. card.Target.EachFriendlyCreature.
	ConstantAbility = engine.ConstantAbility
	// Restrictions are the continuous "cannot" rules a card imposes on its
	// controller while in play (e.g. card.Restrictions{CannotPlay: card.Type.Creature}).
	Restrictions = engine.Restrictions
	// PlayCardLimit caps cards a relative player may play in a turn; use it as
	// card.Restrictions{PlayCardLimit: card.PlayCardLimit{Player: card.Opponent, Amount: 2}}.
	PlayCardLimit = engine.PlayCardLimit
	// Toll is Æmber a card makes its controller's opponent pay to play or use an
	// artifact; use it as card.Restrictions{Toll: card.Toll{Action: card.TollOn.PlayArtifact, Amount: 1}}.
	Toll = engine.Toll
	// PlayPermission is a continuous grant letting a card's controller play cards of
	// a house off-turn; use it as card.WithPlayPermission(card.PlayPermission{House:
	// card.House.Untamed}).
	PlayPermission = engine.PlayPermission
	// TollAction names the action a Toll charges for; the ready-made values live in
	// card.TollOn.
	TollAction = engine.TollAction
	// AttackDamage customizes the damage a creature deals when it fights (Valdr's
	// flank bonus, Ether Spider dealing none); pass it to card.WithAttackDamage.
	AttackDamage = engine.AttackDamage
)

// TollOn groups the actions a Toll can charge for, e.g. card.TollOn.PlayArtifact.
var TollOn = struct {
	PlayArtifact TollAction
	UseArtifact  TollAction
}{
	PlayArtifact: engine.TollPlayArtifact,
	UseArtifact:  engine.TollUseArtifact,
}

// KeyCostChange builds a continuous key-cost change a card imposes while in play,
// e.g. card.KeyCostChange(card.Opponent, 1). The affected player — card.Controller,
// card.Opponent, or card.EachPlayer — is a required argument: there is no default,
// so a key-cost change can never be authored without stating whose keys it
// changes. Pass it to card.WithKeyCost, or to a StaticModifier's KeyCostChange
// field for an upgrade that grants the change.
func KeyCostChange(player Player, amount int) engine.KeyCostChange {
	return engine.NewKeyCostChange(player, amount)
}

// Option configures a card at construction: either a gameplay option (the With*
// helpers below, which wrap the engine's options) or a Provenance tag. New splits
// them — gameplay options build the definition, provenance is recorded in the
// registry (see register.go).
type Option func(*builder)

// builder accumulates the pieces New assembles into a registered card.
type builder struct {
	opts []engine.CardOption
	prov []provenance.Ref
}

func gameplay(o engine.CardOption) Option { return func(b *builder) { b.opts = append(b.opts, o) } }
