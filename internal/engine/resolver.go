package engine

// Resolver is the complete interface an effect uses to inspect and change the
// game. Effects hold only a Resolver (through EffectContext) — never the *Game or
// its GameState — so every change a card is able to make goes through one of
// these methods, the full explicit catalogue of what an effect may do; *Game
// implements it.
//
// The catalogue is deliberately wide, so it is composed from the focused role
// interfaces below (reads, economy, creature state, combat, zones, turn-scoped
// grants, choices, logging). An effect or a test double can depend on just the
// role it needs, and each new mechanic adds its method to the matching role —
// which keeps the capability clusters, and any gaps worth refactoring, visible in
// one place. Add a method to the role it belongs to, not to a flat list.
type Resolver interface {
	StateReader
	EconomyResolver
	CreatureResolver
	CombatResolver
	ZoneResolver
	TurnResolver
	ChoiceResolver
	Logger
}

// StateReader is the read-only view an effect inspects: pools, card stats and
// status, zone contents, and the turn's house and play counts. None of its
// methods change the game. It is the read model, composed from the same domain
// axes as the mutating roles below (economy, per-card creature state, zones,
// turn), so a consumer or test double can depend on just the reads it needs.
type StateReader interface {
	EconomyReader
	CreatureReader
	ZoneReader
	TurnReader
	// PlayerName returns a player's display name.
	PlayerName(player int) string
	// PlayerHasHouse reports whether house is one of the player's identity houses.
	PlayerHasHouse(player int, house House) bool
}

// EconomyReader reads the scoring economy: Æmber pools and forged keys. It mirrors
// EconomyResolver.
type EconomyReader interface {
	// Aember returns a player's Æmber pool.
	Aember(player int) int
	// AemberProtected reports whether a card the player controls makes their Æmber
	// immune to being stolen (The Vaultkeeper).
	AemberProtected(player int) bool
	// AemberTakenFromSupply reports whether Æmber a steal or capture takes from the
	// player's pool is drawn from the common supply instead (Po's Pixies).
	AemberTakenFromSupply(player int) bool
	// Keys returns the number of keys a player has forged.
	Keys(player int) int
	// KeyColors returns the colours of the keys a player has forged, in forge order.
	KeyColors(player int) []KeyColor
	// TideIsHigh reports whether the tide is high for the given player.
	TideIsHigh(player int) bool
	// TideIsLow reports whether the tide is low for the given player.
	TideIsLow(player int) bool
	// TurnHistory returns a player's running tally for a TurnStat.
	TurnHistory(player int, of TurnStat) int
}

// CreatureReader reads the in-play state carried on a single card — its stats,
// status, Æmber, house, owner, type, and traits. It mirrors CreatureResolver.
type CreatureReader interface {
	// Name returns a card's printed name.
	Name(id LocalID) string
	// Owner returns the player who owns a card.
	Owner(id LocalID) int
	// Controller returns the player a card in play currently answers to, which is
	// its owner unless an effect took control of it.
	Controller(id LocalID) int
	// Power returns a creature's current power (including upgrades).
	Power(id LocalID) int
	// Damage returns the damage currently on a creature.
	Damage(id LocalID) int
	// AmberOn returns the Æmber sitting on a card (from capture, exalt, ...).
	AmberOn(id LocalID) int
	// CountersOn returns how many generic counters of a kind sit on a card — the
	// card-placed markers a card reads (a doom counter for Wretched Doll).
	CountersOn(id LocalID, kind CounterKind) int
	// AemberBonus returns the number of Æmber pips printed on a card.
	AemberBonus(id LocalID) int
	// Exhausted reports whether a creature is exhausted.
	Exhausted(id LocalID) bool
	// InPlay reports whether a card is still on the board, as opposed to having been
	// destroyed, returned, or purged partway through an effect that is still
	// resolving.
	InPlay(id LocalID) bool
	// InBattleline reports whether a card currently sits on a battleline, as opposed
	// to having left play or being an artifact — only a battleline creature can be
	// moved to a flank.
	InBattleline(id LocalID) bool
	// InCenterOfBattleline reports whether a creature sits in the exact center of
	// its controller's battleline — the middle creature of an odd-sized line, with
	// equal creatures to its left and right. An even-sized line has no center.
	InCenterOfBattleline(id LocalID) bool
	// CurrentlyFighting reports whether a creature is one of the two combatants in
	// the fight resolving right now — read by a "while fighting" self-grant (Nizak,
	// The Forgotten gains invulnerable). It is false when no fight is in progress.
	CurrentlyFighting(id LocalID) bool
	// Armor is a creature's armor value, before any of it is spent or stripped.
	Armor(id LocalID) int
	// ArmorStripped is how much armor an effect has taken off a creature this turn,
	// as opposed to how much it spent absorbing damage.
	ArmorStripped(id LocalID) int
	// Hazardous is a creature's Hazardous value: an attacker takes this much damage
	// before fight damage is exchanged.
	Hazardous(id LocalID) int
	// Stunned reports whether a creature is stunned.
	Stunned(id LocalID) bool
	// Enraged reports whether a creature is enraged.
	Enraged(id LocalID) bool
	// Warded reports whether a creature has a ward.
	Warded(id LocalID) bool
	// TimesUsedThisTurn reports how many times a creature has been used this turn.
	TimesUsedThisTurn(id LocalID) int
	// IsCreature reports whether a card is a creature.
	IsCreature(id LocalID) bool
	// TypeOf returns a card's type.
	TypeOf(id LocalID) CardType
	// HasTrait reports whether a card has a trait.
	HasTrait(id LocalID, trait Trait) bool
	// TraitCount reports how many traits a card has.
	TraitCount(id LocalID) int
	// SharesTrait reports whether two cards share at least one trait.
	SharesTrait(a, b LocalID) bool
	// HasKeyword reports whether a creature has a keyword (printed or granted).
	HasKeyword(id LocalID, k Keyword) bool
	// HasTrigger reports whether a card has an ability under the trigger, whether
	// printed on it, granted by an attached upgrade, or granted by a constant
	// ability.
	HasTrigger(id LocalID, trigger Trigger) bool
	// ConsideredFlank reports whether a creature is treated as a flank creature for
	// the remainder of the turn regardless of its battleline position (Spectral
	// Tunneler).
	ConsideredFlank(id LocalID) bool
	// House returns a card's house.
	House(id LocalID) House
}

