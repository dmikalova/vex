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
	// Keys returns the number of keys a player has forged.
	Keys(player int) int
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
	// Exhausted reports whether a creature is exhausted.
	Exhausted(id LocalID) bool
	// InPlay reports whether a card is still on the board, as opposed to having been
	// destroyed, returned, or purged partway through an effect that is still
	// resolving.
	InPlay(id LocalID) bool
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
	// TimesUsedThisTurn reports how many times a creature has been used this turn.
	TimesUsedThisTurn(id LocalID) int
	// IsCreature reports whether a card is a creature.
	IsCreature(id LocalID) bool
	// TypeOf returns a card's type.
	TypeOf(id LocalID) CardType
	// HasTrait reports whether a card has a trait.
	HasTrait(id LocalID, trait Trait) bool
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
	// GainAember adds Æmber from the common supply to a player's pool, allowing
	// in-play replacements such as Ether Spider to capture it instead. It returns
	// the capturer and true when the gain was replaced.
	GainAember(player, amount int) (LocalID, bool)
	// ForgeKeyAtExtraCost has a player forge one key at the current cost plus a
	// surcharge for this forge only, if affordable (Key of Darkness forges at +6, an
	// unmodified forge at +0).
	ForgeKeyAtExtraCost(player, extra int)
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
	// PreventDamage marks a creature immune to damage for the remainder of the turn.
	PreventDamage(id LocalID)
	// SetExhausted sets a creature's exhausted status.
	SetExhausted(id LocalID, exhausted bool)
	// AddAmberOn changes the Æmber sitting on a card.
	AddAmberOn(id LocalID, delta int)
	// AddPowerCounter changes the net power counters on a creature, adjusting its
	// power for as long as it stays in play.
	AddPowerCounter(id LocalID, delta int)
	// BelongToHouseForRemainderOfTurn makes a card belong to house until its
	// controller's turn ends.
	BelongToHouseForRemainderOfTurn(id LocalID, house House)
	// SetLastingHouse makes a card belong to house until it leaves play.
	SetLastingHouse(id LocalID, house House)
	// SetNamedHouse records the house a card named as it entered play, which its
	// HouseLock then constrains for as long as the card stays in play.
	SetNamedHouse(id LocalID, house House)
	// TakeControl moves a creature into controller's battleline without changing
	// ownership; when it later leaves play it still goes to its owner's zone. source
	// is the card whose lasting effect holds the control (reverted when it leaves play).
	TakeControl(id LocalID, controller int, source LocalID)
	// TakeControlOfArtifact moves an artifact into controller's artifact row for
	// good, without changing ownership. There is no reverting source.
	TakeControlOfArtifact(id LocalID, controller int)
	// SwapBattlelinePositions exchanges two creatures' positions in the same
	// battleline without moving any state between the creatures.
	SwapBattlelinePositions(a, b LocalID)
	// MoveToFlank moves one creature to a flank of its own controller's battleline:
	// the right flank when right is true, otherwise the left.
	MoveToFlank(id LocalID, right bool)
	// LoseKeyword takes a keyword away from every creature in play for the
	// remainder of the turn.
	LoseKeyword(k Keyword)
	// GrantKeyword gives one creature a keyword for the remainder of the turn
	// (Scout grants Skirmish).
	GrantKeyword(id LocalID, k Keyword)
	// ConsiderFlank makes one creature count as a flank creature for the remainder
	// of the turn regardless of its position (Spectral Tunneler).
	ConsiderFlank(id LocalID)
}

