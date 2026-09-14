package engine

// Selection is how a zone-movement effect picks the cards it acts on from a
// zone: the controller chooses one (Chosen), one is uniformly random (Random),
// every matching card is taken (Each), or a named card is pinned (Named). Each mode filters and picks its own
// cards and renders its own object phrase ("a card" / "a random card" / "each
// creature"), so a movement verb varies along this one axis instead of spawning
// a node per (verb x mode). See ADR 0031.
type Selection interface {
	// pick returns the cards to act on, given every card in the source zone. It
	// may prompt the controller (Chosen) or draw on the RNG (Random).
	pick(ctx *EffectContext, cands []LocalID) []LocalID
	// candidates returns the cards this Selection could act on — the filtered set,
	// before any choice or draw. A verb uses it to tell whether a zone holds
	// anything to act on (PurgeCard picks a discard pile only among those that do).
	candidates(ctx *EffectContext, cands []LocalID) []LocalID
	// noun renders the bare kind of card the verb acts on, without article or
	// count, e.g. "card" / "random card" / "non-Mars creature". A count-bearing
	// verb pluralizes it; object decorates it for the single case.
	noun() string
	// object renders the noun phrase the verb acts on, e.g. "a Sanctum card".
	object() string
}

// declinableSelection is a Selection the controller may pass without acting — an
// Optional Chosen. A Selection that does not implement it can never be
// declined, so its verb never reads "you may".
type declinableSelection interface {
	declinable() bool
}

// selectionDeclinable reports whether a Selection can be passed.
func selectionDeclinable(s Selection) bool {
	d, ok := s.(declinableSelection)
	return ok && d.declinable()
}

// positionalSelection is a Selection that picks by position in an ordered zone —
// the top or bottom card — rather than by identity or filter. edge names which
// end ("top" / "bottom"). Only ordered zones (Deck, Discard) have a top and
// bottom, so a node pairing a positional selection with an unordered zone fails
// validation, and the node hands cands top-first (see topFirst).
type positionalSelection interface {
	positional() bool
	edge() string
}

// selectionPositional reports whether a Selection picks by position.
func selectionPositional(s Selection) bool {
	p, ok := s.(positionalSelection)
	return ok && p.positional()
}

// positionalObject renders a positional selection's count-bearing noun phrase,
// e.g. "the top card" or "the top 2 cards".
func positionalObject(s Selection, n int) string {
	if n == 1 {
		return s.object()
	}
	edge := s.(positionalSelection).edge()
	return "the " + edge + " " + countNoun(n, s.noun())
}

// topFirst returns a zone's cards ordered top-to-bottom, so a positional
// selection can read the top at index 0. A deck is stored top-first already; a
// discard pile is stored bottom-first (its top is the most recently discarded,
// at the end), so it is reversed.
func topFirst(z Zone, cards []LocalID) []LocalID {
	if z != Discard {
		return cards
	}
	out := make([]LocalID, len(cards))
	for i, id := range cards {
		out[len(cards)-1-i] = id
	}
	return out
}

// ownerActsSelection is a Selection where the hand's owner is the one who
// discards, so a discard from an opponent's hand reads "your opponent discards …"
// rather than the controller-directed "discard … from your opponent's hand". A
// Random pick from a hidden hand is attributed to its owner this way.
type ownerActsSelection interface {
	ownerActs() bool
}

// selectionOwnerActs reports whether the hand's owner performs the discard.
func selectionOwnerActs(s Selection) bool {
	o, ok := s.(ownerActsSelection)
	return ok && o.ownerActs()
}

// filterIDs keeps the ids the predicate admits, preserving order.
func filterIDs(ids []LocalID, keep func(LocalID) bool) []LocalID {
	var out []LocalID
	for _, id := range ids {
		if keep(id) {
			out = append(out, id)
		}
	}
	return out
}

