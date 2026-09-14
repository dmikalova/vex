package engine

import "testing"

func TestPurge(t *testing.T) {
	// Text variants.
	if got := (PurgeCard{Player: ChosenPlayer, Selection: Chosen{Optional: true}, Amount: 2}).Text(); got != "purge up to 2 cards from a discard pile" {
		t.Errorf("up-to text = %q", got)
	}
	if got := (PurgeCard{Player: ChosenPlayer, Selection: Chosen{Type: Creature}}).Text(); got != "purge a creature from a discard pile" {
		t.Errorf("single text = %q", got)
	}
	if got := (PurgeCard{Player: ChosenPlayer, Selection: Chosen{}, Amount: 2}).Text(); got != "purge 2 cards from a discard pile" {
		t.Errorf("count text = %q", got)
	}
	if got := (PurgeCard{Player: ChosenPlayer, Selection: Chosen{House: namedHouse(Dis)}}).Text(); got != "purge a Dis card from a discard pile" {
		t.Errorf("house text = %q", got)
	}
	if got := (PurgeCard{Player: EachPlayer, Selection: Each{House: namedHouse(Untamed), Type: Creature}, GainOwnerAember: true}).Text(); got != "purge each Untamed creature from each player's discard pile. For each card purged this way, its owner gains 1 Æmber" {
		t.Errorf("each text = %q", got)
	}

	// Standalone up-to purge: one zone, purge 2 of 3 (default chooser takes 0).
	g := NewGame("A", "B", 1)
	a := g.Register(testCreature("a", 1), 0)
	b := g.Register(testCreature("b", 1), 0)
	c := g.Register(testCreature("c", 1), 0)
	g.State.Discard[0].add(a)
	g.State.Discard[0].add(b)
	g.State.Discard[0].add(c)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	if !(PurgeCard{Player: ChosenPlayer, Selection: Chosen{Optional: true}, Amount: 2}).resolveGate(
		ctx,
	) {
		t.Error("purging cards should report success")
	}
	if len(g.Purge(0)) != 2 || len(g.Discard(0)) != 1 {
		t.Errorf(
			"purged=%d discard=%d, want 2/1",
			len(g.Purge(0)),
			len(g.Discard(0)),
		)
	}

	// Two zones: the controller picks one (default: their own), which empties
	// before Count is reached.
	g2 := NewGame("A", "B", 1)
	m := g2.Register(testCreature("m", 1), 0)
	n := g2.Register(testCreature("n", 1), 1)
	g2.State.Discard[0].add(m)
	g2.State.Discard[1].add(n)
	ctx2 := &EffectContext{Resolver: g2, Controller: 0}
	PurgeCard{
		Player:    ChosenPlayer,
		Selection: Chosen{Optional: true},
		Amount:    2,
	}.Resolve(
		ctx2,
	)
	if got := g2.Purge(0); len(got) != 1 || got[0] != m {
		t.Errorf("own zone purge = %v, want [m]", got)
	}
	if len(g2.Discard(1)) != 1 {
		t.Error("the opponent's zone should be untouched")
	}

	// Declining ("Done") purges nothing and reports failure.
	g3 := NewGame("A", "B", 1)
	x := g3.Register(testCreature("x", 1), 0)
	g3.State.Discard[0].add(x)
	ctx3 := &EffectContext{Resolver: g3, Controller: 0}
	g3.SetChooser(0, optionPicker{idx: 1}) // options [x, Done] -> idx 1 is Done
	if (PurgeCard{Player: ChosenPlayer, Selection: Chosen{Optional: true}, Amount: 2}).resolveGate(
		ctx3,
	) {
		t.Error("declining should report no purge")
	}
	if len(g3.Purge(0)) != 0 {
		t.Error("declining should purge nothing")
	}

	// Type filter: a creature is purged and a non-creature is left; reports success.
	g4 := NewGame("A", "B", 1)
	crea := g4.Register(testCreature("crea", 3), 1)
	act := g4.Register(NewCard("act", Dis, Tactic, Common), 1)
	g4.State.Discard[1].add(act)
	g4.State.Discard[1].add(crea)
	ctx4 := &EffectContext{Resolver: g4, Controller: 0}
	if !(PurgeCard{Player: ChosenPlayer, Selection: Chosen{Type: Creature}}).resolveGate(
		ctx4,
	) {
		t.Error("purging a creature should report success")
	}
	if got := g4.Purge(1); len(got) != 1 || got[0] != crea {
		t.Errorf("purge = %v, want [crea]", got)
	}
	if len(g4.Discard(1)) != 1 {
		t.Error("the non-creature should be left in the discard")
	}

	// No matching card: nothing is purged and it reports failure.
	g5 := NewGame("A", "B", 1)
	g5.State.Discard[0].add(
		g5.Register(NewCard("act2", Dis, Tactic, Common), 0),
	)
	ctx5 := &EffectContext{Resolver: g5, Controller: 0}
	if (PurgeCard{Player: ChosenPlayer, Selection: Chosen{Type: Creature}}).resolveGate(
		ctx5,
	) {
		t.Error("no creature to purge should report failure")
	}

	// House filter: only the matching-house card is eligible.
	g6 := NewGame("A", "B", 1)
	dis := g6.Register(NewCard("dis", Dis, Creature, Common, WithPower(3)), 1)
	logos := g6.Register(
		NewCard("logos", Logos, Creature, Common, WithPower(3)),
		1,
	)
	g6.State.Discard[1].add(logos)
	g6.State.Discard[1].add(dis)
	ctx6 := &EffectContext{Resolver: g6, Controller: 0}
	if !(PurgeCard{Player: ChosenPlayer, Selection: Chosen{House: namedHouse(Dis)}}).resolveGate(
		ctx6,
	) {
		t.Error("purging a Dis card should report success")
	}
	if got := g6.Purge(1); len(got) != 1 || got[0] != dis {
		t.Errorf("purge = %v, want [dis]", got)
	}
}

