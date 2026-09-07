package card

import "github.com/dmikalova/vactrol/internal/engine"

// The card effect "AST" re-exported for authoring: the effect nodes you nest
// inside an ability, plus the CreatureVerb, Condition, and Count helpers some of
// them take. Each mirrors a type in the engine's effect_*.go — see there for what
// it does. They are named directly, e.g. card.DealDamage{Amount: 2, ...}.
type (
	// Effect is a card effect (see DealDamage, GainAember, ...).
	Effect = engine.Effect
	// CreatureVerb is an action applied to a chosen creature (ReadyVerb, FightVerb).
	CreatureVerb = engine.CreatureVerb
	// Condition gates a Conditional effect (see OpponentAemberAtLeast).
	Condition = engine.Condition
)

// Æmber effects.
type (
	// GainAember moves Æmber from the common supply into a player's pool.
	GainAember = engine.GainAember
	// GainAemberEqualTo gains Æmber equal to a running count (The Flex gains half a
	// chosen creature's power).
	GainAemberEqualTo = engine.GainAemberEqualTo
	// MoveAemberFromPool banks Æmber out of your pool onto a card.
	MoveAemberFromPool = engine.MoveAemberFromPool
	// PlaceAemberOnThis places Æmber from the common supply on this card.
	PlaceAemberOnThis = engine.PlaceAemberOnThis
	// LoseAember returns Æmber from a player's pool to the supply (see By: Half, AllBut).
	LoseAember = engine.LoseAember
	// StealAember moves Æmber from the opponent's pool into yours.
	StealAember = engine.StealAember
	// CaptureAember moves Æmber from a pool onto a capturing creature.
	CaptureAember = engine.CaptureAember
	// MoveAemberToSupply removes Æmber sitting on a creature to the common supply.
	MoveAemberToSupply = engine.MoveAemberToSupply
	// Exalt places Æmber from the common supply onto a chosen card.
	Exalt = engine.Exalt
	// ExaltToRepeat resolves Do, then lets the controller exalt a creature to repeat it.
	ExaltToRepeat = engine.ExaltToRepeat
	// Loss says how much Æmber a LoseAember removes (Half, AllBut).
	Loss = engine.Loss
	// MoveAember moves Æmber off a card into a pool or onto another card.
	MoveAember = engine.MoveAember
)

// Damage and combat.
type (
	// DealDamage deals damage to each creature its Target selects.
	DealDamage = engine.DealDamage
	// DamageThenIfDestroyed deals damage, then runs Then only if the creature left play.
	DamageThenIfDestroyed = engine.DamageThenIfDestroyed
	// DamageThenIfSurvives deals damage, then runs Then only if the creature survives.
	DamageThenIfSurvives = engine.DamageThenIfSurvives
	// DamageThen deals damage to a creature, then runs Then on it whether it lived or not.
	DamageThen = engine.DamageThen
	// Spread is a DealDamage strategy that hits several related creatures at once.
	Spread = engine.Spread
	// PerTarget is a DealDamage strategy scaling the damage per creature hit.
	PerTarget = engine.PerTarget
	// CreatureAndNeighbor (a Spread) damages a chosen creature and one of its neighbors.
	CreatureAndNeighbor = engine.CreatureAndNeighbor
	// CreatureAndNeighbors (a Spread) damages a chosen creature and each of its neighbors.
	CreatureAndNeighbors = engine.CreatureAndNeighbors
	// DifferentCreatures (a Spread) damages a chosen creature and a different chosen creature.
	DifferentCreatures = engine.DifferentCreatures
	// UpToCreatures (a Spread) deals damage to up to N chosen creatures.
	UpToCreatures = engine.UpToCreatures
	// DivideDamage (a Spread) divides a damage pool among any number of creatures.
	DivideDamage = engine.DivideDamage
	// FlankWalk (a Spread) deals decreasing damage inward from a chosen flank creature.
	FlankWalk = engine.FlankWalk
	// RedirectFightDamage is a Before Fight effect redirecting this creature's fight damage.
	RedirectFightDamage = engine.RedirectFightDamage
	// Heal takes damage tokens off a creature — a fixed amount or all of them.
	Heal = engine.Heal
	// LoseArmor takes all the remaining armor off each creature it targets, and
	// tallies it for a following ArmorLostThisWay.
	LoseArmor = engine.LoseArmor
	// GainStats gives each targeted creature power and/or armor for the remainder
	// of the turn (Abond the Armorsmith grants +1 armor).
	GainStats = engine.GainStats
)

