package deckgen

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/dmikalova/vex/internal/engine"
)

// firePod fills a Brobnar/Dis/Shadows pod and plants one member in slot 0 so a
// pod-local cluster fires when expandPodClusters runs.
func firePod(g *generator, set Set, house engine.House, member string) HousePod {
	g.deckHouses[0] = house
	pod := g.fillPod(house)
	pod.Slots[0] = Slot{
		Rarity: engine.Rare,
		Card:   set.byName[member].Def,
	}
	return pod
}

func distinctMembers(pod HousePod, ci clusterIndex) int {
	seen := map[string]bool{}
	for i := range pod.Slots {
		s := pod.Slots[i]
		if inCluster(ci, s.Card.Name) {
			seen[s.Card.Name] = true
		}
	}
	return len(seen)
}

// A fired WholePool cluster places every member into the pod; without its lead it
// stays dormant.
func TestWholePoolPlacesAllMembers(t *testing.T) {
	whole := ClusterMembership{
		Name:     "H",
		Strategy: WholePool,
		Trigger:  ByLead,
	}
	lead := clusterMember("Lead", engine.Brobnar, whole)
	lead.Profile.Cluster.Lead = true
	set := NewSet("S", []Card{
		lead,
		clusterMember("R1", engine.Brobnar, whole),
		clusterMember("R2", engine.Brobnar, whole),
		mkCard("FB", engine.Brobnar, engine.Common),
	}, Tuning{RarityWeights: map[engine.Rarity]float64{engine.Common: 1}})

	g := &generator{
		set:    set,
		r:      rand.New(rand.NewSource(1)),
		placed: map[string]bool{},
	}
	pod := g.expandPodClusters(firePod(g, set, engine.Brobnar, "Lead"))
	for _, n := range []string{"Lead", "R1", "R2"} {
		if !podHas(pod, n) {
			t.Fatalf("WholePool did not place %s", n)
		}
	}

	// No lead planted: the ByLead cluster does not fire (only Common rolls).
	g = &generator{
		set:    set,
		r:      rand.New(rand.NewSource(2)),
		placed: map[string]bool{},
	}
	g.deckHouses[0] = engine.Brobnar
	pod = g.expandPodClusters(g.fillPod(engine.Brobnar))
	if distinctMembers(pod, set.clusters["H"]) != 0 {
		t.Fatal("WholePool fired without its lead")
	}
}

// A fired RandomCount cluster places between Min and all of its members, for both
// a range (Max > Min) and a fixed count (Max == Min).
func TestRandomCountPlacesMembers(t *testing.T) {
	for _, tc := range []struct {
		name     string
		min, max int
	}{
		{"range", 1, 3},
		{"fixed", 2, 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			sins := ClusterMembership{
				Name: "Sins", Strategy: RandomCount, Trigger: ByAnyMember,
				Min: tc.min, Max: tc.max,
			}
			set := NewSet("S", []Card{
				clusterMember("S1", engine.Dis, sins),
				clusterMember("S2", engine.Dis, sins),
				clusterMember("S3", engine.Dis, sins),
				mkCard("FD", engine.Dis, engine.Common),
			}, Tuning{RarityWeights: map[engine.Rarity]float64{engine.Common: 1}})

			for seed := range int64(20) {
				g := &generator{
					set:    set,
					r:      rand.New(rand.NewSource(seed)),
					placed: map[string]bool{},
				}
				pod := g.expandPodClusters(firePod(g, set, engine.Dis, "S1"))
				if n := distinctMembers(pod, set.clusters["Sins"]); n < tc.min || n > 3 {
					t.Fatalf("seed %d placed %d members, want [%d,3]", seed, n, tc.min)
				}
			}
		})
	}
}

