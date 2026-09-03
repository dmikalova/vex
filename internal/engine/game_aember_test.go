package engine

import "testing"

func testEtherSpider() CardDefinition {
	return NewCard("Ether Spider", Mars, Creature, Uncommon,
		WithPower(7),
		WithReplaces(Instead{Of: EventAemberAddedToPool, Player: Opponent, With: Capture}))
}

func TestCaptureOpponentAemberReplacement(t *testing.T) {
	t.Run("captures gain effects instead of adding Æmber to the opponent pool", func(t *testing.T) {
		g := started(t)
		src := g.AddToBattleline(testCreature("src", 1), 0)
		spider := g.AddToBattleline(testEtherSpider(), 1)

		GainAember{
			Player: Controller,
			Amount: 2,
		}.Resolve(
			&EffectContext{Resolver: g, Source: src, Controller: 0},
		)

		if g.Aember(0) != 0 {
			t.Errorf("player Æmber = %d, want 0", g.Aember(0))
		}
		if g.AmberOn(spider) != 2 {
			t.Errorf("spider Æmber = %d, want 2", g.AmberOn(spider))
		}
	})

	t.Run(
		"only the incoming gain is captured; the pool owner keeps existing Æmber",
		func(t *testing.T) {
			g := started(t)
			src := g.AddToBattleline(testCreature("src", 1), 0)
			g.AddToBattleline(testEtherSpider(), 1)
			g.State.Aember[0] = 5 // Æmber the player already holds

			GainAember{
				Player: Controller,
				Amount: 2,
			}.Resolve(
				&EffectContext{Resolver: g, Source: src, Controller: 0},
			)

			if g.Aember(0) != 5 {
				t.Errorf(
					"player Æmber = %d, want 5 (existing pool untouched, only the gain captured)",
					g.Aember(0),
				)
			}
		},
	)

	t.Run("captures reap Æmber instead of adding it to the opponent pool", func(t *testing.T) {
		g := started(t)
		reaper := g.AddToBattleline(testCreature("reaper", 2), 0)
		spider := g.AddToBattleline(testEtherSpider(), 1)

		g.reapWith(reaper)

		if g.Aember(0) != 0 {
			t.Errorf("player Æmber = %d, want 0", g.Aember(0))
		}
		if g.AmberOn(spider) != 1 {
			t.Errorf("spider Æmber = %d, want 1", g.AmberOn(spider))
		}
	})

	t.Run("captures Æmber bonus instead of adding it to the opponent pool", func(t *testing.T) {
		g := started(t)
		spider := g.AddToBattleline(testEtherSpider(), 1)
		g.AddToHand(NewCard("Bonus", Brobnar, Tactic, Common, WithAemberBonus(2)), 0)

		if err := g.PlayAction(0, handIdx(g, 0, "Bonus")); err != nil {
			t.Fatalf("PlayAction: %v", err)
		}

		if g.Aember(0) != 0 {
			t.Errorf("player Æmber = %d, want 0", g.Aember(0))
		}
		if g.AmberOn(spider) != 2 {
			t.Errorf("spider Æmber = %d, want 2", g.AmberOn(spider))
		}
	})

	t.Run(
		"captures lasting reaction gains instead of adding them to the opponent pool",
		func(t *testing.T) {
			g := started(t)
			played := g.AddToBattleline(testCreature("played", 2), 0)
			spider := g.AddToBattleline(testEtherSpider(), 1)
			g.AddLasting(
				LastingEffect{On: EventCreaturePlayed, Do: actGainAember, Controller: 0, Amount: 2},
			)

			g.emitLasting(EventCreaturePlayed, 0, played)

			if g.Aember(0) != 0 {
				t.Errorf("player Æmber = %d, want 0", g.Aember(0))
			}
			if g.AmberOn(spider) != 2 {
				t.Errorf("spider Æmber = %d, want 2", g.AmberOn(spider))
			}
		},
	)

	t.Run("does not redirect steal or capture movements", func(t *testing.T) {
		g := started(t)
		src := g.AddToBattleline(testCreature("src", 1), 0)
		spider := g.AddToBattleline(testEtherSpider(), 1)
		g.State.Aember[1] = 3

		StealAember{Amount: 1}.Resolve(&EffectContext{Resolver: g, Source: src, Controller: 0})
		CaptureAember{Amount: 1, Target: Target{Kind: TargetThisCreature}, Source: Opponent}.
			Resolve(&EffectContext{Resolver: g, Source: src, Controller: 0})

		if g.Aember(0) != 1 {
			t.Errorf("player Æmber after steal = %d, want 1", g.Aember(0))
		}
		if g.Aember(1) != 1 {
			t.Errorf("opponent Æmber after steal and capture = %d, want 1", g.Aember(1))
		}
		if g.AmberOn(src) != 1 {
			t.Errorf("source captured Æmber = %d, want 1", g.AmberOn(src))
		}
		if g.AmberOn(spider) != 0 {
			t.Errorf("spider Æmber = %d, want 0", g.AmberOn(spider))
		}
	})
}

// A pool doubled past int16 (Binate Rupture, repeatedly) can be captured whole, so
// the card's counter saturates instead of wrapping into negative Æmber.
func TestAddAmberOnSaturates(t *testing.T) {
	g := started(t)
	id := g.AddToBattleline(testCreature("vault", 1), 0)

	g.AddAmberOn(id, maxCardAember-1)
	g.AddAmberOn(id, 100)
	if got := g.State.Cards[id].Amber; got != maxCardAember {
		t.Fatalf("Æmber on card = %d, want the %d ceiling", got, maxCardAember)
	}
	if err := g.InvariantError(); err != nil {
		t.Fatalf("saturating should keep the state sound, got %v", err)
	}

	g.AddAmberOn(id, -maxCardAember)
	if got := g.State.Cards[id].Amber; got != 0 {
		t.Fatalf("Æmber on card = %d, want 0", got)
	}
}
