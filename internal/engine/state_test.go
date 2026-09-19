package engine

import (
	"reflect"
	"testing"
)

func TestCardList(t *testing.T) {
	var z deckList
	z.add(1)
	z.add(2)
	z.add(3)
	z.addFront(0) // [0 1 2 3]
	if z.indexOf(2) != 2 {
		t.Errorf("indexOf(2) = %d", z.indexOf(2))
	}
	if z.indexOf(9) != -1 {
		t.Errorf("indexOf(9) = %d, want -1", z.indexOf(9))
	}
	if !z.contains(3) || z.contains(9) {
		t.Error("contains check failed")
	}
	if got := z.removeAt(0); got != 0 { // [1 2 3]
		t.Errorf("removeAt(0) = %d, want 0", got)
	}
	if !z.remove(2) { // [1 3]
		t.Error("remove(2) should succeed")
	}
	if z.remove(9) {
		t.Error("remove(9) should fail")
	}
	got := z.slice()
	if len(got) != 2 || got[0] != 1 || got[1] != 3 {
		t.Errorf("zone slice = %v, want [1 3]", got)
	}
}

func TestCatalogAddPanicsOverCapacity(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	def := testCreature("Filler", 1)
	for range maxCards {
		g.Register(def, 0)
	}
	defer func() {
		if recover() == nil {
			t.Error("Register past maxCards should panic")
		}
	}()
	g.Register(def, 0)
}

func TestFastCopyIsIndependent(t *testing.T) {
	g := NewGame("A", "B", 1)
	def := NewCard("c", Brobnar, Creature, Common, WithPower(5))
	id := g.AddToBattleline(def, 0)
	g.State.Cards[id].Damage = 3
	g.State.Aember[0] = 2

	clone := g.State.FastCopy()
	clone.Cards[id].Damage = 99
	clone.Aember[0] = 99
	clone.Battleline[0].add(id)

	if g.State.Cards[id].Damage != 3 {
		t.Errorf("original damage mutated: %d", g.State.Cards[id].Damage)
	}
	if g.State.Aember[0] != 2 {
		t.Errorf("original aember mutated: %d", g.State.Aember[0])
	}
	if g.State.Battleline[0].Count != 1 {
		t.Errorf("original battleline mutated: count %d", g.State.Battleline[0].Count)
	}
}

func TestTurnLogSaturates(t *testing.T) {
	var log turnLog
	for i := range turnLogCap + 5 {
		log.add(LocalID(i%100 + 1))
	}
	if int(log.Count) != turnLogCap {
		t.Errorf("count = %d, want %d", log.Count, turnLogCap)
	}
	if got := len(log.slice()); got != turnLogCap {
		t.Errorf("slice length = %d, want %d", got, turnLogCap)
	}
	log.reset()
	if log.Count != 0 {
		t.Errorf("count after reset = %d, want 0", log.Count)
	}
}

// narration is how the log covers changes to one GameState field.
type narration int

const (
	// narratedDirectly: a change to the field has a log entry of its own, so a
	// player watching the log sees the change itself.
	narratedDirectly narration = iota
	// narratedByCause: the field is bookkeeping for a lasting effect or a turn
	// bar. The log names the card that armed it and narrates the outcome when it
	// bites; the flag flipping is not itself an outcome (ADR 0011).
	narratedByCause
	// narratedNever: deliberately silent.
	narratedNever
)

