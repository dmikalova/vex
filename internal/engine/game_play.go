package engine

import "slices"

// This file holds playing a card: creatures, artifacts, actions, and upgrades,
// plus discarding from hand and the shared checks (house match, taking a card of
// the right type, granting its Æmber pips).

// PlayCreature plays a creature from hand onto the battleline. flankLeft places it
// on the left flank; otherwise it goes to the right flank.
func (g *Game) PlayCreature(player, handIndex int, flankLeft bool) (LocalID, error) {
	id, err := g.validateHandPlay(player, handIndex, Creature)
	if err != nil {
		return 0, err
	}
	def := g.cat.def(id)
	// A creature played as an upgrade is not playing a creature, so a creature ban
	// bars it only when it cannot become an upgrade — no host to attach to.
	if g.cannotPlayCreatures(player) && (!def.PlayableAsUpgrade || !g.hasUpgradeHost()) {
		return 0, ErrCannotPlayCreature
	}
	result, err := g.playCardFromZone(
		player,
		id,
		func() { g.State.Hand[player].removeAt(handIndex) },
		playCardOptions{
			flank:                 flankFor(flankLeft),
			consumePlayPermission: g.usesPlayPermission(player, def),
		},
	)
	if err == nil {
		g.endStepIfOmega(player, def)
	}
	return result, err
}

// PlayArtifact plays an artifact from hand into the artifact row. If the opponent
// controls a card that tolls artifact plays (Customs Office), the player must pay
// that toll first and cannot play the artifact if they cannot.
func (g *Game) PlayArtifact(player, handIndex int) (LocalID, error) {
	id, err := g.validateHandPlay(player, handIndex, Artifact)
	if err != nil {
		return 0, err
	}
	def := g.cat.def(id)
	result, err := g.playCardFromZone(
		player,
		id,
		func() { g.State.Hand[player].removeAt(handIndex) },
		playCardOptions{
			consumePlayPermission: g.usesPlayPermission(player, def),
		},
	)
	if err == nil {
		g.endStepIfOmega(player, def)
	}
	return result, err
}

// chargeToll makes player give every toll an opponent's in-play card imposes for
// action (Customs Office, Tentacus), moving that Æmber to the opponent. It checks
// the full amount owed up front and returns ErrCannotPayToll (changing nothing)
// when player cannot pay it, then gives each toll card's share under that card's
// frame so the log names the card. It is the single cost gate the play and use
// sites share.
func (g *Game) chargeToll(player int, action TollAction) error {
	owed := g.tollOwed(player, action)
	if owed == 0 {
		return nil
	}
	if g.State.Aember[player] < owed {
		return ErrCannotPayToll
	}
	payee := 1 - player
	for _, id := range g.allInPlay(payee) {
		t := g.cat.def(id).Restricts.Toll
		if t.Amount <= 0 || t.Action != action {
			continue
		}
		g.SetAember(player, g.State.Aember[player]-t.Amount)
		g.SetAember(payee, g.State.Aember[payee]+t.Amount)
		closeFrame := g.openFrame(Frame{Actor: payee, Source: id, HasSource: true})
		g.record(AemberGiven{Giver: player, Receiver: payee, Amount: t.Amount, Reason: action})
		closeFrame()
	}
	return nil
}

// tollOwed totals the Æmber player must hand the opponent to take action.
func (g *Game) tollOwed(player int, action TollAction) int {
	owed := 0
	for _, id := range g.allInPlay(1 - player) {
		if t := g.cat.def(id).Restricts.Toll; t.Amount > 0 && t.Action == action {
			owed += t.Amount
		}
	}
	return owed
}

// PlayAction plays an action card: its Æmber bonus and "Play:" abilities resolve,
// then it goes to the discard pile.
func (g *Game) PlayAction(player, handIndex int) error {
	id, err := g.validateHandPlay(player, handIndex, Tactic)
	if err != nil {
		return err
	}
	def := g.cat.def(id)
	_, err = g.playCardFromZone(
		player,
		id,
		func() { g.State.Hand[player].removeAt(handIndex) },
		playCardOptions{
			consumePlayPermission: g.usesPlayPermission(player, def),
		},
	)
	if err == nil {
		g.endStepIfOmega(player, def)
	}
	return err
}

