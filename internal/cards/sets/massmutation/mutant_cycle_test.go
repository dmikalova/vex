package massmutation

import (
	"strings"
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	"github.com/dmikalova/vactrol/internal/engine"
)

// realHouses is every playable KeyForge house, in enum order.
func realHouses() []engine.House {
	var houses []engine.House
	for h := engine.HouseNone + 1; int(h) < engine.NumHouses; h++ {
		houses = append(houses, h)
	}
	return houses
}

// TestMutantHousesComplete enforces the cycle's promise: every house — including
// Brobnar and Mars, which Mass Mutation never prints — carries a full mutant
// contribution, so a later set can mutate any house.
func TestMutantHousesComplete(t *testing.T) {
	var noTrait card.Trait
	for _, h := range realHouses() {
		m, ok := mutantHouses[h]
		if !ok {
			t.Errorf("house %v has no mutant contribution", h)
			continue
		}
		if m.prefix == "" || m.suffix == "" {
			t.Errorf("house %v is missing a name fragment: %+v", h, m)
		}
		if m.trait == noTrait {
			t.Errorf("house %v is missing a creature trait", h)
		}
		if m.power <= 0 {
			t.Errorf("house %v is missing a power contribution", h)
		}
	}
}

// TestMutantComposition spot-checks that a mutant is the sum of its two houses:
// power and armor add, keywords and abilities union, and the trait and house come
// from the suffix.
func TestMutantComposition(t *testing.T) {
	cases := []struct {
		name           string
		prefix, suffix engine.House
		wantName       string
		wantPower      int
		wantArmor      int
		wantHouse      engine.House
		wantTrait      card.Trait
		wantKeywords   []card.KeywordValue
		wantTriggers   []engine.Trigger
	}{
		{
			name: "Daemo-Bot", prefix: card.House.Dis, suffix: card.House.Logos,
			wantName: "Daemo-Bot", wantPower: 3, wantHouse: card.House.Logos,
			wantTrait:    card.Traits.Scientist,
			wantTriggers: []engine.Trigger{card.Trigger.Reap, card.Trigger.Destroyed},
		},
		{
			name: "Dino-Knight", prefix: card.House.Saurian, suffix: card.House.Sanctum,
			wantName: "Dino-Knight", wantPower: 6, wantArmor: 2, wantHouse: card.House.Sanctum,
			wantTrait:    card.Traits.Knight,
			wantTriggers: []engine.Trigger{card.Trigger.Play},
		},
		{
			name: "Xeno-Beast", prefix: card.House.StarAlliance, suffix: card.House.Untamed,
			wantName: "Xeno-Beast", wantPower: 4, wantHouse: card.House.Untamed,
			wantTrait:    card.Traits.Beast,
			wantKeywords: []card.KeywordValue{card.Keyword.Skirmish},
			wantTriggers: []engine.Trigger{card.Trigger.Fight},
		},
		{
			name: "Umbra-Bot", prefix: card.House.Shadows, suffix: card.House.Logos,
			wantName: "Umbra-Bot", wantPower: 3, wantHouse: card.House.Logos,
			wantTrait:    card.Traits.Scientist,
			wantKeywords: []card.KeywordValue{card.Keyword.Elusive},
			wantTriggers: []engine.Trigger{card.Trigger.Reap},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			def := card.Build(
				mutantName(tc.prefix, tc.suffix),
				tc.suffix,
				card.Type.Creature,
				card.Rarity.Uncommon,
				mutantOptions(tc.prefix, tc.suffix)...,
			)
			if def.Name != tc.wantName {
				t.Errorf("name = %q, want %q", def.Name, tc.wantName)
			}
			if def.Power != tc.wantPower {
				t.Errorf("power = %d, want %d", def.Power, tc.wantPower)
			}
			if def.Armor != tc.wantArmor {
				t.Errorf("armor = %d, want %d", def.Armor, tc.wantArmor)
			}
			if def.House != tc.wantHouse {
				t.Errorf("house = %v, want %v", def.House, tc.wantHouse)
			}
			if !hasTrait(def.Traits, card.Traits.Mutant) {
				t.Error("missing the Mutant trait")
			}
			if !hasTrait(def.Traits, tc.wantTrait) {
				t.Errorf("missing the suffix trait %v", tc.wantTrait)
			}
			for _, kw := range tc.wantKeywords {
				if !hasKeyword(def.Keywords, kw) {
					t.Errorf("missing keyword %v", kw)
				}
			}
			for _, tr := range tc.wantTriggers {
				if !hasTrigger(def.Abilities, tr) {
					t.Errorf("missing an ability triggered on %v", tr)
				}
			}
		})
	}
}

