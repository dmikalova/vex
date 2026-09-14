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
	// GainAember moves Æmber from the common supply into a player's pool (or gains
	// Æmber equal to a running count via EqualTo, e.g. The Flex gains half a chosen
	// creature's power).
	GainAember = engine.GainAember
	// MoveAemberFromPool banks Æmber out of your pool onto a card.
	MoveAemberFromPool = engine.MoveAemberFromPool
	// PlaceAemberOnThis places Æmber from the common supply on this card.
	PlaceAemberOnThis = engine.PlaceAemberOnThis
	// LoseAember returns Æmber from a player's pool to the supply (see By:
	// HalfRoundedDown, AllBut; or EqualTo for a running count, e.g. Power of Fire
	// loses half the sacrificed creature's power).
	LoseAember = engine.LoseAember
	// StealAember moves Æmber from the opponent's pool into yours.
	StealAember = engine.StealAember
	// GiveAember moves Æmber from your opponent's pool into yours — the opponent
	// gives you a fixed Amount (a toll) or All their remaining Æmber.
	GiveAember = engine.GiveAember
	// CaptureAember moves Æmber from a pool onto a capturing creature.
	CaptureAember = engine.CaptureAember
	// CaptureFromAnyPlayer captures Æmber onto this creature from both pools in any split.
	CaptureFromAnyPlayer = engine.CaptureFromAnyPlayer
	// MoveAemberToSupply removes Æmber sitting on a creature to the common supply.
	MoveAemberToSupply = engine.MoveAemberToSupply
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
	// ForEachHouse resolves Do once per house, binding that house so a nested
	// "of that house" target (card.Houses.Each) reads it (Gleeful Mayhem).
	ForEachHouse = engine.ForEachHouse
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
	// GainAssault gives each targeted creature Assault equal to a count for the
	// remainder of the turn (Creed of Nature grants assault equal to its power).
	GainAssault = engine.GainAssault
	// GainKeywordForTurn gives each targeted creature a keyword for the remainder
	// of the turn (Creed of Nature grants skirmish), unlike GainKeyword's
	// until-next-turn duration.
	GainKeywordForTurn = engine.GainKeywordForTurn
	// GainTextBox gives the creature its Target selects the printed text box of the
	// creature its Source selects — that card's traits, keywords, and triggered
	// abilities. RemainderOfTurn makes the gain last the turn; otherwise it lasts
	// until the recipient leaves play (Mimic Gel copies a chosen creature).
	GainTextBox = engine.GainTextBox
	// LendTextBoxFromHand reveals a creature from your hand and gives a chosen
	// creature in play that revealed card's text box for the remainder of the turn
	// (Creed of Nurture).
	LendTextBoxFromHand = engine.LendTextBoxFromHand
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

// ByActivePlayer resolves its inner effect as the active player (the chooser),
// not the ability's controller.
type ByActivePlayer = engine.ByActivePlayer

// Selection is how a purge or zone-movement effect picks its cards: the
// controller chooses one (Chosen), one is random (Random), or every match is
// taken (Each).
type (
	// Selection is the axis a movement verb varies along; set it on PurgeFromHand.
	Selection = engine.Selection
	// Chosen has the controller pick one card, narrowed by House and an identity
	// filter (Type, Trait, Name, or an Or disjunction); it is mandatory by default,
	// and Optional makes it a "you may".
	Chosen = engine.Chosen
	// Random takes one uniformly random card.
	Random = engine.Random
	// Each takes every card the House and identity filters (Type, Trait, Name, Or)
	// admit.
	Each = engine.Each
	// Named pins the pick to the first card of a given name (Hyde archives Velum).
	Named = engine.Named
	// Top pins the pick to the top card of an ordered zone (deck or discard pile).
	Top = engine.Top
	// Bottom pins the pick to the bottom card of an ordered zone.
	Bottom = engine.Bottom
)

