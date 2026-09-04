package deckgen

import (
	"math/rand"
	"testing"

	"github.com/dmikalova/vactrol/internal/engine"
)

func mkCard(name string, h engine.House, rr engine.Rarity) Card {
	return Card{Def: engine.NewCard(name, h, engine.Creature, rr, engine.WithPower(3))}
}

type matFunc func(SlotContext, *rand.Rand) engine.CardDefinition

func (f matFunc) Materialize(
	ctx SlotContext,
	r *rand.Rand,
) engine.CardDefinition {
	return f(ctx, r)
}

// A Houseless Special card, always rolled, is stamped with the pod's house.
func TestSpecialOverlay(t *testing.T) {
	special := Card{
		Def: engine.NewCard(
			"Special",
			engine.HouseNone,
			engine.Creature,
			engine.Special,
			engine.WithPower(3),
		),
		Profile: GenerationProfile{Houseless: true},
	}
	set := NewSet("S", []Card{
		special,
		mkCard("FB", engine.Brobnar, engine.Common),
		mkCard("FD", engine.Dis, engine.Common),
		mkCard("FL", engine.Logos, engine.Common),
	}, Tuning{
		RarityWeights: map[engine.Rarity]float64{engine.Common: 1},
		SpecialRate:   1,
	})
	for _, pod := range Generate(set, 1).Pods {
		for _, s := range pod.Slots {
			if !s.Special || s.Card.Name != "Special" || s.Card.House != pod.House {
				t.Fatalf("special slot = %+v (pod %v)", s, pod.House)
			}
		}
	}
}

// An always-maverick draw pulls from another house and rehouses to the pod.
func TestMaverickDraw(t *testing.T) {
	set := NewSet("S", []Card{
		mkCard("B", engine.Brobnar, engine.Common),
		mkCard("D", engine.Dis, engine.Common),
		mkCard("L", engine.Logos, engine.Common),
		mkCard("M", engine.Mars, engine.Common),
	}, Tuning{RarityWeights: map[engine.Rarity]float64{engine.Common: 1}, MaverickRate: 1})
	for _, pod := range Generate(set, 1).Pods {
		for _, s := range pod.Slots {
			if s.Card.House != pod.House {
				t.Fatalf("maverick not rehoused: %v in pod %v", s.Card.House, pod.House)
			}
		}
	}
}

// A single-house set cannot produce mavericks: otherHouse has no alternative.
func TestMaverickSingleHouse(t *testing.T) {
	set := NewSet("S", []Card{
		mkCard("B1", engine.Brobnar, engine.Common),
		mkCard("B2", engine.Brobnar, engine.Common),
	}, Tuning{RarityWeights: map[engine.Rarity]float64{engine.Common: 1}, MaverickRate: 1})
	for _, s := range Generate(set, 1).Pods[0].Slots {
		if s.Maverick {
			t.Fatal("single-house set produced a maverick")
		}
	}
}

// The duplicate-pull copies an already-placed same-rarity card.
func TestDuplicatePull(t *testing.T) {
	set := NewSet("S", []Card{
		mkCard("A", engine.Brobnar, engine.Common),
		mkCard("B", engine.Brobnar, engine.Common),
	}, Tuning{
		RarityWeights: map[engine.Rarity]float64{engine.Common: 1},
		DuplicateRate: map[engine.Rarity]float64{engine.Common: 1},
	})
	if len(Generate(set, 1).Pods[0].Slots) != PodSize {
		t.Fatal("pod not filled")
	}
}

// The duplicate-pull only copies a card the pod could legally hold twice, so a
// placed card of another rarity, an empty slot, and a one-copy-per-deck card are
// each passed over — leaving nothing to copy, and a fresh draw to fill the slot.
func TestDuplicatePullSkipsIneligible(t *testing.T) {
	unique := mkCard("U", engine.Brobnar, engine.Common)
	unique.Profile.OneCopyPerDeck = true
	g := gen(NewSet("S", []Card{unique}, Tuning{
		DuplicateRate: map[engine.Rarity]float64{engine.Common: 1},
	}))
	placed := []placedCard{
		{card: mkCard("R", engine.Brobnar, engine.Rare), rarity: engine.Rare},
		{rarity: engine.Common},
		{card: unique, rarity: engine.Common},
	}
	if c, ok := g.tryDuplicate(engine.Common, placed); ok {
		t.Fatalf("duplicated the ineligible %q", c.Def.Name)
	}
}

