package engine

import (
	"slices"
	"testing"
)

func TestCanPlay(t *testing.T) {
	g := started(t) // player 0 active, Brobnar
	creat := g.AddToHand(testCreature("c", 3), 0)
	if err := g.CanPlay(0, creat); err != nil {
		t.Errorf("playable creature = %v, want nil", err)
	}
	if err := g.CanPlay(1, creat); err != ErrNotActivePlayer {
		t.Errorf("wrong player = %v, want ErrNotActivePlayer", err)
	}

	off := g.AddToHand(NewCard("off", Sanctum, Creature, Common, WithPower(1)), 0)
	if err := g.CanPlay(0, off); err != ErrWrongHouse {
		t.Errorf("off-house = %v, want ErrWrongHouse", err)
	}

	up := g.AddToHand(NewCard("up", Brobnar, Upgrade, Common), 0)
	if err := g.CanPlay(0, up); err != ErrNoTarget {
		t.Errorf("hostless upgrade = %v, want ErrNoTarget", err)
	}

	g.State.Winner = 1
	if err := g.CanPlay(0, creat); err != ErrGameOver {
		t.Errorf("game over = %v, want ErrGameOver", err)
	}
}

func TestCanPlayRestrictions(t *testing.T) {
	// Creatures barred: a card that stops the player playing creatures.
	g := started(t)
	g.AddToBattleline(NewCard("Blocker", Brobnar, Creature, Common, WithPower(1),
		WithRestrictions(Restrictions{CannotPlay: Creature})), 0)
	c := g.AddToHand(testCreature("c", 3), 0)
	if err := g.CanPlay(0, c); err != ErrCannotPlayCreature {
		t.Errorf("creatures barred = %v, want ErrCannotPlayCreature", err)
	}

	// Card-play limit reached this turn.
	g2 := started(t)
	g2.AddToBattleline(NewCard(
		"Imp",
		Brobnar,
		Creature,
		Common,
		WithPower(1),
		WithRestrictions(
			Restrictions{PlayCardLimit: PlayCardLimit{Player: Controller, Amount: 2}},
		),
	), 0)
	g2.State.PlayedThisTurn[0].Count = 2
	c2 := g2.AddToHand(testCreature("c2", 3), 0)
	if err := g2.CanPlay(0, c2); err != ErrCardPlayLimit {
		t.Errorf("limit reached = %v, want ErrCardPlayLimit", err)
	}
}

func TestCardPlayLimit(t *testing.T) {
	g := started(t) // player 0 active, Brobnar
	// Player 1 controls an Ember-Imp-like card limiting player 0 to two plays.
	g.AddToBattleline(NewCard(
		"imp",
		Dis,
		Creature,
		Common,
		WithPower(2),
		WithRestrictions(
			Restrictions{PlayCardLimit: PlayCardLimit{Player: Opponent, Amount: 2}},
		),
	), 1)
	g.AddToHand(testCreature("c0", 3), 0)
	g.AddToHand(testCreature("c1", 3), 0)
	g.AddToHand(testCreature("c2", 3), 0)
	g.AddToHand(exBruteStrength(), 0) // an upgrade

	if _, err := g.PlayCreature(0, handIdx(g, 0, "c0"), false); err != nil {
		t.Fatalf("first play: %v", err)
	}
	if _, err := g.PlayCreature(0, handIdx(g, 0, "c1"), false); err != nil {
		t.Fatalf("second play: %v", err)
	}
	// The third card play is barred, whatever its type.
	if _, err := g.PlayCreature(0, handIdx(g, 0, "c2"), false); err != ErrCardPlayLimit {
		t.Errorf("third creature play = %v, want ErrCardPlayLimit", err)
	}
	if _, err := g.PlayUpgrade(0, handIdx(g, 0, "Brute Strength")); err != ErrCardPlayLimit {
		t.Errorf("upgrade play = %v, want ErrCardPlayLimit", err)
	}

	// A new turn resets the count.
	g.EndPlayPhase(0)
	g.StartTurn(0)
	if err := g.ChooseHouse(0, Brobnar); err != nil {
		t.Fatal(err)
	}
	if _, err := g.PlayCreature(0, handIdx(g, 0, "c2"), false); err != nil {
		t.Errorf("play after new turn = %v, want nil (count reset)", err)
	}
}

