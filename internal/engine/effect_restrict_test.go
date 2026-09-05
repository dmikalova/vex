package engine

import "testing"

func TestCannotPlay(t *testing.T) {
	if got := (CannotPlay{Player: Opponent, Type: Creature, Duration: NextTurn}).Text(); got != "your opponent cannot play creatures during their next turn" {
		t.Errorf("creature text = %q", got)
	}
	if got := (CannotPlay{Player: Controller, Type: Tactic, Duration: NextTurn}).Text(); got != "you cannot play Tactics during your next turn" {
		t.Errorf("tactic text = %q", got)
	}
	if got := (CannotPlay{Player: Controller, Duration: EndOfTurn}).Text(); got != "you cannot play cards for the remainder of the turn" {
		t.Errorf("blanket text = %q", got)
	}
	if (CannotPlay{Type: Creature, Duration: NextTurn}).validate() == nil {
		t.Error("unset player should be invalid")
	}
	// An unset Type is deliberately legal: it bars every type (Treasure Map's "you
	// cannot play cards"), which the AnyType wildcard carries into the bar.
	if (CannotPlay{Player: Opponent, Duration: NextTurn}).validate() != nil {
		t.Error("unset card type should mean every type, not be invalid")
	}
	if (CannotPlay{Player: Opponent, Type: Creature}).validate() == nil {
		t.Error("unset duration should be invalid")
	}
	if (CannotPlay{Player: Opponent, Type: Creature, Duration: NextTurn}).validate() != nil {
		t.Error("a fully set effect should be valid")
	}

	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	if err := g.ChooseHouse(0, Brobnar); err != nil {
		t.Fatal(err)
	}
	// Resolve arms the opponent's next turn.
	CannotPlay{
		Player:   Opponent,
		Type:     Creature,
		Duration: NextTurn,
	}.Resolve(
		&EffectContext{Resolver: g, Controller: 0},
	)
	if g.State.CannotPlayTypeNext[1].Value != Creature {
		t.Fatal("the bar should arm the opponent's next turn")
	}

	// Activate the bar and confirm the play path rejects a barred creature.
	g.State.CannotPlayTypeThis[0].Value = Creature
	idx := int(g.State.Hand[0].Count)
	g.AddToHand(NewCard("beast", Brobnar, Creature, Common, WithPower(3)), 0)
	if _, err := g.PlayCreature(0, idx, false); err != ErrCannotPlayType {
		t.Errorf("playing a barred creature = %v, want ErrCannotPlayType", err)
	}
}

func TestGrantFightForChosenHouse(t *testing.T) {
	if got := (GrantFightForChosenHouse{}).Text(); got != "for the remainder of the turn, each friendly creature of the chosen house may fight" {
		t.Errorf("text = %q", got)
	}
	g := NewGame("A", "B", 1)
	GrantFightForChosenHouse{}.Resolve(
		&EffectContext{Resolver: g, Controller: 0, ChosenHouse: Untamed},
	)
	if g.State.MayFightHouse[0] != Untamed {
		t.Errorf("MayFightHouse[0] = %v, want Untamed", g.State.MayFightHouse[0])
	}
}

func TestGrantFightForFriendlyHouse(t *testing.T) {
	if err := (GrantFightForFriendlyHouse{}).validate(); err == nil {
		t.Error("an unset house should be rejected")
	}
	e := GrantFightForFriendlyHouse{House: Brobnar}
	if err := e.validate(); err != nil {
		t.Errorf("validate = %v, want nil", err)
	}
	if got := e.Text(); got != "for the remainder of the turn, each friendly Brobnar creature may fight" {
		t.Errorf("text = %q", got)
	}
	g := NewGame("A", "B", 1)
	e.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if g.State.MayFightHouse[0] != Brobnar {
		t.Errorf("MayFightHouse[0] = %v, want Brobnar", g.State.MayFightHouse[0])
	}
}

