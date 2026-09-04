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
	// MoveAemberFromPool banks Æmber out of your pool onto a card.
	MoveAemberFromPool = engine.MoveAemberFromPool
	// LoseAember returns Æmber from a player's pool to the supply (see By: Half, AllBut).
	LoseAember = engine.LoseAember
	// StealAember moves Æmber from the opponent's pool into yours.
	StealAember = engine.StealAember
	// CaptureAember moves Æmber from a pool onto a capturing creature.
	CaptureAember = engine.CaptureAember
	// Exalt places Æmber from the common supply onto a chosen card.
	Exalt = engine.Exalt
	// Loss says how much Æmber a LoseAember removes (Half, AllBut).
	Loss = engine.Loss
	// MoveAember moves Æmber off a card into a pool or onto another card.
	MoveAember = engine.MoveAember
)

// Damage and combat.
type (
	// DealDamage deals damage to each creature its Target selects.
	DealDamage = engine.DealDamage
	// DamageIfDestroyed deals damage, then runs Then only if the creature left play.
	DamageIfDestroyed = engine.DamageIfDestroyed
	// DamageIfSurvives deals damage, then runs Then only if the creature survives.
	DamageIfSurvives = engine.DamageIfSurvives
	// Spread is a DealDamage strategy that hits several related creatures at once.
	Spread = engine.Spread
	// PerTarget is a DealDamage strategy scaling the damage per creature hit.
	PerTarget = engine.PerTarget
	// CreatureAndNeighbor (a Spread) damages a chosen creature and one of its neighbors.
	CreatureAndNeighbor = engine.CreatureAndNeighbor
	// CreatureAndNeighbors (a Spread) damages a chosen non-flank creature and each of its neighbors.
	CreatureAndNeighbors = engine.CreatureAndNeighbors
	// DifferentCreatures (a Spread) damages a chosen creature and a different chosen creature.
	DifferentCreatures = engine.DifferentCreatures
	// FlankWalk (a Spread) deals decreasing damage inward from a chosen flank creature.
	FlankWalk = engine.FlankWalk
	// RedirectFightDamage is a Before Fight effect redirecting this creature's fight damage.
	RedirectFightDamage = engine.RedirectFightDamage
	// Heal takes damage tokens off a creature — a fixed amount or all of them.
	Heal = engine.Heal
	// LoseArmor takes all the remaining armor off each creature it targets, and
	// tallies it for a following ArmorLostThisWay.
	LoseArmor = engine.LoseArmor
)

// Destruction and purging.
type (
	// Destroy removes the creatures its Target selects from play.
	Destroy = engine.Destroy
	// PurgeCard sets cards aside out of the game, from a named zone.
	PurgeCard = engine.PurgeCard
	// PurgeFromHand purges one card the controller chooses from a player's hand.
	PurgeFromHand = engine.PurgeFromHand
	// PurgeEachFromHand purges every matching card from a player's hand.
	PurgeEachFromHand = engine.PurgeEachFromHand
	// PurgeCreature purges each creature its Target selects from play.
	PurgeCreature = engine.PurgeCreature
	// LoseKeyword takes a keyword from each creature for the remainder of the turn.
	LoseKeyword = engine.LoseKeyword
	// PurgeCreatureFromHand purges a chosen creature from your hand and puts it in context.
	PurgeCreatureFromHand = engine.PurgeCreatureFromHand
)

// Creature state (stun, exhaust, power counters).
type (
	// Stun places a stun on the creatures its Target selects.
	Stun = engine.Stun
	// Unstun removes the stun from each creature its Target selects.
	Unstun = engine.Unstun
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
)

