package engine

import "testing"

// TestQuantityObject covers the noun phrase each quantity renders, including the
// defaults a verb gets when it leaves Quantity unset.
func TestQuantityObject(t *testing.T) {
	shards := CardsInPlay{Player: Controller, Trait: Shard}
	for _, tc := range []struct {
		name string
		q    Quantity
		want string
	}{
		{"unset takes one", nil, "a card"},
		{"fixed one takes one", Takes{N: Fixed(1)}, "a card"},
		{"fixed many counts the noun", Takes{N: Fixed(2)}, "2 cards"},
		{"board count reads in the lead-in", Takes{N: shards}, "a card"},
		{"up to one is just a may", UpTo{N: Fixed(1)}, "a card"},
		{"up to many caps the noun", UpTo{N: Fixed(3)}, "up to 3 cards"},
		{"up to a board count", UpTo{N: shards}, "a card"},
		{"any number pluralizes", AnyNumber{}, "any number of cards"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := quantityObject(tc.q, "card", "a card"); got != tc.want {
				t.Errorf("object = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestQuantityLeadIn covers which quantities open a sentence with a "for each"
// clause: only the ones whose number is read off the board, because a constant
// has somewhere else to print.
func TestQuantityLeadIn(t *testing.T) {
	shards := CardsInPlay{Player: Controller, Trait: Shard}
	for _, tc := range []struct {
		name string
		q    Quantity
		want Count
	}{
		{"unset", nil, nil},
		{"fixed", Takes{N: Fixed(2)}, nil},
		{"board count", Takes{N: shards}, shards},
		{"up to a board count", UpTo{N: shards}, shards},
		{"up to a constant", UpTo{N: Fixed(2)}, nil},
		{"any number", AnyNumber{}, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := quantityLeadIn(tc.q); got != tc.want {
				t.Errorf("leadIn = %v, want %v", got, tc.want)
			}
		})
	}
}

// TestQuantityPicks covers the loop bound, and the distinction the bounded flag
// exists for: a board count of zero asks nothing, while an unbounded quantity
// keeps asking until a pick comes back empty.
func TestQuantityPicks(t *testing.T) {
	g := started(t)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	shards := CardsInPlay{Player: Controller, Trait: Shard}
	for _, tc := range []struct {
		name    string
		q       Quantity
		want    int
		bounded bool
	}{
		{"unset takes one", nil, 1, true},
		{"fixed", Takes{N: Fixed(3)}, 3, true},
		{"empty board counts zero", Takes{N: shards}, 0, true},
		{"up to a ceiling", UpTo{N: Fixed(2)}, 2, true},
		{"any number is unbounded", AnyNumber{}, 0, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, bounded := quantityPicks(tc.q, ctx)
			if got != tc.want || bounded != tc.bounded {
				t.Errorf("picks = (%d, %v), want (%d, %v)", got, bounded, tc.want, tc.bounded)
			}
		})
	}
}

// TestQuantityOptional covers which quantities let the controller stop short,
// which is what makes a prompt declinable for a verb with no Selection to carry
// the choice (PutChosen).
func TestQuantityOptional(t *testing.T) {
	for _, tc := range []struct {
		name string
		q    Quantity
		want bool
	}{
		{"unset", nil, false},
		{"takes", Takes{N: Fixed(2)}, false},
		{"up to", UpTo{N: Fixed(2)}, true},
		{"any number", AnyNumber{}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := quantityOptional(tc.q); got != tc.want {
				t.Errorf("optional = %v, want %v", got, tc.want)
			}
		})
	}
}

// TestQuantitySingle covers which quantities are known before the game to take
// exactly one card — the test a verb uses to render as one clickable "you may"
// rather than its own cycle (PurgeCard.declinable).
func TestQuantitySingle(t *testing.T) {
	shards := CardsInPlay{Player: Controller, Trait: Shard}
	for _, tc := range []struct {
		name string
		q    Quantity
		want bool
	}{
		{"unset", nil, true},
		{"fixed one", Takes{N: Fixed(1)}, true},
		{"fixed many", Takes{N: Fixed(2)}, false},
		{"board count", Takes{N: shards}, false},
		{"up to one", UpTo{N: Fixed(1)}, true},
		{"up to many", UpTo{N: Fixed(2)}, false},
		{"any number", AnyNumber{}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := quantitySingle(tc.q); got != tc.want {
				t.Errorf("single = %v, want %v", got, tc.want)
			}
		})
	}
}

// TestFixedCardCount covers the numeral the web renderer badges a glyph with: a
// constant above one, and nothing for a single card or a board-scaled number.
func TestFixedCardCount(t *testing.T) {
	shards := CardsInPlay{Player: Controller, Trait: Shard}
	for _, tc := range []struct {
		name string
		q    Quantity
		want int
	}{
		{"unset is one, so no badge", nil, 0},
		{"fixed one is no badge", Takes{N: Fixed(1)}, 0},
		{"fixed many badges the count", Takes{N: Fixed(3)}, 3},
		{"up to many badges the cap", UpTo{N: Fixed(2)}, 2},
		{"board count has no numeral yet", Takes{N: shards}, 0},
		{"any number has no numeral", AnyNumber{}, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := FixedCardCount(tc.q); got != tc.want {
				t.Errorf("FixedCardCount = %d, want %d", got, tc.want)
			}
		})
	}
}

// TestQuantityValidate covers the half-written quantity a card author can no
// longer smuggle past init: a Takes or UpTo whose own number was never set reads
// as "one card" to nobody, so it is rejected rather than defaulted.
func TestQuantityValidate(t *testing.T) {
	for _, tc := range []struct {
		name    string
		q       Quantity
		wantErr bool
	}{
		{"omitted entirely is the singular default", nil, false},
		{"takes with no count", Takes{}, true},
		{"takes of none", Takes{N: Fixed(0)}, true},
		{"takes of one", Takes{N: Fixed(1)}, false},
		{"takes of a board count", Takes{N: CardsInPlay{Player: Controller, Trait: Shard}}, false},
		{"up to with no ceiling", UpTo{}, true},
		{"up to none", UpTo{N: Fixed(0)}, true},
		{"up to two", UpTo{N: Fixed(2)}, false},
		{"any number has no count to set", AnyNumber{}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := quantityValidate(tc.q); (err != nil) != tc.wantErr {
				t.Errorf("validate = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
