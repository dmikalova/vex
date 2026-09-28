package engine

import "testing"

// TestItAttachedToThisOrNeighborCondText covers the condition Commander Dhrxgar
// reads, naming the source card via the self placeholder.
func TestItAttachedToThisOrNeighborCondText(t *testing.T) {
	want := "if it is attached to " + SelfName + " or one of its neighbors"
	if got := (ItAttachedToThisOrNeighbor{}).CondText(); got != want {
		t.Errorf("CondText = %q, want %q", got, want)
	}
}

// TestItIsNamed covers Chain Gang's by-name gate, now an axis of the one ItIs
// condition: a Filter naming a card renders the name outright, with no article,
// and is met only when a card of that name is in context.
func TestItIsNamed(t *testing.T) {
	if got := (ItIs{Filter: Filter{Name: "Subtle Chain"}}).CondText(); got != "if it is Subtle Chain" {
		t.Errorf("CondText() = %q", got)
	}

	g := NewGame("Alice", "Bob", 1)
	chain := g.AddToDeck(NewCard("Subtle Chain", Dis, Tactic, Common), 0)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}

	// With no card in context the condition is never met.
	if (ItIs{Filter: Filter{Name: "Subtle Chain"}}).Met(ctx) {
		t.Error("Met with no context card should be false")
	}
	ctx.It, ctx.HasIt = chain, true
	if !(ItIs{Filter: Filter{Name: "Subtle Chain"}}).Met(ctx) {
		t.Error("the named card should meet the by-name condition")
	}
	if (ItIs{Filter: Filter{Name: "Mind Barb"}}).Met(ctx) {
		t.Error("a differently-named card should not meet the condition")
	}
}
