package engine

import "testing"

func TestOnChooseCreatureEnemyAndNoTarget(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	enemy := g.AddToBattleline(testCreature("enemy", 1), 1)
	g.State.Cards[enemy].Exhausted = true
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	onEnemy := OnChooseCreature{
		Target: Target{Kind: TargetChosenEnemyCreature},
		Verbs:  []CreatureVerb{ReadyVerb{}},
	}
	if onEnemy.Text() != "ready an enemy creature" {
		t.Errorf("text = %q", onEnemy.Text())
	}
	onEnemy.Resolve(ctx)
	if g.State.Cards[enemy].Exhausted {
		t.Error("enemy should have been readied")
	}

	// No candidates: remove the enemy and resolve again (logs, no panic).
	g.DestroyEach(0, []LocalID{enemy})
	onEnemy.Resolve(ctx)
}

func TestFightVerbNoEnemy(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 2), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	FightVerb{}.Apply(ctx, src) // no enemies -> logs and returns
	if g.State.Cards[src].Exhausted {
		t.Error("no fight should have occurred")
	}
}

// A taunter shields its neighbors from an ability-driven fight too, so a forced
// fight (Sergeant Zakiel) cannot reach past a taunter to the creature beside it.
func TestFightVerbRespectsTaunt(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 3), 0)
	taunter := g.AddToBattleline(testCreature("taunter", 5, WithKeywords(Taunt)), 1)
	shielded := g.AddToBattleline(testCreature("shielded", 3), 1)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	// Aiming at the shielded neighbor: it is filtered out, so no fight lands.
	g.SetChooser(0, idChooser{id: shielded})
	FightVerb{}.Apply(ctx, src)
	if g.Damage(shielded) != 0 {
		t.Errorf("shielded neighbor took %d damage; taunt should shield it", g.Damage(shielded))
	}

	// The taunter itself stays reachable.
	g.SetChooser(0, idChooser{id: taunter})
	FightVerb{}.Apply(ctx, src)
	if g.Damage(taunter) != 3 {
		t.Errorf("taunter took %d damage, want 3", g.Damage(taunter))
	}
}

// actor builds a creature with an "Action:" ability that gains 5 Æmber.
func actor(g *Game) LocalID {
	return g.AddToBattleline(NewCard("actor", Brobnar, Creature, Common, WithPower(3),
		WithAbility(TriggerAction, GainAember{Player: Controller, Amount: 5})), 0)
}

func TestUseVerb(t *testing.T) {
	ctxFor := func(g *Game, id LocalID) *EffectContext {
		return &EffectContext{Resolver: g, Source: id, Controller: 0}
	}

	// Reap (default option 0): +1 Æmber, no action fired.
	g := NewGame("A", "B", 1)
	target := actor(g)
	UseVerb{}.Apply(ctxFor(g, target), target)
	if g.Aember(0) != 1 {
		t.Errorf("reap use: aember = %d, want 1", g.Aember(0))
	}

	// Fight (option 1, present because there is an enemy).
	g = NewGame("A", "B", 1)
	target = actor(g)
	foe := g.AddToBattleline(testCreature("foe", 5), 1)
	g.SetChooser(0, optionPicker{idx: 1})
	UseVerb{}.Apply(ctxFor(g, target), target)
	if g.Damage(foe) != 3 {
		t.Errorf("fight use: foe damage = %d, want 3", g.Damage(foe))
	}

	// Use its action (option 2: reap, fight, use its action).
	g = NewGame("A", "B", 1)
	target = actor(g)
	g.AddToBattleline(testCreature("foe", 5), 1)
	g.SetChooser(0, optionPicker{idx: 2})
	UseVerb{}.Apply(ctxFor(g, target), target)
	if g.Aember(0) != 5 {
		t.Errorf("action use: aember = %d, want 5", g.Aember(0))
	}

	// Out-of-range option: nothing happens.
	g = NewGame("A", "B", 1)
	target = actor(g)
	g.SetChooser(0, optionPicker{idx: 99})
	UseVerb{}.Apply(ctxFor(g, target), target)
	if g.Aember(0) != 0 || g.Exhausted(target) {
		t.Error("out-of-range use should do nothing")
	}
}

// A creature can only be used while ready. An exhausted creature may still be
// chosen to be used, but reaping, fighting, or using its action does nothing.
func TestUsingExhaustedCreatureDoesNothing(t *testing.T) {
	g := NewGame("A", "B", 1)
	c := actor(g) // has an "Action:" that would gain 5 Æmber
	foe := g.AddToBattleline(testCreature("foe", 5), 1)
	g.State.Cards[c].Exhausted = true

	g.ReapWith(c)
	g.FightWith(c, foe)
	g.UseActionOf(0, c)

	if g.Aember(0) != 0 {
		t.Errorf("aember = %d, want 0 (an exhausted creature cannot be used)", g.Aember(0))
	}
	if g.Damage(foe) != 0 {
		t.Errorf("foe damage = %d, want 0 (an exhausted creature cannot fight)", g.Damage(foe))
	}
}

