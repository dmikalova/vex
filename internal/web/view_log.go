package web

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/engine"
)

// This file draws the game log: the grouping of entries into per-phase bubbles,
// the rules that demarcate a turn and a phase, and the card and player names
// inside a line.

// logPanel renders the log as one bubble per phase. The engine narrates a turn's
// shape as typed entries (TurnBegan, PhaseBegan — ADR 0012), so the grouping is
// read off the log itself rather than tracked alongside it by the client.
func (g *game) logPanel() app.UI {
	blocks := g.logBlocks()
	return app.Div().Class("log").Body(
		app.Div().Class("log-list").ID("gamelog").Body(
			app.Range(blocks).Slice(func(i int) app.UI { return g.logBlockView(blocks[i]) }),
		),
	)
}

// logRule is how heavily a record rules a line across the log: not at all, as a
// phase header, or as the heavier turn header.
type logRule int

const (
	ruleNone logRule = iota
	rulePhase
	ruleTurn
)

// restoredRule is a persisted header read back from local storage. A typed entry
// does not survive JSON, so the saved line carries the rule it drew and this
// wrapper hands it back — a resumed match keeps its dividers.
type restoredRule struct {
	engine.RestoredEntry
	Rule   logRule
	Player int
}

// ruleOf reports whether a record opens a new block, how it rules, and whose turn
// it announces. The client asks the entry what it is rather than matching a
// prefix on its prose (ADR 0011).
func ruleOf(rec engine.Record) (logRule, int) {
	switch e := rec.Entry.(type) {
	case engine.TurnBegan:
		return ruleTurn, e.Player
	case engine.PhaseBegan:
		return rulePhase, e.Player
	case restoredRule:
		return e.Rule, e.Player
	}
	return ruleNone, -1
}

// logBlock is one drawn piece of the log: either a rule (a turn or phase header
// ruled across the panel, standing between bubbles) or a bubble of the lines one
// root action produced. header is set on a rule, lines on a bubble.
type logBlock struct {
	header engine.LogEntry
	rule   logRule
	lines  []engine.Record
	player int
	newest bool
}

// logBlocks splits the log into the rules that demarcate turns and phases and the
// bubbles between them, one per root action. The turn shape comes from the log's
// own typed entries; where one action stops and the next starts is the client's
// own knowledge, since the engine frames abilities rather than player intent.
func (g *game) logBlocks() []logBlock {
	starts := make(map[int]int, len(g.logGroups))
	for _, m := range g.logGroups {
		starts[m.Start] = m.Player
	}
	var out []logBlock
	cur := logBlock{player: -1}
	flush := func(player int) {
		if len(cur.lines) > 0 {
			out = append(out, cur)
		}
		cur = logBlock{player: player}
	}
	for i, rec := range g.g.Log {
		if rule, player := ruleOf(rec); rule != ruleNone {
			flush(player)
			out = append(out, logBlock{header: rec.Entry, rule: rule, player: player})
			continue
		}
		if player, ok := starts[i]; ok {
			flush(player)
		}
		cur.lines = append(cur.lines, rec)
	}
	flush(-1)
	for i := len(out) - 1; i >= 0; i-- {
		if out[i].header == nil {
			out[i].newest = true
			break
		}
	}
	return out
}

// logBlockView draws one block: a rule on its own, or a bubble of lines.
func (g *game) logBlockView(b logBlock) app.UI {
	if b.header != nil {
		return app.Div().
			Class(cx("log-rule",
				ifCls(b.player == 0, "log-rule--p0"),
				ifCls(b.player == 1, "log-rule--p1"),
				ifCls(b.rule == ruleTurn, "log-rule--turn"),
			)).
			Body(app.Span().Class("log-rule-label").Body(g.logSegments(b.header)...))
	}
	cls := cx("log-group",
		ifCls(b.player == 0, "log-group--p0"),
		ifCls(b.player == 1, "log-group--p1"),
		ifCls(b.newest, "log-group--new"),
	)
	body := make([]app.UI, 0, len(b.lines))
	for _, rec := range b.lines {
		body = append(body, app.Div().Class("log-line").Body(g.logSegments(rec.Entry)...))
	}
	return app.Div().Class(cls).Body(body...)
}

// logSegments draws one log entry, turning the card names, player names, and
// keywords the entry itself reported into clickable spans, tinted names, and
// emblems. The engine hands back what every marked span stands for (ADR 0011),
// so nothing here matches prose against a card index or a word list.
func (g *game) logSegments(entry engine.LogEntry) []app.UI {
	// PlayerStanding carries the actual forged colours, so the end-of-turn tally
	// draws three coloured key slots instead of a plain count.
	if ps, ok := entry.(engine.PlayerStanding); ok {
		return g.playerStandingSegments(ps)
	}
	var out []app.UI
	for _, seg := range engine.RenderEntry(entry, g.g) {
		switch {
		case seg.HasCard:
			out = append(out, app.Span().
				Class("log-card").
				DataSet("card", seg.Text).
				OnMouseEnter(g.onLogCardHover).
				OnMouseLeave(g.onCardHoverOut).
				Text(seg.Text))
		case seg.HasPlayer:
			out = append(out, app.Span().
				Class("log-player log-player--p"+strconv.Itoa(seg.Player)).
				Text(seg.Text))
		case seg.Icon != "":
			out = append(out, logIcon(seg.Icon), app.Text(seg.Text))
		default:
			out = append(out, app.Text(seg.Text))
		}
	}
	return out
}

// playerStandingSegments draws a PlayerStanding entry: the player name, the
// Æmber count with its icon, and the key count with its three slots coloured
// by KeyColors rather than a plain "N keys" number.
func (g *game) playerStandingSegments(e engine.PlayerStanding) []app.UI {
	return []app.UI{
		app.Span().
			Class("log-player log-player--p" + strconv.Itoa(e.Player)).
			Text(g.g.PlayerName(e.Player)),
		app.Text(fmt.Sprintf(" has %d ", e.Aember)),
		logIcon("aember"),
		app.Text(fmt.Sprintf(" Æmber and %d/%d ", len(e.KeyColors), engine.KeysToWin)),
		keysTally(e.KeyColors),
		app.Text(" keys"),
	}
}

// logIcon draws the emblem the engine flagged a keyword with, giving house
// emblems the outline that keeps them legible against the log's background.
func logIcon(key string) app.UI {
	if strings.HasPrefix(key, "house-") {
		return icon(key, "icon-inline", "icon-outline")
	}
	return icon(key, "icon-inline")
}
