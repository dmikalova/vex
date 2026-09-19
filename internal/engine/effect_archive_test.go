package engine

import "testing"

func TestArchiveEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	c1 := g.AddToHand(testCreature("c1", 1), 0)
	c2 := g.AddToHand(testCreature("c2", 1), 0)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}

	one := ArchiveCard{
		Zone:      Hand,
		Selection: Chosen{},
	}
	if one.Text() != "archive a card from your hand" {
		t.Errorf("archive text = %q", one.Text())
	}
	two := ArchiveCard{
		Zone:      Hand,
		Selection: Chosen{},
		Quantity:  Takes{N: Fixed(2)},
	}
	if two.Text() != "archive 2 cards from your hand" {
		t.Errorf("archive plural text = %q", two.Text())
	}

	// The default chooser archives the first hand card (c1).
	one.Resolve(ctx)
	if g.State.Archives[0].Count != 1 || g.State.Archives[0].IDs[0] != c1 {
		t.Errorf("archives = %v, want [%d]", g.State.Archives[0].slice(), c1)
	}
	if len(g.Hand(0)) != 1 || g.Hand(0)[0] != c2 {
		t.Errorf("hand = %v, want [%d]", g.Hand(0), c2)
	}

	// Archiving more than the hand holds stops when the hand empties.
	(ArchiveCard{
		Zone:      Hand,
		Selection: Chosen{},
		Quantity:  Takes{N: Fixed(5)},
	}).Resolve(ctx)
	if len(g.Hand(0)) != 0 {
		t.Errorf("hand should be empty, got %v", g.Hand(0))
	}
	if g.State.Archives[0].Count != 2 {
		t.Errorf("archives count = %d, want 2", g.State.Archives[0].Count)
	}
}

// TestArchiveCardBind covers Bind recording the archived card as the choice
// context so a later effect can read it (Blast from the Past).
func TestArchiveCardBind(t *testing.T) {
	g := NewGame("A", "B", 1)
	c := g.AddToDiscard(testCreature("c", 3), 0)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}

	(ArchiveCard{
		Zone:      Discard,
		Selection: Chosen{},
		Bind:      true,
	}).Resolve(ctx)
	if !ctx.HasIt || ctx.It != c {
		t.Errorf("bound It = %d (has %v), want %d", ctx.It, ctx.HasIt, c)
	}
}

// TestArchiveEnemyHand covers archiving from the opponent's hand into the
// caster's own archives — the abduction-from-hand path (Hidden Stash).
func TestArchiveEnemyHand(t *testing.T) {
	e := ArchiveCard{
		From:      Opponent,
		Zone:      Hand,
		Selection: Chosen{},
	}
	if got := e.Text(); got != "archive a card from your opponent's hand" {
		t.Errorf("text = %q", got)
	}
	// An opponent archive is only defined for the hand.
	if (ArchiveCard{
		From:      Opponent,
		Zone:      Discard,
		Selection: Chosen{},
	}).validate() == nil {
		t.Error("an opponent archive from the discard pile should be invalid")
	}
	if e.validate() != nil {
		t.Error("an opponent hand archive should be valid")
	}

	g := NewGame("A", "B", 1)
	enemy := g.AddToHand(testCreature("enemy", 1), 1)
	e.Resolve(&EffectContext{
		Resolver:   g,
		Controller: 0,
	})
	if g.State.Archives[0].Count != 1 || g.State.Archives[0].IDs[0] != enemy {
		t.Errorf("caster archives = %v, want [%d]", g.State.Archives[0].slice(), enemy)
	}
	if len(g.Hand(1)) != 0 {
		t.Errorf("opponent hand = %v, want empty", g.Hand(1))
	}
}

// TestArchiveCardValidate covers ArchiveCard's guards: a selection must be set,
// the source zone must be the hand or discard pile, and a written quantity must
// name a real number of cards.
func TestArchiveCardValidate(t *testing.T) {
	if (ArchiveCard{Zone: Hand}).validate() == nil {
		t.Error("a nil selection should not validate")
	}
	if (ArchiveCard{Selection: Chosen{}}).validate() == nil {
		t.Error("an unset zone should not validate")
	}
	if (ArchiveCard{
		Zone:      Deck,
		Selection: Chosen{},
	}).validate() == nil {
		t.Error("a Deck source should not validate")
	}
	if (ArchiveCard{
		Zone:      Hand,
		Selection: Chosen{},
		Quantity:  Takes{},
	}).validate() == nil {
		t.Error("a quantity with an unset count should not validate")
	}
	if (ArchiveCard{
		Zone:      Hand,
		Selection: Chosen{},
	}).validate() != nil {
		t.Error("a hand archive should validate")
	}
	if (ArchiveCard{
		Zone:      Discard,
		Selection: Named{Name: "x"},
	}).validate() != nil {
		t.Error("a discard archive should validate")
	}
}

