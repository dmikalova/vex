package engine

import (
	"fmt"
	"strings"
)

// PoolAember gates on one player's Æmber pool: Player names whose pool (Controller
// or Opponent), Is the comparison, and Amount the threshold it compares against
// (unused by the relative MoreThanYou / MoreThanOpponent comparisons, which compare
// the two pools). It replaces the mirror opponent-pool / your-pool conditions with
// one node.
type PoolAember struct {
	Player Player
	Is     Comparison
	Amount int
}

// validate requires a pool-owning player and a named comparison, and ties each
// relative comparison to the side it reads from.
func (c PoolAember) validate() error {
	if c.Player != Controller && c.Player != Opponent {
		return fmt.Errorf("PoolAember: Player must be Controller or Opponent")
	}
	switch c.Is {
	case AtLeast, AtMost, Exactly, Even, Odd:
		return nil
	case MoreThanYou:
		if c.Player != Opponent {
			return fmt.Errorf("PoolAember: MoreThanYou requires Player Opponent")
		}
		return nil
	case MoreThanOpponent:
		if c.Player != Controller {
			return fmt.Errorf("PoolAember: MoreThanOpponent requires Player Controller")
		}
		return nil
	default:
		return fmt.Errorf(
			"PoolAember: Is must be AtLeast, AtMost, Exactly, Even, Odd, " +
				"MoreThanYou, or MoreThanOpponent",
		)
	}
}

// CondText renders the condition, e.g. "if your opponent has 7 Æmber or more" or
// "if you have more Æmber than your opponent".
func (c PoolAember) CondText() string {
	switch c.Is {
	case MoreThanYou:
		return "if your opponent has more Æmber than you"
	case MoreThanOpponent:
		return "if you have more Æmber than your opponent"
	case Even:
		if c.Player == Opponent {
			return "if your opponent has an even amount of Æmber"
		}
		return "if you have an even amount of Æmber"
	case Odd:
		if c.Player == Opponent {
			return "if your opponent has an odd amount of Æmber"
		}
		return "if you have an odd amount of Æmber"
	}
	if c.Player == Opponent {
		switch {
		case c.Is == Exactly && c.Amount == 0:
			return "if your opponent has no Æmber"
		case c.Is == Exactly:
			return fmt.Sprintf("if your opponent has exactly %d Æmber", c.Amount)
		case c.Is == AtMost:
			return fmt.Sprintf("if your opponent has %d Æmber or fewer", c.Amount)
		default:
			return fmt.Sprintf("if your opponent has %d Æmber or more", c.Amount)
		}
	}
	switch {
	case c.Is == Exactly && c.Amount == 0:
		return "if you have no Æmber"
	case c.Is == Exactly:
		return fmt.Sprintf("if you have exactly %d Æmber", c.Amount)
	case c.Is == AtMost:
		return fmt.Sprintf("if you have %d Æmber or fewer", c.Amount)
	default:
		return fmt.Sprintf("if you have %d Æmber or more", c.Amount)
	}
}

// Met reports whether the named player's pool satisfies the comparison.
func (c PoolAember) Met(ctx *EffectContext) bool {
	mine := ctx.Resolver.Aember(ctx.PlayerFor(c.Player))
	switch c.Is {
	case Exactly:
		return mine == c.Amount
	case AtMost:
		return mine <= c.Amount
	case Even:
		return mine%2 == 0
	case Odd:
		return mine%2 != 0
	case MoreThanYou, MoreThanOpponent:
		other := ctx.Controller
		if c.Player == Controller {
			other = ctx.Opponent()
		}
		return mine > ctx.Resolver.Aember(other)
	default:
		return mine >= c.Amount
	}
}

// ControlsMoreCreatures is met while the controller has more creatures in play
// than the opponent. Trait, when set, restricts the comparison to creatures with
// that trait (Pismire compares Mutant counts). It is the excess-creature Count
// read as a threshold: "more than the opponent" is an excess of at least one.
type ControlsMoreCreatures struct {
	Trait Trait
}

// excess is the Count this condition is a threshold on — how many more creatures
// the controller has than the opponent, narrowed to Trait when set. It also
// supplies the counted noun both wordings repeat.
func (c ControlsMoreCreatures) excess() ExcessCreatures {
	return ExcessCreatures{Player: Controller, Trait: c.Trait}
}