// Destruction and purging.
type (
	// Destroy removes the creatures its Target selects from play.
	Destroy = engine.Destroy
	// DestroyChosen destroys any number of creatures the controller picks from its Target.
	DestroyChosen = engine.DestroyChosen
	// DestroyMostPowerfulUnlessReadyHouse destroys the most powerful creature of
	// each player who does not control a ready creature of House (Quicksand).
	DestroyMostPowerfulUnlessReadyHouse = engine.DestroyMostPowerfulUnlessReadyHouse
	// DestroyFriendlyCreaturesToForge destroys any number of friendly creatures
	// totalling a power threshold to trigger a follow-up effect (Might Makes Right
	// forges free).
	DestroyFriendlyCreaturesToForge = engine.DestroyFriendlyCreaturesToForge
	// PurgeCard sets cards aside out of the game, from a named zone.
	PurgeCard = engine.PurgeCard
	// PurgeFromHand purges one card the controller chooses from a player's hand.
	PurgeFromHand = engine.PurgeFromHand
	// PurgeEachFromHand purges every matching card from a player's hand.
	PurgeEachFromHand = engine.PurgeEachFromHand
	// PurgeEachFromDiscard purges every matching card from both discard piles.
	PurgeEachFromDiscard = engine.PurgeEachFromDiscard
	// PurgeCreature purges each creature its Target selects from play.
	PurgeCreature = engine.PurgeCreature
	// PurgeSource purges the card whose ability this is (Library Access purges itself).
	PurgeSource = engine.PurgeSource
	// LoseKeyword takes a keyword from each creature for the remainder of the turn.
	LoseKeyword = engine.LoseKeyword
	// LoseKeywords takes one or more keywords from each targeted creature for the
	// remainder of the turn (Niffle Grounds strips taunt and elusive).
	LoseKeywords = engine.LoseKeywords
	// GainKeyword gives each targeted creature a keyword until the start of your
	// next turn (Hideaway Hole grants your creatures elusive).
	GainKeyword = engine.GainKeyword
	// PurgeCreatureFromHand purges a chosen creature from your hand and puts it in context.
	PurgeCreatureFromHand = engine.PurgeCreatureFromHand
	// PurgeArchivesForDamage purges any number of cards from your archives to deal
	// damage to a creature for each card purged.
	PurgeArchivesForDamage = engine.PurgeArchivesForDamage
	// PurgeArchivedCardThen optionally purges a card from your archives to pay for
	// a follow-up effect (Yzphyz Knowdrone purges to stun a creature).
	PurgeArchivedCardThen = engine.PurgeArchivedCardThen
)

// Creature state (stun, exhaust, power counters).
type (
	// Stun places a stun on the creatures its Target selects.
	Stun = engine.Stun
	// Unstun removes the stun from each creature its Target selects.
	Unstun = engine.Unstun
	// Enrage places an enrage on the creatures its Target selects. An enraged
	// creature must be used to fight on its controller's turn if it is able to.
	Enrage = engine.Enrage
	// Ward places a one-shot shield on the creatures its Target selects; it
	// absorbs the next instance of damage or the next time the creature leaves
	// play, then is spent.
	Ward = engine.Ward
	// Exhaust turns the targeted creatures sideways so they cannot be used.
	Exhaust = engine.Exhaust
	// ExhaustCreatures exhausts up to Max creatures the controller chooses.
	ExhaustCreatures = engine.ExhaustCreatures
	// Ready turns the targeted creatures upright so they can be used again.
	Ready = engine.Ready
	// ReadyIfFirstUse readies a creature only when the current use is its first this turn.
	ReadyIfFirstUse = engine.ReadyIfFirstUse
	// ReadyCreatures readies up to Max creatures the controller chooses.
	ReadyCreatures = engine.ReadyCreatures
	// AddPowerCounter places permanent +1/-1 power counters on a creature.
	AddPowerCounter = engine.AddPowerCounter
	// PlaceCounter puts generic counters (doom and its kin) on each card its target selects.
	PlaceCounter = engine.PlaceCounter
)