func TestPlayedThisTurn(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	g.StartTurn(0)
	creature := g.AddToHand(NewCard("brobnar creature", Brobnar, Creature, Common, WithPower(3)), 0)
	artifact := g.AddToHand(NewCard("sanctum artifact", Sanctum, Artifact, Common), 0)
	action := g.AddToHand(NewCard("sanctum action", Sanctum, Tactic, Common), 0)
	upgrade := g.AddToHand(NewCard("mars upgrade", Mars, Upgrade, Common), 0)

	if _, err := g.PlayCreature(0, handIdxByID(g, 0, creature), false); err != nil {
		t.Fatalf("PlayCreature: %v", err)
	}
	if _, err := g.PlayArtifact(0, handIdxByID(g, 0, artifact)); err != nil {
		t.Fatalf("PlayArtifact: %v", err)
	}
	if err := g.PlayAction(0, handIdxByID(g, 0, action)); err != nil {
		t.Fatalf("PlayAction: %v", err)
	}
	if _, err := g.PlayUpgrade(0, handIdxByID(g, 0, upgrade)); err != nil {
		t.Fatalf("PlayUpgrade: %v", err)
	}

	playsOf := func(house House) int {
		n := 0
		for _, id := range g.PlayedThisTurn(0) {
			if g.House(id) == house {
				n++
			}
		}
		return n
	}
	if got := playsOf(Brobnar); got != 1 {
		t.Errorf("Brobnar plays = %d, want 1", got)
	}
	if got := playsOf(Sanctum); got != 2 {
		t.Errorf("Sanctum plays = %d, want 2", got)
	}
	if got := playsOf(Mars); got != 1 {
		t.Errorf("Mars plays = %d, want 1", got)
	}
	if got := len(g.PlayedThisTurn(0)); got != 4 {
		t.Errorf("total plays = %d, want 4", got)
	}

	g.StartTurn(0)
	if got := playsOf(Sanctum); got != 0 {
		t.Errorf("Sanctum plays after reset = %d, want 0", got)
	}
	if got := len(g.PlayedThisTurn(0)); got != 0 {
		t.Errorf("total plays after reset = %d, want 0", got)
	}
}

// TestTypeUnlimitedPlayPermission covers Matter Maker: a house-agnostic, unlimited
// waiver for a card type lets its controller play any number of off-house upgrades
// without spending any per-turn play grant, but frees no other card type.
func TestTypeUnlimitedPlayPermission(t *testing.T) {
	maker := NewCard("Matter Maker", StarAlliance, Artifact, Rare,
		WithPlayPermission(PlayPermission{Types: CardTypesOf(Upgrade)}))

	t.Run("frees any number of off-house upgrades without a counter", func(t *testing.T) {
		g := started(t) // Brobnar active
		g.AddArtifact(maker, 0)
		g.AddToBattleline(testCreature("host", 3), 0)
		up1 := g.AddToHand(
			NewCard("Bolt", Logos, Upgrade, Common, WithStatic(StaticModifier{PowerBonus: 1})), 0)
		up2 := g.AddToHand(
			NewCard("Coil", Logos, Upgrade, Common, WithStatic(StaticModifier{PowerBonus: 1})), 0)

		if err := g.CanPlay(0, up1); err != nil {
			t.Fatalf("CanPlay off-house upgrade with Matter Maker = %v, want nil", err)
		}
		if _, err := g.PlayUpgrade(0, handIdxByID(g, 0, up1)); err != nil {
			t.Fatalf("first upgrade: %v", err)
		}
		if _, err := g.PlayUpgrade(0, handIdxByID(g, 0, up2)); err != nil {
			t.Fatalf("second upgrade: %v", err)
		}
		if got := g.State.NonActivePlaysUsedThisTurn[0]; got != 0 {
			t.Errorf("type waiver consumed a non-active counter: %d", got)
		}
	})

	t.Run("frees no other card type", func(t *testing.T) {
		g := started(t)
		g.AddArtifact(maker, 0)
		creature := g.AddToHand(NewCard("Logos Bot", Logos, Creature, Common, WithPower(3)), 0)
		if err := g.CanPlay(0, creature); err != ErrWrongHouse {
			t.Fatalf("CanPlay off-house creature = %v, want ErrWrongHouse", err)
		}
	})
}

