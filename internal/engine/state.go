package engine

// LocalID identifies one physical card within a single match. During setup every
// card (including duplicates) is registered in the Catalog and assigned a stable
// LocalID; the flat GameState references cards only by this id.
type LocalID uint8

// Capacities for the fixed-size, pointerless state arrays. They are generous
// upper bounds — a real match uses far fewer.
const (
	// maxCards is the LocalID space for one match.
	maxCards = 128
	// deckCap bounds a zone that can only ever hold a player's own deck: deck, hand,
	// discard, and purge only receive their owner's cards, and a KeyForge deck is 36.
	deckCap = 36
	// wideCap bounds a zone that can hold cards from both decks at once: a player can
	// control an opponent's creatures and artifacts (battle line, artifact row) and
	// archive cards from either deck (archives), so up to both 36-card decks combined.
	wideCap = 72
	// turnLogCap bounds one turn's play or discard log. A card can be played more
	// than once in a turn (returned to hand and replayed), so this is not bounded by
	// the deck; it sits well past any real turn and the log saturates there.
	turnLogCap = 64
)

// CardCore is the mutable per-match state of a single card, stored purely by
// value. It carries no pointers so the whole GameState copies flat.
type CardCore struct {
	Exhausted bool
	Stunned   bool
	// DamageImmune, while set, prevents any damage from being dealt to this creature.
	// It lasts until end of turn (EndTurn clears it) — Shield of Justice, Protectrix.
	DamageImmune bool
	// ElusiveUsedThisTurn records that this creature has already been chosen to be
	// fought this turn, so its Elusive keyword no longer stops pending fight damage.
	// BeginTurn clears it for every creature in play.
	ElusiveUsedThisTurn bool
	// TimesUsedThisTurn counts how many times this creature has been USED this
	// turn — to reap, fight, or use an Action: ability. BeginTurn clears every
	// creature's count; leaving play clears it through resetCore.
	TimesUsedThisTurn int16
	Damage            int16
	ArmorRemaining    int16
	// Amber is Æmber sitting on the card (e.g. placed by exalt or capture). It
	// belongs to no player's pool while it stays here.
	Amber int16
	// PowerCounters is the net power from +1/-1 power counters placed on the card;
	// it adds to the creature's power for as long as it stays in play.
	PowerCounters int16
	// TempHouse is the house this in-play card belongs to until its controller's
	// turn ends. HouseNone means it belongs to its printed house.
	TempHouse House
	// LastingHouse is the house this in-play card belongs to until it leaves play
	// (rather than only until end of turn). HouseNone means none. It is cleared by
	// resetCore when the card leaves play.
	LastingHouse House
	// Upgrades attached to a creature form an intrusive singly-linked list threaded
	// through these three bytes, so a creature carries any number of upgrades with no
	// per-card fixed array — KeyForge sets no limit on how many upgrades a creature
	// may hold. All three use +1 encoding (0 means "none") like ControlPlus, so the
	// zero value is cleanly "unattached": FirstUpgradePlus is the head of a host's
	// chain, NextUpgradePlus is the next upgrade on the same host, and HostPlus is the
	// creature an upgrade is attached to — a back-link that makes detaching and
	// "which creature am I on?" O(1).
	FirstUpgradePlus uint8
	NextUpgradePlus  uint8
	HostPlus         uint8
	// ControlPlus is a temporary control override: 0 means the owner controls the
	// card, otherwise the controller is ControlPlus-1. Ownership never changes.
	ControlPlus uint8
	// ControlSource is the card whose lasting effect took control of this creature
	// "until it leaves play" (Collar of Subordination). When that source leaves
	// play, the control override is reverted. 0 means no such source.
	ControlSource LocalID
}

// A zone is an ordered, fixed-capacity collection of card ids (hand, deck, battle
// line, etc.), stored purely by value so copying it copies its contents. There are
// two sizes — deckList for zones bounded by a single deck, wideList for zones that
// can hold cards from both decks — because Go generics cannot slice a type parameter
// whose array sizes differ, so the shared logic lives in the list* free functions
// below and each type is a thin wrapper over them.
type deckList struct {
	IDs   [deckCap]LocalID
	Count uint8
}