func TestCannotFight(t *testing.T) {
	if got := (CannotFight{Player: Opponent}).Text(); got != "your opponent cannot use creatures to fight during their next turn" {
		t.Errorf("opponent text = %q", got)
	}
	if got := (CannotFight{Player: Controller}).Text(); got != "you cannot use creatures to fight during your next turn" {
		t.Errorf("self text = %q", got)
	}

	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	if err := g.ChooseHouse(0, Brobnar); err != nil {
		t.Fatal(err)
	}
	att := g.AddToBattleline(testCreature("att", 4), 0)
	def := g.AddToBattleline(testCreature("def", 4), 1)

	// Player 0 arms the bar on the opponent during player 0's own turn.
	CannotFight{
		Player:   Opponent,
		Duration: NextTurn,
	}.Resolve(
		&EffectContext{Resolver: g, Controller: 0},
	)
	if !g.State.CannotFightNext[1].Value {
		t.Fatal("CannotFight should arm the opponent's next turn")
	}
	if g.State.CannotFight[0].Value || g.State.CannotFight[1].Value {
		t.Error("no bar should be active yet")
	}

	// Player 0 takes an extra turn; the armed bar keeps waiting for player 1.
	g.EndPlayPhase(0)
	g.StartTurn(0)
	if err := g.ChooseHouse(0, Brobnar); err != nil {
		t.Fatal(err)
	}
	if g.State.CannotFight[0].Value {
		t.Error("the caster's own turns must never become restricted")
	}
	if !g.State.CannotFightNext[1].Value {
		t.Error("the opponent's armed bar must survive the caster's extra turns")
	}
	g.EndPlayPhase(0)

	// Player 1's turn: the bar activates and blocks fighting.
	g.StartTurn(1)
	if err := g.ChooseHouse(1, Brobnar); err != nil {
		t.Fatal(err)
	}
	if !g.State.CannotFight[1].Value || g.State.CannotFightNext[1].Value {
		t.Fatal("the bar should be active and disarmed on the opponent's turn")
	}
	if err := g.Fight(1, def, att); err != ErrCannotFight {
		t.Errorf("restricted Fight = %v, want ErrCannotFight", err)
	}
	// It lifts when player 1 ends the turn.
	g.EndPlayPhase(1)
	if g.State.CannotFight[1].Value {
		t.Error("the ready phase should lift the active bar")
	}
}

func TestCannotFightConstant(t *testing.T) {
	g := started(t) // player 0 active, Brobnar
	att := g.AddToBattleline(testCreature("att", 4), 0)
	def := g.AddToBattleline(testCreature("def", 4), 1)
	if g.cannotFight(0) {
		t.Fatal("no restriction yet")
	}
	// A card in play with a constant Fighting restriction bars its controller.
	g.AddToBattleline(
		NewCard(
			"Curse",
			Brobnar,
			Creature,
			Common,
			WithPower(1),
			WithRestrictions(Restrictions{Fighting: true}),
		),
		0,
	)
	if !g.cannotFight(0) {
		t.Error("a constant Fighting restriction should bar the controller")
	}
	if err := g.Fight(0, att, def); err != ErrCannotFight {
		t.Errorf("Fight = %v, want ErrCannotFight", err)
	}
}

func TestCannotPlayCreatures(t *testing.T) {
	if got := restrictionText(Restrictions{Fighting: true, CannotPlay: Creature}); len(got) != 2 ||
		got[0] != "You cannot use creatures to fight." || got[1] != "You cannot play creatures." {
		t.Errorf("restrictionText = %v", got)
	}
	if got := restrictionText(
		Restrictions{PlayCardLimit: PlayCardLimit{Player: Controller, Amount: 1}},
	); len(
		got,
	) != 1 ||
		got[0] != "You cannot play more than 1 cards each turn." {
		t.Errorf("controller card limit text = %v", got)
	}
	if got := restrictionText(
		Restrictions{PlayCardLimit: PlayCardLimit{Player: EachPlayer, Amount: 1}},
	); len(
		got,
	) != 1 ||
		got[0] != "Each player cannot play more than 1 cards each turn." {
		t.Errorf("each-player card limit text = %v", got)
	}

	g := started(t) // player 0 active, Brobnar
	if g.cannotPlayCreatures(0) {
		t.Fatal("no restriction yet")
	}
	g.AddToBattleline(
		NewCard(
			"Blocker",
			Brobnar,
			Creature,
			Common,
			WithPower(1),
			WithRestrictions(Restrictions{CannotPlay: Creature}),
		),
		0,
	)
	if !g.cannotPlayCreatures(0) {
		t.Fatal("the restriction should bar creature plays")
	}
	g.AddToHand(testCreature("newbie", 2), 0)
	if _, err := g.PlayCreature(0, handIdx(g, 0, "newbie"), false); err != ErrCannotPlayCreature {
		t.Errorf("PlayCreature = %v, want ErrCannotPlayCreature", err)
	}
	// Non-creature plays are unaffected.
	g.AddToHand(NewCard("act", Brobnar, Tactic, Common), 0)
	if err := g.PlayAction(0, handIdx(g, 0, "act")); err != nil {
		t.Errorf("actions should still be playable: %v", err)
	}
}

