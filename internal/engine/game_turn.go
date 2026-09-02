package engine

import (
	"errors"
	"slices"
)

// Engine action errors.
var (
	ErrNotActivePlayer       = errors.New("not the active player")
	ErrCardNotInHand         = errors.New("card index is not in hand")
	ErrWrongType             = errors.New("card is the wrong type for this action")
	ErrWrongHouse            = errors.New("card's house is not the active house")
	ErrCardExhausted         = errors.New("card is exhausted")
	ErrNoTarget              = errors.New("no legal target")
	ErrGameOver              = errors.New("the game is already over")
	ErrCannotFight           = errors.New("cannot use creatures to fight this turn")
	ErrCannotPlayCreature    = errors.New("cannot play creatures")
	ErrCannotPlayType        = errors.New("cannot play cards of this type this turn")
	ErrCardPlayLimit         = errors.New("card-play limit reached this turn")
	ErrCannotPayToll         = errors.New("cannot pay the toll for this action")
	ErrPlayRequirement       = errors.New("not enough Æmber to play this card")
	ErrMustChooseForcedHouse = errors.New("must choose the forced active house this turn")
	ErrCannotUse             = errors.New("card's use condition is not met")
)

// This file holds the turn lifecycle — begin (with the forge step), choose the
// active house, and end (ready cards, refresh armor, draw back up) — together
// with the key-forging step that decides the game.

// BeginTurn starts a player's turn: it becomes the active player, the active
// house is cleared, and the forge step runs (keys are forged if affordable).
func (g *Game) BeginTurn(player int) {
	if g.State.Winner >= 0 {
		return
	}
	g.State.ActivePlayer = player
	g.State.ActiveHouse = HouseNone
	g.State.Turn++
	g.State.PlayedThisTurn[player].reset()
	g.State.DiscardedThisTurn[player].reset()
	g.State.PlayPermissionsUsedThisTurn[player] = [NumHouses]uint8{}
	for p := 0; p < 2; p++ {
		for _, id := range g.State.Battleline[p].slice() {
			g.State.Cards[id].TimesUsedThisTurn = 0
			g.State.Cards[id].ElusiveUsedThisTurn = false
		}
	}
	// Each bar armed on a previous turn becomes active for this player now, taking
	// the card that imposed it along so a reminder can name the reason.
	g.State.CannotFight[player] = g.State.CannotFightNext[player]
	g.State.CannotFightNext[player] = Bar[bool]{}
	g.State.CannotPlayTypeThis[player] = g.State.CannotPlayTypeNext[player]
	g.State.CannotPlayTypeNext[player] = Bar[CardType]{}
	g.State.CannotUse[player] = g.State.CannotUseNext[player]
	g.State.CannotUseNext[player] = Bar[bool]{}
	g.State.ForcedHouse[player] = g.State.ForcedHouseNext[player]
	g.State.ForcedHouseNext[player] = Bar[House]{}
	g.State.SkipForge[player] = g.State.SkipForgeNext[player]
	g.State.SkipForgeNext[player] = Bar[bool]{}
	g.logf("--- %s begins turn %d ---", g.names[player], g.State.Turn)
	if g.State.SkipForge[player].Value {
		g.logf("%s skips their forge a key step", g.names[player])
	} else {
		g.forgeKey(player)
	}
	g.assertInvariants()
}

// Choose a house: pick one of your deck's three houses to be your active house
// for this turn. For the rest of the turn you may play from hand and use only
// cards of that house, except cards that ignore the restriction such as Versatile
// ones.
//
//rulebook:turn Turn structure / 2. Choose a house

// ChooseHouse sets ActiveHouse, then offers to draw the player's archived cards
// into hand now that a house is locked in.
func (g *Game) ChooseHouse(player int, house House) error {
	if g.State.ActivePlayer != player {
		return ErrNotActivePlayer
	}
	// A forced house (Control the Weak) only binds when the player actually has it;
	// if they cannot choose it, cannot overrides must and any house is allowed.
	if fh := g.State.ForcedHouse[player].Value; fh != HouseNone &&
		house != fh &&
		g.playerHasHouse(player, fh) {
		return ErrMustChooseForcedHouse
	}
	g.State.ActiveHouse = house
	g.logf("%s chooses house %s", g.names[player], house)
	g.offerArchives(player)
	for _, id := range g.allInPlay(player) {
		g.triggerAbilities(id, TriggerAfterChooseHouse, 0, false)
	}
	return nil
}

