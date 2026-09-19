package engine

import (
	"strings"
	"testing"
)

func TestSequenceEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	seq := Sequence{
		Effects: []Effect{
			GainAember{Player: Controller, Amount: 1},
			GainAember{Player: Controller, Amount: 2},
		},
	}
	if got := (Sequence{}).Text(); got != "" {
		t.Errorf("empty sequence text = %q, want empty", got)
	}
	if seq.Text() != "gain 1 Æmber, and gain 2 Æmber" {
		t.Errorf("sequence text = %q", seq.Text())
	}
	seq.Resolve(ctx)
	if g.State.Aember[0] != 3 {
		t.Errorf("aember = %d, want 3", g.State.Aember[0])
	}
}

func TestSentencesRendersEachChildAsItsOwnSentence(t *testing.T) {
	seq := Sentences{Effects: []Effect{
		DiscardTop{Amount: 1, Player: Opponent},
		RevealHand{Player: Opponent},
		GainAember{
			Player: Controller,
			Amount: 1,
			Per:    CardsInHand{Player: Opponent, House: TheContextualHouse},
		},
	}}
	want := "discard the top card of your opponent's deck. Reveal your opponent's hand. For each card of the discarded card's house revealed this way, gain 1 Æmber."
	if got := seq.Text(); got != want {
		t.Errorf("sequence text = %q, want %q", got, want)
	}
	if got := (Sentences{}).Text(); got != "" {
		t.Errorf("empty text = %q, want empty", got)
	}
}

func TestSequenceCombinesSameTarget(t *testing.T) {
	// Consecutive combinable effects on the same target fold into one phrase.
	both := Sequence{Effects: []Effect{
		Stun{Target: Target{Kind: TargetThisCreature}},
		Exhaust{Target: Target{Kind: TargetThisCreature}},
	}}
	if got := both.Text(); got != "stun and exhaust "+SelfName {
		t.Errorf("combined text = %q", got)
	}

	// Different targets do not fold, and a trailing non-combinable stands alone.
	mixed := Sequence{Effects: []Effect{
		Stun{Target: Target{Kind: TargetThisCreature}},
		Exhaust{Target: Target{Kind: TargetEachEnemyCreature}},
		GainAember{Player: Controller, Amount: 1},
	}}
	want := "stun " + SelfName + ", exhaust each enemy creature, and gain 1 Æmber"
	if got := mixed.Text(); got != want {
		t.Errorf("mixed text = %q, want %q", got, want)
	}
}

func TestSequenceCombinesSameVerb(t *testing.T) {
	// Consecutive combinable effects sharing a verb fold their targets.
	both := Sequence{Effects: []Effect{
		Destroy{Target: Target{Kind: TargetChosenEnemyCreature}},
		Destroy{Target: Target{Kind: TargetChosenFriendlyCreature}},
	}}
	if got := both.Text(); got != "destroy an enemy creature and a friendly creature" {
		t.Errorf("combined text = %q", got)
	}

	// A run of three shared-verb effects folds every target.
	three := Sequence{Effects: []Effect{
		Destroy{Target: Target{Kind: TargetChosenEnemyCreature}},
		Destroy{Target: Target{Kind: TargetChosenFriendlyCreature}},
		Destroy{Target: Target{Kind: TargetEachCreature}},
	}}
	want := "destroy an enemy creature, a friendly creature, and each creature"
	if got := three.Text(); got != want {
		t.Errorf("three text = %q, want %q", got, want)
	}
}

