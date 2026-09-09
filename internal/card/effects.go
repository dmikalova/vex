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
	// Condition gates a Conditional effect (see PoolAember).
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
	// LoseAemberEqualTo loses Æmber equal to a running count (Power of Fire loses
	// half the sacrificed creature's power).
	LoseAemberEqualTo = engine.LoseAemberEqualTo
	// StealAember moves Æmber from the opponent's pool into yours.
	StealAember = engine.StealAember
	// CaptureAember moves Æmber from a pool onto a capturing creature.
	CaptureAember = engine.CaptureAember
	// CaptureFromAnyPlayer captures Æmber onto this creature from both pools in any split.
	CaptureFromAnyPlayer = engine.CaptureFromAnyPlayer
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
	// DealDamagePerHouse deals damage to one chosen creature of each house
	// (Gleeful Mayhem).
	DealDamagePerHouse = engine.DealDamagePerHouse
	// DamageThen deals damage to a creature, then runs Then on it under After
	// (card.Always / card.IfDestroyed / card.IfSurvives).
	DamageThen = engine.DamageThen
	// DamageAftermath is the After axis of DamageThen (see card.IfDestroyed).
	DamageAftermath = engine.DamageAftermath
	// Spread is a DealDamage strategy that hits several related creatures at once.
	Spread = engine.Spread
	// PerTarget is a DealDamage strategy scaling the damage per creature hit.
	PerTarget = engine.PerTarget
	// CreatureAndNeighbors (a Spread) damages a chosen creature and its neighbors —
	// each of them, or one chosen when Scope is card.OneNeighbor.
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

// After branches for card.DamageThen.
const (
	// Always runs the DamageThen follow-up whether or not the creature was destroyed.
	Always = engine.Always
	// IfDestroyed runs the follow-up only if the damage destroyed the creature.
	IfDestroyed = engine.IfDestroyed
	// IfSurvives runs the follow-up only if the creature survives the damage.
	IfSurvives = engine.IfSurvives
)