// CombatResolver resolves damage, destruction, and the fights, reaps, and actions
// an ability drives directly, including the mid-fight adjustments a Before Fight
// ability makes.
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
	// PutOnTopOfDeck moves a card from play to the top of its owner's deck.
	PutOnTopOfDeck(id LocalID)
	// PutIntoHand moves a card from play to its owner's hand.
	PutIntoHand(id LocalID)
	// PutIntoArchives moves a card from play to its owner's archives.
	PutIntoArchives(id LocalID)
	// PutIntoYourArchives moves a card from play into player's own archives, which
	// may hold an enemy card (an abduction).
	PutIntoYourArchives(id LocalID, player int)
	// PutIntoDeckShuffled moves a card from play into its owner's deck and shuffles.
	PutIntoDeckShuffled(id LocalID)
	// BeginShuffleBatch starts collecting the cards shuffled into a deck until
	// EndShuffleBatch, so an effect that shuffles several creatures at once narrates
	// them as one grouped line per owner attributed to source.
	BeginShuffleBatch()
	// EndShuffleBatch closes the batch opened by BeginShuffleBatch, narrating the
	// collected cards grouped by owner as CardsShuffledIntoDeckBy from source.
	EndShuffleBatch(source LocalID)
	// ArchiveFromHand moves a card from its owner's hand to their archives.
	ArchiveFromHand(id LocalID)
	// ArchiveFromDiscard moves a card from a player's discard pile to their archives.
	ArchiveFromDiscard(owner int, id LocalID)
	// ArchiveTopOfDeck moves the top card of a player's deck to their archives,
	// reporting whether a card was available.
	ArchiveTopOfDeck(player int) bool
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
	// PurgeFromPlay moves a card from play to its owner's purge pile (set aside out
	// of the game).
	PurgeFromPlay(id LocalID)
	// MarkPlayedActionPurged marks a resolving action card to be set aside out of
	// the game when its play completes, rather than going to the discard pile
	// (Library Access purges itself).
	MarkPlayedActionPurged(id LocalID)
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
	// PlayFromHand plays a specific card from a player's hand, bypassing the
	// active-house gate (Phase Shift's off-house card).
	PlayFromHand(player int, id LocalID)
	// PlayFromUnder plays a specific card from under whatever host it sits under
	// (Masterplan's and Jargogle's own "play the card under me"). It does nothing
	// when the card is not currently placed under anything.
	PlayFromUnder(player int, id LocalID)
	// PutCardUnder removes a card from a player's hand and places it under host,
	// face up or face down (Masterplan, Jargogle).
	PutCardUnder(owner int, id, host LocalID, faceDown bool)
	// GraftUnder moves a card from play to faceup under host, out of play
	// (rulebook: Graft; Spangler Box).
	GraftUnder(id, host LocalID)
	// PutUnderIntoPlay puts every card placed under host into play under its
	// owner's control (Spangler Box's Destroyed ability).
	PutUnderIntoPlay(host LocalID)
	// ShuffleZonesIntoDeck moves each named zone's cards into a player's deck and
	// shuffles once (discard, hand, archives).
	ShuffleZonesIntoDeck(player int, zones []Zone)
	// SwapDeckAndDiscard exchanges a player's deck with their discard pile and
	// shuffles the new deck (Reverse Time).
	SwapDeckAndDiscard(player int)
	// MoveFromDiscardToTopOfDeck moves a card from its owner's discard to the top
	// of their deck.
	MoveFromDiscardToTopOfDeck(id LocalID)
	// DiscardCardFromHand moves a specific card from a player's hand to their discard
	// zone.
	DiscardCardFromHand(owner int, id LocalID)
	// DiscardRandomFromHand discards one uniformly random card from a player's hand,
	// doing nothing if the hand is empty.
	DiscardRandomFromHand(owner int)
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
	// SkipForgeStepNextTurn makes a player skip their "forge a key" step at the start
	// of their next turn (Miasma).
	SkipForgeStepNextTurn(player int, source LocalID)
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
	// AddLasting registers a "for the remainder of the turn" effect (Full Moon,
	// Charge!, Crystal Hive reactions; Dimension Door's replacement) on a game event,
	// instead of the effect hardcoding itself into the play or reap path. The record's
	// own fields narrow when it fires: Once for a one-shot (Blypyp), House/Type for a
	// matching subject, Except for the card that armed it (Library Access).
	AddLasting(le LastingEffect)
	// ForceActiveHouseNextTurn makes a player have to choose the given house as their
	// active house on their next turn.
	ForceActiveHouseNextTurn(player int, house House, source LocalID)
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

// The read accessors Aember, AmberOn, Damage, Power, Name, PlayerName,
// Battleline, and Artifacts are defined in game.go; the remaining Resolver
// methods follow.

// Owner returns the player who owns a card.
func (g *Game) Owner(id LocalID) int { return g.owner(id) }

// Controller returns the player a card in play currently answers to.
func (g *Game) Controller(id LocalID) int { return g.controller(id) }

// HasTrait reports whether a card has a trait.
func (g *Game) HasTrait(id LocalID, trait Trait) bool { return g.cat.def(id).hasTrait(trait) }

