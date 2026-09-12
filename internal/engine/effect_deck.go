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
			Source: ctx.Source,
			Card:   id,
			Player: ctx.Controller,
		})
		ctx.Resolver.PlayFromDeck(ctx.Controller, id)
	}
}

// DiscardTopOfDeck discards the top card of a deck and puts it in context (ctx.It)
// so a following effect can react to it — Evasion Sigil cancels the fight when it
// is of the active house; A Fair Game gains Æmber for hand cards of its house.
//
// Player picks whose deck, and its rendering follows the two ways a card names a
// player. Left unset it speaks from a granted ability's card-neutral perspective
// ("its controller's deck"), resolving against the creature's controller — the
// only sensible default for an ability every creature gains (Evasion Sigil).
// Controller and Opponent are the direct first/second-person perspectives ("your
// deck", "your opponent's deck") a card played from hand uses (A Fair Game). An
// empty deck discards nothing.
type DiscardTopOfDeck struct {
	Player Player
}

// Text renders the effect from the chosen player's perspective.
func (e DiscardTopOfDeck) Text() string {
	switch e.Player {
	case Controller:
		return "discard the top card of your deck"
	case Opponent:
		return "discard the top card of your opponent's deck"
	default:
		return "discard the top card of its controller's deck"
	}
}

// Resolve discards the top card of the chosen deck (the controller's when Player
// is unset), putting it in context.
func (e DiscardTopOfDeck) Resolve(ctx *EffectContext) {
	player := ctx.Controller
	if e.Player.valid() {
		player = ctx.PlayerFor(e.Player)
	}
	id, ok := ctx.Resolver.DiscardTopOfDeck(player)
	ctx.It, ctx.HasIt = id, ok
}

// DiscardTopOfEachDeck discards the top card of each player's deck — the
// controller's first, then the opponent's — and records each discarded card on
// the context so a following ForEachDiscarded can act on it. An empty deck
// contributes no card. Bonkers Killing Machine pairs it with ForEachDiscarded.
// Amount discards that many top cards of each deck (Rigged Lottery discards five);
// the zero value discards one.
type DiscardTopOfEachDeck struct {
	// Amount is how many top cards of each deck to discard; the zero value is one.
	Amount int
}

// count is Amount with the zero value treated as one.
func (e DiscardTopOfEachDeck) count() int {
	if e.Amount < 1 {
		return 1
	}
	return e.Amount
}

// Text renders the effect.
func (e DiscardTopOfEachDeck) Text() string {
	if e.count() == 1 {
		return "discard the top card of each player's deck"
	}
	return fmt.Sprintf(
		"discard the top %d cards of each player's deck", e.count())
}

// Resolve discards the controller's top deck cards, then the opponent's,
// recording the discarded cards on the context.
func (e DiscardTopOfEachDeck) Resolve(ctx *EffectContext) {
	ctx.Produced.Discarded = nil
	for _, player := range []int{ctx.Controller, ctx.Opponent()} {
		for range e.count() {
			if discarded, ok := ctx.Resolver.DiscardTopOfDeck(player); ok {
				ctx.Produced.Discarded = append(ctx.Produced.Discarded, discarded)
			}
		}
	}
}

// ForEachDiscarded resolves Do once for each card a preceding DiscardTopOfEachDeck
// discarded, putting that card in context (ctx.It) so Do can refer to it — Bonkers
// Killing Machine destroys a creature or artifact of each discarded card's house
// (Do targets Target.OfContextualHouse). A House filter narrows the iteration to
// the discarded cards of one house — Fetchdrones acts "for each Logos card
// discarded this way".
type ForEachDiscarded struct {
	// House, when set, restricts the iteration to discarded cards of that house.
	House House
	Do    Effect
}

// validate surfaces a configuration error from Do.
func (e ForEachDiscarded) validate() error { return validateEffect(e.Do) }

// Text renders the effect, leading with the iteration clause.
func (e ForEachDiscarded) Text() string {
	card := "card"
	if e.House != HouseNone {
		card = e.House.String() + " card"
	}
	return "for each " + card + " discarded this way, " + e.Do.Text()
}

// Resolve runs Do for each discarded card (of the House filter when set), in
// context as ctx.It.
func (e ForEachDiscarded) Resolve(ctx *EffectContext) {
	for _, id := range ctx.Produced.Discarded {
		if e.House != HouseNone && ctx.Resolver.House(id) != e.House {
			continue
		}
		ctx.It, ctx.HasIt = id, true
		e.Do.Resolve(ctx)
	}
}

