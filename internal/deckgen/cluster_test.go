package deckgen

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/dmikalova/vactrol/internal/engine"
)

// shardCluster is a OnePerHouse / ByAnyMember cluster, the Shards' shape.
var shardCluster = ClusterMembership{Name: "Shard", Strategy: OnePerHouse, Trigger: ByAnyMember}

func shardMember(name string, h engine.House) Card {
	c := mkCard(name, h, engine.Rare)
	c.Profile.Cluster = shardCluster
	c.Profile.OneCopyPerDeck = true
	return c
}

// A drawn cluster member pulls its whole OnePerHouse cycle in: every pod ends up
// holding its own House's member.
func TestOnePerHousePlacesEveryHouse(t *testing.T) {
	set := NewSet("S", []Card{
		shardMember("Shard-B", engine.Brobnar),
		shardMember("Shard-D", engine.Dis),
		shardMember("Shard-L", engine.Logos),
		mkCard("FB", engine.Brobnar, engine.Common),
		mkCard("FD", engine.Dis, engine.Common),
		mkCard("FL", engine.Logos, engine.Common),
	}, Tuning{RarityWeights: map[engine.Rarity]float64{engine.Common: 1}})

	houses := []engine.House{engine.Brobnar, engine.Dis, engine.Logos}
	want := map[engine.House]string{
		engine.Brobnar: "Shard-B",
		engine.Dis:     "Shard-D",
		engine.Logos:   "Shard-L",
	}
	g := &generator{set: set, r: rand.New(rand.NewSource(1)), placed: map[string]bool{}}
	deck := Deck{Set: "S"}
	for i, h := range houses {
		g.deckHouses[i] = h
		deck.Pods[i] = g.fillPod(h)
	}
	// Only Common rolls, so no member is drawn; plant one to fire the cycle.
	deck.Pods[0].Slots[0] = Slot{Rarity: engine.Rare, Card: set.byName["Shard-B"].Def}

	g.expandClusters(&deck)

	for i, h := range houses {
		if !podHas(deck.Pods[i], want[h]) {
			t.Fatalf("pod %d (%v) is missing %s", i, h, want[h])
		}
	}
}

// With no member in the deck the cluster does not fire: no member is placed.
func TestOnePerHouseDormantWithoutMember(t *testing.T) {
	set := NewSet("S", []Card{
		shardMember("Shard-B", engine.Brobnar),
		shardMember("Shard-D", engine.Dis),
		shardMember("Shard-L", engine.Logos),
		mkCard("FB", engine.Brobnar, engine.Common),
		mkCard("FD", engine.Dis, engine.Common),
		mkCard("FL", engine.Logos, engine.Common),
	}, Tuning{RarityWeights: map[engine.Rarity]float64{engine.Common: 1}})

	houses := []engine.House{engine.Brobnar, engine.Dis, engine.Logos}
	g := &generator{set: set, r: rand.New(rand.NewSource(2)), placed: map[string]bool{}}
	deck := Deck{Set: "S"}
	for i, h := range houses {
		g.deckHouses[i] = h
		deck.Pods[i] = g.fillPod(h)
	}
	g.expandClusters(&deck)

	for i, name := range []string{"Shard-B", "Shard-D", "Shard-L"} {
		if podHas(deck.Pods[i], name) {
			t.Fatalf("pod %d holds %s but no member fired the cluster", i, name)
		}
	}
}

// A OnePerHouse cluster missing a member for a House the set can deck fails the
// build at NewSet — the complete-by-construction gate.
func TestOnePerHouseCompletenessGate(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected NewSet to panic on a House with no Shard")
		}
		if msg, _ := r.(string); !strings.Contains(msg, "OnePerHouse") {
			t.Fatalf("panic = %v, want it to name the gate", r)
		}
	}()
	// Logos has cards but no Shard, so the cycle is incomplete.
	NewSet("S", []Card{
		shardMember("Shard-B", engine.Brobnar),
		shardMember("Shard-D", engine.Dis),
		mkCard("FL", engine.Logos, engine.Common),
	}, Tuning{RarityWeights: map[engine.Rarity]float64{engine.Common: 1}})
}