// Destruction and purging.
type (
	// Destroy removes the creatures its Target selects from play.
	Destroy = engine.Destroy
	// DestroyChosen destroys any number of creatures the controller picks from its Target.
	DestroyChosen = engine.DestroyChosen
	// BatchDestroy destroys a computed set of creatures in one simultaneous batch, so
	// their Destroyed abilities see each other still in play. Its Gather strategy
	// picks and names the set (e.g. card.EachPlayerUnless).
	BatchDestroy = engine.BatchDestroy
	// EachPlayerUnless gathers, from each player not spared by a per-player board
	// condition (Spare), the creatures its Take refinement keeps — Quicksand takes
	// each unspared player's most powerful creature.
	EachPlayerUnless = engine.EachPlayerUnless
	// DestroyEachCreatureAtEndOfTurn schedules "destroy each creature" to resolve in
	// the end-of-turn phase rather than now (Ragnarok).
	DestroyEachCreatureAtEndOfTurn = engine.DestroyEachCreatureAtEndOfTurn
	// PurgeCard sets cards aside out of the game, from a discard pile, with a
	// Selection deciding which cards leave and Player choosing the pile(s):
	// card.ChosenPlayer for one the controller picks, card.EachPlayer for both.
	PurgeCard = engine.PurgeCard
	// PurgeFromHand purges cards from a player's hand, with a Selection deciding how
	// they are picked (chosen / random / each).
	PurgeFromHand = engine.PurgeFromHand
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
	// PurgeArchives purges any number of cards from your archives, recording the
	// tally a following CardsPurged scales by (Destructive Analysis).
	PurgeArchives = engine.PurgeArchives
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
	// RemoveWard takes the ward off a creature — any creature, warded or not (Hunter or Hunted?).
	RemoveWard = engine.RemoveWard
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
	// PutFromDiscard moves cards from your discard pile to a Destination, with a
	// Selection (Chosen or Each) deciding which cards.
	PutFromDiscard = engine.PutFromDiscard
	// Filter is a predicate selecting cards by type, trait, and/or name, with Or
	// alternatives — e.g. an upgrade or a Robot card. It is a Chosen/Each Or element.
	Filter = engine.CardFilter
	// PutFromHand puts a chosen card from your hand directly into play.
	PutFromHand = engine.PutFromHand
	// ReturnNamedToHand returns a chosen card of a given name to its owner's hand.
	ReturnNamedToHand = engine.ReturnNamedToHand
	// ReturnItToHand returns the creature in context ("it") to its owner's hand,
	// recovering it from the discard pile when it was already destroyed (Nizak, The
	// Forgotten returns an enemy destroyed fighting it).
	ReturnItToHand = engine.ReturnItToHand
	// SearchForName searches your deck and discard pile for a named card.
	SearchForName = engine.SearchForName
	// SearchDeck searches your deck for a card (any card, or one of a given house),
	// puts it into your hand, and shuffles your deck.
	SearchDeck = engine.SearchDeck
	// Shuffle shuffles your deck, optionally folding whole zones of your cards into
	// it first: nothing (the bare "shuffle your deck" a search ends on, also a
	// RevealTopOfDeck/LookAtTopOfDeck routing terminal — Borr Nit) or whole Zones
	// (hand, discard, archives). To fold every friendly card in play into the deck
	// instead, use ShuffleFriendlyCardsIntoDeck.
	Shuffle = engine.Shuffle
	// ShuffleFriendlyCardsIntoDeck folds every friendly card in play plus its
	// upgrades into its owner's deck, tallying how many returned to each owner's deck
	// so a following Draw{Per: CardsShuffledIntoDeck} draws one card for each card
	// that returned to your own deck (Timequake).
	ShuffleFriendlyCardsIntoDeck = engine.ShuffleFriendlyCardsIntoDeck
	// ShuffleFromDiscard shuffles the cards a Selection picks from your discard pile
	// into your deck — each match, any number of a chosen kind, or a counted number.
	ShuffleFromDiscard = engine.ShuffleFromDiscard
	// ShuffleChosenCreaturesFromZones shuffles any number of chosen friendly creatures from your hand, discard pile, or battleline into your deck.
	ShuffleChosenCreaturesFromZones = engine.ShuffleChosenCreaturesFromZones
	// ShuffleNamedFromDiscardIntoDeck shuffles one card of a given name from your
	// discard pile into your deck.
	ShuffleNamedFromDiscardIntoDeck = engine.ShuffleNamedFromDiscardIntoDeck
	// SwapDeckAndDiscard exchanges the controller's deck with their discard pile,
	// then shuffles.
	SwapDeckAndDiscard = engine.SwapDeckAndDiscard
	// ArchiveCard sets cards aside into the controller's own archives, with a
	// Selection deciding how each is picked (Chosen, Random, Named, or Top), Zone
	// the source (Hand, Discard, or Deck), and Amount / Revealed / Per / Or the rest.
	ArchiveCard = engine.ArchiveCard
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
	// DiscardHand discards a player's whole hand, one card at a time (Player may be
	// card.EachPlayer). RefillHand redraws it as if the turn had ended. Punctuated
	// Equilibrium composes both over EachPlayer.
	DiscardHand = engine.DiscardHand
	// RefillHand — see DiscardHand.
	RefillHand = engine.RefillHand
	// DiscardCard discards cards from a player's hand or archives, with a Selection
	// deciding how each is picked (Chosen or Random), Zone the source, and
	// Amount / AnyNumber how many.
	DiscardCard = engine.DiscardCard
	// DiscardUntil discards from the top of your deck until it turns up a
	// card the filters admit, recording the run and putting that card in context.
	DiscardUntil = engine.DiscardUntil
	// PutDiscardedIntoHand puts the card in context from the discard pile into
	// its owner's hand.
	PutDiscardedIntoHand = engine.PutDiscardedIntoHand
	// DiscardTop discards the top cards of one or both decks (Player EachPlayer for
	// both, unset for a granted ability's own deck), records each discarded card,
	// and binds a lone discard as context for a single-card follow-up.
	DiscardTop = engine.DiscardTop
	// ForEachDiscarded resolves Do once for each card a preceding discard removed.
	ForEachDiscarded = engine.ForEachDiscarded
	// ArchiveDiscardedThisWay archives every card a preceding deck dig discarded.
	ArchiveDiscardedThisWay = engine.ArchiveDiscardedThisWay
	// RevealTopOfDeck reveals the top Amount cards of a deck to both players, binds
	// the top one in context, and routes them through the ordered Then steps. Set
	// ChooseWhoseDeck to have the controller pick whose deck (Borr Nit). Revealing
	// one card with no steps is the inspect-and-play primitive (Chaos Portal).
	RevealTopOfDeck = engine.RevealTopOfDeck
	// ChangeActiveHouse changes the active player's active house for the rest of
	// the turn to the house To names — TheContextualHouse, the card in context
	// (Book of leQ).
	ChangeActiveHouse = engine.ChangeActiveHouse
	// EndTurn ends the active player's turn in place, running the turn out the way
	// the Omega keyword does (Book of leQ).
	EndTurn = engine.EndTurn
	// PlayRevealedCard plays the card a preceding reveal put in context.
	PlayRevealedCard = engine.PlayRevealedCard
	// PutRevealedCard moves the card a preceding reveal put in context from its
	// owner's deck to To (card.Into.Archives → "archive it", card.Into.Discard →
	// "discard it").
	PutRevealedCard = engine.PutRevealedCard
	// PlayTopOfDeck plays the top card of the controller's deck outright.
	PlayTopOfDeck = engine.PlayTopOfDeck
	// LookAtTopOfDeck looks privately at the top Amount cards of your deck and
	// routes them through the ordered Then steps (Eyegor, Philophosaurus, Navigator
	// Ali, Lay of the Land). With no steps it is a pure peek.
	LookAtTopOfDeck = engine.LookAtTopOfDeck
	// TopAct is one routing step of a LookAtTopOfDeck or RevealTopOfDeck Then list.
	TopAct = engine.TopAct
	// ChooseAndMove takes Count of the read cards and sends them to Dest (see
	// card.Into).
	ChooseAndMove = engine.ChooseAndMove
	// DeckDest is where a ChooseAndMove step sends the cards it takes (see
	// card.Into).
	DeckDest = engine.DeckDest
	// ReorderRest puts the cards no earlier step took back in any order; it must be
	// the last step.
	ReorderRest = engine.ReorderRest
	// (From), ignoring the active house. Set Except to make House the house that
	// may not be played.
	PlayFrom = engine.PlayFrom
	// PlayOrUse immediately either plays a matching card from the controller's hand
	// or uses a matching card they have in play, in one prompt (CXO Taber's
	// non-Star Alliance card). Except makes House the house that may not be chosen.
	PlayOrUse = engine.PlayOrUse
	// PlayFromOpponent plays a card from the opponent's deck (From: card.Deck, its
	// top card) or archives (From: card.Archives, a random card) as your own play
	// (Murkens).
	PlayFromOpponent = engine.PlayFromOpponent
	// DiscardOpponentArchivesOrDeckTop discards one card from a source you pick
	// between the opponent's archives (a random card) and their deck top, binding it
	// as "it" for a following effect (Fidgit).
	DiscardOpponentArchivesOrDeckTop = engine.DiscardOpponentArchivesOrDeckTop
	// PlayItFromOpponentDiscard plays the card in context (put there by a preceding
	// discard) from the opponent's discard pile as your own (Fidgit plays it when it
	// is a Tactic).
	PlayItFromOpponentDiscard = engine.PlayItFromOpponentDiscard
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
	// CancelForge makes the opponent's key forge in progress not occur, with no
	// Æmber spent (Keyforgery). It is the forge counterpart to CancelFight.
	CancelForge = engine.CancelForge
	// RevealHand shows the cards in a player's hand to both players and records them.
	RevealHand = engine.RevealHand
	// RevealRandomFromHand reveals a random card from your hand and puts it in
	// context for a following effect (Keyforgery).
	RevealRandomFromHand = engine.RevealRandomFromHand
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
	// TurnIntoCreature turns the targeted card into a creature and moves it to a
	// flank of its controller's battleline (Auto-Legionary).
	TurnIntoCreature = engine.TurnIntoCreature
)

