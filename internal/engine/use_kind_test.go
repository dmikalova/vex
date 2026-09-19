package engine

import (
	"errors"
	"slices"
	"testing"
)

// A card barred from one way of being used stays open to the others: it cannot
// reap, but it fights, and CanUse still offers it while a fight is available.
func TestCannotBeUsedToReap(t *testing.T) {
	g := started(t)
	crocag := g.AddToBattleline(
		testCreature("Crocag", 7, WithCannotBeUsedTo(ReapUse)), 0)
	foe := g.AddToBattleline(testCreature("foe", 3), 1)

	if err := g.CanUse(0, crocag); err != nil {
		t.Fatalf("CanUse = %v, want nil while a fight is available", err)
	}
	if err := g.CanUseTo(0, crocag, ReapUse); !errors.Is(err, ErrCannotUse) {
		t.Errorf("CanUseTo(reap) = %v, want ErrCannotUse", err)
	}
	if err := g.CanUseTo(0, crocag, FightUse); err != nil {
		t.Errorf("CanUseTo(fight) = %v, want nil", err)
	}
	if err := g.Reap(0, crocag); !errors.Is(err, ErrCannotUse) {
		t.Errorf("Reap = %v, want ErrCannotUse", err)
	}
	if g.Exhausted(crocag) {
		t.Error("a refused reap must not exhaust the creature")
	}
	if err := g.Fight(0, crocag, foe); err != nil {
		t.Fatalf("Fight: %v", err)
	}
}

// CannotBeUsedTo is the restriction-only check a client uses to omit an illegal
// verb from the buttons it offers: it reports the bar without CanUseTo's house or
// exhaustion gate, so it stays true for a creature barred by its own text and
// false for a verb that is open.
func TestCannotBeUsedToRestrictionOnly(t *testing.T) {
	g := started(t)
	crocag := g.AddToBattleline(
		testCreature("Crocag", 7, WithCannotBeUsedTo(ReapUse)), 0)

	if !g.CannotBeUsedTo(crocag, ReapUse) {
		t.Error("CannotBeUsedTo(reap) = false, want true")
	}
	if g.CannotBeUsedTo(crocag, FightUse) {
		t.Error("CannotBeUsedTo(fight) = true, want false")
	}
}

// An upgrade whose static modifier bars a use kind bars its host from that use —
// Access Denied ("This creature cannot reap"), Detention Coil ("cannot fight") —
// while leaving the other ways of using the host open.
func TestCannotBeUsedToFromUpgradeStatic(t *testing.T) {
	g := started(t)
	host := g.AddToBattleline(testCreature("host", 7), 0)
	attachUpgrade(g, host, NewCard("Detention Coil", StarAlliance, Upgrade, Common,
		WithStatic(StaticModifier{CannotBeUsedTo: []UseKind{FightUse}})))

	if !g.CannotBeUsedTo(host, FightUse) {
		t.Error("host with a cannot-fight upgrade should be barred from fighting")
	}
	if g.CannotBeUsedTo(host, ReapUse) {
		t.Error("host should still be able to reap")
	}
}

// The upgrade-granted cannot-use restriction renders in the upgrade's own voice,
// both on its card face and while attached to a host.
func TestUpgradeCannotBeUsedToText(t *testing.T) {
	def := NewCard("Access Denied", StarAlliance, Upgrade, Common,
		WithStatic(StaticModifier{CannotBeUsedTo: []UseKind{ReapUse}}))
	if got := staticText(def.Static); got != "This creature cannot reap." {
		t.Errorf("staticText = %q", got)
	}
	lines := upgradeStaticLines(&def, true)
	if len(lines) != 1 || lines[0] != "This creature cannot reap." {
		t.Errorf("hosted lines = %v, want [This creature cannot reap.]", lines)
	}
}

