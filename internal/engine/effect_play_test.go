package engine

import "testing"

func TestPlayFromHandText(t *testing.T) {
	cases := []struct {
		name   string
		effect PlayFromHand
		want   string
	}{
		{"any card", PlayFromHand{}, "play a card"},
		{"excluded house", PlayFromHand{House: Logos, Except: true}, "play a non-Logos card"},
		{"named house", PlayFromHand{House: Mars}, "play a Mars card"},
		{"typed", PlayFromHand{Type: Creature}, "play a creature"},
		{"any type", PlayFromHand{Type: AnyType}, "play a card"},
		{
			"house and type",
			PlayFromHand{House: Untamed, Type: Artifact},
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

func TestPlayFromHandValidate(t *testing.T) {
	if err := (PlayFromHand{Except: true}).validate(); err == nil {
		t.Error("Except without a house should not validate")
	}
	if err := (PlayFromHand{House: Logos, Except: true}).validate(); err != nil {
		t.Errorf("validate = %v, want nil", err)
	}
}

func TestPlayFromHandPlaysAChosenCard(t *testing.T) {
	g := started(t) // Brobnar is the active house.
	off := g.AddToHand(NewCard("Off House", Logos, Creature, Common, WithPower(2)), 0)
	g.AddToHand(NewCard("On House", Brobnar, Creature, Common, WithPower(2)), 0)

	PlayFromHand{House: Brobnar, Except: true}.Resolve(&EffectContext{Resolver: g, Controller: 0})

	if got := g.Battleline(0); len(got) != 1 || got[0] != off {
		t.Errorf("battleline = %v, want the off-house creature %d", got, off)
	}
}

func TestPlayFromHandWithNoCandidate(t *testing.T) {
	g := started(t)
	g.AddToHand(NewCard("On House", Brobnar, Creature, Common, WithPower(2)), 0)

	PlayFromHand{House: Brobnar, Except: true}.Resolve(&EffectContext{Resolver: g, Controller: 0})

	if got := len(g.Battleline(0)); got != 0 {
		t.Errorf("battleline holds %d creatures, want none", got)
	}
}

func TestPlayFromHandDeclined(t *testing.T) {
	g := started(t)
	g.AddToHand(NewCard("First", Logos, Creature, Common, WithPower(2)), 0)
	g.AddToHand(NewCard("Second", Logos, Creature, Common, WithPower(2)), 0)
	g.SetChooser(0, orderRejectChooser{})

	PlayFromHand{Type: Creature}.Resolve(&EffectContext{Resolver: g, Controller: 0})

	if got := len(g.Battleline(0)); got != 0 {
		t.Errorf("battleline holds %d creatures, want none after declining", got)
	}
}

func TestPlayFromHandFiltersByType(t *testing.T) {
	g := started(t)
	artifact := g.AddToHand(NewCard("Gadget", Logos, Artifact, Common), 0)
	g.AddToHand(NewCard("Thinker", Logos, Creature, Common, WithPower(2)), 0)

	PlayFromHand{Type: Artifact}.Resolve(&EffectContext{Resolver: g, Controller: 0})

	if got := g.Artifacts(0); len(got) != 1 || got[0] != artifact {
		t.Errorf("artifacts = %v, want [%d]", got, artifact)
	}
}

func TestGamePlayFromHandIgnoresACardElsewhere(t *testing.T) {
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