// Drawing, moving, and revealing cards between zones.
type (
	// Draw puts the top Amount cards of your deck into your hand.
	Draw = engine.Draw
	// PutFromPlay takes each targeted card out of play into a chosen Destination.
	PutFromPlay = engine.PutFromPlay
	// PutChosen moves Count cards the controller chooses into a Destination,
	// declinably when UpTo is set.
	PutChosen = engine.PutChosen
	// PutFromDiscard moves a chosen card from your discard pile to a Destination.
	PutFromDiscard = engine.PutFromDiscard
	// ReturnNamedToHand returns a chosen card of a given name to its owner's hand.
	ReturnNamedToHand = engine.ReturnNamedToHand
	// SearchForName searches your deck and discard pile for a named card.
	SearchForName = engine.SearchForName
	// ShuffleIntoDeck shuffles the controller's named zones (hand, discard, archives) into their deck.
	ShuffleIntoDeck = engine.ShuffleIntoDeck
	// SwapDeckAndDiscard exchanges the controller's deck with their discard pile,
	// then shuffles.
	SwapDeckAndDiscard = engine.SwapDeckAndDiscard
	// ArchiveFromHand moves cards from a hand into the controller's archives.
	ArchiveFromHand = engine.ArchiveFromHand
	// ArchiveFromDiscard moves a chosen card from the discard pile into archives.
	ArchiveFromDiscard = engine.ArchiveFromDiscard
	// ArchiveTopOfDeck moves the top Count cards of your deck into archives.
	ArchiveTopOfDeck = engine.ArchiveTopOfDeck
	// ArchiveFromPlay moves each targeted in-play card into its owner's archives.
	ArchiveFromPlay = engine.ArchiveFromPlay
	// DiscardArchives moves all of a player's archived cards into their discard pile.
	DiscardArchives = engine.DiscardArchives
	// DiscardHand discards cards from a player's hand.
	DiscardHand = engine.DiscardHand
	// DiscardFromHand has the controller choose and discard Count cards.
	DiscardFromHand = engine.DiscardFromHand
	// DiscardRandomFromHand discards one uniformly random card from a player's hand.
	DiscardRandomFromHand = engine.DiscardRandomFromHand
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
	// ForEachDiscarded resolves Do once for each card a preceding discard removed.
	ForEachDiscarded = engine.ForEachDiscarded
	// RevealTopOfDeck reveals the top card of the controller's deck.
	RevealTopOfDeck = engine.RevealTopOfDeck
	// PlayRevealedCard plays the card a preceding reveal put in context.
	PlayRevealedCard = engine.PlayRevealedCard
	// PlayTopOfDeck plays the top card of the controller's deck outright.
	PlayTopOfDeck = engine.PlayTopOfDeck
	// PlayFrom plays a card the controller chooses from their hand or discard pile
	// (From), ignoring the active house. Set Except to make House the house that
	// may not be played.
	PlayFrom = engine.PlayFrom
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
	// Use uses up to Max cards the controller chooses from Target.
	Use = engine.Use
	// TriggerAbility fires another card's ability as if the controller controlled it.
	TriggerAbility = engine.TriggerAbility
	// TakeControl moves a card to the controller's play area and makes them its controller.
	TakeControl = engine.TakeControl
	// PutIntoPlay puts each targeted card into play without playing it.
	PutIntoPlay = engine.PutIntoPlay
	// Swap exchanges this creature's battleline position with another.
	Swap = engine.Swap
)

// Composites and control flow.
type (
	// Sequence resolves several effects in order.
	Sequence = engine.Sequence
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
)

// Conditions gate a Conditional, RepeatWhile, or MayRepeat.
type (
	// OpponentAember gates on the opponent's Æmber pool (Is + Amount).
	OpponentAember = engine.OpponentAember
	// CardsDestroyedFewerThan is met when fewer than Amount cards were destroyed this way.
	CardsDestroyedFewerThan = engine.CardsDestroyedFewerThan
	// CountIs gates on any Count compared against a threshold (Count + Is + Amount).
	CountIs = engine.CountIs
	// ControlsMoreCreatures is met while you control more creatures than the opponent.
	ControlsMoreCreatures = engine.ControlsMoreCreatures
	// FirstCreaturePlayedThisTurn is met when the card in context is the first
	// creature played this turn — a once-per-turn charge (Speed Sigil).
	FirstCreaturePlayedThisTurn = engine.FirstCreaturePlayedThisTurn
	// Overwhelmed is met while the opponent controls more creatures than you.
	Overwhelmed = engine.Overwhelmed
	// ItIsOfHouse is met when the card in context belongs to a referenced house.
	ItIsOfHouse = engine.ItIsOfHouse
	// ItIs is met when the card in context matches a concrete House and/or Type.
	ItIs = engine.ItIs
	// ItIsOffIdentity is met when the card in context is off your identity houses.
	ItIsOffIdentity = engine.ItIsOffIdentity
	// ChoseHouse is met when the controller's active house is House.
	ChoseHouse = engine.ChoseHouse
)

