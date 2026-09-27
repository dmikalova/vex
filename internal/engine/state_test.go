package engine

import (
	"reflect"
	"strings"
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
	// narratedByCause: the field is bookkeeping for a lasting effect, a turn bar,
	// or a step whose outcome is narrated elsewhere. The log names the card that
	// armed it and narrates the outcome when it bites; the field changing is not
	// itself an outcome (ADR 0011).
	narratedByCause
	// narratedNever: deliberately silent.
	narratedNever
)

// narrationLeafTypes are the struct types fieldNarration does not look inside.
// Each names one reason its fields share a single narration, so the escape hatch
// from per-field classification stays visible in review: the alternative is
// hundreds of rows restating the same answer. CardCore is deliberately absent —
// its 40 fields are the state a player watches, and classifying them one by one
// is the point of the recursion.
var narrationLeafTypes = map[string]string{
	"deckList": "a zone's storage: IDs and Count move together, and the zone " +
		"they belong to is what a move entry names.",
	"wideList": "a zone's storage, as deckList.",
	"turnLog": "a play/discard log's storage, as the zone lists: the play or " +
		"discard that appended to it is what the log narrates.",
	"Destination": "one opaque value — a zone plus whose it is — written and " +
		"read as a unit, with no field a card can change on its own.",
	"Bar": "a turn bar: Value is the restriction and Source the card that " +
		"imposed it, armed and lifted as one (ADR 0011).",
	"ScheduledEffect": "a registry entry: the whole row is one armed effect, " +
		"narrated by the card that armed it and again when it fires.",
	"LastingEffect":         "a registry entry, as ScheduledEffect.",
	"ContinuousEffect":      "a registry entry, as ScheduledEffect.",
	"LastingAlsoTriggersOn": "a registry entry, as ScheduledEffect.",
}

