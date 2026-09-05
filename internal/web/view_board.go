package web

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

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
		phaseMain:        "Main phase",
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
		ids   []engine.LocalID
	}{
		{"zone-hand", "Hand", g.g.Hand(player)},
		{"zone-deck", "Deck", g.g.Deck(player)},
		{"zone-discard", "Discard", g.g.Discard(player)},
		{"zone-archives", "Archives", g.g.Archives(player)},
		{"zone-purge", "Purge", g.g.Purge(player)},
	}
	out := make([]app.UI, 0, len(zones))
	for _, z := range zones {
		// A destroyed or discarded card cannot pulse where it was — it is off the
		// board — so its destination pulses in its place.
		pulse := ""
		if z.label == "Discard" && g.discardFlash[player] {
			pulse = pulseClass(true, g.discardParity[player], "gain")
		}
		body := []app.UI{icon(z.name, "icon-stat"), app.Text(strconv.Itoa(len(z.ids)))}
		body = append(body, g.flightsInto(player, z.name)...)
		// The tip names the cards in the zone when this player is allowed to see
		// them — a face-up pile, or their own hand — and otherwise just labels it.
		cls := cx("zone-count", "tip", pulse)
		tip := z.label
		if names := g.zoneNames(player, z.label, z.ids); len(names) > 0 {
			cls = cx(cls, "tip-multi")
			tip = z.label + "\n" + strings.Join(names, "\n")
		}
		out = append(out,
			app.Span().Class(cls).DataSet("tip", tip).Body(body...),
		)
	}
	return out
}

// zoneNames lists the names of the cards in a zone, sorted, but only for the
// zones this player may read: the face-up discard and purge piles of either
// player, and their own hand. A hidden zone (a deck's order, face-down archives,
// an opponent's hand) returns nothing, so hovering it never leaks its contents.
func (g *game) zoneNames(player int, label string, ids []engine.LocalID) []string {
	switch label {
	case "Discard", "Purge":
	case "Hand":
		if player != g.active() {
			return nil
		}
	default:
		return nil
	}
	names := make([]string, 0, len(ids))
	for _, id := range ids {
		names = append(names, g.g.Def(id).Name)
	}
	sort.Strings(names)
	return names
}

