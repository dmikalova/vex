package engine

import (
	"slices"
	"testing"
)

// TestTargetKindsAreRealAndRendered walks every target kind so a newly added one
// cannot slip in unrendered. TargetKind has no String, so the rendering under test
// is the phrase a card would print.
func TestTargetKindsAreRealAndRendered(t *testing.T) {
	if len(TargetKinds()) != int(targetKindCount)-1 {
		t.Errorf("TargetKinds len = %d, want %d", len(TargetKinds()), targetKindCount-1)
	}
	for _, k := range TargetKinds() {
		if k == targetUnset || k >= targetKindCount {
			t.Errorf("TargetKinds included the sentinel %d", k)
		}
		target := Target{Kind: k}
		phrase, special := target.specialKindText()
		if !special {
			phrase = target.quantifiedPhrase("creature")
		}
		if phrase == "" {
			t.Errorf("target kind %d renders as empty", k)
		}
	}
}

// TestContextReferencesAreExactlyTheSpecialTextKinds pins the partition
// isContextReference names: a kind resolves from the resolution context if and
// only if it renders as a definite phrase through specialKindText. Three switches
// in target.go and target_select.go draw this line independently, so a new kind
// added to one and not the others fails here rather than rendering "a creature"
// for a card that meant "the chosen creature".
func TestContextReferencesAreExactlyTheSpecialTextKinds(t *testing.T) {
	for _, k := range TargetKinds() {
		target := Target{Kind: k}
		_, special := target.specialKindText()
		if got := isContextReference(k); got != special {
			t.Errorf(
				"target kind %d: isContextReference = %v, specialKindText ok = %v",
				k, got, special,
			)
		}
		// A card the context already holds is never a card the player picks.
		if special && target.isChosen() {
			t.Errorf("target kind %d is both a context reference and chosen", k)
		}
	}
}

func TestTargetSelect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	enemy := g.AddToBattleline(testCreature("enemy", 1), 1)
	myArt := g.AddArtifact(NewCard("myrelic", Brobnar, Artifact, Rare), 0)
	enemyArt := g.AddArtifact(NewCard("enemyrelic", Brobnar, Artifact, Rare), 1)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	if ids := (Target{Kind: TargetThisCreature}).Select(ctx); len(ids) != 1 || ids[0] != src {
		t.Errorf("this-creature select = %v", ids)
	}
	if ids := (Target{Kind: TargetTriggeringCreature}).Select(ctx); ids != nil {
		t.Errorf("triggering-creature without It should be nil, got %v", ids)
	}
	ctx.HasIt, ctx.It = true, enemy
	if ids := (Target{Kind: TargetTriggeringCreature}).Select(
		ctx,
	); len(ids) != 1 ||
		ids[0] != enemy {
		t.Errorf("triggering-creature select = %v", ids)
	}
	if ids := (Target{Kind: TargetEachEnemyCreature}).Select(
		ctx,
	); len(ids) != 1 ||
		ids[0] != enemy {
		t.Errorf("each-enemy select = %v", ids)
	}
	if ids := (Target{Kind: TargetEachCreature}).Select(
		ctx,
	); len(ids) != 2 || ids[0] != src ||
		ids[1] != enemy {
		t.Errorf("each-creature select = %v", ids)
	}
	if ids := (Target{Kind: TargetEachArtifact}).Select(
		ctx,
	); len(ids) != 2 || ids[0] != myArt ||
		ids[1] != enemyArt {
		t.Errorf("each-artifact select = %v", ids)
	}
	if ids := (Target{Kind: TargetKind(99)}).Select(ctx); ids != nil {
		t.Errorf("default select should be nil, got %v", ids)
	}

	// TargetChosenCreature asks the chooser to pick one creature from either side.
	if ids := (Target{Kind: TargetChosenCreature}).Select(ctx); len(ids) != 1 || ids[0] != src {
		t.Errorf("chosen-creature (first chooser) = %v, want [%d]", ids, src)
	}
	g.SetChooser(0, orderRejectChooser{})
	if ids := (Target{Kind: TargetChosenCreature}).Select(ctx); ids != nil {
		t.Errorf("chosen-creature (reject) = %v, want nil", ids)
	}
	empty := &EffectContext{Resolver: NewGame("A", "B", 1), Controller: 0}
	if ids := (Target{Kind: TargetChosenCreature}).Select(empty); ids != nil {
		t.Errorf("chosen-creature (no candidates) = %v, want nil", ids)
	}

	// TargetChosenEnemyCreature only offers enemy creatures. With a single enemy
	// the choice is forced and taken automatically — a chooser that would decline
	// is never consulted.
	g.SetChooser(0, orderRejectChooser{})
	if ids := (Target{Kind: TargetChosenEnemyCreature}).Select(
		ctx,
	); len(ids) != 1 ||
		ids[0] != enemy {
		t.Errorf("single-candidate chosen-enemy = %v, want [%d] (auto-selected)", ids, enemy)
	}
	// With two enemies the chooser decides, and may decline.
	g.AddToBattleline(testCreature("enemy2", 1), 1)
	if ids := (Target{Kind: TargetChosenEnemyCreature}).Select(ctx); ids != nil {
		t.Errorf("two-candidate chosen-enemy (reject) = %v, want nil", ids)
	}
	g.SetChooser(0, nil)

	// Damaged filter keeps only creatures with damage on them.
	g.State.Cards[src].Damage = 1
	if ids := (Target{Kind: TargetEachCreature}).Damaged().
		Select(ctx); len(ids) != 1 ||
		ids[0] != src {
		t.Errorf("damaged filter = %v, want [%d]", ids, src)
	}
	g.State.Cards[src].Damage = 0

	// OnFlank keeps only the leftmost/rightmost creatures of a battleline.
	mid := g.AddToBattleline(testCreature("mid", 1), 0)
	right := g.AddToBattleline(testCreature("right", 1), 0)
	// Player 0's battleline is now [src, mid, right]; only src and right are flanks.
	if ids := (Target{Kind: TargetEachFriendlyCreature}).OnFlank().
		Select(ctx); len(ids) != 2 || ids[0] != src ||
		ids[1] != right {
		t.Errorf("flank filter = %v, want [%d %d]", ids, src, right)
	}
	// NotOnFlank keeps only the interior creatures (here, just mid).
	if ids := (Target{Kind: TargetEachFriendlyCreature}).NotOnFlank().
		Select(ctx); len(ids) != 1 ||
		ids[0] != mid {
		t.Errorf("not-on-flank filter = %v, want [%d]", ids, mid)
	}
}

func TestTargetText(t *testing.T) {
	cases := map[TargetKind]string{
		TargetThisCreature:         SelfName,
		TargetTriggeringCreature:   "it",
		TargetEachCreature:         "each creature",
		TargetEachFriendlyCreature: "each friendly creature",
		TargetEachEnemyCreature:    "each enemy creature",
		TargetEachArtifact:         "each artifact",
		TargetTheChosenCreature:    "the chosen creature",
		TargetKind(99):             "a creature",
	}
	for kind, want := range cases {
		if got := (Target{Kind: kind}).Text(); got != want {
			t.Errorf("Text(%d) = %q, want %q", kind, got, want)
		}
	}
	if got := (Target{Kind: TargetChosenEnemyCreature}).Text(); got != "an enemy creature" {
		t.Errorf("chosen-enemy text = %q", got)
	}
	if got := (Target{Kind: TargetChosenCreature}).Damaged().Text(); got != "a damaged creature" {
		t.Errorf("damaged text = %q", got)
	}
	if got := (Target{Kind: TargetChosenCreature}).OnFlank().Text(); got != "a flank creature" {
		t.Errorf("flank text = %q", got)
	}
	if got := (Target{Kind: TargetChosenCreature}).NotOnFlank().
		Text(); got != "a creature that is not on a flank" {
		t.Errorf("not-on-flank text = %q", got)
	}
	if got := (Target{Kind: TargetEachEnemyCreature}).NotOnFlank().
		Text(); got != "each enemy creature that is not on a flank" {
		t.Errorf("not-on-flank each text = %q", got)
	}
	if got := (Target{Kind: TargetChosenArtifact}).Text(); got != "an artifact" {
		t.Errorf("chosen-artifact text = %q", got)
	}
	if got := (Target{Kind: TargetEachCreature}.House(exceptHouse(Mars))).
		Text(); got != "each non-Mars creature" {
		t.Errorf("except-house text = %q", got)
	}
	if got := (Target{Kind: TargetCreatureFought}).Text(); got !=
		"the creature "+SelfName+" fought" {
		t.Errorf("creature-fought text = %q", got)
	}
	if got := (Target{Kind: TargetTheFoughtCreature}).Text(); got != "the fought creature" {
		t.Errorf("fought-creature text = %q", got)
	}
	if got := (Target{Kind: TargetTheFoughtCreature}.NeighborsOf()).
		Text(); got != "each neighbor of the fought creature" {
		t.Errorf("fought-creature neighbors text = %q", got)
	}
	if got := (Target{Kind: TargetTheChosenCreature}).Text(); got != "the chosen creature" {
		t.Errorf("chosen-creature text = %q", got)
	}
	if got := (Target{Kind: TargetTheChosenCreature}.NeighborsOf()).
		Text(); got != "each neighbor of the chosen creature" {
		t.Errorf("chosen-creature neighbors text = %q", got)
	}
}

