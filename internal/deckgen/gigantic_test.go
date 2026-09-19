package deckgen

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/engine"
)

// giganticBaseCard builds a pool entry for a gigantic base half that carries a
// synthetic art half on its profile, the way card.Gigantic registers one.
func giganticBaseCard(name string, h engine.House) Card {
	art := engine.NewCard(
		name, h, engine.Creature, engine.Special,
		engine.WithGiganticRole(engine.GiganticArt),
	)
	base := engine.NewCard(
		name, h, engine.Creature, engine.Rare,
		engine.WithPower(20), engine.WithGiganticRole(engine.GiganticBase),
	)
	c := Card{Def: base}
	c.Profile.GiganticArt = &art
	return c
}

func countRoles(pod HousePod) (bases, arts int) {
	for _, s := range &pod.Slots {
		switch s.Card.GiganticRole {
		case engine.GiganticBase:
			bases++
		case engine.GiganticArt:
			arts++
		}
	}
	return bases, arts
}

// A pod holding a gigantic base half gains its art half in a free slot.
func TestPlaceGiganticArtPlacesBothHalves(t *testing.T) {
	base := giganticBaseCard("Deusillus", engine.Saurian)
	set := NewSet("S", []Card{base, mkCard("F", engine.Saurian, engine.Common)},
		Tuning{RarityWeights: map[engine.Rarity]float64{engine.Common: 1}})
	g := gen(set)

	pod := HousePod{House: engine.Saurian}
	pod.Slots[0] = Slot{Card: base.Def}
	for i := 1; i < PodSize; i++ {
		pod.Slots[i] = Slot{Card: mkCard("F", engine.Saurian, engine.Common).Def}
	}

	pod = g.placeGiganticArt(pod)
	if bases, arts := countRoles(pod); bases != 1 || arts != 1 {
		t.Fatalf("got %d base, %d art halves; want 1 and 1", bases, arts)
	}
}

// A base-role slot whose name is unknown to the set carries no art half, so
// nothing is placed.
func TestPlaceGiganticArtSkipsWhenNoArt(t *testing.T) {
	set := NewSet("S", []Card{mkCard("F", engine.Saurian, engine.Common)},
		Tuning{RarityWeights: map[engine.Rarity]float64{engine.Common: 1}})
	g := gen(set)

	orphan := engine.NewCard("Ghost", engine.Saurian, engine.Creature, engine.Rare,
		engine.WithGiganticRole(engine.GiganticBase))
	pod := HousePod{House: engine.Saurian}
	pod.Slots[0] = Slot{Card: orphan}
	for i := 1; i < PodSize; i++ {
		pod.Slots[i] = Slot{Card: mkCard("F", engine.Saurian, engine.Common).Def}
	}

	pod = g.placeGiganticArt(pod)
	if _, arts := countRoles(pod); arts != 0 {
		t.Fatalf("placed %d art halves for a base with no carried art; want 0", arts)
	}
}

// A pod with no ordinary slot to spend leaves the art half unplaced.
func TestPlaceGiganticArtSkipsWhenPodFull(t *testing.T) {
	base := giganticBaseCard("Deusillus", engine.Saurian)
	set := NewSet("S", []Card{base},
		Tuning{RarityWeights: map[engine.Rarity]float64{engine.Rare: 1}})
	g := gen(set)

	pod := HousePod{House: engine.Saurian}
	for i := range pod.Slots {
		pod.Slots[i] = Slot{Card: base.Def}
	}

	pod = g.placeGiganticArt(pod)
	if _, arts := countRoles(pod); arts != 0 {
		t.Fatalf("placed %d art halves with no free slot; want 0", arts)
	}
}