// fieldNarration classifies every GameState field, keyed by dotted path through
// the struct-valued fields (Cards.Damage, Controls.Controller). It exists to
// force a decision rather than to describe one: TestEveryStateFieldDeclaresItsNarration
// fails when a field is added without an entry here, so "does this change need a
// log entry?" is answered when the field is written instead of being discovered
// as a missing line in a game (Exhaust and ReadyIfFirstUse were both found that
// way — both CardCore fields, which is why the table recurses rather than
// classifying Cards as one row).
var fieldNarration = map[string]narration{
	// CardCore, the per-card state a player watches. The comment on each group
	// names the entry that narrates it.

	// Use and status: CreatureExhausted and CreatureReadied; CreatureStunned,
	// CreaturesUnstunned and StunRecovered; CreatureEnraged and
	// CreatureEnrageRemoved; CreatureWarded, WardRemoved and WardAbsorbed.
	"Cards.Exhausted": narratedDirectly,
	"Cards.Stunned":   narratedDirectly,
	"Cards.Enraged":   narratedDirectly,
	"Cards.Warded":    narratedDirectly,

	// Granted and lost text: CreatureGainedKeyword, CreatureLostKeyword,
	// KeywordLostByAll, CreatureGainedTrait, CreatureConsideredFlank,
	// CreatureGainedTextBox, CreatureCopiedStats.
	"Cards.GrantedKeywords":           narratedDirectly,
	"Cards.LostKeywords":              narratedDirectly,
	"Cards.KeywordsUntilNextTurn":     narratedDirectly,
	"Cards.LostKeywordsUntilNextTurn": narratedDirectly,
	"Cards.TraitUntilNextTurn":        narratedDirectly,
	"Cards.ConsideredFlank":           narratedDirectly,
	"Cards.TextBoxSourcePlus":         narratedDirectly,
	"Cards.TextBoxTurnSourcePlus":     narratedDirectly,
	"Cards.CopiedStatsSourcePlus":     narratedDirectly,

	// Damage and armor: DamageTaken, AbilityDamageDealt, AssaultDealt,
	// HazardousDealt and ArmorAbsorbed for damage marked and armor spent;
	// CreatureHealed for damage taken off; ArmorLost for armor an effect strips,
	// which also tallies what it took.
	"Cards.Damage":         narratedDirectly,
	"Cards.ArmorRemaining": narratedDirectly,
	"Cards.ArmorStripped":  narratedDirectly,

	// Power counters, which change what a creature's power reads as for as long
	// as it stays in play: PowerCountersPlaced.
	"Cards.PowerCounters": narratedDirectly,

	// The house an in-play card belongs to, which decides whether its controller
	// may use it this turn: CardChangedHouse.
	"Cards.TempHouse":    narratedDirectly,
	"Cards.LastingHouse": narratedDirectly,

	// Stat bonuses: CreatureGainedStats for power and armor,
	// CreatureGainedAssault for assault.
	"Cards.TempPowerBonus":       narratedDirectly,
	"Cards.TempArmorBonus":       narratedDirectly,
	"Cards.TempAssaultBonus":     narratedDirectly,
	"Cards.AssaultUntilNextTurn": narratedDirectly,

	// Æmber on the card: AemberMovedToCard, AemberExalted, AemberCaptured and
	// AemberOnCardReleased.
	"Cards.Amber": narratedDirectly,

	// Runtime type conversion: TurnedIntoCreature and RevertedToArtifact
	// (ADR 0033).
	"Cards.LastingType": narratedDirectly,

	// Attachment and control: UpgradeAttached and UpgradeDiscarded for the
	// upgrade chain, CardPutUnder and CardGrafted for the under chain (whose
	// entry carries FaceDown), ControlTaken and ControlReturned for the cached
	// controller.
	"Cards.FirstUpgradePlus": narratedDirectly,
	"Cards.NextUpgradePlus":  narratedDirectly,
	"Cards.HostPlus":         narratedDirectly,
	"Cards.FirstUnderPlus":   narratedDirectly,
	"Cards.NextUnderPlus":    narratedDirectly,
	"Cards.UnderHostPlus":    narratedDirectly,
	"Cards.UnderFaceDown":    narratedDirectly,
	"Cards.ControlPlus":      narratedDirectly,

	// A card mid-play redirecting itself is silent; the arrival it redirects to
	// is narrated when the play completes (CardPutIntoArchives, CardPurged).
	"Cards.ResolvingDest": narratedByCause,

	// Use bookkeeping, narrated by the use itself: Reaped, Fought and
	// ActionAbilityUsed for TimesUsedThisTurn, Fought for the elusive that the
	// fight spent.
	"Cards.TimesUsedThisTurn":   narratedByCause,
	"Cards.ElusiveUsedThisTurn": narratedByCause,

	// The duration marker on an Animator conversion, not a second conversion:
	// TurnedIntoCreature names the window and RevertedToArtifact its expiry.
	"Cards.CreatureUntilTurnEnd": narratedByCause,

	// The pairing of a gigantic's two halves follows the half entering play and
	// is narrated by that play entry (ADR 0042).
	"Cards.GiganticPartnerPlus": narratedByCause,

	// The house or card a permanent names as it enters play is the reader half of
	// a lock (Restringuntus, Etan's Jar). Naming is silent; the lock narrates when
	// it bites, as a house the opponent cannot choose or a card they cannot play.
	"Cards.NamedHouse":    narratedByCause,
	"Cards.NamedCardPlus": narratedByCause,

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

	// The control stack: ControlTaken and ControlReturned narrate every push and
	// pop, naming the card, the new controller, and the effect that holds it.
	"Controls.Card":       narratedDirectly,
	"Controls.Controller": narratedDirectly,
	"Controls.Source":     narratedDirectly,
	"ControlCount":        narratedDirectly,

	// The generic-counter side-table (ADR 0024): CountersPlaced and
	// CountersRemoved narrate every marker that lands on or leaves a card.
	"Counters.Card": narratedDirectly,
	"Counters.Kind": narratedDirectly,
	"Counters.N":    narratedDirectly,
	"CounterCount":  narratedDirectly,

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
	"OffHousePermitCount":         narratedByCause,
	"NonActivePlaysUsedThisTurn":  narratedByCause,
	"FirstTurnPlayLimit":          narratedByCause,
	"HouseConstraintCount":        narratedByCause,
	"HouseConstraintCountNext":    narratedByCause,
	"FightDamageRedirect":         narratedByCause,
	"FightCancelled":              narratedByCause,
	"FightersPlus":                narratedByCause,

	// The delayed constraint table binding a house choice (ADR 0035).
	// HouseForcedNextTurn, HouseForbiddenNextTurn and HouseWagerArmed narrate an
	// entry being armed; the choice it binds is narrated by HouseChosen when it
	// resolves.
	"HouseConstraints.Kind":          narratedByCause,
	"HouseConstraints.House":         narratedByCause,
	"HouseConstraints.Creature":      narratedByCause,
	"HouseConstraints.Amount":        narratedByCause,
	"HouseConstraints.Predictor":     narratedByCause,
	"HouseConstraints.Source":        narratedByCause,
	"HouseConstraintsNext.Kind":      narratedByCause,
	"HouseConstraintsNext.House":     narratedByCause,
	"HouseConstraintsNext.Creature":  narratedByCause,
	"HouseConstraintsNext.Amount":    narratedByCause,
	"HouseConstraintsNext.Predictor": narratedByCause,
	"HouseConstraintsNext.Source":    narratedByCause,

	// An off-house permit's terms (ADR 0037): MayPlayOrUseGranted narrates the
	// grant, and the play or use that spends it is narrated in its own right.
	"OffHousePermits.Except":     narratedByCause,
	"OffHousePermits.Controlled": narratedByCause,
	"OffHousePermits.Types":      narratedByCause,
	"OffHousePermits.Grant":      narratedByCause,
	"OffHousePermits.Remaining":  narratedByCause,

	// The match RNG is state so a snapshot replays bit-exact (ADR 0039); its
	// advancing is not an outcome anyone can observe.
	"PRNG.State": narratedNever,
}

