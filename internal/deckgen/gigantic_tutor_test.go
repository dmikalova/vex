package deckgen

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/engine"
)

// tutorCard builds a Houseless PerGigantic tutor pool entry, the way the tutors
// register: no House of its own, pulled into a gigantic's pod and stamped there.
func tutorCard(name string) Card {
	c := Card{Def: engine.NewCard(name, engine.HouseNone, engine.Tactic, engine.Special)}
	c.Profile.Houseless = true
	c.Profile.Cluster = ClusterMembership{
		Name: "Tutors", Strategy: PerGigantic, Trigger: ByAnyMember,
	}
	return c
}

func countNamed(pod HousePod, name string) (n int, house engine.House) {
	for _, s := range pod.Slots {
		if s.Card.Name == name {
			n++
			house = s.Card.House
		}
	}
	return n, house
}

// A pod holding a gigantic base gains a tutor in a free slot, stamped to the pod's
// House, resolved from the set's own PerGigantic cluster.
func TestPlaceGiganticTutorPullsTutor(t *testing.T) {
	base := giganticBaseCard("Deusillus", engine.Saurian)
	tutor := tutorCard("It's Coming...")
	set := NewSet("S", []Card{base, tutor, mkCard("F", engine.Saurian, engine.Common)},
		Tuning{RarityWeights: map[engine.Rarity]float64{engine.Common: 1}})
	g := gen(set)

	pod := HousePod{House: engine.Saurian}
	pod.Slots[0] = Slot{Card: base.Def}
	for i := 1; i < PodSize; i++ {
		pod.Slots[i] = Slot{Card: mkCard("F", engine.Saurian, engine.Common).Def}
	}
	pod = g.placeGiganticArt(pod)
	pod = g.placeGiganticTutor(pod)

	if n, house := countNamed(pod, "It's Coming..."); n != 1 {
		t.Fatalf("pulled %d tutors, want 1", n)
	} else if house != engine.Saurian {
		t.Errorf("tutor House = %v, want stamped to pod House Saurian", house)
	}
}

// Two gigantic bases in one pod each pull their own tutor into a distinct slot,
// resolved from the attached catalog cluster pool (the cross-set path).
func TestPlaceGiganticTutorPullsPerBaseFromPool(t *testing.T) {
	base := giganticBaseCard("Deusillus", engine.Saurian)
	tutor := tutorCard("It's Coming...")
	pool := NewClusterPool([]Card{tutor})
	set := NewSet("S", []Card{base, mkCard("F", engine.Saurian, engine.Common)},
		Tuning{RarityWeights: map[engine.Rarity]float64{engine.Common: 1}}).
		WithClusters(pool)
	g := gen(set)

	pod := HousePod{House: engine.Saurian}
	pod.Slots[0] = Slot{Card: base.Def}
	pod.Slots[1] = Slot{Card: base.Def}
	for i := 2; i < PodSize; i++ {
		pod.Slots[i] = Slot{Card: mkCard("F", engine.Saurian, engine.Common).Def}
	}
	pod = g.placeGiganticTutor(pod)

	if n, _ := countNamed(pod, "It's Coming..."); n != 2 {
		t.Fatalf("pulled %d tutors for 2 bases, want 2", n)
	}
}

// A set with no PerGigantic cluster pulls no tutor.
func TestPlaceGiganticTutorSkipsWithoutCluster(t *testing.T) {
	base := giganticBaseCard("Deusillus", engine.Saurian)
	set := NewSet("S", []Card{base, mkCard("F", engine.Saurian, engine.Common)},
		Tuning{RarityWeights: map[engine.Rarity]float64{engine.Common: 1}})
	g := gen(set)

	pod := HousePod{House: engine.Saurian}
	pod.Slots[0] = Slot{Card: base.Def}
	for i := 1; i < PodSize; i++ {
		pod.Slots[i] = Slot{Card: mkCard("F", engine.Saurian, engine.Common).Def}
	}
	pod = g.placeGiganticTutor(pod)

	if n, _ := countNamed(pod, "It's Coming..."); n != 0 {
		t.Fatalf("pulled %d tutors with no cluster, want 0", n)
	}
}

// A pod with no free ordinary slot leaves the tutor unpulled.
func TestPlaceGiganticTutorSkipsWhenPodFull(t *testing.T) {
	base := giganticBaseCard("Deusillus", engine.Saurian)
	tutor := tutorCard("It's Coming...")
	set := NewSet("S", []Card{base, tutor},
		Tuning{RarityWeights: map[engine.Rarity]float64{engine.Rare: 1}})
	g := gen(set)

	pod := HousePod{House: engine.Saurian}
	for i := range pod.Slots {
		pod.Slots[i] = Slot{Card: base.Def}
	}
	pod = g.placeGiganticTutor(pod)

	if n, _ := countNamed(pod, "It's Coming..."); n != 0 {
		t.Fatalf("pulled %d tutors with no free slot, want 0", n)
	}
}