// TestMutantAlienDigsToBottom checks the Star Alliance mutant ability renders the
// bottom-of-deck routing this cycle added to the engine.
func TestMutantAlienDigsToBottom(t *testing.T) {
	def := card.Build(
		mutantName(card.House.StarAlliance, card.House.Logos),
		card.House.Logos,
		card.Type.Creature,
		card.Rarity.Uncommon,
		mutantOptions(card.House.StarAlliance, card.House.Logos)...,
	)
	var fight card.Ability
	for _, a := range def.Abilities {
		if a.Trigger == card.Trigger.Fight {
			fight = a
		}
	}
	if fight.Effect == nil {
		t.Fatal("Xeno-Bot has no Fight ability")
	}
	if got := fight.Effect.Text(); !strings.Contains(got, "on the bottom of your deck") {
		t.Errorf("Fight text = %q, want it to mention the bottom of the deck", got)
	}
}

// TestMutantProvenanceGrid checks the printed grid is exactly the 42 native
// mutants: seven suffix houses, each paired with the six other native houses, and
// every collector number distinct.
func TestMutantProvenanceGrid(t *testing.T) {
	seen := map[string]string{}
	total := 0
	for suffix, byPrefix := range mutantProvenance {
		if len(byPrefix) != 6 {
			t.Errorf("suffix %v has %d prefixes, want 6", suffix, len(byPrefix))
		}
		for prefix, number := range byPrefix {
			total++
			if prefix == suffix {
				t.Errorf("suffix %v paired with itself", suffix)
			}
			if _, ok := mutantProvenance[prefix]; !ok {
				t.Errorf("prefix %v is not a native mutant house", prefix)
			}
			if prior, dup := seen[number]; dup {
				t.Errorf("collector number %q used by %s and %s",
					number, prior, mutantName(prefix, suffix))
			}
			seen[number] = mutantName(prefix, suffix)
		}
	}
	if total != 42 {
		t.Errorf("printed %d mutants, want 42", total)
	}
}

// TestMutantCycleBuildsEveryHousePair builds a mutant for every ordered pair of
// distinct houses — all nine, in both positions — so every house's contribution,
// including Brobnar's and Mars's dormant abilities, is exercised and validated.
func TestMutantCycleBuildsEveryHousePair(t *testing.T) {
	houses := realHouses()
	for _, prefix := range houses {
		for _, suffix := range houses {
			if prefix == suffix {
				continue
			}
			def := card.Build(
				mutantName(prefix, suffix),
				suffix,
				card.Type.Creature,
				card.Rarity.Uncommon,
				mutantOptions(prefix, suffix)...,
			)
			wantPower := mutantHouses[prefix].power + mutantHouses[suffix].power
			if def.Power != wantPower {
				t.Errorf("%s power = %d, want %d", def.Name, def.Power, wantPower)
			}
			wantArmor := mutantHouses[prefix].armor + mutantHouses[suffix].armor
			if def.Armor != wantArmor {
				t.Errorf("%s armor = %d, want %d", def.Name, def.Armor, wantArmor)
			}
			if def.House != suffix {
				t.Errorf("%s house = %v, want %v", def.Name, def.House, suffix)
			}
			for _, a := range def.Abilities {
				if a.Effect.Text() == "" {
					t.Errorf("%s has an ability with empty text", def.Name)
				}
			}
		}
	}
}

func hasTrait(traits []card.Trait, want card.Trait) bool {
	for _, tr := range traits {
		if tr == want {
			return true
		}
	}
	return false
}

func hasKeyword(keywords []card.KeywordValue, want card.KeywordValue) bool {
	for _, kw := range keywords {
		if kw == want {
			return true
		}
	}
	return false
}

func hasTrigger(abilities []card.Ability, want engine.Trigger) bool {
	for _, a := range abilities {
		if a.Trigger == want {
			return true
		}
	}
	return false
}
