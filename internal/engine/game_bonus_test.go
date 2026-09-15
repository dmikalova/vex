package engine

import (
	"strings"
	"testing"
)

// resolveBonusIcons is the play-time bonus-icon step. These tests drive it
// directly on cards placed in play, scripting the target choices.

func TestBonusDamagePicksAndHits(t *testing.T) {
	g := started(t)
	src := g.AddToBattleline(
		NewCard("Zapper", Brobnar, Creature, Common, WithPower(5), WithBonus(BonusDamage)), 0)
	foe := g.AddToBattleline(NewCard("Foe", Brobnar, Creature, Common, WithPower(3)), 1)
	g.SetChooser(0, &idQueueChooser{ids: []LocalID{foe}})
	g.resolveBonusIcons(0, src)
	if got := g.Damage(foe); got != 1 {
		t.Fatalf("foe damage = %d, want 1", got)
	}
}

func TestBonusDamageNoCreaturesDoesNothing(t *testing.T) {
	g := started(t)
	src := g.AddArtifact(NewCard("Zap", Brobnar, Artifact, Common, WithBonus(BonusDamage)), 0)
	g.SetChooser(0, orderRejectChooser{})
	g.resolveBonusIcons(0, src) // no creatures in play: nothing to hit
}

func TestBonusDamageDeclineDoesNothing(t *testing.T) {
	g := started(t)
	src := g.AddArtifact(NewCard("Zap", Brobnar, Artifact, Common, WithBonus(BonusDamage)), 0)
	foe := g.AddToBattleline(NewCard("Foe", Brobnar, Creature, Common, WithPower(3)), 1)
	g.AddToBattleline(NewCard("Foe2", Brobnar, Creature, Common, WithPower(3)), 1)
	g.SetChooser(0, orderRejectChooser{})
	g.resolveBonusIcons(0, src)
	if got := g.Damage(foe); got != 0 {
		t.Fatalf("foe damage = %d, want 0", got)
	}
}

func TestBonusCaptureMovesAember(t *testing.T) {
	g := started(t)
	src := g.AddToBattleline(
		NewCard("Captor", Brobnar, Creature, Common, WithPower(5), WithBonus(BonusCapture)), 0)
	g.State.Aember[1] = 2
	g.SetChooser(0, &idQueueChooser{ids: []LocalID{src}})
	g.resolveBonusIcons(0, src)
	if amt := g.AmberOn(src); amt != 1 {
		t.Fatalf("captured Æmber = %d, want 1", amt)
	}
	if g.State.Aember[1] != 1 {
		t.Fatalf("opponent pool = %d, want 1", g.State.Aember[1])
	}
}

func TestBonusCaptureNeedsOpponentAember(t *testing.T) {
	g := started(t)
	src := g.AddToBattleline(
		NewCard("Captor", Brobnar, Creature, Common, WithPower(5), WithBonus(BonusCapture)), 0)
	g.resolveBonusIcons(0, src) // opponent pool empty: nothing captured
	if amt := g.AmberOn(src); amt != 0 {
		t.Fatalf("captured Æmber = %d, want 0", amt)
	}
}

func TestBonusCaptureNeedsFriendlyCreature(t *testing.T) {
	g := started(t)
	src := g.AddArtifact(NewCard("Vault", Brobnar, Artifact, Common, WithBonus(BonusCapture)), 0)
	g.State.Aember[1] = 2
	g.resolveBonusIcons(0, src) // no friendly creature to hold the Æmber
	if g.State.Aember[1] != 2 {
		t.Fatalf("opponent pool = %d, want 2 (nothing captured)", g.State.Aember[1])
	}
}

func TestBonusCaptureDeclineDoesNothing(t *testing.T) {
	g := started(t)
	src := g.AddToBattleline(
		NewCard("Captor", Brobnar, Creature, Common, WithPower(5), WithBonus(BonusCapture)), 0)
	g.AddToBattleline(NewCard("Ally", Brobnar, Creature, Common, WithPower(5)), 0)
	g.State.Aember[1] = 2
	g.SetChooser(0, orderRejectChooser{})
	g.resolveBonusIcons(0, src)
	if amt := g.AmberOn(src); amt != 0 {
		t.Fatalf("captured Æmber = %d, want 0", amt)
	}
	if g.State.Aember[1] != 2 {
		t.Fatalf("opponent pool = %d, want 2", g.State.Aember[1])
	}
}

func TestBonusDrawDrawsACard(t *testing.T) {
	g := started(t)
	g.AddToDeck(NewCard("Top", Brobnar, Creature, Common, WithPower(3)), 0)
	src := g.AddArtifact(NewCard("Drawer", Brobnar, Artifact, Common, WithBonus(BonusDraw)), 0)
	before := len(g.Hand(0))
	g.resolveBonusIcons(0, src)
	if got := len(g.Hand(0)); got != before+1 {
		t.Fatalf("hand size = %d, want %d", got, before+1)
	}
}