func TestTargetSharingTrait(t *testing.T) {
	if got := (Target{Kind: TargetEachCreature}).SharingTrait().
		Text(); got != "each creature that shares a trait with it" {
		t.Errorf("shares-trait text = %q", got)
	}

	g := NewGame("A", "B", 1)
	kin := g.AddToBattleline(testCreature("kin", 3, WithTraits(Beast)), 0)
	prey := g.AddToBattleline(testCreature("prey", 5, WithTraits(Beast)), 1)
	g.AddToBattleline(testCreature("spared", 5, WithTraits(Robot)), 1)
	target := Target{Kind: TargetEachCreature}.SharingTrait()

	// Without a context card the filter matches nothing.
	noIt := &EffectContext{Resolver: g, Controller: 0}
	if ids := target.Select(noIt); len(ids) != 0 {
		t.Errorf("shares-trait without It = %v, want empty", ids)
	}

	// With the Beast in context, only trait-sharing creatures pass.
	ctx := &EffectContext{Resolver: g, Controller: 0, It: kin, HasIt: true}
	ids := target.Select(ctx)
	if len(ids) != 2 || ids[0] != kin || ids[1] != prey {
		t.Errorf("shares-trait select = %v, want [%d %d]", ids, kin, prey)
	}
}

func TestTargetPowerFilters(t *testing.T) {
	g := NewGame("A", "B", 1)
	p2 := g.AddToBattleline(testCreature("p2", 2), 0)
	p4 := g.AddToBattleline(testCreature("p4", 4), 0)
	p6 := g.AddToBattleline(testCreature("p6", 6), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	if ids := (Target{Kind: TargetEachCreature}).PowerAtMost(3).
		Select(ctx); len(ids) != 1 ||
		ids[0] != p2 {
		t.Errorf("PowerAtMost(3) = %v, want [%d]", ids, p2)
	}
	if ids := (Target{Kind: TargetEachCreature}).PowerAtLeast(5).
		Select(ctx); len(ids) != 1 ||
		ids[0] != p6 {
		t.Errorf("PowerAtLeast(5) = %v, want [%d]", ids, p6)
	}
	if ids := (Target{Kind: TargetEachCreature}).PowerExactly(4).
		Select(ctx); len(ids) != 1 ||
		ids[0] != p4 {
		t.Errorf("PowerExactly(4) = %v, want [%d]", ids, p4)
	}

	if got := (Target{Kind: TargetChosenCreature}).PowerAtLeast(5).
		Text(); got != "a creature with power 5 or higher" {
		t.Errorf("PowerAtLeast text = %q", got)
	}
	if got := (Target{Kind: TargetChosenCreature}).PowerExactly(1).
		Text(); got != "a creature with power 1" {
		t.Errorf("PowerExactly text = %q", got)
	}
}

func TestTargetPowerParity(t *testing.T) {
	g := NewGame("A", "B", 1)
	p2 := g.AddToBattleline(testCreature("p2", 2), 0)
	p3 := g.AddToBattleline(testCreature("p3", 3), 0)
	p4 := g.AddToBattleline(testCreature("p4", 4), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	if ids := (Target{Kind: TargetEachCreature}).OddPower().
		Select(ctx); len(ids) != 1 ||
		ids[0] != p3 {
		t.Errorf("OddPower = %v, want [%d]", ids, p3)
	}
	if ids := (Target{Kind: TargetEachCreature}).EvenPower().
		Select(ctx); len(ids) != 2 ||
		ids[0] != p2 ||
		ids[1] != p4 {
		t.Errorf("EvenPower = %v, want [%d %d]", ids, p2, p4)
	}
	if got := (Target{Kind: TargetEachCreature}).OddPower().
		Text(); got != "each creature with odd power" {
		t.Errorf("OddPower text = %q", got)
	}
	if got := (Target{Kind: TargetEachCreature}).EvenPower().
		Text(); got != "each creature with even power" {
		t.Errorf("EvenPower text = %q", got)
	}
}

func TestTargetUndamagedAndOther(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 3), 0)
	hurt := g.AddToBattleline(testCreature("hurt", 3), 0)
	g.State.Cards[hurt].Damage = 1
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	if ids := (Target{Kind: TargetEachCreature}).Undamaged().
		Select(ctx); len(ids) != 1 ||
		ids[0] != src {
		t.Errorf("Undamaged filter = %v, want [%d]", ids, src)
	}
	if ids := (Target{Kind: TargetEachCreature}).Other().
		Select(ctx); len(ids) != 1 ||
		ids[0] != hurt {
		t.Errorf("Other filter = %v, want [%d]", ids, hurt)
	}
	if got := (Target{Kind: TargetEachCreature}).Other().
		Undamaged().
		Text(); got != "each other undamaged creature" {
		t.Errorf("other+undamaged text = %q", got)
	}
}

func TestTargetReady(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 3), 0)
	spent := g.AddToBattleline(testCreature("spent", 3), 0)
	g.State.Cards[spent].Exhausted = true
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	if ids := (Target{Kind: TargetEachCreature}).Ready().
		Select(ctx); len(ids) != 1 ||
		ids[0] != src {
		t.Errorf("Ready filter = %v, want [%d]", ids, src)
	}
	if got := (Target{Kind: TargetEachFriendlyCreature}).Ready().
		Text(); got != "each friendly ready creature" {
		t.Errorf("ready text = %q", got)
	}
}

