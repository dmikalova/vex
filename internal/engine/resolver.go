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
// methods change the game.
type StateReader interface {
	// Aember returns a player's Æmber pool.
	Aember(player int) int
	// AmberOn returns the Æmber sitting on a card (from capture, exalt, ...).
	AmberOn(id LocalID) int
	// Damage returns the damage currently on a creature.
	Damage(id LocalID) int
	// Power returns a creature's current power (including upgrades).
	Power(id LocalID) int
	// Name returns a card's printed name.
	Name(id LocalID) string
	// PlayerName returns a player's display name.
	PlayerName(player int) string
	// Owner returns the player who owns a card.
	Owner(id LocalID) int
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
	// IsCreature reports whether a card is a creature.
	IsCreature(id LocalID) bool
	// TypeOf returns a card's type.
	TypeOf(id LocalID) CardType
	// HasTrait reports whether a card has a trait.
	HasTrait(id LocalID, trait Trait) bool
	// House returns a card's house.
	House(id LocalID) House
	// ActiveHouse returns the house chosen for the current turn.
	ActiveHouse() House
	// Keys returns the number of keys a player has forged.
	Keys(player int) int
	// CardsPlayedOfHouseThisTurn returns the house-specific play count for this turn.
	CardsPlayedOfHouseThisTurn(player int, house House) int
	// TopOfDeck returns the top card of a player's deck without moving it,
	// reporting whether the deck holds a card.
	TopOfDeck(player int) (LocalID, bool)
	// HasAction reports whether a card has an "Action:" ability.
	HasAction(id LocalID) bool
	// Exhausted reports whether a creature is exhausted.
	Exhausted(id LocalID) bool
	// Stunned reports whether a creature is stunned.
	Stunned(id LocalID) bool
	// TimesUsedThisTurn reports how many times a creature has been used this turn.
	TimesUsedThisTurn(id LocalID) int
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
	// ForgeKey has a player forge one key at the current cost, if affordable.
	ForgeKey(player int)
	// ForgeKeyFree has a player forge one key without paying its current cost.
	ForgeKeyFree(player int)
	// GainChains adds chains to a player, penalizing their future draws.
	GainChains(controller, amount int)
}

// CreatureResolver changes the in-play state carried on a single card — its
// damage, status, Æmber, counters, house, controller, and battleline position.
type CreatureResolver interface {
	// SetDamage sets the damage on a creature (never below zero).
	SetDamage(id LocalID, amount int)
	// SetStunned sets a creature's stun status.
	SetStunned(id LocalID, stunned bool)
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
	// TakeControl moves a creature into controller's battleline without changing
	// ownership; when it later leaves play it still goes to its owner's zone. source
	// is the card whose lasting effect holds the control (reverted when it leaves play).
	TakeControl(id LocalID, controller int, source LocalID)
	// SwapBattlelinePositions exchanges two creatures' positions in the same
	// battleline without moving any state between the creatures.
	SwapBattlelinePositions(a, b LocalID)
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
	// ReapWith reaps with a creature (ability-driven, ignoring active player and
	// house). A creature can only be used while ready, so an exhausted creature may
	// be chosen but does nothing.
	ReapWith(id LocalID)
	// UseActionOf fires a card's "Action:" ability (ability-driven, ignoring active
	// player and house). A card can only be used while ready, so an exhausted card
	// may be chosen but does nothing.
	UseActionOf(id LocalID)
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
	// PutIntoDeckShuffled moves a card from play into its owner's deck and shuffles.
	PutIntoDeckShuffled(id LocalID)
	// ArchiveFromHand moves a card from its owner's hand to their archives.
	ArchiveFromHand(id LocalID)
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
	// PutFromDiscardIntoHand moves a card from its owner's discard to their hand.
	PutFromDiscardIntoHand(id LocalID)
	// MoveFromDeckToHand moves a card from its owner's deck to their hand.
	MoveFromDeckToHand(id LocalID)
	// PlayFromDeck plays a specific card from a player's deck, removing it from the
	// deck as it is played (Chaos Portal plays the card it revealed).
	PlayFromDeck(player int, id LocalID)
	// ShuffleDiscardIntoDeck moves a player's whole discard pile into their deck
	// and shuffles it.
	ShuffleDiscardIntoDeck(player int)
	// MoveFromDiscardToTopOfDeck moves a card from its owner's discard to the top
	// of their deck.
	MoveFromDiscardToTopOfDeck(id LocalID)
	// DiscardCardFromHand moves a specific card from a player's hand to their discard
	// zone.
	DiscardCardFromHand(owner int, id LocalID)
}