// A fired RandomCount ByLead cluster places a subset of its non-lead members and
// never the lead again: the lead is already planted in the pod (its roll fired the
// cluster), so only its Connected partners ride in (Dark Harbinger pulls its
// Mutations, not itself).
func TestRandomCountByLeadExcludesLead(t *testing.T) {
	dh := ClusterMembership{
		Name: "Harbinger", Strategy: RandomCount, Trigger: ByLead, Min: 1, Max: 3,
	}
	lead := clusterMember("Lead", engine.Untamed, dh)
	lead.Profile.Cluster.Lead = true
	set := NewSet("S", []Card{
		lead,
		clusterMember("M1", engine.Untamed, dh),
		clusterMember("M2", engine.Untamed, dh),
		clusterMember("M3", engine.Untamed, dh),
		mkCard("FU", engine.Untamed, engine.Common),
	}, Tuning{RarityWeights: map[engine.Rarity]float64{engine.Common: 1}})

	ci := set.clusters["Harbinger"]
	for seed := range int64(30) {
		g := &generator{
			set:    set,
			r:      rand.New(rand.NewSource(seed)),
			placed: map[string]bool{},
		}
		pod := g.expandPodClusters(firePod(g, set, engine.Untamed, "Lead"))
		if got := countMember(pod, "Lead"); got != 1 {
			t.Fatalf("seed %d placed the lead %d times, want only the planted one", seed, got)
		}
		// distinctMembers counts the planted lead too, so its partners are one fewer.
		if partners := distinctMembers(pod, ci) - 1; partners < 1 || partners > 3 {
			t.Fatalf("seed %d placed %d partners, want [1,3]", seed, partners)
		}
	}
}

// A fired SelfPull cluster places at least Min copies of its single member and
// never more than a full pod; without a member present it stays dormant.
func TestSelfPullPlacesCopies(t *testing.T) {
	pull := ClusterMembership{
		Name:     "Rat",
		Strategy: SelfPull,
		Trigger:  ByAnyMember,
		Min:      3,
		Mean:     5,
	}
	set := NewSet("S", []Card{
		clusterMember("Rat", engine.Shadows, pull),
		mkCard("FS", engine.Shadows, engine.Common),
	}, Tuning{RarityWeights: map[engine.Rarity]float64{engine.Common: 1}})

	for seed := range int64(50) {
		g := &generator{
			set:    set,
			r:      rand.New(rand.NewSource(seed)),
			placed: map[string]bool{},
		}
		pod := g.expandPodClusters(firePod(g, set, engine.Shadows, "Rat"))
		if got := countMember(pod, "Rat"); got < 3 || got > PodSize {
			t.Fatalf("seed %d placed %d Rats, want [3,%d]", seed, got, PodSize)
		}
	}

	// No Rat present: the cluster does not fire.
	g := &generator{
		set:    set,
		r:      rand.New(rand.NewSource(9)),
		placed: map[string]bool{},
	}
	g.deckHouses[0] = engine.Shadows
	pod := g.expandPodClusters(g.fillPod(engine.Shadows))
	if countMember(pod, "Rat") != 0 {
		t.Fatal("SelfPull fired without a member")
	}
}

// selfPullCount is at least Min, returns exactly Min when Mean equals Min, and is
// capped at PodSize on the tail.
func TestSelfPullCount(t *testing.T) {
	g := &generator{r: rand.New(rand.NewSource(7))}

	flat := clusterIndex{
		strategy: SelfPull,
		min:      4,
		mean:     4,
	}
	for range 20 {
		if n := g.selfPullCount(flat); n != 4 {
			t.Fatalf("Mean==Min returned %d, want 4", n)
		}
	}

	// A high Min at the pod ceiling drives the cap; every roll stays in range.
	capped := clusterIndex{
		strategy: SelfPull,
		min:      PodSize - 1,
		mean:     PodSize,
	}
	for range 500 {
		if n := g.selfPullCount(capped); n < PodSize-1 || n > PodSize {
			t.Fatalf("capped count %d out of range", n)
		}
	}
}

