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
	// ResolvingDest is where a card mid-play goes when its play completes, when its
	// own ability redirected it there (Sucker Punch archives itself, Library Access
	// purges itself). The unset zero value sends it to its owner's discard pile.
	//
	// A card mid-play has left the zone it was played from and has not yet reached a
	// destination, so it is in no zone at all. That is why the redirect names no
	// source zone and works whatever the card was played from: Wild Wormhole plays
	// Causal Loop off the deck and Causal Loop still archives itself. It is per-card
	// rather than one field on GameState because a resolving Tactic can play another
	// card while it resolves, and each needs its own destination
	// (TestResolvingCardRedirectIsPerCard).
	ResolvingDest Destination
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
	// own start-of-turn phase clears it, before any start-of-turn ability resolves, so
	// a defensive keyword survives the opponent's turn.
	KeywordsUntilNextTurn uint16
	// LostKeywordsUntilNextTurn is the set of keywords this creature has lost until
	// the start of its controller's next turn, as a bitmask of Keyword.bit() values
	// (Reckless Rizzo loses elusive after stealing). Unlike LostKeywords, only the
	// controller's own start-of-turn phase clears it, before any start-of-turn ability
	// resolves, so the loss survives the opponent's turn.
	LostKeywordsUntilNextTurn uint16
	// TraitUntilNextTurn is a trait this creature has gained until the start of its
	// controller's next turn — the Mutation cycle grants a chosen creature the Mutant
	// trait. Only the controller's own start-of-turn phase clears it, before any
	// start-of-turn ability resolves, so the trait survives the opponent's turn. It
	// holds a single trait (the most recently granted); the
	// only cards that grant it, the three Mutations, all grant Mutant, so overwrite
	// never loses a distinct trait in practice. traitUnset means none.
	TraitUntilNextTurn Trait
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
	// AssaultUntilNextTurn is Assault a creature gained until the start of its
	// controller's next turn — the Mutation cycle grants a chosen creature assault 3.
	// Unlike TempAssaultBonus, only the controller's own start-of-turn phase clears it,
	// before any start-of-turn ability resolves, so it matches the keyword and trait a
	// mutation grants in the same breath. It adds to the creature's Assault value.
	AssaultUntilNextTurn int16
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
	// CreatureUntilTurnEnd marks a card whose LastingType conversion to a creature
	// lasts only the current turn (Animator), as opposed to the permanent conversion
	// Auto-Legionary makes. The ready phase reverts every card carrying it to an
	// artifact in its controller's row, so a turn-scoped conversion lifts at end of
	// turn like every other RemainderOfPlayerTurn effect. resetCore clears it when the
	// card leaves play.
	CreatureUntilTurnEnd bool
	// TextBoxSourcePlus records that this creature has gained the printed text box
	// (traits, keywords, and triggered abilities) of another card until it leaves
	// play — Mimic Gel copies a chosen creature. It stores that source's LocalID+1
	// (0 means "no gained text box") so the zero value is cleanly "none", like
	// HostPlus. resetCore clears it when the card leaves play. The source's def is
	// read from the immutable catalog, so it stays available even after the source
	// itself leaves play.
	TextBoxSourcePlus uint8
	// TextBoxTurnSourcePlus records that this creature has gained the printed text
	// box of another card for the remainder of the turn — Creed of Nurture lends a
	// hand creature's text box to a creature in play. It stores that source's
	// LocalID+1 like TextBoxSourcePlus; the ready phase clears it for every creature.
	TextBoxTurnSourcePlus uint8
	// CopiedStatsSourcePlus records that this creature copies another card's printed
	// stats until it leaves play: its power becomes that card's printed power, and it
	// gains that card's printed armor, keywords, and traits — Cyber-Clone copies a
	// creature it purges. It stores that source's LocalID+1 (0 means "no copy") like
	// TextBoxSourcePlus, and resetCore clears it when this card leaves play. The
	// source's printed stats are read from the immutable catalog, so they stay
	// available even after the source card leaves play.
	CopiedStatsSourcePlus uint8
	// NamedHouse is a house this card named as it entered play and holds for as long
	// as it stays there, for a HouseLock that constrains that house rather than one
	// printed on the card — Restringuntus bars the house it named. It is the card's
	// choice, not the house the card belongs to.
	NamedHouse House
	// NamedCardPlus records a card this permanent named as it entered play, matched
	// by name against every card attempted to be played while this permanent stays
	// in play — Etan's Jar bars the card it named. It stores that card's
	// representative LocalID+1 (0 means "named nothing") so the zero value is cleanly
	// "none", like HostPlus; the gate compares by the named def's Name, so any copy
	// of that name in either deck is barred. resetCore clears it when the card leaves
	// play, lifting the bar.
	NamedCardPlus uint8
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
	// GiganticPartnerPlus links the two halves of a gigantic creature while it is
	// in play: the base half (which holds the battleline slot and is the creature)
	// and the slot-less art half (which lends only its bonus icons) each point at
	// the other. It uses +1 encoding (0 means "no partner") like the upgrade and
	// under links above, so the zero value is a lone, unpaired half. See ADR 0042.
	GiganticPartnerPlus uint8
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
	// UsagesThisTurn counts each card's usages this turn, indexed by LocalID — a
	// play, a discard, a use (reap/fight/action), each repeat iteration past the
	// first, each Destroyed: resolution, and each chained Replicator-style trigger
	// past the free first (charged to the card that started the chain). The
	// Rule of Six caps the usages of a whole card *name* (summed across every copy
	// regardless of owner, since only the active player uses cards), not a single
	// card, so it lives here beside Cards rather than in
	// CardCore, which resetCore would clear when a copy leaves play. StartTurn zeroes
	// the whole array.
	UsagesThisTurn [maxCards]uint8
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

	// Aember[p] is the Æmber in player p's pool, capped at maxAember like the
	// Æmber on a card, so one maximum governs every Æmber holder.
	Aember [2]int16

	// KeyColors[p] holds the colour of each key player p has forged, in forge order,
	// as a non-empty prefix followed by KeyColorNone. It is the only record of how
	// many keys a player has — KeyCount derives the count from it, so there is no
	// second copy of the same fact that could drift from the colours. A player picks
	// the colour as they forge (see pickKeyColor); a fourth key is KeyColorColorless.
	KeyColors [2][MaxKeys]KeyColor

	// ForgePrevented means a "when your opponent would forge a key" ability
	// cancelled the forge in progress (Keyforgery). It is set during the before-forge
	// window by CancelForge and read and cleared at the forge site before any Æmber
	// leaves the pool, mirroring FightCancelled; it is false outside that window.
	ForgePrevented bool

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

	// Scheduled holds the effects armed to resolve in the active player's end-of-turn
	// window (Ragnarok's board wipe), fired alongside the in-play "at the end of your
	// turn" abilities (ADR 0013); ScheduledCount is how many of the fixed array are in
	// use. Each is a one-shot armed during the play phase; unlike the turn bars it must
	// survive the ready phase (which runs before end of turn), so it is cleared only as
	// the window fires.
	Scheduled      [maxScheduled]ScheduledEffect
	ScheduledCount uint8

	// Key surcharges. KeyCostBump[p] raises player p's key cost for the current turn;
	// KeyCostBumpNext[p] arms that raise for p's next turn (Lash of Broken Dreams
	// makes keys cost +3 during the opponent's next turn). Unlike a card's
	// KeyCostChange, which lives as long as the card is in play, this is a one-turn
	// surcharge, promoted by StartTurn and lifted by the ready phase like the other
	// bars.
	KeyCostBump     [2]Bar[int]
	KeyCostBumpNext [2]Bar[int]

	// Counted key surcharges. KeyCostPerHouse[p] raises player p's key cost by Per
	// for each creature of House in play, recomputed at each forge; KeyCostPerHouseNext[p]
	// arms it for p's next turn (Waking Nightmare taxes +1 per Dis creature during
	// the opponent's next turn). Unlike KeyCostBump, the surcharge is not a frozen
	// amount — a creature of that house entering or leaving during the taxed turn
	// changes what a key costs. Promoted by StartTurn and lifted by the ready phase.
	KeyCostPerHouse     [2]Bar[perHouseKeySurcharge]
	KeyCostPerHouseNext [2]Bar[perHouseKeySurcharge]

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

	// MayUseTrait[p] is a trait whose creatures player p may fully use this turn
	// (fight, reap, or Action:) even when they are not in the active house —
	// Mutagenic Serum's "you may use friendly Mutant creatures". traitUnset grants
	// nothing; the ready phase clears it.
	MayUseTrait [2]Trait

	// TurnHistory holds the small tallies of what each player did during a turn —
	// several cards ask that rather than what is on the board ("if your opponent
	// forged a key on their previous turn", "for each enemy creature destroyed in a
	// fight this turn"). Keeping them as one array indexed by TurnStat leaves the
	// state flat and comparable, and makes a new tally one more enum value rather
	// than another pair of fields. The ready phase rolls each "this turn" tally
	// into its "last turn" twin, so "their previous turn" always means that
	// player's own last completed turn.
	TurnHistory [2][turnStatCount]int8

	// Lasting holds the "for the remainder of the turn" effects active now (Full
	// Moon, Charge!, Crystal Hive reactions; Dimension Door's replacement), fired or
	// queried by game_lasting.go when their event occurs; LastingCount is how many of
	// the fixed array are in use. The ready phase drops a player's entries.
	Lasting      [maxLasting]LastingEffect
	LastingCount uint8

	// Continuous holds the duration-scoped, read-live modifiers a resolving effect
	// installed on the board — damage immunity, lost keywords, blanked text, stat
	// overrides (game_continuous.go); ContinuousCount is how many of the fixed array
	// are in use. Unlike a ConstantAbility (a property of a card in play), a
	// continuous effect outlives its source and is dropped by turn window, not by
	// the source leaving play. The end-of-turn sweep drops entries whose window has
	// closed.
	Continuous      [maxContinuous]ContinuousEffect
	ContinuousCount uint8

	// AlsoTriggers holds the "for the remainder of the turn" also-triggers-on rules
	// active now (Livia the Elder's fight/reap fuse), queried by game_abilities.go
	// when a creature's abilities are gathered; AlsoTriggersCount is how many of the
	// fixed array are in use. The ready phase drops a player's entries.
	AlsoTriggers      [maxAlsoTriggers]LastingAlsoTriggersOn
	AlsoTriggersCount uint8

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

	// HouseConstraints[p] is the delayed constraint table binding player p's active
	// house choice THIS turn — the musts, cannots, and armed wagers that resolve
	// when p chooses a house (Control the Weak, Tezmal, Snag and its Mirror,
	// Snaglet). Entries accumulate (a must and a cannot stack; cannot overrides
	// must), and HouseConstraintCount[p] is how many of the fixed array are in use.
	// HouseConstraintsNext[p] arms entries for p's next turn; StartTurn promotes the
	// armed entries and clears them, so a constraint lands on the bound player's own
	// next choice. ChooseHouse resolves the whole table at once (see ADR 0035 and
	// allowedHouses in game_read.go).
	HouseConstraints         [2][maxHouseConstraints]HouseConstraint
	HouseConstraintCount     [2]uint8
	HouseConstraintsNext     [2][maxHouseConstraints]HouseConstraint
	HouseConstraintCountNext [2]uint8

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

	// PRNG is the match's random state (ADR 0039). Keeping it in the flat state — a
	// single counter word — is what makes a snapshot self-contained and replay
	// bit-exact: FastCopy captures it and re-running the command log reproduces
	// every shuffle and random pick. See prng.go.
	PRNG PRNG
}