// TestPurgeEachFromBothPiles covers the EachPlayer path folded in from Soldiers to
// Flowers: every matching card is purged from both discard piles at once and each
// purged card's owner gains 1 Æmber.
func TestPurgeEachFromBothPiles(t *testing.T) {
	if got := (PurgeCard{Player: EachPlayer, Selection: Each{}}).Text(); got != "purge each card from each player's discard pile" {
		t.Errorf("bare text = %q", got)
	}

	g := NewGame("A", "B", 1)
	mine := g.Register(
		NewCard("mine", Untamed, Creature, Common, WithAemberBonus(1)),
		0,
	)
	minesafe := g.Register(NewCard("minesafe", Mars, Creature, Common), 0)
	theirs := g.Register(NewCard("theirs", Untamed, Creature, Common), 1)
	g.State.Discard[0].add(mine)
	g.State.Discard[0].add(minesafe)
	g.State.Discard[1].add(theirs)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	(PurgeCard{
		Player:          EachPlayer,
		Selection:       Each{House: namedHouse(Untamed), Type: Creature},
		GainOwnerAember: true,
	}).Resolve(ctx)

	if got := g.Purge(0); len(got) != 1 || got[0] != mine {
		t.Errorf("own purge = %v, want [mine]", got)
	}
	if got := g.Purge(1); len(got) != 1 || got[0] != theirs {
		t.Errorf("opponent purge = %v, want [theirs]", got)
	}
	if len(g.Discard(0)) != 1 {
		t.Error("the non-Untamed card should be left in the discard")
	}
	if g.State.Aember[0] != 1 || g.State.Aember[1] != 1 {
		t.Errorf(
			"owner Æmber = %d/%d, want 1/1",
			g.State.Aember[0],
			g.State.Aember[1],
		)
	}
	if got := ctx.Produced.Purged[0] + ctx.Produced.Purged[1]; got != 2 {
		t.Errorf("purged tally = %d, want 2", got)
	}
}

