package engine

// Control names whose control a card enters under when it is put into play: its
// owner's (the default) or the resolving player's. It renders its own text so
// PutIntoPlay does not branch on a bool.
type Control uint8

const (
	// ControlOwner puts the card under its owner's control (the default).
	ControlOwner Control = iota
	// ControlYours puts the card under the resolving player's control (Overlord
	// Greking reanimates a destroyed enemy "under your control").
	ControlYours
)

// suffix renders the "under your control" clause, empty for owner control.
func (c Control) suffix() string {
	if c == ControlYours {
		return " under your control"
	}
	return ""
}

// PutIntoPlay puts each targeted card into play without playing it. Putting a
// card into play is distinct from playing it: bonus icons and Play: abilities do
// not resolve (only "enters play" reactions do), which is what lets an effect put
// an opponent's card into play without making that player's play decisions.
// Control puts the card under the resolving player's control (Overlord Greking
// reanimates a destroyed enemy "into play under your control") or, by default,
// under its owner's. Ownership never changes.
type PutIntoPlay struct {
	Target  Target
	Control Control
}

// validate requires an explicit target.
func (e PutIntoPlay) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("PutIntoPlay")
	}
	return nil
}

// Text renders the effect, e.g. "put it into play under your control".
func (e PutIntoPlay) Text() string {
	return "put " + e.Target.Text() + " into play" + e.Control.suffix()
}

// Resolve puts each selected card into play, under the resolving player's control
// when Control is ControlYours and under its owner's otherwise.
func (e PutIntoPlay) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		controller := ctx.Resolver.Owner(id)
		if e.Control == ControlYours {
			controller = ctx.Controller
		}
		ctx.Resolver.PutIntoPlay(id, controller)
	}
}
