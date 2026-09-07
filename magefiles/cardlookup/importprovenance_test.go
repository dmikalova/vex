package main

import (
	"reflect"
	"testing"
)

func TestAsciiFold(t *testing.T) {
	cases := map[string]string{
		"Coward\u2019s End":  "Coward's End",
		"\u201cQuoted\u201d": `"Quoted"`,
		"\u00c6mber Skies":   "Aember Skies",
		"caf\u00e9":          "cafe",
		"4\u00d74":           "4x4",
		"en\u2013dash":       "en-dash",
		"em\u2014dash":       "em-dash",
		"and so on\u2026":    "and so on...",
		"plain ascii":        "plain ascii",
	}
	for in, want := range cases {
		if got := asciiFold(in); got != want {
			t.Errorf("asciiFold(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestExpandMarkup(t *testing.T) {
	cases := map[string]string{
		"Play: Place 2<A> on an enemy creature.": "Play: Place 2 Aember on an enemy creature.",
		"Gain 1<A>.":                             "Gain 1 Aember.",
		"Deal 3<D> to a creature.":               "Deal 3 Damage to a creature.",
		"Ready and fight.":                       "Ready and fight.",
		"Line one\u000bLine two":                 "Line one\nLine two",
	}
	for in, want := range cases {
		got, err := expandMarkup(in)
		if err != nil {
			t.Errorf("expandMarkup(%q) error: %v", in, err)
			continue
		}
		if got != want {
			t.Errorf("expandMarkup(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestExpandMarkupUnknownTokenFailsLoud(t *testing.T) {
	if _, err := expandMarkup("Gain 1<Z>."); err == nil {
		t.Error("expandMarkup with unknown token <Z> should error")
	}
	if _, err := expandMarkup("Gain 1<A."); err == nil {
		t.Error("expandMarkup with unterminated token should error")
	}
}

func TestNormalizeCatalogNumber(t *testing.T) {
	cases := map[string]string{
		"4":   "004",
		"151": "151",
		"012": "012",
		"S01": "S01",
		"A21": "A21",
		"P07": "P07",
	}
	for in, want := range cases {
		if got := normalizeCatalogNumber(in); got != want {
			t.Errorf("normalizeCatalogNumber(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestFoldTraits(t *testing.T) {
	if got := foldTraits(
		"Human \u2022 Scientist",
	); !reflect.DeepEqual(
		got,
		[]string{"human", "scientist"},
	) {
		t.Errorf("foldTraits = %v", got)
	}
	if got := foldTraits(""); got != nil {
		t.Errorf("foldTraits(empty) = %v, want nil", got)
	}
}

func TestTransformCardCreature(t *testing.T) {
	got, err := transformCard(mvCard{
		CardTitle:  "Troll",
		House:      "Brobnar",
		CardType:   "Creature",
		CardText:   "After Reap: Gain 2<A>.",
		Traits:     "Giant",
		Power:      "12",
		Armor:      "1",
		Rarity:     "Common",
		CardNumber: "37",
	})
	if err != nil {
		t.Fatal(err)
	}
	want := catalogCard{
		Number: "037",
		Name:   "Troll",
		House:  "brobnar",
		Type:   "creature",
		Rarity: "Common",
		Traits: []string{"giant"},
		Power:  12,
		Armor:  1,
		Amber:  0,
		Text:   "After Reap: Gain 2 Aember.",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("transformCard creature =\n %+v\nwant\n %+v", got, want)
	}
}

func TestTransformCardAnomaly(t *testing.T) {
	got, err := transformCard(mvCard{
		CardTitle:  "Akugyo",
		House:      "Unfathomable",
		CardType:   "Creature",
		Rarity:     "FIXED",
		Power:      "5",
		CardNumber: "S01",
		IsAnomaly:  true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.House != "brobnar" {
		t.Errorf("anomaly house = %q, want brobnar", got.House)
	}
	if got.Rarity != "Special" {
		t.Errorf("anomaly rarity = %q, want Special", got.Rarity)
	}
	if got.Number != "S01" {
		t.Errorf("anomaly number = %q, want S01", got.Number)
	}
}

func TestTransformCardPrintedName(t *testing.T) {
	got, err := transformCard(mvCard{
		CardTitle:  "Bair\u2019s Blessing",
		House:      "Untamed",
		CardType:   "Action",
		Rarity:     "Common",
		CardNumber: "100",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "Bair's Blessing" {
		t.Errorf("folded name = %q", got.Name)
	}
	if got.Printed != "Bair\u2019s Blessing" {
		t.Errorf("printed = %q, want the original title", got.Printed)
	}
	// A non-creature omits power and armor.
	if got.Power != 0 || got.Armor != 0 {
		t.Errorf("non-creature carried power/armor: %+v", got)
	}
}

func TestTransformCardNoPrintedWhenEqual(t *testing.T) {
	got, err := transformCard(mvCard{
		CardTitle:  "Anger",
		House:      "Brobnar",
		CardType:   "Action",
		Rarity:     "Common",
		CardNumber: "001",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.Printed != "" {
		t.Errorf("printed should be empty when name equals title, got %q", got.Printed)
	}
}

func TestHasInteriorGap(t *testing.T) {
	contiguous := []mvCard{{CardNumber: "001"}, {CardNumber: "002"}, {CardNumber: "003"}}
	if hasInteriorGap(contiguous) {
		t.Error("contiguous numbers should not report a gap")
	}
	gapped := []mvCard{{CardNumber: "001"}, {CardNumber: "003"}}
	if !hasInteriorGap(gapped) {
		t.Error("a missing interior number should report a gap")
	}
	lettered := []mvCard{{CardNumber: "001"}, {CardNumber: "002"}, {CardNumber: "S01"}}
	if hasInteriorGap(lettered) {
		t.Error("lettered reference numbers should not count as gaps")
	}
}
