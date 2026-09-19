package card

import (
	"sort"

	"github.com/dmikalova/vactrol/internal/cards/provenance"
	"github.com/dmikalova/vactrol/internal/deckgen"
	"github.com/dmikalova/vactrol/internal/engine"
)

// RegisteredCard is a built card together with the provenance tags and
// deck-generation metadata it was declared with (see Provenance, Template).
type RegisteredCard struct {
	Def          Definition
	Provenance   []provenance.Ref
	Set          provenance.SourceSet
	Profile      deckgen.GenerationProfile
	Materializer Materializer
}

// registry holds every card built with New. A set package declares each card as
// a package-level var, `var X = card.New(...)`, so simply importing the set
// enrolls its cards at package-initialization time — there is no per-card init
// and no central list to maintain (see package cards).
var registry []RegisteredCard

// New builds a card Definition and enrolls it in the global database, returning
// it so it reads as the initializer of a package-level var:
//
//	var Anger = card.New("Anger", card.House.Brobnar, card.Type.Action, card.Rarity.Common, ...)
//
// For a throwaway card that should not join the database (tests, scripted demos)
// use the engine's engine.NewCard directly.
func New(
	name string,
	house engine.House,
	ct engine.CardType,
	rarity engine.Rarity,
	opts ...Option,
) Definition {
	var b builder
	for _, o := range opts {
		o(&b)
	}
	d := engine.NewCard(name, house, ct, rarity, b.opts...)
	registry = append(registry, RegisteredCard{
		Def:          d,
		Provenance:   b.prov,
		Set:          nativeSet(b),
		Profile:      b.profile,
		Materializer: b.materializer,
	})
	return d
}

// Gigantic builds a gigantic creature — a creature printed as two cards, a base
// half and an art half, that share a name and are played as one creature
// (ADR 0042). It is authored as a single call, like New, and the options describe
// the whole creature. Registration splits it: the base half carries the power,
// armor, traits, keywords, and abilities and is the card enrolled in the database
// (with the provenance and generation options); the art half is synthetic,
// carries only the creature's bonus icons, has no provenance, and rides on the
// base's generation profile so deck generation can place it alongside the base.
// The two halves share a name — which the database's unique-name rule forbids for
// two registered cards — so only the base is registered. Gigantic returns the
// base half, so a set declares it as `var X = card.Gigantic(...)`.
func Gigantic(
	name string,
	house engine.House,
	rarity engine.Rarity,
	opts ...Option,
) Definition {
	var b builder
	for _, o := range opts {
		o(&b)
	}
	base := engine.NewCard(
		name, house, engine.Creature, rarity,
		append(b.opts, engine.WithGiganticRole(engine.GiganticBase))...,
	)
	// The art half carries only the creature's bonus icons; move them off the base.
	art := engine.NewCard(
		name, house, engine.Creature, rarity,
		engine.WithGiganticRole(engine.GiganticArt),
		engine.WithBonus(base.Bonuses...),
	)
	base.Bonuses = nil
	prof := b.profile
	prof.GiganticArt = &art
	registry = append(registry, RegisteredCard{
		Def:          base,
		Provenance:   b.prov,
		Set:          nativeSet(b),
		Profile:      prof,
		Materializer: b.materializer,
	})
	return base
}

// GiganticArt returns the synthetic art half of a gigantic creature whose base
// half is base, for tests and tools that need both halves in hand (deck
// generation places the art half automatically). It panics if base is not a
// registered gigantic base half.
func GiganticArt(base Definition) Definition {
	for _, e := range registry {
		if e.Def.Name == base.Name && e.Def.GiganticRole == engine.GiganticBase &&
			e.Profile.GiganticArt != nil {
			return *e.Profile.GiganticArt
		}
	}
	panic("card: GiganticArt called on a non-gigantic card " + base.Name)
}

// nativeSet returns a card's home set for deck generation: the set declared with
// InSet, which a set package's registrar (set.New) always stamps. Deck generation
// groups a card by this set alone and never infers it from provenance (ADR 0003);
// a card registered through the bare New without InSet belongs to no pool (the
// zero set) and is filtered out downstream.
func nativeSet(b builder) provenance.SourceSet {
	return b.set
}

// Build assembles a card Definition WITHOUT registering it, for a template's
// Materializer to produce concrete variants at deck-generation time. It applies
// only the gameplay options; provenance and generation options are ignored.
func Build(
	name string,
	house engine.House,
	ct engine.CardType,
	rarity engine.Rarity,
	opts ...Option,
) Definition {
	var b builder
	for _, o := range opts {
		o(&b)
	}
	return engine.NewCard(name, house, ct, rarity, b.opts...)
}

// Registered returns a copy of every registered card definition, sorted by house
// then name so the database is deterministic regardless of package initialization
// order. The cards aggregator (package cards) imports the set packages that
// populate the registry, so callers should reach it through cards.All().
func Registered() []Definition {
	out := make([]Definition, len(registry))
	for i, e := range registry {
		out[i] = e.Def
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].House != out[j].House {
			return out[i].House < out[j].House
		}
		return out[i].Name < out[j].Name
	})
	return out
}

// Cards returns every registered card with its provenance tags, sorted by house
// then name. Used for coverage reporting against the source catalogs.
func Cards() []RegisteredCard {
	out := make([]RegisteredCard, len(registry))
	copy(out, registry)
	sort.Slice(out, func(i, j int) bool {
		if out[i].Def.House != out[j].Def.House {
			return out[i].Def.House < out[j].Def.House
		}
		return out[i].Def.Name < out[j].Def.Name
	})
	return out
}
