package engine

import (
	"strings"
	"testing"
)

func TestRenderPrompt(t *testing.T) {
	prompt := "fully heal " + SelfName
	if got := renderPrompt("Chuff Ape", prompt); got != "fully heal Chuff Ape" {
		t.Errorf("renderPrompt = %q, want %q", got, "fully heal Chuff Ape")
	}
	// An unattributed prompt has no name to substitute, so it is left as it is.
	if got := renderPrompt("", prompt); got != prompt {
		t.Errorf("unattributed renderPrompt = %q, want %q", got, prompt)
	}
}

func TestEntersStunnedText(t *testing.T) {
	def := NewCard(
		"Chuff",
		Mars,
		Creature,
		Common,
		WithPower(3),
		WithEntersPlay(Stun{Target: Target{Kind: TargetThisCreature}}),
	)
	want := "Chuff enters play stunned."
	found := false
	for _, line := range cardRules(&def, false) {
		if line == want {
			found = true
		}
	}
	if !found {
		t.Errorf("cardRules missing %q; got %v", want, cardRules(&def, false))
	}
}

// TestCannotBeUsedWhileText covers the rules line for a card barred from use while
// a condition holds — Valoocanth's "While the tide is low, ... cannot be used."
func TestCannotBeUsedWhileText(t *testing.T) {
	def := NewCard("Valoocanth", Untamed, Creature, Common, WithPower(6),
		WithCannotBeUsedWhile(TideIsLow{}))
	want := "While the tide is low, Valoocanth cannot be used."
	found := false
	for _, line := range cardRules(&def, false) {
		if line == want {
			found = true
		}
	}
	if !found {
		t.Errorf("cardRules missing %q; got %v", want, cardRules(&def, false))
	}
}

func TestEntersPlayAbilityText(t *testing.T) {
	// A Stun effect renders as the "stunned" state word.
	stun := RenderAbility(
		Ability{
			Trigger: TriggerEntersPlay,
			Effect:  Stun{Target: Target{Kind: TargetThisCreature}},
		},
	)
	if want := SelfName + " enters play stunned."; stun != want {
		t.Errorf("enters-play stun = %q, want %q", stun, want)
	}
	// A Ready effect renders as the "ready" state word.
	ready := RenderAbility(
		Ability{
			Trigger: TriggerEntersPlay,
			Effect:  Ready{Target: Target{Kind: TargetThisCreature}},
		},
	)
	if want := SelfName + " enters play ready."; ready != want {
		t.Errorf("enters-play ready = %q, want %q", ready, want)
	}
	// An Enrage effect renders as the "enraged" state word.
	enrage := RenderAbility(
		Ability{
			Trigger: TriggerEntersPlay,
			Effect:  Enrage{Target: Target{Kind: TargetThisCreature}},
		},
	)
	if want := SelfName + " enters play enraged."; enrage != want {
		t.Errorf("enters-play enrage = %q, want %q", enrage, want)
	}
	// A Sequence of state effects joins its words with "and" (Gizelhart's Zealot).
	both := RenderAbility(
		Ability{
			Trigger: TriggerEntersPlay,
			Effect: Sequence{Effects: []Effect{
				Ready{Target: Target{Kind: TargetThisCreature}},
				Enrage{Target: Target{Kind: TargetThisCreature}},
			}},
		},
	)
	if want := SelfName + " enters play ready and enraged."; both != want {
		t.Errorf("enters-play ready+enrage = %q, want %q", both, want)
	}
	// An effect without a dedicated enter word falls back to its ordinary text.
	other := RenderAbility(
		Ability{
			Trigger: TriggerEntersPlay,
			Effect: GainAember{
				Player: Controller,
				Amount: 1,
			},
		},
	)
	if want := SelfName + " enters play gain 1 Æmber."; other != want {
		t.Errorf("enters-play fallback = %q, want %q", other, want)
	}
}

func TestTriggerPrefixDefault(t *testing.T) {
	// An unknown trigger renders with no prefix and a capitalized effect.
	got := RenderAbility(
		Ability{
			Trigger: Trigger(99),
			Effect: GainAember{
				Player: Controller,
				Amount: 1,
			},
		},
	)
	if got != "Gain 1 Æmber." {
		t.Errorf("RenderAbility(unknown) = %q", got)
	}
}

func TestEndOfTurnAbilityText(t *testing.T) {
	got := RenderAbility(
		Ability{
			Trigger: TriggerEndOfTurn,
			Effect: LoseAember{
				Player: Opponent,
				Amount: 1,
			},
		},
	)
	if want := "At the end of your turn, your opponent loses 1 Æmber."; got != want {
		t.Errorf("end-of-turn text = %q, want %q", got, want)
	}
}

func TestAfterYouPlayFolding(t *testing.T) {
	// A Conditional{ItIs} folds into the natural "after you play a <shape>" wording.
	folded := RenderAbility(
		Ability{
			Trigger: TriggerAfterCardPlayed,
			Effect: Conditional{
				Cond: ItIs{Type: Artifact},
				Then: StealAember{Amount: 1},
			},
		},
	)
	if want := "After you play an artifact, steal 1 Æmber."; folded != want {
		t.Errorf("folded = %q, want %q", folded, want)
	}
	// A Conditional{ItIsNamed} folds into the natural "after you play <Name>" wording.
	named := RenderAbility(
		Ability{
			Trigger: TriggerAfterCardPlayed,
			Effect: Conditional{
				Cond: ItIsNamed{Name: "Subtle Chain"},
				Then: GainAember{
					Player: Controller,
					Amount: 1,
				},
			},
		},
	)
	if want := "After you play Subtle Chain, gain 1 Æmber."; named != want {
		t.Errorf("named = %q, want %q", named, want)
	}
	// A Conditional{ItIsOfTrait} folds into "after you play a <Trait> creature"
	// (Dark Æmber Vault), so a treachery-gifted Mutant still triggers "you played".
	trait := RenderAbility(
		Ability{
			Trigger: TriggerAfterCardPlayed,
			Effect: Conditional{
				Cond: ItIsOfTrait{Trait: Mutant},
				Then: Draw{Amount: 1},
			},
		},
	)
	if want := "After you play a Mutant creature, draw a card."; trait != want {
		t.Errorf("trait = %q, want %q", trait, want)
	}
	// A non-Conditional reaction keeps the broad prefix.
	plain := RenderAbility(
		Ability{
			Trigger: TriggerAfterCardPlayed,
			Effect: GainAember{
				Player: Controller,
				Amount: 1,
			},
		},
	)
	if want := "After you play a card, gain 1 Æmber."; plain != want {
		t.Errorf("plain = %q, want %q", plain, want)
	}
	// A Conditional on something other than the played card's shape stays literal.
	stateGated := RenderAbility(
		Ability{
			Trigger: TriggerAfterCardPlayed,
			Effect: Conditional{
				Cond: PoolAember{
					Player: Opponent,
					Is:     AtLeast,
					Amount: 1,
				},
				Then: GainAember{
					Player: Controller,
					Amount: 1,
				},
			},
		},
	)
	if want := "After you play a card, if your opponent has 1 Æmber or more, gain 1 Æmber."; stateGated != want {
		t.Errorf("state-gated = %q, want %q", stateGated, want)
	}
	// A Conditional{ItHasBonusIcon} folds into the natural "after you play a card
	// with a bonus icon" wording (Adaptoid).
	bonusIcon := RenderAbility(
		Ability{
			Trigger: TriggerAfterCardPlayed,
			Effect: Conditional{
				Cond: ItHasBonusIcon{},
				Then: GainAember{
					Player: Controller,
					Amount: 1,
				},
			},
		},
	)
	if want := "After you play a card with a bonus icon, gain 1 Æmber."; bonusIcon != want {
		t.Errorf("bonus-icon = %q, want %q", bonusIcon, want)
	}
}