// Composites and control flow.
type (
	// Sequence resolves several effects in order.
	Sequence = engine.Sequence
	// ForDuration applies several timed effects sharing one duration and
	// renders their shared "for the remainder of the turn, ..." clause once.
	ForDuration = engine.ForDuration
	// ForEach resolves an effect once for each of a running count, choosing
	// afresh each time.
	ForEach = engine.ForEach
	// Sentences resolves several effects in order, each rendered as its own
	// sentence rather than joined with ", and".
	Sentences = engine.Sentences // ChooseOne offers the controller a set of alternative effects to pick from.
	ChooseOne = engine.ChooseOne
	// ChooseHouseThen asks the controller to choose a house, then resolves Then.
	ChooseHouseThen = engine.ChooseHouseThen
	// OpponentNamesHouse has the opponent name a house, stored for a following
	// ItIsNotOfNamedHouse to read (Keyforgery).
	OpponentNamesHouse = engine.OpponentNamesHouse
	// Conditional resolves Then only when Cond is met.
	Conditional = engine.Conditional
	// Repeat resolves Do and repeats it as its Gate allows — While (automatically
	// while a condition holds), MayWhile (optionally at the controller's choice), or
	// ByExalting (once, paid by exalting a creature).
	Repeat = engine.Repeat
	// While repeats automatically while its Cond holds (Numquid the Fair).
	While = engine.While
	// MayWhile repeats at the controller's choice while its Cond holds (Bouncing
	// Deathquark).
	MayWhile = engine.MayWhile
	// ByExalting repeats once if the controller exalts its Creature to pay for it
	// (Phalanx Strike, Tribute).
	ByExalting = engine.ByExalting
	// May makes an effect optional — the controller chooses whether to resolve it.
	May = engine.May
	// Then is the A -> B result gate: resolves Result only when First did something.
	Then = engine.Then
	// OrAmount switches an effect's amount to Amount when When holds, printing the
	// linear "<base>, or <alt> if <cond>" form instead of an Otherwise branch (rule
	// 22). Used as StealAember{Amount: 1, Or: card.OrAmount{Amount: 2, When: ...}}.
	OrAmount = engine.OrAmount
)