// A house-scoped reap bar (Seismo-entangler) refuses the reap of a creature of
// that house while leaving it free to fight.
func TestCannotReapHouseLeavesFightUsable(t *testing.T) {
	g := started(t)
	g.State.CannotReapHouse[0].Value = Brobnar
	reaper := g.AddToBattleline(testCreature("reaper", 4), 0)
	g.AddToBattleline(testCreature("foe", 3), 1)

	if err := g.CanUseTo(0, reaper, ReapUse); !errors.Is(err, ErrCannotUse) {
		t.Errorf("CanUseTo(reap) barred by house = %v, want ErrCannotUse", err)
	}
	if err := g.CanUseTo(0, reaper, FightUse); err != nil {
		t.Errorf("CanUseTo(fight) = %v, want nil", err)
	}
}

// A card barred from fighting or from its Action: ability is refused those uses,
// and an effect that reaps with a card which cannot reap does nothing.
func TestCannotBeUsedToFightAndAction(t *testing.T) {
	g := started(t)
	pacifist := g.AddToBattleline(
		testCreature("pacifist", 4,
			WithCannotBeUsedTo(FightUse),
			WithAbility(TriggerAction, GainAember{
				Player: Controller,
				Amount: 1,
			}),
		), 0)
	foe := g.AddToBattleline(testCreature("foe", 3), 1)
	if err := g.Fight(0, pacifist, foe); !errors.Is(err, ErrCannotUse) {
		t.Errorf("Fight = %v, want ErrCannotUse", err)
	}

	idle := g.AddToBattleline(
		testCreature("idle", 4,
			WithCannotBeUsedTo(ActionUse),
			WithAbility(TriggerAction, GainAember{
				Player: Controller,
				Amount: 1,
			}),
		), 0)
	if err := g.UseAction(0, idle); !errors.Is(err, ErrCannotUse) {
		t.Errorf("UseAction = %v, want ErrCannotUse", err)
	}

	before := g.Aember(0)
	quiet := g.AddToBattleline(
		testCreature("quiet", 4, WithCannotBeUsedTo(ReapUse)), 0)
	g.reapWith(quiet)
	if g.Aember(0) != before {
		t.Errorf("aember = %d, want %d (a barred reap gains nothing)", g.Aember(0), before)
	}
}

// A creature whose only remaining way of being used has no target is not usable
// at all, so CanUse never promises a use the player cannot make.
func TestHasAnyUse(t *testing.T) {
	g := started(t)
	crocag := g.AddToBattleline(
		testCreature("Crocag", 7, WithCannotBeUsedTo(ReapUse)), 0)
	if err := g.CanUse(0, crocag); !errors.Is(err, ErrCannotUse) {
		t.Errorf("CanUse with an empty enemy battleline = %v, want ErrCannotUse", err)
	}

	// An Action: ability is a use of its own, so it keeps the card usable.
	actor := g.AddToBattleline(
		testCreature("actor", 7,
			WithCannotBeUsedTo(ReapUse),
			WithAbility(TriggerAction, GainAember{
				Player: Controller,
				Amount: 1,
			}),
		), 0)
	if err := g.CanUse(0, actor); err != nil {
		t.Errorf("CanUse with an Action ability = %v, want nil", err)
	}

	// A player-wide fight ban closes the last door on a creature that cannot reap.
	g.AddToBattleline(testCreature("foe", 3), 1)
	g.State.CannotFight[0].Value = true
	if err := g.CanUse(0, crocag); !errors.Is(err, ErrCannotUse) {
		t.Errorf("CanUse while fighting is banned = %v, want ErrCannotUse", err)
	}
}