// DiscardFromHand discards a card of the active house from a player's hand,
// moving it to the discard pile and firing any "after you discard a card from
// your hand" reactions (Baron Mengevin captures Æmber when you discard a Sanctum
// card).
func (g *Game) DiscardFromHand(player, handIndex int) error {
	if g.State.Winner >= 0 {
		return ErrGameOver
	}
	if g.State.ActivePlayer != player {
		return ErrNotActivePlayer
	}
	if g.barredByFirstTurn(player) {
		return ErrFirstTurnOneCard
	}
	hand := &g.State.Hand[player]
	if handIndex < 0 || handIndex >= int(hand.Count) {
		return ErrCardNotInHand
	}
	id := hand.IDs[handIndex]
	if !g.inActiveHouse(g.cat.def(id)) {
		return ErrWrongHouse
	}
	g.DiscardCardFromHand(player, id)
	return nil
}

// CanDiscard reports whether the player can discard the given hand card right now:
// nil if discardable, otherwise the reason. It mirrors the checks DiscardFromHand
// enforces, the way CanPlay mirrors the PlayX methods, so a caller can list the
// legal discards without duplicating the rule.
func (g *Game) CanDiscard(player int, id LocalID) error {
	if g.State.Winner >= 0 {
		return ErrGameOver
	}
	if g.State.ActivePlayer != player {
		return ErrNotActivePlayer
	}
	if !slices.Contains(g.State.Hand[player].slice(), id) {
		return ErrCardNotInHand
	}
	if g.barredByFirstTurn(player) {
		return ErrFirstTurnOneCard
	}
	if !g.inActiveHouse(g.cat.def(id)) {
		return ErrWrongHouse
	}
	return nil
}

// PlayUpgrade plays an upgrade from hand, attaching it to a chosen creature
// (friendly or enemy). It returns the host creature it was attached to.
func (g *Game) PlayUpgrade(player, handIndex int) (LocalID, error) {
	if g.State.Winner >= 0 {
		return 0, ErrGameOver
	}
	if g.State.ActivePlayer != player {
		return 0, ErrNotActivePlayer
	}
	if g.cannotPlayCard(player) {
		return 0, ErrCardPlayLimit
	}
	if g.barredByFirstTurn(player) {
		return 0, ErrFirstTurnOneCard
	}
	hand := &g.State.Hand[player]
	if handIndex < 0 || handIndex >= int(hand.Count) {
		return 0, ErrCardNotInHand
	}
	id := hand.IDs[handIndex]
	def := g.cat.def(id)
	if def.Type != Upgrade {
		return 0, ErrWrongType
	}
	if g.barredByAlpha(player, def) {
		return 0, ErrAlphaNotFirst
	}
	if !g.mayPlayFromHand(player, def) {
		return 0, ErrWrongHouse
	}
	host, err := g.playCardFromZone(
		player,
		id,
		func() { hand.removeAt(handIndex) },
		playCardOptions{
			consumePlayPermission: g.usesPlayPermission(player, def),
		},
	)
	if err == nil {
		g.endStepIfOmega(player, def)
	}
	return host, err
}

// TopOfDeck returns the top card of a player's deck without moving it, reporting
// whether the deck holds a card.
func (g *Game) TopOfDeck(player int) (LocalID, bool) {
	deck := &g.State.Deck[player]
	if deck.Count == 0 {
		return 0, false
	}
	return deck.IDs[0], true
}

// PlayFromDeck plays a specific card from a player's deck, removing it from the
// deck as it is played (Chaos Portal plays the card it revealed). It does nothing
// when the card is not in that deck.
func (g *Game) PlayFromDeck(player int, id LocalID) {
	g.playFromPile(player, id, &g.State.Deck[player])
}

// PlayFromHand plays a specific card out of a player's hand, bypassing the
// active-house gate — an effect that plays a card has already decided it may be
// played (Phase Shift's off-house card). It does nothing when the card is not in
// that hand.
func (g *Game) PlayFromHand(player int, id LocalID) {
	g.playFromPile(player, id, &g.State.Hand[player])
}

// PlayFromDiscard plays a specific card out of a player's discard pile, the way
// Sacrificial Altar brings a creature back. It does nothing when the card is not
// in that discard pile.
func (g *Game) PlayFromDiscard(player int, id LocalID) {
	g.playFromPile(player, id, &g.State.Discard[player])
}