// TestArchiveCardDeclinableFlags covers declinable(): only a single Optional
// Chosen archive is one clickable card; every other shape keeps its own cycle.
func TestArchiveCardDeclinableFlags(t *testing.T) {
	if !(ArchiveCard{
		Zone:      Hand,
		Selection: Chosen{Optional: true},
	}).declinable() {
		t.Error("a single Optional Chosen archive should be declinable")
	}
	if (ArchiveCard{
		Zone:      Hand,
		Selection: Chosen{},
	}).declinable() {
		t.Error("a mandatory archive should not be declinable")
	}
	if (ArchiveCard{
		Zone:      Hand,
		Selection: Chosen{Optional: true},
		Quantity:  UpTo{N: Fixed(2)},
	}).declinable() {
		t.Error("a multi-count (up-to) archive should not be declinable")
	}
	if (ArchiveCard{
		Zone:      Hand,
		Selection: Random{},
	}).declinable() {
		t.Error("a random archive is not a card choice, so not declinable")
	}
	if (ArchiveCard{
		Zone:      Discard,
		Selection: Named{Name: "x"},
	}).declinable() {
		t.Error("a named archive is not a card choice, so not declinable")
	}
}

func TestArchiveEffectDeclined(t *testing.T) {
	g := NewGame("A", "B", 1)
	// Two cards, so declining is a real choice (a sole card would be auto-archived).
	g.AddToHand(testCreature("c", 1), 0)
	g.AddToHand(testCreature("d", 1), 0)
	g.SetChooser(0, &cardDecliner{decline: true})
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}
	(ArchiveCard{
		Zone:      Hand,
		Selection: Chosen{Optional: true},
	}).Resolve(ctx)
	if g.State.Archives[0].Count != 0 {
		t.Error("a declined archive choice should archive nothing")
	}
}

// TestArchiveFromHandUpTo covers Mobius Scroll's "up to 2 cards": the controller
// archives one and stops, leaving the rest in hand.
func TestArchiveFromHandUpTo(t *testing.T) {
	e := ArchiveCard{
		Zone:      Hand,
		Selection: Chosen{Optional: true},
		Quantity:  UpTo{N: Fixed(2)},
	}
	if got := e.Text(); got != "archive up to 2 cards from your hand" {
		t.Errorf("text = %q", got)
	}

	g := NewGame("A", "B", 1)
	g.AddToHand(testCreature("c", 1), 0)
	g.AddToHand(testCreature("d", 1), 0)
	g.SetChooser(0, &cardDecliner{})
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}

	e.Resolve(ctx)
	if g.State.Archives[0].Count != 2 {
		t.Errorf("archives = %d, want 2", g.State.Archives[0].Count)
	}

	g2 := NewGame("A", "B", 1)
	g2.AddToHand(testCreature("c", 1), 0)
	g2.AddToHand(testCreature("d", 1), 0)
	g2.SetChooser(0, &cardDecliner{decline: true})
	e.Resolve(&EffectContext{
		Resolver:   g2,
		Controller: 0,
	})
	if g2.State.Archives[0].Count != 0 {
		t.Error("declining the up-to archive should archive nothing")
	}
}

func TestArchiveFromDiscardEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	buried := g.AddToDiscard(testCreature("buried", 1), 0)
	g.AddToDiscard(testCreature("other", 1), 0)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}

	e := ArchiveCard{
		Zone:      Discard,
		Selection: Chosen{},
	}
	if e.Text() != "archive a card from your discard pile" {
		t.Errorf("text = %q", e.Text())
	}

	// The default chooser archives the first discard card.
	e.Resolve(ctx)
	if g.State.Archives[0].Count != 1 || g.State.Archives[0].IDs[0] != buried {
		t.Errorf(
			"archives = %v, want [%d]",
			g.State.Archives[0].slice(),
			buried,
		)
	}
	if len(g.Discard(0)) != 1 {
		t.Errorf("discard = %v, want one card left", g.Discard(0))
	}

	// An empty discard pile archives nothing.
	empty := &EffectContext{
		Resolver:   g,
		Controller: 1,
	}
	e.Resolve(empty)
	if g.State.Archives[1].Count != 0 {
		t.Error("archiving from an empty discard pile should do nothing")
	}
}

