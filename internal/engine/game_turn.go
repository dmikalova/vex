package engine

import (
	"errors"
	"slices"
)

// Engine action errors.
var (
	ErrNotActivePlayer    = errors.New("not the active player")
	ErrCardNotInHand      = errors.New("card index is not in hand")
	ErrWrongType          = errors.New("card is the wrong type for this action")
	ErrWrongHouse         = errors.New("card's house is not the active house")
	ErrCardExhausted      = errors.New("card is exhausted")
	ErrNoTarget           = errors.New("no legal target")
	ErrGameOver           = errors.New("the game is already over")
	ErrCannotFight        = errors.New("cannot use creatures to fight this turn")
	ErrCannotPlayCreature = errors.New("cannot play creatures")
	ErrCannotPlayType     = errors.New("cannot play cards of this type this turn")
	ErrCardPlayLimit      = errors.New("card-play limit reached this turn")
	ErrFirstTurnOneCard   = errors.New(
		"the first player may play or discard only one card on their first turn",
	)
	ErrCannotPayToll         = errors.New("cannot pay the toll for this action")
	ErrPlayRequirement       = errors.New("not enough Æmber to play this card")
	ErrMustChooseForcedHouse = errors.New("must choose the forced active house this turn")
	ErrHouseLocked           = errors.New("a card in play locks the active house choice")
	ErrHouseForbidden        = errors.New("a card forbids choosing this active house this turn")
	ErrCannotUse             = errors.New("card's use condition is not met")
)

// This file holds the turn lifecycle entry points — the three points at which a
// turn waits for a player and then resumes: begin the turn, choose the active
// house, end the turn. Each resets what it owns and hands back to the phase loop
// in game_phase.go. The key-forging step that decides the game lives here too.

// StartTurn starts a player's turn: it becomes the active player, the active
// house is cleared, and the phase loop runs from the start-of-turn phase through
// the forge phase, stopping when it needs the player's house choice.
func (g *Game) StartTurn(player int) {
	if g.State.Winner >= 0 {
		return
	}
	g.State.ActivePlayer = player
	g.State.ActiveHouse = HouseNone
	g.State.Turn++
	g.State.PlayedThisTurn[player].reset()
	g.State.DiscardedThisTurn[player].reset()
	g.State.PlayPermissionsUsedThisTurn[player] = [NumHouses]uint8{}
	g.State.NonActivePlaysUsedThisTurn[player] = 0
	g.State.FirstTurnPlayLimit[player] = false
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
	g.State.CannotReap[player] = g.State.CannotReapNext[player]
	g.State.CannotReapNext[player] = Bar[bool]{}
	g.State.CannotReapHouse[player] = g.State.CannotReapHouseNext[player]
	g.State.CannotReapHouseNext[player] = Bar[House]{}
	g.State.CreaturesCannot[player] = g.State.CreaturesCannotNext[player]
	g.State.CreaturesCannotNext[player] = Bar[CreatureBar]{}
	g.State.ForcedHouse[player] = g.State.ForcedHouseNext[player]
	g.State.ForcedHouseNext[player] = Bar[House]{}
	g.State.ForbiddenHouse[player] = g.State.ForbiddenHouseNext[player]
	g.State.ForbiddenHouseNext[player] = Bar[House]{}
	g.State.HouseWager[player] = g.State.HouseWagerNext[player]
	g.State.HouseWagerNext[player] = HouseWager{}
	g.State.SkipForge[player] = g.State.SkipForgeNext[player]
	g.State.SkipForgeNext[player] = Bar[bool]{}
	g.State.KeyCostBump[player] = g.State.KeyCostBumpNext[player]
	g.State.KeyCostBumpNext[player] = Bar[int]{}
	g.record(TurnBegan{Player: player, Turn: g.State.Turn})
	g.enterPhase(PhaseStartOfTurn)
	g.runPhases()
	g.assertInvariants()
}