// CondText renders the condition, e.g. "if you control more Mutant creatures than
// your opponent".
func (c ControlsMoreCreatures) CondText() string {
	return "if you control more " + c.excess().filter().noun() +
		"s than your opponent"
}

// symmetricCondText renders the board-wide third-person form a
// ConditionalPlayBar needs, e.g. "has more creatures in play than their
// opponent" (Quixxle Stone).
func (c ControlsMoreCreatures) symmetricCondText() string {
	return "has more " + c.excess().filter().noun() +
		"s in play than their opponent"
}

// Met reports whether the controller has more creatures in play than the opponent.
func (c ControlsMoreCreatures) Met(ctx *EffectContext) bool {
	return CountIs{Count: c.excess(), Is: AtLeast, Amount: 1}.Met(ctx)
}

// ControlsNamed is met when the controller has a card of a given printed name in
// play — Hyde draws an extra card while it controls Velum.
type ControlsNamed struct {
	Name string
}

// CondText renders the condition, e.g. "if you control Velum".
func (c ControlsNamed) CondText() string {
	return "if you control " + c.Name
}

// Met reports whether the controller has the named card in play.
func (c ControlsNamed) Met(ctx *EffectContext) bool {
	return InPlay{Player: Controller, Name: c.Name}.Met(ctx)
}

// Overwhelmed reports whether the controller is overwhelmed — their opponent
// controls more creatures than they do. "Overwhelmed" is the keyword form of that
// board state; Numquid the Fair repeats its destruction while overwhelmed.
type Overwhelmed struct{}

// CondText renders the condition.
func (Overwhelmed) CondText() string { return "if you are overwhelmed" }

// Met reports whether the opponent controls more creatures than the controller.
func (Overwhelmed) Met(ctx *EffectContext) bool {
	return CountIs{
		Count:  ExcessCreatures{Player: Opponent},
		Is:     AtLeast,
		Amount: 1,
	}.Met(ctx)
}

// HousesRepresented is met when the distinct houses represented among a chosen
// set of in-play cards compare (Is) to Amount — Galactic Census pays out more as
// more houses share the board. Among carries no Max, so the raw house count is
// compared.
type HousesRepresented struct {
	Among  HousesAmong
	Is     Comparison
	Amount int
}

// validate requires a comparison the condition supports.
func (c HousesRepresented) validate() error {
	switch c.Is {
	case AtLeast, AtMost, Exactly:
		return nil
	default:
		return fmt.Errorf("HousesRepresented: Is must be AtLeast, AtMost, or Exactly")
	}
}

// Met compares the surveyed house count against Amount by Is.
func (c HousesRepresented) Met(ctx *EffectContext) bool {
	n := c.Among.Value(ctx)
	switch c.Is {
	case AtMost:
		return n <= c.Amount
	case Exactly:
		return n == c.Amount
	default:
		return n >= c.Amount
	}
}

// CondText renders the condition, e.g. "if there are 3 or more houses represented
// among creatures in play".
func (c HousesRepresented) CondText() string {
	qty := fmt.Sprintf("%d or more", c.Amount)
	switch c.Is {
	case AtMost:
		qty = fmt.Sprintf("%d or fewer", c.Amount)
	case Exactly:
		qty = fmt.Sprintf("exactly %d", c.Amount)
	}
	return fmt.Sprintf("if there are %s houses represented among %s", qty, c.Among.scope())
}

// CounterInPlay is met while at least one card in play carries a generic counter
// of Kind — Wretched Doll destroys every doom-marked creature when there is one,
// and otherwise marks a fresh one.
type CounterInPlay struct {
	// Kind is the counter to look for.
	Kind CounterKind
}

// CondText renders the condition.
func (c CounterInPlay) CondText() string {
	return "if there is a " + c.Kind.noun() + " in play"
}

// Met reports whether any creature in either battleline carries the counter.
func (c CounterInPlay) Met(ctx *EffectContext) bool {
	for player := 0; player < 2; player++ {
		for _, id := range ctx.Resolver.Battleline(player) {
			if ctx.Resolver.CountersOn(id, c.Kind) > 0 {
				return true
			}
		}
	}
	return false
}

// NamedCardPurged is met by whether a card of a given name sits in the
// controller's purge pile — Igon the Terrible destroys itself unless Igon the
// Green has already been purged (wrap in Not for "has not been purged"). It names
// the other card by its printed name, not the source.
type NamedCardPurged struct {
	// Name is the card name to look for in the purge pile.
	Name string
}

