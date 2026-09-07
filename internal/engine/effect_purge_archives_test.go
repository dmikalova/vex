package engine

import "testing"

func TestPurgeArchivesForDamageText(t *testing.T) {
	e := PurgeArchivesForDamage{Amount: 2, Target: Target{Kind: TargetTriggeringCreature}}
	want := "purge any number of cards from your archives to deal an additional 2 " +
		"damage to it for each card purged this way"
	if got := e.Text(); got != want {
		t.Errorf("Text() = %q, want %q", got, want)
	}
}

func TestPurgeArchivesForDamageValidate(t *testing.T) {
	if (PurgeArchivesForDamage{Amount: 2}).validate() == nil {
		t.Error("an unset target should not validate")
	}
	if (PurgeArchivesForDamage{Target: Target{Kind: TargetTriggeringCreature}}).validate() == nil {
		t.Error("a zero amount should not validate")
	}
	if (PurgeArchivesForDamage{
		Amount: 2, Target: Target{Kind: TargetTriggeringCreature},
	}).validate() != nil {
		t.Error("a positive amount with a target should validate")
	}
}

func TestPurgeArchivesForDamageResolve(t *testing.T) {
	g := started(t)
	foe := g.AddToBattleline(testCreature("Foe", 9), 1)
	one := g.AddToArchives(NewCard("Archived One", Logos, Creature, Common), 0)
	two := g.AddToArchives(NewCard("Archived Two", Logos, Creature, Common), 0)
	g.SetChooser(0, &declineAfterChooser{ids: []LocalID{one, two}})

	ctx := &EffectContext{Resolver: g, Controller: 0, It: foe, HasIt: true}
	PurgeArchivesForDamage{Amount: 2, Target: Target{Kind: TargetTriggeringCreature}}.
		Resolve(ctx)

	if got := g.Damage(foe); got != 4 {
		t.Errorf("damage = %d, want 4 (2 per card, 2 purged)", got)
	}
	if g.State.Archives[0].contains(one) || g.State.Archives[0].contains(two) {
		t.Error("both archived cards should have been purged")
	}
	if !g.State.Purge[0].contains(one) || !g.State.Purge[0].contains(two) {
		t.Error("both purged cards should be in the purge pile")
	}
}

func TestPurgeArchivesForDamagePurgingNone(t *testing.T) {
	g := started(t)
	foe := g.AddToBattleline(testCreature("Foe", 9), 1)
	g.AddToArchives(NewCard("Archived", Logos, Creature, Common), 0)
	g.SetChooser(0, &declineAfterChooser{}) // decline immediately

	ctx := &EffectContext{Resolver: g, Controller: 0, It: foe, HasIt: true}
	PurgeArchivesForDamage{Amount: 2, Target: Target{Kind: TargetTriggeringCreature}}.
		Resolve(ctx)

	if got := g.Damage(foe); got != 0 {
		t.Errorf("damage = %d, want 0 when nothing is purged", got)
	}
}

func TestPurgeArchivesForDamageEmptyArchives(t *testing.T) {
	g := started(t)
	foe := g.AddToBattleline(testCreature("Foe", 9), 1)

	ctx := &EffectContext{Resolver: g, Controller: 0, It: foe, HasIt: true}
	PurgeArchivesForDamage{Amount: 2, Target: Target{Kind: TargetTriggeringCreature}}.
		Resolve(ctx)

	if got := g.Damage(foe); got != 0 {
		t.Errorf("damage = %d, want 0 with empty archives", got)
	}
}