// Drawing, moving, and revealing cards between zones.
type (
	// Draw puts the top Amount cards of your deck into your hand.
	Draw = engine.Draw
	// AttachSelfTo moves the resolving upgrade onto the friendly creature with the
	// given printed name (a Star Alliance "blaster" homing to its signature creature).
	AttachSelfTo = engine.AttachSelfTo
	// PutFromPlay takes each targeted card out of play into a chosen Destination.
	PutFromPlay = engine.PutFromPlay
	// PutChosen moves Amount cards the controller chooses into a Destination,
	// declinably when UpTo is set.
	PutChosen = engine.PutChosen
	// PutFromDiscard moves a chosen card from your discard pile to a Destination.
	PutFromDiscard = engine.PutFromDiscard
	// PutFromHand puts a chosen card from your hand directly into play.
	PutFromHand = engine.PutFromHand
	// ReturnNamedToHand returns a chosen card of a given name to its owner's hand.
	ReturnNamedToHand = engine.ReturnNamedToHand
	// SearchForName searches your deck and discard pile for a named card.
	SearchForName = engine.SearchForName
	// ShuffleIntoDeck shuffles the controller's named zones (hand, discard, archives) into their deck.
	ShuffleIntoDeck = engine.ShuffleIntoDeck
	// ShuffleChosenCreaturesFromDiscard shuffles any number of chosen creatures from your discard pile into your deck.
	ShuffleChosenCreaturesFromDiscard = engine.ShuffleChosenCreaturesFromDiscard
	// ShuffleChosenCreaturesFromZones shuffles any number of chosen friendly creatures from your hand, discard pile, or battleline into your deck.
	ShuffleChosenCreaturesFromZones = engine.ShuffleChosenCreaturesFromZones
	// ShuffleCardsFromDiscard shuffles a counted number of chosen cards from your discard pile into your deck.
	ShuffleCardsFromDiscard = engine.ShuffleCardsFromDiscard
	// SwapDeckAndDiscard exchanges the controller's deck with their discard pile,
	// then shuffles.
	SwapDeckAndDiscard = engine.SwapDeckAndDiscard
	// ArchiveFromHand moves cards from a hand into the controller's archives.
	ArchiveFromHand = engine.ArchiveFromHand
	// ArchiveRandomFromHand archives Amount random cards from your hand (Eureka!).
	ArchiveRandomFromHand = engine.ArchiveRandomFromHand
	// ArchiveFromDiscard moves a chosen card from the discard pile into archives.
	ArchiveFromDiscard = engine.ArchiveFromDiscard
	// ArchiveTopOfDeck moves the top Amount cards of your deck into archives.
	ArchiveTopOfDeck = engine.ArchiveTopOfDeck
	// ArchiveTopOfDiscard moves the top Amount cards of your discard pile into archives.
	ArchiveTopOfDiscard = engine.ArchiveTopOfDiscard
	// ArchiveFromPlay moves each targeted in-play card into its owner's archives.
	ArchiveFromPlay = engine.ArchiveFromPlay
	// ArchiveSource archives the card whose ability this is (Sucker Punch).
	ArchiveSource = engine.ArchiveSource
	// DiscardArchives moves all of a player's archived cards into their discard pile.
	DiscardArchives = engine.DiscardArchives
	// DiscardHand discards cards from a player's hand.
	DiscardHand = engine.DiscardHand
	// EachPlayerDiscardsAndRefillsHand makes both players discard and redraw their hand.
	EachPlayerDiscardsAndRefillsHand = engine.EachPlayerDiscardsAndRefillsHand
	// DiscardFromHand has the controller choose and discard Amount cards.
	DiscardFromHand = engine.DiscardFromHand
	// DiscardRandomFromHand discards one uniformly random card from a player's hand.
	DiscardRandomFromHand = engine.DiscardRandomFromHand
	// DiscardRandomFromArchives discards one uniformly random card from a player's archives.
	DiscardRandomFromArchives = engine.DiscardRandomFromArchives
	// DiscardTopOfDeck discards the top card of a deck and puts it in context.
	DiscardTopOfDeck = engine.DiscardTopOfDeck
	// DiscardDeckUntil discards from the top of your deck until it turns up a
	// card the filters admit, putting that card in context.
	DiscardDeckUntil = engine.DiscardDeckUntil
	// PutDiscardedIntoHand puts the card in context from the discard pile into
	// its owner's hand.
	PutDiscardedIntoHand = engine.PutDiscardedIntoHand
	// DiscardTopOfEachDeck discards the top card of each player's deck.
	DiscardTopOfEachDeck = engine.DiscardTopOfEachDeck
	// DiscardTop discards the top Amount cards of one player's deck.
	DiscardTop = engine.DiscardTop
	// ForEachDiscarded resolves Do once for each card a preceding discard removed.
	ForEachDiscarded = engine.ForEachDiscarded
	// RevealTopOfDeck reveals the top card of the controller's deck.
	RevealTopOfDeck = engine.RevealTopOfDeck
	// PlayRevealedCard plays the card a preceding reveal put in context.
	PlayRevealedCard = engine.PlayRevealedCard
	// PlayTopOfDeck plays the top card of the controller's deck outright.
	PlayTopOfDeck = engine.PlayTopOfDeck
	// LookAtTop looks at the top Amount cards of your deck, puts one into your
	// hand, and discards the others.
	LookAtTop = engine.LookAtTop
	// ReorderTop looks at the top Amount cards of your deck and puts them back in
	// any order you choose.
	ReorderTop = engine.ReorderTop
	// (From), ignoring the active house. Set Except to make House the house that
	// may not be played.
	PlayFrom = engine.PlayFrom
	// PlayRandomFromOpponentArchives plays a random card from the opponent's
	// archives as your own (a Murkens option).
	PlayRandomFromOpponentArchives = engine.PlayRandomFromOpponentArchives
	// PlayTopOfOpponentDeck plays the top card of the opponent's deck as your own
	// (a Murkens option).
	PlayTopOfOpponentDeck = engine.PlayTopOfOpponentDeck
	// PutUnderFromHand puts a card the controller chooses from their hand under
	// the resolving card, face up or face down.
	PutUnderFromHand = engine.PutUnderFromHand
	// PlayCardUnder plays the card placed under the resolving card.
	PlayCardUnder = engine.PlayCardUnder
	// Graft moves a target card in play faceup under the resolving card, out of
	// play (rulebook: Graft).
	Graft = engine.Graft
	// PutUnderIntoPlay puts every card under the resolving card into play under
	// its owner's control.
	PutUnderIntoPlay = engine.PutUnderIntoPlay
	// ArchiveCardUnder archives the card placed under the resolving card.
	ArchiveCardUnder = engine.ArchiveCardUnder
	// CancelFight makes the fight in progress not occur (a Before Fight effect).
	CancelFight = engine.CancelFight
	// RevealHand shows the cards in a player's hand to both players and records them.
	RevealHand = engine.RevealHand
)

