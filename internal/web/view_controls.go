package web

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/cards"
	"github.com/dmikalova/vactrol/internal/engine"
)

// This file draws the sidebar: everything the player acts through rather than
// looks at — the prompts (card picker, house picker, option chooser), the action
// bar for the selected card, and the manual-mode panel.

// promptSourceHeader shows the face of the card driving the current prompt, so
// the player can read the ability they are resolving without hunting for the card
// on the board. The card's name is not repeated here — it rides inside the prompt
// line itself as a green token (see promptLine) — so the face stands alone on
// desktop and is dropped for room on mobile, where the token is the only source
// affordance.
func (g *game) promptSourceHeader() app.UI {
	return app.If(g.promptSource != "", func() app.UI {
		def := g.defByName[g.promptSource]
		if def == nil {
			return app.Div()
		}
		house, changed := g.promptSourceHouse(def)
		return app.Div().Class("prompt-source-block").Body(
			app.Div().Class("prompt-card").Body(&cardView{
				Title:        def.Name,
				HouseCls:     houseClasses(house),
				Emblem:       houseIconName(house),
				HouseChanged: changed,
				TypeIcon:     typeIconName(def.Type),
				Stat:         handStat(def),
				Rules:        displayRules(faceText(def, engine.RenderCardRules(def))),
				Kind:         kindLabel(def),
				Trait:        traitLabel(def),
				Rarity:       rarityMarkOf(def.Rarity),
				Bonuses:      def.Bonuses,
			}),
		)
	})
}

// promptLine draws a prompt's text, turning every mention of the card driving the
// prompt into a green log-card token — tap or hover to enlarge the card, the same
// affordance a log mention gives — so the prompt reads as one line ("Exalt
// <Centurion Stenopius>") instead of repeating the name as a heading above it.
// It carries the full text as a data-tip so the mobile dock, which clips a long
// prompt to one line (app.css), can still surface the whole sentence on hover/tap.
func (g *game) promptLine(text string) app.UI {
	return app.Div().Class("prompt").DataSet("tip", text).Body(g.promptTextSegments(text)...)
}

// promptTextSegments splits a prompt's text around each mention of its source
// card, wrapping every mention in a log-card token and leaving the rest plain.
// With no source card, or none named in the text, the whole line is plain text.
func (g *game) promptTextSegments(text string) []app.UI {
	name := g.promptSource
	if name == "" || !strings.Contains(text, name) {
		return []app.UI{app.Text(text)}
	}
	var out []app.UI
	for {
		i := strings.Index(text, name)
		if i < 0 {
			if text != "" {
				out = append(out, app.Text(text))
			}
			return out
		}
		if i > 0 {
			out = append(out, app.Text(text[:i]))
		}
		out = append(out, app.Span().
			Class("log-card").
			DataSet("card", name).
			OnMouseEnter(g.onLogCardHover).
			OnMouseLeave(g.onCardHoverOut).
			OnClick(g.onLogCardTap).
			Text(name))
		text = text[i+len(name):]
	}
}

// promptCardButtons lists a bounded out-of-play card prompt's candidates as
// action-bar buttons — a "look at the top N cards" pick (Navigator Ali, Lay of
// the Land) reads as a few named buttons instead of a modal over the deck. Each
// button names its card, previews it on hover (the same preview a log mention
// opens), and answers the prompt on click. A reorder-the-top prompt adds the note
// that the first card picked ends up on the bottom, since each pick is placed on
// top of the one before it and the card left over rides on top.
func (g *game) promptCardButtons() app.UI {
	body := make([]app.UI, 0, len(g.chooserCandidates)+1)
	// The reorder prompt places each pick on top of the previous, so the first pick
	// finishes deepest; the note keys off that prompt's wording rather than the card.
	if strings.Contains(g.chooserPrompt, "on top") {
		body = append(body,
			app.Div().Class("hint").Text("The first card you pick ends up on the bottom."))
	}
	for _, id := range g.chooserCandidates {
		def := g.g.Def(id)
		cursor := ifCls(g.hasCursor && g.promptCursor == id, "btn-cursor")
		body = append(body, app.Button().
			Class(cx("btn-secondary", "prompt-pick", cursor)).
			DataSet("id", strconv.Itoa(int(id))).
			DataSet("card", def.Name).
			OnMouseEnter(g.onLogCardHover).
			OnMouseLeave(g.onCardHoverOut).
			OnClick(g.onPromptButtonPick).
			Body(
				houseIcon(g.g.House(id), "icon-inline"),
				app.Span().Class("prompt-pick-name").Text(def.Name),
			))
	}
	return app.Div().Class("btn-col", "prompt-picks").Body(body...)
}

// promptSourceHouse looks up the live house of the card driving a prompt, so a
// maverick card (played out of its printed house) shows the house it is actually
// resolving as instead of the one printed on it — the engine's Chooser only names
// the source by its card name, with no id, so this matches by name among the
// cards actually in play. A source with no match in play (its effect fires from
// hand, discard, or another zone) falls back to the printed house.
func (g *game) promptSourceHouse(def *engine.CardDefinition) (house engine.House, changed bool) {
	for p := range 2 {
		for _, ids := range [][]engine.LocalID{g.g.Battleline(p), g.g.Artifacts(p)} {
			for _, id := range ids {
				if g.g.Def(id).Name != def.Name {
					continue
				}
				h := g.g.House(id)
				return h, h != def.House
			}
		}
	}
	return def.House, false
}