func TestAfterYouUseFolding(t *testing.T) {
	// Veylan Analyst: a Conditional{ItIs} on an AfterUse reaction folds into the
	// natural "after you use a <shape>" wording.
	folded := RenderAbility(
		Ability{
			Trigger: TriggerAfterUse,
			Effect: Conditional{
				Cond: ItIs{Type: Artifact},
				Then: GainAember{
					Player: Controller,
					Amount: 1,
				},
			},
		},
	)
	if want := "After you use an artifact, gain 1 Æmber."; folded != want {
		t.Errorf("folded = %q, want %q", folded, want)
	}
	// A non-Conditional use reaction keeps the broad prefix.
	plain := RenderAbility(
		Ability{
			Trigger: TriggerAfterUse,
			Effect: GainAember{
				Player: Controller,
				Amount: 1,
			},
		},
	)
	if want := "After you use a card, gain 1 Æmber."; plain != want {
		t.Errorf("plain = %q, want %q", plain, want)
	}
}

func TestAfterYouDiscardFolding(t *testing.T) {
	// Baron Mengevin: a Conditional{ItIs} on an AfterDiscardFromHand reaction folds
	// into the natural "after you discard a <shape>" wording, with a house-only ItIs
	// rendering the house noun.
	folded := RenderAbility(
		Ability{
			Trigger: TriggerAfterDiscardFromHand,
			Effect: Conditional{
				Cond: ItIs{House: namedHouse(Sanctum)},
				Then: CaptureAember{
					Target: Target{Kind: TargetThisCreature},
					Amount: 1,
					Source: Opponent,
				},
			},
		},
	)
	if want := "After you discard a Sanctum card, " + SelfName + " captures 1 Æmber from your opponent."; folded != want {
		t.Errorf("folded = %q, want %q", folded, want)
	}
	// A non-Conditional discard reaction keeps the broad prefix.
	plain := RenderAbility(
		Ability{
			Trigger: TriggerAfterDiscardFromHand,
			Effect: GainAember{
				Player: Controller,
				Amount: 1,
			},
		},
	)
	if want := "After you discard a card from your hand, gain 1 Æmber."; plain != want {
		t.Errorf("plain = %q, want %q", plain, want)
	}
}

func TestAfterCreaturePlayedAdjacentFolding(t *testing.T) {
	// Stilt-Kin: a Conditional{ItIsOfTrait} on an AfterCreaturePlayedAdjacent
	// reaction folds the trait into the trigger phrase.
	folded := RenderAbility(
		Ability{
			Trigger: TriggerAfterCreaturePlayedAdjacent,
			Effect: Conditional{
				Cond: ItIsOfTrait{Trait: Giant},
				Then: OnChooseCreature{
					Target: Target{Kind: TargetThisCreature},
					Verbs:  []CreatureVerb{ReadyVerb{}, FightVerb{}},
				},
			},
		},
	)
	if want := "After a Giant creature is played adjacent to " + SelfName +
		", ready and fight with " + SelfName + "."; folded != want {
		t.Errorf("folded = %q, want %q", folded, want)
	}
	// A non-Conditional played-adjacent reaction keeps the broad prefix.
	plain := RenderAbility(
		Ability{
			Trigger: TriggerAfterCreaturePlayedAdjacent,
			Effect:  Draw{Amount: 1},
		},
	)
	if want := "After a creature is played adjacent to " + SelfName +
		", draw a card."; plain != want {
		t.Errorf("plain = %q, want %q", plain, want)
	}
	// A Conditional gated on something other than the played creature's trait stays
	// literal, keeping the broad prefix.
	stateGated := RenderAbility(
		Ability{
			Trigger: TriggerAfterCreaturePlayedAdjacent,
			Effect: Conditional{
				Cond: PoolAember{
					Player: Opponent,
					Is:     AtLeast,
					Amount: 1,
				},
				Then: Draw{Amount: 1},
			},
		},
	)
	if want := "After a creature is played adjacent to " + SelfName +
		", if your opponent has 1 Æmber or more, draw a card."; stateGated != want {
		t.Errorf("state-gated = %q, want %q", stateGated, want)
	}
}

func TestAfterEnemyPlaysCreatureOnFlankFolding(t *testing.T) {
	// Dexus / Sinestra: a Conditional{OnFlank{OfIt}} on an AfterEnemyCardPlayed
	// reaction folds the flank into the trigger phrase and names the creature the
	// position already requires.
	for _, tc := range []struct {
		where FlankPosition
		want  string
	}{
		{RightFlank, "After your opponent plays a creature on their right flank, your opponent loses 1 Æmber."},
		{LeftFlank, "After your opponent plays a creature on their left flank, your opponent loses 1 Æmber."},
	} {
		folded := RenderAbility(
			Ability{
				Trigger: TriggerAfterEnemyCardPlayed,
				Effect: Conditional{
					Cond: OnFlank{
						OfIt:  true,
						Where: tc.where,
					},
					Then: LoseAember{
						Player: Opponent,
						Amount: 1,
					},
				},
			},
		)
		if folded != tc.want {
			t.Errorf("folded = %q, want %q", folded, tc.want)
		}
	}
	// A non-Conditional enemy-play reaction keeps the broad prefix.
	plain := RenderAbility(
		Ability{
			Trigger: TriggerAfterEnemyCardPlayed,
			Effect:  Draw{Amount: 1},
		},
	)
	if want := "After your opponent plays a card, draw a card."; plain != want {
		t.Errorf("plain = %q, want %q", plain, want)
	}
	// A Conditional with an Else, an any-flank position, a source-flank position, or
	// a non-flank condition stays literal with the broad prefix.
	literal := []struct {
		name string
		cond Condition
		els  Effect
	}{
		{"else", OnFlank{
			OfIt:  true,
			Where: RightFlank,
		}, Draw{Amount: 1}},
		{"any-flank", OnFlank{
			OfIt:  true,
			Where: AnyFlank,
		}, nil},
		{"source-flank", OnFlank{
			OfIt:  false,
			Where: RightFlank,
		}, nil},
		{"non-flank", ItIsOfTrait{Trait: Giant}, nil},
	}
	for _, tc := range literal {
		got := RenderAbility(
			Ability{
				Trigger: TriggerAfterEnemyCardPlayed,
				Effect: Conditional{
					Cond: tc.cond,
					Then: LoseAember{
						Player: Opponent,
						Amount: 1,
					},
					Else: tc.els,
				},
			},
		)
		if strings.HasPrefix(got, "After your opponent plays a creature") {
			t.Errorf("%s: got folded %q, want broad prefix", tc.name, got)
		}
	}
}

