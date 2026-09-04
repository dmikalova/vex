package web

import (
	"strconv"
	"strings"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/engine"
)

// This file holds the small helpers the views share: the pieces of a card face
// (labels, stat lines, rules text), the questions asked of a hand card, and the
// odds and ends (btn, containsID) too small to have a home of their own.

// boardCardID and handCardID name a card's element in the page. The same card is
// drawn in both places over its life, so the two are kept distinct and the play
// animation can measure the hand slot it left and the board slot it arrived in.
func boardCardID(id engine.LocalID) string { return "card-" + strconv.Itoa(int(id)) }
func handCardID(id engine.LocalID) string  { return "hand-" + strconv.Itoa(int(id)) }

func btn(label string, h app.EventHandler, class string) app.UI {
	return app.Button().Class(class).Text(label).OnClick(h)
}

// cardFace builds the plain face of a card that is on the table — everything the
// card itself says, and none of the board's interaction. The board's own cards
// (renderCard) add selection, targeting, and handlers on top; the readers that
// only show a card — the hover preview and the lifted copy — use it as it stands.
func (g *game) cardFace(id engine.LocalID) *cardView {
	def := g.g.Def(id)
	house := g.g.House(id)
	return &cardView{
		Title:         def.Name,
		HouseCls:      houseClasses(house),
		Emblem:        houseIconName(house),
		HouseChanged:  house != def.House,
		TypeIcon:      typeIconName(def.Type),
		Stat:          g.statLine(id),
		Rules:         g.faceRules(id),
		Kind:          kindLabel(def),
		Trait:         traitLabel(def),
		Rarity:        rarityMarkOf(def.Rarity),
		Maverick:      g.isMaverick(id),
		Stunned:       g.g.Stunned(id),
		Exhausted:     g.g.Exhausted(id),
		Bar:           g.barKeywords(id),
		TauntShielded: def.Type == engine.Creature && g.g.TauntShielded(id),
	}
}

// kindLabel is a card's foot label: its type (e.g. "Creature"). Traits render
// separately as their own body line (traitLabel).
func kindLabel(def *engine.CardDefinition) string {
	return def.Type.String()
}

// traitLabel renders a card's traits in KeyForge order (e.g. "Human • Knight"),
// or "" when it has none. It is shown as its own line in the card body, under the
// stat line and above the rules text.
func traitLabel(def *engine.CardDefinition) string {
	parts := make([]string, 0, len(def.Traits))
	for _, t := range def.Traits {
		parts = append(parts, t.String())
	}
	return strings.Join(parts, " • ")
}

func (g *game) statLine(id engine.LocalID) []app.UI {
	def := g.g.Def(id)
	f := g.flashes[id]
	var segs []app.UI
	if def.Type == engine.Creature {
		segs = append(segs, statSeg(g.g.Power(id), "power", pulseClass(f.power, f.odd, "pow")))
		if d := g.g.Damage(id); d > 0 {
			segs = append(segs, statSeg(d, "damage", pulseClass(f.damage, f.odd, "dmg")))
		}
		if a := g.g.Armor(id); a > 0 {
			segs = append(segs, statSeg(a, "shield"))
		}
	}
	if a := g.g.AmberOn(id); a > 0 {
		segs = append(segs, statSeg(a, "aember", pulseClass(f.amber, f.odd, "gain")))
	}
	// Stun and exhaustion show as tokens on the face (see cardView), so they need
	// no stat icon here.
	return segs
}

// pulseClass returns the alternating one-shot animation class for a stat segment
// (or "" when it is not flashing). kind selects the colour: dmg (red), pow (cyan),
// gain (gold). The -a/-b pair alternates so the animation replays on repeats.
func pulseClass(on, odd bool, kind string) string {
	switch {
	case !on:
		return ""
	case odd:
		return "stat-seg--" + kind + "-b"
	default:
		return "stat-seg--" + kind + "-a"
	}
}

func handStat(def *engine.CardDefinition) []app.UI {
	var segs []app.UI
	if def.Type == engine.Creature {
		segs = append(segs, statSeg(def.Power, "power"))
		if def.Armor > 0 {
			segs = append(segs, statSeg(def.Armor, "shield"))
		}
	}
	if def.AemberBonus > 0 {
		segs = append(segs, statSeg(def.AemberBonus, "aember"))
	}
	return segs
}

// faceRules is the rules text shown on a card in play: its own text followed by
// each attached upgrade's name and text, since an upgrade's ability applies to
// this creature and the upgrade has no face of its own to read. The upgrade's
// text is rendered as it reads on its host, so it says `Reap: …` rather than
// repeating "This creature gains" on the creature it is already sitting on.
func (g *game) faceRules(id engine.LocalID) string {
	var lines []string
	if s := engine.RenderCardRules(g.g.Def(id)); s != "" {
		lines = append(lines, displayRules(s))
	}
	for _, up := range g.g.Upgrades(id) {
		def := g.g.Def(up)
		line := "↳ " + def.Name
		if s := engine.RenderUpgradeOnCreature(def); s != "" {
			line += ": " + displayRules(s)
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

// displayRules re-renders a result gate's "->" as an arrow glyph, the way a
// physical card would print it. The plain "->" is the canonical text — it drives
// generated doc comments and the rulebook (docs/card-wording-rules.md §5) and
// stays that way — so the glyph swap happens here, once, for display only.
func displayRules(rules string) string {
	return strings.ReplaceAll(rules, " -> ", " → ")
}

// playableFromHand reports whether the active player can play the given hand card
// right now, so unplayable cards are not draggable.
func (g *game) playableFromHand(id engine.LocalID) bool {
	return g.g.CanPlay(g.active(), id) == nil
}

// discardableFromHand reports whether the active player may discard the given
// hand card: discarding needs only that the card is of the active house.
func (g *game) discardableFromHand(id engine.LocalID) bool {
	h := g.g.State.ActiveHouse
	return h != engine.HouseNone && g.g.Def(id).House == h
}

// usableFromHand reports whether a hand card can be acted on at all this turn —
// played or, failing that, discarded — so a card that is only discardable still
// reads as live rather than being lowlighted with the dead ones.
func (g *game) usableFromHand(id engine.LocalID) bool {
	return g.playableFromHand(id) || g.discardableFromHand(id)
}

func containsID(ids []engine.LocalID, id engine.LocalID) bool {
	for _, x := range ids {
		if x == id {
			return true
		}
	}
	return false
}

func indexOfID(ids []engine.LocalID, id engine.LocalID) int {
	for i, x := range ids {
		if x == id {
			return i
		}
	}
	return -1
}
