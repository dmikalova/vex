package engine

import "fmt"

// This file and its effect_*.go / target.go siblings make up the card effect
// "AST": the small tree of nodes that both prints a card's rules text and
// carries it out. Each node type lives with related nodes in a file grouped by
// mechanic (effect_aember.go, effect_stun.go, ...), and every type's doc comment
// explains the mechanic in rulebook terms before the code shows how it is
// modelled. They all live in package engine because effects reach deep into the
// engine (dealing damage, destroying, drawing, choosing creatures); splitting
// them into a separate package would force the whole engine to be exported.

// Effect is one node in a card's effect tree. Each node knows how to render
// itself to English (Text) and how to carry itself out against a live game
// (Resolve). Because the same node does both, a card's printed rules text can
// never drift from what the card actually does.
type Effect interface {
	Text() string
	Resolve(ctx *EffectContext)
}

// validator is implemented by effects that can be misconfigured in a card
// definition. NewCard checks it when a card is built so a bad definition fails at
// startup instead of silently misbehaving when the ability resolves.
type validator interface {
	validate() error
}

// validateEffect returns an effect's configuration error, or nil if it has none.
// Composite effects (Sequence, Conditional) implement validator by descending
// into their children through this helper.
func validateEffect(e Effect) error {
	if v, ok := e.(validator); ok {
		return v.validate()
	}
	return nil
}

// errUnsetPlayer is the configuration error a player-taking effect returns when
// its Player was left as the invalid zero value.
func errUnsetPlayer(effect string) error {
	return fmt.Errorf("%s: player must be set (Controller, Opponent, or EachPlayer)", effect)
}

// errUnsetTarget is the configuration error a target-taking effect returns when
// its Target was left as the invalid zero value.
func errUnsetTarget(effect string) error {
	return fmt.Errorf("%s: target must be set", effect)
}

// errUnsetDuration is the configuration error a timed effect returns when its
// Duration was left as the invalid zero value.
func errUnsetDuration(effect string) error {
	return fmt.Errorf("%s: duration must be set", effect)
}

// EffectContext carries the state an effect needs while resolving. It exposes the
// game only through a Resolver, so an effect can inspect and change the game only
// via that interface — never by reaching into the state directly. Cards are
// referenced by LocalID, keeping the context flat.
type EffectContext struct {
	Resolver   Resolver
	Source     LocalID // the card whose ability is resolving
	Controller int     // the player who controls the ability
	// It is the card in context: the creature that fired a trigger, or a card an
	// earlier effect in this resolution put in focus (a revealed or discarded top
	// card). HasIt reports whether one is set.
	It    LocalID
	HasIt bool
	// Upgrade is the attached Upgrade whose own ability is resolving, when one is —
	// an Upgrade's "Play:" fires with Source set to its host creature, so Upgrade
	// lets that effect still refer to the Upgrade itself (e.g. as the source of a
	// control change that lasts until the Upgrade leaves play).
	Upgrade LocalID
	// ChosenHouse is a house picked by a ChooseHouseThen, read by
	// Target.OfChosenHouse targets nested inside it.
	ChosenHouse House
	// Produced holds the "... this way" tallies an effect records for a following
	// effect in the same resolution to read.
	Produced Produced
}

// Produced holds the tallies one effect records for a following effect in the
// same resolution to consume — the KeyForge "... this way" counts (creatures
// healed, cards revealed, cards destroyed, houses discarded). Each ability
// resolves with a fresh EffectContext, so these never leak between abilities;
// grouping them keeps this producer/consumer channel in one place as new "this
// way" effects are added (add the tally here, set it in the producer, read it in
// the consuming Count/Condition).
type Produced struct {
	// Healed is how many creatures the most recent Heal healed, read by a
	// CreaturesHealed count in a following effect of the same resolution.
	Healed int
	// DamageHealed is how much damage the most recent Heal actually removed, read by
	// a DamageHealed count in a following effect of the same resolution (Guardian
	// Demon deals that much damage on).
	DamageHealed int
	// Revealed is how many cards the most recent Reveal showed, read by a
	// CardsRevealed count in a following effect of the same resolution.
	Revealed int
	// Destroyed[p] is how many cards player p controlled that this resolution has
	// destroyed, read whole by CardsDestroyed / CardsDestroyedFewerThan and per
	// side by CreaturesDestroyedThisWay (Hecatomb pays each player for their
	// own dead).
	Destroyed [2]int
	// Purged is how many cards the most recent purge removed, read by a CardsPurged
	// count in a following effect of the same resolution (One Last Job steals for
	// each creature it purged).
	Purged int
	// Discarded holds the cards a DiscardTopOfEachDeck discarded, read by a
	// following ForEachDiscarded that acts on each (Bonkers Killing Machine
	// destroys a creature or artifact of each discarded card's house).
	Discarded []LocalID
	// Moved[p] is how many cards player p controlled that a PutFromPlay took out of
	// play — sent home rather than destroyed — this resolution, read by
	// CreaturesShuffledIntoDeckThisWay (Mating Season).
	Moved [2]int
	// AemberLost[p] is how much Æmber a LoseAember has taken from player p's pool
	// this resolution, read by AemberLostThisWay (Shatter Storm drains the opponent
	// for triple what its controller lost).
	AemberLost [2]int
}

