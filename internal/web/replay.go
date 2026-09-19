package web

import (
	"fmt"
	"slices"

	"github.com/dmikalova/vactrol/internal/engine"
	"github.com/dmikalova/vactrol/internal/match"
)

// This file is the client's event-sourced core (ADR 0039, 0040). Every player
// input — each root action and each chooser answer — is recorded, in order,
// into one command log. State and the typed game log are projections of it:
// replaying the log from a fresh deal reproduces the exact match. That is what
// lets a page reload rebuild the game from {seed, sets, log} rather than a
// serialized GameState, and what lets undo rewind by replaying a prefix.

// inputKind tags what a recorded input is. The values below inPick are root
// actions (a play the player initiated); the rest are the answers the chooser
// returned for a decision the engine raised while an action resolved.
type inputKind uint8

const (
	inHouse inputKind = iota // choose the turn's active house
	inPlayCreature
	inPlayArtifact
	inPlayTactic
	inPlayUpgrade
	inDiscard
	inReap
	inUnstun
	inUseAction
	inFight
	inEndTurn

	// manual-mode roots (a playtester edit; each maps to one engine Manual* call):
	inSetManual        // toggle manual mode (OK: on)
	inManualMove       // move a card to a resting zone (ID, Index: zone)
	inManualReady      // clear a card's exhausted flag (ID)
	inManualExhaust    // set a card's exhausted flag (ID)
	inManualAttach     // thread a card under a host (Card: host, ID, Left: down)
	inManualPlace      // drop a card into play (ID, Index: battleline position)
	inManualDetach     // send an attached card to hand (ID)
	inManualAmber      // adjust a player's Æmber (Player, Delta)
	inManualUnforge    // remove a player's last key (Player)
	inManualForgeColor // forge a key of a colour (Player, Index: colour)
	inManualChains     // adjust a player's chains (Player, Delta)
	inManualHouse      // set the active player's active house (House)
	inManualAddCard    // register a card into a hand (Name, Player)

	// chooser answers (everything from here down is a decision result):
	inPick     // ChooseCreature / ChooseCardOrDecline result (ID + OK)
	inOption   // ChooseOption result (Index)
	inPosition // ChoosePosition result (Index)
	inReaction // ChooseReaction result (Index)
	inOrder    // OrderCreatures result (Order)
)

// isRoot reports whether the input is a root action rather than a chooser
// answer. The replay driver dispatches root inputs itself; chooser answers are
// pulled by the replayChooser as the engine asks for them.
func (k inputKind) isRoot() bool { return k < inPick }

// input is one recorded player input. Only the fields its Kind names carry
// meaning; the struct is flat and JSON-safe so the whole log persists as data.
type input struct {
	Kind inputKind

	House engine.House // inHouse / inManualHouse
	Hand  int          // inPlay*/inDiscard: hand index
	Left  bool         // inPlayCreature: flank side; inManualAttach: face down

	Card  engine.LocalID // inReap/Unstun/UseAction: card; inFight: attacker; inManualAttach: host
	Card2 engine.LocalID // inFight: defender

	ID    engine.LocalID   // inPick: chosen card; manual roots: the card acted on
	OK    bool             // inPick: declined?; inSetManual: on; inManualExhaust: exhausted
	Index int              // inOption/inPosition/inReaction; manual: zone/pos/colour
	Order []engine.LocalID // inOrder

	Player int    // manual roots: the target player
	Delta  int    // inManualAmber/inManualChains: signed delta
	Name   string // inManualAddCard: the card definition name
}

// record appends an input to the command log during live play. A root action
// also marks where it begins so an undo can truncate the log at a root
// boundary. It is a no-op while replaying, so a rebuild does not re-record the
// inputs it feeds back.
func (g *game) record(in input) {
	if g.replaying {
		return
	}
	if in.Kind.isRoot() {
		g.rootMarks = append(g.rootMarks, len(g.inputs))
	}
	g.inputs = append(g.inputs, in)
}

// replayChooser answers the engine's choice requests from a recorded command
// log instead of a player. It shares pos with the replay driver: the driver
// advances pos over root inputs, this advances it over the chooser answers the
// engine pulls between them, so one cursor walks the whole interleaved log in
// order.
type replayChooser struct {
	inputs []input
	pos    int
	// err holds the first divergence a choice pull hit; a Chooser method cannot
	// return an error, so it is stashed here and surfaced by the driver.
	err error
}

// replayChooser answers every prompt the engine can raise from the recorded log,
// so it must satisfy every answering capability; it omits BadgeChooser, which is
// display-only and carries no recorded answer (ADR 0045, step 1).
var (
	_ engine.Chooser           = (*replayChooser)(nil)
	_ engine.OptionChooser     = (*replayChooser)(nil)
	_ engine.PositionChooser   = (*replayChooser)(nil)
	_ engine.DeclinableChooser = (*replayChooser)(nil)
	_ engine.Orderer           = (*replayChooser)(nil)
	_ engine.ReactionChooser   = (*replayChooser)(nil)
)

