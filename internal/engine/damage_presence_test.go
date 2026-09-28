package engine

import "testing"

// TestDamagePresenceIsTotal pins what each damage presence admits and the
// adjective it prints, including the zero value, which asks nothing about damage
// and prints no adjective.
func TestDamagePresenceIsTotal(t *testing.T) {
	want := map[DamagePresence]struct {
		undamaged bool
		damaged   bool
		adjective string
	}{
		damageAny:  {true, true, ""},
		DamageSome: {false, true, "damaged"},
		DamageNone: {true, false, "undamaged"},
	}
	for _, presence := range append(allDamagePresences(), damageAny) {
		expected, ok := want[presence]
		if !ok {
			t.Fatalf("damage presence %d has no pinned row; add one", presence)
		}
		if got := presence.admits(0); got != expected.undamaged {
			t.Errorf("presence %d admits(0) = %v, want %v", presence, got, expected.undamaged)
		}
		if got := presence.admits(2); got != expected.damaged {
			t.Errorf("presence %d admits(2) = %v, want %v", presence, got, expected.damaged)
		}
		if got := presence.adjective(); got != expected.adjective {
			t.Errorf("presence %d adjective = %q, want %q", presence, got, expected.adjective)
		}
		if got := presence.filters(); got != (presence != damageAny) {
			t.Errorf("presence %d filters = %v", presence, got)
		}
	}
}
