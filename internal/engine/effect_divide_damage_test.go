package engine

import "testing"

func TestSpreadDivideDamage(t *testing.T) {
	t.Run("text with a Per count front-loads the source", func(t *testing.T) {
		e := DealDamage{Spread: DivideDamage{
			Amount: 2,
			Per:    InPlay{Player: Controller, House: Brobnar, Type: Creature},
		}}
		want := "deal 2 damage for each friendly Brobnar creature, " +
			"divided among any number of creatures"
		if e.Text() != want {
			t.Errorf("text = %q, want %q", e.Text(), want)
		}
	})

	t.Run("text without a Per count states the flat amount", func(t *testing.T) {
		e := DealDamage{Spread: DivideDamage{Amount: 3}}
		want := "deal 3 damage, divided among any number of creatures"
		if e.Text() != want {
			t.Errorf("text = %q, want %q", e.Text(), want)
		}
	})

	t.Run("validate", func(t *testing.T) {
		if (DealDamage{Spread: DivideDamage{Amount: 0}}).validate() == nil {
			t.Error("a zero amount should be invalid")
		}
		if (DealDamage{Spread: DivideDamage{Amount: 1}}).validate() != nil {
			t.Error("a positive amount should be valid")
		}
	})

	t.Run("divides the pool among the chosen creatures", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		a := g.AddToBattleline(testCreature("a", 9), 0)
		b := g.AddToBattleline(testCreature("b", 9), 0)
		g.SetChooser(0, &idQueueChooser{ids: []LocalID{a, b, a, b}})
		DealDamage{Spread: DivideDamage{Amount: 4}}.Resolve(
			&EffectContext{Resolver: g, Controller: 0},
		)
		if g.Damage(a) != 2 || g.Damage(b) != 2 {
			t.Errorf("damage = %d/%d, want 2/2", g.Damage(a), g.Damage(b))
		}
	})

	t.Run("scales the pool by the Per count", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		a := g.AddToBattleline(testCreature("a", 9), 0)
		g.AddToBattleline(testCreature("b", 9), 0)
		// Two friendly Brobnar creatures, 2 each, all placed on the first by the
		// default chooser.
		DealDamage{Spread: DivideDamage{
			Amount: 2,
			Per:    InPlay{Player: Controller, House: Brobnar, Type: Creature},
		}}.Resolve(&EffectContext{Resolver: g, Controller: 0})
		if g.Damage(a) != 4 {
			t.Errorf("damage on a = %d, want 4", g.Damage(a))
		}
	})

	t.Run("a zero pool deals nothing", func(_ *testing.T) {
		g := NewGame("A", "B", 1)
		g.AddToBattleline(testCreature("a", 9), 0)
		DealDamage{Spread: DivideDamage{
			Amount: 2,
			Per:    InPlay{Player: Opponent, House: Brobnar, Type: Creature},
		}}.Resolve(&EffectContext{Resolver: g, Controller: 0})
	})

	t.Run("with no creatures in play deals nothing", func(_ *testing.T) {
		g := NewGame("A", "B", 1)
		DealDamage{Spread: DivideDamage{Amount: 3}}.Resolve(
			&EffectContext{Resolver: g, Controller: 0},
		)
	})
}