// Chosen has the controller pick one card, optionally restricted to a house or a
// type. A Chosen is mandatory by default: the controller must pick when a card
// matches (Greater Oxtet). Optional makes it a "you may" — the controller can
// always decline, and an empty candidate set picks nothing. An Optional Chosen is
// also what drives an "up to N" effect: a node archiving or purging Amount cards
// with an Optional Selection lets the controller stop early, which reads as "up to
// N" (Mobius Scroll, Creeping Oblivion). This means Optional is the single knob
// for "may do fewer" — a card that must do either none or exactly N (all-or-
// nothing) cannot be expressed this way; no card currently needs that.
type Chosen struct {
	// House restricts the choice to cards the matcher admits; the zero value
	// (any house) allows any card (Information Officer Gray reveals a non-Star
	// Alliance card).
	House HouseMatcher
	// Type restricts the choice to cards of this type; the zero value allows any.
	Type CardType
	// Trait restricts the choice to cards carrying this trait; the zero value
	// allows any (Horseman of Death recovers a Horseman creature).
	Trait Trait
	// Name restricts the choice to cards of this exact name; the zero value allows
	// any (Igon the Green recovers an Igon the Terrible).
	Name string
	// Or lists alternative identity filters: a card also qualifies if it satisfies
	// any of them (Chief Engineer Walls recovers an upgrade or a Robot card).
	Or []CardFilter
	// Optional makes the pick a "you may" the controller can decline; the default is
	// a mandatory pick that forces the choice when a card matches.
	Optional bool
}

// filter is the identity predicate the choice narrows by, conjoining Type, Trait,
// and Name and admitting any Or alternative.
func (s Chosen) filter() CardFilter {
	return CardFilter{Type: s.Type, Trait: s.Trait, Name: s.Name, Or: s.Or}
}

// noun renders the bare kind of card chosen, qualified by the identity filter and
// by house.
func (s Chosen) noun() string {
	return s.House.qualify(s.filter().noun())
}

// object renders the single card chosen, e.g. "a Sanctum creature".
func (s Chosen) object() string { return indefinite(s.noun()) }

// declinable reports that an Optional Chosen may be passed.
func (s Chosen) declinable() bool { return s.Optional }

// candidates keeps the cards the house and identity filters admit.
func (s Chosen) candidates(ctx *EffectContext, cands []LocalID) []LocalID {
	return filterIDs(cands, func(id LocalID) bool {
		return s.House.matches(ctx, id) && s.filter().admits(ctx.Resolver, id)
	})
}

// pick offers the matching cards as a card choice — declinable only when Optional —
// and returns the single chosen card, or none when the controller declines or no
// card matches.
func (s Chosen) pick(ctx *EffectContext, cands []LocalID) []LocalID {
	matching := s.candidates(ctx, cands)
	choose := ctx.ChooseCardOptional
	if !s.Optional {
		if len(matching) == 0 {
			return nil
		}
		choose = ctx.ChooseCard
	}
	id, ok := choose("Choose a card", matching)
	if !ok {
		return nil
	}
	return []LocalID{id}
}

// Random takes one uniformly random card, so the acting player does not choose
// which card leaves (Impspector).
type Random struct{}

// noun renders the bare kind of card taken at random.
func (Random) noun() string { return "random card" }

// object renders the random card the verb acts on.
func (Random) object() string { return "a random card" }

// candidates returns every card — a random pick applies no filter.
func (Random) candidates(
	_ *EffectContext,
	cands []LocalID,
) []LocalID {
	return cands
}

// ownerActs reports that a random pick from a hidden hand is attributed to the
// hand's owner, so an opponent's random discard reads "your opponent discards …".
func (Random) ownerActs() bool { return true }

// pick draws one uniformly random card from the candidates, or none when empty.
func (s Random) pick(ctx *EffectContext, cands []LocalID) []LocalID {
	id, ok := ctx.ChooseRandom(s.candidates(ctx, cands))
	if !ok {
		return nil
	}
	return []LocalID{id}
}

// Each takes every card the filters admit, with no choice: the filters decide
// (Martians Make Bad Allies purges each non-Mars creature). A following effect
// can scale with the tally the verb records.
type Each struct {
	// House restricts to cards the matcher admits; the zero value (any house) admits
	// any card (Soldiers to Flowers purges each Untamed creature, Martians Make Bad
	// Allies each non-Mars creature, Deep Probe each creature of the chosen house).
	House HouseMatcher
	// Type restricts to cards of this type; the zero value admits any.
	Type CardType
	// Trait restricts to cards carrying this trait; the zero value admits any
	// (Troop Call recovers each Niffle creature).
	Trait Trait
	// Name restricts to cards of this exact name; the zero value admits any
	// (Ortannu the Chained recovers each Ortannu's Binding).
	Name string
	// Or lists alternative identity filters: a card also qualifies if it satisfies
	// any of them.
	Or []CardFilter
}