// A WholePool member native to another House is rehoused as a maverick in the
// pod it rides into.
func TestWholePoolRehousesMaverickMember(t *testing.T) {
	whole := ClusterMembership{
		Name:     "H",
		Strategy: WholePool,
		Trigger:  ByLead,
	}
	lead := clusterMember("Lead", engine.Brobnar, whole)
	lead.Profile.Cluster.Lead = true
	set := NewSet("S", []Card{
		lead,
		clusterMember("Rider", engine.Dis, whole),
		mkCard("FB", engine.Brobnar, engine.Common),
		mkCard("FD", engine.Dis, engine.Common),
	}, Tuning{RarityWeights: map[engine.Rarity]float64{engine.Common: 1}})

	g := &generator{
		set:    set,
		r:      rand.New(rand.NewSource(1)),
		placed: map[string]bool{},
	}
	pod := g.expandPodClusters(firePod(g, set, engine.Brobnar, "Lead"))
	found := false
	for i := range pod.Slots {
		s := pod.Slots[i]
		if s.Card.Name != "Rider" {
			continue
		}
		found = true
		if !s.Maverick {
			t.Fatal("cross-House member not marked maverick")
		}
		if s.Card.House != engine.Brobnar {
			t.Fatal("maverick member not rehoused to the pod's House")
		}
	}
	if !found {
		t.Fatal("Rider was not placed")
	}
}

// A SelfPull member marked OneCopyPerDeck is recorded as placed once it lands.
func TestSelfPullRecordsOneCopyPerDeck(t *testing.T) {
	pull := ClusterMembership{
		Name:     "Rat",
		Strategy: SelfPull,
		Trigger:  ByAnyMember,
		Min:      3,
		Mean:     3,
	}
	rat := clusterMember("Rat", engine.Shadows, pull)
	rat.Profile.OneCopyPerDeck = true
	set := NewSet("S", []Card{
		rat,
		mkCard("FS", engine.Shadows, engine.Common),
	}, Tuning{RarityWeights: map[engine.Rarity]float64{engine.Common: 1}})

	g := &generator{
		set:    set,
		r:      rand.New(rand.NewSource(1)),
		placed: map[string]bool{},
	}
	g.expandPodClusters(firePod(g, set, engine.Shadows, "Rat"))
	if !g.placed["Rat"] {
		t.Fatal("OneCopyPerDeck member not recorded as placed")
	}
}

// With more members than a pod has slots, WholePool fills the pod and stops when
// no cluster-free slot remains.
func TestWholePoolStopsWhenPodFull(t *testing.T) {
	whole := ClusterMembership{
		Name:     "H",
		Strategy: WholePool,
		Trigger:  ByLead,
	}
	lead := clusterMember("M0", engine.Brobnar, whole)
	lead.Profile.Cluster.Lead = true
	cards := []Card{lead}
	for i := 1; i <= PodSize+1; i++ {
		cards = append(cards, clusterMember(fmt.Sprintf("M%d", i), engine.Brobnar, whole))
	}
	cards = append(cards, mkCard("FB", engine.Brobnar, engine.Common))
	set := NewSet("S", cards, Tuning{RarityWeights: map[engine.Rarity]float64{engine.Common: 1}})

	g := &generator{
		set:    set,
		r:      rand.New(rand.NewSource(1)),
		placed: map[string]bool{},
	}
	pod := g.expandPodClusters(firePod(g, set, engine.Brobnar, "M0"))
	for i := range pod.Slots {
		s := pod.Slots[i]
		if !inCluster(set.clusters["H"], s.Card.Name) {
			t.Fatalf("slot holds non-member %q; the pod should be full of members", s.Card.Name)
		}
	}
}

// A fired PullExact cluster places one partner per lead instance in the pod.
func TestPullExactPlacesPerLead(t *testing.T) {
	pull := ClusterMembership{
		Name:     "TT",
		Strategy: PullExact,
		Trigger:  ByLead,
	}
	lead := clusterMember("Lead", engine.Logos, pull)
	lead.Profile.Cluster.Lead = true
	set := NewSet("S", []Card{
		lead,
		clusterMember("Partner", engine.Logos, pull),
		mkCard("FL", engine.Logos, engine.Common),
	}, Tuning{RarityWeights: map[engine.Rarity]float64{engine.Common: 1}})

	// One lead pulls one partner.
	g := &generator{
		set:    set,
		r:      rand.New(rand.NewSource(1)),
		placed: map[string]bool{},
	}
	pod := g.expandPodClusters(firePod(g, set, engine.Logos, "Lead"))
	if got := countMember(pod, "Partner"); got != 1 {
		t.Fatalf("one lead pulled %d partners, want 1", got)
	}

	// Two leads pull two partners.
	g = &generator{
		set:    set,
		r:      rand.New(rand.NewSource(1)),
		placed: map[string]bool{},
	}
	g.deckHouses[0] = engine.Logos
	pod = g.fillPod(engine.Logos)
	pod.Slots[0] = Slot{
		Rarity: engine.Rare,
		Card:   set.byName["Lead"].Def,
	}
	pod.Slots[1] = Slot{
		Rarity: engine.Rare,
		Card:   set.byName["Lead"].Def,
	}
	pod = g.expandPodClusters(pod)
	if got := countMember(pod, "Partner"); got != 2 {
		t.Fatalf("two leads pulled %d partners, want 2", got)
	}
}