// flightsInto renders the cards that just left the board for this zone as faces
// arcing into its pill and shrinking away, so a card that leaves play is seen
// going somewhere rather than only bumping a counter. The opposing player's bar
// sits above their battleline rather than below it, so their arc is mirrored
// (card-flight--opposing) to fly up into the pill instead of down.
func (g *game) flightsInto(player int, zone string) []app.UI {
	opposing := player != g.active()
	var out []app.UI
	for _, f := range g.flights {
		if f.player != player || f.zone != zone {
			continue
		}
		out = append(out, app.Div().
			Class(cx("card-flight",
				ifCls(!g.flightParity, "card-flight--a"),
				ifCls(g.flightParity, "card-flight--b"),
				ifCls(opposing, "card-flight--opposing"))).
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
	case !g.previewUp():
		return app.Div()
	case g.hoverLive():
		card = g.cardFace(g.hoverID)
	case g.hoverDef != nil:
		def := g.hoverDef
		card = &cardView{
			Title:    def.Name,
			HouseCls: houseClasses(def.House),
			Emblem:   houseIconName(def.House),
			TypeIcon: typeIconName(def.Type),
			Stat:     handStat(def),
			Rules:    displayRules(engine.RenderCardRules(def)),
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
// minus/plus buttons that adjust it (usable on both players from one seat). A
// player at check — holding enough to afford a key — gets a soft glow, the
// client's stand-in for the tabletop's "Check!" callout.
func (g *game) aemberSeg(player int) app.UI {
	count := app.Text(strconv.Itoa(g.g.Aember(player)))
	ic := icon("aember", "icon-stat")
	// A pool gain pulses the segment; the -a/-b pair alternates so it replays.
	gain := cx(
		ifCls(g.poolFlash[player] && !g.poolParity[player], "stat-seg--gain-a"),
		ifCls(g.poolFlash[player] && g.poolParity[player], "stat-seg--gain-b"),
		ifCls(g.g.Aember(player) >= g.g.CurrentKeyCost(player), "stat-seg--check"),
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

// keysTally draws a static row of a player's three key slots for the game log:
// a coloured icon for each colour in colors (forge order) and a dimmed key for
// each still to forge — the same colouring as keysDisplay, without its
// manual-mode forge/unforge buttons, so a past turn's colours never change.
func keysTally(colors []engine.KeyColor) app.UI {
	icons := make([]app.UI, 0, engine.KeysToWin)
	for _, c := range colors {
		name := keyColorIconName(c)
		if name == "" {
			name = "key"
		}
		icons = append(icons, icon(name, "icon-inline"))
	}
	for i := len(colors); i < engine.KeysToWin; i++ {
		icons = append(icons, icon("key", "icon-inline", "key-unforged"))
	}
	return app.Span().Class("score-keys").Body(icons...)
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
	face := &cardView{
		ID:            id,
		DOMID:         boardCardID(id),
		Title:         def.Name,
		HouseCls:      houseClasses(house),
		Emblem:        houseIconName(house),
		HouseChanged:  house != def.House,
		TypeIcon:      typeIconName(def.Type),
		Stat:          g.statLine(id),
		Rules:         g.faceRules(id),
		Kind:          kindLabel(def),
		Trait:         traitLabel(def),
		Rarity:        rarityMarkOf(def.Rarity),
		Maverick:      g.isMaverick(id),
		Stunned:       g.g.Stunned(id),
		Exhausted:     g.g.Exhausted(id),
		PowerCounters: int(g.g.State.Cards[id].PowerCounters),
		Bar:           g.barKeywords(id),
		BarBottom:     opposing,
		TauntShielded: def.Type == engine.Creature && g.g.TauntShielded(id),
		Enter:         flash.enter,
		Fight:         flash.fight,
		FightDown:     opposing,
		Hit:           flash.damage || flash.fight,
		Reap:          flash.reap,
		Act:           flash.act,
		StunFlash:     flash.stun,
		ExhaustFlash:  flash.exhaust,
		PowerFlash:    flash.power,
		FlashOdd:      flash.odd,
		Selected:      g.isSelected(id),
		Targetable:    targetable,
		Dimmed:        dimmed,
		Jiggle:        g.jiggling(id, boardKind),
		OnActivate:    activate,
		OnHover:       g.hoverCard,
		OnHoverOut:    g.hoverClear,
	}
	return g.hostWithTabs(id, face)
}

// hostWithTabs wraps a rendered face in the peeking-tab host when the card
// carries upgrades or under-cards, so the board and the Style gallery build the
// tabbed layout from the same code. A card with nothing attached is returned
// unwrapped.
//
// The tab strips are siblings drawn before the face in a shared, non-clipping
// host, so the face — later in the DOM, same stacking context — paints over
// their inner edge and only a sliver of each peeks out. Nesting them inside
// cardView itself would not work: .card clips its own children to draw the ogee
// name-banner frame, which would hide the peeking part too.
func (g *game) hostWithTabs(id engine.LocalID, face app.UI) app.UI {
	left, right := g.underTabs(id), g.upgradeTabs(id)
	if len(left) == 0 && len(right) == 0 {
		return face
	}
	// Attached cards dim with their host: an exhausted creature has already acted,
	// so its upgrades and under-cards read as spent alongside it rather than
	// standing out beside a greyed face.
	dim := ifCls(g.inPlay(id) && g.g.Exhausted(id), "card-tabs--dim")
	return app.Div().Class("card-host").
		Style("--under-tabs", strconv.Itoa(len(left))).
		Style("--up-tabs", strconv.Itoa(len(right))).
		Body(
			app.Div().Class(cx("card-tabs", "card-tabs--left", dim)).Body(left...),
			app.Div().Class(cx("card-tabs", "card-tabs--right", dim)).Body(right...),
			face,
		)
}

// upgradeTabs renders each upgrade attached to id as a peeking tab along its
// right edge, in attach order. An upgrade is never facedown, so every tab shows
// its own house colour and hovers into the full preview.
func (g *game) upgradeTabs(id engine.LocalID) []app.UI {
	ups := g.g.Upgrades(id)
	tabs := make([]app.UI, 0, len(ups))
	for _, up := range ups {
		tabs = append(tabs, g.cardTab(up))
	}
	return tabs
}

// underTabs renders each card placed under id as a peeking tab along its left
// edge, in the order they were placed: its own house colour when revealed
// (faceup, or facedown but the active player controls id and may Peek), a plain
// card back otherwise, which previews nothing.
func (g *game) underTabs(id engine.LocalID) []app.UI {
	buried := g.g.Under(id)
	tabs := make([]app.UI, 0, len(buried))
	for _, u := range buried {
		if g.g.UnderFaceDown(u) && !g.g.Peekable(g.active(), id) {
			tabs = append(tabs, app.Div().Class("card-tab card-tab--back").
				Body(app.Span().Class("card-tab-title").Text("VEX")))
			continue
		}
		tabs = append(tabs, g.cardTab(u))
	}
	return tabs
}

// cardTab renders one revealed peeking tab: a house-tinted sliver, its card's
// name banner turned on its side, so the whole title reads down the tab and each
// further tab fans out past the last rather than hiding below it. It previews the
// full card the same way hovering the card itself does. The id is read back off
// the element's own dataset, the same way onScorePillClick reads its player,
// since a tab is a plain element rather than a component that could carry it as a
// field.
func (g *game) cardTab(id engine.LocalID) app.UI {
	def := g.g.Def(id)
	return app.Div().
		Class(cx("card-tab", houseClasses(g.g.House(id)))).
		DataSet("id", strconv.Itoa(int(id))).
		OnMouseEnter(g.onCardTabHover).
		OnMouseLeave(g.onCardTabHoverOut).
		Body(app.Span().Class("card-tab-title").Text(def.Name))
}

// barKeywordOrder is the printed keywords the stripe shows, in the order it
// stacks them. Only the combat keywords are included — they decide whether a
// fight is legal and what it costs, so they must be readable without stopping
// to read the rules text.
var barKeywordOrder = []engine.Keyword{
	engine.Taunt,
	engine.Elusive,
	engine.Skirmish,
	engine.Poison,
}

// barKeywords lists the stripe entries a card in play currently has: its
// keywords, granted ones included, in barKeywordOrder, then Hazardous last —
// Hazardous is a magnitude rather than a boolean keyword, so it counts as
// present whenever its value (including upgrades) is greater than zero.
func (g *game) barKeywords(id engine.LocalID) []string {
	var out []string
	for _, k := range barKeywordOrder {
		// A creature that has spent its Elusive this turn is no longer elusive for
		// the rest of the turn, so its stripe drops the Elusive colour.
		if k == engine.Elusive && g.g.ElusiveSpent(id) {
			continue
		}
		if g.g.HasKeyword(id, k) {
			out = append(out, k.String())
		}
	}
	if g.g.Hazardous(id) > 0 {
		out = append(out, "Hazardous")
	}
	return out
}

// barKeywordsOf lists the stripe entries a definition prints, for a card face
// built without a board behind it.
func barKeywordsOf(def *engine.CardDefinition) []string {
	var out []string
	for _, k := range barKeywordOrder {
		if hasKeyword(def, k) {
			out = append(out, k.String())
		}
	}
	if def.Hazardous > 0 {
		out = append(out, "Hazardous")
	}
	return out
}

// boardInert reports whether nothing on the board can be acted on at all right
// now: an option prompt (yes/no, which key colour to forge) is up, the board is
// between turns (choosing a house, game over), or the manual key-forge picker is
// open. It is shared by every place a card is drawn (cardVisual, renderZoneCard)
// so the board reads as inert consistently rather than each place deciding it on
// its own — an option prompt in particular blocks the whole board exactly like
// these already did, but until now fell through to g.busy's plain "leave it as
// it was" instead, which left cards lit as if still actionable.
func (g *game) boardInert() bool {
	return g.choosingOption ||
		g.phase == phaseHouse ||
		g.phase == phaseOver ||
		g.forgingKey >= 0
}

// cardVisual decides how a card (in hand or in play) responds and looks in the
// current mode. It returns the click handler (nil when the card is not
// clickable), whether the card is a highlighted action target, and whether it is
// lowlighted (dimmed) as an invalid choice. During a chooser or fight-target
// prompt only the eligible cards are highlighted and the rest dimmed; whenever no
// card can be acted on at all (boardInert) the whole board dims; and in ordinary
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
	case g.boardInert():
		// Nothing can be acted on until the prompt in front of the player is
		// answered, so the board reads as inert rather than inviting a click that
		// would be rejected. Cards stay clickable to inspect.
		if kind == selHand {
			return g.selectHandID, false, true
		}
		return g.selectBoardID, false, true
	case g.busy:
		return nil, false, false
	case g.phase == phaseFightTarget:
		if containsID(g.g.FightTargets(g.active(), g.attacker), id) {
			return g.fightTargetID, true, false
		}
		return nil, false, true
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
// cards this turn, so the ones they cannot use are lowlighted. Both creatures and
// artifacts defer to the engine (CanUse / CanUseArtifact) rather than
// reimplementing the house check here, so a Versatile artifact (Lifeward) is
// correctly offered out of the active house.
func (g *game) actionable(id engine.LocalID, kind selKind) bool {
	switch kind {
	case selYourCreature:
		return g.g.CanUse(g.active(), id) == nil
	case selYourArtifact:
		return g.g.CanUseArtifact(g.active(), id) == nil
	default:
		return true
	}
}

// jiggling reports whether a card should play the end-turn attention wobble: the
// end-turn confirm is armed and this is one of the cards the player could still
// act with, so the confirm points at exactly what it is warning about. It mirrors
// hasMoves, which decides whether the confirm arms at all.
func (g *game) jiggling(id engine.LocalID, kind selKind) bool {
	if !g.confirmEndTurn {
		return false
	}
	switch kind {
	case selHand:
		return g.usableFromHand(id)
	case selYourCreature, selYourArtifact:
		return g.actionable(id, kind)
	}
	return false
}