type wideList struct {
	IDs   [wideCap]LocalID
	Count uint8
}

// A turnLog records, in order, the cards a player played or discarded this turn.
// Cards filter it themselves (by house, trait, type, …) rather than the engine
// keeping a separate tally per axis. Unlike a zone the same card can appear more
// than once, so it has its own cap and saturates instead of overflowing.
type turnLog struct {
	IDs   [turnLogCap]LocalID
	Count uint8
}

func (l *turnLog) slice() []LocalID { return listSlice(l.IDs[:], l.Count) }
func (l *turnLog) reset()           { *l = turnLog{} }

// add appends an id, dropping it once the log is full — a turn past turnLogCap
// cards is beyond anything a card counts.
func (l *turnLog) add(id LocalID) {
	if int(l.Count) == turnLogCap {
		return
	}
	listAdd(l.IDs[:], &l.Count, id)
}

// listSlice returns the live ids as a read-only slice header into the backing array.
func listSlice(ids []LocalID, count uint8) []LocalID { return ids[:count] }

// listAdd appends an id at the end of the zone.
func listAdd(ids []LocalID, count *uint8, id LocalID) {
	ids[*count] = id
	*count++
}

// listAddFront inserts an id at the front (the left flank / top of deck).
func listAddFront(ids []LocalID, count *uint8, id LocalID) {
	copy(ids[1:*count+1], ids[:*count])
	ids[0] = id
	*count++
}

// listIndexOf returns the position of id, or -1 if absent.
func listIndexOf(ids []LocalID, count uint8, id LocalID) int {
	for i := 0; i < int(count); i++ {
		if ids[i] == id {
			return i
		}
	}
	return -1
}

// listRemoveAt removes the id at position i, preserving order, and returns it.
func listRemoveAt(ids []LocalID, count *uint8, i int) LocalID {
	id := ids[i]
	copy(ids[i:], ids[i+1:*count])
	*count--
	ids[*count] = 0
	return id
}

// listRemove deletes id if present, reporting whether it was found.
func listRemove(ids []LocalID, count *uint8, id LocalID) bool {
	i := listIndexOf(ids, *count, id)
	if i < 0 {
		return false
	}
	listRemoveAt(ids, count, i)
	return true
}

func (z *deckList) slice() []LocalID         { return listSlice(z.IDs[:], z.Count) }
func (z *deckList) add(id LocalID)           { listAdd(z.IDs[:], &z.Count, id) }
func (z *deckList) addFront(id LocalID)      { listAddFront(z.IDs[:], &z.Count, id) }
func (z *deckList) indexOf(id LocalID) int   { return listIndexOf(z.IDs[:], z.Count, id) }
func (z *deckList) contains(id LocalID) bool { return z.indexOf(id) >= 0 }
func (z *deckList) removeAt(i int) LocalID   { return listRemoveAt(z.IDs[:], &z.Count, i) }
func (z *deckList) remove(id LocalID) bool   { return listRemove(z.IDs[:], &z.Count, id) }

func (z *wideList) slice() []LocalID         { return listSlice(z.IDs[:], z.Count) }
func (z *wideList) add(id LocalID)           { listAdd(z.IDs[:], &z.Count, id) }
func (z *wideList) addFront(id LocalID)      { listAddFront(z.IDs[:], &z.Count, id) }
func (z *wideList) indexOf(id LocalID) int   { return listIndexOf(z.IDs[:], z.Count, id) }
func (z *wideList) contains(id LocalID) bool { return z.indexOf(id) >= 0 }
func (z *wideList) remove(id LocalID) bool   { return listRemove(z.IDs[:], &z.Count, id) }

