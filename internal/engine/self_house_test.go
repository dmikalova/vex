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
			RevealHand{Player: Controller, House: SelfHouse},
			Stun{Target: Target{Kind: TargetEachCreature}.
				OfHouse(SelfHouse).
				Selector(ExceptMostPowerful)},
			Exhaust{Target: Target{Kind: TargetEachCreature}.ExceptHouse(SelfHouse)},
		}}),
		WithPlayPermission(PlayPermission{House: SelfHouse, Count: 1}),
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
	out := selfHouseResolved(reflect.ValueOf(in), Dis).Interface().(selfHouseProbe)

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