func TestToll(t *testing.T) {
	// Text renders for both actions a toll can charge for.
	if got := restrictionText(
		Restrictions{Toll: Toll{Action: TollPlayArtifact, Amount: 1}},
	); len(
		got,
	) != 1 ||
		got[0] != "Your opponent must give you 1 Æmber in order to play an artifact." {
		t.Errorf("play-toll text = %v", got)
	}
	if got := restrictionText(
		Restrictions{Toll: Toll{Action: TollUseArtifact, Amount: 2}},
	); len(
		got,
	) != 1 ||
		got[0] != "Your opponent must give you 2 Æmber in order to use an artifact." {
		t.Errorf("use-toll text = %v", got)
	}

	g := started(t) // player 0 active, Brobnar
	// Player 0 controls a play-toll; player 1 will be charged to play an artifact.
	g.AddArtifact(NewCard("Customs", Brobnar, Artifact, Common, WithTraits(Location),
		WithRestrictions(Restrictions{Toll: Toll{Action: TollPlayArtifact, Amount: 1}})), 0)

	g.EndPlayPhase(0)
	g.StartTurn(1)
	if err := g.ChooseHouse(1, Brobnar); err != nil {
		t.Fatalf("ChooseHouse: %v", err)
	}
	g.AddToHand(NewCard("Widget", Brobnar, Artifact, Common, WithTraits(Item)), 1)
	idx := handIdx(g, 1, "Widget")

	// Too poor to pay the toll: the play is rejected and the card stays in hand, and
	// CanPlay says so up front rather than letting the caller find out on the click.
	if err := g.CanPlay(1, g.Hand(1)[idx]); err != ErrCannotPayToll {
		t.Fatalf("CanPlay (broke) = %v, want ErrCannotPayToll", err)
	}
	if _, err := g.PlayArtifact(1, idx); err != ErrCannotPayToll {
		t.Fatalf("PlayArtifact (broke) = %v, want ErrCannotPayToll", err)
	}

	// With Æmber to spare, the toll transfers to the toll card's owner.
	g.State.Aember[1] = 2
	if err := g.CanPlay(1, g.Hand(1)[idx]); err != nil {
		t.Fatalf("CanPlay (funded) = %v, want nil", err)
	}
	if _, err := g.PlayArtifact(1, idx); err != nil {
		t.Fatalf("PlayArtifact: %v", err)
	}
	if g.Aember(1) != 1 || g.Aember(0) != 1 {
		t.Errorf("after play toll: p1=%d p0=%d, want 1/1", g.Aember(1), g.Aember(0))
	}

	// The same gate tolls using an artifact: player 0 controls a use-toll, and
	// player 1 must pay it to fire their own artifact's action ability.
	g.AddArtifact(NewCard("Gatekeeper", Brobnar, Artifact, Common, WithTraits(Item),
		WithRestrictions(Restrictions{Toll: Toll{Action: TollUseArtifact, Amount: 1}})), 0)
	gadget := g.AddArtifact(NewCard("Gadget", Brobnar, Artifact, Common, WithTraits(Item),
		WithAbility(TriggerAction, GainAember{Player: Controller, Amount: 3})), 1)

	g.State.Aember[1] = 0
	if err := g.UseAction(1, gadget); err != ErrCannotPayToll {
		t.Fatalf("UseAction (broke) = %v, want ErrCannotPayToll", err)
	}
	g.State.Aember[1] = 1
	before := g.Aember(0)
	if err := g.UseAction(1, gadget); err != nil {
		t.Fatalf("UseAction: %v", err)
	}
	if g.Aember(1) != 3 { // paid 1, then the action gained 3
		t.Errorf("player 1 aember = %d, want 3", g.Aember(1))
	}
	if g.Aember(0) != before+1 {
		t.Errorf("player 0 aember = %d, want %d (received the use toll)", g.Aember(0), before+1)
	}
}

