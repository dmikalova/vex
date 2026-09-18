package engine

import "testing"

func TestReturnNamedToHand(t *testing.T) {
	e := PutNamedIntoHand{Name: "Urchin"}
	if e.Text() != "put an Urchin from play or from your discard pile into your hand" {
		t.Errorf("text = %q", e.Text())
	}

	urchin := func(g *Game, player int) LocalID {
		return g.Register(
			NewCard("Urchin", Shadows, Creature, Common, WithPower(1), WithTraits(Elf, Thief)),
			player,
		)
	}

	t.Run("from play", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		src := g.AddToBattleline(testCreature("faygin", 3), 0)
		urch := g.AddToBattleline(
			NewCard("Urchin", Shadows, Creature, Common, WithPower(1), WithTraits(Elf, Thief)),
			0,
		)
		g.AddToBattleline(
			testCreature("other", 4),
			0,
		) // different name: filtered out
		g.State.Discard[0].add(
			g.Register(NewCard("junk", Dis, Tactic, Common), 0),
		) // different name in discard: filtered out
		ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

		e.Resolve(ctx) // the sole Urchin candidate is auto-chosen
		if g.inPlay(urch) {
			t.Error("the Urchin should leave play")
		}
		if !g.State.Hand[0].contains(urch) {
			t.Error("the Urchin should be in the controller's hand")
		}
	})

	t.Run("from discard", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		src := g.AddToBattleline(testCreature("faygin", 3), 0)
		urch := urchin(g, 0)
		g.State.Discard[0].add(urch)
		ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

		e.Resolve(ctx)
		if g.State.Discard[0].contains(urch) {
			t.Error("the Urchin should leave the discard pile")
		}
		if !g.State.Hand[0].contains(urch) {
			t.Error("the Urchin should be in the controller's hand")
		}
	})

	t.Run("declined", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		src := g.AddToBattleline(testCreature("faygin", 3), 0)
		u1 := g.AddToBattleline(
			NewCard("Urchin", Shadows, Creature, Common, WithPower(1), WithTraits(Elf, Thief)),
			0,
		)
		g.AddToBattleline(
			NewCard("Urchin", Shadows, Creature, Common, WithPower(1), WithTraits(Elf, Thief)),
			0,
		)
		g.SetChooser(0, orderRejectChooser{})
		ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

		e.Resolve(ctx)
		if !g.inPlay(u1) {
			t.Error("nothing should move when the choice is declined")
		}
	})
}

func TestMoveFromPlayToDeck(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	myArt := g.AddArtifact(NewCard("myrelic", Brobnar, Artifact, Rare), 0)
	enemyArt := g.AddArtifact(NewCard("enemyrelic", Brobnar, Artifact, Rare), 1)
	g.State.Cards[myArt].Exhausted = true
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	e := PutFromPlay{Target: Target{Kind: TargetEachArtifact}, Destination: ToTopOfDeck}
	if e.Text() != "put each artifact on top of its owner's deck" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)

	if g.State.Artifacts[0].Count != 0 || g.State.Artifacts[1].Count != 0 {
		t.Errorf(
			"artifact rows not cleared: %d %d",
			g.State.Artifacts[0].Count,
			g.State.Artifacts[1].Count,
		)
	}
	if g.State.Deck[0].Count != 1 || g.State.Deck[0].IDs[0] != myArt {
		t.Errorf("player 0 deck top = %v, want %d", g.State.Deck[0].IDs[0], myArt)
	}
	if g.State.Deck[1].Count != 1 || g.State.Deck[1].IDs[0] != enemyArt {
		t.Errorf("player 1 deck top = %v, want %d", g.State.Deck[1].IDs[0], enemyArt)
	}
	if g.State.Cards[myArt].Exhausted {
		t.Errorf("returned artifact should be readied")
	}
}

func TestMoveFromPlayToHand(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 3), 0)
	g.State.Cards[src].Damage = 2
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	e := PutFromPlay{Target: Target{Kind: TargetThisCreature}, Destination: ToHand}
	if e.Text() != "put "+SelfName+" into its owner's hand" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)

	if g.State.Battleline[0].Count != 0 {
		t.Errorf("battleline not cleared: %d", g.State.Battleline[0].Count)
	}
	if g.State.Hand[0].Count != 1 || g.State.Hand[0].IDs[0] != src {
		t.Errorf(
			"hand = count %d id %v, want 1 / %d",
			g.State.Hand[0].Count,
			g.State.Hand[0].IDs[0],
			src,
		)
	}
	// Moving to hand clears the per-match state the card accrued in play.
	if g.State.Cards[src].Damage != 0 {
		t.Errorf("damage after move = %d, want 0", g.State.Cards[src].Damage)
	}
}

