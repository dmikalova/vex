package web

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/engine"
)

// Render draws the whole client. It runs on both the server (prerender) and the
// client; before OnMount seeds the match on the client, g is nil, so a lightweight
// placeholder is shown.
func (g *game) Render() app.UI {
	if g.g == nil {
		return app.Div().Class("loading").Text("Loading Vactrol…")
	}

	return app.Div().Class(cx("app", ifCls(g.sidebarCollapsed, "app--sidebar-collapsed"))).Body(
		app.Raw(iconOutlineFilter),
		app.Div().Class("board-area").Body(g.boardArea()...),
		app.If(!g.sidebarCollapsed, func() app.UI {
			return app.Div().Class("sidebar").Body(
				g.brandBar(),
				g.turnHud(),
				g.logPanel(),
				app.If(g.status != "", func() app.UI { return g.statusBanner() }),
				g.controls(),
			)
		}),
		app.If(g.sidebarCollapsed, func() app.UI {
			return app.Button().Class("sidebar-reveal").Title("Show sidebar").
				Text("«").OnClick(g.toggleSidebar)
		}),
		app.If(g.hoverID != 0 || g.hoverDef != nil, func() app.UI { return g.hoverPreview() }),
		app.If(g.zonesPlayer >= 0, func() app.UI { return g.zonesOverlay() }),
		app.If(g.pickerOpen, func() app.UI { return g.cardPicker() }),
		app.If(g.phase == phaseOver, func() app.UI { return g.overBanner() }),
	)
}

// brandBar is the slim top of the sidebar: the title, a busy badge, and the
// navigation buttons (undo/redo, sandbox, new game).
func (g *game) brandBar() app.UI {
	return app.Div().Class("brandbar").Body(
		app.Span().Class("brand-title").Text("Vactrol"),
		app.If(g.busy && !g.choosing && !g.choosingOption, func() app.UI {
			return app.Span().Class("badge-busy").Text("resolving…")
		}),
		app.Div().Class("spacer"),
		app.Button().Class("btn-nav btn-icon").Title("Undo").
			Body(icon("undo", "icon-nav")).
			Disabled(!g.canUndo()).OnClick(g.undoAction),
		app.Button().Class("btn-nav btn-icon").Title("Redo").
			Body(icon("redo", "icon-nav")).
			Disabled(!g.canRedo()).OnClick(g.redoAction),
		app.Button().Class(cx("btn-nav", "btn-icon", ifCls(g.g.Manual(), "btn-nav-on"))).
			Title("Manual mode").
			Body(icon("wrench", "icon-nav")).
			Disabled(g.busy || g.choosing || g.choosingOption).
			OnClick(g.toggleManual),
		app.Button().Class("btn-nav btn-icon").Title("New game").
			Body(icon("restart", "icon-nav")).
			Disabled(g.busy || g.choosing || g.choosingOption).
			OnClick(g.askRestart),
		app.Button().Class("btn-nav btn-icon").Title("Hide sidebar").
			Text("»").OnClick(g.toggleSidebar),
	)
}

// statusBanner shows the transient status (usually a play error) as a red pill in
// the controls area. It fades out over 5s (setStatus also clears the message after
// 5s); statusGen parity alternates the class so the fade replays when a new error
// arrives while one is still showing.
func (g *game) statusBanner() app.UI {
	cls := cx("status-banner",
		ifCls(g.statusGen%2 == 0, "status-banner--a"),
		ifCls(g.statusGen%2 == 1, "status-banner--b"),
	)
	return app.Div().Class(cls).Text(g.status)
}

// ---- board ----

// boardArea renders both battlelines facing each other (opponent on top), each
// player's score, and the active player's hand. Card rows squeeze to the height
// available and scroll horizontally when full. The play rows between the two
// score pills form the drop zone for playing a card dragged from hand.
func (g *game) boardArea() []app.UI {
	p := g.active()
	opp := 1 - p
	playZone := app.Div().
		Class(cx("play-zone", ifCls(g.dragging, "play-zone--drop"))).
		OnDrop(g.dropOnBoard).
		Body(
			g.renderRow(g.g.PlayerName(opp)+"'s artifacts", g.g.Artifacts(opp), selOther),
			g.renderRow(g.g.PlayerName(opp)+"'s creatures", g.g.Battleline(opp), selOther),
			app.Div().Class("midline"),
			g.renderRow(g.g.PlayerName(p)+"'s creatures", g.g.Battleline(p), selYourCreature),
			g.renderRow(g.g.PlayerName(p)+"'s artifacts", g.g.Artifacts(p), selYourArtifact),
		)
	return []app.UI{
		g.scorePill(opp),
		playZone,
		g.scorePill(p),
		g.renderHand(),
	}
}

// turnHud is a slim status line: the turn number, whose turn it is, the current
// step, and the active house.
func (g *game) turnHud() app.UI {
	p := g.active()
	steps := map[phase]string{
		phaseHouse:       "Choose a house",
		phaseMain:        "Play, use, discard",
		phaseFlank:       "Placing a creature",
		phaseFightTarget: "Choosing a fight target",
		phaseOver:        "Game over",
	}
	items := []app.UI{
		app.Span().Class("hud-turn").Text(fmt.Sprintf("Turn %d", g.g.State.Turn)),
		app.Span().Class("hud-player").Text(g.g.PlayerName(p)),
		app.Span().Class("hud-step").Text(steps[g.phase]),
	}
	if h := g.g.State.ActiveHouse; h != engine.HouseNone {
		items = append(items, app.Span().Class(cx("hud-house", houseAccent(h))).Body(
			houseIcon(h, "icon-inline"), app.Text(h.String()),
		))
	}
	return app.Div().Class("hud").Body(items...)
}

