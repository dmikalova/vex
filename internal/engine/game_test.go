package engine

import "testing"

// ---- playing cards ----

func TestPlayCreatureEntersStunned(t *testing.T) {
	g := started(t)
	g.AddToHand(
		NewCard(
			"Stunner",
			Brobnar,
			Creature,
			Common,
			WithPower(3),
			WithEntersPlay(Stun{Target: Target{Kind: TargetThisCreature}}),
		),
		0,
	)
	id, err := g.PlayCreature(0, handIdx(g, 0, "Stunner"), false)
	if err != nil {
		t.Fatalf("PlayCreature: %v", err)
	}
	if !g.State.Cards[id].Stunned {
		t.Error("EntersStunned creature should enter play stunned")
	}

	// A creature without the flag enters unstunned.
	g.AddToHand(testCreature("Plain", 3), 0)
	id2, err := g.PlayCreature(0, handIdx(g, 0, "Plain"), false)
	if err != nil {
		t.Fatalf("PlayCreature: %v", err)
	}
	if g.State.Cards[id2].Stunned {
		t.Error("creature without EntersStunned should not enter stunned")
	}
}

func TestPlayCreatureFlankLeft(t *testing.T) {
	g := started(t)
	first := g.AddToHand(testCreature("first", 2), 0)
	second := g.AddToHand(testCreature("second", 2), 0)
	if _, err := g.PlayCreature(0, 0, false); err != nil { // right flank: [first]
		t.Fatal(err)
	}
	if _, err := g.PlayCreature(0, 0, true); err != nil { // left flank: [second, first]
		t.Fatal(err)
	}
	bl := g.Battleline(0)
	if bl[0] != second || bl[1] != first {
		t.Errorf("battleline order = %v, want [%d %d]", bl, second, first)
	}
}

func TestPlayTypeMismatchErrors(t *testing.T) {
	g := started(t)
	creatureID := g.AddToHand(testCreature("c", 2), 0)
	_ = creatureID
	// Playing a creature as an artifact/action/upgrade should fail on type.
	if _, err := g.PlayArtifact(0, 0); err != ErrWrongType {
		t.Errorf("PlayArtifact err = %v, want ErrWrongType", err)
	}
	if err := g.PlayAction(0, 0); err != ErrWrongType {
		t.Errorf("PlayAction err = %v, want ErrWrongType", err)
	}
	if _, err := g.PlayUpgrade(0, 0); err != ErrWrongType {
		t.Errorf("PlayUpgrade err = %v, want ErrWrongType", err)
	}
	// Now play it correctly.
	if _, err := g.PlayCreature(0, 0, false); err != nil {
		t.Fatalf("PlayCreature: %v", err)
	}
}

func TestTakeFromHandErrors(t *testing.T) {
	g := started(t)
	// Not active player.
	if _, err := g.PlayCreature(1, 0, false); err != ErrNotActivePlayer {
		t.Errorf("err = %v, want ErrNotActivePlayer", err)
	}
	// Bad hand index.
	if _, err := g.PlayCreature(0, 5, false); err != ErrCardNotInHand {
		t.Errorf("err = %v, want ErrCardNotInHand", err)
	}
	// Wrong house.
	dis := NewCard("Dissonant", Dis, Creature, Common, WithPower(2))
	g.AddToHand(dis, 0)
	if _, err := g.PlayCreature(0, 0, false); err != ErrWrongHouse {
		t.Errorf("err = %v, want ErrWrongHouse", err)
	}
	// Game over.
	g.State.Winner = 0
	if _, err := g.PlayCreature(0, 0, false); err != ErrGameOver {
		t.Errorf("err = %v, want ErrGameOver", err)
	}
}