func TestForceActiveHouseNextTurn(t *testing.T) {
	if got := (ForceOpponentActiveHouse{}).Text(); got != "your opponent must choose that house as their active house during their next turn" {
		t.Errorf("text = %q", got)
	}
	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	if err := g.ChooseHouse(0, Brobnar); err != nil {
		t.Fatal(err)
	}

	ForceOpponentActiveHouse{}.Resolve(
		&EffectContext{Resolver: g, Controller: 0, ChosenHouse: Mars},
	)
	if g.State.ForcedHouseNext[1].Value != Mars {
		t.Fatalf("armed = %v, want Mars", g.State.ForcedHouseNext[1].Value)
	}

	// Player 0's own choice is unaffected this turn.
	g.EndPlayPhase(0)
	g.StartTurn(1)
	if g.State.ForcedHouse[1].Value != Mars {
		t.Errorf("promoted = %v, want Mars", g.State.ForcedHouse[1].Value)
	}
	if err := g.ChooseHouse(1, Sanctum); err != ErrMustChooseForcedHouse {
		t.Errorf("wrong house = %v, want ErrMustChooseForcedHouse", err)
	}
	if err := g.ChooseHouse(1, Mars); err != nil {
		t.Fatalf("forced house: %v", err)
	}

	// The restriction lasts only that one turn.
	g.EndPlayPhase(1)
	g.StartTurn(1)
	if g.State.ForcedHouse[1].Value != HouseNone {
		t.Errorf("still forced = %v, want none", g.State.ForcedHouse[1].Value)
	}
	if err := g.ChooseHouse(1, Sanctum); err != nil {
		t.Errorf("free choice = %v, want nil", err)
	}
}

func TestRestrictionSources(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	weak := g.AddToBattleline(testCreature("Control the Weak", 1), 0)
	fog := g.AddToBattleline(testCreature("Fogbank", 1), 0)

	g.ForceActiveHouseNextTurn(1, Mars, weak)
	// A card imposing three bars is named once.
	g.CannotFightNextTurn(1, fog)
	g.SkipForgeStepNextTurn(1, fog)
	g.CannotPlayTypeNextTurn(1, Creature, fog)
	// The armed cards are not binding anyone until the affected player's turn.
	if got := g.RestrictionSources(1); len(got) != 0 {
		t.Errorf("sources before promotion = %v, want none", got)
	}

	g.EndPlayPhase(0)
	g.StartTurn(1)
	got := g.RestrictionSources(1)
	if len(got) != 2 || got[0] != fog || got[1] != weak {
		t.Errorf("promoted sources = %v, want [%d %d]", got, fog, weak)
	}

	// Lifting a bar stops it naming its card.
	g.EndPlayPhase(1)
	if got := g.RestrictionSources(1); len(got) != 2 {
		t.Errorf("sources after the fight bar lifts = %v, want the other two", got)
	}
}