func TestAfterCreatureScopeFolding(t *testing.T) {
	draw := Draw{Amount: 1}
	enemyTurn := And{Conditions: []Condition{ItIsEnemy{}, ItIsYourTurn{}}}

	// A board-wide reap or destroyed trigger gated only on whose creature acted
	// folds the scope into the trigger phrase.
	for _, tc := range []struct {
		name string
		a    Ability
		want string
	}{
		{
			"enemy reaps",
			Ability{
				Trigger: TriggerAfterCreatureReaps,
				Effect: Conditional{
					Cond: ItIsEnemy{},
					Then: draw,
				},
			},
			"After an enemy creature reaps, draw a card.",
		},
		{
			"friendly reaps",
			Ability{
				Trigger: TriggerAfterCreatureReaps,
				Effect: Conditional{
					Cond: ItIsFriendly{},
					Then: draw,
				},
			},
			"After a friendly creature reaps, draw a card.",
		},
		{
			"friendly destroyed",
			Ability{
				Trigger: TriggerAfterCreatureDestroyed,
				Effect: Conditional{
					Cond: ItIsFriendly{},
					Then: draw,
				},
			},
			"After a friendly creature is destroyed, draw a card.",
		},
		{
			"enemy destroyed during your turn",
			Ability{
				Trigger: TriggerAfterCreatureDestroyed,
				Effect: Conditional{
					Cond: enemyTurn,
					Then: draw,
				},
			},
			"After an enemy creature is destroyed during your turn, draw a card.",
		},
		{
			"friendly fights",
			Ability{
				Trigger: TriggerAfterCreatureFights,
				Effect: Conditional{
					Cond: ItIsFriendly{},
					Then: draw,
				},
			},
			"After a friendly creature is used to fight, draw a card.",
		},
	} {
		if got := RenderAbility(tc.a); got != tc.want {
			t.Errorf("%s: RenderAbility = %q, want %q", tc.name, got, tc.want)
		}
	}

	// Shapes that do not match the fold keep the broad trigger prefix.
	for _, tc := range []struct {
		name   string
		a      Ability
		prefix string
	}{
		{
			"plain reap",
			Ability{
				Trigger: TriggerAfterCreatureReaps,
				Effect:  draw,
			},
			"After a creature reaps, ",
		},
		{
			"reap with else",
			Ability{
				Trigger: TriggerAfterCreatureReaps,
				Effect: Conditional{
					Cond: ItIsEnemy{},
					Then: draw,
					Else: draw,
				},
			},
			"After a creature reaps, ",
		},
		{
			"reap gated on a non-scope condition",
			Ability{
				Trigger: TriggerAfterCreatureReaps,
				Effect: Conditional{
					Cond: ItIsOfTrait{Trait: Giant},
					Then: draw,
				},
			},
			"After a creature reaps, ",
		},
		{
			"plain destroyed",
			Ability{
				Trigger: TriggerAfterCreatureDestroyed,
				Effect:  draw,
			},
			"After a creature is destroyed, ",
		},
		{
			"destroyed with else",
			Ability{
				Trigger: TriggerAfterCreatureDestroyed,
				Effect: Conditional{
					Cond: ItIsFriendly{},
					Then: draw,
					Else: draw,
				},
			},
			"After a creature is destroyed, ",
		},
		{
			"destroyed enemy without your turn",
			Ability{
				Trigger: TriggerAfterCreatureDestroyed,
				Effect: Conditional{
					Cond: And{Conditions: []Condition{ItIsEnemy{}}},
					Then: draw,
				},
			},
			"After a creature is destroyed, ",
		},
		{
			"destroyed and but not enemy-your-turn",
			Ability{
				Trigger: TriggerAfterCreatureDestroyed,
				Effect: Conditional{
					Cond: And{Conditions: []Condition{ItIsFriendly{}, ItIsEnemy{}}},
					Then: draw,
				},
			},
			"After a creature is destroyed, ",
		},
		{
			"plain fight",
			Ability{
				Trigger: TriggerAfterCreatureFights,
				Effect:  draw,
			},
			"After a creature is used to fight, ",
		},
		{
			"fight with else",
			Ability{
				Trigger: TriggerAfterCreatureFights,
				Effect: Conditional{
					Cond: ItIsFriendly{},
					Then: draw,
					Else: draw,
				},
			},
			"After a creature is used to fight, ",
		},
		{
			"fight gated on a non-scope condition",
			Ability{
				Trigger: TriggerAfterCreatureFights,
				Effect: Conditional{
					Cond: ItIsOfTrait{Trait: Giant},
					Then: draw,
				},
			},
			"After a creature is used to fight, ",
		},
	} {
		if got := RenderAbility(tc.a); !strings.HasPrefix(got, tc.prefix) {
			t.Errorf("%s: RenderAbility = %q, want broad prefix %q", tc.name, got, tc.prefix)
		}
	}
}

func TestIsFightReapPair(t *testing.T) {
	ready := Conditional{
		Cond: SourceFirstUseThisTurn{},
		Then: Ready{Target: Target{Kind: TargetThisCreature}},
	}
	reap := Ability{
		Trigger: TriggerAfterReap,
		Effect:  ready,
	}
	fight := Ability{
		Trigger: TriggerAfterFight,
		Effect:  ready,
	}
	if !isFightReapPair(reap, fight) {
		t.Error("Reap+Fight sharing one effect should pair")
	}
	if !isFightReapPair(fight, reap) {
		t.Error("Fight+Reap (reversed order) should also pair")
	}
	if isFightReapPair(reap, Ability{
		Trigger: TriggerAfterReap,
		Effect:  ready,
	}) {
		t.Error("Reap+Reap is not a Fight/Reap pair")
	}
	if isFightReapPair(reap, Ability{
		Trigger: TriggerAfterFight,
		Effect:  StealAember{Amount: 1},
	}) {
		t.Error("differing effects should not pair")
	}
}

func TestTargetTextDefault(t *testing.T) {
	cases := map[TargetKind]string{
		TargetThisCreature:           SelfName,
		TargetTriggeringCreature:     "it",
		TargetTheOtherCreature:       "the other creature",
		TargetTheSameCreature:        "the same creature",
		TargetEachEnemyCreature:      "each enemy creature",
		TargetEachFriendlyCardInPlay: "each friendly card",
		TargetKind(99):               "a creature",
	}
	for kind, want := range cases {
		if got := (Target{Kind: kind}).Text(); got != want {
			t.Errorf("Target{%d}.Text() = %q, want %q", kind, got, want)
		}
	}
}

func TestAllTriggerPrefixes(t *testing.T) {
	cases := map[Trigger]string{
		TriggerAfterPlay:                "Play: Gain 1 Æmber.",
		TriggerAfterReap:                "Reap: Gain 1 Æmber.",
		TriggerAfterFight:               "Fight: Gain 1 Æmber.",
		TriggerBeforeFight:              "Before Fight: Gain 1 Æmber.",
		TriggerAction:                   "Action: Gain 1 Æmber.",
		TriggerDestroyed:                "Destroyed: Gain 1 Æmber.",
		TriggerAfterForgeKey:            "After you forge a key, gain 1 Æmber.",
		TriggerAfterCreatureEnters:      "After a creature enters play, gain 1 Æmber.",
		TriggerAfterDestroyedFighting:   "After a creature is destroyed in a fight with {self}, gain 1 Æmber.",
		TriggerAfterAssaultDestroys:     "After a creature is destroyed by {self}'s assault damage, gain 1 Æmber.",
		TriggerAfterArmorPrevents:       "After {self} prevents damage with its armor, gain 1 Æmber.",
		TriggerAfterNeighborFights:      "After a neighbor of {self} is used to fight, gain 1 Æmber.",
		TriggerAfterCardPlayed:          "After you play a card, gain 1 Æmber.",
		TriggerAfterCreaturePlayed:      "After a creature is played, gain 1 Æmber.",
		TriggerAfterAemberStolenFromYou: "After Æmber is stolen from you, gain 1 Æmber.",
	}
	for tr, want := range cases {
		got := RenderAbility(
			Ability{
				Trigger: tr,
				Effect: GainAember{
					Player: Controller,
					Amount: 1,
				},
			},
		)
		if got != want {
			t.Errorf("trigger %d prefix = %q, want %q", tr, got, want)
		}
	}
}