// next returns the input at the cursor and advances it, requiring the recorded
// kind to match what the engine is asking for — a mismatch means the log has
// diverged from the engine, which must fail loudly rather than misreplay.
func (c *replayChooser) next(want inputKind) (input, error) {
	if c.pos >= len(c.inputs) {
		return input{}, fmt.Errorf("web replay: log exhausted, engine still asked for %v", want)
	}
	in := c.inputs[c.pos]
	if in.Kind != want {
		return input{}, fmt.Errorf(
			"web replay: log holds %v at %d, engine asked for %v", in.Kind, c.pos, want)
	}
	c.pos++
	return in, nil
}

// The engine records the first failing next() so the driver can surface it: a
// chooser method cannot return an error, so it stores one and returns a benign
// answer to let resolution unwind.
func (c *replayChooser) fail(err error) {
	if c.err == nil {
		c.err = err
	}
}

// Matching kinds is not enough to prove a replay is faithful: a recorded answer
// that names a card the engine is not offering, or an index outside the range it
// is offering, means the replay has drifted from the game that produced the log.
// The three checks below catch that, so a drifted replay fails where it drifted
// instead of quietly reconstructing a different match.

// requireCandidate rejects a recorded pick that is not among the cards the engine
// is offering. A decline names no card, so it is always legal.
func (c *replayChooser) requireCandidate(in input, cands []engine.LocalID) error {
	if !in.OK || slices.Contains(cands, in.ID) {
		return nil
	}
	return fmt.Errorf(
		"web replay: log picked card %d at %d, which the engine did not offer",
		in.ID, c.pos-1)
}

// requireIndex rejects a recorded index outside the range the engine is offering.
func (c *replayChooser) requireIndex(in input, n int) error {
	if in.Index >= 0 && in.Index < n {
		return nil
	}
	return fmt.Errorf("web replay: log chose index %d at %d, but the engine offered %d",
		in.Index, c.pos-1, n)
}

// requirePermutation rejects a recorded ordering that is not a rearrangement of
// the cards the engine asked to have ordered.
func (c *replayChooser) requirePermutation(in input, ids []engine.LocalID) error {
	if len(in.Order) == len(ids) &&
		slices.Equal(slices.Sorted(slices.Values(in.Order)), slices.Sorted(slices.Values(ids))) {
		return nil
	}
	return fmt.Errorf("web replay: log ordered %v at %d, but the engine asked to order %v",
		in.Order, c.pos-1, ids)
}

func (c *replayChooser) ChooseCreature(
	_, _ string,
	cands []engine.LocalID,
) (engine.LocalID, bool) {
	in, err := c.next(inPick)
	if err == nil {
		err = c.requireCandidate(in, cands)
	}
	if err != nil {
		c.fail(err)
		return 0, false
	}
	return in.ID, in.OK
}

func (c *replayChooser) ChooseCardOrDecline(
	_, _ string,
	cands []engine.LocalID,
) (engine.LocalID, bool) {
	in, err := c.next(inPick)
	if err == nil {
		err = c.requireCandidate(in, cands)
	}
	if err != nil {
		c.fail(err)
		return 0, false
	}
	return in.ID, in.OK
}

func (c *replayChooser) ChooseOption(_, _ string, options []string) int {
	in, err := c.next(inOption)
	if err == nil {
		err = c.requireIndex(in, len(options))
	}
	if err != nil {
		c.fail(err)
		return 0
	}
	return in.Index
}

func (c *replayChooser) ChoosePosition(_, _ string, _ []engine.LocalID) int {
	in, err := c.next(inPosition)
	if err != nil {
		c.fail(err)
		return 0
	}
	return in.Index
}

func (c *replayChooser) ChooseReaction(_ string, reactions []engine.OrderableReaction) int {
	in, err := c.next(inReaction)
	if err == nil {
		err = c.requireIndex(in, len(reactions))
	}
	if err != nil {
		c.fail(err)
		return 0
	}
	return in.Index
}

func (c *replayChooser) OrderCreatures(_, _ string, ids []engine.LocalID) []engine.LocalID {
	in, err := c.next(inOrder)
	if err == nil {
		err = c.requirePermutation(in, ids)
	}
	if err != nil {
		c.fail(err)
		return ids
	}
	return in.Order
}

// driveReplay runs setup then dispatches every root action in the log, pulling
// each action's choices from the replayChooser as the engine raises them. It
// reconstructs the per-action log groups the same way beginAction does during
// live play — one mark per root, at the log length it began — since replay does
// not go through beginAction. It stops at the first divergence.
func driveReplay(
	eg *engine.Game,
	rc *replayChooser,
	defs map[string]*engine.CardDefinition,
) ([]logMark, error) {
	eg.StartGame(0)
	if rc.err != nil {
		return nil, rc.err
	}
	var groups []logMark
	for rc.pos < len(rc.inputs) {
		in := rc.inputs[rc.pos]
		if !in.Kind.isRoot() {
			return nil, fmt.Errorf(
				"web replay: expected a root action at %d, found %v",
				rc.pos,
				in.Kind,
			)
		}
		groups = append(
			groups,
			logMark{
				Start:  len(eg.Log),
				Player: eg.State.ActivePlayer,
			},
		)
		rc.pos++
		if err := dispatchRoot(eg, in, defs); err != nil {
			return nil, err
		}
		if rc.err != nil {
			return nil, rc.err
		}
		handOffEndedTurnOn(eg)
	}
	return groups, nil
}