func podHas(pod HousePod, name string) bool {
	for _, s := range pod.Slots {
		if s.Card.Name == name {
			return true
		}
	}
	return false
}

func clusterMember(name string, h engine.House, m ClusterMembership) Card {
	c := mkCard(name, h, engine.Rare)
	c.Profile.Cluster = m
	return c
}

// A catalog-wide ClusterPool resolves a deck-wide cluster across every House a
// deck can reach, including a foreign House an errant pod brings in. The pool
// carries a Saurian Shard the drawing set itself does not print; an errant Saurian
// pod still receives it. A non-OnePerHouse cluster in the pool (the Horsemen) is
// ignored by the deck-wide gate and pass, and a nil pool is a no-op.
func TestCrossClusterErrantHouse(t *testing.T) {
	horseman := ClusterMembership{Name: "Horsemen", Strategy: WholePool, Trigger: ByAnyMember}
	pool := NewClusterPool([]Card{
		shardMember("Shard-B", engine.Brobnar),
		shardMember("Shard-D", engine.Dis),
		shardMember("Shard-S", engine.Saurian),
		clusterMember("Horse-B", engine.Brobnar, horseman),
		clusterMember("Horse-D", engine.Dis, horseman),
	})
	set := NewSet("S", []Card{
		shardMember("Shard-B", engine.Brobnar),
		shardMember("Shard-D", engine.Dis),
		mkCard("FB", engine.Brobnar, engine.Common),
		mkCard("FD", engine.Dis, engine.Common),
	}, Tuning{RarityWeights: map[engine.Rarity]float64{engine.Common: 1}})
	// A legacy card of a foreign House makes Saurian an errant House. The gate
	// skips the WholePool Horsemen (not deck-wide) and passes on the Shards.
	set = set.WithLegacy(NewLegacy([]LegacyEntry{
		{Card: mkCard("Sau", engine.Saurian, engine.Common), Set: "Other"},
	})).WithClusters(pool)
	// A nil pool is a no-op: neither gate nor deck-wide resolution.
	set.WithClusters(nil)

	g := &generator{set: set, r: rand.New(rand.NewSource(1)), placed: map[string]bool{}}
	deck := Deck{Set: "S"}
	houses := []engine.House{engine.Brobnar, engine.Dis, engine.Saurian}
	for i, h := range houses {
		g.deckHouses[i] = h
		deck.Pods[i] = g.fillPod(h)
	}
	// Plant a member to fire the cross-set cycle.
	deck.Pods[0].Slots[0] = Slot{
		Rarity: engine.Rare,
		Card:   pool.clusters["Shard"].byHouse[engine.Brobnar].Def,
	}

	g.expandClusters(&deck)

	if !podHas(deck.Pods[2], "Shard-S") {
		t.Fatal("errant Saurian pod is missing its cross-set Shard")
	}
}

// The cross-set cluster gate fails the build when the pool has no member for a
// House the set can deck as an errant House.
func TestCrossClusterGateMissingHouse(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected WithClusters to panic on an errant House with no member")
		}
		if msg, _ := r.(string); !strings.Contains(msg, "cross-set OnePerHouse") {
			t.Fatalf("panic = %v, want it to name the cross-set gate", r)
		}
	}()
	pool := NewClusterPool([]Card{
		shardMember("Shard-B", engine.Brobnar),
		shardMember("Shard-D", engine.Dis),
	})
	set := NewSet("S", []Card{
		shardMember("Shard-B", engine.Brobnar),
		shardMember("Shard-D", engine.Dis),
		mkCard("FB", engine.Brobnar, engine.Common),
		mkCard("FD", engine.Dis, engine.Common),
	}, Tuning{RarityWeights: map[engine.Rarity]float64{engine.Common: 1}})
	// Saurian is an errant House but the pool has no Saurian Shard.
	set.WithLegacy(NewLegacy([]LegacyEntry{
		{Card: mkCard("Sau", engine.Saurian, engine.Common), Set: "Other"},
	})).WithClusters(pool)
}

