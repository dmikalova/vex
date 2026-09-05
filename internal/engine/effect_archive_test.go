package engine

import "testing"

func TestArchiveEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	c1 := g.AddToHand(testCreature("c1", 1), 0)
	c2 := g.AddToHand(testCreature("c2", 1), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	if (ArchiveFromHand{Amount: 1}).Text() != "archive a card from your hand" {
		t.Errorf("archive text = %q", (ArchiveFromHand{Amount: 1}).Text())
	}
	if (ArchiveFromHand{Amount: 2}).Text() != "archive 2 cards from your hand" {
		t.Errorf("archive plural text = %q", (ArchiveFromHand{Amount: 2}).Text())
	}

	// The default chooser archives the first hand card (c1).
	(ArchiveFromHand{Amount: 1}).Resolve(ctx)
	if g.State.Archives[0].Count != 1 || g.State.Archives[0].IDs[0] != c1 {
		t.Errorf("archives = %v, want [%d]", g.State.Archives[0].slice(), c1)
	}
	if len(g.Hand(0)) != 1 || g.Hand(0)[0] != c2 {
		t.Errorf("hand = %v, want [%d]", g.Hand(0), c2)
	}

	// Archiving more than the hand holds stops when the hand empties.
	(ArchiveFromHand{Amount: 5}).Resolve(ctx)
	if len(g.Hand(0)) != 0 {
		t.Errorf("hand should be empty, got %v", g.Hand(0))
	}
	if g.State.Archives[0].Count != 2 {
		t.Errorf("archives count = %d, want 2", g.State.Archives[0].Count)
	}
}

func TestArchiveEffectDeclined(t *testing.T) {
	g := NewGame("A", "B", 1)
	// Two cards, so declining is a real choice (a sole card would be auto-archived).
	g.AddToHand(testCreature("c", 1), 0)
	g.AddToHand(testCreature("d", 1), 0)
	g.SetChooser(0, orderRejectChooser{})
	ctx := &EffectContext{Resolver: g, Controller: 0}
	(ArchiveFromHand{Amount: 1}).Resolve(ctx)
	if g.State.Archives[0].Count != 0 {
		t.Error("a declined archive choice should archive nothing")
	}
}

// TestArchiveFromHandUpTo covers Mobius Scroll's "up to 2 cards": the controller
// archives one and stops, leaving the rest in hand.
func TestArchiveFromHandUpTo(t *testing.T) {
	e := ArchiveFromHand{Amount: 2, UpTo: true}
	if got := e.Text(); got != "archive up to 2 cards from your hand" {
		t.Errorf("text = %q", got)
	}

	g := NewGame("A", "B", 1)
	g.AddToHand(testCreature("c", 1), 0)
	g.AddToHand(testCreature("d", 1), 0)
	g.SetChooser(0, &cardDecliner{})
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e.Resolve(ctx)
	if g.State.Archives[0].Count != 2 {
		t.Errorf("archives = %d, want 2", g.State.Archives[0].Count)
	}

	g2 := NewGame("A", "B", 1)
	g2.AddToHand(testCreature("c", 1), 0)
	g2.AddToHand(testCreature("d", 1), 0)
	g2.SetChooser(0, &cardDecliner{decline: true})
	e.Resolve(&EffectContext{Resolver: g2, Controller: 0})
	if g2.State.Archives[0].Count != 0 {
		t.Error("declining the up-to archive should archive nothing")
	}
}

func TestArchiveFromDiscardEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	buried := g.AddToDiscard(testCreature("buried", 1), 0)
	g.AddToDiscard(testCreature("other", 1), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	if (ArchiveFromDiscard{}).Text() != "archive a card from your discard pile" {
		t.Errorf("text = %q", (ArchiveFromDiscard{}).Text())
	}

	// The default chooser archives the first discard card.
	(ArchiveFromDiscard{}).Resolve(ctx)
	if g.State.Archives[0].Count != 1 || g.State.Archives[0].IDs[0] != buried {
		t.Errorf("archives = %v, want [%d]", g.State.Archives[0].slice(), buried)
	}
	if len(g.Discard(0)) != 1 {
		t.Errorf("discard = %v, want one card left", g.Discard(0))
	}

	// An empty discard pile archives nothing.
	empty := &EffectContext{Resolver: g, Controller: 1}
	(ArchiveFromDiscard{}).Resolve(empty)
	if g.State.Archives[1].Count != 0 {
		t.Error("archiving from an empty discard pile should do nothing")
	}
}