// A card barred while a condition holds (Valoocanth while the tide is low) is open
// every way while the condition is false and closed every way once it holds.
func TestCannotBeUsedWhile(t *testing.T) {
	g := started(t)
	valoo := g.AddToBattleline(
		testCreature("Valoo", 6, WithCannotBeUsedWhile(TideIsLow{})), 0)
	g.AddToBattleline(testCreature("foe", 3), 1)

	// Neutral tide: usable every way.
	if err := g.CanUseTo(0, valoo, ReapUse); err != nil {
		t.Errorf("CanUseTo(reap) while the tide is neutral = %v, want nil", err)
	}
	if err := g.CanUseTo(0, valoo, FightUse); err != nil {
		t.Errorf("CanUseTo(fight) while the tide is neutral = %v, want nil", err)
	}

	// Tide high for player 1 is low for player 0: barred every way.
	g.State.Tide = TideHighForP1
	if err := g.CanUseTo(0, valoo, ReapUse); !errors.Is(err, ErrCannotUse) {
		t.Errorf("CanUseTo(reap) while the tide is low = %v, want ErrCannotUse", err)
	}
	if err := g.CanUseTo(0, valoo, FightUse); !errors.Is(err, ErrCannotUse) {
		t.Errorf("CanUseTo(fight) while the tide is low = %v, want ErrCannotUse", err)
	}
}

// A creature with a DestroyedWhen condition survives while the condition is false
// and dies the moment the board makes it true, wherever destruction next settles.
func TestDestroyedWhen(t *testing.T) {
	g := started(t)
	g.AddToBattleline(
		testCreature("Crocag", 7, WithDestroyedWhen(CardsInPlay{
			Player: Opponent,
			Type:   Creature,
			None:   true,
		})), 0)
	foe := g.AddToBattleline(testCreature("foe", 3), 1)

	g.settleDestroyed(0)
	if len(g.Battleline(0)) != 1 {
		t.Fatal("Crocag should survive while the opponent has a creature")
	}
	g.destroyEach(0, []LocalID{foe})
	if len(g.Battleline(0)) != 0 {
		t.Error("Crocag should be destroyed once the opponent has no creatures")
	}
}

// Both new card options are rejected at registration when they are malformed, so
// a broken definition can never reach a game.
func TestCannotBeUsedToRejectsUnsetKind(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("NewCard should reject an unset use kind")
		}
	}()
	NewCard("Bad", Brobnar, Creature, Common, WithPower(1),
		WithCannotBeUsedTo(UseKind(0)))
}

// A constant ability that grants an unset use kind is rejected at registration,
// just like a printed one, so a broken grant can never reach a game.
func TestConstantCannotBeUsedToRejectsUnsetKind(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("NewCard should reject an unset use kind on a constant ability")
		}
	}()
	NewCard("Bad", Brobnar, Creature, Common, WithPower(1),
		WithConstantAbility(ConstantAbility{
			Target:         Target{Kind: TargetEachCreature}.Neighboring(),
			CannotBeUsedTo: []UseKind{UseKind(0)},
		}))
}

// A static modifier that grants an unset use kind is rejected at registration,
// just like a printed or constant-granted one.
func TestStaticCannotBeUsedToRejectsUnsetKind(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("NewCard should reject an unset use kind on a static modifier")
		}
	}()
	NewCard("Bad", StarAlliance, Upgrade, Common,
		WithStatic(StaticModifier{CannotBeUsedTo: []UseKind{UseKind(0)}}))
}

func TestDestroyedWhenRejectsInvalidCondition(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("NewCard should reject an invalid DestroyedWhen condition")
		}
	}()
	NewCard("Bad", Brobnar, Creature, Common, WithPower(1),
		WithDestroyedWhen(PoolAember{Player: Opponent}))
}

// The rules text a card prints for each new option reads the way the card does.
func TestCannotBeUsedToText(t *testing.T) {
	def := NewCard("Crocag", Brobnar, Creature, Common, WithPower(1),
		WithCannotBeUsedTo(ReapUse, FightUse, ActionUse),
		WithDestroyedWhen(CardsInPlay{
			Player: Opponent,
			Type:   Creature,
			None:   true,
		}),
	)
	got := cardRules(&def, false)
	want := []string{
		"Crocag cannot reap.",
		"Crocag cannot fight.",
		"Crocag cannot use its Action ability.",
		"If there are no enemy creatures in play, destroy Crocag.",
	}
	for _, w := range want {
		if !containsLine(got, w) {
			t.Errorf("cardRules = %v, want a line %q", got, w)
		}
	}
	if UseKind(0).valid() {
		t.Error("the zero use kind must not be valid")
	}
	if !ReapUse.valid() || !FightUse.valid() || !ActionUse.valid() {
		t.Error("the real use kinds must be valid")
	}
}

