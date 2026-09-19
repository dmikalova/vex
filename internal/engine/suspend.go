package engine

import "runtime"

// This file holds the engine's suspendable step function (ADR 0040): the seam
// that lets resolution YIELD for a decision instead of PULLING one from a blocking
// Chooser. It is additive — the synchronous Chooser path stays for MCTS rollouts
// and tests, which drive the engine over FastCopy copies and must not allocate a
// goroutine per step. The suspendable path is for interactive and networked play,
// where an answer may arrive much later and must be recordable as a Command.
//
// The mechanism reuses every existing effect: a suspendChooser implements the
// Chooser capabilities, but instead of computing an answer it sends a Request out
// and blocks for a Command back. A Stepper runs one action on a goroutine and
// bridges the two channels, so Advance feeds the parked resolution its next
// answer. This coroutine is the pragmatic realization of the ADR's target
// `Advance(state, command)`; threading the state as a pure value is a later step.

// CommandKind tags which sort of answer a Command carries.
type CommandKind uint8

const (
	// CommandPickCard names one card (a creature or card pick). Card holds its id.
	CommandPickCard CommandKind = iota
	// CommandDecline passes on an optional pick.
	CommandDecline
	// CommandOption selects a labeled option by Index.
	CommandOption
	// CommandPosition selects a battleline insertion point by Index (0 the left
	// flank, len(line) the right flank).
	CommandPosition
	// CommandReaction selects which pending reaction resolves next by Index.
	CommandReaction

	// The kinds below are ROOT ACTIONS — the other half of ADR 0039's "a root
	// action OR one answer to a choice". A root Command is a play the active player
	// initiates rather than an answer the engine pulled mid-resolution; ApplyAction
	// performs one and LegalActions enumerates the legal set (command_action.go).
	// They are appended after the answer kinds so the answer kinds' persisted values
	// do not shift. Only the field each names is meaningful.

	// CommandChooseHouse names the turn's active house in House.
	CommandChooseHouse
	// CommandPlayCreature plays the hand creature at Hand onto a flank (Left).
	CommandPlayCreature
	// CommandPlayArtifact plays the hand artifact at Hand.
	CommandPlayArtifact
	// CommandPlayTactic plays the hand Tactic at Hand.
	CommandPlayTactic
	// CommandPlayUpgrade plays the hand upgrade at Hand.
	CommandPlayUpgrade
	// CommandDiscardFromHand discards the active-house hand card at Hand.
	CommandDiscardFromHand
	// CommandReap reaps with the creature Card.
	CommandReap
	// CommandUnstun spends the stunned creature Card's use to shed the stun.
	CommandUnstun
	// CommandUseAction uses the "Action:" ability of the creature or artifact Card.
	CommandUseAction
	// CommandFight fights the enemy creature Card2 with the creature Card.
	CommandFight
	// CommandEndTurn ends the active player's play phase and hands off the turn.
	CommandEndTurn

	// The kinds below are MANUAL/DEBUG roots — a playtester's force-edit that
	// deliberately bypasses the rules (ADR 0039: "a saved match can contain
	// force-edits"). Each maps to one engine Manual* method through ApplyManual
	// (command_manual.go); they are recorded and replayed like any other command so
	// a match holding force-edits still rebuilds from its log. They are NOT root
	// actions: LegalActions never offers one, so the sim never drives manual mode
	// (TestLegalActionsNeverOffersManualKinds). Appended last so no earlier kind's
	// persisted value shifts. Only the field each names is meaningful.

	// CommandSetManual toggles manual mode; Left is whether to turn it on.
	CommandSetManual
	// CommandManualMove moves the card Card to the resting zone Index (a ManualZone).
	CommandManualMove
	// CommandManualReady clears the exhausted flag on the card Card.
	CommandManualReady
	// CommandManualExhaust sets the exhausted flag on the card Card.
	CommandManualExhaust
	// CommandManualAttach threads the card Card2 under the host Card, face down when
	// Left.
	CommandManualAttach
	// CommandManualPlace drops the card Card into play at battleline position Index.
	CommandManualPlace
	// CommandManualDetach sends the attached card Card to its owner's hand.
	CommandManualDetach
	// CommandManualAmber adjusts Player's Æmber pool by the signed Delta.
	CommandManualAmber
	// CommandManualUnforge removes Player's most recently forged key.
	CommandManualUnforge
	// CommandManualForgeColor forges a key of colour Index (a KeyColor) for Player.
	CommandManualForgeColor
	// CommandManualChains adjusts Player's chain count by the signed Delta.
	CommandManualChains
	// CommandManualHouse sets the active player's active house to House.
	CommandManualHouse
	// CommandManualAddCard registers the card named Name into Player's hand.
	CommandManualAddCard
)