func TestOffHousePlayGrant(t *testing.T) {
	witch := NewCard(
		"Witch",
		Untamed,
		Creature,
		Rare,
		WithPower(4),
		WithPlayPermission(PlayPermission{House: Untamed, Amount: 1}),
	)

	t.Run("allows one off-house play and consumes it", func(t *testing.T) {
		g := started(t)
		g.AddToBattleline(witch, 0)
		first := g.AddToHand(NewCard("bear", Untamed, Creature, Common, WithPower(3)), 0)
		second := g.AddToHand(NewCard("wolf", Untamed, Creature, Common, WithPower(3)), 0)

		if err := g.CanPlay(0, first); err != nil {
			t.Fatalf("CanPlay first off-house Untamed = %v, want nil", err)
		}
		if got := g.State.PlayPermissionsUsedThisTurn[0][Untamed]; got != 0 {
			t.Fatalf("CanPlay consumed off-house grant: used = %d, want 0", got)
		}
		if _, err := g.PlayCreature(0, handIdxByID(g, 0, first), false); err != nil {
			t.Fatalf("first off-house play: %v", err)
		}
		if got := g.State.PlayPermissionsUsedThisTurn[0][Untamed]; got != 1 {
			t.Fatalf("off-house plays used = %d, want 1", got)
		}
		if err := g.CanPlay(0, second); err != ErrWrongHouse {
			t.Fatalf("CanPlay second off-house Untamed = %v, want ErrWrongHouse", err)
		}
		if _, err := g.PlayCreature(0, handIdxByID(g, 0, second), false); err != ErrWrongHouse {
			t.Fatalf("second off-house play = %v, want ErrWrongHouse", err)
		}

		g.EndPlayPhase(0)
		g.StartTurn(0)
		if err := g.ChooseHouse(0, Brobnar); err != nil {
			t.Fatal(err)
		}
		if _, err := g.PlayCreature(0, handIdxByID(g, 0, second), false); err != nil {
			t.Fatalf("off-house play after reset: %v", err)
		}
	})

	t.Run("does not matter when Untamed is active", func(t *testing.T) {
		g := NewGame("Alice", "Bob", 1)
		g.StartTurn(0)
		if err := g.ChooseHouse(0, Untamed); err != nil {
			t.Fatal(err)
		}
		g.AddToBattleline(witch, 0)
		first := g.AddToHand(NewCard("boar", Untamed, Creature, Common, WithPower(3)), 0)
		second := g.AddToHand(NewCard("fox", Untamed, Creature, Common, WithPower(3)), 0)

		if _, err := g.PlayCreature(0, handIdxByID(g, 0, first), false); err != nil {
			t.Fatalf("first active-house play: %v", err)
		}
		if _, err := g.PlayCreature(0, handIdxByID(g, 0, second), false); err != nil {
			t.Fatalf("second active-house play: %v", err)
		}
		if got := g.State.PlayPermissionsUsedThisTurn[0][Untamed]; got != 0 {
			t.Fatalf("active-house plays used off-house grant: got %d, want 0", got)
		}
	})

	t.Run("rejects off-house play without a grant", func(t *testing.T) {
		g := started(t)
		untamed := g.AddToHand(NewCard("badger", Untamed, Creature, Common, WithPower(3)), 0)

		if err := g.CanPlay(0, untamed); err != ErrWrongHouse {
			t.Fatalf("CanPlay without Witch = %v, want ErrWrongHouse", err)
		}
		if _, err := g.PlayCreature(0, handIdxByID(g, 0, untamed), false); err != ErrWrongHouse {
			t.Fatalf("PlayCreature without Witch = %v, want ErrWrongHouse", err)
		}
	})

	t.Run("counts each controlled grant", func(t *testing.T) {
		g := started(t)
		g.AddToBattleline(witch, 0)
		g.AddToBattleline(witch, 0)
		first := g.AddToHand(NewCard("hare", Untamed, Creature, Common, WithPower(3)), 0)
		second := g.AddToHand(NewCard("lynx", Untamed, Creature, Common, WithPower(3)), 0)
		third := g.AddToHand(NewCard("mole", Untamed, Creature, Common, WithPower(3)), 0)

		if _, err := g.PlayCreature(0, handIdxByID(g, 0, first), false); err != nil {
			t.Fatalf("first off-house play with two grants: %v", err)
		}
		if _, err := g.PlayCreature(0, handIdxByID(g, 0, second), false); err != nil {
			t.Fatalf("second off-house play with two grants: %v", err)
		}
		if _, err := g.PlayCreature(0, handIdxByID(g, 0, third), false); err != ErrWrongHouse {
			t.Fatalf("third off-house play with two grants = %v, want ErrWrongHouse", err)
		}
	})
}

