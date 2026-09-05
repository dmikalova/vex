// Package web is the browser client for the Vactrol card game. It renders an
// interactive two-player hotseat match on the Vactrol engine,
// compiled to WebAssembly via the go-app framework: the playtester sees
// the whole board and drives both sides.
package web

import (
	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/cards"
	"github.com/dmikalova/vactrol/internal/engine"
)

// phase is the interaction state of the client, distinct from the engine's own
// turn structure: it tracks what the UI is currently asking the player to do.
type phase int

const (
	phaseHouse       phase = iota // start of turn: pick the active house
	phaseMain                     // select cards and act
	phaseFlank                    // a creature is selected to play; pick a flank
	phaseFightTarget              // a fight was declared; pick an enemy to hit
	phaseOver                     // the game is finished
)

// selKind labels what the currently selected card is, so the action bar can
// offer the right verbs.
type selKind int

const (
	selNone         selKind = iota
	selHand                 // a card in the active player's hand
	selYourCreature         // a creature the active player controls
	selYourArtifact         // an artifact the active player controls
	selOther                // any other card (read-only: opponent's, an upgrade)
)

// NewGame returns the root component for a fresh browser client session. The
// match itself is seeded on the client in OnMount.
func NewGame() app.Composer {
	return &game{selHand: -1, zonesPlayer: -1, forgingKey: -1, handSlot: -1}
}

