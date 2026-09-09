package engine

import (
	"strings"
	"testing"
)

func TestGainAemberEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	e := GainAember{Player: Controller, Amount: 2}
	if e.Text() != "gain 2 Æmber" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if g.State.Aember[0] != 2 {
		t.Errorf("aember = %d, want 2", g.State.Aember[0])
	}

	foe := GainAember{Player: Opponent, Amount: 1}
	if foe.Text() != "your opponent gains 1 Æmber" {
		t.Errorf("enemy text = %q", foe.Text())
	}
	foe.Resolve(ctx)
	if g.State.Aember[1] != 1 {
		t.Errorf("opponent aember = %d, want 1", g.State.Aember[1])
	}

	owner := GainAember{Player: ItsOwner, Amount: 1}
	if owner.Text() != "its owner gains 1 Æmber" {
		t.Errorf("owner text = %q", owner.Text())
	}
}

func TestGainAemberItsController(t *testing.T) {
	g := NewGame("A", "B", 1)
	foe := g.AddToBattleline(testCreature("foe", 1), 1)
	_ = foe
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := GainAember{Player: ItsController, Amount: 1}
	if e.Text() != "its controller gains 1 Æmber" {
		t.Errorf("its-controller text = %q", e.Text())
	}

	// Destroy the enemy creature, then its controller (player 1) gains 1 Æmber.
	// The controller is captured before the creature leaves play, so the gate pays
	// the right side even though the destroyed card would revert to its owner.
	gate := Then{First: Destroy{Target: Target{Kind: TargetChosenCreature}}, Result: e}
	gate.Resolve(ctx)
	if g.State.Aember[1] != 1 {
		t.Errorf("controller aember = %d, want 1", g.State.Aember[1])
	}
	if g.State.Aember[0] != 0 {
		t.Errorf("own aember = %d, want 0", g.State.Aember[0])
	}
}

func TestGainAemberPerCount(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.State.Keys[1] = 2 // opponent has forged 2 keys
	ctx := &EffectContext{Resolver: g, Controller: 0}
	e := GainAember{Player: Controller, Amount: 1, Per: OpponentForgedKeys{}}
	if e.Text() != "for each forged key your opponent has, gain 1 Æmber" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if g.Aember(0) != 2 { // 1 per each of the 2 forged keys
		t.Errorf("aember = %d, want 2", g.Aember(0))
	}
}

func TestGainAemberPerArchivedCards(t *testing.T) {
	g := NewGame("A", "B", 1)
	for i := 0; i < 3; i++ {
		g.State.Archives[0].add(g.Register(testCreature("a", 1), 0))
	}
	ctx := &EffectContext{Resolver: g, Controller: 0}
	e := GainAember{Player: Controller, Amount: 1, Per: CardsInArchives{Player: Controller}}
	if e.Text() != "for each card in your archives, gain 1 Æmber" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if g.Aember(0) != 3 { // 1 per each of the 3 archived cards
		t.Errorf("aember = %d, want 3", g.Aember(0))
	}
}

func TestGainAemberMax(t *testing.T) {
	g := NewGame("A", "B", 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	e := GainAember{Player: Controller, Amount: 1, Per: Fixed(5)}
	e.Resolve(ctx)
	if g.Aember(0) != 5 {
		t.Errorf("aember = %d, want 5 (uncapped)", g.Aember(0))
	}
}

func TestLoseAemberEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	g.State.Aember[0] = 3
	g.State.Aember[1] = 1

	self := LoseAember{Player: Controller, Amount: 2}
	if self.Text() != "lose 2 Æmber" {
		t.Errorf("self text = %q", self.Text())
	}
	self.Resolve(ctx)
	if g.State.Aember[0] != 1 {
		t.Errorf("self aember = %d, want 1", g.State.Aember[0])
	}

	foe := LoseAember{Player: Opponent, Amount: 4}
	if foe.Text() != "your opponent loses 4 Æmber" {
		t.Errorf("enemy text = %q", foe.Text())
	}
	foe.Resolve(ctx) // opponent has only 1; floors at 0
	if g.State.Aember[1] != 0 {
		t.Errorf("opponent aember = %d, want 0", g.State.Aember[1])
	}
}

