package engine

import "fmt"

// PlayRevealedCard plays the card in context (put there by a preceding
// RevealTopOfDeck) from the controller's deck — Chaos Portal. It does nothing when
// no card is in context.
type PlayRevealedCard struct{}

// Text renders the effect.
func (PlayRevealedCard) Text() string { return "play it" }

// Resolve plays the context card from the deck.
func (PlayRevealedCard) Resolve(ctx *EffectContext) {
	if ctx.HasIt {
		ctx.Resolver.PlayFromDeck(ctx.Controller, ctx.It)
	}
}

// PutRevealedCard moves the card in context (put there by a preceding
// RevealTopOfDeck) from its owner's deck to To — Wormhole Technician archives the
// revealed card when it is not a Logos card, Vespilon Theorist discards it when it
// is not of the chosen house. It renders the terse verb its destination reads
// ("archive it", "discard it") and does nothing when no card is in context.
type PutRevealedCard struct {
	// To names where the revealed card goes: its owner's archives, discard pile,
	// hand, or the purge pile.
	To DeckDest
}

// validate rejects an unknown destination.
func (e PutRevealedCard) validate() error { return e.To.validate() }

// Text renders the effect.
func (e PutRevealedCard) Text() string { return e.To.itClause() }

// Resolve moves the context card from its owner's deck to To.
func (e PutRevealedCard) Resolve(ctx *EffectContext) {
	if ctx.HasIt {
		e.To.mover(ctx, ctx.Resolver.Owner(ctx.It))(ctx.It)
	}
}

// PlayTopOfDeck plays the top card of the controller's deck outright (Wild
// Wormhole), resolving that card's own play effect. It does nothing when the deck
// is empty.
type PlayTopOfDeck struct{}

// Text renders the effect.
func (PlayTopOfDeck) Text() string { return "play the top card of your deck" }

// Resolve plays the top card of the controller's deck, if any.
func (PlayTopOfDeck) Resolve(ctx *EffectContext) {
	if id, ok := ctx.Resolver.TopOfDeck(ctx.Controller); ok {
		ctx.Resolver.Record(PlayedFromTopOfDeck{
			Card:   id,
			Player: ctx.Controller,
		})
		ctx.Resolver.PlayFromDeck(ctx.Controller, id)
	}
}

// DiscardTop discards the top cards of one or both decks. It records every card
// discarded this way on the context so a following ForEachDiscarded can act once
// per card (Fetchdrones, Bonkers Killing Machine), and when exactly one card is
// discarded it also binds that card as ctx.It so a single-card follow-up can react
// to it — Evasion Sigil cancels the fight when it is of the active house, A Fair
// Game gains Æmber for hand cards of its house, Gebuk swaps with it.
//
// Player picks whose deck. Controller and Opponent are the direct
// first/second-person perspectives ("your deck", "your opponent's deck") a card
// played from hand uses; EachPlayer discards from both decks at once, the
// controller's first (Rigged Lottery, Bonkers Killing Machine). Left unset it
// speaks from a granted ability's card-neutral perspective ("its controller's
// deck"), resolving against the creature's controller — the only sensible default
// for an ability every creature gains (Evasion Sigil). Amount discards that many
// top cards per deck; the zero value discards one. An empty deck discards nothing.
type DiscardTop struct {
	Player Player
	// Amount is how many top cards of each named deck to discard; the zero value
	// discards one.
	Amount int
}

// count is Amount with the zero value treated as one.
func (e DiscardTop) count() int {
	if e.Amount < 1 {
		return 1
	}
	return e.Amount
}

// decks names the decks the discard acts on: both under EachPlayer, the
// controller's when Player is unset (a granted ability's perspective), else the
// named one.
func (e DiscardTop) decks(ctx *EffectContext) []int {
	switch e.Player {
	case EachPlayer:
		return []int{ctx.Controller, ctx.Opponent()}
	case Controller, Opponent:
		return []int{ctx.PlayerFor(e.Player)}
	default:
		return []int{ctx.Controller}
	}
}

