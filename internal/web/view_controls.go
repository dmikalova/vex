package web

import (
	"strings"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/engine"
)

// This file draws the sidebar: everything the player acts through rather than
// looks at — the prompts (card picker, house picker, option chooser), the action
// bar for the selected card, and the manual-mode panel.

// promptSourceHeader names the card driving the current prompt and shows its
// face, so the player can read the ability they are resolving without hunting for
// the card on the board. The face already carries the card's name, so the name is
// spelled out only when there is no face to show.
func (g *game) promptSourceHeader() app.UI {
	return app.If(g.promptSource != "", func() app.UI {
		def := g.defByName[g.promptSource]
		if def == nil {
			return app.Div().Class("prompt-source-block").Body(
				app.Div().Class("prompt-source").Text(g.promptSource),
			)
		}
		return app.Div().Class("prompt-source-block").Body(
			app.Div().Class("prompt-card").Body(&cardView{
				Title:    def.Name,
				HouseCls: houseClasses(def.House),
				Emblem:   houseIconName(def.House),
				TypeIcon: typeIconName(def.Type),
				Stat:     handStat(def),
				Rules:    engine.RenderCardRules(def),
				Kind:     kindLabel(def),
				Trait:    traitLabel(def),
				Rarity:   rarityMarkOf(def.Rarity),
			}),
		)
	})
}

// endTurnButton is the End turn control. Once a confirm is armed (the player could
// still act this turn) it becomes a red "Confirm end turn" that ends on the next
// click, mirroring pressing E a second time.
func (g *game) endTurnButton() app.UI {
	cursor := ifCls(g.isEndTurnCursor(), "btn-cursor")
	if g.confirmEndTurn {
		return btn("Confirm end turn", g.endTurn, cx("btn-danger", cursor))
	}
	return btn("End turn", g.endTurn, cx("btn-secondary", cursor))
}