func (g *game) scorePill(player int) app.UI {
	active := player == g.active()
	detail := []app.UI{
		app.Text(" • "),
		g.aemberSeg(player),
		app.Text(" • "),
		g.keysDisplay(player),
		app.Text(" • "),
		g.keyCostSeg(player),
	}
	if houses := g.deckHouses[player]; len(houses) > 0 {
		detail = append(detail, app.Text(" • "), g.houseStrip(player, houses))
	}
	if g.g.State.Chains[player] > 0 || g.g.Manual() {
		detail = append(detail, app.Text(" • "), g.chainsSeg(player))
	}
	zones := fmt.Sprintf("deck %d • discard %d • archives %d • purge %d",
		len(g.g.Deck(player)), len(g.g.Discard(player)),
		len(g.g.Archives(player)), len(g.g.Purge(player)))
	cls := cx("score-pill", ifCls(active, "score-pill-active"), ifCls(!active, "score-pill-idle"))
	return app.Div().Class(cls).
		Body(
			app.Span().Class("score-name").Text(g.g.PlayerName(player)),
			app.Span().Class("score-detail").Body(detail...),
			// Only the zone counts open the viewer, so misclicking a key or stepper
			// in the detail does not.
			app.Span().Class("score-zones").
				DataSet("player", strconv.Itoa(player)).
				OnClick(g.onScorePillClick).
				Text(zones),
		)
}

// keyCostSeg shows the Æmber a player must spend to forge their next key.
func (g *game) keyCostSeg(player int) app.UI {
	return app.Span().Class("stat-seg").Body(
		app.Text("forge "),
		app.Text(strconv.Itoa(g.g.CurrentKeyCost(player))),
		icon("aember", "icon-stat"),
	)
}

// hoverPreview renders the hovered card enlarged: a live board/hand card over the
// log, or a printed card (from a log mention) just left of the log.
func (g *game) hoverPreview() app.UI {
	var card app.UI
	switch {
	case g.hoverID != 0:
		def := g.g.Def(g.hoverID)
		card = &cardView{
			Title:    def.Name,
			HouseCls: houseClasses(def.House),
			Emblem:   houseIconName(def.House),
			TypeIcon: typeIconName(def.Type),
			Stat:     g.statLine(g.hoverID),
			Rules:    g.faceRules(g.hoverID),
			Kind:     string(def.Type),
			Stunned:  g.g.Stunned(g.hoverID),
		}
	case g.hoverDef != nil:
		def := g.hoverDef
		card = &cardView{
			Title:    def.Name,
			HouseCls: houseClasses(def.House),
			Emblem:   houseIconName(def.House),
			TypeIcon: typeIconName(def.Type),
			Stat:     handStat(def),
			Rules:    engine.RenderCardRules(def),
			Kind:     string(def.Type),
		}
	default:
		return app.Div()
	}
	pos := "card-preview--board"
	if g.hoverInLog {
		pos = "card-preview--log"
	}
	return app.Div().Class(cx("card-preview", pos)).Body(card)
}

// houseStrip shows the player's three deck houses in their score pill. Once that
// player has chosen an active house, the other houses are lowlighted.
func (g *game) houseStrip(player int, houses []engine.House) app.UI {
	active := engine.HouseNone
	if player == g.active() {
		active = g.g.State.ActiveHouse
	}
	// In sandbox mode the active player can switch their active house by clicking.
	clickable := g.g.Manual() && player == g.active()
	items := make([]app.UI, 0, len(houses))
	for _, h := range houses {
		dim := active != engine.HouseNone && h != active
		cls := cx("score-house", houseAccent(h), ifCls(dim, "score-house-dim"))
		if clickable {
			items = append(items, app.Button().Class(cx(cls, "score-house-btn")).
				OnClick(g.manualSetHouse(h)).
				Body(houseIcon(h, "icon-inline"), app.Text(h.String())))
		} else {
			items = append(items, app.Span().Class(cls).
				Body(houseIcon(h, "icon-inline"), app.Text(h.String())))
		}
	}
	return app.Span().Class("score-houses").Body(items...)
}

// aemberSeg shows a player's Æmber; in sandbox mode it flanks the count with
// minus/plus buttons that adjust it (usable on both players from one seat).
func (g *game) aemberSeg(player int) app.UI {
	count := app.Text(strconv.Itoa(g.g.Aember(player)))
	ic := icon("aember", "icon-stat")
	// A pool gain pulses the segment; the -a/-b pair alternates so it replays.
	gain := cx(
		ifCls(g.poolFlash[player] && !g.poolParity[player], "stat-seg--gain-a"),
		ifCls(g.poolFlash[player] && g.poolParity[player], "stat-seg--gain-b"),
	)
	if !g.g.Manual() {
		return app.Span().Class(cx("stat-seg", gain)).Body(count, ic)
	}
	return app.Span().Class(cx("stat-seg", "amber-manual", gain)).Body(
		g.stepBtn(g.manualAmberDelta(player, -1), false),
		count, ic,
		g.stepBtn(g.manualAmberDelta(player, 1), true),
	)
}

// chainsSeg shows a player's chains; in sandbox mode it always shows (even at 0)
// with minus/plus steppers.
func (g *game) chainsSeg(player int) app.UI {
	count := app.Text(strconv.Itoa(g.g.State.Chains[player]))
	ic := icon("chains", "icon-stat")
	if !g.g.Manual() {
		return app.Span().Class("stat-seg").Body(count, ic)
	}
	return app.Span().Class("stat-seg amber-manual").Body(
		g.stepBtn(g.manualChainsDelta(player, -1), false),
		count, ic,
		g.stepBtn(g.manualChainsDelta(player, 1), true),
	)
}

