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

// TestWardAmount covers the choose-N ward: the controller picks Amount distinct
// creatures from the target pool, and the effect renders the plural quantity.
func TestWardAmount(t *testing.T) {
	e := Ward{Target: Target{Kind: TargetEachFriendlyCreature}, Amount: 2}
	if got := e.Text(); got != "ward 2 friendly creatures" {
		t.Errorf("ward amount text = %q", got)
	}

	g := NewGame("A", "B", 1)
	a := g.AddToBattleline(testCreature("a", 3), 0)
	b := g.AddToBattleline(testCreature("b", 3), 0)
	c := g.AddToBattleline(testCreature("c", 3), 0)
	g.SetChooser(0, &idQueueChooser{ids: []LocalID{a, c}})
	e.Resolve(&EffectContext{Resolver: g, Source: a, Controller: 0})
	if !g.Warded(a) || !g.Warded(c) {
		t.Error("chosen creatures should be warded")
	}
	if g.Warded(b) {
		t.Error("unchosen creature should not be warded")
	}
}

// TestWardAmountRunsOut: when fewer creatures are available than Amount, the ward
// stops early once the target pool is exhausted.
func TestWardAmountRunsOut(t *testing.T) {
	g := NewGame("A", "B", 1)
	a := g.AddToBattleline(testCreature("a", 3), 0)
	Ward{Target: Target{Kind: TargetEachFriendlyCreature}, Amount: 2}.Resolve(
		&EffectContext{Resolver: g, Source: a, Controller: 0},
	)
	if !g.Warded(a) {
		t.Error("the only friendly creature should be warded")
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
