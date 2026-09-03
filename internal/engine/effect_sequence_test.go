package engine

import "testing"

func TestSequenceEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	seq := Sequence{
		Effects: []Effect{
			GainAember{Player: Controller, Amount: 1},
			GainAember{Player: Controller, Amount: 2},
		},
	}
	if got := (Sequence{}).Text(); got != "" {
		t.Errorf("empty sequence text = %q, want empty", got)
	}
	if seq.Text() != "gain 1 Æmber, and gain 2 Æmber" {
		t.Errorf("sequence text = %q", seq.Text())
	}
	seq.Resolve(ctx)
	if g.State.Aember[0] != 3 {
		t.Errorf("aember = %d, want 3", g.State.Aember[0])
	}
}

func TestSentencesRendersEachChildAsItsOwnSentence(t *testing.T) {
	seq := Sentences{Effects: []Effect{
		DiscardTopOfDeck{Player: Opponent},
		RevealHand{Player: Opponent},
		GainAember{
			Player: Controller,
			Amount: 1,
			Per:    CardsInHand{Player: Opponent, House: TheContextualHouse},
		},
	}}
	want := "discard the top card of your opponent's deck. Reveal your opponent's hand. For each card of the discarded card's house revealed this way, gain 1 Æmber."
	if got := seq.Text(); got != want {
		t.Errorf("sequence text = %q, want %q", got, want)
	}
	if got := (Sentences{}).Text(); got != "" {
		t.Errorf("empty text = %q, want empty", got)
	}
}

func TestSequenceCombinesSameTarget(t *testing.T) {
	// Consecutive combinable effects on the same target fold into one phrase.
	both := Sequence{Effects: []Effect{
		Stun{Target: Target{Kind: TargetThisCreature}},
		Exhaust{Target: Target{Kind: TargetThisCreature}},
	}}
	if got := both.Text(); got != "stun and exhaust "+SelfName {
		t.Errorf("combined text = %q", got)
	}

	// Different targets do not fold, and a trailing non-combinable stands alone.
	mixed := Sequence{Effects: []Effect{
		Stun{Target: Target{Kind: TargetThisCreature}},
		Exhaust{Target: Target{Kind: TargetEachEnemyCreature}},
		GainAember{Player: Controller, Amount: 1},
	}}
	want := "stun " + SelfName + ", and exhaust each enemy creature, and gain 1 Æmber"
	if got := mixed.Text(); got != want {
		t.Errorf("mixed text = %q, want %q", got, want)
	}
}

func TestSequenceCombinesSameVerb(t *testing.T) {
	// Consecutive combinable effects sharing a verb fold their targets.
	both := Sequence{Effects: []Effect{
		Destroy{Target: Target{Kind: TargetChosenEnemyCreature}},
		Destroy{Target: Target{Kind: TargetChosenFriendlyCreature}},
	}}
	if got := both.Text(); got != "destroy an enemy creature and a friendly creature" {
		t.Errorf("combined text = %q", got)
	}

	// A run of three shared-verb effects folds every target.
	three := Sequence{Effects: []Effect{
		Destroy{Target: Target{Kind: TargetChosenEnemyCreature}},
		Destroy{Target: Target{Kind: TargetChosenFriendlyCreature}},
		Destroy{Target: Target{Kind: TargetEachCreature}},
	}}
	want := "destroy an enemy creature and a friendly creature and each creature"
	if got := three.Text(); got != want {
		t.Errorf("three text = %q, want %q", got, want)
	}
}