// Destruction and purging.
type (
	// Destroy removes the creatures its Target selects from play.
	Destroy = engine.Destroy
	// DestroyChosen destroys any number of creatures the controller picks from its Target.
	DestroyChosen = engine.DestroyChosen
	// DestroyAllExceptChosen keeps a chosen number of friendly and enemy creatures and destroys every other creature.
	DestroyAllExceptChosen = engine.DestroyAllExceptChosen
	// Tertiate destroys one third of all enemy creatures and one third of all friendly creatures (rounding up each time).
	Tertiate = engine.Tertiate
	// DestroyMostPowerfulUnlessReadyHouse destroys the most powerful creature of
	// each player who does not control a ready creature of House (Quicksand).
	DestroyMostPowerfulUnlessReadyHouse = engine.DestroyMostPowerfulUnlessReadyHouse
	// DestroyEachCreatureAtEndOfTurn schedules "destroy each creature" to resolve in
	// the end-of-turn phase rather than now (Ragnarok).
	DestroyEachCreatureAtEndOfTurn = engine.DestroyEachCreatureAtEndOfTurn
	// DestroyFriendlyCreaturesToForge destroys any number of friendly creatures
	// totalling a power threshold to trigger a follow-up effect (Might Makes Right
	// forges free).
	DestroyFriendlyCreaturesToForge = engine.DestroyFriendlyCreaturesToForge
	// SacrificeToForge destroys any number of friendly creatures then may forge a
	// key at a surcharge reduced per creature destroyed, destroying the source
	// artifact on forge (Obsidian Forge).
	SacrificeToForge = engine.SacrificeToForge
	// PurgeCard sets cards aside out of the game, from a named zone.
	PurgeCard = engine.PurgeCard
	// PurgeFromHand purges one card the controller chooses from a player's hand.
	PurgeFromHand = engine.PurgeFromHand
	// PurgeEachOfChosenTrait purges every card of a chosen trait, paying each player
	// for their losses (Harvest Time).
	PurgeEachOfChosenTrait = engine.PurgeEachOfChosenTrait
	// PurgeRandomFromHand purges one uniformly random card from a player's hand.
	PurgeRandomFromHand = engine.PurgeRandomFromHand
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
	// MoveWard takes the ward off one warded creature and places it on another (Hunter or Hunted?).
	MoveWard = engine.MoveWard
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
	// RemoveCounters takes every generic counter of one kind off each targeted card.
	RemoveCounters = engine.RemoveCounters
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
	// SearchDeck searches your deck for a card (any card, or one of a given house),
	// puts it into your hand, and shuffles your deck.
	SearchDeck = engine.SearchDeck
	// ShuffleIntoDeck shuffles the controller's named zones (hand, discard, archives) into their deck.
	ShuffleIntoDeck = engine.ShuffleIntoDeck
	// ShuffleChosenCreaturesFromDiscard shuffles any number of chosen creatures from your discard pile into your deck.
	ShuffleChosenCreaturesFromDiscard = engine.ShuffleChosenCreaturesFromDiscard
	// ShuffleChosenCreaturesFromZones shuffles any number of chosen friendly creatures from your hand, discard pile, or battleline into your deck.
	ShuffleChosenCreaturesFromZones = engine.ShuffleChosenCreaturesFromZones
	// ShuffleCardsFromDiscard shuffles a counted number of chosen cards from your discard pile into your deck.
	ShuffleCardsFromDiscard = engine.ShuffleCardsFromDiscard
	// ShuffleMatchingFromDiscardIntoDeck shuffles each matching card from your
	// discard pile into your deck.
	ShuffleMatchingFromDiscardIntoDeck = engine.ShuffleMatchingFromDiscardIntoDeck
	// ShuffleFriendlyCardsInPlayIntoDeck shuffles every friendly card in play into your deck, then draws a card for each shuffled this way.
	ShuffleFriendlyCardsInPlayIntoDeck = engine.ShuffleFriendlyCardsInPlayIntoDeck
	// SwapDeckAndDiscard exchanges the controller's deck with their discard pile,
	// then shuffles.
	SwapDeckAndDiscard = engine.SwapDeckAndDiscard
	// ArchiveFromHand moves cards from a hand into the controller's archives.
	ArchiveFromHand = engine.ArchiveFromHand
	// ArchiveRandomFromHand archives Amount random cards from your hand (Eureka!).
	ArchiveRandomFromHand = engine.ArchiveRandomFromHand
	// ArchiveFromDiscard moves a chosen card from the discard pile into archives.
	ArchiveFromDiscard = engine.ArchiveFromDiscard
	// ArchiveTop moves the top Amount cards of a pile named by From (card.Deck or
	// card.Discard) into archives.
	ArchiveTop = engine.ArchiveTop
	// ArchiveFromPlay moves each targeted in-play card into its owner's archives.
	ArchiveFromPlay = engine.ArchiveFromPlay
	// ArchiveSource archives the card whose ability this is (Sucker Punch).
	ArchiveSource = engine.ArchiveSource
	// ArchiveGrantingUpgrade archives the upgrade whose granted ability this is
	// (Ghostform archives itself off its host).
	ArchiveGrantingUpgrade = engine.ArchiveGrantingUpgrade
	// ArchivePurgedCard archives a card the controller chooses from their own purge pile (Universal Recycle Bin).
	ArchivePurgedCard = engine.ArchivePurgedCard
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
	// RevealDeckUntilHouse reveals and archives cards from the top of your deck
	// until you reveal a card of a house or choose to stop, reporting whether one
	// was revealed.
	RevealDeckUntilHouse = engine.RevealDeckUntilHouse
	// RevealTopOfDeck reveals the top card of the controller's deck.
	RevealTopOfDeck = engine.RevealTopOfDeck
	// RevealPurgeShuffleDeck reveals the top cards of a chosen player's deck, purges
	// one revealed card, then shuffles that deck (Borr Nit).
	RevealPurgeShuffleDeck = engine.RevealPurgeShuffleDeck
	// PlayRevealedCard plays the card a preceding reveal put in context.
	PlayRevealedCard = engine.PlayRevealedCard
	// PlayTopOfDeck plays the top card of the controller's deck outright.
	PlayTopOfDeck = engine.PlayTopOfDeck
	// LookAtTop looks at the top Amount cards of your deck, puts one into your
	// hand, and discards the others.
	LookAtTop = engine.LookAtTop
	// LookAtTopSort looks at the top 3 cards of your deck, archiving one, putting
	// one into your hand, and discarding one.
	LookAtTopSort = engine.LookAtTopSort
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
	// PlayDiscardedTacticFromOpponent is Fidgit's reap: discard a random card from
	// the opponent's archives or their deck top, then play it as your own if it is a
	// Tactic.
	PlayDiscardedTacticFromOpponent = engine.PlayDiscardedTacticFromOpponent
	// PutUnderFromHand puts a card the controller chooses from their hand under
	// the resolving card, face up or face down.
	PutUnderFromHand = engine.PutUnderFromHand
	// PlayCardUnder plays the card placed under the resolving card.
	PlayCardUnder = engine.PlayCardUnder
	// TriggerGraftedPlayEffect triggers the play effect of an action card grafted
	// faceup under the resolving card, leaving it grafted (Infomancer, Memolith).
	TriggerGraftedPlayEffect = engine.TriggerGraftedPlayEffect
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
	// ReapVerb makes the chosen creature reap.
	ReapVerb = engine.ReapVerb
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
	// RearrangeBattleline reorders one player's battleline by swapping pairs of creatures.
	RearrangeBattleline = engine.RearrangeBattleline
	// MoveToFlank moves the targeted creature to either flank of its controller's battleline.
	MoveToFlank = engine.MoveToFlank // MoveWithinBattleline repositions the targeted creature anywhere in its
	// controller's battleline and leaves it in context (Malison).
	MoveWithinBattleline = engine.MoveWithinBattleline
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
	Sentences = engine.Sentences // ChooseOne offers the controller a set of alternative effects to pick from.
	ChooseOne = engine.ChooseOne
	// ChooseHouseThen asks the controller to choose a house, then resolves Then.
	ChooseHouseThen = engine.ChooseHouseThen
	// Conditional resolves Then only when Cond is met.
	Conditional = engine.Conditional
	// SaveFromDestruction is a creature's own "Destroyed:" replacement: it stays in
	// play and Do resolves on it instead of being destroyed.
	SaveFromDestruction = engine.SaveFromDestruction
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
	// PoolAember gates on one player's Æmber pool (Player + Is + Amount).
	PoolAember = engine.PoolAember
	// PlayerControlsFewerHousesThan is met while a player controls creatures from fewer than Amount houses.
	PlayerControlsFewerHousesThan = engine.PlayerControlsFewerHousesThan
	// CardsDestroyedFewerThan is met when fewer than Amount cards were destroyed this way.
	CardsDestroyedFewerThan = engine.CardsDestroyedFewerThan
	// CountIs gates on any Count compared against a threshold (Count + Is + Amount).
	CountIs = engine.CountIs
	// ControlsMoreCreatures is met while you control more creatures than the opponent.
	ControlsMoreCreatures = engine.ControlsMoreCreatures
	// SourceOnFlank gates on the source card's flank position (Not inverts it).
	SourceOnFlank = engine.SourceOnFlank
	// HasOtherFriendlyCreatures is met when the controller has any creature in play
	// besides the source.
	HasOtherFriendlyCreatures = engine.HasOtherFriendlyCreatures
	// CardsInDeckAtMost is met when the controller's deck holds at most Amount cards.
	CardsInDeckAtMost = engine.CardsInDeckAtMost
	// CardsInDiscardAtLeast is met when the controller's discard pile holds at least
	// Amount cards matching House and Type.
	CardsInDiscardAtLeast = engine.CardsInDiscardAtLeast
	// Haunted is met while the controller has 10 or more cards in their discard pile.
	Haunted = engine.Haunted
	// ControlsNamed is met when the controller has a card of a given name in play.
	ControlsNamed = engine.ControlsNamed // ItIsOnFlank gates on whether the context creature (ctx.It) is on a flank.
	ItIsOnFlank   = engine.ItIsOnFlank
	// ItIsOnNamedFlank gates on whether the context creature is on the left flank,
	// or the right when Right is set (Sinestra, Dexus).
	ItIsOnNamedFlank = engine.ItIsOnNamedFlank
	// SourceInCenterOfBattleline is met while the source card sits in the center
	// of its controller's battleline (an even-sized line has no center).
	SourceInCenterOfBattleline = engine.SourceInCenterOfBattleline
	// SourceReady is met while the source card is ready (Bellowing Patrizate's gate).
	SourceReady = engine.SourceReady
	// SourceNeighborsAllOfHouse is met while every neighbor of the source card
	// belongs to House (Xanthyx Harvester's use gate).
	SourceNeighborsAllOfHouse = engine.SourceNeighborsAllOfHouse
	// ArchivedCreaturesShareHouse is met when the creatures a preceding
	// ArchiveFromPlay set aside all share one house (Code Monkey's payoff).
	ArchivedCreaturesShareHouse = engine.ArchivedCreaturesShareHouse
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
	// ItIsFriendly is met when the card in context is controlled by you.
	ItIsFriendly = engine.ItIsFriendly
	// ItIs is met when the card in context matches a concrete House and/or Type.
	ItIs = engine.ItIs
	// ItIsOfTrait is met when the creature in context has the named trait.
	ItIsOfTrait = engine.ItIsOfTrait
	// ItHasAember is met when the creature in context has Æmber on it.
	ItHasAember = engine.ItHasAember
	// Or is met when any of its Conditions is met, composing conditions (e.g. a
	// Dinosaur creature or one with Æmber, for Guji Dinosaur Hunter).
	Or = engine.Or
	// ItIsOffIdentity is met when the card in context is off your identity houses.
	ItIsOffIdentity = engine.ItIsOffIdentity
	// ItIsStunned is met when the creature in context is already stunned.
	ItIsStunned = engine.ItIsStunned
	// ChoseHouse is met when the controller's active house is House.
	ChoseHouse = engine.ChoseHouse
	// ActiveHouseMatchesNoCardsInPlay is met when no card in play, across both
	// players and every card type, belongs to the chosen active house.
	ActiveHouseMatchesNoCardsInPlay = engine.ActiveHouseMatchesNoCardsInPlay
)