// Æmber-pool comparisons for card.OpponentAember{Is: ..., Amount: n}.
var (
	// AtLeast is met when the quantity is at least Amount.
	AtLeast = engine.AtLeast
	// Exactly is met when the quantity is exactly Amount.
	Exactly = engine.Exactly
	// MoreThanYou is met when the opponent's pool holds more Æmber than yours.
	MoreThanYou = engine.MoreThanYou
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
	// ExcessCreatures counts how many more creatures one player controls than the other.
	ExcessCreatures = engine.ExcessCreatures
	// CardsInArchives counts the cards in a player's archives.
	CardsInArchives = engine.CardsInArchives
	// CardsRevealed counts the cards the most recent Reveal showed.
	CardsRevealed = engine.CardsRevealed
	// CardsDestroyed counts the cards the most recent destruction removed "this way".
	CardsDestroyed = engine.CardsDestroyed
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
	// CardsInHand counts the cards in a player's hand of a referenced house.
	CardsInHand = engine.CardsInHand
	// CreaturesHealed counts the creatures the most recent Heal healed.
	CreaturesHealed = engine.CreaturesHealed
	// DamageHealed counts the damage the most recent Heal removed (for DealDamage.AmountFrom).
	DamageHealed = engine.DamageHealed
	// UnforgedKeys counts the keys a player has still to forge.
	UnforgedKeys = engine.UnforgedKeys
	// AemberOnThis counts the Æmber sitting on the source card.
	AemberOnThis = engine.AemberOnThis
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
	// SkipForgeStep makes a player skip their forge-a-key step next turn.
	SkipForgeStep = engine.SkipForgeStep
	// PreventDamage marks the targeted creatures immune to damage for a Duration.
	PreventDamage = engine.PreventDamage
	// MayUseFriendlyHouse lets the controller use their House creatures this turn.
	MayUseFriendlyHouse = engine.MayUseFriendlyHouse
	// GrantFightForChosenHouse lets your chosen-house creatures fight this turn.
	GrantFightForChosenHouse = engine.GrantFightForChosenHouse
	// GrantFightAnyHouse lets every friendly creature fight this turn.
	GrantFightAnyHouse = engine.GrantFightAnyHouse
	// BelongToHouse makes the targeted creatures belong to a House for a Duration.
	BelongToHouse = engine.BelongToHouse
	// NameHouse remembers the house an enclosing ChooseHouseThen picked on this card,
	// feeding the card's HouseLock for as long as it stays in play.
	NameHouse = engine.NameHouse
	// ForceOpponentActiveHouse forces the opponent's active house next turn.
	ForceOpponentActiveHouse = engine.ForceOpponentActiveHouse
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

// Event groups the game events a lasting "for the remainder of the turn" effect
// attaches to (see ForRemainderOfTurn and Instead), e.g.
// card.ForRemainderOfTurn{On: card.Event.CreaturePlayed, Do: card.GainAember{...}}.
var Event = events{
	CreaturePlayed:         engine.EventCreaturePlayed,
	Reap:                   engine.EventReap,
	Fight:                  engine.EventFight,
	EnemyCreatureDestroyed: engine.EventEnemyCreatureDestroyed,
	ReapAember:             engine.EventReapAember,
	Destroyed:              engine.EventCreatureDestroyed,
	AemberAddedToPool:      engine.EventAemberAddedToPool,
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
	CardPlayed engine.Event
}

// Steal is the replacement that makes gaining Æmber steal it from the opponent
// instead, for card.Instead (Dimension Door).
var Steal = engine.Steal

// Capture is the replacement that makes Æmber added to the opponent's pool be
// captured by the source creature instead, for card.WithReplaces (Ether Spider).
var Capture = engine.Capture

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

// ArmorLostThisWay is the PerTarget that scales damage by the armor an effect has
// stripped off each creature hit (Red-Hot Armor).
var ArmorLostThisWay = engine.ArmorLostThisWay