func TestArchiveFromDiscardHouseFilter(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.AddToDiscard(NewCard("logos", Logos, Creature, Common, WithPower(1)), 0)
	mars := g.AddToDiscard(
		NewCard("mars", Mars, Creature, Common, WithPower(1)),
		0,
	)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}

	e := ArchiveCard{
		Zone:      Discard,
		Selection: Chosen{House: namedHouse(Mars)},
	}
	if e.Text() != "archive a Mars card from your discard pile" {
		t.Errorf("text = %q", e.Text())
	}
	// Only the Mars card is a candidate, so the default chooser archives it.
	e.Resolve(ctx)
	if g.State.Archives[0].Count != 1 || g.State.Archives[0].IDs[0] != mars {
		t.Errorf("archives = %v, want [%d]", g.State.Archives[0].slice(), mars)
	}
	// A discard pile with no card of the house archives nothing.
	empty := &EffectContext{
		Resolver:   g,
		Controller: 1,
	}
	g.AddToDiscard(NewCard("logos2", Logos, Creature, Common, WithPower(1)), 1)
	e.Resolve(empty)
	if g.State.Archives[1].Count != 0 {
		t.Error("no Mars card in discard should archive nothing")
	}
}

func TestArchiveTopOfDeckEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	top := g.AddToDeck(testCreature("top", 1), 0)
	g.AddToDeck(testCreature("next", 1), 0)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}

	// A positional selection needs an ordered zone, and a deck archive must be
	// positional (a deck has no chooseable cards).
	if (ArchiveCard{
		Zone:      Deck,
		Selection: Top{},
	}).validate() != nil {
		t.Error("a top-of-deck archive should validate")
	}
	if (ArchiveCard{
		Zone:      Hand,
		Selection: Top{},
	}).validate() == nil {
		t.Error("a positional archive from the hand should not validate")
	}
	if (ArchiveCard{
		Zone:      Deck,
		Selection: Chosen{},
	}).validate() == nil {
		t.Error("a non-positional deck archive should not validate")
	}

	if got := (ArchiveCard{
		Zone:      Deck,
		Selection: Top{},
	}).Text(); got != "archive the top card of your deck" {
		t.Errorf("text = %q", got)
	}
	if got := (ArchiveCard{
		Zone:      Deck,
		Selection: Top{},
		Quantity:  Takes{N: Fixed(2)},
	}).Text(); got != "archive the top 2 cards of your deck" {
		t.Errorf("plural text = %q", got)
	}

	(ArchiveCard{
		Zone:      Deck,
		Selection: Top{},
	}).Resolve(ctx)
	if g.State.Archives[0].Count != 1 || g.State.Archives[0].IDs[0] != top {
		t.Errorf(
			"archived %v, want the top card %d",
			g.State.Archives[0].slice(),
			top,
		)
	}

	// Archiving more than the deck holds stops when the deck empties.
	(ArchiveCard{
		Zone:      Deck,
		Selection: Top{},
		Quantity:  Takes{N: Fixed(5)},
	}).Resolve(ctx)
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
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}

	if got := (ArchiveCard{
		Zone:      Discard,
		Selection: Top{},
	}).Text(); got != "archive the top card of your discard pile" {
		t.Errorf("text = %q", got)
	}
	if got := (ArchiveCard{
		Zone:      Discard,
		Selection: Top{},
		Quantity:  Takes{N: Fixed(2)},
	}).Text(); got != "archive the top 2 cards of your discard pile" {
		t.Errorf("plural text = %q", got)
	}

	// The most recently discarded card is the top and archives first.
	(ArchiveCard{
		Zone:      Discard,
		Selection: Top{},
	}).Resolve(ctx)
	if g.State.Archives[0].Count != 1 || g.State.Archives[0].IDs[0] != top {
		t.Errorf(
			"archived %v, want the top card %d",
			g.State.Archives[0].slice(),
			top,
		)
	}

	// Archiving more than the discard holds stops when the pile empties.
	(ArchiveCard{
		Zone:      Discard,
		Selection: Top{},
		Quantity:  Takes{N: Fixed(5)},
	}).Resolve(ctx)
	if g.State.Discard[0].Count != 0 {
		t.Errorf("discard should be empty, got %d", g.State.Discard[0].Count)
	}
	if g.State.Archives[0].Count != 2 {
		t.Errorf("archives count = %d, want 2", g.State.Archives[0].Count)
	}
}

