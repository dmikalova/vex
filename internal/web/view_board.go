package web

import (
	"fmt"
	"strconv"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/engine"
)

// This file draws the board: the two battlelines, the score pills above and
// below them, and the cards in them.

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
			g.renderRow(opp, "artifacts", g.g.Artifacts(opp), selOther, true),
			g.renderRow(opp, "battleline", g.g.Battleline(opp), selOther, true),
			app.Div().Class("midline"),
			g.renderRow(p, "battleline", g.g.Battleline(p), selYourCreature, false),
			g.renderRow(p, "artifacts", g.g.Artifacts(p), selYourArtifact, false),
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
	// Whose turn it is leads, since it is the thing a player re-reads most often.
	items := []app.UI{
		app.Span().Class("hud-player").Text(g.g.PlayerName(p)),
		app.Span().Class("hud-turn").Text(fmt.Sprintf("Turn %d", g.g.State.Turn)),
		app.Span().Class("hud-step").Text(steps[g.phase]),
	}
	hud := "hud"
	if h := g.g.State.ActiveHouse; h != engine.HouseNone {
		// The emblem alone names the house, and the line takes its colour, so the
		// active house reads without spending the width its name would cost.
		hud = cx("hud", "hud--house", houseAccent(h))
		items = append(items, app.Span().Class("hud-house tip").
			DataSet("tip", h.String()).
			Body(houseIcon(h, "icon-inline")))
	}
	return app.Div().Class(hud).Body(items...)
}

func (g *game) scorePill(player int) app.UI {
	active := player == g.active()
	detail := []app.UI{
		app.Text(" • "),
		g.aemberSeg(player),
		app.Text(" / "),
		g.keyCostSeg(player),
		app.Text(" • "),
		g.keysDisplay(player),
	}
	if houses := g.deckHouses[player]; len(houses) > 0 {
		detail = append(detail, app.Text(" • "), g.houseStrip(player, houses))
	}
	if g.g.State.Chains[player] > 0 || g.g.Manual() {
		detail = append(detail, app.Text(" • "), g.chainsSeg(player))
	}
	cls := cx("score-pill",
		"score-pill--p"+strconv.Itoa(player),
		ifCls(active, "score-pill-active"), ifCls(!active, "score-pill-idle"))
	return app.Div().Class(cls).
		Body(
			// Name and detail are one group so that a narrow bar wraps the zone counts
			// onto their own line instead of reflowing the stats one icon at a time.
			app.Span().Class("score-main").Body(
				app.Span().Class("score-name").Text(g.g.PlayerName(player)),
				app.Span().Class("score-detail").Body(detail...),
			),
			// Only the zone counts open the viewer, so misclicking a key or stepper
			// in the detail does not.
			app.Span().Class("score-zones").
				DataSet("player", strconv.Itoa(player)).
				OnClick(g.onScorePillClick).
				Body(g.zoneCounts(player)...),
		)
}

// zoneCounts renders a player's out-of-play zone sizes as icon-and-count pairs.
// The hand is included because knowing how many cards an opponent is holding is
// public information that a physical game makes obvious and a screen does not.
func (g *game) zoneCounts(player int) []app.UI {
	zones := []struct {
		name  string
		label string
		n     int
	}{
		{"zone-hand", "Hand", len(g.g.Hand(player))},
		{"zone-deck", "Deck", len(g.g.Deck(player))},
		{"zone-discard", "Discard", len(g.g.Discard(player))},
		{"zone-archives", "Archives", len(g.g.Archives(player))},
		{"zone-purge", "Purge", len(g.g.Purge(player))},
	}
	out := make([]app.UI, 0, len(zones))
	for _, z := range zones {
		// A destroyed or discarded card cannot pulse where it was — it is off the
		// board — so its destination pulses in its place.
		pulse := ""
		if z.label == "Discard" && g.discardFlash[player] {
			pulse = pulseClass(true, g.discardParity[player], "gain")
		}
		body := []app.UI{icon(z.name, "icon-stat"), app.Text(strconv.Itoa(z.n))}
		body = append(body, g.flightsInto(player, z.name)...)
		out = append(out,
			app.Span().Class(cx("zone-count", "tip", pulse)).DataSet("tip", z.label).Body(body...),
		)
	}
	return out
}

// flightsInto renders the cards that just left the board for this zone as faces
// arcing into its pill and shrinking away, so a card that leaves play is seen
// going somewhere rather than only bumping a counter.
func (g *game) flightsInto(player int, zone string) []app.UI {
	var out []app.UI
	for _, f := range g.flights {
		if f.player != player || f.zone != zone {
			continue
		}
		out = append(out, app.Div().
			Class(cx("card-flight",
				ifCls(!g.flightParity, "card-flight--a"),
				ifCls(g.flightParity, "card-flight--b"))).
			Body(g.printedCard(f.id)))
	}
	return out
}

