package web

import (
	"strconv"
	"time"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vex/internal/engine"
)

// badgeFadeDur is how long the selection badges linger on the board, growing and
// fading, after an effect's choose loop ends. It matches sel-badge--clearing in
// app.css.
const badgeFadeDur = 450 * time.Millisecond

// PreviewBadge implements the engine's BadgeChooser: an effect's choose loop
// previews the status each pick lands, so the board can badge each chosen
// creature and the cursor marker can show what a click does. A non-zero badge
// opens the preview (resetting the running totals); the zero badge closes it,
// leaving the badges on screen to grow and fade.
func (c *webChooser) PreviewBadge(badge engine.SelectionBadge) {
	c.g.dispatch(func(app.Context) {
		if c.stale() {
			return
		}
		if badge.Icon == engine.NoStatusIcon {
			c.g.endBadgePreview()
			return
		}
		c.g.selBadge = badge
		c.g.badgeTotals = map[engine.LocalID]int{}
		c.g.badgeClearing = false
	})
}

// recordBadge accumulates the badge amount a pick lands on a creature, so a
// creature chosen twice sums its badges. A zero-amount badge (a ward) still
// records the key, so the creature draws a numberless icon.
func (g *game) recordBadge(id engine.LocalID) {
	if g.selBadge.Icon == engine.NoStatusIcon {
		return
	}
	if g.badgeTotals == nil {
		g.badgeTotals = map[engine.LocalID]int{}
	}
	g.badgeTotals[id] += g.selBadge.Amount
}

// endBadgePreview leaves the badges the selection landed on screen and grows and
// fades them away, then clears the preview once the animation has run. A new
// preview beginning mid-fade advances badgeGen, retiring this timer.
func (g *game) endBadgePreview() {
	if len(g.badgeTotals) == 0 {
		g.selBadge = engine.SelectionBadge{}
		return
	}
	g.badgeClearing = true
	g.badgeGen++
	gen := g.badgeGen
	time.AfterFunc(badgeFadeDur, func() {
		g.dispatch(func(app.Context) {
			if g.badgeGen != gen {
				return
			}
			g.selBadge = engine.SelectionBadge{}
			g.badgeTotals = nil
			g.badgeClearing = false
		})
	})
}

// badgeIconName maps a status icon to its asset stem, or "" for one with no
// badge to draw.
func badgeIconName(s engine.StatusIcon) string {
	switch s {
	case engine.DamageIcon:
		return "damage"
	case engine.WardIcon:
		return "ward"
	}
	return ""
}

// cardBadge is the large status icon a board card draws while it is a picked
// candidate of a badge preview — the running damage or a numberless ward — or nil
// when this card carries no badge.
func (g *game) cardBadge(id engine.LocalID) app.UI {
	amount, ok := g.badgeTotals[id]
	if !ok {
		return nil
	}
	name := badgeIconName(g.selBadge.Icon)
	if name == "" {
		return nil
	}
	var body []app.UI
	if amount > 0 {
		body = append(body, app.Span().Class("sel-badge-num").Text(strconv.Itoa(amount)))
	}
	body = append(body, icon(name, "sel-badge-icon"))
	return app.Div().
		Class(cx("sel-badge", ifCls(g.badgeClearing, "sel-badge--clearing"))).
		Body(body...)
}

// selCursor is the persistent marker the pointer carries during a badge preview:
// a status icon centred under the cursor showing what a click will do. It is
// a fixed sibling (like tip-float) whose position the pointermove listener sets;
// it draws empty and hidden when no preview is active.
func (g *game) selCursor() app.UI {
	var body []app.UI
	name := badgeIconName(g.selBadge.Icon)
	on := name != "" && !g.badgeClearing
	if on {
		body = append(body, icon(name, "sel-cursor-icon"))
	}
	return app.Div().ID("sel-cursor").
		Class(cx("sel-cursor", ifCls(on, "sel-cursor--on"))).
		Body(body...)
}

// installSelCursor follows the pointer with the selection-badge marker while a
// preview is active, positioning the fixed #sel-cursor element at the cursor. It
// shares the tip listeners' passive registration and is released alongside them.
func (g *game) installSelCursor(doc app.Value, passive map[string]any) {
	g.selCursorFunc = app.FuncOf(func(_ app.Value, args []app.Value) any {
		if len(args) == 0 {
			return nil
		}
		el := doc.Call("getElementById", "sel-cursor")
		if !el.Truthy() {
			return nil
		}
		style := el.Get("style")
		style.Set("left", strconv.Itoa(int(args[0].Get("clientX").Float()))+"px")
		style.Set("top", strconv.Itoa(int(args[0].Get("clientY").Float()))+"px")
		// Over the player bars and sidebar the real cursor is left in place, so the
		// marker hides there; over a selectable card it wears a ring.
		target := args[0].Get("target")
		overBar := target.Truthy() &&
			(target.Call("closest", ".score-pill").Truthy() ||
				target.Call("closest", ".sidebar").Truthy())
		overCard := target.Truthy() &&
			target.Call("closest", ".card--targetable, .card--clickable, .card-tab--target").
				Truthy()
		cl := el.Get("classList")
		cl.Call("toggle", "sel-cursor--hidden", overBar)
		cl.Call("toggle", "sel-cursor--bordered", overCard)
		return nil
	})
	doc.Call("addEventListener", "pointermove", g.selCursorFunc, passive)
}
