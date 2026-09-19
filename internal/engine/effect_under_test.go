package engine

import "testing"

func TestPutUnderFromHandText(t *testing.T) {
	cases := []struct {
		name   string
		effect PutUnderFromHand
		want   string
	}{
		{"faceup", PutUnderFromHand{}, "put a card from your hand faceup under {self}"},
		{
			"facedown",
			PutUnderFromHand{FaceDown: true},
			"put a card from your hand facedown under {self}",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.effect.Text(); got != tc.want {
				t.Errorf("text = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestPutUnderFromHandChoosesAndAttaches(t *testing.T) {
	g := started(t)
	host := g.AddArtifact(NewCard("Host", Brobnar, Artifact, Common), 0)
	buried := g.AddToHand(NewCard("Buried", Brobnar, Creature, Common, WithPower(2)), 0)

	PutUnderFromHand{FaceDown: true}.Resolve(
		&EffectContext{
			Resolver:   g,
			Controller: 0,
			Source:     host,
		},
	)

	if got := g.Under(host); len(got) != 1 || got[0] != buried {
		t.Errorf("under = %v, want [%d]", got, buried)
	}
	if !g.UnderFaceDown(buried) {
		t.Error("the buried card should be facedown")
	}
}

func TestPutUnderFromHandWithEmptyHand(t *testing.T) {
	g := started(t)
	host := g.AddArtifact(NewCard("Host", Brobnar, Artifact, Common), 0)

	PutUnderFromHand{}.Resolve(&EffectContext{
		Resolver:   g,
		Controller: 0,
		Source:     host,
	})

	if got := g.Under(host); len(got) != 0 {
		t.Errorf("under = %v, want none", got)
	}
}

func TestPutUnderFromHandDeclined(t *testing.T) {
	g := started(t)
	host := g.AddArtifact(NewCard("Host", Brobnar, Artifact, Common), 0)
	g.AddToHand(NewCard("First", Brobnar, Creature, Common, WithPower(2)), 0)
	g.AddToHand(NewCard("Second", Brobnar, Creature, Common, WithPower(2)), 0)
	g.SetChooser(0, orderRejectChooser{})

	PutUnderFromHand{}.Resolve(&EffectContext{
		Resolver:   g,
		Controller: 0,
		Source:     host,
	})

	if got := g.Under(host); len(got) != 0 {
		t.Errorf("under = %v, want none after declining", got)
	}
}

func TestPlayCardUnderText(t *testing.T) {
	if got := (PlayCardUnder{}).Text(); got != "play the card under {self}" {
		t.Errorf("text = %q, want %q", got, "play the card under {self}")
	}
}

func TestPlayCardUnderWithOneCandidate(t *testing.T) {
	g := started(t)
	host := g.AddArtifact(NewCard("Host", Brobnar, Artifact, Common), 0)
	buried := g.Register(NewCard("Buried", Brobnar, Creature, Common, WithPower(2)), 0)
	g.AttachUnder(host, buried, true)

	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
		Source:     host,
	}
	PlayCardUnder{}.Resolve(ctx)

	if got := g.Battleline(0); len(got) != 1 || got[0] != buried {
		t.Errorf("battleline = %v, want the played creature %d", got, buried)
	}
	if !ctx.HasIt || ctx.It != buried {
		t.Errorf("ctx.It = %d (has=%v), want %d", ctx.It, ctx.HasIt, buried)
	}
}

func TestPlayCardUnderChoosesAmongSeveral(t *testing.T) {
	g := started(t)
	host := g.AddArtifact(NewCard("Host", Brobnar, Artifact, Common), 0)
	first := g.Register(NewCard("First", Brobnar, Creature, Common, WithPower(1)), 0)
	second := g.Register(NewCard("Second", Brobnar, Creature, Common, WithPower(1)), 0)
	g.AttachUnder(host, first, true)
	g.AttachUnder(host, second, true)
	g.SetChooser(0, &idQueueChooser{ids: []LocalID{second}})

	PlayCardUnder{}.Resolve(&EffectContext{
		Resolver:   g,
		Controller: 0,
		Source:     host,
	})

	if got := g.Battleline(0); len(got) != 1 || got[0] != second {
		t.Errorf("battleline = %v, want the chosen creature %d", got, second)
	}
	if got := g.Under(host); len(got) != 1 || got[0] != first {
		t.Errorf("under = %v, want the remaining card [%d]", got, first)
	}
}

func TestPlayCardUnderDeclined(t *testing.T) {
	g := started(t)
	host := g.AddArtifact(NewCard("Host", Brobnar, Artifact, Common), 0)
	first := g.Register(NewCard("First", Brobnar, Creature, Common, WithPower(1)), 0)
	second := g.Register(NewCard("Second", Brobnar, Creature, Common, WithPower(1)), 0)
	g.AttachUnder(host, first, true)
	g.AttachUnder(host, second, true)
	g.SetChooser(0, orderRejectChooser{})

	PlayCardUnder{}.Resolve(&EffectContext{
		Resolver:   g,
		Controller: 0,
		Source:     host,
	})

	if got := len(g.Battleline(0)); got != 0 {
		t.Errorf("battleline holds %d creatures, want none after declining", got)
	}
}

func TestPlayCardUnderWithNoCandidate(t *testing.T) {
	g := started(t)
	host := g.AddArtifact(NewCard("Host", Brobnar, Artifact, Common), 0)

	PlayCardUnder{}.Resolve(&EffectContext{
		Resolver:   g,
		Controller: 0,
		Source:     host,
	})

	if got := len(g.Battleline(0)); got != 0 {
		t.Errorf("battleline holds %d creatures, want none", got)
	}
}

func TestGraftText(t *testing.T) {
	got := (Graft{Target: Target{Kind: TargetChosenCreature}}).Text()
	if got != "graft a creature from play" {
		t.Errorf("text = %q, want %q", got, "graft a creature from play")
	}
}

func TestGraftValidate(t *testing.T) {
	if (Graft{}).validate() == nil {
		t.Error("Graft with no target should fail validation")
	}
	if err := (Graft{Target: Target{Kind: TargetChosenCreature}}).validate(); err != nil {
		t.Errorf("valid Graft should pass validation, got %v", err)
	}
}

// TestGraftResolvePlacesTargetUnderSource grafts a creature in play onto the
// resolving card: it leaves play, lands faceup under the source, and logs the
// graft.
func TestGraftResolvePlacesTargetUnderSource(t *testing.T) {
	g := started(t)
	host := g.AddArtifact(NewCard("Host", Brobnar, Artifact, Common), 0)
	victim := g.Register(NewCard("Victim", Brobnar, Creature, Common, WithPower(3)), 0)
	g.putIntoPlay(victim, 0)

	Graft{Target: Target{Kind: TargetChosenCreature}}.Resolve(
		&EffectContext{
			Resolver:   g,
			Controller: 0,
			Source:     host,
		},
	)

	if g.inPlay(victim) {
		t.Error("the grafted creature should have left play")
	}
	if got := g.Under(host); len(got) != 1 || got[0] != victim {
		t.Errorf("under = %v, want [%d]", got, victim)
	}
	if g.UnderFaceDown(victim) {
		t.Error("graft places the card faceup")
	}
	last := g.Log[len(g.Log)-1].Entry
	gr, ok := last.(CardGrafted)
	if !ok {
		t.Fatalf("last log entry = %T, want CardGrafted", last)
	}
	if gr.Card != victim || gr.Host != host {
		t.Errorf("log entry = %+v, want Card=%d Host=%d", gr, victim, host)
	}
}

func TestPutUnderIntoPlayText(t *testing.T) {
	got := (PutUnderIntoPlay{}).Text()
	want := "put each card under {self} into play under its owner's control"
	if got != want {
		t.Errorf("text = %q, want %q", got, want)
	}
}

// TestPutUnderIntoPlayResolveReturnsToOwners puts every card under the source
// into play under its own owner's control, whichever player that is.
func TestPutUnderIntoPlayResolveReturnsToOwners(t *testing.T) {
	g := started(t)
	host := g.AddArtifact(NewCard("Host", Brobnar, Artifact, Common), 0)
	mine := g.Register(NewCard("Mine", Brobnar, Creature, Common, WithPower(2)), 0)
	theirs := g.Register(NewCard("Theirs", Brobnar, Creature, Common, WithPower(2)), 1)
	g.AttachUnder(host, mine, false)
	g.AttachUnder(host, theirs, false)

	PutUnderIntoPlay{}.Resolve(&EffectContext{
		Resolver:   g,
		Controller: 0,
		Source:     host,
	})

	if got := g.Under(host); len(got) != 0 {
		t.Errorf("under = %v, want empty", got)
	}
	if g.controller(mine) != 0 || !g.inPlay(mine) {
		t.Errorf(
			"mine controller = %d in-play = %v, want 0/true",
			g.controller(mine),
			g.inPlay(mine),
		)
	}
	if g.controller(theirs) != 1 || !g.inPlay(theirs) {
		t.Errorf(
			"theirs controller = %d in-play = %v, want 1/true",
			g.controller(theirs),
			g.inPlay(theirs),
		)
	}
}

// TestArchiveCardUnderText renders the effect.
func TestArchiveCardUnderText(t *testing.T) {
	if got := (ArchiveCardUnder{}).Text(); got != "archive the card under {self}" {
		t.Errorf("text = %q, want %q", got, "archive the card under {self}")
	}
}

// TestArchiveCardUnderResolveMovesToOwnerArchives archives every card under the
// source into its own owner's archives, whichever player that is.
func TestArchiveCardUnderResolveMovesToOwnerArchives(t *testing.T) {
	g := started(t)
	host := g.AddArtifact(NewCard("Host", Brobnar, Artifact, Common), 0)
	mine := g.Register(NewCard("Mine", Brobnar, Creature, Common, WithPower(2)), 0)
	theirs := g.Register(NewCard("Theirs", Brobnar, Creature, Common, WithPower(2)), 1)
	g.AttachUnder(host, mine, true)
	g.AttachUnder(host, theirs, true)

	ArchiveCardUnder{}.Resolve(&EffectContext{
		Resolver:   g,
		Controller: 0,
		Source:     host,
	})

	if got := g.Under(host); len(got) != 0 {
		t.Errorf("under = %v, want empty", got)
	}
	if got := g.Archives(0); len(got) != 1 || got[0] != mine {
		t.Errorf("player 0 archives = %v, want [%d]", got, mine)
	}
	if got := g.Archives(1); len(got) != 1 || got[0] != theirs {
		t.Errorf("player 1 archives = %v, want [%d]", got, theirs)
	}
}

// TestPutUnderFromHandTypeText renders the "tactic card" noun a Tactic filter
// gives (Infomancer, Memolith graft a Tactic).
func TestPutUnderFromHandTypeText(t *testing.T) {
	got := (PutUnderFromHand{Type: Tactic}).Text()
	want := "put a tactic card from your hand faceup under {self}"
	if got != want {
		t.Errorf("text = %q, want %q", got, want)
	}
}

// TestPutUnderFromHandTypeOnlyOffersMatchingType grafts only a Tactic when
// Type is Tactic, leaving a creature in hand untouched even though it is the only
// other hand card.
func TestPutUnderFromHandTypeOnlyOffersMatchingType(t *testing.T) {
	g := started(t)
	host := g.AddArtifact(NewCard("Host", Brobnar, Artifact, Common), 0)
	g.AddToHand(NewCard("Creature", Brobnar, Creature, Common, WithPower(2)), 0)
	tactic := g.AddToHand(NewCard("Tactic", Brobnar, Tactic, Common), 0)

	PutUnderFromHand{Type: Tactic}.Resolve(
		&EffectContext{
			Resolver:   g,
			Controller: 0,
			Source:     host,
		},
	)

	if got := g.Under(host); len(got) != 1 || got[0] != tactic {
		t.Errorf("under = %v, want the Tactic [%d]", got, tactic)
	}
	if g.UnderFaceDown(tactic) {
		t.Error("a graft from hand places the card faceup")
	}
}

// TestTriggerGraftedPlayEffectText renders the effect.
func TestTriggerGraftedPlayEffectText(t *testing.T) {
	got := (TriggerGraftedPlayEffect{}).Text()
	want := "trigger the play effect of a Tactic grafted onto {self}"
	if got != want {
		t.Errorf("text = %q, want %q", got, want)
	}
}

// TestTriggerGraftedPlayEffectFiresAndStaysGrafted fires the Play effect of a
// grafted action and leaves the action grafted under the source — it does not move
// to play or discard.
func TestTriggerGraftedPlayEffectFiresAndStaysGrafted(t *testing.T) {
	g := started(t)
	host := g.AddArtifact(NewCard("Host", Brobnar, Artifact, Common), 0)
	gift := g.Register(NewCard("Gift", Brobnar, Tactic, Common,
		WithAbility(TriggerAfterPlay, GainAember{
			Amount: 2,
			Player: Controller,
		})), 0)
	g.AttachUnder(host, gift, false)

	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
		Source:     host,
	}
	TriggerGraftedPlayEffect{}.Resolve(ctx)

	if g.Aember(0) != 2 {
		t.Errorf("Æmber = %d, want 2 from the grafted action's play effect", g.Aember(0))
	}
	if got := g.Under(host); len(got) != 1 || got[0] != gift {
		t.Errorf("under = %v, want the action still grafted [%d]", got, gift)
	}
	if g.inPlay(gift) {
		t.Error("the grafted action should not enter play")
	}
	if !ctx.HasIt || ctx.It != gift {
		t.Errorf("ctx.It = %d (has=%v), want the grafted action %d", ctx.It, ctx.HasIt, gift)
	}
}

// TestTriggerGraftedPlayEffectChoosesAmongSeveral triggers the grafted action the
// controller chooses when several are grafted, leaving both under the source.
func TestTriggerGraftedPlayEffectChoosesAmongSeveral(t *testing.T) {
	g := started(t)
	host := g.AddArtifact(NewCard("Host", Brobnar, Artifact, Common), 0)
	first := g.Register(NewCard("First", Brobnar, Tactic, Common,
		WithAbility(TriggerAfterPlay, GainAember{
			Amount: 1,
			Player: Controller,
		})), 0)
	second := g.Register(NewCard("Second", Brobnar, Tactic, Common,
		WithAbility(TriggerAfterPlay, GainAember{
			Amount: 3,
			Player: Controller,
		})), 0)
	g.AttachUnder(host, first, false)
	g.AttachUnder(host, second, false)
	g.SetChooser(0, &idQueueChooser{ids: []LocalID{second}})

	TriggerGraftedPlayEffect{}.Resolve(&EffectContext{
		Resolver:   g,
		Controller: 0,
		Source:     host,
	})

	if g.Aember(0) != 3 {
		t.Errorf("Æmber = %d, want 3 from the chosen action", g.Aember(0))
	}
	if got := g.Under(host); len(got) != 2 {
		t.Errorf("under = %v, want both actions still grafted", got)
	}
}

// TestTriggerGraftedPlayEffectVacuousWithNothingGrafted does nothing when no
// action is grafted under the source.
func TestTriggerGraftedPlayEffectVacuousWithNothingGrafted(t *testing.T) {
	g := started(t)
	host := g.AddArtifact(NewCard("Host", Brobnar, Artifact, Common), 0)

	TriggerGraftedPlayEffect{}.Resolve(&EffectContext{
		Resolver:   g,
		Controller: 0,
		Source:     host,
	})

	if g.Aember(0) != 0 {
		t.Errorf("Æmber = %d, want 0 with nothing grafted", g.Aember(0))
	}
}

// TestTriggerGraftedPlayEffectSkipsFaceDown ignores a facedown card under the
// source: only a faceup grafted action is triggered.
func TestTriggerGraftedPlayEffectSkipsFaceDown(t *testing.T) {
	g := started(t)
	host := g.AddArtifact(NewCard("Host", Brobnar, Artifact, Common), 0)
	hidden := g.Register(NewCard("Hidden", Brobnar, Tactic, Common,
		WithAbility(TriggerAfterPlay, GainAember{
			Amount: 5,
			Player: Controller,
		})), 0)
	g.AttachUnder(host, hidden, true)

	TriggerGraftedPlayEffect{}.Resolve(&EffectContext{
		Resolver:   g,
		Controller: 0,
		Source:     host,
	})

	if g.Aember(0) != 0 {
		t.Errorf("Æmber = %d, want 0 — a facedown card is not a graft", g.Aember(0))
	}
}

// TestTriggerGraftedPlayEffectDeclined fires nothing when the controller declines
// the choice among several grafted actions.
func TestTriggerGraftedPlayEffectDeclined(t *testing.T) {
	g := started(t)
	host := g.AddArtifact(NewCard("Host", Brobnar, Artifact, Common), 0)
	first := g.Register(NewCard("First", Brobnar, Tactic, Common,
		WithAbility(TriggerAfterPlay, GainAember{
			Amount: 1,
			Player: Controller,
		})), 0)
	second := g.Register(NewCard("Second", Brobnar, Tactic, Common,
		WithAbility(TriggerAfterPlay, GainAember{
			Amount: 3,
			Player: Controller,
		})), 0)
	g.AttachUnder(host, first, false)
	g.AttachUnder(host, second, false)
	g.SetChooser(0, orderRejectChooser{})

	TriggerGraftedPlayEffect{}.Resolve(&EffectContext{
		Resolver:   g,
		Controller: 0,
		Source:     host,
	})

	if g.Aember(0) != 0 {
		t.Errorf("Æmber = %d, want 0 after declining the choice", g.Aember(0))
	}
}

// TestPutUnderFromHandNonTacticTypeText renders the general "<type> card" noun for
// a non-Tactic type filter.
func TestPutUnderFromHandNonTacticTypeText(t *testing.T) {
	got := (PutUnderFromHand{Type: Creature}).Text()
	want := "put a creature card from your hand faceup under {self}"
	if got != want {
		t.Errorf("text = %q, want %q", got, want)
	}
}