// keyCostSeg shows the Æmber a player must spend to forge their next key: the
// cost followed by the forge icon.
func (g *game) keyCostSeg(player int) app.UI {
	return app.Span().Class("stat-seg tip").DataSet("tip", "Key cost").Body(
		app.Text(strconv.Itoa(g.g.CurrentKeyCost(player))),
		icon("forge", "icon-stat"),
	)
}

// hoverPreview renders the hovered card enlarged: a live board/hand card over the
// log, or a printed card (from a log mention) just left of the log.
func (g *game) hoverPreview() app.UI {
	var card app.UI
	switch {
	case g.hoverLive():
		def := g.g.Def(g.hoverID)
		house := g.g.House(g.hoverID)
		card = &cardView{
			Title:        def.Name,
			HouseCls:     houseClasses(house),
			Emblem:       houseIconName(house),
			HouseChanged: house != def.House,
			TypeIcon:     typeIconName(def.Type),
			Stat:         g.statLine(g.hoverID),
			Rules:        g.faceRules(g.hoverID),
			Kind:         kindLabel(def),
			Trait:        traitLabel(def),
			Rarity:       rarityMarkOf(def.Rarity),
			Maverick:     g.isMaverick(g.hoverID),
			Stunned:      g.g.Stunned(g.hoverID),
			Exhausted:    g.g.Exhausted(g.hoverID),
			Bar:          g.barKeywords(g.hoverID),
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
			Kind:     kindLabel(def),
			Trait:    traitLabel(def),
			Rarity:   rarityMarkOf(def.Rarity),
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
	// In manual mode the active player can switch their active house by clicking.
	clickable := g.g.Manual() && player == g.active()
	items := make([]app.UI, 0, len(houses))
	for _, h := range houses {
		dim := active != engine.HouseNone && h != active
		cls := cx("score-house", "tip", houseAccent(h), ifCls(dim, "score-house-dim"))
		// The icon alone identifies the house; the name is redundant here and costs
		// the width that pushes the pill's zone counts off small screens.
		if clickable {
			items = append(items, app.Button().Class(cx(cls, "score-house-btn")).
				DataSet("tip", h.String()).
				OnClick(g.manualSetHouse(h)).
				Body(houseIcon(h, "icon-inline")))
		} else {
			items = append(items, app.Span().Class(cls).DataSet("tip", h.String()).
				Body(houseIcon(h, "icon-inline")))
		}
	}
	return app.Span().Class("score-houses").Body(items...)
}

// aemberSeg shows a player's Æmber; in manual mode it flanks the count with
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
		return app.Span().Class(cx("stat-seg", "tip", gain)).DataSet("tip", "Æmber").Body(count, ic)
	}
	return app.Span().Class(cx("stat-seg", "amber-manual", "tip", gain)).
		DataSet("tip", "Æmber").
		Body(
			g.stepBtn(g.manualAmberDelta(player, -1), false),
			count, ic,
			g.stepBtn(g.manualAmberDelta(player, 1), true),
		)
}

// chainsSeg shows a player's chains; in manual mode it always shows (even at 0)
// with minus/plus steppers.
func (g *game) chainsSeg(player int) app.UI {
	count := app.Text(strconv.Itoa(g.g.State.Chains[player]))
	ic := icon("chains", "icon-stat", "icon-outline")
	if !g.g.Manual() {
		return app.Span().Class("stat-seg tip").DataSet("tip", "Chains").Body(count, ic)
	}
	return app.Span().Class("stat-seg amber-manual tip").DataSet("tip", "Chains").Body(
		g.stepBtn(g.manualChainsDelta(player, -1), false),
		count, ic,
		g.stepBtn(g.manualChainsDelta(player, 1), true),
	)
}

// stepBtn is a green plus or red minus stepper for the manual-mode Æmber/chains
// adjusters.
func (g *game) stepBtn(onClick app.EventHandler, plus bool) app.UI {
	label, cls := "−", "amber-btn amber-btn-minus"
	if plus {
		label, cls = "+", "amber-btn amber-btn-plus"
	}
	return app.Button().Class(cls).Text(label).OnClick(onClick)
}

// keyForgePanel is the manual-mode key-forge picker, shown inline in the controls
// space: pick a colour for the new key, or cancel. Once a choice is made
// forgingKey resets, so controls() falls back to the previous buttons on its own.
func (g *game) keyForgePanel() app.UI {
	player := g.forgingKey
	remaining := g.remainingKeyColors(player)
	return app.Div().Class("btn-col").Body(
		app.Div().Class("section-title").Text("Forge a key for "+g.g.PlayerName(player)),
		app.Range(remaining).Slice(func(i int) app.UI {
			c := remaining[i]
			return keyChoiceButton(c, c.String(), g.isButtonCursor(i), g.pickForgeColor(c))
		}),
		btn("Cancel", g.cancelForgeKey, "btn-secondary"),
	)
}

// keysDisplay shows a player's three key slots: a coloured key icon for each key
// forged (in forge order), and a dimmed key for each still to forge. In manual
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

