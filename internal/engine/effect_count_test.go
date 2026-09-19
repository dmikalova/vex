package engine

import (
	"strings"
	"testing"
)

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
		in   CardsInPlay
		want int
	}{
		{"friendly creatures", CardsInPlay{Player: Controller, Type: Creature}, 3},
		{
			"friendly Mars creatures",
			CardsInPlay{Player: Controller, Type: Creature, House: namedHouse(Mars)},
			2,
		},
		{"friendly artifacts", CardsInPlay{Player: Controller, Type: Artifact}, 2},
		{"friendly cards, any type", CardsInPlay{Player: Controller}, 5},
		{"friendly Shards", CardsInPlay{Player: Controller, Trait: Shard}, 1},
		{"enemy creatures", CardsInPlay{Player: Opponent, Type: Creature}, 1},
	}
	for _, tc := range values {
		if got := tc.in.Value(ctx); got != tc.want {
			t.Errorf("%s: Value = %d, want %d", tc.name, got, tc.want)
		}
	}

	// CountText — the singular "for each" noun.
	texts := []struct {
		in   CardsInPlay
		want string
	}{
		{CardsInPlay{Player: Controller, Type: Creature}, "friendly creature in play"},
		{
			CardsInPlay{Player: Controller, Type: Creature, House: namedHouse(Mars)},
			"friendly Mars creature",
		},
		{CardsInPlay{Player: Controller, Trait: Shard}, "friendly Shard"},
		{CardsInPlay{Player: Controller, Type: Creature, Trait: Thief}, "friendly Thief creature"},
		{CardsInPlay{Player: Controller, Type: Artifact, Trait: Shard}, "friendly Shard artifact"},
		{CardsInPlay{Player: Opponent, Type: Creature}, "enemy creature in play"},
		{CardsInPlay{Player: Controller, Type: Artifact}, "friendly artifact in play"},
		{CardsInPlay{Player: Controller}, "friendly card in play"},
	}
	for _, tc := range texts {
		if got := tc.in.CountText(); got != tc.want {
			t.Errorf("CountText = %q, want %q", got, tc.want)
		}
	}

	// CondText — singular and plural.
	if got := (CardsInPlay{Player: Controller, Type: Creature}).CondText(); got != "if there is a friendly creature in play" {
		t.Errorf("singular CondText = %q", got)
	}
	if got := (CardsInPlay{Player: Controller, Type: Creature, Amount: 2}).CondText(); got != "if there are 2 or more friendly creatures in play" {
		t.Errorf("plural CondText = %q", got)
	}
	if got := (CardsInPlay{Player: Controller, Type: Creature, Other: true}).CondText(); got != "if there is another friendly creature in play" {
		t.Errorf("other CondText = %q", got)
	}

	// Met — Amount defaults to one; a higher threshold may not be reached.
	if !(CardsInPlay{Player: Controller, Type: Creature}).Met(ctx) {
		t.Error("default threshold should be met with 3 creatures")
	}
	if !(CardsInPlay{Player: Controller, Type: Creature, Amount: 3}).Met(ctx) {
		t.Error("threshold 3 should be met with 3 creatures")
	}
	if (CardsInPlay{Player: Controller, Type: Creature, Amount: 4}).Met(ctx) {
		t.Error("threshold 4 should not be met with 3 creatures")
	}

	// Met with Other counts friendly creatures besides the source.
	src := g.AddToBattleline(NewCard("src", Mars, Creature, Common, WithPower(2)), 0)
	octx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	if !(CardsInPlay{Player: Controller, Type: Creature, Other: true}).Met(octx) {
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
			CardsInPlay{Player: Controller, Type: Creature, House: namedHouse(Mars)},
			"the number of friendly Mars creatures you control",
		},
		{
			CardsInPlay{Player: Opponent, Type: Creature},
			"the number of enemy creatures your opponent controls",
		},
		{
			CardsInPlay{Player: EachPlayer, Type: Creature},
			"the number of creatures in play",
		},
		{ForgedKeys{Player: Opponent}, "the number of forged key your opponent has"},
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

	strong := CardsInPlay{Player: Controller, Type: Creature, MinPower: 5}
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

	withAember := CardsInPlay{Player: Controller, Type: Creature, WithAember: true}
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
	byHouse := CardsInPlay{Player: EachPlayer, Type: Creature, House: namedHouse(Brobnar)}
	if got := byHouse.CountText(); got != "Brobnar creature in play" {
		t.Errorf("count text = %q, want %q", got, "Brobnar creature in play")
	}
	if got := byHouse.Value(ctx); got != 2 {
		t.Errorf("value = %d, want 2 (one creature per player)", got)
	}

	if got := (CardsInPlay{Player: EachPlayer, Type: Creature}).CountText(); got != "creature in play" {
		t.Errorf("no-house count text = %q, want %q", got, "creature in play")
	}
}