// Æmber-pool comparisons for card.PoolAember{Player: ..., Is: ..., Amount: n}.
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
	// CountersOnThis counts the generic counters of one kind on the source card.
	CountersOnThis = engine.CountersOnThis
	// PurgedCards counts every card in the purge pile across both players.
	PurgedCards = engine.PurgedCards
	// TurnCount counts one of the engine's turn-history tallies (Player + Of).
	TurnCount = engine.TurnCount // ForgedKey gates on whether a player forged a key this turn or their previous one.
	ForgedKey = engine.ForgedKey
	// OpponentHasMoreKeys is met when your opponent has forged more keys than you.
	OpponentHasMoreKeys = engine.OpponentHasMoreKeys
	// KeyColorForged is met while a player has forged a key of a given colour.
	KeyColorForged = engine.KeyColorForged
	// AemberStolenFromYou is met if your opponent stole Æmber from you last turn.
	AemberStolenFromYou = engine.AemberStolenFromYou
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
	// CountersOnThisAtLeast is met when this card carries at least N counters of the given kind.
	CountersOnThisAtLeast = engine.CountersOnThisAtLeast
	// NamedCardPurged is met by whether a card of the given name is in your purge pile.
	NamedCardPurged = engine.NamedCardPurged
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
	// AemberBonusOf counts the Æmber pips on the card its Target names, once it has
	// left play (Rustgnawer gains the destroyed artifact's Æmber bonus via Target:
	// Triggering).
	AemberBonusOf = engine.AemberBonusOf
	// CardsPurged counts the creatures the most recent purge removed "this way".
	CardsPurged = engine.CardsPurged
	// PurgedAemberBonus totals the printed Æmber bonus of the cards the most recent
	// purge removed (Infurnace).
	PurgedAemberBonus = engine.PurgedAemberBonus

	// ProducedThisWay counts a "... this way" tally an earlier effect in the same
	// resolution recorded — creatures destroyed or shuffled home, Æmber lost, or
	// cards returned. Tally names which (card.Tally.*); Player names whose share.
	ProducedThisWay = engine.ProducedThisWay
	// AemberInPool counts the Æmber currently in a player's pool.
	AemberInPool = engine.AemberInPool
	// AemberOnFriendlyCreatures counts the Æmber sitting on your creatures.
	AemberOnFriendlyCreatures = engine.AemberOnFriendlyCreatures
	// NeighborsOfThis counts the battleline neighbors of the source creature (0-2).
	NeighborsOfThis = engine.NeighborsOfThis
	// NeighborsSharingHouse counts the neighbors of the context creature (ctx.It)
	// that share its house (0-2).
	NeighborsSharingHouse = engine.NeighborsSharingHouse
	// CardsInHand counts the cards in a player's hand of a referenced house.
	CardsInHand = engine.CardsInHand
	// CreaturesHealed counts the creatures the most recent Heal healed.
	CreaturesHealed = engine.CreaturesHealed
	// DamageHealed counts the damage the most recent Heal removed (for DealDamage.AmountFrom).
	DamageHealed = engine.DamageHealed
	// DamagePrevented counts the damage a creature just prevented with its own armor.
	DamagePrevented = engine.DamagePrevented
	// AemberStolenThisEvent counts the Æmber taken in the theft that fired an After Æmber Is Stolen From You ability.
	AemberStolenThisEvent = engine.AemberStolenThisEvent
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
	// CreaturesCannot bars every creature in play — both players' — from fighting
	// or reaping until the caster's next turn, sparing an excepted house (Into the
	// Night, Sow Salt).
	CreaturesCannot = engine.CreaturesCannot
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
	// ForceOpponentActiveHouseOfFought forces the opponent to the house of the
	// creature this card fought as their active house next turn (Snag).
	ForceOpponentActiveHouseOfFought = engine.ForceOpponentActiveHouseOfFought
	// ForbidOpponentActiveHouse bars the opponent's chosen house next turn (Tezmal).
	ForbidOpponentActiveHouse = engine.ForbidOpponentActiveHouse
	// ForbidSameActiveHouseNextTurn bars the chooser's opponent from matching the
	// just-chosen house next turn (Snag's Mirror).
	ForbidSameActiveHouseNextTurn = engine.ForbidSameActiveHouseNextTurn
	// WagerOpponentChoosesChosenHouse steals if the opponent picks the chosen house
	// as their active house next turn (Snaglet).
	WagerOpponentChoosesChosenHouse = engine.WagerOpponentChoosesChosenHouse
	// ForgeKey has the controller forge a key outside the normal step.
	ForgeKey = engine.ForgeKey
	// UnforgeKey takes a forged key back off a player (Key Hammer).
	UnforgeKey = engine.UnforgeKey
	// RaiseKeyCost makes keys cost more throughout a player's next turn.
	RaiseKeyCost = engine.RaiseKeyCost
	// LowerKeyCost drops keys' cost (a negative bump) for a Duration; may be EachPlayer.
	LowerKeyCost = engine.LowerKeyCost
	// GiveRemainingAemberAfterOpponentForgeKey arms Interdimensional Graft's delayed gift.
	GiveRemainingAemberAfterOpponentForgeKey = engine.GiveRemainingAemberAfterOpponentForgeKey
	// GainChains gives a player chains (a draw penalty).
	GainChains = engine.GainChains
)

