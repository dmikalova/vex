package engine

import "testing"

func TestEffectValidation(t *testing.T) {
	bad := Heal{Fully: true, Amount: 1, Target: Target{Kind: TargetThisCreature}}
	good := Heal{Fully: true, Target: Target{Kind: TargetThisCreature}}

	if err := validateEffect(GainAember{Player: Controller, Amount: 1}); err != nil {
		t.Errorf("non-validating effect should be nil, got %v", err)
	}
	if err := (Sequence{Effects: []Effect{good, bad}}).validate(); err == nil {
		t.Error("sequence should surface a bad child")
	}
	if err := (Sentences{Effects: []Effect{good, bad}}).validate(); err == nil {
		t.Error("sentences should surface a bad child")
	}
	if err := (Sentences{Effects: []Effect{good}}).validate(); err != nil {
		t.Errorf("sentences of valid effects should pass, got %v", err)
	}
	if err := (Sequence{Effects: []Effect{good, GainAember{Player: Controller, Amount: 1}}}).validate(); err != nil {
		t.Errorf("sequence of valid effects should pass, got %v", err)
	}
	if err := (Conditional{Then: bad}).validate(); err == nil {
		t.Error("conditional should surface a bad gated effect")
	}
	if err := (Conditional{Then: GainAember{Player: Controller, Amount: 1}}).validate(); err != nil {
		t.Errorf("conditional with a valid effect should pass, got %v", err)
	}
	if err := validateEffect(PutFromDiscard{Destination: ToBottomOfDeck}); err == nil {
		t.Error("PutFromDiscard to an unsupported destination should be rejected")
	}
	if err := validateEffect(PutFromDiscard{Destination: ToTopOfDeck}); err != nil {
		t.Errorf("PutFromDiscard to the top of the deck should pass, got %v", err)
	}
	// Purge must name the zone it pulls from.
	if err := validateEffect(PurgeCard{}); err == nil {
		t.Error("a Purge with no zone should be rejected")
	}
	if err := validateEffect(PurgeCard{Zone: Discard, Type: Creature}); err != nil {
		t.Errorf("a Purge naming its zone should pass, got %v", err)
	}
	// A result gate surfaces a bad first action or a bad follow-up.
	if err := validateEffect(
		Then{
			First:  PurgeCard{},
			Result: AddPowerCounter{Target: Target{Kind: TargetThisCreature}, Amount: 1},
		},
	); err == nil {
		t.Error("result gate should surface a bad first action")
	}
	if err := validateEffect(Then{First: PurgeCard{Zone: Discard}, Result: bad}); err == nil {
		t.Error("result gate should surface a bad follow-up")
	}
	if err := validateEffect(
		Then{
			First:  PurgeCard{Zone: Discard, Type: Creature},
			Result: AddPowerCounter{Target: Target{Kind: TargetThisCreature}, Amount: 1},
		},
	); err != nil {
		t.Errorf("result gate with valid halves should pass, got %v", err)
	}
}