// TestVersatileIgnoresActiveHouse verifies that Versatile relaxes only USING an
// in-play card, not playing it: a Versatile creature already in play may be used
// out of house (keyword printed or granted by an upgrade), but a Versatile card
// still cannot be played from hand while another house is active.
func TestVersatileIgnoresActiveHouse(t *testing.T) {
	g := started(t) // Brobnar is the active house.

	// A Versatile creature already in play, of another house, can still be used.
	versatile := NewCard("Free Agent", Dis, Creature, Common, WithPower(3), WithKeywords(Versatile))
	vid := g.AddToBattleline(versatile, 0)
	if err := g.Reap(0, vid); err != nil {
		t.Fatalf("reap Versatile in play out of house: %v", err)
	}

	// Versatile does not let it be PLAYED from hand out of house, though.
	hid := g.AddToHand(versatile, 0)
	if _, err := g.PlayCreature(0, handIdxByID(g, 0, hid), false); err != ErrWrongHouse {
		t.Errorf("play Versatile out of house err = %v, want ErrWrongHouse", err)
	}

	// A non-Versatile creature of the wrong house cannot be used either.
	dis := g.AddToBattleline(NewCard("Dissonant", Dis, Creature, Common, WithPower(2)), 0)
	if err := g.Reap(0, dis); err != ErrWrongHouse {
		t.Errorf("non-Versatile reap err = %v, want ErrWrongHouse", err)
	}

	// Versatile granted by an attached upgrade also relaxes using the host.
	mantle := g.AddToHand(NewCard("Mantle", Sanctum, Upgrade, Rare,
		WithStatic(StaticModifier{Keywords: []Keyword{Versatile}})), 0)
	g.AttachUpgrade(dis, mantle)
	if !g.usableInActiveHouse(dis) {
		t.Error("granted Versatile should let the host be used out of house")
	}

	// With no house chosen, any card may be played or used.
	g.State.ActiveHouse = HouseNone
	if !g.inActiveHouse(&versatile) {
		t.Error("no active house should allow playing any card")
	}
	if !g.usableInActiveHouse(dis) {
		t.Error("no active house should allow using any card")
	}
}

func TestPlayArtifactAndAction(t *testing.T) {
	g := started(t)
	g.AddToHand(exAutocannon(), 0)
	if _, err := g.PlayArtifact(0, 0); err != nil {
		t.Fatalf("PlayArtifact: %v", err)
	}
	if g.Aember(0) != 1 {
		t.Errorf("aember after artifact = %d, want 1", g.Aember(0))
	}

	// Action with a friendly creature to ready-and-fight.
	g.AddToBattleline(testCreature("friend", 4), 0)
	g.AddToBattleline(testCreature("foe", 2), 1)
	g.AddToHand(exBattleFury(), 0)
	if err := g.PlayAction(0, handIdx(g, 0, "Battle Fury")); err != nil {
		t.Fatalf("PlayAction: %v", err)
	}
	if len(g.Battleline(1)) != 0 {
		t.Error("foe should be destroyed by the ready-and-fight")
	}
	if len(g.Discard(0)) == 0 {
		t.Error("action card should be in discard")
	}
}