func TestGeneratedCardText(t *testing.T) {
	cases := []struct {
		def  CardDefinition
		want string
	}{
		{
			exGiant(),
			"House:  Brobnar\nType:   Creature\nRarity: Rare\nPower:  5\nTraits: Giant\n\nAfter you forge a key, deal 2 damage to each enemy creature.",
		},
		{
			exBruteStrength(),
			"House:  Brobnar\nType:   Upgrade\nRarity: Uncommon\nBonus:  Æmber\n\nThis creature gains +5 power.",
		},
		{
			exBattleFury(),
			"House:  Brobnar\nType:   Tactic\nRarity: Common\nBonus:  Æmber\n\nPlay: Ready and fight with a friendly creature.",
		},
		{
			exAutocannon(),
			"House:  Brobnar\nType:   Artifact\nRarity: Rare\nBonus:  Æmber\nTraits: Weapon\n\nAfter a creature enters play, deal 1 damage to it.",
		},
		{
			NewCard(
				"Asp",
				Shadows,
				Creature,
				Uncommon,
				WithPower(3),
				WithKeywords(Skirmish, Poison),
			),
			"House:  Shadows\nType:   Creature\nRarity: Uncommon\nPower:  3\n\nSkirmish, Poison.",
		},
		{
			NewCard(
				"Anaphiel",
				Sanctum,
				Creature,
				Common,
				WithPower(6),
				WithArmor(1),
				WithTraits(Knight),
				WithKeywords(Taunt),
			),
			"House:  Sanctum\nType:   Creature\nRarity: Common\nPower:  6\nArmor:  1\nTraits: Knight\n\nTaunt.",
		},
		{
			NewCard(
				"Camouflage",
				Untamed,
				Upgrade,
				Uncommon,
				WithStatic(StaticModifier{ProtectsFromNonFlank: true}),
			),
			"House:  Untamed\nType:   Upgrade\nRarity: Uncommon\n\nCreatures not on a flank cannot fight this creature.",
		},
		{
			NewCard(
				"Heart of the Forest",
				Untamed,
				Artifact,
				Rare,
				WithRestrictions(Restrictions{NoForgeWhileAheadOnKeys: true}),
			),
			"House:  Untamed\nType:   Artifact\nRarity: Rare\n\nEach player cannot forge keys while they have more forged keys than their opponent.",
		},
		{
			NewCard(
				"Po's Pixies",
				Untamed,
				Creature,
				Rare,
				WithPower(1),
				WithKeywords(Elusive),
				WithReplaces(Instead{
					Of: EventAemberTakenFromPool, Player: Controller, With: FromCommonSupply,
				}),
			),
			"House:  Untamed\nType:   Creature\nRarity: Rare\nPower:  1\n\nElusive.\nÆmber stolen or captured from your pool is taken from the common supply instead.",
		},
		{
			NewCard(
				"Amphora Captura",
				Saurian,
				Artifact,
				Rare,
				WithTraits(Item),
				WithBonusInstead(BonusInstead{
					May: true,
					As:  BonusCapture,
				}),
			),
			"House:  Saurian\nType:   Artifact\nRarity: Rare\nTraits: Item\n\nWhen resolving a bonus icon, you may resolve it as a Capture bonus icon instead.",
		},
		{
			NewCard(
				"Scrivener Favian",
				Sanctum,
				Creature,
				Uncommon,
				WithPower(3),
				WithTraits(Mutant),
				WithBonusInstead(BonusInstead{
					From:    BonusCapture,
					Instead: StealAember{Amount: 1},
				}),
			),
			"House:  Sanctum\nType:   Creature\nRarity: Uncommon\nPower:  3\nTraits: Mutant\n\nWhen you resolve a Capture bonus icon, steal 1 Æmber instead.",
		},
		{
			NewCard(
				"Tabris",
				Sanctum,
				Creature,
				Uncommon,
				WithPower(6),
				WithAbility(
					TriggerAfterFight,
					CaptureAember{
						Amount: 1,
						Target: Target{Kind: TargetThisCreature},
						Source: Opponent,
					},
				),
			),
			"House:  Sanctum\nType:   Creature\nRarity: Uncommon\nPower:  6\n\nFight: Tabris captures 1 Æmber from your opponent.",
		},
		{
			NewCard(
				"Bear",
				Untamed,
				Creature,
				Common,
				WithPower(5),
				WithTraits(Beast),
				WithAssault(2),
			),
			"House:  Untamed\nType:   Creature\nRarity: Common\nPower:  5\nTraits: Beast\n\nAssault 2.",
		},
		{
			NewCard("Grub", Untamed, Creature, Rare, WithPower(2), WithHazardous(5)),
			"House:  Untamed\nType:   Creature\nRarity: Rare\nPower:  2\n\nHazardous 5.",
		},
		{
			NewCard("Cowfyne", Brobnar, Creature, Common, WithPower(5), WithSplashAttack(2)),
			"House:  Brobnar\nType:   Creature\nRarity: Common\nPower:  5\n\nSplash-attack 2.",
		},
		{
			NewCard(
				"Valdr",
				Brobnar,
				Creature,
				Common,
				WithPower(6),
				WithTraits(Giant),
				WithAttackDamage(AttackDamage{
					Amount:    2,
					FlankOnly: true,
				}),
			),
			"House:  Brobnar\nType:   Creature\nRarity: Common\nPower:  6\nTraits: Giant\n\nValdr deals +2 damage while attacking an enemy creature on the flank.",
		},
		{
			NewCard(
				"Spider",
				Mars,
				Creature,
				Common,
				WithPower(7),
				WithAttackDamage(AttackDamage{
					Fixed:  true,
					Amount: 0,
				}),
			),
			"House:  Mars\nType:   Creature\nRarity: Common\nPower:  7\n\nSpider deals no damage when fighting.",
		},
		{
			NewCard(
				"Ether Spider",
				Mars,
				Creature,
				Uncommon,
				WithPower(7),
				WithAttackDamage(AttackDamage{
					Fixed:  true,
					Amount: 0,
				}),
				WithReplaces(Instead{
					Of:     EventAemberAddedToPool,
					Player: Opponent,
					With:   Capture,
				}),
			),
			"House:  Mars\nType:   Creature\nRarity: Uncommon\nPower:  7\n\nEther Spider deals no damage when fighting.\nIf Æmber would be added to your opponent's pool, instead Ether Spider captures it.",
		},
		{
			NewCard(
				"Bruiser",
				Brobnar,
				Creature,
				Common,
				WithPower(8),
				WithAttackDamage(AttackDamage{
					Fixed:  true,
					Amount: 5,
				}),
			),
			"House:  Brobnar\nType:   Creature\nRarity: Common\nPower:  8\n\nBruiser deals 5 damage when fighting.",
		},
		{
			NewCard(
				"Basher",
				Brobnar,
				Creature,
				Common,
				WithPower(4),
				WithAttackDamage(AttackDamage{Amount: 2}),
			),
			"House:  Brobnar\nType:   Creature\nRarity: Common\nPower:  4\n\nBasher deals +2 damage when fighting.",
		},
		{
			NewCard(
				"Runner",
				Shadows,
				Upgrade,
				Uncommon,
				WithStatic(
					StaticModifier{
						Granted: []Ability{
							{Trigger: TriggerAfterReap, Effect: StealAember{Amount: 1}},
						},
					},
				),
			),
			"House:  Shadows\nType:   Upgrade\nRarity: Uncommon\n\nThis creature gains, \"Reap: Steal 1 Æmber.\"",
		},
		{
			NewCard(
				"Boots",
				Logos,
				Upgrade,
				Uncommon,
				WithStatic(
					StaticModifier{
						Granted: []Ability{
							{
								Trigger: TriggerAfterReap,
								Effect: Conditional{
									Cond: SourceFirstUseThisTurn{},
									Then: Ready{Target: Target{Kind: TargetThisCreature}},
								},
							},
							{
								Trigger: TriggerAfterFight,
								Effect: Conditional{
									Cond: SourceFirstUseThisTurn{},
									Then: Ready{Target: Target{Kind: TargetThisCreature}},
								},
							},
						},
					},
				),
			),
			"House:  Logos\nType:   Upgrade\nRarity: Uncommon\n\nThis creature gains, \"Fight/Reap: If this is the first time this creature has been used this turn, ready it.\"",
		},
		{
			NewCard(
				"Boots",
				Logos,
				Upgrade,
				Uncommon,
				WithStatic(StaticModifier{Keywords: []Keyword{Versatile}}),
				WithAbility(
					TriggerAfterPlay,
					Sequence{
						Effects: []Effect{
							Stun{Target: Target{Kind: TargetThisCreature}},
							Exhaust{Target: Target{Kind: TargetThisCreature}},
						},
					},
				),
			),
			"House:  Logos\nType:   Upgrade\nRarity: Uncommon\n\nThis creature gains versatile.\nPlay: Stun and exhaust this creature.",
		},
		{
			NewCard(
				"Jammer",
				Mars,
				Creature,
				Common,
				WithPower(4),
				WithArmor(1),
				WithTraits(Robot),
				WithKeyCost(NewKeyCostChange(Opponent, 1)),
				WithAbility(
					TriggerAfterReap,
					CaptureAember{
						Amount: 1,
						Target: Target{Kind: TargetThisCreature},
						Source: Opponent,
					},
				),
				WithAbility(
					TriggerAfterFight,
					CaptureAember{
						Amount: 1,
						Target: Target{Kind: TargetThisCreature},
						Source: Opponent,
					},
				),
			),
			"House:  Mars\nType:   Creature\nRarity: Common\nPower:  4\nArmor:  1\nTraits: Robot\n\nYour opponent's keys cost +1 Æmber.\nFight/Reap: Jammer captures 1 Æmber from your opponent.",
		},
		{
			NewCard(
				"Pack",
				Mars,
				Upgrade,
				Uncommon,
				WithBonus(BonusAember),
				WithStatic(StaticModifier{KeyCostChange: NewKeyCostChange(Opponent, 2)}),
			),
			"House:  Mars\nType:   Upgrade\nRarity: Uncommon\nBonus:  Æmber\n\nThis creature gains, \"Your opponent's keys cost +2 Æmber.\"",
		},
		{
			NewCard(
				"Shield",
				Sanctum,
				Upgrade,
				Rare,
				WithStatic(
					StaticModifier{
						Replaces: Replace{
							When: EventCreatureDestroyed,
							With: Sequence{
								Effects: []Effect{
									Heal{
										Fully:  true,
										Target: Target{Kind: TargetTriggeringCreature},
									},
									Destroy{Target: Target{Kind: TargetThisCreature}},
								},
							},
						},
					},
				),
			),
			"House:  Sanctum\nType:   Upgrade\nRarity: Rare\n\nThis creature gains, \"If this creature would be destroyed, instead fully heal it, and destroy Shield.\"",
		},
		{
			NewCard(
				"Cloak",
				Sanctum,
				Upgrade,
				Rare,
				WithStatic(
					StaticModifier{
						HazardousBonus: 2,
						Replaces: Replace{
							When: EventCreatureDestroyed,
							With: Sequence{
								Effects: []Effect{
									Heal{
										Fully:  true,
										Target: Target{Kind: TargetTriggeringCreature},
									},
									Destroy{Target: Target{Kind: TargetThisCreature}},
								},
							},
						},
					},
				),
			),
			"House:  Sanctum\nType:   Upgrade\nRarity: Rare\n\nThis creature gains +2 hazardous and, \"If this creature would be destroyed, instead fully heal it, and destroy Cloak.\"",
		},
		{
			NewCard(
				"Antenna",
				Mars,
				Upgrade,
				Rare,
				WithStatic(
					StaticModifier{
						Granted: []Ability{
							{
								Trigger: TriggerAfterCardPlayed,
								Effect: Conditional{
									Cond: ItIs{
										House: namedHouse(Mars),
										Type:  Creature,
									},
									Then: Sequence{
										Effects: []Effect{
											Ready{Target: Target{Kind: TargetThisCreature}},
											BelongToHouse{
												Target:   Target{Kind: TargetThisCreature},
												House:    Mars,
												Duration: RemainderOfPlayerTurn,
											},
										},
									},
								},
							},
						},
					},
				),
			),
			"House:  Mars\nType:   Upgrade\nRarity: Rare\n\nThis creature gains, \"After you play a Mars creature, ready this creature. For the remainder of the turn, this creature belongs to house Mars.\"",
		},
		{
			NewCard(
				"SelfTax",
				Mars,
				Creature,
				Common,
				WithPower(3),
				WithKeyCost(NewKeyCostChange(Controller, 1)),
			),
			"House:  Mars\nType:   Creature\nRarity: Common\nPower:  3\n\nYour keys cost +1 Æmber.",
		},
		{
			NewCard(
				"Tax",
				Mars,
				Creature,
				Common,
				WithPower(3),
				WithKeyCost(NewKeyCostChange(EachPlayer, 1)),
			),
			"House:  Mars\nType:   Creature\nRarity: Common\nPower:  3\n\nEach player's keys cost +1 Æmber.",
		},
		{
			NewCard(
				"Imp",
				Dis,
				Creature,
				Common,
				WithPower(2),
				WithTraits(Imp),
				WithRestrictions(
					Restrictions{PlayCardLimit: PlayCardLimit{
						Player: Opponent,
						Amount: 2,
					}},
				),
			),
			"House:  Dis\nType:   Creature\nRarity: Common\nPower:  2\nTraits: Imp\n\nYour opponent cannot play more than 2 cards each turn.",
		},
		{
			NewCard(
				"Witch",
				Untamed,
				Creature,
				Rare,
				WithPower(4),
				WithPlayPermission(PlayPermission{
					House:  Untamed,
					Amount: 1,
				}),
			),
			"House:  Untamed\nType:   Creature\nRarity: Rare\nPower:  4\n\nEach turn you may play one Untamed card.",
		},
		{
			NewCard(
				"Witch2",
				Untamed,
				Creature,
				Rare,
				WithPower(4),
				WithPlayPermission(PlayPermission{
					House:  Untamed,
					Amount: 2,
				}),
			),
			"House:  Untamed\nType:   Creature\nRarity: Rare\nPower:  4\n\nEach turn you may play 2 Untamed cards.",
		},
		{
			NewCard(
				"Twig",
				Untamed,
				Creature,
				Common,
				WithPower(7),
				WithTraits(Beast),
				WithFightRestriction(Target{Kind: TargetEachCreature}.Stunned()),
			),
			"House:  Untamed\nType:   Creature\nRarity: Common\nPower:  7\nTraits: Beast\n\nTwig can only fight stunned creatures.",
		},
		{
			NewCard(
				"Ritual",
				Dis,
				Artifact,
				Rare,
				WithTraits(Power),
				WithConstantAbility(
					ConstantAbility{
						Target: Target{Kind: TargetEachCreature},
						Granted: []Ability{
							{
								Trigger: TriggerDestroyed,
								Effect:  PurgeCreature{Target: Target{Kind: TargetThisCreature}},
							},
						},
					},
				),
			),
			"House:  Dis\nType:   Artifact\nRarity: Rare\nTraits: Power\n\nEach creature gains, \"Destroyed: Purge this creature.\"",
		},
	}
	for _, tc := range cases {
		if got := RenderCardText(&tc.def); got != tc.want {
			t.Errorf("%s text mismatch:\n got:\n%s\nwant:\n%s", tc.def.Name, got, tc.want)
		}
	}
}

