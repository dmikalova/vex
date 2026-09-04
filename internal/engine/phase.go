package engine

// Phase names one of the eight ordered parts of a turn (ADR 0012). KeyForge's
// rulebook calls these divisions "steps"; Vactrol calls them phases everywhere —
// engine identifiers, generated rulebook, rendered card text, and the game log —
// rather than carrying two words for one concept.
//
// The zero value is invalid (ADR 0010): a game that has not begun a turn is in no
// phase at all, which is a different thing from being in the first one.
type Phase uint8

// The phases of a turn, in the order they run.
const (
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

// String renders the phase as it is named to a player.
func (p Phase) String() string {
	switch p {
	case PhaseStartOfTurn:
		return "start of turn"
	case PhaseForge:
		return "forge a key"
	case PhaseChooseHouse:
		return "choose a house"
	case PhaseArchives:
		return "archives"
	case PhasePlay:
		return "main"
	case PhaseReady:
		return "ready"
	case PhaseDraw:
		return "draw"
	case PhaseEndOfTurn:
		return "end of turn"
	default:
		return "no phase"
	}
}

// waitsForInput reports whether the phase blocks for a frontend decision rather
// than running to completion the moment it is entered. Choosing a house and
// playing are the two open phases; every other phase is engine-driven.
func (p Phase) waitsForInput() bool {
	return p == PhaseChooseHouse || p == PhasePlay
}