// SharesTrait reports whether two cards have at least one trait in common.
func (g *Game) SharesTrait(a, b LocalID) bool {
	other := g.cat.def(b)
	for _, tr := range g.cat.def(a).Traits {
		if other.hasTrait(tr) {
			return true
		}
	}
	return false
}

// HasKeyword reports whether a creature has a keyword, printed or granted.
func (g *Game) HasKeyword(id LocalID, k Keyword) bool { return g.hasKeyword(id, k) }

// ProtectedByTaunt reports whether target is shielded from attacker by a
// neighboring taunter (the exported CombatResolver port method).
func (g *Game) ProtectedByTaunt(attacker, target LocalID) bool {
	return g.protectedByTaunt(attacker, target)
}

// LoseKeyword takes a keyword away from every creature in play for the remainder
// of the turn (Sniffer).
func (g *Game) LoseKeyword(k Keyword) {
	g.State.KeywordsLost |= k.bit()
	g.record(KeywordLostByAll{Keyword: k})
}

// GrantKeyword gives one creature a keyword for the remainder of the turn (Scout).
func (g *Game) GrantKeyword(id LocalID, k Keyword) {
	if g.State.Cards[id].GrantedKeywords&k.bit() != 0 {
		return
	}
	g.State.Cards[id].GrantedKeywords |= k.bit()
	g.record(CreatureGainedKeyword{Creature: id, Keyword: k})
}

// ConsideredFlank reports whether a creature counts as a flank creature for the
// turn regardless of its battleline position (Spectral Tunneler).
func (g *Game) ConsideredFlank(id LocalID) bool { return g.State.Cards[id].ConsideredFlank }

// ConsiderFlank makes one creature count as a flank creature for the remainder of
// the turn (Spectral Tunneler).
func (g *Game) ConsiderFlank(id LocalID) {
	if g.State.Cards[id].ConsideredFlank {
		return
	}
	g.State.Cards[id].ConsideredFlank = true
	g.record(CreatureConsideredFlank{Creature: id})
}

// ForgeKeyAtExtraCost has a player forge one key at its current cost plus extra.
func (g *Game) ForgeKeyAtExtraCost(player, extra int) { g.forgeKeyAtExtraCost(player, extra) }

// ForgeKeyFree has a player forge one key without paying its current cost.
func (g *Game) ForgeKeyFree(player int) { g.forgeKeyFree(player) }

// IsCreature reports whether a card is a creature.
func (g *Game) IsCreature(id LocalID) bool { return g.cat.def(id).Type == Creature }

// TypeOf returns a card's type.
func (g *Game) TypeOf(id LocalID) CardType { return g.cat.def(id).Type }

// SetAember sets a player's Æmber pool, clamped at zero.
func (g *Game) SetAember(player, amount int) {
	if amount < 0 {
		amount = 0
	}
	g.State.Aember[player] = amount
}

// stateOf returns a card's mutable in-play state, or nil once it has left play.
// An ability keeps resolving after its source or target dies — Zyzzix the Many
// adds power counters to itself after its own upgrade's damage destroyed it — and
// KeyForge lets the rest of it resolve, so the parts that need a card on the board
// have to land on nothing. Writing anyway leaves a card in the discard pile
// carrying counters or a stun that nothing will ever clear.
func (g *Game) stateOf(id LocalID) *CardCore {
	if !g.inPlay(id) {
		return nil
	}
	return &g.State.Cards[id]
}

// SetDamage sets the damage on a creature, clamped at zero.
func (g *Game) SetDamage(id LocalID, amount int) {
	if amount < 0 {
		amount = 0
	}
	if c := g.stateOf(id); c != nil {
		c.Damage = int16(amount)
	}
}

// StripArmor empties a creature's remaining armor and tallies what was taken, so
// a following effect can scale with it. It adds to any earlier strip this turn
// rather than replacing it, so two effects that each strip armor both count.
func (g *Game) StripArmor(id LocalID) {
	c := g.stateOf(id)
	if c == nil {
		return
	}
	c.ArmorStripped += c.ArmorRemaining
	c.ArmorRemaining = 0
}

// SetStunned sets a creature's stun status.
func (g *Game) SetStunned(id LocalID, stunned bool) {
	if c := g.stateOf(id); c != nil {
		c.Stunned = stunned
	}
}

// PreventDamage marks a creature immune to damage for the remainder of the turn.
func (g *Game) PreventDamage(id LocalID) {
	if c := g.stateOf(id); c != nil {
		c.DamageImmune = true
	}
}