// DiscardTop discards the top Amount cards of one player's deck, recording each on
// the context so a following ForEachDiscarded can act on it — Fetchdrones discards
// the top two cards of your deck. An empty deck contributes no card.
type DiscardTop struct {
	Player Player
	Amount int
}

// validate rejects a non-positive Amount.
func (e DiscardTop) validate() error {
	if e.Amount < 1 {
		return fmt.Errorf("DiscardTop: Amount must be at least 1")
	}
	return nil
}

// Text renders the effect from the chosen player's perspective.
func (e DiscardTop) Text() string {
	noun := "cards"
	if e.Amount == 1 {
		noun = "card"
	}
	deck := "your deck"
	if e.Player == Opponent {
		deck = "your opponent's deck"
	}
	return fmt.Sprintf("discard the top %d %s of %s", e.Amount, noun, deck)
}

// Resolve discards the top Amount cards of the chosen deck (the controller's when
// Player is unset), recording the discarded cards on the context.
func (e DiscardTop) Resolve(ctx *EffectContext) {
	ctx.Produced.Discarded = nil
	player := ctx.Controller
	if e.Player.valid() {
		player = ctx.PlayerFor(e.Player)
	}
	for range e.Amount {
		if id, ok := ctx.Resolver.DiscardTopOfDeck(player); ok {
			ctx.Produced.Discarded = append(ctx.Produced.Discarded, id)
		}
	}
}

// chooseFromTopOfDeck is the shared core behind LookAtTopOfDeck and RevealTopOfDeck.
// It reads the top amount cards of a deck and runs an ordered list of routing steps
// over them, each step seeing only the cards earlier steps left behind. chooseDeck
// has the controller pick whose deck to read (their own or the opponent's); public
// reveals the cards to both players and binds the top one in context (ctx.It) so a
// following effect can inspect or play it. It reads as many as remain when the deck
// holds fewer than amount, and does nothing on an empty deck.
type chooseFromTopOfDeck struct {
	amount     int
	chooseDeck bool
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
)

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
	default:
		return fmt.Sprintf("put %d into your hand", a.Count)
	}
}

// validate rejects a non-positive Count.
func (a ChooseAndMove) validate() error { return positiveCount("ChooseAndMove", a.Count) }

// terminal reports false: a move takes only Count cards, not the rest.
func (ChooseAndMove) terminal() bool { return false }

// apply moves the chosen cards from the deck to Dest.
func (a ChooseAndMove) apply(ctx *EffectContext, tr *topRead) {
	switch a.Dest {
	case IntoArchives:
		tr.choose(ctx, a.Count, "Choose a card to archive", ctx.Resolver.ArchiveFromDeck)
	case IntoDiscard:
		tr.choose(ctx, a.Count, "Choose a card to discard", ctx.Resolver.MoveFromDeckToDiscard)
	case IntoPurge:
		tr.choose(ctx, a.Count, "Choose a revealed card to purge", func(id LocalID) {
			ctx.Resolver.PurgeFromDeck(tr.player, id)
		})
	default:
		tr.choose(
			ctx,
			a.Count,
			"Choose a card to put into your hand",
			ctx.Resolver.MoveFromDeckToHand,
		)
	}
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

// apply has the controller place the remaining read cards on top in a chosen order;
// fewer than two cards leaves nothing to reorder, and a declined choice keeps the
// original order.
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
	ctx.Resolver.SetDeckTop(tr.player, order)
}

// ShuffleDeck also serves as a routing terminal (the type is the "shuffle your
// deck" effect in effect_shuffle.go): as a terminal it shuffles the whole read
// deck, since the cards revealed this way are still in that deck. It must be the
// last step.

// clause renders the step.
func (ShuffleDeck) clause() string { return "shuffle that deck" }

// validate always passes: a shuffle carries no count.
func (ShuffleDeck) validate() error { return nil }

// terminal reports true: it consumes whatever earlier steps left.
func (ShuffleDeck) terminal() bool { return true }