// controls is the bottom of the sidebar: the contextual controls (house picker or
// action bar) plus End turn. House selection has no End turn — a house must be
// chosen first.
func (g *game) controls() app.UI {
	// A finished game replaces every control with the result: nothing else can be
	// done, and a modal over the board would hide the position that ended it.
	if g.phase == phaseOver {
		return app.Div().Class("controls").Body(g.overPanel())
	}
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
	// A pending manual key forge takes over the controls until a colour is picked
	// or cancelled; once forgingKey resets, the previous buttons return on their own.
	if g.forgingKey >= 0 {
		return app.Div().Class("controls").Body(g.keyForgePanel())
	}
	// A labeled option prompt (e.g. "take archives?") shows its choices as buttons.
	if g.choosingOption {
		return app.Div().Class("controls").Body(g.promptSourceHeader(), g.optionChooser())
	}
	// While an engine chooser waits, the controls become the prompt itself: a
	// green call to action to click one of the highlighted cards.
	if g.choosing {
		body := []app.UI{
			g.promptSourceHeader(),
			app.Div().Class("prompt").Text(g.chooserPrompt),
		}
		// An optional prompt ("you may", "up to N") is passed on with Done; a
		// mandatory one has no cancel, except in manual mode where the player may
		// need to escape a prompt with no clickable candidate.
		switch {
		case g.chooserDeclinable:
			body = append(body, btn("Done", g.declineChooser,
				cx("btn-secondary", ifCls(g.isDoneCursor(), "btn-cursor"))))
		case g.g.Manual():
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
	// In manual mode the manual controls sit above the normal action bar, which
	// still works (house restrictions lifted) so cards can be used out of house.
	body := []app.UI{g.actionBar()}
	if g.g.Manual() {
		body = append([]app.UI{g.manualPanel()}, body...)
	}
	// A selected card with something to do owns the controls; End turn comes back
	// once nothing is selected, so it cannot be hit while mid-decision.
	if !g.selHasActions() {
		body = append(body, g.endTurnButton())
	}
	return app.Div().Class("controls").Body(body...)
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

// pickerInputID marks the card picker's search box, so opening the picker can put
// the caret straight in it.
const pickerInputID = "pickerinput"

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
			app.Input().ID(pickerInputID).Class("picker-input").Type("text").
				Placeholder("Search cards…").
				Value(g.pickerQuery).OnInput(g.pickerInput),
			app.Div().Class("picker-list").Body(
				app.Range(matches).Slice(func(i int) app.UI {
					d := matches[i]
					return app.Button().Class("picker-item").
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

// pickableHouses is the set of houses the active player may choose from. When a
// card has forced a house on them, only that house is offered — showing all
// three and then rejecting two of them makes the restriction look like a bug.
func (g *game) pickableHouses() []engine.House {
	p := g.active()
	houses := g.deckHouses[p]
	if forced := g.g.State.ForcedHouse[p].Value; forced != engine.HouseNone && !g.g.Manual() {
		if containsHouse(houses, forced) {
			return []engine.House{forced}
		}
	}
	return houses
}

// housePicker offers the houses the active player may choose from.
func (g *game) housePicker() app.UI {
	houses := g.pickableHouses()
	return app.Div().Class("btn-col").Body(
		app.Div().Class("section-title").Text("Choose a house"),
		app.Range(houses).Slice(func(i int) app.UI {
			h := houses[i]
			return app.Button().
				Class(cx("house-btn", houseAccent(h), ifCls(g.isButtonCursor(i), "btn-cursor"))).
				OnClick(g.pickHouse(h)).
				Body(houseIcon(h, "icon-inline"), app.Text(h.String()))
		}),
	)
}

func containsHouse(houses []engine.House, h engine.House) bool {
	for _, x := range houses {
		if x == h {
			return true
		}
	}
	return false
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
				return keyChoiceButton(
					c,
					g.optionLabels[i],
					g.isButtonCursor(i),
					g.chooseOptionIdx(i),
				)
			}),
		)
	}
	if g.houseOptions() {
		return app.Div().Class("btn-col").Body(
			app.Div().Class("section-title").Text("Choose a house:"),
			app.Range(g.optionLabels).Slice(func(i int) app.UI {
				h, _ := engine.ParseHouse(g.optionLabels[i])
				return app.Button().
					Class(cx("house-btn", houseAccent(h), ifCls(g.isButtonCursor(i), "btn-cursor"))).
					OnClick(g.chooseOptionIdx(i)).
					Body(houseIcon(h, "icon-inline"), app.Text(g.optionLabels[i]))
			}),
		)
	}
	return app.Div().Class("btn-col").Body(
		app.Div().Class("prompt").Text(g.optionPrompt),
		app.Range(g.optionLabels).Slice(func(i int) app.UI {
			return btn(g.optionLabels[i], g.chooseOptionIdx(i),
				cx("btn-primary", ifCls(g.isButtonCursor(i), "btn-cursor")))
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
			btn("Left flank", g.playFlank(true), cx("btn-primary", "btn-flank", "btn-flank--left",
				ifCls(g.isButtonCursor(0), "btn-cursor"))),
			btn(
				"Right flank",
				g.playFlank(false),
				cx("btn-primary", "btn-flank", "btn-flank--right",
					ifCls(g.isButtonCursor(1), "btn-cursor")),
			),
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

// selCardFace shows the selected card above its actions, the way a prompt shows
// the card that asked it. The face carries the name, so nothing else has to.
func (g *game) selCardFace() app.UI {
	def := g.g.Def(g.sel)
	house := g.g.House(g.sel)
	return app.Div().Class("prompt-card").Body(&cardView{
		Title:        def.Name,
		HouseCls:     houseClasses(house),
		Emblem:       houseIconName(house),
		HouseChanged: house != def.House,
		TypeIcon:     typeIconName(def.Type),
		Stat:         g.statLine(g.sel),
		Rules:        g.faceRules(g.sel),
		Kind:         kindLabel(def),
		Trait:        traitLabel(def),
		Rarity:       rarityMarkOf(def.Rarity),
		Maverick:     g.isMaverick(g.sel),
		Stunned:      g.g.Stunned(g.sel),
		Exhausted:    g.g.Exhausted(g.sel),
		Bar:          g.barKeywords(g.sel),
	})
}

// selHasActions reports whether the action bar is offering the selected card
// something to do.
func (g *game) selHasActions() bool {
	if !g.hasSel {
		return false
	}
	switch g.selKind {
	case selHand:
		return g.playableFromHand(g.sel) || g.discardableFromHand(g.sel)
	case selYourCreature:
		return g.g.CanUse(g.active(), g.sel) == nil
	case selYourArtifact:
		return g.g.HasAction(g.sel) && !g.g.Exhausted(g.sel)
	}
	return false
}

func (g *game) handActions() app.UI {
	items := []app.UI{g.selCardFace()}
	if err := g.g.CanPlay(g.active(), g.sel); err != nil {
		items = append(items, app.Div().Class("hint").Text("Cannot play: "+err.Error()+"."))
	} else {
		items = append(items, btn("Play", g.play, "btn-primary"))
	}
	// Discarding needs only that the card is of the active house.
	if g.discardableFromHand(g.sel) {
		items = append(items, btn("Discard", g.discard, "btn-danger"))
	}
	return app.Div().Class("btn-col").Body(items...)
}

func (g *game) creatureActions() app.UI {
	if err := g.g.CanUse(g.active(), g.sel); err != nil {
		return app.Div().Class("btn-col").Body(
			g.selCardFace(),
			app.Div().Class("hint").Text("Cannot act: "+err.Error()+"."),
		)
	}
	// A stunned creature recovers from stun instead of acting, so any use just
	// removes the stun: offer a single Unstun rather than Reap/Fight/Action.
	if g.g.Stunned(g.sel) {
		return app.Div().Class("btn-col").Body(
			g.selCardFace(),
			app.Div().Class("hint").Text("Stunned"),
			btn("Unstun", g.reap, "btn-primary"),
		)
	}
	items := []app.UI{
		g.selCardFace(),
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
	items := []app.UI{g.selCardFace()}
	switch {
	case !g.g.HasAction(g.sel):
		items = append(items, app.Div().Class("hint").Text("No action ability."))
	case g.g.Exhausted(g.sel):
		items = append(items, app.Div().Class("hint").Text("Exhausted."))
	default:
		items = append(items, btn("Action", g.useAction, "btn-primary"))
	}
	return app.Div().Class("btn-col").Body(items...)
}
