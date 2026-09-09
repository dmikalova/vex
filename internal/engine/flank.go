package engine

// flank is where a creature enters its controller's battleline. The zero value is
// invalid on purpose: a placement is either dictated by the card that puts the
// creature into play or chosen by its controller, and the engine never silently
// assumes one. deployPosition prompts the controller for a flank whenever a play
// leaves the flank unset and the battleline already holds a creature.
type flank int

const (
	// flankUnset is the invalid zero value: the controller is prompted for a flank
	// (or, for a Deploy creature being played, a full position) when a play does not
	// dictate where the creature enters.
	flankUnset flank = iota
	// flankLeftmost enters the creature on the left flank.
	flankLeftmost
	// flankRightmost enters the creature on the right flank.
	flankRightmost
)

// FlankLeftLabel and FlankRightLabel are the option labels the flank and position
// choosers offer for the two ends of the battleline.
const (
	FlankLeftLabel  = "Left flank"
	FlankRightLabel = "Right flank"
)

// FlankPromptPrefix begins the prompt asking a controller which flank a creature
// entering a non-empty battleline takes. It is exported so a driver that
// auto-answers this incidental placement prompt — the card-test harness — can
// recognise it without hard-coding the whole string, and without colliding with
// the move-a-creature-to-a-flank effects that share the flank option labels.
const FlankPromptPrefix = "Choose a flank for "

// flankFor turns the client-facing flankLeft bool of a hand play into an explicit
// flank: the player already chose the side, so the placement is dictated, never
// prompted.
func flankFor(flankLeft bool) flank {
	if flankLeft {
		return flankLeftmost
	}
	return flankRightmost
}
