package engine

// NeighborMode is how a Target reaches the battleline neighbors of the cards it
// selects. It is a set expansion rather than a per-card test — it runs after the
// selection is made and changes which cards the set holds — so it sits beside a
// Filter on the Target rather than inside it.
//
// The three modes are one field because they are alternatives: a target either
// leaves the selection alone, widens it to the neighbors as well, or replaces it
// with them. Held as two independent flags they could say both at once, which no
// card means and the renderer could not print.
type NeighborMode int

const (
	// NeighborsNone is the zero value: the target selects what it selects, and the
	// battleline around it is not read.
	NeighborsNone NeighborMode = iota
	// NeighborsIncluded keeps each selected creature and adds its neighbors, so an
	// effect applies to the creature and each of its neighbors (Tremor).
	NeighborsIncluded
	// NeighborsOnly replaces the selection with the neighbors of what it selects,
	// dropping the selected creature itself — Lord Golgotha damages each neighbor of
	// the creature it fights, but not that creature.
	NeighborsOnly
)

// NeighborModes returns every real mode, so the census enumerates them rather
// than keeping a list a new mode would fall out of.
func NeighborModes() []NeighborMode {
	return []NeighborMode{NeighborsIncluded, NeighborsOnly}
}

// expands reports whether the mode reads the battleline at all.
func (m NeighborMode) expands() bool { return m != NeighborsNone }

// keepsSelected reports whether the selected card stays in the set alongside its
// neighbors.
func (m NeighborMode) keepsSelected() bool { return m != NeighborsOnly }

// decorate wraps a rendered noun phrase in the mode's wording: NeighborsIncluded
// reads "<phrase> and each of its neighbors", NeighborsOnly reads "each neighbor
// of <phrase>".
func (m NeighborMode) decorate(phrase string) string {
	switch m {
	case NeighborsIncluded:
		return phrase + " and each of its neighbors"
	case NeighborsOnly:
		return "each neighbor of " + phrase
	}
	return phrase
}