// CondText renders the condition naming the card it looks for.
func (c NamedCardPurged) CondText() string {
	return "if " + c.Name + " has been purged"
}

// negatedText renders the not-purged clause a Not wrapper prints.
func (c NamedCardPurged) negatedText() string {
	return "if " + c.Name + " has not been purged"
}

// Met reports whether a card of the name is in the controller's purge pile.
func (c NamedCardPurged) Met(ctx *EffectContext) bool {
	for _, id := range ctx.Resolver.Purge(ctx.Controller) {
		if (CardFilter{Name: c.Name}).admits(ctx.Resolver, id) {
			return true
		}
	}
	return false
}

// ForgedKey is the condition on whether a player forged a key in a given window —
// this turn (Smiling Ruth) or on their own previous turn (Tendrils of Pain, Key
// Hammer). It reads the turn history rather than the running key total, so a key
// forged several turns ago does not keep the condition true. Wrap in Not for "has
// not forged a key" (Nightforge).
type ForgedKey struct {
	Player   Player
	Previous bool
}

// validate requires the condition to name whose key it asks about.
func (c ForgedKey) validate() error {
	if !c.Player.valid() {
		return fmt.Errorf("ForgedKey: Player must be set")
	}
	return nil
}

// subject names the player and their possessive, e.g. "you"/"your".
func (c ForgedKey) subject() (string, string) {
	if c.Player == Opponent {
		return "your opponent", "their"
	}
	return "you", "your"
}

// window names the turn the condition asks about.
func (c ForgedKey) window(possessive string) string {
	if c.Previous {
		return "during " + possessive + " previous turn"
	}
	return "this turn"
}

// CondText renders the clause, e.g. "if your opponent forged a key during their
// previous turn".
func (c ForgedKey) CondText() string {
	subject, possessive := c.subject()
	return fmt.Sprintf("if %s forged a key %s", subject, c.window(possessive))
}

// negatedText renders the not-forged clause a Not wrapper prints.
func (c ForgedKey) negatedText() string {
	subject, possessive := c.subject()
	return fmt.Sprintf("if %s have not forged a key %s", subject, c.window(possessive))
}

// Met reports whether the named player forged at least one key in the window.
func (c ForgedKey) Met(ctx *EffectContext) bool {
	return ctx.Resolver.TurnHistory(ctx.PlayerFor(c.Player), c.stat()) > 0
}

// stat picks the tally the window corresponds to.
func (c ForgedKey) stat() TurnStat {
	if c.Previous {
		return KeysForgedLastTurn
	}
	return KeysForgedThisTurn
}

// HasMoreForgedKeys is met when Player has forged strictly more keys than the
// other player — Hugger Mugger steals only when the opponent is ahead on keys.
type HasMoreForgedKeys struct {
	// Player is the side that must be ahead on forged keys.
	Player Player
}

// CondText renders the clause, naming whose keys lead.
func (c HasMoreForgedKeys) CondText() string {
	if c.Player == Opponent {
		return "if your opponent has more forged keys than you"
	}
	return "if you have more forged keys than your opponent"
}

// Met reports whether Player's forged-key count exceeds the other player's.
func (c HasMoreForgedKeys) Met(ctx *EffectContext) bool {
	mine := ctx.PlayerFor(c.Player)
	return ctx.Resolver.Keys(mine) > ctx.Resolver.Keys(1-mine)
}

// KeyColorForged is met while the named player has forged a key of a given colour
// — The Red Baron gains a reap while your red key is forged, and gains elusive
// while your opponent's red key is forged.
type KeyColorForged struct {
	// Player is whose forged keys to look at: Controller or Opponent.
	Player Player
	// Color is the key colour that must be among that player's forged keys.
	Color KeyColor
}

// CondText renders the clause, e.g. "if your red key is forged" or "if your
// opponent's red key is forged".
func (c KeyColorForged) CondText() string {
	possessive := "your"
	if c.Player == Opponent {
		possessive = "your opponent's"
	}
	return fmt.Sprintf("if %s %s key is forged", possessive, strings.ToLower(c.Color.String()))
}

// Met reports whether the named player has forged a key of Color.
func (c KeyColorForged) Met(ctx *EffectContext) bool {
	for _, col := range ctx.Resolver.KeyColors(ctx.PlayerFor(c.Player)) {
		if col == c.Color {
			return true
		}
	}
	return false
}