// TestArchiveBottomSelection covers the Bottom positional selection, which picks
// the far end of an ordered zone.
func TestArchiveBottomSelection(t *testing.T) {
	g := NewGame("A", "B", 1)
	bottom := g.Register(testCreature("bottom", 1), 0)
	g.State.Discard[0].add(bottom)
	g.State.Discard[0].add(g.Register(testCreature("top", 1), 0))
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}

	if got := (ArchiveCard{
		Zone:      Discard,
		Selection: Bottom{},
	}).Text(); got != "archive the bottom card of your discard pile" {
		t.Errorf("text = %q", got)
	}
	if got := (ArchiveCard{
		Zone:      Deck,
		Selection: Bottom{},
		Quantity:  Takes{N: Fixed(2)},
	}).Text(); got != "archive the bottom 2 cards of your deck" {
		t.Errorf("plural text = %q", got)
	}

	// The bottom of the discard pile is the least recently discarded card.
	(ArchiveCard{
		Zone:      Discard,
		Selection: Bottom{},
	}).Resolve(ctx)
	if g.State.Archives[0].Count != 1 || g.State.Archives[0].IDs[0] != bottom {
		t.Errorf(
			"archived %v, want the bottom card %d",
			g.State.Archives[0].slice(),
			bottom,
		)
	}

	// Archiving more than the pile holds stops when it empties.
	(ArchiveCard{
		Zone:      Discard,
		Selection: Bottom{},
		Quantity:  Takes{N: Fixed(5)},
	}).Resolve(ctx)
	if g.State.Discard[0].Count != 0 {
		t.Errorf("discard should be empty, got %d", g.State.Discard[0].Count)
	}
}
func TestArchiveFromPlayEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddArtifact(NewCard("quest", Sanctum, Artifact, Rare), 0)
	knight := g.AddToBattleline(
		NewCard(
			"knight",
			Sanctum,
			Creature,
			Common,
			WithPower(4),
			WithTraits(Knight),
		),
		0,
	)
	nonKnight := g.AddToBattleline(
		NewCard(
			"cleric",
			Sanctum,
			Creature,
			Common,
			WithPower(4),
			WithTraits(Cleric),
		),
		0,
	)
	enemyKnight := g.AddToBattleline(
		NewCard(
			"enemy",
			Sanctum,
			Creature,
			Common,
			WithPower(4),
			WithTraits(Knight),
		),
		1,
	)
	ctx := &EffectContext{
		Resolver:   g,
		Source:     src,
		Controller: 0,
	}

	e := ArchiveFromPlay{
		Target: Target{Kind: TargetEachFriendlyCreature}.WithTrait(Knight),
	}
	if e.Text() != "archive each friendly Knight creature from play" {
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

// TestArchiveFromPlayArchivesBufferAndBuffedTogether checks that archiving a
// snapshot of creatures happens simultaneously: a creature buffing a damaged
// neighbor and that neighbor both go to archives at once, so the neighbor is
// archived rather than destroyed for the power it loses (Epic Quest + "Lion"
// Bautrem's neighbor).
func TestArchiveFromPlayArchivesBufferAndBuffedTogether(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddArtifact(NewCard("quest", Sanctum, Artifact, Rare), 0)
	buffer := g.AddToBattleline(
		NewCard("buffer", Sanctum, Creature, Common, WithPower(4),
			WithConstantAbility(ConstantAbility{
				PowerBonus: 2,
				Target:     Target{Kind: TargetEachCreature}.Neighboring(),
			})),
		0,
	)
	neighbor := g.AddToBattleline(
		NewCard("neighbor", Sanctum, Creature, Common, WithPower(3)),
		0,
	)
	// 4 damage is lethal at base power 3 but survivable at 5 with the buff.
	g.applyRawDamage(DamageTarget{
		ID:          neighbor,
		Amount:      4,
		IgnoreArmor: true,
	})
	if !g.inPlay(neighbor) {
		t.Fatal("neighbor should survive while the buffer is in play")
	}
	ctx := &EffectContext{
		Resolver:   g,
		Source:     src,
		Controller: 0,
	}

	e := ArchiveFromPlay{Target: Target{Kind: TargetEachFriendlyCreature}}
	e.Resolve(ctx)

	if g.inPlay(buffer) || g.inPlay(neighbor) {
		t.Error("both creatures should have left play")
	}
	if !g.State.Archives[0].contains(neighbor) {
		t.Error("the buffed neighbor should be archived, not destroyed")
	}
	if g.State.Discard[0].contains(neighbor) {
		t.Error("the neighbor must not also be in the discard pile")
	}
}

func TestArchiveFromPlayFriendlyInPlay(t *testing.T) {
	g := NewGame("A", "B", 1)
	art := g.AddArtifact(NewCard("relic", Mars, Artifact, Common), 0)
	g.AddToBattleline(NewCard("enemy", Mars, Creature, Common, WithPower(3)), 1)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}

	e := ArchiveFromPlay{
		Target: Target{Kind: TargetChosenFriendlyCreatureOrArtifact},
	}
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
	e := May{
		Do: ArchiveFromPlay{
			Target: Target{Kind: TargetChosenFriendlyCreatureOrArtifact},
		},
	}
	if !e.Do.(declinableEffect).declinable() {
		t.Fatal("a chosen-target ArchiveFromPlay should be declinable")
	}

	g := NewGame("A", "B", 1)
	ch := &cardDecliner{}
	g.SetChooser(0, ch)
	keep := g.AddArtifact(NewCard("keep", Mars, Artifact, Common), 0)
	doomed := g.AddArtifact(NewCard("doomed", Mars, Artifact, Common), 0)
	e.Resolve(&EffectContext{
		Resolver:   g,
		Controller: 0,
	})
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
	e.Resolve(&EffectContext{
		Resolver:   declined,
		Controller: 0,
	})
	if !declined.inPlay(other) {
		t.Error("a declined May should archive nothing")
	}
}

