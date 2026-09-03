package web

import (
	"strconv"
	"strings"

	"github.com/dmikalova/vactrol/internal/cards"
	"github.com/dmikalova/vactrol/internal/engine"
)

// This file picks which cards the Style gallery shows. Every specimen is the
// answer to a query over the card catalog, never a name written down here: a
// handwritten list is a list that stops being true the moment a set is added,
// and the whole point of the gallery is to show what the client actually has to
// render. The queries are what the captions say out loud, so a specimen always
// explains why it was chosen.

// specimen is one displayed card together with the query that selected it.
type specimen struct {
	// Caption is the query in words, e.g. "Skirmish" or "Armor > 0".
	Caption string
	// Def is the card the query found. A nil Def is a gap: the query matched
	// nothing in the loaded sets, which the gallery draws rather than skips.
	Def *engine.CardDefinition
}

// found reports whether the query matched a card.
func (s specimen) found() bool { return s.Def != nil }

// firstMatch returns the first card in the catalog satisfying want, or a gap
// specimen if none does. The catalog is sorted by house then name, so the same
// query always answers with the same card and the gallery does not reshuffle
// between renders.
func firstMatch(caption string, want func(*engine.CardDefinition) bool) specimen {
	all := cards.All()
	for i := range all {
		if want(&all[i]) {
			return specimen{Caption: caption, Def: &all[i]}
		}
	}
	return specimen{Caption: caption}
}

// houseGrid returns the House-by-Card-type table of specimens, row-major with
// one row per playable house. Its gaps are the point: an empty cell says the
// loaded sets have no card of that house and type at all, which is a fact about
// the sets that no other view of the client shows.
func houseGrid() (houses []engine.House, types []engine.CardType, cells []specimen) {
	for h := engine.HouseNone + 1; int(h) < engine.NumHouses; h++ {
		houses = append(houses, h)
	}
	types = []engine.CardType{engine.Creature, engine.Tactic, engine.Artifact, engine.Upgrade}
	for _, h := range houses {
		for _, ct := range types {
			cells = append(cells, firstMatch(
				h.String()+" "+ct.String(),
				func(d *engine.CardDefinition) bool { return d.House == h && d.Type == ct },
			))
		}
	}
	return houses, types, cells
}

// raritySpecimens returns one card per rarity, so the marks at a card's foot can
// be compared side by side rather than found by playing until one turns up.
func raritySpecimens() []specimen {
	rarities := []engine.Rarity{
		engine.Common, engine.Uncommon, engine.Rare, engine.Special, engine.Connected,
	}
	out := make([]specimen, 0, len(rarities))
	for _, r := range rarities {
		out = append(out, firstMatch(string(r), func(d *engine.CardDefinition) bool {
			return d.Rarity == r
		}))
	}
	return out
}

// featureSpecimens returns one card per visual feature a face can carry. The
// keyword rows range over engine.Keywords() rather than a list written here, so
// a keyword added to the engine shows up as a new row — as a gap, if no loaded
// card has it yet.
func featureSpecimens() []specimen {
	out := make([]specimen, 0, len(engine.Keywords())+7)
	for _, kw := range engine.Keywords() {
		out = append(out, firstMatch(kw.String(), func(d *engine.CardDefinition) bool {
			return hasKeyword(d, kw)
		}))
	}
	out = append(
		out,
		firstMatch("Assault", func(d *engine.CardDefinition) bool { return d.Assault > 0 }),
		firstMatch("Hazardous", func(d *engine.CardDefinition) bool { return d.Hazardous > 0 }),
		firstMatch("Armor > 0", func(d *engine.CardDefinition) bool { return d.Armor > 0 }),
		firstMatch("Æmber bonus", func(d *engine.CardDefinition) bool { return d.AemberBonus > 0 }),
		firstMatch(
			"Two or more traits",
			func(d *engine.CardDefinition) bool { return len(d.Traits) > 1 },
		),
		firstMatch("Constant ability", func(d *engine.CardDefinition) bool {
			return len(d.ConstantAbilities) > 0
		}),
		firstMatch("Longest rules text", longestRules()),
	)
	return out
}

// hasKeyword reports whether a printed definition carries a keyword. It reads
// the definition rather than the board, because a specimen is a printed card and
// never enters play.
func hasKeyword(d *engine.CardDefinition, kw engine.Keyword) bool {
	for _, k := range d.Keywords {
		if k == kw {
			return true
		}
	}
	return false
}

// longestRules builds a predicate matching only the card with the most rules
// text in the catalog. It is the worst case for the card face's text box, which
// is the one specimen worth choosing by extremity rather than by feature.
func longestRules() func(*engine.CardDefinition) bool {
	all := cards.All()
	longest, at := 0, -1
	for i := range all {
		if n := len(engine.RenderCardRules(&all[i])); n > longest {
			longest, at = n, i
		}
	}
	return func(d *engine.CardDefinition) bool { return at >= 0 && d.Name == all[at].Name }
}

// specimenLabel describes a specimen's card in one line for its caption, so the
// gallery can say both what was asked for and what came back.
func specimenLabel(s specimen) string {
	if !s.found() {
		return "no card"
	}
	parts := []string{s.Def.House.String(), s.Def.Type.String()}
	if s.Def.Type == engine.Creature {
		parts = append(parts, strconv.Itoa(s.Def.Power)+" power")
	}
	return strings.Join(parts, " · ")
}
