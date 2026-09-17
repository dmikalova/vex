package cardtest

// CardExpect is a fluent set of assertions about a single card. Each method
// reports a mismatch with t.Errorf and returns the receiver so assertions chain:
//
//	h.Expect(troll).Exhausted().Damage(6).At(ct.PlayArea)
type CardExpect struct {
	h *Harness
	c Card
}

// name returns the card's name for messages.
func (e CardExpect) name() string { return e.h.g.Name(e.c.id) }

// Damage asserts the creature's current damage.
func (e CardExpect) Damage(n int) CardExpect {
	e.h.t.Helper()
	if got := e.c.Damage(); got != n {
		e.h.t.Errorf("%s damage = %d, want %d", e.name(), got, n)
	}
	return e
}

// Power asserts the creature's current power.
func (e CardExpect) Power(n int) CardExpect {
	e.h.t.Helper()
	if got := e.c.Power(); got != n {
		e.h.t.Errorf("%s power = %d, want %d", e.name(), got, n)
	}
	return e
}

// Armor asserts the creature's current armor.
func (e CardExpect) Armor(n int) CardExpect {
	e.h.t.Helper()
	if got := e.c.Armor(); got != n {
		e.h.t.Errorf("%s armor = %d, want %d", e.name(), got, n)
	}
	return e
}

// AmberOn asserts the Æmber sitting on the card.
func (e CardExpect) AmberOn(n int) CardExpect {
	e.h.t.Helper()
	if got := e.c.AmberOn(); got != n {
		e.h.t.Errorf("%s Æmber-on-card = %d, want %d", e.name(), got, n)
	}
	return e
}

// Exhausted asserts the card is exhausted.
func (e CardExpect) Exhausted() CardExpect {
	e.h.t.Helper()
	if !e.c.Exhausted() {
		e.h.t.Errorf("%s is ready, want exhausted", e.name())
	}
	return e
}

// Ready asserts the card is not exhausted.
func (e CardExpect) Ready() CardExpect {
	e.h.t.Helper()
	if e.c.Exhausted() {
		e.h.t.Errorf("%s is exhausted, want ready", e.name())
	}
	return e
}

// Stunned asserts the creature's stun status.
func (e CardExpect) Stunned(want bool) CardExpect {
	e.h.t.Helper()
	if got := e.c.Stunned(); got != want {
		e.h.t.Errorf("%s stunned = %v, want %v", e.name(), got, want)
	}
	return e
}

// Enraged asserts the creature's enrage status.
func (e CardExpect) Enraged(want bool) CardExpect {
	e.h.t.Helper()
	if got := e.c.Enraged(); got != want {
		e.h.t.Errorf("%s enraged = %v, want %v", e.name(), got, want)
	}
	return e
}

// At asserts the card is in the given zone.
func (e CardExpect) At(z Zone) CardExpect {
	e.h.t.Helper()
	if got := e.c.Location(); got != z {
		e.h.t.Errorf("%s is in %s, want %s", e.name(), got, z)
	}
	return e
}
