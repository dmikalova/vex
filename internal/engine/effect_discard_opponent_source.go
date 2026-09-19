package engine

import "fmt"

// DiscardFromOpponent discards one card from a source the controller picks among
// Sources — Fidgit offers the opponent's facedown archives (a uniformly random
// card) and the top of their deck — and binds the discarded card as ctx.It so a
// following effect can react to it (Fidgit plays it when it is a Tactic). Sources
// reads each Zone the way PlayFromOpponent does: Archives is a random facedown
// card, Deck is the top card. One source discards from it with no prompt; several
// offer the choice. An empty chosen source discards nothing and binds nothing.
type DiscardFromOpponent struct {
	Sources []Zone
}

// validate requires at least one source, each a zone this effect can discard an
// opponent's card from — their archives or their deck top.
func (e DiscardFromOpponent) validate() error {
	if len(e.Sources) == 0 {
		return fmt.Errorf("DiscardFromOpponent: Sources must not be empty")
	}
	for _, z := range e.Sources {
		if z != Archives && z != Deck {
			return fmt.Errorf("DiscardFromOpponent: unsupported source %d", z)
		}
	}
	return nil
}

// Text renders the source choice as one clause, e.g. "discard a random card from
// your opponent's archives or the top card of their deck". The first source names
// "your opponent's"; a later one says "their", so the joined clause does not repeat
// it.
func (e DiscardFromOpponent) Text() string {
	clause := "discard "
	for i, z := range e.Sources {
		if i > 0 {
			clause += " or "
		}
		clause += discardSourcePhrase(z, i == 0)
	}
	return clause
}

// Resolve discards from the chosen source and binds the discarded card in context.
// A single source discards from it directly; several are offered as a choice.
func (e DiscardFromOpponent) Resolve(ctx *EffectContext) {
	opp := ctx.Opponent()
	source := e.Sources[0]
	if len(e.Sources) > 1 {
		labels := make([]string, len(e.Sources))
		for i, z := range e.Sources {
			labels[i] = discardSourceLabel(z)
		}
		source = e.Sources[ctx.ChooseOption("Which card to discard?", labels)]
	}
	id, ok := discardFromOpponentSource(ctx, opp, source)
	if !ok {
		return
	}
	ctx.It, ctx.HasIt = id, true
}

// discardFromOpponentSource discards from one opponent zone and reports the card:
// the archives give a uniformly random card, the deck its top card.
func discardFromOpponentSource(ctx *EffectContext, opp int, source Zone) (LocalID, bool) {
	if source == Archives {
		return discardRandomArchivesCard(ctx, opp)
	}
	return ctx.Resolver.DiscardTopOfDeck(opp)
}

// discardSourcePhrase renders a source as printed card text names it. first picks
// the fuller "your opponent's" over the "their" back-reference a later clause uses.
func discardSourcePhrase(source Zone, first bool) string {
	whose := "their "
	if first {
		whose = "your opponent's "
	}
	if source == Archives {
		return "a random card from " + whose + "archives"
	}
	return "the top card of " + whose + "deck"
}

// discardSourceLabel renders a source as a prompt option.
func discardSourceLabel(source Zone) string {
	if source == Archives {
		return "your opponent's archives"
	}
	return "the top card of their deck"
}

// discardRandomArchivesCard picks a random card from player's archives, discards
// it, and reports it. Empty archives discard nothing.
func discardRandomArchivesCard(ctx *EffectContext, player int) (LocalID, bool) {
	id, ok := ctx.Resolver.ChooseRandom(ctx.Resolver.Archives(player))
	if !ok {
		return 0, false
	}
	ctx.Resolver.DiscardCardFromArchives(player, id)
	return id, true
}