// Conditions gate a Conditional or a Repeat's While/MayWhile gate.
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
	// OnFlank gates on a creature's flank position — the source card, or ctx.It when
	// OfIt is set; Where picks any flank or the left/right flank. Wrap in card.Not
	// for the off-flank sense.
	OnFlank = engine.OnFlank
	// CardsInDiscardAtLeast is met when the controller's discard pile holds at least
	// Amount cards matching House and Type.
	CardsInDiscardAtLeast = engine.CardsInDiscardAtLeast
	// Haunted is met while the controller has 10 or more cards in their discard pile.
	Haunted = engine.Haunted
	// ControlsNamed is met when the controller has a card of a given name in play.
	ControlsNamed = engine.ControlsNamed
	// SourceInCenterOfBattleline is met while the source card sits in the center
	// of its controller's battleline (an even-sized line has no center).
	SourceInCenterOfBattleline = engine.SourceInCenterOfBattleline
	// SourceReady is met while the source card is ready (Bellowing Patrizate's gate).
	SourceReady = engine.SourceReady
	// SourceIsFighting is met while the source creature is one of the two combatants
	// of the fight resolving right now (Nizak, The Forgotten's "while fighting" gate).
	SourceIsFighting = engine.SourceIsFighting
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
	// Overwhelmed is met while the opponent controls more creatures than you.
	Overwhelmed = engine.Overwhelmed
	// ItIsFriendly is met when the card in context is controlled by you.
	ItIsFriendly = engine.ItIsFriendly
	// TideIsLow is met when the tide is low for you.
	TideIsLow = engine.TideIsLow
	// TideIsHigh is met when the tide is high for you.
	TideIsHigh = engine.TideIsHigh
	// ItIs is met when the card in context matches a House (named, non-<house>, or
	// the chosen/active house) and/or Type filter.
	ItIs = engine.ItIs
	// ItIsNamed is met when the card in context carries a given printed name.
	ItIsNamed = engine.ItIsNamed
	// ItIsOfTrait is met when the creature in context has the named trait.
	ItIsOfTrait = engine.ItIsOfTrait
	// ItHasAember is met when the creature in context has Æmber on it.
	ItHasAember = engine.ItHasAember
	// Or is met when any of its Conditions is met, composing conditions (e.g. a
	// Dinosaur creature or one with Æmber, for Guji Dinosaur Hunter).
	Or = engine.Or
	// Not is met when its inner Condition is not met, composing negation instead of
	// a per-condition flag (card.Not{Cond: card.OnFlank{}} → "if it is not on a
	// flank"). The inner condition must render its own negated clause.
	Not = engine.Not
	// ItIsOffIdentity is met when the card in context is off your identity houses.
	ItIsOffIdentity = engine.ItIsOffIdentity
	// ItIsStunned is met when the creature in context is already stunned.
	ItIsStunned = engine.ItIsStunned
	// ItIsNotOfNamedHouse is met when a card is in context and is not of the house
	// a player named earlier in this ability (Keyforgery).
	ItIsNotOfNamedHouse = engine.ItIsNotOfNamedHouse
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
	// HasMoreForgedKeys is met when the named Player has forged more keys than the
	// other player.
	HasMoreForgedKeys = engine.HasMoreForgedKeys
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
	// CardsInZone counts the cards in one of a player's zones (deck, hand,
	// archives, or discard pile).
	CardsInZone = engine.CardsInZone
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
	// PowerDestroyedThisWay totals the board power of the creatures the most recent
	// destruction removed "this way", each measured just before it left play, for a
	// CountIs threshold (Might Makes Right forges only above 25 power).
	PowerDestroyedThisWay = engine.PowerDestroyedThisWay
	// CardsShuffledIntoDeck counts the cards the most recent from-play shuffle
	// returned to your own deck "this way" (Timequake draws one card for each). A
	// card you played but do not own returns to its owner's deck and is not counted.
	CardsShuffledIntoDeck = engine.CardsShuffledIntoDeck
	// AemberBonusOf counts the Æmber pips on the card its Target names, once it has
	// left play (Rustgnawer gains the destroyed artifact's Æmber bonus via Target:
	// Triggering).
	AemberBonusOf = engine.AemberBonusOf
	// UpgradesOn counts the upgrades attached to the creature its Target names
	// (Walls' Blaster stuns a creature for each upgrade on Chief Engineer Walls).
	UpgradesOn = engine.UpgradesOn
	// CardsPurged counts the cards the most recent purge removed "this way", both
	// sides together; its Type names the noun (unset "card", card.Type.Creature
	// "creature").
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
	// of it (Of: HalfRoundedDown — The Flex).
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
	// ForOpponentNextTurn installs a reaction that fires during your opponent's next turn.
	ForOpponentNextTurn = engine.ForOpponentNextTurn
	// Instead installs a replacement that changes an event's outcome for the turn.
	Instead = engine.Instead
	// Replace is a continuous replacement an Upgrade applies to a game event.
	Replace = engine.Replace
	// NextPlayed makes the next creature of a house you play do something.
	NextPlayed = engine.NextPlayed
	// FuseTriggersForTurn makes each friendly creature's A and B effects each fire on
	// the other for the rest of your turn (Livia the Elder fuses fight and reap).
	FuseTriggersForTurn = engine.FuseTriggersForTurn
	// AlsoTriggersOn declares that an ability under one trigger also fires on
	// another, for a ConstantAbility.AlsoTriggers (Kompsos Haruspex makes a play
	// effect also fire on reap).
	AlsoTriggersOn = engine.AlsoTriggersOn
)

