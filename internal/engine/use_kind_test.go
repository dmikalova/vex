package engine

import "testing"

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
	if err := g.CanUseTo(0, crocag, ReapUse); err != ErrCannotUse {
		t.Errorf("CanUseTo(reap) = %v, want ErrCannotUse", err)
	}
	if err := g.CanUseTo(0, crocag, FightUse); err != nil {
		t.Errorf("CanUseTo(fight) = %v, want nil", err)
	}
	if err := g.Reap(0, crocag); err != ErrCannotUse {
		t.Errorf("Reap = %v, want ErrCannotUse", err)
	}
	if g.Exhausted(crocag) {
		t.Error("a refused reap must not exhaust the creature")
	}
	if err := g.Fight(0, crocag, foe); err != nil {
		t.Fatalf("Fight: %v", err)
	}
}

// A card barred from fighting or from its Action: ability is refused those uses,
// and an effect that reaps with a card which cannot reap does nothing.
func TestCannotBeUsedToFightAndAction(t *testing.T) {
	g := started(t)
	pacifist := g.AddToBattleline(
		testCreature("pacifist", 4,
			WithCannotBeUsedTo(FightUse),
			WithAbility(TriggerAction, GainAember{Player: Controller, Amount: 1}),
		), 0)
	foe := g.AddToBattleline(testCreature("foe", 3), 1)
	if err := g.Fight(0, pacifist, foe); err != ErrCannotUse {
		t.Errorf("Fight = %v, want ErrCannotUse", err)
	}

	idle := g.AddToBattleline(
		testCreature("idle", 4,
			WithCannotBeUsedTo(ActionUse),
			WithAbility(TriggerAction, GainAember{Player: Controller, Amount: 1}),
		), 0)
	if err := g.UseAction(0, idle); err != ErrCannotUse {
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
	if err := g.CanUse(0, crocag); err != ErrCannotUse {
		t.Errorf("CanUse with an empty enemy battleline = %v, want ErrCannotUse", err)
	}

	// An Action: ability is a use of its own, so it keeps the card usable.
	actor := g.AddToBattleline(
		testCreature("actor", 7,
			WithCannotBeUsedTo(ReapUse),
			WithAbility(TriggerAction, GainAember{Player: Controller, Amount: 1}),
		), 0)
	if err := g.CanUse(0, actor); err != nil {
		t.Errorf("CanUse with an Action ability = %v, want nil", err)
	}

	// A player-wide fight ban closes the last door on a creature that cannot reap.
	g.AddToBattleline(testCreature("foe", 3), 1)
	g.State.CannotFight[0].Value = true
	if err := g.CanUse(0, crocag); err != ErrCannotUse {
		t.Errorf("CanUse while fighting is banned = %v, want ErrCannotUse", err)
	}
}

// A creature with a DestroyedWhen condition survives while the condition is false
// and dies the moment the board makes it true, wherever destruction next settles.
func TestDestroyedWhen(t *testing.T) {
	g := started(t)
	g.AddToBattleline(
		testCreature("Crocag", 7, WithDestroyedWhen(InPlay{
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

func TestDestroyedWhenRejectsInvalidCondition(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("NewCard should reject an invalid DestroyedWhen condition")
		}
	}()
	NewCard("Bad", Brobnar, Creature, Common, WithPower(1),
		WithDestroyedWhen(OpponentAember{}))
}

// The rules text a card prints for each new option reads the way the card does.
func TestCannotBeUsedToText(t *testing.T) {
	def := NewCard("Crocag", Brobnar, Creature, Common, WithPower(1),
		WithCannotBeUsedTo(ReapUse, FightUse, ActionUse),
		WithDestroyedWhen(InPlay{Player: Opponent, Type: Creature, None: true}),
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

// containsLine reports whether lines holds exactly the line want.
func containsLine(lines []string, want string) bool {
	for _, l := range lines {
		if l == want {
			return true
		}
	}
	return false
}
