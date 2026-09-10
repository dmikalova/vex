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
	// Exhausted is whether the card is turned sideways from being used, so it
	// cannot be used again until it readies.
	Exhausted bool
	// Stunned is whether the creature is stunned: the next time it is used, that use
	// removes the stun instead of reaping, fighting, or firing an Action.
	Stunned bool
	// Enraged is whether the creature is enraged: while set, its controller must use
	// it to fight on their turn if it is able to. Enrage persists across turns until
	// an effect removes it, so nothing in the ready phase clears it.
	Enraged bool
	// Warded is whether the creature has a ward: a one-shot shield that absorbs the
	// next instance of damage or the next time it would leave play, then is spent.
	// Ward persists until it is spent or an effect removes it; the ready phase does
	// not clear it.
	Warded bool
	// DamageImmune, while set, prevents any damage from being dealt to this creature.
	// It lasts until end of turn (the ready phase clears it) — Shield of Justice,
	// Protectrix.
	DamageImmune bool
	// GrantedKeywords is the set of keywords this creature has gained for the
	// remainder of the turn, as a bitmask of Keyword.bit() values (Scout grants
	// Skirmish). The ready phase clears it for every creature.
	GrantedKeywords uint16
	// LostKeywords is the set of keywords this creature has lost for the remainder
	// of the turn, as a bitmask of Keyword.bit() values (Niffle Grounds strips one
	// creature of taunt and elusive). The ready phase clears it for every creature.
	LostKeywords uint16
	// KeywordsUntilNextTurn is the set of keywords this creature has gained until
	// the start of its controller's next turn, as a bitmask of Keyword.bit() values
	// (Hideaway Hole grants elusive). Unlike GrantedKeywords, only the controller's
	// own ready phase clears it, so a defensive keyword survives the opponent's turn.
	KeywordsUntilNextTurn uint16
	// ConsideredFlank, while set, makes this creature count as a flank creature no
	// matter where it sits in its battleline (Spectral Tunneler). It lasts until the
	// remainder of the turn; the ready phase clears it for every creature.
	ConsideredFlank bool
	// ElusiveUsedThisTurn records that this creature has already been chosen to be
	// fought this turn, so its Elusive keyword no longer stops pending fight damage.
	// StartTurn clears it for every creature in play.
	ElusiveUsedThisTurn bool
	// TimesUsedThisTurn counts how many times this creature has been USED this
	// turn — to reap, fight, or use an Action: ability. StartTurn clears every
	// creature's count; leaving play clears it through resetCore.
	TimesUsedThisTurn int16
	// Damage is the damage marked on the creature; it is destroyed once this reaches
	// its power.
	Damage int16
	// ArmorRemaining is the armor left to absorb damage; it refreshes to full when
	// the creature's controller readies.
	ArmorRemaining int16
	// ArmorStripped is how much armor an effect took off this creature, as opposed
	// to how much it spent absorbing damage — the "for each point of armor it lost
	// this way" tally (Red-Hot Armor). The controller's ready phase clears it along
	// with refreshing ArmorRemaining.
	ArmorStripped int16
	// Amber is Æmber sitting on the card (e.g. placed by exalt or capture). It
	// belongs to no player's pool while it stays here.
	Amber int16
	// PowerCounters is the net power from +1/-1 power counters placed on the card;
	// it adds to the creature's power for as long as it stays in play.
	PowerCounters int16
	// TempPowerBonus is power a creature gained for the remainder of the turn
	// (Abond the Armorsmith's Action grants armor the same way). The ready phase
	// clears it for every creature.
	TempPowerBonus int16
	// TempArmorBonus is armor a creature gained for the remainder of the turn —
	// Abond the Armorsmith gives other friendly creatures +1 armor until end of
	// turn. Adding it also tops up ArmorRemaining so the extra armor can absorb
	// damage this turn; the ready phase clears it for every creature.
	TempArmorBonus int16
	// TempAssaultBonus is Assault a creature gained for the remainder of the turn —
	// Creed of Nature grants a chosen creature assault equal to its power. It adds to
	// the creature's Assault value; the ready phase clears it for every creature.
	TempAssaultBonus int16
	// TempHouse is the house this in-play card belongs to until its controller's
	// turn ends. HouseNone means it belongs to its printed house.
	TempHouse House
	// LastingHouse is the house this in-play card belongs to until it leaves play
	// (rather than only until end of turn). HouseNone means none. It is cleared by
	// resetCore when the card leaves play.
	LastingHouse House
	// LastingType is the card type this in-play card has taken on until it leaves
	// play, overriding its printed type — Auto-Legionary is an artifact that turns
	// itself into a creature on the battleline. TypeUnset means it keeps its printed
	// type. It is cleared by resetCore when the card leaves play. Upgrade mode is not
	// stored here: a card in an upgrade chain reads as an Upgrade from its
	// attachment (HostPlus), not this field, because the discard path for a shed
	// upgrade does not reset the core.
	LastingType CardType
	// NamedHouse is a house this card named as it entered play and holds for as long
	// as it stays there, for a HouseLock that constrains that house rather than one
	// printed on the card — Restringuntus bars the house it named. It is the card's
	// choice, not the house the card belongs to.
	NamedHouse House
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
	// ControlPlus caches a card's current controller: 0 means the owner controls
	// the card, otherwise the controller is ControlPlus-1. Ownership never changes.
	// It is the fast read of the top of the card's control stack in the global
	// Controls table, re-derived whenever a control effect is added or removed
	// (see control.go).
	ControlPlus uint8
	// Cards placed under a host form an intrusive singly-linked list threaded
	// through these three bytes, mirroring FirstUpgradePlus/NextUpgradePlus/HostPlus
	// above (see game_under.go) — but unlike an upgrade, a card placed under a host
	// is out of play: Masterplan and Jargogle set a card aside this way rather than
	// leaving it in play. FirstUnderPlus is the head of a host's chain,
	// NextUnderPlus is the next card under the same host, and UnderHostPlus is the
	// back-link to the host.
	FirstUnderPlus uint8
	NextUnderPlus  uint8
	UnderHostPlus  uint8
	// UnderFaceDown records whether this particular card, while placed under a
	// host, is facedown (Masterplan, Jargogle) rather than faceup (Graft) — so only
	// the host's controller may look at it (see Peekable). Named specifically for
	// the Under mechanic rather than a generic "FaceDown" to avoid colliding with
	// the unrelated facedown-in-play token-creature mechanic (Winds of Exchange).
	UnderFaceDown bool
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

