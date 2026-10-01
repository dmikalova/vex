package web

import (
	"fmt"
	"testing"

	"github.com/dmikalova/vex/internal/engine"
)

// TestConditionOpensWithGateMark pins the condition family's one structural
// rule: a condition draws the gate mark that opens a hypothetical clause, so the
// strip says "what follows holds only if this does" rather than running the test
// and the consequence together into one sentence (ADR 0047).
func TestConditionOpensWithGateMark(t *testing.T) {
	gs, covered := conditionGlyphs(engine.PoolAember{
		Player: engine.Opponent,
		Is:     engine.AtLeast,
		Amount: 6,
	})
	if !covered {
		t.Fatalf("PoolAember drew no real glyph")
	}
	if gs[0].punct != punctGate {
		t.Errorf("condition opens with punct %v, want the gate mark", gs[0].punct)
	}
}

// TestCountedConditionDrawsItsCountWithAThreshold checks that a condition
// comparing a count against a figure reuses the count's own picture rather than
// restating it: "if your opponent has 6 or more Æmber" is the Æmber-in-pool count
// with an enemy tint, wearing "≥6".
func TestCountedConditionDrawsItsCountWithAThreshold(t *testing.T) {
	want, _ := countGlyphs(engine.AemberInPool{Player: engine.Opponent})
	got, _ := conditionGlyphs(engine.PoolAember{
		Player: engine.Opponent,
		Is:     engine.AtLeast,
		Amount: 6,
	})
	if len(got) != len(want) {
		t.Fatalf("condition drew %d glyphs, want the count's %d", len(got), len(want))
	}
	if got[0].asset != want[0].asset || got[0].decor != want[0].decor {
		t.Errorf("condition drew %+v, want the count's noun %+v", got[0], want[0])
	}
	if got[0].comparison != engine.AtLeast || got[0].qty != 6 {
		t.Errorf("condition drew %+v, want an at-least 6 badge", got[0])
	}
}

// TestThresholdBadgesTheCountedNoun checks that a threshold rides the noun being
// totalled rather than the run's head. "3 or more Untamed creatures in your
// discard pile" opens with the house emblem, and a "≥3" there would read as three
// houses rather than three creatures.
func TestThresholdBadgesTheCountedNoun(t *testing.T) {
	gs, covered := conditionGlyphs(engine.CardsInDiscardAtLeast{
		House:  engine.HouseMatcher{Kind: engine.MatchNamedHouse, House: engine.Untamed},
		Type:   engine.Creature,
		Amount: 3,
	})
	if !covered {
		t.Fatalf("CardsInDiscardAtLeast drew no real glyph")
	}
	if gs[0].qty != 0 || gs[0].comparison != comparisonNone {
		t.Errorf("house emblem %+v wears the threshold; it should narrow, not count", gs[0])
	}
	badged := gs[countedNoun(gs)]
	if badged.asset != "type-creature" || badged.qty != 3 ||
		badged.comparison != engine.AtLeast {
		t.Errorf("threshold rides %+v, want the creature noun with an at-least 3", badged)
	}
}

// TestParityThresholdDrawsNoFigure checks that a comparison the engine's Met
// ignores Amount for draws its mark alone. Even Ivan steals while the opponent's
// pool is even, which compares against no number, so printing the struct's unused
// Amount would show a figure the rule never reads.
func TestParityThresholdDrawsNoFigure(t *testing.T) {
	gs, _ := conditionGlyphs(engine.PoolAember{
		Player: engine.Opponent,
		Is:     engine.Even,
		Amount: 4,
	})
	if gs[0].comparison != engine.Even || gs[0].qty != 0 {
		t.Errorf("parity condition drew %+v, want the even mark and no figure", gs[0])
	}
}

// TestNotNegatesItsChild pins the one thing the Not combinator may not fold away.
// Not carries no test of its own, so it draws its child — but drawing the child
// unmarked would draw the opposite of what the card says, so it closes the run
// with the refusal mark the vocabulary already uses for everything barred.
func TestNotNegatesItsChild(t *testing.T) {
	child := engine.ForgedKey{Player: engine.Controller}
	inner, _ := conditionGlyphs(child)
	gs, covered := conditionGlyphs(engine.Not{Cond: child})
	if !covered {
		t.Fatalf("Not drew no real glyph")
	}
	if len(gs) != len(inner)+1 {
		t.Fatalf("Not drew %d glyphs, want its child's %d plus a mark", len(gs), len(inner))
	}
	if gs[len(gs)-1].asset != "glyph-ban" {
		t.Errorf("Not closes with %q, want the refusal mark", gs[len(gs)-1].asset)
	}
}

// TestOrDividesItsChildren checks that the two combinators that fold into their
// children are still told apart: an Or's children are alternatives and take the
// divider mark, an And's are a run and take none.
func TestOrDividesItsChildren(t *testing.T) {
	conds := []engine.Condition{
		engine.ForgedKey{Player: engine.Controller},
		engine.HasAember{Subject: engine.This},
	}
	or, _ := conditionGlyphs(engine.Or{Conditions: conds})
	and, _ := conditionGlyphs(engine.And{Conditions: conds})
	if fmt.Sprintf("%+v", or) == fmt.Sprintf("%+v", and) {
		t.Fatalf("Or and And drew the same strip %+v", or)
	}
	if or[len(or)-1].punct != punctDivider {
		t.Errorf("Or's second child opens with %v, want the divider mark", or[len(or)-1].punct)
	}
	if and[len(and)-1].punct != punctNone {
		t.Errorf("And's second child opens with %v, want no mark", and[len(and)-1].punct)
	}
}

// TestCardsInPlayDrawsOnePictureInBothSlots pins the drawing of the one Go value
// that appears verbatim in both CountCatalog() and ConditionCatalog(), because
// the type implements both interfaces. It draws one picture, read in whichever
// slot it sits: as a count it is "for each friendly creature", as a condition "if
// you control a friendly creature", and the cards surveyed are the same cards.
func TestCardsInPlayDrawsOnePictureInBothSlots(t *testing.T) {
	node := engine.CardsInPlay{Player: engine.Controller, Type: engine.Creature}
	asCount, _ := countGlyphs(node)
	asCondition, _ := conditionTestGlyphs(node)
	if fmt.Sprintf("%+v", asCount) != fmt.Sprintf("%+v", asCondition) {
		t.Errorf(
			"CardsInPlay drew %+v as a count and %+v as a condition, want one picture",
			asCount,
			asCondition,
		)
	}
}

// TestCardsInPlayNoneIsBarredRatherThanZero checks that "if there are no Ancient
// Bears in play" draws the refusal mark rather than a "≥0" badge. The engine
// reads None as its own word rather than as the Amount 0 case, and a "0 or more"
// threshold is met by every board.
func TestCardsInPlayNoneIsBarredRatherThanZero(t *testing.T) {
	gs, covered := conditionGlyphs(engine.CardsInPlay{
		Player: engine.Controller,
		Name:   "Ancient Bear",
		None:   true,
	})
	if !covered {
		t.Fatalf("CardsInPlay drew no real glyph")
	}
	if gs[len(gs)-1].asset != "glyph-ban" {
		t.Errorf("a none-in-play condition drew %+v, want it to close with the refusal mark", gs)
	}
	for _, g := range gs {
		if g.comparison != comparisonNone {
			t.Errorf("a none-in-play condition wears the threshold %+v; None is not Amount 0", g)
		}
	}
}
