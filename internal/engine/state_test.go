package engine

import "testing"

func TestCardList(t *testing.T) {
	var z CardList
	z.add(1)
	z.add(2)
	z.add(3)
	z.addFront(0) // [0 1 2 3]
	if z.indexOf(2) != 2 {
		t.Errorf("indexOf(2) = %d", z.indexOf(2))
	}
	if z.indexOf(9) != -1 {
		t.Errorf("indexOf(9) = %d, want -1", z.indexOf(9))
	}
	if !z.contains(3) || z.contains(9) {
		t.Error("contains check failed")
	}
	if got := z.removeAt(0); got != 0 { // [1 2 3]
		t.Errorf("removeAt(0) = %d, want 0", got)
	}
	if !z.remove(2) { // [1 3]
		t.Error("remove(2) should succeed")
	}
	if z.remove(9) {
		t.Error("remove(9) should fail")
	}
	got := z.slice()
	if len(got) != 2 || got[0] != 1 || got[1] != 3 {
		t.Errorf("zone slice = %v, want [1 3]", got)
	}
}

func TestCatalogAddPanicsOverCapacity(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	def := testCreature("Filler", 1)
	for range maxCards {
		g.Register(def, 0)
	}
	defer func() {
		if recover() == nil {
			t.Error("Register past maxCards should panic")
		}
	}()
	g.Register(def, 0)
}

func TestFastCopyIsIndependent(t *testing.T) {
	g := NewGame("A", "B", 1)
	def := NewCard("c", Brobnar, Creature, Common, WithPower(5))
	id := g.AddToBattleline(def, 0)
	g.State.Cards[id].Damage = 3
	g.State.Aember[0] = 2

	clone := g.State.FastCopy()
	clone.Cards[id].Damage = 99
	clone.Aember[0] = 99
	clone.Battleline[0].add(id)

	if g.State.Cards[id].Damage != 3 {
		t.Errorf("original damage mutated: %d", g.State.Cards[id].Damage)
	}
	if g.State.Aember[0] != 2 {
		t.Errorf("original aember mutated: %d", g.State.Aember[0])
	}
	if g.State.Battleline[0].Count != 1 {
		t.Errorf("original battleline mutated: count %d", g.State.Battleline[0].Count)
	}
}