// deckPhrase renders whose deck the cards come off, from the chosen perspective.
func (e DiscardTop) deckPhrase() string {
	switch e.Player {
	case Controller:
		return "your deck"
	case Opponent:
		return "your opponent's deck"
	case EachPlayer:
		return "each player's deck"
	default:
		return "its controller's deck"
	}
}

// Text renders the effect from the chosen player's perspective.
func (e DiscardTop) Text() string {
	if e.count() == 1 {
		return "discard the top card of " + e.deckPhrase()
	}
	return fmt.Sprintf(
		"discard the top %d cards of %s", e.count(), e.deckPhrase())
}

// Resolve discards the top cards of each named deck, recording them on the context
// and, when exactly one card is discarded, binding it as ctx.It for a single-card
// follow-up.
func (e DiscardTop) Resolve(ctx *EffectContext) {
	resetDiscardedThisWay(ctx)
	for _, player := range e.decks(ctx) {
		for range e.count() {
			if id, ok := ctx.Resolver.DiscardTopOfDeck(player); ok {
				recordDiscardedThisWay(ctx, id)
			}
		}
	}
	if run := discardedThisWay(ctx); len(run) == 1 {
		ctx.It, ctx.HasIt = run[0], true
	} else {
		ctx.It, ctx.HasIt = 0, false
	}
}

// The "discarded this way" binding is the run of cards a dig (DiscardTop,
// DiscardUntil) or a hand/archives discard (DiscardCard) sets aside for a
// following ForEachDiscarded or ArchiveDiscardedThisWay to act on. These four
// helpers are the only code that names the underlying ctx.Produced.Discarded field,
// so producers and consumers share one seam.

// resetDiscardedThisWay clears the recorded run at the start of a dig or discard.
func resetDiscardedThisWay(ctx *EffectContext) { ctx.Produced.Discarded = nil }

// recordDiscardedThisWay adds a card to the run the current dig or discard is
// building.
func recordDiscardedThisWay(ctx *EffectContext, id LocalID) {
	ctx.Produced.Discarded = append(ctx.Produced.Discarded, id)
}

// discardedThisWay returns the cards recorded on the context so far.
func discardedThisWay(ctx *EffectContext) []LocalID { return ctx.Produced.Discarded }

// forEachDiscardedThisWay runs fn for each card a preceding dig or discard recorded.
func forEachDiscardedThisWay(ctx *EffectContext, fn func(id LocalID)) {
	for _, id := range ctx.Produced.Discarded {
		fn(id)
	}
}

// ForEachDiscarded resolves Do once for each card a preceding DiscardTop
// discarded, putting that card in context (ctx.It) so Do can refer to it — Bonkers
// Killing Machine destroys a creature or artifact of each discarded card's house
// (Do targets Target.OfContextualHouse). A House filter narrows the iteration to
// the discarded cards of one house — Fetchdrones acts "for each Logos card
// discarded this way". A Type filter narrows it to one card type — Saurian Egg
// reanimates only the Saurian creatures it discarded, not any Saurian artifacts.
type ForEachDiscarded struct {
	// House, when it filters, restricts the iteration to discarded cards it admits.
	House HouseMatcher
	// Type, when set, restricts the iteration to discarded cards of that type.
	Type CardType
	Do   Effect
}

// validate surfaces a configuration error from Do.
func (e ForEachDiscarded) validate() error { return validateEffect(e.Do) }

// Text renders the effect, leading with the iteration clause.
func (e ForEachDiscarded) Text() string {
	return "for each " + e.House.qualify(typeNoun(e.Type)) + " discarded this way, " + e.Do.Text()
}

// Resolve runs Do for each discarded card (of the House and Type filters when
// set), in context as ctx.It.
func (e ForEachDiscarded) Resolve(ctx *EffectContext) {
	forEachDiscardedThisWay(ctx, func(id LocalID) {
		if !e.House.matches(ctx, id) {
			return
		}
		if e.Type != TypeUnset && ctx.Resolver.TypeOf(id) != e.Type {
			return
		}
		ctx.It, ctx.HasIt = id, true
		e.Do.Resolve(ctx)
	})
}