func TestTargetWithAemberAndLeastPowerful(t *testing.T) {
	g := NewGame("A", "B", 1)
	rich := g.AddToBattleline(testCreature("rich", 5), 0)
	g.State.Cards[rich].Amber = 2
	weak := g.AddToBattleline(testCreature("weak", 1), 1)
	g.AddToBattleline(testCreature("mid", 3), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	if ids := (Target{Kind: TargetEachCreature}).WithAember().
		Select(ctx); len(ids) != 1 ||
		ids[0] != rich {
		t.Errorf("WithAember = %v, want [%d]", ids, rich)
	}
	if got := (Target{Kind: TargetEachCreature}).WithAember().
		Text(); got != "each creature with Æmber on it" {
		t.Errorf("WithAember text = %q", got)
	}
	if ids := (Target{Kind: TargetEachCreature}).Refine(LeastPowerful).
		Select(ctx); len(ids) != 1 ||
		ids[0] != weak {
		t.Errorf("LeastPowerful = %v, want [%d]", ids, weak)
	}
	if got := (Target{Kind: TargetEachCreature}).Refine(LeastPowerful).
		Text(); got != "the least powerful creature" {
		t.Errorf("LeastPowerful text = %q", got)
	}
	// An empty set selects nothing.
	empty := &EffectContext{Resolver: NewGame("A", "B", 1), Controller: 0}
	if ids := (Target{Kind: TargetEachCreature}).Refine(LeastPowerful).Select(empty); ids != nil {
		t.Errorf("LeastPowerful empty = %v, want nil", ids)
	}
}

func TestTargetWithoutAember(t *testing.T) {
	g := NewGame("A", "B", 1)
	rich := g.AddToBattleline(testCreature("rich", 5), 0)
	g.State.Cards[rich].Amber = 2
	bare := g.AddToBattleline(testCreature("bare", 3), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	if ids := (Target{Kind: TargetEachCreature}).WithoutAember().
		Select(ctx); len(ids) != 1 || ids[0] != bare {
		t.Errorf("WithoutAember = %v, want [%d]", ids, bare)
	}
	if got := (Target{Kind: TargetChosenCreature}).WithoutAember().
		Text(); got != "a creature with no Æmber on it" {
		t.Errorf("WithoutAember text = %q", got)
	}
}

func TestTargetWithoutBonusIcons(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.AddToBattleline(testCreature("iconed", 5, WithBonus(BonusAember)), 0)
	bare := g.AddToBattleline(testCreature("bare", 3), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	if ids := (Target{Kind: TargetEachCreature}).WithoutBonusIcons().
		Select(ctx); len(ids) != 1 || ids[0] != bare {
		t.Errorf("WithoutBonusIcons = %v, want [%d]", ids, bare)
	}
	if got := (Target{Kind: TargetChosenCreature}).WithoutBonusIcons().
		Text(); got != "a creature with no bonus icons" {
		t.Errorf("WithoutBonusIcons text = %q", got)
	}
}

func TestLeastPowerfulTieChoice(t *testing.T) {
	g := NewGame("A", "B", 1)
	a := g.AddToBattleline(testCreature("a", 2), 1)
	b := g.AddToBattleline(testCreature("b", 2), 1)
	g.AddToBattleline(testCreature("big", 5), 1)
	g.SetChooser(0, idChooser{id: b})
	ctx := &EffectContext{Resolver: g, Controller: 0}
	ids := (Target{Kind: TargetEachEnemyCreature}).Refine(LeastPowerful).Select(ctx)
	if len(ids) != 1 || ids[0] != b {
		t.Errorf("tie choice = %v, want [%d]; a=%d", ids, b, a)
	}
}

func TestAnyOfPowerTiers(t *testing.T) {
	tiers := AnyOf(LowestPower, HighestPower)

	// Text names both extremes.
	if got := (Target{Kind: TargetEachCreature}).Refine(tiers).
		Text(); got != "each creature with the lowest or highest power" {
		t.Errorf("text = %q", got)
	}

	// An empty set selects nothing.
	empty := &EffectContext{Resolver: NewGame("A", "B", 1), Controller: 0}
	if ids := (Target{Kind: TargetEachCreature}).Refine(tiers).
		Select(empty); ids != nil {
		t.Errorf("empty = %v, want nil", ids)
	}

	// The extremes are kept (ties included) and the middle dropped. The middle
	// creature is added first so a later, lower-power creature exercises the
	// new-minimum branch.
	g := NewGame("A", "B", 1)
	g.AddToBattleline(testCreature("mid", 4), 0)
	lowA := g.AddToBattleline(testCreature("lowA", 2), 0)
	lowB := g.AddToBattleline(testCreature("lowB", 2), 1)
	high := g.AddToBattleline(testCreature("high", 6), 1)
	got := (Target{Kind: TargetEachCreature}).Refine(tiers).
		Select(&EffectContext{Resolver: g, Controller: 0})
	if len(got) != 3 || !containsID(got, lowA) || !containsID(got, lowB) ||
		!containsID(got, high) {
		t.Errorf("AnyOf(LowestPower, HighestPower) = %v, want [%d %d %d]", got, lowA, lowB, high)
	}

	// When every creature shares one power the whole set is both extremes.
	g1 := NewGame("A", "B", 1)
	only := g1.AddToBattleline(testCreature("only", 3), 0)
	all := (Target{Kind: TargetEachCreature}).Refine(tiers).
		Select(&EffectContext{Resolver: g1, Controller: 0})
	if len(all) != 1 || all[0] != only {
		t.Errorf("single-power set = %v, want [%d]", all, only)
	}
}

// framedRank is a Refinement framed around a frame of its own, so AnyOf's
// mismatched-frame path has two framed members that cannot share one frame.
type framedRank struct{ word string }

func (r framedRank) clause(phrase string) string { return framedClauseText(r, phrase) }

func (framedRank) clauseFrame(phrase string) (head, tail string) {
	return phrase + " of", "rank"
}

func (r framedRank) clauseWord() string { return r.word }

func (framedRank) refine(_ *EffectContext, ids []LocalID) []LocalID { return ids }

// TestAnyOfFoldsSharedClauseFrame pins when AnyOf states its frame once. Only
// members framed alike fold; anything else repeats the whole noun phrase, which
// is why the fold is a capability a Refinement opts into rather than a string
// comparison over rendered clauses.
func TestAnyOfFoldsSharedClauseFrame(t *testing.T) {
	each := Target{Kind: TargetEachCreature}

	// One member has nothing to share a frame with.
	if text := each.Refine(AnyOf(LowestPower)).
		Text(); text != "each creature with the lowest power" {
		t.Errorf("single member text = %q", text)
	}

	// A framed refinement still renders its own clause when it stands alone.
	if text := each.Refine(HighestPower).
		Text(); text != "each creature with the highest power" {
		t.Errorf("standalone text = %q", text)
	}

	// Members framed differently each render in full, joined with "and" because
	// each is already a complete noun phrase.
	if text := each.Refine(AnyOf(LowestPower, framedRank{word: "high"})).
		Text(); text != "each creature with the lowest power and each creature of high rank" {
		t.Errorf("mismatched frame text = %q", text)
	}
}

func TestMostPowerfulN(t *testing.T) {
	// Text pluralizes the noun.
	if got := (Target{Kind: TargetEachCreature}).Refine(MostPowerfulN(3)).
		Text(); got != "the 3 most powerful creatures" {
		t.Errorf("text = %q", got)
	}

	// The singular MostPowerful reads without a count.
	if got := (Target{Kind: TargetEachCreature}).Refine(MostPowerful).
		Text(); got != "the most powerful creature" {
		t.Errorf("singular text = %q", got)
	}

	// Fewer creatures than n keeps them all.
	g0 := NewGame("A", "B", 1)
	g0.AddToBattleline(testCreature("only", 3), 1)
	ids := (Target{Kind: TargetEachEnemyCreature}).Refine(MostPowerfulN(3)).
		Select(&EffectContext{Resolver: g0, Controller: 0})
	if len(ids) != 1 {
		t.Errorf("MostPowerfulN(3) of one creature = %v, want the single creature", ids)
	}

	// A clean cutoff: the tied group exactly fills the last slot.
	g1 := NewGame("A", "B", 1)
	a := g1.AddToBattleline(testCreature("a", 5), 1)
	b := g1.AddToBattleline(testCreature("b", 4), 1)
	c := g1.AddToBattleline(testCreature("c", 3), 1)
	g1.AddToBattleline(testCreature("d", 2), 1)
	got := (Target{Kind: TargetEachEnemyCreature}).Refine(MostPowerfulN(3)).
		Select(&EffectContext{Resolver: g1, Controller: 0})
	if len(got) != 3 || !containsID(got, a) || !containsID(got, b) || !containsID(got, c) {
		t.Errorf("MostPowerfulN(3) = %v, want the top three [%d %d %d]", got, a, b, c)
	}

	// A tie at the cutoff: the controller chooses which tied creature to include.
	g2 := NewGame("A", "B", 1)
	top := g2.AddToBattleline(testCreature("top", 5), 1)
	t1 := g2.AddToBattleline(testCreature("t1", 3), 1)
	t2 := g2.AddToBattleline(testCreature("t2", 3), 1)
	g2.AddToBattleline(testCreature("t3", 3), 1)
	g2.SetChooser(0, idChooser{id: t2})
	chosen := (Target{Kind: TargetEachEnemyCreature}).Refine(MostPowerfulN(2)).
		Select(&EffectContext{Resolver: g2, Controller: 0})
	if len(chosen) != 2 || !containsID(chosen, top) || !containsID(chosen, t2) {
		t.Errorf("MostPowerfulN(2) tie = %v, want [%d %d]; t1=%d", chosen, top, t2, t1)
	}

	// A declined tie choice falls back to the first tied creature.
	g3 := NewGame("A", "B", 1)
	hi := g3.AddToBattleline(testCreature("hi", 5), 1)
	lo1 := g3.AddToBattleline(testCreature("lo1", 3), 1)
	g3.AddToBattleline(testCreature("lo2", 3), 1)
	g3.AddToBattleline(testCreature("lo3", 3), 1)
	g3.SetChooser(0, orderRejectChooser{})
	fallback := (Target{Kind: TargetEachEnemyCreature}).Refine(MostPowerfulN(2)).
		Select(&EffectContext{Resolver: g3, Controller: 0})
	if len(fallback) != 2 || !containsID(fallback, hi) || !containsID(fallback, lo1) {
		t.Errorf("declined tie = %v, want [%d %d]", fallback, hi, lo1)
	}
}

// TestCandidatesAppliesRefinement covers candidates' refinement branch: it narrows
// to the refined set without making the choice.
func TestCandidatesAppliesRefinement(t *testing.T) {
	g := NewGame("A", "B", 1)
	top := g.AddToBattleline(testCreature("top", 5), 1)
	g.AddToBattleline(testCreature("low", 2), 1)
	got := Target{Kind: TargetEachEnemyCreature}.Refine(MostPowerful).
		candidates(&EffectContext{Resolver: g, Controller: 0})
	if len(got) != 1 || got[0] != top {
		t.Errorf("candidates(MostPowerful) = %v, want [%d]", got, top)
	}
}

func TestHouseWithAtLeast(t *testing.T) {
	// Text renders the "belongs to a house" clause with the threshold.
	want := "each creature that belongs to a house that has 3 or more creatures in play"
	if got := (Target{Kind: TargetEachCreature}).Refine(HouseWithAtLeast(3)).
		Text(); got != want {
		t.Errorf("text = %q", got)
	}

	// Mars has three creatures split across both players; Sanctum has one. The
	// refinement keeps the three Mars creatures and drops the lone Sanctum creature,
	// proving both players' creatures count toward one house's total.
	g := NewGame("A", "B", 1)
	m1 := g.AddToBattleline(NewCard("m1", Mars, Creature, Common, WithPower(3)), 0)
	m2 := g.AddToBattleline(NewCard("m2", Mars, Creature, Common, WithPower(3)), 0)
	m3 := g.AddToBattleline(NewCard("m3", Mars, Creature, Common, WithPower(3)), 1)
	g.AddToBattleline(NewCard("s1", Sanctum, Creature, Common, WithPower(3)), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	got := (Target{Kind: TargetEachCreature}).Refine(HouseWithAtLeast(3)).Select(ctx)
	if len(got) != 3 || !containsID(got, m1) || !containsID(got, m2) || !containsID(got, m3) {
		t.Errorf("HouseWithAtLeast(3) = %v, want the three Mars creatures", got)
	}

	// Raising the threshold above every house's count keeps nothing.
	none := (Target{Kind: TargetEachCreature}).Refine(HouseWithAtLeast(4)).Select(ctx)
	if len(none) != 0 {
		t.Errorf("HouseWithAtLeast(4) = %v, want nothing", none)
	}
}

func TestWithoutSharedTrait(t *testing.T) {
	// Text renders the "does not share a trait" clause.
	want := "each creature that does not share a trait with another creature in its controller's battleline"
	if got := (Target{Kind: TargetEachCreature}).Refine(WithoutSharedTrait()).
		Text(); got != want {
		t.Errorf("text = %q", got)
	}

	// P0 has two Beasts (they share a trait, so neither is a loner) and one
	// Human whose only trait-sharer sits in the ENEMY battleline. P1 has that
	// lone Human. The refinement keeps the two Humans (each a loner in its own
	// battleline) and drops the two Beasts.
	g := NewGame("A", "B", 1)
	g.AddToBattleline(NewCard("b1", Mars, Creature, Common, WithPower(3), WithTraits(Beast)), 0)
	g.AddToBattleline(NewCard("b2", Mars, Creature, Common, WithPower(3), WithTraits(Beast)), 0)
	h0 := g.AddToBattleline(
		NewCard("h0", Mars, Creature, Common, WithPower(3), WithTraits(Human)),
		0,
	)
	h1 := g.AddToBattleline(
		NewCard("h1", Sanctum, Creature, Common, WithPower(3), WithTraits(Human)),
		1,
	)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	got := (Target{Kind: TargetEachCreature}).Refine(WithoutSharedTrait()).Select(ctx)
	if len(got) != 2 || !containsID(got, h0) || !containsID(got, h1) {
		t.Errorf("WithoutSharedTrait = %v, want the two Human loners", got)
	}
}

func TestTargetKeyword(t *testing.T) {
	g := NewGame("A", "B", 1)
	elusive := g.AddToBattleline(
		NewCard("elu", Brobnar, Creature, Common, WithKeywords(Elusive)),
		0,
	)
	g.AddToBattleline(testCreature("plain", 3), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	if ids := (Target{Kind: TargetEachCreature}).Keyword(Elusive).
		Select(ctx); len(ids) != 1 ||
		ids[0] != elusive {
		t.Errorf("Keyword(Elusive) = %v, want [%d]", ids, elusive)
	}
	if got := (Target{Kind: TargetEachCreature}).Keyword(Elusive).
		Text(); got != "each elusive creature" {
		t.Errorf("keyword text = %q", got)
	}
}

func TestTargetOfHouse(t *testing.T) {
	g := NewGame("A", "B", 1)
	mars := g.AddToBattleline(NewCard("m", Mars, Creature, Common, WithPower(3)), 0)
	g.AddToBattleline(NewCard("s", Sanctum, Creature, Common, WithPower(3)), 0)
	ctx := &EffectContext{Resolver: g, Source: mars, Controller: 0}

	ids := (Target{Kind: TargetEachFriendlyCreature}).House(namedHouse(Mars)).Select(ctx)
	if len(ids) != 1 || ids[0] != mars {
		t.Errorf("OfHouse(Mars) = %v, want [%d] (Sanctum creature filtered out)", ids, mars)
	}
	if got := (Target{Kind: TargetEachCreature}).House(namedHouse(Mars)).
		Text(); got != "each Mars creature" {
		t.Errorf("OfHouse text = %q", got)
	}
}

// TestTargetOfActiveHouse covers Techivore Pulpate's target: only artifacts of
// the player's active house are selected, and the phrase reads "of that house".
func TestTargetOfActiveHouse(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.State.ActiveHouse = Mars
	mars := g.AddArtifact(NewCard("m", Mars, Artifact, Common), 0)
	g.AddArtifact(NewCard("s", Sanctum, Artifact, Common), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	ids := (Target{Kind: TargetEachArtifact}).House(activeHouse).Select(ctx)
	if len(ids) != 1 || ids[0] != mars {
		t.Errorf("OfActiveHouse = %v, want [%d] (Sanctum artifact filtered out)", ids, mars)
	}
	if got := (Target{Kind: TargetEachArtifact}).House(activeHouse).
		Text(); got != "each artifact of that house" {
		t.Errorf("OfActiveHouse text = %q", got)
	}
}

// TestTargetMatchingAny covers EMP Blast's target: house and trait disjoin into
// one set, so a Mars Robot is one member of it rather than a member of two
// sequenced targets — which matters because the same renderer produces "deal 1
// damage to each Mars or Robot creature", where being hit twice would be wrong.
func TestTargetMatchingAny(t *testing.T) {
	g := NewGame("A", "B", 1)
	martian := g.AddToBattleline(
		NewCard("m", Mars, Creature, Common, WithPower(3), WithTraits(Martian)), 0)
	robot := g.AddToBattleline(
		NewCard("r", Logos, Creature, Common, WithPower(3), WithTraits(Robot)), 0)
	marsRobot := g.AddToBattleline(
		NewCard("mr", Mars, Creature, Common, WithPower(3), WithTraits(Robot)), 0)
	neither := g.AddToBattleline(
		NewCard("n", Logos, Creature, Common, WithPower(3)), 0)
	ctx := &EffectContext{Resolver: g, Source: martian, Controller: 0}

	either := (Target{Kind: TargetEachCreature}).House(namedHouse(Mars)).
		WithTrait(Robot).MatchingAny()
	got := either.Select(ctx)
	if len(got) != 3 || containsID(got, neither) {
		t.Errorf("MatchingAny selected %v, want the three matching creatures once each", got)
	}
	hits := 0
	for _, id := range got {
		if id == marsRobot {
			hits++
		}
	}
	if hits != 1 {
		t.Errorf("a creature matching both halves appears %d times in %v, want 1", hits, got)
	}
	if !containsID(got, martian) || !containsID(got, robot) {
		t.Errorf("MatchingAny selected %v, want both halves included", got)
	}
	if text := either.Text(); text != "each Mars or Robot creature" {
		t.Errorf("MatchingAny text = %q", text)
	}

	both := (Target{Kind: TargetEachCreature}).House(namedHouse(Mars)).WithTrait(Robot)
	if conj := both.Select(ctx); len(conj) != 1 || conj[0] != marsRobot {
		t.Errorf("without MatchingAny the axes still conjoin, got %v", conj)
	}
}

// TestTargetMatchingAnyNeedsTwoAxes covers the degenerate case: with only one
// identity axis carrying an adjective there is nothing for "or" to join, so the
// target narrows conjunctively and its text cannot promise a union it does not
// select.
func TestTargetMatchingAnyNeedsTwoAxes(t *testing.T) {
	lone := (Target{Kind: TargetEachCreature}).House(namedHouse(Mars)).MatchingAny()
	if lone.disjoins() {
		t.Error("a house alone has no second axis to disjoin")
	}
	if got := lone.Text(); got != "each Mars creature" {
		t.Errorf("text = %q, want the plain conjunctive noun", got)
	}
	suffix := (Target{Kind: TargetEachCreature}).
		House(HouseMatcher{Kind: MatchChosenHouse}).WithTrait(Robot).MatchingAny()
	if suffix.disjoins() {
		t.Error("a suffix-kind house renders after the noun and cannot be an \"or\" branch")
	}
}

func TestTargetExceptTrait(t *testing.T) {
	g := NewGame("A", "B", 1)
	agent := g.AddToBattleline(
		NewCard("a", Mars, Creature, Common, WithPower(3), WithTraits(Agent)),
		0,
	)
	martian := g.AddToBattleline(
		NewCard("m", Mars, Creature, Common, WithPower(3), WithTraits(Martian)),
		0,
	)
	ctx := &EffectContext{Resolver: g, Source: agent, Controller: 0}

	ids := (Target{Kind: TargetEachCreature}).House(namedHouse(Mars)).ExceptTrait(Agent).Select(ctx)
	if len(ids) != 1 || ids[0] != martian {
		t.Errorf("ExceptTrait(Agent) = %v, want [%d] (Agent filtered out)", ids, martian)
	}
	if got := (Target{Kind: TargetChosenCreature}).House(namedHouse(Mars)).
		ExceptTrait(Agent).
		Text(); got != "a non-Agent Mars creature" {
		t.Errorf("ExceptTrait text = %q", got)
	}
}

// TestOfHouseWithMostCreatures covers Etaromme's target: it keeps only creatures
// of the house with the most creatures in play, counting both battlelines, and
// keeps every tied house eligible on a tie.
func TestOfHouseWithMostCreatures(t *testing.T) {
	if got := (Target{Kind: TargetChosenCreature}).OfHouseWithMostCreatures().
		Text(); got != "a creature of the house with the most creatures in play" {
		t.Errorf("text = %q", got)
	}

	t.Run("keeps only the most populous house", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		// Mars leads with three creatures; Brobnar has two, Dis has one.
		for range 3 {
			g.AddToBattleline(NewCard("m", Mars, Creature, Common, WithPower(3)), 0)
		}
		g.AddToBattleline(NewCard("b", Brobnar, Creature, Common, WithPower(3)), 0)
		g.AddToBattleline(NewCard("b2", Brobnar, Creature, Common, WithPower(3)), 1)
		dis := g.AddToBattleline(NewCard("d", Dis, Creature, Common, WithPower(3)), 1)
		ctx := &EffectContext{Resolver: g, Source: dis, Controller: 0}

		ids := (Target{Kind: TargetEachCreature}).OfHouseWithMostCreatures().Select(ctx)
		if len(ids) != 3 {
			t.Fatalf("selected %v, want the 3 Mars creatures", ids)
		}
		for _, id := range ids {
			if g.House(id) != Mars {
				t.Errorf("selected %d of house %v, want Mars", id, g.House(id))
			}
		}
	})

	t.Run("keeps every tied house on a tie", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		m1 := g.AddToBattleline(NewCard("m", Mars, Creature, Common, WithPower(3)), 0)
		g.AddToBattleline(NewCard("m2", Mars, Creature, Common, WithPower(3)), 0)
		g.AddToBattleline(NewCard("b", Brobnar, Creature, Common, WithPower(3)), 1)
		g.AddToBattleline(NewCard("b2", Brobnar, Creature, Common, WithPower(3)), 1)
		ctx := &EffectContext{Resolver: g, Source: m1, Controller: 0}

		ids := (Target{Kind: TargetEachCreature}).OfHouseWithMostCreatures().Select(ctx)
		if len(ids) != 4 {
			t.Fatalf("selected %v, want all 4 creatures of the two tied houses", ids)
		}
	})
}

func TestNeighbors(t *testing.T) {
	g := NewGame("A", "B", 1)
	a := g.AddToBattleline(testCreature("a", 1), 0)
	b := g.AddToBattleline(testCreature("b", 1), 0)
	c := g.AddToBattleline(testCreature("c", 1), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	// The left flank has only a right neighbor; the right flank only a left one.
	if got := neighbors(ctx, a); len(got) != 1 || got[0] != b {
		t.Errorf("left-flank neighbors = %v, want [%d]", got, b)
	}
	if got := neighbors(ctx, c); len(got) != 1 || got[0] != b {
		t.Errorf("right-flank neighbors = %v, want [%d]", got, b)
	}
	if got := neighbors(ctx, b); len(got) != 2 || got[0] != a || got[1] != c {
		t.Errorf("middle neighbors = %v, want [%d %d]", got, a, c)
	}
	// A card not in a battleline has no neighbors.
	art := g.AddArtifact(exAutocannon(), 0)
	if got := neighbors(ctx, art); got != nil {
		t.Errorf("non-battleline neighbors = %v, want nil", got)
	}
}

// TestTargetEachNeighbor covers the source's-own-neighbors target: it renders as
// "each of <self>'s neighbors" and selects the source's live battleline neighbors.
func TestTargetEachNeighbor(t *testing.T) {
	if got := (Target{Kind: TargetEachNeighbor}).Text(); got != "each of "+SelfName+"'s neighbors" {
		t.Errorf("Text = %q", got)
	}
	g := NewGame("A", "B", 1)
	left := g.AddToBattleline(testCreature("left", 1), 0)
	mid := g.AddToBattleline(testCreature("mid", 1), 0)
	right := g.AddToBattleline(testCreature("right", 1), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0, Source: mid}

	ids := (Target{Kind: TargetEachNeighbor}).Select(ctx)
	if len(ids) != 2 || ids[0] != left || ids[1] != right {
		t.Errorf("Select = %v, want [%d %d]", ids, left, right)
	}
}

// TestTargetInCenter covers the center-of-battleline filter: it renders the
// qualifier and selects only the creature in the center of a battleline (Beware
// the Ides).
func TestTargetInCenter(t *testing.T) {
	want := "a creature in the center of its controller's battleline"
	if got := (Target{Kind: TargetChosenCreature}).InCenter().Text(); got != want {
		t.Errorf("Text = %q, want %q", got, want)
	}
	g := NewGame("A", "B", 1)
	g.AddToBattleline(testCreature("left", 1), 0)
	mid := g.AddToBattleline(testCreature("mid", 1), 0)
	g.AddToBattleline(testCreature("right", 1), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	ids := (Target{Kind: TargetChosenCreature}).InCenter().Select(ctx)
	if len(ids) != 1 || ids[0] != mid {
		t.Errorf("Select = %v, want [%d] (only the center creature)", ids, mid)
	}
}

// TestTargetEachUpgradeOnThis covers the source's-own-upgrades target: it renders
// as "each upgrade on <self>" and selects the upgrades attached to the source
// (Away Team).
func TestTargetEachUpgradeOnThis(t *testing.T) {
	if got := (Target{Kind: TargetEachUpgradeOnThis}).Text(); got != "each upgrade on "+SelfName {
		t.Errorf("Text = %q", got)
	}
	g := NewGame("A", "B", 1)
	host := g.AddToBattleline(testCreature("host", 3), 0)
	up1 := attachUpgrade(g, host, NewCard("coil", Mars, Upgrade, Common))
	up2 := attachUpgrade(g, host, NewCard("plate", Mars, Upgrade, Common))
	ctx := &EffectContext{Resolver: g, Controller: 0, Source: host}

	ids := (Target{Kind: TargetEachUpgradeOnThis}).Select(ctx)
	if len(ids) != 2 || ids[0] != up1 || ids[1] != up2 {
		t.Errorf("Select = %v, want [%d %d]", ids, up1, up2)
	}
}

func TestTargetWithUpgrade(t *testing.T) {
	g := NewGame("A", "B", 1)
	upgraded := g.AddToBattleline(testCreature("up", 3), 0)
	bare := g.AddToBattleline(testCreature("bare", 3), 0)
	attachUpgrade(g, upgraded, NewCard("plating", Mars, Upgrade, Common))
	ctx := &EffectContext{Resolver: g, Controller: 0}

	if ids := (Target{Kind: TargetEachCreature}).WithUpgrade().
		Select(ctx); len(ids) != 1 || ids[0] != upgraded {
		t.Errorf("WithUpgrade = %v, want [%d] (the bare creature %d filtered out)",
			ids, upgraded, bare)
	}
	if got := (Target{Kind: TargetEachCreature}).WithUpgrade().
		Text(); got != "each creature with an upgrade" {
		t.Errorf("with-upgrade text = %q", got)
	}
}

func TestTargetSharesHouseWithNeighbors(t *testing.T) {
	g := NewGame("A", "B", 1)
	left := g.AddToBattleline(NewCard("l", Mars, Creature, Common, WithPower(3)), 0)
	mid := g.AddToBattleline(NewCard("m", Mars, Creature, Common, WithPower(3)), 0)
	right := g.AddToBattleline(NewCard("r", Mars, Creature, Common, WithPower(3)), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	// Every creature shares its house with at least one neighbor.
	if ids := (Target{Kind: TargetEachCreature}).SharesHouseWithNeighbors(1).
		Select(ctx); len(ids) != 3 ||
		ids[0] != left || ids[1] != mid || ids[2] != right {
		t.Errorf("SharesHouseWithNeighbors(1) = %v, want [%d %d %d]", ids, left, mid, right)
	}
	// Only the middle creature has two same-house neighbors; the flanks have one.
	if ids := (Target{Kind: TargetEachCreature}).SharesHouseWithNeighbors(2).
		Select(ctx); len(ids) != 1 || ids[0] != mid {
		t.Errorf("SharesHouseWithNeighbors(2) = %v, want [%d]", ids, mid)
	}
	if got := (Target{Kind: TargetEachCreature}).SharesHouseWithNeighbors(1).
		Text(); got != "each creature that shares a house with at least 1 of its neighbors" {
		t.Errorf("shares-house(1) text = %q", got)
	}
	if got := (Target{Kind: TargetEachCreature}).SharesHouseWithNeighbors(2).
		Text(); got != "each creature that shares a house with 2 of its neighbors" {
		t.Errorf("shares-house(2) text = %q", got)
	}
}

func TestNotMostPowerful(t *testing.T) {
	if got := (Target{Kind: TargetEachEnemyCreature}.Refine(Except(MostPowerful))).Text(); got != "each enemy creature except the most powerful enemy creature" {
		t.Errorf("enemy text = %q", got)
	}
	if got := (Target{Kind: TargetEachFriendlyCreature}.Refine(Except(MostPowerful))).Text(); got != "each friendly creature except the most powerful friendly creature" {
		t.Errorf("friendly text = %q", got)
	}

	// Unique most-powerful (added after a weaker one so the running max updates):
	// only the most powerful is spared.
	g := NewGame("A", "B", 1)
	weak := g.AddToBattleline(testCreature("weak", 3), 0)
	strong := g.AddToBattleline(testCreature("strong", 7), 0)
	mid := g.AddToBattleline(testCreature("mid", 5), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	got := Target{Kind: TargetEachFriendlyCreature}.Refine(Except(MostPowerful)).Select(ctx)
	if len(got) != 2 || !containsID(got, weak) || !containsID(got, mid) || containsID(got, strong) {
		t.Errorf("select = %v, want [weak mid] (most powerful spared)", got)
	}

	// One creature (or none) is its own most powerful, so nothing is selected.
	g2 := NewGame("A", "B", 1)
	g2.AddToBattleline(testCreature("lone", 3), 0)
	ctx2 := &EffectContext{Resolver: g2, Controller: 0}
	if got := (Target{Kind: TargetEachFriendlyCreature}.Refine(Except(MostPowerful))).Select(
		ctx2,
	); got != nil {
		t.Errorf("lone select = %v, want nil", got)
	}
	if got := (Target{Kind: TargetEachEnemyCreature}.Refine(Except(MostPowerful))).Select(
		ctx2,
	); got != nil {
		t.Errorf("empty select = %v, want nil", got)
	}

	// Tied most-powerful: the controller chooses which to keep.
	g3 := NewGame("A", "B", 1)
	a := g3.AddToBattleline(testCreature("a", 5), 0)
	b := g3.AddToBattleline(testCreature("b", 5), 0)
	small := g3.AddToBattleline(testCreature("small", 2), 0)
	g3.SetChooser(0, orderLastChooser{}) // keep the last tied creature (b)
	ctx3 := &EffectContext{Resolver: g3, Controller: 0}
	got = Target{Kind: TargetEachFriendlyCreature}.Refine(Except(MostPowerful)).Select(ctx3)
	if len(got) != 2 || !containsID(got, a) || !containsID(got, small) || containsID(got, b) {
		t.Errorf("tie select = %v, want [a small] (b kept)", got)
	}

	// A rejected tie choice keeps the first tied creature.
	g4 := NewGame("A", "B", 1)
	first := g4.AddToBattleline(testCreature("first", 5), 0)
	second := g4.AddToBattleline(testCreature("second", 5), 0)
	g4.SetChooser(0, orderRejectChooser{})
	ctx4 := &EffectContext{Resolver: g4, Controller: 0}
	got = Target{Kind: TargetEachFriendlyCreature}.Refine(Except(MostPowerful)).Select(ctx4)
	if len(got) != 1 || got[0] != second || containsID(got, first) {
		t.Errorf("rejected tie select = %v, want [second] (first kept)", got)
	}
}

// containsID reports whether ids contains id.
func containsID(ids []LocalID, id LocalID) bool {
	return slices.Contains(ids, id)
}

// TestPowerLessThan covers the refinement Exterminate! Exterminate! uses: keep the
// creatures whose power is below a running count, and render the cardinal clause.
func TestPowerLessThan(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.AddToBattleline(NewCard("m1", Mars, Creature, Common, WithPower(2)), 0)
	g.AddToBattleline(NewCard("m2", Mars, Creature, Common, WithPower(2)), 0)
	weak := g.AddToBattleline(testCreature("weak", 1), 1)
	equal := g.AddToBattleline(testCreature("equal", 2), 1)
	strong := g.AddToBattleline(testCreature("strong", 3), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	// Threshold is the two friendly Mars creatures, so only power < 2 is kept:
	// power == 2 and power 3 both survive.
	limit := CardsInPlay{Player: Controller, Type: Creature, House: namedHouse(Mars)}
	got := Target{Kind: TargetEachEnemyCreature}.Refine(PowerLessThan(limit)).Select(ctx)
	if len(got) != 1 || got[0] != weak || containsID(got, equal) || containsID(got, strong) {
		t.Errorf("PowerLessThan = %v, want [weak]", got)
	}

	if text := (Target{Kind: TargetEachEnemyCreature}).Refine(PowerLessThan(limit)).
		Text(); text !=
		"each enemy creature with power less than the number of friendly Mars creatures you control" {
		t.Errorf("PowerLessThan text = %q", text)
	}

	// An empty set keeps nothing.
	empty := &EffectContext{Resolver: NewGame("A", "B", 1), Controller: 0}
	if ids := (Target{Kind: TargetEachEnemyCreature}).Refine(PowerLessThan(limit)).
		Select(empty); len(
		ids,
	) != 0 {
		t.Errorf("PowerLessThan empty = %v, want none", ids)
	}

	// The SelfHouse sentinel in the count resolves to the card's own house, even
	// though it lives in the refinement's unexported field.
	selfLimit := CardsInPlay{Player: Controller, Type: Creature, House: namedHouse(SelfHouse)}
	resolved := replacedIn(
		(Target{Kind: TargetEachEnemyCreature}).Refine(PowerLessThan(selfLimit)),
		SelfHouse,
		Mars,
	)
	if text := resolved.Text(); text !=
		"each enemy creature with power less than the number of friendly Mars creatures you control" {
		t.Errorf("resolved SelfHouse text = %q", text)
	}
}

// TestPowerLessThanSource covers the source-relative refinement: it keeps the
// creatures whose power is below the source card's own power (Dreadbone Decimus
// destroys a creature with lower power than itself) and renders "... with lower
// power than <self>".
func TestPowerLessThanSource(t *testing.T) {
	g := NewGame("A", "B", 1)
	source := g.AddToBattleline(testCreature("source", 3), 0)
	weak := g.AddToBattleline(testCreature("weak", 1), 1)
	equal := g.AddToBattleline(testCreature("equal", 3), 1)
	strong := g.AddToBattleline(testCreature("strong", 5), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0, Source: source}

	got := Target{Kind: TargetEachEnemyCreature}.Refine(PowerLessThanSource()).Select(ctx)
	if len(got) != 1 || got[0] != weak || containsID(got, equal) || containsID(got, strong) {
		t.Errorf("PowerLessThanSource = %v, want [weak]", got)
	}

	if text := (Target{Kind: TargetChosenEnemyCreature}).Refine(PowerLessThanSource()).
		Text(); text != "an enemy creature with lower power than "+SelfName {
		t.Errorf("PowerLessThanSource text = %q", text)
	}
}

// TestUnionableAxisRefinements covers the trait and power-floor refinements that
// restate a Target axis so AnyOf can union them — Regrettable Meteor destroys the
// Dinosaurs and the power-6-or-higher creatures as one set, so a creature matching
// both halves is destroyed once.
func TestUnionableAxisRefinements(t *testing.T) {
	g := NewGame("A", "B", 1)
	dino := g.AddToBattleline(
		NewCard("dino", Untamed, Creature, Common, WithPower(2), WithTraits(Dinosaur)), 1)
	big := g.AddToBattleline(testCreature("big", 6), 1)
	bigDino := g.AddToBattleline(
		NewCard("bigDino", Untamed, Creature, Common, WithPower(7), WithTraits(Dinosaur)), 1)
	small := g.AddToBattleline(testCreature("small", 3), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	each := Target{Kind: TargetEachEnemyCreature}
	union := each.Refine(AnyOf(OfTrait(Dinosaur), PowerAtLeast(6)))
	got := union.Select(ctx)
	if len(got) != 3 || !containsID(got, dino) || !containsID(got, big) ||
		!containsID(got, bigDino) || containsID(got, small) {
		t.Errorf("union = %v, want dino, big, and bigDino once each", got)
	}

	if text := union.Text(); text !=
		"each enemy Dinosaur creature and each enemy creature with power 6 or higher" {
		t.Errorf("union text = %q", text)
	}

	// Each half also stands on its own.
	if ids := each.Refine(OfTrait(Dinosaur)).Select(ctx); len(ids) != 2 {
		t.Errorf("OfTrait = %v, want the two Dinosaurs", ids)
	}
	if ids := each.Refine(PowerAtLeast(6)).Select(ctx); len(ids) != 2 {
		t.Errorf("PowerAtLeast = %v, want the two power-6+ creatures", ids)
	}
}

// TestQualifyLastNoun covers the one-word phrase, which has no space to insert an
// adjective before.
func TestQualifyLastNoun(t *testing.T) {
	if got := qualifyLastNoun("creature", "Dinosaur"); got != "Dinosaur creature" {
		t.Errorf("qualifyLastNoun = %q", got)
	}
}

func TestTargetChosenOtherFriendly(t *testing.T) {
	if got := (Target{Kind: TargetChosenOtherFriendlyCreature}).Text(); got != "another friendly creature" {
		t.Errorf("text = %q, want %q", got, "another friendly creature")
	}

	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 3), 0)
	other := g.AddToBattleline(testCreature("other", 3), 0)
	g.AddToBattleline(testCreature("enemy", 3), 1)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	// The source is excluded, leaving one candidate that is auto-selected.
	if ids := (Target{Kind: TargetChosenOtherFriendlyCreature}).Select(
		ctx,
	); len(ids) != 1 ||
		ids[0] != other {
		t.Errorf("chosen-other-friendly = %v, want [%d]", ids, other)
	}

	// With two other friendly creatures the chooser decides, and may decline.
	g.AddToBattleline(testCreature("other2", 3), 0)
	g.SetChooser(0, orderRejectChooser{})
	if ids := (Target{Kind: TargetChosenOtherFriendlyCreature}).Select(ctx); ids != nil {
		t.Errorf("chosen-other-friendly (reject) = %v, want nil", ids)
	}

	// A lone source has no other friendly creatures to choose.
	g2 := NewGame("A", "B", 1)
	lone := g2.AddToBattleline(testCreature("lone", 3), 0)
	ctx2 := &EffectContext{Resolver: g2, Source: lone, Controller: 0}
	if ids := (Target{Kind: TargetChosenOtherFriendlyCreature}).Select(ctx2); ids != nil {
		t.Errorf("lone source chosen-other-friendly = %v, want nil", ids)
	}
}

func TestTargetNamed(t *testing.T) {
	g := NewGame("A", "B", 1)
	bear := g.Register(NewCard("Ancient Bear", Untamed, Creature, Common, WithPower(6)), 0)
	g.State.Battleline[0].add(bear)
	other := g.AddToBattleline(testCreature("Chuff Ape", 6), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	tgt := Target{Kind: TargetEachFriendlyCreature}.Named("Ancient Bear")
	// A named card needs no describing, so the name replaces the noun outright.
	if want := "each friendly Ancient Bear"; tgt.Text() != want {
		t.Errorf("text = %q, want %q", tgt.Text(), want)
	}
	got := tgt.Select(ctx)
	if len(got) != 1 || got[0] != bear {
		t.Errorf("selected %v, want just the bear (not %v)", got, other)
	}
}

func TestTargetChosenEnemyCreatureOrArtifact(t *testing.T) {
	tgt := Target{Kind: TargetChosenEnemyCreatureOrArtifact}
	if want := "an enemy creature or artifact"; tgt.Text() != want {
		t.Errorf("text = %q, want %q", tgt.Text(), want)
	}

	g := NewGame("A", "B", 1)
	g.AddToBattleline(testCreature("mine", 3), 0)
	foe := g.AddToBattleline(testCreature("theirs", 3), 1)
	art := g.AddArtifact(NewCard("Their Relic", Logos, Artifact, Common), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	// Both enemy halves are candidates; the chooser takes the first.
	if got := tgt.Select(ctx); len(got) != 1 || got[0] != foe {
		t.Errorf("selected %v, want the enemy creature %v (candidates include %v)", got, foe, art)
	}
}

// TestTargetChosenUpgrade renders "an upgrade" and selects from every upgrade
// attached to any creature in play (Destroy Them All!).
func TestTargetChosenUpgrade(t *testing.T) {
	if got := (Target{Kind: TargetChosenUpgrade}).Text(); got != "an upgrade" {
		t.Errorf("chosen-upgrade text = %q, want %q", got, "an upgrade")
	}
	if !(Target{Kind: TargetChosenUpgrade}).isChosen() {
		t.Error("chosen-upgrade should be a chosen target")
	}

	g := started(t)
	mine := g.AddToBattleline(testCreature("mine", 3), 0)
	theirs := g.AddToBattleline(testCreature("theirs", 3), 1)
	up1 := g.Register(NewCard("up1", Brobnar, Upgrade, Common), 0)
	up2 := g.Register(NewCard("up2", Brobnar, Upgrade, Common), 1)
	g.AttachUpgrade(mine, up1)
	g.AttachUpgrade(theirs, up2)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	if ids := (Target{Kind: TargetChosenUpgrade}).selectBase(ctx); len(ids) != 2 ||
		ids[0] != up1 || ids[1] != up2 {
		t.Errorf("chosen-upgrade selectBase = %v, want [%d %d]", ids, up1, up2)
	}
}

// TestKeepPerSideRefinement covers the KeepPerSide leftover selection: its lead
// and clause wording, an empty board, and keeping the chosen number on each side
// (Unnatural Selection).
func TestKeepPerSideRefinement(t *testing.T) {
	tgt := (Target{Kind: TargetEachCreature}).Refine(KeepPerSide(3))

	// The clause renders the leftover set; the lead renders the pair of choices.
	if got := tgt.Text(); got != "each other creature" {
		t.Errorf("text = %q", got)
	}
	if lead, ok := tgt.leadIn(); !ok ||
		lead != "choose 3 friendly creatures and 3 enemy creatures" {
		t.Errorf("leadIn = %q, %v", lead, ok)
	}

	// An empty board selects nothing.
	empty := &EffectContext{Resolver: NewGame("A", "B", 1), Controller: 0}
	if ids := tgt.Select(empty); ids != nil {
		t.Errorf("empty = %v, want nil", ids)
	}

	// Keep 3 of 4 on each side; the unchosen creature on each side is left over.
	g := NewGame("A", "B", 1)
	f0 := g.AddToBattleline(testCreature("f0", 3), 0)
	f1 := g.AddToBattleline(testCreature("f1", 3), 0)
	f2 := g.AddToBattleline(testCreature("f2", 3), 0)
	f3 := g.AddToBattleline(testCreature("f3", 3), 0)
	e0 := g.AddToBattleline(testCreature("e0", 3), 1)
	e1 := g.AddToBattleline(testCreature("e1", 3), 1)
	e2 := g.AddToBattleline(testCreature("e2", 3), 1)
	e3 := g.AddToBattleline(testCreature("e3", 3), 1)
	g.SetChooser(0, &idQueueChooser{ids: []LocalID{f0, f1, f2, e0, e1, e2}})
	got := tgt.Select(&EffectContext{Resolver: g, Controller: 0})
	if len(got) != 2 || !containsID(got, f3) || !containsID(got, e3) {
		t.Errorf("KeepPerSide(3) leftover = %v, want [%d %d]", got, f3, e3)
	}
}

// TestKeepPerSideRefinementFewerThanKeepCount covers a side no larger than the
// keep count: every creature is kept with no prompt, so nothing is left over.
func TestKeepPerSideRefinementFewerThanKeepCount(t *testing.T) {
	tgt := (Target{Kind: TargetEachCreature}).Refine(KeepPerSide(3))

	// Two friendly and one enemy creature, all at or below the keep count of 3: the
	// choice is vacuous, so no creature is picked and nothing is left over.
	g := NewGame("A", "B", 1)
	g.AddToBattleline(testCreature("f0", 3), 0)
	g.AddToBattleline(testCreature("f1", 3), 0)
	g.AddToBattleline(testCreature("e0", 3), 1)
	spy := &countingChooser{}
	g.SetChooser(0, spy)
	if ids := tgt.Select(&EffectContext{Resolver: g, Controller: 0}); ids != nil {
		t.Errorf("fewer-than-keep leftover = %v, want nil", ids)
	}
	if spy.calls != 0 {
		t.Errorf("vacuous keep prompted %d times, want 0", spy.calls)
	}
}

// TestPortionPerSideRefinement covers the PortionPerSide selection: its wording,
// an empty board, and choosing the fraction of each side, the enemy side first
// (Tertiate).
func TestPortionPerSideRefinement(t *testing.T) {
	tgt := (Target{Kind: TargetEachCreature}).Refine(PortionPerSide(ThirdRoundedUp))

	want := "one third of all enemy creatures and one third of all friendly " +
		"creatures, rounding up each time"
	if got := tgt.Text(); got != want {
		t.Errorf("text = %q", got)
	}

	// A Portion refinement carries no lead.
	if lead, ok := tgt.leadIn(); ok {
		t.Errorf("leadIn = %q, %v, want no lead", lead, ok)
	}

	// An empty board selects nothing.
	empty := &EffectContext{Resolver: NewGame("A", "B", 1), Controller: 0}
	if ids := tgt.Select(empty); ids != nil {
		t.Errorf("empty = %v, want nil", ids)
	}

	// ceil(4/3)=2 enemy (chosen first) and ceil(4/3)=2 friendly are selected.
	g := NewGame("A", "B", 1)
	f0 := g.AddToBattleline(testCreature("f0", 3), 0)
	f1 := g.AddToBattleline(testCreature("f1", 3), 0)
	g.AddToBattleline(testCreature("f2", 3), 0)
	g.AddToBattleline(testCreature("f3", 3), 0)
	e0 := g.AddToBattleline(testCreature("e0", 3), 1)
	e1 := g.AddToBattleline(testCreature("e1", 3), 1)
	g.AddToBattleline(testCreature("e2", 3), 1)
	g.AddToBattleline(testCreature("e3", 3), 1)
	g.SetChooser(0, &idQueueChooser{ids: []LocalID{e0, e1, f0, f1}})
	got := tgt.Select(&EffectContext{Resolver: g, Controller: 0})
	if len(got) != 4 || !containsID(got, e0) || !containsID(got, e1) ||
		!containsID(got, f0) || !containsID(got, f1) {
		t.Errorf("PortionPerSide(third) = %v, want [%d %d %d %d]", got, e0, e1, f0, f1)
	}
}

// TestPortionPerSideRefinementRounding covers ceil(1/3)=1 on a lone-creature side
// with an empty other side.
func TestPortionPerSideRefinementRounding(t *testing.T) {
	tgt := (Target{Kind: TargetEachCreature}).Refine(PortionPerSide(ThirdRoundedUp))

	// ceil(1/3)=1 on the friendly side; the enemy side is empty.
	g := NewGame("A", "B", 1)
	f0 := g.AddToBattleline(testCreature("f0", 3), 0)
	g.SetChooser(0, &idQueueChooser{ids: []LocalID{f0}})
	got := tgt.Select(&EffectContext{Resolver: g, Controller: 0})
	if len(got) != 1 || got[0] != f0 {
		t.Errorf("ceil(1/3) = %v, want [%d]", got, f0)
	}
}

// TestTargetGrantingCard covers the granting-card target: it renders the {card}
// placeholder so the granted-text renderer names the granting card, resolves to
// the card that granted the ability (ctx.Grantor) when one is set, and selects
// nothing when no grantor is set.
func TestTargetGrantingCard(t *testing.T) {
	if got := (Target{Kind: TargetGrantingCard}).Text(); got != CardName {
		t.Errorf("Text() = %q, want %q", got, CardName)
	}

	g := NewGame("A", "B", 1)
	artifact := g.AddArtifact(NewCard("Grantor", StarAlliance, Artifact, Rare), 0)

	withGrantor := &EffectContext{Resolver: g, Grantor: artifact, HasGrantor: true}
	got := Target{Kind: TargetGrantingCard}.selectBase(withGrantor)
	if len(got) != 1 || got[0] != artifact {
		t.Errorf("selectBase with a grantor = %v, want [%d]", got, artifact)
	}

	noGrantor := &EffectContext{Resolver: g}
	if got := (Target{Kind: TargetGrantingCard}).selectBase(noGrantor); got != nil {
		t.Errorf("selectBase without a grantor = %v, want nil", got)
	}
}

// TestMostPowerfulIncludesEarlyReturns covers the tie-inclusive membership test's
// two early returns: an id absent from the set is never included, and a set no
// larger than n includes every member without any power comparison.
func TestMostPowerfulIncludesEarlyReturns(t *testing.T) {
	m := mostPowerfulN{n: 1}
	ctx := &EffectContext{}
	if m.includes(ctx, []LocalID{1, 2}, 3) {
		t.Error("includes: an id absent from the set should be excluded")
	}
	if !m.includes(ctx, []LocalID{5}, 5) {
		t.Error("includes: a single-member set should include its member")
	}
}

// TestItIsAmongNonMembershipRefinement covers couldSelect's fallback path for a
// refinement that is not a membershipRefiner (LeastPowerful): membership falls
// back to the refinement's own concrete selection.
func TestItIsAmongNonMembershipRefinement(t *testing.T) {
	cond := ItIsAmong{
		Target: Target{Kind: TargetEachEnemyCreature}.Refine(LeastPowerful),
		Noun:   FoughtCreature,
	}
	g := started(t)
	weak := g.AddToBattleline(testCreature("weak", 2), 1)
	g.AddToBattleline(testCreature("strong", 6), 1)
	if !cond.Met(&EffectContext{Resolver: g, Controller: 0, It: weak, HasIt: true}) {
		t.Error("the least powerful enemy should be among the candidates")
	}
}