// Draftable is the reservoir-set predicate: a housed non-Connected card is
// draftable; a Houseless, Connected, or HouseNone card is not.
func TestDraftable(t *testing.T) {
	cases := []struct {
		name string
		card Card
		want bool
	}{
		{"housed common", mkCard("C", engine.Brobnar, engine.Common), true},
		{"connected", mkCard("K", engine.Brobnar, engine.Connected), false},
		{"houseless card", mkCard("N", engine.HouseNone, engine.Common), false},
		{
			"special",
			Card{
				Def:     engine.NewCard("S", engine.Brobnar, engine.Creature, engine.Common),
				Profile: GenerationProfile{Houseless: true},
			},
			false,
		},
	}
	for _, tc := range cases {
		if got := Draftable(tc.card); got != tc.want {
			t.Errorf("Draftable(%s) = %v, want %v", tc.name, got, tc.want)
		}
	}
}

// connectedMember is a cluster member of Rarity.Connected: it never rolls in the
// pool, so a cluster made only of these can never fire.
func connectedMember(name string, h engine.House, m ClusterMembership) Card {
	c := Card{Def: engine.NewCard(name, h, engine.Creature, engine.Connected, engine.WithPower(3))}
	c.Profile.Cluster = m
	return c
}

// NewSet fails loudly on a cluster whose strategy, trigger, lead, or RandomCount
// range is unsatisfiable.
func TestClusterValidation(t *testing.T) {
	cases := []struct {
		name  string
		cards []Card
		want  string
	}{
		{
			name: "no strategy",
			cards: []Card{
				clusterMember(
					"A",
					engine.Brobnar,
					ClusterMembership{Name: "C", Trigger: ByAnyMember},
				),
			},
			want: "no strategy",
		},
		{
			name: "no trigger",
			cards: []Card{
				clusterMember(
					"A",
					engine.Brobnar,
					ClusterMembership{Name: "C", Strategy: WholePool},
				),
			},
			want: "no trigger",
		},
		{
			name: "ByLead without lead",
			cards: []Card{
				clusterMember(
					"A",
					engine.Brobnar,
					ClusterMembership{Name: "C", Strategy: WholePool, Trigger: ByLead},
				),
			},
			want: "no lead member",
		},
		{
			name: "RandomCount range",
			cards: []Card{clusterMember("A", engine.Brobnar, ClusterMembership{
				Name: "C", Strategy: RandomCount, Trigger: ByAnyMember, Min: 0, Max: 0,
			})},
			want: "RandomCount",
		},
		{
			name: "SelfPull two members",
			cards: []Card{
				clusterMember("A", engine.Brobnar, ClusterMembership{
					Name: "C", Strategy: SelfPull, Trigger: ByAnyMember, Min: 1, Mean: 2,
				}),
				clusterMember("B", engine.Brobnar, ClusterMembership{
					Name: "C", Strategy: SelfPull, Trigger: ByAnyMember, Min: 1, Mean: 2,
				}),
			},
			want: "SelfPull",
		},
		{
			name: "SelfPull min below one",
			cards: []Card{clusterMember("A", engine.Brobnar, ClusterMembership{
				Name: "C", Strategy: SelfPull, Trigger: ByAnyMember, Min: 0, Mean: 2,
			})},
			want: "SelfPull",
		},
		{
			name: "SelfPull mean below min",
			cards: []Card{clusterMember("A", engine.Brobnar, ClusterMembership{
				Name: "C", Strategy: SelfPull, Trigger: ByAnyMember, Min: 3, Mean: 2,
			})},
			want: "SelfPull",
		},
		{
			name: "SelfPull mean above pod",
			cards: []Card{clusterMember("A", engine.Brobnar, ClusterMembership{
				Name: "C", Strategy: SelfPull, Trigger: ByAnyMember, Min: 3, Mean: PodSize + 1,
			})},
			want: "SelfPull",
		},
		{
			name: "ByLead lead cannot roll",
			cards: []Card{connectedMember("A", engine.Brobnar, ClusterMembership{
				Name: "C", Strategy: WholePool, Trigger: ByLead, Lead: true,
			})},
			want: "can never fire",
		},
		{
			name: "ByAnyMember all connected",
			cards: []Card{connectedMember("A", engine.Brobnar, ClusterMembership{
				Name: "C", Strategy: WholePool, Trigger: ByAnyMember,
			})},
			want: "can never fire",
		},
		{
			name: "PullExact needs ByLead",
			cards: []Card{
				clusterMember("A", engine.Brobnar, ClusterMembership{
					Name: "C", Strategy: PullExact, Trigger: ByAnyMember,
				}),
				clusterMember("B", engine.Brobnar, ClusterMembership{
					Name: "C", Strategy: PullExact, Trigger: ByAnyMember,
				}),
			},
			want: "PullExact",
		},
		{
			name: "PullExact needs a partner",
			cards: []Card{clusterMember("A", engine.Brobnar, ClusterMembership{
				Name: "C", Strategy: PullExact, Trigger: ByLead, Lead: true,
			})},
			want: "PullExact",
		},
		{
			name: "Pull needs ByLead",
			cards: []Card{
				clusterMember("A", engine.Brobnar, ClusterMembership{
					Name: "C", Strategy: Pull, Trigger: ByAnyMember,
				}),
				clusterMember("B", engine.Brobnar, ClusterMembership{
					Name: "C", Strategy: Pull, Trigger: ByAnyMember,
				}),
			},
			want: "Pull cluster",
		},
		{
			name: "Pull needs a partner",
			cards: []Card{clusterMember("A", engine.Brobnar, ClusterMembership{
				Name: "C", Strategy: Pull, Trigger: ByLead, Lead: true,
			})},
			want: "Pull cluster",
		},
		{
			name: "Pull mean below min",
			cards: []Card{
				clusterMember("A", engine.Brobnar, ClusterMembership{
					Name: "C", Strategy: Pull, Trigger: ByLead, Lead: true,
				}),
				clusterMember("B", engine.Brobnar, ClusterMembership{
					Name: "C", Strategy: Pull, Trigger: ByLead, Min: 3, Mean: 1,
				}),
			},
			want: "Pull cluster",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if r == nil {
					t.Fatal("expected NewSet to panic")
				}
				if msg, _ := r.(string); !strings.Contains(msg, tc.want) {
					t.Fatalf("panic = %v, want it to mention %q", r, tc.want)
				}
			}()
			NewSet("S", append(tc.cards, mkCard("F", engine.Brobnar, engine.Common)),
				Tuning{RarityWeights: map[engine.Rarity]float64{engine.Common: 1}})
		})
	}
}

