package web

import (
	"strings"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

// This file draws the game log, including the grouping of lines into per-action
// bubbles and the linking of card names within a line to their card face.

// logPanel renders the log as a flat list of lines. The grouping into per-action
// bubbles is expressed with modifier classes on the lines themselves (tint,
// rounded first/last line, newest-flash) rather than wrapper elements, so an
// already-logged line never changes structure as the log grows.
func (g *game) logPanel() app.UI {
	lines := g.logLineViews()
	return app.Div().Class("log").Body(
		app.Div().Class("log-list").ID("gamelog").Body(
			app.Range(lines).Slice(func(i int) app.UI {
				ln := lines[i]
				cls := cx("log-line",
					ifCls(ln.player == 0, "log-line--p0"),
					ifCls(ln.player == 1, "log-line--p1"),
					ifCls(ln.first, "log-line--group-start"),
					ifCls(ln.last, "log-line--group-end"),
					ifCls(ln.newest, "log-line--new"),
				)
				return app.Div().Class(cls).Body(g.logSegments(ln.text)...)
			}),
		),
	)
}

// logLineView is one rendered log line: its text, the player whose turn produced
// it (-1 for setup), whether it opens or closes its action group, and whether it
// belongs to the newest action.
type logLineView struct {
	text                string
	player              int
	first, last, newest bool
}

// logLineViews flattens the per-action groups into lines tagged with their
// position in the group, which is what the stylesheet needs to draw the bubbles.
func (g *game) logLineViews() []logLineView {
	var out []logLineView
	for _, grp := range g.logGroupViews() {
		for i, line := range grp.lines {
			out = append(out, logLineView{
				text:   line,
				player: grp.player,
				first:  i == 0,
				last:   i == len(grp.lines)-1,
				newest: grp.newest,
			})
		}
	}
	return out
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
		// Never re-cover lines an earlier segment already emitted: clamp the start
		// forward to the previous segment's end. Guards against non-monotonic marks,
		// which would otherwise overlap and duplicate log lines on screen.
		if start < first {
			start = first
		}
		if end < start {
			end = start
		}
		if start < end {
			segs = append(segs, logSeg{start, end, m.Player})
			first = end
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
