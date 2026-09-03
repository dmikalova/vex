package engine

import (
	"fmt"
	"strings"
)

// This file holds the authoring effects for lasting "for the remainder of the
// turn" behavior: ForRemainderOfTurn installs a reaction (do something after an
// event), Instead installs a replacement (change an event's own outcome). Both are
// thin: they translate a composed effect into a flat lasting record via AddLasting
// (see game_lasting.go for how the records are fired and queried).

// ForRemainderOfTurn installs a reaction that runs for the rest of the controller's
// turn each time On occurs — Full Moon (gain Æmber when you play a creature),
// Charge! (deal damage when you play a creature), Crystal Hive (gain Æmber after a
// creature reaps), Library Access (draw a card each time you play another card).
// Do is the effect that runs each time; the supported effects are GainAember
// (crediting the controller), DealDamage to an enemy creature, CaptureAember, and
// Draw. The card that installs the reaction never triggers it itself, so an event
// phrased as "another card" excludes the play that armed it.
type ForRemainderOfTurn struct {
	On Event
	Do Effect
}

// validate rejects a non-reaction event or a Do the reaction cannot carry.
func (e ForRemainderOfTurn) validate() error {
	if !e.On.isReaction() {
		return fmt.Errorf("ForRemainderOfTurn: On must be a reaction event")
	}
	if _, _, ok := lastingActionOf(e.Do); !ok {
		return fmt.Errorf("ForRemainderOfTurn: unsupported Do %T", e.Do)
	}
	if d, ok := e.Do.(DealDamage); ok && d.Target != (Target{Kind: TargetChosenEnemyCreature}) {
		return fmt.Errorf("ForRemainderOfTurn: DealDamage must target an enemy creature")
	}
	return validateEffect(e.Do)
}

// Text renders the effect, e.g. "for the remainder of the turn, each time you play
// a creature, gain 1 Æmber".
func (e ForRemainderOfTurn) Text() string {
	return "for the remainder of the turn, " + e.On.clause() + ", " + e.Do.Text()
}

// Resolve registers the reaction on the controller for the rest of the turn. On
// EventCardPlayed the installing card is excepted, so "each time you play another
// card" does not count the play that armed it.
func (e ForRemainderOfTurn) Resolve(ctx *EffectContext) {
	action, amount, _ := lastingActionOf(e.Do)
	ctx.Resolver.AddLasting(LastingEffect{
		On:         e.On,
		Do:         action,
		Controller: int8(ctx.Controller),
		Amount:     int8(amount),
		Except:     ctx.Source,
		HasExcept:  e.On == EventCardPlayed,
	})
}

// lastingActionOf maps a reaction's Do effect to the flat action and amount stored
// in the registry, reporting whether the effect is one the registry can carry.
func lastingActionOf(e Effect) (lastingAction, int, bool) {
	switch d := e.(type) {
	case GainAember:
		return actGainAember, d.Amount, true
	case DealDamage:
		return actDealDamage, d.Amount, true
	case CaptureAember:
		return actCapture, d.Amount, true
	case Draw:
		return actDraw, d.Amount, true
	}
	return 0, 0, false
}

// Replacement is a lasting change to an event's own outcome, used by Instead.
type Replacement uint8

const (
	// replacementUnset is the invalid zero value: an Instead must name its
	// replacement rather than leave it unset.
	replacementUnset Replacement = iota
	// Steal replaces gaining Æmber with stealing that much from the opponent.
	Steal
	// Capture replaces adding Æmber to a pool with the source creature capturing it
	// (Ether Spider). It is applied continuously at the add-to-pool site, not through
	// the turn-scoped lasting registry, so it has no lastingAction.
	Capture
)

// valid reports whether r names a real replacement (not the unset zero value).
func (r Replacement) valid() bool { return r != replacementUnset }

// action maps the replacement to the flat action stored in the registry.
func (Replacement) action() lastingAction { return actSteal }

// text renders the replacement clause, e.g. "steal the same amount".
func (Replacement) text() string { return "steal the same amount" }

// Replace is a continuous replacement an Upgrade applies to a game event for its
// host while attached: when the event When would happen to the host, the effect
// With resolves in its place. Unlike the turn-scoped Instead — a flat outcome swap
// kept in the pointerless game state — a Replace lives in the card definition, so
// its With is a full effect tree. Armageddon Cloak replaces its host's destruction
// (EventCreatureDestroyed) with "fully heal it and destroy Armageddon Cloak", the
// self-destruction spelled out as an effect rather than implied by the event site.
type Replace struct {
	When Event
	With Effect
}

