// Package web is the browser client for the Vex card game. It renders an
// interactive two-player hotseat match on the Vex engine,
// compiled to WebAssembly via the go-app framework: the playtester sees
// the whole board and drives both sides.
package web

import (
	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vex/internal/cards"
	"github.com/dmikalova/vex/internal/engine"
	"github.com/dmikalova/vex/internal/match"
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

// upgradePlayChoice pre-answers the engine's "play this creature as a creature or
// an upgrade?" option prompt for a creature that may be played either way. The
// lifted card offers explicit Play creature / Play upgrade buttons, so the choice
// is made before the play runs and the engine's prompt answers itself rather than
// asking again in the sidebar.
type upgradePlayChoice int

const (
	choiceNone     upgradePlayChoice = iota // no creature-as-upgrade choice armed
	choiceCreature                          // play it as a creature (keep the flank step)
	choiceUpgrade                           // play it as an upgrade (skip the flank step)
)

// NewGame returns the root component for a fresh browser client session. The
// match itself is seeded on the client in OnMount.
func NewGame() app.Composer {
	return &game{
		selHand:     -1,
		zonesPlayer: -1,
		forgingKey:  -1,
		handSlot:    -1,
	}
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
	// legacy holds the LocalID of every Legacy card dealt this match (a card drawn
	// from an earlier set's pool), so its face shows the legacy emblem.
	legacy map[engine.LocalID]bool
	// rosters holds each player's static deck list — the generated deck kept as a
	// read-only roster so the deck-list popover shows the exact cards dealt (ADR
	// 0025), never the live draw order.
	rosters [2]match.Roster
	// deckOpen[p] is whether player p's deck-list popover is pinned open by a tap.
	// Each side toggles independently, so opening one leaves the other alone. Desktop
	// hover opens either without this; a touch tap toggles its own side and a tap
	// outside clears both (see onDeckToggle and installTips's outside-close).
	deckOpen [2]bool

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
	// swipeOnStrip marks that the touch began inside a horizontally-scrollable card
	// row (.card-strip). Scrolling such a row is a left/right drag that would
	// otherwise read as an open-swipe, so a strip-anchored swipe never opens the
	// sidebar; closing it (swipe right while open) is left alone.
	swipeOnStrip bool
	// isTouch marks that the client has fired a touch, so a log mention's synthetic
	// mouseenter (which a tap raises just before its click) is ignored — otherwise
	// the enter opens the preview and the click toggles it shut in the same tap,
	// making it take two taps. Once set, a single tap opens through onLogCardTap.
	isTouch bool

	// tipOverFunc/tipOutFunc raise the floating tip label on a mouse hover;
	// tipDownFunc/tipMoveFunc/tipUpFunc drive it from a finger dragged along the
	// player bar. All are released on dismount. tipTracking marks a press in
	// progress, which also suppresses the sidebar edge-swipe.
	tipOverFunc app.Func
	tipOutFunc  app.Func
	tipDownFunc app.Func
	tipMoveFunc app.Func
	tipUpFunc   app.Func
	tipTracking bool
	// selCursorFunc follows the pointer with the selection-badge marker (the small
	// icon at the cursor's corner showing what a click will do during a badge
	// preview). Released on dismount with the tip funcs.
	selCursorFunc app.Func

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
	// upgradeChoice pre-answers the engine's "as a creature or an upgrade?" prompt
	// for a creature the player is playing that may go down either way. It is armed
	// on the UI goroutine by the Play creature / Play upgrade buttons before the
	// play runs, then read by the webChooser off that goroutine; g.play resets it so
	// an ordinary play never inherits a stale choice.
	upgradeChoice upgradePlayChoice

	// engine chooser overlay
	choosing          bool
	chooserPrompt     string
	chooserCandidates []engine.LocalID
	// promptAsButtons offers a bounded card prompt's candidates as a short list of
	// action-bar buttons rather than opening the zone viewer — a "look at the top N
	// cards" pick (Navigator Ali, Lay of the Land) reads as a few named buttons
	// instead of a modal over the deck. Set for a small, mandatory out-of-play pick;
	// an unbounded pick (declinable — shuffle any number) keeps the viewer.
	promptAsButtons bool
	// chooserDeclinable marks a prompt the player may pass on — a "you may" or an
	// "up to N". It adds the Done button and lets Escape answer the prompt instead
	// of being swallowed.
	chooserDeclinable bool
	// chooserOrdering marks an ordering prompt (arranging several abilities' or
	// cards' resolution order). It adds the Auto-resolve button, which answers with
	// a random order instead of picking each in turn.
	chooserOrdering bool
	promptSource    string // card driving the current chooser/option prompt, if any
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
	// useTarget is the creature the current action last chose to use, so a "choose
	// how to use X" verb prompt (Universal Translator) can lift that creature and
	// put the reap/fight/action buttons on it — like an ordinary use — rather than
	// listing them in the sidebar. hasUseTarget guards it; it is cleared at the
	// start of each action and set by chooseCandidate when a creature is picked.
	useTarget    engine.LocalID
	hasUseTarget bool

	// selection badge preview (engine BadgeChooser): while an effect's choose loop
	// runs it previews the status each pick lands — the damage a Festering Touch
	// pick deals, the ward an Imperium "ward N" places. selBadge is the active badge
	// (its zero clears the preview); badgeTotals accumulates the amount landed on
	// each creature as it is picked, so a creature chosen twice shows the sum (a
	// zero-amount badge like a ward still records the key, drawing a numberless
	// icon); badgeClearing keeps the last badges on screen for the grow-and-fade
	// once the loop ends; badgeGen retires a stale fade timer when a new preview
	// begins mid-fade.
	selBadge      engine.SelectionBadge
	badgeTotals   map[engine.LocalID]int
	badgeClearing bool
	badgeGen      int

	// engine position chooser: placing a Deploy creature. The creature is lifted
	// while the prompt is up and its placement verbs sit on it (deployActions). The
	// placement is a two-step: first choose a side with the Deploy left / Deploy
	// right pair (which arms positionRight and sets positionSideChosen), then the
	// battleline lights up and a click on one of its creatures lands the new
	// creature on the chosen side of it. positionLine is the battleline being placed
	// into; positionRight is the armed side (false = left, true = right);
	// positionSideChosen gates the creature-picking step until a side is picked.
	choosingPosition   bool
	positionLine       []engine.LocalID
	positionRight      bool
	positionSideChosen bool

	// cancelling marks a manual-mode Cancel in flight: the current prompt (and any
	// that follow it as the effect drains) answers itself, and when the action
	// goroutine returns runAction rolls the whole action back to its start snapshot.
	cancelling bool

	// hostTargeting is set in manual mode after the player picks Graft or Place
	// under for the selected card: the board's in-play cards light up as hosts and
	// clicking one threads the selected card under it. hostFaceDown carries which
	// was picked — false grafts face up, true places under face down.
	hostTargeting bool
	hostFaceDown  bool

	// manualPlacing is set in manual mode after the player picks Put into play for
	// a selected hand creature: it reuses the Deploy placement picker (choosingPosition
	// with positionLine set to the owner's battleline) to drop the creature anywhere
	// in the line, but a clicked position calls ManualPlaceInPlay instead of answering
	// a real prompt goroutine.
	manualPlacing bool

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
	// The picker is also the client's card-name typeahead: pickerNaming marks it
	// open to answer a name-a-card prompt (Etan's Jar) rather than to add a card to
	// hand, which makes it non-dismissible and narrows its pool to the names the
	// prompt offered.
	pickerOpen    bool
	pickerNaming  bool
	pickerFocused bool
	pickerQuery   string
	// pickerCursor is the index of the highlighted row in the filtered picker
	// list; Enter adds it and Tab/arrows move it. It resets to 0 (the first row)
	// whenever the picker opens or the query filters the list.
	pickerCursor int
	allDefs      []engine.CardDefinition

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
	// toastSwipeStart marks whether the in-flight touch began on the toast, so the
	// global swipe listener flicks the toast away rather than moving the sidebar.
	toastSwipeStart bool

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

	// notice is a standing message that stays up until the condition behind it is
	// resolved, unlike status, which fades after 5s. It reports a fault the player
	// has to act on rather than a rejected click they can simply retry — the save
	// being lost, or a set name the deck generator does not know.
	notice string

	// redoLog holds the input segments undo has peeled off the command log, newest
	// last, so redo can splice one back on. Each segment is a root action plus the
	// chooser answers it raised. A new live action clears it.
	redoLog [][]input

	// inputs is the event-sourced command log (ADR 0039): every player input, in
	// order — each root action and each chooser answer. State and the typed log are
	// projections of it, so replaying inputs from a fresh deal reproduces the exact
	// match. record appends to it during live play; a replay drives the engine from
	// it instead of recording. rootMarks[i] is the index in inputs where the i-th
	// root action begins, so an undo truncates the log at a root boundary.
	inputs    []input
	rootMarks []int
	// replaying suppresses recording while the engine is driven from the log, so a
	// replay does not append the very inputs it is feeding back.
	replaying bool

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
	// it (a board card). hoverOverSidebar draws the mention preview over the sidebar
	// instead of beside it, when the window is too narrow to show a whole card to the
	// sidebar's left; hoverAtBottom anchors it to the bottom of the viewport when the
	// tapped log line sits in the top half, so the preview never covers its source.
	// hoverBack previews a plain card back rather than a face: what an opponent
	// sees when hovering a facedown Under-card they may not peek.
	hoverID          engine.LocalID
	hasHover         bool
	hoverBack        bool
	hoverDef         *engine.CardDefinition
	hoverInLog       bool
	hoverOverSidebar bool
	hoverAtBottom    bool

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
	// prevState is the state as it stood before the running root action, captured
	// by beginAction; computeFlashes diffs the resolved state against it to decide
	// what animates. prevValid guards it before the first action of a game.
	prevState  engine.GameState
	prevValid  bool
	poolFlash  [2]bool
	poolParity [2]bool
	keyFlash   [2]bool
	keyParity  [2]bool
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

	// inspecting marks the lift as a read-only peek raised by a long press or
	// right-click rather than a selection: it enlarges a card to read while a prompt
	// (a chooser, a house choice) owns the board, so it shows the card's face with
	// no action buttons and is dropped by the next tap or Escape.
	inspecting bool

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

// logMark records where one root action's log lines begin and whose turn it was.
type logMark struct {
	Start  int
	Player int
}

// persistKey names the local-storage slot holding the in-progress match, so a
// hot-reload of the wasm — or leaving the page and coming back — resumes it
// instead of dealing a new game.
const persistKey = "vex.match"

// snapshotVersion tags persisted state; bump it when an engine or command-log
// change makes older snapshots invalid so a stale one is flushed instead of
// replayed. Because a resume replays the command log rather than deserializing
// state, a change to how any recorded action resolves also dates every snapshot
// and counts as such a change.
const snapshotVersion = 21

// snapshot is the persisted match (ADR 0039): the seed and sets deterministically
// re-deal the same cards, and the ordered command log replays to the exact state
// and typed log. Everything else (state, choosers, deck houses, card index, UI
// phase) is a projection reconstructed from these, so nothing mutable is stored.
type snapshot struct {
	Version int
	Seed    int64
	// SetNames records each player's chosen deck-generation set, so a resume
	// re-deals the same decks the seed alone would not pin down once sets can
	// differ between players.
	SetNames [2]string
	// Inputs is the ordered command log: every root action and chooser answer,
	// which replaying from a fresh deal turns back into the exact match.
	Inputs []input
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

// isMaverick reports whether the dealt card with this LocalID is a Maverick (a
// card played out of its printed house), so its face can show the emblem.
func (g *game) isMaverick(id engine.LocalID) bool {
	return g.mavericks[id]
}

// isLegacy reports whether the dealt card with this LocalID is a Legacy card (a
// card drawn from an earlier set's pool), so its face can show the emblem.
func (g *game) isLegacy(id engine.LocalID) bool {
	return g.legacy[id]
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
	g.inspecting = false
	g.attacker = 0
}

// clearPrompts tears down the display state of any chooser prompt — creature,
// option, or position — that is currently up. A goroutine abandoned by a new deal
// cannot clear its own prompt (its dispatched teardown is stale-guarded so it
// cannot clobber the game that replaced it), so the deal that abandons it clears
// the stale prompt here. Left set, a stale choosingOption makes the next deal's
// await match the old flag instead of the fresh prompt.
func (g *game) clearPrompts() {
	g.choosing = false
	g.chooserDeclinable = false
	g.chooserOrdering = false
	g.chooserPrompt = ""
	g.chooserCandidates = nil
	g.promptAsButtons = false
	g.promptCursor, g.hasCursor = 0, false
	g.btnCursor, g.hasBtnCursor = 0, false
	g.choosingOption = false
	g.optionPrompt = ""
	g.optionLabels = nil
	g.choosingPosition = false
	g.positionLine = nil
	g.positionRight = false
	g.positionSideChosen = false
	g.promptSource = ""
}
