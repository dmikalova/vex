package engine

import "fmt"

// ArchivedCreaturesShareHouse is met when the creatures a preceding
// ArchiveFromPlay set aside (ctx.Produced.Archived) all belong to one house —
// Code Monkey gains 2 Æmber when the neighbors it archived share a house. Fewer
// than two creatures cannot share a house, so it is not met.
type ArchivedCreaturesShareHouse struct{}

// CondText renders the condition naming the just-archived creatures.
func (ArchivedCreaturesShareHouse) CondText() string {
	return "if those creatures share a house"
}

// Met reports whether every archived creature belongs to the same house.
func (ArchivedCreaturesShareHouse) Met(ctx *EffectContext) bool {
	ids := ctx.Produced.Archived
	if len(ids) < 2 {
		return false
	}
	first := ctx.Resolver.House(ids[0])
	for _, id := range ids[1:] {
		if ctx.Resolver.House(id) != first {
			return false
		}
	}
	return true
}

// MovedAnyAember is met when a preceding MoveAember relocated at least one Æmber
// this resolution (ctx.Produced.AemberMoved) — Shadowsaurus takes control of the
// enemy creature it emptied only when there was Æmber to move.
type MovedAnyAember struct{}

// CondText renders the condition as a back-reference to the Æmber just moved.
func (MovedAnyAember) CondText() string {
	return "if you moved any \u00c6mber this way"
}

// Met reports whether the most recent MoveAember moved any Æmber.
func (MovedAnyAember) Met(ctx *EffectContext) bool {
	return ctx.Produced.AemberMoved > 0
}

// FirstCreaturePlayedThisTurn is met when the card in context (ctx.It, the
// creature that fired the trigger) is the first creature its player played this
// turn — Speed Sigil readies it. It is a once-per-turn charge that needs no state
// of its own: the turn's play record already says whether the charge is spent, and
// the record is cleared when the next turn begins.
//
// A creature put into play by an effect rather than played never matches, so it
// neither benefits nor spends the charge.
type FirstCreaturePlayedThisTurn struct{}

// CondText renders the condition.
func (FirstCreaturePlayedThisTurn) CondText() string {
	return "if it is the first creature played this turn"
}

// Met reports whether the context card is the earliest creature in the active
// player's plays this turn.
func (FirstCreaturePlayedThisTurn) Met(ctx *EffectContext) bool {
	if !ctx.HasIt {
		return false
	}
	for _, id := range ctx.Resolver.PlayedThisTurn(ctx.Resolver.ActivePlayer()) {
		if ctx.Resolver.TypeOf(id) == Creature {
			return id == ctx.It
		}
	}
	return false
}

// NoCreaturesPlayedThisTurn is met when the controller has not played any
// creatures during the current turn — Redlock pays out at end of turn only on a
// turn its controller played no creatures. A creature put into play by an effect
// rather than played does not count, so it does not spoil the payout.
type NoCreaturesPlayedThisTurn struct{}

// CondText renders the condition.
func (NoCreaturesPlayedThisTurn) CondText() string {
	return "if you did not play any creatures this turn"
}

// Met reports whether none of the controller's plays this turn were creatures.
func (NoCreaturesPlayedThisTurn) Met(ctx *EffectContext) bool {
	for _, id := range ctx.Resolver.PlayedThisTurn(ctx.Controller) {
		if ctx.Resolver.TypeOf(id) == Creature {
			return false
		}
	}
	return true
}

// ItIsYourTurn is met when the ability's controller is the active player —
// Jargogle plays the card under it when destroyed on its controller's turn, and
// archives it otherwise.
type ItIsYourTurn struct{}

// CondText renders the condition.
func (ItIsYourTurn) CondText() string { return "if it is your turn" }

// Met reports whether the controller is the active player.
func (ItIsYourTurn) Met(ctx *EffectContext) bool {
	return ctx.Resolver.ActivePlayer() == ctx.Controller
}

// ChoseHouse is met when the controller's active house is House. It is the
// condition behind an "After you choose <House> as your active house, ..."
// ability (Jehu the Bureaucrat): the AfterChooseHouse trigger fires for the
// active player as they pick their house, and this checks whether they picked
// the house the ability watches for.
type ChoseHouse struct {
	House House
}

// CondText renders the condition clause.
func (c ChoseHouse) CondText() string {
	return "you choose " + c.House.String() + " as your active house"
}

// Met reports whether the active house is the one the ability watches for.
func (c ChoseHouse) Met(ctx *EffectContext) bool {
	return ctx.Resolver.ActiveHouse() == c.House
}

// AemberStolenFromYou is met when the controller had Æmber stolen from them on
// their opponent's previous turn — Information Exchange steals more if it was.
type AemberStolenFromYou struct{}

// CondText renders the condition.
func (AemberStolenFromYou) CondText() string {
	return "if your opponent stole Æmber from you on their previous turn"
}

// Met reports whether any Æmber was stolen from the controller last turn.
func (AemberStolenFromYou) Met(ctx *EffectContext) bool {
	return ctx.Resolver.TurnHistory(ctx.Controller, AemberStolenFromLastTurn) > 0
}