// chooseFromTopOfDeck is the shared core behind LookAtTopOfDeck and RevealTopOfDeck.
// It reads the top amount cards of a deck and runs an ordered list of routing steps
// over them, each step seeing only the cards earlier steps left behind. chooseDeck
// has the controller pick whose deck to read (their own or the opponent's); player
// instead fixes the deck to one side (Vandalize reads the opponent's), and is unset
// when the deck is simply the controller's. public reveals the cards to both players
// and binds the top one in context (ctx.It) so a
// following effect can inspect or play it. It reads as many as remain when the deck
// holds fewer than amount, and does nothing on an empty deck.
type chooseFromTopOfDeck struct {
	amount     int
	chooseDeck bool
	player     Player
	public     bool
	then       []TopAct
}

// validate rejects a non-positive amount, validates each step, and enforces that a
// terminal step is last: a terminal step consumes whatever earlier steps left, so
// nothing may follow it (and, since it must be last, only one may appear).
func (e chooseFromTopOfDeck) validate() error {
	if e.amount < 1 {
		return fmt.Errorf("chooseFromTopOfDeck: amount must be at least 1")
	}
	for i, act := range e.then {
		if err := act.validate(); err != nil {
			return err
		}
		if act.terminal() && i != len(e.then)-1 {
			return fmt.Errorf("chooseFromTopOfDeck: a terminal step must be last")
		}
	}
	return nil
}

// resolve reads the top amount cards of the chosen deck and runs each routing step
// over them so a step sees only the cards earlier steps left behind.
func (e chooseFromTopOfDeck) resolve(ctx *EffectContext) {
	player := ctx.Controller
	if e.player.valid() {
		player = ctx.PlayerFor(e.player)
	}
	if e.chooseDeck && ctx.ChooseOption(
		"Whose deck to reveal from?",
		[]string{"your deck", "your opponent's deck"},
	) == 1 {
		player = ctx.Opponent()
	}
	top := append([]LocalID(nil), deckTop(ctx, player, e.amount)...)
	if e.public {
		if len(top) > 0 {
			ctx.It, ctx.HasIt = top[0], true
			ctx.Resolver.Record(CardsRevealedToAll{Player: player, Cards: top})
		} else {
			ctx.HasIt = false
		}
	}
	if len(top) == 0 {
		return
	}
	tr := &topRead{player: player, remaining: top}
	for _, act := range e.then {
		act.apply(ctx, tr)
	}
}

// deckTop returns the top n cards of player's deck, or as many as remain.
func deckTop(ctx *EffectContext, player, n int) []LocalID {
	deck := ctx.Resolver.Deck(player)
	return deck[:min(n, len(deck))]
}

// topRead tracks the cards a chooseFromTopOfDeck read that are still in the deck as
// routing steps consume them, alongside whose deck they are — so a step moves,
// purges, reorders, or shuffles the right player's deck.
type topRead struct {
	player    int
	remaining []LocalID
}

// choose moves up to count of the still-read cards, one the controller chooses at a
// time, via move; an empty run or a declined choice stops it early.
func (tr *topRead) choose(ctx *EffectContext, count int, prompt string, move func(LocalID)) {
	for range count {
		if len(tr.remaining) == 0 {
			return
		}
		id, ok := ctx.ChooseCard(prompt, tr.remaining)
		if !ok {
			return
		}
		move(id)
		tr.remaining = withoutID(tr.remaining, id)
	}
}

// positiveCount rejects a non-positive routing count: routing zero cards is meaningless.
func positiveCount(name string, count int) error {
	if count < 1 {
		return fmt.Errorf("%s: Count must be at least 1", name)
	}
	return nil
}

// TopAct is one routing step over the cards a chooseFromTopOfDeck read. It renders
// its own text clause, moves (or reorders/shuffles) some of the still-read cards
// seeing only what earlier steps left, and reports whether it is terminal — a
// terminal step consumes the rest, so it must be the last step.
type TopAct interface {
	// clause renders this step's lowercase text, e.g. "put 1 into your hand".
	clause() string
	validate() error
	terminal() bool
	// apply routes this step's share of the cards still in tr.
	apply(ctx *EffectContext, tr *topRead)
}