// FastCopy returns an independent copy of the state. Because every field is a
// value type this is a single flat copy; mutating the result never affects the
// original.
func (s GameState) FastCopy() GameState { return s }

// KeyCount reports how many keys a player has forged, as the length of the
// non-empty prefix of the colours they forged. Deriving it means an unforge is a
// single write — blanking the last colour — and a count that disagrees with the
// colours is unrepresentable rather than merely invalid.
func (s GameState) KeyCount(player int) int {
	for i, c := range s.KeyColors[player] {
		if c == KeyColorNone {
			return i
		}
	}
	return MaxKeys
}

// ForgeCanonicalKeys sets a player's forged keys to n keys in canonical colour
// order, replacing whatever they had. The key count is derived from the colours,
// so a caller that only cares how many keys a player holds still has to name
// colours; this picks them so it does not have to.
func (s *GameState) ForgeCanonicalKeys(player, n int) {
	s.KeyColors[player] = [MaxKeys]KeyColor{}
	copy(s.KeyColors[player][:], firstKeyColors(n))
}

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

// count returns how many cards are registered in the match.
func (c *catalog) count() int { return len(c.defs) }

// def returns the definition for an id.
func (c *catalog) def(id LocalID) *CardDefinition { return c.defs[id] }

// owner returns the owning player index for an id.
func (c *catalog) owner(id LocalID) int { return int(c.owners[id]) }