func TestPurgeFromHand(t *testing.T) {
	// validate rejects an unset player or an unset selection.
	if err := (PurgeFromHand{}).validate(); err == nil {
		t.Error("unset player should fail validation")
	}
	if err := (PurgeFromHand{Player: Opponent}).validate(); err == nil {
		t.Error("unset selection should fail validation")
	}
	if err := (PurgeFromHand{Player: Opponent, Selection: Chosen{}}).validate(); err != nil {
		t.Errorf("valid player and selection should pass validation: %v", err)
	}

	// Text and object variants across the three selections.
	if got := (PurgeFromHand{Player: Opponent, Selection: Chosen{House: namedHouse(Sanctum), Optional: true}}).Text(); got != "you may purge a Sanctum card from your opponent's hand" {
		t.Errorf("chosen house text = %q", got)
	}
	if got := (PurgeFromHand{Player: Controller, Selection: Chosen{Optional: true}}).Text(); got != "you may purge a card from your hand" {
		t.Errorf("chosen any-card text = %q", got)
	}
	if got := (PurgeFromHand{Player: Opponent, Selection: Random{}}).Text(); got != "purge a random card from your opponent's hand" {
		t.Errorf("random text = %q", got)
	}
	if got := (PurgeFromHand{Player: Controller, Selection: Each{Type: Creature, House: exceptHouse(Mars)}}).Text(); got != "purge each non-Mars creature from your hand" {
		t.Errorf("each text = %q", got)
	}

	// Resolve: only the Sanctum card is eligible; the default chooser purges it.
	g := NewGame("A", "B", 1)
	sanctum := g.Register(
		NewCard("holy", Sanctum, Creature, Common, WithPower(3)),
		1,
	)
	other := g.Register(
		NewCard("dark", Shadows, Creature, Common, WithPower(3)),
		1,
	)
	g.State.Hand[1].add(sanctum)
	g.State.Hand[1].add(other)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	PurgeFromHand{
		Player:    Opponent,
		Selection: Chosen{House: namedHouse(Sanctum), Optional: true},
	}.Resolve(
		ctx,
	)
	if got := g.Purge(1); len(got) != 1 || got[0] != sanctum {
		t.Errorf("purge = %v, want [sanctum]", got)
	}
	if len(g.Hand(1)) != 1 || g.Hand(1)[0] != other {
		t.Errorf("hand = %v, want [other] (non-Sanctum left)", g.Hand(1))
	}

	// Declining ("Done") purges nothing.
	g2 := NewGame("A", "B", 1)
	holy := g2.Register(
		NewCard("holy", Sanctum, Creature, Common, WithPower(3)),
		1,
	)
	g2.State.Hand[1].add(holy)
	ctx2 := &EffectContext{Resolver: g2, Controller: 0}
	g2.SetChooser(
		0,
		optionPicker{idx: 1},
	) // options [holy, Done] -> idx 1 is Done
	PurgeFromHand{
		Player:    Opponent,
		Selection: Chosen{House: namedHouse(Sanctum), Optional: true},
	}.Resolve(
		ctx2,
	)
	if len(g2.Purge(1)) != 0 || len(g2.Hand(1)) != 1 {
		t.Error("declining should purge nothing")
	}

	// No matching card: nothing happens.
	g3 := NewGame("A", "B", 1)
	g3.State.Hand[1].add(
		g3.Register(
			NewCard("dark", Shadows, Creature, Common, WithPower(3)),
			1,
		),
	)
	ctx3 := &EffectContext{Resolver: g3, Controller: 0}
	PurgeFromHand{
		Player:    Opponent,
		Selection: Chosen{House: namedHouse(Sanctum), Optional: true},
	}.Resolve(
		ctx3,
	)
	if len(g3.Purge(1)) != 0 {
		t.Error("no matching card should purge nothing")
	}

	// Under a May a Chosen purge is offered as its own single optional choice.
	if !(PurgeFromHand{Player: Opponent, Selection: Chosen{Optional: true}}).declinable() {
		t.Error("an Optional Chosen purge should be declinable")
	}
	g4 := NewGame("A", "B", 1)
	card4 := g4.Register(
		NewCard("dark", Shadows, Creature, Common, WithPower(3)),
		1,
	)
	g4.State.Hand[1].add(card4)
	ctx4 := &EffectContext{Resolver: g4, Controller: 0}
	May{
		Do: PurgeFromHand{Player: Opponent, Selection: Chosen{Optional: true}},
	}.Resolve(
		ctx4,
	)
	if got := g4.Purge(1); len(got) != 1 || got[0] != card4 {
		t.Errorf("May purge = %v, want [card4]", got)
	}

	// Default: no "you may", cannot be declined, and forces the purge when a
	// card is in hand; an empty hand still purges nothing (Greater Oxtet).
	if got := (PurgeFromHand{Player: Controller, Selection: Chosen{}}).Text(); got != "purge a card from your hand" {
		t.Errorf("mandatory text = %q", got)
	}
	if (PurgeFromHand{Player: Controller, Selection: Chosen{}}).declinable() {
		t.Error("a mandatory purge should not be declinable")
	}
	g5 := NewGame("A", "B", 1)
	forced := g5.Register(
		NewCard("fodder", Shadows, Creature, Common, WithPower(3)),
		0,
	)
	g5.State.Hand[0].add(forced)
	PurgeFromHand{
		Player:    Controller,
		Selection: Chosen{},
	}.Resolve(
		&EffectContext{Resolver: g5, Controller: 0},
	)
	if got := g5.Purge(0); len(got) != 1 || got[0] != forced {
		t.Errorf("mandatory purge = %v, want [forced]", got)
	}
	g6 := NewGame("A", "B", 1)
	PurgeFromHand{
		Player:    Controller,
		Selection: Chosen{},
	}.Resolve(
		&EffectContext{Resolver: g6, Controller: 0},
	)
	if len(g6.Purge(0)) != 0 {
		t.Error("mandatory purge with empty hand should purge nothing")
	}
}

