package web

import (
	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

// This file is the client's outermost frame: the whole-page layout, the brand
// bar above it, and the status banner. The regions it composes are drawn by its
// view_*.go siblings.

// Render draws the whole client. It runs on both the server (prerender) and the
// client; before OnMount seeds the match on the client, g is nil, so a lightweight
// placeholder is shown.
func (g *game) Render() app.UI {
	if g.g == nil {
		return app.Div().Class("")
	}

	return app.Div().Class(cx("app", ifCls(g.sidebarCollapsed, "app--sidebar-collapsed"))).Body(
		app.Raw(iconOutlineFilter),
		app.Div().Class("board-area").OnClick(g.clickAway).Body(g.boardArea()...),
		// The lift and the preview draw themselves empty rather than sitting behind an
		// app.If, so they hold a fixed place among the root's children: a conditional
		// sibling coming and going would shift them, and go-app would rebuild the
		// element underneath them — restarting the lift's grow every time a hover
		// started or ended.
		g.cardFocus(),
		g.hoverPreview(),
		app.If(!g.sidebarCollapsed, func() app.UI {
			return app.Div().Class("sidebar").Body(
				g.brandBar(),
				g.logPanel(),
				g.restrictionNotes(),
				g.turnHud(),
				g.controlDock(),
			)
		}),
		// Hiding the sidebar costs the player the log, never the game: the control
		// dock leaves with it and floats over the board instead.
		app.If(g.sidebarCollapsed, func() app.UI { return g.controlDock() }),
		app.If(g.sidebarCollapsed, func() app.UI {
			return app.Button().Class("btn-nav btn-icon sidebar-reveal").Title("Show sidebar").
				Text("«").OnClick(g.toggleSidebar)
		}),
		app.If(g.sidebarCollapsed && len(g.toastBubbles) > 0, func() app.UI {
			return g.logToast()
		}),
		app.If(g.zonesPlayer >= 0, func() app.UI { return g.zonesOverlay() }),
		app.If(g.pickerOpen, func() app.UI { return g.cardPicker() }),
		app.If(g.keysOpen, func() app.UI { return g.keysOverlay() }),
	)
}

// controlDock is everything the player answers with — the prompt, the action bar,
// the house picker, the flank buttons, end turn — plus the status banner that
// reports a rejected one. It is the same subtree in both places: a block at the
// foot of the sidebar, or a floating panel over the board once the sidebar is
// hidden, which is what keeps the game playable with the sidebar away.
func (g *game) controlDock() app.UI {
	return app.Div().Class(cx("control-dock",
		ifCls(g.sidebarCollapsed, "control-dock--floating"))).Body(
		app.If(g.status != "", func() app.UI { return g.statusBanner() }),
		g.controls(),
	)
}

// restrictionNotes names the cards currently restricting the active player, right
// above the turn HUD. A rule like Control the Weak's forced house otherwise only
// shows up as a rejected click; naming the card (hoverable, like a log mention)
// lets the player read the restriction off the card itself, and sitting above the
// HUD it is read before the step it constrains rather than after.
func (g *game) restrictionNotes() app.UI {
	sources := g.g.RestrictionSources(g.active())
	if len(sources) == 0 {
		return app.Div()
	}
	return app.Div().Class("restrictions").Body(
		app.Range(sources).Slice(func(i int) app.UI {
			name := g.g.Def(sources[i]).Name
			return app.Span().Class("restriction log-card").
				DataSet("card", name).
				OnMouseEnter(g.onLogCardHover).
				OnMouseLeave(g.onCardHoverOut).
				Text(name)
		}),
	)
}

// brandBar is the slim top of the sidebar: the title, a busy badge, the menu the
// game's own controls live behind, and the sidebar toggle. Manual mode is the one
// control that also sits outside the menu, but only while it is on: a mode that
// rewrites the rules should be visibly on and one click from off.
func (g *game) brandBar() app.UI {
	return app.Div().Class("brandbar").Body(
		app.Span().Class("brand-title").Text("Vactrol"),
		// The server publishes the short build id of the bundle it served.
		app.Span().Class("brand-version").Text(app.Getenv("VACTROL_BUILD")),
		app.If(g.busy && !g.choosing && !g.choosingOption, func() app.UI {
			return app.Span().Class("badge-busy").Text("resolving…")
		}),
		app.Div().Class("spacer"),
		app.If(g.g.Manual(), func() app.UI {
			return app.Button().Class("btn-nav btn-icon btn-nav-on").
				Title("Manual mode is on — click to turn it off").
				Disabled(g.busy || g.choosing || g.choosingOption).
				OnClick(g.toggleManual).
				Body(icon("wrench", "icon-nav"))
		}),
		g.brandMenu(),
		app.Button().Class("btn-nav btn-icon").Title("Hide sidebar").
			Text("»").OnClick(g.toggleSidebar),
	)
}

// brandMenu holds the controls that frame a game rather than play it — undo,
// redo, manual mode, a new game, the keyboard sheet — behind one hamburger, so
// the top of the sidebar is a title and not a row of icons competing with the
// board. A transparent backdrop under the panel closes it on the next click
// anywhere else.
func (g *game) brandMenu() app.UI {
	return app.Div().Class("menu").Body(
		app.Button().Class(cx("btn-nav", "btn-icon", ifCls(g.menuOpen, "btn-nav-on"))).
			Title("Menu").Text("☰").OnClick(g.toggleMenu),
		app.If(g.menuOpen, func() app.UI {
			return app.Div().Class("menu-backdrop").OnClick(g.closeMenu).Body(
				app.Div().Class("menu-panel").OnClick(g.stopClick).Body(
					menuItem("undo", "Undo", g.undoMenu, !g.canUndo(), false),
					menuItem("redo", "Redo", g.redoMenu, !g.canRedo(), false),
					menuItem("wrench", "Manual mode", g.manualMenu,
						g.busy || g.choosing || g.choosingOption, g.g.Manual()),
					menuItem("restart", "New game", g.restartMenu,
						g.busy || g.choosing || g.choosingOption, false),
					menuItem("", "Keyboard shortcuts", g.keysMenu, false, false),
					app.Hr().Class("menu-divider"),
					menuLink("Rulebook", "/rulebook"),
					menuLink("Glossary", "/glossary"),
				),
			)
		}),
	)
}

// menuLink is a menu row that navigates to a reference page rather than acting on
// the game — the rulebook and glossary.
func menuLink(label, href string) app.UI {
	return app.A().Class("menu-item").Href(href).Body(
		app.Span().Class("menu-glyph").Text("›"),
		app.Span().Class("menu-label").Text(label),
	)
}

// menuItem is one line of the menu: an icon (or the ? glyph for the shortcut
// sheet), its name, and the on-state for a mode that is currently engaged.
func menuItem(iconName, label string, onClick app.EventHandler, disabled, on bool) app.UI {
	glyph := app.UI(app.Span().Class("menu-glyph").Text("?"))
	if iconName != "" {
		glyph = icon(iconName, "icon-nav")
	}
	return app.Button().
		Class(cx("menu-item", ifCls(on, "menu-item-on"))).
		Disabled(disabled).
		OnClick(onClick).
		Body(glyph, app.Span().Class("menu-label").Text(label))
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