// "Use a friendly creature" offers only ready creatures: with every friendly
// creature exhausted there is no decision to make, so nothing is asked.
func TestUseVerbOffersOnlyReadyCreatures(t *testing.T) {
	e := OnChooseCreature{
		Target: Target{Kind: TargetChosenFriendlyCreature},
		Verbs:  []CreatureVerb{UseVerb{}},
	}

	g := NewGame("A", "B", 1)
	spent := g.AddToBattleline(testCreature("spent", 3), 0)
	ready := g.AddToBattleline(testCreature("ready", 3), 0)
	g.State.Cards[spent].Exhausted = true
	src := g.AddToBattleline(testCreature("src", 3), 0)
	g.SetChooser(0, FirstChooser{}) // takes the first candidate offered
	e.Resolve(&EffectContext{Resolver: g, Source: src, Controller: 0})
	if !g.Exhausted(ready) || g.Aember(0) != 1 {
		t.Errorf("the ready creature should have been the one used (aember = %d)", g.Aember(0))
	}

	// All exhausted: nothing to use, and nothing is asked.
	g2 := NewGame("A", "B", 1)
	only := g2.AddToBattleline(testCreature("only", 3), 0)
	g2.State.Cards[only].Exhausted = true
	src2 := g2.AddToBattleline(testCreature("src", 3), 0)
	g2.State.Cards[src2].Exhausted = true
	asked := &promptRecorder{}
	g2.SetChooser(0, asked)
	e.Resolve(&EffectContext{Resolver: g2, Source: src2, Controller: 0})
	if asked.asked != 0 || g2.Aember(0) != 0 {
		t.Errorf("nothing to use should ask nothing, asked %d times", asked.asked)
	}
}

// promptRecorder counts every question put to the player.
type promptRecorder struct {
	FirstChooser
	asked int
}

func (p *promptRecorder) ChooseCreature(a, b string, cands []LocalID) (LocalID, bool) {
	p.asked++
	return p.FirstChooser.ChooseCreature(a, b, cands)
}

func (p *promptRecorder) ChooseOption(_, _ string, _ []string) int {
	p.asked++
	return 0
}

func TestOnChooseCreatureExcludeHouse(t *testing.T) {
	g := NewGame("A", "B", 1)
	sanc := g.AddToBattleline(NewCard("s", Sanctum, Creature, Common, WithPower(3)), 0)
	mars := g.AddToBattleline(NewCard("m", Mars, Creature, Common, WithPower(3)), 0)
	g.State.Cards[mars].Exhausted = true
	ctx := &EffectContext{Resolver: g, Source: sanc, Controller: 0}

	e := OnChooseCreature{
		Target: Target{Kind: TargetChosenFriendlyCreature}.ExceptHouse(Sanctum),
		Verbs:  []CreatureVerb{ReadyVerb{}},
	}
	if e.Text() != "ready a friendly non-Sanctum creature" {
		t.Errorf("text = %q", e.Text())
	}
	// The Sanctum creature is excluded, so the readied creature must be the Mars one.
	e.Resolve(ctx)
	if g.Exhausted(mars) {
		t.Error("the non-Sanctum creature should have been readied")
	}
}

// idChooser picks a specific creature by id (and no "choose one" option).
type idChooser struct {
	FirstChooser
	id LocalID
}

func (c idChooser) ChooseCreature(_, _ string, cands []LocalID) (LocalID, bool) {
	for _, x := range cands {
		if x == c.id {
			return x, true
		}
	}
	return 0, false
}

func TestUseVerbNesting(t *testing.T) {
	g := NewGame("A", "B", 1)
	// A's "Reap:" uses a friendly creature; B is the one it uses (to reap).
	a := g.AddToBattleline(NewCard(
		"A",
		Brobnar,
		Creature,
		Common,
		WithPower(3),
		WithAbility(
			TriggerAfterReap,
			OnChooseCreature{
				Target: Target{Kind: TargetChosenFriendlyCreature},
				Verbs:  []CreatureVerb{UseVerb{}},
			},
		),
	), 0)
	b := g.AddToBattleline(testCreature("B", 3), 0)
	g.SetChooser(0, idChooser{id: b}) // use B, not A itself

	// Reaping A resolves A's window, which nests B's reap fully before returning.
	g.reapWith(a)
	if g.Aember(0) != 2 { // A's reap (+1) and B's nested reap (+1)
		t.Errorf("aember = %d, want 2 (A's reap + B's nested reap)", g.Aember(0))
	}
}