func TestPlayUpgradeErrorsAndSuccess(t *testing.T) {
	g := started(t)
	// No creature to attach to.
	g.AddToHand(exBruteStrength(), 0)
	if _, err := g.PlayUpgrade(0, 0); err != ErrNoTarget {
		t.Errorf("err = %v, want ErrNoTarget", err)
	}
	// Not active player.
	if _, err := g.PlayUpgrade(1, 0); err != ErrNotActivePlayer {
		t.Errorf("err = %v, want ErrNotActivePlayer", err)
	}
	// Bad index.
	if _, err := g.PlayUpgrade(0, 9); err != ErrCardNotInHand {
		t.Errorf("err = %v, want ErrCardNotInHand", err)
	}
	// Wrong house.
	disUp := NewCard("Dis Upgrade", Dis, Upgrade, Common, WithStatic(StaticModifier{PowerBonus: 1}))
	g.AddToHand(disUp, 0)
	if _, err := g.PlayUpgrade(0, handIdx(g, 0, "Dis Upgrade")); err != ErrWrongHouse {
		t.Errorf("err = %v, want ErrWrongHouse", err)
	}
	// Success: attach +5 power upgrade to a creature.
	host := g.AddToBattleline(testCreature("host", 3), 0)
	if _, err := g.PlayUpgrade(0, handIdx(g, 0, "Brute Strength")); err != nil {
		t.Fatalf("PlayUpgrade: %v", err)
	}
	if g.Power(host) != 8 {
		t.Errorf("host power = %d, want 8", g.Power(host))
	}
	if ups := g.Upgrades(host); len(ups) != 1 {
		t.Errorf("Upgrades(host) = %v, want 1 upgrade", ups)
	}
	// Game over branch.
	g.State.Winner = 0
	if _, err := g.PlayUpgrade(0, 0); err != ErrGameOver {
		t.Errorf("err = %v, want ErrGameOver", err)
	}
}

func TestPlayUpgradeFiresPlayAbility(t *testing.T) {
	g := started(t)
	host := g.AddToBattleline(testCreature("host", 5), 0)
	g.State.Cards[host].Damage = 3
	// The upgrade carries a non-Play ability (must be skipped on attach) and a
	// Play ability that acts on its host.
	up := NewCard("Test Boots", Brobnar, Upgrade, Rare,
		WithAbility(TriggerAfterReap, GainAember{Player: Controller, Amount: 1}),
		WithAbility(TriggerAfterPlay, Heal{Fully: true, Target: Target{Kind: TargetThisCreature}}))
	g.AddToHand(up, 0)
	before := g.Aember(0)
	if _, err := g.PlayUpgrade(0, 0); err != nil {
		t.Fatalf("PlayUpgrade: %v", err)
	}
	if g.Damage(host) != 0 {
		t.Errorf("host damage = %d, want 0 (upgrade Play fully heals its host)", g.Damage(host))
	}
	if g.Aember(0) != before {
		t.Errorf("non-Play ability must not fire on attach; aember %d -> %d", before, g.Aember(0))
	}
}

func TestDiscardFromHand(t *testing.T) {
	g := started(t)
	// Wrong house.
	dis := NewCard("Dis", Dis, Creature, Common, WithPower(2))
	g.AddToHand(dis, 0)
	if err := g.DiscardFromHand(0, 0); err != ErrWrongHouse {
		t.Errorf("err = %v, want ErrWrongHouse", err)
	}
	// Not the active player.
	if err := g.DiscardFromHand(1, 0); err != ErrNotActivePlayer {
		t.Errorf("err = %v, want ErrNotActivePlayer", err)
	}
	// Bad index.
	if err := g.DiscardFromHand(0, 9); err != ErrCardNotInHand {
		t.Errorf("err = %v, want ErrCardNotInHand", err)
	}
	// Success: discard a Brobnar (active-house) card.
	brob := g.AddToHand(testCreature("brob", 3), 0)
	if err := g.DiscardFromHand(0, handIdxByID(g, 0, brob)); err != nil {
		t.Fatalf("DiscardFromHand: %v", err)
	}
	if d := g.Discard(0); len(d) != 1 || d[0] != brob {
		t.Errorf("discard = %v, want [%d]", d, brob)
	}
	// Game over.
	g.State.Winner = 0
	if err := g.DiscardFromHand(0, 0); err != ErrGameOver {
		t.Errorf("err = %v, want ErrGameOver", err)
	}
}

// ---- turn flow ----