// game is the root component: it owns the live engine.Game and all UI state.
type game struct {
	app.Compo

	g          *engine.Game
	chooser    *webChooser
	seed       int64             // deal seed; persisted so a hot-reload can rebuild the match
	deckHouses [2][]engine.House // each player's three deck houses (house choices)
	// mavericks holds the LocalID of every Maverick card dealt this match (a card
	// played out of its printed house), so its face shows the maverick emblem.
	mavericks map[engine.LocalID]bool

	// dispatch schedules a mutation on the UI goroutine (captured from a Context).
	// It lets the background chooser update fields safely.
	dispatch func(func(app.Context))

	// keyFunc is the document-level keydown listener backing the keyboard
	// shortcuts; it is released on dismount.
	keyFunc app.Func

	// scrollFunc is the document-level scroll listener that keeps the lifted card
	// copy over its card; it is released on dismount.
	scrollFunc app.Func

	// touchStartFunc and touchEndFunc back the edge-swipe that reveals or hides the
	// sidebar; both are released on dismount. swipeStartX/swipeStartY/swipeTracking
	// hold the in-progress touch: a swipe that begins near the right edge and travels
	// left opens the sidebar, one that travels right closes it.
	touchStartFunc app.Func
	touchEndFunc   app.Func
	swipeStartX    float64
	swipeStartY    float64
	swipeTracking  bool

	phase phase
	busy  bool // an action goroutine is resolving; input is ignored

	// awaitingSetup shows the new-game set picker in the action bar over the
	// current board. It is on from a New game until both players have chosen a set
	// (or the same-sets shortcut is taken), at which point the match is dealt.
	awaitingSetup bool
	// setPick is the player (0 or 1) currently choosing a set in the action bar.
	setPick int
	// setNames holds each player's chosen deck-generation set for the current
	// match, persisted so a resume rebuilds the same decks.
	setNames [2]string
	// prevSetNames snapshots the running match's sets when the picker opens, so its
	// same-sets shortcut is stable while a player picks.
	prevSetNames [2]string

	// selection
	sel      engine.LocalID
	selKind  selKind
	selHand  int // hand index when selKind == selHand, else -1
	hasSel   bool
	attacker engine.LocalID // creature declared to fight, during phaseFightTarget

	// engine chooser overlay
	choosing          bool
	chooserPrompt     string
	chooserCandidates []engine.LocalID
	// chooserDeclinable marks a prompt the player may pass on — a "you may" or an
	// "up to N". It adds the Done button and lets Escape answer the prompt instead
	// of being swallowed.
	chooserDeclinable bool
	promptSource      string // card driving the current chooser/option prompt, if any
	// promptCursor is the candidate Tab has stepped to while a card prompt is up.
	// It draws as selected and is what Enter answers the prompt with, but only once
	// hasCursor says Tab has moved: LocalID 0 is a real card, so a zero cursor
	// cannot stand for "no cursor".
	promptCursor engine.LocalID
	hasCursor    bool

	// engine option chooser: a labeled multiple choice (e.g. take archives?)
	choosingOption bool
	optionPrompt   string
	optionLabels   []string

	// zonesPlayer, when >= 0, opens the out-of-play zone viewer (discard, archives,
	// and purge piles) for that player. -1 keeps the viewer closed.
	zonesPlayer int
	// promptZone names the zone row a prompt's candidates live in ("Discard"), when
	// the viewer was opened by that prompt rather than by the player. It is also the
	// row scrolled to, once, when the viewer appears.
	promptZone         string
	promptZoneScrolled bool

	// cursorScrolled is the DOM id the Tab cursor was last scrolled to, so the
	// scroll happens once when the cursor moves rather than on every render.
	cursorScrolled string

	// forgingKey, when >= 0, opens the manual-mode key-forge colour picker for that
	// player. -1 keeps it closed.
	forgingKey int

	// manual mode lives on the engine (g.g.Manual()); these back its UI:
	// the fuzzy card picker's open state and query, and the cached card pool it
	// searches. pickerFocused records that the search box has already been focused
	// for this opening, so typing is not interrupted on every later render.
	pickerOpen    bool
	pickerFocused bool
	pickerQuery   string
	allDefs       []engine.CardDefinition
	// manualAdds records, in registration order, the cards manual mode put into a
	// hand. The catalog is rebuilt from the seed, which knows nothing about them, so
	// a reload has to replay them or the saved state holds ids the catalog lacks.
	manualAdds []manualAdd

	// dragging is set while a hand card is being dragged, so the board shows as a
	// drop zone. It is cleared when the drag ends or the card is dropped.
	dragging bool

	// sidebarCollapsed hides the whole sidebar so the board area has more space; a
	// floating button reopens it.
	sidebarCollapsed bool

	// toastBubbles are the recent log bubbles surfaced over the board while the
	// sidebar (and its log) is hidden, so a minimized log still tells the player
	// what just happened. Each holds one root action's log lines and expires on its
	// own timer, oldest first, so only the latest bubbles linger rather than the
	// whole backlog. toastSeen is how far refreshToast has caught up, toastOpen
	// whether the last bubble still takes new lines. toastHover and toastPinned hold
	// every bubble open while the pointer is over the toast or after a click;
	// toastGen mints the per-bubble timer ids.
	toastBubbles []toastBubble
	toastSeen    int
	toastOpen    bool
	toastHover   bool
	toastPinned  bool
	toastGen     int

	// keysOpen shows the keyboard shortcut sheet (the ? key).
	keysOpen bool

	// menuOpen shows the sidebar's hamburger menu (undo/redo, manual mode, new
	// game, keyboard shortcuts).
	menuOpen bool

	// confirmEndTurn is armed when the player asks to end the turn while they could
	// still act; a second end-turn (e or the button) confirms. Any other action
	// (via beginAction), selecting another card, or an undo/redo disarms it.
	confirmEndTurn bool

	// confirmScrolled tracks whether the end-turn confirm's usable rows have already
	// been scrolled over once, so the strips scroll to reveal a jiggling card the
	// moment the confirm arms rather than every render while it stays armed.
	confirmScrolled bool

	// btnCursor is the prompt button Tab has stepped to — an index into whatever
	// buttons the current prompt offers (an option prompt's labels, the house
	// picker, the manual key-forge colours). It draws as selected and is what Enter
	// or Space presses, but only once hasBtnCursor says Tab has moved: index 0 is a
	// real button, so a zero cursor cannot stand for "no cursor".
	btnCursor    int
	hasBtnCursor bool

	// handSlot is where in the drawn hand the selected card sat when the running
	// action began, so the selection can move on to whatever fills that slot once
	// the card leaves. -1 means the action did not start from a card in hand.
	handSlot int

	// keyboardAction says the root action now running was invoked by a keyboard
	// shortcut rather than a click or drag, so advanceSelection only auto-selects
	// the next card when the player is driving from the keyboard — a mouse player
	// who just let go of a card would find the selection jumping to another one
	// surprising. Set by onKey right before it calls into a handler, and consumed
	// (and cleared) by advanceSelection once the action resolves.
	keyboardAction bool

	// defByName maps every card's name to its definition, for log-mention lookups.
	defByName map[string]*engine.CardDefinition

	status string // transient message (usually an action error)

	// undo/redo hold snapshots taken before each root action so the player can step
	// back and forward; each is a flat GameState value copy plus the log and marks.
	undo []undoEntry
	redo []undoEntry

	// logGroups marks where each root action's log lines begin (and whose turn), so
	// the log renders one bubble per action, tinted by player. The engine narrates
	// the turn's shape but frames only abilities, so where one player action stops
	// and the next starts is the client's own knowledge.
	logGroups []logMark

	// logScrollHeight is the log's scrollHeight as of the previous render, so the
	// next one can tell whether the player was reading at the bottom before new
	// lines were appended.
	logScrollHeight float64

	// hover preview: hoverID is a live board/hand card, held only when hasHover is
	// set — LocalID 0 is a real card (the first one dealt), so a zero id cannot mean
	// "nothing hovered". hoverDef is a printed card looked up from a log mention.
	// hoverInLog positions the preview left of the log (a mention) rather than over
	// it (a board card).
	hoverID    engine.LocalID
	hasHover   bool
	hoverDef   *engine.CardDefinition
	hoverInLog bool

	// flashes queues one-shot board animations computed after each action: which
	// cards took damage, gained Æmber, changed power counters, were stunned or
	// exhausted, or entered play, so their faces pulse. flashParity keeps a per-card
	// parity that flips on every flash so the CSS animation replays even on
	// back-to-back hits — go-app patches the persistent card element, and a CSS
	// animation only restarts when its animation-name changes. poolFlash/poolParity
	// and keyFlash/keyParity do the same for a player gaining pool Æmber or forging a
	// key. inPlayPrev is the set of in-play card ids after the previous action, used
	// to detect a card entering play.
	flashes     map[engine.LocalID]cardFlash
	flashParity map[engine.LocalID]bool
	inPlayPrev  map[engine.LocalID]bool
	poolFlash   [2]bool
	poolParity  [2]bool
	keyFlash    [2]bool
	keyParity   [2]bool
	// discardFlash/discardParity pulse a player's discard count when cards land
	// there.
	discardFlash  [2]bool
	discardParity [2]bool
	// flights are the cards that left the board this action, each shown flying into
	// the zone pill it landed in so a card going away is seen going somewhere.
	// flightParity restarts the animation the same way flashParity does.
	flights      []flight
	flightParity bool
	// fighters holds the attacker and defender of the fight the running action is
	// resolving, so both cards can clash when it lands. It is armed by the fight
	// handlers and consumed by computeFlashes; fighting says whether it is armed,
	// since 0 is a real card id.
	fighters [2]engine.LocalID
	fighting bool
	// reapID/reaping arms the yellow reap animation for one card the same way
	// fighters/fighting arms the fight clash; consumed by computeFlashes.
	reapID  engine.LocalID
	reaping bool
	// actID/acting arms the green action-ability animation (an artifact's action, or
	// a creature's own) for one card, the same way reapID/reaping does for reap.
	actID  engine.LocalID
	acting bool
	// takeoff is where a card being played sat in hand when the player let go of
	// it, so the board card it becomes can be animated from that spot rather than
	// appearing in place. Armed at play time and spent by the next render.
	takeoff   cardRect
	takeoffID engine.LocalID
	takingOff bool

	// focusRect is where the selected card sits on screen, so the lifted copy of it
	// (cardFocus) can be laid out over its slot and grown from there. focusID is the
	// card that was measured and focusParity flips whenever the lift moves to another
	// one, so the grow animation replays; hasFocus says a measurement has been taken
	// at all, since a card really can sit at the origin. focusViewW/focusViewH are
	// the window size as of that measurement, which is what the copy is kept inside
	// of and which edge it is anchored to.
	focusRect              cardRect
	focusID                engine.LocalID
	hasFocus               bool
	focusParity            bool
	focusViewW, focusViewH float64

	// focusExit keeps the lifted copy on screen after a card is deselected so it
	// shrinks back to its slot — the grow-in played backwards — rather than blinking
	// away. focusShown is the copy's placement snapshotted from the last live
	// measurement, kept because the selection focusBox reads from is already gone by
	// the time the exit plays; focusExitGen tags the clear timer so a fresh
	// selection cancels a pending exit.
	focusExit    bool
	focusShown   focusSnapshot
	focusExitGen int

	// statusGen tags the current status message so a scheduled auto-clear only
	// clears the message it was armed for, not a newer one.
	statusGen int
}