// DeckDest is where a ChooseAndMove step sends the cards it takes from the deck.
type DeckDest int

const (
	// IntoHand puts the chosen cards into the controller's hand.
	IntoHand DeckDest = iota
	// IntoArchives archives the chosen cards.
	IntoArchives
	// IntoDiscard puts the chosen cards into the discard pile.
	IntoDiscard
	// IntoPurge purges the chosen cards.
	IntoPurge
	// IntoBottomOfDeck puts the chosen cards on the bottom of the same deck (the
	// leftover of a look stays on top) — the Star Alliance mutants' Alien ability.
	IntoBottomOfDeck
)

// validate rejects a DeckDest outside the known routing destinations.
func (d DeckDest) validate() error {
	if d < IntoHand || d > IntoBottomOfDeck {
		return fmt.Errorf("DeckDest: unknown destination %d", d)
	}
	return nil
}

// mover returns the resolver call that moves one card from player's deck to d —
// the single deck-routing dispatch shared by ChooseAndMove and PutRevealedCard.
func (d DeckDest) mover(ctx *EffectContext, player int) func(LocalID) {
	switch d {
	case IntoArchives:
		return ctx.Resolver.ArchiveFromDeck
	case IntoDiscard:
		return ctx.Resolver.MoveFromDeckToDiscard
	case IntoPurge:
		return func(id LocalID) { purgeFrom(ctx, Deck, player, id) }
	case IntoBottomOfDeck:
		return ctx.Resolver.PutDeckCardOnBottom
	default: // IntoHand
		return ctx.Resolver.MoveFromDeckToHand
	}
}

// choosePrompt is the pick prompt a ChooseAndMove step shows for routing to d.
func (d DeckDest) choosePrompt() string {
	switch d {
	case IntoArchives:
		return "Choose a card to archive"
	case IntoDiscard:
		return "Choose a card to discard"
	case IntoPurge:
		return "Choose a revealed card to purge"
	case IntoBottomOfDeck:
		return "Choose a card to put on the bottom of your deck"
	default: // IntoHand
		return "Choose a card to put into your hand"
	}
}

// itClause is the terse verb PutRevealedCard prints for the single context card,
// e.g. "archive it" or "discard it".
func (d DeckDest) itClause() string {
	switch d {
	case IntoArchives:
		return "archive it"
	case IntoDiscard:
		return "discard it"
	case IntoPurge:
		return "purge it"
	default: // IntoHand
		return "put it into your hand"
	}
}

// act renders the verb phrase routing object to d, e.g. "archive each card of the
// chosen house" or "put the others into your hand".
func (d DeckDest) act(object string) string {
	switch d {
	case IntoArchives:
		return "archive " + object
	case IntoDiscard:
		return "discard " + object
	case IntoPurge:
		return "purge " + object
	case IntoBottomOfDeck:
		return "put " + object + " on the bottom of your deck"
	default: // IntoHand
		return "put " + object + " into your hand"
	}
}

// ChooseAndMove takes Count of the read cards the controller chooses and sends them
// to Dest — their hand, archives, discard pile, or the purge pile.
type ChooseAndMove struct {
	Count int
	Dest  DeckDest
}

// clause renders the step, phrased per destination the way the printed cards read.
func (a ChooseAndMove) clause() string {
	switch a.Dest {
	case IntoArchives:
		return fmt.Sprintf("archive %d", a.Count)
	case IntoDiscard:
		return fmt.Sprintf("discard %d", a.Count)
	case IntoPurge:
		if a.Count == 1 {
			return "purge a card revealed this way"
		}
		return fmt.Sprintf("purge %d cards revealed this way", a.Count)
	case IntoBottomOfDeck:
		return fmt.Sprintf("put %d on the bottom of your deck", a.Count)
	default:
		return fmt.Sprintf("put %d into your hand", a.Count)
	}
}

// validate rejects a non-positive Count.
func (a ChooseAndMove) validate() error { return positiveCount("ChooseAndMove", a.Count) }