// A run of plain type-only PutFromDiscard effects sharing a head and tail folds
// into one article-led noun list (Look What I Found!).
func TestSequenceFoldsNounList(t *testing.T) {
	seq := Sequence{Effects: []Effect{
		PutCard{
			Zones:       []Zone{Discard},
			Selection:   Chosen{Type: Tactic},
			Destination: ToHand,
		},
		PutCard{
			Zones:       []Zone{Discard},
			Selection:   Chosen{Type: Artifact},
			Destination: ToHand,
		},
		PutCard{
			Zones:       []Zone{Discard},
			Selection:   Chosen{Type: Creature},
			Destination: ToHand,
		},
		PutCard{
			Zones:       []Zone{Discard},
			Selection:   Chosen{Type: Upgrade},
			Destination: ToHand,
		},
	}}
	want := "put a tactic, artifact, creature, and upgrade " +
		"from your discard pile into your hand"
	if got := seq.Text(); got != want {
		t.Errorf("text = %q, want %q", got, want)
	}
}

// A single foldable effect stays as its own sentence, and a differing tail or a
// non-plain selection breaks the run rather than folding into it.
func TestSequenceNounListDeclines(t *testing.T) {
	// One qualifying effect on its own does not fold.
	single := Sequence{Effects: []Effect{
		PutCard{
			Zones:       []Zone{Discard},
			Selection:   Chosen{Type: Tactic},
			Destination: ToHand,
		},
	}}
	if got := single.Text(); got != "put a tactic from your discard pile into your hand" {
		t.Errorf("single text = %q", got)
	}

	// A differing tail (top of deck vs hand) stops the run after the first.
	tails := Sequence{Effects: []Effect{
		PutCard{
			Zones:       []Zone{Discard},
			Selection:   Chosen{Type: Tactic},
			Destination: ToHand,
		},
		PutCard{
			Zones:       []Zone{Discard},
			Selection:   Chosen{Type: Artifact},
			Destination: ToTopOfDeck,
		},
	}}
	want := "put a tactic from your discard pile into your hand, and " +
		"put an artifact from your discard pile on top of your deck"
	if got := tails.Text(); got != want {
		t.Errorf("tails text = %q, want %q", got, want)
	}

	// A non-plain selection (named, not a bare type) does not qualify.
	named := Sequence{Effects: []Effect{
		PutCard{
			Zones:       []Zone{Discard},
			Selection:   Named{Name: "Velum"},
			Destination: ToHand,
		},
		PutCard{
			Zones:       []Zone{Discard},
			Selection:   Chosen{Type: Artifact},
			Destination: ToHand,
		},
	}}
	if got := named.Text(); !strings.Contains(got, "Velum") {
		t.Errorf("named text = %q, want it to keep Velum unfolded", got)
	}
}

// A single exalt folds with a neighbour on the same target, but "exalt N times"
// keeps its own shape (foldable opt-out) and stands as its own clause.
func TestSequenceExaltFoldsOnlyWhenSingle(t *testing.T) {
	target := Target{Kind: TargetEachCreature}
	single := Sequence{Effects: []Effect{
		Ready{Target: target},
		Exalt{Target: target, Amount: 1},
	}}
	if got, want := single.Text(), "ready and exalt each creature"; got != want {
		t.Errorf("single exalt text = %q, want %q", got, want)
	}

	repeated := Sequence{Effects: []Effect{
		Ready{Target: target},
		Exalt{Target: target, Amount: 2},
	}}
	want := "ready each creature, and exalt each creature 2 times"
	if got := repeated.Text(); got != want {
		t.Errorf("repeated exalt text = %q, want %q", got, want)
	}
}

// A sequence that leads with a single clickable choice is declinable, so a May or
// a Repeat's MayWhile gate wrapping it is driven by that click.
func TestSequenceDeclinable(t *testing.T) {
	led := Sequence{Effects: []Effect{
		Destroy{Target: Target{Kind: TargetChosenEnemyCreature}},
		Destroy{Target: Target{Kind: TargetChosenFriendlyCreature}},
	}}
	if !led.declinable() {
		t.Error("a sequence leading with a chosen Destroy should be declinable")
	}
	if (Sequence{}).declinable() {
		t.Error("an empty sequence should not be declinable")
	}
	untargeted := Sequence{Effects: []Effect{
		Destroy{Target: Target{Kind: TargetEachCreature}},
	}}
	if untargeted.declinable() {
		t.Error("a sequence leading with an untargeted effect should not be declinable")
	}
}