// gameEffect is a test-only effect that runs an arbitrary closure over the game.
type gameEffect struct{ fn func() }

func (gameEffect) Text() string             { return "test" }
func (e gameEffect) Resolve(*EffectContext) { e.fn() }

func TestTriggerWindowIsFixedWhenTheEventHappens(t *testing.T) {
	g := NewGame("A", "B", 1)
	var a LocalID
	// A's printed Reap ability attaches an upgrade to A that grants a NEW
	// "Reap: gain 3 Æmber". That ability did not trigger on this reap: a window's
	// abilities are collected when the event happens, before any of them resolves,
	// which is what lets the active player order the whole window (ADR 0013).
	attach := gameEffect{fn: func() {
		up := g.Register(NewCard(
			"boost",
			Brobnar,
			Upgrade,
			Common,
			WithStatic(
				StaticModifier{
					Granted: []Ability{
						{
							Trigger: TriggerAfterReap,
							Effect:  GainAember{Player: Controller, Amount: 3},
						},
					},
				},
			),
		), 0)
		g.AttachUpgrade(a, up)
	}}
	a = g.AddToBattleline(NewCard("A", Brobnar, Creature, Common, WithPower(2),
		WithAbility(TriggerAfterReap, attach)), 0)

	g.reapWith(a)
	if g.Aember(0) != 1 { // the reap itself; the newly granted ability waits for the next one
		t.Errorf("aember = %d, want 1 (reap only)", g.Aember(0))
	}

	// It does fire on the next reap, once it is in play when the event happens.
	g.State.Cards[a].Exhausted = false
	g.reapWith(a)
	if g.Aember(0) != 5 { // 1 + (1 reap + 3 granted)
		t.Errorf("aember = %d, want 5 (granted ability fires on the next reap)", g.Aember(0))
	}
}

func TestOnChooseCreatureNeighbors(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.AddToBattleline(testCreature("far", 5), 0) // NOT a neighbor of src
	g.AddToBattleline(testCreature("left", 3), 0)
	src := g.AddToBattleline(testCreature("src", 3), 0)
	g.AddToBattleline(testCreature("right", 3), 0)
	foe := g.AddToBattleline(testCreature("foe", 10), 1)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	e := OnChooseCreature{
		Target: Target{Kind: TargetChosenCreature}.Neighboring(),
		Verbs:  []CreatureVerb{ReadyVerb{}, FightVerb{}},
	}
	if e.Text() != "ready and fight with a neighboring creature" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	// The chosen neighbor is `left` (the first neighbor), not the non-neighbor
	// `far`: the foe takes left's 3 power, not far's 5.
	if g.Damage(foe) != 3 {
		t.Errorf("foe damage = %d, want 3 (a neighbor fought, not the far creature)", g.Damage(foe))
	}
}

func TestStunExhaustVerbs(t *testing.T) {
	if got := (StunVerb{}).VerbText(); got != "stun" {
		t.Errorf("StunVerb text = %q", got)
	}
	if got := (ExhaustVerb{}).VerbText(); got != "exhaust" {
		t.Errorf("ExhaustVerb text = %q", got)
	}
	if got := (UseVerb{}).VerbText(); got != "use" {
		t.Errorf("UseVerb text = %q", got)
	}
	g := started(t)
	id := g.AddToBattleline(testCreature("c", 3), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	StunVerb{}.Apply(ctx, id)
	ExhaustVerb{}.Apply(ctx, id)
	if !g.State.Cards[id].Stunned || !g.State.Cards[id].Exhausted {
		t.Errorf(
			"verbs did not apply: stunned=%v exhausted=%v",
			g.State.Cards[id].Stunned,
			g.State.Cards[id].Exhausted,
		)
	}
}

// An earlier sentence can destroy the very creature a later one names —
// Transposition Sandals swaps a creature off the flank that was keeping it alive,
// then says to use it — so a creature that has left play takes no verbs.
func TestOnChooseCreatureSkipsCreatureThatLeftPlay(t *testing.T) {
	g := started(t)
	gone := g.Register(testCreature("gone", 3), 0)
	g.State.Discard[0].add(gone)

	e := OnChooseCreature{Verbs: []CreatureVerb{ReadyVerb{}, StunVerb{}}}
	if e.applyTo(&EffectContext{Resolver: g, Controller: 0}, []LocalID{gone}) {
		t.Error("applyTo acted on a creature that has left play")
	}
	if core := g.State.Cards[gone]; core != (CardCore{}) {
		t.Errorf("core = %+v, want zero", core)
	}
}
