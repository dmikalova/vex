package web

import (
	"strings"
	"unicode/utf8"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/engine"
)

// cardView is a presentational component for a single card face. It carries no
// game logic: the parent supplies already-rendered strings, visual flags, and a
// click handler, so the same component renders a hand card, a creature in play,
// an artifact, or a targeting candidate.
type cardView struct {
	app.Compo

	ID engine.LocalID
	// DOMID names the card's element so a later render can find it in the page (the
	// play animation measures where the card landed); "" leaves it unnamed.
	DOMID    string
	Title    string
	HouseCls string // house-derived border/background classes
	Emblem   string // house emblem asset stem ("" for none)
	// HouseChanged marks that the card's current house differs from its printed
	// house (a control/"belongs to house" effect); the emblem is highlighted.
	HouseChanged bool
	TypeIcon     string   // card-type icon asset stem
	Stat         []app.UI // compact stat nodes (power, damage, Æmber… with icons)
	Rules        string   // rules/ability text for the face
	// Trait is the card's trait line (e.g. "Human • Knight"), shown in the body
	// under the stat line and above the rules; "" when the card has no traits.
	Trait string
	Kind  string // card type label shown at the foot
	// Rarity is the card's rarity mark at the foot: diamonds for Common…Special, a
	// "+" for Connected, nothing for Fixed/none. Maverick shows the maverick emblem
	// beside it (a card rehoused off its printed house).
	Rarity   rarityMark
	Maverick bool
	Stunned  bool // shows a stun token on the face
	// Exhausted shows an exhausted token on the face. Rotating the card the way a
	// physical one turns would break the strip's grid, so the token stands in.
	Exhausted bool
	// Bar lists the keywords that get a coloured stripe along the card's edge
	// (Taunt, Elusive, Hazardous…), since those change how the card can be attacked
	// and are worth seeing without reading the rules text. BarBottom moves the
	// stripe to the bottom edge, for the rows facing the active player across the
	// midline.
	Bar       []engine.Keyword
	BarBottom bool
	// Enter pulses the whole card as it comes into play; Fight shakes it as it
	// attacks or is attacked; Hit washes it red as it takes damage; StunFlash and
	// ExhaustFlash pulse their token as it is first applied. FlashOdd alternates
	// their animation class each time so the CSS animation replays even on
	// back-to-back triggers.
	Enter bool
	Fight bool
	Hit   bool
	// FightDown lunges the fight animation downward instead of up, for the cards in
	// the top battleline, so the two combatants move towards each other.
	FightDown    bool
	StunFlash    bool
	ExhaustFlash bool
	FlashOdd     bool
	Selected     bool
	Targetable   bool
	Dimmed       bool
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
		ifCls(c.Fight && !c.FlashOdd, "card--fight-a"),
		ifCls(c.Fight && c.FlashOdd, "card--fight-b"),
		ifCls(c.Fight && c.FightDown, "card--fight-down"),
		ifCls(c.Hit && !c.FlashOdd, "card--hit-a"),
		ifCls(c.Hit && c.FlashOdd, "card--hit-b"),
	)

	div := app.Div().Class(cls)
	if c.DOMID != "" {
		div = div.ID(c.DOMID)
	}
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
		app.If(len(c.Bar) > 0, func() app.UI {
			return app.Div().
				Class(cx("card-keybar", ifCls(c.BarBottom, "card-keybar--bottom"))).
				Body(
					app.Range(c.Bar).Slice(func(i int) app.UI {
						kw := c.Bar[i].String()
						return app.Div().
							Class("card-keybar-seg card-keybar--" + strings.ToLower(kw)).
							Title(kw)
					}),
				)
		}),
		app.Div().Class("card-name").Body(
			app.If(c.Emblem != "", func() app.UI {
				return icon(
					c.Emblem,
					"icon-house",
					"icon-outline",
					ifCls(c.HouseChanged, "icon-house--changed"),
				)
			}),
			app.Span().Class(cx("card-name-text", squeezeClass(c.Title))).Text(c.Title),
		),
		app.Div().Class("card-body").Body(
			app.If(len(c.Stat) > 0 || c.Stunned || c.Exhausted, func() app.UI {
				return app.Div().Class("card-stat").Body(
					app.Range(c.Stat).Slice(func(i int) app.UI { return c.Stat[i] }),
					// Stun and exhaustion read as more of the card's current condition, so
					// they sit at the end of the stat line rather than in its name banner.
					app.If(c.Stunned || c.Exhausted, func() app.UI {
						return app.Div().Class("card-tokens").Body(
							app.If(c.Stunned, func() app.UI {
								return icon("stun", "icon-token", "icon-outline",
									ifCls(c.StunFlash && !c.FlashOdd, "icon--pulse-a"),
									ifCls(c.StunFlash && c.FlashOdd, "icon--pulse-b"))
							}),
							app.If(c.Exhausted, func() app.UI {
								return icon("exhausted", "icon-token", "icon-outline",
									ifCls(c.ExhaustFlash && !c.FlashOdd, "icon--pulse-a"),
									ifCls(c.ExhaustFlash && c.FlashOdd, "icon--pulse-b"))
							}),
						)
					}),
				)
			}),
			app.If(c.Trait != "", func() app.UI {
				return app.Div().Class("card-traits").Text(c.Trait)
			}),
			app.If(c.Rules != "", func() app.UI {
				return app.Div().Class("card-rules").Text(c.Rules)
			}),
		),
		app.Div().Class("card-kind").Body(
			app.If(c.TypeIcon != "", func() app.UI { return icon(c.TypeIcon, "icon-kind", "icon-outline") }),
			app.Span().Text(c.Kind),
			app.If(c.Maverick || c.Rarity != rarityNone, func() app.UI {
				return app.Div().Class("card-marks").Body(
					app.If(c.Maverick, func() app.UI { return icon("maverick", "icon-mark", "icon-outline") }),
					app.If(c.Rarity.diamonds() > 0, func() app.UI {
						return app.Div().
							Class("rarity-diamonds").
							Body(rarityDiamonds(c.Rarity.diamonds())...)
					}),
					app.If(c.Rarity.isConnected(), func() app.UI {
						return app.Span().Class("rarity-plus").Text("+")
					}),
				)
			}),
		),
	)
}

// squeezeClass picks how hard to condense a title that would otherwise be cut
// off. The banner is a fixed width, so the step is chosen from the title's length
// rather than measured — close enough at these few sizes, and it costs no layout
// read per card.
func squeezeClass(title string) string {
	switch n := utf8.RuneCountInString(title); {
	case n > 26:
		return "card-name-text--squeeze-3"
	case n > 22:
		return "card-name-text--squeeze-2"
	case n > 18:
		return "card-name-text--squeeze-1"
	}
	return ""
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
