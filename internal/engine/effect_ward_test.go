package engine

import "testing"

func TestWardEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 3), 0)
	foe1 := g.AddToBattleline(testCreature("foe1", 3), 1)
	foe2 := g.AddToBattleline(testCreature("foe2", 3), 1)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	e := Ward{Target: Target{Kind: TargetEachEnemyCreature}}
	if e.Text() != "ward each enemy creature" {
		t.Errorf("ward text = %q", e.Text())
	}
	e.Resolve(ctx)
	if !g.Warded(foe1) || !g.Warded(foe2) {
		t.Error("ward should ward each enemy creature")
	}
	if g.Warded(src) {
		t.Error("warding enemy creatures should not touch a friendly creature")
	}

	// A ward that finds its target already warded still logs the choice, just
	// without a state change.
	entries := len(g.Log)
	e.Resolve(ctx)
	if len(g.Log) == entries {
		t.Error("re-warding an already-warded creature should still log the choice")
	}
}

func TestWardValidate(t *testing.T) {
	if (Ward{}).validate() == nil {
		t.Error("Ward with no target should fail validation")
	}
	if (Ward{Target: Target{Kind: TargetEachEnemyCreature}}).validate() != nil {
		t.Error("Ward with a target should validate")
	}
}

// TestWardAbsorbsDamage: a warded creature refuses the next instance of damage,
// takes none of it, and loses its ward — even a single point spends the whole
// ward.
func TestWardAbsorbsDamage(t *testing.T) {
	g := started(t)
	c := g.AddToBattleline(testCreature("c", 3), 0)
	g.SetWarded(c, true)

	g.DealDamage(0, []DamageTarget{{ID: c, Amount: 5}})

	if !g.InPlay(c) {
		t.Fatal("ward should absorb lethal damage, leaving the creature in play")
	}
	if g.Damage(c) != 0 {
		t.Errorf("damage = %d, want 0: ward absorbs the whole instance", g.Damage(c))
	}
	if g.Warded(c) {
		t.Error("absorbing damage should spend the ward")
	}
}

// TestWardAbsorbsDestruction: a warded creature survives a destroy effect, and
// its Destroyed ability never fires because it was never destroyed.
func TestWardAbsorbsDestruction(t *testing.T) {
	g := started(t)
	c := g.AddToBattleline(
		testCreature(
			"c",
			3,
			WithAbility(TriggerDestroyed, GainAember{Player: Controller, Amount: 1}),
		),
		0,
	)
	g.SetWarded(c, true)
	before := g.Aember(0)

	g.DestroyEach(0, []LocalID{c})

	if !g.InPlay(c) {
		t.Fatal("ward should absorb destruction, leaving the creature in play")
	}
	if g.Warded(c) {
		t.Error("absorbing destruction should spend the ward")
	}
	if g.Aember(0) != before {
		t.Error("a warded creature is never destroyed, so its Destroyed ability must not fire")
	}
}

// TestWardAbsorbsRelocation: every relocation out of play is absorbed by a ward —
// the creature stays and the ward is spent.
func TestWardAbsorbsRelocation(t *testing.T) {
	relocations := map[string]func(*Game, LocalID){
		"toHand":         func(g *Game, id LocalID) { g.putIntoHand(id) },
		"toArchives":     func(g *Game, id LocalID) { g.putIntoArchives(id) },
		"toTopOfDeck":    func(g *Game, id LocalID) { g.putOnTopOfDeck(id) },
		"shuffledInDeck": func(g *Game, id LocalID) { g.putIntoDeckShuffled(id) },
		"toYourArchives": func(g *Game, id LocalID) { g.PutIntoYourArchives(id, 1) },
		"purge":          func(g *Game, id LocalID) { g.purgeFromPlay(id) },
		"graftUnder":     func(g *Game, id LocalID) { g.GraftUnder(id, id) },
	}
	for name, remove := range relocations {
		t.Run(name, func(t *testing.T) {
			g := started(t)
			c := g.AddToBattleline(testCreature("c", 3), 0)
			g.SetWarded(c, true)

			remove(g, c)

			if !g.InPlay(c) {
				t.Fatalf("%s: ward should absorb the removal, leaving the creature in play", name)
			}
			if g.Warded(c) {
				t.Errorf("%s: absorbing the removal should spend the ward", name)
			}
		})
	}
}

// TestUnwardedCreatureStillLeavesPlay guards the fall-through: an unwarded
// creature is removed as usual.
func TestUnwardedCreatureStillLeavesPlay(t *testing.T) {
	g := started(t)
	c := g.AddToBattleline(testCreature("c", 3), 0)

	g.putIntoHand(c)

	if g.InPlay(c) {
		t.Error("an unwarded creature should leave play as usual")
	}
}

// TestAbsorbedByWardOnCardOutOfPlay covers the nil-state branch: a card no longer
// in play has no ward to spend.
func TestAbsorbedByWardOnCardOutOfPlay(t *testing.T) {
	g := started(t)
	gone := g.Register(testCreature("gone", 3), 0)
	g.State.Discard[0].add(gone)

	if g.absorbedByWard(gone) {
		t.Error("a card out of play has no ward to absorb a removal")
	}
}