func TestPurgeFromHandRandom(t *testing.T) {
	// Resolve: the sole hand card is purged, records the tally, and gates a Then.
	g := NewGame("A", "B", 1)
	only := g.Register(
		NewCard("dark", Shadows, Creature, Common, WithPower(3)),
		1,
	)
	g.State.Hand[1].add(only)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	if !(PurgeFromHand{Player: Opponent, Selection: Random{}}).resolveGate(
		ctx,
	) {
		t.Error("purging a card should report true")
	}
	if got := g.Purge(1); len(got) != 1 || got[0] != only {
		t.Errorf("purge = %v, want [only]", got)
	}
	if got := ctx.Produced.Purged[0] + ctx.Produced.Purged[1]; got != 1 {
		t.Errorf("tally = %d, want 1", got)
	}

	// An empty hand purges nothing and reports false.
	g2 := NewGame("A", "B", 1)
	ctx2 := &EffectContext{Resolver: g2, Controller: 0}
	if (PurgeFromHand{Player: Opponent, Selection: Random{}}).resolveGate(
		ctx2,
	) {
		t.Error("empty hand should report false")
	}
	if len(g2.Purge(1)) != 0 {
		t.Error("empty hand should purge nothing")
	}
}

func TestPurgeFromHandChosenCreature(t *testing.T) {
	// A mandatory Chosen restricted to creatures is Custom Virus's "purge a
	// creature from your hand", which puts the purged card in context (ctx.It).
	e := PurgeFromHand{Player: Controller, Selection: Chosen{Type: Creature}}
	if e.Text() != "purge a creature from your hand" {
		t.Errorf("text = %q", e.Text())
	}

	t.Run(
		"purges a chosen hand creature and sets it in context",
		func(t *testing.T) {
			g := NewGame("A", "B", 1)
			kin := g.AddToHand(testCreature("kin", 3, WithTraits(Beast)), 0)
			ctx := &EffectContext{Resolver: g, Controller: 0}

			e.Resolve(ctx)

			if !g.State.Purge[0].contains(kin) {
				t.Error("the chosen creature should be purged from hand")
			}
			if !ctx.HasIt || ctx.It != kin {
				t.Errorf(
					"context card = %d (has %v), want the purged creature %d",
					ctx.It,
					ctx.HasIt,
					kin,
				)
			}
		},
	)

	t.Run(
		"with no creature in hand, does nothing and leaves It unset",
		func(t *testing.T) {
			g := NewGame("A", "B", 1)
			g.AddToHand(
				NewCard("Action", Brobnar, Tactic, Common),
				0,
			) // not a creature
			ctx := &EffectContext{Resolver: g, Controller: 0}
			e.Resolve(ctx)
			if len(g.Purge(0)) != 0 {
				t.Error("nothing should be purged when no creature is in hand")
			}
			if ctx.HasIt {
				t.Error("It should stay unset when nothing is purged")
			}
		},
	)

	t.Run("a declined choice purges nothing", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.AddToHand(testCreature("kin1", 3, WithTraits(Beast)), 0)
		g.AddToHand(testCreature("kin2", 3, WithTraits(Beast)), 0)
		g.SetChooser(0, orderRejectChooser{})
		ctx := &EffectContext{Resolver: g, Controller: 0}
		e.Resolve(ctx)
		if len(g.Purge(0)) != 0 {
			t.Error("declining should purge nothing")
		}
		if ctx.HasIt {
			t.Error("It should stay unset when the choice is declined")
		}
	})
}

func TestCardsPurgedCount(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 3), 0)
	g.AddToBattleline(NewCard("a", Shadows, Creature, Common, WithPower(3)), 0)
	g.AddToBattleline(NewCard("b", Shadows, Creature, Common, WithPower(3)), 0)
	g.AddToBattleline(
		NewCard("off-house", Mars, Creature, Common, WithPower(3)),
		0,
	)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	if got := (CardsPurged{}).CountText(); got != "card purged this way" {
		t.Errorf("count text = %q", got)
	}
	if got := (CardsPurged{Type: Creature}).CountText(); got != "creature purged this way" {
		t.Errorf("creature count text = %q", got)
	}
	PurgeCreature{
		Target: Target{Kind: TargetEachFriendlyCreature}.House(namedHouse(Shadows)),
	}.Resolve(ctx)

	if got := (CardsPurged{}).Value(ctx); got != 2 {
		t.Errorf("purged = %d, want 2", got)
	}
}