// PlayFromOpponentDiscard plays a specific card out of the opponent's discard
// pile as the given player's own play — Mimicry copies an action out of the
// other player's discard, so the play counts against the active player's own
// card-play limit (Ember Imp can block it) yet the card is owned by, and returns
// to, the opponent. It does nothing when the card is not in that discard pile.
func (g *Game) PlayFromOpponentDiscard(player int, id LocalID) {
	g.playFromPile(player, id, &g.State.Discard[1-player])
}

// PlayFromOpponentHand plays a specific card out of the opponent's hand as the
// given player's own play — Lateral Shift plays a card out of the other player's
// hand "as if it were yours." The play counts against the active player's own
// card-play limit and, for a creature or artifact, the player takes control of it
// (via playForeign) while its owner stays the opponent. It does nothing when the
// card is not in that hand.
func (g *Game) PlayFromOpponentHand(player int, id LocalID) {
	hand := &g.State.Hand[1-player]
	i := hand.indexOf(id)
	if i < 0 {
		return
	}
	g.playForeign(player, id, func() { hand.removeAt(i) })
}

// PlayFromUnder plays a specific card from under whatever host it sits under,
// bypassing the active-house gate the same way PlayFromHand/PlayFromDiscard do —
// Masterplan's and Jargogle's own "play the card under me." A card under a host
// sits in the intrusive Under chain (game_under.go), so it is detached from there
// rather than removed from a deckList. It routes through playForeign so that a
// buried card the player does not own (one that reached their hand and was buried
// there) is played under the player's control, like any other foreign play. It
// does nothing if the card is not currently placed under anything.
func (g *Game) PlayFromUnder(player int, id LocalID) {
	if _, ok := g.underHostOf(id); !ok {
		return
	}
	g.playForeign(player, id, func() { g.detachUnder(id) })
}

// PlayFromArchives plays a specific card out of a player's archives, bypassing the
// active-house gate the same way PlayFromHand does — Project Z.Y.X. plays an
// archived card "as if it were in your hand and in the active house." Archives is
// a wideList rather than a deckList, so it does not reuse playFromPile. It does
// nothing when the card is not in that player's archives.
func (g *Game) PlayFromArchives(player int, id LocalID) {
	arc := &g.State.Archives[player]
	if arc.indexOf(id) < 0 {
		return
	}
	_, _ = g.playCardFromZone(player, id, func() { arc.remove(id) }, playCardOptions{})
}

// PlayFromOpponent plays a card from a zone of player's opponent as player's own
// play (Murkens): the top card of their deck, or a uniformly random card from
// their facedown archives (drawn at random because the archives are hidden). An
// empty source zone does nothing.
func (g *Game) PlayFromOpponent(player int, from Zone) {
	switch from {
	case Deck:
		deck := &g.State.Deck[1-player]
		if deck.Count == 0 {
			return
		}
		g.playForeign(player, deck.IDs[0], func() { deck.removeAt(0) })
	case Archives:
		arc := &g.State.Archives[1-player]
		if arc.Count == 0 {
			return
		}
		id := arc.IDs[g.State.PRNG.Intn(int(arc.Count))]
		g.playForeign(player, id, func() { arc.remove(id) })
	}
}

// playForeign plays a card that player may not own as player's own play. When the
// card is owned by the other player, player takes control of it — held as a
// permanent control naming the card itself as source, so it lasts until the card
// leaves play and the card still returns to its owner's zone. Control is set
// before the play so the card's controller reads correctly while its Play:
// ability resolves; a rejected play (card-play limit, unmet play requirement)
// reverts it.
func (g *Game) playForeign(player int, id LocalID, remove func()) {
	// Only a card that stays in play under an independent controller — a creature
	// or an artifact — needs a control entry so it reads as player's after the
	// play. A Tactic resolves and leaves play at once, so taking control of it
	// would leak in-play state onto a card no longer in play.
	t := g.cat.def(id).Type
	foreign := g.owner(id) != player && (t == Creature || t == Artifact)
	if foreign {
		g.State.Cards[id].ControlPlus = uint8(player + 1)
		g.pushControl(id, player, id)
	}
	if _, err := g.playCardFromZone(player, id, remove, playCardOptions{}); err != nil && foreign {
		g.State.Cards[id].ControlPlus = 0
		g.clearControls(id)
	}
}