// ZoneReader reads the contents of a player's zones. It mirrors ZoneResolver.
type ZoneReader interface {
	// Battleline returns a copy of a player's creatures, safe to hold across
	// mutations.
	Battleline(player int) []LocalID
	// Artifacts returns a copy of a player's artifacts.
	Artifacts(player int) []LocalID
	// Upgrades returns the upgrades attached to a host, in order.
	Upgrades(host LocalID) []LocalID
	// HostOf returns the creature an attached upgrade is on, reporting ok=false
	// when the id is not attached to any creature.
	HostOf(upgrade LocalID) (LocalID, bool)
	// Hand returns a copy of a player's hand.
	Hand(player int) []LocalID
	// Discard returns a copy of a player's discard pile.
	Discard(player int) []LocalID
	// Deck returns a copy of a player's deck, top card first.
	Deck(player int) []LocalID
	// Archives returns a copy of a player's archived cards.
	Archives(player int) []LocalID
	// Purge returns a copy of a player's purged cards.
	Purge(player int) []LocalID
	// TopOfDeck returns the top card of a player's deck without moving it,
	// reporting whether the deck holds a card.
	TopOfDeck(player int) (LocalID, bool)
	// Under returns the ids of the cards placed under a host, face up or face
	// down, in the order they were placed (Masterplan, Jargogle, Graft).
	Under(host LocalID) []LocalID
	// UnderFaceDown reports whether a card placed under a host is facedown, as
	// opposed to faceup (Graft always places its card faceup).
	UnderFaceDown(id LocalID) bool
}

// TurnReader reads turn-scoped state: the active house and the cards played and
// discarded so far this turn. It mirrors TurnResolver.
type TurnReader interface {
	// ActiveHouse returns the house chosen for the current turn.
	ActiveHouse() House
	// ActivePlayer returns the player whose turn it is — the only player who can be
	// playing a card, so it is who an effect asks about "the cards played this turn"
	// when the effect belongs to neither player in particular.
	ActivePlayer() int
	// PlayedThisTurn returns the cards a player has played this turn, in play order.
	// Callers filter it themselves — by house, trait, or type — so the engine keeps
	// one record rather than a tally per axis.
	PlayedThisTurn(player int) []LocalID
	// DiscardedThisTurn returns the cards a player has discarded from hand this turn,
	// in discard order.
	DiscardedThisTurn(player int) []LocalID
}

