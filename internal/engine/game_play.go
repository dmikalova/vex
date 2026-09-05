package engine

import "slices"

// This file holds playing a card: creatures, artifacts, actions, and upgrades,
// plus discarding from hand and the shared checks (house match, taking a card of
// the right type, granting its Æmber pips).

// PlayCreature plays a creature from hand onto the battleline. flankLeft places it
// on the left flank; otherwise it goes to the right flank.
func (g *Game) PlayCreature(player, handIndex int, flankLeft bool) (LocalID, error) {
	if g.cannotPlayCreatures(player) {
		return 0, ErrCannotPlayCreature
	}
	id, err := g.validateHandPlay(player, handIndex, Creature)
	if err != nil {
		return 0, err
	}
	def := g.cat.def(id)
	result, err := g.playCardFromZone(
		player,
		id,
		func() { g.State.Hand[player].removeAt(handIndex) },
		playCardOptions{
			flankLeft:             flankLeft,
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

// chargeToll makes player pay every toll an opponent's in-play card imposes for
// action (Customs Office, Tentacus), moving that Æmber to the opponent. It is the
// single cost gate the play and use sites share: it changes nothing and returns
// ErrCannotPayToll when player cannot pay the full amount owed.
func (g *Game) chargeToll(player int, action TollAction) error {
	owed := g.tollOwed(player, action)
	if owed == 0 {
		return nil
	}
	if g.State.Aember[player] < owed {
		return ErrCannotPayToll
	}
	payee := 1 - player
	g.State.Aember[player] -= owed
	g.State.Aember[payee] += owed
	g.record(TollPaid{Player: player, Payee: payee, Amount: owed, Action: action})
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
// moving it to the discard pile. It performs no other effect.
func (g *Game) DiscardFromHand(player, handIndex int) error {
	if g.State.Winner >= 0 {
		return ErrGameOver
	}
	if g.State.ActivePlayer != player {
		return ErrNotActivePlayer
	}
	hand := &g.State.Hand[player]
	if handIndex < 0 || handIndex >= int(hand.Count) {
		return ErrCardNotInHand
	}
	id := hand.IDs[handIndex]
	if !g.inActiveHouse(g.cat.def(id)) {
		return ErrWrongHouse
	}
	hand.removeAt(handIndex)
	g.State.Discard[player].add(id)
	g.record(CardDiscarded{Player: player, Card: id})
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

// PlayFromUnder plays a specific card from under whatever host it sits under,
// bypassing the active-house gate the same way PlayFromHand/PlayFromDiscard do —
// Masterplan's and Jargogle's own "play the card under me." Where playFromPile
// removes the card from a deckList, a card under a host sits in the intrusive
// Under chain (game_under.go) instead, so this has its own small body rather
// than reusing playFromPile. It does nothing if the card is not currently placed
// under anything.
func (g *Game) PlayFromUnder(player int, id LocalID) {
	if _, ok := g.underHostOf(id); !ok {
		return
	}
	_, _ = g.playCardFromZone(player, id, func() { g.detachUnder(id) }, playCardOptions{})
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
	flankLeft             bool
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
	if !def.PlayRequirement.met(g.State.Aember[player]) {
		return 0, ErrPlayRequirement
	}
	switch def.Type {
	case Creature:
		if g.cannotPlayCreatures(player) {
			return 0, ErrCannotPlayCreature
		}
		g.recordCardPlayed(player, id, opts)
		remove()
		g.playCreatureCard(player, id, opts.flankLeft)
		return id, nil
	case Artifact:
		if err := g.chargeToll(player, TollPlayArtifact); err != nil {
			return 0, err
		}
		g.recordCardPlayed(player, id, opts)
		remove()
		g.playArtifactCard(player, id)
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

// playCreatureCard places a creature on a flank — or, for a Deploy creature,
// anywhere in the battleline its controller chooses — and fires the standard play
// sequence for a creature already removed from its previous zone.
func (g *Game) playCreatureCard(player int, id LocalID, flankLeft bool) {
	core := &g.State.Cards[id]
	core.Exhausted = true // enters play exhausted; readies during the end-of-turn ready step
	core.ArmorRemaining = int16(g.armor(id))
	pos, interior := g.deployPosition(player, id, flankLeft)
	g.State.Battleline[player].insertAt(pos, id)
	g.record(CardPlayedToBattleline{
		Player:    player,
		Card:      id,
		FlankLeft: pos == 0,
		Interior:  interior,
	})
	g.applyAemberBonus(id)
	g.triggerAbilities(id, TriggerAfterPlay, 0, false)
	g.emitCreatureEnters(id)
	g.emitCardPlayed(player, id)
	g.emitLasting(EventCreaturePlayed, player, id)
	g.emitLasting(EventCardEntersPlay, player, id)
	// The arrival may have walked into a power-reducing constant ability, or been one.
	g.settleDestroyed(player)
}

// putIntoPlay puts a card into play under controller's control without playing
// it: the card is removed from wherever it rests and enters directly — a creature
// onto the right flank of controller's battleline, an artifact into controller's
// artifact row. Because it is not played, its bonus icons and Play: abilities do
// not resolve; only "enters play" reactions fire. Ownership is unchanged, so the
// card still returns to its owner's zone when it later leaves play. Control is
// held permanently: the card is its own control source, and releaseControlHeldBy
// only ever reverts control anchored to a leaving upgrade, never to the card.
func (g *Game) putIntoPlay(id LocalID, controller int) {
	if g.inPlay(id) {
		return
	}
	g.removeFromAnyZone(id)
	core := &g.State.Cards[id]
	if controller != g.owner(id) {
		core.ControlPlus = uint8(controller + 1)
	}
	core.ControlSource = id
	switch g.cat.def(id).Type {
	case Creature:
		core.Exhausted = true
		core.ArmorRemaining = int16(g.armor(id))
		g.State.Battleline[controller].add(id)
		g.record(CardPutIntoPlay{Player: controller, Card: id})
		g.emitCreatureEnters(id)
	case Artifact:
		g.State.Artifacts[controller].add(id)
		g.record(CardPutIntoPlay{Player: controller, Card: id})
	}
	g.settleDestroyed(controller)
}

// playArtifactCard places an artifact and fires the standard play sequence for an
// artifact already removed from its previous zone.
func (g *Game) playArtifactCard(player int, id LocalID) {
	g.State.Cards[id].Exhausted = true // enters play exhausted; readies during the end-of-turn ready step
	g.State.Artifacts[player].add(id)
	g.record(ArtifactPlayed{Player: player, Card: id})
	g.applyAemberBonus(id)
	g.triggerAbilities(id, TriggerAfterPlay, 0, false)
	g.emitCardPlayed(player, id)
	g.emitLasting(EventCardEntersPlay, player, id)
	g.settleDestroyed(player)
}

// playActionCard resolves and discards an action already removed from its previous
// zone.
func (g *Game) playActionCard(player int, id LocalID) {
	g.record(ActionPlayed{Player: player, Card: id})
	g.applyAemberBonus(id)
	// The Play: ability resolves under the control of the player who played the
	// card, not the card's owner. They differ only when one player plays another's
	// card (Mimicry copies an action out of the opponent's discard pile).
	g.triggerAbilitiesAs(player, id, TriggerAfterPlay, 0, false)
	g.emitCardPlayed(player, id)
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
	g.State.Discard[owner].add(id)
}

// playUpgradeCard attaches an upgrade to host and fires its standard play sequence
// after the upgrade has been removed from its previous zone.
func (g *Game) playUpgradeCard(player int, id, host LocalID, def *CardDefinition) {
	g.applyAemberBonus(id)
	g.AttachUpgrade(host, id)
	g.record(UpgradeAttached{Player: player, Upgrade: id, Host: host})
	g.resolveUpgradePlay(host, id, def)
	g.settleDestroyed(player)
}

// DiscardCardFromHand moves a specific card from a player's hand to their discard
// zone, with no active-player or house checks (an effect may discard from either
// hand). It does nothing if the card is not in that hand.
func (g *Game) DiscardCardFromHand(owner int, id LocalID) {
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

// DiscardRandomFromHand discards one uniformly random card from a player's hand,
// doing nothing if the hand is empty.
func (g *Game) DiscardRandomFromHand(owner int) {
	hand := &g.State.Hand[owner]
	if hand.Count == 0 {
		return
	}
	g.DiscardCardFromHand(owner, hand.IDs[g.rng.Intn(int(hand.Count))])
}

// inActiveHouse reports whether a card of the given definition matches the
// active house for the purpose of PLAYING or discarding it from hand: true when
// no house has been chosen or the card's own house is the active house. Manual
// mode lifts the restriction. Versatile does not apply here — it only
// relaxes using a card already in play (see usableInActiveHouse).
func (g *Game) inActiveHouse(def *CardDefinition) bool {
	return g.manual ||
		g.State.ActiveHouse == HouseNone ||
		def.House == g.State.ActiveHouse
}

// mayPlayFromHand reports whether player may play a hand card now by active-house
// match or by a remaining continuous off-house play grant they control.
func (g *Game) mayPlayFromHand(player int, def *CardDefinition) bool {
	if g.inActiveHouse(def) {
		return true
	}
	if def.House == g.State.MayPlayHouse[player] {
		return true
	}
	return g.playPermissionRemaining(player, def.House) > 0
}

// usesPlayPermission reports whether playing def from hand would spend one of
// player's continuous off-house permissions instead of matching the active house.
func (g *Game) usesPlayPermission(player int, def *CardDefinition) bool {
	return !g.inActiveHouse(def) && g.playPermissionRemaining(player, def.House) > 0
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
		g.State.Aember[player] -= r.Aember
		g.record(AemberSpentToPlay{Player: player, Card: id, Amount: r.Aember})
	}
	if opts.consumePlayPermission {
		g.State.PlayPermissionsUsedThisTurn[player][g.cat.def(id).House]++
	}
	g.State.PlayedThisTurn[player].add(id)
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
		return ErrCannotPlayCreature
	}
	if g.barredFromPlaying(player, def.Type) {
		return ErrCannotPlayType
	}
	if g.cannotPlayCard(player) {
		return ErrCardPlayLimit
	}
	if g.barredByAlpha(player, def) {
		return ErrAlphaNotFirst
	}
	if !g.mayPlayFromHand(player, def) {
		return ErrWrongHouse
	}
	if !def.PlayRequirement.met(g.State.Aember[player]) {
		return ErrPlayRequirement
	}
	if def.Type == Artifact &&
		g.State.Aember[player] < g.tollOwed(player, TollPlayArtifact) {
		return ErrCannotPayToll
	}
	if def.Type == Upgrade &&
		len(
			g.State.Battleline[player].slice(),
		) == 0 && len(g.State.Battleline[1-player].slice()) == 0 {
		return ErrNoTarget
	}
	return nil
}

// applyAemberBonus grants a card's Æmber pips to its controller.
func (g *Game) applyAemberBonus(id LocalID) {
	def := g.cat.def(id)
	if def.AemberBonus > 0 {
		o := g.owner(id)
		if capturer, ok := g.gainAember(o, def.AemberBonus); ok {
			g.record(AemberBonusCaptured{
				Creature: capturer,
				Card:     id,
				Amount:   def.AemberBonus,
			})
			return
		}
		g.record(AemberBonusGained{Player: o, Card: id, Amount: def.AemberBonus})
	}
}
