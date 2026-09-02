package web

import (
	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/engine"
)

// This file holds the client's own small pieces of state, which no engine rule
// touches: the hover preview, the restart confirmation, and the sidebar toggle.

// hoverCard previews a live board or hand card over the log.
func (g *game) hoverCard(_ app.Context, id engine.LocalID) {
	g.hoverID, g.hasHover, g.hoverDef, g.hoverInLog = id, true, nil, false
}

// hoverClear hides the hover preview (a card leave).
func (g *game) hoverClear(_ app.Context) { g.hasHover, g.hoverDef = false, nil }

// onLogCardHover previews the printed card named by a log mention, to the left of
// the log.
func (g *game) onLogCardHover(ctx app.Context, _ app.Event) {
	if def, ok := g.defByName[ctx.JSSrc().Get("dataset").Get("card").String()]; ok {
		g.hasHover, g.hoverDef, g.hoverInLog = false, def, true
	}
}

// onCardHoverOut hides the hover preview (a log-mention leave).
func (g *game) onCardHoverOut(_ app.Context, _ app.Event) { g.hasHover, g.hoverDef = false, nil }

// remainingKeyColors lists the key colours player has not yet forged, in the
// canonical order.
func (g *game) remainingKeyColors(player int) []engine.KeyColor {
	used := map[engine.KeyColor]bool{}
	for _, c := range g.g.KeyColors(player) {
		used[c] = true
	}
	var out []engine.KeyColor
	for _, c := range []engine.KeyColor{engine.KeyColorRed, engine.KeyColorBlue, engine.KeyColorYellow} {
		if !used[c] {
			out = append(out, c)
		}
	}
	return out
}

func (g *game) restart(ctx app.Context, _ app.Event) {
	if g.busy || g.choosing {
		return
	}
	g.confirmRestart = false
	g.newMatch()
	g.save(ctx)
}

// askRestart opens the restart confirmation in the controls area rather than
// restarting immediately.
func (g *game) askRestart(_ app.Context, _ app.Event) {
	if g.busy || g.choosing || g.choosingOption {
		return
	}
	g.confirmRestart = true
}

// cancelRestart dismisses the restart confirmation, leaving the game running.
func (g *game) cancelRestart(_ app.Context, _ app.Event) {
	g.confirmRestart = false
}

// toggleSidebar hides or shows the whole sidebar so the board area can use the
// full width.
func (g *game) toggleSidebar(_ app.Context, _ app.Event) {
	g.sidebarCollapsed = !g.sidebarCollapsed
}

// toggleMenu opens or closes the sidebar's hamburger menu.
func (g *game) toggleMenu(_ app.Context, _ app.Event) {
	g.menuOpen = !g.menuOpen
}

// closeMenu shuts the hamburger menu — what a click outside it, or picking one of
// its items, means.
func (g *game) closeMenu(_ app.Context, _ app.Event) {
	g.menuOpen = false
}

// The menu's items close the menu and then do their thing, so the panel does not
// hang open over the change it just made.

func (g *game) undoMenu(ctx app.Context, e app.Event) {
	g.menuOpen = false
	g.undoAction(ctx, e)
}

func (g *game) redoMenu(ctx app.Context, e app.Event) {
	g.menuOpen = false
	g.redoAction(ctx, e)
}

func (g *game) manualMenu(ctx app.Context, e app.Event) {
	g.menuOpen = false
	g.toggleManual(ctx, e)
}

func (g *game) restartMenu(ctx app.Context, e app.Event) {
	g.menuOpen = false
	g.askRestart(ctx, e)
}

func (g *game) keysMenu(ctx app.Context, e app.Event) {
	g.menuOpen = false
	g.toggleKeys(ctx, e)
}

// toggleKeys opens or closes the keyboard shortcut sheet.
func (g *game) toggleKeys(_ app.Context, _ app.Event) {
	g.keysOpen = !g.keysOpen
}

// closeKeys dismisses the keyboard shortcut sheet.
func (g *game) closeKeys(_ app.Context, _ app.Event) {
	g.keysOpen = false
}