// GameState is the complete mutable state of a match, laid out as a flat value.
// It contains no pointers, slices, or maps, so a copy is a pure value copy with
// no heap allocation or garbage-collector pressure — the property MCTS rollouts
// depend on. Read-only card definitions live in the separate catalog.
type GameState struct {
	Cards      [maxCards]CardCore
	Battleline [2]wideList
	Hand       [2]deckList
	Deck       [2]deckList
	Discard    [2]deckList
	Artifacts  [2]wideList
	Archives   [2]wideList
	// Purge holds cards set aside out of the game ("purged"); they never return.
	Purge [2]deckList

	Aember [2]int
	Keys   [2]int

	// KeyColors[p] holds the colour of each key player p has forged, in forge order;
	// entries [0:Keys[p]] are set, the rest are KeyColorNone. A player picks the
	// colour as they forge (see pickKeyColor).
	KeyColors [2][KeysToWin]KeyColor

	// Chains[p] is player p's chain count. Chains penalize a player by reducing how
	// many cards they draw at the end of their turn — one fewer card for every 6
	// chains — and a player sheds a single chain on a turn where that reduction
	// actually blocked a draw (see Game.drawStep).
	Chains [2]int

	ActivePlayer int
	ActiveHouse  House
	Turn         int
	Winner       int // -1 while the game is ongoing

	// Fight bars. CannotFight[p] blocks player p from using creatures to fight on
	// the current turn; CannotFightNext[p] arms that block for p's next turn. An
	// effect (Fogbank) arms the bar, BeginTurn promotes it to active for the
	// affected player, and EndTurn lifts it — so it always lands on that player's
	// own next turn, whoever plays in between.
	CannotFight     [2]Bar[bool]
	CannotFightNext [2]Bar[bool]

	// Play-type bars. CannotPlayTypeThis[p] blocks player p from playing cards of
	// that type this turn; CannotPlayTypeNext[p] arms that block for p's next turn
	// (Lifeward bars creatures, Scrambler Storm bars action cards). The zero value
	// (an unset CardType) bars nothing. Like the fight bar, an effect arms it and
	// BeginTurn promotes it to that player's own next turn.
	CannotPlayTypeThis [2]Bar[CardType]
	CannotPlayTypeNext [2]Bar[CardType]

	// Use bars. CannotUse[p] blocks player p from using any card this turn — reaping,
	// fighting, or an "Action:" ability (Skippy Timehog); playing and discarding are
	// untouched. CannotUseNext[p] arms that block for p's next turn, and like the
	// fight bar BeginTurn promotes it and EndTurn lifts it.
	CannotUse     [2]Bar[bool]
	CannotUseNext [2]Bar[bool]

	// SkipForge bars. SkipForgeNext[p] makes player p skip their "forge a key" step
	// at the start of their next turn (Miasma); BeginTurn promotes it to SkipForge[p]
	// and forges accordingly, so it lands on that player's own next turn.
	SkipForge     [2]Bar[bool]
	SkipForgeNext [2]Bar[bool]

	// MayFightHouse[p] is a house whose creatures player p may use to fight this
	// turn even when it is not the active house — Brothers in Battle's "each
	// friendly creature of that house may fight." HouseNone (the zero value) grants
	// nothing. EndTurn clears it, so the grant lasts only the turn it was made.
	MayFightHouse [2]House

	// MayFightAny[p] lets every creature player p controls fight this turn whatever
	// its house — Follow the Leader, Horseman of War, the unrestricted form of the
	// MayFightHouse grant. EndTurn clears it.
	MayFightAny [2]bool

	// MayUseHouse[p] is a house whose creatures player p may fully use this turn
	// (fight, reap, or Action:) even when it is not the active house — Sigil of
	// Brotherhood, Ritual of the Hunt. HouseNone grants nothing; EndTurn clears it.
	MayUseHouse [2]House

	// TurnHistory holds the small tallies of what each player did during a turn —
	// several cards ask that rather than what is on the board ("if your opponent
	// forged a key on their previous turn", "for each enemy creature destroyed in a
	// fight this turn"). Keeping them as one array indexed by TurnStat leaves the
	// state flat and comparable, and makes a new tally one more enum value rather
	// than another pair of fields. EndTurn rolls each "this turn" tally into its
	// "last turn" twin, so "their previous turn" always means that player's own last
	// completed turn.
	TurnHistory [2][turnStatCount]int8

	// KeywordsLost is the set of keywords every creature in play has lost for the
	// remainder of the turn — Sniffer takes elusive away from each creature. It is a
	// bitmask over keywordBit so the state stays flat and comparable; EndTurn clears
	// it.
	KeywordsLost uint8

	// Lasting holds the "for the remainder of the turn" effects active now (Full
	// Moon, Charge!, Crystal Hive reactions; Dimension Door's replacement), fired or
	// queried by game_lasting.go when their event occurs; LastingCount is how many of
	// the fixed array are in use. EndTurn drops a player's entries.
	Lasting      [maxLasting]LastingEffect
	LastingCount uint8

	// PlayedThisTurn[p] and DiscardedThisTurn[p] record, in order, the cards player p
	// has played and discarded this turn; BeginTurn clears both. Cards filter them
	// themselves — by house for Epic Quest's "7 or more Sanctum cards this turn" and
	// Giant Sloth — and the play log's length is the limit Ember Imp reads. Only hand
	// discards are logged: rule "a player discarding a card" means from their hand.
	PlayedThisTurn    [2]turnLog
	DiscardedThisTurn [2]turnLog
	// PlayPermissionsUsedThisTurn[p][h] counts how many off-house play permissions for
	// house h player p has spent this turn (Witch of the Wilds). BeginTurn resets it.
	PlayPermissionsUsedThisTurn [2][NumHouses]int

	// ForcedHouse[p] is the house player p must choose as their active house this
	// turn (Control the Weak); ForcedHouseNext[p] arms that for p's next turn.
	// BeginTurn promotes the armed house to active for the player, so it lands on
	// their own next turn. HouseNone means no house is forced.
	ForcedHouse     [2]Bar[House]
	ForcedHouseNext [2]Bar[House]

	// FightDamageRedirect is the creature a "Before Fight" ability chose to receive
	// the attacker's fight damage instead of the defender (Gabos Longarms). It is
	// set during the fight in progress and read and cleared by the combat step; 0
	// means the attacker's fight damage hits the creature it is fighting as usual.
	FightDamageRedirect LocalID
	// FightCancelled means a "Before Fight" ability made the fight not occur
	// (Evasion Sigil). It is set during the fight in progress and read and cleared
	// before Assault, Hazardous, fight damage, and Fight: abilities would resolve.
	FightCancelled bool
}