func TestForgeKeyWins(t *testing.T) {
	g := started(t)
	// Giant on board so the after-forge trigger fires; a tough enemy survives.
	g.AddToBattleline(exGiant(), 0)
	enemy := g.AddToBattleline(testCreature("enemy", 10), 1)

	// Only one key is forged at the start of a turn, even with Æmber to spare.
	g.State.Aember[0] = 3 * KeyCost
	g.StartTurn(0)
	if g.Keys(0) != 1 {
		t.Errorf("keys = %d, want 1 (one key forged per turn)", g.Keys(0))
	}
	if g.Aember(0) != 2*KeyCost {
		t.Errorf("aember = %d, want %d (paid for exactly one key)", g.Aember(0), 2*KeyCost)
	}
	// The giant's after-forge trigger fired once: 2 damage.
	if g.Damage(enemy) != 2 {
		t.Errorf("enemy damage = %d, want 2", g.Damage(enemy))
	}

	// Forging the third key wins the game.
	g.State.Keys[0] = 2
	g.State.Aember[0] = KeyCost
	g.forgeKey(0)
	if g.Keys(0) != KeysToWin || g.Winner() != 0 {
		t.Errorf("keys = %d, winner = %d, want %d and 0", g.Keys(0), g.Winner(), KeysToWin)
	}

	// StartTurn after a win is a no-op.
	turn := g.State.Turn
	g.StartTurn(1)
	if g.State.Turn != turn {
		t.Error("StartTurn should be a no-op once the game is over")
	}
}

func TestKeyCost(t *testing.T) {
	g := NewGame("A", "B", 1)
	if c := g.keyCost(0); c != KeyCost {
		t.Errorf("base keyCost = %d, want %d", c, KeyCost)
	}
	if c := g.CurrentKeyCost(0); c != KeyCost {
		t.Errorf("CurrentKeyCost = %d, want %d", c, KeyCost)
	}

	// An opponent-affecting change on player 1's card raises player 0's key cost,
	// not player 1's own.
	g.AddToBattleline(NewCard("Jammer", Mars, Creature, Common, WithPower(4),
		WithKeyCost(NewKeyCostChange(Opponent, 1))), 1)
	if c := g.keyCost(0); c != KeyCost+1 {
		t.Errorf("keyCost(0) = %d, want %d", c, KeyCost+1)
	}
	if c := g.keyCost(1); c != KeyCost {
		t.Errorf(
			"keyCost(1) = %d, want %d (opponent-only change spares its controller)",
			c,
			KeyCost,
		)
	}

	// An upgrade granting the change on player 1's creature stacks with it.
	host := g.AddToBattleline(NewCard("Host", Mars, Creature, Common, WithPower(3)), 1)
	attachUpgrade(g, host, NewCard("Pack", Mars, Upgrade, Uncommon,
		WithStatic(StaticModifier{KeyCostChange: NewKeyCostChange(Opponent, 2)})))
	if c := g.keyCost(0); c != KeyCost+3 {
		t.Errorf("keyCost(0) with jammer + pack = %d, want %d", c, KeyCost+3)
	}

	// EachPlayer affects both players; Controller affects only the card's owner.
	g2 := NewGame("A", "B", 1)
	g2.AddToBattleline(NewCard("Toll", Dis, Creature, Common, WithPower(3),
		WithKeyCost(NewKeyCostChange(EachPlayer, 1))), 0)
	if g2.keyCost(0) != KeyCost+1 || g2.keyCost(1) != KeyCost+1 {
		t.Errorf(
			"each-player change = %d/%d, want %d each",
			g2.keyCost(0),
			g2.keyCost(1),
			KeyCost+1,
		)
	}
	g2.AddToBattleline(NewCard("SelfTax", Dis, Creature, Common, WithPower(3),
		WithKeyCost(NewKeyCostChange(Controller, 2))), 1) // Controller = the card's owner
	if g2.keyCost(1) != KeyCost+1+2 || g2.keyCost(0) != KeyCost+1 {
		t.Errorf(
			"unset (controller) change = %d/%d, want %d/%d",
			g2.keyCost(0),
			g2.keyCost(1),
			KeyCost+1,
			KeyCost+3,
		)
	}

	// Forging respects the raised cost: one Æmber short forges nothing.
	g.State.Aember[0] = KeyCost + 2
	g.forgeKey(0)
	if g.Keys(0) != 0 {
		t.Errorf("keys = %d, want 0 (below the raised cost)", g.Keys(0))
	}
	// At exactly the cost it forges one key, paying the full raised amount.
	g.State.Aember[0] = KeyCost + 3
	g.forgeKey(0)
	if g.Keys(0) != 1 {
		t.Errorf("keys = %d, want 1", g.Keys(0))
	}
	if g.Aember(0) != 0 {
		t.Errorf("aember = %d, want 0 (paid the full raised cost)", g.Aember(0))
	}
}