// A valid WholePool cluster and a valid RandomCount cluster pass validation, and
// the deck-wide pass leaves them alone (they resolve in the pod pass, not here).
func TestNonOnePerHouseClustersSkipDeckWide(_ *testing.T) {
	whole := ClusterMembership{Name: "Horsemen", Strategy: WholePool, Trigger: ByLead}
	lead := clusterMember("Lead", engine.Brobnar, whole)
	lead.Profile.Cluster.Lead = true
	sins := ClusterMembership{
		Name:     "Sins",
		Strategy: RandomCount,
		Trigger:  ByAnyMember,
		Min:      1,
		Max:      2,
	}
	set := NewSet("S", []Card{
		lead,
		clusterMember("Rider", engine.Brobnar, whole),
		clusterMember("Sin1", engine.Dis, sins),
		clusterMember("Sin2", engine.Dis, sins),
		mkCard("FB", engine.Brobnar, engine.Common),
		mkCard("FD", engine.Dis, engine.Common),
	}, Tuning{RarityWeights: map[engine.Rarity]float64{engine.Common: 1}})

	g := &generator{set: set, r: rand.New(rand.NewSource(1)), placed: map[string]bool{}}
	deck := Deck{Set: "S"}
	for i, h := range []engine.House{engine.Brobnar, engine.Dis} {
		g.deckHouses[i] = h
		deck.Pods[i] = g.fillPod(h)
	}
	g.expandClusters(&deck) // non-OnePerHouse clusters are skipped, no panic.
}

