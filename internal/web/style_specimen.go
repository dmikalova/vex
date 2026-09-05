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

// attachSpecimen is a host card built on the harness carrying a fixed set of
// upgrades and under-cards, captioned with the combination it shows. The gallery
// renders it through the board's own tabbed-host path (hostWithTabs), so what it
// shows is exactly what the board draws.
type attachSpecimen struct {
	caption string
	host    engine.LocalID
}

// cardCursor hands out catalog cards in catalog order, so successive picks are
// distinct without a written-down list. It wraps at the end rather than running
// out, so a small loaded catalog still fills every slot.
type cardCursor struct {
	all []engine.CardDefinition
	i   int
}

func newCardCursor() *cardCursor { return &cardCursor{all: cards.All()} }

func (c *cardCursor) next() engine.CardDefinition {
	d := c.all[c.i%len(c.all)]
	c.i++
	return d
}

// nextOf returns the next card of a given type, falling back to any card when
// the loaded catalog has none, so the upgrade tabs read as upgrades whenever the
// sets provide them.
func (c *cardCursor) nextOf(ct engine.CardType) engine.CardDefinition {
	for range c.all {
		if d := c.next(); d.Type == ct {
			return d
		}
	}
	return c.next()
}

// underMaskLabel names one faceup/facedown combination for n under-cards, bit i
// set meaning the i-th card is facedown.
func underMaskLabel(n, mask int) string {
	parts := make([]string, n)
	for i := range n {
		if mask&(1<<i) != 0 {
			parts[i] = "down"
		} else {
			parts[i] = "up"
		}
	}
	return strings.Join(parts, " · ")
}

// buildAttachments populates the harness with the attachment specimens the
// gallery shows — 1/2/3 upgrades, 1/2/3 under-cards in every faceup/facedown
// combination, and the two composed together — and returns them in display
// order. A facedown under-card only renders as a card-back to a player who does
// not control the host, so the facedown combinations sit on the opponent's cards
// (owner 1), where the active player (owner 0) cannot peek; the upgrade rows and
// the "peeked" example sit on the active player's own cards.
func buildAttachments(g *engine.Game) []attachSpecimen {
	cur := newCardCursor()
	host := func(owner int) engine.LocalID {
		return g.AddToBattleline(cur.nextOf(engine.Creature), owner)
	}
	var out []attachSpecimen

	// Upgrades: 1, 2, 3 on the active player's own creature.
	for n := 1; n <= 3; n++ {
		h := host(0)
		for range n {
			g.AttachUpgrade(h, g.Register(cur.nextOf(engine.Upgrade), 0))
		}
		out = append(out, attachSpecimen{
			caption: strconv.Itoa(n) + " upgrade" + plural(n),
			host:    h,
		})
	}

	// Under-cards: every faceup/facedown combination of 1, 2, 3 cards, on the
	// opponent's creature so a facedown card reads as a card-back.
	for n := 1; n <= 3; n++ {
		for mask := range 1 << n {
			h := host(1)
			for i := range n {
				g.AttachUnder(h, g.Register(cur.next(), 1), mask&(1<<i) != 0)
			}
			out = append(out, attachSpecimen{
				caption: strconv.Itoa(n) + " under: " + underMaskLabel(n, mask),
				host:    h,
			})
		}
	}

	// Your own facedown under-card: you control the host, so you see it revealed
	// rather than as a back — the one case the opponent rows cannot show.
	yours := host(0)
	g.AttachUnder(yours, g.Register(cur.next(), 0), true)
	out = append(out, attachSpecimen{caption: "1 under, facedown (you peek)", host: yours})

	// Combined: upgrades and under-cards on the same opponent host.
	for _, n := range []int{2, 3} {
		h := host(1)
		for range n {
			g.AttachUpgrade(h, g.Register(cur.nextOf(engine.Upgrade), 1))
		}
		for i := range n {
			g.AttachUnder(h, g.Register(cur.next(), 1), i%2 == 1)
		}
		out = append(out, attachSpecimen{
			caption: strconv.Itoa(n) + " upgrades + " + strconv.Itoa(n) + " under",
			host:    h,
		})
	}
	return out
}

// plural returns the "s" that turns a count's noun plural, or "" for one.
func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