func TestUseConditionRestriction(t *testing.T) {
	cond := CardsDiscarded{Player: Controller, House: Untamed, Amount: 1}
	want := "You cannot use this card unless you have discarded an Untamed card " +
		"from your hand this turn."
	if got := restrictionText(Restrictions{UseCondition: cond}); len(got) != 1 || got[0] != want {
		t.Errorf("restrictionText = %v", got)
	}

	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	if err := g.ChooseHouse(0, Untamed); err != nil {
		t.Fatalf("ChooseHouse: %v", err)
	}
	sloth := g.AddToBattleline(NewCard(
		"Giant Sloth",
		Untamed,
		Creature,
		Rare,
		WithPower(6),
		WithRestrictions(Restrictions{UseCondition: cond}),
		WithAbility(TriggerAction, GainAember{Player: Controller, Amount: 3}),
	), 0)
	g.State.Cards[sloth].Exhausted = false
	if err := g.CanUse(0, sloth); err != ErrCannotUse {
		t.Errorf("unmet use condition = %v, want ErrCannotUse", err)
	}

	discarded := g.AddToHand(NewCard("Untamed Action", Untamed, Tactic, Common), 0)
	g.DiscardCardFromHand(0, discarded)
	if err := g.CanUse(0, sloth); err != nil {
		t.Errorf("met use condition = %v, want nil", err)
	}
}

// TestCannotPlayBlanketThisTurn covers the AnyType blanket bar armed for the rest
// of the current turn (Treasure Map), which every play path consults.
func TestCannotPlayBlanketThisTurn(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	if err := g.ChooseHouse(0, Brobnar); err != nil {
		t.Fatal(err)
	}
	beast := g.AddToHand(NewCard("beast", Brobnar, Creature, Common, WithPower(3)), 0)

	CannotPlay{Player: Controller, Duration: EndOfTurn}.Resolve(
		&EffectContext{Resolver: g, Controller: 0},
	)
	if g.State.CannotPlayTypeThis[0].Value != AnyType {
		t.Fatalf("bar = %q, want the AnyType wildcard", g.State.CannotPlayTypeThis[0].Value)
	}
	if err := g.CanPlay(0, beast); err != ErrCannotPlayType {
		t.Errorf("CanPlay = %v, want ErrCannotPlayType", err)
	}
	if _, err := g.PlayCreature(0, int(g.State.Hand[0].Count)-1, false); err != ErrCannotPlayType {
		t.Errorf("PlayCreature = %v, want ErrCannotPlayType", err)
	}

	g.EndPlayPhase(0)
	if g.State.CannotPlayTypeThis[0].Value != TypeUnset {
		t.Error("the blanket bar should lift at end of turn")
	}
}

// TestCannotUse covers the bar that stops a player reaping, fighting, or firing an
// "Action:" throughout their next turn (Skippy Timehog).
func TestCannotUse(t *testing.T) {
	if got := (CannotUse{Player: Opponent, Duration: NextTurn}).Text(); got != "your opponent cannot use any cards during their next turn" {
		t.Errorf("opponent text = %q", got)
	}
	if got := (CannotUse{Player: Controller, Duration: NextTurn}).Text(); got != "you cannot use any cards during your next turn" {
		t.Errorf("controller text = %q", got)
	}
	if (CannotUse{Duration: NextTurn}).validate() == nil {
		t.Error("unset player should be invalid")
	}
	if (CannotUse{Player: Opponent}).validate() == nil {
		t.Error("unset duration should be invalid")
	}
	if (CannotUse{Player: Opponent, Duration: NextTurn}).validate() != nil {
		t.Error("a fully set effect should be valid")
	}

	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	CannotUse{Player: Opponent, Duration: NextTurn}.Resolve(
		&EffectContext{Resolver: g, Controller: 0},
	)
	// A duration the effect does not handle arms nothing.
	CannotUse{Player: Controller, Duration: EndOfTurn}.Resolve(
		&EffectContext{Resolver: g, Controller: 0},
	)
	if g.State.CannotUse[0].Value {
		t.Error("only NextTurn arms the use bar")
	}
	g.EndPlayPhase(0)

	g.StartTurn(1)
	if err := g.ChooseHouse(1, Brobnar); err != nil {
		t.Fatal(err)
	}
	beast := g.AddToBattleline(NewCard("beast", Brobnar, Creature, Common, WithPower(3)), 1)
	if err := g.Reap(1, beast); err != ErrCannotUse {
		t.Errorf("Reap = %v, want ErrCannotUse", err)
	}
	g.EndPlayPhase(1)
	if g.State.CannotUse[1].Value {
		t.Error("the use bar should lift at end of turn")
	}
}