func TestBonusAemberGainedToPool(t *testing.T) {
	g := started(t)
	src := g.AddArtifact(NewCard("Coin", Brobnar, Artifact, Common, WithBonus(BonusAember)), 0)
	before := g.State.Aember[0]
	g.resolveBonusIcons(0, src)
	if g.State.Aember[0] != before+1 {
		t.Fatalf("pool = %d, want %d", g.State.Aember[0], before+1)
	}
}

func TestBonusAemberInterceptedByCaptor(t *testing.T) {
	g := started(t)
	spider := g.AddToBattleline(NewCard("Spider", Brobnar, Creature, Common, WithPower(3),
		WithReplaces(Instead{Of: EventAemberAddedToPool, Player: Opponent, With: Capture})), 1)
	src := g.AddArtifact(NewCard("Coin", Brobnar, Artifact, Common, WithBonus(BonusAember)), 0)
	g.resolveBonusIcons(0, src)
	if g.State.Aember[0] != 0 {
		t.Fatalf("pool = %d, want 0 (intercepted)", g.State.Aember[0])
	}
	if amt := g.AmberOn(spider); amt != 1 {
		t.Fatalf("captor Æmber = %d, want 1", amt)
	}
}

func TestBonusIconsStopWhenCreatureLeavesPlay(t *testing.T) {
	g := started(t)
	// A 1-power creature that damages itself with its first icon dies before its
	// second (Æmber) icon can resolve — the Vactrol leaves-play divergence.
	src := g.AddToBattleline(NewCard("Fragile", Brobnar, Creature, Common,
		WithPower(1), WithBonus(BonusDamage, BonusAember)), 0)
	g.SetChooser(0, &idQueueChooser{ids: []LocalID{src}})
	before := g.State.Aember[0]
	g.resolveBonusIcons(0, src)
	if g.inPlay(src) {
		t.Fatalf("source should be destroyed by its own bonus damage")
	}
	if g.State.Aember[0] != before {
		t.Fatalf("pool = %d, want %d (second icon must not resolve)", g.State.Aember[0], before)
	}
}

func TestBonusIconsSkippedIfCreatureAlreadyGone(t *testing.T) {
	g := started(t)
	src := g.AddToBattleline(
		NewCard("Ghost", Brobnar, Creature, Common, WithPower(3), WithBonus(BonusAember)), 0)
	g.removeFromAnyZone(src)
	before := g.State.Aember[0]
	g.resolveBonusIcons(0, src) // not in play: no icon resolves
	if g.State.Aember[0] != before {
		t.Fatalf("pool = %d, want %d", g.State.Aember[0], before)
	}
}

func TestBonusIconHelpers(t *testing.T) {
	cases := map[BonusIcon]string{
		BonusAember:  "Æmber",
		BonusCapture: "Capture",
		BonusDamage:  "Damage",
		BonusDraw:    "Draw",
		bonusUnset:   "Unset",
	}
	for ic, want := range cases {
		if got := ic.String(); got != want {
			t.Errorf("%d.String() = %q, want %q", ic, got, want)
		}
	}
	if bonusUnset.valid() {
		t.Error("bonusUnset should be invalid")
	}
	if !BonusAember.valid() {
		t.Error("BonusAember should be valid")
	}
	if got := bonusIconsText(
		[]BonusIcon{BonusAember, BonusAember, BonusDraw},
	); got != "Æmber Æmber Draw" {
		t.Errorf("bonusIconsText = %q", got)
	}
	def := NewCard("Pips", Brobnar, Tactic, Common, WithBonus(BonusAember, BonusAember))
	if got := def.AemberBonus(); got != 2 {
		t.Errorf("AemberBonus() = %d, want 2", got)
	}
}

func TestNewCardRejectsUnsetBonus(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic on unset bonus icon")
		}
	}()
	NewCard("Bad", Brobnar, Tactic, Common, WithBonus(bonusUnset))
}

func TestNewCardRejectsUnsetEnhance(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic on unset enhance icon")
		}
	}()
	NewCard("Bad", Brobnar, Tactic, Common, WithEnhance(bonusUnset))
}

func TestEnhanceRendersAsRulesLine(t *testing.T) {
	def := NewCard("Splinterish", Brobnar, Creature, Common, WithPower(1),
		WithEnhance(BonusDamage, BonusDamage))
	if got := RenderCardText(&def); !strings.Contains(got, "Enhance Damage Damage.") {
		t.Errorf("text = %q, want the Enhance line", got)
	}
}

func TestWithoutEnhancementBarsAKind(t *testing.T) {
	def := NewCard("OptOut", Brobnar, Creature, Common, WithPower(3),
		WithoutEnhancement(BonusCapture))
	if !def.BarsEnhanceIcon(BonusCapture) {
		t.Error("WithoutEnhancement(Capture) should bar the Capture kind")
	}
	if def.BarsEnhanceIcon(BonusAember) {
		t.Error("a kind that was not barred should still be allowed")
	}
}

func TestNewCardRejectsUnsetNoEnhance(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic on unset WithoutEnhancement icon")
		}
	}()
	NewCard("Bad", Brobnar, Tactic, Common, WithoutEnhancement(bonusUnset))
}
