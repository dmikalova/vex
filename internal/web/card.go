package web

import (
	"strings"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/engine"
)

// cardView is a presentational component for a single card face. It carries no
// game logic: the parent supplies already-rendered strings, visual flags, and a
// click handler, so the same component renders a hand card, a creature in play,
// an artifact, or a targeting candidate.
type cardView struct {
	app.Compo

	ID       engine.LocalID
	Title    string
	HouseCls string // house-derived border/background classes
	Emblem   string // house emblem asset stem ("" for none)
	// HouseChanged marks that the card's current house differs from its printed
	// house (a control/"belongs to house" effect); the emblem is highlighted.
	HouseChanged bool
	TypeIcon     string   // card-type icon asset stem
	Stat         []app.UI // compact stat nodes (power, damage, Æmber… with icons)
	Rules        string   // rules/ability text for the face
	Kind         string   // card type label shown at the foot
	Stunned      bool     // shows a stun token on the face
	// Enter pulses the whole card as it comes into play; StunFlash pulses the stun
	// token as it is first applied. FlashOdd alternates their animation class each
	// time so the CSS animation replays even on back-to-back triggers.
	Enter      bool
	StunFlash  bool
	FlashOdd   bool
	Selected   bool
	Targetable bool
	Dimmed     bool
	// OnActivate is called with ID when the card is clicked; nil means the card is
	// not clickable. The id is passed rather than captured in the handler because
	// go-app compares event handlers by function pointer and would not refresh a
	// captured id when the board re-renders.
	OnActivate func(app.Context, engine.LocalID)
	// Draggable makes the card an HTML5 drag source (a playable hand card). When
	// set, OnDragStart fires with ID as the drag begins and OnDragEnd as it ends.
	Draggable   bool
	OnDragStart func(app.Context, engine.LocalID)
	OnDragEnd   func(app.Context, engine.LocalID)
	// OnHover fires with ID when the pointer enters the card, OnHoverOut when it
	// leaves — they drive the hover card preview.
	OnHover    func(app.Context, engine.LocalID)
	OnHoverOut func(app.Context)
}

// onClick is a stable method (unlike a per-card closure) so go-app keeps it bound
// across re-renders; it reads the up-to-date ID and OnActivate fields at click
// time.
func (c *cardView) onClick(ctx app.Context, _ app.Event) {
	if c.OnActivate != nil {
		c.OnActivate(ctx, c.ID)
	}
}

func (c *cardView) onDragStart(ctx app.Context, _ app.Event) {
	if c.OnDragStart != nil {
		c.OnDragStart(ctx, c.ID)
	}
}

func (c *cardView) onDragEnd(ctx app.Context, _ app.Event) {
	if c.OnDragEnd != nil {
		c.OnDragEnd(ctx, c.ID)
	}
}

func (c *cardView) onMouseEnter(ctx app.Context, _ app.Event) {
	if c.OnHover != nil {
		c.OnHover(ctx, c.ID)
	}
}

func (c *cardView) onMouseLeave(ctx app.Context, _ app.Event) {
	if c.OnHoverOut != nil {
		c.OnHoverOut(ctx)
	}
}

func (c *cardView) Render() app.UI {
	clickable := c.OnActivate != nil
	cls := cx(
		"card",
		c.HouseCls,
		ifCls(c.Selected, "card--selected"),
		ifCls(c.Targetable, "card--targetable"),
		ifCls(c.Dimmed, "card--dimmed"),
		ifCls(clickable && !c.Targetable, "card--clickable"),
		ifCls(c.Enter && !c.FlashOdd, "card--enter-a"),
		ifCls(c.Enter && c.FlashOdd, "card--enter-b"),
	)

	div := app.Div().Class(cls)
	if c.Draggable {
		div = div.Draggable(true).OnDragStart(c.onDragStart).OnDragEnd(c.onDragEnd)
	}
	if clickable {
		div = div.OnClick(c.onClick)
	}
	if c.OnHover != nil {
		div = div.OnMouseEnter(c.onMouseEnter).OnMouseLeave(c.onMouseLeave)
	}

	return div.Body(
		app.Div().Class("card-name").Body(
			app.If(c.Emblem != "", func() app.UI {
				return icon(
					c.Emblem,
					"icon-house",
					"icon-outline",
					ifCls(c.HouseChanged, "icon-house--changed"),
				)
			}),
			app.Span().Class("card-name-text").Text(c.Title),
			app.If(c.Stunned, func() app.UI {
				return icon("stun", "icon-token",
					ifCls(c.StunFlash && !c.FlashOdd, "icon--pulse-a"),
					ifCls(c.StunFlash && c.FlashOdd, "icon--pulse-b"))
			}),
		),
		app.Div().Class("card-body").Body(
			app.If(len(c.Stat) > 0, func() app.UI {
				return app.Div().Class("card-stat").Body(c.Stat...)
			}),
			app.If(c.Rules != "", func() app.UI {
				return app.Div().Class("card-rules").Text(c.Rules)
			}),
		),
		app.Div().Class("card-kind").Body(
			app.If(c.TypeIcon != "", func() app.UI { return icon(c.TypeIcon, "icon-kind", "icon-outline") }),
			app.Span().Text(c.Kind),
		),
	)
}

// cx joins non-empty class fragments with spaces.
func cx(parts ...string) string {
	kept := parts[:0]
	for _, p := range parts {
		if p != "" {
			kept = append(kept, p)
		}
	}
	return strings.Join(kept, " ")
}

// ifCls returns cls when cond holds, otherwise the empty string.
func ifCls(cond bool, cls string) string {
	if cond {
		return cls
	}
	return ""
}