// SetExhausted sets a creature's exhausted status.
func (g *Game) SetExhausted(id LocalID, exhausted bool) {
	if c := g.stateOf(id); c != nil {
		c.Exhausted = exhausted
	}
}

// BelongToHouseForRemainderOfTurn makes a card belong to house until its
// controller's turn ends.
func (g *Game) BelongToHouseForRemainderOfTurn(id LocalID, house House) {
	if c := g.stateOf(id); c != nil {
		c.TempHouse = house
	}
}

// SetLastingHouse makes a card belong to house until it leaves play.
func (g *Game) SetLastingHouse(id LocalID, house House) {
	if c := g.stateOf(id); c != nil {
		c.LastingHouse = house
	}
}

// SetNamedHouse records the house a card named as it entered play, which its
// HouseLock then constrains for as long as the card stays in play.
func (g *Game) SetNamedHouse(id LocalID, house House) {
	if c := g.stateOf(id); c != nil {
		c.NamedHouse = house
	}
}

// SetFightDamageRedirect redirects the attacker's fight damage in the current
// fight to another creature; the combat step reads and clears it.
func (g *Game) SetFightDamageRedirect(id LocalID) { g.State.FightDamageRedirect = id }

// CancelCurrentFight makes the fight in progress not occur; the combat step reads
// and clears it before Assault, Hazardous, and fight damage.
func (g *Game) CancelCurrentFight() { g.State.FightCancelled = true }

// AddAmberOn changes the Æmber sitting on a card.
func (g *Game) AddAmberOn(id LocalID, delta int) { g.addAmberOn(id, delta) }

// DealDamage is the Resolver entry point for the internal dealDamage.
func (g *Game) DealDamage(controller int, targets []DamageTarget) {
	g.dealDamage(controller, targets...)
}

// DestroyEach is the Resolver entry point for destroyEach.
func (g *Game) DestroyEach(controller int, ids []LocalID) { g.destroyEach(controller, ids) }

// DestroyEachFrom credits a source card for the destruction, so the batch it
// targets narrates as one grouped line.
func (g *Game) DestroyEachFrom(controller int, source LocalID, ids []LocalID) {
	prevS, prevH := g.destroyingSource, g.hasDestroyingSource
	g.destroyingSource, g.hasDestroyingSource = source, true
	g.destroyEach(controller, ids)
	g.destroyingSource, g.hasDestroyingSource = prevS, prevH
}

// TakeControl is the Resolver entry point for takeControl.
func (g *Game) TakeControl(id LocalID, controller int, source LocalID) {
	g.takeControl(id, controller, source)
}

// TakeControlOfArtifact is the Resolver entry point for takeControlOfArtifact.
func (g *Game) TakeControlOfArtifact(id LocalID, controller int) {
	g.takeControlOfArtifact(id, controller)
}

// PutIntoPlay is the Resolver entry point for putIntoPlay.
func (g *Game) PutIntoPlay(id LocalID, controller int) {
	g.putIntoPlay(id, controller)
}

// PlayerHasHouse reports whether house is one of the player's identity houses.
func (g *Game) PlayerHasHouse(player int, house House) bool {
	return g.playerHasHouse(player, house)
}

// Draw is the Resolver entry point for the internal draw.
func (g *Game) Draw(controller, count int) { g.draw(controller, count) }

// PutOnTopOfDeck is the Resolver entry point for putOnTopOfDeck.
func (g *Game) PutOnTopOfDeck(id LocalID) { g.putOnTopOfDeck(id) }

// PutIntoHand is the Resolver entry point for putIntoHand.
func (g *Game) PutIntoHand(id LocalID) { g.putIntoHand(id) }

// PutIntoArchives is the Resolver entry point for putIntoArchives.
func (g *Game) PutIntoArchives(id LocalID) { g.putIntoArchives(id) }

// PutIntoDeckShuffled is the Resolver entry point for putIntoDeckShuffled.
func (g *Game) PutIntoDeckShuffled(id LocalID) { g.putIntoDeckShuffled(id) }

// BeginShuffleBatch opens a shuffle batch: cards shuffled into a deck until
// EndShuffleBatch are collected rather than narrated one by one.
func (g *Game) BeginShuffleBatch() {
	g.shuffleBatch, g.batchingShuffle = nil, true
}

