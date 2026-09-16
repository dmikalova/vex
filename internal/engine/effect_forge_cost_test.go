package engine

import "testing"

func TestForgeKeyExtraCost(t *testing.T) {
	cases := []struct {
		name   string
		effect ForgeKey
		want   string
	}{
		{"current cost", ForgeKey{}, "forge a key at current cost -> purge {self}"},
		{"free", ForgeKey{FreeOfCost: true}, "forge a key at no cost -> purge {self}"},
		{"surcharge", ForgeKey{Extra: 6}, "forge a key at +6 Æmber current cost -> purge {self}"},
		{
			"reduced surcharge",
			ForgeKey{Extra: 9, ReducedBy: CardsInHand{Player: Controller, House: AnyHouse}},
			"forge a key at +9 Æmber current cost, reduced by 1 Æmber for each card in your hand -> purge {self}",
		},
		{
			"discount",
			ForgeKey{Discount: true, ReducedBy: CardsInHand{Player: Controller, House: AnyHouse}},
			"forge a key at current cost, reduced by 1 Æmber for each card in your hand -> purge {self}",
		},
		{"kept", ForgeKey{Keep: true}, "forge a key at current cost"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.effect.Text(); got != tc.want {
				t.Errorf("text = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestForgeKeyValidate(t *testing.T) {
	if err := (ForgeKey{ReducedBy: CardsInHand{}}).validate(); err == nil {
		t.Error("a reduction with no Extra should not validate")
	}
	if err := (ForgeKey{FreeOfCost: true, Extra: 2}).validate(); err == nil {
		t.Error("a free forge with an Extra should not validate")
	}
	if err := (ForgeKey{Extra: 6}).validate(); err != nil {
		t.Errorf("validate = %v, want nil", err)
	}
	if err := (ForgeKey{Discount: true}).validate(); err == nil {
		t.Error("a Discount with no ReducedBy should not validate")
	}
	if err := (ForgeKey{
		Discount:  true,
		Extra:     2,
		ReducedBy: CardsInHand{Player: Controller, House: AnyHouse},
	}).validate(); err == nil {
		t.Error("a Discount forge with an Extra should not validate")
	}
	if err := (ForgeKey{
		Discount:  true,
		ReducedBy: CardsInHand{Player: Controller, House: AnyHouse},
	}).validate(); err != nil {
		t.Errorf("validate = %v, want nil", err)
	}
}

func TestForgeKeyPaysTheSurcharge(t *testing.T) {
	g := started(t)
	g.State.Aember[0] = KeyCost + 2
	ctx := &EffectContext{Resolver: g, Controller: 0}

	ForgeKey{Extra: 6}.Resolve(ctx)
	if g.State.Keys[0] != 0 {
		t.Fatalf("keys = %d, want 0 — the surcharge is unaffordable", g.State.Keys[0])
	}

	ForgeKey{Extra: 2}.Resolve(ctx)
	if g.State.Keys[0] != 1 {
		t.Errorf("keys = %d, want 1", g.State.Keys[0])
	}
	if g.State.Aember[0] != 0 {
		t.Errorf("Æmber = %d, want 0", g.State.Aember[0])
	}
}

func TestForgeKeyReducedBelowTheSurcharge(t *testing.T) {
	g := started(t)
	g.State.Aember[0] = KeyCost
	for i := 0; i < 12; i++ {
		g.AddToHand(NewCard("Filler", Brobnar, Tactic, Common), 0)
	}

	// A +9 surcharge reduced by 12 cards in hand cannot drop below the key cost.
	ForgeKey{
		Extra:     9,
		ReducedBy: CardsInHand{Player: Controller, House: AnyHouse},
	}.Resolve(&EffectContext{Resolver: g, Controller: 0})

	if g.State.Keys[0] != 1 {
		t.Errorf("keys = %d, want 1", g.State.Keys[0])
	}
	if g.State.Aember[0] != 0 {
		t.Errorf("Æmber = %d, want 0 — the cost floor is the key cost", g.State.Aember[0])
	}
}

func TestForgeKeyDiscountFloorsAndKeeps(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddArtifact(NewCard("Forge", Dis, Artifact, Common), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	// A discount larger than the current key cost floors the whole cost at 0, so the
	// forge lands with an empty pool, and Keep leaves the source in play.
	for i := 0; i < 8; i++ {
		g.AddToHand(NewCard("Filler", Brobnar, Tactic, Common), 0)
	}
	ForgeKey{
		Discount:  true,
		Keep:      true,
		ReducedBy: CardsInHand{Player: Controller, House: AnyHouse},
	}.Resolve(ctx)

	if g.Keys(0) != 1 {
		t.Errorf("keys = %d, want 1 — the discount floors the cost at 0", g.Keys(0))
	}
	if !g.inPlay(src) {
		t.Error("source should stay in play with Keep")
	}
}

func TestRaiseKeyCost(t *testing.T) {
	effect := RaiseKeyCost{Player: Opponent, Amount: 3, Duration: OpponentNextTurn}
	if got := effect.Text(); got != "keys cost +3 Æmber during your opponent's next turn" {
		t.Errorf("text = %q", got)
	}
	if got := (RaiseKeyCost{Player: Controller, Amount: 1, Duration: OpponentNextTurn}).Text(); got !=
		"keys cost +1 Æmber during your next turn" {
		t.Errorf("controller text = %q", got)
	}
	if err := (RaiseKeyCost{Amount: 3, Duration: OpponentNextTurn}).validate(); err == nil {
		t.Error("an unset player should not validate")
	}
	if err := (RaiseKeyCost{Player: Opponent}).validate(); err == nil {
		t.Error("a zero raise should not validate")
	}
	if err := effect.validate(); err != nil {
		t.Errorf("validate = %v, want nil", err)
	}
	if err := (RaiseKeyCost{Player: Opponent, Amount: 3}).validate(); err == nil {
		t.Error("an unset duration should not validate")
	}
	if err := (RaiseKeyCost{
		Player:   Opponent,
		Amount:   3,
		Duration: Forever,
	}).validate(); err == nil {
		t.Error("a duration the surcharge cannot express should not validate")
	}
}

// TestRaiseKeyCostThisTurn checks the RemainderOfPlayerTurn duration bites
// immediately and lifts when the current turn ends, rather than waiting for a
// turn boundary.
func TestRaiseKeyCostThisTurn(t *testing.T) {
	effect := RaiseKeyCost{Player: Controller, Amount: 2, Duration: RemainderOfPlayerTurn}
	if got := effect.Text(); got != "your keys cost +2 Æmber for the remainder of the turn" {
		t.Errorf("text = %q", got)
	}
	if got := (RaiseKeyCost{
		Player:   Opponent,
		Amount:   2,
		Duration: RemainderOfPlayerTurn,
	}).Text(); got != "your opponent's keys cost +2 Æmber for the remainder of the turn" {
		t.Errorf("opponent text = %q", got)
	}

	g := started(t)
	source := g.AddToBattleline(NewCard("Lash", Dis, Creature, Common, WithPower(1)), 0)
	effect.Resolve(&EffectContext{Resolver: g, Controller: 0, Source: source})

	if got := g.CurrentKeyCost(0); got != KeyCost+2 {
		t.Errorf("key cost = %d, want %d right away", got, KeyCost+2)
	}
	g.EndPlayPhase(0)
	if got := g.CurrentKeyCost(0); got != KeyCost {
		t.Errorf("key cost = %d, want %d after the turn ends", got, KeyCost)
	}
}

func TestRaiseKeyCostLandsOnTheNextTurn(t *testing.T) {
	g := started(t)
	source := g.AddToBattleline(NewCard("Lash", Dis, Creature, Common, WithPower(1)), 0)

	RaiseKeyCost{Player: Opponent, Amount: 3, Duration: OpponentNextTurn}.
		Resolve(&EffectContext{Resolver: g, Controller: 0, Source: source})

	if got := g.CurrentKeyCost(1); got != KeyCost {
		t.Errorf("key cost = %d, want %d before the raise lands", got, KeyCost)
	}
	g.EndPlayPhase(0)
	g.StartTurn(1)
	if got := g.CurrentKeyCost(1); got != KeyCost+3 {
		t.Errorf("key cost = %d, want %d", got, KeyCost+3)
	}
	if sources := g.KeyCostSources(1); len(sources) == 0 {
		t.Error("the raise should name its source on the key-cost pill")
	}
	if reasons := g.RestrictionSources(1); len(reasons) != 0 {
		t.Errorf("a key-cost raise should not be a restriction note, got %v", reasons)
	}
	g.EndPlayPhase(1)
	if got := g.CurrentKeyCost(1); got != KeyCost {
		t.Errorf("key cost = %d, want %d after the turn ends", got, KeyCost)
	}
}

// TestKeyCostSourcesNamesContinuousModifier checks the key-cost pill's reader
// names an in-play card whose continuous key-cost change binds a player, and
// names nothing for a player no modifier touches.
func TestKeyCostSourcesNamesContinuousModifier(t *testing.T) {
	g := NewGame("A", "B", 1)
	if got := g.KeyCostSources(1); len(got) != 0 {
		t.Errorf("sources with no modifier = %v, want none", got)
	}
	jammer := g.AddArtifact(
		NewCard("Test Jammer", Logos, Artifact, Common,
			WithKeyCost(NewKeyCostChange(Opponent, 1))),
		0,
	)
	if got := g.KeyCostSources(1); len(got) != 1 || got[0] != jammer {
		t.Errorf("sources against the opponent = %v, want [%d]", got, jammer)
	}
	// The controller's own key cost is untouched, so nothing is named for them.
	if got := g.KeyCostSources(0); len(got) != 0 {
		t.Errorf("sources for the controller = %v, want none", got)
	}
}

func TestRaiseKeyCostStacks(t *testing.T) {
	g := started(t)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	RaiseKeyCost{Player: Opponent, Amount: 3, Duration: OpponentNextTurn}.Resolve(ctx)
	RaiseKeyCost{Player: Opponent, Amount: 3, Duration: OpponentNextTurn}.Resolve(ctx)
	if got := g.State.KeyCostBumpNext[1].Value; got != 6 {
		t.Errorf("armed raise = %d, want 6", got)
	}
}

func TestRaiseKeyCostPerHouseCreature(t *testing.T) {
	effect := RaiseKeyCostPerHouseCreature{
		Player:   Opponent,
		Amount:   1,
		House:    Dis,
		Duration: OpponentNextTurn,
	}
	if got := effect.Text(); got !=
		"keys cost +1 Æmber for each Dis creature in play during your opponent's next turn" {
		t.Errorf("text = %q", got)
	}
	if got := (RaiseKeyCostPerHouseCreature{
		Player:   Controller,
		Amount:   2,
		House:    Dis,
		Duration: OpponentNextTurn,
	}).Text(); got !=
		"keys cost +2 Æmber for each Dis creature in play during your next turn" {
		t.Errorf("controller text = %q", got)
	}
	if err := (RaiseKeyCostPerHouseCreature{
		Amount:   1,
		House:    Dis,
		Duration: OpponentNextTurn,
	}).validate(); err == nil {
		t.Error("an unset player should not validate")
	}
	if err := (RaiseKeyCostPerHouseCreature{
		Player:   Opponent,
		House:    Dis,
		Duration: OpponentNextTurn,
	}).validate(); err == nil {
		t.Error("a zero raise should not validate")
	}
	if err := (RaiseKeyCostPerHouseCreature{
		Player:   Opponent,
		Amount:   1,
		Duration: OpponentNextTurn,
	}).validate(); err == nil {
		t.Error("an unset house should not validate")
	}
	if err := (RaiseKeyCostPerHouseCreature{
		Player: Opponent,
		Amount: 1,
		House:  Dis,
	}).validate(); err == nil {
		t.Error("an unset duration should not validate")
	}
	if err := (RaiseKeyCostPerHouseCreature{
		Player:   Opponent,
		Amount:   1,
		House:    Dis,
		Duration: RemainderOfPlayerTurn,
	}).validate(); err == nil {
		t.Error("an unsupported duration should not validate")
	}
	if err := effect.validate(); err != nil {
		t.Errorf("validate = %v, want nil", err)
	}
}

// TestRaiseKeyCostPerHouseCreatureCountsLive proves the surcharge is dormant until
// the taxed player's turn, then reads the number of Dis creatures in play live at
// each forge, and lifts when that turn ends.
func TestRaiseKeyCostPerHouseCreatureCountsLive(t *testing.T) {
	g := started(t)
	source := g.AddToBattleline(NewCard("Nightmare", Dis, Creature, Common, WithPower(1)), 0)
	g.AddToBattleline(NewCard("dis-a", Dis, Creature, Common, WithPower(1)), 0)
	g.AddToBattleline(NewCard("dis-b", Dis, Creature, Common, WithPower(1)), 1)
	g.AddToBattleline(NewCard("logos", Logos, Creature, Common, WithPower(1)), 1) // not counted

	RaiseKeyCostPerHouseCreature{
		Player:   Opponent,
		Amount:   1,
		House:    Dis,
		Duration: OpponentNextTurn,
	}.Resolve(&EffectContext{Resolver: g, Controller: 0, Source: source})

	if got := g.CurrentKeyCost(1); got != KeyCost {
		t.Errorf("key cost = %d, want %d before the surcharge lands", got, KeyCost)
	}
	g.EndPlayPhase(0)
	g.StartTurn(1)
	// Three Dis creatures in play (source, dis-a, dis-b) -> +3.
	if got := g.CurrentKeyCost(1); got != KeyCost+3 {
		t.Errorf("key cost = %d, want %d for 3 Dis creatures", got, KeyCost+3)
	}
	if sources := g.KeyCostSources(1); len(sources) == 0 {
		t.Error("the surcharge should name its source on the key-cost pill")
	}
	// Recomputed live: another Dis creature enters -> +4.
	g.AddToBattleline(NewCard("dis-c", Dis, Creature, Common, WithPower(1)), 1)
	if got := g.CurrentKeyCost(1); got != KeyCost+4 {
		t.Errorf("key cost = %d, want %d after another Dis creature enters", got, KeyCost+4)
	}
	g.EndPlayPhase(1)
	if got := g.CurrentKeyCost(1); got != KeyCost {
		t.Errorf("key cost = %d, want %d after the taxed turn ends", got, KeyCost)
	}
}

// TestRaiseKeyCostPerHouseCreatureStacks proves two surcharges of the same house
// sum their per-creature amount, while a different house replaces the surcharge.
func TestRaiseKeyCostPerHouseCreatureStacks(t *testing.T) {
	g := started(t)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	RaiseKeyCostPerHouseCreature{
		Player:   Opponent,
		Amount:   1,
		House:    Dis,
		Duration: OpponentNextTurn,
	}.Resolve(ctx)
	RaiseKeyCostPerHouseCreature{
		Player:   Opponent,
		Amount:   2,
		House:    Dis,
		Duration: OpponentNextTurn,
	}.Resolve(ctx)
	if got := g.State.KeyCostPerHouseNext[1].Value; got.Per != 3 || got.House != Dis {
		t.Errorf("armed surcharge = %+v, want {Dis 3}", got)
	}
	RaiseKeyCostPerHouseCreature{
		Player:   Opponent,
		Amount:   5,
		House:    Logos,
		Duration: OpponentNextTurn,
	}.Resolve(ctx)
	if got := g.State.KeyCostPerHouseNext[1].Value; got.Per != 5 || got.House != Logos {
		t.Errorf("armed surcharge = %+v, want {Logos 5}", got)
	}
}

// TestRaiseKeyCostUntilEndOfNextTurnIsLiveNow proves the EndOfPlayerNextTurn
// window bites the moment it resolves — so a forge forced on the controller's own
// turn (Keyfrog) already pays the surcharge — and stays live through the affected
// player's next turn.
func TestRaiseKeyCostUntilEndOfNextTurnIsLiveNow(t *testing.T) {
	g := started(t)
	source := g.AddToBattleline(NewCard("Lash", Dis, Creature, Common, WithPower(1)), 0)

	RaiseKeyCost{Player: Opponent, Amount: 3, Duration: EndOfPlayerNextTurn}.
		Resolve(&EffectContext{Resolver: g, Controller: 0, Source: source})

	if got := g.CurrentKeyCost(1); got != KeyCost+3 {
		t.Errorf("key cost = %d, want %d live on the controller's turn", got, KeyCost+3)
	}
	g.EndPlayPhase(0)
	g.StartTurn(1)
	if got := g.CurrentKeyCost(1); got != KeyCost+3 {
		t.Errorf("key cost = %d, want %d during the opponent's next turn", got, KeyCost+3)
	}
	g.EndPlayPhase(1)
	if got := g.CurrentKeyCost(1); got != KeyCost {
		t.Errorf("key cost = %d, want %d after the opponent's turn ends", got, KeyCost)
	}
}

// TestLowerKeyCostText covers the rendered sentence and validation.
func TestLowerKeyCostText(t *testing.T) {
	each := LowerKeyCost{Player: EachPlayer, Amount: 2, Duration: EndOfPlayerNextTurn}
	if got := each.Text(); got != "each player's keys cost -2 Æmber until the end of your next turn" {
		t.Errorf("text = %q", got)
	}
	if got := (LowerKeyCost{
		Player:   Controller,
		Amount:   1,
		Duration: RemainderOfPlayerTurn,
	}).Text(); got != "your keys cost -1 Æmber for the remainder of the turn" {
		t.Errorf("controller text = %q", got)
	}
	if got := (LowerKeyCost{
		Player:   Opponent,
		Amount:   3,
		Duration: OpponentNextTurn,
	}).Text(); got != "your opponent's keys cost -3 Æmber during your next turn" {
		t.Errorf("opponent text = %q", got)
	}
	if err := each.validate(); err != nil {
		t.Errorf("validate = %v, want nil", err)
	}
	if err := (LowerKeyCost{Amount: 2, Duration: EndOfPlayerNextTurn}).validate(); err == nil {
		t.Error("an unset player should not validate")
	}
	if err := (LowerKeyCost{Player: EachPlayer, Duration: EndOfPlayerNextTurn}).validate(); err == nil {
		t.Error("a zero drop should not validate")
	}
	if err := (LowerKeyCost{Player: EachPlayer, Amount: 2}).validate(); err == nil {
		t.Error("an unset duration should not validate")
	}
	if err := (LowerKeyCost{
		Player:   EachPlayer,
		Amount:   2,
		Duration: Forever,
	}).validate(); err == nil {
		t.Error("a duration the drop cannot express should not validate")
	}
}

// TestLowerKeyCostEachPlayer checks the drop is live for both players the moment
// it resolves and lasts through the controller's next turn.
func TestLowerKeyCostEachPlayer(t *testing.T) {
	g := started(t)
	source := g.AddToBattleline(NewCard("Win", StarAlliance, Creature, Common, WithPower(1)), 0)

	LowerKeyCost{Player: EachPlayer, Amount: 2, Duration: EndOfPlayerNextTurn}.
		Resolve(&EffectContext{Resolver: g, Controller: 0, Source: source})

	// Both players see the -2 immediately.
	if got := g.CurrentKeyCost(0); got != KeyCost-2 {
		t.Errorf("controller key cost = %d, want %d right away", got, KeyCost-2)
	}
	if got := g.CurrentKeyCost(1); got != KeyCost-2 {
		t.Errorf("opponent key cost = %d, want %d right away", got, KeyCost-2)
	}

	// The opponent's next turn (the one between) still carries the drop.
	g.EndPlayPhase(0)
	g.StartTurn(1)
	if got := g.CurrentKeyCost(1); got != KeyCost-2 {
		t.Errorf("opponent key cost = %d, want %d on their turn", got, KeyCost-2)
	}

	// The controller's next turn still carries the drop, then it lifts at its end.
	g.EndPlayPhase(1)
	g.StartTurn(0)
	if got := g.CurrentKeyCost(0); got != KeyCost-2 {
		t.Errorf("controller key cost = %d, want %d on their next turn", got, KeyCost-2)
	}
	g.EndPlayPhase(0)
	if got := g.CurrentKeyCost(0); got != KeyCost {
		t.Errorf("controller key cost = %d, want %d after the window", got, KeyCost)
	}
}

// TestLowerKeyCostSumsWithRaiseAndFloorsAtZero checks a lower and a raise on the
// same player sum, and a heavier lower never drives the key cost below 0.
func TestLowerKeyCostSumsWithRaiseAndFloorsAtZero(t *testing.T) {
	g := started(t)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	LowerKeyCost{Player: Controller, Amount: 2, Duration: RemainderOfPlayerTurn}.Resolve(ctx)
	RaiseKeyCost{Player: Controller, Amount: 3, Duration: RemainderOfPlayerTurn}.Resolve(ctx)
	if got := g.CurrentKeyCost(0); got != KeyCost+1 {
		t.Errorf("key cost = %d, want %d (a -2 and a +3 sum)", got, KeyCost+1)
	}

	// A drop larger than the base cost floors the key cost at 0, never negative.
	LowerKeyCost{Player: Controller, Amount: 20, Duration: RemainderOfPlayerTurn}.Resolve(ctx)
	if got := g.CurrentKeyCost(0); got != 0 {
		t.Errorf("key cost = %d, want 0 (floored, never negative)", got)
	}
}

func TestConditionalElse(t *testing.T) {
	effect := Conditional{
		Cond: PoolAember{Player: Opponent, Is: Exactly},
		Then: ForgeKey{Extra: 2},
		Else: ForgeKey{Extra: 6},
	}
	want := "if your opponent has no Æmber, forge a key at +2 Æmber current cost -> purge {self}. " +
		"Otherwise, forge a key at +6 Æmber current cost -> purge {self}"
	if got := effect.Text(); got != want {
		t.Errorf("text = %q, want %q", got, want)
	}
	if err := effect.validate(); err != nil {
		t.Errorf("validate = %v, want nil", err)
	}
	if err := (Conditional{
		Cond: PoolAember{Player: Opponent, Is: Exactly},
		Then: ForgeKey{},
		Else: ForgeKey{FreeOfCost: true, Extra: 1},
	}).validate(); err == nil {
		t.Error("an invalid Else should not validate")
	}

	g := started(t)
	g.State.Aember[0] = 100
	g.State.Aember[1] = 1
	effect.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if got := g.State.Aember[0]; got != 100-(KeyCost+6) {
		t.Errorf("Æmber = %d, want the Else branch's +6 cost paid", got)
	}
}

func TestCardsInHandAnyHouse(t *testing.T) {
	g := started(t)
	g.AddToHand(NewCard("A", Brobnar, Tactic, Common), 0)
	g.AddToHand(NewCard("B", Logos, Tactic, Common), 0)
	g.AddToHand(NewCard("C", Dis, Tactic, Common), 1)

	mine := CardsInHand{Player: Controller, House: AnyHouse}
	theirs := CardsInHand{Player: Opponent, House: AnyHouse}
	ctx := &EffectContext{Resolver: g, Controller: 0}

	if got := mine.Value(ctx); got != 2 {
		t.Errorf("controller hand = %d, want 2", got)
	}
	if got := theirs.Value(ctx); got != 1 {
		t.Errorf("opponent hand = %d, want 1", got)
	}
	if got := mine.CountText(); got != "card in your hand" {
		t.Errorf("controller text = %q", got)
	}
	if got := theirs.CountText(); got != "card in your opponent's hand" {
		t.Errorf("opponent text = %q", got)
	}
}
