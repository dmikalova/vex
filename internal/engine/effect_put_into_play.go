package engine

// Control names whose control a card enters under when it is put into play: its
// owner's (the default) or the resolving player's. It renders its own text so
// PutIntoPlay does not branch on a bool.
type Control uint8

const (
	// ControlOwner puts the card under its owner's control (the default).
	ControlOwner Control = iota
	// ControlYours puts the card under the resolving player's control (Overlord
	// Greking reanimates a destroyed enemy "under your control").
	ControlYours
)

// suffix renders the "under your control" clause, empty for owner control.
func (c Control) suffix() string {
	if c == ControlYours {
		return " under your control"
	}
	return ""
}

// PutIntoPlay puts each targeted card into play without playing it. Putting a
// card into play is distinct from playing it: bonus icons and Play: abilities do
// not resolve (only "enters play" reactions do), which is what lets an effect put
// an opponent's card into play without making that player's play decisions.
// Control puts the card under the resolving player's control (Overlord Greking
// reanimates a destroyed enemy "into play under your control") or, by default,
// under its owner's. Ownership never changes. Ready puts a creature into play
// ready rather than exhausted — Saurian Egg hatches its Saurian creatures ready.
// A gigantic half is never put into play — a lone half cannot enter play, so it
// stays where it came from.
type PutIntoPlay struct {
	Target  Target
	Control Control
	// Ready, when set, puts the card into play ready instead of exhausted.
	Ready bool
}

// validate requires an explicit target.
func (e PutIntoPlay) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("PutIntoPlay")
	}
	return nil
}

// Text renders the effect, e.g. "put it into play under your control".
func (e PutIntoPlay) Text() string {
	ready := ""
	if e.Ready {
		ready = " ready"
	}
	return "put " + e.Target.Text() + " into play" + ready + e.Control.suffix()
}

// Resolve puts each selected card into play, under the resolving player's control
// when Control is ControlYours and under its owner's otherwise.
func (e PutIntoPlay) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		controller := ctx.Resolver.Owner(id)
		if e.Control == ControlYours {
			controller = ctx.Controller
		}
		ctx.Resolver.PutIntoPlay(id, controller)
		if e.Ready {
			ctx.Resolver.SetExhausted(id, false)
		}
	}
}

// EachPlayerPutsHandCreaturesIntoPlay has both players reveal their hand and put
// every creature from it into play, in a player order the active player chooses
// (Aemberlution). The creatures are put into play, not played, so their Play:
// abilities never fire — which is what lets the active player resolve this over
// the opponent's hand without making the opponent's play decisions. Ready puts
// the creatures into play ready rather than exhausted. There is no hand-creature
// Target, so the node walks each player's hand directly rather than through the
// Target machinery.
type EachPlayerPutsHandCreaturesIntoPlay struct {
	// Ready, when set, puts the creatures into play ready instead of exhausted.
	Ready bool
}

// Text renders the effect, e.g. "each player reveals their hand and puts each
// creature from their hand into play ready".
func (e EachPlayerPutsHandCreaturesIntoPlay) Text() string {
	ready := ""
	if e.Ready {
		ready = " ready"
	}
	return "each player reveals their hand and puts each creature from their hand into play" + ready
}

// Resolve reveals each player's hand and puts its creatures into play, letting the
// active player choose which player resolves first.
func (e EachPlayerPutsHandCreaturesIntoPlay) Resolve(ctx *EffectContext) {
	first := ctx.Controller
	choice := ctx.ChooseOption(
		"Choose which player reveals and puts their hand creatures into play first",
		[]string{
			ctx.Resolver.PlayerName(ctx.Controller) + " first",
			ctx.Resolver.PlayerName(ctx.Opponent()) + " first",
		},
	)
	if choice == 1 {
		first = ctx.Opponent()
	}
	e.forPlayer(ctx, first)
	e.forPlayer(ctx, 1-first)
}

// forPlayer reveals one player's hand and puts each creature in it into play under
// that player's control.
func (e EachPlayerPutsHandCreaturesIntoPlay) forPlayer(ctx *EffectContext, player int) {
	sub := *ctx
	sub.Controller = player
	RevealHand{Player: Controller}.Resolve(&sub)
	for _, id := range creaturesInHand(ctx.Resolver, player) {
		ctx.Resolver.PutIntoPlay(id, ctx.Resolver.Owner(id))
		if e.Ready {
			ctx.Resolver.SetExhausted(id, false)
		}
	}
}

// creaturesInHand snapshots the creatures currently in a player's hand, so putting
// them into play (which empties them from the hand) does not disturb the walk.
func creaturesInHand(r Resolver, player int) []LocalID {
	var creatures []LocalID
	for _, id := range r.Hand(player) {
		if r.IsCreature(id) {
			creatures = append(creatures, id)
		}
	}
	return creatures
}