// Using and choosing creatures.
type (
	// OnChooseCreature picks a creature named by its Target and applies Verbs to it.
	OnChooseCreature = engine.OnChooseCreature
	// ChooseCreatureThen asks the controller to choose a creature, then resolves
	// Then unconditionally (unlike Then, a result gate).
	ChooseCreatureThen = engine.ChooseCreatureThen
	// OneAtATime repeats a chosen-creature action over several different creatures,
	// resolving each pass fully before offering the next choice.
	OneAtATime = engine.OneAtATime
	// RepeatedFight readies and fights with a creature several times, each fight
	// against a different enemy creature.
	RepeatedFight = engine.RepeatedFight
	// ReadyVerb readies the chosen creature.
	ReadyVerb = engine.ReadyVerb
	// FightVerb makes the chosen creature fight an enemy creature.
	FightVerb = engine.FightVerb
	// UseVerb uses the chosen creature (reap, fight, or Action:).
	UseVerb = engine.UseVerb
	// StunVerb stuns the chosen creature.
	StunVerb = engine.StunVerb
	// ExhaustVerb exhausts the chosen creature.
	ExhaustVerb = engine.ExhaustVerb
	// GainKeywordVerb gives the chosen creature a keyword for the remainder of the
	// turn (Scout grants Skirmish).
	GainKeywordVerb = engine.GainKeywordVerb
	// ConsiderFlank makes the chosen creature count as a flank creature for the
	// remainder of the turn (Spectral Tunneler).
	ConsiderFlank = engine.ConsiderFlank
	// GainAbility grants the chosen creature a triggered ability for the remainder
	// of the turn (Spectral Tunneler grants "Reap: Draw a card").
	GainAbility = engine.GainAbility
	// TakesExtraDamage makes the chosen creature take additional damage whenever it
	// takes damage, for the remainder of the turn (Lethal Distraction).
	TakesExtraDamage = engine.TakesExtraDamage
	// RedistributeCapturedAember moves all the Æmber on one side's creatures back
	// among that side's creatures however the controller chooses (Equalize).
	RedistributeCapturedAember = engine.RedistributeCapturedAember
	// RedistributeDamage lets the controller choose a player and optionally move
	// all the damage on that player's creatures among that player's creatures
	// (Entropic Manipulator).
	RedistributeDamage = engine.RedistributeDamage
	// Use uses up to Max cards the controller chooses from Target.
	Use = engine.Use
	// TriggerAbility fires another card's ability as if the controller controlled it.
	TriggerAbility = engine.TriggerAbility
	// TakeControl moves a card to the controller's play area and makes them its controller.
	TakeControl = engine.TakeControl
	// PutIntoPlay puts each targeted card into play without playing it.
	PutIntoPlay = engine.PutIntoPlay
	// Control names whose control a card enters under when put into play (see
	// PutIntoPlay): card.Owner (the default) or card.Yours.
	Control = engine.Control
	// Swap exchanges this creature's battleline position with another.
	Swap = engine.Swap
	// SwapChosen swaps the positions of two creatures chosen from one battleline.
	SwapChosen = engine.SwapChosen
	// MoveToFlank moves the targeted creature to either flank of its controller's battleline.
	MoveToFlank = engine.MoveToFlank
)

