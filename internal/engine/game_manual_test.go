package engine

import "testing"

func TestManualZoneString(t *testing.T) {
	cases := map[ManualZone]string{
		ManualHand: "hand", ManualDeckTop: "top of deck", ManualDeckBottom: "bottom of deck",
		ManualDiscard: "discard", ManualArchives: "archives", ManualPurge: "purge",
		ManualZone(99): "unknown",
	}
	for z, want := range cases {
		if got := z.String(); got != want {
			t.Errorf("ManualZone(%d).String() = %q, want %q", z, got, want)
		}
	}
}

func TestManualModeLiftsHouse(t *testing.T) {
	g := started(t) // active house Brobnar
	off := g.AddToHand(NewCard("off", Sanctum, Creature, Common, WithPower(1)), 0)
	if err := g.CanPlay(0, off); err != ErrWrongHouse {
		t.Fatalf("off-house without manual = %v, want ErrWrongHouse", err)
	}
	if g.Manual() {
		t.Fatal("manual should be off by default")
	}
	g.SetManual(true)
	if !g.Manual() {
		t.Fatal("SetManual(true) should turn manual on")
	}
	if err := g.CanPlay(0, off); err != nil {
		t.Errorf("off-house in manual = %v, want nil (house lifted)", err)
	}
}

func TestManualAddAmber(t *testing.T) {
	g := started(t)
	g.ManualAddAmber(0, 3)
	if g.Aember(0) != 3 {
		t.Errorf("amber = %d, want 3", g.Aember(0))
	}
	g.ManualAddAmber(0, -5) // clamps at zero
	if g.Aember(0) != 0 {
		t.Errorf("amber = %d, want 0 (clamped)", g.Aember(0))
	}
}

func TestManualAddChains(t *testing.T) {
	g := started(t)
	g.ManualAddChains(0, 4)
	if g.State.Chains[0] != 4 {
		t.Errorf("chains = %d, want 4", g.State.Chains[0])
	}
	g.ManualAddChains(0, -9) // clamps at zero
	if g.State.Chains[0] != 0 {
		t.Errorf("chains = %d, want 0 (clamped)", g.State.Chains[0])
	}
}

func TestManualSetActiveHouse(t *testing.T) {
	g := started(t)
	g.ManualSetActiveHouse(Dis)
	if g.State.ActiveHouse != Dis {
		t.Errorf("active house = %v, want Dis", g.State.ActiveHouse)
	}
}

func TestManualForgeAndUnforgeKey(t *testing.T) {
	g := started(t)
	g.ManualUnforgeKey(0) // no-op with no keys forged
	if g.Keys(0) != 0 {
		t.Fatalf("keys = %d, want 0", g.Keys(0))
	}
	for i := 1; i <= KeysToWin; i++ {
		g.ManualForgeKey(0)
		if g.Keys(0) != i {
			t.Fatalf("after forge %d: keys = %d", i, g.Keys(0))
		}
	}
	g.ManualForgeKey(0)                   // no-op: no colours remain
	g.ManualForgeKeyColor(0, KeyColorRed) // no-op: already at KeysToWin
	if g.Keys(0) != KeysToWin {
		t.Errorf("keys = %d, want %d (capped)", g.Keys(0), KeysToWin)
	}
	if got := len(g.KeyColors(0)); got != KeysToWin {
		t.Errorf("colors = %d, want %d", got, KeysToWin)
	}
	g.ManualUnforgeKey(0)
	if g.Keys(0) != KeysToWin-1 || len(g.KeyColors(0)) != KeysToWin-1 {
		t.Errorf("after unforge: keys = %d, colors = %d", g.Keys(0), len(g.KeyColors(0)))
	}
	// A specific colour can be chosen for the freed slot.
	g.ManualForgeKeyColor(0, KeyColorYellow)
	if got := g.KeyColors(0); got[len(got)-1] != KeyColorYellow {
		t.Errorf("last key = %v, want Yellow", got[len(got)-1])
	}
}