func TestRenderCardDetail(t *testing.T) {
	def := NewCard("Dr. Escotera", Logos, Creature, Rare, WithPower(4))
	want := "Name:   Dr. Escotera\nHouse:  Logos\nType:   Creature\nRarity: Rare\nPower:  4"
	if got := RenderCardDetail(&def); got != want {
		t.Errorf("detail mismatch:\n got:\n%s\nwant:\n%s", got, want)
	}
}

func TestRenderCardRules(t *testing.T) {
	cases := []struct {
		def  CardDefinition
		want string
	}{
		// Vanilla creature: no rules text.
		{NewCard("Grunt", Brobnar, Creature, Common, WithPower(4)), ""},
		// Keyword line only.
		{
			NewCard(
				"Asp",
				Shadows,
				Creature,
				Uncommon,
				WithPower(3),
				WithKeywords(Skirmish, Poison),
			),
			"Skirmish, Poison.",
		},
		// Upgrade static modifier (no own ability): the face still shows its text.
		{exBruteStrength(), "This creature gains +5 power."},
		// Triggered ability with a self-reference resolved to the card's name.
		{
			NewCard(
				"Tabris",
				Sanctum,
				Creature,
				Uncommon,
				WithPower(6),
				WithAbility(
					TriggerAfterFight,
					CaptureAember{
						Amount: 1,
						Target: Target{Kind: TargetThisCreature},
						Source: Opponent,
					},
				),
			),
			"Fight: Tabris captures 1 Æmber from your opponent.",
		},
		// A card that takes damage for other creatures renders that shield line.
		{
			NewCard(
				"Ward",
				Sanctum,
				Creature,
				Common,
				WithPower(5),
				WithTakesDamageFor(Target{Kind: TargetEachCreature}.Neighboring()),
			),
			"Damage dealt to each neighboring creature is dealt to Ward instead.",
		},
		// A card that shares its neighbors' fight damage renders that line.
		{
			NewCard(
				"Drecker",
				Dis,
				Creature,
				Common,
				WithPower(4),
				WithAlsoTakesNeighborFightDamage(),
			),
			"Damage dealt to Drecker's neighbors during fights is also dealt to Drecker.",
		},
		// A card that gains a keyword only while attacking renders that clause.
		{
			NewCard(
				"Spyyyder",
				Shadows,
				Creature,
				Common,
				WithPower(4),
				WithAttackKeywords(AttackKeywords{Keywords: []Keyword{Poison}}),
			),
			"Spyyyder gains poison while attacking.",
		},
	}
	for _, tc := range cases {
		if got := RenderCardRules(&tc.def); got != tc.want {
			t.Errorf("%s rules mismatch:\n got:  %q\n want: %q", tc.def.Name, got, tc.want)
		}
	}
}

