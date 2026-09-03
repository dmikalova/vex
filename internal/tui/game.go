package tui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dmikalova/vactrol/internal/engine"
	"github.com/dmikalova/vactrol/internal/match"
)

// ---- messages ----

// chooseRequestMsg is posted by the engine (via the chooser bridge) when an
// effect needs the active player to pick a creature or a labeled option.
type chooseRequestMsg struct {
	prompt     string
	candidates []engine.LocalID
	options    []string // set when the engine asked for a labeled option choice
	isOption   bool
}

// actionDoneMsg reports that a play/use/fight action goroutine has finished.
type actionDoneMsg struct{ err error }

// beginDoneMsg reports that a StartTurn goroutine has finished (keys forged).
type beginDoneMsg struct{}

// ---- chooser bridge ----

// chooseReply carries the UI's answer to a chooser request: a chosen id, or
// ok=false when the player cancels the selection.
type chooseReply struct {
	id engine.LocalID
	ok bool
}

// teaChooser adapts the engine's synchronous Chooser to bubbletea's event loop:
// it posts a chooseRequestMsg and blocks until the UI replies with the id.
type teaChooser struct {
	snd         *sender
	reply       chan chooseReply
	optionReply chan int
}

func (c *teaChooser) ChooseCreature(
	_, prompt string,
	candidates []engine.LocalID,
) (engine.LocalID, bool) {
	if len(candidates) == 0 {
		return 0, false
	}
	c.snd.send(chooseRequestMsg{prompt: prompt, candidates: candidates})
	r := <-c.reply
	return r.id, r.ok
}

// ChooseOption posts a labeled option choice (e.g. which key colour to forge) and
// blocks until the UI picks one. Implementing it means the TUI prompts rather
// than silently taking the first option.
func (c *teaChooser) ChooseOption(_, prompt string, options []string) int {
	if len(options) <= 1 {
		return 0
	}
	c.snd.send(chooseRequestMsg{prompt: prompt, options: options, isOption: true})
	return <-c.optionReply
}

// ---- model ----

type gamePhase int

const (
	phaseHouse        gamePhase = iota // choosing a house (board stays visible)
	phaseMain                          // navigating the board/hand and acting
	phaseFlank                         // choosing a flank for a creature being played
	phaseCreatureMenu                  // reap/fight/action menu for one of your creatures
	phaseChoose                        // engine asked to pick a creature (on the board)
	phaseChooseOption                  // engine asked to pick a labeled option (e.g. key colour)
	phaseTargetFight                   // picking an enemy to fight (on the board)
	phaseDiscard                       // viewing the discard piles
	phaseBusy                          // an action is resolving
	phaseOver                          // game finished
	phaseConfirmExit                   // confirming leaving the game for the menu
)

// slotKind labels an entry in the flat, navigable list of on-screen cards.
type slotKind int

const (
	slotOppCreature slotKind = iota
	slotOppArtifact
	slotYourCreature
	slotYourArtifact
	slotUpgrade // an upgrade attached to a creature (read-only)
	slotHand
)

// slot is one navigable card: a creature/artifact/upgrade in play or a card in hand.
type slot struct {
	id      engine.LocalID
	kind    slotKind
	handIdx int // real hand index for slotHand, else -1
}

// creatureAction is one option in the reap/fight/action menu for a creature.
type creatureAction int

const (
	caReap creatureAction = iota
	caFight
	caAction
)

func (ca creatureAction) label() string {
	switch ca {
	case caReap:
		return "Reap"
	case caFight:
		return "Fight"
	case caAction:
		return "Action"
	default:
		return "?"
	}
}

type gameModel struct {
	snd     *sender
	g       *engine.Game
	chooser *teaChooser

	player     int
	phase      gamePhase
	deckHouses [2][]engine.House // each player's three deck houses (house choices)

	cur      int // index into slots() during phaseMain
	houseCur int

	playHandPos int // hand position to restore the cursor to after a play (-1 = none)
	pendingHand int // hand index of a creature awaiting its flank choice

	choosePrompt  string
	chooseCands   []engine.LocalID
	chooseOptions []string
	chooseCur     int

	fightAttacker engine.LocalID
	fightCands    []engine.LocalID
	fightCur      int

	actionCreature engine.LocalID
	actionMenu     []creatureAction
	actionMenuCur  int

	discardCur int

	suppressErr bool // discard the error of a just-cancelled action

	exitReturnPhase gamePhase // phase to restore if the player cancels an exit

	status string
	width  int
	height int
}