// Taking the leading choice resolves the whole sequence; declining it passes on
// everything.
func TestSequenceResolveOptional(t *testing.T) {
	seq := Sequence{Effects: []Effect{
		Destroy{Target: Target{Kind: TargetChosenEnemyCreature}},
		Destroy{Target: Target{Kind: TargetChosenFriendlyCreature}},
	}}

	accepted := NewGame("A", "B", 1)
	accepted.SetChooser(0, &cardDecliner{})
	foe := accepted.AddToBattleline(testCreature("Foe", 3), 1)
	ally := accepted.AddToBattleline(testCreature("Ally", 3), 0)
	if !seq.resolveOptional(&EffectContext{Resolver: accepted, Controller: 0}) {
		t.Error("taking the leading choice should report the sequence resolved")
	}
	if onAnyLine(accepted, foe) || onAnyLine(accepted, ally) {
		t.Error("both creatures should have been destroyed")
	}

	declined := NewGame("A", "B", 1)
	declined.SetChooser(0, &cardDecliner{decline: true})
	survivor := declined.AddToBattleline(testCreature("Foe", 3), 1)
	if seq.resolveOptional(&EffectContext{Resolver: declined, Controller: 0}) {
		t.Error("declining the leading choice should report nothing resolved")
	}
	if !onAnyLine(declined, survivor) {
		t.Error("a declined sequence should destroy nothing")
	}
}

// onAnyLine reports whether a creature is still in either player's battleline.
func onAnyLine(g *Game, id LocalID) bool {
	for _, p := range []int{0, 1} {
		for _, x := range g.Battleline(p) {
			if x == id {
				return true
			}
		}
	}
	return false
}

// A sequence with no declinable lead has nothing to offer optionally.
func TestSequenceResolveOptionalWithoutAChoice(t *testing.T) {
	empty := Sequence{}
	if empty.resolveOptional(&EffectContext{}) {
		t.Error("an empty sequence should resolve nothing optionally")
	}
	untargeted := Sequence{Effects: []Effect{
		Destroy{Target: Target{Kind: TargetEachCreature}},
	}}
	if untargeted.resolveOptional(&EffectContext{}) {
		t.Error("an untargeted lead should resolve nothing optionally")
	}
}

// Sentences carries the same leading-choice rule as Sequence: "you may destroy a
// creature. Gain 1 Æmber" is answered by clicking the creature, and declining
// passes on the later sentences too.
func TestSentencesDeclinable(t *testing.T) {
	led := Sentences{Effects: []Effect{
		Destroy{Target: Target{Kind: TargetChosenEnemyCreature}},
		GainAember{Player: Controller, Amount: 1},
	}}
	if !led.declinable() {
		t.Error("sentences leading with a chosen Destroy should be declinable")
	}
	if (Sentences{}).declinable() {
		t.Error("empty sentences should not be declinable")
	}

	accepted := NewGame("A", "B", 1)
	accepted.SetChooser(0, &cardDecliner{})
	foe := accepted.AddToBattleline(testCreature("Foe", 3), 1)
	if !led.resolveOptional(&EffectContext{Resolver: accepted, Controller: 0}) {
		t.Error("taking the leading choice should report the sentences resolved")
	}
	if onAnyLine(accepted, foe) || accepted.Aember(0) != 1 {
		t.Error("both sentences should have resolved")
	}

	declined := NewGame("A", "B", 1)
	declined.SetChooser(0, &cardDecliner{decline: true})
	survivor := declined.AddToBattleline(testCreature("Foe", 3), 1)
	if led.resolveOptional(&EffectContext{Resolver: declined, Controller: 0}) {
		t.Error("declining the leading choice should report nothing resolved")
	}
	if !onAnyLine(declined, survivor) || declined.Aember(0) != 0 {
		t.Error("a declined lead should pass on the later sentences too")
	}
}