// filter is the identity predicate the take narrows by, conjoining Type, Trait,
// and Name and admitting any Or alternative.
func (s Each) filter() CardFilter {
	return CardFilter{Type: s.Type, Trait: s.Trait, Name: s.Name, Or: s.Or}
}

// noun renders the bare kind of card taken, e.g. "non-Mars creature" or "creature
// of the chosen house".
func (s Each) noun() string {
	return s.House.qualify(s.filter().noun())
}

// object renders the kind of card taken, e.g. "each non-Mars creature".
func (s Each) object() string { return "each " + s.noun() }

// candidates returns every card the filters admit — Each takes all of them.
func (s Each) candidates(ctx *EffectContext, cands []LocalID) []LocalID {
	return filterIDs(cands, func(id LocalID) bool {
		return s.House.matches(ctx, id) && s.filter().admits(ctx.Resolver, id)
	})
}

// pick returns every candidate the filters admit.
func (s Each) pick(ctx *EffectContext, cands []LocalID) []LocalID {
	return s.candidates(ctx, cands)
}

// Named pins the pick to the first card of a given name, taken without a choice —
// Hyde archives Velum from the discard pile, Velum archives Hyde. It renders as
// the bare card name.
type Named struct {
	// Name is the exact card name to pick; the first matching card is taken.
	Name string
}

// noun renders the pinned card's name.
func (s Named) noun() string { return s.Name }

// object renders the pinned card's name — a proper name takes no article.
func (s Named) object() string { return s.Name }

// candidates keeps the cards whose name matches.
func (s Named) candidates(ctx *EffectContext, cands []LocalID) []LocalID {
	return filterIDs(cands, func(id LocalID) bool {
		return ctx.Resolver.Name(id) == s.Name
	})
}

// pick takes the first card of the name, or none when none is present.
func (s Named) pick(ctx *EffectContext, cands []LocalID) []LocalID {
	if matching := s.candidates(ctx, cands); len(matching) > 0 {
		return matching[:1]
	}
	return nil
}

// Top pins the pick to the top card of an ordered zone (Deck or Discard), taken
// without a choice. A node's Amount loops it for the top N, each pass taking the
// new top after the last moved — ArchiveCard{Zone: Deck, Selection: Top{}}
// archives the top card of the deck (Random Access Archives). The node hands
// cands top-first (see topFirst), so the top is index 0.
type Top struct{}

// noun renders the bare kind of card taken.
func (Top) noun() string { return "card" }

// object renders the positioned card the verb acts on.
func (Top) object() string { return "the top card" }

// positional reports that Top picks by position, not by identity or filter.
func (Top) positional() bool { return true }

// edge names the end of the zone Top picks from.
func (Top) edge() string { return "top" }

// candidates returns the top card, or none when the zone is empty.
func (Top) candidates(_ *EffectContext, cands []LocalID) []LocalID {
	if len(cands) == 0 {
		return nil
	}
	return cands[:1]
}

// pick takes the top card.
func (s Top) pick(ctx *EffectContext, cands []LocalID) []LocalID {
	return s.candidates(ctx, cands)
}

// Bottom pins the pick to the bottom card of an ordered zone (Deck or Discard),
// taken without a choice — the positional sibling of Top. The node hands cands
// top-first (see topFirst), so the bottom is the last card.
type Bottom struct{}

// noun renders the bare kind of card taken.
func (Bottom) noun() string { return "card" }

// object renders the positioned card the verb acts on.
func (Bottom) object() string { return "the bottom card" }

// positional reports that Bottom picks by position, not by identity or filter.
func (Bottom) positional() bool { return true }

// edge names the end of the zone Bottom picks from.
func (Bottom) edge() string { return "bottom" }

// candidates returns the bottom card, or none when the zone is empty.
func (Bottom) candidates(_ *EffectContext, cands []LocalID) []LocalID {
	if len(cands) == 0 {
		return nil
	}
	return cands[len(cands)-1:]
}

// pick takes the bottom card.
func (s Bottom) pick(ctx *EffectContext, cands []LocalID) []LocalID {
	return s.candidates(ctx, cands)
}

// whoseHand renders the possessive for a player's hand from the controller's
// point of view: "your hand" or "your opponent's hand".
func whoseHand(p Player) string {
	if p == Opponent {
		return "your opponent's hand"
	}
	return "your hand"
}