// endTurnButton is the End turn control. Once a confirm is armed (the player could
// still act this turn) it becomes a red "Confirm end turn" that ends on the next
// click, mirroring pressing E a second time.
func (g *game) endTurnButton() app.UI {
	cursor := ifCls(g.isEndTurnCursor(), "btn-cursor")
	if g.confirmEndTurn {
		return btn("Confirm end turn", g.endTurn, cx("btn-danger", cursor))
	}
	// With nothing left to do, the button fades green to invite ending the turn.
	ready := ifCls(!g.hasMoves(), "btn-endturn--ready")
	return btn("End turn", g.endTurn, cx("btn-secondary", ready, cursor))
}

// undoIcon is the icon-only Undo that rides in the top-right of the action area,
// so a card played by mistake is one click from being taken back without opening
// the menu. It is disabled when there is nothing to undo.
func (g *game) undoIcon() app.UI {
	return app.Button().
		Class(cx("btn-secondary", "btn-icon")).
		Title("Undo").
		Disabled(!g.canUndo()).
		OnClick(g.undoAction).
		Body(icon("undo", "icon-nav"))
}

// actionHeader is the top row of the action area: an optional section title on the
// left and the Undo icon pinned to the right, with any extra right-side icons (the
// manual-mode wrench) sitting beside Undo.
func (g *game) actionHeader(title string, extra ...app.UI) app.UI {
	left := app.UI(app.Div().Class("action-head-spacer"))
	if title != "" {
		left = app.Div().Class("section-title").Text(title)
	}
	return g.headerRow(left, extra)
}

// promptHeader is actionHeader with the prompt itself as the left slot, so the
// question and the Undo that backs it out share one row instead of stacking.
func (g *game) promptHeader(text string, extra ...app.UI) app.UI {
	return g.headerRow(g.promptLine(text), extra)
}

// headerRow lays out the action area's top row: left slot, then Undo and any extra
// icons hard right.
func (g *game) headerRow(left app.UI, extra []app.UI) app.UI {
	icons := make([]app.UI, 0, len(extra)+1)
	icons = append(icons, extra...)
	icons = append(icons, g.undoIcon())
	return app.Div().Class("action-head").Body(
		left,
		app.Div().Class("action-head-icons").Body(icons...),
	)
}

// endTurnBar is the resting End turn control with the Undo icon in the top-right
// above it. In manual mode a green wrench rides beside Undo, so the mode that is
// on is one click from off without opening the menu.
func (g *game) endTurnBar() app.UI {
	var extra []app.UI
	if g.g.Manual() {
		extra = append(extra, app.Button().
			Class(cx("btn-nav", "btn-icon", "btn-nav-on")).
			Title("Manual mode is on — click to turn it off").
			OnClick(g.toggleManual).
			Body(icon("wrench", "icon-nav")))
	}
	return app.Div().Class("btn-col", "end-turn-bar").Body(
		g.actionHeader("", extra...),
		g.endTurnButton(),
	)
}

// disabledEndTurnBar draws the resting Undo header + a greyed-out, non-clickable
// End turn. It fills the dock during a mid-action step (placing a creature) whose
// own controls live on the lifted card, so the dock reads as the controls paused
// rather than a blank box. Undo stays live, since it backs the whole step out.
func (g *game) disabledEndTurnBar() app.UI {
	end := app.Button().Class("btn-secondary").Disabled(true).Text("End turn")
	return app.Div().Class("btn-col", "end-turn-bar").Body(g.actionHeader(""), end)
}

// controls is the bottom of the sidebar: the contextual controls (house picker or
// action bar) plus End turn. House selection has no End turn — a house must be
// chosen first.
func (g *game) controls() app.UI {
	// A pending new game takes over the controls with the set picker until sets are
	// chosen or it is cancelled — even over a finished game, which is how "New game"
	// on the end-of-game panel opens the picker rather than redrawing the result.
	if g.awaitingSetup {
		return app.Div().Class("controls").Body(g.setChooser())
	}
	// A finished game replaces every control with the result: nothing else can be
	// done, and a modal over the board would hide the position that ended it.
	if g.phase == phaseOver {
		return app.Div().Class("controls").Body(g.overPanel())
	}
	// A pending manual key forge takes over the controls until a colour is picked
	// or cancelled; once forgingKey resets, the previous buttons return on their own.
	if g.forgingKey >= 0 {
		return app.Div().Class("controls").Body(g.keyForgePanel())
	}
	// An engine prompt takes over the dock. Every prompt kind routes through
	// promptControls, the single seam a totality test proves renders them all
	// (ADR 0045). The option and position kinds are dispatched here, before the
	// manual Graft / Place-under target below; the card kinds after it, so their
	// prior dock priority is unchanged.
	if g.choosingOption {
		ui, _ := g.promptControls(engine.PromptOption)
		return ui
	}
	if g.choosingPosition {
		ui, _ := g.promptControls(engine.PromptPosition)
		return ui
	}
	// A manual Graft / Place under is waiting for the player to click the in-play
	// host to thread the selected card under; the board lights its hosts and the
	// dock shows the pick's prompt with a Cancel.
	if g.hostTargeting {
		return g.hostTargetingControls()
	}
	// While an engine chooser waits, the controls become the prompt itself: a
	// green call to action to click one of the highlighted cards.
	if g.choosing {
		ui, _ := g.promptControls(g.cardPromptKind())
		return ui
	}
	if g.phase == phaseHouse {
		if g.g.Manual() {
			return app.Div().Class("controls").Body(g.manualPanel(), g.housePicker())
		}
		return app.Div().Class("controls").Body(g.housePicker())
	}
	// Mid-action selections (choosing a flank or a fight target) show only their
	// own controls — End turn is withheld until the action finishes or is cancelled.
	// A flank is asked on the lifted card, so for that one the dock is simply empty.
	if g.phase == phaseFlank || g.phase == phaseFightTarget {
		return app.Div().Class("controls").Body(g.targetingPrompt())
	}
	return g.restingControls()
}

