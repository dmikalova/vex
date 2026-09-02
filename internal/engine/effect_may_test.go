package engine

import "testing"

func TestMayText(t *testing.T) {
	e := May{Do: GainAember{Player: Controller, Amount: 1}}
	if got := e.Text(); got != "you may gain 1 Æmber" {
		t.Errorf("text = %q", got)
	}
}

func TestMayResolvesWhenAccepted(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.SetChooser(0, optionPicker{idx: 0}) // "Yes"
	ctx := &EffectContext{Resolver: g, Controller: 0}

	May{Do: GainAember{Player: Controller, Amount: 2}}.Resolve(ctx)
	if g.Aember(0) != 2 {
		t.Errorf("accepted: aember = %d, want 2", g.Aember(0))
	}
}

func TestMayDeclined(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.SetChooser(0, optionPicker{idx: 1}) // "No"
	ctx := &EffectContext{Resolver: g, Controller: 0}

	May{Do: GainAember{Player: Controller, Amount: 2}}.Resolve(ctx)
	if g.Aember(0) != 0 {
		t.Errorf("declined: aember = %d, want 0", g.Aember(0))
	}
}

// cardDecliner answers a declinable prompt: it takes the last candidate, or
// declines when decline is set.
type cardDecliner struct {
	FirstChooser
	decline bool
	asked   int
}

func (c *cardDecliner) ChooseCardOrDecline(
	_, _ string,
	candidates []LocalID,
) (LocalID, bool) {
	c.asked++
	if c.decline {
		return 0, false
	}
	return candidates[len(candidates)-1], true
}

