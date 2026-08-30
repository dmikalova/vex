package engine

import "testing"

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
	if ids := (Target{Kind: TargetTriggeringCreature}).Select(ctx); len(ids) != 1 || ids[0] != enemy {
		t.Errorf("triggering-creature select = %v", ids)
	}
	if ids := (Target{Kind: TargetEachEnemyCreature}).Select(ctx); len(ids) != 1 || ids[0] != enemy {
		t.Errorf("each-enemy select = %v", ids)
	}
	if ids := (Target{Kind: TargetEachCreature}).Select(ctx); len(ids) != 2 || ids[0] != src || ids[1] != enemy {
		t.Errorf("each-creature select = %v", ids)
	}
	if ids := (Target{Kind: TargetEachArtifact}).Select(ctx); len(ids) != 2 || ids[0] != myArt || ids[1] != enemyArt {
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
	if ids := (Target{Kind: TargetChosenEnemyCreature}).Select(ctx); len(ids) != 1 || ids[0] != enemy {
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
	if ids := (Target{Kind: TargetEachCreature}).Damaged().Select(ctx); len(ids) != 1 || ids[0] != src {
		t.Errorf("damaged filter = %v, want [%d]", ids, src)
	}
	g.State.Cards[src].Damage = 0

	// OnFlank keeps only the leftmost/rightmost creatures of a battleline.
	mid := g.AddToBattleline(testCreature("mid", 1), 0)
	right := g.AddToBattleline(testCreature("right", 1), 0)
	// Player 0's battleline is now [src, mid, right]; only src and right are flanks.
	if ids := (Target{Kind: TargetEachFriendlyCreature}).OnFlank().Select(ctx); len(ids) != 2 || ids[0] != src || ids[1] != right {
		t.Errorf("flank filter = %v, want [%d %d]", ids, src, right)
	}
	// NotOnFlank keeps only the interior creatures (here, just mid).
	if ids := (Target{Kind: TargetEachFriendlyCreature}).NotOnFlank().Select(ctx); len(ids) != 1 || ids[0] != mid {
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
	if got := (Target{Kind: TargetChosenCreature}).NotOnFlank().Text(); got != "a creature that is not on a flank" {
		t.Errorf("not-on-flank text = %q", got)
	}
	if got := (Target{Kind: TargetEachEnemyCreature}).NotOnFlank().Text(); got != "each enemy creature that is not on a flank" {
		t.Errorf("not-on-flank each text = %q", got)
	}
	if got := (Target{Kind: TargetChosenArtifact}).Text(); got != "an artifact" {
		t.Errorf("chosen-artifact text = %q", got)
	}
}

func TestTargetPowerFilters(t *testing.T) {
	g := NewGame("A", "B", 1)
	p2 := g.AddToBattleline(testCreature("p2", 2), 0)
	p4 := g.AddToBattleline(testCreature("p4", 4), 0)
	p6 := g.AddToBattleline(testCreature("p6", 6), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	if ids := (Target{Kind: TargetEachCreature}).PowerAtMost(3).Select(ctx); len(ids) != 1 || ids[0] != p2 {
		t.Errorf("PowerAtMost(3) = %v, want [%d]", ids, p2)
	}
	if ids := (Target{Kind: TargetEachCreature}).PowerAtLeast(5).Select(ctx); len(ids) != 1 || ids[0] != p6 {
		t.Errorf("PowerAtLeast(5) = %v, want [%d]", ids, p6)
	}
	if ids := (Target{Kind: TargetEachCreature}).PowerExactly(4).Select(ctx); len(ids) != 1 || ids[0] != p4 {
		t.Errorf("PowerExactly(4) = %v, want [%d]", ids, p4)
	}

	if got := (Target{Kind: TargetChosenCreature}).PowerAtLeast(5).Text(); got != "a creature with power 5 or higher" {
		t.Errorf("PowerAtLeast text = %q", got)
	}
	if got := (Target{Kind: TargetChosenCreature}).PowerExactly(1).Text(); got != "a creature with power 1" {
		t.Errorf("PowerExactly text = %q", got)
	}
}

func TestTargetUndamagedAndOther(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 3), 0)
	hurt := g.AddToBattleline(testCreature("hurt", 3), 0)
	g.State.Cards[hurt].Damage = 1
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	if ids := (Target{Kind: TargetEachCreature}).Undamaged().Select(ctx); len(ids) != 1 || ids[0] != src {
		t.Errorf("Undamaged filter = %v, want [%d]", ids, src)
	}
	if ids := (Target{Kind: TargetEachCreature}).Other().Select(ctx); len(ids) != 1 || ids[0] != hurt {
		t.Errorf("Other filter = %v, want [%d]", ids, hurt)
	}
	if got := (Target{Kind: TargetEachCreature}).Other().Undamaged().Text(); got != "each other undamaged creature" {
		t.Errorf("other+undamaged text = %q", got)
	}
}

func TestTargetWithAemberAndLeastPowerful(t *testing.T) {
	g := NewGame("A", "B", 1)
	rich := g.AddToBattleline(testCreature("rich", 5), 0)
	g.State.Cards[rich].Amber = 2
	weak := g.AddToBattleline(testCreature("weak", 1), 1)
	g.AddToBattleline(testCreature("mid", 3), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	if ids := (Target{Kind: TargetEachCreature}).WithAember().Select(ctx); len(ids) != 1 || ids[0] != rich {
		t.Errorf("WithAember = %v, want [%d]", ids, rich)
	}
	if got := (Target{Kind: TargetEachCreature}).WithAember().Text(); got != "each creature with Æmber on it" {
		t.Errorf("WithAember text = %q", got)
	}
	if ids := (Target{Kind: TargetEachCreature}).Selector(LeastPowerful).Select(ctx); len(ids) != 1 || ids[0] != weak {
		t.Errorf("LeastPowerful = %v, want [%d]", ids, weak)
	}
	if got := (Target{Kind: TargetEachCreature}).Selector(LeastPowerful).Text(); got != "the least powerful creature" {
		t.Errorf("LeastPowerful text = %q", got)
	}
	// An empty set selects nothing.
	empty := &EffectContext{Resolver: NewGame("A", "B", 1), Controller: 0}
	if ids := (Target{Kind: TargetEachCreature}).Selector(LeastPowerful).Select(empty); ids != nil {
		t.Errorf("LeastPowerful empty = %v, want nil", ids)
	}
}

func TestLeastPowerfulTieChoice(t *testing.T) {
	g := NewGame("A", "B", 1)
	a := g.AddToBattleline(testCreature("a", 2), 1)
	b := g.AddToBattleline(testCreature("b", 2), 1)
	g.AddToBattleline(testCreature("big", 5), 1)
	g.SetChooser(0, idChooser{id: b})
	ctx := &EffectContext{Resolver: g, Controller: 0}
	ids := (Target{Kind: TargetEachEnemyCreature}).Selector(LeastPowerful).Select(ctx)
	if len(ids) != 1 || ids[0] != b {
		t.Errorf("tie choice = %v, want [%d]; a=%d", ids, b, a)
	}
}

func TestTargetKeyword(t *testing.T) {
	g := NewGame("A", "B", 1)
	elusive := g.AddToBattleline(NewCard("elu", Brobnar, Creature, Common, WithKeywords(Elusive)), 0)
	g.AddToBattleline(testCreature("plain", 3), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	if ids := (Target{Kind: TargetEachCreature}).Keyword(Elusive).Select(ctx); len(ids) != 1 || ids[0] != elusive {
		t.Errorf("Keyword(Elusive) = %v, want [%d]", ids, elusive)
	}
	if got := (Target{Kind: TargetEachCreature}).Keyword(Elusive).Text(); got != "each elusive creature" {
		t.Errorf("keyword text = %q", got)
	}
}

func TestTargetOfHouse(t *testing.T) {
	g := NewGame("A", "B", 1)
	mars := g.AddToBattleline(NewCard("m", Mars, Creature, Common, WithPower(3)), 0)
	g.AddToBattleline(NewCard("s", Sanctum, Creature, Common, WithPower(3)), 0)
	ctx := &EffectContext{Resolver: g, Source: mars, Controller: 0}

	ids := (Target{Kind: TargetEachFriendlyCreature}).OfHouse(Mars).Select(ctx)
	if len(ids) != 1 || ids[0] != mars {
		t.Errorf("OfHouse(Mars) = %v, want [%d] (Sanctum creature filtered out)", ids, mars)
	}
	if got := (Target{Kind: TargetEachCreature}).OfHouse(Mars).Text(); got != "each Mars creature" {
		t.Errorf("OfHouse text = %q", got)
	}
}

func TestTargetExceptTrait(t *testing.T) {
	g := NewGame("A", "B", 1)
	agent := g.AddToBattleline(NewCard("a", Mars, Creature, Common, WithPower(3), WithTraits("Agent")), 0)
	martian := g.AddToBattleline(NewCard("m", Mars, Creature, Common, WithPower(3), WithTraits("Martian")), 0)
	ctx := &EffectContext{Resolver: g, Source: agent, Controller: 0}

	ids := (Target{Kind: TargetEachCreature}).OfHouse(Mars).ExceptTrait("Agent").Select(ctx)
	if len(ids) != 1 || ids[0] != martian {
		t.Errorf("ExceptTrait(Agent) = %v, want [%d] (Agent filtered out)", ids, martian)
	}
	if got := (Target{Kind: TargetChosenCreature}).OfHouse(Mars).ExceptTrait("Agent").Text(); got != "a non-Agent trait Mars creature" {
		t.Errorf("ExceptTrait text = %q", got)
	}
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

func TestTargetExceptMostPowerful(t *testing.T) {
	if got := (Target{Kind: TargetEachEnemyCreature}.Selector(ExceptMostPowerful)).Text(); got != "each enemy creature except the most powerful enemy creature" {
		t.Errorf("enemy text = %q", got)
	}
	if got := (Target{Kind: TargetEachFriendlyCreature}.Selector(ExceptMostPowerful)).Text(); got != "each friendly creature except the most powerful friendly creature" {
		t.Errorf("friendly text = %q", got)
	}

	// Unique most-powerful (added after a weaker one so the running max updates):
	// only the most powerful is spared.
	g := NewGame("A", "B", 1)
	weak := g.AddToBattleline(testCreature("weak", 3), 0)
	strong := g.AddToBattleline(testCreature("strong", 7), 0)
	mid := g.AddToBattleline(testCreature("mid", 5), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	got := Target{Kind: TargetEachFriendlyCreature}.Selector(ExceptMostPowerful).Select(ctx)
	if len(got) != 2 || !containsID(got, weak) || !containsID(got, mid) || containsID(got, strong) {
		t.Errorf("select = %v, want [weak mid] (most powerful spared)", got)
	}

	// One creature (or none) is its own most powerful, so nothing is selected.
	g2 := NewGame("A", "B", 1)
	g2.AddToBattleline(testCreature("lone", 3), 0)
	ctx2 := &EffectContext{Resolver: g2, Controller: 0}
	if got := (Target{Kind: TargetEachFriendlyCreature}.Selector(ExceptMostPowerful)).Select(ctx2); got != nil {
		t.Errorf("lone select = %v, want nil", got)
	}
	if got := (Target{Kind: TargetEachEnemyCreature}.Selector(ExceptMostPowerful)).Select(ctx2); got != nil {
		t.Errorf("empty select = %v, want nil", got)
	}

	// Tied most-powerful: the controller chooses which to keep.
	g3 := NewGame("A", "B", 1)
	a := g3.AddToBattleline(testCreature("a", 5), 0)
	b := g3.AddToBattleline(testCreature("b", 5), 0)
	small := g3.AddToBattleline(testCreature("small", 2), 0)
	g3.SetChooser(0, orderLastChooser{}) // keep the last tied creature (b)
	ctx3 := &EffectContext{Resolver: g3, Controller: 0}
	got = Target{Kind: TargetEachFriendlyCreature}.Selector(ExceptMostPowerful).Select(ctx3)
	if len(got) != 2 || !containsID(got, a) || !containsID(got, small) || containsID(got, b) {
		t.Errorf("tie select = %v, want [a small] (b kept)", got)
	}

	// A rejected tie choice keeps the first tied creature.
	g4 := NewGame("A", "B", 1)
	first := g4.AddToBattleline(testCreature("first", 5), 0)
	second := g4.AddToBattleline(testCreature("second", 5), 0)
	g4.SetChooser(0, orderRejectChooser{})
	ctx4 := &EffectContext{Resolver: g4, Controller: 0}
	got = Target{Kind: TargetEachFriendlyCreature}.Selector(ExceptMostPowerful).Select(ctx4)
	if len(got) != 1 || got[0] != second || containsID(got, first) {
		t.Errorf("rejected tie select = %v, want [second] (first kept)", got)
	}
}

// containsID reports whether ids contains id.
func containsID(ids []LocalID, id LocalID) bool {
	for _, x := range ids {
		if x == id {
			return true
		}
	}
	return false
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
	if ids := (Target{Kind: TargetChosenOtherFriendlyCreature}).Select(ctx); len(ids) != 1 || ids[0] != other {
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