// Houses, keys, chains, and restrictions.
type (
	// Restrict bars a player from an action for a Duration — using creatures to
	// fight, using creatures to reap, or using any cards at all. Action selects
	// which (card.Restricted.Fighting/Reaping/Use); RestrictFighting is only valid
	// for the player's next turn.
	Restrict = engine.Restrict
	// CannotPlay bars a player from playing cards of a Type for a Duration.
	CannotPlay = engine.CannotPlay
	// PlayersCannotPlay bars both players from playing cards of a Type until the
	// end of the caster's next turn (Stealth Mode stops either player playing
	// Tactics).
	PlayersCannotPlay = engine.PlayersCannotPlay
	// CreaturesCannot bars every creature in play — both players' — from fighting
	// or reaping until the caster's next turn, reaching the houses its matcher
	// admits (Into the Night, Sow Salt).
	CreaturesCannot = engine.CreaturesCannot
	// BlankEnemyText blanks enemy creatures' text boxes until your next turn (Shadow of Dis).
	BlankEnemyText = engine.BlankEnemyText
	// SkipForgePhase makes a player skip their forge-a-key phase next turn.
	SkipForgePhase = engine.SkipForgePhase
	// CannotBeDealtDamage marks the targeted creatures unable to be dealt damage
	// for a Duration.
	CannotBeDealtDamage = engine.CannotBeDealtDamage
	// MayPlayOrUse lets the controller act with cards outside their active house for
	// the remainder of the turn — the one node for every out-of-house permission
	// grant. Houses selects whose cards it frees (card.GrantHouses.Named/Chosen/Any/
	// Except/Controlled), Grant the verbs (card.GrantPlay, card.GrantUse,
	// card.GrantFight, or a combination), Types narrows the card types (the zero
	// value frees all — card.Types.Of(card.Type.Artifact, ...)), and Count
	// bounds how many cards (zero is unlimited).
	MayPlayOrUse = engine.MayPlayOrUse
	// BelongToHouse makes the targeted creatures belong to a House for a Duration.
	BelongToHouse = engine.BelongToHouse
	// NameHouse remembers the house an enclosing ChooseHouseThen picked on this card,
	// feeding the card's HouseLock for as long as it stays in play.
	NameHouse = engine.NameHouse
	// MustChooseHouse forces a player's next active-house choice, read from
	// Reference (card.ChosenActiveHouse — Control the Weak; card.FoughtActiveHouse —
	// Snag). Player names whose choice is forced (card.Opponent, or card.Controller
	// for a card that binds its own next turn).
	MustChooseHouse = engine.MustChooseHouse
	// CannotChooseHouse bars a player's next active-house choice, read from Reference
	// (card.ChosenActiveHouse — Tezmal; card.JustChosenActiveHouse — Snag's Mirror).
	// Player names whose choice is barred.
	CannotChooseHouse = engine.CannotChooseHouse
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
	CardsPurged:               engine.TallyCardsPurged,
}

