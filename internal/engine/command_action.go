package engine

import "errors"

// This file gives the 11 root actions a Command representation's other half: a
// single place to PERFORM one (ApplyAction). ADR 0039 defines a Command as "a root
// action or one answer to a choice", but the answer kinds in suspend.go left the
// root half unexpressed — so a driver that already speaks Command (internal/session)
// could answer a choice but not initiate a play. ApplyAction closes that: it maps
// each root Command onto the engine method the matching player click calls,
// mirroring the web client's proven dispatchRoot so that hand-rolled driver can
// fold into this one (ADR 0040). It is the seam the canonical turn loop applies a
// player's chosen action through.

// ErrNotRootAction is returned by ApplyAction for a Command whose Kind is a choice
// answer rather than one of the root actions it performs.
var ErrNotRootAction = errors.New("command is not a root action")

// ApplyAction performs one root-action Command, dispatching to the same engine
// method the corresponding player click calls. The action's player is the active
// player, as it was when the command was chosen. It returns the engine's own error
// for an illegal play (wrong house, no target, exhausted, ...) so a driver can
// surface it, and ErrNotRootAction for a Command that is not a root action.
func (g *Game) ApplyAction(cmd Command) error {
	p := g.State.ActivePlayer
	switch cmd.Kind {
	case CommandChooseHouse:
		return g.ChooseHouse(p, cmd.House)
	case CommandPlayCreature:
		_, err := g.PlayCreature(p, cmd.Hand, cmd.Left)
		return err
	case CommandPlayArtifact:
		_, err := g.PlayArtifact(p, cmd.Hand)
		return err
	case CommandPlayTactic:
		return g.PlayTactic(p, cmd.Hand)
	case CommandPlayUpgrade:
		_, err := g.PlayUpgrade(p, cmd.Hand)
		return err
	case CommandDiscardFromHand:
		return g.DiscardFromHand(p, cmd.Hand)
	case CommandReap:
		return g.Reap(p, cmd.Card)
	case CommandUnstun:
		return g.Unstun(p, cmd.Card)
	case CommandUseAction:
		return g.UseAction(p, cmd.Card)
	case CommandFight:
		return g.Fight(p, cmd.Card, cmd.Card2)
	case CommandEndTurn:
		// Ending the turn hands off to the other player, mirroring the web client's
		// end-turn root: a won game's StartTurn is a no-op, so this is safe after a
		// key forge decides the match.
		g.EndPlayPhase(p)
		g.StartTurn(1 - p)
		return nil
	default:
		return ErrNotRootAction
	}
}

// LegalActions returns every root-action Command the player may take right now, so
// a driver (the canonical turn loop, a UI, or the sim) can offer exactly the legal
// set instead of re-deriving it. It is the enumeration half of the root-action seam
// ApplyAction performs one command of; the contract that binds them —
// TestLegalActionsAreAllApplicable — is that every command LegalActions offers,
// ApplyAction on the same state accepts. Only the active player has root actions;
// any other player, or a finished game, has none.
//
// The set follows the phase: PhaseChooseHouse offers a house choice, PhasePlay
// offers plays, discards, uses, and ending the turn, and the engine-driven phases
// in between offer nothing (they auto-resolve or prompt through Requests, not root
// actions).
func (g *Game) LegalActions(player int) []Command {
	if g.State.Winner >= 0 || g.State.ActivePlayer != player {
		return nil
	}
	switch g.State.Phase {
	case PhaseChooseHouse:
		return g.legalHouseChoices(player)
	case PhasePlay:
		return g.legalPlayActions(player)
	default:
		return nil
	}
}

// legalHouseChoices lists the house the active player may name. An empty allowed
// set is not a dead end: the only legal choice is then HouseNone ("No House"),
// exactly what ChooseHouse accepts (ADR 0035).
func (g *Game) legalHouseChoices(player int) []Command {
	allowed := g.allowedHouses(player)
	if len(allowed) == 0 {
		return []Command{{Kind: CommandChooseHouse, House: HouseNone}}
	}
	out := make([]Command, len(allowed))
	for i, h := range allowed {
		out[i] = Command{
			Kind:  CommandChooseHouse,
			House: h,
		}
	}
	return out
}