// stepBtn is a green plus or red minus stepper for the sandbox Æmber/chains
// adjusters.
func (g *game) stepBtn(onClick app.EventHandler, plus bool) app.UI {
	label, cls := "−", "amber-btn amber-btn-minus"
	if plus {
		label, cls = "+", "amber-btn amber-btn-plus"
	}
	return app.Button().Class(cls).Text(label).OnClick(onClick)
}

// keyForgePanel is the sandbox key-forge picker, shown inline in the controls
// space: pick a colour for the new key, or cancel. Once a choice is made
// forgingKey resets, so controls() falls back to the previous buttons on its own.
func (g *game) keyForgePanel() app.UI {
	player := g.forgingKey
	remaining := g.remainingKeyColors(player)
	return app.Div().Class("btn-col").Body(
		app.Div().Class("section-title").Text("Forge a key for "+g.g.PlayerName(player)),
		app.Range(remaining).Slice(func(i int) app.UI {
			c := remaining[i]
			return app.Button().Class("house-btn").OnClick(g.pickForgeColor(c)).Body(
				icon(keyColorIconName(c), "icon-inline"), app.Text(c.String()),
			)
		}),
		btn("Cancel", g.cancelForgeKey, "btn-secondary"),
	)
}

// keysDisplay shows a player's three key slots: a coloured key icon for each key
// forged (in forge order), and a dimmed key for each still to forge. In sandbox
// mode each slot is a button that forges or unforges that key.
func (g *game) keysDisplay(player int) app.UI {
	colors := g.g.KeyColors(player)
	manual := g.g.Manual()
	slots := make([]app.UI, 0, engine.KeysToWin)
	for _, c := range colors {
		// A forged key with no recorded colour (e.g. a legacy snapshot) still counts,
		// so show the neutral key rather than a broken image from an empty icon name.
		name := keyColorIconName(c)
		if name == "" {
			name = "key"
		}
		slots = append(
			slots,
			g.keySlot(icon(name, "icon-stat"), manual, g.manualUnforgeKey(player)),
		)
	}
	for i := len(colors); i < engine.KeysToWin; i++ {
		slots = append(
			slots,
			g.keySlot(icon("key", "icon-stat", "key-unforged"), manual, g.manualForgeKey(player)),
		)
	}
	gain := cx(
		ifCls(g.keyFlash[player] && !g.keyParity[player], "stat-seg--gain-a"),
		ifCls(g.keyFlash[player] && g.keyParity[player], "stat-seg--gain-b"),
	)
	return app.Span().Class(cx("score-keys", gain)).Body(slots...)
}

// keySlot renders a key icon as a clickable forge/unforge button in sandbox mode,
// or a plain icon otherwise.
func (g *game) keySlot(ic app.UI, manual bool, onClick app.EventHandler) app.UI {
	if !manual {
		return ic
	}
	return app.Button().Class("key-btn").OnClick(onClick).Body(ic)
}

func (g *game) renderRow(label string, ids []engine.LocalID, boardKind selKind) app.UI {
	// The row always reserves card height (board-row--fill) so playing or losing a
	// card does not resize the row and shift the rest of the board.
	return app.Div().Class("board-row board-row--fill").Body(
		app.Div().Class("row-label").Text(fmt.Sprintf("%s (%d)", label, len(ids))),
		app.If(len(ids) == 0, func() app.UI {
			return app.Div().Class("row-empty").Text("—")
		}).Else(func() app.UI {
			return app.Div().Class("card-strip").Body(
				app.Range(ids).Slice(func(i int) app.UI { return g.renderCard(ids[i], boardKind) }),
			)
		}),
	)
}

func (g *game) renderCard(id engine.LocalID, boardKind selKind) app.UI {
	def := g.g.Def(id)
	activate, targetable, dimmed := g.cardVisual(id, boardKind)
	house := g.g.House(
		id,
	) // effective house: a control/"belongs to house" effect may override the printed one
	flash := g.flashes[id]
	return &cardView{
		ID:           id,
		Title:        def.Name,
		HouseCls:     houseClasses(house),
		Emblem:       houseIconName(house),
		HouseChanged: house != def.House,
		TypeIcon:     typeIconName(def.Type),
		Stat:         g.statLine(id),
		Rules:        g.faceRules(id),
		Kind:         string(def.Type),
		Stunned:      g.g.Stunned(id),
		Enter:        flash.enter,
		StunFlash:    flash.stun,
		FlashOdd:     flash.odd,
		Selected:     g.hasSel && g.sel == id,
		Targetable:   targetable,
		Dimmed:       dimmed,
		OnActivate:   activate,
		OnHover:      g.hoverCard,
		OnHoverOut:   g.hoverClear,
	}
}

// cardVisual decides how a card (in hand or in play) responds and looks in the
// current mode. It returns the click handler (nil when the card is not
// clickable), whether the card is a highlighted action target, and whether it is
// lowlighted (dimmed) as an invalid choice. During a chooser or fight-target
// prompt only the eligible cards are highlighted and the rest dimmed; in ordinary
// play, cards the active player cannot act with (wrong house, exhausted, or
// unplayable from hand) and the opponent's read-only cards are dimmed so the
// usable ones stand out.
func (g *game) cardVisual(
	id engine.LocalID,
	kind selKind,
) (activate func(app.Context, engine.LocalID), targetable, dimmed bool) {
	switch {
	case g.choosing:
		// A chooser runs on a background goroutine, so g.busy is also set; the
		// choosing case must come first or the candidates would not be clickable.
		if containsID(g.chooserCandidates, id) {
			return g.chooseCandidate, true, false
		}
		return nil, false, true
	case g.busy:
		return nil, false, false
	case g.phase == phaseFightTarget:
		if kind == selOther && g.isEnemyCreature(id) {
			return g.fightTargetID, true, false
		}
		return nil, false, true
	case kind == selHand:
		return g.selectHandID, false, !g.playableFromHand(id)
	case kind == selYourCreature, kind == selYourArtifact:
		return g.selectBoardID, false, !g.actionable(id, kind)
	default:
		// Opponent's cards are read-only in ordinary play: dimmed like the active
		// player's non-actionable cards, but still clickable to inspect.
		return g.selectBoardID, false, true
	}
}

