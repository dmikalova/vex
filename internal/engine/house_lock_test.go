package engine

import (
	"strings"
	"testing"
)

// lockedCreature is a creature definition carrying a house lock, used to stand a
// lock up in play without depending on a real card.
func lockedCreature(name string, lock HouseLock) CardDefinition {
	def := testCreature(name, 3)
	def.HouseLock = lock
	return def
}

func TestHouseLockText(t *testing.T) {
	cases := []struct {
		name string
		lock HouseLock
		want string
	}{
		{"unset", HouseLock{}, ""},
		{
			"named on play prints nothing",
			HouseLock{Player: Opponent, Bars: true},
			"",
		}, {
			"controller must",
			HouseLock{Player: Controller, House: Dis},
			"While {self} is in play you must choose Dis as your active house.",
		}, {
			"opponent cannot",
			HouseLock{Player: Opponent, House: Mars, Bars: true},
			"While {self} is in play your opponent cannot choose Mars as their active house.",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.lock.text(); got != tc.want {
				t.Errorf("text() = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestHouseLockLocked covers both sources of the locked house: printed on the
// card, or named by the card as it entered play.
func TestHouseLockLocked(t *testing.T) {
	if got := (HouseLock{House: Dis}).locked(Mars); got != Dis {
		t.Errorf("printed house = %v, want Dis", got)
	}
	if got := (HouseLock{}).locked(Mars); got != Mars {
		t.Errorf("named house = %v, want Mars", got)
	}
}

// TestChooseHouseLockRequires covers a "must choose" lock (Pitlord): it binds its
// controller, leaves the opponent free, and yields when the player lacks the house.
func TestChooseHouseLockRequires(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.SetPlayerHouses(0, []House{Dis, Logos, Untamed})
	g.SetPlayerHouses(1, []House{Brobnar, Mars, Shadows})
	g.AddToBattleline(lockedCreature("Pit", HouseLock{Player: Controller, House: Dis}), 0)
	g.State.ActivePlayer = 0
	if err := g.ChooseHouse(0, Logos); err != ErrHouseLocked {
		t.Errorf("choosing another house = %v, want ErrHouseLocked", err)
	}
	if err := g.ChooseHouse(0, Dis); err != nil {
		t.Errorf("choosing the locked house = %v, want nil", err)
	}

	g.State.ActivePlayer = 1
	if err := g.ChooseHouse(1, Brobnar); err != nil {
		t.Errorf("the opponent is unconstrained = %v, want nil", err)
	}
}

// TestChooseHouseLockRequiresYieldsWhenUnavailable covers cannot-overrides-must for
// a lock, matching how a forced house yields: a house the player's deck lacks
// cannot be demanded.
func TestChooseHouseLockRequiresYieldsWhenUnavailable(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.SetPlayerHouses(0, []House{Brobnar, Logos, Untamed})
	g.AddToBattleline(lockedCreature("Pit", HouseLock{Player: Controller, House: Dis}), 0)
	g.State.ActivePlayer = 0
	if err := g.ChooseHouse(0, Brobnar); err != nil {
		t.Errorf("choosing an available house = %v, want nil", err)
	}
}

// TestChooseHouseLockBars covers a "cannot choose" lock whose house is named on
// play (Restringuntus): it binds the controller's opponent, and does nothing until
// a house has actually been named.
func TestChooseHouseLockBars(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.SetPlayerHouses(1, []House{Brobnar, Mars, Shadows})
	id := g.AddToBattleline(lockedCreature("Restr", HouseLock{Player: Opponent, Bars: true}), 0)
	g.State.ActivePlayer = 1
	if err := g.ChooseHouse(1, Mars); err != nil {
		t.Errorf("before a house is named = %v, want nil", err)
	}

	g.SetNamedHouse(id, Mars)
	if err := g.ChooseHouse(1, Mars); err != ErrHouseLocked {
		t.Errorf("choosing the barred house = %v, want ErrHouseLocked", err)
	}
	if err := g.ChooseHouse(1, Brobnar); err != nil {
		t.Errorf("choosing another house = %v, want nil", err)
	}
}

// TestWithHouseLockRenders covers the option and the printed rule line it adds,
// with {self} substituted for the card's own name.
func TestWithHouseLockRenders(t *testing.T) {
	pit := NewCard("Pitlord", Dis, Creature, Rare,
		WithPower(9),
		WithHouseLock(HouseLock{Player: Controller, House: Dis}))
	want := "While Pitlord is in play you must choose Dis as your active house."
	if got := RenderCardText(&pit); !strings.Contains(got, want) {
		t.Errorf("RenderCardText() = %q, want it to contain %q", got, want)
	}
}

func TestNameHouse(t *testing.T) {
	t.Run("validate requires a player", func(t *testing.T) {
		if err := (NameHouse{}).validate(); err == nil {
			t.Error("validate() = nil, want an error for the unset player")
		}
		if err := (NameHouse{Player: Opponent}).validate(); err != nil {
			t.Errorf("validate() = %v, want nil", err)
		}
	})

	t.Run("text names whose choice the house binds", func(t *testing.T) {
		want := "your opponent cannot choose that house as their active house until " +
			"{self} leaves play"
		if got := (NameHouse{Player: Opponent}).Text(); got != want {
			t.Errorf("Text() = %q, want %q", got, want)
		}
		want = "you cannot choose that house as your active house until {self} leaves play"
		if got := (NameHouse{Player: Controller}).Text(); got != want {
			t.Errorf("Text() = %q, want %q", got, want)
		}
	})

	t.Run("resolve records the chosen house on its source", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		id := g.AddToBattleline(testCreature("Restr", 3), 0)
		ctx := &EffectContext{Resolver: g, Source: id, ChosenHouse: Mars}
		NameHouse{Player: Opponent}.Resolve(ctx)
		if got := g.State.Cards[id].NamedHouse; got != Mars {
			t.Errorf("named house = %v, want Mars", got)
		}
	})
}