// TurnResolver installs turn-scoped and lasting effects: restrictions and grants
// for this or the next turn, and the "remainder of the turn" reaction/replacement
// registry that keeps such behavior out of the play and reap paths.
type TurnResolver interface {
	// CannotFightNextTurn bars a player from using creatures to fight throughout
	// their next turn.
	CannotFightNextTurn(player int)
	// GrantFightForHouse lets a player use creatures of the given house to fight
	// this turn even out of the active house.
	GrantFightForHouse(player int, house House)
	// AddLasting registers a "for the remainder of the turn" effect (Full Moon,
	// Charge!, Crystal Hive reactions; Dimension Door's replacement) on a game event,
	// instead of the effect hardcoding itself into the play or reap path.
	AddLasting(on Event, do lastingAction, controller, amount int)
	// AddLastingOnce registers a one-shot reaction that fires (and self-removes) the
	// next time its event occurs for a subject of the given house (Blypyp).
	AddLastingOnce(on Event, do lastingAction, controller, amount int, house House)
	// ForceActiveHouseNextTurn makes a player have to choose the given house as their
	// active house on their next turn.
	ForceActiveHouseNextTurn(player int, house House)
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
	// ChooseOption asks a player to pick one of several labeled options, returning
	// its index (0 when the player's chooser expresses no preference). source is the
	// card whose ability is asking (usually ctx.Source), for prompt attribution.
	ChooseOption(player int, source LocalID, prompt string, options []string) int
}

// Logger writes to the game log.
type Logger interface {
	// Logf writes a line to the game log.
	Logf(format string, args ...any)
}

// The read accessors Aember, AmberOn, Damage, Power, Name, PlayerName,
// Battleline, and Artifacts are defined in game.go; the remaining Resolver
// methods follow.

// Owner returns the player who owns a card.
func (g *Game) Owner(id LocalID) int { return g.owner(id) }

// HasTrait reports whether a card has a trait.
func (g *Game) HasTrait(id LocalID, trait Trait) bool { return g.cat.def(id).hasTrait(trait) }

// ForgeKey has a player forge one key at its current cost, if affordable.
func (g *Game) ForgeKey(player int) { g.forgeKey(player) }

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

// SetDamage sets the damage on a creature, clamped at zero.
func (g *Game) SetDamage(id LocalID, amount int) {
	if amount < 0 {
		amount = 0
	}
	g.State.Cards[id].Damage = int16(amount)
}

// SetStunned sets a creature's stun status.
func (g *Game) SetStunned(id LocalID, stunned bool) { g.State.Cards[id].Stunned = stunned }

// SetExhausted sets a creature's exhausted status.
func (g *Game) SetExhausted(id LocalID, exhausted bool) { g.State.Cards[id].Exhausted = exhausted }

// BelongToHouseForRemainderOfTurn makes a card belong to house until its
// controller's turn ends.
func (g *Game) BelongToHouseForRemainderOfTurn(id LocalID, house House) {
	g.State.Cards[id].TempHouse = house
}

// SetLastingHouse makes a card belong to house until it leaves play.
func (g *Game) SetLastingHouse(id LocalID, house House) {
	g.State.Cards[id].LastingHouse = house
}

// SetFightDamageRedirect redirects the attacker's fight damage in the current
// fight to another creature; the combat step reads and clears it.
func (g *Game) SetFightDamageRedirect(id LocalID) { g.State.FightDamageRedirect = id }

// CancelCurrentFight makes the fight in progress not occur; the combat step reads
// and clears it before Assault, Hazardous, and fight damage.
func (g *Game) CancelCurrentFight() { g.State.FightCancelled = true }

// AddAmberOn changes the Æmber sitting on a card.
func (g *Game) AddAmberOn(id LocalID, delta int) { g.State.Cards[id].Amber += int16(delta) }

