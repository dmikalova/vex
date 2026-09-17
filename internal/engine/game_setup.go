package engine

// This file holds registering cards into the catalog and placing them into a
// player's zones — the setup and test helpers used to build a game state.

// Register adds a definition to the catalog for an owner and returns its id.
func (g *Game) Register(def CardDefinition, owner int) LocalID {
	d := def
	return g.cat.add(&d, owner)
}

// AddToHand registers a card and places it in a player's hand.
func (g *Game) AddToHand(def CardDefinition, owner int) LocalID {
	id := g.Register(def, owner)
	g.State.Hand[owner].add(id)
	return id
}

// AddToDeck registers a card and places it on the bottom of a player's deck.
func (g *Game) AddToDeck(def CardDefinition, owner int) LocalID {
	id := g.Register(def, owner)
	g.State.Deck[owner].add(id)
	return id
}

// AddToBattleline registers a creature and places it on a player's battleline.
func (g *Game) AddToBattleline(def CardDefinition, owner int) LocalID {
	id := g.Register(def, owner)
	g.State.Cards[id].ArmorRemaining = int16(def.Armor)
	g.State.Battleline[owner].add(id)
	return id
}

// AddArtifact registers an artifact and places it in a player's artifact row.
func (g *Game) AddArtifact(def CardDefinition, owner int) LocalID {
	id := g.Register(def, owner)
	g.State.Artifacts[owner].add(id)
	return id
}

// AddToDiscard registers a card and places it in a player's discard pile.
func (g *Game) AddToDiscard(def CardDefinition, owner int) LocalID {
	id := g.Register(def, owner)
	g.State.Discard[owner].add(id)
	return id
}

// AddToArchives registers a card and places it in a player's archives.
func (g *Game) AddToArchives(def CardDefinition, owner int) LocalID {
	id := g.Register(def, owner)
	g.State.Archives[owner].add(id)
	return id
}

// FirstPlayerBonusCards is how many more cards the first player draws in their
// opening hand than the second player, who draws a full HandSize.
const FirstPlayerBonusCards = 1

// StartGame runs KeyForge setup and begins the first turn. It shuffles both
// decks, deals each player their opening hand (the first player draws one more
// than the second), offers each player one mulligan starting with the first,
// then starts the first player's turn and arms the first-turn rule. It expects
// both players' whole decks already registered in their deck zones; match
// populates them and calls this to begin play.
func (g *Game) StartGame(firstPlayer int) {
	g.record(GameStarted{FirstPlayer: firstPlayer})
	order := [2]int{firstPlayer, 1 - firstPlayer}
	for _, player := range order {
		g.Shuffle(player)
		g.record(DeckShuffled{Player: player})
	}
	for _, player := range order {
		base := HandSize
		if player == firstPlayer {
			base += FirstPlayerBonusCards
		}
		g.dealOpeningHand(player, base)
	}
	for _, player := range order {
		g.offerMulligan(player)
	}
	g.StartTurn(firstPlayer)
	// Armed after StartTurn, which clears the flag at the start of every turn, so it
	// lands on this opening turn only and is gone by the first player's next one.
	g.State.FirstTurnPlayLimit[firstPlayer] = true
}

// dealOpeningHand draws player up to base cards for their opening hand, cut by
// their chains the same way the draw phase cuts a refill (one fewer for every 6
// chains), shedding a single chain when that cut actually kept them from a card.
func (g *Game) dealOpeningHand(player, base int) {
	chains := g.State.Chains[player]
	target := base - (chains+5)/6
	g.drawTo(player, target)
	hand := int(g.State.Hand[player].Count)
	g.record(CardsDrawn{Player: player, Count: hand, Hand: hand})
	if chains > 0 && hand < base && g.canDraw(player) {
		g.State.Chains[player]--
		g.record(ChainShed{Player: player, Remaining: g.State.Chains[player]})
	}
}

// offerMulligan gives player their one setup mulligan: on choosing it they
// shuffle their opening hand back into the deck and draw one fewer card than they
// held. A mulligan sheds no further chain — the opening draw already shed one, so
// a chained player is down the same single chain whether or not they mulligan.
func (g *Game) offerMulligan(player int) {
	// Make player the active one while they decide, so a client that renders the
	// active player's hand shows this player their own opening hand to judge.
	g.State.ActivePlayer = player
	had := int(g.State.Hand[player].Count)
	if g.chooseOption(player, "", "Mulligan your opening hand?",
		[]string{"Keep", "Mulligan"}) != 1 {
		return
	}
	g.shuffleZonesIntoDeck(player, []Zone{Hand})
	g.record(DeckShuffled{Player: player})
	g.drawTo(player, had-1)
	g.record(Mulliganed{Player: player, Hand: int(g.State.Hand[player].Count)})
}
