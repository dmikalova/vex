package engine

// This file holds combat: an attacker fighting a defender, dealing damage (with
// armor, Skirmish, Assault, and Hazardous), and deciding which creatures the
// damage destroys. The destruction itself is carried out in game_destroy.go.

// Combat: use one of your ready creatures to fight an enemy creature. Using it to
// fight exhausts it. First, any "Before Fight" abilities and the Assault and
// Hazardous keywords resolve; if these destroy either creature, the fight does not
// occur. Otherwise both creatures deal damage equal to their power at the same
// time — armor reduces the damage a creature takes, and a Skirmish attacker takes
// no damage back. A creature with damage equal to or greater than its power is
// destroyed, and both deaths are resolved together, so neither fighter's
// destruction changes the damage the other deals. Finally, if the attacker
// survived, its "Fight" abilities resolve.
//
//rulebook:combat Combat

// fight resolves combat between an attacker and a defender. Both deal damage equal
// to their power simultaneously; Skirmish prevents the attacker taking damage back.
func (g *Game) fight(attacker, defender LocalID) {
	g.recordUse(attacker)
	// Read both controllers before combat, which may take either creature off the
	// board and with it the controller the fight-kill tally is recorded against.
	attackerSide, defenderSide := g.controller(attacker), g.controller(defender)
	if g.recoverFromStun(attacker) {
		return
	}
	// Using a creature to fight exhausts it before anything else resolves, so a
	// "Before Fight" ability already sees the attacker exhausted. The defender is
	// put in context, so an ability can act on "the creature this fights".
	g.State.Cards[attacker].Exhausted = true
	g.triggerAbilities(attacker, TriggerBeforeFight, defender, true)

	// A "Before Fight" ability may redirect the attacker's fight damage to another
	// creature (Gabos Longarms) or make the fight not occur (Evasion Sigil). Read
	// and clear those per-fight flags before the combat step.
	redirect := g.State.FightDamageRedirect
	g.State.FightDamageRedirect = 0
	cancelled := g.State.FightCancelled
	g.State.FightCancelled = false
	if cancelled {
		g.logf("%s's fight does not occur", g.Name(attacker))
		return
	}

	// Assault and Hazardous deal their damage before fight damage: the attacker's
	// Assault hits the defender, the defender's Hazardous hits the attacker.
	// Either can destroy a fighter before combat. Skipped if a "Before Fight"
	// effect already removed the defender.
	if g.inPlay(defender) {
		var pre []DamageTarget
		if a := g.assault(attacker); a > 0 {
			pre = append(pre, DamageTarget{ID: defender, Amount: a})
		}
		if h := g.hazardous(defender); h > 0 {
			pre = append(pre, DamageTarget{ID: attacker, Amount: h})
		}
		if len(pre) > 0 {
			g.dealDamage(g.controller(attacker), pre...)
		}
	}

	// Elusive replaces the fight itself, so it is read only once the whole pre-fight
	// sequence — "Before Fight" abilities, Assault, Hazardous — has resolved. It is
	// still spent when that sequence leaves a fighter dead and no damage is
	// exchanged: the creature was attacked.
	elusive := g.spendElusive(attacker, defender)

	// Combat damage is exchanged only while both fighters are still in play (a
	// "Before Fight" effect, Assault, or Hazardous can remove one first).
	if g.inPlay(attacker) && g.inPlay(defender) {
		ap, dp := g.Power(attacker), g.Power(defender)
		g.logf("%s (%d power) fights %s (%d power)", g.Name(attacker), ap, g.Name(defender), dp)
		if elusive {
			g.logf("%s is elusive — no fight damage is dealt", g.Name(defender))
		} else {
			// Both fighters take their damage simultaneously, then destruction is
			// resolved together as part of dealing it — so each dying creature is
			// already in the discard before the other's "Destroyed:" ability (or the
			// attacker's "After Fight") fires, and neither death changes the other's.
			dmgTarget := defender
			if redirect != 0 {
				dmgTarget = redirect
			}
			targets := []DamageTarget{{ID: dmgTarget, Amount: g.fightDamage(attacker, defender)}}
			if !g.hasKeyword(attacker, Skirmish) {
				targets = append(targets, DamageTarget{ID: attacker, Amount: dp})
			}
			g.dealDamage(g.controller(attacker), targets...)
		}
	}
	g.triggerAbilities(attacker, TriggerAfterFight, 0, false)

	// "After a creature is destroyed fighting X": when exactly one combatant is
	// removed by the fight, the survivor's ability fires with the destroyed
	// creature as `it`.
	attackerDead, defenderDead := !g.inPlay(attacker), !g.inPlay(defender)
	// A creature destroyed in a fight is an enemy kill from the other side's point of
	// view, which is the tally The Warchest is paid for.
	if attackerDead {
		g.State.TurnHistory[1-attackerSide][EnemyCreaturesFightKilled]++
	}
	if defenderDead {
		g.State.TurnHistory[1-defenderSide][EnemyCreaturesFightKilled]++
	}
	if defenderDead && !attackerDead {
		g.triggerAbilities(attacker, TriggerAfterDestroyedFighting, defender, true)
	}
	if attackerDead && !defenderDead {
		g.triggerAbilities(defender, TriggerAfterDestroyedFighting, attacker, true)
	}
	g.emitCardUsed(g.controller(attacker), attacker)
	g.emitLasting(EventFight, g.controller(attacker), attacker)
}

