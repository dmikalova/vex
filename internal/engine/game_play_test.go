package engine

import "testing"

func TestCanPlay(t *testing.T) {
	g := started(t) // player 0 active, Brobnar
	creat := g.AddToHand(testCreature("c", 3), 0)
	if err := g.CanPlay(0, creat); err != nil {
		t.Errorf("playable creature = %v, want nil", err)
	}
	if err := g.CanPlay(1, creat); err != ErrNotActivePlayer {
		t.Errorf("wrong player = %v, want ErrNotActivePlayer", err)
	}

	off := g.AddToHand(NewCard("off", Sanctum, Creature, Common, WithPower(1)), 0)
	if err := g.CanPlay(0, off); err != ErrWrongHouse {
		t.Errorf("off-house = %v, want ErrWrongHouse", err)
	}

	up := g.AddToHand(NewCard("up", Brobnar, Upgrade, Common), 0)
	if err := g.CanPlay(0, up); err != ErrNoTarget {
		t.Errorf("hostless upgrade = %v, want ErrNoTarget", err)
	}

	g.State.Winner = 1
	if err := g.CanPlay(0, creat); err != ErrGameOver {
		t.Errorf("game over = %v, want ErrGameOver", err)
	}
}

func TestCanPlayRestrictions(t *testing.T) {
	// Creatures barred: a card that stops the player playing creatures.
	g := started(t)
	g.AddToBattleline(NewCard("Blocker", Brobnar, Creature, Common, WithPower(1),
		WithRestrictions(Restrictions{CannotPlay: Creature})), 0)
	c := g.AddToHand(testCreature("c", 3), 0)
	if err := g.CanPlay(0, c); err != ErrCannotPlayCreature {
		t.Errorf("creatures barred = %v, want ErrCannotPlayCreature", err)
	}

	// Card-play limit reached this turn.
	g2 := started(t)
	g2.AddToBattleline(NewCard(
		"Imp",
		Brobnar,
		Creature,
		Common,
		WithPower(1),
		WithRestrictions(
			Restrictions{PlayCardLimit: PlayCardLimit{Player: Controller, Amount: 2}},
		),
	), 0)
	g2.State.PlayedThisTurn[0].Count = 2
	c2 := g2.AddToHand(testCreature("c2", 3), 0)
	if err := g2.CanPlay(0, c2); err != ErrCardPlayLimit {
		t.Errorf("limit reached = %v, want ErrCardPlayLimit", err)
	}
}

func TestCardPlayLimit(t *testing.T) {
	g := started(t) // player 0 active, Brobnar
	// Player 1 controls an Ember-Imp-like card limiting player 0 to two plays.
	g.AddToBattleline(NewCard(
		"imp",
		Dis,
		Creature,
		Common,
		WithPower(2),
		WithRestrictions(
			Restrictions{PlayCardLimit: PlayCardLimit{Player: Opponent, Amount: 2}},
		),
	), 1)
	g.AddToHand(testCreature("c0", 3), 0)
	g.AddToHand(testCreature("c1", 3), 0)
	g.AddToHand(testCreature("c2", 3), 0)
	g.AddToHand(exBruteStrength(), 0) // an upgrade

	if _, err := g.PlayCreature(0, handIdx(g, 0, "c0"), false); err != nil {
		t.Fatalf("first play: %v", err)
	}
	if _, err := g.PlayCreature(0, handIdx(g, 0, "c1"), false); err != nil {
		t.Fatalf("second play: %v", err)
	}
	// The third card play is barred, whatever its type.
	if _, err := g.PlayCreature(0, handIdx(g, 0, "c2"), false); err != ErrCardPlayLimit {
		t.Errorf("third creature play = %v, want ErrCardPlayLimit", err)
	}
	if _, err := g.PlayUpgrade(0, handIdx(g, 0, "Brute Strength")); err != ErrCardPlayLimit {
		t.Errorf("upgrade play = %v, want ErrCardPlayLimit", err)
	}

	// A new turn resets the count.
	g.EndTurn(0)
	g.BeginTurn(0)
	if err := g.ChooseHouse(0, Brobnar); err != nil {
		t.Fatal(err)
	}
	if _, err := g.PlayCreature(0, handIdx(g, 0, "c2"), false); err != nil {
		t.Errorf("play after new turn = %v, want nil (count reset)", err)
	}
}

