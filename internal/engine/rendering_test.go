package engine

import "testing"

func TestEntersStunnedText(t *testing.T) {
	def := NewCard("Chuff", Mars, Creature, Common, WithPower(3), WithEntersPlay(Stun{Target: Target{Kind: TargetThisCreature}}))
	want := "Chuff enters play stunned."
	found := false
	for _, line := range cardRules(&def) {
		if line == want {
			found = true
		}
	}
	if !found {
		t.Errorf("cardRules missing %q; got %v", want, cardRules(&def))
	}
}

func TestEntersPlayAbilityText(t *testing.T) {
	// A Stun effect renders as the "stunned" state word.
	stun := RenderAbility(Ability{Trigger: TriggerEntersPlay, Effect: Stun{Target: Target{Kind: TargetThisCreature}}})
	if want := SelfName + " enters play stunned."; stun != want {
		t.Errorf("enters-play stun = %q, want %q", stun, want)
	}
	// A Ready effect renders as the "ready" state word.
	ready := RenderAbility(Ability{Trigger: TriggerEntersPlay, Effect: Ready{Target: Target{Kind: TargetThisCreature}}})
	if want := SelfName + " enters play ready."; ready != want {
		t.Errorf("enters-play ready = %q, want %q", ready, want)
	}
	// An effect without a dedicated enter word falls back to its ordinary text.
	other := RenderAbility(Ability{Trigger: TriggerEntersPlay, Effect: GainAember{Player: Controller, Amount: 1}})
	if want := SelfName + " enters play gain 1 Æmber."; other != want {
		t.Errorf("enters-play fallback = %q, want %q", other, want)
	}
}

func TestTriggerPrefixDefault(t *testing.T) {
	// An unknown trigger renders with no prefix and a capitalized effect.
	got := RenderAbility(Ability{Trigger: Trigger(99), Effect: GainAember{Player: Controller, Amount: 1}})
	if got != "Gain 1 Æmber." {
		t.Errorf("RenderAbility(unknown) = %q", got)
	}
}

func TestEndOfTurnAbilityText(t *testing.T) {
	got := RenderAbility(Ability{Trigger: TriggerEndOfTurn, Effect: LoseAember{Player: Opponent, Amount: 1}})
	if want := "At the end of your turn, your opponent loses 1 Æmber."; got != want {
		t.Errorf("end-of-turn text = %q, want %q", got, want)
	}
}

func TestAfterYouPlayFolding(t *testing.T) {
	// A Conditional{ItIs} folds into the natural "after you play a <shape>" wording.
	folded := RenderAbility(Ability{Trigger: TriggerAfterCardPlayed, Effect: Conditional{Cond: ItIs{Type: Artifact}, Then: StealAember{Amount: 1}}})
	if want := "After you play an artifact, steal 1 Æmber."; folded != want {
		t.Errorf("folded = %q, want %q", folded, want)
	}
	// A non-Conditional reaction keeps the broad prefix.
	plain := RenderAbility(Ability{Trigger: TriggerAfterCardPlayed, Effect: GainAember{Player: Controller, Amount: 1}})
	if want := "After you play a card, gain 1 Æmber."; plain != want {
		t.Errorf("plain = %q, want %q", plain, want)
	}
	// A Conditional on something other than the played card's shape stays literal.
	stateGated := RenderAbility(Ability{Trigger: TriggerAfterCardPlayed, Effect: Conditional{Cond: OpponentAember{Is: AtLeast, Amount: 1}, Then: GainAember{Player: Controller, Amount: 1}}})
	if want := "After you play a card, if your opponent has 1 Æmber or more, gain 1 Æmber."; stateGated != want {
		t.Errorf("state-gated = %q, want %q", stateGated, want)
	}
}