func TestMoveFromPlayToArchives(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 3), 0)
	attachUpgrade(g, src, NewCard("plating", Mars, Upgrade, Common))
	g.State.Cards[src].Damage = 1
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	e := PutFromPlay{Target: Target{Kind: TargetThisCreature}, Destination: ToArchives}
	if e.Text() != "put "+SelfName+" into its owner's archives" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)

	if g.inPlay(src) {
		t.Error("the creature should have left play")
	}
	if g.State.Archives[0].Count != 1 || g.State.Archives[0].IDs[0] != src {
		t.Errorf(
			"creature should be archived, got %v",
			g.State.Archives[0].IDs[:g.State.Archives[0].Count],
		)
	}
	if len(g.Discard(0)) != 1 { // the upgrade sheds to the discard pile
		t.Errorf("upgrade should be discarded; discard = %v", g.Discard(0))
	}
	if g.State.Cards[src].Damage != 0 {
		t.Error("archiving should clear the creature's in-play state")
	}
}

func TestMoveFromPlayToDeckShuffled(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("chrono", 2), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	e := PutFromPlay{Target: Target{Kind: TargetThisCreature}, Destination: ToDeckShuffled}
	if e.Text() != "shuffle {self} into its owner's deck" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if g.inPlay(src) {
		t.Error("the creature should leave play")
	}
	if g.State.Deck[0].Count != 1 || !g.State.Deck[0].contains(src) {
		t.Errorf("the creature should be shuffled into its owner's deck")
	}
}

func TestPutFromPlayGate(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	other := g.AddToBattleline(testCreature("other", 1), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	e := PutFromPlay{Target: Target{Kind: TargetThisCreature}, Destination: ToHand}
	if !e.resolveGate(ctx) {
		t.Error("resolveGate should report true when a card moved")
	}
	if !ctx.HasIt || ctx.It != src {
		t.Errorf("ctx.It = %v (HasIt %v), want %d", ctx.It, ctx.HasIt, src)
	}

	ctx = &EffectContext{Resolver: g, Source: other, Controller: 0}
	e = PutFromPlay{Target: Target{Kind: TargetThisCreature}.Damaged(), Destination: ToHand}
	if e.resolveGate(ctx) {
		t.Error("resolveGate should report false when nothing moved")
	}
	if ctx.HasIt {
		t.Error("ctx.It should not be set when nothing moved")
	}
}

func TestMoveFromPlayValidate(t *testing.T) {
	this := Target{Kind: TargetThisCreature}
	for _, d := range []Destination{ToHand, ToTopOfDeck, ToDeckShuffled, ToArchives} {
		if err := (PutFromPlay{Target: this, Destination: d}).validate(); err != nil {
			t.Errorf("destination %d should be valid, got %v", d.zone, err)
		}
	}
	if err := (PutFromPlay{Target: this}).validate(); err == nil {
		t.Error("an unset destination should be rejected")
	}
	if err := (PutFromPlay{Target: this, Destination: ToBottomOfDeck}).validate(); err == nil {
		t.Error("an unsupported destination should be rejected")
	}
	if err := (PutFromPlay{Destination: ToHand}).validate(); err == nil {
		t.Error("an unset target should be rejected")
	}
	if err := (PutFromPlay{Target: this, Destination: ToArchives, WithUpgrades: true}).validate(); err == nil {
		t.Error("WithUpgrades should be rejected for a non-hand destination")
	}
	if err := (PutFromPlay{Target: this, Destination: ToHand, WithUpgrades: true}).validate(); err != nil {
		t.Errorf("WithUpgrades to hand should be valid, got %v", err)
	}
}

func TestPutFromPlayWithUpgrades(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 3), 0)
	up := attachUpgrade(g, src, NewCard("plating", Mars, Upgrade, Common))
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	e := PutFromPlay{
		Target:       Target{Kind: TargetThisCreature},
		Destination:  ToHand,
		WithUpgrades: true,
	}
	if e.Text() != "put "+SelfName+" and each upgrade attached to it into its owner's hand" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)

	if !g.State.Hand[0].contains(src) {
		t.Error("the creature should be in its owner's hand")
	}
	if !g.State.Hand[0].contains(up) {
		t.Error("the upgrade should follow the creature to hand, not the discard pile")
	}
	if len(g.Discard(0)) != 0 {
		t.Errorf("nothing should be discarded; discard = %v", g.Discard(0))
	}
}