// A fired Pull cluster places each partner's own rate — at least its Min copies —
// when the lead rolls in, and stays dormant without the lead. Mean equal to Min
// makes the count exactly Min, so the per-partner rates are checkable.
func TestPullPlacesPerPartnerCount(t *testing.T) {
	pull := ClusterMembership{
		Name:     "Troop",
		Strategy: Pull,
		Trigger:  ByLead,
	}
	lead := clusterMember("Lead", engine.Untamed, pull)
	lead.Profile.Cluster.Lead = true
	ape := pull
	ape.Min, ape.Mean = 2, 2
	queen := pull
	queen.Min, queen.Mean = 1, 1
	set := NewSet("S", []Card{
		lead,
		clusterMember("Ape", engine.Untamed, ape),
		clusterMember("Queen", engine.Untamed, queen),
		mkCard("FL", engine.Untamed, engine.Common),
	}, Tuning{RarityWeights: map[engine.Rarity]float64{engine.Common: 1}})

	g := &generator{
		set:    set,
		r:      rand.New(rand.NewSource(1)),
		placed: map[string]bool{},
	}
	pod := g.expandPodClusters(firePod(g, set, engine.Untamed, "Lead"))
	if got := countMember(pod, "Ape"); got != 2 {
		t.Fatalf("pulled %d Apes, want 2", got)
	}
	if got := countMember(pod, "Queen"); got != 1 {
		t.Fatalf("pulled %d Queens, want 1", got)
	}

	// No lead present: the ByLead cluster does not fire.
	g = &generator{
		set:    set,
		r:      rand.New(rand.NewSource(2)),
		placed: map[string]bool{},
	}
	g.deckHouses[0] = engine.Untamed
	pod = g.expandPodClusters(g.fillPod(engine.Untamed))
	if distinctMembers(pod, set.clusters["Troop"]) != 0 {
		t.Fatal("Pull fired without its lead")
	}
}

// pullCount returns exactly Min when Mean equals Min, and is capped at PodSize.
func TestPullCount(t *testing.T) {
	g := &generator{r: rand.New(rand.NewSource(7))}

	flat := ClusterMembership{
		Min:  2,
		Mean: 2,
	}
	for range 20 {
		if n := g.pullCount(flat); n != 2 {
			t.Fatalf("Mean==Min returned %d, want 2", n)
		}
	}

	// A Min past the pod ceiling always caps.
	capped := ClusterMembership{
		Min:  PodSize + 1,
		Mean: PodSize + 1,
	}
	if n := g.pullCount(capped); n != PodSize {
		t.Fatalf("capped count %d, want %d", n, PodSize)
	}
}

// The pod-local pass leaves a deck-wide OnePerHouse cluster untouched.
func TestExpandPodClustersSkipsOnePerHouse(t *testing.T) {
	set := NewSet("S", []Card{
		shardMember("Shard-B", engine.Brobnar),
		shardMember("Shard-D", engine.Dis),
		mkCard("FB", engine.Brobnar, engine.Common),
		mkCard("FD", engine.Dis, engine.Common),
	}, Tuning{RarityWeights: map[engine.Rarity]float64{engine.Common: 1}})

	g := &generator{
		set:    set,
		r:      rand.New(rand.NewSource(1)),
		placed: map[string]bool{},
	}
	pod := firePod(g, set, engine.Brobnar, "Shard-B")
	before := countMember(pod, "Shard-B")
	pod = g.expandPodClusters(pod)
	if countMember(pod, "Shard-B") != before {
		t.Fatal("pod-local pass touched a OnePerHouse cluster")
	}
}