// cardFlash is a queued one-shot face animation for a card after an action. A
// single odd bit alternates all of the card's pulses at once (a card flashes at
// most once per action) so the CSS animations restart on repeats.
type cardFlash struct {
	damage  bool // took new damage
	amber   bool // gained Æmber on the card
	power   bool // power counters changed
	exhaust bool // became exhausted
	stun    bool // became stunned
	enter   bool // entered play this action
	fight   bool // was the attacker or defender of a fight
	reap    bool // reaped
	act     bool // used an action ability (an artifact's, or a creature's own)
	odd     bool // alternates each flash to restart the CSS animation
}

// flight is a card that just left the board, on its way into the zone pill it
// landed in: the owning player's pill and that pill's zone icon name.
type flight struct {
	id     engine.LocalID
	player int
	zone   string
}

// cardRect is a card's place on screen in viewport pixels, as the browser
// reports it.
type cardRect struct {
	x, y, w, h float64
}

// undoEntry is a restorable snapshot: the flat GameState (a pure value copy), the
// log, and the log group marks at that moment.
type undoEntry struct {
	state  engine.GameState
	log    []engine.Record
	groups []logMark
}

// logMark records where one root action's log lines begin and whose turn it was.
type logMark struct {
	Start  int
	Player int
}

// persistKey names the local-storage slot holding the in-progress match, so a
// hot-reload of the wasm — or leaving the page and coming back — resumes it
// instead of dealing a new game.
const persistKey = "vactrol.match"

