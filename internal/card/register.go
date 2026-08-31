package card

import (
	"sort"

	"github.com/dmikalova/vactrol/internal/cards/provenance"
	"github.com/dmikalova/vactrol/internal/engine"
)

// RegisteredCard is a built card together with the provenance tags it was
// declared with (see Provenance).
type RegisteredCard struct {
	Def        Definition
	Provenance []provenance.Ref
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
	registry = append(registry, RegisteredCard{Def: d, Provenance: b.prov})
	return d
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