// EconomyResolver changes the scoring economy: Æmber pools, forged keys, and
// chains.
type EconomyResolver interface {
	// SetAember sets a player's Æmber pool (never below zero).
	SetAember(player, amount int)
	// NoteAemberStolenFrom tallies Æmber stolen from a player this turn, so a card
	// can later ask whether they were robbed on their opponent's previous turn
	// (Information Exchange).
	NoteAemberStolenFrom(player, amount int)
	// EmitAemberStolenFrom fires each of the victim's in-play After Æmber Is Stolen
	// From You abilities, carrying the amount just stolen in this single theft
	// (Molephin deals damage scaled by it).
	EmitAemberStolenFrom(victim, amount int)
	// GainAember adds Æmber from the common supply to a player's pool, allowing
	// in-play replacements such as Ether Spider to capture it instead. It returns
	// the capturer and true when the gain was replaced.
	GainAember(player, amount int) (LocalID, bool)
	// StolenAemberCaptor returns a creature player controls that captures Æmber a
	// steal would otherwise add to player's pool, and true, when an in-play card
	// redirects stolen Æmber (Gargantodon) and player has a creature to hold it;
	// ok is false otherwise, leaving the steal to land in the pool as usual.
	StolenAemberCaptor(player int) (LocalID, bool)
	// ForgeKeyAtExtraCost has a player forge one key at the current cost plus a
	// surcharge for this forge only, if affordable (Key of Darkness forges at +6, an
	// unmodified forge at +0).
	ForgeKeyAtExtraCost(player, extra int)
	// ForgeKeyAtExtraCostReport forges one key at the current cost plus a surcharge
	// and reports whether a key was forged, so Obsidian Forge can destroy itself only
	// when a key was actually forged.
	ForgeKeyAtExtraCostReport(player, extra int) bool
	// RaiseKeyCostNextTurn raises what a player's keys cost throughout their next
	// turn (Lash of Broken Dreams).
	RaiseKeyCostNextTurn(player, amount int, source LocalID)
	// RaiseKeyCostThisTurn raises what a player's keys cost for the remainder of
	// the current turn, biting immediately rather than waiting for a turn boundary.
	RaiseKeyCostThisTurn(player, amount int, source LocalID)
	// ForgeKeyFree has a player forge one key without paying its current cost.
	ForgeKeyFree(player int)
	// UnforgeKey takes one forged key back off a player (Key Hammer).
	UnforgeKey(player int)
	// GainChains adds chains to a player, penalizing their future draws.
	GainChains(controller, amount int)
}

// CreatureResolver changes the in-play state carried on a single card — its
// damage, status, Æmber, counters, house, controller, and battleline position.
type CreatureResolver interface {
	// SetDamage sets the damage on a creature (never below zero).
	SetDamage(id LocalID, amount int)
	// StripArmor takes all of a creature's remaining armor away and records how much
	// was taken, so a following effect can scale with it (Red-Hot Armor).
	StripArmor(id LocalID)
	// SetStunned sets a creature's stun status.
	SetStunned(id LocalID, stunned bool)
	// SetEnraged sets a creature's enrage status.
	SetEnraged(id LocalID, enraged bool)
	// SetWarded sets a creature's ward status.
	SetWarded(id LocalID, warded bool)
	// SetDamageImmune marks a creature unable to be dealt damage for the remainder
	// of the turn.
	SetDamageImmune(id LocalID)
	// SetSideDamageImmune makes every creature player controls unable to be dealt
	// damage for the remainder of the turn, read live so creatures gained after it
	// resolves are covered too (Shield of Justice).
	SetSideDamageImmune(player int)
	// SetExhausted sets a creature's exhausted status.
	SetExhausted(id LocalID, exhausted bool)
	// AddAmberOn changes the Æmber sitting on a card.
	AddAmberOn(id LocalID, delta int)
	// AddPowerCounter changes the net power counters on a creature, adjusting its
	// power for as long as it stays in play.
	AddPowerCounter(id LocalID, delta int)
	// PlaceCounter puts n generic counters of a kind on an in-play card. Generic
	// counters are card-placed markers that only matter to cards that read them
	// (a doom counter for Wretched Doll); they are shed when the card leaves play.
	PlaceCounter(id LocalID, kind CounterKind, n int)
	// RemoveCounters drops every generic counter of one kind from a card, leaving
	// its other kinds untouched (Vineapple Tree sheds its growth counters after a
	// key is forged).
	RemoveCounters(id LocalID, kind CounterKind)
	// RemoveCountersN takes just n generic counters of one kind off a card,
	// dropping the entry once it reaches zero (The Colosseum removes six glory
	// counters to forge a key).
	RemoveCountersN(id LocalID, kind CounterKind, n int)
	// BelongToHouseForRemainderOfTurn makes a card belong to house until its
	// controller's turn ends.
	BelongToHouseForRemainderOfTurn(id LocalID, house House)
	// SetLastingHouse makes a card belong to house until it leaves play.
	SetLastingHouse(id LocalID, house House)
	// PutIntoBattlelineAsCreature turns an in-play card (an artifact, Auto-Legionary)
	// into a creature and moves it onto a flank of its controller's battleline, the
	// right flank when right is true. The card keeps its exhaustion and any power
	// counters, and reads as a creature until it leaves play.
	PutIntoBattlelineAsCreature(id LocalID, right bool)
	// SetNamedHouse records the house a card named as it entered play, which its
	// HouseLock then constrains for as long as the card stays in play.
	SetNamedHouse(id LocalID, house House)
	// TakeControl moves a card into controller's play area — a creature into their
	// battleline, an artifact into their artifact row — without changing ownership;
	// when it later leaves play it still goes to its owner's zone. source is the
	// card whose lasting effect holds the control, reverted when source leaves play;
	// a permanent (Forever) control names the seized card itself as source. Control
	// stacks LIFO, so a later take takes precedence and removing it falls back to
	// the one beneath.
	TakeControl(id LocalID, controller int, source LocalID)
	// SwapBattlelinePositions exchanges two creatures' positions in the same
	// battleline without moving any state between the creatures.
	SwapBattlelinePositions(a, b LocalID)
	// MoveToFlank moves one creature to a flank of its own controller's battleline:
	// the right flank when right is true, otherwise the left.
	MoveToFlank(id LocalID, right bool)
	// MoveWithinBattleline repositions one creature anywhere in its own
	// controller's battleline, chooser picking the destination slot (which may be
	// the creature's opponent — Malison moves an enemy creature).
	MoveWithinBattleline(chooser int, id LocalID)
	// SaveFromDestruction marks a creature whose own "Destroyed:" ability replaced
	// its destruction, so the current batch's discard step leaves it in play
	// (Reassembling Automaton).
	SaveFromDestruction(id LocalID)
	// LoseKeyword takes a keyword away from every creature in play for the
	// remainder of the turn.
	LoseKeyword(k Keyword)
	// GrantKeyword gives one creature a keyword for the remainder of the turn
	// (Scout grants Skirmish).
	GrantKeyword(id LocalID, k Keyword)
	// LoseKeywordFrom takes a keyword away from one creature for the remainder of
	// the turn (Niffle Grounds strips taunt and elusive).
	LoseKeywordFrom(id LocalID, k Keyword)
	// GrantKeywordUntilNextTurn gives one creature a keyword until the start of its
	// controller's next turn (Hideaway Hole grants elusive), surviving the
	// opponent's turn.
	GrantKeywordUntilNextTurn(id LocalID, k Keyword)
	// ConsiderFlank makes one creature count as a flank creature for the remainder
	// of the turn regardless of its position (Spectral Tunneler).
	ConsiderFlank(id LocalID)
	// GainStats gives one creature power and/or armor for the remainder of the turn
	// (Abond the Armorsmith grants +1 armor). Added armor also tops up the armor
	// left to absorb damage this turn.
	GainStats(id LocalID, power, armor int)
	// GainAssault gives one creature Assault for the remainder of the turn (Creed of
	// Nature grants assault equal to a chosen creature's power).
	GainAssault(id LocalID, amount int)
	// GrantTextBox gives creature recipient the printed text box of source — its
	// traits, keywords, and triggered abilities — either until recipient leaves
	// play (Mimic Gel) or for the remainder of the turn (Creed of Nurture) when
	// remainderOfTurn is set. It does not copy name, power, armor, type, or house.
	GrantTextBox(recipient, source LocalID, remainderOfTurn bool)
}