// actionable reports whether the active player can act with one of their own
// cards this turn, so the ones they cannot use are lowlighted. Creatures defer to
// the engine's CanUse (house, readiness, ownership); an artifact needs an action
// ability, readiness, and the active house.
func (g *game) actionable(id engine.LocalID, kind selKind) bool {
	switch kind {
	case selYourCreature:
		return g.g.CanUse(g.active(), id) == nil
	case selYourArtifact:
		def := g.g.Def(id)
		inHouse := g.g.State.ActiveHouse == engine.HouseNone || def.House == g.g.State.ActiveHouse
		return inHouse && g.g.HasAction(id) && !g.g.Exhausted(id)
	default:
		return true
	}
}

// ---- hand ----

func (g *game) renderHand() app.UI {
	p := g.active()
	ids := g.sortedHand(p)
	cls := "board-row board-row--fill"
	return app.Div().Class(cls).Body(
		app.Div().
			Class("row-label").
			Text(fmt.Sprintf("%s (%d)", g.g.PlayerName(p)+"'s hand", len(ids))),
		app.Div().Class("card-strip").Body(
			app.Range(ids).Slice(func(i int) app.UI { return g.renderHandCard(ids[i]) }),
		),
	)
}

// sortedHand returns the player's hand ids ordered by house, then card type, then
// name, so the hand reads consistently. It sorts a copy — the engine's own hand
// order (which play/discard index into) is untouched, so selection still maps to
// the right card.
func (g *game) sortedHand(p int) []engine.LocalID {
	return g.sortByHouseTypeName(g.g.Hand(p))
}

// sortByHouseTypeName returns a copy of ids ordered by house, then card type,
// then name — the stable reading order shared by the hand and the deck view. The
// deck in particular must not reveal its shuffled order, so it is always sorted.
func (g *game) sortByHouseTypeName(ids []engine.LocalID) []engine.LocalID {
	ids = append([]engine.LocalID(nil), ids...)
	sort.SliceStable(ids, func(i, j int) bool {
		a, b := g.g.Def(ids[i]), g.g.Def(ids[j])
		if a.House != b.House {
			return a.House < b.House
		}
		if a.Type != b.Type {
			return a.Type < b.Type
		}
		return a.Name < b.Name
	})
	return ids
}

func (g *game) renderHandCard(id engine.LocalID) app.UI {
	def := g.g.Def(id)
	activate, targetable, dimmed := g.cardVisual(id, selHand)
	draggable := !g.busy && !g.choosing && !g.choosingOption &&
		g.phase == phaseMain && g.playableFromHand(id)
	return &cardView{
		ID:          id,
		Title:       def.Name,
		HouseCls:    houseClasses(def.House),
		Emblem:      houseIconName(def.House),
		TypeIcon:    typeIconName(def.Type),
		Stat:        handStat(def),
		Rules:       engine.RenderCardRules(def),
		Kind:        string(def.Type),
		Selected:    g.hasSel && g.sel == id,
		Targetable:  targetable,
		Dimmed:      dimmed,
		OnActivate:  activate,
		Draggable:   draggable,
		OnDragStart: g.startHandDrag,
		OnDragEnd:   g.endHandDrag,
		OnHover:     g.hoverCard,
		OnHoverOut:  g.hoverClear,
	}
}

// ---- sidebar controls ----

// promptSourceHeader names the card driving the current prompt, shown above the
// prompt's buttons so the player knows which card they are resolving.
func (g *game) promptSourceHeader() app.UI {
	return app.If(g.promptSource != "", func() app.UI {
		return app.Div().Class("prompt-source").Text(g.promptSource)
	})
}

// endTurnButton is the End turn control. Once a confirm is armed (the player could
// still act this turn) it becomes a red "Confirm end turn" that ends on the next
// click, mirroring pressing E/Y a second time.
func (g *game) endTurnButton() app.UI {
	if g.confirmEndTurn {
		return btn("Confirm end turn", g.endTurn, "btn-danger")
	}
	return btn("End turn", g.endTurn, "btn-secondary")
}

