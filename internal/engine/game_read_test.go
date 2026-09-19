package engine

import (
	"errors"
	"slices"
	"strings"
	"testing"
)

// DeckList reads the whole catalog, so it lists a player's entire deck — draw
// pile included — deduplicated with an "xN" count and sorted.
func TestDeckList(t *testing.T) {
	g := started(t)
	// Player 0 owns two copies of Ape (an "xN" entry) and one Ogre (a single).
	g.AddToDeck(testCreature("Ape", 4), 0)
	g.AddToDeck(testCreature("Ape", 4), 0)
	g.AddToDeck(testCreature("Ogre", 6), 0)
	// Player 1 owns a single Imp, which must not appear in player 0's list.
	g.AddToDeck(testCreature("Imp", 1), 1)

	list := g.DeckList(0)
	want := []string{"Ape x2", "Ogre"}
	if !slices.Equal(list, want) {
		t.Errorf("DeckList(0) = %v, want %v", list, want)
	}
	if !slices.IsSorted(list) {
		t.Errorf("DeckList(0) is not sorted: %v", list)
	}
	if strings.Join(g.DeckList(1), ",") != "Imp" {
		t.Errorf("DeckList(1) = %v, want [Imp]", g.DeckList(1))
	}
}

// Quixxle Stone: while it is in play, any player who controls more creatures than
// their opponent cannot play creatures — symmetric, evaluated per attempting
// player, and only barring creatures.
func TestConditionalPlayBarText(t *testing.T) {
	def := NewCard(
		"Quixxle Stone",
		StarAlliance,
		Artifact,
		Rare,
		WithCannotPlayWhile(ConditionalPlayBar{
			Type: Creature,
			When: ControlsMoreCreatures{},
		}),
	)
	rules := cardRules(&def, false)
	want := "If a player has more creatures in play than their opponent, they cannot play creatures."
	found := false
	for _, r := range rules {
		if r == want {
			found = true
		}
	}
	if !found {
		t.Errorf("cardRules = %v, want a line %q", rules, want)
	}
}

// symmetricCondText falls back to the condition's own CondText (minus its "if "
// prefix) for any condition other than ControlsMoreCreatures.
func TestSymmetricCondTextFallback(t *testing.T) {
	if got := symmetricCondText(
		HasMoreForgedKeys{Player: Opponent},
	); got != "your opponent has more forged keys than you" {
		t.Errorf("symmetricCondText fallback = %q", got)
	}
}

func TestConditionalPlayBarBarsAheadPlayer(t *testing.T) {
	g := started(t) // player 0 active, Brobnar
	g.AddArtifact(
		NewCard(
			"Quixxle Stone",
			StarAlliance,
			Artifact,
			Rare,
			WithCannotPlayWhile(ConditionalPlayBar{
				Type: Creature,
				When: ControlsMoreCreatures{},
			}),
		),
		1, // controlled by the opponent, yet it bars whichever side is ahead
	)

	// Even counts (0 vs 0): the active player is not ahead, so creatures are legal.
	if g.cannotPlayCreatures(0) {
		t.Fatal("with equal creature counts, plays should not be barred")
	}

	// Player 0 pulls ahead: one creature to none.
	g.AddToBattleline(testCreature("ahead", 3), 0)
	if !g.cannotPlayCreatures(0) {
		t.Fatal("the player with more creatures should be barred from playing creatures")
	}
	g.AddToHand(testCreature("newbie", 2), 0)
	if _, err := g.PlayCreature(
		0,
		handIdx(g, 0, "newbie"),
		false,
	); !errors.Is(
		err,
		ErrCannotPlayCreature,
	) {
		t.Errorf("Playcreature = %v, want ErrCannotPlaycreature", err)
	}

	// The bar only blocks creatures: a non-creature is still playable.
	g.AddToHand(NewCard("act", Brobnar, Tactic, Common), 0)
	if err := g.PlayTactic(0, handIdx(g, 0, "act")); err != nil {
		t.Errorf("actions should still be playable: %v", err)
	}

	// The opponent draws even (1 vs 1): the active player is no longer ahead.
	g.AddToBattleline(testCreature("evener", 1), 1)
	if g.cannotPlayCreatures(0) {
		t.Fatal("with counts even again, the bar should lift")
	}
	g.AddToHand(testCreature("second", 2), 0)
	if _, err := g.PlayCreature(0, handIdx(g, 0, "second"), false); err != nil {
		t.Errorf("Playcreature once even = %v, want nil", err)
	}
}

// A directional constant ability reaches only the creatures on the named side of
// its source in the battleline: Panpaca, Anga buffs those to its right, Jaga
// those to its left.
func TestConstantAbilityDirectional(t *testing.T) {
	g := NewGame("A", "B", 1)
	left := g.AddToBattleline(testCreature("left", 4), 0)
	anga := g.AddToBattleline(NewCard("Anga", Untamed, Creature, Common,
		WithPower(5), WithConstantAbility(ConstantAbility{
			Target:     Target{Kind: TargetEachCreature}.ToRightOfSource(),
			PowerBonus: 2,
		})), 0)
	right := g.AddToBattleline(testCreature("right", 4), 0)

	if got := g.Power(right); got != 6 {
		t.Errorf("right creature power = %d, want 6", got)
	}
	if got := g.Power(left); got != 4 {
		t.Errorf("left creature power = %d, want 4 (not to Anga's right)", got)
	}
	if got := g.Power(anga); got != 5 {
		t.Errorf("Anga power = %d, want 5 (does not buff itself)", got)
	}

	def := g.cat.def(anga)
	if !strings.Contains(RenderCardRules(def),
		"to the right of Anga") {
		t.Errorf("Anga rules should name the right side: %q", RenderCardRules(def))
	}
}

// The left-facing variant reaches only creatures to its source's left, and a
// creature off the source's battleline is on neither side.
func TestConstantAbilityDirectionalLeft(t *testing.T) {
	g := NewGame("A", "B", 1)
	left := g.AddToBattleline(testCreature("left", 4), 0)
	jaga := g.AddToBattleline(NewCard("Jaga", Untamed, Creature, Common,
		WithPower(3), WithConstantAbility(ConstantAbility{
			Target:   Target{Kind: TargetEachCreature}.ToLeftOfSource(),
			Keywords: []Keyword{Skirmish},
		})), 0)
	right := g.AddToBattleline(testCreature("right", 4), 0)
	enemy := g.AddToBattleline(testCreature("enemy", 4), 1)

	if !g.hasKeyword(left, Skirmish) {
		t.Error("creature to Jaga's left should gain skirmish")
	}
	if g.hasKeyword(right, Skirmish) {
		t.Error("creature to Jaga's right should not gain skirmish")
	}
	if g.hasKeyword(enemy, Skirmish) {
		t.Error("enemy in another battleline is on neither side of Jaga")
	}
	if !strings.Contains(RenderCardRules(g.cat.def(jaga)), "to the left of Jaga") {
		t.Error("Jaga rules should name the left side")
	}
}