// playFromPile plays a card out of one of a player's face-down piles. Where the
// card comes from is the only thing that differs between playing from deck, hand,
// and discard pile — every gate, the log, and the card's Play: ability are the
// shared playCardFromZone — so the zone is a parameter rather than three bodies.
func (g *Game) playFromPile(player int, id LocalID, pile *deckList) {
	i := pile.indexOf(id)
	if i < 0 {
		return
	}
	_, _ = g.playCardFromZone(player, id, func() { pile.removeAt(i) }, playCardOptions{})
}

type playCardOptions struct {
	// flank is where the creature enters the battleline. It is flankUnset for an
	// effect-play that does not dictate a side, so the controller is prompted; a
	// hand play sets it from the player's own choice (flankFor).
	flank                 flank
	consumePlayPermission bool
}

// playCardFromZone resolves a card play after the source zone has identified the
// card. It performs the shared play gates, removes the card only once it can be
// played, records the play, then dispatches to the card-type-specific placement.
func (g *Game) playCardFromZone(
	player int,
	id LocalID,
	remove func(),
	opts playCardOptions,
) (LocalID, error) {
	if g.State.Winner >= 0 {
		return 0, ErrGameOver
	}
	if g.State.ActivePlayer != player {
		return 0, ErrNotActivePlayer
	}
	if g.cannotPlayCard(player) {
		return 0, ErrCardPlayLimit
	}
	def := g.cat.def(id)
	if !g.meetsPlayRequirement(player, def) {
		return 0, ErrPlayRequirement
	}
	switch def.Type {
	case Creature:
		if def.GiganticRole != GiganticNone {
			return g.playGigantic(player, id, remove, opts)
		}
		banned := g.cannotPlayCreatures(player)
		if def.PlayableAsUpgrade && g.hasUpgradeHost() {
			// Playing a creature as an upgrade is not playing a creature, so it is
			// allowed even while creatures are barred — a ban then forces upgrade mode.
			asUpgrade := banned
			if !banned {
				asUpgrade = g.chooseOption(player, g.Name(id),
					"Play "+SelfName+" as a creature or an upgrade?",
					[]string{"Creature", "Upgrade"}) == 1
			}
			if asUpgrade {
				candidates := append(g.battlelineCopy(player), g.battlelineCopy(1-player)...)
				host, ok := g.pickCreature(
					player,
					g.Name(id),
					"Choose a creature to attach "+SelfName+" to",
					candidates,
				)
				if !ok {
					return 0, ErrNoTarget
				}
				g.recordCardPlayed(player, id, opts)
				remove()
				g.playUpgradeCard(player, id, host, def)
				return host, nil
			}
		}
		if banned {
			return 0, ErrCannotPlayCreature
		}
		g.recordCardPlayed(player, id, opts)
		remove()
		g.playCreatureCard(player, id, opts.flank)
		g.applyTreachery(player, id)
		return id, nil
	case Artifact:
		if err := g.chargeToll(player, TollPlayArtifact); err != nil {
			return 0, err
		}
		g.recordCardPlayed(player, id, opts)
		remove()
		g.playArtifactCard(player, id)
		g.applyTreachery(player, id)
		return id, nil
	case Tactic:
		g.recordCardPlayed(player, id, opts)
		remove()
		g.playActionCard(player, id)
		return id, nil
	case Upgrade:
		candidates := append(g.battlelineCopy(player), g.battlelineCopy(1-player)...)
		host, ok := g.pickCreature(
			player,
			g.Name(id),
			"Choose a creature to attach "+SelfName+" to",
			candidates,
		)
		if !ok {
			return 0, ErrNoTarget
		}
		g.recordCardPlayed(player, id, opts)
		remove()
		g.playUpgradeCard(player, id, host, def)
		return host, nil
	default:
		return 0, ErrWrongType
	}
}

// applyTreachery hands a just-played Treachery card to its player's opponent,
// permanently — the keyword means the card enters play under your opponent's
// control. The card is placed and its Play abilities and Æmber bonus resolve for
// the player first (they still take the bonus); only then does control pass. The
// card itself is the control source, so the control lasts until it leaves play.
func (g *Game) applyTreachery(player int, id LocalID) {
	if g.hasKeyword(id, Treachery) {
		g.takeControl(id, 1-player, id)
		// The active player chose to play it, so they choose which flank of the
		// opponent's battleline it enters (a control change; the active player always
		// places). The handoff re-forms neighbors on both battlelines, so a creature
		// that just lost a flank bonus can drop to or below its damage; settle that
		// power shift here, the boundary of the handoff (ADR 0029).
		g.placeGainedOnFlank(id, 1-player)
		g.settleDestroyed(player)
	}
}