// controls is the bottom of the sidebar: the contextual controls (house picker or
// action bar) plus End turn. House selection has no End turn — a house must be
// chosen first.
func (g *game) controls() app.UI {
	// A pending restart takes over the controls until confirmed or cancelled.
	if g.confirmRestart {
		return app.Div().Class("controls").Body(
			app.Div().Class("btn-col").Body(
				app.Div().Class("section-title").Text("Restart game?"),
				app.Div().Class("hint").Text("This ends the current game and deals a new one."),
				btn("Restart", g.restart, "btn-danger"),
				btn("Cancel", g.cancelRestart, "btn-secondary"),
			),
		)
	}
	// A pending sandbox key forge takes over the controls until a colour is picked
	// or cancelled; once forgingKey resets, the previous buttons return on their own.
	if g.forgingKey >= 0 {
		return app.Div().Class("controls").Body(g.keyForgePanel())
	}
	// A labeled option prompt (e.g. "take archives?") shows its choices as buttons.
	if g.choosingOption {
		return app.Div().Class("controls").Body(g.promptSourceHeader(), g.optionChooser())
	}
	// While an engine chooser waits, the controls become the prompt itself: a
	// green call to action to click one of the highlighted cards. The choice is
	// mandatory, so there is no cancel.
	if g.choosing {
		body := []app.UI{
			g.promptSourceHeader(),
			app.Div().Class("prompt").Text(g.chooserPrompt),
		}
		// A mandatory chooser has no cancel; in sandbox mode add one so the player
		// can escape a prompt with no clickable candidate.
		if g.g.Manual() {
			body = append(body, btn("Cancel", g.cancelChooser, "btn-secondary"))
		}
		return app.Div().Class("controls").Body(body...)
	}
	if g.phase == phaseHouse {
		if g.g.Manual() {
			return app.Div().Class("controls").Body(g.manualPanel(), g.housePicker())
		}
		return app.Div().Class("controls").Body(g.housePicker())
	}
	// Mid-action selections (choosing a flank or a fight target) show only their
	// own controls — End turn is withheld until the action finishes or is cancelled.
	if g.phase == phaseFlank || g.phase == phaseFightTarget {
		return app.Div().Class("controls").Body(g.actionBar())
	}
	// In sandbox mode the manual controls sit above the normal action bar, which
	// still works (house restrictions lifted) so cards can be used out of house.
	if g.g.Manual() {
		return app.Div().Class("controls").Body(
			g.manualPanel(),
			g.actionBar(),
			g.endTurnButton(),
		)
	}
	return app.Div().Class("controls").Body(
		g.actionBar(),
		g.endTurnButton(),
	)
}

// manualPanel is the sandbox control block: add an arbitrary card, and (for a
// selected card) ready/exhaust it and move it between zones.
func (g *game) manualPanel() app.UI {
	items := []app.UI{
		app.Div().Class("section-title").Text("Manual"),
		btn("Add card…", g.openPicker, "btn-secondary"),
	}
	if g.hasSel {
		name := g.g.Def(g.sel).Name
		if g.isInPlay(g.sel) {
			if g.g.Exhausted(g.sel) {
				items = append(items, btn("Ready "+name, g.manualReady, "btn-secondary"))
			} else {
				items = append(items, btn("Exhaust "+name, g.manualExhaust, "btn-secondary"))
			}
		}
		items = append(items,
			app.Div().Class("hint").Text("Move "+name+" to:"),
			g.moveButtons(),
		)
	}
	return app.Div().Class("btn-col manual-panel").Body(items...)
}

// moveButtons is the row of destination buttons for manually relocating the
// selected card.
func (g *game) moveButtons() app.UI {
	dests := []struct {
		label string
		dest  engine.ManualZone
	}{
		{"Hand", engine.ManualHand},
		{"Deck top", engine.ManualDeckTop},
		{"Deck bottom", engine.ManualDeckBottom},
		{"Discard", engine.ManualDiscard},
		{"Archives", engine.ManualArchives},
		{"Purge", engine.ManualPurge},
	}
	buttons := make([]app.UI, len(dests))
	for i, d := range dests {
		buttons[i] = btn(d.label, g.manualMove(d.dest), "btn-mini")
	}
	return app.Div().Class("btn-wrap").Body(buttons...)
}

// cardPicker is the fuzzy, text-only card picker for adding an arbitrary card
// from the pool to hand. It filters the pool by a case-insensitive name substring.
func (g *game) cardPicker() app.UI {
	q := strings.ToLower(strings.TrimSpace(g.pickerQuery))
	var matches []engine.CardDefinition
	for _, d := range g.allDefs {
		if q == "" || strings.Contains(strings.ToLower(d.Name), q) {
			matches = append(matches, d)
		}
	}
	return app.Div().Class("over-backdrop").OnClick(g.closePicker).Body(
		app.Div().Class("picker-panel").OnClick(g.stopClick).Body(
			app.Button().Class("zones-close").Text("✕").OnClick(g.closePicker),
			app.Div().Class("over-title").Text("Add a card to hand"),
			app.Input().Class("picker-input").Type("text").Placeholder("Search cards…").
				AutoFocus(true).Value(g.pickerQuery).OnInput(g.pickerInput),
			app.Div().Class("picker-list").Body(
				app.Range(matches).Slice(func(i int) app.UI {
					d := matches[i]
					return app.Button().Class("picker-item").OnClick(g.addPickedCard(d)).Body(
						houseIcon(d.House, "icon-inline"),
						app.Span().Class("picker-name").Text(d.Name),
						app.Span().Class("picker-kind").Text(string(d.Type)),
					)
				}),
			),
		),
	)
}

func (g *game) housePicker() app.UI {
	houses := g.deckHouses[g.active()]
	return app.Div().Class("btn-col").Body(
		app.Div().Class("section-title").Text("Choose a house"),
		app.Range(houses).Slice(func(i int) app.UI {
			h := houses[i]
			return app.Button().
				Class(cx("house-btn", houseAccent(h))).
				OnClick(g.pickHouse(h)).
				Body(houseIcon(h, "icon-inline"), app.Text(h.String()))
		}),
	)
}

