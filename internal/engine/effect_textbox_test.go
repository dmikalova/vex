package engine

import "testing"

// These tests cover GAINING A TEXT BOX — the shared mechanic behind Mimic Gel
// (permanent copy) and Creed of Nurture (remainder-of-turn loan): a creature
// takes on another card's printed traits, keywords, and triggered abilities.

// TestGrantedTextBoxSources reports the permanent copy, the turn loan, and both.
func TestGrantedTextBoxSources(t *testing.T) {
	g := started(t)
	src := g.AddToBattleline(testCreature("src", 3), 0)
	turnSrc := g.AddToBattleline(testCreature("turnsrc", 3), 0)
	recipient := g.AddToBattleline(testCreature("recipient", 2), 0)

	if got := g.grantedTextBoxSources(recipient); len(got) != 0 {
		t.Errorf("no gain should list no sources, got %v", got)
	}
	g.GrantTextBox(recipient, src, UntilCardLeavesPlay)
	if got := g.grantedTextBoxSources(recipient); len(got) != 1 || got[0] != src {
		t.Errorf("permanent gain = %v, want [%d]", got, src)
	}
	g.GrantTextBox(recipient, turnSrc, RemainderOfPlayerTurn)
	got := g.grantedTextBoxSources(recipient)
	if len(got) != 2 || got[0] != src || got[1] != turnSrc {
		t.Errorf("both gains = %v, want [%d %d]", got, src, turnSrc)
	}
}

// TestGrantTextBoxGone grants nothing when the recipient has left play.
func TestGrantTextBoxGone(t *testing.T) {
	g := started(t)
	src := g.AddToBattleline(testCreature("src", 3), 0)
	gone := g.AddToBattleline(testCreature("gone", 2), 0)
	g.DestroyEach(0, []LocalID{gone})
	g.GrantTextBox(gone, src, UntilCardLeavesPlay) // must not panic and must record nothing
	if g.State.Cards[gone].TextBoxSourcePlus != 0 {
		t.Error("granting to a creature no longer in play should do nothing")
	}
}

// TestGainTextBoxValidate rejects a missing recipient or source.
func TestGainTextBoxValidate(t *testing.T) {
	if (GainTextBox{Source: Target{Kind: TargetTheChosenCreature}}).validate() == nil {
		t.Error("an unset target should be invalid")
	}
	if (GainTextBox{Target: Target{Kind: TargetThisCreature}}).validate() == nil {
		t.Error("an unset source should be invalid")
	}
	if (GainTextBox{
		Target: Target{Kind: TargetThisCreature},
		Source: Target{Kind: TargetTheChosenCreature},
	}).validate() != nil {
		t.Error("a set target and source should be valid")
	}
}

// TestGainTextBoxText renders the printed clause.
func TestGainTextBoxText(t *testing.T) {
	e := GainTextBox{
		Target: Target{Kind: TargetThisCreature},
		Source: Target{Kind: TargetTheChosenCreature},
	}
	if got := e.Text(); got != SelfName+" gains the text box of the chosen creature" {
		t.Errorf("text = %q", got)
	}
}

// TestGainTextBoxResolve folds the source's traits, keywords, and triggered
// abilities into the recipient, and grants nothing when the source selects none.
func TestGainTextBoxResolve(t *testing.T) {
	g := started(t)
	g.SetRecording(true)
	source := g.AddToBattleline(
		testCreature(
			"source",
			3,
			WithTraits(Beast),
			WithKeywords(Skirmish),
			WithAbility(TriggerAfterReap, GainAember{
				Player: Controller,
				Amount: 1,
			}),
		),
		0,
	)
	recipient := g.AddToBattleline(testCreature("recipient", 2), 0)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
		It:         recipient,
		HasIt:      true,
	}
	g.SetChooser(0, idChooser{id: source})

	GainTextBox{
		Target: Target{Kind: TargetTriggeringCreature},
		Source: Target{Kind: TargetChosenCreature},
	}.Resolve(ctx)

	if !g.HasTrait(recipient, Beast) {
		t.Error("the recipient should gain the source's trait")
	}
	if !g.hasKeyword(recipient, Skirmish) {
		t.Error("the recipient should gain the source's keyword")
	}
	if !hasLogLine(g, "recipient gains the text box of source") {
		t.Error("the gain should be logged")
	}
	before := g.Aember(0)
	if err := g.Reap(0, recipient); err != nil {
		t.Fatalf("Reap: %v", err)
	}
	// +1 from the reap itself, +1 from the gained reap ability.
	if got := g.Aember(0); got != before+2 {
		t.Errorf("aember = %d, want %d (reap plus gained reap ability)", got, before+2)
	}

	// A source that selects nothing grants nothing.
	other := g.AddToBattleline(testCreature("other", 2), 0)
	g.SetChooser(0, orderRejectChooser{})
	GainTextBox{
		Target: Target{Kind: TargetTriggeringCreature},
		Source: Target{Kind: TargetChosenCreature},
	}.Resolve(&EffectContext{
		Resolver:   g,
		Controller: 0,
		It:         other,
		HasIt:      true,
	})
	if g.State.Cards[other].TextBoxSourcePlus != 0 {
		t.Error("a source that selects nothing should grant no text box")
	}
}

