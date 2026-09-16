package deckgen

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/engine"
)

// isUpgradeOrRobot mirrors Chief Engineer Walls's pull predicate for tests.
func isUpgradeOrRobot(d engine.CardDefinition) bool {
	if d.Type == engine.Upgrade {
		return true
	}
	for _, tr := range d.Traits {
		if tr == engine.Robot {
			return true
		}
	}
	return false
}

func upgradeCard(name string, h engine.House) Card {
	return Card{Def: engine.NewCard(name, h, engine.Upgrade, engine.Common)}
}

func robotCard(name string, h engine.House) Card {
	return Card{
		Def: engine.NewCard(
			name, h, engine.Creature, engine.Common,
			engine.WithPower(3), engine.WithTraits(engine.Robot),
		),
	}
}

func leadCard(name string, h engine.House, fc FilteredCluster) Card {
	c := mkCard(name, h, engine.Common)
	c.Profile.Leads = &fc
	return c
}

// richFilteredSet is a set whose lead card pulls Upgrades and Robots to a floor of
// two, with a matching pool spread across several Houses plus vanilla filler.
func richFilteredSet() Set {
	fc := FilteredCluster{Name: "F", Floor: 2, Match: isUpgradeOrRobot}
	return NewSet("S", []Card{
		leadCard("Lead", engine.Brobnar, fc),
		mkCard("VB", engine.Brobnar, engine.Common),
		mkCard("VD1", engine.Dis, engine.Common),
		mkCard("VD2", engine.Dis, engine.Common),
		mkCard("VL1", engine.Logos, engine.Common),
		mkCard("VL2", engine.Logos, engine.Common),
		upgradeCard("UpD", engine.Dis),
		upgradeCard("UpL", engine.Logos),
		upgradeCard("UpM", engine.Mars),
		robotCard("BotD", engine.Dis),
	}, Tuning{RarityWeights: map[engine.Rarity]float64{engine.Common: 1}})
}

// vanillaDeck fills each pod with a distinct non-matching filler creature so a
// filtered pull has room to overwrite slots.
func vanillaDeck(houses [PodCount]engine.House) Deck {
	var d Deck
	for i, h := range houses {
		d.Pods[i].House = h
		for s := range d.Pods[i].Slots {
			def := engine.NewCard("V", h, engine.Creature, engine.Common, engine.WithPower(3))
			d.Pods[i].Slots[s] = Slot{Rarity: engine.Common, Card: def}
		}
	}
	return d
}

func TestBuildFilteredClustersStampsLead(t *testing.T) {
	s := richFilteredSet()
	fc, ok := s.filtered["F"]
	if !ok {
		t.Fatal("filtered cluster F missing")
	}
	if fc.Lead != "Lead" {
		t.Errorf("lead = %q, want Lead", fc.Lead)
	}
	if pool := s.matchingPool(fc); len(pool) != 4 {
		t.Errorf("matching pool = %d, want 4", len(pool))
	}
}

// filteredTarget with a zero Mean pulls exactly to the Floor; a Mean above the
// Floor rolls a Poisson-tailed total that is always at least the Floor and
// averages about the Mean.
func TestFilteredTarget(t *testing.T) {
	g := gen(richFilteredSet())
	flat := FilteredCluster{Name: "F", Floor: 3, Match: isUpgradeOrRobot}
	for i := 0; i < 200; i++ {
		if got := g.filteredTarget(flat); got != 3 {
			t.Fatalf("zero-Mean target = %d, want exactly 3", got)
		}
	}
	spread := FilteredCluster{Name: "F", Floor: 4, Mean: 6, Match: isUpgradeOrRobot}
	sum := 0
	const n = 2000
	for i := 0; i < n; i++ {
		got := g.filteredTarget(spread)
		if got < 4 {
			t.Fatalf("target %d below floor 4", got)
		}
		sum += got
	}
	if avg := float64(sum) / n; avg < 5.5 || avg > 6.5 {
		t.Errorf("mean target = %.2f, want about 6", avg)
	}
}

func TestValidateFilteredClustersPanics(t *testing.T) {
	cases := []struct {
		name string
		fc   FilteredCluster
	}{
		{"nil predicate", FilteredCluster{Name: "F", Floor: 2, Match: nil}},
		{"floor below one", FilteredCluster{Name: "F", Floor: 0, Match: isUpgradeOrRobot}},
		{"pool too small", FilteredCluster{Name: "F", Floor: 5, Match: isUpgradeOrRobot}},
		{
			"mean below floor",
			FilteredCluster{Name: "F", Floor: 2, Mean: 1, Match: isUpgradeOrRobot},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("expected panic")
				}
			}()
			NewSet("S", []Card{
				leadCard("Lead", engine.Brobnar, tc.fc),
				mkCard("VB", engine.Brobnar, engine.Common),
				upgradeCard("UpD", engine.Dis),
			}, Tuning{RarityWeights: map[engine.Rarity]float64{engine.Common: 1}})
		})
	}
}