// A "you may destroy a creature" is one card choice, so the player picks the
// creature directly instead of first answering Yes.
func TestMayDeclinableIsAskedAsACardChoice(t *testing.T) {
	g := NewGame("A", "B", 1)
	ch := &cardDecliner{}
	g.SetChooser(0, ch)
	keep := g.AddToBattleline(testCreature("Keep", 3), 0)
	doomed := g.AddToBattleline(testCreature("Doomed", 3), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := May{Do: Destroy{Target: Target{Kind: TargetChosenFriendlyCreature}}}
	if !e.Do.(declinableEffect).declinable() {
		t.Fatal("a chosen-target Destroy should be declinable")
	}
	e.Resolve(ctx)
	if ch.asked != 1 {
		t.Errorf("declinable prompts = %d, want 1", ch.asked)
	}
	if stillInPlay(g, doomed) {
		t.Error("the chosen creature should have been destroyed")
	}
	if !stillInPlay(g, keep) {
		t.Error("the unchosen creature should have survived")
	}
}

// Declining the card choice is declining the whole "you may".
func TestMayDeclinableDeclined(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.SetChooser(0, &cardDecliner{decline: true})
	doomed := g.AddToBattleline(testCreature("Doomed", 3), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	May{Do: Destroy{Target: Target{Kind: TargetChosenFriendlyCreature}}}.Resolve(ctx)
	if !stillInPlay(g, doomed) {
		t.Error("a declined May should destroy nothing")
	}
}

// A gate's follow-up still hangs off its first half happening.
func TestMayDeclinableGate(t *testing.T) {
	e := May{Do: Then{
		First:  Destroy{Target: Target{Kind: TargetChosenFriendlyCreature}},
		Result: GainAember{Player: Controller, Amount: 2},
	}}
	if !e.Do.(declinableEffect).declinable() {
		t.Fatal("a gate on a chosen-target Destroy should be declinable")
	}

	g := NewGame("A", "B", 1)
	g.SetChooser(0, &cardDecliner{})
	g.AddToBattleline(testCreature("Doomed", 3), 0)
	e.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if g.Aember(0) != 2 {
		t.Errorf("accepted gate: aember = %d, want 2", g.Aember(0))
	}

	declined := NewGame("A", "B", 1)
	declined.SetChooser(0, &cardDecliner{decline: true})
	declined.AddToBattleline(testCreature("Doomed", 3), 0)
	e.Resolve(&EffectContext{Resolver: declined, Controller: 0})
	if declined.Aember(0) != 0 {
		t.Errorf("declined gate: aember = %d, want 0", declined.Aember(0))
	}
}

// The verbs of a chosen creature are offered the same way, so Sergeant Zakiel is
// answered by clicking the neighbor.
func TestMayDeclinableChosenCreatureVerbs(t *testing.T) {
	e := May{Do: OnChooseCreature{
		Target: Target{Kind: TargetChosenFriendlyCreature},
		Verbs:  []CreatureVerb{ReadyVerb{}},
	}}
	if !e.Do.(declinableEffect).declinable() {
		t.Fatal("a chosen-target OnChooseCreature should be declinable")
	}

	g := NewGame("A", "B", 1)
	g.SetChooser(0, &cardDecliner{})
	ally := g.AddToBattleline(testCreature("Ally", 3), 0)
	g.SetExhausted(ally, true)
	e.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if g.Exhausted(ally) {
		t.Error("the chosen creature should have been readied")
	}

	declined := NewGame("A", "B", 1)
	declined.SetChooser(0, &cardDecliner{decline: true})
	other := declined.AddToBattleline(testCreature("Ally", 3), 0)
	declined.SetExhausted(other, true)
	e.Resolve(&EffectContext{Resolver: declined, Controller: 0})
	if !declined.Exhausted(other) {
		t.Error("a declined May should ready nothing")
	}
}

// "You may destroy each creature" names no card to click, so it stays a Yes/No.
func TestMayWithoutACardChoiceStaysYesNo(t *testing.T) {
	e := May{Do: Destroy{Target: Target{Kind: TargetEachCreature}}}
	if e.Do.(declinableEffect).declinable() {
		t.Fatal("an untargeted Destroy should not be declinable")
	}

	g := NewGame("A", "B", 1)
	ch := &cardDecliner{}
	g.SetChooser(0, ch)
	doomed := g.AddToBattleline(testCreature("Doomed", 3), 0)
	e.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if ch.asked != 0 {
		t.Errorf("declinable prompts = %d, want 0", ch.asked)
	}
	if stillInPlay(g, doomed) {
		t.Error("the Yes answer should still have destroyed the creature")
	}
}

// A target that chooses nothing has no decision to decline, so SelectOptional
// behaves like Select.
func TestSelectOptionalWithoutAChoice(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.SetChooser(0, &cardDecliner{decline: true})
	id := g.AddToBattleline(testCreature("Ally", 3), 0)
	got := Target{Kind: TargetEachCreature}.SelectOptional(&EffectContext{
		Resolver:   g,
		Controller: 0,
	})
	if len(got) != 1 || got[0] != id {
		t.Errorf("selected = %v, want [%d]", got, id)
	}
}

// optionRecorder counts the Yes/No questions asked and always answers Yes.
type optionRecorder struct {
	FirstChooser
	asked int
}

func (o *optionRecorder) ChooseOption(_, _ string, _ []string) int {
	o.asked++
	return 0
}

// "You may destroy each Mars creature" with no Mars creature in play is not a
// decision, so the question is never asked.
func TestMayWithNothingToDoIsNotOffered(t *testing.T) {
	e := May{Do: Destroy{Target: Target{Kind: TargetEachCreature}.OfHouse(Mars)}}

	g := NewGame("A", "B", 1)
	ch := &optionRecorder{}
	g.SetChooser(0, ch)
	safe := g.AddToBattleline(testCreature("Bystander", 3), 0)
	e.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if ch.asked != 0 {
		t.Errorf("Yes/No prompts = %d, want 0", ch.asked)
	}
	if !stillInPlay(g, safe) {
		t.Error("an off-house creature should not have been destroyed")
	}

	mars := testCreature("Raider", 3)
	mars.House = Mars
	g.AddToBattleline(mars, 0)
	e.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if ch.asked != 1 {
		t.Errorf("Yes/No prompts with a target = %d, want 1", ch.asked)
	}
}

func TestMayValidate(t *testing.T) {
	bad := May{Do: Heal{Fully: true, Amount: 1, Target: Target{Kind: TargetThisCreature}}}
	if validateEffect(bad) == nil {
		t.Error("May wrapping an invalid effect should fail validation")
	}
	good := May{Do: GainAember{Player: Controller, Amount: 1}}
	if validateEffect(good) != nil {
		t.Error("May wrapping a valid effect should pass validation")
	}
}

// stillInPlay reports whether an id is still in player 0's battleline.
func stillInPlay(g *Game, id LocalID) bool {
	for _, x := range g.Battleline(0) {
		if x == id {
			return true
		}
	}
	return false
}
