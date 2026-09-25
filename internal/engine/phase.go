package engine

// Phase names one of the eight ordered parts of a turn (ADR 0012). KeyForge's
// rulebook calls these divisions "steps"; Vex calls them phases everywhere —
// engine identifiers, generated rulebook, rendered card text, and the game log —
// rather than carrying two words for one concept.
//
// The zero value is invalid (ADR 0010): a game that has not begun a turn is in no
// phase at all, which is a different thing from being in the first one.
type Phase uint8

// The phases of a turn, in the order they run.
const (
	// phaseUnset is the invalid zero value: a game between turns is in no phase.
	phaseUnset Phase = iota
	// PhaseStartOfTurn resolves "at the start of your turn" abilities, before the
	// active player forges.
	PhaseStartOfTurn
	// PhaseForge is the mandatory forge-a-key phase, after start-of-turn abilities
	// have had their chance to change what a key costs.
	PhaseForge
	// PhaseChooseHouse waits for the active player to name their active house.
	PhaseChooseHouse
	// PhaseArchives offers the active player's archived cards into their hand.
	PhaseArchives
	// PhasePlay is the open phase, shown to players as the "main" phase: the
	// active player plays, discards, and uses cards until they end their turn.
	PhasePlay
	// PhaseReady readies the active player's cards and refreshes creature armor.
	PhaseReady
	// PhaseDraw refills the active player's hand and sheds a chain if one blocked a
	// draw.
	PhaseDraw
	// PhaseEndOfTurn resolves "at the end of your turn" abilities, last of all, so
	// they see the board and hand the turn actually ends with (ADR 0013).
	PhaseEndOfTurn
)

// valid reports whether p names a real phase (not the unset zero value).
func (p Phase) valid() bool { return p != phaseUnset }

// phaseInfo is everything the engine knows about one phase, in one place: how it
// is named to a player, how the rulebook titles it, whether it blocks for a
// frontend decision, and what it does when it runs.
type phaseInfo struct {
	name string
	// rulebookStep is the phase's heading under the "Turn structure" rulebook term
	// — a numbered "N. Name" subtitle naming the phase as the rulebook presents it,
	// so the play phase reads "Main phase". This is the single source for those
	// subtitles: ruleterms_turn.go builds each term's Subtitle from it, and the
	// completeness check requires every phase to have a term keyed by it.
	rulebookStep string
	// waitsForInput marks a phase that blocks for a frontend decision rather than
	// running to completion the moment it is entered. Choosing a house and playing
	// are the two open phases; every other phase is engine-driven.
	waitsForInput bool
	// run carries out the phase. It is nil for an open phase, whose body is the
	// frontend's own play loop rather than anything the engine runs.
	run func(g *Game, player int)
}

// phases describes every phase. One table replaces the four parallel switches
// these fields used to live in, so a phase added to the enum cannot be named in
// one of them and forgotten in the others; TestPhaseTableIsTotal proves the
// table covers exactly Phases().
var phases = map[Phase]phaseInfo{
	PhaseStartOfTurn: {
		name:         "start of turn",
		rulebookStep: "1. Start of turn",
		run:          (*Game).startOfTurnPhase,
	},
	PhaseForge: {
		name:         "forge a key",
		rulebookStep: "2. Forge a key",
		run:          (*Game).forgePhase,
	},
	PhaseChooseHouse: {
		name:          "choose a house",
		rulebookStep:  "3. Choose a house",
		waitsForInput: true,
	},
	PhaseArchives: {
		name:         "archives",
		rulebookStep: "4. Archives",
		run:          (*Game).offerArchives,
	},
	PhasePlay: {
		name:          "main",
		rulebookStep:  "5. Main phase",
		waitsForInput: true,
	},
	PhaseReady: {
		name:         "ready",
		rulebookStep: "6. Ready",
		run:          (*Game).readyPhase,
	},
	PhaseDraw: {
		name:         "draw",
		rulebookStep: "7. Draw",
		run:          (*Game).drawStep,
	},
	PhaseEndOfTurn: {
		name:         "end of turn",
		rulebookStep: "8. End of turn",
		run:          (*Game).endOfTurnPhase,
	},
}

// String renders the phase as it is named to a player.
func (p Phase) String() string {
	if info, ok := phases[p]; ok {
		return info.name
	}
	return "no phase"
}

// waitsForInput reports whether the phase blocks for a frontend decision rather
// than running to completion the moment it is entered.
func (p Phase) waitsForInput() bool { return phases[p].waitsForInput }

// Phases lists every real phase in turn order (excluding the unset zero value).
// It is the canonical enumeration: the rulebook completeness check ranges over
// this, so a phase added above cannot go silently undescribed.
func Phases() []Phase {
	return []Phase{
		PhaseStartOfTurn,
		PhaseForge,
		PhaseChooseHouse,
		PhaseArchives,
		PhasePlay,
		PhaseReady,
		PhaseDraw,
		PhaseEndOfTurn,
	}
}

// rulebookStep returns the phase's numbered rulebook heading, or "" for a phase
// that is not a real one.
func (p Phase) rulebookStep() string { return phases[p].rulebookStep }