// keySlot renders a key icon as a clickable forge/unforge button in manual mode,
// or a plain icon otherwise.
func (g *game) keySlot(ic app.UI, manual bool, onClick app.EventHandler) app.UI {
	if !manual {
		return ic
	}
	return app.Button().Class("key-btn").OnClick(onClick).Body(ic)
}

// renderRow draws one line of the board. opposing marks the rows across the
// midline from the active player, whose cards face the other way. The label is
// built from separate pieces so a short window can drop the zone word for its
// icon, and then the player name, rather than clipping the whole label.
func (g *game) renderRow(
	player int,
	zone string,
	ids []engine.LocalID,
	boardKind selKind,
	opposing bool,
) app.UI {
	zoneIcon := "type-artifact"
	if zone == "battleline" {
		zoneIcon = "type-creature"
	}
	// The board's rows are fixed tracks, so playing or losing a card does not
	// resize the row and shift the rest of the board. An opposing row hangs its
	// cards from the bottom edge, so the board is a mirror about the midline and
	// both players' cards sit the same distance from it.
	return app.Div().
		Class(cx("board-row", ifCls(opposing, "board-row--opposing"))).
		Body(
			app.Div().Class("row-label").Body(
				app.Span().Class("row-label-name").Text(g.g.PlayerName(player)+" "),
				app.Span().Class("row-label-zone").Text(zone+" "),
				app.Text(fmt.Sprintf("(%d", len(ids))),
				icon(zoneIcon, "row-label-icon"),
				app.Text(")"),
			),
			app.Div().Class("card-strip").Body(
				app.Range(ids).Slice(func(i int) app.UI {
					return g.renderCard(ids[i], boardKind, opposing)
				}),
			),
		)
}

func (g *game) renderCard(id engine.LocalID, boardKind selKind, opposing bool) app.UI {
	def := g.g.Def(id)
	activate, targetable, dimmed := g.cardVisual(id, boardKind)
	house := g.g.House(
		id,
	) // effective house: a control/"belongs to house" effect may override the printed one
	flash := g.flashes[id]
	return &cardView{
		ID:           id,
		DOMID:        boardCardID(id),
		Title:        def.Name,
		HouseCls:     houseClasses(house),
		Emblem:       houseIconName(house),
		HouseChanged: house != def.House,
		TypeIcon:     typeIconName(def.Type),
		Stat:         g.statLine(id),
		Rules:        g.faceRules(id),
		Kind:         kindLabel(def),
		Trait:        traitLabel(def),
		Rarity:       rarityMarkOf(def.Rarity),
		Maverick:     g.isMaverick(id),
		Stunned:      g.g.Stunned(id),
		Exhausted:    g.g.Exhausted(id),
		Bar:          g.barKeywords(id),
		Enter:        flash.enter,
		Fight:        flash.fight,
		FightDown:    opposing,
		Hit:          flash.damage || flash.fight,
		StunFlash:    flash.stun,
		ExhaustFlash: flash.exhaust,
		FlashOdd:     flash.odd,
		Selected:     g.isSelected(id),
		Targetable:   targetable,
		Dimmed:       dimmed,
		OnActivate:   activate,
		OnHover:      g.hoverCard,
		OnHoverOut:   g.hoverClear,
	}
}

// barKeywords lists the keywords a card in play shows as a coloured stripe on its
// top edge. Only the combat keywords are included — they decide whether a fight
// is legal and what it costs, so they must be readable without stopping to read
// the rules text.
func (g *game) barKeywords(id engine.LocalID) []engine.Keyword {
	var out []engine.Keyword
	for _, k := range []engine.Keyword{
		engine.Taunt,
		engine.Elusive,
		engine.Skirmish,
		engine.Poison,
	} {
		if g.g.HasKeyword(id, k) {
			out = append(out, k)
		}
	}
	return out
}

// cardVisual decides how a card (in hand or in play) responds and looks in the
// current mode. It returns the click handler (nil when the card is not
// clickable), whether the card is a highlighted action target, and whether it is
// lowlighted (dimmed) as an invalid choice. During a chooser or fight-target
// prompt only the eligible cards are highlighted and the rest dimmed; whenever no
// card can be acted on at all (choosing a house, picking a key colour, game over)
// the whole board dims; and in ordinary play, cards the active player cannot act
// with (wrong house, exhausted, or unplayable from hand) and the opponent's
// read-only cards are dimmed so the usable ones stand out.
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
		if containsID(g.g.FightTargets(g.active(), g.attacker), id) {
			return g.fightTargetID, true, false
		}
		return nil, false, true
	case g.phase == phaseHouse, g.phase == phaseOver, g.forgingKey >= 0:
		// Nothing can be acted on until the board is back in the main phase, so it
		// reads as inert rather than inviting a click that would be rejected. Cards
		// stay clickable to inspect.
		return g.selectBoardID, false, true
	case kind == selHand:
		return g.selectHandID, false, !g.usableFromHand(id)
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