// terminal reports false: a move takes only Count cards, not the rest.
func (ChooseAndMove) terminal() bool { return false }

// apply moves the chosen cards from the deck to Dest through the shared dispatch.
func (a ChooseAndMove) apply(ctx *EffectContext, tr *topRead) {
	tr.choose(ctx, a.Count, a.Dest.choosePrompt(), a.Dest.mover(ctx, tr.player))
}

// ReorderRest puts the read cards no earlier step took back on top in any order the
// controller chooses. It is terminal, so it must be the last step.
type ReorderRest struct{}

// clause renders the step.
func (ReorderRest) clause() string { return "put them back in any order" }

// validate always passes: a reorder carries no count.
func (ReorderRest) validate() error { return nil }

// terminal reports true: it consumes whatever earlier steps left.
func (ReorderRest) terminal() bool { return true }

// apply has the controller place the remaining read cards back one at a time,
// each on top of the last, so the card chosen first ends up deepest and the card
// left over rides on top (drawn next). Fewer than two cards leaves nothing to
// reorder, and a declined choice keeps the original order.
func (ReorderRest) apply(ctx *EffectContext, tr *topRead) {
	if len(tr.remaining) < 2 {
		return
	}
	order := make([]LocalID, 0, len(tr.remaining))
	for len(tr.remaining) > 1 {
		id, ok := ctx.ChooseCard(
			"Choose the next card to place on top of your deck", tr.remaining,
		)
		if !ok {
			return
		}
		order = append(order, id)
		tr.remaining = withoutID(tr.remaining, id)
	}
	order = append(order, tr.remaining[0])
	// Each pick was placed on top of the previous, so the pick order runs
	// top-of-deck-last; reverse it into deck order (order[0] is the top card).
	for i, j := 0, len(order)-1; i < j; i, j = i+1, j-1 {
		order[i], order[j] = order[j], order[i]
	}
	ctx.Resolver.SetDeckTop(tr.player, order)
}

// MayDiscardLookedAt lets the controller optionally discard the single card a
// preceding look read (Scout Pete). It reads "you may discard that card", so it is
// meant to follow a one-card look.
type MayDiscardLookedAt struct{}

// clause renders the step.
func (MayDiscardLookedAt) clause() string { return "you may discard that card" }

// validate always passes: the optional discard carries no count.
func (MayDiscardLookedAt) validate() error { return nil }

// terminal reports false: it discards at most the one looked-at card.
func (MayDiscardLookedAt) terminal() bool { return false }

// apply offers the controller the optional discard of the read card.
func (MayDiscardLookedAt) apply(ctx *EffectContext, tr *topRead) {
	if len(tr.remaining) == 0 {
		return
	}
	id, ok := ctx.ChooseCardOptional("Discard the top card?", tr.remaining)
	if !ok {
		return
	}
	ctx.Resolver.MoveFromDeckToDiscard(id)
	tr.remaining = withoutID(tr.remaining, id)
}

// PartitionByChosenHouse splits the read cards by the house an enclosing
// ChooseHouseThen picked: cards of that house go to Matching, the rest go to Rest.
// New Frontiers archives each card of the chosen house and discards the others. It
// is terminal, so it must be the last step.
type PartitionByChosenHouse struct {
	Matching DeckDest
	Rest     DeckDest
}

// clause renders the split, e.g. "archive each card of the chosen house and
// discard the others".
func (p PartitionByChosenHouse) clause() string {
	return p.Matching.act("each card of the chosen house") + " and " +
		p.Rest.act("the others")
}

// validate rejects an unknown destination on either side.
func (p PartitionByChosenHouse) validate() error {
	if err := p.Matching.validate(); err != nil {
		return err
	}
	return p.Rest.validate()
}

// terminal reports true: the split routes every read card, so nothing follows it.
func (PartitionByChosenHouse) terminal() bool { return true }