// A "you may purge a neighboring creature" is one clickable creature, so it is
// asked declinably and asks nothing when the source has no neighbor (Buzzle at a
// flank).
func TestMayPurgeCreatureDeclinable(t *testing.T) {
	neighboring := func() Target { return Target{Kind: TargetChosenCreature}.Neighboring() }

	if !(PurgeCreature{Target: neighboring()}).declinable() {
		t.Fatal("a chosen-target PurgeCreature should be declinable")
	}

	t.Run("accepted purges the clicked neighbor", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		ch := &cardDecliner{}
		g.SetChooser(0, ch)
		src := g.AddToBattleline(testCreature("src", 3), 0)
		neighbor := g.AddToBattleline(testCreature("neighbor", 3), 0)
		ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

		May{Do: PurgeCreature{Target: neighboring()}}.Resolve(ctx)

		if ch.asked != 1 {
			t.Errorf("declinable prompts = %d, want 1", ch.asked)
		}
		if !g.State.Purge[0].contains(neighbor) {
			t.Error("the clicked neighbor should have been purged")
		}
	})

	t.Run("declined purges nothing", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.SetChooser(0, &cardDecliner{decline: true})
		src := g.AddToBattleline(testCreature("src", 3), 0)
		neighbor := g.AddToBattleline(testCreature("neighbor", 3), 0)
		ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

		May{Do: PurgeCreature{Target: neighboring()}}.Resolve(ctx)

		if g.State.Purge[0].contains(neighbor) {
			t.Error("a declined May should purge nothing")
		}
	})

	t.Run("no neighbor is vacuous and asks nothing", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		ch := &cardDecliner{}
		g.SetChooser(0, ch)
		src := g.AddToBattleline(testCreature("src", 3), 0)
		ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

		if !(PurgeCreature{Target: neighboring()}).vacuous(ctx) {
			t.Fatal("a lone creature's neighboring target should be vacuous")
		}
		May{Do: PurgeCreature{Target: neighboring()}}.Resolve(ctx)
		if ch.asked != 0 {
			t.Errorf("prompts with no neighbor = %d, want 0", ch.asked)
		}
	})
}

func TestPurgedAemberBonusCount(t *testing.T) {
	if got := (PurgedAemberBonus{}).CountText(); got != "the total Æmber bonus of the purged cards" {
		t.Errorf("count text = %q", got)
	}

	g := NewGame("A", "B", 1)
	two := g.Register(
		NewCard("two", Shadows, Creature, Common, WithAemberBonus(2)),
		0,
	)
	one := g.Register(
		NewCard("one", Shadows, Creature, Common, WithAemberBonus(1)),
		0,
	)
	none := g.Register(NewCard("none", Shadows, Creature, Common), 0)
	g.State.Discard[0].add(two)
	g.State.Discard[0].add(one)
	g.State.Discard[0].add(none)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	// Default chooser purges the first two (bonuses 2 and 1).
	PurgeCard{
		Player:    ChosenPlayer,
		Selection: Chosen{},
		Amount:    2,
	}.Resolve(
		ctx,
	)
	if got := (PurgedAemberBonus{}).Value(ctx); got != 3 {
		t.Errorf("purged bonus = %d, want 3", got)
	}
}

func TestPurgeSource(t *testing.T) {
	if got := (PurgeSource{}).Text(); got != "purge "+SelfName {
		t.Errorf("text = %q", got)
	}

	t.Run("purges an in-play source from play", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		src := g.AddToBattleline(testCreature("src", 3), 0)

		PurgeSource{}.Resolve(
			&EffectContext{Resolver: g, Controller: 0, Source: src},
		)

		if !g.State.Purge[0].contains(src) {
			t.Error("an in-play source should be purged from play")
		}
	})

	t.Run(
		"sets a resolving action aside instead of discarding it",
		func(t *testing.T) {
			g := started(t)
			idx := int(g.State.Hand[0].Count)
			id := g.AddToHand(
				NewCard(
					"Self Purge", Brobnar, Tactic, Common,
					WithAbility(TriggerAfterPlay, PurgeSource{}),
				),
				0,
			)

			if err := g.PlayAction(0, idx); err != nil {
				t.Fatalf("PlayAction: %v", err)
			}

			if !g.State.Purge[0].contains(id) {
				t.Error(
					"a self-purging action should be set aside in the purge pile",
				)
			}
			if g.State.Discard[0].contains(id) {
				t.Error(
					"a self-purging action should not go to the discard pile",
				)
			}
		},
	)
}