func TestIsFightReapPair(t *testing.T) {
	ready := ReadyIfFirstUse{Target: Target{Kind: TargetThisCreature}}
	reap := Ability{Trigger: TriggerAfterReap, Effect: ready}
	fight := Ability{Trigger: TriggerAfterFight, Effect: ready}
	if !isFightReapPair(reap, fight) {
		t.Error("Reap+Fight sharing one effect should pair")
	}
	if !isFightReapPair(fight, reap) {
		t.Error("Fight+Reap (reversed order) should also pair")
	}
	if isFightReapPair(reap, Ability{Trigger: TriggerAfterReap, Effect: ready}) {
		t.Error("Reap+Reap is not a Fight/Reap pair")
	}
	if isFightReapPair(reap, Ability{Trigger: TriggerAfterFight, Effect: StealAember{Amount: 1}}) {
		t.Error("differing effects should not pair")
	}
}

func TestTargetTextDefault(t *testing.T) {
	cases := map[TargetKind]string{
		TargetThisCreature:       SelfName,
		TargetTriggeringCreature: "it",
		TargetTheOtherCreature:   "the other creature",
		TargetEachEnemyCreature:  "each enemy creature",
		TargetEachFriendlyInPlay: "each friendly card",
		TargetKind(99):           "a creature",
	}
	for kind, want := range cases {
		if got := (Target{Kind: kind}).Text(); got != want {
			t.Errorf("Target{%d}.Text() = %q, want %q", kind, got, want)
		}
	}
}

