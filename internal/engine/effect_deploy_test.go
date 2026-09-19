package engine

import "testing"

// This test covers the Deploy keyword: a creature with Deploy may enter play at
// any position in its controller's battleline, not only on a flank.

func TestDeployChoosesPosition(t *testing.T) {
	deployCreature := func() CardDefinition {
		return NewCard("Ranger", Brobnar, Creature, Common,
			WithPower(3), WithKeywords(Deploy))
	}

	// Deploy onto an empty battleline needs no choice: it lands at position 0.
	empty := started(t)
	d := empty.AddToHand(deployCreature(), 0)
	if _, err := empty.PlayCreature(0, handIdxByID(empty, 0, d), false); err != nil {
		t.Fatalf("deploy onto empty line: %v", err)
	}
	if got := empty.Battleline(0); len(got) != 1 || empty.Name(got[0]) != "Ranger" {
		t.Fatalf("empty-line deploy = %v, want [Ranger]", names(empty, got))
	}

	// Deploy between two existing creatures lands interior.
	mid := started(t)
	a := mid.AddToBattleline(testCreature("A", 3), 0)
	b := mid.AddToBattleline(testCreature("B", 3), 0)
	mid.SetChooser(0, optionPicker{idx: 1}) // "Between A and B"
	dm := mid.AddToHand(deployCreature(), 0)
	if _, err := mid.PlayCreature(0, handIdxByID(mid, 0, dm), false); err != nil {
		t.Fatalf("deploy between: %v", err)
	}
	line := mid.Battleline(0)
	if got := names(mid, line); got != "A Ranger B" {
		t.Errorf("interior deploy order = %q, want %q", got, "A Ranger B")
	}
	_ = a
	_ = b

	// Deploy to the left flank (choice 0) lands leftmost.
	left := started(t)
	left.AddToBattleline(testCreature("A", 3), 0)
	left.AddToBattleline(testCreature("B", 3), 0)
	left.SetChooser(0, optionPicker{idx: 0})
	dl := left.AddToHand(deployCreature(), 0)
	if _, err := left.PlayCreature(0, handIdxByID(left, 0, dl), false); err != nil {
		t.Fatalf("deploy left: %v", err)
	}
	if got := names(left, left.Battleline(0)); got != "Ranger A B" {
		t.Errorf("left-flank deploy order = %q, want %q", got, "Ranger A B")
	}

	// Deploy to the right flank (choice n) lands rightmost.
	right := started(t)
	right.AddToBattleline(testCreature("A", 3), 0)
	right.AddToBattleline(testCreature("B", 3), 0)
	right.SetChooser(0, optionPicker{idx: 2})
	dr := right.AddToHand(deployCreature(), 0)
	if _, err := right.PlayCreature(0, handIdxByID(right, 0, dr), false); err != nil {
		t.Fatalf("deploy right: %v", err)
	}
	if got := names(right, right.Battleline(0)); got != "A B Ranger" {
		t.Errorf("right-flank deploy order = %q, want %q", got, "A B Ranger")
	}
}

// positionPicker is a PositionChooser that lands a Deploy creature at a fixed
// position, and records the line it was asked about.
type positionPicker struct {
	pos  int
	line []LocalID
}

func (positionPicker) ChooseCreature(_, _ string, _ []LocalID) (LocalID, bool) {
	return 0, false
}

func (p *positionPicker) ChoosePosition(_, _ string, line []LocalID) int {
	p.line = line
	return p.pos
}

// TestDeployPositionChooser covers the PositionChooser capability: a chooser that
// speaks the battleline directly is pointed at the line and returns a position
// index, bypassing the labeled-option fallback.
func TestDeployPositionChooser(t *testing.T) {
	g := started(t)
	g.AddToBattleline(testCreature("A", 3), 0)
	g.AddToBattleline(testCreature("B", 3), 0)
	pick := &positionPicker{pos: 1} // between A and B
	g.SetChooser(0, pick)
	d := g.AddToHand(NewCard("Ranger", Brobnar, Creature, Common,
		WithPower(3), WithKeywords(Deploy)), 0)
	if _, err := g.PlayCreature(0, handIdxByID(g, 0, d), false); err != nil {
		t.Fatalf("deploy via PositionChooser: %v", err)
	}
	if got := names(g, g.Battleline(0)); got != "A Ranger B" {
		t.Errorf("position-chooser deploy order = %q, want %q", got, "A Ranger B")
	}
	if len(pick.line) != 2 {
		t.Errorf("chooser saw line of %d, want 2", len(pick.line))
	}
}

// names joins the printed names of ids into a space-separated string for readable
// battleline assertions.
func names(g *Game, ids []LocalID) string {
	out := ""
	for i, id := range ids {
		if i > 0 {
			out += " "
		}
		out += g.Name(id)
	}
	return out
}

// TestDeployPositionMeasuresTheLineAfterThePrompt pins that a flank position is
// measured against the battleline as it stands after the chooser answers, not the
// line the prompt was drawn from. Asking crosses a resolution boundary that
// settles destroyed creatures (ADR 0029), so a creature already at lethal damage
// reaches its discard pile while the prompt is up. A right flank measured before
// the prompt would then name a slot one past the end of the line, which the
// insert cannot reach.
func TestDeployPositionMeasuresTheLineAfterThePrompt(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.SetChooser(0, optionPicker{idx: 1}) // the right flank
	doomed := g.AddToBattleline(testCreature("Doomed", 1), 0)
	g.State.Cards[doomed].Damage = 1
	entering := g.Register(testCreature("Entering", 3), 0)

	pos, interior := g.deployPosition(0, entering, flankUnset, false)

	if g.State.Battleline[0].Count != 0 {
		t.Fatalf("battleline holds %d creatures, want the doomed one settled away by the prompt",
			g.State.Battleline[0].Count)
	}
	if pos != 0 || interior {
		t.Errorf("deployPosition = (%d, %v), want (0, false) — the line is empty by now",
			pos, interior)
	}
}
