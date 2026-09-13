package engine

import (
	"slices"
	"testing"
)

// exRover is an Explo-rover analog: a creature that may be played as an upgrade,
// granting its host skirmish.
func exRover() CardDefinition {
	return NewCard("Rover", Brobnar, Creature, Common,
		WithPower(2),
		WithKeywords(Skirmish),
		WithStatic(StaticModifier{Keywords: []Keyword{Skirmish}}),
		WithPlayableAsUpgrade())
}

// exCalv is a CALV-1N analog: a creature whose own Fight/Reap draws, and which may
// be played as an upgrade granting its host that same Fight/Reap draw.
func exCalv() CardDefinition {
	return NewCard("CALV", Brobnar, Creature, Rare,
		WithPower(4),
		WithAbility(TriggerAfterReap, Draw{Amount: 1}),
		WithAbility(TriggerAfterFight, Draw{Amount: 1}),
		WithStatic(StaticModifier{Granted: []Ability{
			{Trigger: TriggerAfterReap, Effect: Draw{Amount: 1}},
			{Trigger: TriggerAfterFight, Effect: Draw{Amount: 1}},
		}}),
		WithPlayableAsUpgrade())
}

func TestPlayableAsUpgradeChoosesCreature(t *testing.T) {
	g := started(t)
	host := g.AddToBattleline(testCreature("host", 3), 0) // a host, so mode is a real choice
	rid := g.AddToHand(exRover(), 0)
	// The default chooser answers option 0 ("Creature").
	id, err := g.PlayCreature(0, handIdxByID(g, 0, rid), false)
	if err != nil {
		t.Fatalf("PlayCreature: %v", err)
	}
	if !slices.Contains(g.Battleline(0), id) {
		t.Error("creature mode: the card should be on the battleline")
	}
	if ups := g.Upgrades(host); len(ups) != 0 {
		t.Errorf("host should carry no upgrade in creature mode, got %v", ups)
	}
}

func TestPlayableAsUpgradeChoosesUpgrade(t *testing.T) {
	g := started(t)
	g.SetChooser(0, optionPicker{idx: 1}) // "Upgrade"
	host := g.AddToBattleline(testCreature("host", 3), 0)
	rid := g.AddToHand(exRover(), 0)
	got, err := g.PlayCreature(0, handIdxByID(g, 0, rid), false)
	if err != nil {
		t.Fatalf("PlayCreature: %v", err)
	}
	if got != host {
		t.Errorf("attached to %d, want host %d", got, host)
	}
	if !g.HasKeyword(host, Skirmish) {
		t.Error("host should gain skirmish from the attached creature-as-upgrade")
	}
	if ups := g.Upgrades(host); len(ups) != 1 || ups[0] != rid {
		t.Errorf("Upgrades(host) = %v, want [%d]", ups, rid)
	}
	if slices.Contains(g.Battleline(0), rid) {
		t.Error("a creature played as an upgrade must not be on the battleline")
	}
	// While attached it reads as an Upgrade, not a creature (ADR 0026): a
	// creature-reaching effect must treat it as the upgrade it now is.
	if got := g.TypeOf(rid); got != Upgrade {
		t.Errorf("attached creature-as-upgrade TypeOf = %v, want Upgrade", got)
	}
	if g.IsCreature(rid) {
		t.Error("an attached creature-as-upgrade must not read as a creature")
	}
	// "Destroy an upgrade" reaches it (it is now an upgrade).
	cands := Target{Kind: TargetChosenUpgrade}.
		selectBase(&EffectContext{Resolver: g, Controller: 0})
	if !slices.Contains(cands, rid) {
		t.Errorf("Destroy-upgrade candidates = %v, want to include %d", cands, rid)
	}
}

// upgradeThenDecline chooses upgrade mode but then declines every host, so the
// play finds no target.
type upgradeThenDecline struct{}

func (upgradeThenDecline) ChooseOption(_, _ string, _ []string) int { return 1 }

func (upgradeThenDecline) ChooseCreature(_, _ string, _ []LocalID) (LocalID, bool) {
	return 0, false
}

func TestPlayableAsUpgradeDeclineHostErrors(t *testing.T) {
	g := started(t)
	g.SetChooser(0, upgradeThenDecline{})
	g.AddToBattleline(testCreature("h1", 3), 0)
	g.AddToBattleline(testCreature("h2", 3), 0) // two hosts, so the pick is a real choice
	rid := g.AddToHand(exRover(), 0)
	if _, err := g.PlayCreature(0, handIdxByID(g, 0, rid), false); err != ErrNoTarget {
		t.Errorf("err = %v, want ErrNoTarget", err)
	}
}