// hostTargetingControls draws the dock for a manual Graft / Place under waiting
// for the player to click the in-play host to thread the selected card under.
func (g *game) hostTargetingControls() app.UI {
	verb := "graft"
	if g.hostFaceDown {
		verb = "place"
	}
	return app.Div().Class("controls").Body(
		app.Div().Class("btn-col").Body(
			app.Div().Class("prompt").Text(
				"Click a card to "+verb+" "+g.g.Def(g.sel).Name+" under it"),
			btn("Cancel", g.cancelHostTargeting, "btn-secondary"),
		),
	)
}

// restingControls is the default dock when no prompt is up: End turn, with the
// manual-mode panel above it when manual mode is on. A selected card's verbs are
// drawn on the card itself, so nothing here competes with End turn for the dock.
func (g *game) restingControls() app.UI {
	body := []app.UI{g.endTurnBar()}
	if g.g.Manual() {
		body = append([]app.UI{g.manualPanel()}, body...)
	}
	return app.Div().Class("controls").Body(body...)
}

// promptControls renders the dock cluster for one engine prompt kind, returning
// the widget and whether the kind is covered. It is the single seam a totality
// test drives to prove every engine.PromptKind has a rendering (ADR 0045): a new
// prompt route is a new PromptKind and a case here together, so the test fails
// until the case is added. controls dispatches its live prompt branches through
// this, so the mapping the test checks is the one the client actually renders.
func (g *game) promptControls(kind engine.PromptKind) (app.UI, bool) {
	switch kind {
	case engine.PromptOption:
		return g.optionPromptControls(), true
	case engine.PromptPosition:
		return g.positionPromptControls(), true
	// The card kinds share one cluster; the modifier flags (declinable, ordering,
	// as-buttons) it reads decide which extra buttons it grows.
	case engine.PromptCreature,
		engine.PromptCardOrDecline,
		engine.PromptReaction,
		engine.PromptOrder:
		return g.cardPromptControls(), true
	case engine.PromptBadge:
		// A badge preview decorates the candidate card it previews; it grows no dock
		// cluster of its own, so it is covered by rendering nothing here.
		return nil, true
	}
	return nil, false
}

// cardPromptKind names which card-prompt kind is up, from the modifier flags the
// raised prompt set. All card kinds render through cardPromptControls, so this is
// for routing and honest naming, not a rendering fork.
func (g *game) cardPromptKind() engine.PromptKind {
	switch {
	case g.chooserOrdering:
		return engine.PromptOrder
	case g.chooserDeclinable:
		return engine.PromptCardOrDecline
	default:
		return engine.PromptCreature
	}
}

// optionPromptControls draws a labeled option prompt (e.g. "take archives?") as
// its choices.
func (g *game) optionPromptControls() app.UI {
	// A "choose how to use X" verb prompt lifts the chosen creature and draws its
	// use buttons on it (like an ordinary use), so the dock shows the resting
	// controls disabled rather than repeating the buttons in the sidebar.
	if _, ok := g.liftUseTarget(); ok {
		return app.Div().Class("controls").Body(g.disabledEndTurnBar())
	}
	body := []app.UI{
		g.promptHeader(g.optionPrompt),
		g.promptSourceHeader(),
		g.optionChooser(),
	}
	// Manual mode adds a Cancel that backs the whole action out — an option prompt
	// has no decline of its own, so this is the only way out of a stuck one.
	if g.g.Manual() {
		body = append(body, btn("Cancel", g.cancelChooser, "btn-secondary"))
	}
	return app.Div().Class("controls").Body(body...)
}

// positionPromptControls draws the dock while a Deploy creature is being placed:
// the battleline lights as click targets and the where question rides on the
// lifted card (deployActions), exactly like the flank question, so the dock shows
// the resting controls disabled rather than a panel of its own.
func (g *game) positionPromptControls() app.UI {
	return app.Div().Class("controls").Body(g.disabledEndTurnBar())
}

// cardPromptControls draws a card decision: the controls become the prompt
// itself, a green call to action to click one of the highlighted cards, plus the
// buttons the modifier flags add.
func (g *game) cardPromptControls() app.UI {
	body := []app.UI{
		g.promptHeader(g.chooserPrompt),
		g.promptSourceHeader(),
	}
	// A bounded out-of-play pick (look at the top N cards) lists its candidates as
	// buttons here rather than opening the zone viewer, so the whole choice reads
	// in the action bar.
	if g.promptAsButtons {
		body = append(body, g.promptCardButtons())
	}
	// An optional prompt ("you may", "up to N") is passed on with Done. Manual mode
	// adds a Cancel on every prompt — optional or mandatory — that backs the whole
	// action out, the way out of a prompt with no clickable candidate.
	if g.chooserDeclinable {
		body = append(body, btn("Done", g.declineChooser,
			cx("btn-primary", ifCls(g.isDoneCursor(), "btn-cursor"))))
	}
	// An ordering prompt offers Auto-resolve, which answers with a random order
	// so the player need not arrange abilities whose order does not matter to them.
	if g.chooserOrdering {
		body = append(body, btn("Auto-resolve", g.autoResolveOrder, "btn-secondary"))
	}
	if g.g.Manual() {
		body = append(body, btn("Cancel", g.cancelChooser, "btn-secondary"))
	}
	return app.Div().Class("controls").Body(body...)
}