func TestPutChosen(t *testing.T) {
	eachArt := Target{Kind: TargetEachArtifact}
	// Text renders "up to N" with the plural noun for each destination.
	cases := map[Destination]string{
		ToHand:         "put up to 3 artifacts into their owners' hands",
		ToTopOfDeck:    "put up to 3 artifacts on top of their owners' decks",
		ToDeckShuffled: "shuffle up to 3 artifacts into their owners' decks",
		ToArchives:     "put up to 3 artifacts into their owners' archives",
	}
	for dest, want := range cases {
		e := PutChosen{Amount: 3, UpTo: true, Target: eachArt, Destination: dest}
		if got := e.Text(); got != want {
			t.Errorf("text(%d) = %q, want %q", dest.zone, got, want)
		}
	}

	// Without UpTo the count is mandatory, so "up to" drops out; a single card
	// reads as the indefinite noun.
	mandatory := PutChosen{Amount: 2, Target: eachArt, Destination: ToDeckShuffled}
	if got := mandatory.Text(); got != "shuffle 2 artifacts into their owners' decks" {
		t.Errorf("mandatory text = %q", got)
	}
	one := PutChosen{Amount: 1, Target: eachArt, Destination: ToHand}
	if got := one.Text(); got != "put an artifact into its owner's hand" {
		t.Errorf("single text = %q", got)
	}

	// validate rejects an unset target, a non-positive Count, and a bad destination.
	if err := (PutChosen{Amount: 3, Destination: ToHand}).validate(); err == nil {
		t.Error("unset target should be rejected")
	}
	if err := (PutChosen{Target: eachArt, Destination: ToHand}).validate(); err == nil {
		t.Error("non-positive Count should be rejected")
	}
	if err := (PutChosen{Amount: 1, Target: eachArt, Destination: ToBottomOfDeck}).validate(); err == nil {
		t.Error("unsupported destination should be rejected")
	}
	if err := (PutChosen{Amount: 3, Target: eachArt, Destination: ToHand}).validate(); err != nil {
		t.Errorf("valid PutChosen = %v", err)
	}

	g := NewGame("A", "B", 1)
	a1 := g.AddArtifact(exAutocannon(), 0)
	a2 := g.AddArtifact(exAutocannon(), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	// Only two artifacts exist, so the loop stops when none remain (below Count).
	PutChosen{Amount: 3, UpTo: true, Target: eachArt, Destination: ToHand}.Resolve(ctx)
	if g.inPlay(a1) || g.inPlay(a2) {
		t.Error("both artifacts should have left play")
	}
	if g.State.Hand[0].Count != 1 || g.State.Hand[1].Count != 1 {
		t.Errorf(
			"hands = %d/%d, want 1/1 (each returned to its owner)",
			g.State.Hand[0].Count,
			g.State.Hand[1].Count,
		)
	}

	// A mandatory count is not declinable, and stops when the pool empties.
	g3 := NewGame("A", "B", 1)
	solo := g3.AddArtifact(exAutocannon(), 0)
	mandatory.Resolve(&EffectContext{Resolver: g3, Controller: 0})
	if g3.inPlay(solo) {
		t.Error("a mandatory choice should have moved the only artifact")
	}

	// Choosing "Done" (the option past the sole artifact) stops early, leaving it.
	g2 := NewGame("A", "B", 1)
	art := g2.AddArtifact(exAutocannon(), 0)
	g2.SetChooser(0, optionPicker{idx: 1}) // index 0 is the artifact, 1 is "Done"
	PutChosen{
		Amount:      3,
		UpTo:        true,
		Target:      eachArt,
		Destination: ToHand,
	}.Resolve(
		&EffectContext{Resolver: g2, Controller: 0},
	)
	if !g2.inPlay(art) {
		t.Error("choosing Done should leave the artifact in play")
	}
}

// A card the settle-before-choice boundary (ADR 0029) destroys after the pool was
// gathered — its buff left with an earlier pick — is skipped rather than moved into
// a second zone (ADR 0030), exactly as PutFromPlay skips one an earlier move took
// out. Two enemy creatures make the choice present (so it settles); the 0-power one
// is swept there, and the stale pool must not abduct it on top of its discard.
func TestPutChosenSkipsACardSettledOutOfPlay(t *testing.T) {
	g := NewGame("A", "B", 1)
	dead := g.AddToBattleline(testCreature("dead", 0), 1)
	alive := g.AddToBattleline(testCreature("alive", 3), 1)
	g.SetChooser(0, FirstChooser{})
	PutChosen{
		Amount:      2,
		Target:      Target{Kind: TargetEachEnemyCreature},
		Destination: ToArchives.Yours(),
	}.Resolve(&EffectContext{Resolver: g, Controller: 0})

	if g.State.Archives[0].contains(dead) {
		t.Error("a creature settled out of play must not be abducted into archives")
	}
	if !g.State.Discard[1].contains(dead) {
		t.Error("the 0-power creature should have been destroyed to its owner's discard")
	}
	if !g.State.Archives[0].contains(alive) {
		t.Error("the surviving creature should have been abducted into archives")
	}
}

// Shuffling several creatures into their owners' decks with one effect narrates
// one grouped, source-attributed line per owner instead of a passive line each.
func TestPutChosenGroupsShufflesByOwnerInLog(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddArtifact(exAutocannon(), 0)
	a1 := g.AddToBattleline(testCreature("a1", 3), 0)
	a2 := g.AddToBattleline(testCreature("a2", 3), 0)
	b1 := g.AddToBattleline(testCreature("b1", 3), 1)
	g.SetChooser(0, FirstChooser{})
	PutChosen{
		Amount:      3,
		Target:      Target{Kind: TargetEachCreature},
		Destination: ToDeckShuffled,
	}.Resolve(&EffectContext{Resolver: g, Source: src, Controller: 0})

	if g.inPlay(a1) || g.inPlay(a2) || g.inPlay(b1) {
		t.Fatal("all three creatures should have left play")
	}
	var grouped int
	for _, rec := range g.Log {
		if _, ok := rec.Entry.(CardsShuffledIntoDeckBy); ok {
			grouped++
		}
		if _, ok := rec.Entry.(CardShuffledIntoDeck); ok {
			t.Error("a batched shuffle should not also narrate a passive per-card line")
		}
	}
	// One line for P0's two creatures, one for P1's one creature.
	if grouped != 2 {
		t.Errorf("grouped shuffle lines = %d, want 2 (one per owner)", grouped)
	}
}

// A "you may put a creature into its owner's hand" is one clickable creature, so
// May drives it by the click rather than by a Yes/No.
func TestPutFromPlayDeclinable(t *testing.T) {
	chosen := PutFromPlay{Target: Target{Kind: TargetChosenCreature}, Destination: ToHand}
	if !chosen.declinable() {
		t.Error("a chosen PutFromPlay should be declinable")
	}
	if (PutFromPlay{Target: Target{Kind: TargetEachCreature}, Destination: ToHand}).
		declinable() {
		t.Error("an untargeted PutFromPlay should not be declinable")
	}

	empty := NewGame("A", "B", 1)
	if !chosen.vacuous(&EffectContext{Resolver: empty, Controller: 0}) {
		t.Error("a PutFromPlay with no creature to move should be vacuous")
	}

	taken := NewGame("A", "B", 1)
	taken.SetChooser(0, &cardDecliner{})
	foe := taken.AddToBattleline(testCreature("foe", 3), 1)
	if !chosen.resolveOptional(&EffectContext{Resolver: taken, Controller: 0}) {
		t.Error("clicking the creature should report it moved")
	}
	if onAnyLine(taken, foe) {
		t.Error("the clicked creature should have left play")
	}

	declined := NewGame("A", "B", 1)
	declined.SetChooser(0, &cardDecliner{decline: true})
	stayed := declined.AddToBattleline(testCreature("stayed", 3), 1)
	if chosen.resolveOptional(&EffectContext{Resolver: declined, Controller: 0}) {
		t.Error("declining should report nothing moved")
	}
	if !onAnyLine(declined, stayed) {
		t.Error("a declined PutFromPlay should move nothing")
	}
}

// An abduction puts an enemy creature into the abductor's archives. Your archives
// are one of the three zones that may hold an enemy card, so the card needs no
// rider: leaving them, it goes to its owner's matching zone.

func TestAbductionText(t *testing.T) {
	single := PutFromPlay{
		Target:      Target{Kind: TargetChosenEnemyCreature},
		Destination: ToArchives.Yours(),
	}
	if got := single.Text(); got != "put an enemy creature into your archives" {
		t.Errorf("text = %q", got)
	}

	many := PutChosen{
		Amount:      3,
		UpTo:        true,
		Target:      Target{Kind: TargetEachEnemyCreature},
		Destination: ToArchives.Yours(),
	}
	if got := many.Text(); got != "put up to 3 enemy creatures into your archives" {
		t.Errorf("plural text = %q", got)
	}
}

func TestAbductionResolve(t *testing.T) {
	abduct := PutFromPlay{
		Target:      Target{Kind: TargetChosenEnemyCreature},
		Destination: ToArchives.Yours(),
	}

	// setup abducts player 1's sole creature into player 0's archives.
	setup := func(t *testing.T) (*Game, LocalID) {
		t.Helper()
		g := NewGame("A", "B", 1)
		prey := g.AddToBattleline(testCreature("prey", 3), 1)
		abduct.Resolve(&EffectContext{Resolver: g, Controller: 0})
		if !g.State.Archives[0].contains(prey) {
			t.Fatalf(
				"prey should sit in the abductor's archives, got %v",
				g.State.Archives[0].slice(),
			)
		}
		return g, prey
	}

	t.Run("taken into hand", func(t *testing.T) {
		g, prey := setup(t)
		g.offerArchives(0) // the default chooser answers "Yes"
		if !g.State.Hand[1].contains(prey) {
			t.Error("prey should go to its owner's hand, not the abductor's")
		}
		if g.State.Hand[0].contains(prey) {
			t.Error("prey should not enter the abductor's hand")
		}
	})

	t.Run("archives discarded", func(t *testing.T) {
		g, prey := setup(t)
		g.discardArchives(0)
		if !g.State.Discard[1].contains(prey) {
			t.Error("a discarded abductee should go to its owner's discard pile")
		}
		if g.State.Discard[0].contains(prey) {
			t.Error("prey should not enter the abductor's discard pile")
		}
	})

	t.Run("ordinary archived card is untouched", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		own := g.AddToHand(testCreature("own", 3), 0)
		g.archiveFromHand(0, own)
		g.offerArchives(0)
		if !g.State.Hand[0].contains(own) {
			t.Error("a card archived normally returns to its own controller's hand")
		}
	})
}