func TestPlayableAsUpgradeNoHostPlaysAsCreature(t *testing.T) {
	g := started(t)
	g.SetChooser(0, optionPicker{idx: 1}) // would pick "Upgrade", but no host exists
	rid := g.AddToHand(exRover(), 0)
	id, err := g.PlayCreature(0, handIdxByID(g, 0, rid), false)
	if err != nil {
		t.Fatalf("PlayCreature: %v", err)
	}
	if !slices.Contains(g.Battleline(0), id) {
		t.Error("with no host, the card plays as a creature")
	}
}

func TestPlayableAsUpgradeForcedWhenCreaturesBanned(t *testing.T) {
	g := started(t)
	// The default chooser would answer "Creature", but the creature ban forces
	// upgrade mode. The Blocker is the only creature, so it is the forced host.
	host := g.AddToBattleline(NewCard("Blocker", Brobnar, Creature, Common,
		WithPower(1), WithRestrictions(Restrictions{CannotPlay: Creature})), 0)
	rid := g.AddToHand(exRover(), 0)
	got, err := g.PlayCreature(0, handIdxByID(g, 0, rid), false)
	if err != nil {
		t.Fatalf("PlayCreature: %v", err)
	}
	if got != host {
		t.Errorf("attached to %d, want %d", got, host)
	}
	if !g.HasKeyword(host, Skirmish) {
		t.Error("forced upgrade should still grant skirmish")
	}
}

func TestPlayableAsUpgradeBannedNoHostErrors(t *testing.T) {
	g := started(t)
	// Ban creatures from an artifact, so no creature is in play to host the upgrade.
	g.AddArtifact(NewCard("Ban", Brobnar, Artifact, Common,
		WithRestrictions(Restrictions{CannotPlay: Creature})), 0)
	rid := g.AddToHand(exRover(), 0)
	if _, err := g.PlayCreature(0, handIdxByID(g, 0, rid), false); err != ErrCannotPlayCreature {
		t.Errorf("err = %v, want ErrCannotPlayCreature", err)
	}
}

func TestCanPlayPlayableAsUpgradeUnderCreatureBan(t *testing.T) {
	g := started(t)
	rid := g.AddToHand(exRover(), 0)
	g.AddArtifact(NewCard("Ban", Brobnar, Artifact, Common,
		WithRestrictions(Restrictions{CannotPlay: Creature})), 0)
	if err := g.CanPlay(0, rid); err != ErrCannotPlayCreature {
		t.Errorf("CanPlay with no host = %v, want ErrCannotPlayCreature", err)
	}
	g.AddToBattleline(testCreature("host", 3), 0)
	if err := g.CanPlay(0, rid); err != nil {
		t.Errorf("CanPlay with a host = %v, want nil", err)
	}
}

func TestPlayableAsUpgradeRejectsNonCreature(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("NewCard should reject WithPlayableAsUpgrade on a non-creature")
		}
	}()
	NewCard("BadType", Brobnar, Upgrade, Common,
		WithStatic(StaticModifier{PowerBonus: 1}), WithPlayableAsUpgrade())
}

func TestPlayableAsUpgradeRejectsEmptyStatic(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("NewCard should reject a creature-as-upgrade whose Static grants nothing")
		}
	}()
	NewCard("BadStatic", Brobnar, Creature, Common,
		WithPower(2), WithPlayableAsUpgrade())
}

func TestPlayableAsUpgradeText(t *testing.T) {
	cases := []struct {
		def        CardDefinition
		rules      string
		onCreature string
	}{
		{
			exRover(),
			"Skirmish.\n" +
				`Rover may be played as an Upgrade instead of a Creature, ` +
				`with the text: "This Creature gains skirmish."`,
			"Skirmish.",
		},
		{
			exCalv(),
			"Fight/Reap: Draw a card.\n" +
				`CALV may be played as an Upgrade instead of a Creature, ` +
				`with the text: "This Creature gains, 'Fight/Reap: Draw a card.'"`,
			"Fight/Reap: Draw a card.",
		},
		{
			// A creature-as-upgrade granting non-flank fight protection.
			NewCard("Scout", Brobnar, Creature, Common,
				WithPower(2),
				WithStatic(StaticModifier{ProtectsFromNonFlank: true}),
				WithPlayableAsUpgrade()),
			`Scout may be played as an Upgrade instead of a Creature, ` +
				`with the text: "Creatures not on a flank cannot fight this Creature."`,
			"Creatures not on a flank cannot fight this Creature.",
		},
	}
	for _, tc := range cases {
		if got := RenderCardRules(&tc.def); got != tc.rules {
			t.Errorf("%s rules:\n got:  %q\n want: %q", tc.def.Name, got, tc.rules)
		}
		if got := RenderUpgradeOnCreature(&tc.def); got != tc.onCreature {
			t.Errorf("%s on creature:\n got:  %q\n want: %q", tc.def.Name, got, tc.onCreature)
		}
	}
}