// InCenterOfBattleline is true only for the single middle creature of an
// odd-sized battleline; an even-sized line has no center, and a lone creature is
// its own center.
func TestInCenterOfBattleline(t *testing.T) {
	g := NewGame("A", "B", 1)

	// A lone creature is its own center.
	lone := g.AddToBattleline(testCreature("lone", 2), 0)
	if !g.InCenterOfBattleline(lone) {
		t.Error("a lone creature should be its own center")
	}

	// A second creature makes the line even: no center.
	second := g.AddToBattleline(testCreature("second", 2), 0)
	if g.InCenterOfBattleline(lone) || g.InCenterOfBattleline(second) {
		t.Error("an even-sized battleline should have no center")
	}

	// A third makes it odd again: only the middle creature is centered.
	third := g.AddToBattleline(testCreature("third", 2), 0)
	if !g.InCenterOfBattleline(second) {
		t.Error("the middle creature of an odd line should be centered")
	}
	if g.InCenterOfBattleline(lone) || g.InCenterOfBattleline(third) {
		t.Error("a flank creature of an odd line should not be centered")
	}
}

// A WhileInCenter constant self-buff applies only while the source is the middle
// creature of an odd-sized battleline, covering Kaloch Stonefather.
func TestConstantAbilityWhileInCenter(t *testing.T) {
	g := NewGame("A", "B", 1)
	def := NewCard("Kaloch", Brobnar, Creature, Rare, WithPower(6),
		WithConstantAbility(ConstantAbility{
			Target:        Target{Kind: TargetEachFriendlyCreature},
			Keywords:      []Keyword{Skirmish},
			WhileInCenter: true,
		}))

	if !strings.Contains(
		RenderCardRules(&def),
		"While Kaloch is in the center of your battleline, each friendly creature gains skirmish.",
	) {
		t.Error("card rules should render the while-in-center line")
	}

	left := g.AddToBattleline(testCreature("left", 3), 0)
	kaloch := g.AddToBattleline(def, 0)
	right := g.AddToBattleline(testCreature("right", 3), 0)

	// Kaloch is the middle creature, so the grant is active.
	if !g.hasKeyword(left, Skirmish) {
		t.Error("centered Kaloch should grant friendly creatures skirmish")
	}

	// Remove the right creature so the line is even: the grant lifts.
	g.State.Battleline[0].remove(right)
	if g.hasKeyword(left, Skirmish) {
		t.Error("off-center Kaloch should not grant skirmish")
	}
	_ = kaloch
}

// A WhileInCenter granted ability renders and reaches its target only while
// centered, covering The Shadow Council and Eldest Bear.
func TestConstantAbilityWhileInCenterGranted(t *testing.T) {
	g := NewGame("A", "B", 1)
	def := NewCard("Council", Shadows, Creature, Rare, WithPower(3),
		WithConstantAbility(ConstantAbility{
			Target:        Target{Kind: TargetThisCreature},
			WhileInCenter: true,
			Granted: []Ability{{
				Trigger: TriggerAction,
				Effect:  StealAember{Amount: 2},
			}},
		}))

	if !strings.Contains(RenderCardRules(&def),
		`While Council is in the center of your battleline, it gains, "Action: Steal 2 Æmber."`) {
		t.Errorf("card rules should render the while-in-center granted line, got %q",
			RenderCardRules(&def))
	}

	left := g.AddToBattleline(testCreature("left", 3), 0)
	council := g.AddToBattleline(def, 0)
	right := g.AddToBattleline(testCreature("right", 3), 0)
	_ = left

	if !g.HasTrigger(council, TriggerAction) {
		t.Error("centered Council should have the granted action")
	}

	g.State.Battleline[0].remove(right)
	if g.HasTrigger(council, TriggerAction) {
		t.Error("off-center Council should not have the granted action")
	}
}

// SourceInCenterOfBattleline is met only while the source is the middle creature
// of an odd-sized battleline.
func TestSourceInCenterOfBattlelineCondition(t *testing.T) {
	g := NewGame("A", "B", 1)
	left := g.AddToBattleline(testCreature("left", 2), 0)
	mid := g.AddToBattleline(testCreature("mid", 2), 0)
	g.AddToBattleline(testCreature("right", 2), 0)

	c := SourceInCenterOfBattleline{}
	if c.CondText() != "if "+SelfName+" is in the center of your battleline" {
		t.Errorf("CondText = %q", c.CondText())
	}
	if !c.Met(&EffectContext{
		Resolver: g,
		Source:   mid,
	}) {
		t.Error("middle creature should satisfy SourceInCenterOfBattleline")
	}
	if c.Met(&EffectContext{
		Resolver: g,
		Source:   left,
	}) {
		t.Error("flank creature should not satisfy SourceInCenterOfBattleline")
	}
}

// A WhileOffFlank constant self-buff applies only while the source sits off a
// flank (in the interior of its battleline), covering Gub.
func TestConstantAbilityWhileOffFlank(t *testing.T) {
	g := NewGame("A", "B", 1)
	def := NewCard("Gub", Dis, Creature, Common, WithPower(1),
		WithConstantAbility(ConstantAbility{
			Target:        Target{Kind: TargetThisCreature},
			PowerBonus:    5,
			Keywords:      []Keyword{Taunt},
			WhileOffFlank: true,
		}))

	if !strings.Contains(RenderCardRules(&def),
		"Gub gains +5 power and taunt while it is not on a flank.") {
		t.Error("card rules should render the while-off-flank line")
	}

	left := g.AddToBattleline(testCreature("left", 3), 0)
	gub := g.AddToBattleline(def, 0)
	right := g.AddToBattleline(testCreature("right", 3), 0)
	_ = left
	_ = right

	// Gub is the middle creature: off a flank, so the buff is active.
	if got := g.Power(gub); got != 6 {
		t.Errorf("off-flank power = %d, want 6 (1 + 5)", got)
	}
	if !g.hasKeyword(gub, Taunt) {
		t.Error("off-flank Gub should gain taunt")
	}

	// Remove the right creature so Gub becomes the right flank; the buff lifts.
	g.State.Battleline[0].remove(right)
	if got := g.Power(gub); got != 1 {
		t.Errorf("on-flank power = %d, want 1 (buff suspended)", got)
	}
	if g.hasKeyword(gub, Taunt) {
		t.Error("on-flank Gub should not have taunt")
	}
}

// A constant self-buff gated on the source being damaged applies only while it
// carries damage — Gron Nine-Toes's "+4 power while it is damaged."
func TestConstantAbilityWhileDamaged(t *testing.T) {
	g := NewGame("A", "B", 1)
	def := NewCard("Gron", Brobnar, Creature, Rare, WithPower(8),
		WithConstantAbility(ConstantAbility{
			Target:     Target{Kind: TargetThisCreature}.Damaged(),
			PowerBonus: 4,
		}))

	if !strings.Contains(RenderCardRules(&def),
		"Gron gains +4 power while it is damaged.") {
		t.Errorf("card rules missing the while-damaged line:\n%s", RenderCardRules(&def))
	}

	gron := g.AddToBattleline(def, 0)

	// Undamaged: the buff is suspended.
	if got := g.Power(gron); got != 8 {
		t.Errorf("undamaged power = %d, want 8 (buff suspended)", got)
	}

	// One point of damage turns the buff on.
	g.applyRawDamage(DamageTarget{
		ID:          gron,
		Amount:      1,
		IgnoreArmor: true,
	})
	if got := g.Power(gron); got != 12 {
		t.Errorf("damaged power = %d, want 12 (8 + 4)", got)
	}
}