func TestLoseAemberItsOwner(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	foe := g.AddToBattleline(testCreature("foe", 1), 1)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0, It: foe, HasIt: true}
	g.State.Aember[1] = 3

	e := LoseAember{Player: ItsOwner, Amount: 1}
	if e.Text() != "its controller loses 1 Æmber" {
		t.Errorf("its-owner text = %q", e.Text())
	}
	e.Resolve(ctx)
	if g.State.Aember[1] != 2 {
		t.Errorf("owner aember = %d, want 2", g.State.Aember[1])
	}
}

func TestLoseAemberAllBut(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.State.Aember[0] = 8 // over the cap: loses 3, left with 5
	g.State.Aember[1] = 3 // at or below the cap: loses 0
	ctx := &EffectContext{Resolver: g, Controller: 0}

	e := LoseAember{Player: EachPlayer, By: AllBut(5)}
	if e.Text() != "each player with 6 Æmber or more loses all but 5 Æmber" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if g.State.Aember[0] != 5 {
		t.Errorf("player 0 aember = %d, want 5 (reduced)", g.State.Aember[0])
	}
	if g.State.Aember[1] != 3 {
		t.Errorf("player 1 aember = %d, want 3 (unchanged)", g.State.Aember[1])
	}
}

func TestLoseAemberHalfAndEachPlayer(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.State.Aember[0] = 5 // loses 2 (floor 5/2), keeps 3
	g.State.Aember[1] = 4 // loses 2, keeps 2
	ctx := &EffectContext{Resolver: g, Controller: 0}

	each := LoseAember{Player: EachPlayer, By: Half}
	if each.Text() != "each player loses half of their Æmber, rounded down" {
		t.Errorf("each text = %q", each.Text())
	}
	each.Resolve(ctx)
	if g.State.Aember[0] != 3 || g.State.Aember[1] != 2 {
		t.Errorf("aember after each-half = %d/%d, want 3/2", g.State.Aember[0], g.State.Aember[1])
	}

	// A controller losing half uses the "your" possessive.
	me := LoseAember{Player: Controller, By: Half}
	if me.Text() != "lose half of your Æmber, rounded down" {
		t.Errorf("controller-half text = %q", me.Text())
	}
	me.Resolve(ctx) // player 0 has 3 -> loses 1 -> 2
	if g.State.Aember[0] != 2 {
		t.Errorf("controller-half aember = %d, want 2", g.State.Aember[0])
	}
}

// TestGainAemberEqualToAndHalfPower covers gaining Æmber equal to a count that
// reads the chosen creature's power, halved and rounded down (The Flex).
func TestGainAemberEqualToAndHalfPower(t *testing.T) {
	g := NewGame("A", "B", 1)
	beefy := g.AddToBattleline(testCreature("beefy", 5), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0, It: beefy, HasIt: true}

	e := GainAemberEqualTo{Player: Controller, Count: PowerOfChosen{Of: Half}}
	if got := e.Text(); got != "gain Æmber equal to half its power, rounded down" {
		t.Errorf("text = %q", got)
	}
	e.Resolve(ctx)
	if g.State.Aember[0] != 2 {
		t.Errorf("aember = %d, want 2 (half of power 5, rounded down)", g.State.Aember[0])
	}

	// With no creature in context the count is zero and nothing is gained.
	if got := (PowerOfChosen{Of: Half}).Value(
		&EffectContext{Resolver: g, Controller: 0},
	); got != 0 {
		t.Errorf("PowerOfChosen{Of: Half} with no It = %d, want 0", got)
	}
	if got := (PowerOfChosen{Of: Half}).CountText(); got != "half its power, rounded down" {
		t.Errorf("count text = %q", got)
	}

	// The opponent form uses the "your opponent gains" verb.
	if got := (GainAemberEqualTo{Player: Opponent, Count: PowerOfChosen{Of: Half}}).Text(); got !=
		"your opponent gains Æmber equal to half its power, rounded down" {
		t.Errorf("opponent text = %q", got)
	}
	// The each-player form uses the "each player gains" verb.
	if got := (GainAemberEqualTo{Player: EachPlayer, Count: PowerOfChosen{Of: Half}}).Text(); got !=
		"each player gains Æmber equal to half its power, rounded down" {
		t.Errorf("each-player text = %q", got)
	}
	// Validation rejects an unset player or count, and accepts a fully set effect.
	if (GainAemberEqualTo{Count: PowerOfChosen{Of: Half}}).validate() == nil {
		t.Error("unset player should be rejected")
	}
	if (GainAemberEqualTo{Player: Controller}).validate() == nil {
		t.Error("unset count should be rejected")
	}
	if (GainAemberEqualTo{Player: Controller, Count: PowerOfChosen{Of: Half}}).validate() != nil {
		t.Error("a fully set effect should be valid")
	}
	// With no creature in context the count is zero and nothing is gained.
	before := g.State.Aember[0]
	GainAemberEqualTo{Player: Controller, Count: PowerOfChosen{Of: Half}}.Resolve(
		&EffectContext{Resolver: g, Controller: 0},
	)
	if g.State.Aember[0] != before {
		t.Errorf("zero count gained Æmber: %d, want %d", g.State.Aember[0], before)
	}
}

