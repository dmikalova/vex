package engine

// This file holds combat: an attacker fighting a defender, dealing damage (with
// armor, Skirmish, Assault, Hazardous, and Splash-attack), and deciding which
// creatures the damage destroys. The destruction itself is carried out in
// game_destroy.go.

// Combat: use one of your ready creatures to fight an enemy creature. Using it to
// fight exhausts it. First, any "Before Fight" abilities and the Assault and
// Hazardous keywords resolve; if these destroy either creature, the fight does not
// occur. Otherwise both creatures deal damage equal to their power at the same
// time — armor reduces the damage a creature takes, and a Skirmish attacker takes
// no damage back. A creature with damage equal to or greater than its power is
// destroyed, and both deaths are resolved together, so neither fighter's
// destruction changes the damage the other deals. Finally, if the attacker
// survived, its "Fight" abilities resolve.
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
	// While the fight resolves, both combatants count as "fighting", which a
	// creature's "while fighting" self-grant reads (Nizak, The Forgotten gains
	// invulnerable). Restored after the fight so a nested fight sees only its own
	// pair and the grant lifts once combat ends.
	prevFighters := g.State.FightersPlus
	g.State.FightersPlus = [2]LocalID{attacker + 1, defender + 1}
	defer func() { g.State.FightersPlus = prevFighters }()
	g.State.TurnHistory[attackerSide][CreaturesFoughtThisTurn]++
	// Using a creature to fight exhausts it before anything else resolves, so a
	// "Before Fight" ability already sees the attacker exhausted. The defender is
	// put in context, so an ability can act on "the creature this fights".
	g.State.Cards[attacker].Exhausted = true
	g.triggerAbilities(attacker, TriggerBeforeFight, defender, true)
	g.fireLastingBeforeFight(attacker)

	// A "Before Fight" ability may redirect the attacker's fight damage to another
	// creature (Gabos Longarms) or make the fight not occur (Evasion Sigil). Read
	// and clear those per-fight flags before the combat step.
	redirect := g.State.FightDamageRedirect
	g.State.FightDamageRedirect = 0
	cancelled := g.State.FightCancelled
	g.State.FightCancelled = false
	if cancelled {
		g.record(FightCancelled{Attacker: attacker})
		return
	}

	// Snapshot the attacker's battleline neighbors before combat, so an "after a
	// neighbor is used to fight" reaction (Little Niff) fires for the creatures
	// beside the attacker even if the fight destroys the attacker and shifts the
	// battleline.
	neighborsAtFight := neighbors(&EffectContext{Resolver: g}, attacker)

	// Assault and Hazardous deal their damage before fight damage: the attacker's
	// Assault hits the defender, the defender's Hazardous hits the attacker.
	// Either can destroy a fighter before combat. Skipped if a "Before Fight"
	// effect already removed the defender.
	if g.inPlay(defender) {
		var pre []DamageTarget
		if a := g.assault(attacker); a > 0 {
			pre = append(pre, DamageTarget{
				ID:            defender,
				Amount:        a,
				Source:        attacker,
				SourceKeyword: assaultDamage,
			})
		}
		if h := g.hazardous(defender); h > 0 {
			pre = append(pre, DamageTarget{
				ID:            attacker,
				Amount:        h,
				Source:        defender,
				SourceKeyword: hazardousDamage,
			})
		}
		if len(pre) > 0 {
			g.dealDamage(g.controller(attacker), pre...)
			// Skoll's Assault destroying the creature it attacks — before the fight
			// itself — fires "after a creature is destroyed by this creature's assault
			// damage" on the attacker, with the destroyed defender as "it". Assault is
			// the only pre-fight damage aimed at the defender, so a defender no longer
			// in play here was destroyed by it.
			if a := g.assault(attacker); a > 0 && g.inPlay(attacker) && !g.inPlay(defender) {
				g.triggerAbilities(attacker, TriggerAfterAssaultDestroys, defender, true)
			}
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
		g.record(Fought{
			Attacker:      attacker,
			AttackerPower: ap,
			Defender:      defender,
			DefenderPower: dp,
		})
		if elusive {
			g.record(ElusiveAvoidedFight{Defender: defender})
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
			// Retaliation damage equals the defender's power, unless a Fixed
			// AttackDamage replaces it: Shadow Self and Ether Spider "deal no damage
			// when fighting", so they deal none back to an attacker either. An
			// additive AttackDamage bonus (Valdr's flank +2) is an attack-only bonus
			// and never adds to retaliation.
			retaliation := dp
			if ad := g.cat.def(defender).AttackDamage; ad.Fixed {
				retaliation = ad.Amount
			}
			skirmish := g.hasKeyword(attacker, Skirmish)
			if !skirmish &&
				!g.cat.def(defender).DealsNoDamageWhenAttacked {
				targets = append(targets, DamageTarget{ID: attacker, Amount: retaliation})
			} else if skirmish && retaliation > 0 {
				g.record(SkirmishAvoidedReturn{Attacker: attacker})
			}
			// Splash-attack deals its damage to each neighbor of the creature the
			// attacker fights, at the same time as fight damage.
			if s := g.splashAttack(attacker); s > 0 {
				for _, n := range neighbors(&EffectContext{Resolver: g}, defender) {
					targets = append(targets, DamageTarget{ID: n, Amount: s})
				}
			}
			// Snapshot each target's damage so poison can tell which creatures the
			// fight's damage actually landed on, once armor and immunity had their say.
			before := make([]int16, len(targets))
			for i, t := range targets {
				before[i] = g.State.Cards[t.ID].Damage
			}
			g.dealDamage(g.controller(attacker), targets...)
			g.applyFightPoison(attacker, defender, dmgTarget, targets, before)

		}
	}
	// The whole post-combat reaction set resolves as one ordered window (ADR 0013):
	// the attacker's own "Fight:", the after-destroyed-in-a-fight reactions, the
	// after-use reactions, the board-wide "after a creature fights", and the
	// neighbor reactions are gathered before any resolves so the active player orders
	// them together. Tallies are read first because they are bookkeeping, not
	// reactions, and a later reaction must see them settled.
	attackerDead, defenderDead := !g.inPlay(attacker), !g.inPlay(defender)
	// A creature destroyed in a fight is an enemy kill from the other side's point of
	// view, which is the tally The Warchest is paid for.
	if attackerDead {
		g.State.TurnHistory[1-attackerSide][EnemyCreaturesFightKilled]++
	}
	if defenderDead {
		g.State.TurnHistory[1-defenderSide][EnemyCreaturesFightKilled]++
	}
	g.resolveWindow(g.orderTriggered(
		attackerSide,
		append(
			g.fightReactions(attacker, defender, attackerSide, defenderSide, neighborsAtFight),
			g.lastingReactions(EventFight, attackerSide, attacker)...,
		),
	))
	// Foggify's stun-fighter bar, armed against this player for this turn, stuns
	// each creature they use to fight — read after the fight window so it lands on
	// an attacker that survived combat and its own "Fight:" abilities.
	if g.State.StunFighter[attackerSide].Value && g.inPlay(attacker) {
		g.SetStunned(attacker, true)
	}
}