// EndShuffleBatch closes the batch and narrates the collected cards grouped by
// owner as one CardsShuffledIntoDeckBy line each, attributed to source, in the
// order the owners were first shuffled.
func (g *Game) EndShuffleBatch(source LocalID) {
	batch := g.shuffleBatch
	g.shuffleBatch, g.batchingShuffle = nil, false
	byOwner := map[int][]LocalID{}
	var owners []int
	for _, id := range batch {
		o := g.owner(id)
		if _, seen := byOwner[o]; !seen {
			owners = append(owners, o)
		}
		byOwner[o] = append(byOwner[o], id)
	}
	for _, o := range owners {
		g.record(CardsShuffledIntoDeckBy{Source: source, Owner: o, Cards: byOwner[o]})
	}
}

// ArchiveFromHand moves a card from its owner's hand to their archives.
func (g *Game) ArchiveFromHand(id LocalID) { g.archiveFromHand(g.owner(id), id) }

// ArchiveFromDiscard moves a card from a player's discard pile to their archives.
func (g *Game) ArchiveFromDiscard(owner int, id LocalID) { g.archiveFromDiscard(owner, id) }

// ArchiveTopOfDeck moves the top card of a player's deck to their archives.
func (g *Game) ArchiveTopOfDeck(player int) bool { return g.archiveTopOfDeck(player) }

// DiscardArchives moves all of a player's archived cards to their discard pile.
func (g *Game) DiscardArchives(owner int) { g.discardArchives(owner) }

// PurgeFromDiscard moves a card from a player's discard pile to their purge pile.
func (g *Game) PurgeFromDiscard(owner int, id LocalID) { g.purgeFromDiscard(owner, id) }

// PurgeFromHand moves a card from a player's hand to their purge pile.
func (g *Game) PurgeFromHand(owner int, id LocalID) { g.purgeFromHand(owner, id) }

// PurgeFromPlay is the Resolver entry point for purgeFromPlay.
func (g *Game) PurgeFromPlay(id LocalID) { g.purgeFromPlay(id) }

// MarkPlayedActionPurged marks a resolving action to be purged instead of
// discarded when its play completes (Library Access).
func (g *Game) MarkPlayedActionPurged(id LocalID) {
	g.State.PurgePlayedAction = id
	g.State.PurgePlayedActionSet = true
}

// AddPowerCounter changes the net power counters on a creature.
func (g *Game) AddPowerCounter(id LocalID, delta int) {
	if c := g.stateOf(id); c != nil {
		c.PowerCounters += int16(delta)
	}
}

// PutFromDiscardIntoHand moves a card from its owner's discard pile to their hand.
func (g *Game) PutFromDiscardIntoHand(id LocalID) {
	o := g.owner(id)
	g.State.Discard[o].remove(id)
	g.State.Hand[o].add(id)
	g.record(CardReturnedFromDiscardToHand{Player: o, Card: id})
}

// MoveFromDeckToHand moves a card from its owner's deck to their hand.
func (g *Game) MoveFromDeckToHand(id LocalID) {
	o := g.owner(id)
	g.State.Deck[o].remove(id)
	g.State.Hand[o].add(id)
	g.record(CardPutFromDeckIntoHand{Player: o, Card: id})
}

// MoveFromDeckToDiscard moves a card from its owner's deck to their discard pile.
func (g *Game) MoveFromDeckToDiscard(id LocalID) {
	o := g.owner(id)
	g.State.Deck[o].remove(id)
	g.State.Discard[o].add(id)
	g.record(CardDiscardedFromDeck{Player: o, Card: id})
}

// ShuffleZonesIntoDeck moves each named zone's cards into a player's deck and
// shuffles once.
func (g *Game) ShuffleZonesIntoDeck(player int, zones []Zone) {
	g.shuffleZonesIntoDeck(player, zones)
}

// MoveFromDiscardToTopOfDeck moves a card from its owner's discard pile to the
// top of their deck.
func (g *Game) MoveFromDiscardToTopOfDeck(id LocalID) {
	o := g.owner(id)
	g.State.Discard[o].remove(id)
	g.State.Deck[o].addFront(id)
	g.record(CardPutFromDiscardOnTopOfDeck{Player: o, Card: id})
}

// GainChains adds chains to a player, which reduce their draws until shed.
func (g *Game) GainChains(controller, amount int) {
	g.State.Chains[controller] += amount
	g.record(ChainsGained{
		Player: controller,
		Amount: amount,
		Total:  g.State.Chains[controller],
	})
}

