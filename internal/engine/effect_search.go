package engine

import (
	"errors"
	"fmt"
	"strings"
)

// errSearchShuffleNotTopOfDeck rejects a mid-search shuffle on a search that does
// not place its finds on top of the deck: shuffling only matters when the found
// cards go back onto the deck, so the flag is a definition error elsewhere.
var errSearchShuffleNotTopOfDeck = errors.New(
	"Search: ShuffleBeforePlacing requires Dest ToTopOfDeck")

// Search is the KeyForge "search" keyword: the controller looks through one or
// more of their own hidden or private zones for cards matching a filter,
// optionally reveals what they take, and moves it to a destination. It
// generalizes searching across the axes cards vary on: which zones to look
// through (Sources), what to look for (House and Filter — name, trait, or type),
// how many to take (one, up to Max, or every match when Any is set), whether the
// taken card is revealed (Reveal), and where the found cards go (Dest). A search
// never shuffles; a card follows it with a separate Shuffle, enforced by a card
// lint.
type Search struct {
	// Sources are the zones searched, in order. A search must name at least one
	// zone; there is no assumed default (validate rejects an empty Sources).
	Sources []Zone
	// House restricts the search to cards a matcher admits; the zero value (any
	// house) does not narrow by house. House and Filter conjoin.
	House HouseMatcher
	// Filter restricts the search by name, trait, or type; the zero value admits
	// any card.
	Filter CardFilter
	// Any takes every matching card (any number). Otherwise the controller chooses
	// one, or up to Max when Max is set.
	Any bool
	// Max caps how many cards the controller may take, each choice optional (the
	// gigantic tutors take up to two halves). The zero value takes exactly one; Any
	// overrides it to take every match.
	Max int
	// Reveal shows the taken card to both players. Every card states this
	// explicitly: whether a search reveals is never inferred from its filter.
	Reveal bool
	// ShuffleBeforePlacing shuffles the controller's deck after the found cards are
	// chosen and revealed but before they are placed, so a search that puts its
	// finds on top of the deck lands them atop an already-shuffled deck (Digging Up
	// the Monster). It is only meaningful with Dest ToTopOfDeck.
	ShuffleBeforePlacing bool
	// Dest is where found cards go; it must be set to one of ToHand, ToTopOfDeck,
	// or ToArchives.
	Dest Destination
}

// validate requires at least one source zone and a supported destination: a search
// must name the zones it looks through (Search with no Sources is a definition
// error rather than a silent "search the deck") and where its finds go (an unset
// destination is a definition error rather than a silent "put them into the hand").
func (e Search) validate() error {
	if len(e.Sources) == 0 {
		return errUnsetZone("Search")
	}
	if !e.Dest.valid() {
		return errUnsetDestination("Search")
	}
	if e.Dest != ToHand && e.Dest != ToTopOfDeck && e.Dest != ToArchives {
		return fmt.Errorf("Search: unsupported destination %d", e.Dest.zone)
	}
	if e.ShuffleBeforePlacing && e.Dest != ToTopOfDeck {
		return errSearchShuffleNotTopOfDeck
	}
	return nil
}

// revealsFound reports whether a taken card is shown to both players.
func (e Search) revealsFound() bool { return e.Reveal }

// noun renders the kind of card the search takes, e.g. "card", "Niffle creature",
// or "Saurian card".
func (e Search) noun() string { return e.House.qualify(e.Filter.noun()) }

// zonesPhrase renders the searched zones, e.g. "deck" or "deck and discard pile".
func (e Search) zonesPhrase() string {
	names := map[Zone]string{
		Deck:    "deck",
		Discard: "discard pile",
		Hand:    "hand",
	}
	parts := make([]string, 0, len(e.Sources))
	for _, z := range e.Sources {
		parts = append(parts, names[z])
	}
	return strings.Join(parts, " and ")
}

// destPhrase renders where the search puts the cards it takes, always one of the
// controller's own zones.
func (e Search) destPhrase() string {
	switch e.Dest {
	case ToArchives:
		return "your archives"
	case ToTopOfDeck:
		return "the top of your deck"
	default:
		return "your hand"
	}
}

// plural reports whether the search takes more than one card, so its text and
// pronouns read in the plural ("reveal them").
func (e Search) plural() bool { return e.Any || e.Max >= 2 }

// objectPhrase renders the searched-for object with its count, e.g. "a
// Timetraveller", "any number of Ancient Bears", "either half of a gigantic
// creature", or "two halves of a gigantic creature".
func (e Search) objectPhrase() string {
	if e.Any {
		return "any number of " + plural(2, e.noun())
	}
	if e.Filter.Gigantic {
		if e.Max >= 2 {
			return "two halves of a gigantic creature"
		}
		return "either half of a gigantic creature"
	}
	return indefinite(e.noun())
}