// hasUpgradeHost reports whether any creature is in play on either battleline to
// attach an upgrade to — the host an Upgrade, or a creature played as an upgrade,
// needs.
func (g *Game) hasUpgradeHost() bool {
	return g.State.Battleline[0].Count > 0 || g.State.Battleline[1].Count > 0
}

// playCreatureCard places a creature on a flank — or, for a Deploy creature,
// anywhere in the battleline its controller chooses — and fires the standard play
// sequence for a creature already removed from its previous zone.
func (g *Game) playCreatureCard(player int, id LocalID, fl flank) {
	core := &g.State.Cards[id]
	core.Exhausted = true // enters play exhausted; readies during the end-of-turn ready step
	if g.entersPlayReady(g.controller(id), Creature, g.cat.def(id).House) {
		core.Exhausted = false
	}
	core.ArmorRemaining = int16(g.armor(id))
	pos, interior := g.deployPosition(player, id, fl, true)
	g.State.Battleline[player].insertAt(pos, id)
	g.record(CardPlayedToBattleline{
		Player:    player,
		Card:      id,
		FlankLeft: pos == 0,
		Interior:  interior,
	})
	// The arrival can push a neighbor off a flank, dropping it below the power its
	// flank bonus was keeping it at. That is a state-based death that settles at
	// once — before bonus icons and before the after-play window — so the played
	// creature's own ability never sees the doomed neighbor still in play.
	g.settleDestroyed(player)
	g.resolveBonusIcons(player, id)
	// Playing a creature opens one window: its own "Play:" and "enters play"
	// abilities, every bystander's "after a creature enters play / is played /
	// is played adjacent", every "after you play a card" reaction, and the
	// duration reactions on playing a card, a creature, and a card entering play
	// all trigger at once, so the active player orders the whole set (ADR 0013).
	pending := g.playCreatureReactions(player, id)
	pending = append(pending, g.lastingReactions(EventCardPlayed, player, id)...)
	pending = append(pending, g.lastingReactions(EventCreaturePlayed, player, id)...)
	pending = append(pending, g.lastingReactions(EventCardEntersPlay, player, id)...)
	g.resolveWindow(g.orderTriggered(player, pending))
	// The arrival may have walked into a power-reducing constant ability, or been one.
	g.settleDestroyed(player)
}

// putIntoPlay puts a card into play under controller's control without playing
// it: the card is removed from wherever it rests and enters directly — a creature
// onto a flank of controller's battleline that controller is prompted to choose,
// an artifact into controller's artifact row. Because it is not played, its bonus
// icons and Play: abilities do not resolve; only "enters play" reactions fire.
// Ownership is unchanged, so the card still returns to its owner's zone when it
// later leaves play. When it enters under a player other than its owner, that
// control is held permanently: the card names itself as the source of a
// control-stack entry, which only clears when the card itself leaves play (never
// reverted by a leaving source).
func (g *Game) putIntoPlay(id LocalID, controller int) {
	if g.inPlay(id) {
		return
	}
	g.removeFromAnyZone(id)
	core := &g.State.Cards[id]
	if controller != g.owner(id) {
		core.ControlPlus = uint8(controller + 1)
		g.pushControl(id, controller, id)
	}
	switch g.cat.def(id).Type {
	case Creature:
		core.Exhausted = true
		core.ArmorRemaining = int16(g.armor(id))
		pos, _ := g.deployPosition(controller, id, flankUnset, false)
		g.State.Battleline[controller].insertAt(pos, id)
		g.record(CardPutIntoPlay{Player: controller, Card: id})
		g.emitCreatureEnters(id)
	case Artifact:
		g.State.Artifacts[controller].add(id)
		g.record(CardPutIntoPlay{Player: controller, Card: id})
	}
}

