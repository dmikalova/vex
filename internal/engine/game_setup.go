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
