package engine

// Duration says how long a timed effect lasts. Timed effects (see Restrict)
// take a Duration, so one effect can serve several windows as more are added —
// the same way a Target says which cards an effect reaches.
type Duration int

const (
	// durationUnset is the invalid zero value: a timed effect must name its
	// duration rather than leave it unset.
	durationUnset Duration = iota
	// RemainderOfPlayerTurn lasts through the rest of the current turn, then lifts
	// at end of turn (Brain Stem Antenna's host belongs to Mars for the remainder of
	// the turn).
	RemainderOfPlayerTurn
	// OpponentNextTurn is dormant for the rest of this turn and bites only during
	// the affected player's next turn — typically the opponent's — lifting when that
	// turn ends (Fogbank: the opponent cannot fight during their next turn). It
	// waits for that player's next turn however many turns intervene, so an extra
	// turn taken by someone else never consumes it. Contrast StartOfPlayerNextTurn,
	// which is live now rather than dormant.
	OpponentNextTurn
	// StartOfPlayerNextTurn is live the moment it is established and stays live
	// through the opponent's turn, lifting at the start of the caster's next turn
	// (their ready phase) — so a defensive effect survives the opponent's turn (Into
	// the Night, Sow Salt, Diplomacy). It lifts at the same boundary as
	// OpponentNextTurn but, unlike it, is in force from the moment it resolves.
	StartOfPlayerNextTurn
	// EndOfPlayerNextTurn lasts from now through the end of the affected player's
	// own next turn, then lifts. "Next turn" is that player's next turn — not the
	// next turn any player takes — so for the controller the window covers the rest
	// of this turn, the opponent's intervening turn, and the whole of the
	// controller's next turn, lifting only when that turn ends. Unlike
	// OpponentNextTurn — which waits for that next turn before it bites — this
	// window is live the moment it is established and stays live across the
	// intervening turns, so it suits a reduction a player wants in force immediately
	// (We Can ALL Win's key-cost drop).
	EndOfPlayerNextTurn
	// UntilThisLeavesPlay lasts until the card whose effect established it leaves
	// play (Collar of Subordination's control change), rather than expiring at a
	// turn boundary. The leave-play teardown honors it.
	UntilThisLeavesPlay
	// Forever never lifts: the effect it establishes lasts for the rest of the
	// game (Sneklifter's control of a seized artifact). It registers no teardown.
	//
	// Design rule — "the latest ability wins": when two effects change the same
	// thing on the same card (control, house), the most recently applied one is in
	// force. Control is a single last-write-wins field, so a later take-control
	// simply overrides a Forever one. If a *timed* override is ever layered over a
	// Forever effect, its expiry must fall back to the Forever effect rather than to
	// the owner's default — the Forever effect still governs once the timed one
	// lifts. (Today artifact control is only ever Forever, so no such timed override
	// exists yet; this rule is the invariant to preserve when one is added.)
	Forever
	// durationCount is the exclusive upper bound Durations ranges to. It is not a
	// real duration; keep it last.
	durationCount
)

// valid reports whether d names a real duration (not the unset zero value).
func (d Duration) valid() bool { return d != durationUnset }

// Durations lists every real duration in declaration order, so a caller can
// enumerate the timing windows a timed effect can take (the /style gallery shows
// one card per duration). The unset zero value is excluded.
func Durations() []Duration {
	all := make([]Duration, 0, int(durationCount)-1)
	for d := durationUnset + 1; d < durationCount; d++ {
		all = append(all, d)
	}
	return all
}

// String names the timing window in the Rules voice — a short canonical label,
// not the clause a card prints. There is no single printed phrase for a duration:
// the same window renders differently per effect and flips between prefix and
// suffix (durationClause and the timed effects render those), so this names the
// window itself. The unset zero and the count sentinel render empty.
func (d Duration) String() string {
	switch d {
	case RemainderOfPlayerTurn:
		return "Rest of this turn"
	case OpponentNextTurn:
		return "Opponent's next turn"
	case StartOfPlayerNextTurn:
		return "Until your next turn"
	case EndOfPlayerNextTurn:
		return "Through your next turn"
	case UntilThisLeavesPlay:
		return "Until this leaves play"
	case Forever:
		return "Forever"
	default:
		return ""
	}
}
