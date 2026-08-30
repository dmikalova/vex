package cardtest

import (
	"fmt"
	"sync/atomic"

	"github.com/dmikalova/vactrol/internal/engine"
)

// DefaultHouse is the house given to a vanilla card when none is specified, so a
// scenario that only needs "a body" can write ct.Creature() with no ceremony.
var DefaultHouse = engine.Brobnar

// vanillaCount names each generated vanilla card uniquely ("Vanilla Creature 1",
// ...). Unique names matter: the harness resolves a card definition to its single
// in-play/hand instance by name, so two anonymous vanillas must not collide.
var vanillaCount atomic.Uint64

// spec accumulates the options for a vanilla card before it is built into an
// engine.CardDefinition. It is populated by Option values (OfHouse, Power, ...).
type spec struct {
	house       engine.House
	power       int
	armor       int
	traits      []engine.Trait
	keywords    []engine.Keyword
	aemberBonus int
	static      engine.StaticModifier
}

// Option configures a vanilla card built by Creature, Artifact, Tactic, or
// Upgrade. Options are small, composable, and read left to right, e.g.
// ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(6)).
type Option func(*spec)

// OfHouse sets a vanilla card's house (defaults to DefaultHouse).
func OfHouse(h engine.House) Option { return func(s *spec) { s.house = h } }

// Power sets a vanilla creature's power.
func Power(p int) Option { return func(s *spec) { s.power = p } }

// Armor sets a vanilla creature's armor.
func Armor(a int) Option { return func(s *spec) { s.armor = a } }

// Traits sets a vanilla card's traits.
func Traits(traits ...engine.Trait) Option {
	return func(s *spec) { s.traits = append(s.traits, traits...) }
}

// Keywords gives a vanilla card the given keywords (Skirmish, Poison, ...).
func Keywords(keywords ...engine.Keyword) Option {
	return func(s *spec) { s.keywords = append(s.keywords, keywords...) }
}

// AemberBonus sets a vanilla card's Æmber pips.
func AemberBonus(n int) Option { return func(s *spec) { s.aemberBonus = n } }

// PowerBonus sets the power a vanilla Upgrade grants its host creature.
func PowerBonus(n int) Option { return func(s *spec) { s.static.PowerBonus = n } }

// ArmorBonus sets the armor a vanilla Upgrade grants its host creature.
func ArmorBonus(n int) Option { return func(s *spec) { s.static.ArmorBonus = n } }

// build turns the accumulated spec into a definition of the given type, applying
// a sensible default power for creatures and a unique generated name.
func build(kind string, ct engine.CardType, defaultPower int, opts []Option) engine.CardDefinition {
	s := spec{house: DefaultHouse, power: defaultPower}
	for _, o := range opts {
		o(&s)
	}
	name := fmt.Sprintf("Vanilla %s %d", kind, vanillaCount.Add(1))
	cardOpts := []engine.CardOption{
		engine.WithPower(s.power),
		engine.WithArmor(s.armor),
		engine.WithAemberBonus(s.aemberBonus),
		engine.WithStatic(s.static),
	}
	if len(s.traits) > 0 {
		cardOpts = append(cardOpts, engine.WithTraits(s.traits...))
	}
	if len(s.keywords) > 0 {
		cardOpts = append(cardOpts, engine.WithKeywords(s.keywords...))
	}
	return engine.NewCard(name, s.house, ct, engine.Common, cardOpts...)
}

// Creature builds a plain vanilla creature (no abilities). Without options it is
// a DefaultHouse creature of power 4 — a "body with no baggage" for exercising a
// mechanic in isolation.
func Creature(opts ...Option) engine.CardDefinition {
	return build("Creature", engine.Creature, 4, opts)
}

// Artifact builds a plain vanilla artifact (no abilities).
func Artifact(opts ...Option) engine.CardDefinition {
	return build("Artifact", engine.Artifact, 0, opts)
}

// Tactic builds a plain vanilla tactic (no abilities). It does nothing when
// played beyond leaving the hand, useful as filler or an Æmber-pip source.
func Tactic(opts ...Option) engine.CardDefinition {
	return build("Tactic", engine.Tactic, 0, opts)
}

// Upgrade builds a plain vanilla upgrade. Use PowerBonus/ArmorBonus to give it a
// static effect on its host.
func Upgrade(opts ...Option) engine.CardDefinition {
	return build("Upgrade", engine.Upgrade, 0, opts)
}