// narrationTypeName is a struct type's name with any type arguments dropped, so
// Bar[bool] and Bar[House] are the one type Bar that narrationLeafTypes names.
func narrationTypeName(typ reflect.Type) string {
	name := typ.Name()
	if i := strings.IndexByte(name, '['); i >= 0 {
		name = name[:i]
	}
	return name
}

// narrationPaths returns every leaf path of a state struct, descending through
// struct-valued fields (and through arrays of them, since an array narrates the
// same way whatever slot changed) and stopping at a narrationLeafTypes type.
func narrationPaths(typ reflect.Type, prefix string, leavesUsed map[string]bool) []string {
	var paths []string
	for field := range typ.Fields() {
		path := field.Name
		if prefix != "" {
			path = prefix + "." + field.Name
		}
		elem := field.Type
		for elem.Kind() == reflect.Array {
			elem = elem.Elem()
		}
		if elem.Kind() == reflect.Struct {
			name := narrationTypeName(elem)
			if _, leaf := narrationLeafTypes[name]; !leaf {
				paths = append(paths, narrationPaths(elem, path, leavesUsed)...)
				continue
			}
			leavesUsed[name] = true
		}
		paths = append(paths, path)
	}
	return paths
}

// TestEveryStateFieldDeclaresItsNarration is the ratchet behind "every state
// change is logged". It cannot prove the log is complete, but it can stop the
// gap from being introduced silently: a new GameState field — or a new CardCore
// field, which is where the gaps have actually been — fails the build until
// fieldNarration says how the log covers it.
func TestEveryStateFieldDeclaresItsNarration(t *testing.T) {
	leavesUsed := map[string]bool{}
	paths := narrationPaths(reflect.TypeFor[GameState](), "", leavesUsed)
	live := make(map[string]bool, len(paths))
	for _, path := range paths {
		live[path] = true
		if _, ok := fieldNarration[path]; !ok {
			t.Errorf(
				"GameState.%s has no fieldNarration entry; decide whether a change "+
					"to it needs its own log entry", path,
			)
		}
	}
	for path := range fieldNarration {
		if !live[path] {
			t.Errorf("fieldNarration names %q, which GameState no longer has", path)
		}
	}
	for name := range narrationLeafTypes {
		if !leavesUsed[name] {
			t.Errorf(
				"narrationLeafTypes declares %q a leaf, but no GameState field "+
					"holds one", name,
			)
		}
	}
}
