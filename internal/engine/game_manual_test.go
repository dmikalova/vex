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
	id := g.ManualAddCard(NewCard("Import", Logos, Tactic, Common), 0)
	if len(g.Hand(0)) != before+1 {
		t.Errorf("hand size = %d, want %d", len(g.Hand(0)), before+1)
	}
	if g.Def(id).Name != "Import" {
		t.Errorf("added card name = %q, want Import", g.Def(id).Name)
	}
}