// playerHasHouse reports whether house is one the player may choose — a house in
// their declared deck houses. When deck houses are unknown (unset), every house
// is treated as available.
func (g *Game) playerHasHouse(player int, house House) bool {
	if len(g.houses[player]) == 0 {
		return true
	}
	for _, h := range g.houses[player] {
		if h == house {
			return true
		}
	}
	return false
}

// Ready and draw: at the end of your turn every card you control readies (turns
// back upright, ready to act next turn) and you draw back up to a full hand of
// six cards. Creatures and artifacts that entered play exhausted this turn ready
// here too.
//
//rulebook:turn Turn structure / 3. Ready and draw

// EndTurn clears Exhausted on the player's in-play cards, refreshes creature
// armor for the new turn, and draws up to HandSize.
func (g *Game) EndTurn(player int) {
	for _, id := range g.allInPlay(player) {
		g.triggerAbilities(id, TriggerEndOfTurn, 0, false)
	}
	for _, id := range g.allInPlay(player) {
		core := &g.State.Cards[id]
		core.Exhausted = false
		core.TempHouse = HouseNone
		if g.cat.def(id).Type == Creature {
			core.ArmorRemaining = int16(g.armor(id))
		}
	}
	// "Cannot be dealt damage" lasts only the turn, so clear it on every creature,
	// including any enemy one an effect protected (Protectrix).
	for _, id := range append(g.allInPlay(player), g.allInPlay(1-player)...) {
		g.State.Cards[id].DamageImmune = false
	}
	g.State.CannotFight[player] = Bar[bool]{}
	g.State.CannotUse[player] = Bar[bool]{}
	g.State.CannotPlayTypeThis[player] = Bar[CardType]{}
	// Roll this turn's history into "last turn" so the next player can ask what their
	// opponent just did.
	h := &g.State.TurnHistory
	h[player][KeysForgedLastTurn] = h[player][KeysForgedThisTurn]
	h[player][KeysForgedThisTurn] = 0
	h[player][CreaturesPlayedLastTurn] = int8(g.creaturesPlayedThisTurn(player))
	h[0][EnemyCreaturesFightKilled] = 0
	h[1][EnemyCreaturesFightKilled] = 0
	g.State.MayFightHouse[player] = HouseNone
	g.State.MayFightAny[player] = false
	g.State.MayUseHouse[player] = HouseNone
	g.State.KeywordsLost = 0
	g.clearLasting(player)
	g.drawStep(player)
	g.logf("%s ends their turn", g.names[player])
	g.assertInvariants()
}

// drawStep draws the player back up to their hand size, reduced by their chains,
// then sheds one chain only if that reduction actually blocked a draw. A player
// draws one fewer card for every 6 chains they hold (1-6 chains cost one card, 7-12
// cost two, and so on), and removes a single chain only on a turn where the reduced
// draw kept them from taking a card they could otherwise have drawn.
func (g *Game) drawStep(player int) {
	chains := g.State.Chains[player]
	target := HandSize + g.drawModifier(player) - (chains+5)/6
	if target < 0 {
		target = 0
	}
	g.drawTo(player, target)
	// The reduction blocked a draw only when it left the player below a full hand
	// with cards still available to draw.
	if chains > 0 &&
		int(g.State.Hand[player].Count) < HandSize+g.drawModifier(player) &&
		g.canDraw(player) {
		g.State.Chains[player]--
		g.logf("%s sheds a chain (%d remaining)", g.names[player], g.State.Chains[player])
	}
}

// drawModifier sums the end-of-turn hand-refill changes that cards in play impose
// on player, so a full hand becomes HandSize + drawModifier (Mother, Succubus, The
// Howling Pit).
func (g *Game) drawModifier(player int) int {
	total := 0
	for owner := 0; owner < 2; owner++ {
		for _, id := range g.allInPlay(owner) {
			if m := g.cat.def(id).DrawModifier; m.Amount != 0 && m.affects(owner, player) {
				total += m.Amount
			}
		}
	}
	return total
}

// CannotFightNextTurn arms a fight bar on a player for their next turn.
func (g *Game) CannotFightNextTurn(player int, source LocalID) {
	g.State.CannotFightNext[player] = Bar[bool]{Value: true, Source: source}
}

