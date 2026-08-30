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
	// zoneCap is the capacity of a single zone (deck is the largest zone).
	zoneCap = 40
	// maxUpgrades is the most upgrades that can be stacked on one creature.
	maxUpgrades = 6
)

// CardCore is the mutable per-match state of a single card, stored purely by
// value. It carries no pointers so the whole GameState copies flat.
type CardCore struct {
	Exhausted bool
	Stunned   bool
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
	UpgradeCount uint8
	Upgrades     [maxUpgrades]LocalID
	// ControlPlus is a temporary control override: 0 means the owner controls the
	// card, otherwise the controller is ControlPlus-1. Ownership never changes.
	ControlPlus uint8
	// ControlSource is the card whose lasting effect took control of this creature
	// "until it leaves play" (Collar of Subordination). When that source leaves
	// play, the control override is reverted. 0 means no such source.
	ControlSource LocalID
}

// CardList is an ordered, fixed-capacity collection of card ids (hand, deck, battle
// line, etc.). It is a value type: copying a CardList copies its contents.
type CardList struct {
	IDs   [zoneCap]LocalID
	Count uint8
}

// slice returns the live ids as a slice header into the underlying array. The
// result must be treated as read-only and not retained across mutations.
func (z *CardList) slice() []LocalID { return z.IDs[:z.Count] }

// add appends an id to the end of the zone.
func (z *CardList) add(id LocalID) {
	z.IDs[z.Count] = id
	z.Count++
}

// addFront inserts an id at the front of the zone (the left flank / top of deck).
func (z *CardList) addFront(id LocalID) {
	copy(z.IDs[1:z.Count+1], z.IDs[:z.Count])
	z.IDs[0] = id
	z.Count++
}

// indexOf returns the position of id, or -1 if absent.
func (z *CardList) indexOf(id LocalID) int {
	for i := 0; i < int(z.Count); i++ {
		if z.IDs[i] == id {
			return i
		}
	}
	return -1
}

// contains reports whether the zone holds id.
func (z *CardList) contains(id LocalID) bool { return z.indexOf(id) >= 0 }

// removeAt removes the id at position i, preserving order, and returns it.
func (z *CardList) removeAt(i int) LocalID {
	id := z.IDs[i]
	copy(z.IDs[i:], z.IDs[i+1:z.Count])
	z.Count--
	z.IDs[z.Count] = 0
	return id
}

// remove deletes id from the zone if present, reporting whether it was found.
func (z *CardList) remove(id LocalID) bool {
	i := z.indexOf(id)
	if i < 0 {
		return false
	}
	z.removeAt(i)
	return true
}

// GameState is the complete mutable state of a match, laid out as a flat value.
// It contains no pointers, slices, or maps, so a copy is a pure value copy with
// no heap allocation or garbage-collector pressure — the property MCTS rollouts
// depend on. Read-only card definitions live in the separate catalog.
type GameState struct {
	Cards      [maxCards]CardCore
	Battleline [2]CardList
	Hand       [2]CardList
	Deck       [2]CardList
	Discard    [2]CardList
	Artifacts  [2]CardList
	Archives   [2]CardList
	// Purge holds cards set aside out of the game ("purged"); they never return.
	Purge [2]CardList

	Aember [2]int
	Keys   [2]int

	// KeyColors[p] holds the colour of each key player p has forged, in forge order;
	// entries [0:Keys[p]] are set, the rest are KeyColorNone. A player picks the
	// colour as they forge (see chooseKeyColor).
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
	CannotFight     [2]bool
	CannotFightNext [2]bool

	// MayFightHouse[p] is a house whose creatures player p may use to fight this
	// turn even when it is not the active house — Brothers in Battle's "each
	// friendly creature of that house may fight." HouseNone (the zero value) grants
	// nothing. EndTurn clears it, so the grant lasts only the turn it was made.
	MayFightHouse [2]House

	// Lasting holds the "for the remainder of the turn" effects active now (Full
	// Moon, Charge!, Crystal Hive reactions; Dimension Door's replacement), fired or
	// queried by game_lasting.go when their event occurs; LastingCount is how many of
	// the fixed array are in use. EndTurn drops a player's entries.
	Lasting      [maxLasting]LastingEffect
	LastingCount uint8

	// CardsPlayedThisTurn[p] counts the cards player p has played this turn, reset by
	// BeginTurn. It backs card-play limits such as Ember Imp's "your opponent cannot
	// play more than 2 cards each turn."
	CardsPlayedThisTurn [2]int
	// CardsPlayedByHouseThisTurn[p][h] counts how many cards of house h player p has
	// played this turn. It resets with CardsPlayedThisTurn and backs conditions such
	// as Epic Quest's "7 or more Sanctum cards this turn."
	CardsPlayedByHouseThisTurn [2][NumHouses]int
	// PlayPermissionsUsedThisTurn[p][h] counts how many off-house play permissions for
	// house h player p has spent this turn (Witch of the Wilds). BeginTurn resets it.
	PlayPermissionsUsedThisTurn [2][NumHouses]int

	// ForcedHouse[p] is the house player p must choose as their active house this
	// turn (Control the Weak); ForcedHouseNext[p] arms that for p's next turn.
	// BeginTurn promotes the armed house to active for the player, so it lands on
	// their own next turn. HouseNone means no house is forced.
	ForcedHouse     [2]House
	ForcedHouseNext [2]House

	// KeyForgeAemberGive[p] means that after player p forges a key this turn, they
	// give their remaining Æmber to KeyForgeAemberGiveBeneficiary[p]. The Next fields
	// arm that delayed transfer for p's next turn (Interdimensional Graft).
	KeyForgeAemberGive                [2]bool
	KeyForgeAemberGiveNext            [2]bool
	KeyForgeAemberGiveBeneficiary     [2]int8
	KeyForgeAemberGiveBeneficiaryNext [2]int8

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

// add registers a definition for an owner and returns its assigned LocalID.
func (c *catalog) add(def *CardDefinition, owner int) LocalID {
	id := LocalID(len(c.defs))
	c.defs = append(c.defs, def)
	c.owners = append(c.owners, uint8(owner))
	return id
}

// def returns the definition for an id.
func (c *catalog) def(id LocalID) *CardDefinition { return c.defs[id] }

// owner returns the owning player index for an id.
func (c *catalog) owner(id LocalID) int { return int(c.owners[id]) }