// valid reports whether a replacement is set, distinguishing a StaticModifier that
// carries a Replace from the zero value that carries none.
func (r Replace) valid() bool { return r.When != eventUnset }

// validate surfaces a configuration error in the replacement effect, ignoring the
// zero value (a StaticModifier with no replacement).
func (r Replace) validate() error {
	if !r.valid() {
		return nil
	}
	return validateEffect(r.With)
}

// Instead installs a replacement that, for the rest of the controller's turn,
// changes the outcome of the event Of before it happens — Dimension Door replaces
// gaining Æmber from reaping with stealing it. As a plain {Of, With} value it also
// describes a continuous replacement a card applies while in play (CardDefinition.Replaces,
// Ether Spider capturing Æmber added to its opponent's pool); in that use it is read,
// never resolved. Player scopes an event that names a pool — which player's pool the
// replacement watches (Ether Spider watches its Opponent's).
type Instead struct {
	Of     Event
	With   Replacement
	Player Player
}

// valid reports whether a replacement is set, distinguishing a card that carries a
// continuous Instead from the zero value that carries none.
func (e Instead) valid() bool { return e.Of != eventUnset }

// validate rejects an Of that is not a replacement event or an unset With, and
// requires the pool-scoping Player when Of names a pool.
func (e Instead) validate() error {
	if e.Of.isReaction() {
		return fmt.Errorf("Instead: Of must be a replacement event")
	}
	if !e.With.valid() {
		return fmt.Errorf("Instead: replacement must be set")
	}
	if e.Of == EventAemberAddedToPool && !e.Player.valid() {
		return fmt.Errorf("Instead: a pool event needs a Player to scope it")
	}
	return nil
}

// Text renders the effect, e.g. "for the remainder of the turn, instead of gaining
// Æmber from reaping, steal the same amount".
func (e Instead) Text() string {
	return "for the remainder of the turn, instead of " + e.Of.gerund() + ", " + e.With.text()
}

// Resolve registers the replacement on the controller for the rest of the turn.
func (e Instead) Resolve(ctx *EffectContext) {
	ctx.Resolver.AddLasting(LastingEffect{
		On:         e.Of,
		Do:         e.With.action(),
		Controller: int8(ctx.Controller),
	})
}

// NextPlayed makes the next card its controller plays this turn that matches its
// filters enter play with EntersPlay applied to it — Blypyp readying the next Mars
// creature, Soft Landing readying the next creature or artifact. It registers a
// one-shot reaction that applies the effect to that card and then removes itself;
// the turn's end clears it if no matching card is played. EntersPlay must be an
// enter-play effect the flat registry can carry (Ready).
//
// Of narrows the card to one house, and is optional. Type narrows it to one card
// type; AnyType means "creature or artifact", the two types that stay in play.
type NextPlayed struct {
	Of         House
	Type       CardType
	EntersPlay Effect
}

// validate requires a card type and rejects an EntersPlay effect the registry
// cannot carry.
func (e NextPlayed) validate() error {
	if e.Type == TypeUnset {
		return fmt.Errorf("NextPlayed: type must be set")
	}
	if _, ok := enterActionOf(e.EntersPlay); !ok {
		return fmt.Errorf("NextPlayed: unsupported EntersPlay %T", e.EntersPlay)
	}
	return validateEffect(e.EntersPlay)
}

// Text renders the effect, e.g. "the next Mars creature you play this turn enters
// play ready" or "the next creature or artifact you play this turn enters play
// ready".
func (e NextPlayed) Text() string {
	noun := "creature or artifact"
	if e.Type != AnyType {
		noun = strings.ToLower(e.Type.String())
	}
	if e.Of != HouseNone {
		noun = e.Of.String() + " " + noun
	}
	return fmt.Sprintf(
		"the next %s you play this turn enters play %s",
		noun,
		enterStateWord(e.EntersPlay),
	)
}

// Resolve registers the one-shot enter-play reaction on the controller.
func (e NextPlayed) Resolve(ctx *EffectContext) {
	action, _ := enterActionOf(e.EntersPlay)
	ctx.Resolver.AddLasting(LastingEffect{
		On:         EventCardEntersPlay,
		Do:         action,
		Controller: int8(ctx.Controller),
		House:      e.Of,
		Type:       e.Type,
		Once:       true,
	})
}

// enterActionOf maps an enter-play effect to the flat lasting action that applies
// it to the next played creature, reporting whether the registry can carry it.
func enterActionOf(e Effect) (lastingAction, bool) {
	if _, ok := e.(Ready); ok {
		return actReadyPlayed, true
	}
	return 0, false
}
