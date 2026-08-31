package web

import (
	"strconv"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/engine"
)

// assetBase is the URL prefix the dev server maps to web/assets on disk.
const assetBase = "/web/assets/"

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