func newGameModel(snd *sender, w, h int) gameModel {
	seed := time.Now().UnixNano()
	g, houses := match.New("Player 1", "Player 2", seed)
	ch := &teaChooser{snd: snd, reply: make(chan chooseReply, 1), optionReply: make(chan int, 1)}
	g.SetChooser(0, ch)
	g.SetChooser(1, ch)
	g.SetPlayerHouses(0, houses[0])
	g.SetPlayerHouses(1, houses[1])
	g.StartTurn(0) // first turn: no Æmber yet, so forging never triggers here
	m := gameModel{
		snd: snd, g: g, chooser: ch,
		player: 0, deckHouses: houses, playHandPos: -1,
		width: w, height: h,
	}
	return m.enterHousePhase()
}

func (m gameModel) resize(w, h int) gameModel { m.width, m.height = w, h; return m }

func (m gameModel) Init() tea.Cmd { return nil }

// ---- slots: the flat, navigable list of on-screen cards ----

func (m gameModel) slots() []slot {
	opp := 1 - m.player
	var s []slot
	for _, id := range m.g.Battleline(opp) {
		s = append(s, slot{id, slotOppCreature, -1})
		s = m.appendUpgrades(s, id)
	}
	for _, id := range m.g.Artifacts(opp) {
		s = append(s, slot{id, slotOppArtifact, -1})
	}
	for _, id := range m.g.Battleline(m.player) {
		s = append(s, slot{id, slotYourCreature, -1})
		s = m.appendUpgrades(s, id)
	}
	for _, id := range m.g.Artifacts(m.player) {
		s = append(s, slot{id, slotYourArtifact, -1})
	}
	return append(s, m.handSlots()...)
}

// appendUpgrades adds a creature's attached upgrades as read-only slots so they
// can be hovered to read their text.
func (m gameModel) appendUpgrades(s []slot, host engine.LocalID) []slot {
	for _, up := range m.g.Upgrades(host) {
		s = append(s, slot{up, slotUpgrade, -1})
	}
	return s
}

// handSlots returns the player's hand as slots sorted by house then name; the
// real hand index is retained so plays still target the right card.
func (m gameModel) handSlots() []slot {
	hand := m.g.Hand(m.player)
	out := make([]slot, len(hand))
	for i, id := range hand {
		out[i] = slot{id: id, kind: slotHand, handIdx: i}
	}
	sort.SliceStable(out, func(a, b int) bool {
		da, db := m.g.Def(out[a].id), m.g.Def(out[b].id)
		if da.House != db.House {
			return da.House.String() < db.House.String()
		}
		return da.Name < db.Name
	})
	return out
}

func (m gameModel) curSlot() (slot, bool) {
	s := m.slots()
	if m.cur < 0 || m.cur >= len(s) {
		return slot{}, false
	}
	return s[m.cur], true
}

func (m gameModel) firstHandSlot() int {
	s := m.slots()
	for i, sl := range s {
		if sl.kind == slotHand {
			return i
		}
	}
	return clamp(0, 0, len(s)-1)
}

// ---- update ----

