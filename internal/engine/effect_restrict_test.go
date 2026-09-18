package engine

import "testing"

func TestCannotPlay(t *testing.T) {
	if got := (CannotPlay{Player: Opponent, Type: Creature, Duration: OpponentNextTurn}).Text(); got != "your opponent cannot play creatures during their next turn" {
		t.Errorf("creature text = %q", got)
	}
	if got := (CannotPlay{Player: Controller, Type: Tactic, Duration: OpponentNextTurn}).Text(); got != "you cannot play tactics during your next turn" {
		t.Errorf("tactic text = %q", got)
	}
	if got := (CannotPlay{Player: Controller, Duration: RemainderOfPlayerTurn}).Text(); got != "you cannot play cards for the remainder of the turn" {
		t.Errorf("blanket text = %q", got)
	}
	if (CannotPlay{Type: Creature, Duration: OpponentNextTurn}).validate() == nil {
		t.Error("unset player should be invalid")
	}
	// An unset Type is deliberately legal: it bars every type (Treasure Map's "you
	// cannot play cards"), which the AnyType wildcard carries into the bar.
	if (CannotPlay{Player: Opponent, Duration: OpponentNextTurn}).validate() != nil {
		t.Error("unset card type should mean every type, not be invalid")
	}
	if (CannotPlay{Player: Opponent, Type: Creature}).validate() == nil {
		t.Error("unset duration should be invalid")
	}
	if (CannotPlay{Player: Opponent, Type: Creature, Duration: OpponentNextTurn}).validate() != nil {
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
		Duration: OpponentNextTurn,
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

func TestRestrictFighting(t *testing.T) {
	if got := (Restrict{Player: Opponent, Action: RestrictFighting, Duration: OpponentNextTurn}).Text(); got != "your opponent cannot use creatures to fight during their next turn" {
		t.Errorf("opponent text = %q", got)
	}
	if got := (Restrict{Player: Controller, Action: RestrictFighting, Duration: OpponentNextTurn}).Text(); got != "you cannot use creatures to fight during your next turn" {
		t.Errorf("self text = %q", got)
	}
	// A Restrict needs a player, an action, and a duration; fighting has no
	// current-turn gate, so it is only valid for the player's next turn.
	if (Restrict{Action: RestrictFighting, Duration: OpponentNextTurn}).validate() == nil {
		t.Error("unset player should be invalid")
	}
	if (Restrict{Player: Opponent, Duration: OpponentNextTurn}).validate() == nil {
		t.Error("unset action should be invalid")
	}
	if (Restrict{Player: Opponent, Action: RestrictFighting}).validate() == nil {
		t.Error("unset duration should be invalid")
	}
	if (Restrict{Player: Opponent, Action: RestrictFighting, Duration: RemainderOfPlayerTurn}).validate() == nil {
		t.Error("fighting barred for the current turn should be invalid")
	}
	if (Restrict{Player: Opponent, Action: RestrictFighting, Duration: OpponentNextTurn}).validate() != nil {
		t.Error("a fully set fighting bar should be valid")
	}

	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	if err := g.ChooseHouse(0, Brobnar); err != nil {
		t.Fatal(err)
	}
	att := g.AddToBattleline(testCreature("att", 4), 0)
	def := g.AddToBattleline(testCreature("def", 4), 1)

	// Player 0 arms the bar on the opponent during player 0's own turn.
	Restrict{
		Player:   Opponent,
		Action:   RestrictFighting,
		Duration: OpponentNextTurn,
	}.Resolve(
		&EffectContext{Resolver: g, Controller: 0},
	)
	if !g.State.CannotFightNext[1].Value {
		t.Fatal("Restrict should arm the opponent's next turn")
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

func TestPlayersCannotPlay(t *testing.T) {
	if got := (PlayersCannotPlay{Type: Tactic, Duration: EndOfPlayerNextTurn}).Text(); got != "until the end of your next turn, players cannot play tactics" {
		t.Errorf("tactic text = %q", got)
	}
	// An unset Type bars every type, carried by the AnyType wildcard.
	if got := (PlayersCannotPlay{Duration: EndOfPlayerNextTurn}).Text(); got != "until the end of your next turn, players cannot play cards" {
		t.Errorf("blanket text = %q", got)
	}
	if (PlayersCannotPlay{Type: Tactic}).validate() == nil {
		t.Error("a duration other than EndOfPlayerNextTurn should be invalid")
	}
	if (PlayersCannotPlay{Type: Tactic, Duration: EndOfPlayerNextTurn}).validate() != nil {
		t.Error("EndOfPlayerNextTurn should be valid")
	}

	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	if err := g.ChooseHouse(0, Brobnar); err != nil {
		t.Fatal(err)
	}
	// Player 0 arms the board-wide bar: their own this turn and next, the opponent's
	// next turn.
	PlayersCannotPlay{Type: Tactic, Duration: EndOfPlayerNextTurn}.Resolve(
		&EffectContext{Resolver: g, Controller: 0},
	)
	if g.State.CannotPlayTypeThis[0].Value != Tactic {
		t.Error("the caster should be barred this turn")
	}
	if g.State.CannotPlayTypeNext[0].Value != Tactic {
		t.Error("the caster's next turn should be armed")
	}
	if g.State.CannotPlayTypeNext[1].Value != Tactic {
		t.Error("the opponent's next turn should be armed")
	}
}

func TestForOpponentNextTurnStunsFighters(t *testing.T) {
	e := ForOpponentNextTurn{
		On: EventFight,
		Do: Stun{Target: Target{Kind: TargetTriggeringCreature}},
	}
	if got := e.Text(); got != "during your opponent's next turn, after an enemy creature is used to fight, stun it" {
		t.Errorf("text = %q", got)
	}
	if err := e.validate(); err != nil {
		t.Errorf("valid effect rejected: %v", err)
	}

	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	if err := g.ChooseHouse(0, Brobnar); err != nil {
		t.Fatal(err)
	}
	// Player 0 arms the reaction on the opponent for their next turn.
	e.Resolve(&EffectContext{Resolver: g, Controller: 0})

	// It lies dormant for the rest of the caster's own turn: a creature the caster
	// uses to fight is not stunned.
	mine := g.AddToBattleline(testCreature("mine", 6), 0)
	theirs := g.AddToBattleline(testCreature("theirs", 2), 1)
	if err := g.Fight(0, mine, theirs); err != nil {
		t.Fatalf("Fight = %v", err)
	}
	if g.State.Cards[mine].Stunned {
		t.Error("the caster's own fighter should not be stunned this turn")
	}
	g.EndPlayPhase(0)

	// The opponent's turn: each creature they use to fight is stunned after the fight.
	g.StartTurn(1)
	if err := g.ChooseHouse(1, Brobnar); err != nil {
		t.Fatal(err)
	}
	att := g.AddToBattleline(testCreature("att", 6), 1)
	att2 := g.AddToBattleline(testCreature("att2", 6), 1)
	def := g.AddToBattleline(testCreature("def", 2), 0)
	def2 := g.AddToBattleline(testCreature("def2", 2), 0)
	if err := g.Fight(1, att, def); err != nil {
		t.Fatalf("Fight = %v", err)
	}
	if err := g.Fight(1, att2, def2); err != nil {
		t.Fatalf("Fight = %v", err)
	}
	if !g.State.Cards[att].Stunned || !g.State.Cards[att2].Stunned {
		t.Error("each fighting creature should be stunned")
	}

	// The reaction clears at the opponent's ready phase.
	g.EndPlayPhase(1)
	if g.State.LastingCount != 0 {
		t.Error("the reaction should clear at the opponent's ready phase")
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

func TestCannotReapConstant(t *testing.T) {
	// restrictionText renders each relative player.
	if got := restrictionText(Restrictions{Reaping: Controller}, false); len(got) != 1 ||
		got[0] != "Your creatures cannot reap." {
		t.Errorf("controller reaping text = %v", got)
	}
	if got := restrictionText(Restrictions{Reaping: Opponent}, false); len(got) != 1 ||
		got[0] != "Enemy creatures cannot reap." {
		t.Errorf("opponent reaping text = %v", got)
	}
	if got := restrictionText(Restrictions{Reaping: EachPlayer}, false); len(got) != 1 ||
		got[0] != "Creatures cannot reap." {
		t.Errorf("each-player reaping text = %v", got)
	}
	// The bonus-icon bar renders each relative player (Master of the Grey).
	if got := restrictionText(Restrictions{BonusIcons: Controller}, false); len(got) != 1 ||
		got[0] != "You cannot resolve bonus icons on cards you play." {
		t.Errorf("controller bonus-icon text = %v", got)
	}
	if got := restrictionText(Restrictions{BonusIcons: Opponent}, false); len(got) != 1 ||
		got[0] != "Your opponent cannot resolve bonus icons on cards they play." {
		t.Errorf("opponent bonus-icon text = %v", got)
	}
	if got := restrictionText(Restrictions{BonusIcons: EachPlayer}, false); len(got) != 1 ||
		got[0] != "Players cannot resolve bonus icons on cards they play." {
		t.Errorf("each-player bonus-icon text = %v", got)
	}
	// A use-condition phrases against the card, or the host creature on an upgrade.
	cond := Restrictions{UseCondition: CardsDiscarded{Player: Controller, Amount: 1}}
	if got := restrictionText(cond, false); len(got) != 1 ||
		got[0] != "You cannot use this card unless you have discarded a card from your hand this turn." {
		t.Errorf("card use-condition text = %v", got)
	}
	if got := restrictionText(cond, true); len(got) != 1 ||
		got[0] != "This creature cannot be used unless you have discarded a card from your hand this turn." {
		t.Errorf("upgrade use-condition text = %v", got)
	}

	g := started(t) // player 0 active, Brobnar
	reaper := g.AddToBattleline(testCreature("reaper", 4), 0)
	if g.cannotReap(0) || g.cannotReap(1) {
		t.Fatal("no restriction yet")
	}
	// A card player 0 controls with Reaping Opponent bars only player 1.
	joya := g.AddToBattleline(
		NewCard("Joya", Brobnar, Creature, Common, WithPower(5),
			WithRestrictions(Restrictions{Reaping: Opponent})),
		0,
	)
	if g.cannotReap(0) {
		t.Error("Reaping Opponent should not bar the controller")
	}
	if !g.cannotReap(1) {
		t.Error("Reaping Opponent should bar the opponent")
	}
	// EachPlayer bars everyone.
	eachDef := NewCard("Joya", Brobnar, Creature, Common, WithPower(5),
		WithRestrictions(Restrictions{Reaping: EachPlayer}))
	g.cat.defs[joya] = &eachDef
	if !g.cannotReap(0) || !g.cannotReap(1) {
		t.Error("Reaping EachPlayer should bar both players")
	}
	// Controller bars only the card's controller.
	ctrlDef := NewCard("Joya", Brobnar, Creature, Common, WithPower(5),
		WithRestrictions(Restrictions{Reaping: Controller}))
	g.cat.defs[joya] = &ctrlDef
	if !g.cannotReap(0) || g.cannotReap(1) {
		t.Error("Reaping Controller should bar only the controller")
	}
	// The active player cannot reap. With an enemy creature present the reaper can
	// still fight, so canUse passes and canUseTo's ReapUse guard is what stops it.
	g.AddToBattleline(testCreature("enemy", 4), 1)
	if err := g.CanUseTo(0, reaper, ReapUse); err != ErrCannotUse {
		t.Errorf("CanUseTo(reap) while barred = %v, want ErrCannotUse", err)
	}
	if err := g.CanUseTo(0, reaper, FightUse); err != nil {
		t.Errorf("CanUseTo(fight) while only reaping is barred = %v, want nil", err)
	}
	if err := g.Reap(0, reaper); err != ErrCannotUse {
		t.Errorf("Reap while barred = %v, want ErrCannotUse", err)
	}
	// reapWith is a no-op too, so a forced reap gains no Æmber.
	g.reapWith(reaper)
	if g.State.Aember[0] != 0 {
		t.Errorf("forced reap while barred gained %d Æmber, want 0", g.State.Aember[0])
	}
}

func TestCannotPlayCreatures(t *testing.T) {
	if got := restrictionText(
		Restrictions{Fighting: true, CannotPlay: Creature},
		false,
	); len(got) != 2 ||
		got[0] != "You cannot use creatures to fight." ||
		got[1] != "You cannot play creatures." {
		t.Errorf("restrictionText = %v", got)
	}
	if got := restrictionText(
		Restrictions{PlayCardLimit: PlayCardLimit{Player: Controller, Amount: 1}}, false,
	); len(
		got,
	) != 1 ||
		got[0] != "You cannot play more than 1 cards each turn." {
		t.Errorf("controller card limit text = %v", got)
	}
	if got := restrictionText(
		Restrictions{PlayCardLimit: PlayCardLimit{Player: EachPlayer, Amount: 1}}, false,
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
		Restrictions{Toll: Toll{Action: TollPlayArtifact, Amount: 1}}, false,
	); len(
		got,
	) != 1 ||
		got[0] != "In order to play an artifact, your opponent must give you 1 Æmber." {
		t.Errorf("play-toll text = %v", got)
	}
	if got := restrictionText(
		Restrictions{Toll: Toll{Action: TollUseArtifact, Amount: 2}}, false,
	); len(
		got,
	) != 1 ||
		got[0] != "In order to use an artifact, your opponent must give you 2 Æmber." {
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
	e := MustChooseHouse{Player: Opponent, Reference: ChosenActiveHouse}
	if got := e.Text(); got != "your opponent must choose that house as their active house during their next turn" {
		t.Errorf("text = %q", got)
	}
	// Player Controller reads in the first person, so a card can bind its own turn.
	if got := (MustChooseHouse{Player: Controller, Reference: ChosenActiveHouse}).Text(); got !=
		"you must choose that house as your active house during your next turn" {
		t.Errorf("controller text = %q", got)
	}
	if err := e.validate(); err != nil {
		t.Errorf("a Chosen must should validate: %v", err)
	}
	if err := (MustChooseHouse{Player: Opponent, Reference: JustChosenActiveHouse}).validate(); err == nil {
		t.Error("a JustChosen must should not validate")
	}
	if err := (MustChooseHouse{Reference: ChosenActiveHouse}).validate(); err == nil {
		t.Error("an unset mode should not validate")
	}
	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	if err := g.ChooseHouse(0, Brobnar); err != nil {
		t.Fatal(err)
	}

	e.Resolve(
		&EffectContext{Resolver: g, Controller: 0, ChosenHouse: Mars},
	)
	if got := g.State.HouseConstraintsNext[1][0]; g.State.HouseConstraintCountNext[1] != 1 ||
		got.Kind != constraintMustHouse || got.House != Mars {
		t.Fatalf("armed = %+v (count %d), want one must-Mars",
			got, g.State.HouseConstraintCountNext[1])
	}

	// Player 0's own choice is unaffected this turn.
	g.EndPlayPhase(0)
	g.StartTurn(1)
	if got := g.State.HouseConstraints[1][0]; g.State.HouseConstraintCount[1] != 1 ||
		got.Kind != constraintMustHouse || got.House != Mars {
		t.Errorf("promoted = %+v, want must-Mars", got)
	}
	if err := g.ChooseHouse(1, Sanctum); err != ErrHouseNotAllowed {
		t.Errorf("wrong house = %v, want ErrHouseNotAllowed", err)
	}
	if err := g.ChooseHouse(1, Mars); err != nil {
		t.Fatalf("forced house: %v", err)
	}

	// The restriction lasts only that one turn.
	g.EndPlayPhase(1)
	g.StartTurn(1)
	if g.State.HouseConstraintCount[1] != 0 {
		t.Errorf("still forced (count %d), want none", g.State.HouseConstraintCount[1])
	}
	if err := g.ChooseHouse(1, Sanctum); err != nil {
		t.Errorf("free choice = %v, want nil", err)
	}
}

func TestForceActiveHouseOfFoughtNextTurn(t *testing.T) {
	e := MustChooseHouse{Player: Opponent, Reference: FoughtActiveHouse}
	if err := e.validate(); err != nil {
		t.Errorf("a Fought must should validate: %v", err)
	}
	if got := e.Text(); got != "your opponent must choose the house of the creature {self} fights as their active house during their next turn" {
		t.Errorf("text = %q", got)
	}
	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	if err := g.ChooseHouse(0, Dis); err != nil {
		t.Fatal(err)
	}
	foe := g.AddToBattleline(NewCard("foe", Logos, Creature, Common, WithPower(3)), 1)
	e.Resolve(&EffectContext{Resolver: g, Controller: 0, It: foe, HasIt: true})
	if got := g.State.HouseConstraintsNext[1][0]; g.State.HouseConstraintCountNext[1] != 1 ||
		got.Kind != constraintMustCreature || got.Creature != foe {
		t.Errorf(
			"armed = %+v, want a live must referencing the fought creature %d",
			got, foe,
		)
	}

	// With no creature in context there is nothing to force.
	g.State.HouseConstraintCountNext[1] = 0
	e.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if g.State.HouseConstraintCountNext[1] != 0 {
		t.Errorf(
			"armed without a fought creature (count %d), want none",
			g.State.HouseConstraintCountNext[1],
		)
	}
}

// TestForceActiveHouseOfItNextTurn covers Mark of Dis's It reference: the damaged
// creature's own controller is forced to choose that creature's house next turn,
// whichever side the creature is on, and nothing arms without a creature in context.
func TestForceActiveHouseOfItNextTurn(t *testing.T) {
	e := MustChooseHouse{Player: ItsController, Reference: ItActiveHouse}
	if err := e.validate(); err != nil {
		t.Errorf("an It must should validate: %v", err)
	}
	if got := e.Text(); got != "its controller must choose that creature's house as their active house during their next turn" {
		t.Errorf("text = %q", got)
	}
	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	if err := g.ChooseHouse(0, Dis); err != nil {
		t.Fatal(err)
	}
	// The creature belongs to the opponent, so its own controller (player 1) is bound.
	foe := g.AddToBattleline(NewCard("foe", Logos, Creature, Common, WithPower(3)), 1)
	e.Resolve(&EffectContext{Resolver: g, Controller: 0, It: foe, HasIt: true})
	if got := g.State.HouseConstraintsNext[1][0]; g.State.HouseConstraintCountNext[1] != 1 ||
		got.Kind != constraintMustCreature || got.Creature != foe {
		t.Errorf(
			"armed = %+v, want a live must on player 1 referencing creature %d",
			got, foe,
		)
	}

	// With no creature in context there is nothing to force.
	g.State.HouseConstraintCountNext[1] = 0
	e.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if g.State.HouseConstraintCountNext[1] != 0 {
		t.Errorf(
			"armed without a creature (count %d), want none",
			g.State.HouseConstraintCountNext[1],
		)
	}
}

// the opponent's next active house and the payoff when they match it.
func TestWagerOpponentChoosesChosenHouse(t *testing.T) {
	e := WagerOpponentChoosesChosenHouse{Amount: 2}
	want := "if your opponent chooses that house as their active house during their next turn, steal 2 Æmber"
	if got := e.Text(); got != want {
		t.Errorf("text = %q", got)
	}

	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	if err := g.ChooseHouse(0, Brobnar); err != nil {
		t.Fatal(err)
	}
	e.Resolve(&EffectContext{Resolver: g, Controller: 0, ChosenHouse: Logos})
	if got := g.State.HouseConstraintsNext[1][0]; g.State.HouseConstraintCountNext[1] != 1 ||
		got.Kind != constraintWager || got.House != Logos || got.Amount != 2 || got.Predictor != 0 {
		t.Fatalf("armed = %+v, want wager Logos/2/predictor 0", got)
	}

	// The opponent locks in the wagered house next turn, so the predictor steals.
	g.EndPlayPhase(0)
	g.StartTurn(1)
	g.SetAember(1, 5)
	if err := g.ChooseHouse(1, Logos); err != nil {
		t.Fatal(err)
	}
	if g.Aember(0) != 2 || g.Aember(1) != 3 {
		t.Errorf("aember = %d/%d, want 2/3 (predictor stole 2)", g.Aember(0), g.Aember(1))
	}
}

// TestWagerMissed covers the payoff that does not land: the opponent chooses a
// house other than the wagered one, so the wager is spent for nothing.
func TestWagerMissed(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	if err := g.ChooseHouse(0, Brobnar); err != nil {
		t.Fatal(err)
	}
	WagerOpponentChoosesChosenHouse{Amount: 2}.Resolve(
		&EffectContext{Resolver: g, Controller: 0, ChosenHouse: Logos},
	)
	g.EndPlayPhase(0)
	g.StartTurn(1)
	g.SetAember(1, 5)
	if err := g.ChooseHouse(1, Mars); err != nil {
		t.Fatal(err)
	}
	if g.Aember(0) != 0 || g.Aember(1) != 5 {
		t.Errorf("aember = %d/%d, want 0/5 (no steal on a missed wager)", g.Aember(0), g.Aember(1))
	}
}

// TestForbidSameActiveHouseNextTurn covers Snag's Mirror: after a player chooses
// their active house, their opponent cannot choose that same house next turn.
func TestForbidSameActiveHouseNextTurn(t *testing.T) {
	e := CannotChooseHouse{Player: Opponent, Reference: JustChosenActiveHouse}
	if got := e.Text(); got != "their opponent cannot choose the same house as their active house during their next turn" {
		t.Errorf("text = %q", got)
	}
	if err := e.validate(); err != nil {
		t.Errorf("a JustChosen cannot should validate: %v", err)
	}
	if err := (CannotChooseHouse{Player: Opponent, Reference: FoughtActiveHouse}).validate(); err == nil {
		t.Error("a Fought cannot should not validate")
	}
	if err := (CannotChooseHouse{Reference: ChosenActiveHouse}).validate(); err == nil {
		t.Error("an unset player should not validate")
	}

	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	if err := g.ChooseHouse(0, Brobnar); err != nil {
		t.Fatal(err)
	}
	e.Resolve(&EffectContext{Resolver: g})
	if got := g.State.HouseConstraintsNext[1][0]; g.State.HouseConstraintCountNext[1] != 1 ||
		got.Kind != constraintCannotHouse || got.House != Brobnar {
		t.Errorf("armed = %+v, want cannot-Brobnar", got)
	}
}

func TestForbidActiveHouseNextTurn(t *testing.T) {
	e := CannotChooseHouse{Player: Opponent, Reference: ChosenActiveHouse}
	if got := e.Text(); got != "your opponent cannot choose that house as their active house during their next turn" {
		t.Errorf("text = %q", got)
	}
	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	if err := g.ChooseHouse(0, Brobnar); err != nil {
		t.Fatal(err)
	}

	e.Resolve(
		&EffectContext{Resolver: g, Controller: 0, ChosenHouse: Mars},
	)
	if got := g.State.HouseConstraintsNext[1][0]; g.State.HouseConstraintCountNext[1] != 1 ||
		got.Kind != constraintCannotHouse || got.House != Mars {
		t.Fatalf("armed = %+v, want cannot-Mars", got)
	}

	// Player 0's own choice is unaffected this turn.
	g.EndPlayPhase(0)
	g.StartTurn(1)
	if got := g.State.HouseConstraints[1][0]; g.State.HouseConstraintCount[1] != 1 ||
		got.Kind != constraintCannotHouse || got.House != Mars {
		t.Errorf("promoted = %+v, want cannot-Mars", got)
	}
	if err := g.ChooseHouse(1, Mars); err != ErrHouseNotAllowed {
		t.Errorf("forbidden house = %v, want ErrHouseNotAllowed", err)
	}
	if err := g.ChooseHouse(1, Sanctum); err != nil {
		t.Fatalf("allowed house: %v", err)
	}

	// The restriction lasts only that one turn.
	g.EndPlayPhase(1)
	g.StartTurn(1)
	if g.State.HouseConstraintCount[1] != 0 {
		t.Errorf("still forbidden (count %d), want none", g.State.HouseConstraintCount[1])
	}
	if err := g.ChooseHouse(1, Mars); err != nil {
		t.Errorf("free choice = %v, want nil", err)
	}
}

func TestRestrictionSources(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	weak := g.AddToBattleline(testCreature("Control the Weak", 1), 0)
	fog := g.AddToBattleline(testCreature("Fogbank", 1), 0)

	g.MustChooseHouseNextTurn(1, Mars, weak)
	// A card imposing three bars is named once.
	g.CannotFightNextTurn(1, fog)
	g.SkipForgePhaseNextTurn(1, fog)
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

// TestRestrictionSourcesConditionalPlayBar checks that a continuous
// CannotPlayWhile bar (Quixxle Stone) names its card only while its condition
// currently holds against the player, and stops once it lifts.
func TestRestrictionSourcesConditionalPlayBar(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	stone := g.AddArtifact(
		NewCard(
			"Quixxle Stone",
			StarAlliance,
			Artifact,
			Rare,
			WithCannotPlayWhile(ConditionalPlayBar{Type: Creature, When: ControlsMoreCreatures{}}),
		),
		1, // the opponent controls it, yet it bars whichever side is ahead
	)

	// Equal creature counts: nobody is ahead, so the bar names no card.
	if got := g.RestrictionSources(0); len(got) != 0 {
		t.Errorf("sources with equal counts = %v, want none", got)
	}

	// Player 0 pulls ahead: the bar now binds player 0 and names the stone.
	g.AddToBattleline(testCreature("ahead", 3), 0)
	got := g.RestrictionSources(0)
	if len(got) != 1 || got[0] != stone {
		t.Errorf("sources while ahead = %v, want [%d]", got, stone)
	}
	// The opponent, who is behind, is not bound.
	if got := g.RestrictionSources(1); len(got) != 0 {
		t.Errorf("sources for the behind player = %v, want none", got)
	}
}

func TestUseConditionRestriction(t *testing.T) {
	cond := CardsDiscarded{Player: Controller, House: namedHouse(Untamed), Amount: 1}
	want := "You cannot use this card unless you have discarded an Untamed card " +
		"from your hand this turn."
	if got := restrictionText(
		Restrictions{UseCondition: cond},
		false,
	); len(got) != 1 ||
		got[0] != want {
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

	CannotPlay{Player: Controller, Duration: RemainderOfPlayerTurn}.Resolve(
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
// "Action:" — throughout their next turn (Skippy Timehog) or for the rest of the
// current turn (United Action).
func TestRestrictUse(t *testing.T) {
	if got := (Restrict{Player: Opponent, Action: RestrictUse, Duration: OpponentNextTurn}).Text(); got != "your opponent cannot use any cards during their next turn" {
		t.Errorf("opponent text = %q", got)
	}
	if got := (Restrict{Player: Controller, Action: RestrictUse, Duration: OpponentNextTurn}).Text(); got != "you cannot use any cards during your next turn" {
		t.Errorf("controller text = %q", got)
	}
	if got := (Restrict{Player: Controller, Action: RestrictUse, Duration: RemainderOfPlayerTurn}).Text(); got != "you cannot use any cards for the remainder of the turn" {
		t.Errorf("this-turn text = %q", got)
	}
	if got := (Restrict{Player: Opponent, Action: RestrictUse, Duration: RemainderOfPlayerTurn}).Text(); got != "your opponent cannot use any cards for the remainder of the turn" {
		t.Errorf("this-turn opponent text = %q", got)
	}
	if (Restrict{Player: Opponent, Action: RestrictUse, Duration: OpponentNextTurn}).validate() != nil {
		t.Error("a fully set use bar should be valid")
	}

	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	Restrict{Player: Opponent, Action: RestrictUse, Duration: OpponentNextTurn}.Resolve(
		&EffectContext{Resolver: g, Controller: 0},
	)
	// The RemainderOfPlayerTurn form arms the bar for the rest of the current turn.
	Restrict{Player: Controller, Action: RestrictUse, Duration: RemainderOfPlayerTurn}.Resolve(
		&EffectContext{Resolver: g, Controller: 0},
	)
	if !g.State.CannotUse[0].Value {
		t.Error("RemainderOfPlayerTurn should arm the use bar this turn")
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

// TestRestrictReaping covers the timed player-wide bar that stops a player reaping
// with any creature — throughout their next turn (Inky Gloom) or for the rest of
// the current turn (Ragnarok) — while fighting stays open.
func TestRestrictReaping(t *testing.T) {
	if got := (Restrict{Player: Opponent, Action: RestrictReaping, Duration: OpponentNextTurn}).Text(); got != "your opponent cannot use creatures to reap during their next turn" {
		t.Errorf("opponent text = %q", got)
	}
	if got := (Restrict{Player: Controller, Action: RestrictReaping, Duration: OpponentNextTurn}).Text(); got != "you cannot use creatures to reap during your next turn" {
		t.Errorf("controller text = %q", got)
	}
	if got := (Restrict{Player: Controller, Action: RestrictReaping, Duration: RemainderOfPlayerTurn}).Text(); got != "you cannot use creatures to reap for the remainder of the turn" {
		t.Errorf("this-turn text = %q", got)
	}
	if got := (Restrict{Player: Opponent, Action: RestrictReaping, Duration: RemainderOfPlayerTurn}).Text(); got != "your opponent cannot use creatures to reap for the remainder of the turn" {
		t.Errorf("this-turn opponent text = %q", got)
	}
	if (Restrict{Player: Opponent, Action: RestrictReaping, Duration: OpponentNextTurn}).validate() != nil {
		t.Error("a fully set reaping bar should be valid")
	}

	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	if err := g.ChooseHouse(0, Brobnar); err != nil {
		t.Fatal(err)
	}
	Restrict{Player: Opponent, Action: RestrictReaping, Duration: OpponentNextTurn}.Resolve(
		&EffectContext{Resolver: g, Controller: 0},
	)
	// RemainderOfPlayerTurn arms the current-turn reap bar directly on the resolving player.
	Restrict{Player: Controller, Action: RestrictReaping, Duration: RemainderOfPlayerTurn}.Resolve(
		&EffectContext{Resolver: g, Controller: 0},
	)
	if !g.State.CannotReap[0].Value {
		t.Error("RemainderOfPlayerTurn should arm this turn's reap bar")
	}
	// A creature the barred player controls cannot reap while the bar is up.
	mine := g.AddToBattleline(NewCard("mine", Brobnar, Creature, Common, WithPower(3)), 0)
	if err := g.Reap(0, mine); err != ErrCannotUse {
		t.Errorf("Reap while barred this turn = %v, want ErrCannotUse", err)
	}
	foe := g.AddToBattleline(NewCard("foe", Brobnar, Creature, Common, WithPower(3)), 0)
	g.EndPlayPhase(0)
	if g.State.CannotReap[0].Value {
		t.Error("the this-turn reap bar should lift at the end of the turn")
	}

	g.StartTurn(1)
	if err := g.ChooseHouse(1, Brobnar); err != nil {
		t.Fatal(err)
	}
	beast := g.AddToBattleline(NewCard("beast", Brobnar, Creature, Common, WithPower(3)), 1)
	if err := g.Reap(1, beast); err != ErrCannotUse {
		t.Errorf("Reap while barred = %v, want ErrCannotUse", err)
	}
	// The bar stops reaping only — fighting stays open.
	if err := g.Fight(1, beast, foe); err != nil {
		t.Errorf("Fight while reap-barred = %v, want nil", err)
	}
	g.EndPlayPhase(1)
	if g.State.CannotReap[1].Value {
		t.Error("the reap bar should lift at end of turn")
	}
}

// TestRestrictReapHouseNextTurn covers the house-scoped reap bar armed for a
// player's next turn (Seismo-entangler), a Restrict narrowed by a chosen house.
func TestRestrictReapHouseNextTurn(t *testing.T) {
	scoped := Restrict{
		Player:   Opponent,
		Action:   RestrictReaping,
		House:    TheChosenHouse,
		Duration: OpponentNextTurn,
	}
	if got := scoped.Text(); got != "your opponent cannot use creatures of the chosen house to reap during their next turn" {
		t.Errorf("scoped text = %q", got)
	}
	if scoped.validate() != nil {
		t.Error("a fully set house-scoped reap bar should be valid")
	}
	// A house scope only bars reaping on the player's next turn.
	if (Restrict{
		Player:   Opponent,
		Action:   RestrictFighting,
		House:    TheChosenHouse,
		Duration: OpponentNextTurn,
	}).validate() == nil {
		t.Error("a house scope on a fight bar should be invalid")
	}
	if (Restrict{
		Player:   Opponent,
		Action:   RestrictReaping,
		House:    TheChosenHouse,
		Duration: RemainderOfPlayerTurn,
	}).validate() == nil {
		t.Error("a house scope on a current-turn reap bar should be invalid")
	}

	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	scoped.Resolve(&EffectContext{Resolver: g, Controller: 0, ChosenHouse: Brobnar})
	if g.State.CannotReapHouseNext[1].Value != Brobnar {
		t.Fatalf("armed = %v, want Brobnar", g.State.CannotReapHouseNext[1].Value)
	}
	g.EndPlayPhase(0)

	g.StartTurn(1)
	if err := g.ChooseHouse(1, Brobnar); err != nil {
		t.Fatal(err)
	}
	if g.State.CannotReapHouse[1].Value != Brobnar {
		t.Fatalf("promoted = %v, want Brobnar", g.State.CannotReapHouse[1].Value)
	}
	beast := g.AddToBattleline(NewCard("beast", Brobnar, Creature, Common, WithPower(3)), 1)
	if err := g.Reap(1, beast); err != ErrCannotUse {
		t.Errorf("Reap barred house = %v, want ErrCannotUse", err)
	}
	if g.CanUse(1, beast) == nil {
		t.Error("a barred creature with no other use should have no use")
	}
	g.EndPlayPhase(1)
	if g.State.CannotReapHouse[1].Value != HouseNone {
		t.Error("the reap bar should lift at end of turn")
	}
}

// TestCreaturesCannotFight covers the board-wide fight bar with a house exception
// (Into the Night): every non-excepted creature, on either side, is stopped from
// fighting until the caster's next turn, while reaping and the spared house stay
// open.
func TestCreaturesCannotFight(t *testing.T) {
	withHouse := CreaturesCannot{
		Action:   FightUse,
		Houses:   exceptHouse(Shadows),
		Duration: StartOfPlayerNextTurn,
	}
	if got := withHouse.Text(); got != "until the start of your next turn, non-Shadows creatures cannot be used to fight" {
		t.Errorf("house-exception text = %q", got)
	}
	noHouse := CreaturesCannot{Action: ReapUse, Duration: StartOfPlayerNextTurn}
	if got := noHouse.Text(); got != "until the start of your next turn, creatures cannot be used to reap" {
		t.Errorf("no-exception text = %q", got)
	}

	if (CreaturesCannot{Duration: StartOfPlayerNextTurn}).validate() == nil {
		t.Error("unset action should be invalid")
	}
	if (CreaturesCannot{Action: ActionUse, Duration: StartOfPlayerNextTurn}).validate() == nil {
		t.Error("an Action-ability bar should be invalid")
	}
	if (CreaturesCannot{
		Action:   FightUse,
		Houses:   HouseMatcher{Kind: MatchChosenHouse},
		Duration: StartOfPlayerNextTurn,
	}).validate() == nil {
		t.Error("a context-dependent house matcher should be invalid")
	}
	if (CreaturesCannot{
		Action:   FightUse,
		Houses:   HouseMatcher{Kind: MatchExceptHouse},
		Duration: StartOfPlayerNextTurn,
	}).validate() == nil {
		t.Error("a house matcher missing its house should be invalid")
	}
	if (CreaturesCannot{Action: FightUse}).validate() == nil {
		t.Error("unset duration should be invalid")
	}
	if (CreaturesCannot{Action: FightUse, Duration: StartOfPlayerNextTurn}).validate() != nil {
		t.Error("a fully set effect should be valid")
	}

	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	if err := g.ChooseHouse(0, Shadows); err != nil {
		t.Fatal(err)
	}
	withHouse.Resolve(&EffectContext{Resolver: g, Controller: 0})

	// The caster is barred this turn; the opponent's next turn is armed.
	if g.State.CreaturesCannot[0].Value.Action != FightUse {
		t.Fatal("the caster's own turn should be barred immediately")
	}
	if g.State.CreaturesCannotNext[1].Value.Action != FightUse {
		t.Fatal("the opponent's next turn should be armed")
	}

	// The exception spares Shadows creatures and bars every other house; reaping
	// stays open. The bar rides on each player's own bar slot, so the opponent's
	// creatures are reached on the opponent's turn (below), not the caster's.
	myShadows := g.AddToBattleline(NewCard("mine-sh", Shadows, Creature, Common, WithPower(3)), 0)
	myBrobnar := g.AddToBattleline(NewCard("mine-br", Brobnar, Creature, Common, WithPower(3)), 0)
	if g.creaturesGloballyBarred(myShadows, FightUse) {
		t.Error("a Shadows creature should be spared by the exception")
	}
	if !g.creaturesGloballyBarred(myBrobnar, FightUse) {
		t.Error("a non-Shadows creature should be barred from fighting")
	}
	if g.creaturesGloballyBarred(myBrobnar, ReapUse) {
		t.Error("a fight bar must not stop reaping")
	}

	// End to end on the caster's own turn: a spared Shadows creature can fight in
	// the active Shadows house.
	def := g.AddToBattleline(NewCard("def", Brobnar, Creature, Common, WithPower(3)), 1)
	if err := g.Fight(0, myShadows, def); err != nil {
		t.Errorf("spared Fight = %v, want nil", err)
	}

	// The caster's own bar lifts at end of their turn.
	g.EndPlayPhase(0)
	if g.State.CreaturesCannot[0].Value.Action != useKindUnset {
		t.Error("the caster's bar should lift at end of their turn")
	}

	// The opponent's turn: the armed bar activates and blocks fighting.
	g.StartTurn(1)
	if err := g.ChooseHouse(1, Brobnar); err != nil {
		t.Fatal(err)
	}
	if g.State.CreaturesCannot[1].Value.Action != FightUse {
		t.Fatal("the opponent's bar should activate on their turn")
	}
	if g.State.CreaturesCannotNext[1].Value.Action != useKindUnset {
		t.Error("the armed slot should be disarmed once promoted")
	}
	p1c := g.AddToBattleline(NewCard("p1-br", Brobnar, Creature, Common, WithPower(3)), 1)
	target := g.AddToBattleline(NewCard("tgt", Brobnar, Creature, Common, WithPower(3)), 0)
	if err := g.Fight(1, p1c, target); err != ErrCannotUse {
		t.Errorf("opponent barred Fight = %v, want ErrCannotUse", err)
	}
	g.EndPlayPhase(1)

	// The caster's next turn is clean.
	g.StartTurn(0)
	if err := g.ChooseHouse(0, Shadows); err != nil {
		t.Fatal(err)
	}
	if g.State.CreaturesCannot[0].Value.Action != useKindUnset {
		t.Error("the caster's next turn should be free of the bar")
	}
}

// TestCreaturesCannotReap covers the board-wide reap bar with no house exception
// (Sow Salt): every creature on either side is stopped from reaping until the
// caster's next turn, while fighting stays open.
func TestCreaturesCannotReap(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	if err := g.ChooseHouse(0, Brobnar); err != nil {
		t.Fatal(err)
	}
	CreaturesCannot{Action: ReapUse, Duration: StartOfPlayerNextTurn}.Resolve(
		&EffectContext{Resolver: g, Controller: 0},
	)

	mine := g.AddToBattleline(NewCard("mine", Brobnar, Creature, Common, WithPower(3)), 0)
	if !g.creaturesGloballyBarred(mine, ReapUse) {
		t.Error("with no exception, every creature should be barred from reaping")
	}
	if g.creaturesGloballyBarred(mine, FightUse) {
		t.Error("a reap bar must not stop fighting")
	}
	if err := g.Reap(0, mine); err != ErrCannotUse {
		t.Errorf("barred Reap = %v, want ErrCannotUse", err)
	}

	// The opponent's turn: the bar reaches their creatures too.
	g.EndPlayPhase(0)
	g.StartTurn(1)
	if err := g.ChooseHouse(1, Brobnar); err != nil {
		t.Fatal(err)
	}
	foe := g.AddToBattleline(NewCard("foe", Brobnar, Creature, Common, WithPower(3)), 1)
	if err := g.Reap(1, foe); err != ErrCannotUse {
		t.Errorf("opponent barred Reap = %v, want ErrCannotUse", err)
	}
	g.EndPlayPhase(1)

	// The caster's next turn: reaping works again.
	g.StartTurn(0)
	if err := g.ChooseHouse(0, Brobnar); err != nil {
		t.Fatal(err)
	}
	back := g.AddToBattleline(NewCard("back", Brobnar, Creature, Common, WithPower(3)), 0)
	if err := g.Reap(0, back); err != nil {
		t.Errorf("Reap after the bar lifts = %v, want nil", err)
	}
}