// While a "creatures must fight when used, if able" card is in play, a creature
// with a legal fight target may only fight — reap and Action uses are refused —
// but a creature with nothing to fight is free to reap or act.
func TestMustFightIfAble(t *testing.T) {
	g := started(t)
	g.AddArtifact(NewCard("Little Rapscal Rule", Brobnar, Artifact, Common,
		WithRestrictions(Restrictions{MustFightIfAble: true})), 0)
	brute := g.AddToBattleline(
		testCreature("brute", 5,
			WithAbility(TriggerAction, GainAember{
				Player: Controller,
				Amount: 1,
			})), 0)

	// With no enemy creature, the brute may reap or act freely.
	if err := g.CanUseTo(0, brute, ReapUse); err != nil {
		t.Errorf("reap with nothing to fight = %v, want nil", err)
	}

	foe := g.AddToBattleline(testCreature("foe", 3), 1)
	if err := g.CanUseTo(0, brute, ReapUse); !errors.Is(err, ErrCannotUse) {
		t.Errorf("reap while a fight is available = %v, want ErrCannotUse", err)
	}
	if err := g.Reap(0, brute); !errors.Is(err, ErrCannotUse) {
		t.Errorf("Reap = %v, want ErrCannotUse", err)
	}
	if err := g.UseAction(0, brute); !errors.Is(err, ErrCannotUse) {
		t.Errorf("UseAction while a fight is available = %v, want ErrCannotUse", err)
	}
	if err := g.CanUseTo(0, brute, FightUse); err != nil {
		t.Errorf("fight = %v, want nil", err)
	}
	if err := g.Fight(0, brute, foe); err != nil {
		t.Fatalf("Fight: %v", err)
	}
}

// The must-fight rule only bars reap while the creature is truly able to fight;
// when fighting is barred (Fogbank) it may not fight and so is free to reap,
// keeping CanUse consistent with the individual uses.
func TestMustFightIfAbleUnableToFight(t *testing.T) {
	g := started(t)
	g.AddArtifact(NewCard("Little Rapscal Rule", Brobnar, Artifact, Common,
		WithRestrictions(Restrictions{MustFightIfAble: true})), 0)
	brute := g.AddToBattleline(testCreature("brute", 5), 0)
	g.AddToBattleline(testCreature("foe", 3), 1)
	g.State.CannotFight[0] = Bar[bool]{Value: true}

	if err := g.CanUse(0, brute); err != nil {
		t.Errorf("CanUse = %v, want nil", err)
	}
	if err := g.CanUseTo(0, brute, ReapUse); err != nil {
		t.Errorf("reap while unable to fight = %v, want nil", err)
	}
	if err := g.Reap(0, brute); err != nil {
		t.Errorf("Reap while unable to fight = %v, want nil", err)
	}
}

// A card imposing the must-fight rule prints it in its rules text.
func TestMustFightIfAbleText(t *testing.T) {
	def := NewCard("Rapscal", Brobnar, Creature, Common, WithPower(2),
		WithRestrictions(Restrictions{MustFightIfAble: true}))
	if !containsLine(cardRules(&def, false), "Creatures must fight when used, if able.") {
		t.Errorf("cardRules = %v, want the must-fight line", cardRules(&def, false))
	}
}

// containsLine reports whether lines holds exactly the line want.
func containsLine(lines []string, want string) bool {
	return slices.Contains(lines, want)
}