// setChooser is the new-game set picker drawn in the action bar. When a previous
// game exists it leads with a same-sets shortcut so a rematch is one click, then
// asks each player which set to play; Cancel leaves the current game running.
func (g *game) setChooser() app.UI {
	names := cards.DeckSetNames()
	body := []app.UI{app.Div().Class("section-title").Text("New game")}
	if g.hasPrevSets() {
		body = append(body, btn("Same sets — "+g.prevSetLabel(),
			func(ctx app.Context, _ app.Event) { g.continueSameSets(ctx) }, "btn-primary"))
	}
	body = append(body,
		// "Player N" wears that player's colour like every other place a player is
		// named; the rest of the line is ordinary prompt text.
		app.Div().Class("prompt").Body(
			app.Span().
				Class(playerNameCls(g.setPick)).
				Text(fmt.Sprintf("Player %d", g.setPick+1)),
			app.Text(" — choose a set"),
		),
	)
	for _, name := range names {
		body = append(body, app.Button().
			Class(cx("btn-secondary", "set-btn", setAccent(name))).
			OnClick(func(ctx app.Context, _ app.Event) { g.pickSet(ctx, name) }).
			Body(app.Span().Class("set-emblem"), app.Text(name)))
	}
	// Cancel only makes sense when there is a running game to fall back to. On a
	// first-time load the picker is the whole screen with nothing behind it, so
	// there is nothing to cancel to.
	if g.g != nil {
		body = append(body, btn("Cancel", g.cancelSetup, "btn-secondary"))
	}
	return app.Div().Class("btn-col", "set-pick").Body(body...)
}

// setupScreen fills the viewport with the new-game set picker for a first-time
// load, before any match is dealt. It reuses the action-bar picker (setChooser)
// centered on its own so the player chooses their sets rather than being dropped
// into a silently dealt base-set game.
func (g *game) setupScreen() app.UI {
	return app.Div().Class("setup-screen").Body(
		app.Div().Class("setup-panel").Body(g.setChooser()),
		// The shortcut sheet is a fixed overlay, so it works here before any match
		// exists too — ? opens it on the very first load.
		app.If(g.keysOpen, func() app.UI { return g.keysOverlay() }),
	)
}

// manualPanel is the manual-mode control block: add an arbitrary card, and (for a
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
		// An attached card (upgrade or under-card) can only leave to hand by
		// detaching first, so it gets a To hand of its own rather than the Move
		// buttons' Hand, which would leave it attached and duplicated.
		if g.isAttached(g.sel) {
			items = append(items, btn("To hand", g.manualToHand, "btn-secondary"))
		}
		items = append(items,
			app.Div().Class("hint").Text("Move "+name+" to:"),
			g.moveButtons(),
			app.Div().Class("hint").Text("Thread "+name+" under a host:"),
			app.Div().Class("btn-wrap").Body(
				btn("Graft", g.manualGraft, "btn-mini"),
				btn("Place under", g.manualPlaceUnder, "btn-mini"),
			),
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

// pickerInputID marks the card picker's search box, so opening the picker can put
// the caret straight in it.
const pickerInputID = "pickerinput"

// cardPicker is the fuzzy, text-only card picker for adding an arbitrary card
// from the pool to hand. It filters the pool by a case-insensitive name substring.
func (g *game) cardPicker() app.UI {
	matches := g.pickerMatches()
	// Answering a name-a-card prompt, the panel carries the engine's own prompt and
	// has no close: the blocked effect is waiting on a name.
	title := "Add a card to hand"
	if g.pickerNaming {
		title = g.optionPrompt
	}
	return app.Div().Class("over-backdrop").OnClick(g.closePicker).Body(
		app.Div().Class("picker-panel").OnClick(g.stopClick).Body(
			app.If(!g.pickerNaming, func() app.UI {
				return app.Button().Class("zones-close").Text("✕").OnClick(g.closePicker)
			}),
			app.Div().Class("over-title").Text(title),
			app.Input().ID(pickerInputID).Class("picker-input").Type("text").
				AutoFocus(true).
				Placeholder("Search cards…").
				Value(g.pickerQuery).OnInput(g.pickerInput),
			app.Div().Class("picker-list").Body(
				app.Range(matches).Slice(func(i int) app.UI {
					d := matches[i]
					return app.Button().
						Class(cx("picker-item", ifCls(i == g.pickerCursor, "picker-cursor"))).
						DataSet("card", d.Name).
						OnClick(g.addPickedCard).
						Body(
							houseIcon(d.House, "icon-inline"),
							app.Span().Class("picker-name").Text(d.Name),
							app.Span().Class("picker-kind").Text(d.Type.String()),
						)
				}),
			),
		),
	)
}

// pickableHouses is the set of houses the active player may choose from, sorted by
// house name so the choice reads the same every turn. In ordinary play it is the
// engine's allowed set — the delayed constraint table already resolved, so a forced
// house shows alone and a barred house is absent rather than shown and then rejected
// (ADR 0035). When the allowed set is empty the player has no active house; the
// picker offers a single "No House" choice (see houseButtons). In manual mode the
// operator drives both sides, so every deck house is offered.
func (g *game) pickableHouses() []engine.House {
	p := g.active()
	houses := g.deckHouses[p]
	if !g.g.Manual() {
		houses = g.g.AllowedHouses(p)
	}
	sorted := append([]engine.House(nil), houses...)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].String() < sorted[j].String()
	})
	return sorted
}

