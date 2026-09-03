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
			PlayFrom{From: Hand, House: Logos, Except: true},
			"play a non-Logos card",
		},
		{"named house", PlayFrom{From: Hand, House: Mars}, "play a Mars card"},
		{"typed", PlayFrom{From: Hand, Type: Creature}, "play a creature"},
		{"any type", PlayFrom{From: Hand, Type: AnyType}, "play a card"},
		{
			"house and type",
			PlayFrom{From: Hand, House: Untamed, Type: Artifact},
			"play an Untamed artifact",
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
	if err := (PlayFrom{Except: true}).validate(); err == nil {
		t.Error("Except without a house should not validate")
	}
	if err := (PlayFrom{From: Hand, House: Logos, Except: true}).validate(); err != nil {
		t.Errorf("validate = %v, want nil", err)
	}
}

func TestPlayFromPlaysAChosenCard(t *testing.T) {
	g := started(t) // Brobnar is the active house.
	off := g.AddToHand(NewCard("Off House", Logos, Creature, Common, WithPower(2)), 0)
	g.AddToHand(NewCard("On House", Brobnar, Creature, Common, WithPower(2)), 0)

	PlayFrom{
		From:   Hand,
		House:  Brobnar,
		Except: true,
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
		From:   Hand,
		House:  Brobnar,
		Except: true,
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

	PlayFrom{From: Hand, Type: Creature}.Resolve(&EffectContext{Resolver: g, Controller: 0})

	if got := len(g.Battleline(0)); got != 0 {
		t.Errorf("battleline holds %d creatures, want none after declining", got)
	}
}

func TestPlayFromFiltersByType(t *testing.T) {
	g := started(t)
	artifact := g.AddToHand(NewCard("Gadget", Logos, Artifact, Common), 0)
	g.AddToHand(NewCard("Thinker", Logos, Creature, Common, WithPower(2)), 0)

	PlayFrom{From: Hand, Type: Artifact}.Resolve(&EffectContext{Resolver: g, Controller: 0})

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
	e := PlayFrom{From: Discard, Type: Creature}
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
	if err := (PlayFrom{From: Archives}).validate(); err == nil {
		t.Error("the archives are not a pile a card may be played from")
	}
	if err := (PlayFrom{From: Discard}).validate(); err != nil {
		t.Errorf("playing from the discard pile should validate: %v", err)
	}
}
