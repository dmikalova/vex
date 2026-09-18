package engine

// StatMask is one stat fixed to a value for a duration, together with whether it
// is masked at all. The two travel as one value so a masked value can never be
// read without its Set flag, and so a call masking both power and armor passes
// two named stats rather than four positional arguments whose value and flag
// halves could be paired up wrongly.
//
// An unset mask carries no value: Value is meaningful only when Set is true.
type StatMask struct {
	Value int8
	Set   bool
}