// A "you may reveal a creature from your hand and archive it" is one card
// choice, so the player picks the card directly instead of first answering Yes
// (Zyzzix the Many).
func TestMayDeclinableArchiveFromHand(t *testing.T) {
	e := May{
		Do: ArchiveCard{
			Zone: Hand,
			Selection: Chosen{
				Type:     Creature,
				Optional: true,
			},
			Revealed: true,
		},
	}
	if !e.Do.(declinableEffect).declinable() {
		t.Fatal("a single-count ArchiveCard should be declinable")
	}

	g := NewGame("A", "B", 1)
	ch := &cardDecliner{}
	g.SetChooser(0, ch)
	g.AddToHand(testCreature("keep", 1), 0)
	doomed := g.AddToHand(testCreature("doomed", 1), 0)
	e.Resolve(&EffectContext{
		Resolver:   g,
		Controller: 0,
	})
	if ch.asked != 1 {
		t.Errorf("declinable prompts = %d, want 1", ch.asked)
	}
	if g.State.Archives[0].Count != 1 || !g.State.Archives[0].contains(doomed) {
		t.Error("the chosen card should have been archived")
	}

	declined := NewGame("A", "B", 1)
	declined.SetChooser(0, &cardDecliner{decline: true})
	declined.AddToHand(testCreature("keep", 1), 0)
	e.Resolve(&EffectContext{
		Resolver:   declined,
		Controller: 0,
	})
	if declined.State.Archives[0].Count != 0 {
		t.Error("a declined May should archive nothing")
	}

	// A filter admitting nothing is not offered at all (Zyzzix the Many with no
	// creature in hand).
	empty := NewGame("A", "B", 1)
	empty.AddToHand(NewCard("tactic", Mars, Tactic, Common), 0)
	e.Resolve(&EffectContext{
		Resolver:   empty,
		Controller: 0,
	})
	if empty.State.Archives[0].Count != 0 {
		t.Error("an empty candidate set should archive nothing")
	}
}

func TestArchivesOfferedOnChooseHouse(t *testing.T) {
	g := NewGame("A", "B", 1)
	a := g.AddToHand(testCreature("a", 1), 0)
	(ArchiveCard{
		Zone:      Hand,
		Selection: Chosen{},
	}).
		Resolve(&EffectContext{
			Resolver:   g,
			Controller: 0,
		})
	if g.State.Archives[0].Count != 1 {
		t.Fatalf(
			"setup: archives count = %d, want 1",
			g.State.Archives[0].Count,
		)
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
	(ArchiveCard{
		Zone:      Hand,
		Selection: Chosen{},
	}).
		Resolve(&EffectContext{
			Resolver:   g,
			Controller: 0,
		})
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
	e := ArchiveCard{
		Zone: Hand,
		Selection: Chosen{
			Type:  Creature,
			House: namedHouse(Mars),
		},
		Revealed: true,
	}
	want := "reveal a Mars creature from your hand and archive it"
	if e.Text() != want {
		t.Errorf("text = %q, want %q", e.Text(), want)
	}
	if got := (ArchiveCard{
		Zone:      Hand,
		Selection: Chosen{},
		Quantity:  Takes{N: Fixed(2)},
	}).Text(); got != "archive 2 cards from your hand" {
		t.Errorf("plain plural text = %q", got)
	}
	if got := (ArchiveCard{
		Zone:      Hand,
		Selection: Chosen{Type: Artifact},
	}).Text(); got != "archive an artifact from your hand" {
		t.Errorf("artifact text = %q", got)
	}

	g := NewGame("A", "B", 1)
	tactic := g.AddToHand(NewCard("Trick", Mars, Tactic, Common), 0)
	offHouse := g.AddToHand(testCreature("Logos One", 3), 0)
	martian := g.AddToHand(
		NewCard("Martian", Mars, Creature, Common, WithPower(3)),
		0,
	)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}

	if !e.resolveGate(ctx) {
		t.Fatal("archiving a matching creature should report success")
	}
	if g.State.Archives[0].Count != 1 || g.State.Archives[0].IDs[0] != martian {
		t.Errorf(
			"archives = %v, want [%d]",
			g.State.Archives[0].slice(),
			martian,
		)
	}
	if len(g.Hand(0)) != 2 || g.Hand(0)[0] != tactic ||
		g.Hand(0)[1] != offHouse {
		t.Errorf("hand = %v, want the filtered-out cards", g.Hand(0))
	}
	if e.resolveGate(ctx) {
		t.Error(
			"with no matching card left the gate should report nothing archived",
		)
	}
}

