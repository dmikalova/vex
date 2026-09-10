package engine

import "strings"

// Selection is how a zone-movement effect picks the cards it acts on from a
// zone: the controller chooses one (Chosen), one is uniformly random (Random),
// or every matching card is taken (Each). Each mode filters and picks its own
// cards and renders its own object phrase ("a card" / "a random card" / "each
// creature"), so a movement verb varies along this one axis instead of spawning
// a node per (verb x mode). See ADR 0031.
type Selection interface {
	// pick returns the cards to act on, given every card in the source zone. It
	// may prompt the controller (Chosen) or draw on the RNG (Random).
	pick(ctx *EffectContext, cands []LocalID) []LocalID
	// noun renders the bare kind of card the verb acts on, without article or
	// count, e.g. "card" / "random card" / "non-Mars creature". A count-bearing
	// verb pluralizes it; object decorates it for the single case.
	noun() string
	// object renders the noun phrase the verb acts on, e.g. "a Sanctum card".
	object() string
}

// declinableSelection is a Selection the controller may pass without acting — a
// non-mandatory Chosen. A Selection that does not implement it can never be
// declined, so its verb never reads "you may".
type declinableSelection interface {
	declinable() bool
}

// selectionDeclinable reports whether a Selection can be passed.
func selectionDeclinable(s Selection) bool {
	d, ok := s.(declinableSelection)
	return ok && d.declinable()
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
// type. A non-mandatory Chosen is a "you may": the controller can always decline,
// and an empty candidate set picks nothing. Mandatory forces the pick when a card
// matches (Greater Oxtet).
type Chosen struct {
	// House restricts the choice to cards of this house; HouseNone allows any.
	House House
	// Type restricts the choice to cards of this type; the zero value allows any.
	Type CardType
	// Mandatory forces the pick and drops the "you may"; an empty hand still picks
	// nothing.
	Mandatory bool
}

// noun renders the bare kind of card chosen, qualified by type when set and by
// house.
func (s Chosen) noun() string {
	noun := "card"
	if s.Type != TypeUnset {
		noun = strings.ToLower(s.Type.String())
	}
	if s.House != HouseNone {
		noun = s.House.String() + " " + noun
	}
	return noun
}

// object renders the single card chosen, e.g. "a Sanctum creature".
func (s Chosen) object() string { return indefinite(s.noun()) }

// declinable reports that a non-mandatory Chosen may be passed.
func (s Chosen) declinable() bool { return !s.Mandatory }

// pick offers the matching cards as a card choice — declinable unless Mandatory —
// and returns the single chosen card, or none when the controller declines or no
// card matches.
func (s Chosen) pick(ctx *EffectContext, cands []LocalID) []LocalID {
	matching := filterIDs(cands, func(id LocalID) bool {
		if s.House != HouseNone && ctx.Resolver.House(id) != s.House {
			return false
		}
		if s.Type != TypeUnset && ctx.Resolver.TypeOf(id) != s.Type {
			return false
		}
		return true
	})
	choose := ctx.ChooseCardOptional
	if s.Mandatory {
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

// ownerActs reports that a random pick from a hidden hand is attributed to the
// hand's owner, so an opponent's random discard reads "your opponent discards …".
func (Random) ownerActs() bool { return true }

// pick draws one uniformly random card from the candidates, or none when empty.
func (Random) pick(ctx *EffectContext, cands []LocalID) []LocalID {
	id, ok := ctx.ChooseRandom(cands)
	if !ok {
		return nil
	}
	return []LocalID{id}
}

// Each takes every card the filters admit, with no choice: the filters decide
// (Martians Make Bad Allies purges each non-Mars creature). A following effect
// can scale with the tally the verb records.
type Each struct {
	// Type restricts to cards of this type; the zero value admits any.
	Type CardType
	// ExceptHouse spares the cards of that house; HouseNone spares nothing.
	ExceptHouse House
	// OfChosenHouse limits the cards to the house an enclosing ChooseHouseThen
	// picked (Deep Probe discards each creature of the chosen house).
	OfChosenHouse bool
}

// noun renders the bare kind of card taken, e.g. "non-Mars creature" or "creature
// of the chosen house".
func (s Each) noun() string {
	noun := "card"
	if s.Type != TypeUnset {
		noun = strings.ToLower(s.Type.String())
	}
	if s.ExceptHouse != HouseNone {
		noun = "non-" + s.ExceptHouse.String() + " " + noun
	}
	if s.OfChosenHouse {
		noun += " of the chosen house"
	}
	return noun
}

// object renders the kind of card taken, e.g. "each non-Mars creature".
func (s Each) object() string { return "each " + s.noun() }

// pick returns every candidate the filters admit.
func (s Each) pick(ctx *EffectContext, cands []LocalID) []LocalID {
	return filterIDs(cands, func(id LocalID) bool {
		if s.Type != TypeUnset && ctx.Resolver.TypeOf(id) != s.Type {
			return false
		}
		if s.ExceptHouse != HouseNone && ctx.Resolver.House(id) == s.ExceptHouse {
			return false
		}
		if s.OfChosenHouse && ctx.Resolver.House(id) != ctx.ChosenHouse {
			return false
		}
		return true
	})
}

// whoseHand renders the possessive for a player's hand from the controller's
// point of view: "your hand" or "your opponent's hand".
func whoseHand(p Player) string {
	if p == Opponent {
		return "your opponent's hand"
	}
	return "your hand"
}