// apply routes each read card to Matching when its house is the chosen house and to
// Rest otherwise.
func (p PartitionByChosenHouse) apply(ctx *EffectContext, tr *topRead) {
	toMatching := p.Matching.mover(ctx, tr.player)
	toRest := p.Rest.mover(ctx, tr.player)
	for _, id := range append([]LocalID(nil), tr.remaining...) {
		if ctx.Resolver.House(id) == ctx.ChosenHouse {
			toMatching(id)
		} else {
			toRest(id)
		}
		tr.remaining = withoutID(tr.remaining, id)
	}
}

// effect in effect_shuffle.go): as a terminal it shuffles the whole read deck,
// since the cards revealed this way are still in that deck. It must be the last
// step, and only the zero value Shuffle{} is used here (its validate lives with
// the effect in effect_shuffle.go).

// clause renders the step.
func (Shuffle) clause() string { return "shuffle that deck" }

// terminal reports true: it consumes whatever earlier steps left.
func (Shuffle) terminal() bool { return true }

// apply shuffles the read player's deck — the revealed cards among them — and logs it.
func (Shuffle) apply(ctx *EffectContext, tr *topRead) {
	ctx.Resolver.Shuffle(tr.player)
	ctx.Resolver.Record(DeckShuffled{Player: tr.player})
}

// LookAtTopOfDeck looks privately at the top Amount cards of the controller's deck,
// then routes them through the ordered Then steps (see chooseFromTopOfDeck). With no
// steps it is a pure peek — the cards never move. A card no step routes stays on top
// in its original position — Eyegor, Philophosaurus, Navigator Ali, Lay of the Land.
type LookAtTopOfDeck struct {
	Amount int
	Then   []TopAct
}

// core builds the private, own-deck read this sugar wraps.
func (e LookAtTopOfDeck) core() chooseFromTopOfDeck {
	return chooseFromTopOfDeck{amount: e.Amount, then: e.Then}
}

// validate delegates to the core read.
func (e LookAtTopOfDeck) validate() error { return e.core().validate() }

// Text names the peek and folds each routing step's clause into one instruction.
func (e LookAtTopOfDeck) Text() string {
	peek := "look at the top card of your deck"
	if e.Amount != 1 {
		peek = fmt.Sprintf("look at the top %d cards of your deck", e.Amount)
	}
	parts := []string{peek}
	for _, act := range e.Then {
		parts = append(parts, act.clause())
	}
	return serialJoin(parts, " and ")
}

// Resolve carries out the private read.
func (e LookAtTopOfDeck) Resolve(ctx *EffectContext) { e.core().resolve(ctx) }

// RevealTopOfDeck reveals the top Amount cards of a deck to both players and binds
// the top one in context (ctx.It) so a following effect can inspect or play it, then
// routes the revealed cards through the ordered Then steps (see chooseFromTopOfDeck).
// With ChooseWhoseDeck the controller picks whose deck to reveal. With Player set,
// the deck is fixed to that side without a choice (Vandalize reveals the opponent's).
// Revealing a single
// card with no steps is the classic inspect-and-play primitive — Chaos Portal, Book
// of leQ, Wormhole Technician, Vespilon Theorist, Gambling Den; Borr Nit and Borr
// Nit's Touch reveal five from a chosen deck, purge one, and shuffle the rest.
type RevealTopOfDeck struct {
	Amount          int
	ChooseWhoseDeck bool
	Player          Player
	Then            []TopAct
}

// core builds the public read this sugar wraps.
func (e RevealTopOfDeck) core() chooseFromTopOfDeck {
	return chooseFromTopOfDeck{
		amount:     e.Amount,
		chooseDeck: e.ChooseWhoseDeck,
		player:     e.Player,
		public:     true,
		then:       e.Then,
	}
}

// validate delegates to the core read.
func (e RevealTopOfDeck) validate() error { return e.core().validate() }

// Text names the reveal and punctuates each routing step as its own sentence, the
// way the printed cards read.
func (e RevealTopOfDeck) Text() string {
	deck := "your deck"
	switch {
	case e.ChooseWhoseDeck:
		deck = "a player's deck"
	case e.Player == Opponent:
		deck = "your opponent's deck"
	}
	text := "reveal the top card of " + deck
	if e.Amount != 1 {
		text = fmt.Sprintf("reveal the top %d cards of %s", e.Amount, deck)
	}
	for _, act := range e.Then {
		text += ". " + capitalizeFirst(act.clause())
	}
	return text
}