func TestRequiredTargetValidation(t *testing.T) {
	// Player-taking effects reject an unset player and accept an explicit one.
	cases := []struct {
		name       string
		unset, set Effect
	}{
		{"GainAember", GainAember{Amount: 1}, GainAember{Player: Controller, Amount: 1}},
		{"LoseAember", LoseAember{Amount: 1}, LoseAember{Player: Controller, Amount: 1}},
		{"DiscardArchives", DiscardArchives{}, DiscardArchives{Player: Controller}},
		{"DiscardHand", DiscardHand{}, DiscardHand{Player: Controller}},
		{"Reveal", RevealHand{}, RevealHand{Player: Controller}},
	}
	for _, tc := range cases {
		if err := validateEffect(tc.unset); err == nil {
			t.Errorf("%s with an unset player should be rejected", tc.name)
		}
		if err := validateEffect(tc.set); err != nil {
			t.Errorf("%s with a player set should pass, got %v", tc.name, err)
		}
	}

	// Target-taking effects reject an unset target and accept an explicit one.
	this := Target{Kind: TargetThisCreature}
	targetCases := []struct {
		name       string
		unset, set Effect
	}{
		{"DealDamage", DealDamage{Amount: 1}, DealDamage{Amount: 1, Target: this}},
		{"Destroy", Destroy{}, Destroy{Target: this}},
		{"Exalt", Exalt{Amount: 1}, Exalt{Amount: 1, Target: this}},
		{"Exhaust", Exhaust{}, Exhaust{Target: this}},
		{"Ready", Ready{}, Ready{Target: this}},
		{"ReadyIfFirstUse", ReadyIfFirstUse{}, ReadyIfFirstUse{Target: this}},
		{"ReadyCreatures", ReadyCreatures{}, ReadyCreatures{Target: this}},
		{"Stun", Stun{}, Stun{Target: this}},
		{"Unstun", Unstun{}, Unstun{Target: this}},
		{"PurgeCreature", PurgeCreature{}, PurgeCreature{Target: this}},
		{"OnChooseCreature", OnChooseCreature{}, OnChooseCreature{Target: this}},
		{"AddPowerCounter", AddPowerCounter{Amount: 1}, AddPowerCounter{Amount: 1, Target: this}},
		{"RedirectFightDamage", RedirectFightDamage{}, RedirectFightDamage{Target: this}},
	}
	for _, tc := range targetCases {
		if err := validateEffect(tc.unset); err == nil {
			t.Errorf("%s with an unset target should be rejected", tc.name)
		}
		if err := validateEffect(tc.set); err != nil {
			t.Errorf("%s with a target set should pass, got %v", tc.name, err)
		}
	}

	// CannotFight needs both a player and a duration.
	if err := validateEffect(CannotFight{Duration: NextTurn}); err == nil {
		t.Error("CannotFight with an unset player should be rejected")
	}
	if err := validateEffect(CannotFight{Player: Opponent}); err == nil {
		t.Error("CannotFight with an unset duration should be rejected")
	}
	if err := validateEffect(CannotFight{Player: Opponent, Duration: NextTurn}); err != nil {
		t.Errorf("CannotFight fully set should pass, got %v", err)
	}
}

func TestPlayerForRejectsUnsetPlayer(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("PlayerFor(playerUnset) should panic")
		}
	}()
	(&EffectContext{Controller: 0}).PlayerFor(playerUnset)
}

func TestNewCardRejectsUnsetTrigger(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("NewCard should reject an ability whose trigger is unset")
		}
	}()
	NewCard("Bad", Brobnar, Tactic, Common,
		WithAbility(triggerUnset, GainAember{Player: Controller, Amount: 1}))
}

func TestNewCardRejectsConflictingHeal(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("NewCard should panic on a Heal with both Amount and Fully")
		}
	}()
	NewCard(
		"bad",
		Sanctum,
		Creature,
		Common,
		WithAbility(
			TriggerAfterPlay,
			Heal{Amount: 2, Fully: true, Target: Target{Kind: TargetThisCreature}},
		),
	)
}

func TestNewCardRejectsInvalidReplace(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("NewCard should panic on a StaticModifier.Replaces with a bad effect")
		}
	}()
	NewCard("bad", Sanctum, Upgrade, Rare, WithStatic(StaticModifier{
		Replaces: Replace{
			When: EventCreatureDestroyed,
			With: Heal{Amount: 2, Fully: true, Target: Target{Kind: TargetTriggeringCreature}},
		},
	}))
}

func TestNewCardRejectsInvalidReplaces(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("NewCard should panic on a continuous Replaces using a reaction event")
		}
	}()
	NewCard(
		"bad",
		Mars,
		Creature,
		Uncommon,
		WithReplaces(Instead{Of: EventCreaturePlayed, With: Capture}),
	)
}

func TestNewCardRejectsInvalidPlayPermission(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("NewCard should panic on a granted PlayPermission with no count")
		}
	}()
	NewCard("bad", Untamed, Creature, Rare, WithPlayPermission(PlayPermission{House: Untamed}))
}
