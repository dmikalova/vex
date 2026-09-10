package engine

import "strings"

// Spend-as-pool: some cards let the Æmber sitting on a creature be spent to pay
// Æmber costs as if it were in its controller's pool. Senator Bracchus grants it
// to every friendly creature (a ConstantAbility with SpendAsPool set), and The
// Callipygian Ideal grants it to the one creature it upgrades (a StaticModifier
// with SpendAsPool set). The permission is a flat modifier the pay path consults
// (ADR 0005): "which creatures' Æmber may be spent as pool" is a predicate over
// the modifiers in play, not a closure held in state.

// spendAsPoolCreatures returns player's creatures whose Æmber may be spent as if
// it were in player's pool, in battleline order, so the pay path can draw from
// them after the pool.
func (g *Game) spendAsPoolCreatures(player int) []LocalID {
	var out []LocalID
	for _, id := range g.State.Battleline[player].slice() {
		if g.creatureSpendableAsPool(player, id) {
			out = append(out, id)
		}
	}
	return out
}

// creatureSpendableAsPool reports whether the Æmber on creature id, controlled by
// player, may be spent as pool Æmber — either an attached Upgrade grants it
// (StaticModifier.SpendAsPool, The Callipygian Ideal) or an in-play card's active
// ConstantAbility reaches it (Senator Bracchus).
func (g *Game) creatureSpendableAsPool(player int, id LocalID) bool {
	for up, ok := g.firstUpgrade(id); ok; up, ok = g.nextUpgrade(up) {
		if g.cat.def(up).Static.SpendAemberOnCard {
			return true
		}
	}
	for _, src := range g.allInPlay(player) {
		for _, c := range g.cat.def(src).ConstantAbilities {
			if c.SpendAemberOnCard && g.constantActive(src, c) && g.constantAffects(src, c, id) {
				return true
			}
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
	avail := g.State.Aember[player]
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
	if def.Static.SpendAemberOnCard {
		body := "You may spend Æmber on this creature as if it were in your pool."
		if hosted {
			lines = append(lines, body)
		} else {
			lines = append(lines, `This creature gains, "`+body+`"`)
		}
	}
	for _, c := range def.ConstantAbilities {
		if !c.SpendAemberOnCard {
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