// snapshotVersion tags persisted state; bump it when an engine change makes older
// snapshots invalid so a stale one is flushed instead of restored. A log entry is
// saved as the prose it was narrated with, so rewording an entry dates every
// snapshot holding the old wording and counts as such a change.
const snapshotVersion = 10

// snapshot is the persisted match. The seed deterministically rebuilds the
// catalog and card ids; the flat GameState carries everything mutable. All other
// state (choosers, deck houses, card index, UI phase) is reconstructed from these.
type snapshot struct {
	Version int
	Seed    int64
	State   engine.GameState
	// SetNames records each player's chosen deck-generation set, so a resume
	// rebuilds the same decks the seed alone would not pin down once sets can
	// differ between players.
	SetNames [2]string
	// Log persists the game log so a hot-reload does not lose the history. A typed
	// log entry does not survive JSON (it carries an interface, like a
	// CardDefinition's effect nodes), so each entry is saved as the text it was
	// narrated with plus the divider it drew. Groups persists the per-action
	// bubbling alongside it.
	Log    []savedLine
	Groups []logMark
	// Manual replays the cards manual mode added, in registration order, so the
	// rebuilt catalog holds the same ids the state refers to.
	Manual []manualAdd
	// UI carries the view across the reload with the match.
	UI savedUI
}

// savedUI is the view state a reload should put back rather than throw away:
// whether the sidebar was hidden, and which overlay was open over the board.
type savedUI struct {
	SidebarCollapsed bool
	// ZonesPlayer is the zone viewer's player, or -1 for closed, as the field it
	// restores means it.
	ZonesPlayer int
	KeysOpen    bool
	PickerOpen  bool
}

// savedLine is one persisted log entry: the attribution it was recorded under,
// the text it was narrated with, and — for a turn or phase header — the rule it
// drew and the player it announced, which the entry's type no longer carries once
// it has been through JSON.
type savedLine struct {
	Frame  engine.Frame
	Text   string
	Rule   logRule
	Player int
	// Standing holds an end-of-turn PlayerStanding's data, so the coloured key
	// tally survives the reload the typed entry does not. Its fields are all
	// JSON-safe (ints and a KeyColor slice), so the real entry is rebuilt on
	// resume rather than flattened to the plain "N keys" text. Nil on every other
	// line.
	Standing *savedStanding
}

// savedStanding persists a PlayerStanding so a resumed match redraws the same
// coloured key tally it showed before the reload instead of a plain count.
type savedStanding struct {
	Player int
	Aember int
	Keys   []engine.KeyColor
}

// manualAdd is one card manual mode put into a hand, named rather than embedded:
// a CardDefinition carries effect nodes, which do not survive JSON.
type manualAdd struct {
	Name   string
	Player int
}

// isMaverick reports whether the dealt card with this LocalID is a Maverick (a
// card played out of its printed house), so its face can show the emblem.
func (g *game) isMaverick(id engine.LocalID) bool {
	return g.mavericks[id]
}

// cardsByName indexes every registered card's definition by its display name so
// the log can link the cards it mentions. The pointers reference the returned
// slice's backing array, which the map keeps alive.
func cardsByName() map[string]*engine.CardDefinition {
	all := cards.All()
	m := make(map[string]*engine.CardDefinition, len(all))
	for i := range all {
		m[all[i].Name] = &all[i]
	}
	return m
}

// active returns the player whose turn it is; the client always renders from
// this player's perspective (their board and hand at the bottom).
func (g *game) active() int { return g.g.State.ActivePlayer }

func (g *game) clearSelection() {
	g.sel = 0
	g.selKind = selNone
	g.selHand = -1
	g.hasSel = false
	g.attacker = 0
}