func TestArchiveFromDiscardDeclined(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.AddToDiscard(testCreature("c", 1), 0)
	g.AddToDiscard(testCreature("d", 1), 0)
	g.SetChooser(0, orderRejectChooser{})
	ctx := &EffectContext{Resolver: g, Controller: 0}
	(ArchiveFromDiscard{}).Resolve(ctx)
	if g.State.Archives[0].Count != 0 {
		t.Error("a declined archive choice should archive nothing")
	}
}

func TestArchiveTopOfDeckEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	top := g.AddToDeck(testCreature("top", 1), 0)
	g.AddToDeck(testCreature("next", 1), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	if (ArchiveTopOfDeck{Amount: 1}).Text() != "archive the top card of your deck" {
		t.Errorf("text = %q", (ArchiveTopOfDeck{Amount: 1}).Text())
	}
	if (ArchiveTopOfDeck{Amount: 2}).Text() != "archive the top 2 cards of your deck" {
		t.Errorf("plural text = %q", (ArchiveTopOfDeck{Amount: 2}).Text())
	}

	(ArchiveTopOfDeck{Amount: 1}).Resolve(ctx)
	if g.State.Archives[0].Count != 1 || g.State.Archives[0].IDs[0] != top {
		t.Errorf("archived %v, want the top card %d", g.State.Archives[0].slice(), top)
	}

	// Archiving more than the deck holds stops when the deck empties.
	(ArchiveTopOfDeck{Amount: 5}).Resolve(ctx)
	if g.State.Deck[0].Count != 0 {
		t.Errorf("deck should be empty, got %d", g.State.Deck[0].Count)
	}
	if g.State.Archives[0].Count != 2 {
		t.Errorf("archives count = %d, want 2", g.State.Archives[0].Count)
	}
}
func TestArchiveTopOfDiscardEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.State.Discard[0].add(g.Register(testCreature("bottom", 1), 0))
	top := g.Register(testCreature("top", 1), 0)
	g.State.Discard[0].add(top)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	if (ArchiveTopOfDiscard{Amount: 1}).Text() != "archive the top card of your discard pile" {
		t.Errorf("text = %q", (ArchiveTopOfDiscard{Amount: 1}).Text())
	}
	if (ArchiveTopOfDiscard{Amount: 2}).Text() != "archive the top 2 cards of your discard pile" {
		t.Errorf("plural text = %q", (ArchiveTopOfDiscard{Amount: 2}).Text())
	}

	// The most recently discarded card is the top and archives first.
	(ArchiveTopOfDiscard{Amount: 1}).Resolve(ctx)
	if g.State.Archives[0].Count != 1 || g.State.Archives[0].IDs[0] != top {
		t.Errorf("archived %v, want the top card %d", g.State.Archives[0].slice(), top)
	}

	// Archiving more than the discard holds stops when the pile empties.
	(ArchiveTopOfDiscard{Amount: 5}).Resolve(ctx)
	if g.State.Discard[0].Count != 0 {
		t.Errorf("discard should be empty, got %d", g.State.Discard[0].Count)
	}
	if g.State.Archives[0].Count != 2 {
		t.Errorf("archives count = %d, want 2", g.State.Archives[0].Count)
	}
}
func TestArchiveFromPlayEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddArtifact(NewCard("quest", Sanctum, Artifact, Rare), 0)
	knight := g.AddToBattleline(
		NewCard("knight", Sanctum, Creature, Common, WithPower(4), WithTraits(Knight)),
		0,
	)
	nonKnight := g.AddToBattleline(
		NewCard("cleric", Sanctum, Creature, Common, WithPower(4), WithTraits(Cleric)),
		0,
	)
	enemyKnight := g.AddToBattleline(
		NewCard("enemy", Sanctum, Creature, Common, WithPower(4), WithTraits(Knight)),
		1,
	)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	e := ArchiveFromPlay{Target: Target{Kind: TargetEachFriendlyCreature}.WithTrait(Knight)}
	if e.Text() != "archive each friendly Knight trait creature from play" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)

	if g.inPlay(knight) || !g.State.Archives[0].contains(knight) {
		t.Error("friendly Knight should be archived")
	}
	if !g.inPlay(nonKnight) {
		t.Error("friendly non-Knight should stay in play")
	}
	if !g.inPlay(enemyKnight) {
		t.Error("enemy Knight should stay in play")
	}
	if err := validateEffect(e); err != nil {
		t.Errorf("valid archive-from-play effect rejected: %v", err)
	}
	if err := validateEffect(ArchiveFromPlay{}); err == nil {
		t.Error("unset target should be rejected")
	}
}