// TestArchiveFromHandExceptHouse covers the "non-<house>" reveal (Information
// Officer Gray): only cards outside the excluded house are offered, and the text
// carries the "non-" qualifier.
func TestArchiveFromHandExceptHouse(t *testing.T) {
	e := ArchiveCard{
		Zone:      Hand,
		Selection: Chosen{House: exceptHouse(StarAlliance)},
		Revealed:  true,
	}
	want := "reveal a non-Star Alliance card from your hand and archive it"
	if e.Text() != want {
		t.Errorf("text = %q, want %q", e.Text(), want)
	}

	g := NewGame("A", "B", 1)
	ally := g.AddToHand(
		NewCard("Ally", StarAlliance, Creature, Common, WithPower(3)),
		0,
	)
	outsider := g.AddToHand(
		NewCard("Outsider", Logos, Creature, Common, WithPower(3)),
		0,
	)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}

	if !e.resolveGate(ctx) {
		t.Fatal("archiving a non-Star Alliance card should report success")
	}
	if g.State.Archives[0].Count != 1 ||
		g.State.Archives[0].IDs[0] != outsider {
		t.Errorf(
			"archives = %v, want [%d]",
			g.State.Archives[0].slice(),
			outsider,
		)
	}
	if len(g.Hand(0)) != 1 || g.Hand(0)[0] != ally {
		t.Errorf("hand = %v, want the excluded Star Alliance card", g.Hand(0))
	}
	if e.resolveGate(ctx) {
		t.Error(
			"with only excluded-house cards left the gate should archive nothing",
		)
	}
}

func TestArchiveRandomFromHand(t *testing.T) {
	if (ArchiveCard{
		Zone:      Hand,
		Selection: Random{},
	}).Text() != "archive a random card from your hand" {
		t.Errorf(
			"text = %q",
			(ArchiveCard{
				Zone:      Hand,
				Selection: Random{},
			}).Text(),
		)
	}
	two := ArchiveCard{
		Zone:      Hand,
		Selection: Random{},
		Quantity:  Takes{N: Fixed(2)},
	}
	if two.Text() != "archive 2 random cards from your hand" {
		t.Errorf("plural text = %q", two.Text())
	}

	g := NewGame("A", "B", 1)
	g.AddToHand(testCreature("c1", 1), 0)
	g.AddToHand(testCreature("c2", 1), 0)
	g.AddToHand(testCreature("c3", 1), 0)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}

	two.Resolve(ctx)
	if g.State.Archives[0].Count != 2 {
		t.Errorf("archives count = %d, want 2", g.State.Archives[0].Count)
	}
	if len(g.Hand(0)) != 1 {
		t.Errorf("hand = %v, want 1 card left", g.Hand(0))
	}

	// Archiving more than the hand holds stops when the hand empties.
	(ArchiveCard{
		Zone:      Hand,
		Selection: Random{},
		Quantity:  Takes{N: Fixed(5)},
	}).Resolve(ctx)
	if len(g.Hand(0)) != 0 {
		t.Errorf("hand should be empty, got %v", g.Hand(0))
	}
	if g.State.Archives[0].Count != 3 {
		t.Errorf("archives count = %d, want 3", g.State.Archives[0].Count)
	}
}

// TestArchiveFromHandOrAmount covers the alternate-amount tail ("archive a card,
// or 2 cards if …", Velum) — its validate, Text, and the archiveOrObject noun.
func TestArchiveFromHandOrAmount(t *testing.T) {
	e := ArchiveCard{
		Zone:      Hand,
		Selection: Chosen{},
		Or: OrAmount{
			Amount: 2,
			When:   ControlsNamed{Name: "Hyde"},
		},
	}
	if got := e.Text(); got != "archive a card from your hand, or 2 cards if you control Hyde" {
		t.Errorf("text = %q", got)
	}
	if err := e.validate(); err != nil {
		t.Errorf("validate = %v", err)
	}

	// A singular alternate exercises archiveOrObject's n==1 branch.
	one := ArchiveCard{
		Zone:      Hand,
		Selection: Chosen{},
		Quantity:  Takes{N: Fixed(2)},
		Or: OrAmount{
			Amount: 1,
			When:   ControlsNamed{Name: "Hyde"},
		},
	}
	if got := one.Text(); got != "archive 2 cards from your hand, or a card if you control Hyde" {
		t.Errorf("singular alt text = %q", got)
	}

	// A Per count leads the sentence (Dr. Milli).
	per := ArchiveCard{
		Zone:      Hand,
		Selection: Chosen{},
		Quantity:  Takes{N: ExcessCreatures{Player: Opponent}},
	}
	if got := per.Text(); got != "for each creature your opponent controls in excess of you, archive a card from your hand" {
		t.Errorf("per text = %q", got)
	}
}