// optionChooser renders a labeled multiple-choice prompt. When every option is a
// house or a key colour, it shows themed icon buttons; otherwise plain primary
// buttons.
func (g *game) optionChooser() app.UI {
	if g.keyColorOptions() {
		return app.Div().Class("btn-col").Body(
			app.Div().Class("section-title").Text("Choose a key colour:"),
			app.Range(g.optionLabels).Slice(func(i int) app.UI {
				c := keyColorByName(g.optionLabels[i])
				return app.Button().
					Class("house-btn").
					OnClick(g.chooseOptionIdx(i)).
					Body(icon(keyColorIconName(c), "icon-inline"), app.Text(g.optionLabels[i]))
			}),
		)
	}
	if g.houseOptions() {
		return app.Div().Class("btn-col").Body(
			app.Div().Class("section-title").Text("Choose a house:"),
			app.Range(g.optionLabels).Slice(func(i int) app.UI {
				h, _ := engine.ParseHouse(g.optionLabels[i])
				return app.Button().
					Class(cx("house-btn", houseAccent(h))).
					OnClick(g.chooseOptionIdx(i)).
					Body(houseIcon(h, "icon-inline"), app.Text(g.optionLabels[i]))
			}),
		)
	}
	return app.Div().Class("btn-col").Body(
		app.Div().Class("prompt").Text(g.optionPrompt),
		app.Range(g.optionLabels).Slice(func(i int) app.UI {
			return btn(g.optionLabels[i], g.chooseOptionIdx(i), "btn-primary")
		}),
	)
}

// houseOptions reports whether every current option label names a house, so the
// prompt can be shown as the colored house picker.
func (g *game) houseOptions() bool {
	if len(g.optionLabels) == 0 {
		return false
	}
	for _, label := range g.optionLabels {
		if _, ok := engine.ParseHouse(label); !ok {
			return false
		}
	}
	return true
}

// keyColorOptions reports whether every current option label names a key colour,
// so the forge prompt can be shown as coloured key buttons.
func (g *game) keyColorOptions() bool {
	if len(g.optionLabels) == 0 {
		return false
	}
	for _, label := range g.optionLabels {
		if keyColorByName(label) == engine.KeyColorNone {
			return false
		}
	}
	return true
}

func (g *game) actionBar() app.UI {
	switch {
	case g.phase == phaseFlank:
		return app.Div().Class("btn-col").Body(
			app.Div().Class("section-title").Text("Play "+g.g.Def(g.sel).Name+" to which flank?"),
			btn("Left flank", g.playFlank(true), "btn-primary"),
			btn("Right flank", g.playFlank(false), "btn-primary"),
			btn("Cancel", g.cancelTargeting, "btn-secondary"),
		)
	case g.phase == phaseFightTarget:
		return app.Div().Class("btn-col").Body(
			app.Div().Class("prompt").Text("Pick an enemy creature to fight"),
			btn("Cancel", g.cancelTargeting, "btn-secondary"),
		)
	case !g.hasSel:
		return app.Div().Class("hint").Text("Select a card to inspect it and choose an action.")
	}

	switch g.selKind {
	case selHand:
		return g.handActions()
	case selYourCreature:
		return g.creatureActions()
	case selYourArtifact:
		return g.artifactActions()
	default:
		return app.Div().Class("hint").Text("Read-only — this is your opponent's card.")
	}
}

func (g *game) handActions() app.UI {
	def := g.g.Def(g.sel)
	items := []app.UI{app.Div().Class("section-title").Text(def.Name)}
	if err := g.g.CanPlay(g.active(), g.sel); err != nil {
		items = append(items, app.Div().Class("hint").Text("Cannot play: "+err.Error()+"."))
	} else {
		items = append(items, btn("Play", g.play, "btn-primary"))
	}
	// Discarding needs only that the card is of the active house.
	if h := g.g.State.ActiveHouse; h == engine.HouseNone || def.House == h {
		items = append(items, btn("Discard", g.discard, "btn-danger"))
	}
	return app.Div().Class("btn-col").Body(items...)
}

func (g *game) creatureActions() app.UI {
	if err := g.g.CanUse(g.active(), g.sel); err != nil {
		return app.Div().Class("btn-col").Body(
			app.Div().Class("section-title").Text(g.g.Def(g.sel).Name),
			app.Div().Class("hint").Text("Cannot act: "+err.Error()+"."),
		)
	}
	// A stunned creature recovers from stun instead of acting, so any use just
	// removes the stun: offer a single Unstun rather than Reap/Fight/Action.
	if g.g.Stunned(g.sel) {
		return app.Div().Class("btn-col").Body(
			app.Div().Class("section-title").Text(g.g.Def(g.sel).Name),
			app.Div().Class("hint").Text("Stunned"),
			btn("Unstun", g.reap, "btn-primary"),
		)
	}
	items := []app.UI{
		app.Div().Class("section-title").Text(g.g.Def(g.sel).Name),
		btn("Reap", g.reap, "btn-warning"),
	}
	// Fight is offered only when a legal target exists (e.g. with no enemy
	// creatures, a ready Valdr can still reap but has nothing to fight).
	if len(g.g.FightTargets(g.active(), g.sel)) > 0 {
		items = append(items, btn("Fight", g.startFight, "btn-danger"))
	}
	if g.g.HasAction(g.sel) {
		items = append(items, btn("Action", g.useAction, "btn-primary"))
	}
	return app.Div().Class("btn-col").Body(items...)
}

func (g *game) artifactActions() app.UI {
	name := g.g.Def(g.sel).Name
	if !g.g.HasAction(g.sel) {
		return app.Div().Class("hint").Text(name + " has no action ability.")
	}
	items := []app.UI{app.Div().Class("section-title").Text(name)}
	if g.g.Exhausted(g.sel) {
		items = append(items, app.Div().Class("hint").Text("Exhausted."))
	} else {
		items = append(items, btn("Action", g.useAction, "btn-primary"))
	}
	return app.Div().Class("btn-col").Body(items...)
}