// A constant ability scaled PerTarget reads the count separately for each creature
// it reaches — Tribune Pompitus's "+2 power for each Æmber on it."
func TestConstantAbilityPerTarget(t *testing.T) {
	g := NewGame("A", "B", 1)
	def := NewCard("Tribune", Saurian, Creature, Rare, WithPower(4),
		WithConstantAbility(ConstantAbility{
			Target:     Target{Kind: TargetEachFriendlyCreature},
			PowerBonus: 2,
			PerTarget:  AemberOnIt,
		}))

	if !strings.Contains(RenderCardRules(&def),
		"Each friendly creature gains +2 power for each Æmber on it.") {
		t.Errorf("card rules missing the per-target line:\n%s", RenderCardRules(&def))
	}

	tribune := g.AddToBattleline(def, 0)
	friend := g.AddToBattleline(testCreature("friend", 3), 0)

	// No Æmber on either: the bonus scales to zero.
	if got := g.Power(friend); got != 3 {
		t.Errorf("friend power with no Æmber = %d, want 3", got)
	}
	if got := g.Power(tribune); got != 4 {
		t.Errorf("tribune power with no Æmber = %d, want 4", got)
	}

	// Three Æmber on the friend: +2 per Æmber reaches only that creature.
	g.addAmberOn(friend, 3)
	if got := g.Power(friend); got != 9 {
		t.Errorf("friend power with 3 Æmber = %d, want 9 (3 + 2*3)", got)
	}
	if got := g.Power(tribune); got != 4 {
		t.Errorf("tribune power should still be 4 (its own Æmber is 0), got %d", got)
	}
}

// A constant ability scaled Per a source count names the source in its text and
// applies the same bonus to every creature it reaches — Primus Unguis's "+2 power
// for each Æmber on Primus Unguis."
func TestConstantAbilityPerSourceReachingOthers(t *testing.T) {
	g := NewGame("A", "B", 1)
	def := NewCard("Primus Unguis", Saurian, Creature, Rare, WithPower(5),
		WithConstantAbility(ConstantAbility{
			Target:     Target{Kind: TargetEachFriendlyCreature},
			PowerBonus: 2,
			Per:        AemberOnThis{},
		}))

	if !strings.Contains(RenderCardRules(&def),
		"Each friendly creature gains +2 power for each Æmber on Primus Unguis.") {
		t.Errorf("card rules missing the per-source line:\n%s", RenderCardRules(&def))
	}

	primus := g.AddToBattleline(def, 0)
	friend := g.AddToBattleline(testCreature("friend", 3), 0)

	// No Æmber on the source: the bonus scales to zero for everyone.
	if got := g.Power(friend); got != 3 {
		t.Errorf("friend power with no Æmber on source = %d, want 3", got)
	}

	// Two Æmber on the source reaches every friendly creature equally.
	g.addAmberOn(primus, 2)
	if got := g.Power(friend); got != 7 {
		t.Errorf("friend power = %d, want 7 (3 + 2*2)", got)
	}
	if got := g.Power(primus); got != 9 {
		t.Errorf("primus power = %d, want 9 (5 + 2*2)", got)
	}
}

// A constant ability targeting friendly creatures buffs every friendly creature
// and no enemy, and the buff vanishes the moment the source card leaves play.
func TestConstantAbilityFriendlyTarget(t *testing.T) {
	g := NewGame("A", "B", 1)
	banner := g.Register(NewCard(
		"Banner",
		Brobnar,
		Artifact,
		Rare,
		WithConstantAbility(
			ConstantAbility{
				PowerBonus: 1,
				Target:     Target{Kind: TargetEachFriendlyCreature},
			},
		),
	), 0)
	g.State.Artifacts[0].add(banner)
	friend := g.AddToBattleline(testCreature("friend", 3), 0)
	enemy := g.AddToBattleline(testCreature("enemy", 3), 1)

	if got := g.Power(friend); got != 4 {
		t.Errorf("friendly power = %d, want 4 (3 + constant)", got)
	}
	if got := g.Power(enemy); got != 3 {
		t.Errorf("enemy power = %d, want 3 (a friendly target does not reach enemies)", got)
	}

	g.State.Artifacts[0].remove(banner)
	if got := g.Power(friend); got != 3 {
		t.Errorf("friendly power after the source leaves = %d, want 3", got)
	}
}

// A constant ability targeting neighbors reaches only the source's immediate
// battleline neighbors.
func TestConstantAbilityNeighboringTarget(t *testing.T) {
	g := NewGame("A", "B", 1)
	bulwark := NewCard(
		"Bulwark",
		Sanctum,
		Creature,
		Common,
		WithPower(4),
		WithArmor(2),
		WithConstantAbility(
			ConstantAbility{
				ArmorBonus: 2,
				Target:     Target{Kind: TargetEachCreature}.Neighboring(),
			},
		),
	)
	left := g.AddToBattleline(testCreature("left", 3), 0)
	g.AddToBattleline(bulwark, 0)
	right := g.AddToBattleline(testCreature("right", 3), 0)
	far := g.AddToBattleline(testCreature("far", 3), 0)

	if got := g.Armor(left); got != 2 {
		t.Errorf("left neighbor armor = %d, want 2", got)
	}
	if got := g.Armor(right); got != 2 {
		t.Errorf("right neighbor armor = %d, want 2", got)
	}
	if got := g.Armor(far); got != 0 {
		t.Errorf("non-neighbor armor = %d, want 0", got)
	}
}

// A constant ability can grant a "cannot reap" restriction to the creatures its
// Target reaches: Narp bars its neighbors from reaping while distant friends stay
// free to reap, and the restriction reads on the card.
func TestCannotBeUsedToFromConstantAbility(t *testing.T) {
	g := started(t)
	narp := NewCard("Narp", Brobnar, Creature, Common, WithPower(8),
		WithConstantAbility(ConstantAbility{
			Target:         Target{Kind: TargetEachCreature}.Neighboring(),
			CannotBeUsedTo: []UseKind{ReapUse},
		}))
	left := g.AddToBattleline(testCreature("left", 3), 0)
	g.AddToBattleline(narp, 0)
	right := g.AddToBattleline(testCreature("right", 3), 0)
	far := g.AddToBattleline(testCreature("far", 3), 0)

	if err := g.CanUseTo(0, left, ReapUse); !errors.Is(err, ErrCannotUse) {
		t.Errorf("left neighbor reap = %v, want ErrCannotUse", err)
	}
	if err := g.CanUseTo(0, right, ReapUse); !errors.Is(err, ErrCannotUse) {
		t.Errorf("right neighbor reap = %v, want ErrCannotUse", err)
	}
	if err := g.CanUseTo(0, far, ReapUse); err != nil {
		t.Errorf("far friend reap = %v, want nil", err)
	}
	if got := constantText(&narp); got != "Each neighboring creature cannot reap." {
		t.Errorf("Narp constant text = %q", got)
	}
}

// A constant ability with no target reaches every creature in play, including
// the source itself and the enemy.
func TestConstantAbilityNoTargetReachesEveryone(t *testing.T) {
	g := NewGame("A", "B", 1)
	self := g.AddToBattleline(NewCard("Totem", Brobnar, Creature, Common, WithPower(2),
		WithConstantAbility(ConstantAbility{PowerBonus: 1})), 0)
	friend := g.AddToBattleline(testCreature("friend", 3), 0)
	enemy := g.AddToBattleline(testCreature("enemy", 3), 1)

	if got := g.Power(self); got != 3 {
		t.Errorf("source power = %d, want 3 (it buffs itself)", got)
	}
	if got := g.Power(friend); got != 4 {
		t.Errorf("friend power = %d, want 4", got)
	}
	if got := g.Power(enemy); got != 4 {
		t.Errorf("enemy power = %d, want 4 (an untargeted constant reaches everyone)", got)
	}
}