func TestPlayCardLimitTargets(t *testing.T) {
	cases := []struct {
		name         string
		player       Player
		limited      int
		unrestricted int
	}{
		{name: "controller", player: Controller, limited: 0, unrestricted: 1},
		{name: "opponent", player: Opponent, limited: 1, unrestricted: 0},
		{name: "each player", player: EachPlayer, limited: 0, unrestricted: -1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			g := started(t)
			g.AddToBattleline(NewCard(
				"limit",
				Brobnar,
				Creature,
				Common,
				WithPower(1),
				WithRestrictions(
					Restrictions{PlayCardLimit: PlayCardLimit{Player: tc.player, Amount: 2}},
				),
			), 0)
			g.State.PlayedThisTurn[tc.limited].Count = 2
			if !g.cannotPlayCard(tc.limited) {
				t.Errorf("player %d should be limited", tc.limited)
			}
			if tc.unrestricted >= 0 && g.cannotPlayCard(tc.unrestricted) {
				t.Errorf("player %d should not be limited", tc.unrestricted)
			}
		})
	}
	if (PlayCardLimit{}).affects(0, 0) {
		t.Error("an unset play-card limit should affect no player")
	}
}

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
		t.Fatalf("Playcreature: %v", err)
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
		t.Fatalf("Playcreature: %v", err)
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
		t.Errorf("attached creature-as-upgrade TypeOf = %v, want upgrade", got)
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
		t.Fatalf("Playcreature: %v", err)
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
		t.Fatalf("Playcreature: %v", err)
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
		t.Errorf("err = %v, want ErrCannotPlaycreature", err)
	}
}

func TestCanPlayPlayableAsUpgradeUnderCreatureBan(t *testing.T) {
	g := started(t)
	rid := g.AddToHand(exRover(), 0)
	g.AddArtifact(NewCard("Ban", Brobnar, Artifact, Common,
		WithRestrictions(Restrictions{CannotPlay: Creature})), 0)
	if err := g.CanPlay(0, rid); err != ErrCannotPlayCreature {
		t.Errorf("CanPlay with no host = %v, want ErrCannotPlaycreature", err)
	}
	g.AddToBattleline(testCreature("host", 3), 0)
	if err := g.CanPlay(0, rid); err != nil {
		t.Errorf("CanPlay with a host = %v, want nil", err)
	}
}

func TestPlayableAsUpgradeRejectsNonCreature(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("NewCard should reject WithPlayableAsupgrade on a non-creature")
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
				`Rover may be played as an upgrade instead of a creature, ` +
				`with the text: "This creature gains skirmish."`,
			"Skirmish.",
		},
		{
			exCalv(),
			"Fight/Reap: Draw a card.\n" +
				`CALV may be played as an upgrade instead of a creature, ` +
				`with the text: "This creature gains, 'Fight/Reap: Draw a card.'"`,
			"Fight/Reap: Draw a card.",
		},
		{
			// A creature-as-upgrade granting non-flank fight protection.
			NewCard("Scout", Brobnar, Creature, Common,
				WithPower(2),
				WithStatic(StaticModifier{ProtectsFromNonFlank: true}),
				WithPlayableAsUpgrade()),
			`Scout may be played as an upgrade instead of a creature, ` +
				`with the text: "Creatures not on a flank cannot fight this creature."`,
			"Creatures not on a flank cannot fight this creature.",
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