// houseButtons is the buttons the house picker draws, in order. It is the pickable
// houses, or — when no house is choosable — a single HouseNone entry rendered as the
// "No House" choice, so the phase always has a button to advance it. A full lockout
// (the constraint table bars every choosable house) leaves No House the only legal
// choice in manual mode too, even though manual play otherwise offers every deck
// house: ChooseHouse still enforces the allowed set, so offering a barred house
// would only be rejected, and the operator can still force a house from the score
// pill (ManualSetActiveHouse) if they mean to override the rules (ADR 0035).
func (g *game) houseButtons() []engine.House {
	if len(g.g.AllowedHouses(g.active())) == 0 {
		return []engine.House{engine.HouseNone}
	}
	houses := g.pickableHouses()
	if len(houses) == 0 {
		return []engine.House{engine.HouseNone}
	}
	return houses
}

// housePicker offers the houses the active player may choose from, or a single "No
// House" choice when the constraint table leaves no house choosable (ADR 0035).
func (g *game) housePicker() app.UI {
	houses := g.houseButtons()
	return app.Div().Class("btn-col", "house-pick").Body(
		// Undo sits on the title row and steps out of the state that put the player at
		// house selection (the previous turn); it greys out on the first turn, when
		// there is nothing to step back to.
		g.actionHeader("Choose a house"),
		app.Range(houses).Slice(func(i int) app.UI {
			h := houses[i]
			if h == engine.HouseNone {
				return app.Button().
					Class(cx("house-btn", ifCls(g.isButtonCursor(i), "btn-cursor"))).
					OnClick(g.pickHouse(h)).
					Text("No House")
			}
			return app.Button().
				Class(cx("house-btn", houseAccent(h), ifCls(g.isButtonCursor(i), "btn-cursor"))).
				OnClick(g.pickHouse(h)).
				Body(houseIcon(h, "icon-inline"), app.Text(h.String()))
		}),
	)
}

// optionKind classifies a labeled option prompt by how its widget is drawn, so
// the option dispatch is an exhaustive switch over a named kind rather than an
// if-chain falling through to a silent catch-all (ADR 0045). The kind is derived
// from the option labels' shape — the seam that survives until ADR 0040's Request
// carries the kind from the engine — but optionGeneric is an explicit, named
// rendering, so a new option shape is a new kind and a new case together.
type optionKind uint8

const (
	// optionCardName is a "name a card" prompt answered in the typeahead panel, so
	// the dock shows no per-option widget.
	optionCardName optionKind = iota
	// optionKeyColor is a forge prompt shown as themed key-colour buttons.
	optionKeyColor
	// optionUseVerb is a reap/fight/action prompt shown as the standard use buttons.
	optionUseVerb
	// optionFlank is a move-to-a-flank prompt shown as the two flank buttons.
	optionFlank
	// optionHouse is a choose-a-house prompt shown as the house emblem grid.
	optionHouse
	// optionGeneric is any other labeled choice, shown as plain buttons.
	optionGeneric
)

// optionKinds returns every option-widget kind, in declaration order, for a
// totality test that proves each has a rendering.
func optionKinds() []optionKind {
	return []optionKind{
		optionCardName,
		optionKeyColor,
		optionUseVerb,
		optionFlank,
		optionHouse,
		optionGeneric,
	}
}

// classifyOptions names how the current option labels are drawn. optionGeneric is
// the honest default for a shape none of the specific classifiers claim — a named
// case, not an unhandled fall-through.
func (g *game) classifyOptions() optionKind {
	switch {
	case g.cardNameOptions():
		return optionCardName
	case g.keyColorOptions():
		return optionKeyColor
	case g.useVerbOptions():
		return optionUseVerb
	case g.flankOptions():
		return optionFlank
	case g.houseOptions():
		return optionHouse
	default:
		return optionGeneric
	}
}

// optionChooser renders the widget for a labeled multiple-choice prompt. The
// question itself rides in the header row beside Undo (promptHeader), so a card
// that words its question differently is not overwritten by a heading the client
// made up; only the widget varies.
func (g *game) optionChooser() app.UI {
	ui, _ := g.optionControls(g.classifyOptions())
	return ui
}