func TestChooseHouseWrongPlayer(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	if err := g.ChooseHouse(1, Brobnar); err != ErrNotActivePlayer {
		t.Errorf("err = %v, want ErrNotActivePlayer", err)
	}
}

func TestChooseHouseForcedBindsWhenAvailable(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.SetPlayerHouses(0, []House{Mars, Logos, Untamed})
	g.State.ActivePlayer = 0
	g.State.ForcedHouse[0].Value = Mars
	if err := g.ChooseHouse(0, Logos); err != ErrMustChooseForcedHouse {
		t.Errorf("choosing a different house = %v, want ErrMustChooseForcedHouse", err)
	}
	if err := g.ChooseHouse(0, Mars); err != nil {
		t.Errorf("choosing the forced house = %v, want nil", err)
	}
}

// TestChooseHouseForcedIgnoredWhenUnavailable covers cannot overrides must: a
// forced house the player's deck does not contain is ignored, freeing the choice.
func TestChooseHouseForcedIgnoredWhenUnavailable(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.SetPlayerHouses(0, []House{Brobnar, Logos, Untamed})
	g.State.ActivePlayer = 0
	g.State.ForcedHouse[0].Value = Mars // a house player 0 does not have
	if err := g.ChooseHouse(0, Brobnar); err != nil {
		t.Errorf("cannot-overrides-must: an available house = %v, want nil", err)
	}
}

func TestEndTurnReadyDrawAndArmorRefresh(t *testing.T) {
	g := started(t)
	// A creature with an armor upgrade, exhausted and damaged; plus an artifact.
	host := g.AddToBattleline(testCreature("host", 5, WithArmor(1)), 0)
	g.AddToHand(
		NewCard("Plating", Brobnar, Upgrade, Common, WithStatic(StaticModifier{ArmorBonus: 2})),
		0,
	)
	if _, err := g.PlayUpgrade(0, handIdx(g, 0, "Plating")); err != nil {
		t.Fatal(err)
	}
	g.State.Cards[host].Exhausted = true
	g.State.Cards[host].ArmorRemaining = 0
	g.AddArtifact(exAutocannon(), 0)

	// Deck to draw from.
	for i := 0; i < 3; i++ {
		g.AddToDeck(testCreature("deckcard", 1), 0)
	}
	g.Shuffle(0)

	g.EndPlayPhase(0)

	if g.Exhausted(host) {
		t.Error("host should be readied")
	}
	if g.State.Cards[host].ArmorRemaining != 3 { // base 1 + upgrade 2
		t.Errorf("armor refreshed to %d, want 3", g.State.Cards[host].ArmorRemaining)
	}
	if len(g.Hand(0)) != 3 {
		t.Errorf("hand size = %d, want 3 drawn", len(g.Hand(0)))
	}
}

// ---- accessors / choosers / logging ----

