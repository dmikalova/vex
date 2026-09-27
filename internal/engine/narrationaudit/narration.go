package narrationaudit

// narration is how the game log covers one mutating Resolver method — the
// port-side twin of internal/engine's fieldNarration, which asks the same
// question of one GameState field.
type narration int

const (
	// narratedDirectly: the method appends its own log entry, so a call that
	// changes anything leaves a line behind. Auditor holds it to that.
	narratedDirectly narration = iota
	// narratedByCause: the line is appended by the effect that calls the method,
	// or later when the state the method armed bites (ADR 0011). The comment on
	// each group names the entry. Auditor does not check these: the record lands
	// just outside the call, so a per-call check could only fail them.
	narratedByCause
	// narratedNever: deliberately silent — a read that happens to sit on a
	// mutating role, or bookkeeping no player can see.
	narratedNever
)

// methodNarration classifies every method of the mutating Resolver roles. Like
// fieldNarration it exists to force a decision rather than to describe one:
// TestEveryMutatingPortMethodDeclaresItsNarration fails when a method has no
// entry, so "does this change need a log entry?" is answered when the port grows
// the method. Answering narratedDirectly is not free — Auditor then fails any
// test whose call changes the state without a line.
var methodNarration = map[string]narration{
	// ---- narratedDirectly ----
	//
	// Every one of these records in its own body, so the audit enforces it on
	// every call a test makes.

	// Æmber and keys: AemberGained, AemberStolen, AemberGivenAfterForging,
	// KeyForged, KeyUnforged, ChainsGained.
	"ForgeKeyAtExtraCost": narratedDirectly,
	"ForgeKeyFree":        narratedDirectly,
	"ForgeKeyFreeForced":  narratedDirectly,
	"UnforgeKey":          narratedDirectly,
	"GainChains":          narratedDirectly,

	// Creature state that carries its own line: CreatureGainedStats,
	// CreatureGainedAssault, CreatureGainedKeyword, CreatureLostKeyword,
	// KeywordLostByAll, CreatureGainedTrait, CreatureConsideredFlank,
	// CreatureGainedTextBox, CreatureCopiedStats, ControlTaken, TurnedIntoCreature,
	// PositionsSwapped, MovedToFlank, MovedWithinBattleline.
	"ConsiderFlank":               narratedDirectly,
	"CopyStats":                   narratedDirectly,
	"GainAssault":                 narratedDirectly,
	"GainStats":                   narratedDirectly,
	"GrantAssaultUntilNextTurn":   narratedDirectly,
	"GrantKeyword":                narratedDirectly,
	"GrantKeywordUntilNextTurn":   narratedDirectly,
	"GrantTextBox":                narratedDirectly,
	"GrantTraitUntilNextTurn":     narratedDirectly,
	"LoseKeyword":                 narratedDirectly,
	"LoseKeywordFrom":             narratedDirectly,
	"LoseKeywordUntilNextTurn":    narratedDirectly,
	"MoveToFlank":                 narratedDirectly,
	"MoveWithinBattleline":        narratedDirectly,
	"PutIntoBattlelineAsCreature": narratedDirectly,
	"SwapCards":                   narratedDirectly,
	"TakeControl":                 narratedDirectly,

	// Combat and use: DamageTaken, ArmorAbsorbed, CardDestroyed,
	// CardsDestroyedBy, Fought, Reaped, ActionAbilityUsed.
	"DealDamage":      narratedDirectly,
	"DestroyEach":     narratedDirectly,
	"DestroyEachFrom": narratedDirectly,
	"FightWith":       narratedDirectly,
	"ReapWith":        narratedDirectly,
	"UseActionOf":     narratedDirectly,

	// Zone moves, which are one mechanism (ADR 0031) and narrate as one family:
	// CardMoved, CardDiscarded, CardPurged, CardPutIntoHand, CardPutIntoArchives,
	// CardPutOnTopOfDeck, CardsDrawn, CardShuffledIntoDeck, CardsShuffledIntoDeckBy,
	// DeckAndDiscardSwapped, CardPutUnder, CardGrafted, UpgradeAttached,
	// UpgradeDiscarded, CardPlayedToBattleline and its play siblings.
	"ArchiveCardUnder":           narratedDirectly,
	"ArchiveEnemyFromHand":       narratedDirectly,
	"ArchiveFromDeck":            narratedDirectly,
	"ArchiveFromDiscard":         narratedDirectly,
	"ArchiveFromHand":            narratedDirectly,
	"ArchiveFromPurge":           narratedDirectly,
	"ArchiveUpgrade":             narratedDirectly,
	"DiscardArchives":            narratedDirectly,
	"DiscardCardFromArchives":    narratedDirectly,
	"DiscardCardFromHand":        narratedDirectly,
	"DiscardTopOfDeck":           narratedDirectly,
	"Draw":                       narratedDirectly,
	"GraftUnder":                 narratedDirectly,
	"MoveFromDeckToDiscard":      narratedDirectly,
	"MoveFromDeckToHand":         narratedDirectly,
	"MoveFromDeckToTopOfDeck":    narratedDirectly,
	"MoveFromDiscardToTopOfDeck": narratedDirectly,
	"MoveUpgrade":                narratedDirectly,
	"PlayFromArchives":           narratedDirectly,
	"PlayFromDeck":               narratedDirectly,
	"PlayFromDiscard":            narratedDirectly,
	"PlayFromHand":               narratedDirectly,
	"PlayFromOpponent":           narratedDirectly,
	"PlayFromOpponentDiscard":    narratedDirectly,
	"PlayFromOpponentHand":       narratedDirectly,
	"PlayFromUnder":              narratedDirectly,
	"PurgeFromArchives":          narratedDirectly,
	"PurgeFromDeck":              narratedDirectly,
	"PurgeFromDiscard":           narratedDirectly,
	"PurgeFromHand":              narratedDirectly,
	"PurgeFromPlay":              narratedDirectly,
	"PutCardUnder":               narratedDirectly,
	"PutFromDiscardIntoHand":     narratedDirectly,
	"PutIntoArchives":            narratedDirectly,
	"PutIntoArchivesEach":        narratedDirectly,
	"PutIntoDeckShuffled":        narratedDirectly,
	"PutIntoHand":                narratedDirectly,
	"PutIntoPlay":                narratedDirectly,
	"PutIntoYourArchives":        narratedDirectly,
	"PutOnTopOfDeck":             narratedDirectly,
	"PutUnderIntoPlay":           narratedDirectly,
	"RefillHand":                 narratedDirectly,
	"ShuffleFromDiscardIntoDeck": narratedDirectly,
	"ShuffleFromHandIntoDeck":    narratedDirectly,
	"ShuffleZonesIntoDeck":       narratedDirectly,
	"SwapDeckAndDiscard":         narratedDirectly,
	// EndShuffleBatch is the one that narrates a whole batch, as one grouped
	// CardsShuffledIntoDeckBy per owner.
	"EndShuffleBatch": narratedDirectly,

	// Counters, armor, and a card's house. Each of these narrated nothing until
	// the audit found it: the method knows the amount that actually landed (a
	// counter entry saturates at 255, a strip takes only the armor that was left),
	// so the entry is recorded here rather than by the calling effect —
	// CountersPlaced, CountersRemoved, PowerCountersPlaced, ArmorLost,
	// CardChangedHouse.
	"AddPowerCounter":                 narratedDirectly,
	"PlaceCounter":                    narratedDirectly,
	"RemoveCounters":                  narratedDirectly,
	"RemoveCountersN":                 narratedDirectly,
	"StripArmor":                      narratedDirectly,
	"BelongToHouseForRemainderOfTurn": narratedDirectly,
	"SetLastingHouse":                 narratedDirectly,

	// Bonus icons resolve as if the card had just been played, so each icon
	// narrates through its own BonusAemberGained, BonusDamageDealt, BonusCardDrawn.
	"ResolveBonusIconsOn": narratedDirectly,

	// Turn-level changes a player sees happen at once: HouseChosen,
	// HouseForcedNextTurn, HouseForbiddenNextTurn, HouseWagerArmed,
	// MayPlayOrUseGranted, MayUseTraitGranted, PhaseBegan.
	"CannotChooseHouseNextTurn":     narratedDirectly,
	"EndTurnNow":                    narratedDirectly,
	"GrantMayPlayOrUse":             narratedDirectly,
	"GrantMayUseTrait":              narratedDirectly,
	"MustChooseFoughtHouseNextTurn": narratedDirectly,
	"MustChooseHouseNextTurn":       narratedDirectly,
	"SetActiveHouse":                narratedDirectly,
	"WagerOnHouseNextTurn":          narratedDirectly,

	// ---- narratedByCause ----

	// The Æmber primitives. Each is the arithmetic half of a move the effect
	// narrates around it: AemberGained, AemberLost, AemberStolen, AemberCaptured,
	// AemberExalted, AemberMovedToPool, AemberMovedToCard, AemberLostToMaximum.
	"SetAember":            narratedByCause,
	"AddAmberOn":           narratedByCause,
	"GainAember":           narratedByCause,
	"NoteAemberStolenFrom": narratedByCause,

	// The creature-status primitives, each narrated by the effect that sets it:
	// CreatureExhausted and CreatureReadied; CreatureStunned, CreaturesUnstunned
	// and StunRecovered; CreatureEnraged and CreatureEnrageRemoved;
	// CreatureWarded, WardRemoved and WardAbsorbed; DamageTaken when damage is
	// dealt and CreatureHealed when Heal takes it off (SetDamage is the shared
	// write behind both, so only its callers know which way it went).
	"SetDamage":    narratedByCause,
	"SetEnraged":   narratedByCause,
	"SetExhausted": narratedByCause,
	"SetStunned":   narratedByCause,
	"SetWarded":    narratedByCause,

	// The house a card names as it enters play is the reader half of a lock
	// (Restringuntus, Etan's Jar), and narrates when the lock bites — as the
	// house the opponent then cannot choose, or CardCannotBeUsed.
	"SetNamedCard":  narratedByCause,
	"SetNamedHouse": narratedByCause,

	// Turn bars and board-wide restrictions. Arming one is not an outcome
	// (ADR 0011): the log names the card that armed it through the frame of its
	// play, and narrates when it bites — CardCannotBeUsed for a barred use,
	// DamageRefused for damage an immunity refuses, ForgeSkipped for a skipped
	// forge phase, KeyForged showing the cost a surcharge raised.
	"BlankEnemyText":               narratedByCause,
	"CannotFightNextTurn":          narratedByCause,
	"CannotPlayTypeNextTurn":       narratedByCause,
	"CannotPlayTypeThisTurn":       narratedByCause,
	"CannotReapHouseNextTurn":      narratedByCause,
	"CannotReapNextTurn":           narratedByCause,
	"CannotReapThisTurn":           narratedByCause,
	"CannotUseNextTurn":            narratedByCause,
	"CannotUseThisTurn":            narratedByCause,
	"CreaturesCannotUntilNextTurn": narratedByCause,
	"RaiseKeyCostNextTurn":         narratedByCause,
	"RaiseKeyCostPerHouseNextTurn": narratedByCause,
	"RaiseKeyCostThisTurn":         narratedByCause,
	"SetDamageImmune":              narratedByCause,
	"SetSideDamageImmune":          narratedByCause,
	"SetStatOverride":              narratedByCause,
	"SkipForgePhaseNextTurn":       narratedByCause,

	// The lasting and scheduled registries (ADR 0007, ADR 0013). Registering an
	// effect is not an outcome either: it narrates when it fires — LastingDraw,
	// AemberGiven, AemberGivenAfterForging, or whatever the scheduled effect does
	// in its window.
	"AddLasting":             narratedByCause,
	"AddLastingAlsoTriggers": narratedByCause,
	"ScheduleAtEndOfTurn":    narratedByCause,
	"ScheduleOnLeave":        narratedByCause,

	// Cancellations, each recorded by the effect immediately before it cancels:
	// KeyForgePrevented and FightCancelled.
	"CancelCurrentFight": narratedByCause,
	"CancelCurrentForge": narratedByCause,

	// Fight bookkeeping read and cleared by the combat step, whose own entries
	// (Fought, DamageTaken) narrate where the damage landed.
	"SetFightDamageRedirect": narratedByCause,

	// The shuffle the effect asks for; every caller records DeckShuffled right
	// after it (effect_shuffle.go, effect_search.go, effect_deck.go).
	"Shuffle": narratedByCause,
	// BeginShuffleBatch opens the window EndShuffleBatch narrates as one line.
	"BeginShuffleBatch": narratedByCause,

	// Dispatchers, which narrate through what they run rather than themselves:
	// the moves inside a simultaneous batch, the abilities a triggered card
	// resolves, and the victim's "After Æmber Is Stolen From You" abilities. Each
	// resolves whole effects, whose own entries (AemberStolen and the damage,
	// draws, or bars the ability then applies) are the narration. Classifying a
	// dispatcher as narrating directly would misattribute: the audit would name
	// the dispatcher for a silent change made three frames inside it.
	"EmitAemberStolenFrom":   narratedByCause,
	"Simultaneously":         narratedByCause,
	"TriggerAbilityOf":       narratedByCause,
	"TriggerAbilityOfRooted": narratedByCause,

	// A card mid-play redirecting itself is silent; the arrival it redirects to
	// is narrated when the play completes (CardPutIntoArchives, CardPurged).
	"RedirectResolvingCard": narratedByCause,

	// ---- narratedNever ----

	// Reads that sit on a mutating role because they answer a question only the
	// mutation cares about. They change nothing, so there is nothing to narrate.
	"AtRuleOfSix":        narratedNever,
	"ProtectedByTaunt":   narratedNever,
	"StolenAemberCaptor": narratedNever,

	// The Rule-of-Six tally (ADR 0043): a usage pool per card name, which no
	// player can see and no card reads back.
	"RecordUsage": narratedNever,

	// Private deck order. Both move a card the controller looked at within their
	// own hidden deck, so narrating them would leak the deck to the opponent —
	// PutDeckCardOnBottom's doc says so outright.
	"PutDeckCardOnBottom": narratedNever,
	"SetDeckTop":          narratedNever,
}