func TestPlayedThisTurn(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	g.BeginTurn(0)
	creature := g.AddToHand(NewCard("brobnar creature", Brobnar, Creature, Common, WithPower(3)), 0)
	artifact := g.AddToHand(NewCard("sanctum artifact", Sanctum, Artifact, Common), 0)
	action := g.AddToHand(NewCard("sanctum action", Sanctum, Tactic, Common), 0)
	upgrade := g.AddToHand(NewCard("mars upgrade", Mars, Upgrade, Common), 0)

	if _, err := g.PlayCreature(0, handIdxByID(g, 0, creature), false); err != nil {
		t.Fatalf("PlayCreature: %v", err)
	}
	if _, err := g.PlayArtifact(0, handIdxByID(g, 0, artifact)); err != nil {
		t.Fatalf("PlayArtifact: %v", err)
	}
	if err := g.PlayAction(0, handIdxByID(g, 0, action)); err != nil {
		t.Fatalf("PlayAction: %v", err)
	}
	if _, err := g.PlayUpgrade(0, handIdxByID(g, 0, upgrade)); err != nil {
		t.Fatalf("PlayUpgrade: %v", err)
	}

	playsOf := func(house House) int {
		n := 0
		for _, id := range g.PlayedThisTurn(0) {
			if g.House(id) == house {
				n++
			}
		}
		return n
	}
	if got := playsOf(Brobnar); got != 1 {
		t.Errorf("Brobnar plays = %d, want 1", got)
	}
	if got := playsOf(Sanctum); got != 2 {
		t.Errorf("Sanctum plays = %d, want 2", got)
	}
	if got := playsOf(Mars); got != 1 {
		t.Errorf("Mars plays = %d, want 1", got)
	}
	if got := len(g.PlayedThisTurn(0)); got != 4 {
		t.Errorf("total plays = %d, want 4", got)
	}

	g.BeginTurn(0)
	if got := playsOf(Sanctum); got != 0 {
		t.Errorf("Sanctum plays after reset = %d, want 0", got)
	}
	if got := len(g.PlayedThisTurn(0)); got != 0 {
		t.Errorf("total plays after reset = %d, want 0", got)
	}
}

