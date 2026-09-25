package web

import (
	"strings"

	"github.com/dmikalova/vex/internal/engine"
)

// House-specific CSS classes. The actual colours live in web/app.css; these
// helpers just map a house to the .card-<slug> class defined there, which
// supplies the house's colour variables (--nm, --tp, --edge, ...) so the Go
// markup only ever carries labelled class names.
func houseSlug(h engine.House) string {
	if h == engine.HouseNone {
		return "none"
	}
	// Star Alliance's printed name has a space; the class/asset stem drops it.
	return strings.ReplaceAll(strings.ToLower(h.String()), " ", "")
}

// houseClasses returns the class that tints a card face by its house.
func houseClasses(h engine.House) string { return "card-" + houseSlug(h) }

// houseAccent returns the class supplying a house's colour variables to accent
// elements such as a house-selection button.
func houseAccent(h engine.House) string { return "card-" + houseSlug(h) }