// Resolve carries out the public read.
func (e RevealTopOfDeck) Resolve(ctx *EffectContext) { e.core().resolve(ctx) }

// CancelFight makes the fight in progress not occur — a "Before Fight" effect
// (Evasion Sigil, gated on the discarded card's house). The attacker was still used
// to fight, so it stays exhausted; combat reads the cancellation and skips Assault,
// Hazardous, fight damage, and Fight: abilities.
type CancelFight struct{}

// Text renders the effect.
func (CancelFight) Text() string { return "the fight does not occur" }

// Resolve cancels the current fight.
func (CancelFight) Resolve(ctx *EffectContext) { ctx.Resolver.CancelCurrentFight() }

// resolverInPlay reports whether id appears in either player's battleline or
// artifact row using only Resolver reads.
func resolverInPlay(ctx *EffectContext, id LocalID) bool {
	for p := 0; p < 2; p++ {
		for _, candidate := range ctx.Resolver.Battleline(p) {
			if candidate == id {
				return true
			}
		}
		for _, candidate := range ctx.Resolver.Artifacts(p) {
			if candidate == id {
				return true
			}
		}
	}
	return false
}

// DiscardUntil digs through the top of a named deck, discarding as it goes,
// until it turns up a card the filters admit or the deck runs out. Player names
// whose deck is dug: Controller ("your deck"), Opponent ("your opponent's deck"),
// or ItsController ("its controller's deck", the contextual card's controller —
// Purify digs the purged creature's controller). With MayStop the
// controller may also stop the dig before a match, so the terminator reads "or
// choose to stop". Every discarded card is recorded on the context
// (ctx.Produced.Discarded) and the matching card is left in context (ctx.It), so a
// following effect can act on the found card or the whole discarded run — Sound the
// Horns and Invasion Portal pair it with PutDiscardedIntoHand; Old Boomy archives
// the run with ArchiveDiscardedThisWay and, when it discards a card of its house,
// deals itself damage.
type DiscardUntil struct {
	// Player names whose deck is dug; every caller sets it explicitly.
	Player Player
	// House filters what ends the dig; an unset matcher stops at any house. House
	// and Filter conjoin, as they do on Search.
	House HouseMatcher
	// Filter restricts what ends the dig by name, trait, or type; the zero value
	// stops at any card (Angry Mob digs for another Angry Mob, Purify for a
	// non-Mutant creature).
	Filter CardFilter
	// MayStop lets the controller stop the dig before a match; the terminator then
	// reads "or choose to stop" instead of "or run out of cards".
	MayStop bool
}

// deckPhrase renders whose deck is dug, from the chosen perspective.
func (e DiscardUntil) deckPhrase() string {
	switch e.Player {
	case Opponent:
		return "your opponent's deck"
	case ItsController:
		return "its controller's deck"
	default:
		return "your deck"
	}
}

// Text renders the dig and names both ways it can end, as the cards do.
func (e DiscardUntil) Text() string {
	end := "run out of cards"
	if e.MayStop {
		end = "choose to stop"
	}
	return "discard cards from the top of " + e.deckPhrase() + " until you discard " +
		indefinite(e.noun()) + " or " + end
}

// noun names the cards the filters admit, e.g. "card", "Brobnar Creature", a
// printed name ("Angry Mob"), or a trait exclusion ("non-Mutant creature").
func (e DiscardUntil) noun() string {
	if e.Filter.Name != "" {
		return e.Filter.Name
	}
	return e.House.qualifyNoun(e.Filter.noun())
}

// matches reports whether a discarded card is the one the dig was looking for.
func (e DiscardUntil) matches(ctx *EffectContext, id LocalID) bool {
	if !e.Filter.admits(ctx.Resolver, id) {
		return false
	}
	return e.House.matches(ctx, id)
}