// Command is one player input crossing the engine boundary — an answer to a
// Request or a root action (ADR 0039). It is a flat comparable value so a match's
// inputs can be recorded, replayed, and compared. Only the field its Kind names is
// meaningful; the root-action fields stay zero for the answer kinds.
type Command struct {
	Kind  CommandKind
	Card  LocalID
	Index int

	// Root-action fields (see the root CommandKinds above).
	House House   // CommandChooseHouse: the active house to name.
	Hand  int     // the Play*/Discard kinds: the hand index to act on.
	Left  bool    // CommandPlayCreature: place on the left flank.
	Card2 LocalID // CommandFight: the enemy creature Card fights.

	// Manual/debug fields (see the manual CommandKinds above). Left doubles as the
	// on/face-down flag for CommandSetManual and CommandManualAttach; Card/Card2/
	// Index/House carry the manual target the same way they carry a root action's.
	Player int    // the manual kinds that name a player (Æmber, chains, keys, add).
	Delta  int    // CommandManualAmber / CommandManualChains: the signed delta.
	Name   string // CommandManualAddCard: the card definition's name.
}

// RequestKind tags what sort of decision a Request marks.
type RequestKind uint8

const (
	// RequestPickCard asks for one card from Cards.
	RequestPickCard RequestKind = iota
	// RequestPickCardOrDecline asks for one card from Cards, or a decline.
	RequestPickCardOrDecline
	// RequestOption asks for one of the labeled Options.
	RequestOption
	// RequestPosition asks where in the battleline (Cards) a creature enters.
	RequestPosition
	// RequestReaction asks which pending Reactions resolves next.
	RequestReaction

	// RequestAction asks the active player for their next ROOT action, chosen from
	// the legal set the engine gathered into Actions (ADR 0039). It is appended
	// after the mid-resolution kinds so their persisted values do not shift. Where
	// they each yield mid-resolution to answer one choice, this one yields BETWEEN
	// actions, so the canonical turn loop can ask what the player does next.
	RequestAction
)

// Request marks a decision point: which player owes a decision and in what
// context. It is the value the suspendable engine yields in place of pulling a
// Chooser. It carries the candidate context the engine computed at the suspension
// point (Cards/Options/Reactions) so the driving side can render and validate the
// answer; deriving the legal set purely from GameState is a later refinement.
type Request struct {
	Player    int
	Kind      RequestKind
	Source    string
	Prompt    string
	Cards     []LocalID
	Options   []string
	Reactions []OrderableReaction
	// Actions is the legal root-action set a RequestAction offers (ADR 0039). It is
	// the whole legal set, so LegalCommands returns it verbatim.
	Actions []Command
}

// StepInfo reports what an Advance step observed — whether it crossed an
// information barrier (ADR 0039): stepped the PRNG or revealed a hidden zone. The
// undo rule reads this to know which steps are freely reversible. PRNG stepping is
// detected here; hidden-reveal detection is a later refinement.
type StepInfo struct {
	CrossedBarrier bool
}

// LegalCommands returns every answer this request accepts, derived from its
// candidate context. A holder re-derives rather than trusting a transmitted list.
func (r Request) LegalCommands() []Command {
	switch r.Kind {
	case RequestPickCard:
		return pickCommands(r.Cards, false)
	case RequestPickCardOrDecline:
		return pickCommands(r.Cards, true)
	case RequestOption:
		return indexCommands(CommandOption, len(r.Options))
	case RequestPosition:
		return indexCommands(CommandPosition, len(r.Cards)+1)
	case RequestReaction:
		return indexCommands(CommandReaction, len(r.Reactions))
	case RequestAction:
		return r.Actions
	default:
		return nil
	}
}

// pickCommands builds a card-pick command per candidate, plus a decline when the
// prompt allows passing.
func pickCommands(cards []LocalID, declinable bool) []Command {
	cmds := make([]Command, 0, len(cards)+1)
	for _, id := range cards {
		cmds = append(cmds, Command{Kind: CommandPickCard, Card: id})
	}
	if declinable {
		cmds = append(cmds, Command{Kind: CommandDecline})
	}
	return cmds
}

// indexCommands builds one command of kind per index in [0, n).
func indexCommands(kind CommandKind, n int) []Command {
	cmds := make([]Command, n)
	for i := range cmds {
		cmds[i] = Command{Kind: kind, Index: i}
	}
	return cmds
}

// IsLegal reports whether cmd is a valid answer to this request.
func (r Request) IsLegal(cmd Command) bool {
	for _, c := range r.LegalCommands() {
		if c == cmd {
			return true
		}
	}
	return false
}

// suspendChooser is the adapter that turns a pulled choice into a yielded Request:
// each capability sends a Request for its player and blocks for the Command back.
// One instance per player shares the Stepper's channels; because resolution runs
// single-threaded on the Stepper's goroutine, only one call is ever in flight.
type suspendChooser struct {
	player   int
	requests chan<- Request
	commands <-chan Command
	// cancel is closed by Stepper.Close to release this goroutine when its Stepper
	// is abandoned mid-action (an undo or replay deals a fresh game); yield unwinds
	// via runtime.Goexit rather than leaking parked on an answer that never comes.
	cancel <-chan struct{}
}

// A suspendChooser installs every optional capability, including the turn loop's
// ActionChooser, so the driving side is offered every kind of decision.
var _ ActionChooser = (*suspendChooser)(nil)

