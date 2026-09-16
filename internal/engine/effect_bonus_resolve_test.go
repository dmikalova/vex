package engine

import "testing"

func TestResolveBonusIconsText(t *testing.T) {
	e := ResolveBonusIcons{Target: Target{Kind: TargetTriggeringCreature}}
	if got := e.Text(); got != "resolve that card's bonus icons" {
		t.Errorf("text = %q", got)
	}
}

func TestResolveBonusIconsValidate(t *testing.T) {
	if (ResolveBonusIcons{}).validate() == nil {
		t.Error("unset target should fail validation")
	}
	if err := (ResolveBonusIcons{
		Target: Target{Kind: TargetTriggeringCreature},
	}).validate(); err != nil {
		t.Errorf("valid target failed: %v", err)
	}
}

// The node resolves the bonus icons on the card in context ("it"), even one that
// is not in play — a card sitting in a discard pile still resolves its icons.
func TestResolveBonusIconsResolvesForeignCard(t *testing.T) {
	g := started(t)
	relic := g.AddToDiscard(
		NewCard("Relic", Brobnar, Artifact, Common, WithBonus(BonusAember, BonusAember)), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0, It: relic, HasIt: true}
	before := g.State.Aember[0]
	ResolveBonusIcons{Target: Target{Kind: TargetTriggeringCreature}}.Resolve(ctx)
	if got := g.State.Aember[0] - before; got != 2 {
		t.Fatalf("aember gained = %d, want 2", got)
	}
}

// With no card in context the node selects nothing and resolves no icons.
func TestResolveBonusIconsNoContextCard(t *testing.T) {
	g := started(t)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	before := g.State.Aember[0]
	ResolveBonusIcons{Target: Target{Kind: TargetTriggeringCreature}}.Resolve(ctx)
	if g.State.Aember[0] != before {
		t.Fatalf("aember changed with no context card: %d -> %d", before, g.State.Aember[0])
	}
}

// A bar (Master of the Grey) stops a foreign resolution too — the icons on the
// named card resolve none of their effects.
func TestResolveBonusIconsOnBarred(t *testing.T) {
	g := started(t)
	g.AddToBattleline(NewCard("Grey", Sanctum, Creature, Rare, WithPower(4),
		WithRestrictions(Restrictions{BonusIcons: Opponent})), 1)
	relic := g.AddToDiscard(
		NewCard("Relic", Brobnar, Artifact, Common, WithBonus(BonusAember)), 0)
	before := g.State.Aember[0]
	g.ResolveBonusIconsOn(0, relic)
	if g.State.Aember[0] != before {
		t.Fatalf("barred foreign resolution gained aember: %d -> %d", before, g.State.Aember[0])
	}
}

func TestExtraBonusIconResolutionText(t *testing.T) {
	want := "the next time you play a card this turn, resolve each of its bonus " +
		"icons an additional time"
	if got := (ExtraBonusIconResolution{}).Text(); got != want {
		t.Errorf("text = %q", got)
	}
}

// Resolve arms a one-shot boost the play path later consumes.
func TestExtraBonusIconResolutionArms(t *testing.T) {
	g := started(t)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	ExtraBonusIconResolution{}.Resolve(ctx)
	if !g.consumeBonusIconBoost(0) {
		t.Fatal("boost was not armed")
	}
	if g.consumeBonusIconBoost(0) {
		t.Fatal("boost should be consumed after one read")
	}
}

// Consuming a boost that is not the last lasting effect compacts the registry so a
// later effect slides down into its place and survives.
func TestConsumeBonusIconBoostCompactsRegistry(t *testing.T) {
	g := started(t)
	g.AddLasting(LastingEffect{On: EventBonusIconBoost, Do: actResolveBonusAgain, Once: true})
	g.AddLasting(LastingEffect{On: EventReapAember, Do: actSteal, Controller: 0})
	if !g.consumeBonusIconBoost(0) {
		t.Fatal("boost was not armed")
	}
	if g.State.LastingCount != 1 {
		t.Fatalf("lasting count = %d, want 1", g.State.LastingCount)
	}
	if got := g.State.Lasting[0]; got.On != EventReapAember || got.Do != actSteal {
		t.Errorf("surviving lasting = %+v, want the reap replacement", got)
	}
}

// An armed boost makes each icon of the next played card resolve an additional
// time, interleaved with the first — a single Æmber bonus gains 2.
func TestExtraBonusIconResolutionDoublesIcons(t *testing.T) {
	g := started(t)
	src := g.AddToBattleline(
		NewCard("Bud", Untamed, Creature, Common, WithPower(2), WithBonus(BonusAember)), 0)
	g.AddLasting(LastingEffect{On: EventBonusIconBoost, Do: actResolveBonusAgain, Once: true})
	before := g.State.Aember[0]
	g.resolveBonusIcons(0, src)
	if got := g.State.Aember[0] - before; got != 2 {
		t.Fatalf("aember gained = %d, want 2 (icon + boost)", got)
	}
	// The boost fires for one card only: the next play resolves the icon once.
	before2 := g.State.Aember[0]
	g.resolveBonusIcons(0, src)
	if got := g.State.Aember[0] - before2; got != 1 {
		t.Fatalf("second play aember = %d, want 1 (boost spent)", got)
	}
}

// The boost is spent on the next card played even when that card has no icons.
func TestExtraBonusIconResolutionSpentOnIconlessCard(t *testing.T) {
	g := started(t)
	iconless := g.AddToBattleline(
		NewCard("Blank", Untamed, Creature, Common, WithPower(2)), 0)
	withIcon := g.AddToBattleline(
		NewCard("Bud", Untamed, Creature, Common, WithPower(2), WithBonus(BonusAember)), 0)
	g.AddLasting(LastingEffect{On: EventBonusIconBoost, Do: actResolveBonusAgain, Once: true})
	g.resolveBonusIcons(0, iconless) // consumes the boost, does nothing
	before := g.State.Aember[0]
	g.resolveBonusIcons(0, withIcon)
	if got := g.State.Aember[0] - before; got != 1 {
		t.Fatalf("aember gained = %d, want 1 (boost already spent)", got)
	}
}

// A bar arriving after the first resolution stops the boost's extra resolution.
func TestExtraBonusIconResolutionBoostRespectsBar(t *testing.T) {
	g := started(t)
	src := g.AddToBattleline(
		NewCard("Bud", Untamed, Creature, Common, WithPower(2), WithBonus(BonusAember)), 0)
	g.AddToBattleline(NewCard("Grey", Sanctum, Creature, Rare, WithPower(4),
		WithRestrictions(Restrictions{BonusIcons: Opponent})), 1)
	g.AddLasting(LastingEffect{On: EventBonusIconBoost, Do: actResolveBonusAgain, Once: true})
	before := g.State.Aember[0]
	g.resolveBonusIcons(0, src)
	if got := g.State.Aember[0] - before; got != 0 {
		t.Fatalf("aember gained = %d, want 0 (barred)", got)
	}
}
