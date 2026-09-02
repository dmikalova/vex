package web

import (
	"strconv"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/engine"
)

// assetBase is the URL prefix the dev server maps to web/assets on disk.
const assetBase = "/web/assets/"

// iconOutlineFilter is a hidden inline SVG filter, injected once into the page,
// that the .icon-outline CSS rule references by id. It blurs the icon's alpha
// silhouette and then drives that blur back up to solid, which grows the shape by
// the same amount in every direction — a Gaussian is round, where feMorphology's
// square kernel came out √2 thicker on diagonals and curves. The grown shape is
// flooded black and the icon laid back on top, so the ring hugs the real geometry
// instead of the icon's box.
const iconOutlineFilter = `<svg width="0" height="0" aria-hidden="true" focusable="false" style="position:absolute">` +
	`<filter id="icon-outline" x="-25%" y="-25%" width="150%" height="150%" color-interpolation-filters="sRGB">` +
	`<feGaussianBlur in="SourceAlpha" stdDeviation="0.5" result="blur"/>` +
	`<feComponentTransfer in="blur" result="grown">` +
	`<feFuncA type="linear" slope="5"/>` +
	`</feComponentTransfer>` +
	`<feFlood flood-color="#000000" flood-opacity="0.85" result="ink"/>` +
	`<feComposite in="ink" in2="grown" operator="in" result="ring"/>` +
	`<feMerge><feMergeNode in="ring"/><feMergeNode in="SourceGraphic"/></feMerge>` +
	`</filter></svg>`

// icon renders a static SVG asset (by file stem under web/assets) as a small
// inline <img>. extra adds modifier classes for sizing/placement.
func icon(name string, extra ...string) app.UI {
	return app.Img().
		Class(cx(append([]string{"icon"}, extra...)...)).
		Src(assetBase + name + ".svg").
		Alt("")
}

// houseIconName is the asset stem for a house emblem, or "" for HouseNone.
func houseIconName(h engine.House) string {
	if h == engine.HouseNone {
		return ""
	}
	return "house-" + houseSlug(h)
}

// houseIcon renders a house emblem with the outline that keeps it legible on any
// background; extra adds sizing/placement classes.
func houseIcon(h engine.House, extra ...string) app.UI {
	return icon(houseIconName(h), append([]string{"icon-outline"}, extra...)...)
}

// typeIconName is the asset stem for a card type's icon.
func typeIconName(t engine.CardType) string {
	switch t {
	case engine.Creature:
		return "type-creature"
	case engine.Artifact:
		return "type-artifact"
	case engine.Tactic:
		return "type-action"
	case engine.Upgrade:
		return "type-upgrade"
	}
	return ""
}

// rarityMark is how a card's rarity renders at its foot. The diamond marks are
// ordered so a mark's ordinal position is its diamond count (rarityCommon is 1
// … raritySpecial is 4); rarityConnected instead shows a single "+", and
// rarityNone (Fixed and the rest) shows nothing.
type rarityMark int

const (
	rarityNone rarityMark = iota
	rarityCommon
	rarityUncommon
	rarityRare
	raritySpecial
	rarityConnected
)

// rarityMarkOf maps a card's rarity to the mark shown at its foot.
func rarityMarkOf(r engine.Rarity) rarityMark {
	switch r {
	case engine.Common:
		return rarityCommon
	case engine.Uncommon:
		return rarityUncommon
	case engine.Rare:
		return rarityRare
	case engine.Special:
		return raritySpecial
	case engine.Connected:
		return rarityConnected
	}
	return rarityNone
}

// diamonds is how many rarity diamonds the mark shows. The diamond marks are
// consecutive from rarityCommon (1) to raritySpecial (4), so a mark in that
// range is its own count; every other mark shows no diamonds.
func (m rarityMark) diamonds() int {
	if m >= rarityCommon && m <= raritySpecial {
		return int(m)
	}
	return 0
}

// isConnected reports whether the mark is a Connected card's single "+".
func (m rarityMark) isConnected() bool { return m == rarityConnected }

// rarityDiamonds renders n identical rarity diamonds; the count is the rarity
// (one for Common up to four for Special). Each carries the hard outline so the
// diamonds read against any card art.
func rarityDiamonds(n int) []app.UI {
	out := make([]app.UI, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, icon("rarity-diamond", "icon-mark", "icon-outline"))
	}
	return out
}

// keyColorIconName is the asset stem for a forged key's colour.
func keyColorIconName(c engine.KeyColor) string {
	switch c {
	case engine.KeyColorRed:
		return "key-red"
	case engine.KeyColorBlue:
		return "key-blue"
	case engine.KeyColorYellow:
		return "key-yellow"
	}
	return ""
}

// keyColorClass is the modifier class that paints a control in a key's colour,
// so choosing a key colour is done by clicking that colour rather than by
// reading its name.
func keyColorClass(c engine.KeyColor) string {
	switch c {
	case engine.KeyColorRed:
		return "key-choice--red"
	case engine.KeyColorBlue:
		return "key-choice--blue"
	case engine.KeyColorYellow:
		return "key-choice--yellow"
	}
	return ""
}

// keyChoiceButton is one key colour offered as a choice: a button in that
// colour, its key icon sparkling, labelled with the colour's name.
func keyChoiceButton(
	c engine.KeyColor,
	label string,
	cursor bool,
	onClick app.EventHandler,
) app.UI {
	return app.Button().
		Class(cx("house-btn", "key-choice", keyColorClass(c), ifCls(cursor, "btn-cursor"))).
		OnClick(onClick).
		Body(
			app.Span().Class("key-sparkle").Body(icon(keyColorIconName(c), "icon-inline")),
			app.Text(label),
		)
}

// keyColorByName resolves a key-colour label to its value, or KeyColorNone.
func keyColorByName(name string) engine.KeyColor {
	switch name {
	case "Red":
		return engine.KeyColorRed
	case "Blue":
		return engine.KeyColorBlue
	case "Yellow":
		return engine.KeyColorYellow
	}
	return engine.KeyColorNone
}

// statSeg is a stat value followed by its icon (e.g. "3" + the power icon).
// extra adds modifier classes (e.g. a one-shot pulse animation).
func statSeg(n int, iconName string, extra ...string) app.UI {
	return app.Span().Class(cx(append([]string{"stat-seg"}, extra...)...)).Body(
		app.Text(strconv.Itoa(n)),
		icon(iconName, "icon-stat"),
	)
}