func TestManualMoveToEachZone(t *testing.T) {
	g := started(t)
	moves := []struct {
		dest  ManualZone
		count func(*Game) int
	}{
		{ManualHand, func(g *Game) int { return len(g.Hand(0)) }},
		{ManualDeckTop, func(g *Game) int { return len(g.Deck(0)) }},
		{ManualDeckBottom, func(g *Game) int { return len(g.Deck(0)) }},
		{ManualDiscard, func(g *Game) int { return len(g.Discard(0)) }},
		{ManualArchives, func(g *Game) int { return len(g.Archives(0)) }},
		{ManualPurge, func(g *Game) int { return len(g.Purge(0)) }},
	}
	for _, m := range moves {
		id := g.AddToHand(testCreature("c", 3), 0)
		g.ManualMove(id, m.dest)
		if m.count(g) == 0 {
			t.Errorf("ManualMove to %s did not place the card", m.dest)
		}
	}
	// Deck top places on top.
	top := g.AddToHand(testCreature("top", 3), 0)
	g.ManualMove(top, ManualDeckTop)
	if g.State.Deck[0].IDs[0] != top {
		t.Error("ManualMove ManualDeckTop should place on top of the deck")
	}
}

func TestManualPlaceInPlay(t *testing.T) {
	g := started(t)
	// Two creatures already in the line: a new one placed before index 1 lands
	// between them, and an armor value carries onto the board.
	left := g.AddToBattleline(testCreature("left", 3), 0)
	right := g.AddToBattleline(testCreature("right", 3), 0)
	mid := g.AddToHand(testCreature("mid", 3, WithArmor(2)), 0)

	g.ManualPlaceInPlay(mid, 1)

	line := g.Battleline(0)
	if len(line) != 3 || line[0] != left || line[1] != mid || line[2] != right {
		t.Fatalf("battleline = %v, want [left mid right]", line)
	}
	if a := g.State.Cards[mid].ArmorRemaining; a != 2 {
		t.Errorf("armor = %d, want 2", a)
	}
	if len(g.Hand(0)) != 0 {
		t.Errorf("card should have left the hand")
	}

	// An out-of-range index clamps to the right flank.
	far := g.AddToHand(testCreature("far", 3), 0)
	g.ManualPlaceInPlay(far, 99)
	if line = g.Battleline(0); line[len(line)-1] != far {
		t.Errorf("clamped placement should land on the right flank, got %v", line)
	}

	// A negative index clamps to the left flank.
	neg := g.AddToHand(testCreature("neg", 3), 0)
	g.ManualPlaceInPlay(neg, -5)
	if line = g.Battleline(0); line[0] != neg {
		t.Errorf("negative index should land on the left flank, got %v", line)
	}

	// A non-creature enters the artifact row instead of the battleline.
	art := g.AddToHand(testArtifact("relic"), 0)
	g.ManualPlaceInPlay(art, 0)
	if arts := g.Artifacts(0); len(arts) != 1 || arts[0] != art {
		t.Errorf("artifact should enter the artifact row, got %v", arts)
	}
}

func TestManualMoveFromPlayResetsAndShedsUpgrades(t *testing.T) {
	g := started(t)
	host := g.AddToBattleline(testCreature("host", 3), 0)
	g.State.Cards[host].Exhausted = true
	g.State.Cards[host].Damage = 2
	up := g.Register(exBruteStrength(), 0)
	g.AttachUpgrade(host, up)

	g.ManualMove(host, ManualHand)

	if g.inPlay(host) {
		t.Error("card should have left the battleline")
	}
	if len(g.Hand(0)) != 1 || g.Hand(0)[0] != host {
		t.Error("card should be in hand")
	}
	if g.State.Cards[host] != (CardCore{}) {
		t.Errorf("in-play state should be reset, got %+v", g.State.Cards[host])
	}
	if d := g.Discard(0); len(d) != 1 || d[0] != up {
		t.Errorf("upgrade should be discarded, discard = %v", d)
	}
}