// CombatResolver resolves damage, destruction, and the fights, reaps, and actions
// an ability makes.
type CombatResolver interface {
	// DealDamage deals damage to each target simultaneously, then resolves
	// destruction (see the internal dealDamage).
	DealDamage(controller int, targets []DamageTarget)
	// DestroyEach destroys the given creatures as one simultaneous event.
	DestroyEach(controller int, ids []LocalID)
	// DestroyEachFrom is DestroyEach credited to a source card, so the batch it
	// targets narrates as "<source> destroys A, B, and C" in one line.
	DestroyEachFrom(controller int, source LocalID, ids []LocalID)
	// SetFightDamageRedirect redirects the attacker's fight damage in the fight in
	// progress to another creature (Gabos Longarms), read and cleared by combat.
	SetFightDamageRedirect(id LocalID)
	// CancelCurrentFight makes the fight in progress not occur. Combat reads and
	// clears it after Before Fight abilities resolve.
	CancelCurrentFight()
	// FightWith makes attacker fight defender (ability-driven, ignoring active
	// player and house). A creature can only be used while ready, so an exhausted
	// attacker may be chosen but does nothing.
	FightWith(attacker, defender LocalID)
	// ProtectedByTaunt reports whether target cannot be chosen to be fought by
	// attacker because a neighboring taunter shields it — respected by
	// ability-driven fights too, so a forced fight cannot reach past a taunter.
	ProtectedByTaunt(attacker, target LocalID) bool
	// ReapWith reaps with a creature (ability-driven, ignoring active player and
	// house). A creature can only be used while ready, so an exhausted creature may
	// be chosen but does nothing.
	ReapWith(id LocalID)
	// UseActionOf fires a card's "Action:" ability on behalf of actor (ability-driven,
	// ignoring active player and house), so a card can be used "as if it were yours".
	// A card can only be used while ready, so an exhausted card may be chosen but
	// does nothing.
	UseActionOf(actor int, id LocalID)
	// TriggerAbilityOf resolves a card's abilities under one trigger on behalf of
	// actor, without using the card: it does not exhaust and nothing watching for a
	// card being used fires (Replicator triggers another creature's reap effect).
	TriggerAbilityOf(actor int, id LocalID, trigger Trigger)
	// TriggerDepth is how many TriggerAbilityOf resolutions are already open, which
	// the Rule of Six bounds so a chain of them cannot run forever.
	TriggerDepth() int
}