// replayGame rebuilds a match from its command log. It deals a fresh game from
// the same seed and sets, runs setup (mulligans answered from the log), then
// dispatches each root action in order — the choices each one raises pulled from
// the log by the replayChooser — so the exact state and typed log are
// regenerated with no serialized GameState. defs resolves the card names manual
// mode added, so those cards re-register to the same ids.
func replayGame(
	seed int64,
	sets [2]string,
	inputs []input,
	defs map[string]*engine.CardDefinition,
) (*engine.Game, error) {
	eg, _, _, _, _, err := match.NewWithSets("Player 1", "Player 2", seed, sets)
	if err != nil {
		return nil, err
	}
	rc := &replayChooser{inputs: inputs}
	eg.SetChooser(0, rc)
	eg.SetChooser(1, rc)
	if _, err := driveReplay(eg, rc, defs); err != nil {
		return nil, err
	}
	return eg, nil
}

// dispatchRoot performs one recorded root action against the game, mirroring
// the engine call the live handler made. A normal action's player is the active
// one, as it was when recorded; a manual root carries its own target.
func dispatchRoot(
	eg *engine.Game,
	in input,
	defs map[string]*engine.CardDefinition,
) error {
	p := eg.State.ActivePlayer
	switch in.Kind {
	case inHouse:
		return eg.ChooseHouse(p, in.House)
	case inPlayCreature:
		_, err := eg.PlayCreature(p, in.Hand, in.Left)
		return err
	case inPlayArtifact:
		_, err := eg.PlayArtifact(p, in.Hand)
		return err
	case inPlayTactic:
		return eg.PlayTactic(p, in.Hand)
	case inPlayUpgrade:
		_, err := eg.PlayUpgrade(p, in.Hand)
		return err
	case inDiscard:
		return eg.DiscardFromHand(p, in.Hand)
	case inReap:
		return eg.Reap(p, in.Card)
	case inUnstun:
		return eg.Unstun(p, in.Card)
	case inUseAction:
		return eg.UseAction(p, in.Card)
	case inFight:
		return eg.Fight(p, in.Card, in.Card2)
	case inEndTurn:
		eg.EndPlayPhase(p)
		eg.StartTurn(1 - p)
		return nil
	default:
		return dispatchManual(eg, in, defs)
	}
}

// dispatchManual performs one recorded manual-mode root. The engine's Manual*
// methods perform no rule checks and record their own log entries, so replaying
// them reproduces both state and log exactly.
func dispatchManual(
	eg *engine.Game,
	in input,
	defs map[string]*engine.CardDefinition,
) error {
	switch in.Kind {
	case inSetManual:
		eg.SetManual(in.OK)
	case inManualMove:
		eg.ManualMove(in.ID, engine.ManualZone(in.Index))
	case inManualReady:
		eg.ManualSetExhausted(in.ID, false)
	case inManualExhaust:
		eg.ManualSetExhausted(in.ID, true)
	case inManualAttach:
		eg.ManualAttachUnder(in.Card, in.ID, in.Left)
	case inManualPlace:
		eg.ManualPlaceInPlay(in.ID, in.Index)
	case inManualDetach:
		eg.ManualDetachToHand(in.ID)
	case inManualAmber:
		eg.ManualAddAmber(in.Player, in.Delta)
	case inManualUnforge:
		eg.ManualUnforgeKey(in.Player)
	case inManualForgeColor:
		eg.ManualForgeKeyColor(in.Player, engine.KeyColor(in.Index))
	case inManualChains:
		eg.ManualAddChains(in.Player, in.Delta)
	case inManualHouse:
		eg.ManualSetActiveHouse(in.House)
	case inManualAddCard:
		def, ok := defs[in.Name]
		if !ok {
			return fmt.Errorf("web replay: manual card %q not in pool", in.Name)
		}
		eg.ManualAddCard(*def, in.Player)
	default:
		return fmt.Errorf("web replay: unknown root input %v", in.Kind)
	}
	return nil
}

// handOffEndedTurnOn mirrors game.handOffEndedTurn during replay: an Omega card
// ends the play phase without pairing StartTurn, so replay hands the turn on
// the same way live play does after every action.
func handOffEndedTurnOn(eg *engine.Game) {
	if eg.State.Phase == engine.PhaseEndOfTurn && eg.Winner() < 0 {
		eg.StartTurn(1 - eg.State.ActivePlayer)
	}
}