// ManualAttachUnder takes a card from hand or from play and threads it under a
// host; a card from play sheds its state and upgrades on the way under.
func TestManualAttachUnder(t *testing.T) {
	g := started(t)
	host := g.AddToBattleline(testCreature("host", 3), 0)

	// From hand, face down.
	fromHand := g.AddToHand(testCreature("buried", 2), 0)
	g.ManualAttachUnder(host, fromHand, true)
	if u, ok := g.firstUnder(host); !ok || u != fromHand {
		t.Fatalf("firstUnder = %d,%v, want %d", u, ok, fromHand)
	}
	if !g.State.Cards[fromHand].UnderFaceDown {
		t.Error("card from hand should be placed face down")
	}
	if len(g.Hand(0)) != 0 {
		t.Error("card should have left the hand")
	}

	// From play, face up (graft): it sheds its damage and its upgrade.
	inPlay := g.AddToBattleline(testCreature("grafted", 4), 0)
	g.State.Cards[inPlay].Damage = 2
	up := g.Register(exBruteStrength(), 0)
	g.AttachUpgrade(inPlay, up)
	g.ManualAttachUnder(host, inPlay, false)
	if g.inPlay(inPlay) {
		t.Error("card should have left the battleline")
	}
	if g.State.Cards[inPlay].UnderFaceDown {
		t.Error("a graft should be face up")
	}
	if under := g.underOf(host); len(under) != 2 || under[1] != inPlay {
		t.Errorf("underOf = %v, want [%d %d]", under, fromHand, inPlay)
	}
	if d := g.Discard(0); len(d) != 1 || d[0] != up {
		t.Errorf("the upgrade should be discarded, discard = %v", d)
	}
}

// ManualDetachToHand sends a selected upgrade or under-card to hand; a card that
// is neither is left where it is.
func TestManualDetachToHand(t *testing.T) {
	g := started(t)
	host := g.AddToBattleline(testCreature("host", 3), 0)

	// An upgrade detaches to hand.
	up := g.Register(exBruteStrength(), 0)
	g.AttachUpgrade(host, up)
	g.ManualDetachToHand(up)
	if _, ok := g.hostOf(up); ok {
		t.Error("the upgrade should be detached from its host")
	}
	if len(g.Hand(0)) != 1 || g.Hand(0)[0] != up {
		t.Errorf("hand = %v, want [%d]", g.Hand(0), up)
	}

	// An under-card detaches to hand and sheds its state.
	buried := g.Register(testCreature("buried", 2), 0)
	g.AttachUnder(host, buried, true)
	g.ManualDetachToHand(buried)
	if _, ok := g.underHostOf(buried); ok {
		t.Error("the under-card should be detached from its host")
	}
	if !g.State.Hand[0].contains(buried) {
		t.Error("the under-card should be in hand")
	}
	if g.State.Cards[buried] != (CardCore{}) {
		t.Errorf("the under-card's state should be reset, got %+v", g.State.Cards[buried])
	}

	// A card that is neither an upgrade nor under a host is left in play.
	loose := g.AddToBattleline(testCreature("loose", 3), 0)
	g.ManualDetachToHand(loose)
	if !g.inPlay(loose) {
		t.Error("a card that is neither attached nor placed under a host should be left alone")
	}
}

func TestManualSetExhausted(t *testing.T) {
	g := started(t)
	id := g.AddToBattleline(testCreature("c", 3), 0)
	g.State.Cards[id].Exhausted = true
	g.ManualSetExhausted(id, false)
	if g.Exhausted(id) {
		t.Error("ManualSetExhausted(false) should ready the card")
	}
	g.ManualSetExhausted(id, true)
	if !g.Exhausted(id) {
		t.Error("ManualSetExhausted(true) should exhaust the card")
	}
}

func TestManualAddCard(t *testing.T) {
	g := started(t)
	before := len(g.Hand(0))
	id, ok := g.ManualAddCard(NewCard("Import", Logos, Tactic, Common), 0)
	if !ok {
		t.Fatal("ManualAddCard should succeed in a fresh match")
	}
	if len(g.Hand(0)) != before+1 {
		t.Errorf("hand size = %d, want %d", len(g.Hand(0)), before+1)
	}
	if g.Def(id).Name != "Import" {
		t.Errorf("added card name = %q, want Import", g.Def(id).Name)
	}

	for g.cat.hasRoom() {
		g.Register(NewCard("Filler", Logos, Tactic, Common), 0)
	}
	full := len(g.Hand(0))
	if _, ok := g.ManualAddCard(NewCard("Overflow", Logos, Tactic, Common), 0); ok {
		t.Error("ManualAddCard should refuse once the match is full")
	}
	if len(g.Hand(0)) != full {
		t.Errorf("refused add changed the hand: %d, want %d", len(g.Hand(0)), full)
	}
}
