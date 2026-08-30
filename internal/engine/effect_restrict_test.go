package engine

import "testing"

func TestGrantFightForChosenHouse(t *testing.T) {
	if got := (GrantFightForChosenHouse{}).Text(); got != "for the remainder of the turn, each friendly creature of the chosen house may fight" {
		t.Errorf("text = %q", got)
	}
	g := NewGame("A", "B", 1)
	GrantFightForChosenHouse{}.Resolve(&EffectContext{Resolver: g, Controller: 0, ChosenHouse: Untamed})
	if g.State.MayFightHouse[0] != Untamed {
		t.Errorf("MayFightHouse[0] = %v, want Untamed", g.State.MayFightHouse[0])
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
	g.BeginTurn(0)
	if err := g.ChooseHouse(0, Brobnar); err != nil {
		t.Fatal(err)
	}
	att := g.AddToBattleline(testCreature("att", 4), 0)
	def := g.AddToBattleline(testCreature("def", 4), 1)

	// Player 0 arms the bar on the opponent during player 0's own turn.
	CannotFight{Player: Opponent, Duration: NextTurn}.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if !g.State.CannotFightNext[1] {
		t.Fatal("CannotFight should arm the opponent's next turn")
	}
	if g.State.CannotFight[0] || g.State.CannotFight[1] {
		t.Error("no bar should be active yet")
	}

	// Player 0 takes an extra turn; the armed bar keeps waiting for player 1.
	g.EndTurn(0)
	g.BeginTurn(0)
	if err := g.ChooseHouse(0, Brobnar); err != nil {
		t.Fatal(err)
	}
	if g.State.CannotFight[0] {
		t.Error("the caster's own turns must never become restricted")
	}
	if !g.State.CannotFightNext[1] {
		t.Error("the opponent's armed bar must survive the caster's extra turns")
	}
	g.EndTurn(0)

	// Player 1's turn: the bar activates and blocks fighting.
	g.BeginTurn(1)
	if err := g.ChooseHouse(1, Brobnar); err != nil {
		t.Fatal(err)
	}
	if !g.State.CannotFight[1] || g.State.CannotFightNext[1] {
		t.Fatal("the bar should be active and disarmed on the opponent's turn")
	}
	if err := g.Fight(1, def, att); err != ErrCannotFight {
		t.Errorf("restricted Fight = %v, want ErrCannotFight", err)
	}
	// It lifts when player 1 ends the turn.
	g.EndTurn(1)
	if g.State.CannotFight[1] {
		t.Error("EndTurn should lift the active bar")
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
	g.AddToBattleline(NewCard("Curse", Brobnar, Creature, Common, WithPower(1), WithRestrictions(Restrictions{Fighting: true})), 0)
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
	if got := restrictionText(Restrictions{PlayCardLimit: PlayCardLimit{Player: Controller, Amount: 1}}); len(got) != 1 ||
		got[0] != "You cannot play more than 1 cards each turn." {
		t.Errorf("controller card limit text = %v", got)
	}
	if got := restrictionText(Restrictions{PlayCardLimit: PlayCardLimit{Player: EachPlayer, Amount: 1}}); len(got) != 1 ||
		got[0] != "Each player cannot play more than 1 cards each turn." {
		t.Errorf("each-player card limit text = %v", got)
	}

	g := started(t) // player 0 active, Brobnar
	if g.cannotPlayCreatures(0) {
		t.Fatal("no restriction yet")
	}
	g.AddToBattleline(NewCard("Blocker", Brobnar, Creature, Common, WithPower(1), WithRestrictions(Restrictions{CannotPlay: Creature})), 0)
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
	if got := restrictionText(Restrictions{Toll: Toll{Action: TollPlayArtifact, Amount: 1}}); len(got) != 1 ||
		got[0] != "Your opponent must give you 1 Æmber in order to play an artifact." {
		t.Errorf("play-toll text = %v", got)
	}
	if got := restrictionText(Restrictions{Toll: Toll{Action: TollUseArtifact, Amount: 2}}); len(got) != 1 ||
		got[0] != "Your opponent must give you 2 Æmber in order to use an artifact." {
		t.Errorf("use-toll text = %v", got)
	}

	g := started(t) // player 0 active, Brobnar
	// Player 0 controls a play-toll; player 1 will be charged to play an artifact.
	g.AddArtifact(NewCard("Customs", Brobnar, Artifact, Common, WithTraits("Location"),
		WithRestrictions(Restrictions{Toll: Toll{Action: TollPlayArtifact, Amount: 1}})), 0)

	g.EndTurn(0)
	g.BeginTurn(1)
	if err := g.ChooseHouse(1, Brobnar); err != nil {
		t.Fatalf("ChooseHouse: %v", err)
	}
	g.AddToHand(NewCard("Widget", Brobnar, Artifact, Common, WithTraits("Item")), 1)
	idx := handIdx(g, 1, "Widget")

	// Too poor to pay the toll: the play is rejected and the card stays in hand.
	if _, err := g.PlayArtifact(1, idx); err != ErrCannotPayToll {
		t.Fatalf("PlayArtifact (broke) = %v, want ErrCannotPayToll", err)
	}

	// With Æmber to spare, the toll transfers to the toll card's owner.
	g.State.Aember[1] = 2
	if _, err := g.PlayArtifact(1, idx); err != nil {
		t.Fatalf("PlayArtifact: %v", err)
	}
	if g.Aember(1) != 1 || g.Aember(0) != 1 {
		t.Errorf("after play toll: p1=%d p0=%d, want 1/1", g.Aember(1), g.Aember(0))
	}

	// The same gate tolls using an artifact: player 0 controls a use-toll, and
	// player 1 must pay it to fire their own artifact's action ability.
	g.AddArtifact(NewCard("Gatekeeper", Brobnar, Artifact, Common, WithTraits("Item"),
		WithRestrictions(Restrictions{Toll: Toll{Action: TollUseArtifact, Amount: 1}})), 0)
	gadget := g.AddArtifact(NewCard("Gadget", Brobnar, Artifact, Common, WithTraits("Item"),
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
	g.BeginTurn(0)
	if err := g.ChooseHouse(0, Brobnar); err != nil {
		t.Fatal(err)
	}

	ForceOpponentActiveHouse{}.Resolve(&EffectContext{Resolver: g, Controller: 0, ChosenHouse: Mars})
	if g.State.ForcedHouseNext[1] != Mars {
		t.Fatalf("armed = %v, want Mars", g.State.ForcedHouseNext[1])
	}

	// Player 0's own choice is unaffected this turn.
	g.EndTurn(0)
	g.BeginTurn(1)
	if g.State.ForcedHouse[1] != Mars {
		t.Errorf("promoted = %v, want Mars", g.State.ForcedHouse[1])
	}
	if err := g.ChooseHouse(1, Sanctum); err != ErrMustChooseForcedHouse {
		t.Errorf("wrong house = %v, want ErrMustChooseForcedHouse", err)
	}
	if err := g.ChooseHouse(1, Mars); err != nil {
		t.Fatalf("forced house: %v", err)
	}

	// The restriction lasts only that one turn.
	g.EndTurn(1)
	g.BeginTurn(1)
	if g.State.ForcedHouse[1] != HouseNone {
		t.Errorf("still forced = %v, want none", g.State.ForcedHouse[1])
	}
	if err := g.ChooseHouse(1, Sanctum); err != nil {
		t.Errorf("free choice = %v, want nil", err)
	}
}
