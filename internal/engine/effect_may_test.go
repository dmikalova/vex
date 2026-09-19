package engine

import (
	"slices"
	"strings"
	"testing"
)

// bareEffect is an Effect that is not a struct, so it has no fields the shape
// rule could read.
type bareEffect string

func (bareEffect) Text() string           { return "do nothing" }
func (bareEffect) Resolve(*EffectContext) {}

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

// "You may destroy each creature" names no card to click. Off the board — a
// Tactic resolving its Play ability — the offer falls back to the Yes/No.
func TestMayWithoutACardChoiceStaysYesNo(t *testing.T) {
	e := May{Do: Destroy{Target: Target{Kind: TargetEachCreature}}}
	if e.Do.(declinableEffect).declinable() {
		t.Fatal("an untargeted Destroy should not be declinable")
	}

	g := NewGame("A", "B", 1)
	ch := &cardDecliner{}
	g.SetChooser(0, ch)
	doomed := g.AddToBattleline(testCreature("Doomed", 3), 0)
	tactic := g.AddToHand(testCreature("Tactic", 1), 0)
	e.Resolve(&EffectContext{Resolver: g, Source: tactic, Controller: 0})
	if ch.asked != 0 {
		t.Errorf("declinable prompts = %d, want 0", ch.asked)
	}
	if stillInPlay(g, doomed) {
		t.Error("the Yes answer should still have destroyed the creature")
	}
}

// On the board the same offer is opted into by clicking the card that makes it,
// so an ability with no card of its own to choose still reads as a click.
func TestMayWithoutACardChoiceClicksItsSource(t *testing.T) {
	e := May{Do: Destroy{Target: Target{Kind: TargetEachEnemyCreature}}}

	g := NewGame("A", "B", 1)
	g.SetChooser(0, &cardDecliner{})
	source := g.AddToBattleline(testCreature("Source", 3), 0)
	doomed := g.AddToBattleline(testCreature("Doomed", 3), 1)
	e.Resolve(&EffectContext{Resolver: g, Source: source, Controller: 0})
	if g.InPlay(doomed) {
		t.Error("clicking the source should have resolved the effect")
	}

	declined := NewGame("A", "B", 1)
	declined.SetChooser(0, &cardDecliner{decline: true})
	src := declined.AddToBattleline(testCreature("Source", 3), 0)
	spared := declined.AddToBattleline(testCreature("Doomed", 3), 1)
	e.Resolve(&EffectContext{Resolver: declined, Source: src, Controller: 0})
	if !declined.InPlay(spared) {
		t.Error("declining the source click should resolve nothing")
	}

	// A combatant mid-fight is not offered the click: the player just clicked it
	// to fight, so a second click would not read as opting in.
	fighting := NewGame("A", "B", 1)
	fighting.SetChooser(0, &cardDecliner{})
	att := fighting.AddToBattleline(testCreature("Source", 3), 0)
	def := fighting.AddToBattleline(testCreature("Doomed", 3), 1)
	fighting.State.FightersPlus = [2]LocalID{att + 1, def + 1}
	if (May{}).offeredOnSource(&EffectContext{Resolver: fighting, Source: att}) {
		t.Error("a creature resolving its own fight should keep the Yes/No")
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
	e := May{Do: Destroy{Target: Target{Kind: TargetEachCreature}.House(namedHouse(Mars))}}

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
	return slices.Contains(g.Battleline(0), id)
}

// The shape rule fires at init: a May over an effect that decides by choosing a
// card must offer that click. Ward carries a chosen Target but implements no
// declinable path, so wrapping it is rejected rather than silently asking twice.
func TestMayRejectsAnUnclickableCardChoice(t *testing.T) {
	chosen := May{Do: Ward{Target: Target{Kind: TargetChosenFriendlyCreature}}}
	err := chosen.validate()
	if err == nil {
		t.Fatal("a May over an unclickable card choice should be invalid")
	}
	if !strings.Contains(err.Error(), "not declinable") {
		t.Errorf("error = %q, want it to name the missing declinable path", err)
	}

	// An effect with no chosen Target of its own is a Yes/No and stays legal.
	each := May{Do: Ward{Target: Target{Kind: TargetEachFriendlyCreature}}}
	if err := each.validate(); err != nil {
		t.Errorf("a May over an untargeted effect was rejected: %v", err)
	}

	// An effect that already offers the click is legal whatever its Target says.
	clickable := May{Do: Destroy{Target: Target{Kind: TargetChosenEnemyCreature}}}
	if err := clickable.validate(); err != nil {
		t.Errorf("a May over a declinable effect was rejected: %v", err)
	}

	// A non-struct effect has no fields to read and is not a card choice.
	if chosenTargetField(bareEffect("")) {
		t.Error("a non-struct effect should not look like a card choice")
	}
}
