package web

import (
	"fmt"
	"testing"

	"github.com/dmikalova/vex/internal/engine"
)

// conditionRows returns the Condition family's catalogued rows, so a test reads
// the census rather than a second list of conditions that could drift from it.
func conditionRows(t *testing.T) []iconRow {
	t.Helper()
	for _, family := range iconCatalog() {
		if family.name == "Condition" {
			return family.rows
		}
	}
	t.Fatal("the icon catalog walks no Condition family")
	return nil
}

// TestConditionGlyphAssetsExist walks every catalogued condition and checks each
// glyph it draws names a real asset. A glyph with a missing stem renders as a
// broken image rather than failing anything, so without this a typo in a noun
// ships silently.
func TestConditionGlyphAssetsExist(t *testing.T) {
	for _, row := range conditionRows(t) {
		gs, _ := row.draw()
		for _, g := range gs {
			if !assetExists(t, g.asset) {
				t.Errorf("condition %s draws %q, which has no web/assets/%s.svg",
					row.name, g.asset, g.asset)
			}
		}
	}
}

// TestAlwaysMetDrawsNothing pins the family's one silent row: an always-met
// condition adds no clause, so it draws neither a noun nor the gate mark that
// would open one. Drawing a bare gate would say a card is conditional when it is
// not — WithAemberCannotBeStolen() sets AlwaysMet to mean "on, unconditionally".
func TestAlwaysMetDrawsNothing(t *testing.T) {
	gs, covered := conditionGlyphs(engine.AlwaysMet{})
	if !covered {
		t.Error("AlwaysMet should be covered, not drawn as the unknown glyph")
	}
	if len(gs) != 0 {
		t.Errorf("AlwaysMet drew %+v, want nothing", gs)
	}
}

// TestTideReadsOneMarkTwoWays pins the one noun this family adds: the tide is a
// single state seen from the two sides of the table — high for you is low for your
// opponent — so high and low are one mark wearing the increase/decrease sigil
// pair, never two unrelated drawings.
func TestTideReadsOneMarkTwoWays(t *testing.T) {
	high, highCovered := conditionTestGlyphs(engine.TideIsHigh{})
	low, lowCovered := conditionTestGlyphs(engine.TideIsLow{})
	if !highCovered || !lowCovered {
		t.Fatalf("the tide conditions drew no real glyph")
	}
	if len(high) != 1 || len(low) != 1 {
		t.Fatalf("the tide drew %+v and %+v, want one glyph each", high, low)
	}
	if high[0].asset != low[0].asset {
		t.Errorf(
			"the tide drew %q high and %q low, want one mark",
			high[0].asset,
			low[0].asset,
		)
	}
	if high[0].sigil != sigilGain || low[0].sigil != sigilLose {
		t.Errorf(
			"the tide drew sigils %v high and %v low, want the increase/decrease pair",
			high[0].sigil,
			low[0].sigil,
		)
	}
}

// TestItConditionsMarkTheCardInContext pins the "It …" family's grammar: each one
// draws the property's own glyph — the stun mark, the trait tag, the nameplate,
// the bonus icon — marked as being about the card an effect or trigger just put in
// context, rather than inventing a drawing per condition.
func TestItConditionsMarkTheCardInContext(t *testing.T) {
	cases := map[string]struct {
		cond  engine.Condition
		asset string
	}{
		"stunned":    {engine.ItIsStunned{}, "stun"},
		"named":      {engine.ItIsNamed{Name: "Subtle Chain"}, nameGlyph().asset},
		"bonus icon": {engine.ItHasBonusIcon{}, "glyph-bonus"},
	}
	for name, tc := range cases {
		gs, covered := conditionTestGlyphs(tc.cond)
		if !covered {
			t.Errorf("%s drew no real glyph", name)
			continue
		}
		if gs[0].asset != tc.asset {
			t.Errorf("%s drew %q, want the %q glyph", name, gs[0].asset, tc.asset)
		}
		if gs[0].decor&decorChosen == 0 {
			t.Errorf("%s drew %+v, want it marked as the card in context", name, gs[0])
		}
	}
}

// TestFirstUseAndAnyUseShareOneTally pins how "first" is drawn: it is the same
// turn tally "you have done it at all" reads, against a different figure, so the
// two differ only in the comparison on the badge. Drawing "first" as a picture of
// its own would make the strip spell one tally two ways.
func TestFirstUseAndAnyUseShareOneTally(t *testing.T) {
	atAll, _ := conditionTestGlyphs(engine.UsedCreatureToReap{})
	first, _ := conditionTestGlyphs(engine.FirstReapOfTurn{})
	if len(atAll) != len(first) {
		t.Fatalf("the two reap conditions drew %+v and %+v, want one shape", atAll, first)
	}
	for i := range atAll {
		if atAll[i].asset != first[i].asset {
			t.Fatalf("the two reap conditions drew %+v and %+v, want one shape", atAll, first)
		}
	}
	if atAll[countedNoun(atAll)].comparison != engine.AtLeast {
		t.Errorf(
			"a reap this turn wears %+v, want an at-least badge",
			atAll[countedNoun(atAll)],
		)
	}
	if first[countedNoun(first)].comparison != engine.Exactly {
		t.Errorf("the first reap wears %+v, want an exactly badge", first[countedNoun(first)])
	}
}

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