// An untargeted constant ability's scope also reaches artifacts in play, not just
// creatures (creatures have no artifact-relevant stat yet, so this is about the
// scope, verified through constantAffects).
func TestConstantAbilityReachesArtifacts(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(NewCard("Totem", Brobnar, Creature, Common, WithPower(2),
		WithConstantAbility(ConstantAbility{PowerBonus: 1})), 0)
	art := g.AddArtifact(exAutocannon(), 1)

	if !g.constantAffects(src, g.cat.def(src).ConstantAbilities[0], art) {
		t.Error("an untargeted constant ability should reach artifacts in play")
	}
}

// A constant ability that grants a trigger gives it only to the creatures its
// Target reaches: hasTrigger sees it on a reached creature and skips a creature
// the grantor does not affect.
func TestHasTriggerFromConstantAbility(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.AddToBattleline(NewCard("Grantor", Brobnar, Creature, Common, WithPower(3),
		WithConstantAbility(ConstantAbility{
			Target: Target{Kind: TargetEachFriendlyCreature},
			Granted: []Ability{
				{Trigger: TriggerAfterReap, Effect: GainAember{
					Player: Controller,
					Amount: 1,
				}},
			},
		})), 0)
	friend := g.AddToBattleline(testCreature("friend", 3), 0)
	enemy := g.AddToBattleline(testCreature("enemy", 3), 1)

	if !g.hasTrigger(friend, TriggerAfterReap) {
		t.Error("a friendly creature should gain the granted Reap trigger")
	}
	if g.hasTrigger(enemy, TriggerAfterReap) {
		t.Error("an enemy creature is not reached by the friendly constant ability")
	}
}

func TestHasKeywordFromConstantAbility(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.AddToBattleline(NewCard("Grantor", Untamed, Creature, Common, WithPower(3),
		WithConstantAbility(ConstantAbility{
			Target:   Target{Kind: TargetEachFriendlyCreature},
			Keywords: []Keyword{Skirmish},
		})), 0)
	friend := g.AddToBattleline(testCreature("friend", 3), 0)
	enemy := g.AddToBattleline(testCreature("enemy", 3), 1)

	if !g.hasKeyword(friend, Skirmish) {
		t.Error("a friendly creature should gain the granted keyword")
	}
	if g.hasKeyword(friend, Poison) {
		t.Error("only the granted keyword is gained")
	}
	if g.hasKeyword(enemy, Skirmish) {
		t.Error("an enemy creature is not reached by the friendly constant ability")
	}
}

func TestConstantText(t *testing.T) {
	banner := NewCard(
		"Banner",
		Brobnar,
		Artifact,
		Rare,
		WithConstantAbility(
			ConstantAbility{
				PowerBonus: 1,
				Target:     Target{Kind: TargetEachFriendlyCreature},
			},
		),
	)
	if got := constantText(&banner); got != "Each friendly creature gains +1 power." {
		t.Errorf("friendly constant text = %q", got)
	}
	if !strings.Contains(RenderCardText(&banner), "Each friendly creature gains +1 power.") {
		t.Error("RenderCardText should include the constant-ability line")
	}

	bulwark := NewCard(
		"Bulwark",
		Sanctum,
		Creature,
		Common,
		WithPower(4),
		WithConstantAbility(
			ConstantAbility{
				ArmorBonus: 2,
				Target:     Target{Kind: TargetEachCreature}.Neighboring(),
			},
		),
	)
	if got := constantText(&bulwark); got != "Each neighboring creature gains +2 armor." {
		t.Errorf("neighboring constant text = %q", got)
	}

	all := NewCard("Totem", Brobnar, Creature, Common, WithPower(2),
		WithConstantAbility(ConstantAbility{PowerBonus: 1}))
	if got := constantText(&all); got != "Each card in play gains +1 power." {
		t.Errorf("untargeted constant text = %q", got)
	}

	table := NewCard(
		"Table",
		Sanctum,
		Artifact,
		Rare,
		WithConstantAbility(
			ConstantAbility{
				PowerBonus: 1,
				Keywords:   []Keyword{Taunt},
				Target:     Target{Kind: TargetEachFriendlyCreature},
			},
		),
	)
	if got := constantText(&table); got != "Each friendly creature gains +1 power and taunt." {
		t.Errorf("keyword constant text = %q", got)
	}
	kwOnly := NewCard(
		"Halo",
		Untamed,
		Creature,
		Common,
		WithPower(4),
		WithConstantAbility(
			ConstantAbility{
				Keywords: []Keyword{Skirmish},
				Target:   Target{Kind: TargetEachFriendlyCreature},
			},
		),
	)
	if got := constantText(&kwOnly); got != "Each friendly creature gains skirmish." {
		t.Errorf("keyword-only constant text = %q", got)
	}

	flank := NewCard(
		"Staunch Knight",
		Sanctum,
		Creature,
		Common,
		WithPower(4),
		WithConstantAbility(
			ConstantAbility{
				PowerBonus: 2,
				Target:     Target{Kind: TargetThisCreature}.OnFlank(),
			},
		),
	)
	if got := constantText(&flank); got != "Staunch Knight gains +2 power while it is on a flank." {
		t.Errorf("flank self constant text = %q", got)
	}

	plain := NewCard("Plain", Brobnar, Creature, Common, WithPower(3))
	if got := constantText(&plain); got != "" {
		t.Errorf("no-constant text = %q, want empty", got)
	}

	blank := NewCard(
		"Blossom Drake",
		Untamed,
		Creature,
		Rare,
		WithPower(4),
		WithConstantAbility(
			ConstantAbility{
				Target:    Target{Kind: TargetEachArtifact},
				BlankText: true,
			},
		),
	)
	if got := constantText(&blank); got !=
		"Each artifact's text box is considered blank, except for traits." {
		t.Errorf("blank-text constant = %q", got)
	}
}

// TestConstantHazardousGrant covers a constant ability that grants Hazardous to
// the creatures it reaches, both in its printed text and the value it produces.
func TestConstantHazardousGrant(t *testing.T) {
	molina := NewCard(
		"Arms",
		StarAlliance,
		Creature,
		Common,
		WithPower(4),
		WithHazardous(3),
		WithConstantAbility(ConstantAbility{
			HazardousBonus: 3,
			Target:         Target{Kind: TargetEachCreature}.Neighboring(),
		}),
	)
	if got := constantText(&molina); got != "Each neighboring creature gains hazardous 3." {
		t.Errorf("hazardous constant text = %q", got)
	}

	g := NewGame("A", "B", 1)
	left := g.AddToBattleline(testCreature("l", 3), 0)
	g.AddToBattleline(molina, 0)
	g.AddToBattleline(testCreature("r", 3), 0)
	far := g.AddToBattleline(testCreature("f", 3), 1)
	if got := g.Hazardous(left); got != 3 {
		t.Errorf("neighbor hazardous = %d, want 3", got)
	}
	if got := g.Hazardous(far); got != 0 {
		t.Errorf("distant creature hazardous = %d, want 0", got)
	}
}