// TestPutItIntoHand covers each branch of the effect: its rendered text, the
// no-context no-op, recovering a destroyed creature from the discard pile, and
// putting a creature that is still in play into its owner's hand.
func TestPutItIntoHand(t *testing.T) {
	if got := (PutItIntoHand{}).Text(); got != "put it into its owner's hand" {
		t.Errorf("Text = %q, want %q", got, "put it into its owner's hand")
	}

	// No card in context: the effect does nothing.
	g := started(t)
	(PutItIntoHand{}).Resolve(&EffectContext{Resolver: g})

	// A destroyed creature (in the discard) is recovered to its owner's hand.
	dead := g.AddToBattleline(testCreature("dead", 3), 1)
	g.DestroyEach(0, []LocalID{dead})
	if g.inPlay(dead) {
		t.Fatal("dead should have left play")
	}
	(PutItIntoHand{}).Resolve(&EffectContext{Resolver: g, It: dead, HasIt: true})
	if !handContains(g, 1, dead) {
		t.Error("destroyed creature should be recovered to its owner's hand")
	}

	// A creature still in play is returned straight from the battleline.
	live := g.AddToBattleline(testCreature("live", 3), 1)
	(PutItIntoHand{}).Resolve(&EffectContext{Resolver: g, It: live, HasIt: true})
	if g.inPlay(live) {
		t.Fatal("live should have left play")
	}
	if !handContains(g, 1, live) {
		t.Error("in-play creature should be returned to its owner's hand")
	}
}

func handContains(g *Game, player int, id LocalID) bool {
	for _, h := range g.Hand(player) {
		if h == id {
			return true
		}
	}
	return false
}
