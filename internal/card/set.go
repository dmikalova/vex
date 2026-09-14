package card

import (
	"github.com/dmikalova/vactrol/internal/cards/provenance"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Set is a set package's card registrar. Each set's 0set.go declares one — with
// NewSet or ReservoirSet — and every card in that package registers through its
// New method, so the card's home set is prefilled by the set itself and never
// inferred from provenance. Deck generation groups a card by this declared set
// alone (ADR 0003): provenance stays a pure coverage tag.
type Set struct {
	src       provenance.SourceSet
	reservoir bool
}

// NewSet declares a set's registrar. A card registered through its New method
// belongs to src's deck-generation pool.
func NewSet(src provenance.SourceSet) *Set { return &Set{src: src} }

// ReservoirSet declares a reservoir set's registrar (ADR 0036): every card
// registered through its New method is marked undraftable, so the set builds no
// draw pool of its own and its cards reach a deck only through a cross-set
// mechanism such as a cluster. The Anomaly Expansion Shards are its members.
func ReservoirSet(src provenance.SourceSet) *Set { return &Set{src: src, reservoir: true} }

// New registers a card as a member of the set, prefilling its home set (and, for a
// reservoir set, marking it undraftable) before delegating to the package-level
// New. Author a card as `var X = set.New("X", card.House.Y, card.Type.Z,
// card.Rarity.W, …)`.
func (s *Set) New(
	name string,
	house engine.House,
	ct engine.CardType,
	rarity engine.Rarity,
	opts ...Option,
) Definition {
	opts = append(opts, InSet(s.src))
	if s.reservoir {
		opts = append(opts, reservoir())
	}
	return New(name, house, ct, rarity, opts...)
}

// Reprint records that Name — implemented in the set that introduced it — is also
// printed in this set at collector number number, joining this set's pool as a
// full member (ADR 0021). It is the set-scoped form of the package-level Reprint.
func (s *Set) Reprint(number string, name string) { Reprint(s.src, number, name) }

// reservoir marks a card undraftable, so its set builds no draw pool of its own.
// The card can still guest into a deck through a cross-set mechanism (a cluster,
// or the legacy pool when it is housed and non-Connected). It is set by a
// ReservoirSet's New, not authored directly.
func reservoir() Option { return func(b *builder) { b.profile.Reservoir = true } }