// CreatureDestroyedThisTurn is met while at least one creature on the named side
// has been destroyed this turn: Opponent for an enemy creature (Foozle reaps for
// an extra Æmber once the opponent has lost one), Controller for a friendly one
// (Bonesaw enters play ready once one of yours has died).
type CreatureDestroyedThisTurn struct {
	// Player names whose creature must have died — Controller for a friendly
	// creature, Opponent for an enemy one.
	Player Player
}

// creatureDestroyedSides pairs each side the condition can name with the noun
// phrase its text uses and the turn tally that answers it, so the wording and the
// stat cannot drift apart. A side absent here is not a side this condition reads.
var creatureDestroyedSides = map[Player]struct {
	subject string
	stat    TurnStat
}{
	Controller: {"a friendly creature", FriendlyCreaturesDestroyed},
	Opponent:   {"an enemy creature", EnemyCreaturesDestroyed},
}

// validate requires a side the condition can read: a creature is destroyed from
// one player's board or the other's, so EachPlayer and the unset zero value are
// both rejected.
func (c CreatureDestroyedThisTurn) validate() error {
	if _, ok := creatureDestroyedSides[c.Player]; !ok {
		return fmt.Errorf(
			"CreatureDestroyedThisTurn: Player must be Controller or Opponent")
	}
	return nil
}

// CondText renders the condition, e.g. "if an enemy creature has been destroyed
// this turn".
func (c CreatureDestroyedThisTurn) CondText() string {
	return "if " + creatureDestroyedSides[c.Player].subject + " has been destroyed this turn"
}

// Met reports whether the named side has lost a creature this turn.
func (c CreatureDestroyedThisTurn) Met(ctx *EffectContext) bool {
	return ctx.Resolver.TurnHistory(ctx.Controller, creatureDestroyedSides[c.Player].stat) > 0
}

// UsedCreatureToReap is met while the controller has used a creature to reap at
// least once this turn — Bramble Lynx enters play ready once you have reaped.
type UsedCreatureToReap struct{}

// CondText renders the condition.
func (UsedCreatureToReap) CondText() string {
	return "if you have used a creature to reap this turn"
}

// Met reports whether the controller has reaped with a creature this turn.
func (UsedCreatureToReap) Met(ctx *EffectContext) bool {
	return ctx.Resolver.TurnHistory(ctx.Controller, CreaturesReapedThisTurn) > 0
}

// UsedCreatureToFight is met while the controller has used a creature to fight at
// least once this turn — Alaka enters play ready once you have fought.
type UsedCreatureToFight struct{}

// CondText renders the condition.
func (UsedCreatureToFight) CondText() string {
	return "if you have used a creature to fight this turn"
}

// Met reports whether the controller has fought with a creature this turn.
func (UsedCreatureToFight) Met(ctx *EffectContext) bool {
	return ctx.Resolver.TurnHistory(ctx.Controller, CreaturesFoughtThisTurn) > 0
}

// UsedNoCreatures is met while the controller has not used any creature this turn
// by any means — reaping, fighting, or using an Action ability — Sloth rewards a
// turn spent without using a creature.
type UsedNoCreatures struct{}

// CondText renders the condition.
func (UsedNoCreatures) CondText() string {
	return "if you did not use any creatures this turn"
}

// Met reports whether the controller has used no creature this turn.
func (UsedNoCreatures) Met(ctx *EffectContext) bool {
	return ctx.Resolver.TurnHistory(ctx.Controller, CreaturesUsedThisTurn) == 0
}

// FirstReapOfTurn is met when the reap in context is the first time a creature
// has reaped this turn — Aember Conduction Unit stuns only the first enemy
// creature to reap. It reads the reaping creature (ctx.It) so it asks about the
// active player's tally, which counts one once this reap has been tallied.
type FirstReapOfTurn struct{}

// CondText renders the condition.
func (FirstReapOfTurn) CondText() string {
	return "if it is the first time a creature has reaped this turn"
}

// Met reports whether exactly one creature has reaped this turn, the reaping
// creature in context being that one.
func (FirstReapOfTurn) Met(ctx *EffectContext) bool {
	if !ctx.HasIt {
		return false
	}
	return ctx.Resolver.TurnHistory(ctx.Resolver.Controller(ctx.It), CreaturesReapedThisTurn) == 1
}

// SourceFirstUseThisTurn is met when the current use of the source creature is
// its first this turn — Gladiodontus readies and enrages itself only the first
// time it is used. A creature is used when it reaps, fights, or fires an Action
// ability; the use is tallied before the Fight:/Reap: ability resolves, so the
// first use reads as a tally of one.
type SourceFirstUseThisTurn struct{}

// CondText renders the condition, naming the source card.
func (SourceFirstUseThisTurn) CondText() string {
	return "if this is the first time " + SelfName + " has been used this turn"
}

// Met reports whether the source creature's current use is its first this turn.
func (SourceFirstUseThisTurn) Met(ctx *EffectContext) bool {
	return ctx.Resolver.TimesUsedThisTurn(ctx.Source) == 1
}
