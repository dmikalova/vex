package web

import (
	"strconv"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vex/internal/engine"
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
// background; extra adds sizing/placement classes. A card's own emblem stays
// hidden when its house is unset (see houseIconName), but a caller drawing an
// explicit house label — the Style gallery's house rows, a house-picker button —
// always wants something on screen, so HouseNone falls back to its own icon here.
func houseIcon(h engine.House, extra ...string) app.UI {
	name := houseIconName(h)
	if name == "" {
		name = "house-none"
	}
	return icon(name, append([]string{"icon-outline"}, extra...)...)
}

// typeIconName is the asset stem for a card type's icon.
func typeIconName(t engine.CardType) string {
	switch t {
	case engine.Creature:
		return "type-creature"
	case engine.Artifact:
		return "type-artifact"
	case engine.Tactic:
		return "type-tactic"
	case engine.Upgrade:
		return "type-upgrade"
	}
	return ""
}

// bonusIconStem is the asset stem for a bonus-icon kind.
func bonusIconStem(b engine.BonusIcon) string {
	switch b {
	case engine.BonusAember:
		return "aember"
	case engine.BonusCapture:
		return "capture"
	case engine.BonusDamage:
		return "damage"
	case engine.BonusDraw:
		return "draw"
	}
	return ""
}

// rarityMark is how a card's rarity renders wherever it is shown — the card face,
// the deck list, and the gallery filter chips. Each tier maps to its own shape via
// iconName; rarityNone (Fixed and the rest) shows nothing.
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

// iconName is the stem of the single shape a rarity renders as everywhere it is
// shown — the card face, the deck list, and the gallery filter chips: a polygon
// whose side count rises with the tier (triangle Common, square Uncommon, pentagon
// Rare, hexagon Special) and the link glyph for a Connected card. Distinct
// silhouettes, not colour, carry the tier so it stays legible to colour-blind
// players. rarityNone renders nothing.
func (m rarityMark) iconName() string {
	switch m {
	case rarityCommon:
		return "rarity-triangle"
	case rarityUncommon:
		return "rarity-square"
	case rarityRare:
		return "rarity-pentagon"
	case raritySpecial:
		return "rarity-hexagon"
	case rarityConnected:
		return "rarity-connected"
	}
	return ""
}

// deckRarityIcon renders a card's rarity as its single shape for the deck list, or
// nothing for a rarity with no mark.
func deckRarityIcon(r engine.Rarity) app.UI {
	name := rarityMarkOf(r).iconName()
	if name == "" {
		return nil
	}
	return icon(name, "icon-mark", "icon-outline")
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