func TestOffHousePlayGrant(t *testing.T) {
	witch := NewCard(
		"Witch",
		Untamed,
		Creature,
		Rare,
		WithPower(4),
		WithPlayPermission(PlayPermission{House: Untamed, Count: 1}),
	)

	t.Run("allows one off-house play and consumes it", func(t *testing.T) {
		g := started(t)
		g.AddToBattleline(witch, 0)
		first := g.AddToHand(NewCard("bear", Untamed, Creature, Common, WithPower(3)), 0)
		second := g.AddToHand(NewCard("wolf", Untamed, Creature, Common, WithPower(3)), 0)

		if err := g.CanPlay(0, first); err != nil {
			t.Fatalf("CanPlay first off-house Untamed = %v, want nil", err)
		}
		if got := g.State.PlayPermissionsUsedThisTurn[0][Untamed]; got != 0 {
			t.Fatalf("CanPlay consumed off-house grant: used = %d, want 0", got)
		}
		if _, err := g.PlayCreature(0, handIdxByID(g, 0, first), false); err != nil {
			t.Fatalf("first off-house play: %v", err)
		}
		if got := g.State.PlayPermissionsUsedThisTurn[0][Untamed]; got != 1 {
			t.Fatalf("off-house plays used = %d, want 1", got)
		}
		if err := g.CanPlay(0, second); err != ErrWrongHouse {
			t.Fatalf("CanPlay second off-house Untamed = %v, want ErrWrongHouse", err)
		}
		if _, err := g.PlayCreature(0, handIdxByID(g, 0, second), false); err != ErrWrongHouse {
			t.Fatalf("second off-house play = %v, want ErrWrongHouse", err)
		}

		g.EndTurn(0)
		g.BeginTurn(0)
		if err := g.ChooseHouse(0, Brobnar); err != nil {
			t.Fatal(err)
		}
		if _, err := g.PlayCreature(0, handIdxByID(g, 0, second), false); err != nil {
			t.Fatalf("off-house play after reset: %v", err)
		}
	})

	t.Run("does not matter when Untamed is active", func(t *testing.T) {
		g := NewGame("Alice", "Bob", 1)
		g.BeginTurn(0)
		if err := g.ChooseHouse(0, Untamed); err != nil {
			t.Fatal(err)
		}
		g.AddToBattleline(witch, 0)
		first := g.AddToHand(NewCard("boar", Untamed, Creature, Common, WithPower(3)), 0)
		second := g.AddToHand(NewCard("fox", Untamed, Creature, Common, WithPower(3)), 0)

		if _, err := g.PlayCreature(0, handIdxByID(g, 0, first), false); err != nil {
			t.Fatalf("first active-house play: %v", err)
		}
		if _, err := g.PlayCreature(0, handIdxByID(g, 0, second), false); err != nil {
			t.Fatalf("second active-house play: %v", err)
		}
		if got := g.State.PlayPermissionsUsedThisTurn[0][Untamed]; got != 0 {
			t.Fatalf("active-house plays used off-house grant: got %d, want 0", got)
		}
	})

	t.Run("rejects off-house play without a grant", func(t *testing.T) {
		g := started(t)
		untamed := g.AddToHand(NewCard("badger", Untamed, Creature, Common, WithPower(3)), 0)

		if err := g.CanPlay(0, untamed); err != ErrWrongHouse {
			t.Fatalf("CanPlay without Witch = %v, want ErrWrongHouse", err)
		}
		if _, err := g.PlayCreature(0, handIdxByID(g, 0, untamed), false); err != ErrWrongHouse {
			t.Fatalf("PlayCreature without Witch = %v, want ErrWrongHouse", err)
		}
	})

	t.Run("counts each controlled grant", func(t *testing.T) {
		g := started(t)
		g.AddToBattleline(witch, 0)
		g.AddToBattleline(witch, 0)
		first := g.AddToHand(NewCard("hare", Untamed, Creature, Common, WithPower(3)), 0)
		second := g.AddToHand(NewCard("lynx", Untamed, Creature, Common, WithPower(3)), 0)
		third := g.AddToHand(NewCard("mole", Untamed, Creature, Common, WithPower(3)), 0)

		if _, err := g.PlayCreature(0, handIdxByID(g, 0, first), false); err != nil {
			t.Fatalf("first off-house play with two grants: %v", err)
		}
		if _, err := g.PlayCreature(0, handIdxByID(g, 0, second), false); err != nil {
			t.Fatalf("second off-house play with two grants: %v", err)
		}
		if _, err := g.PlayCreature(0, handIdxByID(g, 0, third), false); err != ErrWrongHouse {
			t.Fatalf("third off-house play with two grants = %v, want ErrWrongHouse", err)
		}
	})
}

func TestPlayCardLimitTargets(t *testing.T) {
	cases := []struct {
		name         string
		player       Player
		limited      int
		unrestricted int
	}{
		{name: "controller", player: Controller, limited: 0, unrestricted: 1},
		{name: "opponent", player: Opponent, limited: 1, unrestricted: 0},
		{name: "each player", player: EachPlayer, limited: 0, unrestricted: -1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			g := started(t)
			g.AddToBattleline(NewCard(
				"limit",
				Brobnar,
				Creature,
				Common,
				WithPower(1),
				WithRestrictions(
					Restrictions{PlayCardLimit: PlayCardLimit{Player: tc.player, Amount: 2}},
				),
			), 0)
			g.State.PlayedThisTurn[tc.limited].Count = 2
			if !g.cannotPlayCard(tc.limited) {
				t.Errorf("player %d should be limited", tc.limited)
			}
			if tc.unrestricted >= 0 && g.cannotPlayCard(tc.unrestricted) {
				t.Errorf("player %d should not be limited", tc.unrestricted)
			}
		})
	}
	if (PlayCardLimit{}).affects(0, 0) {
		t.Error("an unset play-card limit should affect no player")
	}
}