// OrderByChoice is the Resolver entry point for orderByChoice.
func (g *Game) OrderByChoice(controller int, prompt string, ids []LocalID) []LocalID {
	return g.orderByChoice(controller, prompt, ids)
}

// ChooseCreature asks a player to choose one creature from candidates, attributing
// the prompt to the source card. A sole candidate is taken automatically.
func (g *Game) ChooseCreature(
	player int,
	source LocalID,
	prompt string,
	candidates []LocalID,
) (LocalID, bool) {
	return g.pickCreature(player, g.sourceName(source), prompt, candidates)
}

// ChooseCard asks a player to choose one card from candidates, attributing the
// prompt to the source card. A sole candidate is taken automatically.
func (g *Game) ChooseCard(
	player int,
	source LocalID,
	prompt string,
	candidates []LocalID,
) (LocalID, bool) {
	return g.pickCard(player, g.sourceName(source), prompt, candidates)
}

// ChooseCardOptional asks a player to choose one card from candidates or to
// decline, attributing the prompt to the source card. A sole candidate is still
// offered rather than forced, because declining is a legal answer.
func (g *Game) ChooseCardOptional(
	player int,
	source LocalID,
	prompt string,
	candidates []LocalID,
) (LocalID, bool) {
	return g.pickOptional(player, g.sourceName(source), prompt, candidates)
}

// ChooseOption asks a player to choose one of several labeled options, attributing
// (does not implement OptionChooser), the first option is taken.
func (g *Game) ChooseOption(player int, source LocalID, prompt string, options []string) int {
	return g.chooseOption(player, g.sourceName(source), prompt, options)
}

// chooseOption is the shared option-choice path: it attributes the prompt to a
// source name (empty for a source-less prompt such as a turn-structure choice)
// and defaults to the first option when the chooser has no preference. A sole
// option is taken automatically without consulting the chooser.
func (g *Game) chooseOption(player int, source, prompt string, options []string) int {
	if len(options) == 1 {
		return 0
	}
	if oc, ok := g.chooserFor(player).(OptionChooser); ok {
		return oc.ChooseOption(source, renderPrompt(source, prompt), options)
	}
	return 0
}

// sourceName returns the name of a source card for prompt attribution, or "" when
// the id is not a registered card (e.g. an unset source in a unit test).
func (g *Game) sourceName(source LocalID) string {
	if int(source) < len(g.cat.defs) {
		return g.cat.defs[source].Name
	}
	return ""
}

// FightWith makes attacker fight defender, ability-driven (ignoring active player
// and house). A creature can only be used while ready, so an exhausted attacker
// does nothing.
func (g *Game) FightWith(attacker, defender LocalID) {
	if g.readyToUse(attacker) {
		g.fight(attacker, defender)
	}
}

// ReapWith reaps with a creature, ability-driven (ignoring active player and
// house). A creature can only be used while ready, so an exhausted creature does
// nothing.
func (g *Game) ReapWith(id LocalID) {
	if g.readyToUse(id) {
		g.reapWith(id)
	}
}

// UseActionOf fires a card's "Action:" ability on behalf of actor, ability-driven
// (ignoring active player and house). A card can only be used while ready, so an
// exhausted card does nothing.
func (g *Game) UseActionOf(actor int, id LocalID) {
	if g.readyToUse(id) {
		g.useActionOf(actor, id)
	}
}

// TriggerAbilityOf resolves a card's abilities under one trigger on behalf of
// actor. Unlike UseActionOf this does not use the card, so readiness is beside
// the point: an exhausted creature's reap effect still triggers, and the creature
// neither exhausts nor counts as used.
func (g *Game) TriggerAbilityOf(actor int, id LocalID, trigger Trigger) {
	g.triggerDepth++
	defer func() { g.triggerDepth-- }()
	g.triggerAbilitiesAs(actor, id, trigger, 0, false)
}

// TriggerDepth is how many TriggerAbilityOf resolutions are currently open.
func (g *Game) TriggerDepth() int { return g.triggerDepth }

// HasTrigger reports whether a card has an ability under the trigger.
func (g *Game) HasTrigger(id LocalID, trigger Trigger) bool {
	return g.hasTrigger(id, trigger)
}

// Record appends one narrated outcome to the game log (ADR 0011).
func (g *Game) Record(e LogEntry) { g.record(e) }
