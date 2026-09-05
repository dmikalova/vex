package web

import (
	"fmt"
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
	// PowerCounters is the net power from +1/-1 counters on the card. A creature
	// carrying counters shows a +1 (or -1) power token beside its stun/exhaust
	// tokens, with the count when more than one, so how many tokens ride the card
	// is legible for the interactions that care about it, not just the net power.
	PowerCounters int
	// (Taunt, Elusive, Hazardous…) — printed keywords plus Hazardous, a magnitude
	// rather than a boolean keyword so it does not live in engine.Keyword — since
	// those change how the card can be attacked and are worth seeing without
	// reading the rules text. BarBottom moves the stripe to the bottom edge, for
	// the rows facing the active player across the midline.
	Bar       []string
	BarBottom bool
	// TauntShielded adds a fainter taunt-coloured segment to the same stripe for a
	// creature that does not have taunt itself but sits beside one that does, so a
	// creature most attackers cannot reach reads as protected without being
	// confused for a taunter.
	TauntShielded bool
	// Enter pulses the whole card as it comes into play; Fight shakes it as it
	// attacks or is attacked; Hit washes it red as it takes damage; Reap squeezes it
	// and washes it yellow as it reaps; Act does the same in green as it uses an
	// action ability (an artifact's, or a creature's own); StunFlash and ExhaustFlash
	// pulse their token as it is first applied. FlashOdd alternates their animation
	// class each time so the CSS animation replays even on back-to-back triggers.
	Enter bool
	Fight bool
	Hit   bool
	Reap  bool
	Act   bool
	// FightDown lunges the fight animation downward instead of up, for the cards in
	// the top battleline, so the two combatants move towards each other.
	FightDown    bool
	StunFlash    bool
	ExhaustFlash bool
	PowerFlash   bool
	FlashOdd     bool
	Selected     bool
	Targetable   bool
	Dimmed       bool
	// Jiggle wobbles the card once to draw the eye: it is played on the cards a
	// player could still act with when they try to end the turn with moves left, so
	// the end-turn confirm points at what it is warning about.
	Jiggle bool
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

// powerCounterToken draws the +1 (or -1) power-counter token in the card's status
// row, with the count when more than a single counter rides the card. The net
// counters carry the sign, so a positive net shows the plus token and a negative
// net the minus one; the pulse replays when the count last changed.
func (c *cardView) powerCounterToken() app.UI {
	name, n := "power-counter-plus", c.PowerCounters
	if n < 0 {
		name, n = "power-counter-minus", -n
	}
	return app.Div().Class("card-counter").Body(
		icon(name, "icon-token", "icon-outline",
			ifCls(c.PowerFlash && !c.FlashOdd, "icon--pulse-a"),
			ifCls(c.PowerFlash && c.FlashOdd, "icon--pulse-b")),
		app.If(n > 1, func() app.UI {
			return app.Span().Class("card-counter-num").Text(fmt.Sprintf("%d", n))
		}),
	)
}

// kwColorVar is the CSS var() reference for a keybar entry's colour, e.g.
// "var(--kw-taunt)" — used to build the keybar's gradient.
func kwColorVar(name string) string {
	return "var(--kw-" + strings.ToLower(name) + ")"
}

// barTitle joins a card's keybar entries for the segment's hover tooltip, e.g.
// "Taunt, Elusive".
func barTitle(bar []string) string {
	return strings.Join(bar, ", ")
}

// keybarBlend is how many percentage points of the keybar's total width blend
// into the neighbouring keyword on either side of a seam — kept small so the
// stripe reads as a solid band per keyword with a tight transition, not a soft
// wash between them.
const keybarBlend = 1.5

// keybarGradient builds one left-to-right gradient spanning every entry in
// bar, each given an equal-width solid band with a tight blended seam where it
// meets its neighbour, so a multi-keyword card's edge reads as a single
// continuous stripe rather than tiled solid blocks. It always returns a
// gradient, even for a single entry — background-image (unlike
// background-color) cannot take a bare color, only an <image>, so a lone
// keyword still needs the degenerate two-stop form to render at all.
func keybarGradient(bar []string) string {
	n := len(bar)
	stops := make([]string, 0, 2*n)
	stops = append(stops, kwColorVar(bar[0])+" 0%")
	for i := range n {
		end := float64(i+1) / float64(n) * 100
		if i == n-1 {
			stops = append(stops, kwColorVar(bar[i])+" 100%")
			break
		}
		stops = append(stops, fmt.Sprintf("%s %.4g%%", kwColorVar(bar[i]), end-keybarBlend))
		stops = append(stops, fmt.Sprintf("%s %.4g%%", kwColorVar(bar[i+1]), end+keybarBlend))
	}
	return "linear-gradient(to right, " + strings.Join(stops, ", ") + ")"
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
		ifCls(c.Reap && !c.FlashOdd, "card--reap-a"),
		ifCls(c.Reap && c.FlashOdd, "card--reap-b"),
		ifCls(c.Act && !c.FlashOdd, "card--act-a"),
		ifCls(c.Act && c.FlashOdd, "card--act-b"),
		ifCls(c.Jiggle, "card--jiggle"),
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
		app.If(len(c.Bar) > 0 || c.TauntShielded, func() app.UI {
			return app.Div().
				Class(cx("card-keybar", ifCls(c.BarBottom, "card-keybar--bottom"))).
				Body(
					app.If(len(c.Bar) > 0, func() app.UI {
						return app.Div().
							Class("card-keybar-seg").
							Style("flex", fmt.Sprintf("%d 1 0%%", len(c.Bar))).
							Style("background-image", keybarGradient(c.Bar)).
							Title(barTitle(c.Bar))
					}),
					app.If(c.TauntShielded, func() app.UI {
						return app.Div().
							Class("card-keybar-seg card-keybar--taunt-shielded").
							Title("Shielded by a neighboring Taunt creature")
					}),
				)
		}),
		app.Div().Class("card-name").Body(
			squeezedTitle(c.Title),
		),
		app.Div().Class("card-body").Body(
			app.If(len(c.Stat) > 0 || c.Stunned || c.Exhausted || c.PowerCounters != 0, func() app.UI {
				return app.Div().Class("card-stat").Body(
					app.Range(c.Stat).Slice(func(i int) app.UI { return c.Stat[i] }),
					// Stun and exhaustion read as more of the card's current condition, so
					// they sit at the end of the stat line rather than in its name banner.
					app.If(c.Stunned || c.Exhausted || c.PowerCounters != 0, func() app.UI {
						return app.Div().Class("card-tokens").Body(
							app.If(c.PowerCounters != 0, func() app.UI {
								return c.powerCounterToken()
							}),
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
			app.If(c.Emblem != "", func() app.UI {
				return icon(
					c.Emblem,
					"icon-house",
					"icon-outline",
					ifCls(c.HouseChanged, "icon-house--changed"),
				)
			}),
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
						return icon(
							"rarity-connected",
							"icon-mark",
							"icon-outline",
							"rarity-connected-icon",
						)
					}),
				)
			}),
		),
		// The selection and targetable rings are drawn inside the card edge and
		// below the keybar, so they slip under the keyword stripe instead of cutting
		// across it and never peek past the card's top or bottom edge. Both use the
		// same .card-selection overlay; targetable is rendered first so a card that
		// is both shows the yellow selection ring over the green targetable one.
		app.If(c.Targetable, func() app.UI {
			return app.Div().Class("card-selection card-selection--targetable")
		}),
		app.If(c.Selected, func() app.UI {
			return app.Div().Class("card-selection")
		}),
	)
}

// squeezeBaseChars is how many runes of the card font fit the name banner
// unsqueezed, and squeezeMinScale is the floor below which a title is left to
// the ellipsis instead of being squeezed illegibly thin. Both are calibrated
// against the rendered banner (128px available at this font) rather than
// measured per card, so a title's condensing is a continuous function of its
// length instead of a few discrete steps that either under- or over-squeeze
// whichever titles land near a step's edge.
const (
	squeezeBaseChars = 20
	squeezeMinScale  = 0.7
)

// squeezeScale reports how much to horizontally condense a title that would
// otherwise be cut off, 1 meaning no condensing at all.
func squeezeScale(title string) float64 {
	n := utf8.RuneCountInString(title)
	if n <= squeezeBaseChars {
		return 1
	}
	if scale := float64(squeezeBaseChars) / float64(n); scale > squeezeMinScale {
		return scale
	}
	return squeezeMinScale
}

// squeezedTitle renders a card's name banner text, horizontally condensed via
// scaleX when squeezeScale calls for it. The width is grown by the inverse of
// the scale so the squeeze has the extra room to condense — flex-shrink must be
// disabled alongside it, or the flex layout claws that extra width back before
// the transform ever sees it.
func squeezedTitle(title string) app.UI {
	span := app.Span().Class("card-name-text").Text(title)
	scale := squeezeScale(title)
	if scale >= 1 {
		return span
	}
	return span.
		Style("transform", fmt.Sprintf("scaleX(%.3f)", scale)).
		Style("width", fmt.Sprintf("%.2f%%", 100/scale)).
		Style("flex-shrink", "0")
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