// yield sends req and blocks until the driving side answers with a Command. If the
// Stepper is closed while yield is parked, cancel fires and the goroutine unwinds
// (Goexit runs its defers), so an abandoned action leaves no goroutine behind.
func (c *suspendChooser) yield(req Request) Command {
	select {
	case c.requests <- req:
	case <-c.cancel:
		runtime.Goexit()
	}
	select {
	case cmd := <-c.commands:
		return cmd
	case <-c.cancel:
		runtime.Goexit()
		return Command{} // unreachable: Goexit does not return
	}
}

// ChooseCreature yields a card pick and returns the answered card.
func (c *suspendChooser) ChooseCreature(source, prompt string, cands []LocalID) (LocalID, bool) {
	cmd := c.yield(Request{
		Player: c.player, Kind: RequestPickCard, Source: source, Prompt: prompt, Cards: cands,
	})
	if cmd.Kind == CommandDecline {
		return 0, false
	}
	return cmd.Card, true
}

// ChooseCardOrDecline yields a declinable card pick.
func (c *suspendChooser) ChooseCardOrDecline(
	source, prompt string, cands []LocalID,
) (LocalID, bool) {
	cmd := c.yield(Request{
		Player: c.player, Kind: RequestPickCardOrDecline,
		Source: source, Prompt: prompt, Cards: cands,
	})
	if cmd.Kind == CommandDecline {
		return 0, false
	}
	return cmd.Card, true
}

// ChooseOption yields a labeled-option pick and returns the chosen index.
func (c *suspendChooser) ChooseOption(source, prompt string, options []string) int {
	return c.yield(Request{
		Player: c.player, Kind: RequestOption, Source: source, Prompt: prompt, Options: options,
	}).Index
}

// ChoosePosition yields a battleline placement and returns the chosen gap index.
func (c *suspendChooser) ChoosePosition(source, prompt string, line []LocalID) int {
	return c.yield(Request{
		Player: c.player, Kind: RequestPosition, Source: source, Prompt: prompt, Cards: line,
	}).Index
}

// ChooseReaction yields a reaction-ordering pick and returns the chosen index.
func (c *suspendChooser) ChooseReaction(prompt string, reactions []OrderableReaction) int {
	return c.yield(Request{
		Player: c.player, Kind: RequestReaction, Prompt: prompt, Reactions: reactions,
	}).Index
}

// ChooseAction yields the active player's legal root-action set and returns the
// root Command they answer with. It is the turn loop's suspension point (ADR 0039):
// where the other capabilities yield mid-resolution to answer one choice, this one
// yields between actions to ask what the player does next. The legal set travels in
// the Request so the driving side can render and validate the answer, which is a
// root Command ApplyAction performs.
func (c *suspendChooser) ChooseAction(actions []Command) Command {
	return c.yield(Request{Player: c.player, Kind: RequestAction, Actions: actions})
}

// Stepper runs one engine action to completion on a goroutine, suspending it at
// each decision so a driver can supply answers one Command at a time. It is the
// coroutine behind the suspendable step function: Start yields the first Request
// (or reports the action already finished), and Advance answers the current
// Request and yields the next.
type Stepper struct {
	g        *Game
	requests chan Request
	commands chan Command
	cancel   chan struct{}
	done     bool
}

// NewStepper installs suspending choosers on g and launches action on a goroutine.
// The goroutine runs until it needs a decision (blocking in a chooser) or returns;
// call Start to receive the first Request or completion.
func NewStepper(g *Game, action func(*Game)) *Stepper {
	s := &Stepper{
		g:        g,
		requests: make(chan Request),
		commands: make(chan Command),
		cancel:   make(chan struct{}),
	}
	g.SetChooser(0, &suspendChooser{
		player: 0, requests: s.requests, commands: s.commands, cancel: s.cancel,
	})
	g.SetChooser(1, &suspendChooser{
		player: 1, requests: s.requests, commands: s.commands, cancel: s.cancel,
	})
	go func() {
		action(g)
		close(s.requests)
	}()
	return s
}

// Close releases the action goroutine if it is still parked at a decision, so a
// Stepper abandoned mid-action (an undo or replay that deals a fresh game) does
// not leak the goroutine and the Game snapshot it captured. It is idempotent and a
// no-op once the action has finished on its own (TestStepperCloseReleasesGoroutine).
func (s *Stepper) Close() {
	if s.done {
		return
	}
	s.done = true
	close(s.cancel)
}

// Start returns the first Request the action reaches, or done=true if the action
// finished without needing any decision.
func (s *Stepper) Start() (Request, bool) {
	req, ok := <-s.requests
	s.done = !ok
	return req, s.done
}

// Advance answers the current Request with cmd and returns the next Request, or
// done=true when the action has finished. StepInfo reports whether resolving cmd
// crossed an information barrier. Calling Advance after done is a no-op that
// reports done again.
func (s *Stepper) Advance(cmd Command) (Request, bool, StepInfo) {
	if s.done {
		return Request{}, true, StepInfo{}
	}
	before := s.g.State.PRNG.State
	s.commands <- cmd
	req, ok := <-s.requests
	s.done = !ok
	return req, s.done, StepInfo{CrossedBarrier: s.g.State.PRNG.State != before}
}