// TestGainTextBoxTurnExpires lends a text box for the remainder of the turn and
// clears it at the ready phase.
func TestGainTextBoxTurnExpires(t *testing.T) {
	g := started(t)
	source := g.AddToBattleline(testCreature("source", 3, WithTraits(Beast)), 0)
	recipient := g.AddToBattleline(testCreature("recipient", 2), 0)

	g.GrantTextBox(recipient, source, RemainderOfPlayerTurn)
	if g.State.Cards[recipient].TextBoxTurnSourcePlus != uint8(source)+1 {
		t.Error("the turn loan should be recorded")
	}
	if !g.HasTrait(recipient, Beast) {
		t.Error("the recipient should gain the source's trait for the turn")
	}

	// The ready phase clears the loan for every creature.
	g.StartTurn(0)
	g.EndPlayPhase(0)
	if g.State.Cards[recipient].TextBoxTurnSourcePlus != 0 {
		t.Error("the ready phase should clear the turn loan")
	}
}

func TestGainTextBoxResolveTurnWindow(t *testing.T) {
	g := started(t)
	source := g.AddToBattleline(testCreature("source", 3, WithTraits(Beast)), 0)
	recipient := g.AddToBattleline(testCreature("recipient", 2), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0, Source: recipient}
	g.SetChooser(0, idChooser{id: source})

	GainTextBox{
		Target:          Target{Kind: TargetThisCreature},
		Source:          Target{Kind: TargetChosenCreature},
		RemainderOfTurn: true,
	}.Resolve(ctx)

	if g.State.Cards[recipient].TextBoxTurnSourcePlus != uint8(source)+1 {
		t.Fatalf(
			"turn loan source = %d, want %d",
			g.State.Cards[recipient].TextBoxTurnSourcePlus-1,
			source,
		)
	}
	if !g.HasTrait(recipient, Beast) {
		t.Fatal("the recipient should gain the source's trait for the turn")
	}
}

// TestLendTextBoxFromHandText renders the whole instruction.
func TestLendTextBoxFromHandText(t *testing.T) {
	want := "reveal a creature from your hand and choose a creature in play. " +
		"For the remainder of the turn, the chosen creature gains the text box of " +
		"the revealed creature"
	if got := (LendTextBoxFromHand{}).Text(); got != want {
		t.Errorf("text = %q, want %q", got, want)
	}
}

// TestLendTextBoxFromHandResolve reveals a hand creature and lends its text box to
// a chosen creature in play for the turn.
func TestLendTextBoxFromHandResolve(t *testing.T) {
	g := started(t)
	g.SetRecording(true)
	recipient := g.AddToBattleline(testCreature("recipient", 4), 0)
	loaner := g.AddToHand(testCreature("loaner", 2, WithTraits(Beast)), 0)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}

	LendTextBoxFromHand{}.Resolve(ctx)

	if g.State.Cards[recipient].TextBoxTurnSourcePlus != uint8(loaner)+1 {
		t.Error("the chosen creature should gain the revealed creature's text box")
	}
	if !g.HasTrait(recipient, Beast) {
		t.Error("the chosen creature should gain the revealed creature's trait")
	}
	if handIdxByID(g, 0, loaner) < 0 {
		t.Error("the revealed creature should stay in hand")
	}
}

// TestLendTextBoxFromHandNoCreature does nothing with no creature in hand to
// reveal.
func TestLendTextBoxFromHandNoCreature(t *testing.T) {
	g := started(t)
	recipient := g.AddToBattleline(testCreature("recipient", 4), 0)
	g.AddToHand(NewCard("Filler", Brobnar, Tactic, Common), 0) // not a creature
	LendTextBoxFromHand{}.Resolve(&EffectContext{
		Resolver:   g,
		Controller: 0,
	})
	if g.State.Cards[recipient].TextBoxTurnSourcePlus != 0 {
		t.Error("with no creature in hand nothing should be lent")
	}
}

// TestLendTextBoxFromHandDecline stops when the controller declines the reveal.
func TestLendTextBoxFromHandDecline(t *testing.T) {
	g := started(t)
	recipient := g.AddToBattleline(testCreature("recipient", 4), 0)
	g.AddToHand(testCreature("a", 2), 0)
	g.AddToHand(testCreature("b", 2), 0)
	g.SetChooser(0, orderRejectChooser{}) // decline the reveal
	LendTextBoxFromHand{}.Resolve(&EffectContext{
		Resolver:   g,
		Controller: 0,
	})
	if g.State.Cards[recipient].TextBoxTurnSourcePlus != 0 {
		t.Error("declining the reveal should lend nothing")
	}
}

// TestLendTextBoxFromHandNoRecipient reveals a creature but lends nothing with no
// creature in play to receive it.
func TestLendTextBoxFromHandNoRecipient(t *testing.T) {
	g := started(t)
	loaner := g.AddToHand(testCreature("loaner", 2), 0)
	// No creature is in play, so the reveal happens but the loan lands on nothing.
	LendTextBoxFromHand{}.Resolve(&EffectContext{
		Resolver:   g,
		Controller: 0,
	})
	if handIdxByID(g, 0, loaner) < 0 {
		t.Error("the revealed creature should stay in hand")
	}
}

// TestCreatureGainedTextBoxText names the creature and the card whose text box it
// gained.
func TestCreatureGainedTextBoxText(t *testing.T) {
	g := started(t)
	a := g.AddToBattleline(testCreature("Mimic", 1), 0)
	b := g.AddToBattleline(testCreature("Troll", 8), 1)
	e := CreatureGainedTextBox{
		Creature: a,
		Source:   b,
	}
	if got := e.Text(g); got != "Mimic gains the text box of Troll" {
		t.Errorf("text = %q", got)
	}
}