// TestArchiveFromDiscardNamed covers pinning the choice to a named card (Hyde
// archives Velum): the text is the bare name, and the first card of that name is
// taken from the discard pile.
func TestArchiveFromDiscardNamed(t *testing.T) {
	e := ArchiveCard{
		Zone:      Discard,
		Selection: Named{Name: "Velum"},
	}
	if got := e.Text(); got != "archive Velum from your discard pile" {
		t.Errorf("named text = %q", got)
	}
	if got := (Named{Name: "Velum"}).noun(); got != "Velum" {
		t.Errorf("named noun = %q", got)
	}

	g := NewGame("A", "B", 1)
	g.AddToDiscard(NewCard("Other", Mars, Creature, Common, WithPower(1)), 0)
	velum := g.AddToDiscard(
		NewCard("Velum", Mars, Creature, Common, WithPower(1)),
		0,
	)
	e.Resolve(&EffectContext{
		Resolver:   g,
		Controller: 0,
	})
	if g.State.Archives[0].Count != 1 || !g.State.Archives[0].contains(velum) {
		t.Errorf(
			"archives = %v, want the named card [%d]",
			g.State.Archives[0].slice(),
			velum,
		)
	}

	// No card of the name archives nothing.
	g.AddToDiscard(NewCard("Nope", Mars, Creature, Common, WithPower(1)), 1)
	e.Resolve(&EffectContext{
		Resolver:   g,
		Controller: 1,
	})
	if g.State.Archives[1].Count != 0 {
		t.Error("no Velum in discard should archive nothing")
	}
}

func TestDiscardArchivesEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.State.ActivePlayer = 0 // player 0 discards player 1's (hidden) archives
	g.State.Archives[1].add(g.Register(testCreature("a", 1), 1))
	g.State.Archives[1].add(g.Register(testCreature("b", 1), 1))
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}

	e := DiscardArchives{Player: Opponent}
	if e.Text() != "your opponent discards each of their archived cards" {
		t.Errorf("text = %q", e.Text())
	}
	if (DiscardArchives{Player: Controller}).Text() != "discard each of your archived cards" {
		t.Errorf("controller text = %q", (DiscardArchives{Player: Controller}).Text())
	}

	e.Resolve(ctx)
	if g.State.Archives[1].Count != 0 {
		t.Error("opponent archives should be emptied")
	}
	if len(g.Discard(1)) != 2 {
		t.Errorf("opponent discard = %d, want 2", len(g.Discard(1)))
	}

	// Empty archives: nothing to move.
	e.Resolve(ctx)
	if len(g.Discard(1)) != 2 {
		t.Error("discarding empty archives should change nothing")
	}
}

// You can see your own archives, so you choose the order the cards are discarded.
func TestDiscardYourOwnArchivesIsOrdered(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.State.ActivePlayer = 0 // player 0 discards their own archives
	c1 := g.Register(testCreature("c1", 1), 0)
	c2 := g.Register(testCreature("c2", 1), 0)
	c3 := g.Register(testCreature("c3", 1), 0)
	g.State.Archives[0].add(c1)
	g.State.Archives[0].add(c2)
	g.State.Archives[0].add(c3)
	g.SetChooser(0, orderLastChooser{}) // choose to discard last-first
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}

	DiscardArchives{Player: Controller}.Resolve(ctx) // your own archives
	if g.State.Archives[0].Count != 0 {
		t.Error("own archives should be emptied")
	}
	got := g.Discard(0)
	want := []LocalID{c3, c2, c1} // the chosen (reversed) order
	if len(got) != 3 || got[0] != want[0] || got[1] != want[1] || got[2] != want[2] {
		t.Errorf("own discard order = %v, want %v (owner chooses the order)", got, want)
	}
}

func TestGainAemberPerOpponentArchivedCards(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.State.Archives[1].add(g.Register(testCreature("x", 1), 1))
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}
	e := GainAember{
		Player: Controller,
		Amount: 1,
		Per: CardsInZone{
			Zone:   Archives,
			Player: Opponent,
		},
	}
	if e.Text() != "for each card in your opponent's archives, gain 1 Æmber" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if g.Aember(0) != 1 {
		t.Errorf("aember = %d, want 1", g.Aember(0))
	}
}

func TestArchiveSource(t *testing.T) {
	if got := (ArchiveSource{}).Text(); got != "archive "+SelfName {
		t.Errorf("text = %q", got)
	}

	t.Run("archives an in-play source from play", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		src := g.AddToBattleline(testCreature("src", 3), 0)

		ArchiveSource{}.Resolve(&EffectContext{
			Resolver:   g,
			Controller: 0,
			Source:     src,
		})

		if !g.State.Archives[0].contains(src) {
			t.Error("an in-play source should be archived from play")
		}
	})

	t.Run("archives a resolving action instead of discarding it", func(t *testing.T) {
		g := started(t)
		idx := int(g.State.Hand[0].Count)
		id := g.AddToHand(
			NewCard(
				"Self Archive", Brobnar, Tactic, Common,
				WithAbility(TriggerAfterPlay, ArchiveSource{}),
			),
			0,
		)

		if err := g.PlayTactic(0, idx); err != nil {
			t.Fatalf("PlayTactic: %v", err)
		}

		if !g.State.Archives[0].contains(id) {
			t.Error("a self-archiving action should go to the archives")
		}
		if g.State.Discard[0].contains(id) {
			t.Error("a self-archiving action should not go to the discard pile")
		}
	})
}