// TestConstantAssaultGrant covers a constant ability that grants Assault to the
// creatures it reaches, both in its printed text and the value it produces.
func TestConstantAssaultGrant(t *testing.T) {
	bullwark := NewCard(
		"Bull",
		Sanctum,
		Creature,
		Common,
		WithPower(4),
		WithAssault(2),
		WithConstantAbility(ConstantAbility{
			AssaultBonus: 2,
			Target:       Target{Kind: TargetEachCreature}.Neighboring(),
		}),
	)
	if got := constantText(&bullwark); got != "Each neighboring creature gains assault 2." {
		t.Errorf("assault constant text = %q", got)
	}

	g := NewGame("A", "B", 1)
	left := g.AddToBattleline(testCreature("l", 3), 0)
	g.AddToBattleline(bullwark, 0)
	g.AddToBattleline(testCreature("r", 3), 0)
	far := g.AddToBattleline(testCreature("f", 3), 1)
	if got := g.assault(left); got != 2 {
		t.Errorf("neighbor assault = %d, want 2", got)
	}
	if got := g.assault(far); got != 0 {
		t.Errorf("distant creature assault = %d, want 0", got)
	}
}

// TestUpgradeStatBonusHonoursWhileOnFlankForEveryStat pins that a WhileOnFlank
// upgrade suspends every stat it grants, not just power and armor. Assault,
// Hazardous, and Splash used to read Static directly and applied their bonuses
// wherever the host stood; no printed card combines WhileOnFlank with those stats
// yet, so nothing caught it.
func TestUpgradeStatBonusHonoursWhileOnFlankForEveryStat(t *testing.T) {
	flankOnly := NewCard(
		"Flank Rig",
		Sanctum,
		Upgrade,
		Common,
		WithStatic(StaticModifier{
			PowerBonus:        1,
			ArmorBonus:        1,
			AssaultBonus:      2,
			HazardousBonus:    3,
			SplashAttackBonus: 4,
			WhileOnFlank:      true,
		}),
	)

	g := NewGame("A", "B", 1)
	host := g.AddToBattleline(testCreature("host", 3), 0)
	g.AttachUpgrade(host, g.Register(flankOnly, 0))

	onFlank := map[string]int{
		"power": g.Power(host), "armor": g.Armor(host), "assault": g.assault(host),
		"hazardous": g.hazardous(host), "splash": g.splashAttack(host),
	}
	want := map[string]int{
		"power": 4, "armor": 1, "assault": 2, "hazardous": 3, "splash": 4,
	}
	for stat, got := range onFlank {
		if got != want[stat] {
			t.Errorf("on flank %s = %d, want %d", stat, got, want[stat])
		}
	}

	// Flanked on both sides, the upgrade is suspended and grants nothing.
	off := NewGame("A", "B", 1)
	off.AddToBattleline(testCreature("left", 3), 0)
	host = off.AddToBattleline(testCreature("host", 3), 0)
	off.AddToBattleline(testCreature("right", 3), 0)
	off.AttachUpgrade(host, off.Register(flankOnly, 0))
	offFlank := map[string]int{
		"power": off.Power(host), "armor": off.Armor(host), "assault": off.assault(host),
		"hazardous": off.hazardous(host), "splash": off.splashAttack(host),
	}
	bare := map[string]int{
		"power": 3, "armor": 0, "assault": 0, "hazardous": 0, "splash": 0,
	}
	for stat, got := range offFlank {
		if got != bare[stat] {
			t.Errorf("off flank %s = %d, want %d", stat, got, bare[stat])
		}
	}
}

// TestConstantPerCount covers a constant ability whose bonus scales with a
// running count, both in its printed text and in the power it produces.
func TestConstantPerCount(t *testing.T) {
	def := NewCard(
		"Mushroom Man",
		Untamed,
		Creature,
		Uncommon,
		WithPower(2),
		WithConstantAbility(ConstantAbility{
			PowerBonus: 3,
			Target:     Target{Kind: TargetThisCreature},
			Per:        UnforgedKeys{Player: Controller},
		}),
	)
	want := "Mushroom Man gains +3 power for each unforged key you have."
	if got := constantText(&def); got != want {
		t.Errorf("constant text = %q, want %q", got, want)
	}
	if got := (UnforgedKeys{Player: Opponent}).CountText(); got != "unforged key your opponent has" {
		t.Errorf("opponent count text = %q", got)
	}

	g := NewGame("A", "B", 1)
	id := g.AddToBattleline(def, 0)
	if got := g.Power(id); got != 11 {
		t.Errorf("power with no keys forged = %d, want 11", got)
	}
	g.ManualForgeKey(0)
	if got := g.Power(id); got != 8 {
		t.Errorf("power with one key forged = %d, want 8", got)
	}

	// AemberOnThis reads the Æmber sitting on the source card.
	onThis := AemberOnThis{}
	if got := onThis.CountText(); got != "Æmber on it" {
		t.Errorf("AemberOnThis count text = %q", got)
	}
	g.AddAmberOn(id, 2)
	if got := onThis.Value(&EffectContext{
		Resolver:   g,
		Source:     id,
		Controller: 0,
	}); got != 2 {
		t.Errorf("AemberOnThis value = %d, want 2", got)
	}

	// DamageOnThis reads the damage sitting on the source card.
	dmgOnThis := DamageOnThis{}
	if got := dmgOnThis.CountText(); got != "damage on it" {
		t.Errorf("DamageOnThis count text = %q", got)
	}
	g.SetDamage(id, 3)
	if got := dmgOnThis.Value(&EffectContext{
		Resolver:   g,
		Source:     id,
		Controller: 0,
	}); got != 3 {
		t.Errorf("DamageOnThis value = %d, want 3", got)
	}
}

// KeyColorForged is met only while the named player has forged a key of the given
// colour, covering The Red Baron.
func TestKeyColorForgedCondition(t *testing.T) {
	g := NewGame("A", "B", 1)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}

	own := KeyColorForged{
		Player: Controller,
		Color:  KeyColorRed,
	}
	opp := KeyColorForged{
		Player: Opponent,
		Color:  KeyColorRed,
	}

	if own.CondText() != "if your red key is forged" {
		t.Errorf("own CondText = %q", own.CondText())
	}
	if opp.CondText() != "if your opponent's red key is forged" {
		t.Errorf("opp CondText = %q", opp.CondText())
	}

	if own.Met(ctx) || opp.Met(ctx) {
		t.Error("no key forged: neither condition should be met")
	}

	g.State.ForgeCanonicalKeys(0, 1)
	g.State.KeyColors[0][0] = KeyColorRed
	if !own.Met(ctx) {
		t.Error("your red key forged: own condition should be met")
	}
	if opp.Met(ctx) {
		t.Error("your red key forged: opponent condition should not be met")
	}

	g.State.ForgeCanonicalKeys(1, 1)
	g.State.KeyColors[1][0] = KeyColorBlue
	if opp.Met(ctx) {
		t.Error("opponent's blue key forged: red condition should not be met")
	}
	g.State.KeyColors[1][0] = KeyColorRed
	if !opp.Met(ctx) {
		t.Error("opponent's red key forged: opponent condition should be met")
	}
}

