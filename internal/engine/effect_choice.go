package engine

import "strings"

// A "choose one" ability offers its controller a set of alternative effects and
// resolves only the one they pick; the options not chosen do nothing.
type ChooseOne struct {
	Options []Effect
}

// Text renders the choice as a "choose one:" header followed by one bulleted,
// capitalized line per option, e.g.:
//
//	choose one:
//	- Destroy each Dis creature
//	- Gain 1 Æmber
func (e ChooseOne) Text() string {
	var b strings.Builder
	b.WriteString("choose one:")
	for _, o := range e.Options {
		b.WriteString("\n- " + capitalizeFirst(o.Text()))
	}
	return b.String()
}

// Resolve asks the controller which option to take, then resolves that one.
func (e ChooseOne) Resolve(ctx *EffectContext) {
	options := make([]string, len(e.Options))
	for i, o := range e.Options {
		options[i] = o.Text()
	}
	idx := ctx.ChooseOption("Choose one", options)
	if idx < 0 || idx >= len(e.Options) {
		return
	}
	e.Options[idx].Resolve(ctx)
}

// validate surfaces the first configuration error among the options.
func (e ChooseOne) validate() error {
	for _, o := range e.Options {
		if err := validateEffect(o); err != nil {
			return err
		}
	}
	return nil
}

// ChooseHouseThen asks the controller to choose a house, records it on the effect
// context, and resolves Then — which typically acts on creatures of that house
// through Target.House(chosenHouse). It models "Choose a house. <do something to that
// house>."
type ChooseHouseThen struct {
	Then Effect
}

// Text renders the effect, e.g. "choose a house. Stun each creature of the
// chosen house".
func (e ChooseHouseThen) Text() string {
	return leadInSentence("choose a house", e.Then.Text())
}

// hedgedText renders the form a May wraps, where the choice is optional and so
// its consequence must say so.
func (e ChooseHouseThen) hedgedText() string {
	return hedgedLeadInSentence("choose a house", e.Then.Text())
}

func (ChooseHouseThen) endsSentence() bool { return true }

// Resolve asks for a house, stores it on the context, then resolves Then.
func (e ChooseHouseThen) Resolve(ctx *EffectContext) {
	options := houseNames[1:] // every house except HouseNone
	idx := ctx.ChooseOption("Choose a house", options)
	if idx < 0 || idx >= len(options) {
		return
	}
	ctx.ChosenHouse = House(idx + 1) // house values start at Brobnar = 1
	e.Then.Resolve(ctx)
}

// validate descends into the wrapped effect.
func (e ChooseHouseThen) validate() error { return validateEffect(e.Then) }

// OpponentNamesHouse has the opponent (the non-controlling player) name a house,
// storing it in ctx.ChosenHouse for a following effect to read — an ItIsNotOfNamedHouse
// condition tests the card in context against it. Keyforgery: the forging opponent
// names a house before the guard reveals a card. It renders as its own sentence, so
// it is a standalone effect rather than a ChooseHouseThen wrapper.
type OpponentNamesHouse struct{}

// Text renders the effect. The opponent is "they" — the antecedent is the trigger
// that names them (Keyforgery's "when your opponent would forge a key").
func (OpponentNamesHouse) Text() string { return "they name a house" }

// Resolve asks the opponent to name a house and stores it on the context.
func (OpponentNamesHouse) Resolve(ctx *EffectContext) {
	options := houseNames[1:] // every house except HouseNone
	idx := ctx.Resolver.ChooseOption(ctx.Opponent(), ctx.Source, "Name a house", options)
	if idx < 0 || idx >= len(options) {
		return
	}
	ctx.ChosenHouse = House(idx + 1) // house values start at Brobnar = 1
}
