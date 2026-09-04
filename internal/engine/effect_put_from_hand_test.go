package engine

import "testing"

func TestPutFromHand(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	returned := g.AddToHand(testCreature("Old Guy", 1), 0)
	mars := g.AddToHand(NewCard("New Guy", Mars, Creature, Common, WithPower(2)), 0)
	sameName := g.AddToHand(NewCard("Old Guy", Mars, Creature, Common, WithPower(2)), 0)
	offHouse := g.AddToHand(NewCard("Off House", Brobnar, Creature, Common, WithPower(2)), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	ctx.It, ctx.HasIt = returned, true

	e := PutFromHand{Type: Creature, House: Mars, ExceptSameName: true}
	if e.Text() != "put a Mars creature with a different name from your hand into play" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)

	if !g.inPlay(mars) {
		t.Error("the differently named Mars creature should have entered play")
	}
	if g.inPlay(sameName) {
		t.Error("a card sharing the returned card's name should not be a candidate")
	}
	if g.inPlay(offHouse) {
		t.Error("a card of another house should not be a candidate")
	}
	if !ctx.HasIt || ctx.It != mars {
		t.Errorf("ctx.It = %v (HasIt %v), want %d", ctx.It, ctx.HasIt, mars)
	}

	// No candidate matches: nothing happens.
	ctx2 := &EffectContext{Resolver: g, Source: src, Controller: 0}
	(PutFromHand{Type: Artifact}).Resolve(ctx2)
	if ctx2.HasIt {
		t.Error("a Resolve with no candidates should leave ctx.It unset")
	}
}