// A WhileCondition constant self-grant applies its granted ability only while the
// condition holds, and renders the "While ..." line, covering The Red Baron's reap.
func TestConstantAbilityWhileConditionGranted(t *testing.T) {
	g := NewGame("A", "B", 1)
	def := NewCard("Baron", Brobnar, Creature, Special, WithPower(4),
		WithConstantAbility(ConstantAbility{
			Target: Target{Kind: TargetThisCreature},
			WhileCondition: KeyColorForged{
				Player: Controller,
				Color:  KeyColorRed,
			},
			Granted: []Ability{{
				Trigger: TriggerAfterReap,
				Effect:  StealAember{Amount: 1},
			}},
		}))

	if !strings.Contains(RenderCardRules(&def),
		`While your red key is forged, Baron gains, "Reap: Steal 1 Æmber."`) {
		t.Errorf("card rules should render the while-condition granted line, got %q",
			RenderCardRules(&def))
	}

	baron := g.AddToBattleline(def, 0)

	if g.HasTrigger(baron, TriggerAfterReap) {
		t.Error("no key forged: Baron should not have the granted reap")
	}

	g.State.ForgeCanonicalKeys(0, 1)
	g.State.KeyColors[0][0] = KeyColorRed
	if !g.HasTrigger(baron, TriggerAfterReap) {
		t.Error("your red key forged: Baron should have the granted reap")
	}
}

// A WhileCondition constant keyword grant applies only while the condition holds,
// and renders the "While ..." line, covering The Red Baron's elusive.
func TestConstantAbilityWhileConditionKeyword(t *testing.T) {
	g := NewGame("A", "B", 1)
	def := NewCard("Baron", Brobnar, Creature, Special, WithPower(4),
		WithConstantAbility(ConstantAbility{
			Target: Target{Kind: TargetThisCreature},
			WhileCondition: KeyColorForged{
				Player: Opponent,
				Color:  KeyColorRed,
			},
			Keywords: []Keyword{Elusive},
		}))

	if !strings.Contains(RenderCardRules(&def),
		"While your opponent's red key is forged, Baron gains elusive.") {
		t.Errorf("card rules should render the while-condition keyword line, got %q",
			RenderCardRules(&def))
	}

	baron := g.AddToBattleline(def, 0)

	if g.hasKeyword(baron, Elusive) {
		t.Error("no key forged: Baron should not have elusive")
	}

	g.State.ForgeCanonicalKeys(1, 1)
	g.State.KeyColors[1][0] = KeyColorRed
	if !g.hasKeyword(baron, Elusive) {
		t.Error("opponent's red key forged: Baron should have elusive")
	}
}

// TestConstantBlankTextBlanksArtifacts checks a while-in-play BlankText constant
// ability blanks the artifacts its Target reaches, suppressing their abilities.
func TestConstantBlankTextBlanksArtifacts(t *testing.T) {
	g := started(t)
	// An enemy artifact grants every creature +2 power via a constant ability.
	buffer := g.AddArtifact(NewCard("buffer", Untamed, Artifact, Rare,
		WithConstantAbility(ConstantAbility{
			PowerBonus: 2,
			Target:     Target{Kind: TargetEachCreature},
		})), 1)
	victim := g.AddToBattleline(testCreature("victim", 3), 1)
	if g.Power(victim) != 5 {
		t.Fatalf("precondition: buffer should give +2 power, got %d", g.Power(victim))
	}

	drake := g.AddToBattleline(NewCard("drake", Untamed, Creature, Rare, WithPower(4),
		WithConstantAbility(ConstantAbility{
			Target:    Target{Kind: TargetEachArtifact},
			BlankText: true,
		})), 0)

	if !g.textBlanked(buffer) {
		t.Error("the artifact's text box should be blanked")
	}
	if g.Power(victim) != 3 {
		t.Errorf("the blanked artifact's buff should vanish, power = %d", g.Power(victim))
	}

	// Blanking the drake itself (registry blank) lifts its constant blank, so the
	// artifact's text returns.
	g.BlankEnemyText(1) // player 1 blanks player 0's creatures, including the drake
	if !g.textBlanked(drake) {
		t.Fatal("the drake should be blanked by the registry")
	}
	if g.textBlanked(buffer) {
		t.Error("a blanked drake grants no blank, so the artifact's text returns")
	}
}

// TestConstantRemovesTraitsStripsTraits checks a while-in-play RemovesTraits
// constant ability strips the traits of every creature its Target reaches, and
// that the traits return once the source leaves play.
func TestConstantRemovesTraitsStripsTraits(t *testing.T) {
	g := started(t)
	victim := g.AddToBattleline(testCreature("victim", 3, WithTraits(Mutant, Beast)), 1)
	kin := g.AddToBattleline(testCreature("kin", 3, WithTraits(Beast)), 1)
	if !g.HasTrait(victim, Mutant) || g.TraitCount(victim) != 2 {
		t.Fatal("precondition: victim should have its two printed traits")
	}
	if !g.SharesTrait(victim, kin) {
		t.Fatal("precondition: victim and kin share the Beast trait")
	}

	aberrant := g.AddToBattleline(NewCard("aberrant", Sanctum, Creature, Common,
		WithPower(3), WithConstantAbility(ConstantAbility{
			Target:        Target{Kind: TargetEachCreature},
			RemovesTraits: true,
		})), 0)

	if g.HasTrait(victim, Mutant) {
		t.Error("victim should lose its traits while the aberrant is in play")
	}
	if g.TraitCount(victim) != 0 {
		t.Errorf("victim trait count = %d, want 0", g.TraitCount(victim))
	}
	if g.SharesTrait(victim, kin) {
		t.Error("with traits stripped, victim shares no trait with kin")
	}

	// Blanking the aberrant lifts its constant, so the traits return.
	g.BlankEnemyText(1) // player 1 blanks player 0's creatures, including the aberrant.
	if !g.textBlanked(aberrant) {
		t.Fatal("the aberrant should be blanked by the registry")
	}
	if !g.HasTrait(victim, Mutant) || g.TraitCount(victim) != 2 {
		t.Error("a blanked aberrant strips nothing, so the traits return")
	}
}

// archivistDef builds a card carrying the SelectiveArchivePickup constant ability
// (The Archivist), for the archive-pickup tests below.
func archivistDef() CardDefinition {
	return NewCard("The Archivist", Logos, Creature, Rare,
		WithPower(3), WithConstantAbility(ConstantAbility{SelectiveArchivePickup: true}))
}

// TestConstantSelectiveArchivePickupScanner checks the scanner reports the
// granting card for its own controller only, and drops out once that card is
// blanked.
func TestConstantSelectiveArchivePickupScanner(t *testing.T) {
	g := started(t)
	if _, ok := g.constantSelectiveArchivePickup(0); ok {
		t.Fatal("no Archivist in play: scanner should report false")
	}

	src := g.AddToBattleline(archivistDef(), 0)
	got, ok := g.constantSelectiveArchivePickup(0)
	if !ok || got != src {
		t.Errorf("scanner = (%v, %v), want (%v, true)", got, ok, src)
	}

	// It is player-scoped: player 1 does not gain player 0's Archivist rule.
	if _, ok := g.constantSelectiveArchivePickup(1); ok {
		t.Error("the Archivist changes only its own controller's pickup rule")
	}

	// A blanked Archivist grants nothing (constantActive gate).
	g.BlankEnemyText(1) // player 1 blanks player 0's creatures, including the Archivist.
	if !g.textBlanked(src) {
		t.Fatal("the Archivist should be blanked by the registry")
	}
	if _, ok := g.constantSelectiveArchivePickup(0); ok {
		t.Error("a blanked Archivist should not grant selective pickup")
	}
}