// ZoneResolver moves cards between zones — drawing, and shuffling a card between
// play, hand, deck, discard, archives, and purge.
type ZoneResolver interface {
	// Draw makes a player draw count cards.
	Draw(controller, count int)
	// RefillHand refills a player's hand as if it were the end of their turn,
	// honoring their chains and draw modifiers (Punctuated Equilibrium).
	RefillHand(player int)
	// PutOnTopOfDeck moves a card from play to the top of its owner's deck.
	PutOnTopOfDeck(id LocalID)
	// PutIntoHand moves a card from play to its owner's hand.
	PutIntoHand(id LocalID)
	// ReturnUpgradesToHand moves each upgrade attached to a host in play to its
	// owner's hand, rather than shedding it to the discard pile.
	ReturnUpgradesToHand(host LocalID)
	// ArchiveUpgrade detaches an attached upgrade from its host and moves it to its
	// owner's archives (Ghostform archives itself off its host).
	ArchiveUpgrade(upgrade LocalID)
	// PutIntoArchives moves a card from play to its owner's archives.
	PutIntoArchives(id LocalID)
	// PutIntoArchivesEach archives a snapshot of in-play cards simultaneously, so
	// one card's leave-play does not destroy another still in the batch.
	PutIntoArchivesEach(controller int, ids []LocalID)
	// PutIntoYourArchives moves a card from play into player's own archives, which
	// may hold an enemy card (an abduction).
	PutIntoYourArchives(id LocalID, player int)
	// PutIntoDeckShuffled moves a card from play into its owner's deck and shuffles.
	PutIntoDeckShuffled(id LocalID)
	// ShuffleFriendlyCardsInPlayIntoDeck moves every card player controls in play —
	// each creature and artifact and their upgrades — into their deck, returning how
	// many cards were shuffled. The caller opens a shuffle batch around it.
	ShuffleFriendlyCardsInPlayIntoDeck(player int) int
	// BeginShuffleBatch starts collecting the cards shuffled into a deck until
	// EndShuffleBatch, so an effect that shuffles several creatures at once narrates
	// them as one grouped line per owner attributed to source.
	BeginShuffleBatch()
	// EndShuffleBatch closes the batch opened by BeginShuffleBatch, narrating the
	// collected cards grouped by owner as CardsShuffledIntoDeckBy from source.
	EndShuffleBatch(source LocalID)
	// ArchiveFromHand moves a card from its owner's hand to their archives.
	ArchiveFromHand(id LocalID)
	// ArchiveRandomFromHand moves one uniformly random card from a player's hand
	// to their archives (Eureka!).
	ArchiveRandomFromHand(owner int)
	// ArchiveFromDiscard moves a card from a player's discard pile to their archives.
	ArchiveFromDiscard(owner int, id LocalID)
	// ArchiveFromPurge moves a card from a player's purge pile to their archives —
	// a card recovered from out of the game (Universal Recycle Bin).
	ArchiveFromPurge(owner int, id LocalID)
	// ArchiveTopOfDeck moves the top card of a player's deck to their archives,
	// reporting whether a card was available.
	ArchiveTopOfDeck(player int) bool
	// ArchiveTopOfDiscard moves the top card of a player's discard pile to their
	// archives, reporting whether a card was available.
	ArchiveTopOfDiscard(player int) bool
	// DiscardTopOfDeck moves the top card of a player's deck to their discard pile,
	// returning that card and whether one was available.
	DiscardTopOfDeck(player int) (LocalID, bool)
	// DiscardArchives moves all of a player's archived cards to their discard pile.
	// The active player performs the discard, so they choose the order for their own
	// archives but get a random order for an opponent's (which they cannot see).
	DiscardArchives(owner int)
	// PurgeFromDiscard moves a card from a player's discard pile to their purge pile
	// (set aside out of the game).
	PurgeFromDiscard(owner int, id LocalID)
	// PurgeFromHand moves a card from a player's hand to their purge pile (set aside
	// out of the game).
	PurgeFromHand(owner int, id LocalID)
	// PurgeRandomFromHand moves one uniformly random card from a player's hand to
	// their purge pile, doing nothing if the hand is empty.
	PurgeRandomFromHand(owner int)
	// PurgeFromArchives moves a card from a player's archives to their purge pile
	// (set aside out of the game).
	PurgeFromArchives(owner int, id LocalID)
	// PurgeFromDeck moves a card from a player's deck to their purge pile (set aside
	// out of the game).
	PurgeFromDeck(owner int, id LocalID)
	// PurgeFromPlay moves a card from play to its owner's purge pile (set aside out
	// of the game).
	PurgeFromPlay(id LocalID)
	// MarkPlayedActionPurged marks a resolving action card to be set aside out of
	// the game when its play completes, rather than going to the discard pile
	// (Library Access purges itself).
	MarkPlayedActionPurged(id LocalID)
	// MarkPlayedActionArchived marks a resolving action card to go to its owner's
	// archives when its play completes, rather than going to the discard pile
	// (Sucker Punch archives itself).
	MarkPlayedActionArchived(id LocalID)
	// PutIntoPlay puts a card into play under controller's control without playing
	// it — no bonus icons and no Play: abilities resolve.
	PutIntoPlay(id LocalID, controller int)
	// PutFromDiscardIntoHand moves a card from its owner's discard to their hand.
	PutFromDiscardIntoHand(id LocalID)
	// MoveFromDeckToHand moves a card from its owner's deck to their hand.
	MoveFromDeckToHand(id LocalID)
	// MoveFromDeckToDiscard moves a card from its owner's deck to their discard
	// pile — a card the controller looked at and chose not to keep (Eyegor).
	MoveFromDeckToDiscard(id LocalID)
	// ArchiveFromDeck moves a card from its owner's deck to their archives — a card
	// the controller looked at and chose to archive (Philophosaurus).
	ArchiveFromDeck(id LocalID)
	// SetDeckTop rewrites the top len(order) cards of a player's deck to the given
	// order (order[0] becomes the new top) — the controller reordering the cards
	// they looked at (Navigator Ali). The ids must be exactly the cards currently
	// in those top positions, permuted.
	SetDeckTop(player int, order []LocalID)
	// PlayFromDeck plays a specific card from a player's deck, removing it from the
	// deck as it is played (Chaos Portal plays the card it revealed).
	PlayFromDeck(player int, id LocalID)
	// PlayFromDiscard plays a specific card from a player's discard pile, bypassing
	// the active-house gate (Sacrificial Altar). It does nothing when the card is
	// not in that discard pile.
	PlayFromDiscard(player int, id LocalID)
	// PlayFromOpponentDiscard plays a card out of the given player's opponent's
	// discard pile as that player's own play (Mimicry copies an action from the
	// other player's discard). It does nothing when the card is not in that discard
	// pile.
	PlayFromOpponentDiscard(player int, id LocalID)
	// PlayFromOpponentHand plays a card out of the given player's opponent's hand
	// as that player's own play (Lateral Shift plays a card out of the other
	// player's hand "as if it were yours"). The play counts against the active
	// player's own card-play limit and, for a creature or artifact, that player
	// takes control of it while its owner stays the opponent. It does nothing when
	// the card is not in that hand.
	PlayFromOpponentHand(player int, id LocalID)
	// PlayFromHand plays a specific card from a player's hand, bypassing the
	// active-house gate (Phase Shift's off-house card).
	PlayFromHand(player int, id LocalID)
	// PlayFromUnder plays a specific card from under whatever host it sits under
	// (Masterplan's and Jargogle's own "play the card under me"). It does nothing
	// when the card is not currently placed under anything.
	PlayFromUnder(player int, id LocalID)
	// PlayFromArchives plays a specific card from a player's archives, bypassing the
	// active-house gate (Project Z.Y.X.). It does nothing when the card is not in
	// that player's archives.
	PlayFromArchives(player int, id LocalID)
	// PlayRandomFromOpponentArchives plays a uniformly random card from player's
	// opponent's archives as player's own play, giving player control of it if it
	// stays in play (Murkens). It does nothing when those archives are empty.
	PlayRandomFromOpponentArchives(player int)
	// PlayTopOfOpponentDeck plays the top card of player's opponent's deck as
	// player's own play, giving player control of it if it stays in play (Murkens).
	// It does nothing when that deck is empty.
	PlayTopOfOpponentDeck(player int)
	// PutCardUnder removes a card from a player's hand and places it under host,
	// face up or face down (Masterplan, Jargogle).
	PutCardUnder(owner int, id, host LocalID, faceDown bool)
	// GraftUnder moves a card from play to faceup under host, out of play
	// (rulebook: Graft; Spangler Box).
	GraftUnder(id, host LocalID)
	// PutUnderIntoPlay puts every card placed under host into play under its
	// owner's control (Spangler Box's Destroyed ability).
	PutUnderIntoPlay(host LocalID)
	// ArchiveCardUnder moves each card placed under host to its owner's archives
	// (Jargogle's Destroyed ability when it is not its controller's turn).
	ArchiveCardUnder(host LocalID)
	// MoveUpgrade relocates an already-attached upgrade from its current host onto
	// newHost, keeping it in play (a movable Star Alliance "blaster" homing to its
	// signature creature). It does nothing when the upgrade is not attached.
	MoveUpgrade(upgrade, newHost LocalID)
	// ShuffleZonesIntoDeck moves each named zone's cards into a player's deck and
	// shuffles once (discard, hand, archives).
	ShuffleZonesIntoDeck(player int, zones []Zone)
	// SwapDeckAndDiscard exchanges a player's deck with their discard pile and
	// shuffles the new deck (Reverse Time).
	SwapDeckAndDiscard(player int)
	// MoveFromDiscardToTopOfDeck moves a card from its owner's discard to the top
	// of their deck.
	MoveFromDiscardToTopOfDeck(id LocalID)
	// ShuffleFromDiscardIntoDeck moves a card from its owner's discard pile into
	// their deck and shuffles, collected into a shuffle batch when one is open.
	ShuffleFromDiscardIntoDeck(id LocalID)
	// ShuffleFromHandIntoDeck moves a card from its owner's hand into their deck and
	// shuffles, collected into a shuffle batch when one is open.
	ShuffleFromHandIntoDeck(id LocalID)
	// Shuffle randomizes the order of a player's deck.
	Shuffle(player int)
	// DiscardCardFromHand moves a specific card from a player's hand to their discard
	// zone.
	DiscardCardFromHand(owner int, id LocalID)
	// DiscardRandomFromHand discards one uniformly random card from a player's hand,
	// doing nothing if the hand is empty.
	DiscardRandomFromHand(owner int)
	// DiscardRandomFromArchives discards one uniformly random card from a player's
	// archives, doing nothing if the archives are empty.
	DiscardRandomFromArchives(owner int)
}