// playArtifactCard places an artifact and fires the standard play sequence for an
// artifact already removed from its previous zone.
func (g *Game) playArtifactCard(player int, id LocalID) {
	g.State.Cards[id].Exhausted = true // enters play exhausted; readies during the end-of-turn ready step
	if g.entersPlayReady(g.controller(id), Artifact, g.cat.def(id).House) {
		g.State.Cards[id].Exhausted = false
	}
	g.State.Artifacts[player].add(id)
	g.record(ArtifactPlayed{Player: player, Card: id})
	g.resolveBonusIcons(player, id)
	// The artifact's own "Play:", every "after you play a card" reaction, and the
	// duration reactions on playing a card and a card entering play trigger at once,
	// so the active player orders the set (ADR 0013).
	pending := g.afterPlayReactions(player, id)
	pending = append(pending, g.lastingReactions(EventCardPlayed, player, id)...)
	pending = append(pending, g.lastingReactions(EventCardEntersPlay, player, id)...)
	g.resolveWindow(g.orderTriggered(player, pending))
	g.settleDestroyed(player)
}

// playActionCard resolves and discards an action already removed from its previous
// zone.
func (g *Game) playActionCard(player int, id LocalID) {
	g.record(ActionPlayed{Player: player, Card: id})
	g.resolveBonusIcons(player, id)
	// A reaction to a Tactic being played resolves before the Tactic's own effect
	// (Encounter Suit wards its host before the Tactic can reach it).
	g.emitActionPlayedBeforeResolve(player, id)
	// The Tactic's own "Play:", every "after you play a card" reaction, and the
	// duration reactions on playing a card trigger at once, so the active player
	// orders the set (ADR 0013). The "Play:" resolves under the control of the player
	// who played the card, not its owner; they differ only when one player plays
	// another's card (Mimicry copies an action out of the opponent's discard pile),
	// which afterPlayReactions carries as the entry actor.
	pending := g.afterPlayReactions(player, id)
	pending = append(pending, g.lastingReactions(EventCardPlayed, player, id)...)
	g.resolveWindow(g.orderTriggered(player, pending))
	// A played action goes to the top of its owner's discard pile — unless its own
	// "Play:" ability purged it (Library Access), in which case it is set aside out
	// of the game instead. Owner and player differ only when one player plays
	// another's card (Mimicry).
	owner := g.owner(id)
	if g.State.PurgePlayedActionSet && g.State.PurgePlayedAction == id {
		g.State.PurgePlayedAction = 0
		g.State.PurgePlayedActionSet = false
		g.State.Purge[owner].add(id)
		g.record(CardPurged{Card: id})
		return
	}
	if g.State.ArchivePlayedActionSet && g.State.ArchivePlayedAction == id {
		g.State.ArchivePlayedAction = 0
		g.State.ArchivePlayedActionSet = false
		g.PutIntoArchives(id)
		return
	}
	g.State.Discard[owner].add(id)
}

// playUpgradeCard attaches an upgrade to host and fires its standard play sequence
// after the upgrade has been removed from its previous zone.
func (g *Game) playUpgradeCard(player int, id, host LocalID, def *CardDefinition) {
	g.AttachUpgrade(host, id)
	g.record(UpgradeAttached{Player: player, Upgrade: id, Host: host})
	// Bonus icons resolve after the upgrade enters play (KeyForge), so a Damage icon
	// that kills the host finds the upgrade already attached and it sheds cleanly
	// rather than attaching to a destroyed host.
	g.resolveBonusIcons(player, id)
	g.resolveUpgradePlay(host, id, def)
	g.emitUpgradeEntered(id)
	g.settleDestroyed(player)
}

// DiscardCardFromHand moves a specific card from a player's hand to their discard
// zone, with no active-player or house checks (an effect may discard from either
// hand). It does nothing if the card is not in that hand. The discard is subjected
// to the card whose ability forced it through the record's frame, or to the player
// when there is none (a player discarding as their turn action).
func (g *Game) DiscardCardFromHand(owner int, id LocalID) {
	g.discardFromHand(owner, id)
}

// discardFromHand carries out a hand-to-discard move and its after-discard
// reactions.
func (g *Game) discardFromHand(owner int, id LocalID) {
	hand := &g.State.Hand[owner]
	i := hand.indexOf(id)
	if i < 0 {
		return
	}
	hand.removeAt(i)
	g.State.Discard[owner].add(id)
	g.State.DiscardedThisTurn[owner].add(id)
	g.record(CardDiscarded{Player: owner, Card: id})
	for _, watcher := range g.allInPlay(owner) {
		g.triggerAbilities(watcher, TriggerAfterDiscardFromHand, id, true)
	}
}

