package engine

import "strings"

// Spend-as-pool: some cards let the Æmber sitting on a creature be spent to pay
// Æmber costs as if it were in a pool. Senator Bracchus grants it to every
// friendly creature (a ConstantAbility), The Callipygian Ideal to the one creature
// it upgrades (a StaticModifier), and Mole grants the creature's OPPONENT the right
// to spend its Æmber (a StaticModifier scoped to the opponent). The permission is a
// flat modifier the pay path consults (ADR 0005): "which creatures' Æmber may this
// player spend as pool" is a predicate over the modifiers in play, not a closure.

// SpendScope says which player a spend-as-pool permission benefits. It has no
// implicit default: a card that grants the permission must name Controller or
// Opponent explicitly (the zero value grants nothing).
type SpendScope uint8

const (
	spendScopeUnset SpendScope = iota
	// SpendByController lets the creature's own controller spend its Æmber (Senator
	// Bracchus, The Callipygian Ideal).
	SpendByController
	// SpendByOpponent lets the creature's controller's opponent spend its Æmber
	// (Mole).
	SpendByOpponent
)

// grants reports whether the scope grants any permission at all.
func (s SpendScope) grants() bool { return s != spendScopeUnset }

// allows reports whether beneficiary may spend the Æmber on a creature this scope
// governs, given that creature's controller.
func (s SpendScope) allows(beneficiary, controller int) bool {
	switch s {
	case SpendByController:
		return beneficiary == controller
	case SpendByOpponent:
		return beneficiary == 1-controller
	default:
		return false
	}
}

// spendAsPoolCreatures returns every creature in play whose Æmber player may spend
// as if it were in their pool — their own (Bracchus, Callipygian) and any enemy
// creature a permission hands them (Mole). Player's own side is listed first so the
// pay path draws from it before reaching across the board.
func (g *Game) spendAsPoolCreatures(player int) []LocalID {
	var out []LocalID
	for _, p := range [2]int{player, 1 - player} {
		for _, id := range g.State.Battleline[p].slice() {
			if g.creatureSpendableAsPool(player, id) {
				out = append(out, id)
			}
		}
	}
	return out
}

// creatureSpendableAsPool reports whether player may spend the Æmber on creature id
// as pool Æmber — an attached Upgrade grants it (StaticModifier.SpendAemberOnCard,
// The Callipygian Ideal to its controller, Mole to the opponent) or an in-play
// card's active ConstantAbility reaches it (Senator Bracchus).
func (g *Game) creatureSpendableAsPool(player int, id LocalID) bool {
	controller := g.controller(id)
	for up, ok := g.firstUpgrade(id); ok; up, ok = g.nextUpgrade(up) {
		if g.cat.def(up).Static.SpendAemberOnCard.allows(player, controller) {
			return true
		}
	}
	// Only the controller's cards are scanned: the permission is granted by the
	// side that owns the creature the Æmber sits on.
	for src, c := range g.constantAbilitiesOf(controller) {
		if c.SpendAemberOnCard.allows(player, controller) &&
			g.constantActive(src, c) && g.constantAffects(src, c, id) {
			return true
		}
	}
	return false
}

// spendAsPoolTotal is the Æmber sitting on player's spend-as-pool creatures — the
// extra a cost may draw beyond the pool itself.
func (g *Game) spendAsPoolTotal(player int) int {
	total := 0
	for _, id := range g.spendAsPoolCreatures(player) {
		total += g.AmberOn(id)
	}
	return total
}

// meetsPlayRequirement reports whether player can currently afford def's play
// requirement, counting Æmber on spend-as-pool creatures toward a requirement
// that spends (Truebaru) but not toward one that only checks the pool (Kelifi
// Dragon).
func (g *Game) meetsPlayRequirement(player int, def *CardDefinition) bool {
	r := def.PlayRequirement
	if !r.required() {
		return true
	}
	avail := g.Aember(player)
	if r.Spend {
		avail += g.spendAsPoolTotal(player)
	}
	return r.met(avail)
}

// drawFromSpendAsPool removes up to want Æmber from player's spend-as-pool
// creatures, in battleline order, and returns how much it took. The caller has
// already drawn what it can from the pool.
func (g *Game) drawFromSpendAsPool(player, want int) int {
	taken := 0
	for _, id := range g.spendAsPoolCreatures(player) {
		if taken >= want {
			break
		}
		t := min(want-taken, g.AmberOn(id))
		g.AddAmberOn(id, -t)
		taken += t
	}
	return taken
}

// spendAsPoolLines renders the spend-as-pool permission a card carries, one line
// per source: a StaticModifier an Upgrade grants its host ("This creature gains,
// …") and each ConstantAbility that reaches other creatures ("You may spend Æmber
// on friendly creatures …"). It is empty when the card grants no such permission.
func spendAsPoolLines(def *CardDefinition, hosted bool) []string {
	var lines []string
	if s := def.Static.SpendAemberOnCard; s.grants() {
		body := "You may spend Æmber on this creature as if it were in your pool."
		if s == SpendByOpponent {
			body = "Your opponent may spend Æmber on this creature as if it were in their pool."
		}
		if hosted {
			lines = append(lines, body)
		} else {
			lines = append(lines, `This creature gains, "`+body+`"`)
		}
	}
	for _, c := range def.ConstantAbilities {
		if !c.SpendAemberOnCard.grants() {
			continue
		}
		lines = append(lines,
			"You may spend Æmber on "+spendAsPoolSubject(c.target())+
				" as if it were in your pool.")
	}
	return lines
}

// spendAsPoolSubject names the creatures a spend-as-pool constant ability reaches
// in the possessive voice its printed text uses ("friendly creatures"), falling
// back to the Target's own phrasing for any other reach.
func spendAsPoolSubject(t Target) string {
	if t.Kind == TargetEachFriendlyCreature {
		return "friendly creatures"
	}
	return strings.TrimSpace(t.Text())
}
