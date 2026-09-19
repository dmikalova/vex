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
		"non\u2011Mars":      "non-Mars",
		"Nine\u2010Toes":     "Nine-Toes",
		"a\u00a0b":           "a b",
		"end.\ufeff":         "end.",
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
		"Enhance <A><C><D><R>.":                  "Enhance Aember Capture Damage Draw.",
		"Ready and fight.":                       "Ready and fight.",
		"Line one\u000bLine two":                 "Line one\nLine two",
		"Line one\rLine two":                     "Line one\nLine two",
		"Line one\r\nLine two":                   "Line one\nLine two",
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

func TestExpandMarkupGlyphs(t *testing.T) {
	cases := map[string]string{
		// Inline resource glyphs read as their word with a space off a digit.
		"Gain 1\uF360.":          "Gain 1 Aember.",
		"Deal 5\uF361 to a foe.": "Deal 5 Damage to a foe.",
		// Enhance lists of adjacent glyphs read as separate words.
		"Enhance \uF360\uF360.":       "Enhance Aember Aember.",
		"Enhance \uF361\uF361\uF361.": "Enhance Damage Damage Damage.",
		"Enhance \uF36E.":             "Enhance Draw.",
		// Capture is an F36F+F560 pair and also a lone F565.
		"Enhance \uF36F\uF560\uF36F\uF560.": "Enhance Capture Capture.",
		"Enhance \uF565.":                   "Enhance Capture.",
		// House enhancements (Aember Skies) read as the house name.
		"Enhance \uF379\uF379.": "Enhance Brobnar Brobnar.",
		"Enhance \uF37A.":       "Enhance Dis.",
		"Enhance \uF38A.":       "Enhance Star Alliance.",
		// Discard and the power counter.
		"Enhance \uF372\uF372.":       "Enhance Discard Discard.",
		"Enhance \uF360\uF392\uF392.": "Enhance Aember +1 power counter +1 power counter.",
		// The tide icon is dropped.
		"Skirmish. \uF566 After Fight: Gain 1\uF360.": "Skirmish.  After Fight: Gain 1 Aember.",
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

func TestExpandMarkupUnknownGlyphFailsLoud(t *testing.T) {
	if _, err := expandMarkup("Gain 1\uF3FF."); err == nil {
		t.Error("expandMarkup with unknown glyph U+F3FF should error")
	}
	if _, err := expandMarkup("Enhance \uF36F."); err == nil {
		t.Error("expandMarkup with an unpaired capture glyph should error")
	}
}

func TestExpandEnhanceLetters(t *testing.T) {
	cases := map[string]string{
		"Enhance AA.":     "Enhance Aember Aember.",
		"Enhance PT.":     "Enhance Capture.",
		"Enhance PTDR.":   "Enhance Capture Damage Draw.",
		"Enhance APTDR.":  "Enhance Aember Capture Damage Draw.",
		"Enhance AADDRR.": "Enhance Aember Aember Damage Damage Draw Draw.",
		// A run that does not follow Enhance, or is already a word, is untouched.
		"Enhance Aember.": "Enhance Aember.",
		"Reap: Deal D.":   "Reap: Deal D.",
	}
	for in, want := range cases {
		if got := expandEnhanceLetters(in); got != want {
			t.Errorf("expandEnhanceLetters(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestExpandAndFoldTextStripsKeywordReminders(t *testing.T) {
	cases := map[string]string{
		// Templated keyword reminders are removed, leaving the keyword.
		"Elusive. (The first time this creature is attacked each turn, no damage is dealt.)\nReap: Gain 1<A>.":        "Elusive.\nReap: Gain 1 Aember.",
		"Skirmish. (When you use this creature to fight, it is dealt no damage in return.)":                           "Skirmish.",
		"Taunt. (This creature's neighbors cannot be attacked unless they have taunt.)":                               "Taunt.",
		"Assault 2.(Before this creature attacks, deal 2<D> to the attacked enemy.)":                                  "Assault 2.",
		"Hazardous 5. (Before this creature is attacked, deal 5<D> to the attacking enemy.)":                          "Hazardous 5.",
		"Alpha. (You can only play this card before doing anything else this step.)\nPlay: Draw a card.":              "Alpha.\nPlay: Draw a card.",
		"Omega. (After you play this card, end this step.)":                                                           "Omega.",
		"Deploy. (This creature can enter play anywhere in your battleline.)":                                         "Deploy.",
		"Enhance Damage Damage. (These icons have already been added to cards in your deck.)":                         "Enhance Damage Damage.",
		"Versatile. (This card may be used as if it belonged to the active house.)":                                   "Versatile.",
		"Treachery. (This card enters play under your opponent's control.)":                                           "Treachery.",
		"Splash-attack 4. (When this creature attacks, also deal 4<D> to each of the attacked creature's neighbors).": "Splash-attack 4.",
		// Card-specific parentheticals are kept.
		"Play: Each player loses half their Aember (rounding down the loss). Gain 1 chain.": "Play: Each player loses half their Aember (rounding down the loss). Gain 1 chain.",
		"(Vanilla)": "(Vanilla)",
		// Trailing and pre-newline whitespace is tidied, and double spaces collapse.
		"Destroyed: Gain 5<A>. ":              "Destroyed: Gain 5 Aember.",
		"First line.  \nSecond line.":         "First line.\nSecond line.",
		"Play:  Enemy creatures cannot reap.": "Play: Enemy creatures cannot reap.",
	}
	for in, want := range cases {
		got, err := expandAndFoldText(in)
		if err != nil {
			t.Errorf("expandAndFoldText(%q) error: %v", in, err)
			continue
		}
		if got != want {
			t.Errorf("expandAndFoldText(%q) = %q, want %q", in, got, want)
		}
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

func TestTransformCardGiganticBase(t *testing.T) {
	got, err := transformCard(mvCard{
		CardTitle:  "Deusillus",
		House:      "Saurian",
		CardType:   "Gigantic Creature Base",
		CardText:   "Play: Capture all of your opponent's <A>.",
		Traits:     "Mutant",
		Power:      "20",
		Armor:      "0",
		Rarity:     "Special",
		CardNumber: "197",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.Type != "gigantic creature base" {
		t.Errorf("type = %q, want gigantic creature base", got.Type)
	}
	if got.Power != 20 {
		t.Errorf("power = %d, want 20 (a gigantic base carries stats)", got.Power)
	}
}

func TestMergeCardPrefersGiganticBase(t *testing.T) {
	base := mvCard{
		CardNumber: "197",
		CardTitle:  "Deusillus",
		CardType:   "Gigantic Creature Base",
		Power:      "20",
	}
	art := mvCard{
		CardNumber: "197",
		CardTitle:  "Deusillus",
		CardType:   "Gigantic Creature Art",
	}

	t.Run("art seen first, base replaces it", func(t *testing.T) {
		seen := map[string]int{}
		var got []mvCard
		var isNew bool
		got, isNew = mergeCard(got, seen, art)
		if !isNew {
			t.Error("the first half should count as a new card")
		}
		got, isNew = mergeCard(got, seen, base)
		if isNew {
			t.Error("the second half shares a key, so it is not a new card")
		}
		if len(got) != 1 || got[0].CardType != "Gigantic Creature Base" {
			t.Errorf("base half should win, got %+v", got)
		}
	})

	t.Run("base seen first, art does not displace it", func(t *testing.T) {
		seen := map[string]int{}
		var got []mvCard
		got, _ = mergeCard(got, seen, base)
		got, _ = mergeCard(got, seen, art)
		if len(got) != 1 || got[0].CardType != "Gigantic Creature Base" {
			t.Errorf("base half should stay, got %+v", got)
		}
	})

	t.Run("distinct cards each collect", func(t *testing.T) {
		seen := map[string]int{}
		var got []mvCard
		got, _ = mergeCard(got, seen, base)
		other := mvCard{
			CardNumber: "198",
			CardTitle:  "Other",
		}
		var isNew bool
		got, isNew = mergeCard(got, seen, other)
		if !isNew || len(got) != 2 {
			t.Errorf("a distinct card should be added, got %+v", got)
		}
	})
}
