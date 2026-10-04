package web

import (
	"fmt"
	"testing"

	"github.com/dmikalova/vex/internal/engine"
)

// TestConditionalPunctuatesItsBranches pins the grammar a gated effect draws: the
// consequence opens with the gate mark and the alternative with the branch mark,
// so "if x, destroy a creature, otherwise draw a card" reads as two clauses rather
// than one run-on sentence of four glyphs (ADR 0047).
func TestConditionalPunctuatesItsBranches(t *testing.T) {
	gs := effectGlyphs(engine.Conditional{
		Cond: engine.PoolAember{Player: engine.Opponent, Is: engine.AtLeast, Amount: 7},
		Then: engine.Destroy{Target: engine.Target{Kind: engine.TargetChosenCreature}},
		Else: engine.Draw{Amount: 1},
	})
	if len(gs) != 3 {
		t.Fatalf("want 3 glyphs (destroy → creature, draw), got %d: %+v", len(gs), gs)
	}
	if gs[0].punct != punctGate {
		t.Errorf("Then head = %+v, want the gate mark opening the clause", gs[0])
	}
	if gs[1].punct != punctArrow {
		t.Errorf("target glyph = %+v, want the flow arrow left alone", gs[1])
	}
	if gs[2].punct != punctBranch {
		t.Errorf("Else head = %+v, want the branch mark", gs[2])
	}
}

// TestConditionalWithNoElseStillOpensItsClause checks the gate mark is the clause
// opener, not half of a pair: a Conditional with no Else still marks where its
// consequence begins.
func TestConditionalWithNoElseStillOpensItsClause(t *testing.T) {
	gs := effectGlyphs(engine.Conditional{
		Cond: engine.PoolAember{Player: engine.Opponent, Is: engine.AtLeast, Amount: 7},
		Then: engine.Draw{Amount: 1},
	})
	if len(gs) == 0 || gs[0].punct != punctGate {
		t.Errorf("glyphs = %+v, want the consequence opened by the gate mark", gs)
	}
	for _, g := range gs {
		if g.punct == punctBranch {
			t.Errorf("glyphs = %+v, want no branch mark without an Else", gs)
		}
	}
}

// TestChooseOneDividesItsOptions checks each option after the first opens with the
// divider, so a choice reads as alternatives rather than as a sequence of things
// that all happen. The choose glyph itself heads the run and takes no mark.
func TestChooseOneDividesItsOptions(t *testing.T) {
	gs := effectGlyphs(engine.ChooseOne{Options: []engine.Effect{
		engine.GainAember{Amount: 1},
		engine.Draw{Amount: 1},
		engine.Stun{Target: engine.Target{Kind: engine.TargetChosenEnemyCreature}},
	}})
	if len(gs) < 4 || gs[0].asset != "glyph-choose" || gs[0].punct != punctNone {
		t.Fatalf("glyphs = %+v, want an unmarked choose glyph heading the run", gs)
	}
	dividers := 0
	for _, g := range gs {
		if g.punct == punctDivider {
			dividers++
		}
	}
	if dividers != 2 {
		t.Errorf("glyphs = %+v, want 2 dividers for 3 options, got %d", gs, dividers)
	}
}

// TestDurationWrapperDrawsItsWindow checks a duration wrapper no longer renders as
// the bare effect it wraps. engine.ForDuration carried a window the strip dropped
// entirely, so a timed effect and an immediate one drew the same picture; the turn
// marker is now appended to the run (ADR 0047).
func TestDurationWrapperDrawsItsWindow(t *testing.T) {
	timed := effectGlyphs(engine.ForDuration{
		Duration: engine.RemainderOfPlayerTurn,
		Effects:  []engine.Effect{engine.GainAember{Amount: 1}},
	})
	immediate := effectGlyphs(engine.GainAember{Amount: 1})
	if len(timed) != len(immediate)+1 {
		t.Fatalf("timed = %+v, immediate = %+v, want one more glyph", timed, immediate)
	}
	if timed[len(timed)-1].asset != "phase-turn" {
		t.Errorf("timed = %+v, want the turn marker closing the run", timed)
	}
}

// TestEveryRepeatGateDraws holds the RepeatGate family's gate from the other side:
// every way a Repeat can decide to run again draws a real strip of existing
// assets. The looping gates are plumbing, so they draw the condition they fold
// into; only the exalting gate owes a picture of its own, and the distinctness
// check covers that one row.
func TestEveryRepeatGateDraws(t *testing.T) {
	assertFamilyDraws(t, "RepeatGate")
	assertFamilyRowsDrawApart(t, "RepeatGate")
}

// TestRepeatDrawsTheGateItLoopsOn checks the gate reaches the strip of the Repeat
// carrying it. A Repeat drew only the effect it repeats, so a free loop, a loop
// the controller may take again, and a loop paid for by exalting a creature all
// drew one picture — the picture of doing the thing once.
func TestRepeatDrawsTheGateItLoopsOn(t *testing.T) {
	once := effectGlyphs(engine.Draw{Amount: 1})
	paid := effectGlyphs(engine.Repeat{
		Do: engine.Draw{Amount: 1},
		Gate: engine.ByExalting{
			Creature: engine.Target{Kind: engine.TargetChosenFriendlyCreature},
		},
	})
	if len(paid) != len(once)+2 {
		t.Fatalf("paid repeat = %+v, want the effect plus the cost it charges", paid)
	}
	cost := paid[len(once)]
	if cost.asset != "aember" || cost.sigil != sigilExalt {
		t.Errorf("paid repeat charges %+v, want the exalted Æmber", cost)
	}
	if payer := paid[len(paid)-1]; payer.sigil != sigilMay {
		t.Errorf("the creature paying = %+v, want the may mark: the cost is declinable", payer)
	}
}

// TestLoopingGatesFoldIntoTheirCondition pins what the census calls these two:
// plumbing. A While adds the condition it loops on behind the gate mark and
// nothing else, and a MayWhileYouDo adds the may mark to say the controller
// chooses each further pass — so the two are told apart by that mark rather than
// by a drawing neither of them owes.
func TestLoopingGatesFoldIntoTheirCondition(t *testing.T) {
	cond := engine.Overwhelmed{}
	free := repeatGateGlyphs(engine.While{Cond: cond})
	chosen := repeatGateGlyphs(engine.MayWhileYouDo{Cond: cond})
	if len(free) == 0 || len(chosen) == 0 {
		t.Fatal("the looping gates drew no glyph")
	}
	plain := conditionGlyphs(cond)
	if len(free) != len(plain) || free[0].punct != punctGate {
		t.Errorf("a while loop = %+v, want the gated condition %+v", free, plain)
	}
	if chosen[0].sigil != sigilMay {
		t.Errorf("a may-while loop = %+v, want the may mark on its condition", chosen)
	}
}

// TestWhileYouDoDrawsLikeWhile pins that the result-gate arrow WhileYouDo's text
// wears for "only while the effect keeps happening" is a wording difference, not
// a second picture: it folds into its condition exactly as While does.
func TestWhileYouDoDrawsLikeWhile(t *testing.T) {
	cond := engine.Overwhelmed{}
	free := repeatGateGlyphs(engine.While{Cond: cond})
	result := repeatGateGlyphs(engine.WhileYouDo{Cond: cond})
	if len(result) == 0 {
		t.Fatal("WhileYouDo drew no glyph")
	}
	if fmt.Sprintf("%+v", free) != fmt.Sprintf("%+v", result) {
		t.Errorf("WhileYouDo drew %+v, want the same strip as While %+v", result, free)
	}
}