func TestRenderUpgradeOnCreature(t *testing.T) {
	cases := []struct {
		def  CardDefinition
		want string
	}{
		// A static bonus reads as the bonus itself; the creature is already there.
		{exBruteStrength(), "+5 power."},
		// A granted ability loses the "This creature gains" quoting.
		{
			NewCard(
				"Runner",
				Shadows,
				Upgrade,
				Uncommon,
				WithStatic(
					StaticModifier{
						Granted: []Ability{
							{Trigger: TriggerAfterReap, Effect: StealAember{Amount: 1}},
						},
					},
				),
			),
			"Reap: Steal 1 Æmber.",
		},
		// A Fight/Reap pair keeps its shorthand, unquoted.
		{
			NewCard(
				"Boots",
				Logos,
				Upgrade,
				Uncommon,
				WithStatic(
					StaticModifier{
						Granted: []Ability{
							{
								Trigger: TriggerAfterReap,
								Effect: Conditional{
									Cond: SourceFirstUseThisTurn{},
									Then: Ready{Target: Target{Kind: TargetThisCreature}},
								},
							},
							{
								Trigger: TriggerAfterFight,
								Effect: Conditional{
									Cond: SourceFirstUseThisTurn{},
									Then: Ready{Target: Target{Kind: TargetThisCreature}},
								},
							},
						},
					},
				),
			),
			"Fight/Reap: If this is the first time this creature has been used this turn, ready it.",
		},
		// A granted key-cost change too.
		{
			NewCard(
				"Pack",
				Mars,
				Upgrade,
				Uncommon,
				WithStatic(StaticModifier{KeyCostChange: NewKeyCostChange(Opponent, 2)}),
			),
			"Your opponent's keys cost +2 Æmber.",
		},
		// A bonus and a destruction replacement each stand on their own line.
		{
			NewCard(
				"Cloak",
				Sanctum,
				Upgrade,
				Rare,
				WithStatic(
					StaticModifier{
						HazardousBonus: 2,
						Replaces: Replace{
							When: EventCreatureDestroyed,
							With: Destroy{Target: Target{Kind: TargetThisCreature}},
						},
					},
				),
			),
			"+2 hazardous.\nIf this creature would be destroyed, instead destroy Cloak.",
		},
		// Non-flank fight protection reads as the rule itself, hosted on the creature.
		{
			NewCard(
				"Camouflage",
				Untamed,
				Upgrade,
				Uncommon,
				WithStatic(StaticModifier{ProtectsFromNonFlank: true}),
			),
			"Creatures not on a flank cannot fight this creature.",
		},
		// A granted Æmber-protection reads as the rule itself, hosted on the creature.
		{
			NewCard(
				"Guard",
				Sanctum,
				Upgrade,
				Uncommon,
				WithStatic(StaticModifier{AemberCannotBeStolen: AlwaysMet{}}),
			),
			"Your Æmber cannot be stolen.",
		},
	}
	for _, tc := range cases {
		if got := RenderUpgradeOnCreature(&tc.def); got != tc.want {
			t.Errorf("%s hosted rules mismatch:\n got:  %q\n want: %q", tc.def.Name, got, tc.want)
		}
	}
}

// TestUpgradeGrantLinesHouseOverride covers a creature-as-upgrade's house-override
// grant line (Academy Training makes its host a Logos creature).
func TestUpgradeGrantLinesHouseOverride(t *testing.T) {
	def := NewCard("Academy Training", Logos, Creature, Uncommon, WithPower(1),
		WithStatic(StaticModifier{HouseOverride: Logos}), WithPlayableAsUpgrade())
	lines := upgradeGrantLines(&def, true)
	want := "This creature belongs to Logos"
	if len(lines) == 0 || lines[0] != want {
		t.Errorf("grant lines = %v, want first %q", lines, want)
	}
}

// TestUpgradeGrantLinesHouseOverrideWithGrant covers a house-override line that
// folds in a granted ability, so it reads as one sentence (Academy Training's
// `belongs to Logos and gains "Reap: Draw a card."`).
func TestUpgradeGrantLinesHouseOverrideWithGrant(t *testing.T) {
	def := NewCard("Academy Training", Logos, Upgrade, Rare,
		WithStatic(StaticModifier{
			HouseOverride: Logos,
			Granted:       []Ability{{Trigger: TriggerAfterReap, Effect: Draw{Amount: 1}}},
		}))
	lines := upgradeGrantLines(&def, true)
	want := `This creature belongs to Logos and gains "Reap: Draw a card."`
	if len(lines) == 0 || lines[0] != want {
		t.Errorf("grant lines = %v, want first %q", lines, want)
	}
}

// TestCardRulesHouseOverride covers an Upgrade card's house-override line rendered
// standalone from its own rules (Academy Training's printed text).
func TestCardRulesHouseOverride(t *testing.T) {
	def := NewCard("Academy Training", Logos, Upgrade, Rare,
		WithStatic(StaticModifier{HouseOverride: Logos}))
	got := RenderCardRules(&def)
	want := "This creature belongs to Logos"
	if !strings.Contains(got, want) {
		t.Errorf("card rules missing the house-override line:\n%s", got)
	}
}