// Counter groups the generic counter kinds a card can place or read, e.g.
// card.PlaceCounter{Kind: card.Counter.Doom, Target: card.Target.ChosenCreature}.
var Counter = counters{
	Doom:       engine.CounterDoom,
	Fuse:       engine.CounterFuse,
	Growth:     engine.CounterGrowth,
	Glory:      engine.CounterGlory,
	Disruption: engine.CounterDisruption,
}

type counters struct {
	Doom       engine.CounterKind
	Fuse       engine.CounterKind
	Growth     engine.CounterKind
	Glory      engine.CounterKind
	Disruption engine.CounterKind
}

// Tally names a per-resolution "... this way" tally for card.ProducedThisWay, e.g.
// card.ProducedThisWay{Tally: card.Tally.AemberLost, Player: card.Controller}.
var Tally = tallies{
	CreaturesDestroyed:        engine.TallyCreaturesDestroyed,
	CreaturesShuffledIntoDeck: engine.TallyCreaturesShuffledIntoDeck,
	AemberLost:                engine.TallyAemberLost,
	CardsReturned:             engine.TallyCardsReturned,
}

type tallies struct {
	CreaturesDestroyed        engine.ProducedTally
	CreaturesShuffledIntoDeck engine.ProducedTally
	AemberLost                engine.ProducedTally
	CardsReturned             engine.ProducedTally
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

// OneNeighbor makes a CreatureAndNeighbors spread hit one chosen neighbor instead
// of every neighbor (Mighty Lance): card.CreatureAndNeighbors{Scope: card.OneNeighbor}.
var OneNeighbor = engine.OneNeighbor

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