// TurnResolver installs turn-scoped and lasting effects: restrictions and grants
// for this or the next turn, and the "remainder of the turn" reaction/replacement
// registry that keeps such behavior out of the play and reap paths.
type TurnResolver interface {
	// CannotFightNextTurn bars a player from using creatures to fight throughout
	// their next turn. source is the card imposing the bar, recorded so a frontend
	// can name it.
	CannotFightNextTurn(player int, source LocalID)
	// CannotPlayTypeNextTurn bars a player from playing cards of the given type
	// throughout their next turn (Lifeward, Scrambler Storm).
	CannotPlayTypeNextTurn(player int, t CardType, source LocalID)
	// CannotPlayTypeThisTurn bars a player from playing cards of the given type for
	// the rest of the current turn (Treasure Map, with the AnyType wildcard).
	CannotPlayTypeThisTurn(player int, t CardType, source LocalID)
	// CannotUseNextTurn bars a player from reaping, fighting, or using an "Action:"
	// ability throughout their next turn (Skippy Timehog).
	CannotUseNextTurn(player int, source LocalID)
	// CannotUseThisTurn bars a player from reaping, fighting, or using an "Action:"
	// ability for the rest of the current turn (United Action).
	CannotUseThisTurn(player int, source LocalID)
	// CannotReapNextTurn bars a player from reaping with any creature throughout
	// their next turn (Inky Gloom); fighting and "Action:" abilities stay open.
	CannotReapNextTurn(player int, source LocalID)
	// CannotReapThisTurn bars a player from reaping with any creature for the rest
	// of the current turn (Ragnarok); fighting and "Action:" abilities stay open.
	CannotReapThisTurn(player int, source LocalID)
	// CannotReapHouseNextTurn bars a player from reaping with creatures of the given
	// house throughout their next turn (Seismo-entangler).
	CannotReapHouseNextTurn(player int, house House, source LocalID)
	// CreaturesCannotUntilNextTurn arms a board-wide bar that stops both players
	// using creatures one way — fighting or reaping — until the caster's next turn,
	// sparing creatures of exceptHouse (HouseNone spares none): Into the Night,
	// Sow Salt.
	CreaturesCannotUntilNextTurn(caster int, action UseKind, exceptHouse House, source LocalID)
	// BlankEnemyText blanks the text box of every creature the given player controls
	// until the card's controller's next turn — its printed keywords, abilities, and
	// constant grants are ignored (Shadow of Dis). Its traits and stats remain.
	BlankEnemyText(player int, source LocalID)
	// SkipForgePhaseNextTurn makes a player skip their "forge a key" phase at the start
	// of their next turn (Miasma).
	SkipForgePhaseNextTurn(player int, source LocalID)
	// ScheduleDestroyEachCreatureAtEndOfTurn arms "destroy each creature" to resolve
	// in the active player's end-of-turn phase (Ragnarok). source is the card that
	// armed it, recorded for attribution.
	ScheduleDestroyEachCreatureAtEndOfTurn(source LocalID)
	// GrantFightForHouse lets a player use creatures of the given house to fight
	// this turn even out of the active house.
	GrantFightForHouse(player int, house House)
	// GrantFightAnyHouse lets every creature a player controls fight this turn,
	// whatever its house (Follow the Leader).
	GrantFightAnyHouse(player int)
	// GrantUseForHouse lets a player fully use (fight, reap, or Action:) creatures of
	// the given house this turn even out of the active house.
	GrantUseForHouse(player int, house House)
	// GrantPlayForHouse lets a player play cards of the given house from hand this
	// turn even out of the active house (the Ambassador cycle).
	GrantPlayForHouse(player int, house House)
	// GrantUseArtifactsAnyHouse lets a player use any friendly artifact this turn as
	// if it belonged to the active house (Scientifical Hack).
	GrantUseArtifactsAnyHouse(player int)
	// GrantOffHousePermit records a this-turn grant letting a player play or use a
	// bounded number of cards outside their active house (Com. Officer Kirby, CXO
	// Taber, United Action).
	GrantOffHousePermit(player int, p OffHousePermit)
	// AddLasting registers a "for the remainder of the turn" effect (Full Moon,
	// Charge!, Crystal Hive reactions; Dimension Door's replacement) on a game event,
	// instead of the effect hardcoding itself into the play or reap path. The record's
	// own fields narrow when it fires: Once for a one-shot (Blypyp), House/Type for a
	// matching subject, Except for the card that armed it (Library Access).
	AddLasting(le LastingEffect)
	// AddLastingMorph registers a "for the remainder of the turn" trigger morph
	// (Livia the Elder's fight/reap fuse) on the owning player, so a creature's
	// abilities under one trigger also fire on another until the turn ends.
	AddLastingMorph(m LastingMorph)
	// ForceActiveHouseNextTurn makes a player have to choose the given house as their
	// active house on their next turn.
	ForceActiveHouseNextTurn(player int, house House, source LocalID)
	// ForbidActiveHouseNextTurn makes a player unable to choose the given house as
	// their active house on their next turn.
	ForbidActiveHouseNextTurn(player int, house House, source LocalID)
	// WagerOnHouseNextTurn arms a bet on a player's next active house: if they
	// choose that house, the predictor steals amount (Snaglet).
	WagerOnHouseNextTurn(player int, house House, amount, predictor int, source LocalID)
	// SetActiveHouse reassigns the active player's active house for the current
	// turn without a house choice (Book of leQ makes its revealed card's house
	// active).
	SetActiveHouse(house House)
	// EndTurnNow ends the active player's turn in place, running the turn out the
	// way the Omega keyword does (Book of leQ ends your turn).
	EndTurnNow()
}