func TestCardDocComment(t *testing.T) {
	cases := []struct {
		def  CardDefinition
		want string
	}{
		{
			exGiant(),
			"// Brobnar Giant\n//\n//\tHouse:  Brobnar\n//\tType:   Creature\n//\tRarity: Rare\n//\tPower:  5\n//\tTraits: Giant\n//\n//\tAfter you forge a key, deal 2 damage to each enemy creature.",
		},
		{
			NewCard(
				"Asp",
				Shadows,
				Creature,
				Uncommon,
				WithPower(3),
				WithKeywords(Skirmish, Poison),
			),
			"// Asp\n//\n//\tHouse:  Shadows\n//\tType:   Creature\n//\tRarity: Uncommon\n//\tPower:  3\n//\n//\tSkirmish, Poison.",
		},
	}
	for _, tc := range cases {
		if got := CardDocComment(&tc.def); got != tc.want {
			t.Errorf("%s doc comment mismatch:\n got:\n%s\nwant:\n%s", tc.def.Name, got, tc.want)
		}
	}
}

func TestCapitalizeFirst(t *testing.T) {
	if capitalizeFirst("") != "" {
		t.Error("capitalizeFirst(\"\") should be empty")
	}
	if capitalizeFirst("hello") != "Hello" {
		t.Errorf("capitalizeFirst(hello) = %q", capitalizeFirst("hello"))
	}
}

func TestLowerFirst(t *testing.T) {
	if lowerFirst("") != "" {
		t.Error("lowerFirst(\"\") should be empty")
	}
	if lowerFirst("Hello") != "hello" {
		t.Errorf("lowerFirst(Hello) = %q", lowerFirst("Hello"))
	}
}

func TestWithGiganticRole(t *testing.T) {
	def := NewCard(
		"Half", Brobnar, Creature, Common, WithPower(6), WithGiganticRole(GiganticBase),
	)
	if def.GiganticRole != GiganticBase {
		t.Errorf("GiganticRole = %d, want %d", def.GiganticRole, GiganticBase)
	}
}

// A gigantic reads "Gigantic Creature" on its type line for both halves, so
// either card in hand shows it is gigantic; an ordinary creature reads plainly.
func TestCardTypeLabelGigantic(t *testing.T) {
	base := NewCard(
		"Colossus", Brobnar, Creature, Common, WithPower(9), WithGiganticRole(GiganticBase),
	)
	art := NewCard(
		"Colossus", Brobnar, Creature, Common, WithPower(9), WithGiganticRole(GiganticArt),
	)
	plain := NewCard("Goon", Brobnar, Creature, Common, WithPower(3))

	if got := CardTypeLabel(&base); got != "Gigantic Creature" {
		t.Errorf("base CardTypeLabel = %q, want %q", got, "Gigantic Creature")
	}
	if got := CardTypeLabel(&art); got != "Gigantic Creature" {
		t.Errorf("art CardTypeLabel = %q, want %q", got, "Gigantic Creature")
	}
	if got := CardTypeLabel(&plain); got != "Creature" {
		t.Errorf("plain CardTypeLabel = %q, want %q", got, "Creature")
	}
	if !strings.Contains(RenderCardText(&base), "Gigantic Creature") {
		t.Errorf("RenderCardText omits gigantic type line:\n%s", RenderCardText(&base))
	}
}