// TestGrantFightAnyHouse covers the house-blind fight grant (Follow the Leader).
func TestGrantFightAnyHouse(t *testing.T) {
	if got := (GrantFightAnyHouse{}).Text(); got != "for the remainder of the turn, each friendly creature may fight" {
		t.Errorf("text = %q", got)
	}

	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	if err := g.ChooseHouse(0, Brobnar); err != nil {
		t.Fatal(err)
	}
	outsider := g.AddToBattleline(NewCard("outsider", Logos, Creature, Common, WithPower(5)), 0)
	enemy := g.AddToBattleline(NewCard("enemy", Dis, Creature, Common, WithPower(2)), 1)
	if err := g.Fight(0, outsider, enemy); err != ErrWrongHouse {
		t.Fatalf("Fight before the grant = %v, want ErrWrongHouse", err)
	}

	GrantFightAnyHouse{}.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if err := g.Fight(0, outsider, enemy); err != nil {
		t.Fatalf("Fight after the grant = %v, want nil", err)
	}
	g.EndPlayPhase(0)
	if g.State.MayFightAny[0] {
		t.Error("the grant should lift at end of turn")
	}
}

func TestMayUseFriendlyHouse(t *testing.T) {
	if got := (MayUseFriendlyHouse{House: Sanctum}).Text(); got != "for the remainder of the turn, you may use friendly Sanctum creatures" {
		t.Errorf("text = %q", got)
	}
	if (MayUseFriendlyHouse{}).validate() == nil {
		t.Error("unset house should be invalid")
	}
	if (MayUseFriendlyHouse{House: Sanctum}).validate() != nil {
		t.Error("a set house should be valid")
	}

	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	if err := g.ChooseHouse(0, Brobnar); err != nil {
		t.Fatal(err)
	}
	c := g.AddToBattleline(NewCard("cleric", Sanctum, Creature, Common, WithPower(3)), 0)
	if g.usableInActiveHouse(c) {
		t.Fatal("an off-house creature should not be usable before the grant")
	}

	MayUseFriendlyHouse{House: Sanctum}.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if g.State.MayUseHouse[0] != Sanctum {
		t.Fatal("the grant should record the house")
	}
	if !g.usableInActiveHouse(c) {
		t.Error("the granted-house creature should be usable")
	}

	g.EndPlayPhase(0)
	if g.State.MayUseHouse[0] != HouseNone {
		t.Error("the grant should clear at end of turn")
	}
}

func TestMayPlayOrUseFriendlyHouse(t *testing.T) {
	if got := (MayPlayOrUseFriendlyHouse{House: Mars}).Text(); got != "you may play or use a Mars card this turn" {
		t.Errorf("text = %q", got)
	}
	if (MayPlayOrUseFriendlyHouse{}).validate() == nil {
		t.Error("unset house should be invalid")
	}
	if (MayPlayOrUseFriendlyHouse{House: Mars}).validate() != nil {
		t.Error("a set house should be valid")
	}

	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	if err := g.ChooseHouse(0, Sanctum); err != nil {
		t.Fatal(err)
	}
	off := NewCard("marauder", Mars, Creature, Common, WithPower(3))
	if g.mayPlayFromHand(0, &off) {
		t.Fatal("an off-house card should not be playable before the grant")
	}
	c := g.AddToBattleline(NewCard("cleric", Mars, Creature, Common, WithPower(3)), 0)
	if g.usableInActiveHouse(c) {
		t.Fatal("an off-house creature should not be usable before the grant")
	}

	MayPlayOrUseFriendlyHouse{House: Mars}.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if g.State.MayPlayHouse[0] != Mars {
		t.Fatal("the grant should record the play house")
	}
	if g.State.MayUseHouse[0] != Mars {
		t.Fatal("the grant should record the use house")
	}
	if !g.mayPlayFromHand(0, &off) {
		t.Error("the granted-house card should be playable")
	}
	if !g.usableInActiveHouse(c) {
		t.Error("the granted-house creature should be usable")
	}

	g.EndPlayPhase(0)
	if g.State.MayPlayHouse[0] != HouseNone {
		t.Error("the play grant should clear at end of turn")
	}
}