// optionControls draws the widget for one option kind, returning it and whether
// the kind is covered. It is the single seam a totality test drives to prove
// every optionKind has a rendering (ADR 0045).
func (g *game) optionControls(kind optionKind) (app.UI, bool) {
	switch kind {
	case optionCardName:
		// A name-a-card prompt is answered in the typeahead panel over the board, so
		// the dock shows no widget at all here — a button per card in the database is
		// not one.
		return app.Div().Class("btn-col"), true
	case optionKeyColor:
		// When every option is a key colour it shows themed key buttons.
		return app.Div().Class("btn-col").Body(
			app.Range(g.optionLabels).Slice(func(i int) app.UI {
				c := keyColorByName(g.optionLabels[i])
				return keyChoiceButton(
					c,
					g.optionLabels[i],
					g.isButtonCursor(i),
					g.chooseOptionIdx(i),
				)
			}),
		), true
	case optionUseVerb:
		// When every option is a way of using a creature (a reap/fight/action prompt
		// another card raised) it shows the standard use buttons, so a triggered use
		// reads like a chosen one.
		var body []app.UI
		for i, label := range g.optionLabels {
			// A verb the chosen creature cannot be used for (Narp bars its neighbors
			// from reaping) is omitted, so an illegal use is never offered.
			if g.useVerbBarred(label) {
				continue
			}
			k, _ := useVerbKindOfLabel(label)
			s := useVerbSpec(k)
			body = append(body, btn(s.text, g.chooseOptionIdx(i),
				cx(s.class, ifCls(g.isButtonCursor(i), "btn-cursor"))))
		}
		return app.Div().Class("btn-col").Body(body...), true
	case optionFlank:
		// A move-to-a-flank prompt (Reassembling Automaton, Harland Mindlock) uses
		// the same flank buttons as placing a creature.
		return app.Div().Class("btn-col").Body(
			btn(engine.FlankLeftLabel, g.chooseOptionIdx(0),
				cx("btn-primary", "btn-flank", "btn-flank--left",
					ifCls(g.isButtonCursor(0), "btn-cursor"))),
			btn(engine.FlankRightLabel, g.chooseOptionIdx(1),
				cx("btn-primary", "btn-flank", "btn-flank--right",
					ifCls(g.isButtonCursor(1), "btn-cursor"))),
		), true
	case optionHouse:
		// When every option is a house it shows a grid of house emblems — a house
		// prompt can offer all seven, and seven full-width rows push the rest of the
		// controls off the screen where a grid does not.
		return app.Div().Class("btn-col").Body(
			app.Div().Class("house-grid").Body(
				app.Range(g.optionLabels).Slice(func(i int) app.UI {
					h, _ := engine.ParseHouse(g.optionLabels[i])
					return app.Button().
						Class(cx("house-btn", "house-btn--icon", houseAccent(h),
							ifCls(g.isButtonCursor(i), "btn-cursor"))).
						Title(g.optionLabels[i]).
						OnClick(g.chooseOptionIdx(i)).
						Body(houseIcon(h, "icon-house"))
				}),
			),
		), true
	case optionGeneric:
		// Anything else falls back to plain primary buttons.
		return app.Div().Class("btn-col").Body(
			app.Range(g.optionLabels).Slice(func(i int) app.UI {
				// A declining "No" or a hand-shedding "Mulligan" is the
				// destructive-looking choice, so it reads red.
				kind := "btn-primary"
				if isDecliningOption(g.optionLabels[i]) {
					kind = "btn-danger"
				}
				return btn(g.optionLabels[i], g.chooseOptionIdx(i),
					cx(kind, ifCls(g.isButtonCursor(i), "btn-cursor")))
			}),
		), true
	}
	return nil, false
}

// isDecliningOption reports whether a label is the "turn this down" answer: the
// No of a yes/no question, or the Mulligan that sheds an opening hand. It is both
// what the n key answers and what reads red.
func isDecliningOption(label string) bool {
	return label == "No" || label == "Mulligan"
}

