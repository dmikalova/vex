package engine

import (
	"errors"
	"fmt"
)

// This file is the manual/debug half of the Command seam: a single place to PERFORM
// one manual force-edit (ApplyManual), the counterpart to command_action.go's
// ApplyAction for the 11 real root actions. A manual command deliberately bypasses
// the rules (ADR 0039), so it maps onto the engine's Manual* methods, which perform
// no checks. Recording manual edits as Commands lets a match holding force-edits
// replay from its command log like any other.

// ErrNotManualAction is returned by ApplyManual for a Command whose Kind is a root
// action or a choice answer rather than one of the manual force-edits it performs.
var ErrNotManualAction = errors.New("command is not a manual action")

// ApplyManual performs one manual/debug Command, dispatching to the engine's
// Manual* method for its Kind. resolveCard maps a card name to its definition for
// CommandManualAddCard — the one manual edit that needs the card pool the engine
// deliberately does not hold (cards import engine, not the reverse), so the caller
// that owns the pool injects the lookup. Every other kind ignores it. A command
// that is not a manual edit returns ErrNotManualAction, and an add naming a card the
// pool does not know returns an error so a diverged replay fails where it diverged.
func (g *Game) ApplyManual(
	cmd Command,
	resolveCard func(name string) (CardDefinition, bool),
) error {
	switch cmd.Kind {
	case CommandSetManual:
		g.SetManual(cmd.Left)
	case CommandManualMove:
		g.ManualMove(cmd.Card, ManualZone(cmd.Index))
	case CommandManualReady:
		g.ManualSetExhausted(cmd.Card, false)
	case CommandManualExhaust:
		g.ManualSetExhausted(cmd.Card, true)
	case CommandManualAttach:
		g.ManualAttachUnder(cmd.Card, cmd.Card2, cmd.Left)
	case CommandManualPlace:
		g.ManualPlaceInPlay(cmd.Card, cmd.Index)
	case CommandManualDetach:
		g.ManualDetachToHand(cmd.Card)
	case CommandManualAmber:
		g.ManualAddAmber(cmd.Player, cmd.Delta)
	case CommandManualUnforge:
		g.ManualUnforgeKey(cmd.Player)
	case CommandManualForgeColor:
		g.ManualForgeKeyColor(cmd.Player, KeyColor(cmd.Index))
	case CommandManualChains:
		g.ManualAddChains(cmd.Player, cmd.Delta)
	case CommandManualHouse:
		g.ManualSetActiveHouse(cmd.House)
	case CommandManualAddCard:
		def, ok := resolveCard(cmd.Name)
		if !ok {
			return fmt.Errorf("manual command names card %q, which is not in the pool", cmd.Name)
		}
		g.ManualAddCard(def, cmd.Player)
	default:
		return ErrNotManualAction
	}
	return nil
}
