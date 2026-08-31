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
	GainAember    = engine.GainAember
	LoseAember    = engine.LoseAember
	StealAember   = engine.StealAember
	CaptureAember = engine.CaptureAember
	Exalt         = engine.Exalt
	Loss          = engine.Loss

	MoveAemberToPool = engine.MoveAemberToPool
)

// Damage and combat.
type (
	DealDamage                = engine.DealDamage
	DamageIfDestroyed         = engine.DamageIfDestroyed
	DamageIfSurvives          = engine.DamageIfSurvives
	DamageCreatureAndNeighbor = engine.DamageCreatureAndNeighbor
	SplashDamage              = engine.SplashDamage
	DamageDifferent           = engine.DamageDifferent
	FlankWalkDamage           = engine.FlankWalkDamage
	RedirectFightDamage       = engine.RedirectFightDamage
	Heal                      = engine.Heal
	HealThenDamage            = engine.HealThenDamage
)

// Destruction and purging.
type (
	Destroy       = engine.Destroy
	Purge         = engine.Purge
	PurgeFromHand = engine.PurgeFromHand
	PurgeCreature = engine.PurgeCreature

	PurgeCreatureFromHand = engine.PurgeCreatureFromHand
)

// Creature state (stun, exhaust, power counters).
type (
	Stun             = engine.Stun
	Unstun           = engine.Unstun
	Exhaust          = engine.Exhaust
	ExhaustCreatures = engine.ExhaustCreatures
	Ready            = engine.Ready
	ReadyIfFirstUse  = engine.ReadyIfFirstUse
	ReadyCreatures   = engine.ReadyCreatures
	AddPowerCounter  = engine.AddPowerCounter
)

// Drawing, moving, and revealing cards between zones.
type (
	Draw                  = engine.Draw
	PutFromPlay           = engine.PutFromPlay
	PutUpTo               = engine.PutUpTo
	PutFromDiscard        = engine.PutFromDiscard
	ReturnNamedToHand     = engine.ReturnNamedToHand
	SearchForName         = engine.SearchForName
	ShuffleDiscard        = engine.ShuffleDiscard
	ShuffleHandAndDiscard = engine.ShuffleHandAndDiscard
	ArchiveFromHand       = engine.ArchiveFromHand
	ArchiveFromDiscard    = engine.ArchiveFromDiscard
	ArchiveTopOfDeck      = engine.ArchiveTopOfDeck
	ArchiveFromPlay       = engine.ArchiveFromPlay
	DiscardArchives       = engine.DiscardArchives
	DiscardHand           = engine.DiscardHand
	DiscardFromHand       = engine.DiscardFromHand
	DiscardRandomFromHand = engine.DiscardRandomFromHand
	DiscardTopOfDeck      = engine.DiscardTopOfDeck
	DiscardTopOfEachDeck  = engine.DiscardTopOfEachDeck
	ForEachDiscarded      = engine.ForEachDiscarded
	RevealTopOfDeck       = engine.RevealTopOfDeck
	PlayRevealedCard      = engine.PlayRevealedCard
	PlayTopOfDeck         = engine.PlayTopOfDeck
	CancelFight           = engine.CancelFight
	Reveal                = engine.Reveal
)

// Using and choosing creatures.
type (
	OnChooseCreature = engine.OnChooseCreature
	ReadyVerb        = engine.ReadyVerb
	FightVerb        = engine.FightVerb
	UseVerb          = engine.UseVerb
	StunVerb         = engine.StunVerb
	ExhaustVerb      = engine.ExhaustVerb
	Use              = engine.Use
	TakeControl      = engine.TakeControl
	PutIntoPlay      = engine.PutIntoPlay
	Swap             = engine.Swap
)

// Composites and control flow.
type (
	Sequence          = engine.Sequence
	Sentence          = engine.Sentence
	ChooseOne         = engine.ChooseOne
	ChooseHouseThen   = engine.ChooseHouseThen
	Conditional       = engine.Conditional
	RepeatWhile       = engine.RepeatWhile
	RepeatOnCondition = engine.RepeatOnCondition
	MayRepeat         = engine.MayRepeat
	May               = engine.May
	Then              = engine.Then
)

// Conditions gate a Conditional, RepeatWhile, or MayRepeat.
type (
	OpponentAember          = engine.OpponentAember
	CardsDestroyedFewerThan = engine.CardsDestroyedFewerThan
	ControlsMoreCreatures   = engine.ControlsMoreCreatures
	Overwhelmed             = engine.Overwhelmed
	ItIsOfHouse             = engine.ItIsOfHouse
	ItIs                    = engine.ItIs
	ItIsOffIdentity         = engine.ItIsOffIdentity
	ChoseHouse              = engine.ChoseHouse
)

// Æmber-pool comparisons for card.OpponentAember{Is: ..., Amount: n}.
var (
	AtLeast     = engine.AtLeast
	Exactly     = engine.Exactly
	MoreThanYou = engine.MoreThanYou
)

// House references for conditions that compare a card's house dynamically.
var (
	TheChosenHouse     = engine.TheChosenHouse
	TheActiveHouse     = engine.TheActiveHouse
	TheContextualHouse = engine.TheContextualHouse
)

// Counts feed an effect's Per, scaling it by a board quantity; InPlay doubles as
// a Condition.
type (
	InPlay                  = engine.InPlay
	CardsPlayed             = engine.CardsPlayed
	OpponentForgedKeys      = engine.OpponentForgedKeys
	OpponentExcessCreatures = engine.OpponentExcessCreatures
	CardsInArchives         = engine.CardsInArchives
	CardsRevealed           = engine.CardsRevealed
	CardsInHand             = engine.CardsInHand
	CreaturesHealed         = engine.CreaturesHealed
)

// Lasting "for the remainder of the turn" effects.
type (
	ForRemainderOfTurn = engine.ForRemainderOfTurn
	Instead            = engine.Instead
	Replace            = engine.Replace
	NextCreaturePlayed = engine.NextCreaturePlayed
)

// Houses, keys, chains, and restrictions.
type (
	CannotFight                              = engine.CannotFight
	CannotPlay                               = engine.CannotPlay
	SkipForgeStep                            = engine.SkipForgeStep
	PreventDamage                            = engine.PreventDamage
	MayUseFriendlyHouse                      = engine.MayUseFriendlyHouse
	GrantFightForChosenHouse                 = engine.GrantFightForChosenHouse
	BelongToHouse                            = engine.BelongToHouse
	ForceOpponentActiveHouse                 = engine.ForceOpponentActiveHouse
	ForgeKey                                 = engine.ForgeKey
	GiveRemainingAemberAfterOpponentForgeKey = engine.GiveRemainingAemberAfterOpponentForgeKey
	GainChains                               = engine.GainChains
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
}

type events struct {
	CreaturePlayed,
	Reap,
	Fight,
	EnemyCreatureDestroyed,
	ReapAember,
	Destroyed,
	AemberAddedToPool engine.Event
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
