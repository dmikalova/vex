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