func TestAllTriggerPrefixes(t *testing.T) {
	cases := map[Trigger]string{
		TriggerAfterPlay:              "Play: Gain 1 Æmber.",
		TriggerAfterReap:              "Reap: Gain 1 Æmber.",
		TriggerAfterFight:             "Fight: Gain 1 Æmber.",
		TriggerBeforeFight:            "Before Fight: Gain 1 Æmber.",
		TriggerAction:                 "Action: Gain 1 Æmber.",
		TriggerDestroyed:              "Destroyed: Gain 1 Æmber.",
		TriggerAfterForgeKey:          "After you forge a key, gain 1 Æmber.",
		TriggerAfterCreatureEnters:    "After a creature enters play, gain 1 Æmber.",
		TriggerAfterDestroyedFighting: "After a creature is destroyed fighting {self}, gain 1 Æmber.",
		TriggerAfterCardPlayed:        "After you play a card, gain 1 Æmber.",
	}
	for tr, want := range cases {
		got := RenderAbility(Ability{Trigger: tr, Effect: GainAember{Player: Controller, Amount: 1}})
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
		{exGiant(), "House:  Brobnar\nType:   Creature\nRarity: Rare\nPower:  5\nTraits: Giant\n\nAfter you forge a key, deal 2 damage to each enemy creature."},
		{exBruteStrength(), "House:  Brobnar\nType:   Upgrade\nRarity: Uncommon\nÆmber:  1\n\nThis creature gains +5 power."},
		{exBattleFury(), "House:  Brobnar\nType:   Tactic\nRarity: Common\nÆmber:  1\n\nPlay: Ready and fight with a friendly creature."},
		{exAutocannon(), "House:  Brobnar\nType:   Artifact\nRarity: Rare\nÆmber:  1\nTraits: Weapon\n\nAfter a creature enters play, deal 1 damage to it."},
		{NewCard("Asp", Shadows, Creature, Uncommon, WithPower(3), WithKeywords(Skirmish, Poison)), "House:  Shadows\nType:   Creature\nRarity: Uncommon\nPower:  3\n\nSkirmish. Poison."},
		{NewCard("Anaphiel", Sanctum, Creature, Common, WithPower(6), WithArmor(1), WithTraits("Knight"), WithKeywords(Taunt)), "House:  Sanctum\nType:   Creature\nRarity: Common\nPower:  6\nArmor:  1\nTraits: Knight\n\nTaunt."},
		{NewCard("Tabris", Sanctum, Creature, Uncommon, WithPower(6), WithAbility(TriggerAfterFight, CaptureAember{Amount: 1, Target: Target{Kind: TargetThisCreature}, Source: Opponent})), "House:  Sanctum\nType:   Creature\nRarity: Uncommon\nPower:  6\n\nFight: Tabris captures 1 Æmber from your opponent."},
		{NewCard("Bear", Untamed, Creature, Common, WithPower(5), WithTraits("Beast"), WithAssault(2)), "House:  Untamed\nType:   Creature\nRarity: Common\nPower:  5\nTraits: Beast\n\nAssault 2."},
		{NewCard("Grub", Untamed, Creature, Rare, WithPower(2), WithHazardous(5)), "House:  Untamed\nType:   Creature\nRarity: Rare\nPower:  2\n\nHazardous 5."},
		{NewCard("Valdr", Brobnar, Creature, Common, WithPower(6), WithTraits("Giant"), WithAttackDamage(AttackDamage{Amount: 2, FlankOnly: true})), "House:  Brobnar\nType:   Creature\nRarity: Common\nPower:  6\nTraits: Giant\n\nValdr deals +2 Damage while attacking an enemy creature on the flank."},
		{NewCard("Spider", Mars, Creature, Common, WithPower(7), WithAttackDamage(AttackDamage{Fixed: true, Amount: 0})), "House:  Mars\nType:   Creature\nRarity: Common\nPower:  7\n\nSpider deals no damage when fighting."},
		{NewCard("Ether Spider", Mars, Creature, Uncommon, WithPower(7), WithAttackDamage(AttackDamage{Fixed: true, Amount: 0}), WithReplaces(Instead{Of: EventAemberAddedToPool, Player: Opponent, With: Capture})), "House:  Mars\nType:   Creature\nRarity: Uncommon\nPower:  7\n\nEther Spider deals no damage when fighting.\nIf Æmber would be added to your opponent's pool, instead Ether Spider captures it."},
		{NewCard("Bruiser", Brobnar, Creature, Common, WithPower(8), WithAttackDamage(AttackDamage{Fixed: true, Amount: 5})), "House:  Brobnar\nType:   Creature\nRarity: Common\nPower:  8\n\nBruiser deals 5 damage when fighting."},
		{NewCard("Basher", Brobnar, Creature, Common, WithPower(4), WithAttackDamage(AttackDamage{Amount: 2})), "House:  Brobnar\nType:   Creature\nRarity: Common\nPower:  4\n\nBasher deals +2 Damage when fighting."},
		{NewCard("Runner", Shadows, Upgrade, Uncommon, WithStatic(StaticModifier{Granted: []Ability{{Trigger: TriggerAfterReap, Effect: StealAember{Amount: 1}}}})), "House:  Shadows\nType:   Upgrade\nRarity: Uncommon\n\nThis creature gains, \"Reap: Steal 1 Æmber.\""},
		{NewCard("Boots", Logos, Upgrade, Uncommon, WithStatic(StaticModifier{Granted: []Ability{{Trigger: TriggerAfterReap, Effect: ReadyIfFirstUse{Target: Target{Kind: TargetThisCreature}}}, {Trigger: TriggerAfterFight, Effect: ReadyIfFirstUse{Target: Target{Kind: TargetThisCreature}}}}})), "House:  Logos\nType:   Upgrade\nRarity: Uncommon\n\nThis creature gains, \"Fight/Reap: If this is the first time this creature was used this turn, ready it.\""},
		{NewCard("Boots", Logos, Upgrade, Uncommon, WithStatic(StaticModifier{Keywords: []Keyword{Versatile}}), WithAbility(TriggerAfterPlay, Sequence{Effects: []Effect{Stun{Target: Target{Kind: TargetThisCreature}}, Exhaust{Target: Target{Kind: TargetThisCreature}}}})), "House:  Logos\nType:   Upgrade\nRarity: Uncommon\n\nThis creature gains versatile.\nPlay: Stun and exhaust this creature."},
		{NewCard("Jammer", Mars, Creature, Common, WithPower(4), WithArmor(1), WithTraits("Robot"), WithKeyCost(NewKeyCostChange(Opponent, 1)), WithAbility(TriggerAfterReap, CaptureAember{Amount: 1, Target: Target{Kind: TargetThisCreature}, Source: Opponent}), WithAbility(TriggerAfterFight, CaptureAember{Amount: 1, Target: Target{Kind: TargetThisCreature}, Source: Opponent})), "House:  Mars\nType:   Creature\nRarity: Common\nPower:  4\nArmor:  1\nTraits: Robot\n\nYour opponent's keys cost +1 Æmber.\nFight/Reap: Jammer captures 1 Æmber from your opponent."},
		{NewCard("Pack", Mars, Upgrade, Uncommon, WithAemberBonus(1), WithStatic(StaticModifier{KeyCostChange: NewKeyCostChange(Opponent, 2)})), "House:  Mars\nType:   Upgrade\nRarity: Uncommon\nÆmber:  1\n\nThis creature gains, \"Your opponent's keys cost +2 Æmber.\""},
		{NewCard("Shield", Sanctum, Upgrade, Rare, WithStatic(StaticModifier{Replaces: Replace{When: EventCreatureDestroyed, With: Sequence{Effects: []Effect{Heal{Fully: true, Target: Target{Kind: TargetTriggeringCreature}}, Destroy{Target: Target{Kind: TargetThisCreature}}}}}})), "House:  Sanctum\nType:   Upgrade\nRarity: Rare\n\nThis creature gains, \"If this creature would be destroyed, instead fully heal it, and destroy Shield.\""},
		{NewCard("Cloak", Sanctum, Upgrade, Rare, WithStatic(StaticModifier{HazardousBonus: 2, Replaces: Replace{When: EventCreatureDestroyed, With: Sequence{Effects: []Effect{Heal{Fully: true, Target: Target{Kind: TargetTriggeringCreature}}, Destroy{Target: Target{Kind: TargetThisCreature}}}}}})), "House:  Sanctum\nType:   Upgrade\nRarity: Rare\n\nThis creature gains +2 hazardous and, \"If this creature would be destroyed, instead fully heal it, and destroy Cloak.\""},
		{NewCard("Antenna", Mars, Upgrade, Rare, WithStatic(StaticModifier{Granted: []Ability{{Trigger: TriggerAfterCardPlayed, Effect: Conditional{Cond: ItIs{House: Mars, Type: Creature}, Then: Sequence{Effects: []Effect{Ready{Target: Target{Kind: TargetThisCreature}}, BelongToHouse{Target: Target{Kind: TargetThisCreature}, House: Mars, Duration: EndOfTurn}}}}}}})), "House:  Mars\nType:   Upgrade\nRarity: Rare\n\nThis creature gains, \"After you play a Mars creature, ready this creature, and for the remainder of the turn this creature belongs to house Mars.\""},
		{NewCard("SelfTax", Mars, Creature, Common, WithPower(3), WithKeyCost(NewKeyCostChange(Controller, 1))), "House:  Mars\nType:   Creature\nRarity: Common\nPower:  3\n\nYour keys cost +1 Æmber."},
		{NewCard("Tax", Mars, Creature, Common, WithPower(3), WithKeyCost(NewKeyCostChange(EachPlayer, 1))), "House:  Mars\nType:   Creature\nRarity: Common\nPower:  3\n\nEach player's keys cost +1 Æmber."},
		{NewCard("Imp", Dis, Creature, Common, WithPower(2), WithTraits("Imp"), WithRestrictions(Restrictions{PlayCardLimit: PlayCardLimit{Player: Opponent, Amount: 2}})), "House:  Dis\nType:   Creature\nRarity: Common\nPower:  2\nTraits: Imp\n\nYour opponent cannot play more than 2 cards each turn."},
		{NewCard("Witch", Untamed, Creature, Rare, WithPower(4), WithPlayPermission(PlayPermission{House: Untamed, Count: 1})), "House:  Untamed\nType:   Creature\nRarity: Rare\nPower:  4\n\nDuring each turn in which Untamed is not your active house, you may play one Untamed card."},
		{NewCard("Witch2", Untamed, Creature, Rare, WithPower(4), WithPlayPermission(PlayPermission{House: Untamed, Count: 2})), "House:  Untamed\nType:   Creature\nRarity: Rare\nPower:  4\n\nDuring each turn in which Untamed is not your active house, you may play 2 Untamed cards."},
		{NewCard("Twig", Untamed, Creature, Common, WithPower(7), WithTraits("Beast"), WithFightRestriction(Target{Kind: TargetEachCreature}.Stunned())), "House:  Untamed\nType:   Creature\nRarity: Common\nPower:  7\nTraits: Beast\n\nTwig can only fight stunned creatures."},
		{NewCard("Ritual", Dis, Artifact, Rare, WithTraits("Power"), WithConstantAbility(ConstantAbility{Target: Target{Kind: TargetEachCreature}, Granted: []Ability{{Trigger: TriggerDestroyed, Effect: PurgeCreature{Target: Target{Kind: TargetThisCreature}}}}})), "House:  Dis\nType:   Artifact\nRarity: Rare\nTraits: Power\n\nEach creature gains, \"Destroyed: Purge this creature.\""},
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
		{NewCard("Asp", Shadows, Creature, Uncommon, WithPower(3), WithKeywords(Skirmish, Poison)), "Skirmish. Poison."},
		// Upgrade static modifier (no own ability): the face still shows its text.
		{exBruteStrength(), "This creature gains +5 power."},
		// Triggered ability with a self-reference resolved to the card's name.
		{NewCard("Tabris", Sanctum, Creature, Uncommon, WithPower(6), WithAbility(TriggerAfterFight, CaptureAember{Amount: 1, Target: Target{Kind: TargetThisCreature}, Source: Opponent})), "Fight: Tabris captures 1 Æmber from your opponent."},
	}
	for _, tc := range cases {
		if got := RenderCardRules(&tc.def); got != tc.want {
			t.Errorf("%s rules mismatch:\n got:  %q\n want: %q", tc.def.Name, got, tc.want)
		}
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
			NewCard("Asp", Shadows, Creature, Uncommon, WithPower(3), WithKeywords(Skirmish, Poison)),
			"// Asp\n//\n//\tHouse:  Shadows\n//\tType:   Creature\n//\tRarity: Uncommon\n//\tPower:  3\n//\n//\tSkirmish. Poison.",
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

func TestIndefinite(t *testing.T) {
	cases := map[string]string{
		"":         "",
		"Urchin":   "an Urchin",
		"elf":      "an elf",
		"Knight":   "a Knight",
		"creature": "a creature",
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
	got := staticText(StaticModifier{PowerBonus: 5, ArmorBonus: 2})
	if got != "This creature gains +5 power and +2 armor." {
		t.Errorf("staticText = %q", got)
	}
	if got := staticText(StaticModifier{AssaultBonus: 2}); got != "This creature gains +2 assault." {
		t.Errorf("assault staticText = %q", got)
	}
	if got := staticText(StaticModifier{PowerBonus: 2, HazardousBonus: 2}); got != "This creature gains +2 power and +2 hazardous." {
		t.Errorf("hazardous staticText = %q", got)
	}
	if got := staticText(StaticModifier{Keywords: []Keyword{Skirmish}}); got != "This creature gains skirmish." {
		t.Errorf("keyword staticText = %q", got)
	}
	if got := staticText(StaticModifier{Keywords: []Keyword{Elusive, Skirmish}}); got != "This creature gains elusive and skirmish." {
		t.Errorf("two-keyword staticText = %q", got)
	}
	if got := staticText(StaticModifier{ArmorBonus: 1, Keywords: []Keyword{Taunt}}); got != "This creature gains +1 armor and taunt." {
		t.Errorf("armor+keyword staticText = %q", got)
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