// TestInvulnerableKeyword covers Ghostform's grant: an invulnerable creature
// takes no damage and cannot be destroyed, whether the keyword is printed on it
// or granted by an attached upgrade.
func TestInvulnerableKeyword(t *testing.T) {
	t.Run("takes no damage", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		ghost := g.AddToBattleline(testCreature("ghost", 3, WithKeywords(Invulnerable)), 0)
		g.dealDamage(0, DamageTarget{
			ID:     ghost,
			Amount: 5,
		})
		if got := g.Damage(ghost); got != 0 {
			t.Errorf("invulnerable creature took %d damage, want 0", got)
		}
		if !g.inPlay(ghost) {
			t.Fatal("invulnerable creature should survive the damage")
		}
	})

	t.Run("cannot be destroyed", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		ghost := g.AddToBattleline(testCreature("ghost", 3, WithKeywords(Invulnerable)), 0)
		g.destroyEach(0, []LocalID{ghost})
		if !g.inPlay(ghost) {
			t.Fatal("invulnerable creature should survive destruction")
		}
	})

	t.Run("granted by an attached upgrade", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		host := g.AddToBattleline(testCreature("host", 3), 0)
		attachUpgrade(g, host, NewCard("cloak", Brobnar, Upgrade, Rare,
			WithStatic(StaticModifier{Keywords: []Keyword{Invulnerable}})))
		g.destroyEach(0, []LocalID{host})
		if !g.inPlay(host) {
			t.Fatal("host granted invulnerable should survive destruction")
		}
	})
}

// TestArchiveGrantingUpgrade covers Ghostform archiving itself: the effect
// renders the {card} placeholder and sends the granting upgrade off its host to
// the owner's archives while the host stays in play.
func TestArchiveGrantingUpgrade(t *testing.T) {
	if got := (ArchiveGrantingUpgrade{}).Text(); got != "archive "+CardName {
		t.Errorf("text = %q", got)
	}

	g := NewGame("A", "B", 1)
	host := g.AddToBattleline(testCreature("host", 3), 0)
	up := attachUpgrade(
		g,
		host,
		NewCard("Ghostform", Brobnar, Upgrade, Rare, WithBonus(BonusAember)),
	)
	ArchiveGrantingUpgrade{}.Resolve(
		&EffectContext{
			Resolver:   g,
			Source:     host,
			Controller: 0,
			Upgrade:    up,
		},
	)
	if !containsID(g.Archives(0), up) {
		t.Errorf("archived upgrade should be in the owner's archives, got %v", g.Archives(0))
	}
	if _, ok := g.hostOf(up); ok {
		t.Error("archived upgrade should be detached from its host")
	}
	if !g.inPlay(host) {
		t.Error("the host should stay in play")
	}
}

// TestArchiveCardFromPurge covers the purge pile as an archive source — Universal
// Recycle Bin recovering a card the game had set aside for good. The purge pile is
// the one zone a card may name as a source but never as a destination.
func TestArchiveCardFromPurge(t *testing.T) {
	e := ArchiveCard{
		Zone:      Purged,
		Selection: Chosen{},
	}
	if err := e.validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}
	if got := e.Text(); got != "archive a card from your purge pile" {
		t.Errorf("Text = %q", got)
	}

	g := NewGame("A", "B", 1)
	id := g.Register(testCreature("p", 1), 0)
	g.State.Purge[0].add(id)
	e.Resolve(&EffectContext{
		Resolver:   g,
		Controller: 0,
	})

	if !g.State.Archives[0].contains(id) {
		t.Error("the purged card should have been archived")
	}
	if g.State.Purge[0].contains(id) {
		t.Error("the archived card should have left the purge pile")
	}
}

// TestArchiveCardFromPurgeEmpty covers an empty purge pile: nothing is archived
// and no choice is asked for.
func TestArchiveCardFromPurgeEmpty(t *testing.T) {
	g := NewGame("A", "B", 1)
	ArchiveCard{
		Zone:      Purged,
		Selection: Chosen{},
	}.
		Resolve(&EffectContext{
			Resolver:   g,
			Controller: 0,
		})
	if g.State.Archives[0].Count != 0 {
		t.Errorf("archives = %v, want empty", g.State.Archives[0].slice())
	}
}