// TestOfferArchivesSelectivePickup checks that with the Archivist in play the
// house-choice pickup lets the player take some archived cards and leave others.
func TestOfferArchivesSelectivePickup(t *testing.T) {
	g := started(t)
	g.AddToBattleline(archivistDef(), 0)
	one := g.AddToArchives(NewCard("Archived One", Logos, Creature, Common), 0)
	two := g.AddToArchives(NewCard("Archived Two", Logos, Creature, Common), 0)
	three := g.AddToArchives(NewCard("Archived Three", Logos, Creature, Common), 0)
	// Take one and two, decline the third.
	g.SetChooser(0, &declineAfterChooser{ids: []LocalID{one, two}})

	g.offerArchives(0)

	if !g.State.Hand[0].contains(one) || !g.State.Hand[0].contains(two) {
		t.Error("the chosen cards should move into hand")
	}
	if g.State.Hand[0].contains(three) {
		t.Error("the declined card should not enter hand")
	}
	if !g.State.Archives[0].contains(three) {
		t.Error("the declined card should remain in archives")
	}
}

// TestOfferArchivesSelectivePickupTakingNone checks that declining every card
// leaves the archives untouched and takes nothing into hand.
func TestOfferArchivesSelectivePickupTakingNone(t *testing.T) {
	g := started(t)
	g.AddToBattleline(archivistDef(), 0)
	one := g.AddToArchives(NewCard("Archived", Logos, Creature, Common), 0)
	g.SetChooser(0, &declineAfterChooser{}) // decline immediately

	g.offerArchives(0)

	if !g.State.Archives[0].contains(one) {
		t.Error("declining should leave the card in archives")
	}
	if g.State.Hand[0].contains(one) {
		t.Error("declining should take nothing into hand")
	}
}

// TestConstantSelectiveArchivePickupText checks the constant ability renders its
// own printed line.
func TestConstantSelectiveArchivePickupText(t *testing.T) {
	def := archivistDef()
	got := constantText(&def)
	want := "Instead of picking up all of your archives, you may pick up any number of cards in your archives."
	if got != want {
		t.Errorf("constantText = %q, want %q", got, want)
	}
}

// TestConstantRemovesTraitsText checks the sibling RemovesTraits branch renders
// its printed line (Grey Aberrant).
func TestConstantRemovesTraitsText(t *testing.T) {
	def := NewCard("Grey Aberrant", Sanctum, Creature, Common, WithPower(3),
		WithConstantAbility(ConstantAbility{
			Target:        Target{Kind: TargetEachCreature},
			RemovesTraits: true,
		}))
	got := constantText(&def)
	want := "Each creature loses each of its traits."
	if got != want {
		t.Errorf("constantText = %q, want %q", got, want)
	}
}

// A card with GrantsEntersReady makes friendly cards of that type enter play
// ready instead of exhausted while it is in play — Duskwitch for creatures, The
// Curator for artifacts.
func TestEntersPlayReady(t *testing.T) {
	t.Run("creature grant readies played creatures", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.StartTurn(0)
		granter := g.AddToHand(
			NewCard("granter", Untamed, Creature, Common,
				WithPower(1), WithFriendlyEntersPlayReady(EntersReadyGrant{Type: Creature})),
			0,
		)
		g.putIntoPlay(granter, 0)

		newbie := g.AddToHand(NewCard("newbie", Untamed, Creature, Common, WithPower(3)), 0)
		if _, err := g.PlayCreature(0, handIdxByID(g, 0, newbie), false); err != nil {
			t.Fatalf("Playcreature: %v", err)
		}
		if g.State.Cards[newbie].Exhausted {
			t.Error("creature should enter play ready under a creature grant")
		}
	})

	t.Run("artifact grant readies played artifacts", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.StartTurn(0)
		granter := g.AddToHand(
			NewCard("granter", Logos, Creature, Common,
				WithPower(1), WithFriendlyEntersPlayReady(EntersReadyGrant{Type: Artifact})),
			0,
		)
		g.putIntoPlay(granter, 0)

		art := g.AddToHand(NewCard("art", Logos, Artifact, Common), 0)
		if _, err := g.PlayArtifact(0, handIdxByID(g, 0, art)); err != nil {
			t.Fatalf("Playartifact: %v", err)
		}
		if g.State.Cards[art].Exhausted {
			t.Error("artifact should enter play ready under an artifact grant")
		}
	})

	t.Run("an opponent's grant does not ready your card", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.StartTurn(0)
		granter := g.AddToHand(
			NewCard("granter", Logos, Creature, Common,
				WithPower(1), WithFriendlyEntersPlayReady(EntersReadyGrant{Type: Artifact})),
			1,
		)
		g.putIntoPlay(granter, 1)

		art := g.AddToHand(NewCard("art", Logos, Artifact, Common), 0)
		if _, err := g.PlayArtifact(0, handIdxByID(g, 0, art)); err != nil {
			t.Fatalf("Playartifact: %v", err)
		}
		if !g.State.Cards[art].Exhausted {
			t.Error("artifact should enter exhausted: the grant belongs to the opponent")
		}
	})

	t.Run("without a grant cards enter play exhausted", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.StartTurn(0)
		newbie := g.AddToHand(NewCard("newbie", Untamed, Creature, Common, WithPower(3)), 0)
		if _, err := g.PlayCreature(0, handIdxByID(g, 0, newbie), false); err != nil {
			t.Fatalf("Playcreature: %v", err)
		}
		if !g.State.Cards[newbie].Exhausted {
			t.Error("creature should enter play exhausted without a grant")
		}
	})

	t.Run(
		"a pool-gated, house-filtered grant readies only off-house cards while rich",
		func(t *testing.T) {
			// Fandangle: while you have 4+ Æmber, your non-Untamed creatures enter ready.
			grant := EntersReadyGrant{
				Type:        Creature,
				MinAember:   4,
				ExceptHouse: Untamed,
			}
			g := NewGame("A", "B", 1)
			g.StartTurn(0)
			granter := g.AddToHand(
				NewCard("Fandangle", Untamed, Creature, Common,
					WithPower(3), WithFriendlyEntersPlayReady(grant)),
				0,
			)
			g.putIntoPlay(granter, 0)

			// Too little Æmber: the grant is dormant even for an off-house creature.
			g.SetAember(0, 3)
			if g.entersPlayReady(0, Creature, Logos) {
				t.Error("grant should be dormant below the Æmber threshold")
			}

			// Enough Æmber, off-house creature: readied.
			g.SetAember(0, 4)
			if !g.entersPlayReady(0, Creature, Logos) {
				t.Error("a non-Untamed creature should be readied at 4 Æmber")
			}

			// Enough Æmber, but an Untamed creature is excluded.
			if g.entersPlayReady(0, Creature, Untamed) {
				t.Error("an Untamed creature should be excluded by the house filter")
			}
		},
	)

	t.Run("renders the printed line for each grant", func(t *testing.T) {
		creat := &CardDefinition{
			Name:             "Duskwitch",
			EntersReadyGrant: EntersReadyGrant{Type: Creature},
		}
		if rules := cardRules(creat, false); len(rules) != 1 ||
			rules[0] != "Your creatures enter play ready." {
			t.Errorf("creature rules = %v", rules)
		}
		art := &CardDefinition{
			Name:             "The Curator",
			EntersReadyGrant: EntersReadyGrant{Type: Artifact},
		}
		if rules := cardRules(art, false); len(rules) != 1 ||
			rules[0] != "Friendly artifacts enter play ready." {
			t.Errorf("artifact rules = %v", rules)
		}
		fan := &CardDefinition{Name: "Fandangle", EntersReadyGrant: EntersReadyGrant{
			Type: Creature, MinAember: 4, ExceptHouse: Untamed,
		}}
		if rules := cardRules(fan, false); len(rules) != 1 ||
			rules[0] != "While you have 4 or more Æmber, your non-Untamed creatures enter play ready." {
			t.Errorf("gated rules = %v", rules)
		}
		filtered := &CardDefinition{Name: "Filtered", EntersReadyGrant: EntersReadyGrant{
			Type: Creature, ExceptHouse: Untamed,
		}}
		if rules := cardRules(filtered, false); len(rules) != 1 ||
			rules[0] != "Your non-Untamed creatures enter play ready." {
			t.Errorf("filtered rules = %v", rules)
		}
	})
}