// spendElusive reports whether the defender's Elusive keyword stops the pending
// fight damage of this fight, marking it spent for the turn. Elusive applies the
// first time an elusive creature is attacked each turn, and is read after the
// pre-fight sequence so a "Before Fight" ability, Assault, or Hazardous resolves
// against the creature normally. Damage from keywords and abilities is unaffected.
func (g *Game) spendElusive(attacker, defender LocalID) bool {
	if g.attackIgnores(attacker, Elusive) ||
		!g.hasKeyword(defender, Elusive) ||
		g.State.Cards[defender].ElusiveUsedThisTurn {
		return false
	}
	g.State.Cards[defender].ElusiveUsedThisTurn = true
	return true
}

// attackIgnores reports whether an attacking creature ignores a defensive keyword
// while it attacks — Niffle Ape ignores taunt and elusive.
func (g *Game) attackIgnores(attacker LocalID, k Keyword) bool {
	for _, kw := range g.cat.def(attacker).AttackIgnores {
		if kw == k {
			return true
		}
	}
	return false
}

// protectedByTaunt reports whether a creature cannot be chosen to be fought by an
// attacker: a taunter shields its neighbors, so a neighbor of one is out of reach
// unless it has taunt itself or the attacker ignores taunt.
func (g *Game) protectedByTaunt(attacker, target LocalID) bool {
	if g.attackIgnores(attacker, Taunt) || g.hasKeyword(target, Taunt) {
		return false
	}
	for _, neighbor := range neighbors(&EffectContext{Resolver: g}, target) {
		if g.hasKeyword(neighbor, Taunt) {
			return true
		}
	}
	return false
}

// onFlankOf reports whether a creature sits on a flank (the leftmost or rightmost
// creature) of its controller's battleline.
func (g *Game) onFlankOf(id LocalID) bool {
	bl := g.State.Battleline[g.controller(id)].slice()
	return len(bl) > 0 &&
		(bl[0] == id || bl[len(bl)-1] == id)
}

// fightDamage returns the damage an attacker deals to the defender it fights: its
// power, unless the card's AttackDamage overrides it (a Fixed amount) or adds a
// bonus (which may be limited to a defender on a flank).
func (g *Game) fightDamage(attacker, defender LocalID) int {
	ad := g.cat.def(attacker).AttackDamage
	if ad.Fixed {
		return ad.Amount
	}
	dmg := g.Power(attacker)
	if !ad.FlankOnly || g.onFlankOf(defender) {
		dmg += ad.Amount
	}
	return dmg
}

// applyRawDamage records damage on a creature, letting armor absorb it first. It
// only updates the damage counters; destruction is resolved by dealDamage, of
// which this is the per-creature step.
func (g *Game) applyRawDamage(id LocalID, amount int, ignoreArmor bool) {
	if amount <= 0 {
		return
	}
	core := &g.State.Cards[id]
	if core.DamageImmune {
		g.logf("%s cannot be dealt damage", g.Name(id))
		return
	}
	if !ignoreArmor {
		if absorbed := min(int(core.ArmorRemaining), amount); absorbed > 0 {
			core.ArmorRemaining -= int16(absorbed)
			amount -= absorbed
			g.logf("%s's armor absorbs %d damage", g.Name(id), absorbed)
		}
	}
	if amount > 0 {
		core.Damage += int16(amount)
		g.logf("%s takes %d damage (%d total)", g.Name(id), amount, core.Damage)
	}
}

// DamageTarget pairs a creature with the amount of damage to deal it within a
// single simultaneous batch (see dealDamage). It is exported because effects
// build damage batches and pass them through the Resolver.
type DamageTarget struct {
	ID     LocalID
	Amount int
	// IgnoreArmor makes this instance of damage bypass the creature's armor.
	IgnoreArmor bool
}

// dealDamage deals damage to every target simultaneously: each creature takes its
// damage (armor absorbs first) before any destruction is resolved, then all that
// were dealt lethal damage are destroyed together, in an order the controller
// chooses. Resolving destruction is part of dealing damage, so once dealDamage
// returns the dead creatures are already in the discard.
func (g *Game) dealDamage(controller int, targets ...DamageTarget) {
	for _, t := range targets {
		g.applyRawDamage(t.ID, t.Amount, t.IgnoreArmor)
	}
	var dying []LocalID
	for _, t := range targets {
		if g.shouldDestroy(t.ID) {
			dying = append(dying, t.ID)
		}
	}
	g.destroyEach(controller, dying)
}

// shouldDestroy reports whether a creature is currently in a destroyable state:
// its damage meets or exceeds its power, its power is zero or less, or it has
// Poison and any damage on it.
func (g *Game) shouldDestroy(id LocalID) bool {
	def := g.cat.def(id)
	if def.Type != Creature || !g.inPlay(id) {
		return false
	}
	core := &g.State.Cards[id]
	poisoned := g.hasKeyword(id, Poison) && core.Damage > 0
	return int(core.Damage) >= g.Power(id) ||
		g.Power(id) <= 0 ||
		poisoned
}