// TotalDestroyed is how many cards this resolution has destroyed, both sides
// together.
func (p Produced) TotalDestroyed() int { return p.Destroyed[0] + p.Destroyed[1] }

// Opponent returns the absolute index of the controller's opponent.
func (ctx *EffectContext) Opponent() int { return 1 - ctx.Controller }

// PlayerFor resolves a relative Player (Controller or Opponent) to an absolute
// player index. Use it for a Player value held by an effect (e.g. e.Player); for
// the two fixed players prefer the plainer ctx.Controller and ctx.Opponent().
func (ctx *EffectContext) PlayerFor(p Player) int {
	switch p {
	case Opponent:
		return ctx.Opponent()
	case Controller, EachPlayer:
		return ctx.Controller
	case ItsOwner:
		return ctx.Resolver.Owner(ctx.It)
	case ItsOpponent:
		return 1 - ctx.Resolver.Controller(ctx.It)
	default:
		panic("engine: effect has no player set (playerUnset)")
	}
}

// ChooseCreature asks the controlling player to pick one creature from candidates,
// attributing the prompt to this ability's source card. It is the common form of
// Resolver.ChooseCreature; call the Resolver directly only when a different player
// makes the choice (e.g. the owner of a creature being used to fight).
func (ctx *EffectContext) ChooseCreature(prompt string, candidates []LocalID) (LocalID, bool) {
	return ctx.Resolver.ChooseCreature(ctx.Controller, ctx.Source, prompt, candidates)
}

// ChooseCard asks the controlling player to pick one card from candidates,
// attributing the prompt to this ability's source card. It is the common form of
// Resolver.ChooseCard; call the Resolver directly only when a different player
// makes the choice.
func (ctx *EffectContext) ChooseCard(prompt string, candidates []LocalID) (LocalID, bool) {
	return ctx.Resolver.ChooseCard(ctx.Controller, ctx.Source, prompt, candidates)
}

// ChooseOption asks the controlling player to pick one labeled option, attributing
// the prompt to this ability's source card.
func (ctx *EffectContext) ChooseOption(prompt string, options []string) int {
	return ctx.Resolver.ChooseOption(ctx.Controller, ctx.Source, prompt, options)
}

// ChooseCardOptional asks the controlling player to pick one card from candidates
// or to decline, attributing the prompt to this ability's source card. Use it for
// every "you may" and "up to N" choice so the player picks a card rather than a
// name off a list; a sole candidate is offered, not forced.
func (ctx *EffectContext) ChooseCardOptional(
	prompt string,
	candidates []LocalID,
) (LocalID, bool) {
	return ctx.Resolver.ChooseCardOptional(ctx.Controller, ctx.Source, prompt, candidates)
}

// OrderByChoice asks the controlling player to arrange ids into a resolution order.
func (ctx *EffectContext) OrderByChoice(prompt string, ids []LocalID) []LocalID {
	return ctx.Resolver.OrderByChoice(ctx.Controller, prompt, ids)
}

// Player selects which player an effect targets, relative to the card's
// controller: Controller is the player who controls the card, Opponent is their
// opponent, and EachPlayer is both. Every effect names its player explicitly:
// there is no default, so the zero value is an invalid placeholder rejected when
// the card is built and when the effect resolves.
type Player int

const (
	// playerUnset is the invalid zero value: an effect must name its player
	// (Controller, Opponent, or EachPlayer) rather than leave it unset.
	playerUnset Player = iota
	// Controller is the player who controls the card/ability.
	Controller
	// Opponent is the controller's opponent.
	Opponent
	// EachPlayer is both players. It is meaningful only for effects that reach
	// everyone at once (e.g. a KeyCostChange on "each player's keys"); the
	// single-target effects use only Controller and Opponent.
	EachPlayer
	// ItsOwner is the owner of the creature currently in context (ctx.It) — the
	// "its owner" referent, for an effect that acts on the owner of a creature a
	// preceding clause touched (Gongoozle's damaged creature).
	ItsOwner
	// ItsOpponent is the opponent of the creature currently in context (ctx.It).
	// It is how one effect can reach a different pool for each side's creatures:
	// Pandemonium's "each undamaged creature captures 1 Æmber from its opponent"
	// takes from your pool for the enemy's creatures and from theirs for yours.
	ItsOpponent
)

// valid reports whether p names a real player (not the unset zero value).
func (p Player) valid() bool { return p != playerUnset }

// SelfName is a placeholder an effect's text uses to refer to its own source
// card; RenderCardText and the game log substitute it with the card's name so
// text like "{self} captures 1 Æmber" prints as "Charette captures 1 Æmber".
const SelfName = "{self}"

// UpgradeName is a placeholder an Upgrade's own Play effect uses when the text
// must name the Upgrade rather than its host creature.
const UpgradeName = "{upgrade}"
