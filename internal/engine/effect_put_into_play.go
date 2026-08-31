package engine

// PutIntoPlay puts each targeted card into play without playing it. Putting a
// card into play is distinct from playing it: bonus icons and Play: abilities do
// not resolve (only "enters play" reactions do), which is what lets an effect put
// an opponent's card into play without making that player's play decisions.
// UnderYourControl puts the card under the ability controller's control (Overlord
// Greking reanimates a destroyed enemy "into play under your control"); otherwise
// it enters under its owner's control. Ownership never changes.
type PutIntoPlay struct {
	Target           Target
	UnderYourControl bool
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
	if e.UnderYourControl {
		return "put " + e.Target.Text() + " into play under your control"
	}
	return "put " + e.Target.Text() + " into play"
}

// Resolve puts each selected card into play, under the controller's control when
// UnderYourControl is set and under its owner's otherwise.
func (e PutIntoPlay) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		controller := ctx.Resolver.Owner(id)
		if e.UnderYourControl {
			controller = ctx.Controller
		}
		ctx.Resolver.PutIntoPlay(id, controller)
	}
}