// TestForgeKeyNumberBarred covers the Key Imps' constant rule: a card in play
// that bars a key ordinal stops every player from forging that key, whoever
// controls the card, and renders its line.
func TestForgeKeyNumberBarred(t *testing.T) {
	g := NewGame("A", "B", 1)
	if g.forgeKeyNumberBarred(0) {
		t.Error("no barring card in play should leave forging open")
	}

	def := NewCard("Bronze Key Imp", Dis, Creature, Common,
		WithPower(2), WithRestrictions(Restrictions{NoForgeKeyNumber: 1}))
	g.AddToBattleline(def, 1) // opponent controls the Imp; it still bars player 0

	g.State.Aember[0] = 3 * KeyCost
	keysBefore := g.Keys(0)
	g.forgeKey(0)
	if g.Keys(0) != keysBefore {
		t.Error("the first key should be barred while a NoForgeKeyNumber:1 card is in play")
	}
	if g.State.Aember[0] != 3*KeyCost {
		t.Error("a barred forge should not spend Æmber")
	}

	g.forgeKeyFree(0)
	if g.Keys(0) != keysBefore {
		t.Error("a free forge of a barred ordinal should also be barred")
	}

	// A second key is not barred by the first-key Imp.
	g.State.ForgeCanonicalKeys(0, 1)
	g.forgeKey(0)
	if g.Keys(0) != 2 {
		t.Error("the second key should forge when only the first is barred")
	}

	if !strings.Contains(RenderCardRules(&def), "Players cannot forge their first key.") {
		t.Error("card rules should render the no-forge line")
	}
}

// TestOrdinalWord covers the ordinal words the Key Imps print.
func TestOrdinalWord(t *testing.T) {
	for n, want := range map[int]string{
		1: "first",
		2: "second",
		3: "third",
		4: "4th",
	} {
		if got := ordinalWord(n); got != want {
			t.Errorf("ordinalWord(%d) = %q, want %q", n, got, want)
		}
	}
}

// TestNoForgeWhileAheadOnKeys checks the constant forge bar Heart of the Forest
// imposes: a player leading on keys skips their forge phase, but a tied or trailing
// player forges normally.
func TestNoForgeWhileAheadOnKeys(t *testing.T) {
	heart := func() CardDefinition {
		return NewCard("heart", Untamed, Artifact, Rare,
			WithRestrictions(Restrictions{NoForgeWhileAheadOnKeys: true}))
	}

	t.Run("bars a player who leads on keys", func(t *testing.T) {
		g := started(t)
		g.AddArtifact(heart(), 0)
		g.State.ForgeCanonicalKeys(0, 1)
		g.State.Aember[0] = 3 * KeyCost
		g.forgePhase(0)
		if g.Keys(0) != 1 {
			t.Errorf("keys = %d, want 1 (barred while ahead)", g.Keys(0))
		}
	})

	t.Run("lets a tied player forge", func(t *testing.T) {
		g := started(t)
		g.AddArtifact(heart(), 0)
		g.State.Aember[0] = 3 * KeyCost
		g.forgePhase(0)
		if g.Keys(0) != 1 {
			t.Errorf("keys = %d, want 1 (tied forges)", g.Keys(0))
		}
	})

	t.Run("lets a trailing player forge", func(t *testing.T) {
		g := started(t)
		g.AddArtifact(heart(), 0)
		g.State.ForgeCanonicalKeys(1, 1)
		g.State.Aember[0] = 3 * KeyCost
		g.forgePhase(0)
		if g.Keys(0) != 1 {
			t.Errorf("keys = %d, want 1 (trailing forges)", g.Keys(0))
		}
	})

	t.Run("does not bar without the restriction in play", func(t *testing.T) {
		g := started(t)
		g.State.ForgeCanonicalKeys(0, 1)
		if g.forgeBarredWhileAhead(0) {
			t.Error("should not bar with no Heart in play")
		}
	})

	t.Run("renders its rule", func(t *testing.T) {
		lines := restrictionText(Restrictions{NoForgeWhileAheadOnKeys: true}, false)
		want := "Each player cannot forge keys while they have more forged keys than their opponent."
		found := false
		for _, l := range lines {
			if l == want {
				found = true
			}
		}
		if !found {
			t.Errorf("lines = %v, want to contain %q", lines, want)
		}
	})
}

// ZoneOf reports the out-of-play pile a card sits in and whose it is, for either
// player, and false for a card that is in play or unknown.
func TestZoneOf(t *testing.T) {
	g := started(t)
	cases := []struct {
		name  string
		add   func(CardDefinition, int) LocalID
		owner int
		want  Zone
	}{
		{"discard", g.AddToDiscard, 0, Discard},
		{"hand", g.AddToHand, 0, Hand},
		{"archives", g.AddToArchives, 1, Archives},
		{"deck", g.AddToDeck, 1, Deck},
	}
	for _, c := range cases {
		id := c.add(testCreature(c.name, 1), c.owner)
		if p, z, ok := g.ZoneOf(id); !ok || p != c.owner || z != c.want {
			t.Errorf("ZoneOf(%s) = (%d, %v, %v), want (%d, %v, true)",
				c.name, p, z, ok, c.owner, c.want)
		}
	}
	// A purged card is reported from the purge pile.
	purged := g.Register(testCreature("purged", 1), 0)
	g.State.Purge[0].add(purged)
	if p, z, ok := g.ZoneOf(purged); !ok || p != 0 || z != Purged {
		t.Errorf("ZoneOf(purged) = (%d, %v, %v), want (0, Purged, true)", p, z, ok)
	}
	// A card in play sits in no out-of-play pile.
	if p, z, ok := g.ZoneOf(g.AddToBattleline(testCreature("inplay", 1), 0)); ok {
		t.Errorf("ZoneOf(in-play) = (%d, %v, %v), want ok false", p, z, ok)
	}
}

// AtCheck reports whether a pool of Æmber meets a player's current key cost.
func TestAtCheck(t *testing.T) {
	g := started(t)
	cost := g.CurrentKeyCost(0)
	if g.AtCheck(0, cost-1) {
		t.Errorf("AtCheck(0, %d) = true, want false at cost %d", cost-1, cost)
	}
	if !g.AtCheck(0, cost) {
		t.Errorf("AtCheck(0, %d) = false, want true at cost %d", cost, cost)
	}
}
