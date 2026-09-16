package engine

import "testing"

func TestPurgedCardsCount(t *testing.T) {
	g := NewGame("A", "B", 1)
	a := g.AddToBattleline(testCreature("a", 1), 0)
	b := g.AddToBattleline(testCreature("b", 1), 1)
	c := g.AddToBattleline(testCreature("c", 1), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	if got := (PurgedCards{}).Value(ctx); got != 0 {
		t.Errorf("empty purge piles: Value = %d, want 0", got)
	}
	g.purgeFromPlay(a)
	g.purgeFromPlay(b)
	g.purgeFromPlay(c)
	if got := (PurgedCards{}).Value(ctx); got != 3 {
		t.Errorf("across both players: Value = %d, want 3", got)
	}
	if got := (PurgedCards{}).CountText(); got != "purged card" {
		t.Errorf("CountText = %q", got)
	}
}

func TestInPlay(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.AddToBattleline(NewCard("m1", Mars, Creature, Common, WithPower(2)), 0)
	g.AddToBattleline(NewCard("m2", Mars, Creature, Common, WithPower(2)), 0)
	g.AddToBattleline(NewCard("b1", Brobnar, Creature, Common, WithPower(2)), 0)
	g.AddArtifact(NewCard("relic", Brobnar, Artifact, Common), 0)
	g.AddArtifact(NewCard("shard", Brobnar, Artifact, Common, WithTraits(Shard)), 0)
	g.AddToBattleline(NewCard("foe", Shadows, Creature, Common, WithPower(2)), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	// Value across type and house filters, and the opposing side.
	values := []struct {
		name string
		in   InPlay
		want int
	}{
		{"friendly creatures", InPlay{Player: Controller, Type: Creature}, 3},
		{"friendly Mars creatures", InPlay{Player: Controller, Type: Creature, House: Mars}, 2},
		{"friendly artifacts", InPlay{Player: Controller, Type: Artifact}, 2},
		{"friendly cards, any type", InPlay{Player: Controller}, 5},
		{"friendly Shards", InPlay{Player: Controller, Trait: Shard}, 1},
		{"enemy creatures", InPlay{Player: Opponent, Type: Creature}, 1},
	}
	for _, tc := range values {
		if got := tc.in.Value(ctx); got != tc.want {
			t.Errorf("%s: Value = %d, want %d", tc.name, got, tc.want)
		}
	}

	// CountText — the singular "for each" noun.
	texts := []struct {
		in   InPlay
		want string
	}{
		{InPlay{Player: Controller, Type: Creature}, "friendly creature in play"},
		{InPlay{Player: Controller, Type: Creature, House: Mars}, "friendly Mars creature"},
		{InPlay{Player: Controller, Trait: Shard}, "friendly Shard"},
		{InPlay{Player: Controller, Type: Creature, Trait: Thief}, "friendly Thief creature"},
		{InPlay{Player: Controller, Type: Artifact, Trait: Shard}, "friendly Shard artifact"},
		{InPlay{Player: Opponent, Type: Creature}, "enemy creature in play"},
		{InPlay{Player: Controller, Type: Artifact}, "friendly artifact in play"},
		{InPlay{Player: Controller}, "friendly card in play"},
	}
	for _, tc := range texts {
		if got := tc.in.CountText(); got != tc.want {
			t.Errorf("CountText = %q, want %q", got, tc.want)
		}
	}

	// CondText — singular and plural.
	if got := (InPlay{Player: Controller, Type: Creature}).CondText(); got != "if there is a friendly creature in play" {
		t.Errorf("singular CondText = %q", got)
	}
	if got := (InPlay{Player: Controller, Type: Creature, Amount: 2}).CondText(); got != "if there are 2 or more friendly creatures in play" {
		t.Errorf("plural CondText = %q", got)
	}
	if got := (InPlay{Player: Controller, Type: Creature, Other: true}).CondText(); got != "if there is another friendly creature in play" {
		t.Errorf("other CondText = %q", got)
	}

	// Met — Amount defaults to one; a higher threshold may not be reached.
	if !(InPlay{Player: Controller, Type: Creature}).Met(ctx) {
		t.Error("default threshold should be met with 3 creatures")
	}
	if !(InPlay{Player: Controller, Type: Creature, Amount: 3}).Met(ctx) {
		t.Error("threshold 3 should be met with 3 creatures")
	}
	if (InPlay{Player: Controller, Type: Creature, Amount: 4}).Met(ctx) {
		t.Error("threshold 4 should not be met with 3 creatures")
	}

	// Met with Other counts friendly creatures besides the source.
	src := g.AddToBattleline(NewCard("src", Mars, Creature, Common, WithPower(2)), 0)
	octx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	if !(InPlay{Player: Controller, Type: Creature, Other: true}).Met(octx) {
		t.Error("Other should be met while another friendly creature is in play")
	}
}

// TestCardinalCountText covers the cardinal "the number of …" form a clause
// compares against, including the InPlay side phrasing and the generic fallback.
func TestCardinalCountText(t *testing.T) {
	cardinals := []struct {
		in   Count
		want string
	}{
		{
			InPlay{Player: Controller, Type: Creature, House: Mars},
			"the number of friendly Mars creatures you control",
		},
		{
			InPlay{Player: Opponent, Type: Creature},
			"the number of enemy creatures your opponent controls",
		},
		{
			InPlay{Player: EachPlayer, Type: Creature},
			"the number of creatures in play",
		},
		{OpponentForgedKeys{}, "the number of forged key your opponent has"},
	}
	for _, tc := range cardinals {
		if got := cardinalCountText(tc.in); got != tc.want {
			t.Errorf("cardinalCountText = %q, want %q", got, tc.want)
		}
	}
}

// TestInPlayMinPower covers the power floor Grump Buggy scales its key-cost change
// by: only creatures at or above MinPower count, and the noun names the threshold.
func TestInPlayMinPower(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.AddToBattleline(NewCard("big", Brobnar, Creature, Common, WithPower(6)), 0)
	g.AddToBattleline(NewCard("exact", Brobnar, Creature, Common, WithPower(5)), 0)
	g.AddToBattleline(NewCard("small", Brobnar, Creature, Common, WithPower(4)), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	strong := InPlay{Player: Controller, Type: Creature, MinPower: 5}
	if got := strong.Value(ctx); got != 2 {
		t.Errorf("MinPower 5 Value = %d, want 2", got)
	}
	if got := strong.CountText(); got != "friendly creature with power 5 or higher" {
		t.Errorf("CountText = %q", got)
	}
}

func TestInPlayWithAember(t *testing.T) {
	g := NewGame("A", "B", 1)
	rich := g.AddToBattleline(testCreature("rich", 4), 0)
	g.AddToBattleline(testCreature("poor", 4), 0)
	g.AddAmberOn(rich, 2)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	withAember := InPlay{Player: Controller, Type: Creature, WithAember: true}
	if got := withAember.Value(ctx); got != 1 {
		t.Errorf("WithAember Value = %d, want 1 (only the Æmber-bearer counts)", got)
	}
	if got := withAember.CountText(); got != "friendly creature with \u00c6mber on it" {
		t.Errorf("CountText = %q", got)
	}
}

func TestInPlayEachPlayer(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.AddToBattleline(testCreature("f", 5), 0)
	g.AddToBattleline(testCreature("e", 5), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	// EachPlayer counts both players' matching cards, with no friendly/enemy
	// qualifier in the rendered noun.
	byHouse := InPlay{Player: EachPlayer, Type: Creature, House: Brobnar}
	if got := byHouse.CountText(); got != "Brobnar creature in play" {
		t.Errorf("count text = %q, want %q", got, "Brobnar creature in play")
	}
	if got := byHouse.Value(ctx); got != 2 {
		t.Errorf("value = %d, want 2 (one creature per player)", got)
	}

	if got := (InPlay{Player: EachPlayer, Type: Creature}).CountText(); got != "creature in play" {
		t.Errorf("no-house count text = %q, want %q", got, "creature in play")
	}
}

func TestInPlayReady(t *testing.T) {
	g := NewGame("A", "B", 1)
	spent := g.AddToBattleline(marsCreature("spent", 3), 0)
	g.AddToBattleline(marsCreature("fresh", 3), 0)
	g.SetExhausted(spent, true)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	ready := InPlay{Player: Controller, Type: Creature, House: Mars, Ready: true}
	if got := ready.CountText(); got != "friendly ready Mars creature" {
		t.Errorf("count text = %q, want %q", got, "friendly ready Mars creature")
	}
	if got := ready.Value(ctx); got != 1 {
		t.Errorf("value = %d, want 1 (the exhausted creature does not count)", got)
	}
}

func TestCardsDestroyedCount(t *testing.T) {
	g := NewGame("A", "B", 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	ctx.Produced.Destroyed = [2]int{2, 0}

	c := CardsDestroyed{}
	if got := c.CountText(); got != "card destroyed this way" {
		t.Errorf("count text = %q", got)
	}
	if got := c.Value(ctx); got != 2 {
		t.Errorf("value = %d, want 2", got)
	}
}

func TestCreaturesDestroyedCount(t *testing.T) {
	g := NewGame("A", "B", 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	ctx.Produced.Destroyed = [2]int{2, 1}

	c := CreaturesDestroyed{}
	if got := c.CountText(); got != "creature destroyed this way" {
		t.Errorf("count text = %q", got)
	}
	if got := c.Value(ctx); got != 3 {
		t.Errorf("value = %d, want 3", got)
	}
}

func TestDestroyBindsDestroyedCard(t *testing.T) {
	g := NewGame("A", "B", 1)
	art := g.AddArtifact(
		NewCard(
			"relic",
			Brobnar,
			Artifact,
			Common,
			WithBonus(BonusAember, BonusAember, BonusAember),
		),
		1,
	)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	Destroy{Target: Target{Kind: TargetChosenEnemyArtifact}}.Resolve(ctx)
	// The destroyed artifact is bound in context so a following effect can act on it.
	if !ctx.HasIt || ctx.It != art {
		t.Fatalf("ctx.It = %v (has %v), want the destroyed artifact %d", ctx.It, ctx.HasIt, art)
	}
}

func TestCardsReturnedThisWayCount(t *testing.T) {
	g := NewGame("A", "B", 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	ctx.Produced.Returned = 3

	c := ProducedThisWay{Tally: TallyCardsReturned}
	if got := c.CountText(); got != "card put into your hand this way" {
		t.Errorf("count text = %q", got)
	}
	if got := c.Value(ctx); got != 3 {
		t.Errorf("value = %d, want 3", got)
	}
}

func TestCardsPurgedThisWayCount(t *testing.T) {
	g := NewGame("A", "B", 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	ctx.Produced.Purged = [2]int{2, 1}

	mine := ProducedThisWay{Tally: TallyCardsPurged, Player: Controller}
	if got := mine.CountText(); got != "card they controlled that was purged this way" {
		t.Errorf("count text = %q", got)
	}
	if got := mine.Value(ctx); got != 2 {
		t.Errorf("value = %d, want 2", got)
	}

	theirs := ProducedThisWay{Tally: TallyCardsPurged, Player: Opponent}
	if got := theirs.CountText(); got != "card your opponent controlled that was purged this way" {
		t.Errorf("opponent count text = %q", got)
	}
	if got := theirs.Value(ctx); got != 1 {
		t.Errorf("opponent value = %d, want 1", got)
	}
}

func TestHousesInPlay(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.AddToBattleline(NewCard("m", Mars, Creature, Common, WithPower(4)), 0)
	g.AddToBattleline(NewCard("l", Logos, Creature, Common, WithPower(4)), 0)
	g.AddToBattleline(NewCard("s", Sanctum, Creature, Common, WithPower(4)), 0)
	g.AddToBattleline(NewCard("b", Brobnar, Creature, Common, WithPower(4)), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	all := HousesInPlay{}
	if got := all.Value(ctx); got != 4 {
		t.Errorf("value = %d, want 4", got)
	}
	if got := all.CountText(); got != "house represented among cards in play" {
		t.Errorf("count text = %q", got)
	}

	exceptSanctum := HousesInPlay{Except: Sanctum}
	if got := exceptSanctum.Value(ctx); got != 3 {
		t.Errorf("except value = %d, want 3", got)
	}
	want := "house represented among cards in play, except for Sanctum"
	if got := exceptSanctum.CountText(); got != want {
		t.Errorf("count text = %q, want %q", got, want)
	}
}

func TestHousesAmong(t *testing.T) {
	g := NewGame("A", "B", 1)
	// Controller's board: three creature houses, one artifact house, one houseless
	// creature (which counts toward no house).
	g.AddToBattleline(NewCard("m0", Mars, Creature, Common, WithPower(4)), 0)
	g.AddToBattleline(NewCard("l0", Logos, Creature, Common, WithPower(4)), 0)
	g.AddToBattleline(NewCard("s0", Sanctum, Creature, Common, WithPower(4)), 0)
	g.AddToBattleline(NewCard("none", HouseNone, Creature, Common, WithPower(4)), 0)
	g.AddArtifact(NewCard("u0", Untamed, Artifact, Common), 0)
	// Opponent's board: two creature houses (one shared with the controller) and one
	// artifact house.
	g.AddToBattleline(NewCard("b1", Brobnar, Creature, Common, WithPower(4)), 1)
	g.AddToBattleline(NewCard("m1", Mars, Creature, Common, WithPower(4)), 1)
	g.AddArtifact(NewCard("sh1", Shadows, Artifact, Common), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	cases := []struct {
		name  string
		count HousesAmong
		value int
		text  string
	}{
		{
			"friendly creatures, uncapped",
			HousesAmong{Player: Controller, Type: Creature},
			3, "house represented among friendly creatures",
		},
		{
			"enemy creatures",
			HousesAmong{Player: Opponent, Type: Creature},
			2, "house represented among enemy creatures",
		},
		{
			"friendly cards, both rows",
			HousesAmong{Player: Controller},
			4, "house represented among friendly cards",
		},
		{
			"every card in play, both rows",
			HousesAmong{Player: EachPlayer},
			6, "house represented among cards in play",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.count.Value(ctx); got != c.value {
				t.Errorf("value = %d, want %d", got, c.value)
			}
			if got := c.count.CountText(); got != c.text {
				t.Errorf("count text = %q, want %q", got, c.text)
			}
		})
	}
}

func TestEachFriendlyArtifactTarget(t *testing.T) {
	g := NewGame("A", "B", 1)
	mine := g.AddArtifact(NewCard("mine", Brobnar, Artifact, Common), 0)
	g.AddArtifact(NewCard("theirs", Brobnar, Artifact, Common), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	target := Target{Kind: TargetEachFriendlyArtifact}
	if got := target.Text(); got != "each friendly artifact" {
		t.Errorf("text = %q", got)
	}
	got := target.Select(ctx)
	if len(got) != 1 || got[0] != mine {
		t.Errorf("selected = %v, want [%v]", got, mine)
	}
}

func TestExcessCreatures(t *testing.T) {
	opponents := ExcessCreatures{Player: Opponent}
	mine := ExcessCreatures{Player: Controller}
	if opponents.CountText() != "creature your opponent controls in excess of you" {
		t.Errorf("count text = %q", opponents.CountText())
	}
	if mine.CountText() != "creature you have in excess of your opponent" {
		t.Errorf("count text = %q", mine.CountText())
	}

	g := NewGame("A", "B", 1)
	g.AddToBattleline(testCreature("o1", 3), 1)
	g.AddToBattleline(testCreature("o2", 3), 1)
	g.AddToBattleline(testCreature("m1", 3), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	if got := opponents.Value(ctx); got != 1 {
		t.Errorf("excess = %d, want 1 (opponent 2, you 1)", got)
	}
	if got := mine.Value(ctx); got != 0 {
		t.Errorf("own excess = %d, want 0 when your opponent controls more", got)
	}

	// When the controller controls at least as many, the opponent's excess floors
	// at zero and the controller's own excess starts counting up.
	g.AddToBattleline(testCreature("m2", 3), 0)
	g.AddToBattleline(testCreature("m3", 3), 0)
	if got := opponents.Value(ctx); got != 0 {
		t.Errorf("excess = %d, want 0 when you control more", got)
	}
	if got := mine.Value(ctx); got != 1 {
		t.Errorf("own excess = %d, want 1 (you 3, opponent 2)", got)
	}
}

func TestExcessCreaturesNotCountingSelf(t *testing.T) {
	// Dr. Milli: "in excess of you, not counting Dr. Milli". The source creature
	// is excluded from its controller's side of the comparison.
	opp := ExcessCreatures{Player: Opponent, NotCountingSelf: true}
	if want := "creature your opponent controls in excess of you, not counting " + SelfName; opp.CountText() != want {
		t.Errorf("count text = %q, want %q", opp.CountText(), want)
	}
	mine := ExcessCreatures{Player: Controller, NotCountingSelf: true}
	if want := "creature you have in excess of your opponent, not counting " + SelfName; mine.CountText() != want {
		t.Errorf("own count text = %q, want %q", mine.CountText(), want)
	}

	g := NewGame("A", "B", 1)
	g.AddToBattleline(testCreature("o1", 3), 1)
	g.AddToBattleline(testCreature("o2", 3), 1)
	g.AddToBattleline(testCreature("m1", 3), 0)
	g.AddToBattleline(testCreature("m2", 3), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	// Opponent side is "more"; the controller (the "less" side) drops one for self.
	if got := opp.Value(ctx); got != 1 {
		t.Errorf("opponent excess = %d, want 1 (opp 2 vs you 2-1)", got)
	}
	// Controller side is "more"; it drops one for self.
	if got := mine.Value(ctx); got != 0 {
		t.Errorf("own excess = %d, want 0 (you 2-1 vs opp 2)", got)
	}
}

func TestNeighborsSharingHouse(t *testing.T) {
	c := NeighborsSharingHouse{}
	if got := c.CountText(); got != "neighbor that shares a house with it" {
		t.Errorf("count text = %q", got)
	}
	if got := c.leadingCountText(); got != "neighbor that shares a house with the chosen creature" {
		t.Errorf("leading count text = %q", got)
	}

	// With no creature in context the count is zero.
	g := NewGame("A", "B", 1)
	if got := c.Value(&EffectContext{Resolver: g, Controller: 0}); got != 0 {
		t.Errorf("value with no context creature = %d, want 0", got)
	}

	left := g.AddToBattleline(NewCard("l", Mars, Creature, Common, WithPower(3)), 0)
	mid := g.AddToBattleline(NewCard("m", Mars, Creature, Common, WithPower(3)), 0)
	right := g.AddToBattleline(NewCard("r", Logos, Creature, Common, WithPower(3)), 0)
	_, _ = left, right
	ctx := &EffectContext{Resolver: g, Controller: 0, It: mid, HasIt: true}
	if got := c.Value(ctx); got != 1 {
		t.Errorf("shared-house neighbors = %d, want 1 (Mars left, Logos right)", got)
	}
}

func TestInPlayByName(t *testing.T) {
	g := NewGame("A", "B", 1)
	bear := g.Register(NewCard("Ancient Bear", Untamed, Creature, Common, WithPower(6)), 0)
	g.AddToBattleline(testCreature("other", 3), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	none := InPlay{Player: EachPlayer, Type: Creature, Name: "Ancient Bear", None: true}
	if want := "if there are no Ancient Bears in play"; none.CondText() != want {
		t.Errorf("cond text = %q, want %q", none.CondText(), want)
	}
	if !none.Met(ctx) {
		t.Error("None should be met while no Ancient Bear is in play")
	}

	g.State.Battleline[0].add(bear)
	if none.Met(ctx) {
		t.Error("None should not be met once an Ancient Bear is in play")
	}
	some := InPlay{Player: EachPlayer, Type: Creature, Name: "Ancient Bear"}
	if n := some.Value(ctx); n != 1 {
		t.Errorf("value = %d, want 1 (only the named card counts)", n)
	}
	if want := "if there is an Ancient Bear in play"; some.CondText() != want {
		t.Errorf("cond text = %q, want %q", some.CondText(), want)
	}
	two := InPlay{Player: EachPlayer, Type: Creature, Name: "Ancient Bear", Amount: 2}
	if want := "if there are 2 or more Ancient Bears in play"; two.CondText() != want {
		t.Errorf("cond text = %q, want %q", two.CondText(), want)
	}
}

func TestAemberInPool(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.SetAember(0, 4)
	g.SetAember(1, 2)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	if got := (AemberInPool{Player: Controller}).Value(ctx); got != 4 {
		t.Errorf("controller pool count = %d, want 4", got)
	}
	if got := (AemberInPool{Player: Opponent}).Value(ctx); got != 2 {
		t.Errorf("opponent pool count = %d, want 2", got)
	}
	if got := (AemberInPool{Player: Controller}).CountText(); got != "Æmber in your pool" {
		t.Errorf("count text = %q", got)
	}
	if got := (AemberInPool{Player: Opponent}).CountText(); got != "Æmber in your opponent's pool" {
		t.Errorf("count text = %q", got)
	}
}

func TestAemberOnFriendlyCreatures(t *testing.T) {
	g := NewGame("A", "B", 1)
	a := g.AddToBattleline(testCreature("a", 3), 0)
	b := g.AddToBattleline(testCreature("b", 3), 0)
	g.AddToBattleline(testCreature("e", 3), 1)
	g.addAmberOn(a, 2)
	g.addAmberOn(b, 3)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	c := AemberOnFriendlyCreatures{}
	if got := c.Value(ctx); got != 5 {
		t.Errorf("value = %d, want 5", got)
	}
	if got := c.CountText(); got != "Æmber on friendly creatures" {
		t.Errorf("count text = %q", got)
	}
}

func TestNeighborsOfThis(t *testing.T) {
	g := NewGame("A", "B", 1)
	left := NewCard("left", Mars, Creature, Common, WithPower(2))
	mid := NewCard("mid", Mars, Creature, Common, WithPower(2))
	right := NewCard("right", Mars, Creature, Common, WithPower(2))
	g.AddToBattleline(left, 0)
	g.AddToBattleline(mid, 0)
	g.AddToBattleline(right, 0)

	// A middle creature has two neighbors; a flank creature has one.
	midID := g.State.Battleline[0].slice()[1]
	flankID := g.State.Battleline[0].slice()[0]

	if got := (NeighborsOfThis{}).Value(&EffectContext{Resolver: g, Source: midID}); got != 2 {
		t.Errorf("middle neighbors = %d, want 2", got)
	}
	if got := (NeighborsOfThis{}).Value(&EffectContext{Resolver: g, Source: flankID}); got != 1 {
		t.Errorf("flank neighbors = %d, want 1", got)
	}
	if got := (NeighborsOfThis{}).CountText(); got != "neighbor it has" {
		t.Errorf("count text = %q", got)
	}
	// Leading a sentence, the clause names the source instead of saying "it".
	if got := (NeighborsOfThis{}).leadingCountText(); got != "neighbor "+SelfName+" has" {
		t.Errorf("leading count text = %q", got)
	}
	// keyCostText renders the leading, named form for a self-referential count.
	kc := NewKeyCostChange(Opponent, 2).Per(NeighborsOfThis{})
	want := "For each neighbor " + SelfName + " has, your opponent's keys cost +2 Æmber."
	if got := keyCostText(kc); got != want {
		t.Errorf("key cost text = %q, want %q", got, want)
	}
}

func TestFixed(t *testing.T) {
	// Fixed yields its constant regardless of game state and leads no clause.
	if got := Fixed(3).Value(nil); got != 3 {
		t.Errorf("Fixed value = %d, want 3", got)
	}
	if got := Fixed(3).CountText(); got != "" {
		t.Errorf("Fixed count text = %q, want empty", got)
	}
}

func TestTraitsOfChosen(t *testing.T) {
	g := NewGame("A", "B", 1)
	id := g.AddToBattleline(
		NewCard("beast", Untamed, Creature, Common, WithPower(3), WithTraits(Beast, Mutant)), 0)
	c := TraitsOfChosen{}
	if got := c.CountText(); got != "trait that creature has" {
		t.Errorf("text = %q", got)
	}
	if got := c.Value(&EffectContext{Resolver: g, Controller: 0, It: id, HasIt: true}); got != 2 {
		t.Errorf("value = %d, want 2", got)
	}
	if got := c.Value(&EffectContext{Resolver: g, Controller: 0}); got != 0 {
		t.Errorf("value with no creature in context = %d, want 0", got)
	}
}

// TestPowerOfChosen covers the full-power count (Mindworm): the value is the
// context creature's power, zero without one, and it renders inside an "equal to"
// clause.
func TestPowerOfChosen(t *testing.T) {
	if got := (PowerOfChosen{}).CountText(); got != "its power" {
		t.Errorf("CountText = %q, want %q", got, "its power")
	}
	if got := (PowerOfChosen{}).Value(&EffectContext{}); got != 0 {
		t.Errorf("Value with no It = %d, want 0", got)
	}

	g := started(t)
	c := g.AddToBattleline(NewCard("brute", Mars, Creature, Common, WithPower(6)), 0)
	if got := (PowerOfChosen{}).Value(&EffectContext{Resolver: g, It: c, HasIt: true}); got != 6 {
		t.Errorf("Value = %d, want 6", got)
	}

	e := DealDamage{
		AmountFrom: PowerOfChosen{},
		Target:     Target{Kind: TargetCreatureFought}.NeighborsOf(),
	}
	want := "deal damage equal to its power to each neighbor of the creature {self} fights"
	if got := e.Text(); got != want {
		t.Errorf("Text = %q, want %q", got, want)
	}
}