// CannotPlayTypeNextTurn arms a play-type bar on a player for their next turn.
func (g *Game) CannotPlayTypeNextTurn(player int, t CardType, source LocalID) {
	g.State.CannotPlayTypeNext[player] = Bar[CardType]{Value: t, Source: source}
}

// CannotPlayTypeThisTurn bars a player from playing cards of the given type for
// the rest of the current turn (Treasure Map bars every type once it pays out).
func (g *Game) CannotPlayTypeThisTurn(player int, t CardType, source LocalID) {
	g.State.CannotPlayTypeThis[player] = Bar[CardType]{Value: t, Source: source}
}

// CannotUseNextTurn arms a use bar on a player for their next turn, stopping them
// reaping, fighting, or using an "Action:" ability (Skippy Timehog).
func (g *Game) CannotUseNextTurn(player int, source LocalID) {
	g.State.CannotUseNext[player] = Bar[bool]{Value: true, Source: source}
}

// SkipForgeStepNextTurn makes a player skip their forge-a-key step at the start of
// their next turn.
func (g *Game) SkipForgeStepNextTurn(player int, source LocalID) {
	g.State.SkipForgeNext[player] = Bar[bool]{Value: true, Source: source}
}

// GrantFightForHouse lets a player use creatures of house h to fight this turn
// even when h is not the active house. EndTurn clears the grant.
func (g *Game) GrantFightForHouse(player int, h House) {
	g.State.MayFightHouse[player] = h
	g.logf("%s's %s creatures may fight this turn", g.names[player], h)
}

// GrantFightAnyHouse lets every creature a player controls fight this turn, whatever
// its house (Follow the Leader). EndTurn clears the grant.
func (g *Game) GrantFightAnyHouse(player int) {
	g.State.MayFightAny[player] = true
	g.logf("%s's creatures may all fight this turn", g.names[player])
}

// GrantUseForHouse lets a player fully use (fight, reap, or Action:) creatures of
// house h this turn even when h is not the active house. EndTurn clears the grant.
func (g *Game) GrantUseForHouse(player int, h House) {
	g.State.MayUseHouse[player] = h
	g.logf("%s may use %s creatures this turn", g.names[player], h)
}

// ForceActiveHouseNextTurn makes a player have to choose house h as their active
// house on their next turn (Control the Weak). BeginTurn promotes the armed house.
func (g *Game) ForceActiveHouseNextTurn(player int, h House, source LocalID) {
	g.State.ForcedHouseNext[player] = Bar[House]{Value: h, Source: source}
	g.logf("%s must choose house %s next turn", g.names[player], h)
}

// RestrictionSources returns the cards imposing a turn-scoped restriction on a
// player right now, so a frontend can remind them which cards are binding them.
// It reads the bars themselves, so a bar that has been lifted stops naming its
// card, and one card imposing two bars is named once.
func (g *Game) RestrictionSources(player int) []LocalID {
	var out []LocalID
	name := func(id LocalID) {
		if !slices.Contains(out, id) {
			out = append(out, id)
		}
	}
	if g.State.CannotFight[player].Value {
		name(g.State.CannotFight[player].Source)
	}
	if g.State.CannotPlayTypeThis[player].Value != TypeUnset {
		name(g.State.CannotPlayTypeThis[player].Source)
	}
	if g.State.ForcedHouse[player].Value != HouseNone {
		name(g.State.ForcedHouse[player].Source)
	}
	if g.State.SkipForge[player].Value {
		name(g.State.SkipForge[player].Source)
	}
	return out
}

// Forge a key: at the start of your turn you forge a single key if you can pay
// its current cost — 6 Aember by default. A player forges at most one key per turn.
// Keys are the win condition — forge your third key and you win the game.
//
//rulebook:turn Turn structure / 1. Forge a key

// forgeKey forges one key when the player can afford the current key cost, paying
// it and firing "after you forge a key" abilities. BeginTurn forges at most one
// key at the start of a turn; cards may forge one more via the ForgeKey effect.
func (g *Game) forgeKey(player int) {
	cost := g.keyCost(player)
	if g.spendableAember(player) < cost {
		return
	}
	// The colour is settled before the Æmber leaves the pool, so a forge is one
	// step: a player looking at the colour prompt has not paid for anything yet.
	color, ok := g.pickKeyColor(player)
	g.payKeyCost(player, cost)
	g.finishForgeKey(player, color, ok)
}