// A mean below the floor is only reached once the pool is large enough to clear the
// pool-size check, so it needs a set with enough matching cards to satisfy the floor.
func TestValidateFilteredClustersMeanBelowFloor(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for a mean below the floor")
		}
	}()
	fc := FilteredCluster{Name: "F", Floor: 2, Mean: 1, Match: isUpgradeOrRobot}
	NewSet("S", []Card{
		leadCard("Lead", engine.Brobnar, fc),
		mkCard("VB", engine.Brobnar, engine.Common),
		upgradeCard("UpD", engine.Dis),
		upgradeCard("UpL", engine.Logos),
	}, Tuning{RarityWeights: map[engine.Rarity]float64{engine.Common: 1}})
}

func TestExpandFilteredClustersLeadAbsent(t *testing.T) {
	g := gen(richFilteredSet())
	deck := vanillaDeck([PodCount]engine.House{engine.Dis, engine.Logos, engine.Mars})
	g.expandFilteredClusters(&deck)
	if n := deckCountMatching(&deck, isUpgradeOrRobot); n != 0 {
		t.Errorf("matches without lead = %d, want 0 (no pull should fire)", n)
	}
}

func TestExpandFilteredClustersAlreadySatisfied(t *testing.T) {
	g := gen(richFilteredSet())
	deck := vanillaDeck([PodCount]engine.House{engine.Brobnar, engine.Dis, engine.Logos})
	deck.Pods[0].Slots[0] = Slot{Card: g.set.byName["Lead"].Def}
	deck.Pods[1].Slots[0] = Slot{Card: g.set.byName["UpD"].Def}
	deck.Pods[1].Slots[1] = Slot{Card: g.set.byName["BotD"].Def}
	g.expandFilteredClusters(&deck)
	if n := deckCountMatching(&deck, isUpgradeOrRobot); n != 2 {
		t.Errorf("matches = %d, want 2 (already met, no top-up)", n)
	}
}

func TestExpandFilteredClustersTopsUp(t *testing.T) {
	g := gen(richFilteredSet())
	deck := vanillaDeck([PodCount]engine.House{engine.Brobnar, engine.Dis, engine.Logos})
	deck.Pods[0].Slots[0] = Slot{Card: g.set.byName["Lead"].Def}
	g.expandFilteredClusters(&deck)
	if n := deckCountMatching(&deck, isUpgradeOrRobot); n < 2 {
		t.Errorf("matches after top-up = %d, want ≥ 2", n)
	}
	if !deckHasCard(&deck, "Lead") {
		t.Error("lead should still be present")
	}
}

func TestPlaceFilteredMatchNative(t *testing.T) {
	g := gen(richFilteredSet())
	deck := vanillaDeck([PodCount]engine.House{engine.Dis, engine.Logos, engine.Brobnar})
	fc := g.set.filtered["F"]
	if !g.placeFilteredMatch(&deck, fc, g.set.byName["UpD"]) {
		t.Fatal("expected native placement to succeed")
	}
	found := false
	for _, s := range deck.Pods[0].Slots {
		if s.Card.Name == "UpD" {
			found = true
			if s.Maverick {
				t.Error("native placement should not be a maverick")
			}
		}
	}
	if !found {
		t.Error("UpD not placed in its native Dis pod")
	}
}

func TestPlaceFilteredMatchMaverick(t *testing.T) {
	g := gen(richFilteredSet())
	deck := vanillaDeck([PodCount]engine.House{engine.Dis, engine.Logos, engine.Brobnar})
	fc := g.set.filtered["F"]
	// UpM's Mars House is not in the deck, so it must land as a maverick.
	if !g.placeFilteredMatch(&deck, fc, g.set.byName["UpM"]) {
		t.Fatal("expected maverick placement to succeed")
	}
	for i := 0; i < PodCount; i++ {
		for _, s := range deck.Pods[i].Slots {
			if s.Card.Name == "UpM" {
				if !s.Maverick || s.Card.House != deck.Pods[i].House {
					t.Errorf("UpM not rehoused as maverick: %+v", s)
				}
				return
			}
		}
	}
	t.Error("UpM not placed anywhere")
}