// legalPlayActions lists every play-phase root action: the hand cards the player
// may play (a creature onto either flank, unless the battleline is empty and the
// flanks coincide) or discard, the creatures they may use (a stunned creature
// offers only Unstun, since any use just sheds the stun), the artifacts whose
// Action they may fire, and ending the turn, which is always open. Each offer is
// gated by the same Can* check ApplyAction's corresponding method performs, so the
// two never disagree.
func (g *Game) legalPlayActions(player int) []Command {
	var out []Command
	for i, id := range g.Hand(player) {
		if g.CanPlay(player, id) == nil {
			out = append(out, g.legalHandPlays(player, i, id)...)
		}
		if g.CanDiscard(player, id) == nil {
			out = append(out, Command{
				Kind: CommandDiscardFromHand,
				Hand: i,
			})
		}
	}
	for _, id := range g.Battleline(player) {
		out = append(out, g.legalCreatureUses(player, id)...)
	}
	for _, id := range g.Artifacts(player) {
		if g.CanUseArtifact(player, id) == nil {
			out = append(out, Command{
				Kind: CommandUseAction,
				Card: id,
			})
		}
	}
	return append(out, Command{Kind: CommandEndTurn})
}

// legalHandPlays lists the ways the hand card at index i may be played, dispatching
// on its printed type the way the sim's playCard does. A creature is offered on
// each flank, except on an empty battleline where both flanks place identically.
func (g *Game) legalHandPlays(player, i int, id LocalID) []Command {
	switch g.Def(id).Type {
	case Creature:
		out := []Command{{Kind: CommandPlayCreature, Hand: i}}
		if len(g.State.Battleline[player].slice()) > 0 {
			out = append(out, Command{
				Kind: CommandPlayCreature,
				Hand: i,
				Left: true,
			})
		}
		return out
	case Artifact:
		return []Command{{Kind: CommandPlayArtifact, Hand: i}}
	case Tactic:
		return []Command{{Kind: CommandPlayTactic, Hand: i}}
	default:
		// Upgrade — CanPlay has already ruled out any non-playable type in hand, so
		// the only type left is Upgrade. Making it the default rather than a fifth
		// case keeps every arm reachable (no dead branch for the coverage gate).
		return []Command{{Kind: CommandPlayUpgrade, Hand: i}}
	}
}

// legalCreatureUses lists the ways the battleline creature may be used: a stunned
// creature offers only Unstun (any use sheds the stun instead of acting), and an
// unstunned one offers reap, a Fight per legal target, and its Action, each gated
// exactly as the web client offers them (CanUseTo plus, for Fight, a legal target).
func (g *Game) legalCreatureUses(player int, id LocalID) []Command {
	if g.Stunned(id) {
		if g.canUnstun(player, id) == nil {
			return []Command{{Kind: CommandUnstun, Card: id}}
		}
		return nil
	}
	var out []Command
	if g.canUseTo(player, id, ReapUse) == nil {
		out = append(out, Command{
			Kind: CommandReap,
			Card: id,
		})
	}
	if g.canUseTo(player, id, FightUse) == nil {
		for _, def := range g.FightTargets(player, id) {
			out = append(out, Command{
				Kind:  CommandFight,
				Card:  id,
				Card2: def,
			})
		}
	}
	if g.hasTrigger(id, TriggerAction) && g.canUseTo(player, id, ActionUse) == nil {
		out = append(out, Command{
			Kind: CommandUseAction,
			Card: id,
		})
	}
	return out
}

// RunMatch is the canonical turn loop (ADR 0039, 0040): it deals the game, then
// drives it to a winner by asking the active player's ActionChooser for their next
// root action from the legal set and performing it, over and over. It is the
// standing Action a real match hands the session, so a client no longer hand-rolls
// its own turn driver — each yielded RequestAction is one click, and ApplyAction
// folds the answer back in.
//
// Because the loop owns setup, the mulligan prompts StartGame raises surface
// through the same choosers as every later decision, so they are recorded answers
// like any other. Between actions it hands off a turn an effect ended without
// pairing a StartTurn — the Omega keyword ends the play phase mid-action, and
// EndTurn (Book of leQ) jumps to the end-of-turn phase — so the loop resumes on the
// next player's turn rather than asking for an action in a phase that has none.
//
// The active player's chooser MUST be an ActionChooser; only an interactive driver
// (the Stepper's suspendChooser) installs one, and the sim, which drives legal play
// its own way, never calls this. The assertion is deliberate: a chooser that cannot
// answer a RequestAction has no business driving the loop.
func (g *Game) RunMatch(firstPlayer int) {
	g.StartGame(firstPlayer)
	for g.State.Winner < 0 {
		if g.State.Phase == PhaseEndOfTurn {
			g.StartTurn(1 - g.State.ActivePlayer)
			continue
		}
		player := g.State.ActivePlayer
		chooser, ok := g.chooserFor(player).(ActionChooser)
		if !ok {
			panic("game: active player chooser is not an ActionChooser")
		}
		// The chooser answers from the legal set (TestLegalActionsAreAllApplicable),
		// which ApplyAction always accepts, so its error is unreachable here.
		_ = g.ApplyAction(chooser.ChooseAction(g.LegalActions(player)))
	}
}
