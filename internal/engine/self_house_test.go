package engine

import (
	"reflect"
	"strings"
	"testing"
)

// TestResolveSelfHouseThroughDefinition checks that building a card fills the
// card's own house in for every SelfHouse sentinel, whichever part of the
// definition holds it: an effect field, a Target's house filters, a play
// permission, a house lock, and the count a key-cost change scales by.
func TestResolveSelfHouseThroughDefinition(t *testing.T) {
	def := NewCard("Probe", Mars, Creature, Common,
		WithAbility(TriggerAfterPlay, Sequence{Effects: []Effect{
			RevealHand{Player: Controller, House: namedHouse(SelfHouse)},
			Stun{Target: Target{Kind: TargetEachCreature}.
				House(namedHouse(SelfHouse)).
				Refine(Except(MostPowerful))},
			Exhaust{Target: Target{Kind: TargetEachCreature}.House(exceptHouse(SelfHouse))},
		}}),
		WithPlayPermission(PlayPermission{House: SelfHouse, Amount: 1}),
		WithHouseLock(HouseLock{Player: Controller, House: SelfHouse}),
		WithKeyCost(NewKeyCostChange(Opponent, 1).Per(InPlay{
			Player: Controller,
			Type:   Creature,
			House:  SelfHouse,
		})),
	)
	text := RenderCardText(&def)
	if strings.Contains(text, SelfHouse.String()) {
		t.Fatalf("SelfHouse survived into printed text:\n%s", text)
	}
	for _, want := range []string{
		"Mars cards from your hand",
		"each Mars creature except the most powerful",
		"each non-Mars creature",
		"you may play one Mars card",
		"must choose Mars",
		"for each friendly Mars creature",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("printed text is missing %q:\n%s", want, text)
		}
	}
}

// TestRehouseMovesEverySelfHouseReference checks that rehousing a built card — a
// Maverick or Special adopting a new pod house — moves its House and every
// reference that resolved to its printed house (target filters, play permission,
// house lock, key-cost count) to the new house, without disturbing an unrelated
// house the card names outright.
func TestRehouseMovesEverySelfHouseReference(t *testing.T) {
	def := NewCard("Probe", Mars, Creature, Common,
		WithAbility(TriggerAfterPlay, Sequence{Effects: []Effect{
			RevealHand{Player: Controller, House: namedHouse(SelfHouse)},
			Stun{Target: Target{Kind: TargetEachCreature}.House(namedHouse(SelfHouse))},
			// A house named outright must survive rehousing untouched.
			Exhaust{Target: Target{Kind: TargetEachCreature}.House(namedHouse(Brobnar))},
		}}),
		WithPlayPermission(PlayPermission{House: SelfHouse, Amount: 1}),
		WithHouseLock(HouseLock{Player: Controller, House: SelfHouse}),
	)

	def = Rehouse(def, Untamed)
	if def.House != Untamed {
		t.Fatalf("rehoused House = %v, want Untamed", def.House)
	}
	text := RenderCardText(&def)
	for _, want := range []string{
		"Untamed cards from your hand",
		"stun each Untamed creature",
		"you may play one Untamed card",
		"must choose Untamed",
		"exhaust each Brobnar creature",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("rehoused text is missing %q:\n%s", want, text)
		}
	}
	if strings.Contains(text, "Mars") {
		t.Errorf("printed Mars survived rehousing:\n%s", text)
	}
}

// TestRehouseToSameHouseIsANoop covers the guard that returns the definition
// untouched when the target house is already its printed house.
func TestRehouseToSameHouseIsANoop(t *testing.T) {
	def := NewCard("Probe", Mars, Creature, Common)
	if out := Rehouse(def, Mars); out.House != Mars {
		t.Fatalf("rehousing to the same house changed House to %v", out.House)
	}
}

// TestSelfHouseStringNamesTheCard covers the sentinel's own rendering, the only
// hint a mis-authored definition would leave.
func TestSelfHouseStringNamesTheCard(t *testing.T) {
	if got := SelfHouse.String(); got != "this card's house" {
		t.Fatalf("SelfHouse.String() = %q", got)
	}
	if got := House(200).String(); got != "Unknown" {
		t.Fatalf("House(200).String() = %q", got)
	}
}

// selfHouseProbe exercises the walker's remaining shapes — a pointer, a nil
// pointer, a nil interface, a nil slice, and a kind it leaves alone — which no
// real card definition happens to hold.
type selfHouseProbe struct {
	Ptr      *House
	NilPtr   *House
	NilFace  Effect
	NilSlice []Effect
	Map      map[string]House
	hidden   House
}

func TestSelfHouseResolvedWalksEveryShape(t *testing.T) {
	sentinel := SelfHouse
	in := selfHouseProbe{Ptr: &sentinel, Map: map[string]House{"k": SelfHouse}, hidden: SelfHouse}
	out := replaceHouse(reflect.ValueOf(in), SelfHouse, Dis).Interface().(selfHouseProbe)

	if *out.Ptr != Dis {
		t.Errorf("through pointer = %v, want Dis", *out.Ptr)
	}
	if sentinel != SelfHouse {
		t.Error("resolving rewrote the input rather than a copy")
	}
	if out.NilPtr != nil || out.NilFace != nil || out.NilSlice != nil {
		t.Error("nil pointer, interface, or slice did not survive")
	}
	if out.Map["k"] != SelfHouse || out.hidden != SelfHouse {
		t.Error("the walker reached past the shapes it handles")
	}
}