// housesRung builds one rung of Galactic Census's ladder: gain 1 Æmber if at
// least n houses are represented among creatures in play.
func housesRung(n int) Conditional {
	return Conditional{
		Cond: HousesRepresented{
			Among:  HousesAmong{Player: EachPlayer, Type: Creature},
			Is:     AtLeast,
			Amount: n,
		},
		Then: GainAember{Player: Controller, Amount: 1},
	}
}

// TestSentencesFoldsThresholdLadder pins Galactic Census's wording: a run of
// rungs over the same board survey names that survey once, and every later rung
// is a bare "gain 1 more if there are N or more".
func TestSentencesFoldsThresholdLadder(t *testing.T) {
	ladder := Sentences{Effects: []Effect{housesRung(3), housesRung(5), housesRung(6)}}
	want := "if there are 3 or more houses represented among creatures in play, gain 1 Æmber. " +
		"Gain 1 more if there are 5 or more. Gain 1 more if there are 6 or more."
	if got := ladder.Text(); got != want {
		t.Errorf("ladder text = %q, want %q", got, want)
	}
}

// TestSentencesLadderDeclines walks every reason a Conditional is not a ladder
// rung, so a card that merely looks like one keeps its full text.
func TestSentencesLadderDeclines(t *testing.T) {
	full := "if there are 3 or more houses represented among creatures in play, gain 1 Æmber."
	otherSubject := housesRung(5)
	otherSubject.Cond = HousesRepresented{
		Among:  HousesAmong{Player: Controller, Type: Creature},
		Is:     AtLeast,
		Amount: 5,
	}
	withElse := housesRung(5)
	withElse.Else = GainAember{Player: Controller, Amount: 1}
	perRung := housesRung(5)
	perRung.Then = GainAember{Player: Controller, Amount: 1, Per: Fixed(2)}
	notRepeating := housesRung(5)
	notRepeating.Then = Draw{Amount: 1}

	cases := []struct {
		name    string
		second  Effect
		wantFix string
	}{
		{"a lone rung", nil, full},
		{"a different survey", otherSubject, "among friendly creatures"},
		{"a two-way branch", withElse, "Otherwise"},
		{"a scaled gain", perRung, "for each"},
		{"an effect with no repeated form", notRepeating, "draw"},
		{
			"an effect that is no Conditional at all",
			GainAember{Player: Controller, Amount: 2},
			"Gain 2 Æmber",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			effects := []Effect{housesRung(3)}
			if tc.second != nil {
				effects = append(effects, tc.second)
			}
			got := Sentences{Effects: effects}.Text()
			if !strings.HasPrefix(got, full) {
				t.Errorf("text = %q, want it to open with the unfolded rung", got)
			}
			if !strings.Contains(got, tc.wantFix) {
				t.Errorf("text = %q, want it to contain %q", got, tc.wantFix)
			}
			if strings.Contains(got, "more if there are") {
				t.Errorf("text = %q, want no ladder fold", got)
			}
		})
	}
}

// TestGainAemberRepeatedText covers the forms a gain can and cannot repeat as.
func TestGainAemberRepeatedText(t *testing.T) {
	cases := []struct {
		name string
		gain GainAember
		want string
	}{
		{"controller", GainAember{Player: Controller, Amount: 1}, "gain 1 more"},
		{"opponent", GainAember{Player: Opponent, Amount: 2}, "your opponent gains 2 more"},
		{"each player", GainAember{Player: EachPlayer, Amount: 1}, "each player gains 1 more"},
		{"equal to a count", GainAember{Player: Controller, EqualTo: Fixed(1)}, ""},
		{"scaled by a count", GainAember{Player: Controller, Amount: 1, Per: Fixed(2)}, ""},
		{"a subject gainVerb cannot name", GainAember{Player: ItsOwner, Amount: 1}, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.gain.repeatedText(); got != tc.want {
				t.Errorf("repeatedText() = %q, want %q", got, tc.want)
			}
		})
	}
}