func TestInPlayReady(t *testing.T) {
	g := NewGame("A", "B", 1)
	spent := g.AddToBattleline(marsCreature("spent", 3), 0)
	g.AddToBattleline(marsCreature("fresh", 3), 0)
	g.SetExhausted(spent, true)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	ready := CardsInPlay{Player: Controller, Type: Creature, House: namedHouse(Mars), Ready: true}
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

func TestNeighborsMatching(t *testing.T) {
	c := NeighborsMatching{House: HouseMatcher{Kind: MatchContextualHouse}}
	if got := c.CountText(); got != "neighbor of that card's house" {
		t.Errorf("count text = %q", got)
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

	none := CardsInPlay{Player: EachPlayer, Type: Creature, Name: "Ancient Bear", None: true}
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
	some := CardsInPlay{Player: EachPlayer, Type: Creature, Name: "Ancient Bear"}
	if n := some.Value(ctx); n != 1 {
		t.Errorf("value = %d, want 1 (only the named card counts)", n)
	}
	if want := "if there is an Ancient Bear in play"; some.CondText() != want {
		t.Errorf("cond text = %q, want %q", some.CondText(), want)
	}
	two := CardsInPlay{Player: EachPlayer, Type: Creature, Name: "Ancient Bear", Amount: 2}
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

// TestEqualToTextSpeaksToWhoIsPaid pins the voice an "equal to" count is read in:
// Binate Rupture pays each player, so the count is about them, not about you.
func TestEqualToTextSpeaksToWhoIsPaid(t *testing.T) {
	pool := AemberInPool{Player: Controller}
	if got := equalToText(pool, Controller); got != "Æmber in your pool" {
		t.Errorf("controller equal-to text = %q", got)
	}
	if got := equalToText(AemberInPool{Player: Opponent}, EachPlayer); got !=
		"the Æmber in their opponent's pool" {
		t.Errorf("mirrored equal-to text = %q", got)
	}
	if got := equalToText(pool, EachPlayer); got != "the Æmber in their pool" {
		t.Errorf("each-player equal-to text = %q", got)
	}
	// A count with no third-person form keeps its own voice rather than blocking
	// an EachPlayer gain.
	if got := equalToText(Fixed(2), EachPlayer); got != (Fixed(2)).CountText() {
		t.Errorf("plain count under EachPlayer = %q, want its own text", got)
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

// TestBonusIconsOfChosen covers Mindfire: the value is the number of bonus icons
// on the card in context (ctx.It), zero without one, and the text names the card.
func TestBonusIconsOfChosen(t *testing.T) {
	c := BonusIconsOfChosen{Noun: DiscardedCard}
	if got := c.CountText(); got != "bonus icon on the discarded card" {
		t.Errorf("CountText = %q", got)
	}
	if got := c.Value(&EffectContext{}); got != 0 {
		t.Errorf("Value with no It = %d, want 0", got)
	}

	g := NewGame("A", "B", 1)
	id := g.Register(
		NewCard("pip", Dis, Tactic, Common, WithBonus(BonusAember, BonusDraw)), 1)
	if got := g.BonusIconCount(id); got != 2 {
		t.Errorf("BonusIconCount = %d, want 2", got)
	}
	if got := c.Value(&EffectContext{Resolver: g, It: id, HasIt: true}); got != 2 {
		t.Errorf("Value = %d, want 2", got)
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

// TestCombinedPowerOfNeighborsWithout covers the count's live sum and its two
// text renderings, and the blank-able "X" power it drives on a card.
func TestCombinedPowerOfNeighborsWithout(t *testing.T) {
	c := CombinedPowerOfNeighborsWithout{Without: Changeling}
	if got := c.CountText(); got != "combined power of {self}'s non-Changeling neighbors" {
		t.Errorf("CountText = %q", got)
	}
	if got := c.cardinalCountText(); got !=
		"the combined power of {self}'s non-Changeling neighbors" {
		t.Errorf("cardinalCountText = %q", got)
	}

	g := NewGame("A", "B", 1)
	left := g.AddToBattleline(testCreature("left", 3), 0)
	picaroon := g.AddToBattleline(
		testCreature("Picaroon", 0, WithTraits(Mutant, Changeling), WithPowerX(c)), 0)
	// A Changeling neighbor is excluded from the tally.
	g.AddToBattleline(testCreature("changeling", 4, WithTraits(Changeling)), 0)
	_ = left

	// 3 (non-Changeling left) + 0 (excluded Changeling right) = 3.
	if got := g.Power(picaroon); got != 3 {
		t.Errorf("Picaroon power = %d, want 3", got)
	}

	// Blanking the card drops its X power to 0: the "X is …" line is part of its text.
	BlankEnemyText{}.Resolve(&EffectContext{Resolver: g, Source: picaroon, Controller: 1})
	if !g.textBlanked(picaroon) {
		t.Fatal("precondition: Picaroon should be blanked")
	}
	if got := g.Power(picaroon); got != 0 {
		t.Errorf("blanked Picaroon power = %d, want 0", got)
	}

	// The passive renders into the card's rules text.
	def := NewCard("Picaroon", Dis, Creature, Uncommon, WithPower(0), WithPowerX(c))
	if got := RenderCardRules(&def); !strings.Contains(got,
		"X is the combined power of Picaroon's non-Changeling neighbors.") {
		t.Errorf("rules = %q", got)
	}
}

// TestCombinedPowerOfNeighborsTerminatesOnCycle guards the mutual-reference
// hazard: two creatures whose X power reads each other (two Picaroons that both
// lost Changeling to Grey Aberrant) would recurse forever without the guard. The
// re-entrant creature contributes 0, so power stays finite and Power never
// overflows the stack.
func TestCombinedPowerOfNeighborsTerminatesOnCycle(t *testing.T) {
	c := CombinedPowerOfNeighborsWithout{Without: Changeling}
	g := NewGame("A", "B", 1)
	// Grey Aberrant has stripped Changeling, so these two X-power creatures now each
	// count the other. Line: anchor(3) — p1 — p2.
	anchor := g.AddToBattleline(testCreature("anchor", 3), 0)
	p1 := g.AddToBattleline(testCreature("p1", 0, WithPowerX(c)), 0)
	p2 := g.AddToBattleline(testCreature("p2", 0, WithPowerX(c)), 0)
	_ = anchor

	// p2's only non-Changeling neighbor is p1, which is mid-computation when p2 is
	// reached, so p1 contributes 0 there: p2 = 0. p1 = anchor(3) + p2(0) = 3.
	if got := g.Power(p2); got != 3 {
		t.Errorf("p2 power = %d, want 3", got)
	}
	if got := g.Power(p1); got != 3 {
		t.Errorf("p1 power = %d, want 3", got)
	}
	// The guard stack is balanced back to empty after each top-level read.
	if len(g.powerComputing) != 0 {
		t.Errorf("powerComputing stack left at %d entries, want 0", len(g.powerComputing))
	}
}

func TestArtifactsInPlayCount(t *testing.T) {
	g := started(t)
	g.AddArtifact(NewCard("a1", Brobnar, Artifact, Common), 0)
	g.AddArtifact(NewCard("a2", Brobnar, Artifact, Common), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	if got := (ArtifactsInPlay{}).Value(ctx); got != 2 {
		t.Errorf("ArtifactsInPlay = %d, want 2", got)
	}
	if got := (ArtifactsInPlay{}).CountText(); got != "artifact in play" {
		t.Errorf("CountText = %q", got)
	}
}

func TestPowerCountersOnThisCount(t *testing.T) {
	g := started(t)
	c := g.AddToBattleline(testCreature("c", 3), 0)
	g.AddPowerCounter(c, 4)
	ctx := &EffectContext{Resolver: g, Source: c, Controller: 0}

	if got := (PowerCountersOnThis{}).Value(ctx); got != 4 {
		t.Errorf("PowerCountersOnThis = %d, want 4", got)
	}
	want := "the number of +1 power counters on " + SelfName
	if got := (PowerCountersOnThis{}).CountText(); got != want {
		t.Errorf("CountText = %q, want %q", got, want)
	}
}

// TestCannotBeDealtDamageOpponentNextTurn covers the cross-turn immunity (Lucky
// Dice): dormant on the caster's turn, active during the opponent's next turn.
func TestCannotBeDealtDamageOpponentNextTurn(t *testing.T) {
	g := started(t)
	c := g.AddToBattleline(testCreature("c", 5), 0)
	e := CannotBeDealtDamage{
		Target:   Target{Kind: TargetEachFriendlyCreature},
		Duration: OpponentNextTurn,
	}
	e.Resolve(&EffectContext{Resolver: g, Controller: 0})

	if g.DamageImmune(c) {
		t.Error("OpponentNextTurn immunity should be dormant on the caster's turn")
	}
	g.EndPlayPhase(0)
	g.StartTurn(1)
	if !g.DamageImmune(c) {
		t.Error("immunity should be active during the opponent's next turn")
	}
	want := "during your opponent's next turn, each friendly creature cannot be dealt damage"
	if got := e.Text(); got != want {
		t.Errorf("text = %q, want %q", got, want)
	}
}

// TestTheirCreaturesThisWayCounts covers the two "... this way" counts that
// read the per-player tallies, and the tally itself as Destroy and
// PutFromPlay record it.
func TestCreaturesRemovedThisWayCounts(t *testing.T) {
	texts := []struct {
		count Count
		want  string
	}{
		{
			ProducedThisWay{Tally: TallyCreaturesDestroyed, Player: Controller},
			"creature they controlled that was destroyed this way",
		},
		{
			ProducedThisWay{Tally: TallyCreaturesShuffledIntoDeck, Player: Controller},
			"creature shuffled into their deck this way",
		},
		{
			ProducedThisWay{Tally: TallyCreaturesDestroyed, Player: Opponent},
			"creature your opponent controlled that was destroyed this way",
		},
		{
			ProducedThisWay{Tally: TallyCreaturesShuffledIntoDeck, Player: Opponent},
			"creature shuffled into your opponent's deck this way",
		},
	}
	for _, tc := range texts {
		if got := tc.count.CountText(); got != tc.want {
			t.Errorf("CountText = %q, want %q", got, tc.want)
		}
	}

	ctx := &EffectContext{Controller: 1}
	ctx.Produced.Destroyed = [2]int{4, 7}
	ctx.Produced.Moved = [2]int{5, 9}
	if got := (ProducedThisWay{Tally: TallyCreaturesDestroyed, Player: Controller}).Value(
		ctx,
	); got != 7 {
		t.Errorf("destroyed Value = %d, want 7", got)
	}
	if got := (ProducedThisWay{Tally: TallyCreaturesShuffledIntoDeck, Player: Controller}).Value(
		ctx,
	); got != 9 {
		t.Errorf("shuffled Value = %d, want 9", got)
	}
	if got := ctx.Produced.TotalDestroyed(); got != 11 {
		t.Errorf("TotalDestroyed = %d, want 11", got)
	}
}

// TestPowerDestroyedThisWayCount covers the summed-power tally Destroy fills and
// the PowerDestroyedThisWay count CountIs reads (Might Makes Right forges above
// 25). CountText names the noun a "for each" clause would repeat; CountClause is
// the clause CountIs puts after "if".
func TestPowerDestroyedThisWayCount(t *testing.T) {
	if got := (PowerDestroyedThisWay{}).CountText(); got != "power of creatures destroyed this way" {
		t.Errorf("CountText = %q", got)
	}
	if got := (PowerDestroyedThisWay{}).CountClause(
		"25 or more",
		false,
	); got != "the total power of creatures destroyed this way is 25 or more" {
		t.Errorf("CountClause = %q", got)
	}
	ctx := &EffectContext{}
	ctx.Produced.DestroyedPower = 31
	if got := (PowerDestroyedThisWay{}).Value(ctx); got != 31 {
		t.Errorf("Value = %d, want 31", got)
	}
}

// TestAemberLostThisWayCount covers the Æmber-lost tally LoseAember fills and a
// ProducedThisWay with TallyAemberLost reads (Shatter Storm).
func TestAemberLostThisWayCount(t *testing.T) {
	if got := (ProducedThisWay{Tally: TallyAemberLost, Player: Controller}).CountText(); got != "Æmber you lost this way" {
		t.Errorf("CountText = %q", got)
	}
	want := "Æmber your opponent lost this way"
	if got := (ProducedThisWay{Tally: TallyAemberLost, Player: Opponent}).CountText(); got != want {
		t.Errorf("opponent CountText = %q", got)
	}

	g := NewGame("A", "B", 1)
	g.SetAember(0, 3)
	g.SetAember(1, 10)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	// Losing everything fills the tally, which the next loss triples.
	LoseAember{Player: Controller, By: AllAember}.Resolve(ctx)
	if got := (ProducedThisWay{Tally: TallyAemberLost, Player: Controller}).Value(ctx); got != 3 {
		t.Errorf("tally = %d, want 3", got)
	}
	LoseAember{
		Player: Opponent,
		Amount: 3,
		Per:    ProducedThisWay{Tally: TallyAemberLost, Player: Controller},
	}.Resolve(ctx)
	if got := g.Aember(1); got != 1 {
		t.Errorf("opponent pool = %d, want 1", got)
	}
}

// TestAllAemberLoss covers the Loss that empties a pool whatever its size.
func TestAllAemberLoss(t *testing.T) {
	e := LoseAember{Player: Controller, By: AllAember}
	if got := e.Text(); got != "lose all your Æmber" {
		t.Errorf("text = %q", got)
	}
	if got := (LoseAember{Player: Opponent, By: AllAember}).Text(); got != "your opponent loses all their Æmber" {
		t.Errorf("opponent text = %q", got)
	}
	if got := AllAember.lose(7); got != 7 {
		t.Errorf("lose(7) = %d, want 7", got)
	}
	if got := AllAember.qualifier(); got != "" {
		t.Errorf("qualifier = %q, want empty (every pool is affected)", got)
	}
}

// TestDestroyTalliesPerController checks Destroy splits its "this way"
// tally by the controller of each creature that actually left play.
func TestDestroyTalliesRemovalsPerController(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.AddToBattleline(NewCard("mine1", Dis, Creature, Common, WithPower(2)), 0)
	g.AddToBattleline(NewCard("mine2", Dis, Creature, Common, WithPower(2)), 0)
	g.AddToBattleline(NewCard("theirs", Dis, Creature, Common, WithPower(2)), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	Destroy{Target: Target{Kind: TargetEachCreature}.House(namedHouse(Dis))}.Resolve(ctx)

	if ctx.Produced.Destroyed != [2]int{2, 1} {
		t.Errorf("Destroyed = %v, want [2 1]", ctx.Produced.Destroyed)
	}
}

// TestPutFromPlayTalliesPerController checks PutFromPlay records the
// same per-controller tally when it shuffles creatures away.
func TestPutFromPlayTalliesRemovalsPerController(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.AddToBattleline(NewCard("mine", Mars, Creature, Common, WithPower(2)), 0)
	g.AddToBattleline(NewCard("theirs1", Mars, Creature, Common, WithPower(2)), 1)
	g.AddToBattleline(NewCard("theirs2", Mars, Creature, Common, WithPower(2)), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	PutFromPlay{
		Target:      Target{Kind: TargetEachCreature}.House(namedHouse(Mars)),
		Destination: ToDeckShuffled,
	}.Resolve(ctx)

	if ctx.Produced.Moved != [2]int{1, 2} {
		t.Errorf("Moved = %v, want [1 2]", ctx.Produced.Moved)
	}
}

// TestPutFromPlayMovesTheWholeSelectionTogether checks every card an effect
// targets moves before any of their "Leaves Play:" abilities fire. The bomb's
// ability would destroy both enemy creatures, but they are in the same selection
// and have already been shuffled away by the time it resolves, so it finds
// nothing and all three cards are counted. Before leave-play windows were held
// to the end of the batch the bomb fired mid-loop and the other two were skipped,
// which made the outcome depend on the order the selection happened to be
// visited.
func TestPutFromPlayMovesTheWholeSelectionTogether(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.AddToBattleline(NewCard("bomb", Mars, Creature, Common, WithPower(2),
		WithAbility(TriggerLeavesPlay, Destroy{
			Target: Target{Kind: TargetEachEnemyCreature},
		})), 0)
	theirs1 := g.AddToBattleline(NewCard("theirs1", Mars, Creature, Common, WithPower(2)), 1)
	theirs2 := g.AddToBattleline(NewCard("theirs2", Mars, Creature, Common, WithPower(2)), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	PutFromPlay{
		Target:      Target{Kind: TargetEachCreature}.House(namedHouse(Mars)),
		Destination: ToDeckShuffled,
	}.Resolve(ctx)

	if ctx.Produced.Moved != [2]int{1, 2} {
		t.Errorf("Moved = %v, want [1 2]: the whole selection moves as one moment",
			ctx.Produced.Moved)
	}
	for _, id := range []LocalID{theirs1, theirs2} {
		if !g.State.Deck[1].contains(id) {
			t.Errorf("card %d was destroyed mid-selection instead of being shuffled away", id)
		}
	}
}

// TestGainAemberEachPlayer checks an EachPlayer gain pays both players, and
// scales each player's share by the count read from their own side.
func TestGainAemberEachPlayer(t *testing.T) {
	gain := GainAember{
		Player: EachPlayer,
		Amount: 1,
		Per:    ProducedThisWay{Tally: TallyCreaturesDestroyed, Player: Controller},
	}
	if got, want := gain.Text(),
		"for each creature they controlled that was destroyed this way, "+
			"each player gains 1 Æmber"; got != want {
		t.Errorf("Text = %q, want %q", got, want)
	}

	g := NewGame("A", "B", 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	ctx.Produced.Destroyed = [2]int{3, 2}
	gain.Resolve(ctx)

	if got := g.State.Aember[0]; got != 3 {
		t.Errorf("controller Æmber = %d, want 3", got)
	}
	if got := g.State.Aember[1]; got != 2 {
		t.Errorf("opponent Æmber = %d, want 2", got)
	}
}

// TestController checks the port exposes which player a card in play answers to.
func TestController(t *testing.T) {
	g := NewGame("A", "B", 1)
	id := g.AddToBattleline(NewCard("theirs", Dis, Creature, Common, WithPower(2)), 1)
	if got := g.Controller(id); got != 1 {
		t.Errorf("Controller = %d, want 1", got)
	}
}

// TestForgedKeysCountText covers both sides of the "for each forged key" noun: the
// controller's own keys and the opponent's.
func TestForgedKeysCountText(t *testing.T) {
	if got := (ForgedKeys{Player: Controller}).CountText(); got != "forged key you have" {
		t.Errorf("CountText = %q, want %q", got, "forged key you have")
	}
	if got := (ForgedKeys{Player: Opponent}).CountText(); got != "forged key your opponent has" {
		t.Errorf("CountText = %q, want %q", got, "forged key your opponent has")
	}
}