// Composites and control flow.
type (
	// Sequence resolves several effects in order.
	Sequence = engine.Sequence
	// ForDuration applies several timed effects sharing one duration and
	// renders their shared "for the remainder of the turn, ..." clause once.
	ForDuration = engine.ForDuration
	// Repeat resolves an effect once for each of a running count, choosing
	// afresh each time.
	Repeat = engine.Repeat
	// Sentences resolves several effects in order, each rendered as its own
	// sentence rather than joined with ", and".
	Sentences = engine.Sentences
	// ChooseOne offers the controller a set of alternative effects to pick from.
	ChooseOne = engine.ChooseOne
	// ChooseHouseThen asks the controller to choose a house, then resolves Then.
	ChooseHouseThen = engine.ChooseHouseThen
	// Conditional resolves Then only when Cond is met.
	Conditional = engine.Conditional
	// RepeatWhile resolves Do again and again while Cond holds.
	RepeatWhile = engine.RepeatWhile
	// RepeatOnCondition resolves Do and repeats it while it succeeds and Cond holds.
	RepeatOnCondition = engine.RepeatOnCondition
	// MayRepeat resolves Do, then lets the controller repeat it.
	MayRepeat = engine.MayRepeat
	// May makes an effect optional — the controller chooses whether to resolve it.
	May = engine.May
	// Then is the A -> B result gate: resolves Result only when First did something.
	Then = engine.Then
	// OrAmount switches an effect's amount to Amount when When holds, printing the
	// linear "<base>, or <alt> if <cond>" form instead of an Otherwise branch (rule
	// 22). Used as StealAember{Amount: 1, Or: card.OrAmount{Amount: 2, When: ...}}.
	OrAmount = engine.OrAmount
)

// Conditions gate a Conditional, RepeatWhile, or MayRepeat.
type (
	// OpponentAember gates on the opponent's Æmber pool (Is + Amount).
	OpponentAember = engine.OpponentAember
	// PlayerControlsFewerHousesThan is met while a player controls creatures from fewer than Amount houses.
	PlayerControlsFewerHousesThan = engine.PlayerControlsFewerHousesThan
	// YourAember gates on the controller's own Æmber pool (Is + Amount).
	YourAember = engine.YourAember
	// CardsDestroyedFewerThan is met when fewer than Amount cards were destroyed this way.
	CardsDestroyedFewerThan = engine.CardsDestroyedFewerThan
	// CountIs gates on any Count compared against a threshold (Count + Is + Amount).
	CountIs = engine.CountIs
	// ControlsMoreCreatures is met while you control more creatures than the opponent.
	ControlsMoreCreatures = engine.ControlsMoreCreatures
	// SourceOnFlank gates on the source card's flank position (Not inverts it).
	SourceOnFlank = engine.SourceOnFlank
	// SourceInCenterOfBattleline is met while the source card sits in the center
	// of its controller's battleline (an even-sized line has no center).
	SourceInCenterOfBattleline = engine.SourceInCenterOfBattleline
	// SourceReady is met while the source card is ready (Bellowing Patrizate's gate).
	SourceReady = engine.SourceReady
	// SourceNeighborsAllOfHouse is met while every neighbor of the source card
	// belongs to House (Xanthyx Harvester's use gate).
	SourceNeighborsAllOfHouse = engine.SourceNeighborsAllOfHouse
	// ControlsCreaturesOfHouses is met while your creatures span at least Amount houses.
	ControlsCreaturesOfHouses = engine.ControlsCreaturesOfHouses
	// FirstCreaturePlayedThisTurn is met when the card in context is the first
	// creature played this turn — a once-per-turn charge (Speed Sigil).
	FirstCreaturePlayedThisTurn = engine.FirstCreaturePlayedThisTurn
	// NoCreaturesPlayedThisTurn is met when you played no creatures this turn (Redlock).
	NoCreaturesPlayedThisTurn = engine.NoCreaturesPlayedThisTurn
	// ItIsYourTurn is met when the ability's controller is the active player.
	ItIsYourTurn = engine.ItIsYourTurn
	// AemberOnThisAtLeast is met when at least Amount Æmber sits on this card.
	AemberOnThisAtLeast = engine.AemberOnThisAtLeast
	// Overwhelmed is met while the opponent controls more creatures than you.
	Overwhelmed = engine.Overwhelmed
	// ItIsOfHouse is met when the card in context belongs to a referenced house.
	ItIsOfHouse = engine.ItIsOfHouse
	// ItIs is met when the card in context matches a concrete House and/or Type.
	ItIs = engine.ItIs
	// ItIsOffIdentity is met when the card in context is off your identity houses.
	ItIsOffIdentity = engine.ItIsOffIdentity
	// ItIsStunned is met when the creature in context is already stunned.
	ItIsStunned = engine.ItIsStunned
	// ChoseHouse is met when the controller's active house is House.
	ChoseHouse = engine.ChoseHouse
)