// fieldNarration classifies every GameState field. It exists to force a decision
// rather than to describe one: TestEveryStateFieldDeclaresItsNarration fails when
// a field is added without an entry here, so "does this change need a log entry?"
// is answered when the field is written instead of being discovered as a missing
// line in a game (Exhaust and ReadyIfFirstUse were both found that way).
var fieldNarration = map[string]narration{
	"Cards":          narratedDirectly,
	"UsagesThisTurn": narratedDirectly,
	"Battleline":     narratedDirectly,
	"Hand":           narratedDirectly,
	"Deck":           narratedDirectly,
	"Discard":        narratedDirectly,
	"Artifacts":      narratedDirectly,
	"Archives":       narratedDirectly,
	"Purge":          narratedDirectly,
	"Aember":         narratedDirectly,
	"KeyColors":      narratedDirectly,
	"Chains":         narratedDirectly,
	"ActivePlayer":   narratedDirectly,
	"ActiveHouse":    narratedDirectly,
	"Tide":           narratedDirectly,
	"Turn":           narratedDirectly,
	"Winner":         narratedDirectly,
	"Phase":          narratedDirectly,
	"Counters":       narratedDirectly,
	"CounterCount":   narratedDirectly,
	"Controls":       narratedDirectly,
	"ControlCount":   narratedDirectly,

	"ForgePrevented":              narratedByCause,
	"PhaseEnded":                  narratedByCause,
	"CannotFight":                 narratedByCause,
	"CannotFightNext":             narratedByCause,
	"CannotPlayTypeThis":          narratedByCause,
	"CannotPlayTypeNext":          narratedByCause,
	"CannotUse":                   narratedByCause,
	"CannotUseNext":               narratedByCause,
	"CannotReap":                  narratedByCause,
	"CannotReapNext":              narratedByCause,
	"CannotReapHouse":             narratedByCause,
	"CannotReapHouseNext":         narratedByCause,
	"CreaturesCannot":             narratedByCause,
	"CreaturesCannotNext":         narratedByCause,
	"SkipForge":                   narratedByCause,
	"SkipForgeNext":               narratedByCause,
	"Scheduled":                   narratedByCause,
	"ScheduledCount":              narratedByCause,
	"KeyCostBump":                 narratedByCause,
	"KeyCostBumpNext":             narratedByCause,
	"KeyCostPerHouse":             narratedByCause,
	"KeyCostPerHouseNext":         narratedByCause,
	"MayFightHouse":               narratedByCause,
	"MayFightAny":                 narratedByCause,
	"MayUseHouse":                 narratedByCause,
	"MayPlayHouse":                narratedByCause,
	"MayUseArtifactsAnyHouse":     narratedByCause,
	"MayUseTrait":                 narratedByCause,
	"TurnHistory":                 narratedByCause,
	"Lasting":                     narratedByCause,
	"LastingCount":                narratedByCause,
	"Continuous":                  narratedByCause,
	"ContinuousCount":             narratedByCause,
	"AlsoTriggers":                narratedByCause,
	"AlsoTriggersCount":           narratedByCause,
	"PlayedThisTurn":              narratedByCause,
	"DiscardedThisTurn":           narratedByCause,
	"PlayPermissionsUsedThisTurn": narratedByCause,
	"OffHousePermits":             narratedByCause,
	"OffHousePermitCount":         narratedByCause,
	"NonActivePlaysUsedThisTurn":  narratedByCause,
	"FirstTurnPlayLimit":          narratedByCause,
	"HouseConstraints":            narratedByCause,
	"HouseConstraintCount":        narratedByCause,
	"HouseConstraintsNext":        narratedByCause,
	"HouseConstraintCountNext":    narratedByCause,
	"FightDamageRedirect":         narratedByCause,
	"FightCancelled":              narratedByCause,
	"FightersPlus":                narratedByCause,

	// The match RNG is state so a snapshot replays bit-exact (ADR 0039); its
	// advancing is not an outcome anyone can observe.
	"PRNG": narratedNever,
}

// TestEveryStateFieldDeclaresItsNarration is the ratchet behind "every state
// change is logged". It cannot prove the log is complete, but it can stop the
// gap from being introduced silently: a new GameState field fails the build until
// fieldNarration says how the log covers it.
func TestEveryStateFieldDeclaresItsNarration(t *testing.T) {
	typ := reflect.TypeFor[GameState]()
	for field := range typ.Fields() {
		name := field.Name
		if _, ok := fieldNarration[name]; !ok {
			t.Errorf(
				"GameState.%s has no fieldNarration entry; decide whether a change "+
					"to it needs its own log entry", name,
			)
		}
	}
	for name := range fieldNarration {
		if _, ok := typ.FieldByName(name); !ok {
			t.Errorf("fieldNarration names %q, which GameState no longer has", name)
		}
	}
}