func TestPlaceFilteredMatchNoSlot(t *testing.T) {
	g := gen(richFilteredSet())
	var deck Deck
	// Pod 0 has no House: the guard must skip it. The other pods are packed with
	// cards that already match, so no overwritable slot exists.
	deck.Pods[0].House = engine.HouseNone
	up := g.set.byName["UpD"].Def
	for i := 1; i < PodCount; i++ {
		deck.Pods[i].House = engine.Dis
		for s := range deck.Pods[i].Slots {
			deck.Pods[i].Slots[s] = Slot{Card: up}
		}
	}
	if g.placeFilteredMatch(&deck, g.set.filtered["F"], g.set.byName["BotD"]) {
		t.Error("expected no placement when every slot is matching or House-less")
	}
}

func TestPlaceFilteredMatchOneCopyPerDeck(t *testing.T) {
	fc := FilteredCluster{Name: "F", Floor: 1, Match: isUpgradeOrRobot}
	once := upgradeCard("UpOnce", engine.Dis)
	once.Profile.OneCopyPerDeck = true
	s := NewSet("S", []Card{
		leadCard("Lead", engine.Brobnar, fc),
		mkCard("VB", engine.Brobnar, engine.Common),
		mkCard("VD", engine.Dis, engine.Common),
		once,
	}, Tuning{RarityWeights: map[engine.Rarity]float64{engine.Common: 1}})
	g := gen(s)
	deck := vanillaDeck([PodCount]engine.House{engine.Dis, engine.Logos, engine.Brobnar})
	if !g.placeFilteredMatch(&deck, g.set.filtered["F"], s.byName["UpOnce"]) {
		t.Fatal("expected placement to succeed")
	}
	if !g.placed["UpOnce"] {
		t.Error("placing a one-copy-per-deck card should mark it placed")
	}
}

func TestFilteredCandidatesSkips(t *testing.T) {
	fc := FilteredCluster{Name: "F", Floor: 1, Match: isUpgradeOrRobot}
	once := upgradeCard("UpOnce", engine.Dis)
	once.Profile.OneCopyPerDeck = true
	s := NewSet("S", []Card{
		leadCard("Lead", engine.Brobnar, fc),
		mkCard("VB", engine.Brobnar, engine.Common),
		mkCard("VD", engine.Dis, engine.Common),
		upgradeCard("UpD", engine.Dis),
		once,
	}, Tuning{RarityWeights: map[engine.Rarity]float64{engine.Common: 1}})
	g := gen(s)
	g.placed["UpOnce"] = true // already placed elsewhere
	deck := vanillaDeck([PodCount]engine.House{engine.Brobnar, engine.Dis, engine.Logos})
	deck.Pods[1].Slots[0] = Slot{Card: s.byName["UpD"].Def} // already in deck
	cands := g.filteredCandidates(&deck, g.set.filtered["F"])
	for _, c := range cands {
		if c.Def.Name == "UpD" || c.Def.Name == "UpOnce" {
			t.Errorf("candidate %q should be excluded", c.Def.Name)
		}
	}
}

func TestProtectedNamesUnion(t *testing.T) {
	whole := ClusterMembership{Name: "Horsemen", Strategy: WholePool, Trigger: ByAnyMember}
	shard := ClusterMembership{Name: "Shard", Strategy: OnePerHouse, Trigger: ByAnyMember}
	fc := FilteredCluster{Name: "F", Floor: 1, Match: isUpgradeOrRobot}
	s := NewSet("S", []Card{
		leadCard("Lead", engine.Brobnar, fc),
		clusterMember("H1", engine.Brobnar, whole),
		clusterMember("H2", engine.Dis, whole),
		clusterMember("Sh1", engine.Brobnar, shard),
		clusterMember("Sh2", engine.Dis, shard),
		upgradeCard("UpD", engine.Dis),
		mkCard("VB", engine.Brobnar, engine.Common),
	}, Tuning{RarityWeights: map[engine.Rarity]float64{engine.Common: 1}})
	g := gen(s)
	protected := g.protectedNames()
	for _, name := range []string{"H1", "H2", "Sh1", "Sh2"} {
		if !protected[name] {
			t.Errorf("%q should be protected", name)
		}
	}
	if protected["UpD"] {
		t.Error("UpD is not a cluster member and should not be protected")
	}
}

func TestFilteredGenerateGuaranteesFloor(t *testing.T) {
	s := richFilteredSet()
	for seed := int64(1); seed <= 40; seed++ {
		deck := Generate(s, seed)
		if !deckHasCard(&deck, "Lead") {
			continue
		}
		if n := deckCountMatching(&deck, isUpgradeOrRobot); n < 2 {
			t.Fatalf("seed %d: deck with lead has %d matches, want ≥ 2", seed, n)
		}
	}
}