type tallies struct {
	CreaturesDestroyed        engine.ProducedTally
	CreaturesShuffledIntoDeck engine.ProducedTally
	AemberLost                engine.ProducedTally
	CardsReturned             engine.ProducedTally
	CardsPurged               engine.ProducedTally
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
	AemberStolen:           engine.EventAemberStolen,
	CardPlayed:             engine.EventCardPlayed,
	Forge:                  engine.EventForgeKey,
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
	AemberStolen,
	CardPlayed,
	Forge engine.Event
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

// HalfRoundedDown, HalfRoundedUp, ThirdRoundedDown, and ThirdRoundedUp are the
// Fractions the game uses — a half or a third of a quantity, with rounding stated
// explicitly. One Fraction serves an Æmber-pool share (card.LoseAember{By:
// card.HalfRoundedDown}), a creature's power (card.PowerOfChosen{Of:
// card.HalfRoundedDown}), and a battleline count
// (card.Target.EachCreature.Refine(card.PortionPerSide(card.ThirdRoundedUp))).
var HalfRoundedDown = engine.HalfRoundedDown

// HalfRoundedUp — see HalfRoundedDown.
var HalfRoundedUp = engine.HalfRoundedUp

// ThirdRoundedDown — see HalfRoundedDown.
var ThirdRoundedDown = engine.ThirdRoundedDown

// ThirdRoundedUp — see HalfRoundedDown.
var ThirdRoundedUp = engine.ThirdRoundedUp

// OneNeighbor makes a CreatureAndNeighbors spread hit one chosen neighbor instead
// of every neighbor (Mighty Lance): card.CreatureAndNeighbors{Scope: card.OneNeighbor}.
var OneNeighbor = engine.OneNeighbor

// GrantPlay, GrantUse, and GrantFight compose a MayPlayOrUse grant's verbs, e.g.
// card.MayPlayOrUse{Houses: ..., Grant: card.GrantPlay | card.GrantUse}.
var GrantPlay = engine.GrantPlay

// GrantUse — see GrantPlay.
var GrantUse = engine.GrantUse

// GrantFight — see GrantPlay. Frees the freed creatures to fight only.
var GrantFight = engine.GrantFight

// HouseSelector is the Houses axis of a MayPlayOrUse; build one with the
// card.GrantHouses helpers.
type HouseSelector = engine.HouseSelector

// GrantHouses builds the HouseSelector for a MayPlayOrUse grant: Named frees one
// house, Chosen the house an enclosing ChooseHouseThen picked, Any every house,
// Except every house but one, and Controlled every house you have a card in play
// for.
var GrantHouses = houseSelectors{
	Chosen:     engine.HouseSelector{Match: engine.HouseMatcher{Kind: engine.MatchChosenHouse}},
	Any:        engine.HouseSelector{Match: engine.HouseMatcher{Kind: engine.MatchAnyHouse}},
	Controlled: engine.HouseSelector{Controlled: true},
}

type houseSelectors struct {
	Chosen     engine.HouseSelector
	Any        engine.HouseSelector
	Controlled engine.HouseSelector
}

// Named frees one named house.
func (houseSelectors) Named(h engine.House) engine.HouseSelector {
	return engine.HouseSelector{Match: engine.HouseMatcher{Kind: engine.MatchNamedHouse, House: h}}
}

// Except frees every house but the named one (card.House.Self for "non-Star
// Alliance").
func (houseSelectors) Except(h engine.House) engine.HouseSelector {
	return engine.HouseSelector{Match: engine.HouseMatcher{Kind: engine.MatchExceptHouse, House: h}}
}

// HouseMatcher is the house filter a per-card effect or target names — which
// houses it admits (ADR 0038). Build one with the card.Houses helpers.
type HouseMatcher = engine.HouseMatcher

// Houses builds the HouseMatcher that narrows a per-card effect or target to a
// house: Named admits one house, Except every house but one, Chosen the house an
// enclosing ChooseHouseThen picked, Active the active house, Each the house an
// enclosing ForEachHouse is on, Contextual the house of the card in context
// (ctx.It), and Any (the zero value) every house.
var Houses = houseMatchers{
	Any:        engine.HouseMatcher{Kind: engine.MatchAnyHouse},
	Chosen:     engine.HouseMatcher{Kind: engine.MatchChosenHouse},
	Active:     engine.HouseMatcher{Kind: engine.MatchActiveHouse},
	Each:       engine.HouseMatcher{Kind: engine.MatchEachHouse},
	Contextual: engine.HouseMatcher{Kind: engine.MatchContextualHouse},
}

type houseMatchers struct {
	Any        engine.HouseMatcher
	Chosen     engine.HouseMatcher
	Active     engine.HouseMatcher
	Each       engine.HouseMatcher
	Contextual engine.HouseMatcher
}

// Restricted names the action a card.Restrict bars: Fighting and Reaping bar one
// verb (creatures can still be used the other ways); Use bars every use of a card
// in play. Fighting is only valid for the player's next turn.
var Restricted = restrictKinds{
	Fighting: engine.RestrictFighting,
	Reaping:  engine.RestrictReaping,
	Use:      engine.RestrictUse,
}

type restrictKinds struct {
	Fighting engine.RestrictKind
	Reaping  engine.RestrictKind
	Use      engine.RestrictKind
}

// Named admits only the named house ("a Mars card").
func (houseMatchers) Named(h engine.House) engine.HouseMatcher {
	return engine.HouseMatcher{Kind: engine.MatchNamedHouse, House: h}
}

// Except admits every house but the named one (card.House.Self for "a non-Star
// Alliance card").
func (houseMatchers) Except(h engine.House) engine.HouseMatcher {
	return engine.HouseMatcher{Kind: engine.MatchExceptHouse, House: h}
}

// Types names the card types an effect admits: card.Types.Of(card.Type.Artifact,
// card.Type.Upgrade, card.Type.Tactic). The zero value (an unset Types field)
// admits every type.
var Types = typeSets{}

type typeSets struct{}

// Of builds the set of card types an effect admits.
func (typeSets) Of(types ...engine.CardType) engine.CardTypes {
	return engine.CardTypesOf(types...)
}

// ChosenActiveHouse, FoughtActiveHouse, and JustChosenActiveHouse name where a
// MustChooseHouse or CannotChooseHouse reads its house: the chosen house, the house
// of the creature this card fought, or the house a player just chose (read from the
// board, for Snag's Mirror).
var ChosenActiveHouse = engine.ChosenActiveHouse

// FoughtActiveHouse — see ChosenActiveHouse.
var FoughtActiveHouse = engine.FoughtActiveHouse

// JustChosenActiveHouse — see ChosenActiveHouse.
var JustChosenActiveHouse = engine.JustChosenActiveHouse

// LeftFlank and RightFlank name a flank for an OnFlank predicate, e.g.
// card.OnFlank{OfIt: true, Where: card.LeftFlank}.
var LeftFlank = engine.LeftFlank

// RightFlank — see LeftFlank.
var RightFlank = engine.RightFlank

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

// UpgradesOnIt is the PerTarget that scales a bonus by the number of upgrades
// attached to a creature (Light of the Archons, a StaticModifier with Per set).
var UpgradesOnIt = engine.UpgradesOnIt