// A ByLead OnePerHouse cluster fires only for its lead: the lead pulls the cycle
// in, a lone non-lead member does not.
func TestOnePerHouseByLead(t *testing.T) {
	cyc := ClusterMembership{Name: "C", Strategy: OnePerHouse, Trigger: ByLead}
	brob := clusterMember("C-B", engine.Brobnar, cyc)
	brob.Profile.Cluster.Lead = true
	set := NewSet("S", []Card{
		brob,
		clusterMember("C-D", engine.Dis, cyc),
		mkCard("FB", engine.Brobnar, engine.Common),
		mkCard("FD", engine.Dis, engine.Common),
	}, Tuning{RarityWeights: map[engine.Rarity]float64{engine.Common: 1}})
	houses := []engine.House{engine.Brobnar, engine.Dis}

	fresh := func() (*generator, Deck) {
		g := &generator{set: set, r: rand.New(rand.NewSource(3)), placed: map[string]bool{}}
		deck := Deck{Set: "S"}
		for i, h := range houses {
			g.deckHouses[i] = h
			deck.Pods[i] = g.fillPod(h)
		}
		return g, deck
	}

	// The lead present fires the cycle.
	g, deck := fresh()
	deck.Pods[0].Slots[0] = Slot{Rarity: engine.Rare, Card: set.byName["C-B"].Def}
	g.expandClusters(&deck)
	if !podHas(deck.Pods[1], "C-D") {
		t.Fatal("lead present but Dis pod missing its member")
	}

	// A non-lead member alone does not fire it.
	g, deck = fresh()
	deck.Pods[1].Slots[0] = Slot{Rarity: engine.Rare, Card: set.byName["C-D"].Def}
	g.expandClusters(&deck)
	if podHas(deck.Pods[0], "C-B") {
		t.Fatal("non-lead member fired a ByLead cycle")
	}
}

// At MaverickRate a pod without its own member receives another House's member,
// rehoused as a maverick. A pod short a House (an unfilled HouseNone pod) is
// skipped.
func TestOnePerHouseMaverickSubstitution(t *testing.T) {
	set := NewSet("S", []Card{
		shardMember("Shard-B", engine.Brobnar),
		shardMember("Shard-D", engine.Dis),
		mkCard("FB", engine.Brobnar, engine.Common),
		mkCard("FD", engine.Dis, engine.Common),
	}, Tuning{RarityWeights: map[engine.Rarity]float64{engine.Common: 1}, MaverickRate: 1})

	g := &generator{set: set, r: rand.New(rand.NewSource(4)), placed: map[string]bool{}}
	deck := Deck{Set: "S"}
	for i, h := range []engine.House{engine.Brobnar, engine.Dis} {
		g.deckHouses[i] = h
		deck.Pods[i] = g.fillPod(h)
	}
	// Pod 2 is left HouseNone, so expandClusters skips it.
	deck.Pods[0].Slots[0] = Slot{Rarity: engine.Rare, Card: set.byName["Shard-B"].Def}

	g.expandClusters(&deck)

	// Pod 1 lacked its member and rolled maverick, so it holds the other House's
	// member as a maverick.
	maverick := false
	for _, s := range deck.Pods[1].Slots {
		if s.Maverick && (s.Card.Name == "Shard-B" || s.Card.Name == "Shard-D") {
			maverick = true
		}
	}
	if !maverick {
		t.Fatalf("Dis pod has no maverick cluster member: %+v", deck.Pods[1])
	}
}