// DealDamage is the Resolver entry point for the internal dealDamage.
func (g *Game) DealDamage(controller int, targets []DamageTarget) {
	g.dealDamage(controller, targets...)
}

// DestroyEach is the Resolver entry point for destroyEach.
func (g *Game) DestroyEach(controller int, ids []LocalID) { g.destroyEach(controller, ids) }

// TakeControl is the Resolver entry point for takeControl.
func (g *Game) TakeControl(id LocalID, controller int, source LocalID) {
	g.takeControl(id, controller, source)
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

// ArchiveFromHand moves a card from its owner's hand to their archives.
func (g *Game) ArchiveFromHand(id LocalID) { g.archiveFromHand(g.owner(id), id) }

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

// AddPowerCounter changes the net power counters on a creature.
func (g *Game) AddPowerCounter(id LocalID, delta int) {
	g.State.Cards[id].PowerCounters += int16(delta)
}

// PutFromDiscardIntoHand moves a card from its owner's discard pile to their hand.
func (g *Game) PutFromDiscardIntoHand(id LocalID) {
	o := g.owner(id)
	g.State.Discard[o].remove(id)
	g.State.Hand[o].add(id)
	g.logf("%s returns %s from their discard to hand", g.names[o], g.Name(id))
}

// MoveFromDeckToHand moves a card from its owner's deck to their hand.
func (g *Game) MoveFromDeckToHand(id LocalID) {
	o := g.owner(id)
	g.State.Deck[o].remove(id)
	g.State.Hand[o].add(id)
	g.logf("%s puts %s from their deck into hand", g.names[o], g.Name(id))
}

// ShuffleDiscardIntoDeck moves a player's whole discard pile into their deck and
// shuffles it.
func (g *Game) ShuffleDiscardIntoDeck(player int) { g.shuffleDiscardIntoDeck(player) }

// MoveFromDiscardToTopOfDeck moves a card from its owner's discard pile to the
// top of their deck.
func (g *Game) MoveFromDiscardToTopOfDeck(id LocalID) {
	o := g.owner(id)
	g.State.Discard[o].remove(id)
	g.State.Deck[o].addFront(id)
	g.logf("%s puts %s from their discard on top of their deck", g.names[o], g.Name(id))
}

// GainChains adds chains to a player, which reduce their draws until shed.
func (g *Game) GainChains(controller, amount int) {
	g.State.Chains[controller] += amount
	g.logf("%s gains %d %s (%d total)", g.names[controller], amount, chainNoun(amount), g.State.Chains[controller])
}

// OrderByChoice is the Resolver entry point for orderByChoice.
func (g *Game) OrderByChoice(controller int, prompt string, ids []LocalID) []LocalID {
	return g.orderByChoice(controller, prompt, ids)
}

// ChooseCreature asks a player to choose one creature from candidates, attributing
// the prompt to the source card. A sole candidate is taken automatically.
func (g *Game) ChooseCreature(player int, source LocalID, prompt string, candidates []LocalID) (LocalID, bool) {
	return g.pickCreature(player, g.sourceName(source), prompt, candidates)
}

// ChooseCard asks a player to choose one card from candidates, attributing the
// prompt to the source card. A sole candidate is taken automatically.
func (g *Game) ChooseCard(player int, source LocalID, prompt string, candidates []LocalID) (LocalID, bool) {
	return g.pickCard(player, g.sourceName(source), prompt, candidates)
}

// ChooseOption asks a player to choose one of several labeled options, attributing
// the prompt to the source card. If the player's chooser expresses no preference
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
		return oc.ChooseOption(source, prompt, options)
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

// UseActionOf fires a card's "Action:" ability, ability-driven (ignoring active
// player and house). A card can only be used while ready, so an exhausted card
// does nothing.
func (g *Game) UseActionOf(id LocalID) {
	if g.readyToUse(id) {
		g.useActionOf(id)
	}
}

// HasAction reports whether a card has an "Action:" ability.
func (g *Game) HasAction(id LocalID) bool { return g.hasTrigger(id, TriggerAction) }

// Logf writes a line to the game log.
func (g *Game) Logf(format string, args ...any) { g.logf(format, args...) }
