package engine

// This file holds REANIMATE-IN-PLACE: a Destroyed reaction that discards the top
// card of the controller's deck and, when that card is a creature, arms it to
// enter play in the source's own former battleline slot once the source has left
// play (Gebuk). The delay matters — the new creature lands only after the source
// has reached the discard pile — so the arming and the firing are split across a
// small flat registry, imitating the lasting-effect pattern (game_lasting.go): the
// effect records a ReanimateInPlace, and destroyTogether fires it as the batch
// clears.

// maxReanimations bounds how many delayed put-into-play records can be armed at
// once — one per source destroyed in a single batch. A fixed size keeps the state
// a flat value; extra arms are dropped silently.
const maxReanimations = 4

// ReanimateInPlace is one armed delayed put-into-play: after Source leaves play,
// Creature enters Controller's battleline at Position — the slot Source held. It
// is a plain comparable value so the flat GameState can hold a fixed array of
// them (ADR 0005).
type ReanimateInPlace struct {
	Source     LocalID
	Creature   LocalID
	Controller int8
	Position   uint8
}

// reanimator is the Game capability this effect needs beyond the common Resolver
// port: arming a delayed put-into-play. The concrete Resolver is always the
// *Game, so the effect upgrades the port in-package to reach it rather than
// widening the shared Resolver interface for a single mechanic.
type reanimator interface {
	armReanimateInPlace(r ReanimateInPlace)
}

// armReanimateInPlace records a delayed put-into-play, dropping it silently when
// the registry is full.
func (g *Game) armReanimateInPlace(r ReanimateInPlace) {
	if int(g.State.ReanimationsCount) >= maxReanimations {
		return
	}
	g.State.Reanimations[g.State.ReanimationsCount] = r
	g.State.ReanimationsCount++
}

// fireReanimations puts into play every armed creature whose source has left
// play, each at its source's former battleline slot, then drops those records. A
// record whose source is still in play is kept — its "after ... leaves play"
// window has not opened. destroyTogether calls it as a destruction batch clears,
// so a source's own reanimation lands only after the source has reached the
// discard pile.
func (g *Game) fireReanimations() {
	n := 0
	for i := 0; i < int(g.State.ReanimationsCount); i++ {
		r := g.State.Reanimations[i]
		if g.inPlay(r.Source) {
			g.State.Reanimations[n] = r
			n++
			continue
		}
		g.putIntoPlayAt(r.Creature, int(r.Controller), int(r.Position))
	}
	for i := n; i < int(g.State.ReanimationsCount); i++ {
		g.State.Reanimations[i] = ReanimateInPlace{}
	}
	g.State.ReanimationsCount = uint8(n)
}

// putIntoPlayAt puts a creature into play at a dictated battleline position,
// without prompting for a flank — the placement is the leaving source's own slot.
// The position is clamped to the current line length, since creatures destroyed
// alongside the source have already left and shortened the line. Like putIntoPlay,
// the creature enters exhausted with full armor.
func (g *Game) putIntoPlayAt(id LocalID, controller, pos int) {
	g.removeFromAnyZone(id)
	core := &g.State.Cards[id]
	core.Exhausted = true
	core.ArmorRemaining = int16(g.armor(id))
	if n := int(g.State.Battleline[controller].Count); pos > n {
		pos = n
	}
	g.State.Battleline[controller].insertAt(pos, id)
	g.record(CardPutIntoPlay{Player: controller, Card: id})
	g.emitCreatureEnters(id)
	g.settleDestroyed(controller)
}

// ReanimateTopOfDeckInPlace discards the top card of the controller's deck and,
// when that card is a creature, arms it to enter play in the source's own former
// battleline slot once the source has left play (Gebuk). The discarded card
// always reaches the discard pile; only a creature is armed to return.
type ReanimateTopOfDeckInPlace struct{}

// Text renders the effect. It is a Destroyed reaction, so it speaks from the
// controller's perspective ("your deck") and names the source with {self}.
func (ReanimateTopOfDeckInPlace) Text() string {
	return "discard the top card of your deck. If it is a creature, after " +
		SelfName + " leaves play, put that creature into play in " + SelfName +
		"'s position in the battleline"
}

// Resolve discards the top card of the controller's deck and, when it is a
// creature, arms it to enter play in the source's battleline slot after the
// source leaves play.
func (ReanimateTopOfDeckInPlace) Resolve(ctx *EffectContext) {
	id, ok := ctx.Resolver.DiscardTopOfDeck(ctx.Controller)
	if !ok || !ctx.Resolver.IsCreature(id) {
		return
	}
	pos := 0
	for i, c := range ctx.Resolver.Battleline(ctx.Controller) {
		if c == ctx.Source {
			pos = i
			break
		}
	}
	ctx.Resolver.(reanimator).armReanimateInPlace(ReanimateInPlace{
		Source:     ctx.Source,
		Creature:   id,
		Controller: int8(ctx.Controller),
		Position:   uint8(pos),
	})
}
