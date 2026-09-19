package engine

import "testing"

// TestItIsNamed covers Chain Gang's by-name gate: the condition renders the card
// name outright and is met only when a card of that name is in context.
func TestItIsNamed(t *testing.T) {
	if got := (ItIsNamed{Name: "Subtle Chain"}).CondText(); got != "if it is Subtle Chain" {
		t.Errorf("CondText() = %q", got)
	}

	g := NewGame("Alice", "Bob", 1)
	chain := g.AddToDeck(NewCard("Subtle Chain", Dis, Tactic, Common), 0)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}

	// With no card in context the condition is never met.
	if (ItIsNamed{Name: "Subtle Chain"}).Met(ctx) {
		t.Error("Met with no context card should be false")
	}
	ctx.It, ctx.HasIt = chain, true
	if !(ItIsNamed{Name: "Subtle Chain"}).Met(ctx) {
		t.Error("the named card should meet the by-name condition")
	}
	if (ItIsNamed{Name: "Mind Barb"}).Met(ctx) {
		t.Error("a differently-named card should not meet the condition")
	}
}