// TestLoseAemberEqualTo covers each player losing Æmber equal to half a creature's
// power, floored, and never below zero (Power of Fire).
func TestLoseAemberEqualTo(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.State.Aember[0] = 1 // loses 1, floored at zero (would lose 2)
	g.State.Aember[1] = 5 // loses 2
	beefy := g.AddToBattleline(testCreature("beefy", 5), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0, It: beefy, HasIt: true}

	e := LoseAemberEqualTo{Player: EachPlayer, Count: PowerOfChosen{Of: Half}}
	if got := e.Text(); got != "each player loses Æmber equal to half its power, rounded down" {
		t.Errorf("text = %q", got)
	}
	e.Resolve(ctx)
	if g.State.Aember[0] != 0 || g.State.Aember[1] != 3 {
		t.Errorf("aember = %d/%d, want 0/3", g.State.Aember[0], g.State.Aember[1])
	}

	// The opponent form uses the "your opponent loses" verb.
	if got := (LoseAemberEqualTo{Player: Opponent, Count: PowerOfChosen{Of: Half}}).Text(); got !=
		"your opponent loses Æmber equal to half its power, rounded down" {
		t.Errorf("opponent text = %q", got)
	}
	// Validation rejects an unset player or count, and accepts a fully set effect.
	if (LoseAemberEqualTo{Count: PowerOfChosen{Of: Half}}).validate() == nil {
		t.Error("unset player should be rejected")
	}
	if (LoseAemberEqualTo{Player: Controller}).validate() == nil {
		t.Error("unset count should be rejected")
	}
	if (LoseAemberEqualTo{Player: Controller, Count: PowerOfChosen{Of: Half}}).validate() != nil {
		t.Error("a fully set effect should be valid")
	}
	// A zero count loses nothing.
	before := g.State.Aember[1]
	LoseAemberEqualTo{Player: Opponent, Count: PowerOfChosen{Of: Half}}.Resolve(
		&EffectContext{Resolver: g, Controller: 0},
	)
	if g.State.Aember[1] != before {
		t.Errorf("zero count lost Æmber: %d, want %d", g.State.Aember[1], before)
	}
}

// TestGainAemberEqualToCaptured covers a gain-equal-to count that a continuous
// replacement (Ether Spider) captures instead of adding to the pool.
func TestGainAemberEqualToCaptured(t *testing.T) {
	g := started(t)
	src := g.AddToBattleline(testCreature("src", 4), 0)
	spider := g.AddToBattleline(testEtherSpider(), 1)
	GainAemberEqualTo{Player: Controller, Count: PowerOfChosen{Of: Half}}.Resolve(
		&EffectContext{Resolver: g, Controller: 0, It: src, HasIt: true},
	)
	if g.Aember(0) != 0 {
		t.Errorf("player Æmber = %d, want 0 (captured)", g.Aember(0))
	}
	if g.AmberOn(spider) != 2 {
		t.Errorf("spider Æmber = %d, want 2 (half of power 4)", g.AmberOn(spider))
	}
}

func TestLoseAemberValidate(t *testing.T) {
	if err := (LoseAember{Player: Controller, Amount: 2, By: Half}).validate(); err == nil {
		t.Error("setting both Amount and By should be rejected")
	}
	if err := (LoseAember{Player: Controller, By: Half}).validate(); err != nil {
		t.Errorf("By alone should be valid, got %v", err)
	}
	if err := (LoseAember{Player: Controller, Amount: 2}).validate(); err != nil {
		t.Errorf("Amount alone should be valid, got %v", err)
	}
}

