package engine

import "testing"

func TestPlayFromText(t *testing.T) {
	cases := []struct {
		name   string
		effect PlayFrom
		want   string
	}{
		{"any card", PlayFrom{From: Hand}, "play a card"},
		{
			"excluded house",
			PlayFrom{From: Hand, House: HouseMatcher{Kind: MatchExceptHouse, House: Logos}},
			"play a non-Logos card",
		},
		{
			"named house",
			PlayFrom{From: Hand, House: HouseMatcher{Kind: MatchNamedHouse, House: Mars}},
			"play a Mars card",
		},
		{"typed", PlayFrom{From: Hand, Types: CardTypesOf(Creature)}, "play a creature"},
		{
			"house and type",
			PlayFrom{
				From:  Hand,
				House: HouseMatcher{Kind: MatchNamedHouse, House: Untamed},
				Types: CardTypesOf(Artifact),
			},
			"play an Untamed artifact",
		},
		{
			"opponent's discard pile",
			PlayFrom{From: Discard, Player: Opponent, Types: CardTypesOf(Tactic)},
			"play a tactic from your opponent's discard pile",
		},
		{
			"opponent's hand",
			PlayFrom{From: Hand, Player: Opponent},
			"play a card from your opponent's hand",
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

func TestPlayFromValidate(t *testing.T) {
	if err := (PlayFrom{From: Hand, House: HouseMatcher{Kind: MatchExceptHouse}}).validate(); err == nil {
		t.Error("an except-house matcher without a house should not validate")
	}
	if err := (PlayFrom{From: Hand, House: HouseMatcher{Kind: MatchExceptHouse, House: Logos}}).validate(); err != nil {
		t.Errorf("validate = %v, want nil", err)
	}
	if err := (PlayFrom{From: Hand, Player: Opponent}).validate(); err != nil {
		t.Errorf("the opponent's hand may be played from (Lateral Shift): %v", err)
	}
	if err := (PlayFrom{From: Archives, Player: Opponent}).validate(); err == nil {
		t.Error("only the opponent's hand or discard pile may be played from, not their archives")
	}
}

// TestPlayFromOpponentDiscard covers Mimicry: a Tactic played out of the
// opponent's discard pile resolves under the controller's control, counts as the
// controller's own play, and returns to the top of its owner's discard pile.
func TestPlayFromOpponentDiscard(t *testing.T) {
	g := started(t)
	copied := g.AddToDiscard(NewCard("Copied", Logos, Tactic, Common,
		WithAbility(TriggerAfterPlay, GainAember{Player: Controller, Amount: 2})), 1)
	buried := g.AddToDiscard(NewCard("Buried", Logos, Creature, Common, WithPower(1)), 1)

	PlayFrom{From: Discard, Player: Opponent, Types: CardTypesOf(Tactic)}.Resolve(
		&EffectContext{Resolver: g, Controller: 0},
	)

	// The Play: ability resolved under player 0's control, so player 0 — not the
	// owner — gained the Æmber.
	if got := g.State.Aember[0]; got != 2 {
		t.Errorf("player 0 Æmber = %d, want 2 from playing the copied action", got)
	}
	if got := g.State.Aember[1]; got != 0 {
		t.Errorf("player 1 Æmber = %d, want 0 — the owner does not gain", got)
	}
	// It counts as player 0's own play (Ember Imp can limit it).
	if played := g.State.PlayedThisTurn[0].slice(); len(played) != 1 || played[0] != copied {
		t.Errorf("player 0's plays = %v, want just the copied action %d", played, copied)
	}
	// The action returns to the top of its owner's (player 1's) discard pile, above
	// the card that was already there.
	if g.State.Discard[0].contains(copied) {
		t.Error("the copied action should not land in the player's own discard pile")
	}
	discard := g.Discard(1)
	if len(discard) != 2 || discard[len(discard)-1] != copied || discard[0] != buried {
		t.Errorf("player 1 discard = %v, want [%d %d]", discard, buried, copied)
	}
}

// TestPlayFromBindsIt covers Imperial Road: PlayFrom binds the played card in
// context (ctx.It) so a chained Stun via Then stuns it.
func TestPlayFromBindsIt(t *testing.T) {
	g := started(t)
	dino := g.AddToHand(NewCard("Dino", Brobnar, Creature, Common, WithPower(3)), 0)

	ctx := &EffectContext{Resolver: g, Controller: 0}
	Then{
		First:  PlayFrom{From: Hand, Types: CardTypesOf(Creature)},
		Result: Stun{Target: Target{Kind: TargetTriggeringCreature}},
	}.Resolve(ctx)

	if !ctx.HasIt || ctx.It != dino {
		t.Fatalf("ctx.It = %d (has %v), want played creature %d", ctx.It, ctx.HasIt, dino)
	}
	if line := g.Battleline(0); len(line) != 1 || line[0] != dino {
		t.Fatalf("battleline = %v, want the played creature %d", line, dino)
	}
	if !g.Stunned(dino) {
		t.Error("the played creature should be stunned by the chained Stun")
	}
}

// TestPlayFromGateFalseWithNoCandidate confirms the gate reports false when no
// eligible card is in the pile, so a Then's follow-up never runs.
func TestPlayFromGateFalseWithNoCandidate(t *testing.T) {
	g := started(t)
	bystander := g.AddToBattleline(NewCard("Bystander", Brobnar, Creature, Common, WithPower(2)), 0)

	Then{
		First:  PlayFrom{From: Hand, Types: CardTypesOf(Creature)},
		Result: Stun{Target: Target{Kind: TargetEachFriendlyCreature}},
	}.Resolve(&EffectContext{Resolver: g, Controller: 0})

	if g.Stunned(bystander) {
		t.Error("no card was played, so the gate is false and the follow-up Stun must not run")
	}
}

// TestPlayFromOpponentHand covers Lateral Shift: a player plays a creature out of
// the opponent's hand as their own — it enters the player's battleline under their
// control and leaves the opponent's hand, while its owner stays the opponent.
func TestPlayFromOpponentHand(t *testing.T) {
	g := started(t)
	foreign := g.AddToHand(NewCard("Borrowed", Brobnar, Creature, Common, WithPower(4)), 1)

	PlayFrom{From: Hand, Player: Opponent}.Resolve(&EffectContext{Resolver: g, Controller: 0})

	if line := g.Battleline(0); len(line) != 1 || line[0] != foreign {
		t.Fatalf("player 0 battleline = %v, want the borrowed creature %d", line, foreign)
	}
	for _, id := range g.Hand(1) {
		if id == foreign {
			t.Error("the borrowed creature should have left the opponent's hand")
		}
	}
	if got := g.Controller(foreign); got != 0 {
		t.Errorf("controller = %d, want player 0 who played it", got)
	}
	if got := g.Owner(foreign); got != 1 {
		t.Errorf("owner = %d, want player 1 whose hand it came from", got)
	}
}

// TestPlayFromOpponentHandIgnoresUnknownCard covers the guard: asking to play a
// card that is not in the opponent's hand does nothing.
func TestPlayFromOpponentHandIgnoresUnknownCard(t *testing.T) {
	g := started(t)
	before := len(g.Battleline(0))

	g.PlayFromOpponentHand(0, LocalID(200))

	if got := len(g.Battleline(0)); got != before {
		t.Errorf("battleline changed to %d cards, want %d", got, before)
	}
}

func TestPlayFromPlaysAChosenCard(t *testing.T) {
	g := started(t) // Brobnar is the active house.
	off := g.AddToHand(NewCard("Off House", Logos, Creature, Common, WithPower(2)), 0)
	g.AddToHand(NewCard("On House", Brobnar, Creature, Common, WithPower(2)), 0)

	PlayFrom{
		From:  Hand,
		House: HouseMatcher{Kind: MatchExceptHouse, House: Brobnar},
	}.Resolve(
		&EffectContext{Resolver: g, Controller: 0},
	)

	if got := g.Battleline(0); len(got) != 1 || got[0] != off {
		t.Errorf("battleline = %v, want the off-house creature %d", got, off)
	}
}

func TestPlayFromWithNoCandidate(t *testing.T) {
	g := started(t)
	g.AddToHand(NewCard("On House", Brobnar, Creature, Common, WithPower(2)), 0)

	PlayFrom{
		From:  Hand,
		House: HouseMatcher{Kind: MatchExceptHouse, House: Brobnar},
	}.Resolve(
		&EffectContext{Resolver: g, Controller: 0},
	)

	if got := len(g.Battleline(0)); got != 0 {
		t.Errorf("battleline holds %d creatures, want none", got)
	}
}

func TestPlayFromDeclined(t *testing.T) {
	g := started(t)
	g.AddToHand(NewCard("First", Logos, Creature, Common, WithPower(2)), 0)
	g.AddToHand(NewCard("Second", Logos, Creature, Common, WithPower(2)), 0)
	g.SetChooser(0, orderRejectChooser{})

	PlayFrom{
		From:  Hand,
		Types: CardTypesOf(Creature),
	}.Resolve(
		&EffectContext{Resolver: g, Controller: 0},
	)

	if got := len(g.Battleline(0)); got != 0 {
		t.Errorf("battleline holds %d creatures, want none after declining", got)
	}
}

// TestPlayFromTypesText covers the multi-type filter's rendered noun — Com.
// Officer Kirby frees a non-creature card.
func TestPlayFromTypesText(t *testing.T) {
	e := PlayFrom{
		From:  Hand,
		House: HouseMatcher{Kind: MatchExceptHouse, House: StarAlliance},
		Types: CardTypesOf(Artifact, Upgrade, Tactic),
	}
	if got, want := e.Text(), "play a non-Star Alliance artifact, upgrade, or tactic"; got != want {
		t.Errorf("Text = %q, want %q", got, want)
	}
}

// TestPlayFromTypesFiltersCandidates covers the multi-type filter dropping cards
// of an excluded type: a creature is not offered when only non-creatures qualify.
func TestPlayFromTypesFiltersCandidates(t *testing.T) {
	g := NewGame("A", "B", 1)
	art := g.AddToHand(NewCard("art", Mars, Artifact, Common), 0)
	g.AddToHand(NewCard("creat", Mars, Creature, Common, WithPower(2)), 0)

	cands := PlayFrom{
		From:  Hand,
		Types: CardTypesOf(Artifact, Upgrade, Tactic),
	}.candidates(&EffectContext{Resolver: g, Controller: 0})

	if len(cands) != 1 || cands[0] != art {
		t.Errorf("candidates = %v, want [%d] (the creature filtered out)", cands, art)
	}
}

func TestPlayFromFiltersByType(t *testing.T) {
	g := started(t)
	artifact := g.AddToHand(NewCard("Gadget", Logos, Artifact, Common), 0)
	g.AddToHand(NewCard("Thinker", Logos, Creature, Common, WithPower(2)), 0)

	PlayFrom{
		From:  Hand,
		Types: CardTypesOf(Artifact),
	}.Resolve(
		&EffectContext{Resolver: g, Controller: 0},
	)

	if got := g.Artifacts(0); len(got) != 1 || got[0] != artifact {
		t.Errorf("artifacts = %v, want [%d]", got, artifact)
	}
}

func TestGamePlayFromIgnoresACardElsewhere(t *testing.T) {
	g := started(t)
	deckCard := g.AddToDeck(NewCard("Not In Hand", Logos, Creature, Common, WithPower(2)), 0)

	g.PlayFromHand(0, deckCard)

	if got := len(g.Battleline(0)); got != 0 {
		t.Errorf("battleline holds %d creatures, want none", got)
	}
}

func TestActivePlayer(t *testing.T) {
	g := started(t)
	if got := g.ActivePlayer(); got != g.State.ActivePlayer {
		t.Errorf("ActivePlayer = %d, want %d", got, g.State.ActivePlayer)
	}
}

func TestPlayFromDiscardPile(t *testing.T) {
	e := PlayFrom{From: Discard, Types: CardTypesOf(Creature)}
	want := "play a creature from your discard pile"
	if e.Text() != want {
		t.Errorf("text = %q, want %q", e.Text(), want)
	}

	g := started(t)
	creature := g.AddToDiscard(NewCard("Risen", Logos, Creature, Common, WithPower(2)), 0)
	g.AddToDiscard(NewCard("Spent Tactic", Logos, Tactic, Common), 0) // wrong type

	e.Resolve(&EffectContext{Resolver: g, Controller: 0})

	if got := g.Battleline(0); len(got) != 1 || got[0] != creature {
		t.Errorf("battleline = %v, want the creature %d from the discard pile", got, creature)
	}
	if g.State.Discard[0].contains(creature) {
		t.Error("the played creature should have left the discard pile")
	}
}

func TestGamePlayFromDiscardIgnoresACardElsewhere(t *testing.T) {
	g := started(t)
	inHand := g.AddToHand(NewCard("Not Discarded", Logos, Creature, Common, WithPower(2)), 0)

	g.PlayFromDiscard(0, inHand)

	if got := len(g.Battleline(0)); got != 0 {
		t.Errorf("battleline holds %d creatures, want none", got)
	}
}

func TestPlayFromValidatesItsSourcePile(t *testing.T) {
	if err := (PlayFrom{}).validate(); err == nil {
		t.Error("an effect with no source pile should be rejected")
	}
	if err := (PlayFrom{From: Discard}).validate(); err != nil {
		t.Errorf("playing from the discard pile should validate: %v", err)
	}
}

// TestPlayFromArchives covers Project Z.Y.X.: a creature plays a card out of its
// controller's archives, bypassing the active-house gate.
func TestPlayFromArchives(t *testing.T) {
	e := PlayFrom{From: Archives, Types: CardTypesOf(Creature)}
	want := "play a creature from your archives"
	if e.Text() != want {
		t.Errorf("text = %q, want %q", e.Text(), want)
	}
	if err := (PlayFrom{From: Archives}).validate(); err != nil {
		t.Errorf("playing from the archives should validate: %v", err)
	}

	g := started(t) // Brobnar is the active house.
	creature := g.AddToArchives(NewCard("Archived", Logos, Creature, Common, WithPower(2)), 0)
	g.AddToArchives(NewCard("Archived Tactic", Logos, Tactic, Common), 0) // wrong type

	e.Resolve(&EffectContext{Resolver: g, Controller: 0})

	if got := g.Battleline(0); len(got) != 1 || got[0] != creature {
		t.Errorf("battleline = %v, want the creature %d from the archives", got, creature)
	}
	if g.State.Archives[0].contains(creature) {
		t.Error("the played creature should have left the archives")
	}
}

func TestGamePlayFromArchivesIgnoresACardElsewhere(t *testing.T) {
	g := started(t)
	inHand := g.AddToHand(NewCard("Not Archived", Logos, Creature, Common, WithPower(2)), 0)

	g.PlayFromArchives(0, inHand)

	if got := len(g.Battleline(0)); got != 0 {
		t.Errorf("battleline holds %d creatures, want none", got)
	}
}

func TestPlayRandomFromOpponentArchivesText(t *testing.T) {
	want := "play a random card from your opponent's archives"
	if got := (PlayFromOpponent{From: Archives}).Text(); got != want {
		t.Errorf("text = %q", got)
	}
}

func TestPlayTopOfOpponentDeckText(t *testing.T) {
	want := "play the top card of your opponent's deck"
	if got := (PlayFromOpponent{From: Deck}).Text(); got != want {
		t.Errorf("text = %q", got)
	}
}

func TestPlayFromOpponentEffectsValidate(t *testing.T) {
	if err := validateEffect(PlayFromOpponent{From: Archives}); err != nil {
		t.Errorf("archives validate = %v", err)
	}
	if err := validateEffect(PlayFromOpponent{From: Deck}); err != nil {
		t.Errorf("deck validate = %v", err)
	}
	if err := validateEffect(PlayFromOpponent{From: Hand}); err == nil {
		t.Error("an unsupported zone should not validate")
	}
}

// TestPlayTopOfOpponentDeckCreature covers Murkens grabbing a creature off the
// top of the opponent's deck: it enters the controller's battleline under their
// control, still owned by the opponent, and its Play: ability resolves for the
// controller.
func TestPlayTopOfOpponentDeckCreature(t *testing.T) {
	g := started(t)
	foe := g.AddToDeck(NewCard("Foe", Logos, Creature, Common, WithPower(3),
		WithAbility(TriggerAfterPlay, GainAember{Player: Controller, Amount: 2})), 1)

	g.PlayFromOpponent(0, Deck)

	if !g.inPlay(foe) {
		t.Fatal("the creature should be in play")
	}
	if got := g.controller(foe); got != 0 {
		t.Errorf("controller = %d, want 0", got)
	}
	if got := g.owner(foe); got != 1 {
		t.Errorf("owner = %d, want 1 (unchanged)", got)
	}
	if got := g.Battleline(0); len(got) != 1 || got[0] != foe {
		t.Errorf("battleline[0] = %v, want [%d]", got, foe)
	}
	if got := g.State.Aember[0]; got != 2 {
		t.Errorf("player 0 Æmber = %d, want 2 from the Play: ability", got)
	}
}

// TestPlayRandomFromOpponentArchives covers Murkens grabbing an action out of the
// opponent's archives: it resolves as the controller's own play and returns to
// its owner's discard pile.
func TestPlayRandomFromOpponentArchives(t *testing.T) {
	g := started(t)
	act := g.AddToArchives(NewCard("Snatched", Logos, Tactic, Common,
		WithAbility(TriggerAfterPlay, GainAember{Player: Controller, Amount: 3})), 1)

	g.PlayFromOpponent(0, Archives)

	if got := g.State.Aember[0]; got != 3 {
		t.Errorf("player 0 Æmber = %d, want 3", got)
	}
	if g.State.Archives[1].contains(act) {
		t.Error("the action should have left the archives")
	}
	if d := g.Discard(1); len(d) != 1 || d[0] != act {
		t.Errorf("owner discard = %v, want [%d]", d, act)
	}
}

func TestPlayRandomFromOpponentArchivesEmpty(t *testing.T) {
	g := started(t)
	g.PlayFromOpponent(0, Archives) // no panic on empty archives
}

func TestPlayTopOfOpponentDeckEmpty(t *testing.T) {
	g := started(t)
	g.PlayFromOpponent(0, Deck) // no panic on empty deck
}

// TestPlayForeignRevertsRejectedPlay covers a foreign play the gates reject: the
// optimistic control taken before the play is reverted, and the card is left
// untouched in its pile.
func TestPlayForeignRevertsRejectedPlay(t *testing.T) {
	g := started(t) // player 0 is active
	foe := g.AddToDeck(NewCard("Foe", Logos, Creature, Common, WithPower(3)), 0)

	// Player 1 is not active, so the play is rejected before the card is removed.
	g.PlayFromOpponent(1, Deck)

	if got := g.controller(foe); got != 0 {
		t.Errorf("controller = %d, want 0 (control reverted)", got)
	}
	if got := g.State.Cards[foe].ControlPlus; got != 0 {
		t.Errorf("ControlPlus = %d, want 0", got)
	}
	if !g.State.Deck[0].contains(foe) {
		t.Error("the card should still be in its deck")
	}
}

// TestPlayRandomFromOpponentArchivesResolve plays a random card out of the
// opponent's archives as the controller's own.
func TestPlayRandomFromOpponentArchivesResolve(t *testing.T) {
	g := started(t)
	act := g.AddToArchives(NewCard("A", Logos, Tactic, Common,
		WithAbility(TriggerAfterPlay, GainAember{Player: Controller, Amount: 1})), 1)

	PlayFromOpponent{From: Archives}.Resolve(&EffectContext{Resolver: g, Controller: 0})

	if g.State.Archives[1].contains(act) {
		t.Error("archives card should have been played")
	}
	if got := g.State.Aember[0]; got != 1 {
		t.Errorf("player 0 Æmber = %d, want 1", got)
	}
}

// TestPlayTopOfOpponentDeckResolve plays the top card of the opponent's deck as
// the controller's own.
func TestPlayTopOfOpponentDeckResolve(t *testing.T) {
	g := started(t)
	top := g.AddToDeck(NewCard("D", Logos, Creature, Common, WithPower(2)), 1)

	PlayFromOpponent{From: Deck}.Resolve(&EffectContext{Resolver: g, Controller: 0})

	if !g.inPlay(top) || g.controller(top) != 0 {
		t.Error("deck top should be in play under player 0")
	}
}