// Æmber-pool comparisons for card.OpponentAember{Is: ..., Amount: n}.
var (
	// AtLeast is met when the quantity is at least Amount.
	AtLeast = engine.AtLeast
	// AtMost is met when the quantity is at most Amount.
	AtMost = engine.AtMost
	// Exactly is met when the quantity is exactly Amount.
	Exactly = engine.Exactly
	// MoreThanYou is met when the opponent's pool holds more Æmber than yours.
	MoreThanYou = engine.MoreThanYou
	// MoreThanOpponent is met when your pool holds more Æmber than the opponent's.
	MoreThanOpponent = engine.MoreThanOpponent
)

// House references for conditions that compare a card's house dynamically.
var (
	// TheChosenHouse is the house picked by an enclosing ChooseHouseThen.
	TheChosenHouse = engine.TheChosenHouse
	// TheActiveHouse is the player's active house this turn.
	TheActiveHouse = engine.TheActiveHouse
	// TheContextualHouse is the house of the card in context (ctx.It).
	TheContextualHouse = engine.TheContextualHouse
	// AnyHouse applies no house filter at all.
	AnyHouse = engine.AnyHouse
)

// Counts feed an effect's Per, scaling it by a board quantity; InPlay doubles as
// a Condition.
type (
	// Fixed is a Count of a constant number, for a Times that does not scale with
	// the board (RepeatedFight fights a fixed number of times).
	Fixed = engine.Fixed
	// InPlay counts (or gates on) the cards a player has in play matching its filters.
	InPlay = engine.InPlay
	// CardsPlayed counts the cards of a house a player has played this turn.
	CardsPlayed = engine.CardsPlayed
	// CreaturesUsed counts the creatures a player has used this turn.
	CreaturesUsed = engine.CreaturesUsed
	// CardsDiscarded is a Condition met when a player has discarded cards of a house this turn.
	CardsDiscarded = engine.CardsDiscarded
	// OpponentForgedKeys counts the keys the opponent has forged.
	OpponentForgedKeys = engine.OpponentForgedKeys
	// TurnCount counts one of the engine's turn-history tallies (Player + Of).
	TurnCount = engine.TurnCount
	// ForgedKey gates on whether a player forged a key this turn or their previous one.
	ForgedKey = engine.ForgedKey
	// EnemyCreatureDestroyed is met once an enemy creature has been destroyed this turn.
	EnemyCreatureDestroyed = engine.EnemyCreatureDestroyed
	// UsedCreatureToReap is met once you have used a creature to reap this turn.
	UsedCreatureToReap = engine.UsedCreatureToReap
	// UsedCreatureToFight is met once you have used a creature to fight this turn.
	UsedCreatureToFight = engine.UsedCreatureToFight
	// FirstReapOfTurn is met when the reap in context is the first this turn.
	FirstReapOfTurn = engine.FirstReapOfTurn
	// CounterInPlay is met while at least one card in play carries a generic counter of the given kind.
	CounterInPlay = engine.CounterInPlay
	// ExcessCreatures counts how many more creatures one player controls than the other.
	ExcessCreatures = engine.ExcessCreatures
	// CardsInArchives counts the cards in a player's archives.
	CardsInArchives = engine.CardsInArchives
	// HousesInPlay counts the distinct houses among all cards in play, optionally
	// excluding one house.
	HousesInPlay = engine.HousesInPlay
	// HousesAmong counts the distinct houses among a player's creatures, all
	// creatures, or all cards in play, optionally capped by Max.
	HousesAmong = engine.HousesAmong
	// HousesRepresented is met when the houses among a chosen set of cards compare to
	// an amount.
	HousesRepresented = engine.HousesRepresented
	// CardsRevealed counts the cards the most recent Reveal showed.
	CardsRevealed = engine.CardsRevealed
	// CardsDestroyed counts the cards the most recent destruction removed "this way".
	CardsDestroyed = engine.CardsDestroyed
	// CreaturesDestroyed counts the cards the most recent destruction removed "this
	// way", rendered as creatures. Use it when only creatures can be destroyed;
	// use CardsDestroyed when artifacts can be too.
	CreaturesDestroyed = engine.CreaturesDestroyed
	// AemberBonusDestroyed counts the total Æmber pips on the cards the most recent
	// destruction removed (Rustgnawer gains the destroyed artifact's Æmber bonus).
	AemberBonusDestroyed = engine.AemberBonusDestroyed
	// CardsPurged counts the creatures the most recent purge removed "this way".
	CardsPurged = engine.CardsPurged

	// CreaturesDestroyedThisWay counts the creatures its Player controlled that an
	// earlier effect in this resolution destroyed. Under an EachPlayer effect,
	// Player: Controller is each player's own dead; use CardsDestroyed for the
	// whole tally.
	CreaturesDestroyedThisWay = engine.CreaturesDestroyedThisWay

	// CreaturesShuffledIntoDeckThisWay counts the creatures its Player controlled
	// that an earlier effect in this resolution put back into a deck.
	CreaturesShuffledIntoDeckThisWay = engine.CreaturesShuffledIntoDeckThisWay

	// AemberLostThisWay counts the Æmber an earlier LoseAember in the same
	// resolution took from its Player's pool.
	AemberLostThisWay = engine.AemberLostThisWay
	// AemberInPool counts the Æmber currently in a player's pool.
	AemberInPool = engine.AemberInPool
	// NeighborsOfThis counts the battleline neighbors of the source creature (0-2).
	NeighborsOfThis = engine.NeighborsOfThis
	// CardsReturnedThisWay counts the cards an earlier PutFromDiscard in the same
	// resolution recovered from the discard pile.
	CardsReturnedThisWay = engine.CardsReturnedThisWay
	// CardsInHand counts the cards in a player's hand of a referenced house.
	CardsInHand = engine.CardsInHand
	// CreaturesHealed counts the creatures the most recent Heal healed.
	CreaturesHealed = engine.CreaturesHealed
	// DamageHealed counts the damage the most recent Heal removed (for DealDamage.AmountFrom).
	DamageHealed = engine.DamageHealed
	// DamagePrevented counts the damage a creature just prevented with its own armor.
	DamagePrevented = engine.DamagePrevented
	// UnforgedKeys counts the keys a player has still to forge.
	UnforgedKeys = engine.UnforgedKeys
	// AemberOnThis counts the Æmber sitting on the source card.
	AemberOnThis = engine.AemberOnThis
	// DamageOnThis counts the damage sitting on the source card.
	DamageOnThis = engine.DamageOnThis
	// PowerOfChosen is the power of the creature in context, optionally a fraction
	// of it (Of: Half — The Flex).
	PowerOfChosen = engine.PowerOfChosen
	// TraitsOfChosen counts the traits of the creature just chosen.
	TraitsOfChosen = engine.TraitsOfChosen
	// CopiesInDiscard counts the copies of this card in your discard pile.
	CopiesInDiscard = engine.CopiesInDiscard
)