// Resolve digs, recording the run and leaving the found card in context.
func (e DiscardUntil) Resolve(ctx *EffectContext) { e.resolveGate(ctx) }

// resolveGate digs, recording every discarded card on the context, and reports
// whether it found a matching card so a Then can hang a follow-up on the dig
// succeeding. When MayStop is set it offers the controller a stop after each
// non-matching discard.
func (e DiscardUntil) resolveGate(ctx *EffectContext) bool {
	ctx.It, ctx.HasIt = 0, false
	resetDiscardedThisWay(ctx)
	player := ctx.PlayerFor(e.Player)
	for {
		id, ok := ctx.Resolver.DiscardTopOfDeck(player)
		if !ok {
			return false
		}
		recordDiscardedThisWay(ctx, id)
		if e.matches(ctx, id) {
			ctx.It, ctx.HasIt = id, true
			return true
		}
		if e.MayStop && ctx.ChooseOption(
			"Discard another card from the top of "+e.deckPhrase()+"?",
			[]string{"Discard another card", "Stop"},
		) == 1 {
			return false
		}
	}
}

// PutDiscardedIntoHand takes the card in context out of the discard pile and
// into its owner's hand. It is the tail of a dig through the deck (DiscardUntil)
// that just discarded the card. Type names what the dig stopped on so the tail
// reads "put the discarded creature into your hand" rather than a bare "it"; the
// zero value stays the generic "card".
type PutDiscardedIntoHand struct {
	// Type names the discarded card the dig stopped on; the zero value is "card".
	Type CardType
}

// Text renders the effect, naming the discarded card the dig stopped on.
func (e PutDiscardedIntoHand) Text() string {
	return "put the discarded " + discardedNoun(e.Type) + " into your hand"
}

// discardedNoun names a discarded card by type for the "put the discarded …"
// tail: a creature, an artifact, or a bare card when the type is unset.
func discardedNoun(t CardType) string {
	switch t {
	case Creature:
		return "creature"
	case Artifact:
		return "artifact"
	default:
		return "card"
	}
}

// Resolve moves the contextual card from the discard pile to hand.
func (e PutDiscardedIntoHand) Resolve(ctx *EffectContext) {
	if ctx.HasIt {
		ctx.Resolver.PutFromDiscardIntoHand(ctx.It)
	}
}

// PutDiscardedIntoPlay puts the card in context (the card a preceding dig stopped
// on) into play under its owner's control. It is the tail of a DiscardUntil dig
// that puts what it found onto the battleline rather than into hand — Purify digs
// the purged creature's controller's deck for a non-Mutant creature and puts that
// creature into play under its owner's control. Type names what the dig stopped on
// so the tail reads "put the discarded creature into play …" rather than a bare
// "it"; the zero value stays the generic "card".
type PutDiscardedIntoPlay struct {
	// Type names the discarded card the dig stopped on; the zero value is "card".
	Type CardType
}

// Text renders the effect, naming the discarded card the dig stopped on.
func (e PutDiscardedIntoPlay) Text() string {
	return "put the discarded " + discardedNoun(e.Type) +
		" into play under its owner's control"
}

// Resolve puts the contextual card into play under its owner's control.
func (e PutDiscardedIntoPlay) Resolve(ctx *EffectContext) {
	if ctx.HasIt {
		ctx.Resolver.PutIntoPlay(ctx.It, ctx.Resolver.Owner(ctx.It))
	}
}

// ArchiveDiscardedThisWay archives every card a preceding deck dig discarded
// (recorded on ctx.Produced.Discarded), moving each from the controller's discard
// pile into their archives — Old Boomy archives the run it dug from the top of its
// deck. An empty run archives nothing.
type ArchiveDiscardedThisWay struct{}

// Text renders the effect.
func (ArchiveDiscardedThisWay) Text() string { return "archive each card discarded this way" }

// Resolve archives each recorded discarded card from the controller's discard pile.
func (ArchiveDiscardedThisWay) Resolve(ctx *EffectContext) {
	forEachDiscardedThisWay(ctx, func(id LocalID) {
		archiveFrom(ctx, Discard, ctx.Controller, id)
	})
}
