package engine

// StatusIcon names the status a selection badge previews on a creature a player
// is about to choose. It is a display hint only — the client draws the matching
// icon over the candidate — so it carries no rules weight and no state depends on
// it.
type StatusIcon uint8

const (
	// NoStatusIcon is the zero badge: no icon, which clears any badge showing.
	NoStatusIcon StatusIcon = iota
	// DamageIcon marks a creature a pick is about to damage (Festering Touch).
	DamageIcon
	// WardIcon marks a creature a pick is about to ward (Imperium's "ward N").
	WardIcon
)

// SelectionBadge is the hint an effect shows for the creature a player is about
// to choose: which status it will get (Icon) and, when it matters, by how much
// (Amount — 3 for "deal 3 damage to each"; 0 when the status has no number, as a
// ward does). The zero SelectionBadge (Icon NoStatusIcon) clears any badge.
type SelectionBadge struct {
	Icon   StatusIcon
	Amount int
}