// inActiveHouse reports whether a card of the given definition matches the
// active house for the purpose of PLAYING or discarding it from hand: true when
// the card's own house is the active house. A player at No House — locked out of
// every house and resolved to no active house for their play phase (ADR 0035) —
// matches nothing, so they play only the cards an out-of-house allowance permits.
// Before the choice resolves (HouseNone outside the play phase: setup and manual
// board building) nothing is house-gated yet. Manual mode lifts the restriction.
// Versatile does not apply here — it only relaxes using a card already in play
// (see usableInActiveHouse).
func (g *Game) inActiveHouse(def *CardDefinition) bool {
	return g.manual ||
		(g.State.ActiveHouse == HouseNone && g.State.Phase != PhasePlay) ||
		def.House == g.State.ActiveHouse
}

// mayPlayFromHand reports whether player may play a hand card now by active-house
// match or by a remaining off-house play grant they control — a continuous
// permission (Witch of the Wilds, Captain Val Jericho) or a this-turn permit
// (Com. Officer Kirby, CXO Taber, United Action).
func (g *Game) mayPlayFromHand(player int, def *CardDefinition) bool {
	if g.inActiveHouse(def) {
		return true
	}
	if g.freesTypeUnlimited(player, def) {
		return true
	}
	if def.House == g.State.MayPlayHouse[player] {
		return true
	}
	if g.playPermissionRemaining(player, def.House) > 0 {
		return true
	}
	if g.offHousePlayPermit(player, def) >= 0 {
		return true
	}
	return g.nonActivePlayRemaining(player) > 0
}

// usesPlayPermission reports whether playing def from hand would spend one of
// player's off-house play grants instead of matching the active house or an
// unlimited Ambassador-style grant.
func (g *Game) usesPlayPermission(player int, def *CardDefinition) bool {
	if g.inActiveHouse(def) || def.House == g.State.MayPlayHouse[player] {
		return false
	}
	return g.playPermissionRemaining(player, def.House) > 0 ||
		g.offHousePlayPermit(player, def) >= 0 ||
		g.nonActivePlayRemaining(player) > 0
}

// consumeOffHousePlay spends the off-house grant that made an off-house hand play
// legal, trying a continuous house-specific permission first (Witch of the Wilds),
// then a this-turn permit (Kirby, Taber, United Action), then a continuous
// "any non-active house" permission (Captain Val Jericho).
func (g *Game) consumeOffHousePlay(player int, def *CardDefinition) {
	if g.playPermissionRemaining(player, def.House) > 0 {
		if int(def.House) < NumHouses {
			g.State.PlayPermissionsUsedThisTurn[player][def.House]++
		}
		return
	}
	if i := g.offHousePlayPermit(player, def); i >= 0 {
		g.consumeOffHousePermit(player, i)
		return
	}
	g.State.NonActivePlaysUsedThisTurn[player]++
}

// freesTypeUnlimited reports whether player controls an in-play card whose
// PlayPermission is a house-agnostic, unlimited waiver for def's card type — Matter
// Maker frees any number of upgrades. It grants without limit, so nothing consumes
// it and no per-turn counter tracks it.
func (g *Game) freesTypeUnlimited(player int, def *CardDefinition) bool {
	for _, id := range g.allInPlay(player) {
		p := g.cat.def(id).PlayPermission
		if p.Types != 0 && p.Types.has(def.Type) {
			return true
		}
	}
	return false
}

// playPermissionRemaining returns how many unspent permissions player controls
// for playing a card of house while house is not their active house.
func (g *Game) playPermissionRemaining(player int, house House) int {
	limit := 0
	if g.State.ActiveHouse != house {
		for _, id := range g.allInPlay(player) {
			if p := g.cat.def(id).PlayPermission; p.granted() && p.House == house {
				limit += p.count()
			}
		}
	}
	used := 0
	if int(house) < NumHouses {
		used = int(g.State.PlayPermissionsUsedThisTurn[player][house])
	}
	if used >= limit {
		return 0
	}
	return limit - used
}