// spendableAember is everything a player can put toward a key: their pool plus
// the Æmber banked on cards that let it be spent when forging (Safe Place).
func (g *Game) spendableAember(player int) int {
	total := g.State.Aember[player]
	for _, id := range g.vaults(player) {
		total += g.AmberOn(id)
	}
	return total
}

// payKeyCost takes the cost out of the player's pool first, falling back to the
// Æmber banked on their vault cards — pool Æmber is the more exposed of the two,
// so spending it first is what a player wants.
func (g *Game) payKeyCost(player, cost int) {
	fromPool := min(cost, g.State.Aember[player])
	g.State.Aember[player] -= fromPool
	cost -= fromPool
	for _, id := range g.vaults(player) {
		if cost == 0 {
			return
		}
		taken := min(cost, g.AmberOn(id))
		g.AddAmberOn(id, -taken)
		cost -= taken
	}
}

// vaults returns the player's in-play cards whose Æmber may be spent on a key.
func (g *Game) vaults(player int) []LocalID {
	var out []LocalID
	for _, id := range g.allInPlay(player) {
		if g.cat.def(id).SpendableAember {
			out = append(out, id)
		}
	}
	return out
}

// forgeKeyFree forges one key without paying its current cost.
func (g *Game) forgeKeyFree(player int) {
	color, ok := g.pickKeyColor(player)
	g.finishForgeKey(player, color, ok)
}

// finishForgeKey records a newly forged key in the colour already picked, fires
// "after you forge a key" abilities, and checks for the win. hasColor is false
// only when every colour is already spent, which leaves the key colourless.
func (g *Game) finishForgeKey(player int, color KeyColor, hasColor bool) {
	g.State.Keys[player]++
	g.State.TurnHistory[player][KeysForgedThisTurn]++
	if hasColor {
		g.State.KeyColors[player][g.State.Keys[player]-1] = color
		g.logf("%s forges a %s key", g.names[player], color)
	}
	g.logf("%s forges a key (%d/%d)", g.names[player], g.State.Keys[player], KeysToWin)
	for _, id := range g.allInPlay(player) {
		g.triggerAbilities(id, TriggerAfterForgeKey, 0, false)
	}
	g.emitLasting(EventForgeKey, player, 0)
	if g.State.Keys[player] >= KeysToWin {
		g.State.Winner = player
		g.logf("%s wins the game!", g.names[player])
	}
}

// pickKeyColor asks the player which colour the key they are forging should be,
// choosing among the colours they have not forged yet, and reports whether one
// was available. The final key's colour is forced (only one remains), so it is
// taken without a prompt. There is no default: every UI is asked whenever more
// than one colour is available.
func (g *Game) pickKeyColor(player int) (KeyColor, bool) {
	remaining := g.remainingKeyColors(player)
	if len(remaining) == 0 {
		return 0, false
	}
	choice := remaining[0]
	if len(remaining) > 1 {
		labels := make([]string, len(remaining))
		for i, c := range remaining {
			labels[i] = c.String()
		}
		if idx := g.chooseOption(
			player,
			"",
			KeyColorPrompt,
			labels,
		); idx >= 0 &&
			idx < len(remaining) {
			choice = remaining[idx]
		}
	}
	return choice, true
}

// remainingKeyColors lists the key colours the player has not yet forged, in
// canonical order.
func (g *Game) remainingKeyColors(player int) []KeyColor {
	var used [4]bool
	for i := 0; i < g.State.Keys[player]; i++ {
		used[g.State.KeyColors[player][i]] = true
	}
	var out []KeyColor
	for _, c := range keyColorOrder {
		if !used[c] {
			out = append(out, c)
		}
	}
	return out
}

// UnforgeKey takes one forged key back off a player (Key Hammer). Unlike a forge
// it is silent: no cost is refunded and no forge abilities fire.
func (g *Game) UnforgeKey(player int) {
	if g.State.Keys[player] == 0 {
		return
	}
	g.State.Keys[player]--
	g.State.KeyColors[player][g.State.Keys[player]] = KeyColor(0)
	g.logf("%s unforges a key (%d/%d)", g.names[player], g.State.Keys[player], KeysToWin)
}
