package engine

import "testing"

// TestFractionOf covers the rounding arithmetic in both directions.
func TestFractionOf(t *testing.T) {
	cases := []struct {
		f    Fraction
		n    int
		want int
	}{
		{ThirdRoundedUp, 0, 0},
		{ThirdRoundedUp, 1, 1},
		{ThirdRoundedUp, 4, 2},
		{ThirdRoundedDown, 4, 1},
		{HalfRoundedUp, 3, 2},
		{HalfRoundedDown, 3, 1},
	}
	for _, c := range cases {
		if got := c.f.of(c.n); got != c.want {
			t.Errorf("%s.of(%d) = %d, want %d", c.f.word(), c.n, got, c.want)
		}
	}
}

// TestFractionText covers valid and the rounded-up / rounding-down phrasings no
// card renders through a Loss or DestroyFraction.
func TestFractionText(t *testing.T) {
	if (Fraction{}).valid() {
		t.Error("the zero Fraction should be invalid")
	}
	if !HalfRoundedUp.valid() {
		t.Error("a built Fraction should be valid")
	}
	if got := HalfRoundedUp.object("their"); got != "half of their Æmber, rounded up" {
		t.Errorf("object = %q", got)
	}
	if got := HalfRoundedUp.countPhrase("its power"); got != "half its power, rounded up" {
		t.Errorf("countPhrase = %q", got)
	}
	if got := HalfRoundedDown.roundingPhrase(); got != "rounding down" {
		t.Errorf("roundingPhrase = %q", got)
	}
	if got := HalfRoundedDown.qualifier(); got != "" {
		t.Errorf("qualifier = %q", got)
	}
}