// apply shuffles the read player's deck — the revealed cards among them — and logs it.
func (ShuffleDeck) apply(ctx *EffectContext, tr *topRead) {
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
	noun := "cards"
	if e.Amount == 1 {
		noun = "card"
	}
	parts := []string{fmt.Sprintf("look at the top %d %s of your deck", e.Amount, noun)}
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
// With ChooseWhoseDeck the controller picks whose deck to reveal. Revealing a single
// card with no steps is the classic inspect-and-play primitive — Chaos Portal, Book
// of leQ, Wormhole Technician, Vespilon Theorist, Gambling Den; Borr Nit and Borr
// Nit's Touch reveal five from a chosen deck, purge one, and shuffle the rest.
type RevealTopOfDeck struct {
	Amount          int
	ChooseWhoseDeck bool
	Then            []TopAct
}

// core builds the public read this sugar wraps.
func (e RevealTopOfDeck) core() chooseFromTopOfDeck {
	return chooseFromTopOfDeck{
		amount:     e.Amount,
		chooseDeck: e.ChooseWhoseDeck,
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
	if e.ChooseWhoseDeck {
		deck = "a player's deck"
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

// DiscardDeckUntil digs through the top of your deck, discarding as it goes,
// until it turns up a card the filters admit or the deck runs out. The card it
// finds stays in the discard pile and goes into context (ctx.It), so what happens
// to it is a separate effect gated on the dig succeeding — Sound the Horns and
// Invasion Portal both pair it with PutDiscardedIntoHand.
type DiscardDeckUntil struct {
	// Type filters what ends the dig; the zero value stops at any card.
	Type CardType
	// House filters what ends the dig; HouseNone stops at any house.
	House House
}

// Text renders the dig and names both ways it can end, as the cards do.
func (e DiscardDeckUntil) Text() string {
	return "discard cards from the top of your deck until you discard " +
		indefinite(e.noun()) + " or run out of cards"
}

// noun names the cards the filters admit, e.g. "card" or "Brobnar creature".
func (e DiscardDeckUntil) noun() string {
	noun := "card"
	switch e.Type {
	case Creature:
		noun = "creature"
	case Artifact:
		noun = "artifact"
	}
	if e.House != HouseNone {
		noun = e.House.String() + " " + noun
	}
	return noun
}

// matches reports whether a discarded card is the one the dig was looking for.
func (e DiscardDeckUntil) matches(ctx *EffectContext, id LocalID) bool {
	if e.Type != TypeUnset && ctx.Resolver.TypeOf(id) != e.Type {
		return false
	}
	return e.House == HouseNone || ctx.Resolver.House(id) == e.House
}

// Resolve digs, leaving the found card in context.
func (e DiscardDeckUntil) Resolve(ctx *EffectContext) { e.resolveGate(ctx) }

// resolveGate digs and reports whether it found a matching card, so a Then can
// hang "put it into your hand" off the dig succeeding.
func (e DiscardDeckUntil) resolveGate(ctx *EffectContext) bool {
	ctx.It, ctx.HasIt = 0, false
	for {
		id, ok := ctx.Resolver.DiscardTopOfDeck(ctx.Controller)
		if !ok {
			return false
		}
		if e.matches(ctx, id) {
			ctx.It, ctx.HasIt = id, true
			return true
		}
	}
}

// PutDiscardedIntoHand takes the card in context out of the discard pile and
// into its owner's hand. It is the tail of a dig through the deck (DiscardDeckUntil)
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

// RevealDeckUntilHouse reveals cards from the top of the controller's deck one at
// a time, archiving each as it is revealed, until it reveals a card of House or
// the controller chooses to stop. It reports whether a card of House was revealed
// (via resolveGate), so a Then can hang a follow-up on that — Old Boomy deals
// itself 2 damage when it turns up a card of its own house.
type RevealDeckUntilHouse struct {
	// House ends the dig: revealing a card of this house stops it. It must be set.
	House House
}

// validate requires a house to stop on.
func (e RevealDeckUntilHouse) validate() error {
	if e.House == HouseNone {
		return fmt.Errorf("RevealDeckUntilHouse: House must be set")
	}
	return nil
}

// Text renders the dig, naming the house it stops on and the stop choice.
func (e RevealDeckUntilHouse) Text() string {
	return "reveal cards from the top of your deck until you reveal " +
		indefinite(e.House.String()+" card") +
		" or choose to stop, archiving each card revealed this way"
}

// Resolve digs, discarding the found-a-house report.
func (e RevealDeckUntilHouse) Resolve(ctx *EffectContext) { e.resolveGate(ctx) }

// resolveGate reveals and archives cards from the top of the deck until a card of
// House turns up or the controller stops, reporting whether a card of House was
// revealed so a Then can act on the dig succeeding.
func (e RevealDeckUntilHouse) resolveGate(ctx *EffectContext) bool {
	for {
		id, ok := ctx.Resolver.TopOfDeck(ctx.Controller)
		if !ok {
			return false
		}
		ctx.Resolver.Record(CardsRevealedToAll{
			Player: ctx.Controller,
			Cards:  []LocalID{id},
		})
		matched := ctx.Resolver.House(id) == e.House
		ctx.Resolver.ArchiveTopOfDeck(ctx.Controller)
		if matched {
			return true
		}
		if ctx.ChooseOption(
			"Reveal another card from the top of your deck?",
			[]string{"Reveal another card", "Stop"},
		) == 1 {
			return false
		}
	}
}