// recordCardPlayed logs the play for the turn, charges any Æmber the card's play
// requirement spends, and includes an off-house grant if the hand play was legal
// only because of that continuous permission.
func (g *Game) recordCardPlayed(player int, id LocalID, opts playCardOptions) {
	def := g.cat.def(id)
	if r := def.PlayRequirement; r.Spend && r.required() {
		// Draw the pool first, then fall back to creatures whose Æmber counts as pool
		// (Senator Bracchus) for any shortfall.
		fromPool := min(r.Aember, g.State.Aember[player])
		g.SetAember(player, g.State.Aember[player]-fromPool)
		if rest := r.Aember - fromPool; rest > 0 {
			g.drawFromSpendAsPool(player, rest)
		}
		g.record(AemberSpentToPlay{Player: player, Card: id, Amount: r.Aember})
	}
	if opts.consumePlayPermission {
		g.consumeOffHousePlay(player, def)
	}
	g.State.PlayedThisTurn[player].add(id)
}

// barredByFirstTurn reports whether the first-turn rule blocks player from
// playing or discarding a card from hand now: it is armed only for the first
// player's first turn, and blocks once they have taken one volitional play or
// discard. A card that a played card lets them play (Wild Wormhole, Phase Shift)
// goes through playCardFromZone directly, not these gates, so it is never barred.
func (g *Game) barredByFirstTurn(player int) bool {
	if g.manual {
		return false
	}
	return g.State.FirstTurnPlayLimit[player] &&
		g.State.PlayedThisTurn[player].Count+g.State.DiscardedThisTurn[player].Count >= 1
}

// validateHandPlay runs the read-only checks that a hand card may be played —
// game not over, active player, no card-play limit reached, index in hand, right
// type, and play permission — and returns the card without mutating the hand. It
// is the shared validation step before the source zone actually removes the card.
func (g *Game) validateHandPlay(player, handIndex int, want CardType) (LocalID, error) {
	if g.State.Winner >= 0 {
		return 0, ErrGameOver
	}
	if g.State.ActivePlayer != player {
		return 0, ErrNotActivePlayer
	}
	if g.cannotPlayCard(player) {
		return 0, ErrCardPlayLimit
	}
	if g.barredByFirstTurn(player) {
		return 0, ErrFirstTurnOneCard
	}
	hand := &g.State.Hand[player]
	if handIndex < 0 || handIndex >= int(hand.Count) {
		return 0, ErrCardNotInHand
	}
	id := hand.IDs[handIndex]
	def := g.cat.def(id)
	if def.Type != want {
		return 0, ErrWrongType
	}
	if g.barredFromPlaying(player, want) {
		return 0, ErrCannotPlayType
	}
	if g.barredByAlpha(player, def) {
		return 0, ErrAlphaNotFirst
	}
	if !g.mayPlayFromHand(player, def) {
		return 0, ErrWrongHouse
	}
	return id, nil
}

// CanPlay reports whether the player can play the given hand card right now: nil
// if playable, otherwise the reason (wrong house, card-play limit, creatures
// barred, or an upgrade with no host). It mirrors the checks the PlayX methods
// enforce, so a UI can dim unplayable cards and explain why before the click.
func (g *Game) CanPlay(player int, id LocalID) error {
	if g.State.Winner >= 0 {
		return ErrGameOver
	}
	if g.State.ActivePlayer != player {
		return ErrNotActivePlayer
	}
	def := g.cat.def(id)
	if def.Type == Creature && g.cannotPlayCreatures(player) {
		if !def.PlayableAsUpgrade || !g.hasUpgradeHost() {
			return ErrCannotPlayCreature
		}
	}
	if def.GiganticRole != GiganticNone {
		if _, _, ok := g.recruitGiganticPartner(
			id, []*deckList{&g.State.Hand[player]},
		); !ok {
			return ErrGiganticNoPartner
		}
	}
	if g.barredFromPlaying(player, def.Type) {
		return ErrCannotPlayType
	}
	if g.cannotPlayCard(player) {
		return ErrCardPlayLimit
	}
	if g.barredByFirstTurn(player) {
		return ErrFirstTurnOneCard
	}
	if g.barredByAlpha(player, def) {
		return ErrAlphaNotFirst
	}
	if !g.mayPlayFromHand(player, def) {
		return ErrWrongHouse
	}
	if !g.meetsPlayRequirement(player, def) {
		return ErrPlayRequirement
	}
	if def.Type == Artifact &&
		g.State.Aember[player] < g.tollOwed(player, TollPlayArtifact) {
		return ErrCannotPayToll
	}
	if def.Type == Upgrade &&
		!g.hasUpgradeHost() {
		return ErrNoTarget
	}
	return nil
}