func TestIndefinite(t *testing.T) {
	cases := map[string]string{
		"":                 "",
		"Urchin":           "an Urchin",
		"elf":              "an elf",
		"Knight":           "a Knight",
		"creature":         "a creature",
		"another creature": "another creature",
	}
	for in, want := range cases {
		if got := indefinite(in); got != want {
			t.Errorf("indefinite(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestStaticText(t *testing.T) {
	if staticText(StaticModifier{}) != "" {
		t.Error("empty static modifier should render empty")
	}
	got := staticText(StaticModifier{
		PowerBonus: 5,
		ArmorBonus: 2,
	})
	if got != "This creature gains +5 power and +2 armor." {
		t.Errorf("staticText = %q", got)
	}
	if got := staticText(
		StaticModifier{AssaultBonus: 2},
	); got != "This creature gains +2 assault." {
		t.Errorf("assault staticText = %q", got)
	}
	if got := staticText(
		StaticModifier{
			PowerBonus:     2,
			HazardousBonus: 2,
		},
	); got != "This creature gains +2 power and +2 hazardous." {
		t.Errorf("hazardous staticText = %q", got)
	}
	if got := staticText(
		StaticModifier{SplashAttackBonus: 3},
	); got != "This creature gains +3 splash-attack." {
		t.Errorf("splash-attack staticText = %q", got)
	}
	if got := staticText(
		StaticModifier{Keywords: []Keyword{Skirmish}},
	); got != "This creature gains skirmish." {
		t.Errorf("keyword staticText = %q", got)
	}
	if got := staticText(
		StaticModifier{Keywords: []Keyword{Elusive, Skirmish}},
	); got != "This creature gains elusive and skirmish." {
		t.Errorf("two-keyword staticText = %q", got)
	}
	if got := staticText(
		StaticModifier{
			ArmorBonus: 1,
			Keywords:   []Keyword{Taunt},
		},
	); got != "This creature gains +1 armor and taunt." {
		t.Errorf("armor+keyword staticText = %q", got)
	}
	if got := staticText(
		StaticModifier{
			PowerBonus: 1,
			ArmorBonus: 1,
			Per:        UpgradesOnIt,
		},
	); got != "This creature gains +1 power and +1 armor for each upgrade attached to it." {
		t.Errorf("per-upgrade staticText = %q", got)
	}
	if got := staticText(
		StaticModifier{
			KeywordGrants: []KeywordGrant{
				{Keywords: []Keyword{Elusive}, Host: true, Neighbors: true},
			},
		},
	); got != "This creature and each of its neighbors gains elusive." {
		t.Errorf("neighbor-keyword staticText = %q", got)
	}
	if got := staticText(
		StaticModifier{KeywordGrants: []KeywordGrant{{Keywords: []Keyword{Elusive}, Host: true}}},
	); got != "This creature gains elusive." {
		t.Errorf("host-only keyword-grant staticText = %q", got)
	}
	if got := staticText(
		StaticModifier{
			KeywordGrants: []KeywordGrant{{Keywords: []Keyword{Elusive}, Neighbors: true}},
		},
	); got != "Each of this creature's neighbors gains elusive." {
		t.Errorf("neighbors-only keyword-grant staticText = %q", got)
	}
	if got := staticText(
		StaticModifier{KeywordGrants: []KeywordGrant{{Keywords: []Keyword{Elusive}}}},
	); got != "" {
		t.Errorf("reachless keyword-grant staticText = %q", got)
	}
	if got := staticText(
		StaticModifier{KeywordGrants: []KeywordGrant{{Host: true, Neighbors: true}}},
	); got != "" {
		t.Errorf("keywordless keyword-grant staticText = %q", got)
	}
	if got := staticText(
		StaticModifier{
			PowerBonus: 1,
			KeywordGrants: []KeywordGrant{
				{Keywords: []Keyword{Elusive}, Host: true, Neighbors: true},
			},
		},
	); got != "This creature gains +1 power. This creature and each of its neighbors gains elusive." {
		t.Errorf("bonus+neighbor-keyword staticText = %q", got)
	}
	dongle := NewCard(
		"Cloaking Dongle",
		StarAlliance,
		Upgrade,
		Common,
		WithStatic(
			StaticModifier{
				KeywordGrants: []KeywordGrant{
					{Keywords: []Keyword{Elusive}, Host: true, Neighbors: true},
				},
			},
		),
	)
	if got := upgradeStaticLines(&dongle, true); len(got) != 1 ||
		got[0] != "This creature and each of its neighbors gains elusive." {
		t.Errorf("hosted neighbor-keyword lines = %v", got)
	}
}

func TestPunctuate(t *testing.T) {
	cases := map[string]string{
		`gains, "Play: x."`: `gains, "Play: x."`, // quoted, already has inner period
		`gains, "Play: x"`:  `gains, "Play: x."`, // quoted, period tucked inside
		"deal damage.":      "deal damage.",      // already ends with a period
		"deal damage":       "deal damage.",      // period appended
	}
	for in, want := range cases {
		if got := punctuate(in); got != want {
			t.Errorf("punctuate(%q) = %q, want %q", in, got, want)
		}
	}
}

// TestNewTriggerPrefixes covers the printed prefix of each persistent reaction
// trigger added for Teliga, Veylan Analyst, Rock-Hurling Giant, and Magda the Rat.
func TestNewTriggerPrefixes(t *testing.T) {
	for trigger, want := range map[Trigger]string{
		TriggerAfterEnemyCardPlayed:   "After your opponent plays a card, ",
		TriggerAfterUse:               "After you use a card, ",
		TriggerAfterDiscardFromHand:   "After you discard a card from your hand, ",
		TriggerAfterUsedSelf:          "After " + SelfName + " is used, ",
		TriggerAfterCreatureReaps:     "After a creature reaps, ",
		TriggerAfterCreatureFights:    "After a creature is used to fight, ",
		TriggerAfterCreatureDestroyed: "After a creature is destroyed, ",
		TriggerLeavesPlay:             "Leaves Play: ",
	} {
		if got, _ := trigger.prefix(); got != want {
			t.Errorf("prefix(%v) = %q, want %q", trigger, got, want)
		}
	}
}

// TestTurnBoundaryTriggerPrefixes covers the two triggers that bracket a turn.
// Start of turn is its own phase (ADR 0012), and its printed prefix mirrors the
// end-of-turn one.
func TestTurnBoundaryTriggerPrefixes(t *testing.T) {
	for trigger, want := range map[Trigger]string{
		TriggerStartOfTurn: "At the start of your turn, ",
		TriggerEndOfTurn:   "At the end of your turn, ",
	} {
		if got, _ := trigger.prefix(); got != want {
			t.Errorf("prefix(%v) = %q, want %q", trigger, got, want)
		}
	}
}

// TestTauntReachesNeighborsNeighborsText covers the rules line for a creature whose
// taunt extends one step further — Lady Loreena's mechanic.
func TestTauntReachesNeighborsNeighborsText(t *testing.T) {
	def := NewCard("Lady Loreena", Sanctum, Creature, Rare, WithPower(4),
		WithTauntReachingNeighborsNeighbors())
	want := "Lady Loreena's taunt also applies to its neighbors' neighbors."
	found := false
	for _, line := range cardRules(&def, false) {
		if line == want {
			found = true
		}
	}
	if !found {
		t.Errorf("cardRules missing %q; got %v", want, cardRules(&def, false))
	}
}

// TestGrantedAemberCannotBeStolenText covers the granted line for an unconditional
// AemberCannotBeStolen (AlwaysMet) on a static modifier.
func TestGrantedAemberCannotBeStolenText(t *testing.T) {
	m := StaticModifier{AemberCannotBeStolen: AlwaysMet{}}
	frame := func(body string) string { return body }
	want := "Your \u00c6mber cannot be stolen."
	found := false
	for _, line := range grantedLines(m, "Upgrade", frame) {
		if line == want {
			found = true
		}
	}
	if !found {
		t.Errorf("grantedLines missing %q; got %v", want, grantedLines(m, "Upgrade", frame))
	}

	// A conditional protection renders the "While … , your Æmber cannot be stolen."
	// form instead.
	cond := StaticModifier{AemberCannotBeStolen: HasAember{Subject: This}}
	wantCond := "While " + SelfName + " has \u00c6mber on it, your \u00c6mber cannot be stolen."
	foundCond := false
	for _, line := range grantedLines(cond, "Upgrade", frame) {
		if line == wantCond {
			foundCond = true
		}
	}
	if !foundCond {
		t.Errorf("grantedLines missing %q; got %v", wantCond, grantedLines(cond, "Upgrade", frame))
	}
}

// TestGrantedAbilityCapitalizesSelfReference pins that a granted body opening on
// the self-reference prints capitalized like every other granted body. The
// substitution runs after RenderAbility's capitalization and "{self}" is not a
// letter, so Wild Spirit and Operations Officer Yshi used to print a lowercase
// "this creature captures …" where Observe-u-Max printed a capital one.
func TestGrantedAbilityCapitalizesSelfReference(t *testing.T) {
	opens := Ability{
		Trigger: TriggerAfterReap,
		Effect: CaptureAember{
			Amount: 1,
			Source: Opponent,
			Target: Target{Kind: TargetThisCreature},
		},
	}
	want := "Reap: This creature captures 1 \u00c6mber from your opponent."
	if got := grantedAbilityText(opens, "Horn"); got != want {
		t.Errorf("granted text = %q, want %q", got, want)
	}

	// A body with no trigger prefix capitalizes its first letter instead.
	noPrefix := Ability{
		Trigger: TriggerEntersPlay,
		Effect:  Draw{Amount: 1},
	}
	if got := grantedAbilityText(noPrefix, "Horn"); got != capitalizeFirst(got) {
		t.Errorf("prefixless granted text = %q, want it capitalized", got)
	}
}

// TestBeforeFightTargetReadsInPresentTense pins that a reference to the fought
// creature follows its trigger's tense, not the neighbor decoration that used to
// stand in for it: a Before Fight: ability resolves before the fight, so Siren
// Horn moves Æmber "to the creature it fights", while a Fight: ability keeps the
// past.
func TestBeforeFightTargetReadsInPresentTense(t *testing.T) {
	fought := Target{Kind: TargetCreatureFought}
	before := Ability{
		Trigger: TriggerBeforeFight,
		Effect:  Stun{Target: fought},
	}
	want := "Before Fight: Stun the creature " + SelfName + " fights."
	if got := RenderAbility(before); got != want {
		t.Errorf("before-fight ability = %q, want %q", got, want)
	}

	after := Ability{
		Trigger: TriggerAfterFight,
		Effect:  Stun{Target: fought},
	}
	wantAfter := "Fight: Stun the creature " + SelfName + " fought."
	if got := RenderAbility(after); got != wantAfter {
		t.Errorf("fight ability = %q, want %q", got, wantAfter)
	}
}

// TestCollapsesRepeatedSubject pins that a sentence names a card once and then
// refers back to it, and pins the three cases that must not collapse: a possessive
// mention (Pain Reaction), a mention in a later sentence, and a mention inside a
// quoted ability, which is its own sentence and may not borrow the outer subject.
func TestCollapsesRepeatedSubject(t *testing.T) {
	cases := []struct{ in, want string }{
		{
			"If there was any Æmber on that creature, take control of it, and that creature belongs to house Shadows.",
			"If there was any Æmber on that creature, take control of it, and it belongs to house Shadows.",
		},
		{
			"After this creature is used, destroy this creature.",
			"After this creature is used, destroy it.",
		},
		{
			"If this damage destroys that creature, deal 2 damage to that creature's neighbors.",
			"If this damage destroys that creature, deal 2 damage to that creature's neighbors.",
		},
		{
			"Ready this creature. Destroy this creature.",
			"Ready this creature. Destroy this creature.",
		},
		{
			`This creature gains, "Reap: Ward this creature."`,
			`This creature gains, "Reap: Ward this creature."`,
		},
	}
	for _, c := range cases {
		if got := collapseRepeatedSubject(c.in); got != c.want {
			t.Errorf("collapseRepeatedSubject(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
