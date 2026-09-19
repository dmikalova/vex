package engine

import "testing"

// namedHouse, exceptHouse, chosenHouse, activeHouse, and contextualHouse build
// the HouseMatcher shapes the effect tests name, mirroring the card.Houses facade
// the cards use.
func namedHouse(h House) HouseMatcher {
	return HouseMatcher{
		Kind:  MatchNamedHouse,
		House: h,
	}
}
func exceptHouse(h House) HouseMatcher {
	return HouseMatcher{
		Kind:  MatchExceptHouse,
		House: h,
	}
}

var (
	anyHouse        = HouseMatcher{Kind: MatchAnyHouse}
	chosenHouse     = HouseMatcher{Kind: MatchChosenHouse}
	activeHouse     = HouseMatcher{Kind: MatchActiveHouse}
	contextualHouse = HouseMatcher{Kind: MatchContextualHouse}
)

func TestHouseMatcherMatches(t *testing.T) {
	g := started(t) // Brobnar is the active house.
	mars := g.AddToHand(NewCard("Mars Card", Mars, Artifact, Common), 0)
	logos := g.AddToHand(NewCard("Logos Card", Logos, Artifact, Common), 0)
	brobnar := g.AddToHand(NewCard("Brobnar Card", Brobnar, Artifact, Common), 0)

	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}

	cases := []struct {
		name    string
		matcher HouseMatcher
		id      LocalID
		want    bool
	}{
		{"any admits every house", HouseMatcher{Kind: MatchAnyHouse}, mars, true},
		{"named admits its house", HouseMatcher{
			Kind:  MatchNamedHouse,
			House: Mars,
		}, mars, true},
		{"named bars other houses", HouseMatcher{
			Kind:  MatchNamedHouse,
			House: Mars,
		}, logos, false},
		{"except bars its house", HouseMatcher{
			Kind:  MatchExceptHouse,
			House: Mars,
		}, mars, false},
		{
			"except admits other houses",
			HouseMatcher{
				Kind:  MatchExceptHouse,
				House: Mars,
			},
			logos,
			true,
		},
		{"active admits the active house", HouseMatcher{Kind: MatchActiveHouse}, brobnar, true},
		{"active bars off-house", HouseMatcher{Kind: MatchActiveHouse}, mars, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.matcher.matches(ctx, tc.id); got != tc.want {
				t.Errorf("matches = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestHouseMatcherMatchesChosenHouse(t *testing.T) {
	g := started(t)
	mars := g.AddToHand(NewCard("Mars Card", Mars, Artifact, Common), 0)
	logos := g.AddToHand(NewCard("Logos Card", Logos, Artifact, Common), 0)

	ctx := &EffectContext{
		Resolver:    g,
		Controller:  0,
		ChosenHouse: Mars,
	}
	m := HouseMatcher{Kind: MatchChosenHouse}
	if !m.matches(ctx, mars) {
		t.Error("chosen-house matcher should admit a card of the chosen house")
	}
	if m.matches(ctx, logos) {
		t.Error("chosen-house matcher should bar a card of another house")
	}
}

func TestHouseMatcherMatchesContextualHouse(t *testing.T) {
	g := started(t)
	focus := g.AddToHand(NewCard("Focus", Mars, Artifact, Common), 0)
	sameHouse := g.AddToHand(NewCard("Same House", Mars, Artifact, Common), 0)
	otherHouse := g.AddToHand(NewCard("Other House", Logos, Artifact, Common), 0)

	m := HouseMatcher{Kind: MatchContextualHouse}

	withIt := &EffectContext{
		Resolver:   g,
		Controller: 0,
		It:         focus,
		HasIt:      true,
	}
	if !m.matches(withIt, sameHouse) {
		t.Error("contextual matcher should admit a card sharing the focused card's house")
	}
	if m.matches(withIt, otherHouse) {
		t.Error("contextual matcher should bar a card of another house")
	}

	noIt := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}
	if m.matches(noIt, sameHouse) {
		t.Error("contextual matcher should admit nothing with no card in context")
	}
}

func TestHouseMatcherQualify(t *testing.T) {
	cases := []struct {
		name    string
		matcher HouseMatcher
		want    string
	}{
		{"any", HouseMatcher{Kind: MatchAnyHouse}, "card"},
		{"named", HouseMatcher{
			Kind:  MatchNamedHouse,
			House: Mars,
		}, "Mars card"},
		{"except", HouseMatcher{
			Kind:  MatchExceptHouse,
			House: Logos,
		}, "non-Logos card"},
		{"chosen", HouseMatcher{Kind: MatchChosenHouse}, "card of the chosen house"},
		{"active", HouseMatcher{Kind: MatchActiveHouse}, "card of that house"},
		{"contextual", HouseMatcher{Kind: MatchContextualHouse}, "card of that card's house"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.matcher.qualify("card"); got != tc.want {
				t.Errorf("qualify = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestHouseMatcherFilters(t *testing.T) {
	if (HouseMatcher{Kind: MatchAnyHouse}).filters() {
		t.Error("the any-house matcher should filter nothing")
	}
	if !(HouseMatcher{
		Kind:  MatchNamedHouse,
		House: Mars,
	}).filters() {
		t.Error("a named-house matcher should filter")
	}
}

func TestHouseMatcherValidate(t *testing.T) {
	if err := (HouseMatcher{Kind: MatchAnyHouse}).validate(); err != nil {
		t.Errorf("the any-house matcher should validate: %v", err)
	}
	if err := (HouseMatcher{
		Kind:  MatchNamedHouse,
		House: Mars,
	}).validate(); err != nil {
		t.Errorf("a named-house matcher with a house should validate: %v", err)
	}
	if err := (HouseMatcher{Kind: MatchNamedHouse}).validate(); err == nil {
		t.Error("a named-house matcher without a house should not validate")
	}
	if err := (HouseMatcher{Kind: MatchExceptHouse}).validate(); err == nil {
		t.Error("an except-house matcher without a house should not validate")
	}
}