// Lasting "for the remainder of the turn" effects.
type (
	// ForRemainderOfTurn installs a reaction that runs for the rest of your turn.
	ForRemainderOfTurn = engine.ForRemainderOfTurn
	// Instead installs a replacement that changes an event's outcome for the turn.
	Instead = engine.Instead
	// Replace is a continuous replacement an Upgrade applies to a game event.
	Replace = engine.Replace
	// NextPlayed makes the next creature of a house you play do something.
	NextPlayed = engine.NextPlayed
)

// Houses, keys, chains, and restrictions.
type (
	// CannotFight bars a player from using creatures to fight for a Duration.
	CannotFight = engine.CannotFight
	// CannotPlay bars a player from playing cards of a Type for a Duration.
	CannotPlay = engine.CannotPlay
	// CannotUse bars a player from reaping, fighting, or using Action: abilities.
	CannotUse = engine.CannotUse
	// CannotReap bars a player from using creatures to reap for a Duration.
	CannotReap = engine.CannotReap
	// ChosenHouseCannotReapNextTurn bars a player from reaping with creatures of the
	// chosen house throughout their next turn.
	ChosenHouseCannotReapNextTurn = engine.ChosenHouseCannotReapNextTurn
	// BlankEnemyText blanks enemy creatures' text boxes until your next turn (Shadow of Dis).
	BlankEnemyText = engine.BlankEnemyText
	// SkipForgePhase makes a player skip their forge-a-key phase next turn.
	SkipForgePhase = engine.SkipForgePhase
	// CannotBeDealtDamage marks the targeted creatures unable to be dealt damage
	// for a Duration.
	CannotBeDealtDamage = engine.CannotBeDealtDamage
	// MayUseFriendlyHouse lets the controller use their House creatures this turn.
	MayUseFriendlyHouse = engine.MayUseFriendlyHouse
	// MayUseFriendlyArtifacts lets the controller use any friendly artifact this turn.
	MayUseFriendlyArtifacts = engine.MayUseFriendlyArtifacts
	// MayPlayOrUseFriendlyHouse lets the controller play and use a House this turn.
	MayPlayOrUseFriendlyHouse = engine.MayPlayOrUseFriendlyHouse
	// GrantFightForChosenHouse lets your chosen-house creatures fight this turn.
	GrantFightForChosenHouse = engine.GrantFightForChosenHouse
	// GrantFightForFriendlyHouse lets your creatures of a named House fight this turn.
	GrantFightForFriendlyHouse = engine.GrantFightForFriendlyHouse
	// GrantFightAnyHouse lets every friendly creature fight this turn.
	GrantFightAnyHouse = engine.GrantFightAnyHouse
	// BelongToHouse makes the targeted creatures belong to a House for a Duration.
	BelongToHouse = engine.BelongToHouse
	// NameHouse remembers the house an enclosing ChooseHouseThen picked on this card,
	// feeding the card's HouseLock for as long as it stays in play.
	NameHouse = engine.NameHouse
	// ForceOpponentActiveHouse forces the opponent's active house next turn.
	ForceOpponentActiveHouse = engine.ForceOpponentActiveHouse
	// ForbidOpponentActiveHouse bars the opponent's chosen house next turn (Tezmal).
	ForbidOpponentActiveHouse = engine.ForbidOpponentActiveHouse
	// ForgeKey has the controller forge a key outside the normal step.
	ForgeKey = engine.ForgeKey
	// UnforgeKey takes a forged key back off a player (Key Hammer).
	UnforgeKey = engine.UnforgeKey
	// RaiseKeyCost makes keys cost more throughout a player's next turn.
	RaiseKeyCost = engine.RaiseKeyCost
	// GiveRemainingAemberAfterOpponentForgeKey arms Interdimensional Graft's delayed gift.
	GiveRemainingAemberAfterOpponentForgeKey = engine.GiveRemainingAemberAfterOpponentForgeKey
	// GainChains gives a player chains (a draw penalty).
	GainChains = engine.GainChains
)

