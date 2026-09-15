package engine

import "testing"

func TestDiscardOpponentArchivesOrDeckTopText(t *testing.T) {
	want := "discard a random card from your opponent's archives or the top card " +
		"of their deck"
	if got := (DiscardOpponentArchivesOrDeckTop{}).Text(); got != want {
		t.Errorf("text = %q, want %q", got, want)
	}
	if got := (PlayItFromOpponentDiscard{}).Text(); got != "play it as if it were yours" {
		t.Errorf("play text = %q", got)
	}
}

// Choosing the archives on an empty archives binds nothing, so a following
// "play it" effect has no card to act on.
func TestDiscardOpponentSourceEmptyArchives(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.State.ActivePlayer = 0
	g.SetChooser(0, optionPicker{idx: 0}) // your opponent's archives
	ctx := &EffectContext{Resolver: g, Controller: 0}

	DiscardOpponentArchivesOrDeckTop{}.Resolve(ctx)
	if ctx.HasIt {
		t.Error("empty archives should bind no card")
	}
}

// Choosing the deck top on an empty deck binds nothing.
func TestDiscardOpponentSourceEmptyDeck(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.State.ActivePlayer = 0
	g.State.Deck[1].Count = 0
	g.SetChooser(0, optionPicker{idx: 1}) // the top card of their deck
	ctx := &EffectContext{Resolver: g, Controller: 0}

	DiscardOpponentArchivesOrDeckTop{}.Resolve(ctx)
	if ctx.HasIt {
		t.Error("empty deck should bind no card")
	}
}

// Discarding the opponent's deck top binds it, so "play it" can reach it.
func TestDiscardOpponentSourceDeckTopBindsIt(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.State.ActivePlayer = 0
	top := g.Register(testCreature("t", 1), 1)
	g.State.Deck[1].add(top)
	g.SetChooser(0, optionPicker{idx: 1})
	ctx := &EffectContext{Resolver: g, Controller: 0}

	DiscardOpponentArchivesOrDeckTop{}.Resolve(ctx)
	if !ctx.HasIt || ctx.It != top {
		t.Errorf("It = %v (has %v), want %v", ctx.It, ctx.HasIt, top)
	}
}

// Choosing the archives discards a random card from them and binds it, so a
// following "play it" effect can reach the discarded card.
func TestDiscardOpponentSourceArchivesBindsIt(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.State.ActivePlayer = 0
	id := g.Register(testCreature("a", 1), 1)
	g.State.Archives[1].add(id)
	g.SetChooser(0, optionPicker{idx: 0}) // your opponent's archives
	ctx := &EffectContext{Resolver: g, Controller: 0}

	DiscardOpponentArchivesOrDeckTop{}.Resolve(ctx)
	if g.State.Archives[1].contains(id) {
		t.Error("the discarded card should have left the opponent's archives")
	}
	if !ctx.HasIt || ctx.It != id {
		t.Errorf("It = %v (has %v), want %v", ctx.It, ctx.HasIt, id)
	}
}

// Both source effects validate without configuration.
func TestDiscardOpponentSourceValidate(t *testing.T) {
	if err := (DiscardOpponentArchivesOrDeckTop{}).validate(); err != nil {
		t.Errorf("DiscardOpponentArchivesOrDeckTop validate = %v", err)
	}
	if err := (PlayItFromOpponentDiscard{}).validate(); err != nil {
		t.Errorf("PlayItFromOpponentDiscard validate = %v", err)
	}
}

// With no card in context, playing "it" does nothing.
func TestPlayItFromOpponentDiscardNoContext(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.State.ActivePlayer = 0
	ctx := &EffectContext{Resolver: g, Controller: 0}
	PlayItFromOpponentDiscard{}.Resolve(ctx) // must not panic
	if ctx.HasIt {
		t.Error("no card should be in context")
	}
}

// A card in the opponent's discard, bound as context, is played from there.
func TestPlayItFromOpponentDiscardPlaysContext(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.State.ActivePlayer = 0
	id := g.Register(testCreature("c", 1), 1)
	g.State.Discard[1].add(id)
	ctx := &EffectContext{Resolver: g, Controller: 0, It: id, HasIt: true}

	PlayItFromOpponentDiscard{}.Resolve(ctx)
	if g.State.Discard[1].indexOf(id) >= 0 {
		t.Error("card should have left the opponent's discard")
	}
}

// typeNoun renders every card type, including Tactic and Upgrade.
func TestTypeNounAllTypes(t *testing.T) {
	cases := map[CardType]string{
		Creature:  "creature",
		Artifact:  "artifact",
		Tactic:    "tactic",
		Upgrade:   "upgrade",
		TypeUnset: "card",
	}
	for typ, want := range cases {
		if got := typeNoun(typ); got != want {
			t.Errorf("typeNoun(%v) = %q, want %q", typ, got, want)
		}
	}
}

// ThatCard names the context card "that card" when "it" would be ambiguous.
func TestSubjectThatCardNoun(t *testing.T) {
	if got := ThatCard.noun(); got != "that card" {
		t.Errorf("ThatCard.noun() = %q", got)
	}
}