func TestAccessorsAndChooser(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	if g.PlayerName(0) != "Alice" || g.PlayerName(1) != "Bob" {
		t.Error("player names wrong")
	}
	id := g.AddToBattleline(testCreature("c", 4), 0)
	if g.Def(id).Name != "c" {
		t.Errorf("Def name = %q", g.Def(id).Name)
	}
	if g.Name(id) != "c" {
		t.Errorf("Name = %q", g.Name(id))
	}

	// Artifacts accessor returns the artifact row.
	art := g.AddArtifact(exAutocannon(), 0)
	if arts := g.Artifacts(0); len(arts) != 1 || arts[0] != art {
		t.Errorf("Artifacts = %v, want [%d]", arts, art)
	}

	// Deck accessor returns the deck from top to bottom.
	top := g.AddToDeck(testCreature("top", 1), 0)
	if deck := g.Deck(0); len(deck) != 1 || deck[0] != top {
		t.Errorf("Deck = %v, want [%d]", deck, top)
	}

	// Nil chooser falls back to the default (FirstChooser).
	g.SetChooser(0, nil)
	got, ok := g.ChooseCreature(0, id, "p", []LocalID{id})
	if !ok || got != id {
		t.Errorf("default chooser = (%d,%v)", got, ok)
	}
	// Empty candidate list.
	if _, ok := g.ChooseCreature(0, id, "p", nil); ok {
		t.Error("empty candidates should return ok=false")
	}
}

func TestVerboseLogging(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.Verbose = true
	g.Record(TurnBegan{Player: 0, Turn: 1})
	if got := g.LogText(); len(got) != 1 || got[0] != "A begins turn 1" {
		t.Errorf("log = %v", got)
	}
	// With recording off the entry is dropped rather than stored, so a bot
	// exploring cloned positions pays nothing for narration it never reads.
	g.SetRecording(false)
	g.Record(TurnBegan{Player: 0, Turn: 2})
	if len(g.Log) != 1 {
		t.Errorf("log = %v, want recording off to drop the entry", g.LogText())
	}
}

// These tests exercise raw engine internals that are not reachable through the
// exported API, so they stay in package engine with vanilla (ability-free) cards.

func vanillaCreature(name string, power int) CardDefinition {
	return NewCard(name, Brobnar, Creature, Common, WithPower(power))
}

func TestShouldDestroyEdges(t *testing.T) {
	g := NewGame("A", "B", 1)
	// An artifact is not a creature, so damage rules never destroy it.
	art := g.AddArtifact(NewCard("relic", Brobnar, Artifact, Common), 0)
	if g.shouldDestroy(art) {
		t.Error("an artifact should not be destroyable via shouldDestroy")
	}
	// A creature no longer in play is not destroyable.
	c := g.AddToBattleline(vanillaCreature("c", 3), 0)
	g.State.Battleline[0].remove(c)
	if g.shouldDestroy(c) {
		t.Error("an out-of-play creature should not be destroyable")
	}
}

func TestDrawReshufflesEmptyDeck(t *testing.T) {
	g := NewGame("A", "B", 1)
	// Deck empty, discard holds 2 cards: drawing reshuffles the discard into deck.
	g.State.Discard[0].add(g.Register(vanillaCreature("a", 1), 0))
	g.State.Discard[0].add(g.Register(vanillaCreature("b", 1), 0))
	g.draw(0, 1)
	if g.State.Hand[0].Count != 1 || g.State.Deck[0].Count != 1 || g.State.Discard[0].Count != 0 {
		t.Errorf("reshuffle draw: hand=%d deck=%d discard=%d, want 1/1/0",
			g.State.Hand[0].Count, g.State.Deck[0].Count, g.State.Discard[0].Count)
	}
	// drawTo pulls the last card, then stops once deck and discard are both empty.
	g.drawTo(0, 5)
	if g.State.Hand[0].Count != 2 {
		t.Errorf("drawTo drained deck: hand = %d, want 2", g.State.Hand[0].Count)
	}
}

