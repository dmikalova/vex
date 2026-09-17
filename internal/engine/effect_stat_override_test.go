package engine

import "testing"

// TestOverrideStatsMasksPowerAndArmor checks the stat override reads as a live
// mask: a creature's power and armor become the masked values regardless of
// counters or bonuses, and the real values are revealed again when it lifts.
func TestOverrideStatsMasksPowerAndArmor(t *testing.T) {
	g := started(t)
	c := g.AddToBattleline(testCreature("c", 5, WithArmor(2)), 0)
	g.AddPowerCounter(c, 3) // real power 8

	OverrideStats{
		Power: 1, HasPower: true, Armor: 0, HasArmor: true,
		Duration: RemainderOfPlayerTurn,
	}.Resolve(&EffectContext{Resolver: g, Controller: 0})

	if got := g.Power(c); got != 1 {
		t.Errorf("masked power = %d, want 1", got)
	}
	if got := g.Armor(c); got != 0 {
		t.Errorf("masked armor = %d, want 0", got)
	}

	// End of turn lifts the mask; the real 8 power / 2 armor are revealed again, and
	// the stored counters were never touched.
	g.EndPlayPhase(0)
	if got := g.Power(c); got != 8 {
		t.Errorf("revealed power = %d, want 8", got)
	}
	if got := g.Armor(c); got != 2 {
		t.Errorf("revealed armor = %d, want 2", got)
	}
}

// TestOverrideStatsMaskedArmorAbsorbsNothing checks the masked armor value governs
// combat absorption without spending the real armor pool.
func TestOverrideStatsMaskedArmorAbsorbsNothing(t *testing.T) {
	g := started(t)
	c := g.AddToBattleline(testCreature("c", 9, WithArmor(3)), 0)

	OverrideStats{
		Armor: 0, HasArmor: true,
		Duration: RemainderOfPlayerTurn,
	}.Resolve(&EffectContext{Resolver: g, Controller: 0})

	g.applyRawDamage(DamageTarget{ID: c, Amount: 2})
	if g.Damage(c) != 2 {
		t.Errorf("masked-armor creature took %d, want 2 (0 armor absorbed)", g.Damage(c))
	}
	// The real armor pool was untouched: it absorbs again once the mask lifts.
	g.EndPlayPhase(0)
	if got := g.Armor(c); got != 3 {
		t.Errorf("revealed armor = %d, want 3", got)
	}
}

// TestOverrideStatsValidate covers the guard rails.
func TestOverrideStatsValidate(t *testing.T) {
	if (OverrideStats{Duration: RemainderOfPlayerTurn}).validate() == nil {
		t.Error("no masked stat should be invalid")
	}
	if (OverrideStats{Power: 1, HasPower: true, Duration: OpponentNextTurn}).validate() == nil {
		t.Error("unsupported duration should be invalid")
	}
	if (OverrideStats{Power: 1, HasPower: true, Duration: RemainderOfPlayerTurn}).validate() != nil {
		t.Error("a masked power with a supported duration should be valid")
	}
	if got := (OverrideStats{Armor: 0, HasArmor: true, Duration: RemainderOfPlayerTurn}).
		Text(); got != "for the remainder of the turn, each creature is considered to have 0 armor" {
		t.Errorf("armor-only text = %q", got)
	}
	if got := (OverrideStats{Power: 1, HasPower: true, Duration: RemainderOfPlayerTurn}).
		Text(); got != "for the remainder of the turn, each creature is considered to have 1 power" {
		t.Errorf("power-only text = %q", got)
	}
}

// TestOverrideStatsDestroysNewlyLethal checks that masking power to 1 destroys a
// creature whose marked damage now meets it, and that a ward absorbs the
// destruction once before the still-lethal state destroys it for good.
func TestOverrideStatsDestroysNewlyLethal(t *testing.T) {
	g := started(t)
	damaged := g.AddToBattleline(testCreature("damaged", 5), 0)
	g.SetDamage(damaged, 3) // alive at 5 power
	warded := g.AddToBattleline(testCreature("warded", 5), 0)
	g.SetDamage(warded, 3)
	g.SetWarded(warded, true)

	OverrideStats{
		Power: 1, HasPower: true,
		Duration: RemainderOfPlayerTurn,
	}.Resolve(&EffectContext{Resolver: g, Controller: 0})
	g.settleDestroyed(0) // the resolution boundary settles the mask (ADR 0029)

	if g.inPlay(damaged) {
		t.Error("a creature at 3 damage should die once masked to 1 power")
	}
	if g.inPlay(warded) {
		t.Error("ward absorbs one destruction, then the still-lethal state destroys it")
	}
}