// A single one-copy-per-deck card is placed at most once; the rest of the pod
// exhausts every draw fallback.
func TestOnePerDeckExhaustion(t *testing.T) {
	unique := Card{
		Def: engine.NewCard(
			"U",
			engine.Brobnar,
			engine.Creature,
			engine.Common,
			engine.WithPower(3),
		),
		Profile: GenerationProfile{OneCopyPerDeck: true},
	}
	set := NewSet(
		"S",
		[]Card{unique},
		Tuning{RarityWeights: map[engine.Rarity]float64{engine.Common: 1}},
	)
	count := 0
	for _, s := range Generate(set, 1).Pods[0].Slots {
		if s.Card.Name == "U" {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("one-per-deck placed %d times", count)
	}
}

// A template's Materializer builds the concrete card.
func TestMaterializerTemplate(t *testing.T) {
	tmpl := Card{
		Def: engine.NewCard(
			"Face",
			engine.Brobnar,
			engine.Creature,
			engine.Common,
			engine.WithPower(3),
		),
		Materializer: matFunc(func(ctx SlotContext, _ *rand.Rand) engine.CardDefinition {
			return engine.NewCard(
				"Materialized",
				ctx.House,
				engine.Creature,
				engine.Common,
				engine.WithPower(7),
			)
		}),
	}
	set := NewSet(
		"S",
		[]Card{tmpl},
		Tuning{RarityWeights: map[engine.Rarity]float64{engine.Common: 1}},
	)
	for _, s := range Generate(set, 1).Pods[0].Slots {
		if s.Card.Name != "Materialized" || s.Card.Power != 7 {
			t.Fatalf("template not materialized: %q pow %d", s.Card.Name, s.Card.Power)
		}
	}
}

func TestNewSetDropsHouseNone(t *testing.T) {
	set := NewSet("S", []Card{
		{
			Def: engine.NewCard(
				"NoHouse",
				engine.HouseNone,
				engine.Creature,
				engine.Common,
				engine.WithPower(3),
			),
		},
		mkCard("B", engine.Brobnar, engine.Common),
	}, DefaultTuning())
	if hs := set.Houses(); len(hs) != 1 || hs[0] != engine.Brobnar {
		t.Fatalf("houses = %v, want [Brobnar]", hs)
	}
}

// House exclusions bar two houses from both being chosen; weights bias the draw.
func TestPickHousesExclusionsAndWeights(t *testing.T) {
	tuning := DefaultTuning()
	tuning.HouseWeights = map[engine.House]float64{engine.Brobnar: 10}
	tuning.HouseExclusions = [][2]engine.House{{engine.Dis, engine.Logos}}
	set := NewSet("S", []Card{
		mkCard("B", engine.Brobnar, engine.Common),
		mkCard("D", engine.Dis, engine.Common),
		mkCard("L", engine.Logos, engine.Common),
		mkCard("M", engine.Mars, engine.Common),
	}, tuning)
	for seed := int64(0); seed < 40; seed++ {
		hasDis, hasLogos := false, false
		for _, h := range Generate(set, seed).Houses() {
			hasDis = hasDis || h == engine.Dis
			hasLogos = hasLogos || h == engine.Logos
		}
		if hasDis && hasLogos {
			t.Fatalf("seed %d chose both excluded houses", seed)
		}
	}
}

// A zero Tuning defaults to DefaultTuning; an empty (non-nil) rarity map rolls
// Common.
func TestTuningDefaults(t *testing.T) {
	zero := NewSet("S", []Card{mkCard("B", engine.Brobnar, engine.Common)}, Tuning{})
	if Generate(zero, 1).Pods[0].House != engine.Brobnar {
		t.Fatal("zero tuning did not default")
	}
	empty := NewSet("S", []Card{mkCard("B", engine.Brobnar, engine.Common)},
		Tuning{RarityWeights: map[engine.Rarity]float64{}})
	for _, s := range Generate(empty, 1).Pods[0].Slots {
		if s.Rarity != engine.Common {
			t.Fatalf("empty weights rolled %v, want Common", s.Rarity)
		}
	}
}

// A house with only some rarities falls back to any card of that house when the
// rolled rarity is absent.
func TestDrawRarityFallback(t *testing.T) {
	set := NewSet("S", []Card{mkCard("B", engine.Brobnar, engine.Common)},
		Tuning{RarityWeights: map[engine.Rarity]float64{engine.Rare: 1}})
	for _, s := range Generate(set, 1).Pods[0].Slots {
		if s.Card.Name != "B" {
			t.Fatalf("fallback failed: %q", s.Card.Name)
		}
	}
}