func TestOrderByChoice(t *testing.T) {
	g := NewGame("A", "B", 1)
	ids := []LocalID{10, 20, 30}
	eq := func(a, b []LocalID) bool {
		if len(a) != len(b) {
			return false
		}
		for i := range a {
			if a[i] != b[i] {
				return false
			}
		}
		return true
	}

	// Nothing to order for 0 or 1 ids.
	if got := g.OrderByChoice(0, "p", nil); len(got) != 0 {
		t.Errorf("nil order = %v", got)
	}
	if got := g.OrderByChoice(0, "p", []LocalID{7}); !eq(got, []LocalID{7}) {
		t.Errorf("single order = %v", got)
	}
	// FirstChooser (default) keeps the original order.
	if got := g.OrderByChoice(0, "p", ids); !eq(got, []LocalID{10, 20, 30}) {
		t.Errorf("first-chooser order = %v", got)
	}
	// Picking the last each time reverses the order.
	g.SetChooser(0, orderLastChooser{})
	if got := g.OrderByChoice(0, "p", ids); !eq(got, []LocalID{30, 20, 10}) {
		t.Errorf("last-chooser order = %v", got)
	}
	// A rejected pick falls back to the given order.
	g.SetChooser(0, orderRejectChooser{})
	if got := g.OrderByChoice(0, "p", ids); !eq(got, []LocalID{10, 20, 30}) {
		t.Errorf("reject order = %v", got)
	}
	// An Orderer chooser arranges the ids in a single call.
	g.SetChooser(0, orderAllChooser{})
	if got := g.OrderByChoice(0, "p", ids); !eq(got, []LocalID{30, 20, 10}) {
		t.Errorf("orderer order = %v", got)
	}
}

// TestKeyCostChangePerAndFlank covers the two refinements a key-cost change can
// carry: a count that scales it, and a flank condition that suspends it.
func TestKeyCostChangePerAndFlank(t *testing.T) {
	obelisk := NewCard("Iron Obelisk", Brobnar, Artifact, Rare,
		WithKeyCost(NewKeyCostChange(Opponent, 1).Per(InPlay{
			Player:  Controller,
			Type:    Creature,
			House:   Brobnar,
			Damaged: true,
		})))
	want := "Your opponent's keys cost +1 Æmber for each friendly damaged Brobnar creature."
	if got := keyCostText(obelisk.KeyCostChange); got != want {
		t.Errorf("per text = %q, want %q", got, want)
	}

	g := NewGame("A", "B", 1)
	g.AddArtifact(obelisk, 0)
	hurt := g.AddToBattleline(NewCard("Brute", Brobnar, Creature, Common, WithPower(9)), 0)
	if got := g.CurrentKeyCost(1); got != KeyCost {
		t.Errorf("key cost with no damaged creature = %d, want %d", got, KeyCost)
	}
	g.State.Cards[hurt].Damage = 1
	if got := g.CurrentKeyCost(1); got != KeyCost+1 {
		t.Errorf("key cost with one damaged creature = %d, want %d", got, KeyCost+1)
	}

	mechanic := NewCard("Titan Mechanic", Logos, Creature, Common, WithPower(6),
		WithKeyCost(NewKeyCostChange(EachPlayer, -1).WhileOnFlank()))
	wantFlank := "While " + SelfName + " is on a flank, each player's keys cost -1 Æmber."
	if got := keyCostText(mechanic.KeyCostChange); got != wantFlank {
		t.Errorf("flank text = %q, want %q", got, wantFlank)
	}

	f := NewGame("A", "B", 1)
	f.AddToBattleline(mechanic, 0)
	if got := f.CurrentKeyCost(0); got != KeyCost-1 {
		t.Errorf("key cost on a flank = %d, want %d", got, KeyCost-1)
	}

	mid := NewGame("A", "B", 1)
	mid.AddToBattleline(testCreature("left", 3), 0)
	mid.AddToBattleline(mechanic, 0)
	mid.AddToBattleline(testCreature("right", 3), 0)
	if got := mid.CurrentKeyCost(0); got != KeyCost {
		t.Errorf("key cost off the flanks = %d, want %d", got, KeyCost)
	}
}