func (g *game) logPanel() app.UI {
	groups := g.logGroupViews()
	return app.Div().Class("log").Body(
		app.Div().Class("panel-label").Text("Log"),
		app.Div().Class("log-list").ID("gamelog").Body(
			app.Range(groups).Slice(func(i int) app.UI {
				grp := groups[i]
				cls := cx("log-group",
					ifCls(grp.player == 0, "log-group--p0"),
					ifCls(grp.player == 1, "log-group--p1"),
					ifCls(grp.newest, "log-group--new"),
				)
				return app.Div().Class(cls).Body(
					app.Range(grp.lines).Slice(func(j int) app.UI {
						return app.Div().Class("log-line").Body(g.logSegments(grp.lines[j])...)
					}),
				)
			}),
		),
	)
}

// logGroupView is one rendered log bubble: the lines of a single root action, the
// player whose turn produced it (-1 for setup), and whether it is the newest.
type logGroupView struct {
	lines  []string
	player int
	newest bool
}

// logSeg is one root action's slice of the log (half-open [start, end)) and whose
// turn recorded it (-1 for the leading setup lines).
type logSeg struct {
	start, end, player int
}

// actionSegments returns the log ranges of each root action (from logGroups) plus
// a leading setup range, each clamped to the current log length.
func (g *game) actionSegments() []logSeg {
	log := g.g.Log
	var segs []logSeg
	first := len(log)
	if len(g.logGroups) > 0 && g.logGroups[0].Start < first {
		first = g.logGroups[0].Start
	}
	if first > 0 {
		segs = append(segs, logSeg{0, first, -1})
	}
	for i, m := range g.logGroups {
		start, end := m.Start, len(log)
		if i+1 < len(g.logGroups) {
			end = g.logGroups[i+1].Start
		}
		if start < 0 {
			start = 0
		}
		if end > len(log) {
			end = len(log)
		}
		if start < end {
			segs = append(segs, logSeg{start, end, m.Player})
		}
	}
	return segs
}

// turnBeginPlayer returns the player a "--- X begins turn N ---" line announces,
// or -1 if the line is not a turn header.
func (g *game) turnBeginPlayer(line string) int {
	if !strings.HasPrefix(line, "--- ") || !strings.Contains(line, "begins turn") {
		return -1
	}
	for p := 0; p < 2; p++ {
		if strings.Contains(line, g.g.PlayerName(p)+" begins turn") {
			return p
		}
	}
	return -1
}

// logGroupViews splits the flat log into per-action bubbles using logGroups, then
// further splits each bubble so a "begins turn" line starts a fresh bubble tinted
// for the new player. Lines before the first action form a leading setup bubble.
func (g *game) logGroupViews() []logGroupView {
	log := g.g.Log
	var out []logGroupView
	emit := func(lines []string, player int) {
		if len(lines) > 0 {
			out = append(out, logGroupView{lines: lines, player: player})
		}
	}
	for _, seg := range g.actionSegments() {
		player, lineStart := seg.player, seg.start
		for i := seg.start; i < seg.end; i++ {
			p := g.turnBeginPlayer(log[i])
			if p < 0 {
				continue
			}
			if i > lineStart { // the turn header opens a new bubble
				emit(log[lineStart:i], player)
				lineStart = i
			}
			player = p
		}
		emit(log[lineStart:seg.end], player)
	}
	if len(out) > 0 {
		out[len(out)-1].newest = true
	}
	return out
}

// logSegments splits a log line into plain text and clickable spans for the card
// names it mentions, so a mentioned card can be opened in the detail panel
// without changing the engine's log strings.
func (g *game) logSegments(line string) []app.UI {
	var out []app.UI
	var plain strings.Builder
	flush := func() {
		if plain.Len() > 0 {
			out = append(out, app.Text(plain.String()))
			plain.Reset()
		}
	}
	for i := 0; i < len(line); {
		if name := g.cardNameAt(line, i); name != "" {
			flush()
			out = append(out, app.Span().
				Class("log-card").
				DataSet("card", name).
				OnMouseEnter(g.onLogCardHover).
				OnMouseLeave(g.onCardHoverOut).
				Text(name))
			i += len(name)
			continue
		}
		// Put the Æmber icon before the word wherever the log mentions it.
		if strings.HasPrefix(line[i:], "Æmber") {
			flush()
			out = append(out, icon("aember", "icon-inline"), app.Text("Æmber"))
			i += len("Æmber")
			continue
		}
		plain.WriteByte(line[i])
		i++
	}
	flush()
	return out
}

// cardNameAt returns the longest known card name that begins at line[i] on word
// boundaries, or "" if none.
func (g *game) cardNameAt(line string, i int) string {
	if i > 0 && isWordByte(line[i-1]) {
		return "" // in the middle of a word — not a name boundary
	}
	best := ""
	for name := range g.defByName {
		if len(name) <= len(best) || !strings.HasPrefix(line[i:], name) {
			continue
		}
		if end := i + len(name); end < len(line) && isWordByte(line[end]) {
			continue // the name is only a prefix of a longer word
		}
		best = name
	}
	return best
}

func isWordByte(b byte) bool {
	return b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' || b >= '0' && b <= '9'
}

// ---- overlays ----

func (g *game) overBanner() app.UI {
	winner := g.g.Winner()
	return app.Div().Class("over-backdrop").Body(
		app.Div().Class("over-card").Body(
			app.Div().Class("over-title").Text(g.g.PlayerName(winner)+" wins!"),
			app.Button().Class("btn-primary").Text("New game").OnClick(g.restart),
		),
	)
}