// ChoiceResolver asks a player to make a decision — ordering a set of cards, or
// picking a creature, card, or labeled option — so an effect can branch on the
// answer. A sole candidate is taken automatically.
type ChoiceResolver interface {
	// OrderByChoice lets a player arrange ids into a resolution order.
	OrderByChoice(controller int, prompt string, ids []LocalID) []LocalID
	// ChooseCreature asks a player to pick one creature from candidates; a sole
	// candidate is taken automatically. source is the card whose ability is asking
	// (usually ctx.Source), for prompt attribution.
	ChooseCreature(player int, source LocalID, prompt string, candidates []LocalID) (LocalID, bool)
	// ChooseCard asks a player to pick one card from candidates; a sole candidate
	// is taken automatically. source is the card whose ability is asking (usually
	// ctx.Source), for prompt attribution.
	ChooseCard(player int, source LocalID, prompt string, candidates []LocalID) (LocalID, bool)
	// ChooseCardOptional asks a player to pick one card from candidates or to
	// decline. Unlike ChooseCard a sole candidate is still offered, because passing
	// is a legal answer. source is the card whose ability is asking (usually
	// ctx.Source), for prompt attribution.
	ChooseCardOptional(
		player int,
		source LocalID,
		prompt string,
		candidates []LocalID,
	) (LocalID, bool)
	// ChooseOption asks a player to pick one of several labeled options, returning
	// its index (0 when the player's chooser expresses no preference). source is the
	// card whose ability is asking (usually ctx.Source), for prompt attribution.
	ChooseOption(player int, source LocalID, prompt string, options []string) int
}

// Logger narrates resolved outcomes to the game log (ADR 0011). An effect does
// not write prose: it records a typed entry describing what actually happened,
// and the entry renders itself for whoever is reading.
type Logger interface {
	// Record appends one narrated outcome to the game log.
	Record(e LogEntry)
}