// listInsertAt inserts an id at position i (0..count), shifting the ids at and
// after i one slot right — the general placement a Deploy creature uses to enter
// anywhere in its battleline. i == 0 is the left flank, i == count the right.
func listInsertAt(ids []LocalID, count *uint8, i int, id LocalID) {
	copy(ids[i+1:*count+1], ids[i:*count])
	ids[i] = id
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

func (z *wideList) slice() []LocalID { return listSlice(z.IDs[:], z.Count) }
func (z *wideList) add(id LocalID)   { listAdd(z.IDs[:], &z.Count, id) }
func (z *wideList) insertAt(i int, id LocalID) {
	listInsertAt(z.IDs[:], &z.Count, i, id)
}
func (z *wideList) indexOf(id LocalID) int   { return listIndexOf(z.IDs[:], z.Count, id) }
func (z *wideList) contains(id LocalID) bool { return z.indexOf(id) >= 0 }
func (z *wideList) remove(id LocalID) bool   { return listRemove(z.IDs[:], &z.Count, id) }

// GameState is the complete mutable state of a match, laid out as a flat value.
// It contains no pointers, slices, or maps, so a copy is a pure value copy with
// no heap allocation or garbage-collector pressure — the property MCTS rollouts
// depend on. Read-only card definitions live in the separate catalog.
type GameState struct {
	// Cards is every card's mutable per-match state, indexed by LocalID.
	Cards [maxCards]CardCore
	// Battleline[p] is player p's row of creatures in play, in flank-to-flank order.
	Battleline [2]wideList
	// Hand[p] is player p's hand.
	Hand [2]deckList
	// Deck[p] is player p's draw pile.
	Deck [2]deckList
	// Discard[p] is player p's discard pile.
	Discard [2]deckList
	// Artifacts[p] is player p's row of artifacts in play.
	Artifacts [2]wideList
	// Archives[p] is player p's archives — cards set aside to be taken into hand.
	Archives [2]wideList
	// Purge holds cards set aside out of the game ("purged"); they never return.
	Purge [2]deckList

	// Aember[p] is the Æmber in player p's pool.
	Aember [2]int
	// Keys[p] is how many keys player p has forged.
	Keys [2]int

	// KeyColors[p] holds the colour of each key player p has forged, in forge order;
	// entries [0:Keys[p]] are set, the rest are KeyColorNone. A player picks the
	// colour as they forge (see pickKeyColor).
	KeyColors [2][KeysToWin]KeyColor

	// Chains[p] is player p's chain count. Chains penalize a player by reducing how
	// many cards they draw at the end of their turn — one fewer card for every 6
	// chains — and a player sheds a single chain on a turn where that reduction
	// actually blocked a draw (see Game.drawStep).
	Chains [2]int

	// ActivePlayer is the index of the player whose turn it is.
	ActivePlayer int
	// ActiveHouse is the house the active player chose to act with this turn.
	ActiveHouse House
	// Tide is the game-wide tide state: neutral at the start of the game, then high
	// for the player who raised it and low for their opponent. No implemented card
	// raises the tide yet (the mechanic arrives in a later set), so the tide stays
	// neutral and any "while the tide is low/high" check reads false.
	Tide Tide
	// Turn is the current turn number.
	Turn   int
	Winner int // -1 while the game is ongoing

	// Phase is the part of the turn now running (ADR 0012). The engine advances it
	// and blocks only on the phases that need a frontend decision.
	Phase Phase
	// PhaseEnded is the "skip the rest of this phase" flag the phase loop reads, so
	// an effect that cuts a phase short (Omega ending the play phase) does not have
	// to special-case any one phase. Entering a phase clears it.
	PhaseEnded bool

	// Fight bars. CannotFight[p] blocks player p from using creatures to fight on
	// the current turn; CannotFightNext[p] arms that block for p's next turn. An
	// effect (Fogbank) arms the bar, StartTurn promotes it to active for the
	// affected player, and the ready phase lifts it — so it always lands on that
	// player's own next turn, whoever plays in between.
	CannotFight     [2]Bar[bool]
	CannotFightNext [2]Bar[bool]

	// Play-type bars. CannotPlayTypeThis[p] blocks player p from playing cards of
	// that type this turn; CannotPlayTypeNext[p] arms that block for p's next turn
	// (Lifeward bars creatures, Scrambler Storm bars action cards). The zero value
	// (an unset CardType) bars nothing. Like the fight bar, an effect arms it and
	// StartTurn promotes it to that player's own next turn.
	CannotPlayTypeThis [2]Bar[CardType]
	CannotPlayTypeNext [2]Bar[CardType]

	// Use bars. CannotUse[p] blocks player p from using any card this turn — reaping,
	// fighting, or an "Action:" ability (Skippy Timehog); playing and discarding are
	// untouched. CannotUseNext[p] arms that block for p's next turn, and like the
	// fight bar StartTurn promotes it and the ready phase lifts it.
	CannotUse     [2]Bar[bool]
	CannotUseNext [2]Bar[bool]

	// SideDamageImmune[p], while set, makes every creature player p controls immune
	// to damage for the current turn — Shield of Justice protects each friendly
	// creature for the remainder of the turn. It is a side-wide mask read live at
	// damage time (not a snapshot of the creatures in play when it resolved), so a
	// creature played or gained after it resolves is protected too, and one taken by
	// the opponent is not. The ready phase clears it, like the turn bars.
	SideDamageImmune [2]bool

	// Reap bars. CannotReap[p] stops player p reaping with any creature this turn
	// (Inky Gloom); CannotReapNext[p] arms that block for p's next turn. Like the
	// use bar, StartTurn promotes it and the ready phase lifts it. This is narrower
	// than CannotUse — fighting and "Action:" abilities stay open.
	CannotReap     [2]Bar[bool]
	CannotReapNext [2]Bar[bool]

	// Reap-by-house bars. CannotReapHouse[p] stops player p reaping with creatures
	// of the named house this turn (Seismo-entangler); CannotReapHouseNext[p] arms
	// that block for p's next turn. HouseNone (the zero value) bars nothing.
	// StartTurn promotes the armed house, so it lands on that player's own turn.
	CannotReapHouse     [2]Bar[House]
	CannotReapHouseNext [2]Bar[House]

	// Board-wide creature bars. CreaturesCannot[p] stops player p using their
	// creatures one way — fighting or reaping — this turn, save for the house the
	// bar spares (Into the Night, Sow Salt); CreaturesCannotNext[p] arms that block
	// for p's next turn. Unlike CannotFight/CannotReap, which are one player's own
	// choice, these are a rule on every creature in play, so the caster arms both
	// entries: their own for this turn and the opponent's for next turn, and it
	// lifts at the start of the caster's next turn. The zero value bars nothing.
	CreaturesCannot     [2]Bar[CreatureBar]
	CreaturesCannotNext [2]Bar[CreatureBar]

	// SkipForge bars. SkipForgeNext[p] makes player p skip their "forge a key" phase
	// at the start of their next turn (Miasma); StartTurn promotes it to SkipForge[p]
	// and forges accordingly, so it lands on that player's own next turn.
	SkipForge     [2]Bar[bool]
	SkipForgeNext [2]Bar[bool]

	// EndOfTurnDestroyAll, while set, schedules "destroy each creature" to resolve
	// in the active player's end-of-turn phase (Ragnarok). It is a one-shot armed
	// during the play phase; unlike the turn bars it must survive the ready phase
	// (which runs before end of turn), so it is cleared only when it fires. Its
	// Source is the card that armed it, for attribution.
	EndOfTurnDestroyAll Bar[bool]

	// Key surcharges. KeyCostBump[p] raises player p's key cost for the current turn;
	// KeyCostBumpNext[p] arms that raise for p's next turn (Lash of Broken Dreams
	// makes keys cost +3 during the opponent's next turn). Unlike a card's
	// KeyCostChange, which lives as long as the card is in play, this is a one-turn
	// surcharge, promoted by StartTurn and lifted by the ready phase like the other
	// bars.
	KeyCostBump     [2]Bar[int]
	KeyCostBumpNext [2]Bar[int]

	// MayFightHouse[p] is a house whose creatures player p may use to fight this
	// turn even when it is not the active house — Brothers in Battle's "each
	// friendly creature of that house may fight." HouseNone (the zero value) grants
	// nothing. The ready phase clears it, so the grant lasts only the turn it was
	// made.
	MayFightHouse [2]House

	// MayFightAny[p] lets every creature player p controls fight this turn whatever
	// its house — Follow the Leader, Horseman of War, the unrestricted form of the
	// MayFightHouse grant. The ready phase clears it.
	MayFightAny [2]bool

	// MayUseHouse[p] is a house whose creatures player p may fully use this turn
	// (fight, reap, or Action:) even when it is not the active house — Sigil of
	// Brotherhood, Ritual of the Hunt. HouseNone grants nothing; the ready phase
	// clears it.
	MayUseHouse [2]House

	// MayPlayHouse[p] is a house whose cards player p may play from hand this turn
	// even when it is not the active house — the Ambassador cycle's "you may play
	// or use a <House> card this turn". HouseNone grants nothing; the ready phase
	// clears it.
	MayPlayHouse [2]House

	// MayUseArtifactsAnyHouse[p] lets player p use any friendly artifact as if it
	// belonged to the active house for the remainder of the turn — Scientifical
	// Hack. The ready phase clears it.
	MayUseArtifactsAnyHouse [2]bool

	// TurnHistory holds the small tallies of what each player did during a turn —
	// several cards ask that rather than what is on the board ("if your opponent
	// forged a key on their previous turn", "for each enemy creature destroyed in a
	// fight this turn"). Keeping them as one array indexed by TurnStat leaves the
	// state flat and comparable, and makes a new tally one more enum value rather
	// than another pair of fields. The ready phase rolls each "this turn" tally
	// into its "last turn" twin, so "their previous turn" always means that
	// player's own last completed turn.
	TurnHistory [2][turnStatCount]int8

	// KeywordsLost is the set of keywords every creature in play has lost for the
	// remainder of the turn — Sniffer takes elusive away from each creature. It is a
	// bitmask over keywordBit so the state stays flat and comparable; the ready
	// phase clears it.
	KeywordsLost uint16

	// TextBlank[p] blanks the text box of every creature player p controls — its
	// printed keywords, abilities, and constant grants are ignored (its traits and
	// stats remain). Shadow of Dis blanks the opponent's creatures "until your next
	// turn": set on the affected player, it persists through that player's own next
	// turn and is lifted by their ready phase, so it always spans exactly through
	// the opponent's turn whoever plays in between.
	TextBlank [2]Bar[bool]

	// Lasting holds the "for the remainder of the turn" effects active now (Full
	// Moon, Charge!, Crystal Hive reactions; Dimension Door's replacement), fired or
	// queried by game_lasting.go when their event occurs; LastingCount is how many of
	// the fixed array are in use. The ready phase drops a player's entries.
	Lasting      [maxLasting]LastingEffect
	LastingCount uint8

	// Reanimations holds the delayed put-into-play records armed by a creature's
	// Destroyed ability (Gebuk), each fired by fireReanimations once its source has
	// left play; ReanimationsCount is how many of the fixed array are in use.
	Reanimations      [maxReanimations]ReanimateInPlace
	ReanimationsCount uint8

	// Morphs holds the "for the remainder of the turn" trigger morphs active now
	// (Livia the Elder's fight/reap fuse), queried by game_abilities.go when a
	// creature's abilities are gathered; MorphCount is how many of the fixed array
	// are in use. The ready phase drops a player's entries.
	Morphs     [maxMorph]LastingMorph
	MorphCount uint8

	// PlayedThisTurn[p] and DiscardedThisTurn[p] record, in order, the cards player p
	// has played and discarded this turn; StartTurn clears both. Cards filter them
	// themselves — by house for Epic Quest's "7 or more Sanctum cards this turn" and
	// Giant Sloth — and the play log's length is the limit Ember Imp reads. Only hand
	// discards are logged: rule "a player discarding a card" means from their hand.
	PlayedThisTurn    [2]turnLog
	DiscardedThisTurn [2]turnLog
	// PlayPermissionsUsedThisTurn[p][h] counts how many off-house play permissions for
	// house h player p has spent this turn (Witch of the Wilds). StartTurn resets it.
	// A turn cannot spend more than a hand's worth, so a byte per house is ample.
	PlayPermissionsUsedThisTurn [2][NumHouses]uint8

	// OffHousePermits[p] are the this-turn grants letting player p play or use a
	// bounded number of cards outside their active house (Com. Officer Kirby, CXO
	// Taber, United Action); OffHousePermitCount is how many of the fixed array are
	// in use. The ready phase clears them.
	OffHousePermits     [2][maxOffHousePermits]OffHousePermit
	OffHousePermitCount [2]uint8

	// NonActivePlaysUsedThisTurn[p] counts the off-house plays player p has made this
	// turn using a continuous "any non-active house" permission (Captain Val
	// Jericho). StartTurn resets it.
	NonActivePlaysUsedThisTurn [2]uint8

	// FirstTurnPlayLimit[p] holds the first-turn rule for player p: on the first
	// player's first turn they may play or discard only one card from hand. StartGame
	// arms it for the first player once the opening turn has begun, and StartTurn
	// clears it at the start of every turn, so it lasts only that one turn. Cards a
	// played card lets its controller play (Wild Wormhole, Phase Shift) bypass the
	// limit because they never pass through the volitional play/discard gates.
	FirstTurnPlayLimit [2]bool

	// ForcedHouse[p] is the house player p must choose as their active house this
	// turn (Control the Weak); ForcedHouseNext[p] arms that for p's next turn.
	// StartTurn promotes the armed house to active for the player, so it lands on
	// their own next turn. HouseNone means no house is forced.
	ForcedHouse     [2]Bar[House]
	ForcedHouseNext [2]Bar[House]

	// ForbiddenHouse[p] is a house player p cannot choose as their active house
	// this turn (Tezmal); ForbiddenHouseNext[p] arms that for p's next turn.
	// StartTurn promotes the armed house, so it lands on the player's own next
	// turn. HouseNone means no house is forbidden.
	ForbiddenHouse     [2]Bar[House]
	ForbiddenHouseNext [2]Bar[House]

	// HouseWager[p] is a bet on player p's active house this turn: if p chooses
	// its House, the Predictor steals Amount (Snaglet). HouseWagerNext[p] arms
	// that for p's next turn, and StartTurn promotes it so the payoff lands when p
	// next chooses a house. A zero Amount means no wager is armed.
	HouseWager     [2]HouseWager
	HouseWagerNext [2]HouseWager

	// FightDamageRedirect is the creature a "Before Fight" ability chose to receive
	// the attacker's fight damage instead of the defender (Gabos Longarms). It is
	// set during the fight in progress and read and cleared by the combat step; 0
	// means the attacker's fight damage hits the creature it is fighting as usual.
	FightDamageRedirect LocalID
	// FightCancelled means a "Before Fight" ability made the fight not occur
	// (Evasion Sigil). It is set during the fight in progress and read and cleared
	// before Assault, Hazardous, fight damage, and Fight: abilities would resolve.
	FightCancelled bool
	// FightersPlus holds the two creatures resolving the fight in progress — the
	// attacker and the creature it fights — each +1 encoded so card 0 is
	// distinguishable from "no fight" (the zero value). A creature counts as
	// "fighting" (Nizak, The Forgotten gains invulnerable while fighting) only while
	// it is one of these. It is set for the span of one fight and restored after, so
	// a nested fight reads only its own pair, and is the zero value outside combat.
	FightersPlus [2]LocalID
	// PurgePlayedAction is the action card whose own "Play:" ability purges it
	// (Library Access): it is set while that ability resolves and read when the
	// played action would go to the discard pile, sending it to the purge pile
	// instead. PurgePlayedActionSet distinguishes "purge card 0" from the unset
	// zero value, since LocalID 0 is a valid card.
	PurgePlayedAction    LocalID
	PurgePlayedActionSet bool
	// ArchivePlayedAction is the action card whose own "Play:" ability archives it
	// (Sucker Punch): it is set while that ability resolves and read when the
	// played action would go to the discard pile, sending it to the archives
	// instead. ArchivePlayedActionSet distinguishes "archive card 0" from the
	// unset zero value, since LocalID 0 is a valid card.
	ArchivePlayedAction    LocalID
	ArchivePlayedActionSet bool
	// Counters is the global side-table of generic counters — the card-placed
	// markers (doom, and its kin) whose meaning is defined entirely by the card
	// that reads them. One entry per (card, kind) pair, its count folded into N, so
	// a creature piled high with one kind is still a single entry (ADR 0024). The
	// table lives here rather than in CardCore so an unused kind costs nothing per
	// snapshot; CounterCount is how many of the entries are live. Compacted on
	// removal and shed when a card leaves play.
	Counters     [maxCounterEntries]CounterEntry
	CounterCount uint8
	// Controls is the global stack of "take control" effects, ordered by when each
	// was applied. A card's current controller is the most recently pushed entry
	// still in effect; removing one falls back to the entry beneath it (LIFO). The
	// table lives here rather than in CardCore so a card that is never seized costs
	// nothing per snapshot; ControlCount is how many entries are live. Compacted on
	// removal (see control.go).
	Controls     [maxControlEntries]ControlEntry
	ControlCount uint8
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