// zonesOverlay shows a single player's zones — their deck, discard, archives,
// and purge piles — as read-only card strips, so cards outside the board can be
// inspected. It is a modal dismissed by the ✕ in the corner or by clicking the
// backdrop outside the panel.
func (g *game) zonesOverlay() app.UI {
	p := g.zonesPlayer
	return app.Div().Class("over-backdrop").OnClick(g.closeZones).Body(
		app.Div().Class("zones-panel").OnClick(g.stopClick).Body(
			app.Div().Class("zones-header").Body(
				app.Button().Class("zones-close").Text("✕").OnClick(g.closeZones),
				app.Div().Class("over-title").Text(g.g.PlayerName(p)+"'s Zones"),
			),
			app.Div().Class("zones-body").Body(
				g.zoneRow("Deck", g.sortByHouseTypeName(g.g.Deck(p))),
				g.zoneRow("Discard", g.g.Discard(p)),
				g.zoneRow("Archives", g.g.Archives(p)),
				g.zoneRow("Purge", g.g.Purge(p)),
			),
		),
	)
}

func (g *game) zoneRow(label string, ids []engine.LocalID) app.UI {
	return app.Div().Class("zone-row").Body(
		app.Div().Class("row-label").Text(fmt.Sprintf("%s (%d)", label, len(ids))),
		app.If(len(ids) == 0, func() app.UI {
			return app.Div().Class("row-empty").Text("—")
		}).Else(func() app.UI {
			return app.Div().Class("card-strip").Body(
				app.Range(ids).Slice(func(i int) app.UI { return g.renderZoneCard(ids[i]) }),
			)
		}),
	)
}

// renderZoneCard renders a read-only card face for a card in an out-of-play zone.
func (g *game) renderZoneCard(id engine.LocalID) app.UI {
	def := g.g.Def(id)
	// In manual mode a zone card can be selected to move it (e.g. deck → hand).
	var activate func(app.Context, engine.LocalID)
	if g.g.Manual() {
		activate = g.selectZoneCard
	}
	return &cardView{
		ID:         id,
		Title:      def.Name,
		HouseCls:   houseClasses(def.House),
		Emblem:     houseIconName(def.House),
		TypeIcon:   typeIconName(def.Type),
		Stat:       handStat(def),
		Rules:      engine.RenderCardRules(def),
		Kind:       string(def.Type),
		OnActivate: activate,
	}
}

// ---- small helpers ----

func btn(label string, h app.EventHandler, class string) app.UI {
	return app.Button().Class(class).Text(label).OnClick(h)
}

func (g *game) statLine(id engine.LocalID) []app.UI {
	def := g.g.Def(id)
	f := g.flashes[id]
	var segs []app.UI
	if def.Type == engine.Creature {
		segs = append(segs, statSeg(g.g.Power(id), "power", pulseClass(f.power, f.odd, "pow")))
		if d := g.g.Damage(id); d > 0 {
			segs = append(segs, statSeg(d, "damage", pulseClass(f.damage, f.odd, "dmg")))
		}
		if a := g.g.Armor(id); a > 0 {
			segs = append(segs, statSeg(a, "shield"))
		}
	}
	if a := g.g.AmberOn(id); a > 0 {
		segs = append(segs, statSeg(a, "aember", pulseClass(f.amber, f.odd, "gain")))
	}
	// Stun shows as a token on the face (see cardView); exhaustion shows as an icon.
	if g.g.Exhausted(id) {
		segs = append(segs, icon("exhausted", "icon-stat",
			ifCls(f.exhaust && !f.odd, "icon--pulse-a"),
			ifCls(f.exhaust && f.odd, "icon--pulse-b")))
	}
	return segs
}

// pulseClass returns the alternating one-shot animation class for a stat segment
// (or "" when it is not flashing). kind selects the colour: dmg (red), pow (cyan),
// gain (gold). The -a/-b pair alternates so the animation replays on repeats.
func pulseClass(on, odd bool, kind string) string {
	switch {
	case !on:
		return ""
	case odd:
		return "stat-seg--" + kind + "-b"
	default:
		return "stat-seg--" + kind + "-a"
	}
}

func handStat(def *engine.CardDefinition) []app.UI {
	var segs []app.UI
	if def.Type == engine.Creature {
		segs = append(segs, statSeg(def.Power, "power"))
		if def.Armor > 0 {
			segs = append(segs, statSeg(def.Armor, "shield"))
		}
	}
	if def.AemberBonus > 0 {
		segs = append(segs, statSeg(def.AemberBonus, "aember"))
	}
	return segs
}

func (g *game) faceRules(id engine.LocalID) string {
	var lines []string
	if s := engine.RenderCardRules(g.g.Def(id)); s != "" {
		lines = append(lines, s)
	}
	for _, up := range g.g.Upgrades(id) {
		lines = append(lines, "↳ "+g.g.Def(up).Name)
	}
	return strings.Join(lines, "\n")
}

// playableFromHand reports whether the active player can play the given hand card
// right now, so unplayable cards are dimmed and not draggable.
func (g *game) playableFromHand(id engine.LocalID) bool {
	return g.g.CanPlay(g.active(), id) == nil
}

func (g *game) isEnemyCreature(id engine.LocalID) bool {
	if g.g.Def(id).Type != engine.Creature {
		return false
	}
	for _, c := range g.g.Battleline(1 - g.active()) {
		if c == id {
			return true
		}
	}
	return false
}

func containsID(ids []engine.LocalID, id engine.LocalID) bool {
	for _, x := range ids {
		if x == id {
			return true
		}
	}
	return false
}

func indexOfID(ids []engine.LocalID, id engine.LocalID) int {
	for i, x := range ids {
		if x == id {
			return i
		}
	}
	return -1
}