// Text renders the effect, e.g. "search your deck and discard pile for a
// Timetraveller, reveal it, and put it into your hand".
func (e Search) Text() string {
	base := "search your " + e.zonesPhrase() + " for " + e.objectPhrase()
	pronoun := "it"
	if e.plural() {
		pronoun = "them"
	}
	putClause := "put " + pronoun + " into " + e.destPhrase()
	switch {
	case e.ShuffleBeforePlacing && e.revealsFound():
		return base + ", reveal " + pronoun + ", shuffle your deck, and " + putClause
	case e.ShuffleBeforePlacing:
		return base + ", shuffle your deck, and " + putClause
	case e.revealsFound():
		return base + ", reveal " + pronoun + ", and " + putClause
	default:
		return base + " and " + putClause
	}
}

// Resolve searches the zones for matching cards and takes them.
func (e Search) Resolve(ctx *EffectContext) { e.resolveGate(ctx) }

// resolveGate searches and reports whether it took anything, so a Then can hang a
// follow-up off the search succeeding (Bear Flute reshuffles only if it did).
func (e Search) resolveGate(ctx *EffectContext) bool {
	mover := crossZoneMover{
		Player:  ctx.Controller,
		Dest:    e.Dest,
		Sources: e.Sources,
	}
	if e.ShuffleBeforePlacing {
		return e.resolveShuffleBeforePlacing(ctx, mover)
	}
	candidates := mover.gather(ctx, func(id LocalID) bool {
		return e.House.matches(ctx, id) && e.Filter.admits(ctx.Resolver, id)
	})
	if e.Any {
		for _, id := range candidates {
			e.take(ctx, mover, id)
		}
		return len(candidates) > 0
	}
	if e.Max > 0 {
		return e.resolveUpToMax(ctx, mover)
	}
	id, ok := ctx.ChooseCard(
		"Choose "+indefinite(e.noun())+" to put into "+e.destPhrase(), candidates)
	if !ok {
		return false
	}
	e.take(ctx, mover, id)
	return true
}

// resolveUpToMax lets the controller take up to Max matching cards, one optional
// choice at a time, re-gathering after each so a taken card drops out. It reports
// whether it took at least one.
func (e Search) resolveUpToMax(ctx *EffectContext, mover crossZoneMover) bool {
	took := false
	for range e.Max {
		cands := mover.gather(ctx, func(id LocalID) bool {
			return e.House.matches(ctx, id) && e.Filter.admits(ctx.Resolver, id)
		})
		if len(cands) == 0 {
			break
		}
		id, ok := ctx.ChooseCardOptional(
			"Choose "+indefinite(e.noun())+" to put into "+e.destPhrase(), cands)
		if !ok {
			break
		}
		e.take(ctx, mover, id)
		took = true
	}
	return took
}

// resolveShuffleBeforePlacing chooses and reveals the found cards without moving
// them, shuffles the controller's deck, then places the finds on top — so a chosen
// card leaves the deck and lands on top of an already-shuffled deck rather than
// being buried by a shuffle that follows placement (Digging Up the Monster). The
// chosen cards stay in their source zone through the shuffle, identified by id, so
// no holding zone is needed.
func (e Search) resolveShuffleBeforePlacing(ctx *EffectContext, mover crossZoneMover) bool {
	chosen := e.chooseFound(ctx, mover)
	if e.revealsFound() && len(chosen) > 0 {
		ctx.Resolver.Record(CardsRevealedToAll{
			Player: ctx.Controller,
			Cards:  chosen,
		})
	}
	ctx.Resolver.Shuffle(ctx.Controller)
	ctx.Resolver.Record(DeckShuffled{Player: ctx.Controller})
	for _, id := range chosen {
		mover.move(ctx, id)
	}
	return len(chosen) > 0
}

// chooseFound collects the cards the controller takes, without moving them: up to
// Max optional picks when Max is set, otherwise a single mandatory pick. It
// re-gathers after each pick and excludes already-picked cards, so each choice is
// distinct.
func (e Search) chooseFound(ctx *EffectContext, mover crossZoneMover) []LocalID {
	picked := map[LocalID]bool{}
	var chosen []LocalID
	single := e.Max == 0
	limit := e.Max
	if single {
		limit = 1
	}
	for range limit {
		cands := mover.gather(ctx, func(id LocalID) bool {
			return !picked[id] && e.House.matches(ctx, id) && e.Filter.admits(ctx.Resolver, id)
		})
		if len(cands) == 0 {
			break
		}
		prompt := "Choose " + indefinite(e.noun()) + " to put into " + e.destPhrase()
		var (
			id LocalID
			ok bool
		)
		if single {
			id, ok = ctx.ChooseCard(prompt, cands)
		} else {
			id, ok = ctx.ChooseCardOptional(prompt, cands)
		}
		if !ok {
			break
		}
		picked[id] = true
		chosen = append(chosen, id)
	}
	return chosen
}

// take reveals one found card (when the search reveals) and moves it to the
// destination from whichever source zone holds it.
func (e Search) take(ctx *EffectContext, mover crossZoneMover, id LocalID) {
	if e.revealsFound() {
		ctx.Resolver.Record(CardsRevealedToAll{
			Player: ctx.Controller,
			Cards:  []LocalID{id},
		})
	}
	mover.move(ctx, id)
}