// Choose a house: pick one of your deck's three houses to be your active house
// for this turn. For the rest of the turn you may play from hand and use only
// cards of that house, except cards that ignore the restriction such as Versatile
// ones.
// ChooseHouse sets ActiveHouse, then hands back to the phase loop, which offers
// the player's archived cards into hand now that a house is locked in and settles
// on the play phase.
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
	if g.State.ForbiddenHouse[player].Value == house && house != HouseNone {
		return ErrHouseForbidden
	}
	if !g.houseLockAllows(player, house) {
		return ErrHouseLocked
	}
	g.State.ActiveHouse = house
	g.record(HouseChosen{Player: player, House: house})
	g.payOffHouseWager(player, house)
	// The snapshot is taken once, but an earlier card's ability can remove a later
	// one from play (Strange Gizmo destroys friendly artifacts); triggerAbilitiesAs
	// drops the trigger of a card that has left play mid-window (ADR 0030).
	for _, id := range g.allInPlay(player) {
		g.triggerAbilities(id, TriggerAfterChooseHouse, 0, false)
	}
	// A house choice is public, so cards on either side may react to it (Snag's
	// Mirror bars the chooser's opponent from the same house next turn). This
	// window fires for both players' cards, unlike the controller-only
	// AfterChooseHouse above.
	for _, p := range [2]int{player, 1 - player} {
		for _, id := range g.allInPlay(p) {
			g.triggerAbilities(id, TriggerAfterAnyPlayerChoosesHouse, 0, false)
		}
	}
	g.enterPhase(PhaseArchives)
	g.runPhases()
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
// EndPlayPhase ends the active player's play phase and runs the turn out: ready,
// draw, then the end-of-turn abilities. It is the third and last point at which
// a turn waits for the player — nothing after it needs a decision, so the turn
// finishes in one call.
func (g *Game) EndPlayPhase(player int) {
	// The named player is the one whose turn runs out: the phases that follow act
	// on the active player, so naming a player here is what makes them active.
	g.State.ActivePlayer = player
	g.enterPhase(PhaseReady)
	g.runPhases()
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
	before := int(g.State.Hand[player].Count)
	g.drawTo(player, target)
	hand := int(g.State.Hand[player].Count)
	g.record(CardsDrawn{Player: player, Count: hand - before, Hand: hand})
	// The reduction blocked a draw only when it left the player below a full hand
	// with cards still available to draw.
	if chains > 0 &&
		int(g.State.Hand[player].Count) < HandSize+g.drawModifier(player) &&
		g.canDraw(player) {
		g.State.Chains[player]--
		g.record(ChainShed{Player: player, Remaining: g.State.Chains[player]})
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
				if m.OnlyWhileOffFlank && g.onFlankOf(id) {
					continue
				}
				if m.OnlyWhileInCenter && !g.InCenterOfBattleline(id) {
					continue
				}
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

// CannotUseThisTurn bars a player from reaping, fighting, or using an "Action:"
// ability for the rest of the current turn (United Action). It sets the use bar
// directly rather than arming the next-turn form, so the ready phase lifts it at
// the end of this turn.
func (g *Game) CannotUseThisTurn(player int, source LocalID) {
	g.State.CannotUse[player] = Bar[bool]{Value: true, Source: source}
}

// CannotReapNextTurn arms a bar that stops a player reaping with any creature
// throughout their next turn (Inky Gloom). StartTurn promotes it and the ready
// phase lifts it. It is narrower than CannotUseNextTurn: fighting and "Action:"
// abilities stay open.
func (g *Game) CannotReapNextTurn(player int, source LocalID) {
	g.State.CannotReapNext[player] = Bar[bool]{Value: true, Source: source}
}

// CannotReapThisTurn bars a player from reaping with any creature for the rest of
// the current turn (Ragnarok). It sets the reap bar directly rather than arming
// the next-turn form, so the ready phase lifts it at the end of this turn. Only
// the active player can reap on their own turn, so barring them stops reaping this
// turn.
func (g *Game) CannotReapThisTurn(player int, source LocalID) {
	g.State.CannotReap[player] = Bar[bool]{Value: true, Source: source}
}

// CannotReapHouseNextTurn arms a bar that stops a player reaping with creatures
// of house h throughout their next turn (Seismo-entangler). StartTurn promotes
// the armed house.
func (g *Game) CannotReapHouseNextTurn(player int, h House, source LocalID) {
	g.State.CannotReapHouseNext[player] = Bar[House]{Value: h, Source: source}
}

// CreaturesCannotUntilNextTurn arms a board-wide bar that stops both players
// using creatures one way — fighting or reaping — until the caster's next turn,
// sparing creatures of exceptHouse (Into the Night, Sow Salt). The caster is
// barred for the rest of this turn and the opponent for their next turn, so the
// bar lifts at the start of the caster's next turn.
func (g *Game) CreaturesCannotUntilNextTurn(
	caster int,
	action UseKind,
	exceptHouse House,
	source LocalID,
) {
	bar := Bar[CreatureBar]{
		Value:  CreatureBar{Action: action, ExceptHouse: exceptHouse},
		Source: source,
	}
	g.State.CreaturesCannot[caster] = bar
	g.State.CreaturesCannotNext[1-caster] = bar
}

// BlankEnemyText blanks the text box of every creature the given player controls
// until the source's controller's next turn (Shadow of Dis). Set on the affected
// player, the blank persists through that player's own next turn and is lifted by
// their ready phase, so it spans exactly through the opponent's turn.
func (g *Game) BlankEnemyText(player int, source LocalID) {
	g.State.TextBlank[player] = Bar[bool]{Value: true, Source: source}
}

// SkipForgePhaseNextTurn makes a player skip their forge-a-key phase at the start of
// their next turn.
func (g *Game) SkipForgePhaseNextTurn(player int, source LocalID) {
	g.State.SkipForgeNext[player] = Bar[bool]{Value: true, Source: source}
}

// ScheduleDestroyEachCreatureAtEndOfTurn arms "destroy each creature" to resolve in
// the active player's end-of-turn phase (Ragnarok). The flag must survive the ready
// phase, which runs before end of turn, so endOfTurnPhase clears it as it fires.
func (g *Game) ScheduleDestroyEachCreatureAtEndOfTurn(source LocalID) {
	g.State.EndOfTurnDestroyAll = Bar[bool]{Value: true, Source: source}
}

// RaiseKeyCostNextTurn raises what a player's keys cost throughout their next turn
// (Lash of Broken Dreams). Successive raises stack.
func (g *Game) RaiseKeyCostNextTurn(player, amount int, source LocalID) {
	g.State.KeyCostBumpNext[player] = Bar[int]{
		Value:  g.State.KeyCostBumpNext[player].Value + amount,
		Source: source,
	}
}

// RaiseKeyCostThisTurn raises what a player's keys cost for the remainder of the
// current turn. Successive raises stack.
func (g *Game) RaiseKeyCostThisTurn(player, amount int, source LocalID) {
	g.State.KeyCostBump[player] = Bar[int]{
		Value:  g.State.KeyCostBump[player].Value + amount,
		Source: source,
	}
}

// GrantMayPlayOrUse records a this-turn grant letting a player act with cards
// outside their active house, folding the whole out-of-house permission family
// (ADR 0037). It dispatches the axes onto the flat this-turn state slots — a fight
// grant onto MayFightHouse/MayFightAny, a use/play grant onto MayUseHouse/
// MayPlayHouse or MayUseArtifactsAnyHouse, and an exclusion or controlled grant
// onto an off-house permit — and records one MayPlayOrUseGranted. The ready phase
// clears every slot.
func (g *Game) GrantMayPlayOrUse(
	player int,
	houses HouseSelector,
	grant HouseGrant,
	types CardTypes,
	count int,
) {
	switch houses.Kind {
	case SelectAny:
		if grant&GrantFight != 0 {
			g.State.MayFightAny[player] = true
		}
		if grant&GrantUse != 0 {
			g.State.MayUseArtifactsAnyHouse[player] = true
		}
	case SelectHouse:
		if grant&GrantFight != 0 {
			g.State.MayFightHouse[player] = houses.House
		}
		if grant&GrantUse != 0 {
			g.State.MayUseHouse[player] = houses.House
		}
		if grant&GrantPlay != 0 {
			g.State.MayPlayHouse[player] = houses.House
		}
	case SelectExcept, SelectControlled:
		rem := permitUnlimited
		if count > 0 {
			rem = uint8(count)
		}
		except := HouseNone
		if houses.Kind == SelectExcept {
			except = houses.House
		}
		g.addOffHousePermit(player, OffHousePermit{
			Except:     except,
			Controlled: houses.Kind == SelectControlled,
			Types:      types,
			Grant:      grant,
			Remaining:  rem,
		})
	}
	g.record(MayPlayOrUseGranted{
		Player: player,
		Houses: houses,
		Grant:  grant,
		Types:  types,
		Count:  count,
	})
}

// ForceActiveHouseNextTurn makes a player have to choose house h as their active
// house on their next turn (Control the Weak). StartTurn promotes the armed house.
func (g *Game) ForceActiveHouseNextTurn(player int, h House, source LocalID) {
	g.State.ForcedHouseNext[player] = Bar[House]{Value: h, Source: source}
	g.record(HouseForcedNextTurn{Player: player, House: h})
}

// ForbidActiveHouseNextTurn makes a player unable to choose house h as their
// active house on their next turn (Tezmal). StartTurn promotes the armed house.
func (g *Game) ForbidActiveHouseNextTurn(player int, h House, source LocalID) {
	g.State.ForbiddenHouseNext[player] = Bar[House]{Value: h, Source: source}
	g.record(HouseForbiddenNextTurn{Player: player, House: h})
}

// WagerOnHouseNextTurn arms a bet on player's next active house: if they choose h
// then, predictor steals amount (Snaglet). StartTurn promotes it so the payoff
// lands when player next chooses a house.
func (g *Game) WagerOnHouseNextTurn(player int, h House, amount, predictor int, source LocalID) {
	g.State.HouseWagerNext[player] = HouseWager{
		House:     h,
		Amount:    amount,
		Predictor: predictor,
		Source:    source,
	}
	g.record(HouseWagerArmed{Predictor: predictor, Player: player, House: h, Amount: amount})
}

// payOffHouseWager settles any wager armed on player once they lock in house: a
// matching choice pays the predictor a steal, and either way the wager is spent.
func (g *Game) payOffHouseWager(player int, house House) {
	w := g.State.HouseWager[player]
	g.State.HouseWager[player] = HouseWager{}
	if w.Amount == 0 || w.House != house {
		return
	}
	ctx := &EffectContext{Resolver: g, Controller: w.Predictor, Source: w.Source}
	StealAember{Amount: w.Amount}.Resolve(ctx)
}

// RestrictionSources returns the cards restricting a player right now, so a
// frontend can remind them which cards are binding them. It reads the bars
// themselves — the turn-scoped State bars plus any continuous CannotPlayWhile bar
// whose condition currently holds — so a bar that has been lifted stops naming its
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
	// A symmetric CannotPlayWhile bar (Quixxle Stone) is continuous, not
	// turn-scoped, so it is not in State; name each in-play card whose bar
	// currently holds against this player.
	for p := 0; p < 2; p++ {
		for _, id := range g.allInPlay(p) {
			bar := g.cat.def(id).CannotPlayWhile
			if bar.When == nil {
				continue
			}
			ctx := &EffectContext{Resolver: g, Source: id, Controller: player}
			if bar.When.Met(ctx) {
				name(id)
			}
		}
	}
	return out
}

// KeyCostSources returns the cards changing a player's key cost right now — the
// turn-scoped bump (RaiseKeyCost) plus every in-play card whose continuous
// key-cost change currently applies to that player — so a frontend can name them
// on the key-cost pill rather than mixing them into the restriction list.
func (g *Game) KeyCostSources(player int) []LocalID {
	var out []LocalID
	name := func(id LocalID) {
		if !slices.Contains(out, id) {
			out = append(out, id)
		}
	}
	if g.State.KeyCostBump[player].Value != 0 {
		name(g.State.KeyCostBump[player].Source)
	}
	for controller := 0; controller < 2; controller++ {
		for _, id := range g.allInPlay(controller) {
			if g.keyCostChangeFor(id, controller, player) != 0 {
				name(id)
			}
		}
	}
	return out
}

// Forge a key: at the start of your turn you forge a single key if you can pay
// its current cost — 6 Æmber by default. A player forges at most one key per turn.
// Keys are the win condition — forge your third key and you win the game.
// forgeKey forges one key when the player can afford the current key cost, paying
// it and firing "after you forge a key" abilities. StartTurn forges at most one
// key at the start of a turn; cards may forge one more via the ForgeKey effect.
func (g *Game) forgeKey(player int) {
	g.forgeKeyAtExtraCost(player, 0)
}

// forgeKeyAtExtraCost forges one key at the current cost plus a surcharge for
// this forge alone, doing nothing when the player cannot afford the total. It
// reports whether a key was forged, so "forge a key … if you do, destroy Obsidian
// Forge" can gate on the forge actually happening.
func (g *Game) forgeKeyAtExtraCost(player, extra int) bool {
	if g.forgeKeyNumberBarred(player) {
		return false
	}
	cost := g.keyCost(player) + extra
	if g.spendableAember(player) < cost {
		return false
	}
	// An opponent's forge guard (Keyforgery) may prevent the forge here, before any
	// Æmber leaves the pool, so a prevented forge costs the player nothing.
	if g.opponentForgeGuarded(player) {
		return false
	}
	// The colour is settled before the Æmber leaves the pool, so a forge is one
	// step: a player looking at the colour prompt has not paid for anything yet.
	color, ok := g.pickKeyColor(player)
	g.payKeyCost(player, cost)
	g.finishForgeKey(player, color, ok)
	return true
}

// spendableAember is everything a player can put toward a key: their pool, the
// Æmber banked on cards that let it be spent when forging (Safe Place), and the
// Æmber on creatures a spend-as-pool permission covers (Senator Bracchus).
func (g *Game) spendableAember(player int) int {
	total := g.State.Aember[player]
	for _, id := range g.vaults(player) {
		total += g.AmberOn(id)
	}
	total += g.spendAsPoolTotal(player)
	return total
}

// payKeyCost takes the cost out of the player's pool first, falling back to the
// Æmber banked on their vault cards — pool Æmber is the more exposed of the two,
// so spending it first is what a player wants. If the opponent controls a card
// that gains forge spending (The Sting), the whole cost goes to them instead of
// vanishing.
func (g *Game) payKeyCost(player, cost int) {
	total := cost
	fromPool := min(cost, g.State.Aember[player])
	g.State.Aember[player] -= fromPool
	cost -= fromPool
	for _, id := range g.vaults(player) {
		if cost == 0 {
			break
		}
		taken := min(cost, g.AmberOn(id))
		g.AddAmberOn(id, -taken)
		cost -= taken
	}
	if cost > 0 {
		g.drawFromSpendAsPool(player, cost)
	}
	if gainer, ok := g.forgeAemberGainer(player); ok && total > 0 {
		beneficiary := g.controller(gainer)
		g.State.Aember[beneficiary] += total
		g.record(AemberGainedFromForging{Card: gainer, From: player, Amount: total})
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

// forgeKeyFree forges one key without paying its current cost, reporting whether
// a key was forged.
func (g *Game) forgeKeyFree(player int) bool {
	if g.forgeKeyNumberBarred(player) {
		return false
	}
	if g.opponentForgeGuarded(player) {
		return false
	}
	color, ok := g.pickKeyColor(player)
	g.finishForgeKey(player, color, ok)
	return true
}

// finishForgeKey records a newly forged key in the colour already picked, fires
// "after you forge a key" abilities, and checks for the win. hasColor is false
// only when every colour is already spent, which leaves the key colourless.
func (g *Game) finishForgeKey(player int, color KeyColor, hasColor bool) {
	g.State.Keys[player]++
	g.State.TurnHistory[player][KeysForgedThisTurn]++
	if hasColor {
		g.State.KeyColors[player][g.State.Keys[player]-1] = color
	}
	g.record(KeyForged{
		Player:   player,
		Color:    color,
		HasColor: hasColor,
		Keys:     g.State.Keys[player],
		Needed:   KeysToWin,
	})
	// Forging opens one window: the forger's own "after you forge a key" abilities,
	// every card's "after a player forges a key" abilities, and the forge-key lasting
	// reactions all trigger at once, so the forger orders the whole set (ADR 0013).
	pending := g.forgeKeyReactions(player)
	pending = append(pending, g.lastingReactions(EventForgeKey, player, 0)...)
	g.resolveWindow(g.orderTriggered(player, pending))
	// Forging changes the unforged-key count some creatures draw their power from.
	g.settleDestroyed(player)
	if g.State.Keys[player] >= KeysToWin {
		g.State.Winner = player
		g.record(GameWon{Player: player})
	}
}

// forgeKeyReactions gathers, as one ordered window, every reaction to the forger
// forging a key: the forger's own "after you forge a key" abilities (Redeemer
// Amara) and every card's "after a player forges a key" abilities (Forgemaster Og),
// the latter resolving for the forger so its "that player" names whoever forged.
// Gathering them lets the forger order the set when several fire at once (ADR 0013);
// the EventForgeKey lasting reactions are folded in by the caller.
func (g *Game) forgeKeyReactions(forger int) []triggeredAbility {
	w := g.window()
	for _, id := range g.allInPlay(forger) {
		w.add(id, TriggerAfterForgeKey, 0, false)
	}
	for _, p := range []int{forger, 1 - forger} {
		for _, id := range g.allInPlay(p) {
			w.addAs(id, TriggerAfterPlayerForgesKey, forger, 0, false)
		}
	}
	return w.pending
}

// Concede forfeits the game for player: their opponent becomes the winner. It is
// a no-op once the game is already decided.
func (g *Game) Concede(player int) {
	if g.State.Winner >= 0 {
		return
	}
	g.record(PlayerConceded{Player: player})
	g.State.Winner = 1 - player
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
	g.record(KeyUnforged{Player: player, Keys: g.State.Keys[player], Needed: KeysToWin})
}