func (m gameModel) Update(msg tea.Msg) (gameModel, tea.Cmd) {
	switch msg := msg.(type) {
	case chooseRequestMsg:
		m.choosePrompt = msg.prompt
		m.chooseCur = 0
		if msg.isOption {
			m.chooseOptions = msg.options
			m.phase = phaseChooseOption
		} else {
			m.chooseCands = msg.candidates
			m.phase = phaseChoose
		}
		return m, nil
	case actionDoneMsg:
		if m.suppressErr {
			m.suppressErr = false
		} else if msg.err != nil {
			m.status = msg.err.Error()
		}
		m.phase = phaseMain
		if m.playHandPos >= 0 {
			m.cur = clamp(m.firstHandSlot()+m.playHandPos, 0, len(m.slots())-1)
			m.playHandPos = -1
		} else {
			m.cur = clamp(m.cur, 0, len(m.slots())-1)
		}
		return m, nil
	case beginDoneMsg:
		if m.g.Winner() >= 0 {
			m.phase = phaseOver
			return m, nil
		}
		return m.enterHousePhase(), nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m gameModel) handleKey(k tea.KeyMsg) (gameModel, tea.Cmd) {
	switch m.phase {
	case phaseHouse:
		return m.handleHouseKey(k)
	case phaseMain:
		return m.handleMainKey(k)
	case phaseFlank:
		return m.handleFlankKey(k)
	case phaseCreatureMenu:
		return m.handleCreatureMenuKey(k)
	case phaseChoose:
		return m.handleChooseKey(k)
	case phaseChooseOption:
		return m.handleChooseOptionKey(k)
	case phaseTargetFight:
		return m.handleFightKey(k)
	case phaseDiscard:
		return m.handleDiscardKey(k)
	case phaseConfirmExit:
		return m.handleConfirmExitKey(k)
	case phaseOver:
		switch k.String() {
		case "enter", " ", "esc", "q":
			return m, gotoScreen(screenMenu)
		}
	case phaseBusy:
		// ignore input while an action resolves
	}
	return m, nil
}

func (m gameModel) handleHouseKey(k tea.KeyMsg) (gameModel, tea.Cmd) {
	switch k.String() {
	case "up", "k":
		m.houseCur = wrap(m.houseCur, -1, len(m.deckHouses[m.player]))
	case "down", "l":
		m.houseCur = wrap(m.houseCur, 1, len(m.deckHouses[m.player]))
	case "enter", " ":
		if len(m.deckHouses[m.player]) > 0 {
			_ = m.g.ChooseHouse(m.player, m.deckHouses[m.player][m.houseCur])
			m.phase = phaseMain
			m.cur = m.firstHandSlot()
			m.status = ""
		}
	case "e":
		return m.endTurn()
	case "esc":
		return m.promptExit(phaseHouse)
	}
	return m, nil
}

func (m gameModel) handleMainKey(k tea.KeyMsg) (gameModel, tea.Cmd) {
	switch k.String() {
	case "up", "k":
		m.cur = wrap(m.cur, -1, len(m.slots()))
		m.status = ""
	case "down", "l":
		m.cur = wrap(m.cur, 1, len(m.slots()))
		m.status = ""
	case "e":
		return m.endTurn()
	case "d":
		if sl, ok := m.curSlot(); ok && sl.kind == slotHand {
			return m.discardHand(sl)
		}
	case "v":
		m.discardCur = 0
		m.phase = phaseDiscard
		m.status = ""
	case "esc":
		return m.promptExit(phaseMain)
	case "enter", "p":
		if sl, ok := m.curSlot(); ok {
			switch sl.kind {
			case slotHand:
				return m.playSlot(sl)
			case slotYourCreature:
				return m.openCreatureMenu(sl.id)
			case slotYourArtifact:
				return m.useArtifact(sl.id)
			}
		}
	case "r":
		if sl, ok := m.curSlot(); ok && sl.kind == slotYourCreature {
			g, player, id := m.g, m.player, sl.id
			return m.runAction(func() error { return g.Reap(player, id) })
		}
	case "a":
		if sl, ok := m.curSlot(); ok {
			switch sl.kind {
			case slotYourCreature:
				g, player, id := m.g, m.player, sl.id
				return m.runAction(func() error { return g.UseAction(player, id) })
			case slotYourArtifact:
				return m.useArtifact(sl.id)
			}
		}
	case "f":
		if sl, ok := m.curSlot(); ok && sl.kind == slotYourCreature {
			return m.startFight(sl.id)
		}
	}
	return m, nil
}

func (m gameModel) handleFlankKey(k tea.KeyMsg) (gameModel, tea.Cmd) {
	switch k.String() {
	case "left", "h":
		return m.playCreature(true)
	case "right", "l", "enter", " ":
		return m.playCreature(false)
	case "esc", "c":
		m.phase = phaseMain
		m.playHandPos = -1
	}
	return m, nil
}

func (m gameModel) handleChooseKey(k tea.KeyMsg) (gameModel, tea.Cmd) {
	switch k.String() {
	case "up", "k":
		m.chooseCur = wrap(m.chooseCur, -1, len(m.chooseCands))
	case "down", "l":
		m.chooseCur = wrap(m.chooseCur, 1, len(m.chooseCands))
	case "enter", " ":
		m.chooser.reply <- chooseReply{id: m.chooseCands[m.chooseCur], ok: true}
		m.phase = phaseBusy
	case "c", "esc":
		m.chooser.reply <- chooseReply{ok: false}
		m.suppressErr = true
		m.phase = phaseBusy
	}
	return m, nil
}

// handleChooseOptionKey drives a labeled option prompt (e.g. which key colour to
// forge). The choice is mandatory, so there is no cancel.
func (m gameModel) handleChooseOptionKey(k tea.KeyMsg) (gameModel, tea.Cmd) {
	switch k.String() {
	case "up", "k":
		m.chooseCur = wrap(m.chooseCur, -1, len(m.chooseOptions))
	case "down", "j":
		m.chooseCur = wrap(m.chooseCur, 1, len(m.chooseOptions))
	case "enter", " ":
		m.chooser.optionReply <- m.chooseCur
		m.phase = phaseBusy
	}
	return m, nil
}

func (m gameModel) handleFightKey(k tea.KeyMsg) (gameModel, tea.Cmd) {
	switch k.String() {
	case "up", "k":
		m.fightCur = wrap(m.fightCur, -1, len(m.fightCands))
	case "down", "l":
		m.fightCur = wrap(m.fightCur, 1, len(m.fightCands))
	case "esc", "c":
		m.phase = phaseMain
	case "enter", " ":
		def := m.fightCands[m.fightCur]
		g, player, att := m.g, m.player, m.fightAttacker
		return m.runAction(func() error { return g.Fight(player, att, def) })
	}
	return m, nil
}

func (m gameModel) handleCreatureMenuKey(k tea.KeyMsg) (gameModel, tea.Cmd) {
	switch k.String() {
	case "up", "k":
		m.actionMenuCur = wrap(m.actionMenuCur, -1, len(m.actionMenu))
	case "down", "l":
		m.actionMenuCur = wrap(m.actionMenuCur, 1, len(m.actionMenu))
	case "esc", "c":
		m.phase = phaseMain
	case "enter", " ":
		if m.actionMenuCur < 0 || m.actionMenuCur >= len(m.actionMenu) {
			m.phase = phaseMain
			return m, nil
		}
		id := m.actionCreature
		g, player := m.g, m.player
		switch m.actionMenu[m.actionMenuCur] {
		case caReap:
			return m.runAction(func() error { return g.Reap(player, id) })
		case caFight:
			return m.startFight(id)
		case caAction:
			return m.runAction(func() error { return g.UseAction(player, id) })
		}
	}
	return m, nil
}

func (m gameModel) handleDiscardKey(k tea.KeyMsg) (gameModel, tea.Cmd) {
	n := len(m.g.Discard(m.player)) + len(m.g.Discard(1-m.player))
	switch k.String() {
	case "up", "k":
		m.discardCur = wrap(m.discardCur, -1, n)
	case "down", "l":
		m.discardCur = wrap(m.discardCur, 1, n)
	case "esc", "v", "q":
		m.phase = phaseMain
	}
	return m, nil
}

// promptExit asks the player to confirm leaving the game for the menu, keeping
// from so the game can be resumed if they cancel.
func (m gameModel) promptExit(from gamePhase) (gameModel, tea.Cmd) {
	m.exitReturnPhase = from
	m.phase = phaseConfirmExit
	m.status = ""
	return m, nil
}

func (m gameModel) handleConfirmExitKey(k tea.KeyMsg) (gameModel, tea.Cmd) {
	switch k.String() {
	case "y", "enter":
		return m, gotoScreen(screenMenu)
	case "n", "esc", "c":
		m.phase = m.exitReturnPhase
	}
	return m, nil
}

// ---- actions ----

// runAction starts an engine mutation in a goroutine so the chooser bridge can
// post choice requests back to the program without blocking the UI thread.
func (m gameModel) runAction(fn func() error) (gameModel, tea.Cmd) {
	m.phase = phaseBusy
	m.status = ""
	return m, func() tea.Msg { return actionDoneMsg{err: fn()} }
}

func (m gameModel) beginTurn() (gameModel, tea.Cmd) {
	m.phase = phaseBusy
	g, player := m.g, m.player
	return m, func() tea.Msg { g.StartTurn(player); return beginDoneMsg{} }
}

func (m gameModel) endTurn() (gameModel, tea.Cmd) {
	m.g.EndPlayPhase(m.player)
	m.player = 1 - m.player
	m.status = ""
	return m.beginTurn() // forge + fire triggers for the next player, then house phase
}

func (m gameModel) enterHousePhase() gameModel {
	m.houseCur = 0
	m.phase = phaseHouse
	m.status = ""
	return m
}

func (m gameModel) playSlot(sl slot) (gameModel, tea.Cmd) {
	def := m.g.Def(sl.id)
	if h := m.g.State.ActiveHouse; h != engine.HouseNone && def.House != h {
		m.status = def.Name + " is not in your active house (" + h.String() + ")"
		return m, nil
	}
	m.playHandPos = m.cur - m.firstHandSlot() // keep the cursor at this hand row after the play
	if def.Type == engine.Creature {
		// With no creatures already in your battleline, left and right are the same
		// spot, so play immediately instead of prompting for a flank.
		if len(m.g.Battleline(m.player)) == 0 {
			g, player, idx := m.g, m.player, sl.handIdx
			return m.runAction(
				func() error { _, err := g.PlayCreature(player, idx, false); return err },
			)
		}
		m.pendingHand = sl.handIdx
		m.phase = phaseFlank
		m.status = ""
		return m, nil
	}
	g, player, idx := m.g, m.player, sl.handIdx
	var fn func() error
	switch def.Type {
	case engine.Artifact:
		fn = func() error { _, err := g.PlayArtifact(player, idx); return err }
	case engine.Tactic:
		fn = func() error { return g.PlayAction(player, idx) }
	case engine.Upgrade:
		fn = func() error { _, err := g.PlayUpgrade(player, idx); return err }
	default:
		m.status = "cannot play that card"
		return m, nil
	}
	return m.runAction(fn)
}

func (m gameModel) playCreature(left bool) (gameModel, tea.Cmd) {
	g, player, idx := m.g, m.player, m.pendingHand
	return m.runAction(func() error { _, err := g.PlayCreature(player, idx, left); return err })
}

// discardHand discards the selected hand card (must be of the active house).
func (m gameModel) discardHand(sl slot) (gameModel, tea.Cmd) {
	def := m.g.Def(sl.id)
	if h := m.g.State.ActiveHouse; h != engine.HouseNone && def.House != h {
		m.status = def.Name + " is not in your active house (" + h.String() + ")"
		return m, nil
	}
	m.playHandPos = m.cur - m.firstHandSlot()
	g, player, idx := m.g, m.player, sl.handIdx
	return m.runAction(func() error { return g.DiscardFromHand(player, idx) })
}

func (m gameModel) startFight(att engine.LocalID) (gameModel, tea.Cmd) {
	if err := m.g.CanUse(m.player, att); err != nil {
		m.status = err.Error()
		m.phase = phaseMain
		return m, nil
	}
	enemy := m.g.Battleline(1 - m.player)
	if len(enemy) == 0 {
		m.status = "no enemy creatures to fight"
		m.phase = phaseMain
		return m, nil
	}
	m.fightAttacker = att
	m.fightCands = enemy
	m.fightCur = 0
	m.phase = phaseTargetFight
	m.status = ""
	return m, nil
}

func (m gameModel) openCreatureMenu(id engine.LocalID) (gameModel, tea.Cmd) {
	m.actionCreature = id
	m.actionMenu = m.creatureActions(id)
	m.actionMenuCur = 0
	m.phase = phaseCreatureMenu
	m.status = ""
	return m, nil
}

// creatureActions lists the actions offered for one of your creatures: Reap
// always, Fight only if the opponent has creatures, and Action only if the
// creature has an "Action:" ability.
func (m gameModel) creatureActions(id engine.LocalID) []creatureAction {
	out := []creatureAction{caReap}
	if len(m.g.Battleline(1-m.player)) > 0 {
		out = append(out, caFight)
	}
	if m.hasActionAbility(id) {
		out = append(out, caAction)
	}
	return out
}

// hasActionAbility reports whether a card has an "Action:" ability.
func (m gameModel) hasActionAbility(id engine.LocalID) bool {
	for _, ab := range m.g.Def(id).Abilities {
		if ab.Trigger == engine.TriggerAction {
			return true
		}
	}
	return false
}

// useArtifact activates an artifact's "Action:" ability if it has one; otherwise
// the artifact is read-only and the attempt is reported in the status line.
func (m gameModel) useArtifact(id engine.LocalID) (gameModel, tea.Cmd) {
	if !m.hasActionAbility(id) {
		m.status = m.g.Def(id).Name + " has no action ability"
		return m, nil
	}
	g, player := m.g, m.player
	return m.runAction(func() error { return g.UseAction(player, id) })
}

// ---- highlight / detail: which card is under the cursor ----

// highlight returns the board id to highlight (with hasSel) and the hand index to
// highlight (-1 for none) for the current phase.
func (m gameModel) highlight() (engine.LocalID, bool, int) {
	switch m.phase {
	case phaseChoose:
		if len(m.chooseCands) > 0 {
			return m.chooseCands[m.chooseCur], true, -1
		}
	case phaseTargetFight:
		if len(m.fightCands) > 0 {
			return m.fightCands[m.fightCur], true, -1
		}
	case phaseFlank:
		return 0, false, m.pendingHand
	case phaseCreatureMenu:
		return m.actionCreature, true, -1
	case phaseMain:
		if sl, ok := m.curSlot(); ok {
			if sl.kind == slotHand {
				return 0, false, sl.handIdx
			}
			return sl.id, true, -1
		}
	}
	return 0, false, -1
}

func (m gameModel) selectedDef() *engine.CardDefinition {
	switch m.phase {
	case phaseChoose:
		if len(m.chooseCands) > 0 {
			return m.g.Def(m.chooseCands[m.chooseCur])
		}
	case phaseTargetFight:
		if len(m.fightCands) > 0 {
			return m.g.Def(m.fightCands[m.fightCur])
		}
	case phaseFlank:
		return m.pendingDef()
	case phaseCreatureMenu:
		return m.g.Def(m.actionCreature)
	case phaseMain:
		if sl, ok := m.curSlot(); ok {
			return m.g.Def(sl.id)
		}
	}
	return nil
}

func (m gameModel) pendingDef() *engine.CardDefinition {
	hand := m.g.Hand(m.player)
	if m.pendingHand >= 0 && m.pendingHand < len(hand) {
		return m.g.Def(hand[m.pendingHand])
	}
	return nil
}

// ---- view ----

func (m gameModel) View() string {
	switch m.phase {
	case phaseOver:
		return m.overView()
	case phaseDiscard:
		return m.discardView()
	default:
		return m.boardView()
	}
}

func (m gameModel) overView() string {
	w := m.g.Winner()
	return titleStyle.Render("Game over") + "\n\n" +
		fmt.Sprintf("%s wins with %d keys!\n\n", m.g.PlayerName(w), m.g.Keys(w)) +
		helpStyle.Render("enter return to menu")
}

func (m gameModel) discardView() string {
	you := m.g.Discard(m.player)
	opp := m.g.Discard(1 - m.player)
	all := append(append([]engine.LocalID{}, you...), opp...)
	var b strings.Builder
	b.WriteString(titleStyle.Render("Discard piles") + "\n\n")
	if len(all) == 0 {
		b.WriteString(faintStyle.Render("  (both zones empty)") + "\n\n")
		b.WriteString(helpStyle.Render("esc/v back"))
		return b.String()
	}
	cur := clamp(m.discardCur, 0, len(all)-1)
	idx := 0
	writePile := func(label string, ids []engine.LocalID) {
		b.WriteString(headerStyle.Render(label) + "\n")
		if len(ids) == 0 {
			b.WriteString(faintStyle.Render("  (empty)") + "\n")
			return
		}
		for _, id := range ids {
			name := m.g.Def(id).Name
			if idx == cur {
				b.WriteString(selectedStyle.Render(cursor(true)+name) + "\n")
			} else {
				b.WriteString(cursor(false) + name + "\n")
			}
			idx++
		}
	}
	writePile(m.g.PlayerName(m.player)+" (you)", you)
	writePile(m.g.PlayerName(1-m.player)+" (opponent)", opp)

	def := m.g.Def(all[cur])
	b.WriteString("\n" + faintStyle.Render("\u2500\u2500 "+def.Name+" \u2500\u2500") + "\n")
	b.WriteString(indentText(engine.RenderCardText(def)) + "\n")
	b.WriteString("\n" + helpStyle.Render("\u2191/\u2193 move \u00b7 esc/v back"))
	return b.String()
}

func (m gameModel) boardView() string {
	g := m.g
	opp := 1 - m.player
	selID, hasSel, selHand := m.highlight()
	var b strings.Builder

	b.WriteString(titleStyle.Render("VACTROL") + "  " +
		faintStyle.Render(
			fmt.Sprintf("turn %d — %s to act", g.State.Turn, g.PlayerName(m.player)),
		) + "\n")
	b.WriteString("\n")
	b.WriteString(m.scoreLine(opp, "opponent") + "\n")
	b.WriteString(headerStyle.Render(g.PlayerName(opp)+"'s board") +
		faintStyle.Render(fmt.Sprintf("  (%d cards in hand)", len(g.Hand(opp)))) + "\n")
	b.WriteString(m.renderCreatures(opp, selID, hasSel))
	b.WriteString(m.renderArtifacts(opp, selID, hasSel))
	b.WriteString("\n")

	b.WriteString(m.scoreLine(m.player, "you") + "\n")
	b.WriteString(headerStyle.Render(g.PlayerName(m.player)+"'s board") + "\n")
	b.WriteString(m.renderCreatures(m.player, selID, hasSel))
	b.WriteString(m.renderArtifacts(m.player, selID, hasSel))
	b.WriteString("\n")

	b.WriteString(headerStyle.Render(g.PlayerName(m.player)+"'s hand") + "\n")
	b.WriteString(m.renderHand(selHand))

	if def := m.selectedDef(); def != nil {
		b.WriteString("\n" + faintStyle.Render("── "+def.Name+" ──") + "\n")
		b.WriteString(indentText(engine.RenderCardText(def)) + "\n")
	}

	b.WriteString("\n" + m.bottom())
	if m.status != "" {
		b.WriteString("\n" + errStyle.Render(m.status))
	}
	b.WriteString(m.logTail(10))
	return b.String()
}

// logTail renders the last n engine log lines as a faint panel at the bottom of
// the screen (newest last); it renders nothing when the log is empty.
func (m gameModel) logTail(n int) string {
	log := m.g.Log
	if len(log) == 0 {
		return ""
	}
	if len(log) > n {
		log = log[len(log)-n:]
	}
	var b strings.Builder
	b.WriteString("\n\n" + faintStyle.Render("── log ──") + "\n")
	for _, line := range log {
		b.WriteString(faintStyle.Render(line.Text(m.g)) + "\n")
	}
	return b.String()
}

// bottom renders the phase-specific prompt/help under the board.
func (m gameModel) bottom() string {
	switch m.phase {
	case phaseHouse:
		var b strings.Builder
		b.WriteString(headerStyle.Render("Choose a house") + "\n")
		for i, h := range m.deckHouses[m.player] {
			if i == m.houseCur {
				b.WriteString(selectedStyle.Render(cursor(true)+h.String()) + "\n")
			} else {
				b.WriteString(cursor(false) + h.String() + "\n")
			}
		}
		b.WriteString(helpStyle.Render("↑/↓ · enter choose · e end/skip turn · esc menu"))
		return b.String()
	case phaseFlank:
		name := "creature"
		if def := m.pendingDef(); def != nil {
			name = def.Name
		}
		return selectedStyle.Render("Place "+name+" on which flank?") + "\n" +
			helpStyle.Render("← left · → / enter right · c cancel")
	case phaseCreatureMenu:
		var b strings.Builder
		b.WriteString(selectedStyle.Render("Use "+m.g.Def(m.actionCreature).Name+":") + "\n")
		for i, ca := range m.actionMenu {
			if i == m.actionMenuCur {
				b.WriteString(selectedStyle.Render(cursor(true)+ca.label()) + "\n")
			} else {
				b.WriteString(cursor(false) + ca.label() + "\n")
			}
		}
		b.WriteString(helpStyle.Render("↑/↓ · enter use · c cancel"))
		return b.String()
	case phaseChoose:
		return selectedStyle.Render(m.choosePrompt) + "\n" +
			helpStyle.Render("↑/↓ pick a highlighted creature · enter choose · c cancel")
	case phaseChooseOption:
		var b strings.Builder
		b.WriteString(selectedStyle.Render(m.choosePrompt) + "\n")
		for i, opt := range m.chooseOptions {
			if i == m.chooseCur {
				b.WriteString(selectedStyle.Render(cursor(true)+opt) + "\n")
			} else {
				b.WriteString(cursor(false) + opt + "\n")
			}
		}
		b.WriteString(helpStyle.Render("↑/↓ · enter choose"))
		return b.String()
	case phaseTargetFight:
		return selectedStyle.Render("Choose an enemy creature to fight") + "\n" +
			helpStyle.Render("↑/↓ pick a highlighted enemy · enter fight · c cancel")
	case phaseBusy:
		return faintStyle.Render("resolving…")
	case phaseConfirmExit:
		return selectedStyle.Render("Leave the game and return to the menu?") + "\n" +
			helpStyle.Render("y confirm · n cancel")
	default:
		return m.mainHelp()
	}
}

func (m gameModel) mainHelp() string {
	sl, ok := m.curSlot()
	if !ok {
		return helpStyle.Render("↑/↓ move · v discards · e end turn · esc menu")
	}
	switch sl.kind {
	case slotHand:
		return helpStyle.Render(
			"↑/↓ move · enter play · d discard · v discards · e end turn · esc menu",
		)
	case slotYourCreature:
		return helpStyle.Render("↑/↓ move · enter/r/f/a use · v discards · e end turn · esc menu")
	case slotYourArtifact:
		if m.hasActionAbility(sl.id) {
			return helpStyle.Render(
				"↑/↓ move · enter/a use action · v discards · e end turn · esc menu",
			)
		}
		return helpStyle.Render("↑/↓ move · (read-only) · v discards · e end turn · esc menu")
	default:
		return helpStyle.Render("↑/↓ move · (read-only) · v discards · e end turn · esc menu")
	}
}

func (m gameModel) scoreLine(player int, tag string) string {
	g := m.g
	line := fmt.Sprintf("%-20s keys %d/%d   Æmber %d",
		g.PlayerName(player)+" ("+tag+")", g.Keys(player), engine.KeysToWin, g.Aember(player))
	if player == m.player {
		if h := g.State.ActiveHouse; h != engine.HouseNone {
			line += "   house " + activeHouseStyle.Render(h.String())
		} else {
			line += "   " + faintStyle.Render("house none")
		}
	}
	return line
}

// renderCardLine colors a card row by house — green for the turn's active
// house, blue for every other — and reverses the selected row so it stands out.
func (m gameModel) renderCardLine(house engine.House, sel bool, line string) string {
	st := otherHouseStyle
	if h := m.g.State.ActiveHouse; h != engine.HouseNone && house == h {
		st = activeHouseStyle
	}
	if sel {
		st = st.Bold(true).Reverse(true)
	}
	return st.Render(line)
}

// nameColWidth returns the width of the card-name column: wide enough for the
// longest name currently on either board or in either hand, so the house/type
// columns line up even when a name is long. It never shrinks below 18.
func (m gameModel) nameColWidth() int {
	w := 18
	consider := func(id engine.LocalID) {
		if n := len(m.g.Def(id).Name); n > w {
			w = n
		}
	}
	for p := 0; p < 2; p++ {
		for _, id := range m.g.Battleline(p) {
			consider(id)
		}
		for _, id := range m.g.Artifacts(p) {
			consider(id)
		}
		for _, id := range m.g.Hand(p) {
			consider(id)
		}
	}
	return w
}

func (m gameModel) renderCreatures(player int, selID engine.LocalID, hasSel bool) string {
	bl := m.g.Battleline(player)
	if len(bl) == 0 {
		return faintStyle.Render("  (no creatures)") + "\n"
	}
	var b strings.Builder
	w := m.nameColWidth()
	for _, id := range bl {
		def := m.g.Def(id)
		state := "ready"
		if m.g.Exhausted(id) {
			state = "exhausted"
		}
		dmg := ""
		if d := m.g.Damage(id); d > 0 {
			dmg = fmt.Sprintf(", %d dmg", d)
		}
		amber := ""
		if a := m.g.AmberOn(id); a > 0 {
			amber = fmt.Sprintf(", %d Æmber", a)
		}
		kw := ""
		if len(def.Keywords) > 0 {
			parts := make([]string, len(def.Keywords))
			for j, k := range def.Keywords {
				parts[j] = string(k)
			}
			kw = " [" + strings.Join(parts, ", ") + "]"
		}
		sel := hasSel && id == selID
		house := m.g.House(id)
		hs := house.String()
		if house != def.House { // not its natural house: mark the change
			hs += "*"
		}
		line := fmt.Sprintf(
			"%s%-*s %-8s %d power (%s%s%s)%s",
			cursor(sel),
			w,
			def.Name,
			hs,
			m.g.Power(id),
			state,
			dmg,
			amber,
			kw,
		)
		line = m.renderCardLine(house, sel, line)
		b.WriteString(line + "\n")
		for _, up := range m.g.Upgrades(id) {
			upLine := "      \u21b3 " + m.g.Def(up).Name + " (upgrade)"
			if hasSel && up == selID {
				upLine = selectedStyle.Render(upLine)
			} else {
				upLine = faintStyle.Render(upLine)
			}
			b.WriteString(upLine + "\n")
		}
	}
	return b.String()
}

func (m gameModel) renderArtifacts(player int, selID engine.LocalID, hasSel bool) string {
	arts := m.g.Artifacts(player)
	var b strings.Builder
	w := m.nameColWidth()
	for _, id := range arts {
		state := "ready"
		if m.g.Exhausted(id) {
			state = "exhausted"
		}
		sel := hasSel && id == selID
		def := m.g.Def(id)
		house := m.g.House(id)
		hs := house.String()
		if house != def.House { // not its natural house: mark the change
			hs += "*"
		}
		line := fmt.Sprintf("%s%-*s %-8s artifact (%s)", cursor(sel), w, def.Name, hs, state)
		line = m.renderCardLine(house, sel, line)
		b.WriteString(line + "\n")
	}
	return b.String()
}

func (m gameModel) renderHand(selHandIdx int) string {
	slots := m.handSlots()
	if len(slots) == 0 {
		return faintStyle.Render("  (empty)") + "\n"
	}
	var b strings.Builder
	w := m.nameColWidth()
	for _, sl := range slots {
		def := m.g.Def(sl.id)
		sel := sl.handIdx == selHandIdx
		line := fmt.Sprintf(
			"%s%-*s %-7s · %s",
			cursor(sel),
			w,
			def.Name,
			def.House.String(),
			def.Type,
		)
		line = m.renderCardLine(def.House, sel, line)
		b.WriteString(line + "\n")
	}
	return b.String()
}

func indentText(s string) string { return "  " + strings.ReplaceAll(s, "\n", "\n  ") }