// FastCopy returns an independent copy of the state. Because every field is a
// value type this is a single flat copy; mutating the result never affects the
// original.
func (s GameState) FastCopy() GameState { return s }

// catalog is the read-only registry of card definitions for a match. It is held
// separately from GameState (by pointer) and never mutated during play, so it is
// shared freely across cloned states.
type catalog struct {
	defs   []*CardDefinition
	owners []uint8
}

// add registers a definition for an owner and returns its assigned LocalID. It
// panics if the match exceeds maxCards, turning a silent LocalID overflow (and the
// cryptic out-of-range access into GameState.Cards that follows) into a clear
// diagnostic at setup.
func (c *catalog) add(def *CardDefinition, owner int) LocalID {
	if len(c.defs) >= maxCards {
		panic("engine: too many cards registered for one match (maxCards exceeded)")
	}
	id := LocalID(len(c.defs))
	c.defs = append(c.defs, def)
	c.owners = append(c.owners, uint8(owner))
	return id
}

// hasRoom reports whether another card can still be registered in this match.
func (c *catalog) hasRoom() bool { return len(c.defs) < maxCards }

// def returns the definition for an id.
func (c *catalog) def(id LocalID) *CardDefinition { return c.defs[id] }

// owner returns the owning player index for an id.
func (c *catalog) owner(id LocalID) int { return int(c.owners[id]) }