// Counter groups the generic counter kinds a card can place or read, e.g.
// card.PlaceCounter{Kind: card.Counter.Doom, Target: card.Target.ChosenCreature}.
var Counter = counters{
	Doom: engine.CounterDoom,
}

type counters struct {
	Doom engine.CounterKind
}

// Event groups the game events a lasting "for the remainder of the turn" effect
// attaches to (see ForRemainderOfTurn and Instead), e.g.
// card.ForRemainderOfTurn{On: card.Event.CreaturePlayed, Do: card.GainAember{...}}.
var Event = events{CreaturePlayed: engine.EventCreaturePlayed,
	Reap:                   engine.EventReap,
	Fight:                  engine.EventFight,
	EnemyCreatureDestroyed: engine.EventEnemyCreatureDestroyed,
	ReapAember:             engine.EventReapAember,
	Destroyed:              engine.EventCreatureDestroyed,
	AemberAddedToPool:      engine.EventAemberAddedToPool,
	AemberTakenFromPool:    engine.EventAemberTakenFromPool,
	CardPlayed:             engine.EventCardPlayed,
}

type events struct {
	CreaturePlayed,
	Reap,
	Fight,
	EnemyCreatureDestroyed,
	ReapAember,
	Destroyed,
	AemberAddedToPool,
	AemberTakenFromPool,
	CardPlayed engine.Event
}

// Steal is the replacement that makes gaining Æmber steal it from the opponent
// instead, for card.Instead (Dimension Door).
var Steal = engine.Steal

// Capture is the replacement that makes Æmber added to the opponent's pool be
// captured by the source creature instead, for card.WithReplaces (Ether Spider).
var Capture = engine.Capture

// FromCommonSupply is the replacement that draws a steal or capture from the
// common supply instead of the target's pool, for card.WithReplaces (Po's Pixies).
var FromCommonSupply = engine.FromCommonSupply

// Owner puts a card into play under its owner's control — the default for
// PutIntoPlay.Control, usually omitted.
var Owner = engine.ControlOwner

// Yours puts a card into play under the resolving player's control:
// card.PutIntoPlay{Control: card.Yours} (Overlord Greking).
var Yours = engine.ControlYours

// Half is the Loss that makes a LoseAember remove half the pool, rounded down:
// card.LoseAember{Player: card.EachPlayer, By: card.Half}.
var Half = engine.Half

// AllBut is the Loss that makes a LoseAember reduce a pool to keep, removing
// everything above it: card.LoseAember{Player: card.EachPlayer, By: card.AllBut(5)}.
var AllBut = engine.AllBut

// AllAember is the Loss that makes a LoseAember empty a pool entirely:
// card.LoseAember{Player: card.Controller, By: card.AllAember}.
var AllAember = engine.AllAember

// AemberOnIt is the PerTarget that scales damage by the Æmber on each creature
// hit: card.DealDamage{Amount: 1, Target: ..., PerTarget: card.AemberOnIt}.
var AemberOnIt = engine.AemberOnIt

// DamageOnIt is the PerTarget that scales damage by the damage already on each
// creature hit (Cauldron Boil).
var DamageOnIt = engine.DamageOnIt

// ArmorLostThisWay is the PerTarget that scales damage by the armor an effect has
// stripped off each creature hit (Red-Hot Armor).
var ArmorLostThisWay = engine.ArmorLostThisWay
