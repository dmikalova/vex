package engine

// GAINING A TEXT BOX lets a creature take on another card's printed text box: its
// traits, its keywords, and its triggered abilities, resolved as if printed on the
// gaining creature. It does not copy the source's name, power, armor, type, house,
// or other stats — only the text box. Mimic Gel copies a chosen creature until
// Mimic Gel leaves play; Creed of Nurture lends a hand creature's text box to a
// creature in play for the remainder of the turn.
//
// A gained ability resolves with its Source set to the gaining creature, so any
// self-referential text ("destroy this creature", "{self} gains…") now refers to
// the creature that gained the text, matching the KeyForge rule that a copied or
// gained text box rebinds its own references. The source's text box is read from
// the immutable card catalog, so it stays available even after the source card
// leaves play — and it is the source's PRINTED text box, never a text box the
// source itself gained, so gaining never chains.
//
// Because the game state is a flat, pointerless value, a gain is stored as the
// source's LocalID+1 on the gaining creature's CardCore: TextBoxSourcePlus lasts
// until the creature leaves play (resetCore clears it) and TextBoxTurnSourcePlus
// lasts the turn (the ready phase clears it). Traits, keywords, and triggered
// abilities each fold the gained sources in at their read seam — HasTrait,
// hasKeyword, and triggeredBy. Constant abilities and numeric combat keywords
// (Assault N, Hazardous N) are not yet part of a gained text box; add them at
// those read seams when a card needs them.

// grantedTextBoxSources returns the cards whose printed text box creature id has
// gained — its permanent copy (Mimic Gel) and its remainder-of-turn loan (Creed
// of Nurture). Each is stored as a LocalID+1, so 0 means none.
func (g *Game) grantedTextBoxSources(id LocalID) []LocalID {
	core := &g.State.Cards[id]
	var out []LocalID
	if p := core.TextBoxSourcePlus; p != 0 {
		out = append(out, LocalID(p-1))
	}
	if p := core.TextBoxTurnSourcePlus; p != 0 {
		out = append(out, LocalID(p-1))
	}
	return out
}

// GrantTextBox gives creature recipient the printed text box of source, either
// until recipient leaves play or, when remainderOfTurn is set, for the remainder
// of the turn. It is the CreatureResolver port method both text-box effects share.
func (g *Game) GrantTextBox(recipient, source LocalID, remainderOfTurn bool) {
	c := g.stateOf(recipient)
	if c == nil {
		return
	}
	if remainderOfTurn {
		c.TextBoxTurnSourcePlus = uint8(source) + 1
	} else {
		c.TextBoxSourcePlus = uint8(source) + 1
	}
	g.record(CreatureGainedTextBox{Creature: recipient, Source: source})
}

// GainTextBox gives the creature its Target selects the printed text box of the
// creature its Source selects — that card's traits, keywords, and triggered
// abilities, as if printed on the recipient. RemainderOfTurn makes the gain last
// only the turn; otherwise it lasts until the recipient leaves play (Mimic Gel
// copies a chosen creature).
type GainTextBox struct {
	Target          Target
	Source          Target
	RemainderOfTurn bool
}

// validate requires both an explicit recipient and an explicit source.
func (e GainTextBox) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("GainTextBox")
	}
	if !e.Source.valid() {
		return errUnsetTarget("GainTextBox source")
	}
	return nil
}

// Text renders the effect, e.g. "{self} gains the text box of the chosen creature".
func (e GainTextBox) Text() string {
	return e.Target.Text() + " gains the text box of " + e.Source.Text()
}

// Resolve gives each recipient the source creature's text box. A source that
// selects nothing grants nothing.
func (e GainTextBox) Resolve(ctx *EffectContext) {
	sources := e.Source.Select(ctx)
	if len(sources) == 0 {
		return
	}
	source := sources[0]
	for _, recipient := range e.Target.Select(ctx) {
		ctx.Resolver.GrantTextBox(recipient, source, e.RemainderOfTurn)
	}
}

// LendTextBoxFromHand reveals a creature from the controller's hand and gives a
// chosen creature in play that revealed card's printed text box for the remainder
// of the turn — Creed of Nurture. Revealing turns the hand card public so the
// loaned text box can be trusted; the revealed card stays in hand.
type LendTextBoxFromHand struct{}

// Text renders the effect as its whole instruction, naming the two creatures the
// controller picks and the duration of the loan.
func (LendTextBoxFromHand) Text() string {
	return "reveal a creature from your hand and choose a creature in play - " +
		"for the remainder of the turn, the chosen creature gains the text box of " +
		"the revealed creature"
}

// Resolve reveals a chosen hand creature, then gives a chosen creature in play its
// text box for the turn. With no creature in hand there is nothing to reveal, so
// nothing happens.
func (LendTextBoxFromHand) Resolve(ctx *EffectContext) {
	var handCreatures []LocalID
	for _, id := range ctx.Resolver.Hand(ctx.Controller) {
		if ctx.Resolver.TypeOf(id) == Creature {
			handCreatures = append(handCreatures, id)
		}
	}
	if len(handCreatures) == 0 {
		return
	}
	source, ok := ctx.ChooseCard("Reveal a creature from your hand", handCreatures)
	if !ok {
		return
	}
	ctx.Resolver.Record(CardsRevealedToAll{Player: ctx.Controller, Cards: []LocalID{source}})
	recipients := (Target{Kind: TargetChosenCreature}).Select(ctx)
	if len(recipients) == 0 {
		return
	}
	ctx.Resolver.GrantTextBox(recipients[0], source, true)
}

// CreatureGainedTextBox narrates a creature gaining another card's text box.
type CreatureGainedTextBox struct {
	Creature LocalID
	Source   LocalID
}

// Text renders the creature and the card whose text box it gained, e.g.
// "Mimic Gel gains the text box of Troll".
func (e CreatureGainedTextBox) Text(n Namer) string {
	return n.Name(e.Creature) + " gains the text box of " + n.Name(e.Source)
}
