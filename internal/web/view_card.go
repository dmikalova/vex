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

// playerNameCls is the colour class that tints a player's name in that player's
// own colour (p0 green, p1 yellow), so who is named is read off the colour the
// same way the log bubbles and score pill already are. Every place a player's
// name is drawn wears it, so a name never renders in a neutral or the wrong
// colour — including the end-of-game "wins!" banner.
func playerNameCls(player int) string { return "player-name--p" + strconv.Itoa(player) }

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
	// The keybar is a fact about a card on the table — its granted keywords
	// included — so a card being read in hand does not draw one.
	var bar []string
	var taunted bool
	if g.inPlay(id) {
		bar = g.barKeywords(id)
		taunted = def.Type == engine.Creature && g.g.TauntShielded(id)
	}
	return &cardView{
		Title:         def.Name,
		HouseCls:      houseClasses(house),
		Emblem:        houseIconName(house),
		HouseChanged:  house != def.House,
		TypeIcon:      typeIconName(def.Type),
		Stat:          g.statLine(id),
		Rules:         g.faceRules(id),
		Icons:         cardGlyphs(def),
		Bonuses:       def.Bonuses,
		Kind:          kindLabel(def),
		Trait:         traitLabel(def),
		Rarity:        rarityMarkOf(def.Rarity),
		Maverick:      g.isMaverick(id),
		Legacy:        g.isLegacy(id),
		Stunned:       g.g.Stunned(id),
		Warded:        g.g.Warded(id),
		Enraged:       g.g.Enraged(id),
		Exhausted:     g.g.Exhausted(id),
		InPlay:        g.inPlay(id),
		Bar:           bar,
		TauntShielded: taunted,
	}
}

// inPlay reports whether a card is on the table — in either battleline or either
// artifact row — as opposed to in a hand or an out-of-play pile.
func (g *game) inPlay(id engine.LocalID) bool {
	for p := range 2 {
		if containsID(g.g.Battleline(p), id) || containsID(g.g.Artifacts(p), id) {
			return true
		}
	}
	return false
}

// kindLabel is a card's foot label: its type (e.g. "Creature"). Traits render
// separately as their own body line (traitLabel).
func kindLabel(def *engine.CardDefinition) string {
	return engine.CardTypeLabel(def)
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
	f := g.flashes[id]
	var segs []app.UI
	if g.g.TypeOf(id) == engine.Creature {
		segs = append(segs, statSeg(g.g.Power(id), "power", pulseClass(f.power, f.odd, "pow")))
		if d := g.g.Damage(id); d > 0 {
			segs = append(segs, statSeg(d, "damage", pulseClass(f.damage, f.odd, "dmg")))
		}
		// The armor a creature has left to absorb damage this turn (its full armor
		// minus what it has already spent), so the shield count falls as hits land
		// and refreshes when the creature readies — not the printed maximum.
		if a := int(g.g.State.Cards[id].ArmorRemaining); a > 0 {
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
	return segs
}

// faceRules is the rules text shown on a card in play: its own text followed by
// each attached upgrade's name and text, since an upgrade's ability applies to
// this creature and the upgrade has no face of its own to read. The upgrade's
// text is rendered as it reads on its host, so it says `Reap: …` rather than
// repeating "This creature gains" on the creature it is already sitting on.
func (g *game) faceRules(id engine.LocalID) string {
	def := g.g.Def(id)
	var lines []string
	if s := faceText(def, engine.RenderCardRules(def)); s != "" {
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

// faceText is a card's rules text ready for the face: the engine's plain text
// with its trailing "Enhance …" line re-encoded so the named bonus icons render
// as inline glyphs (richText). It is a no-op for a card that is not an Enhance
// source.
func faceText(def *engine.CardDefinition, rules string) string {
	if len(def.Enhances) == 0 {
		return rules
	}
	stems := make([]string, len(def.Enhances))
	for i, b := range def.Enhances {
		stems[i] = inlineIcon(bonusIconStem(b))
	}
	// The glyphs abut: unlike the engine's spelled-out bonus names, an icon run
	// reads as one strip, the way it is printed on the card.
	enhance := "Enhance " + strings.Join(stems, "") + "."
	if i := strings.LastIndexByte(rules, '\n'); i >= 0 {
		return rules[:i+1] + enhance
	}
	return enhance
}

// playableFromHand reports whether the active player can play the given hand card
// right now, so unplayable cards are not draggable.
func (g *game) playableFromHand(id engine.LocalID) bool {
	return g.g.CanPlay(g.active(), id) == nil
}

// discardableFromHand reports whether the active player may discard the given
// hand card right now. It asks the engine (CanDiscard) rather than re-deriving the
// rule, so the first-turn one-card restriction bars discarding exactly as it bars
// playing — after the opening volition no hand card reads as live.
func (g *game) discardableFromHand(id engine.LocalID) bool {
	return g.g.CanDiscard(g.active(), id) == nil
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