// fightReactions gathers, in default resolution order, every card-sourced ability
// that reacts to a fight that has just resolved: the attacker's own "Fight:"
// (bound to a surviving defender as "it"), the "after a creature is destroyed in a
// fight" reactions on the survivor and on the enemies of a slain fighter, the
// attacker's "after I am used" and the "after you use a card" reactions on its
// controller's other cards, the board-wide "after a creature fights", and the
// "after a neighbor fights" reactions on the attacker's flankmates. Every entry
// carries its own actor and "it" so the whole set orders as one window while each
// bystander reaction resolves for its owner.
func (g *Game) fightReactions(
	attacker, defender LocalID,
	attackerSide, defenderSide int,
	neighborsAtFight []LocalID,
) []triggeredAbility {
	var pending []triggeredAbility
	add := func(src LocalID, trigger Trigger, it LocalID, hasIt bool) {
		got := g.triggeredBy(src, trigger)
		for i := range got {
			got[i].actor = int8(g.controller(src))
			got[i].it, got[i].hasIt = it, hasIt
		}
		pending = append(pending, got...)
	}
	attackerAlive, defenderAlive := g.inPlay(attacker), g.inPlay(defender)
	// A creature the fight destroyed has left play, so its "After Fight:" ability
	// does not resolve. A defender that survived is passed as "it" so an after-fight
	// ability can name "the creature <self> fights" (Roxador stuns it).
	if attackerAlive {
		if defenderAlive {
			add(attacker, TriggerAfterFight, defender, true)
		} else {
			add(attacker, TriggerAfterFight, 0, false)
		}
	}
	// "After a creature is destroyed in a fight with X": when exactly one combatant
	// is removed by the fight, the survivor's ability fires with the destroyed
	// creature as "it".
	if !defenderAlive && attackerAlive {
		add(attacker, TriggerAfterDestroyedFighting, defender, true)
	}
	if !attackerAlive && defenderAlive {
		add(defender, TriggerAfterDestroyedFighting, attacker, true)
	}
	// "After an enemy creature is destroyed while fighting": a bystander (The
	// Colosseum) whose controller is the enemy of a combatant killed in the fight
	// reacts, with the destroyed creature as "it".
	if !defenderAlive {
		for _, id := range g.allInPlay(1 - defenderSide) {
			add(id, TriggerAfterEnemyDestroyedFighting, defender, true)
		}
	}
	if !attackerAlive {
		for _, id := range g.allInPlay(1 - attackerSide) {
			add(id, TriggerAfterEnemyDestroyedFighting, attacker, true)
		}
	}
	// The used card's own "after I am used", then "after you use a card" on the
	// attacking player's other cards, with the attacker as "it".
	add(attacker, TriggerAfterUsedSelf, 0, false)
	for _, id := range g.allInPlay(attackerSide) {
		if id != attacker {
			add(id, TriggerAfterUse, attacker, true)
		}
	}
	// "After a creature is used to fight" on every in-play card (Shattered Throne,
	// Peace Accord), with the fighting creature as "it".
	for player := 0; player < 2; player++ {
		for _, id := range g.allInPlay(player) {
			add(id, TriggerAfterCreatureFights, attacker, true)
		}
	}
	// "After a neighbor of this is used to fight" on the creatures that flanked the
	// attacker when the fight began (Little Niff), with the attacker as "it". A
	// neighbor destroyed by the fight is skipped.
	for _, neighbor := range neighborsAtFight {
		if g.inPlay(neighbor) {
			add(neighbor, TriggerAfterNeighborFights, attacker, true)
		}
	}
	return pending
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
	// The pre-fight sequence may already have removed the defender, and a card out
	// of play has had its per-match state zeroed — marking it would strand the flag
	// there and make the creature arrive back in play with its elusive spent.
	if g.inPlay(defender) {
		g.State.Cards[defender].ElusiveUsedThisTurn = true
	}
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

// attackGrantsPoison reports whether an attacker gains poison for a fight against
// this defender — Spyyyder gains poison only when the defender is on a flank.
func (g *Game) attackGrantsPoison(attacker, defender LocalID) bool {
	ak := g.cat.def(attacker).AttackKeywords
	if ak.FlankOnly && !g.onFlankOf(defender) {
		return false
	}
	for _, kw := range ak.Keywords {
		if kw == Poison {
			return true
		}
	}
	return false
}

// applyFightPoison destroys each creature a poison combatant dealt landing damage
// to in the fight just resolved. Poison is offensive: the attacker poisons the
// defender it fought — from its base Poison keyword, or one it gains for the fight
// (Spyyyder) — and any creature its splash hit, while a poison defender poisons the
// attacker with its return damage. A creature that merely has poison is not
// destroyed by taking damage; only damage above its power destroys it. A target is
// poisoned only when the fight's damage actually landed on it (its total rose above
// the snapshot), so armor or immunity that soaked the whole hit spares it.
func (g *Game) applyFightPoison(
	attacker, defender, dmgTarget LocalID, targets []DamageTarget, before []int16,
) {
	attackerPoison := g.hasKeyword(attacker, Poison)
	grantsPoison := g.attackGrantsPoison(attacker, defender)
	defenderPoison := g.hasKeyword(defender, Poison)
	var poisoned []LocalID
	for i, t := range targets {
		if !g.inPlay(t.ID) || g.State.Cards[t.ID].Damage <= before[i] {
			continue
		}
		var source LocalID
		switch t.ID {
		case attacker:
			if !defenderPoison {
				continue
			}
			source = defender
		default: // a creature the attacker dealt fight or splash damage to
			if !attackerPoison && (t.ID != dmgTarget || !grantsPoison) {
				continue
			}
			source = attacker
		}
		g.record(PoisonKills{Source: source, Victim: t.ID})
		poisoned = append(poisoned, t.ID)
	}
	g.destroyEach(g.controller(attacker), poisoned)
}

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

// TauntShielded reports whether a creature is shielded by a neighboring
// taunter — sitting beside a creature with taunt, without having taunt itself.
// It answers with no attacker in mind (unlike protectedByTaunt, which a fight
// asks per attacker), so a client can show it as standing appearance rather than
// a fact that only holds for the current fight.
func (g *Game) TauntShielded(id LocalID) bool {
	if g.hasKeyword(id, Taunt) {
		return false
	}
	for _, neighbor := range neighbors(&EffectContext{Resolver: g}, id) {
		if g.hasKeyword(neighbor, Taunt) {
			return true
		}
	}
	return false
}

// onFlankOf reports whether a creature sits on a flank (the leftmost or rightmost
// creature) of its controller's battleline.
func (g *Game) onFlankOf(id LocalID) bool {
	if g.State.Cards[id].ConsideredFlank {
		return true
	}
	bl := g.State.Battleline[g.controller(id)].slice()
	return len(bl) > 0 &&
		(bl[0] == id || bl[len(bl)-1] == id)
}

// protectedFromNonFlank reports whether a creature carries a "creatures not on a
// flank cannot fight this creature" restriction from an attached upgrade
// (Camouflage). A blanked upgrade grants nothing (staticOn honors blanking).
func (g *Game) protectedFromNonFlank(id LocalID) bool {
	for up, ok := g.firstUpgrade(id); ok; up, ok = g.nextUpgrade(up) {
		if g.staticOn(id, up).ProtectsFromNonFlank {
			return true
		}
	}
	return false
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

// applyRawDamage deals damage to a creature and returns the creature that ended
// up marked with it. The damage runs that creature's defenses first — invulnerability
// refuses all of it, then armor absorbs from the front — and only what survives is
// dealt. Redirection is the last step of that chart, not the first: a shield takes
// the damage the creature was actually dealt, and then runs its own defenses over
// it in turn, so a shield's armor absorbs from an already-armored amount. It only
// updates the damage counters; destruction is resolved by dealDamage, of which this
// is the per-creature step.
func (g *Game) applyRawDamage(t DamageTarget) LocalID {
	id, amount, ignoreArmor := t.ID, t.Amount, t.IgnoreArmor
	if amount <= 0 {
		return id
	}
	// A creature an earlier step already removed can still be named by a batch or by
	// a trigger's "it", and damage dealt to it now would sit on a card in hand.
	if !g.inPlay(id) {
		return id
	}
	if amount = g.mitigateDamage(id, amount, ignoreArmor); amount <= 0 {
		return id
	}
	if shield := g.damageRedirect(id); shield != id {
		id = shield
		if amount = g.mitigateDamage(id, amount, ignoreArmor); amount <= 0 {
			return id
		}
	}
	core := &g.State.Cards[id]
	core.Damage += int16(amount)
	g.record(t.damageEntry(id, amount, int(core.Damage)))
	// A creature carrying a Lethal Distraction takes additional damage on top of
	// each instance it takes. The bonus is added once here, so it does not itself
	// re-trigger the effect on the same instance.
	if bonus := g.lastingExtraDamage(id); bonus > 0 {
		core.Damage += int16(bonus)
		g.record(DamageTaken{Creature: id, Amount: bonus, Total: int(core.Damage)})
	}
	return id
}

// mitigateDamage runs a creature's defenses over incoming damage and returns what
// is left for it to be dealt.
func (g *Game) mitigateDamage(id LocalID, amount int, ignoreArmor bool) int {
	core := &g.State.Cards[id]
	if core.DamageImmune || g.State.SideDamageImmune[g.controller(id)] ||
		g.hasKeyword(id, Invulnerable) {
		g.record(DamageRefused{Creature: id})
		return 0
	}
	if g.absorbedByWard(id, wardDamage, amount) {
		return 0
	}
	if !ignoreArmor {
		if absorbed := min(int(core.ArmorRemaining), amount); absorbed > 0 {
			core.ArmorRemaining -= int16(absorbed)
			amount -= absorbed
			g.record(ArmorAbsorbed{Creature: id, Amount: absorbed})
		}
	}
	return amount
}

// DamageTarget pairs a creature with the amount of damage to deal it within a
// single simultaneous batch (see dealDamage). It is exported because effects
// build damage batches and pass them through the Resolver.
type DamageTarget struct {
	// ID is the creature; Amount is the damage dealt to it in the batch.
	ID     LocalID
	Amount int
	// IgnoreArmor makes this instance of damage bypass the creature's armor.
	IgnoreArmor bool
	// Source, when SourceKeyword is set, credits this damage to the card that dealt
	// it — a creature's pre-fight Assault or Hazardous, or the card whose ability
	// dealt it — so the hit narrates "<Source> ... deals N damage to <ID>" instead
	// of the bare "<ID> takes N damage" line, naming where the damage came from.
	Source        LocalID
	SourceKeyword combatKeyword
}

// combatKeyword names the source crediting a DamageTarget's damage, so the hit
// narrates the striking creature — its pre-fight keyword and value, or the card
// whose ability dealt it — rather than a bare "takes N damage" line. The zero
// value credits no source: the hit narrates as the plain DamageTaken line.
type combatKeyword int

const (
	// assaultDamage credits the attacker's Assault striking the creature it attacks.
	assaultDamage combatKeyword = iota + 1
	// hazardousDamage credits the defender's Hazardous striking its attacker.
	hazardousDamage
	// abilityDamage credits the card whose ability dealt the damage (Musthic
	// Murmook deals 4 damage to a creature).
	abilityDamage
)

// dealDamage deals damage to every target simultaneously: each creature takes its
// damage (armor absorbs first) before any destruction is resolved, then all that
// were dealt lethal damage are destroyed together, in an order the controller
// chooses. Resolving destruction is part of dealing damage, so once dealDamage
// returns the dead creatures are already in the discard.
func (g *Game) dealDamage(controller int, targets ...DamageTarget) {
	// Snapshot the armor of any creature that watches for its own armor prevention
	// (Maruck the Marked) so the batch can report how much each just prevented. The
	// snapshot is skipped entirely when no such creature is in play, so the common
	// case adds nothing to the hot path.
	watchers := g.armorPreventWatchers()
	var armorBefore map[LocalID]int
	if len(watchers) > 0 {
		armorBefore = make(map[LocalID]int, len(watchers))
		for _, id := range watchers {
			armorBefore[id] = int(g.State.Cards[id].ArmorRemaining)
		}
	}
	// A redirect moves the damage to a shield, so the creature to test for
	// destruction is whichever one applyRawDamage ended up marking.
	hit := make([]LocalID, len(targets))
	for i, t := range targets {
		hit[i] = g.applyRawDamage(t)
	}
	var dying []LocalID
	for _, id := range hit {
		if g.shouldDestroy(id) {
			dying = append(dying, id)
		}
	}
	g.destroyEach(controller, dying)
	if len(watchers) > 0 {
		g.emitArmorPrevented(watchers, armorBefore)
	}
}

// armorPreventWatchers lists the in-play creatures carrying an After This Creature
// Prevents Damage With Its Armor ability, so dealDamage can skip its bookkeeping
// entirely when none are present.
func (g *Game) armorPreventWatchers() []LocalID {
	var watchers []LocalID
	for player := 0; player < 2; player++ {
		for _, id := range g.allInPlay(player) {
			if len(g.triggeredBy(id, TriggerAfterArmorPrevents)) > 0 {
				watchers = append(watchers, id)
			}
		}
	}
	return watchers
}

// emitArmorPrevented fires each watcher's After This Creature Prevents Damage With
// Its Armor ability, scaled by the armor it just spent absorbing damage (before
// minus what remains). A watcher that spent no armor, or left play in the batch,
// does not fire.
func (g *Game) emitArmorPrevented(watchers []LocalID, armorBefore map[LocalID]int) {
	for _, id := range watchers {
		prevented := armorBefore[id] - int(g.State.Cards[id].ArmorRemaining)
		if prevented <= 0 || !g.inPlay(id) {
			continue
		}
		actor := g.controller(id)
		for _, t := range g.orderTriggered(
			actor, g.triggeredBy(id, TriggerAfterArmorPrevents),
		) {
			closeFrame := g.openFrame(Frame{
				Actor:      actor,
				Source:     id,
				HasSource:  true,
				Trigger:    TriggerAfterArmorPrevents,
				Grantor:    t.grantor,
				HasGrantor: t.grantor != id,
			})
			ctx := &EffectContext{Resolver: g, Source: id, Controller: actor}
			ctx.Produced.ArmorPrevented = prevented
			t.ability.Effect.Resolve(ctx)
			closeFrame()
		}
	}
}

// damageRedirect returns the creature that takes the damage a creature was dealt:
// normally id itself, but a card in play whose TakesDamageFor covers id takes it
// instead (Shadow Self shields its non-Specter neighbors). A redirect never
// chains — the shield's own damage is never redirected again — so two shields
// cannot bounce damage between them.
func (g *Game) damageRedirect(id LocalID) LocalID {
	for player := range 2 {
		for _, shield := range g.allInPlay(player) {
			t := g.cat.def(shield).TakesDamageFor
			if shield == id || !t.valid() {
				continue
			}
			for _, warded := range t.Select(g.constantContext(shield)) {
				if warded == id {
					return shield
				}
			}
		}
	}
	return id
}

// shouldDestroy reports whether a creature is currently in a destroyable state:
// its damage meets or exceeds its power, its power is zero or less, or its own
// text names a board state that destroys it (Tireless Crocag while the opponent
// has no creatures). Poison does not make a creature destroyable — it destroys
// the creatures a poison creature damages in a fight, resolved in applyFightPoison.
func (g *Game) shouldDestroy(id LocalID) bool {
	def := g.cat.def(id)
	if g.TypeOf(id) != Creature || !g.inPlay(id) {
		return false
	}
	if dw := def.DestroyedWhen; dw != nil {
		ctx := &EffectContext{Resolver: g, Source: id, Controller: g.controller(id)}
		if dw.Met(ctx) {
			return true
		}
	}
	core := &g.State.Cards[id]
	return int(core.Damage) >= g.Power(id) ||
		g.Power(id) <= 0
}