// TestMoveAemberFromPoolAndVault covers banking Æmber on a card and spending it
// back out again when its controller forges a key.
func TestMoveAemberFromPoolAndVault(t *testing.T) {
	e := MoveAemberFromPool{Amount: 1, Target: Target{Kind: TargetThisCreature}}
	if got := e.Text(); got != "move 1 Æmber from your pool to "+SelfName {
		t.Errorf("text = %q", got)
	}
	if err := (MoveAemberFromPool{Target: e.Target}).validate(); err == nil {
		t.Error("a move of no Æmber should be rejected")
	}
	if err := (MoveAemberFromPool{Amount: 1}).validate(); err == nil {
		t.Error("a move with no destination should be rejected")
	}
	if err := e.validate(); err != nil {
		t.Errorf("validate = %v, want nil", err)
	}

	vault := NewCard("Safe Place", Shadows, Artifact, Rare, WithSpendableAember())
	if text := RenderCardText(&vault); !strings.Contains(
		text, "You may spend Æmber on Safe Place when forging keys.") {
		t.Errorf("vault text = %q", text)
	}

	g := NewGame("A", "B", 1)
	id := g.AddArtifact(vault, 0)
	ctx := &EffectContext{Resolver: g, Source: id, Controller: 0}

	// An empty pool banks nothing.
	e.Resolve(ctx)
	if g.AmberOn(id) != 0 {
		t.Errorf("banked %d from an empty pool, want 0", g.AmberOn(id))
	}

	g.SetAember(0, 6)
	e.Resolve(ctx)
	if g.Aember(0) != 5 || g.AmberOn(id) != 1 {
		t.Errorf("pool = %d, banked = %d; want 5 and 1", g.Aember(0), g.AmberOn(id))
	}

	// 5 in the pool plus the 1 banked covers the key; the pool empties first.
	g.forgeKey(0)
	if g.Keys(0) != 1 {
		t.Errorf("keys = %d, want 1", g.Keys(0))
	}
	if g.Aember(0) != 0 || g.AmberOn(id) != 0 {
		t.Errorf("pool = %d, banked = %d; want both 0", g.Aember(0), g.AmberOn(id))
	}

	// Banked Æmber alone is still short of a key, so nothing is spent.
	g.AddAmberOn(id, 3)
	g.forgeKey(0)
	if g.Keys(0) != 1 || g.AmberOn(id) != 3 {
		t.Errorf("keys = %d, banked = %d; want 1 and 3", g.Keys(0), g.AmberOn(id))
	}

	// A pool that covers the cost on its own leaves the bank alone.
	g.SetAember(0, KeyCost)
	g.forgeKey(0)
	if g.Keys(0) != 2 || g.AmberOn(id) != 3 {
		t.Errorf("keys = %d, banked = %d; want 2 and 3", g.Keys(0), g.AmberOn(id))
	}
}

// TestPlaceAemberOnThis covers placing Æmber from the common supply on the source
// card, which accrues on the card rather than moving from a pool.
func TestPlaceAemberOnThis(t *testing.T) {
	e := PlaceAemberOnThis{Amount: 1}
	if got := e.Text(); got != "place 1 Æmber from the common supply on "+SelfName {
		t.Errorf("text = %q", got)
	}
	if err := (PlaceAemberOnThis{}).validate(); err == nil {
		t.Error("placing no Æmber should be rejected")
	}
	if err := e.validate(); err != nil {
		t.Errorf("validate = %v, want nil", err)
	}

	vault := NewCard("Safe Place", Shadows, Artifact, Rare)
	g := NewGame("A", "B", 1)
	id := g.AddArtifact(vault, 0)
	ctx := &EffectContext{Resolver: g, Source: id, Controller: 0}

	e.Resolve(ctx)
	e.Resolve(ctx)
	if g.AmberOn(id) != 2 {
		t.Errorf("Æmber on card = %d, want 2", g.AmberOn(id))
	}
	if g.Aember(0) != 0 {
		t.Errorf("pool = %d, want 0 (Æmber comes from the supply)", g.Aember(0))
	}
}