// flankOptions reports whether the current option labels are exactly the two
// battleline flanks, so a "move it to a flank" prompt (Reassembling Automaton,
// Harland Mindlock) is drawn with the same flank buttons as placing a creature.
func (g *game) flankOptions() bool {
	return len(g.optionLabels) == 2 &&
		g.optionLabels[0] == engine.FlankLeftLabel &&
		g.optionLabels[1] == engine.FlankRightLabel
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

// cardNameOptions reports whether every current option label names a card in the
// database, so the prompt is answered through the card-name typeahead instead of
// one button per option — a "name a card" prompt (Etan's Jar) offers the whole
// database, which is thousands of buttons. Two labels are not a card-name prompt
// even when both happen to be card names, so an ordinary Yes/No-shaped choice
// between two cards keeps its buttons.
func (g *game) cardNameOptions() bool {
	if len(g.optionLabels) <= 2 {
		return false
	}
	for _, label := range g.optionLabels {
		if g.defByName[label] == nil {
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

// useVerbOptions reports whether every current option label is a way of using a
// creature (reap, fight, use its action), so a reap/fight/action prompt another
// card raised is drawn as the standard use buttons rather than a plain option list.
func (g *game) useVerbOptions() bool {
	if len(g.optionLabels) == 0 {
		return false
	}
	for _, label := range g.optionLabels {
		if _, ok := useVerbKindOfLabel(label); !ok {
			return false
		}
	}
	return true
}

// liftUseTarget is the creature a "choose how to use X" verb prompt should lift,
// if one is up: the creature this action last chose to use, still in play. When it
// is known the use buttons go on that creature (like an ordinary use) rather than
// in the sidebar; when it is not (the target was auto-picked with no prompt) the
// sidebar option list stands in.
func (g *game) liftUseTarget() (engine.LocalID, bool) {
	if !g.choosingOption || !g.hasUseTarget || !g.useVerbOptions() {
		return 0, false
	}
	if !containsID(g.g.Battleline(0), g.useTarget) &&
		!containsID(g.g.Battleline(1), g.useTarget) {
		return 0, false
	}
	return g.useTarget, true
}

// useVerbCardActions draws the current use-verb option prompt as the standard
// reap/fight/action buttons on the lifted creature. Each button answers the
// engine's option prompt at that label's own index, so the buttons on the card
// drive the same choice the sidebar list would have. Manual mode adds a Cancel
// that backs the whole action out, the option prompt's only way out.
func (g *game) useVerbCardActions() []cardAction {
	acts := make([]cardAction, 0, len(g.optionLabels)+1)
	for i, label := range g.optionLabels {
		// A verb the chosen creature cannot be used for (Narp bars its neighbors from
		// reaping) is omitted, so Universal Translator never offers an illegal use.
		if g.useVerbBarred(label) {
			continue
		}
		k, _ := useVerbKindOfLabel(label)
		s := useVerbSpec(k)
		acts = append(acts, cardAction{
			s.text,
			cx(s.class, ifCls(g.isButtonCursor(i), "btn-cursor")),
			g.chooseOptionIdx(i),
		})
	}
	if g.g.Manual() {
		acts = append(acts, cardAction{"Cancel", "btn-secondary", g.cancelChooser})
	}
	return acts
}

// targetingPrompt is the dock's mid-action question: which enemy to fight. A
// card's own verbs are not here — they are drawn on the lifted copy of the card
// itself (cardFocus), which leaves the dock to the questions whose answer is some
// other card on the board.
func (g *game) targetingPrompt() app.UI {
	if g.phase == phaseFightTarget {
		return app.Div().Class("btn-col").Body(
			app.Div().Class("prompt").Text("Pick an enemy creature to fight"),
			btn("Cancel", g.cancelTargeting, "btn-secondary"),
		)
	}
	// Placing a creature asks its which-end question on the lifted card, so the dock
	// has no prompt of its own — it shows the resting controls disabled rather than
	// an empty box.
	return g.disabledEndTurnBar()
}

// cardAction is one verb the selected card offers. The dock's mid-action prompts
// and the lifted copy over the board both render this list, so what a card can do
// is decided once rather than once per place it is drawn.
type cardAction struct {
	Label string
	Class string
	On    app.EventHandler
}

// selActions is what the selected card can do right now, and the note explaining
// why it can do nothing. The handlers act on g.sel, which is why they can be
// plain methods: the buttons are only ever drawn for the card that is selected.
func (g *game) selActions() ([]cardAction, string) {
	if !g.hasSel {
		return nil, ""
	}
	// A "choose how to use X" verb prompt lifts the just-chosen creature and puts
	// its reap/fight/action buttons on it, so a use another card raised (Universal
	// Translator) reads like an ordinary use. This wins over the inert/boardInert
	// bail below, which the in-flight prompt would otherwise trip.
	if _, ok := g.liftUseTarget(); ok {
		return g.useVerbCardActions(), ""
	}
	// A peek lift, or the lift raised while choosing a house, only enlarges the
	// card to read it — there is no turn action to offer yet, so it carries no verbs.
	if g.inspecting || g.boardInert() {
		return nil, ""
	}
	// A Deploy creature is lifted while its position prompt is up: its placement
	// verbs sit on the card being placed, like the flank question but extended with
	// the interior deploy options.
	if g.choosingPosition {
		return g.deployActions()
	}
	if g.phase == phaseFlank {
		return g.flankActions()
	}
	switch g.selKind {
	case selHand:
		return g.handCardActions()
	case selYourCreature:
		return g.creatureCardActions()
	case selYourArtifact:
		return g.artifactCardActions()
	}
	return nil, "Read-only — this is your opponent's card."
}

// flankActions is the which-end question, asked on the creature being placed: it
// is a verb of that card like Play is, so it belongs where Play was rather than
// in a dock the player has to look away to.
func (g *game) flankActions() ([]cardAction, string) {
	return []cardAction{
		{"Left flank", cx("btn-primary", "btn-flank", "btn-flank--left",
			ifCls(g.isButtonCursor(0), "btn-cursor")), g.playFlank(true)},
		{"Right flank", cx("btn-primary", "btn-flank", "btn-flank--right",
			ifCls(g.isButtonCursor(1), "btn-cursor")), g.playFlank(false)},
		{"Cancel", "btn-secondary", g.cancelTargeting},
	}, ""
}

// deployActions is the Deploy placement question, asked on the creature being
// placed like the flank question but in two steps for a creature that may enter
// anywhere in the line. First it offers the Deploy left / Deploy right pair,
// naming which side of a battleline creature the new one lands on. Once a side is
// chosen the pair gives way to a "click a creature" prompt with a Back button to
// re-pick the side, and the battleline lights its creatures as the click targets.
// The ends of the line are the flanks (left of the leftmost, right of the
// rightmost), so both flanks are reachable without their own buttons. Cancel
// appears only in manual mode, the one place a play can be backed out mid-action.
func (g *game) deployActions() ([]cardAction, string) {
	var acts []cardAction
	note := ""
	if g.positionSideChosen {
		side := "left"
		if g.positionRight {
			side = "right"
		}
		note = "Click a creature to deploy " + side + " of it."
		acts = append(acts, cardAction{"Back", "btn-secondary", g.deploySideBack})
	} else {
		acts = append(acts,
			cardAction{"Deploy left", cx("btn-primary", "btn-flank", "btn-flank--left"),
				g.chooseDeploySide(false)},
			cardAction{"Deploy right", cx("btn-primary", "btn-flank", "btn-flank--right"),
				g.chooseDeploySide(true)},
		)
	}
	if g.g.Manual() {
		cancel := g.cancelChooser
		if g.manualPlacing {
			cancel = g.cancelManualPlace
		}
		acts = append(acts, cardAction{"Cancel", "btn-secondary", cancel})
	}
	return acts, note
}

func (g *game) handCardActions() ([]cardAction, string) {
	var acts []cardAction
	var note string
	if err := g.g.CanPlay(g.active(), g.sel); err != nil {
		note = "Cannot play: " + err.Error() + "."
	} else if g.canPlayAsUpgrade() {
		// A creature that may go down as a creature or an upgrade offers the choice
		// as two explicit buttons, so it is made before the play instead of by a
		// later sidebar prompt. Play upgrade skips the flank step (upgrades take no
		// flank); Play creature keeps the flank question.
		acts = append(acts,
			cardAction{"Play creature", "btn-primary", g.playAsCreature},
			cardAction{"Play upgrade", "btn-primary", g.playAsUpgrade},
		)
	} else {
		acts = append(acts, cardAction{"Play", "btn-primary", g.play})
	}
	// Discarding is offered whenever the engine allows it (active house, and not
	// barred by the first-turn one-card rule).
	if g.discardableFromHand(g.sel) {
		acts = append(acts, cardAction{"Discard", "btn-danger", g.discard})
	}
	// Manual mode adds a "Put into play" that stages the card straight onto the
	// board — no play effects, no bonus Æmber — deploying a creature anywhere in
	// the line.
	if g.g.Manual() {
		acts = append(acts, cardAction{"Put into play", "btn-secondary", g.manualPlay})
	}
	return acts, note
}

func (g *game) creatureCardActions() ([]cardAction, string) {
	// A fight grant (Brothers in Battle) makes a creature usable to fight even when
	// CanUse rejects its house, so bail with the error only when no fight is open
	// either; the per-use gates below then offer Fight alone.
	if err := g.g.CanUse(g.active(), g.sel); err != nil &&
		g.g.CanUseTo(g.active(), g.sel, engine.FightUse) != nil {
		return nil, "Cannot act: " + err.Error() + "."
	}
	// A stunned creature recovers from stun instead of acting, so any use just
	// removes the stun: offer a single Unstun rather than Reap/Fight/Action.
	if g.g.Stunned(g.sel) {
		return []cardAction{{"Unstun", "btn-unstun", g.unstun}}, "Stunned"
	}
	// Each way of using a creature is offered only when the card allows it —
	// Tireless Crocag fights and uses its Action: ability but cannot reap.
	var acts []cardAction
	if g.g.CanUseTo(g.active(), g.sel, engine.ReapUse) == nil {
		s := useVerbSpec(engine.ReapUse)
		acts = append(acts, cardAction{s.text, s.class, g.reap})
	}
	// Fight also needs a legal target (e.g. with no enemy creatures, a ready Valdr
	// can still reap but has nothing to fight).
	if g.g.CanUseTo(g.active(), g.sel, engine.FightUse) == nil &&
		len(g.g.FightTargets(g.active(), g.sel)) > 0 {
		s := useVerbSpec(engine.FightUse)
		acts = append(acts, cardAction{s.text, s.class, g.startFight})
	}
	if g.g.HasTrigger(g.sel, engine.TriggerAction) &&
		g.g.CanUseTo(g.active(), g.sel, engine.ActionUse) == nil {
		s := useVerbSpec(engine.ActionUse)
		acts = append(acts, cardAction{s.text, s.class, g.useAction})
	}
	return acts, ""
}

// useVerbButtonSpec is how one way of using a creature reads and is styled. The
// three specs are the single source for the reap/fight/action buttons, so a use
// prompted by another card (Inspiration's UseVerb) reads exactly like one the
// player reached by selecting the creature — rather than the plain lowercase
// option list a generic prompt would draw.
type useVerbButtonSpec struct {
	text  string
	class string
}

// useVerbSpec returns the button spec for a use kind.
func useVerbSpec(k engine.UseKind) useVerbButtonSpec {
	switch k {
	case engine.FightUse:
		return useVerbButtonSpec{"Fight", "btn-danger"}
	case engine.ActionUse:
		return useVerbButtonSpec{"Action", "btn-primary"}
	default: // ReapUse
		return useVerbButtonSpec{"Reap", "btn-warning"}
	}
}

// useVerbKindOfLabel maps a UseVerb option label (see the engine's effect_creature
// UseVerb.Apply) to its use kind, reporting !ok for a label that is not one of the
// three use verbs. It is how the option prompt recognises a reap/fight/action
// choice and draws it as the standard buttons.
func useVerbKindOfLabel(label string) (engine.UseKind, bool) {
	switch label {
	case "reap":
		return engine.ReapUse, true
	case "fight":
		return engine.FightUse, true
	case "use its action":
		return engine.ActionUse, true
	}
	return 0, false
}

// useVerbBarred reports whether the use-verb option should be omitted because the
// chosen creature cannot be used that way — Narp bars its neighbors from reaping,
// so Universal Translator must not offer Reap on them. It applies only when the
// target creature is known (a use another card chose, set by chooseCandidate);
// with no known target every verb stands, since there is nothing to check against.
func (g *game) useVerbBarred(label string) bool {
	if !g.hasUseTarget {
		return false
	}
	k, ok := useVerbKindOfLabel(label)
	if !ok {
		return false
	}
	return g.g.CannotBeUsedTo(g.useTarget, k)
}

func (g *game) artifactCardActions() ([]cardAction, string) {
	// An out-of-house artifact offers no Action at all, the way an out-of-house
	// creature offers no reap or fight: CanUseArtifact carries the same house check
	// creatures use, so the button is withheld rather than offered and then rejected.
	if err := g.g.CanUseArtifact(g.active(), g.sel); err != nil {
		switch {
		case !g.g.HasTrigger(g.sel, engine.TriggerAction):
			return nil, "No action ability."
		case g.g.Exhausted(g.sel):
			return nil, "Exhausted."
		}
		return nil, "Cannot act: " + err.Error() + "."
	}
	return []cardAction{{"Action", "btn-primary", g.useAction}}, ""
}
