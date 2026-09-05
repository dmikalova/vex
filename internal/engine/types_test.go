package engine

import "testing"

func TestHouseString(t *testing.T) {
	if Brobnar.String() != "Brobnar" {
		t.Errorf("Brobnar.String() = %q", Brobnar.String())
	}
	if House(99).String() != "Unknown" {
		t.Errorf("House(99).String() = %q, want Unknown", House(99).String())
	}
}

func TestParseHouse(t *testing.T) {
	// Case-insensitive match.
	if h, ok := ParseHouse("bRoBnAr"); !ok || h != Brobnar {
		t.Errorf("ParseHouse(bRoBnAr) = %v, %v; want Brobnar, true", h, ok)
	}
	// "None" is not a selectable house.
	if h, ok := ParseHouse("none"); ok || h != HouseNone {
		t.Errorf("ParseHouse(none) = %v, %v; want HouseNone, false", h, ok)
	}
	// Unknown name.
	if h, ok := ParseHouse("nonsense"); ok || h != HouseNone {
		t.Errorf("ParseHouse(nonsense) = %v, %v; want HouseNone, false", h, ok)
	}
}

func TestTraitString(t *testing.T) {
	if Beast.String() != "Beast" {
		t.Errorf("Beast.String() = %q, want Beast", Beast.String())
	}
	if traitUnset.String() != "" {
		t.Errorf("traitUnset.String() = %q, want empty", traitUnset.String())
	}
	if Trait(999).String() != "" {
		t.Errorf("Trait(999).String() = %q, want empty", Trait(999).String())
	}
}

// TestCardTypeReacts covers the three filter modes a lasting entry's card type
// has: unset matches anything, AnyType means creature-or-artifact, and a named
// type matches only itself.
func TestCardTypeReacts(t *testing.T) {
	cases := []struct {
		filter, subject CardType
		want            bool
	}{
		{TypeUnset, Upgrade, true},
		{AnyType, Creature, true},
		{AnyType, Artifact, true},
		{AnyType, Upgrade, false},
		{AnyType, Tactic, false},
		{Creature, Creature, true},
		{Creature, Artifact, false},
	}
	for _, c := range cases {
		if got := c.filter.reacts(c.subject); got != c.want {
			t.Errorf("%v.reacts(%v) = %v, want %v", c.filter, c.subject, got, c.want)
		}
	}
}

// TestCardTypes pins the real-type enumeration: CardTypes lists exactly the four
// types a card can be, each renders its printed word, and neither the TypeUnset
// nor AnyType sentinel leaks in.
func TestCardTypes(t *testing.T) {
	all := CardTypes()
	want := []CardType{Creature, Tactic, Artifact, Upgrade}
	if len(all) != len(want) {
		t.Fatalf("CardTypes() has %d entries, want %d", len(all), len(want))
	}
	for i, ct := range all {
		if ct != want[i] {
			t.Errorf("CardTypes()[%d] = %v, want %v", i, ct, want[i])
		}
		if ct.String() == "" {
			t.Errorf("card type %d renders empty", ct)
		}
		if ct == TypeUnset || ct == AnyType {
			t.Errorf("CardTypes() leaked the sentinel %v", ct)
		}
	}
}

// TestKeywords pins the enum's three derived views against each other: every
// keyword Keywords() lists must name itself and claim a distinct bit, and the
// unset zero must render empty and claim no bit rather than shifting negatively.
func TestKeywords(t *testing.T) {
	if keywordUnset.String() != "" {
		t.Errorf("keywordUnset.String() = %q, want empty", keywordUnset.String())
	}
	if keywordUnset.bit() != 0 {
		t.Errorf("keywordUnset.bit() = %d, want 0", keywordUnset.bit())
	}

	all := Keywords()
	if len(all) != int(keywordCount)-1 {
		t.Fatalf("Keywords() has %d entries, want %d", len(all), keywordCount-1)
	}
	var seen uint8
	for _, k := range all {
		if k.String() == "" {
			t.Errorf("keyword %d renders empty", k)
		}
		if seen&k.bit() != 0 {
			t.Errorf("%v reuses a bit already taken", k)
		}
		seen |= k.bit()
	}
	if all[0] != Skirmish || all[len(all)-1] != Deploy {
		t.Errorf("Keywords() = %v, want Skirmish first and Deploy last", all)
	}
}

// TestTriggers pins the trigger catalog against its rulebook surface: Triggers
// omits the unset zero, every real trigger names itself, and the two implicit
// triggers (EntersPlay, AfterChooseHouse) name themselves but are unprinted, so
// they head no rulebook entry.
func TestTriggers(t *testing.T) {
	all := Triggers()
	if len(all) != int(triggerCount)-1 {
		t.Fatalf("Triggers() has %d entries, want %d", len(all), triggerCount-1)
	}
	for _, tr := range all {
		if tr == triggerUnset {
			t.Errorf("Triggers() leaked the unset zero")
		}
		if tr.String() == "" {
			t.Errorf("trigger %d names itself empty", tr)
		}
	}
	if TriggerEntersPlay.Printed() || TriggerAfterChooseHouse.Printed() {
		t.Errorf("EntersPlay/AfterChooseHouse should be unprinted")
	}
	if TriggerEntersPlay.String() == "" || TriggerAfterChooseHouse.String() == "" {
		t.Errorf("implicit triggers should still name themselves for logging")
	}
	if !TriggerAfterPlay.Printed() || TriggerAfterPlay.String() != "Play" {
		t.Errorf("Play trigger should be printed and name itself %q", "Play")
	}
	if triggerUnset.String() != "" || triggerUnset.Printed() {
		t.Errorf("triggerUnset should be unprinted with no name")
	}
}