func TestArchiveFromPlayFriendlyInPlay(t *testing.T) {
	g := NewGame("A", "B", 1)
	art := g.AddArtifact(NewCard("relic", Mars, Artifact, Common), 0)
	g.AddToBattleline(NewCard("enemy", Mars, Creature, Common, WithPower(3)), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := ArchiveFromPlay{Target: Target{Kind: TargetChosenFriendlyCreatureOrArtifact}}
	if e.Text() != "archive a friendly creature or artifact from play" {
		t.Errorf("text = %q", e.Text())
	}
	// The friendly artifact is the sole candidate, so it is archived automatically.
	e.Resolve(ctx)
	if g.inPlay(art) || !g.State.Archives[0].contains(art) {
		t.Error("friendly artifact should be archived")
	}
}

// A "you may archive a friendly creature or artifact from play" is one card
// choice, so the player picks the card directly instead of first answering Yes
// (Vezyma Thinkdrone).
func TestMayDeclinableArchiveFromPlay(t *testing.T) {
	e := May{Do: ArchiveFromPlay{Target: Target{Kind: TargetChosenFriendlyCreatureOrArtifact}}}
	if !e.Do.(declinableEffect).declinable() {
		t.Fatal("a chosen-target ArchiveFromPlay should be declinable")
	}

	g := NewGame("A", "B", 1)
	ch := &cardDecliner{}
	g.SetChooser(0, ch)
	keep := g.AddArtifact(NewCard("keep", Mars, Artifact, Common), 0)
	doomed := g.AddArtifact(NewCard("doomed", Mars, Artifact, Common), 0)
	e.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if ch.asked != 1 {
		t.Errorf("declinable prompts = %d, want 1", ch.asked)
	}
	if !g.inPlay(keep) {
		t.Error("the unchosen artifact should stay in play")
	}
	if g.inPlay(doomed) || !g.State.Archives[0].contains(doomed) {
		t.Error("the chosen artifact should have been archived")
	}

	declined := NewGame("A", "B", 1)
	declined.SetChooser(0, &cardDecliner{decline: true})
	other := declined.AddArtifact(NewCard("relic", Mars, Artifact, Common), 0)
	e.Resolve(&EffectContext{Resolver: declined, Controller: 0})
	if !declined.inPlay(other) {
		t.Error("a declined May should archive nothing")
	}
}

// A "you may reveal a creature from your hand and archive it" is one card
// choice, so the player picks the card directly instead of first answering Yes
// (Zyzzix the Many).
func TestMayDeclinableArchiveFromHand(t *testing.T) {
	e := May{Do: ArchiveFromHand{Amount: 1, Type: Creature}}
	if !e.Do.(declinableEffect).declinable() {
		t.Fatal("a single-count ArchiveFromHand should be declinable")
	}
	if (ArchiveFromHand{Amount: 1, UpTo: true}).declinable() {
		t.Error("an UpTo ArchiveFromHand should keep its own cycle, not be declinable")
	}
	if (ArchiveFromHand{Amount: 2}).declinable() {
		t.Error("a multi-count ArchiveFromHand should keep its own cycle, not be declinable")
	}

	g := NewGame("A", "B", 1)
	ch := &cardDecliner{}
	g.SetChooser(0, ch)
	g.AddToHand(testCreature("keep", 1), 0)
	doomed := g.AddToHand(testCreature("doomed", 1), 0)
	e.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if ch.asked != 1 {
		t.Errorf("declinable prompts = %d, want 1", ch.asked)
	}
	if g.State.Archives[0].Count != 1 || !g.State.Archives[0].contains(doomed) {
		t.Error("the chosen card should have been archived")
	}

	declined := NewGame("A", "B", 1)
	declined.SetChooser(0, &cardDecliner{decline: true})
	declined.AddToHand(testCreature("keep", 1), 0)
	e.Resolve(&EffectContext{Resolver: declined, Controller: 0})
	if declined.State.Archives[0].Count != 0 {
		t.Error("a declined May should archive nothing")
	}

	// A filter admitting nothing is not offered at all (Zyzzix the Many with no
	// creature in hand).
	empty := NewGame("A", "B", 1)
	empty.AddToHand(NewCard("tactic", Mars, Tactic, Common), 0)
	e.Resolve(&EffectContext{Resolver: empty, Controller: 0})
	if empty.State.Archives[0].Count != 0 {
		t.Error("an empty candidate set should archive nothing")
	}
}

func TestArchivesOfferedOnChooseHouse(t *testing.T) {
	g := NewGame("A", "B", 1)
	a := g.AddToHand(testCreature("a", 1), 0)
	(ArchiveFromHand{Amount: 1}).Resolve(&EffectContext{Resolver: g, Controller: 0})
	if g.State.Archives[0].Count != 1 {
		t.Fatalf("setup: archives count = %d, want 1", g.State.Archives[0].Count)
	}

	// The archives are NOT taken at the start of the turn.
	g.StartTurn(0)
	if g.State.Archives[0].Count != 1 {
		t.Error("archives should not return before the house is chosen")
	}

	// Choosing a house offers them; the default chooser accepts (takes them).
	if err := g.ChooseHouse(0, Logos); err != nil {
		t.Fatalf("ChooseHouse: %v", err)
	}
	if g.State.Archives[0].Count != 0 {
		t.Error("accepting the offer should empty the archives")
	}
	if len(g.Hand(0)) != 1 || g.Hand(0)[0] != a {
		t.Errorf("hand = %v, want the returned card [%d]", g.Hand(0), a)
	}
}

func TestArchivesDeclinedOnChooseHouse(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.AddToHand(testCreature("a", 1), 0)
	(ArchiveFromHand{Amount: 1}).Resolve(&EffectContext{Resolver: g, Controller: 0})
	g.SetChooser(0, optionPicker{idx: 1}) // decline the offer

	g.StartTurn(0)
	if err := g.ChooseHouse(0, Logos); err != nil {
		t.Fatalf("ChooseHouse: %v", err)
	}
	if g.State.Archives[0].Count != 1 {
		t.Error("declining the offer should leave the cards archived")
	}
	if len(g.Hand(0)) != 0 {
		t.Errorf("hand = %v, want empty (offer declined)", g.Hand(0))
	}
}

// TestArchiveFromHandFiltered covers the filtered, revealed archive: only cards
// matching the type and house filters are offered, the effect reports whether it
// archived anything, and the text reads as the reveal it is.
func TestArchiveFromHandFiltered(t *testing.T) {
	e := ArchiveFromHand{Amount: 1, Type: Creature, House: Mars, Revealed: true}
	want := "reveal a Mars creature from your hand and archive it"
	if e.Text() != want {
		t.Errorf("text = %q, want %q", e.Text(), want)
	}
	if got := (ArchiveFromHand{Amount: 2}).Text(); got != "archive 2 cards from your hand" {
		t.Errorf("plain plural text = %q", got)
	}
	if got := (ArchiveFromHand{Amount: 1, Type: Artifact}).Text(); got != "archive an artifact from your hand" {
		t.Errorf("artifact text = %q", got)
	}

	g := NewGame("A", "B", 1)
	tactic := g.AddToHand(NewCard("Trick", Mars, Tactic, Common), 0)
	offHouse := g.AddToHand(testCreature("Logos One", 3), 0)
	martian := g.AddToHand(NewCard("Martian", Mars, Creature, Common, WithPower(3)), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	if !e.resolveGate(ctx) {
		t.Fatal("archiving a matching creature should report success")
	}
	if g.State.Archives[0].Count != 1 || g.State.Archives[0].IDs[0] != martian {
		t.Errorf("archives = %v, want [%d]", g.State.Archives[0].slice(), martian)
	}
	if len(g.Hand(0)) != 2 || g.Hand(0)[0] != tactic || g.Hand(0)[1] != offHouse {
		t.Errorf("hand = %v, want the filtered-out cards", g.Hand(0))
	}
	if e.resolveGate(ctx) {
		t.Error("with no matching card left the gate should report nothing archived")
	}
}
