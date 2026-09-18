package engine

import (
	"strings"
	"testing"
)

func TestForgeKeyEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddArtifact(NewCard("Forge", Dis, Artifact, Common), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	e := ForgeKey{}
	if e.Text() != "forge a key at current cost -> purge {self}" {
		t.Errorf("text = %q", e.Text())
	}

	// Not enough Æmber: no key is forged and the source stays in play.
	e.Resolve(ctx)
	if g.Keys(0) != 0 {
		t.Errorf("keys = %d, want 0 (could not afford)", g.Keys(0))
	}
	if !g.inPlay(src) {
		t.Error("source should stay in play when no key is forged")
	}

	// Enough Æmber: one key is forged, its cost paid, and the source purged.
	g.State.Aember[0] = KeyCost + 2
	e.Resolve(ctx)
	if g.Keys(0) != 1 {
		t.Errorf("keys = %d, want 1", g.Keys(0))
	}
	if g.Aember(0) != 2 {
		t.Errorf("aember = %d, want 2 (paid the key cost)", g.Aember(0))
	}
	if g.inPlay(src) {
		t.Error("source should be purged after forging")
	}
}

func TestForgeKeyFreeEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	forger := g.AddToBattleline(exGiant(), 0)
	foe := g.AddToBattleline(testCreature("foe", 4), 1)
	ctx := &EffectContext{Resolver: g, Source: forger, Controller: 0}
	e := ForgeKey{FreeOfCost: true}
	if e.Text() != "forge a key at no cost -> purge {self}" {
		t.Errorf("text = %q", e.Text())
	}

	e.Resolve(ctx)

	if g.Keys(0) != 1 {
		t.Errorf("keys = %d, want 1", g.Keys(0))
	}
	if g.Aember(0) != 0 {
		t.Errorf("aember = %d, want 0 (no cost paid)", g.Aember(0))
	}
	if g.Damage(foe) != 2 {
		t.Errorf("after-forge ability damage = %d, want 2", g.Damage(foe))
	}
	if g.inPlay(forger) {
		t.Error("source should be purged after forging for free")
	}
}

func TestForOpponentNextTurnForgeGivesAember(t *testing.T) {
	e := ForOpponentNextTurn{
		On: EventForgeKey,
		Do: GiveAember{All: true},
	}
	if e.Text() != "during your opponent's next turn, after forging a key, your opponent gives you all their Æmber" {
		t.Errorf("text = %q", e.Text())
	}

	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	e.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if g.State.LastingCount != 1 || g.State.Lasting[0].Controller != 1 ||
		g.State.Lasting[0].On != EventForgeKey ||
		g.State.Lasting[0].Once {
		t.Fatal("opponent's next-turn key-forge reaction was not armed as a durable reaction")
	}

	g.State.Aember[1] = KeyCost + 4
	g.EndPlayPhase(0) // the opponent-owned reaction survives the controller's turn end
	g.StartTurn(1)

	if g.Keys(1) != 1 {
		t.Errorf("opponent keys = %d, want 1", g.Keys(1))
	}
	if g.Aember(1) != 0 {
		t.Errorf("opponent aember = %d, want 0", g.Aember(1))
	}
	if g.Aember(0) != 4 {
		t.Errorf("controller aember = %d, want 4", g.Aember(0))
	}
	if g.State.LastingCount != 1 {
		t.Error("reaction should persist so further forges this turn also transfer")
	}

	g.EndPlayPhase(1)
	if g.State.LastingCount != 0 {
		t.Error("reaction should clear at the end of the opponent's next turn")
	}
}

func TestForOpponentNextTurnForgeGivesAemberEveryForge(t *testing.T) {
	// A key cheat can forge more than one key in a turn; the opponent must give
	// their remaining Æmber each time.
	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	ForOpponentNextTurn{
		On: EventForgeKey,
		Do: GiveAember{All: true},
	}.
		Resolve(
			&EffectContext{Resolver: g, Controller: 0},
		)
	g.EndPlayPhase(0)
	g.StartTurn(1) // the start-of-turn forge is skipped (no Æmber), reaction still armed

	g.State.Aember[1] = 5
	ForgeKey{FreeOfCost: true}.Resolve(&EffectContext{Resolver: g, Controller: 1})
	if g.Aember(1) != 0 || g.Aember(0) != 5 {
		t.Fatalf("first forge: opponent=%d controller=%d, want 0/5", g.Aember(1), g.Aember(0))
	}

	g.State.Aember[1] = 3
	ForgeKey{FreeOfCost: true}.Resolve(&EffectContext{Resolver: g, Controller: 1})
	if g.Aember(1) != 0 || g.Aember(0) != 8 {
		t.Fatalf("second forge: opponent=%d controller=%d, want 0/8", g.Aember(1), g.Aember(0))
	}
}

func TestForOpponentNextTurnForgeGivesAemberExpiresIfNoForge(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	ForOpponentNextTurn{
		On: EventForgeKey,
		Do: GiveAember{All: true},
	}.
		Resolve(
			&EffectContext{Resolver: g, Controller: 0},
		)

	g.State.Aember[1] = KeyCost - 1
	g.EndPlayPhase(0)
	g.StartTurn(1)
	if g.Aember(1) != KeyCost-1 {
		t.Errorf("opponent aember = %d, want %d", g.Aember(1), KeyCost-1)
	}

	g.EndPlayPhase(1)
	if g.State.LastingCount != 0 {
		t.Error("reaction should expire at the end of the opponent's next turn")
	}

	g.State.Aember[1] = KeyCost + 2
	g.StartTurn(1)
	if g.Aember(1) != 2 {
		t.Errorf("later opponent aember = %d, want 2", g.Aember(1))
	}
	if g.Aember(0) != 0 {
		t.Errorf("controller aember = %d, want 0 (expired transfer)", g.Aember(0))
	}
}

func TestForOpponentNextTurnForgeGivesAemberAppliesToCardForge(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	ForOpponentNextTurn{
		On: EventForgeKey,
		Do: GiveAember{All: true},
	}.
		Resolve(
			&EffectContext{Resolver: g, Controller: 0},
		)
	g.EndPlayPhase(0)
	g.StartTurn(1)

	g.State.Aember[1] = 3
	ForgeKey{FreeOfCost: true}.Resolve(&EffectContext{Resolver: g, Controller: 1})

	if g.Aember(1) != 0 {
		t.Errorf("opponent aember = %d, want 0", g.Aember(1))
	}
	if g.Aember(0) != 3 {
		t.Errorf("controller aember = %d, want 3", g.Aember(0))
	}
}

func TestForgeKeyLastingText(t *testing.T) {
	if got := EventForgeKey.clause(); got != "after forging a key" {
		t.Errorf("clause = %q", got)
	}
	if got := actGiveRemainingAember.describe(); got != "give remaining Æmber" {
		t.Errorf("describe = %q", got)
	}
}

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
	if g.Keys(0) != 0 {
		t.Fatalf("keys = %d, want 0 — the surcharge is unaffordable", g.Keys(0))
	}

	ForgeKey{Extra: 2}.Resolve(ctx)
	if g.Keys(0) != 1 {
		t.Errorf("keys = %d, want 1", g.Keys(0))
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

	if g.Keys(0) != 1 {
		t.Errorf("keys = %d, want 1", g.Keys(0))
	}
	if g.State.Aember[0] != 0 {
		t.Errorf("Æmber = %d, want 0 — the cost floor is the key cost", g.State.Aember[0])
	}
}

func TestForgeKeyDiscountFloorsAndPurges(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddArtifact(NewCard("Forge", Dis, Artifact, Common), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	// A discount larger than the current key cost floors the whole cost at 0, so the
	// forge lands with an empty pool, and the successful forge purges the source.
	for i := 0; i < 8; i++ {
		g.AddToHand(NewCard("Filler", Brobnar, Tactic, Common), 0)
	}
	ForgeKey{
		Discount:  true,
		ReducedBy: CardsInHand{Player: Controller, House: AnyHouse},
	}.Resolve(ctx)

	if g.Keys(0) != 1 {
		t.Errorf("keys = %d, want 1 — the discount floors the cost at 0", g.Keys(0))
	}
	if g.inPlay(src) {
		t.Error("a successful forge should purge the source")
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

// TestRaiseKeyCostEachPlayer checks a raise on both players speaks in the
// controller's frame (like LowerKeyCost) and arms the surcharge on each player.
func TestRaiseKeyCostEachPlayer(t *testing.T) {
	each := RaiseKeyCost{Player: EachPlayer, Amount: 2, Duration: EndOfPlayerNextTurn}
	if got := each.Text(); got != "each player's keys cost +2 Æmber until the end of your next turn" {
		t.Errorf("text = %q", got)
	}
	if err := each.validate(); err != nil {
		t.Errorf("validate = %v, want nil", err)
	}

	g := started(t)
	source := g.AddToBattleline(NewCard("Tax", Dis, Creature, Common, WithPower(1)), 0)
	each.Resolve(&EffectContext{Resolver: g, Controller: 0, Source: source})
	if got := g.CurrentKeyCost(0); got != KeyCost+2 {
		t.Errorf("controller key cost = %d, want %d right away", got, KeyCost+2)
	}
	if got := g.CurrentKeyCost(1); got != KeyCost+2 {
		t.Errorf("opponent key cost = %d, want %d right away", got, KeyCost+2)
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

func TestRaiseKeyCostPerCreature(t *testing.T) {
	effect := RaiseKeyCost{
		Player:   Opponent,
		Amount:   1,
		House:    namedHouse(Dis),
		Duration: OpponentNextTurn,
	}
	if got := effect.Text(); got !=
		"keys cost +1 Æmber for each Dis creature in play during your opponent's next turn" {
		t.Errorf("text = %q", got)
	}
	if got := (RaiseKeyCost{
		Player:   Controller,
		Amount:   2,
		House:    namedHouse(Dis),
		Duration: OpponentNextTurn,
	}).Text(); got !=
		"keys cost +2 Æmber for each Dis creature in play during your next turn" {
		t.Errorf("controller text = %q", got)
	}
	if err := (RaiseKeyCost{
		Amount:   1,
		House:    namedHouse(Dis),
		Duration: OpponentNextTurn,
	}).validate(); err == nil {
		t.Error("an unset player should not validate")
	}
	if err := (RaiseKeyCost{
		Player:   Opponent,
		House:    namedHouse(Dis),
		Duration: OpponentNextTurn,
	}).validate(); err == nil {
		t.Error("a zero raise should not validate")
	}
	if err := (RaiseKeyCost{
		Player:   Opponent,
		Amount:   1,
		House:    HouseMatcher{Kind: MatchChosenHouse},
		Duration: OpponentNextTurn,
	}).validate(); err == nil {
		t.Error("a context-dependent house matcher should not validate")
	}
	if err := (RaiseKeyCost{
		Player:   Opponent,
		Amount:   1,
		House:    HouseMatcher{Kind: MatchNamedHouse},
		Duration: OpponentNextTurn,
	}).validate(); err == nil {
		t.Error("a named house matcher with no house should not validate")
	}
	if err := (RaiseKeyCost{
		Player: Opponent,
		Amount: 1,
		House:  namedHouse(Dis),
	}).validate(); err == nil {
		t.Error("an unset duration should not validate")
	}
	if err := (RaiseKeyCost{
		Player:   Opponent,
		Amount:   1,
		House:    namedHouse(Dis),
		Duration: RemainderOfPlayerTurn,
	}).validate(); err == nil {
		t.Error("an unsupported duration should not validate")
	}
	if err := effect.validate(); err != nil {
		t.Errorf("validate = %v, want nil", err)
	}
}

// TestRaiseKeyCostPerCreatureCountsLive proves the surcharge is dormant until
// the taxed player's turn, then reads the number of Dis creatures in play live at
// each forge, and lifts when that turn ends.
func TestRaiseKeyCostPerCreatureCountsLive(t *testing.T) {
	g := started(t)
	source := g.AddToBattleline(NewCard("Nightmare", Dis, Creature, Common, WithPower(1)), 0)
	g.AddToBattleline(NewCard("dis-a", Dis, Creature, Common, WithPower(1)), 0)
	g.AddToBattleline(NewCard("dis-b", Dis, Creature, Common, WithPower(1)), 1)
	g.AddToBattleline(NewCard("logos", Logos, Creature, Common, WithPower(1)), 1) // not counted

	RaiseKeyCost{
		Player:   Opponent,
		Amount:   1,
		House:    namedHouse(Dis),
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

// TestRaiseKeyCostPerCreatureStacks proves two surcharges of the same house
// sum their per-creature amount, while a different house replaces the surcharge.
func TestRaiseKeyCostPerCreatureStacks(t *testing.T) {
	g := started(t)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	RaiseKeyCost{
		Player:   Opponent,
		Amount:   1,
		House:    namedHouse(Dis),
		Duration: OpponentNextTurn,
	}.Resolve(ctx)
	RaiseKeyCost{
		Player:   Opponent,
		Amount:   2,
		House:    namedHouse(Dis),
		Duration: OpponentNextTurn,
	}.Resolve(ctx)
	if got := g.State.KeyCostPerHouseNext[1].Value; got.Per != 3 || got.House != namedHouse(Dis) {
		t.Errorf("armed surcharge = %+v, want {Dis 3}", got)
	}
	RaiseKeyCost{
		Player:   Opponent,
		Amount:   5,
		House:    namedHouse(Logos),
		Duration: OpponentNextTurn,
	}.Resolve(ctx)
	if got := g.State.KeyCostPerHouseNext[1].Value; got.Per != 5 || got.House != namedHouse(Logos) {
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
	}).Text(); got != "keys cost -3 Æmber during your opponent's next turn" {
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

func TestSkipForgePhase(t *testing.T) {
	if got := (SkipForgePhase{Player: Opponent}).Text(); got != `your opponent skips the "forge a key" phase during their next turn` {
		t.Errorf("opponent text = %q", got)
	}
	if got := (SkipForgePhase{Player: Controller}).Text(); got != `you skip the "forge a key" phase during your next turn` {
		t.Errorf("self text = %q", got)
	}
	if (SkipForgePhase{}).validate() == nil {
		t.Error("unset player should be invalid")
	}
	if (SkipForgePhase{Player: Opponent}).validate() != nil {
		t.Error("a set player should be valid")
	}

	g := NewGame("A", "B", 1)
	g.StartTurn(0)
	if err := g.ChooseHouse(0, Brobnar); err != nil {
		t.Fatal(err)
	}
	g.State.Aember[1] = 6 // enough for player 1 to forge on their turn
	SkipForgePhase{Player: Opponent}.Resolve(&EffectContext{Resolver: g, Controller: 0})
	if !g.State.SkipForgeNext[1].Value {
		t.Fatal("the skip should arm the opponent's next turn")
	}
	g.EndPlayPhase(0)

	keysBefore := g.Keys(1)
	g.StartTurn(1)
	if g.State.SkipForgeNext[1].Value {
		t.Error("the skip should be consumed when the turn begins")
	}
	if g.Keys(1) != keysBefore {
		t.Error("the opponent should not have forged a key on the skipped turn")
	}
}

func TestSkipsForgeConstant(t *testing.T) {
	g := NewGame("A", "B", 1)
	if g.skipsForge(0) {
		t.Error("a player with nothing in play should not skip forging")
	}
	def := NewCard("The Sting", Shadows, Artifact, Rare,
		WithRestrictions(Restrictions{SkipForge: true}))
	sting := g.AddArtifact(def, 0)
	if !g.skipsForge(0) {
		t.Error("a controller of a SkipForge card should skip forging")
	}
	if g.skipsForge(1) {
		t.Error("the opponent should be unaffected")
	}
	if !strings.Contains(RenderCardRules(&def), `You skip your "forge a key" phase.`) {
		t.Error("card rules should render the skip-forge line")
	}

	g.State.Aember[0] = 6
	keysBefore := g.Keys(0)
	g.forgePhase(0)
	if g.Keys(0) != keysBefore {
		t.Error("a player skipping their forge phase should not forge a key")
	}
	if g.State.Aember[0] != 6 {
		t.Error("skipping the forge phase should not spend any Æmber")
	}
	_ = sting
}

func TestForgeAemberGain(t *testing.T) {
	if got := gainsForgeAemberText(&CardDefinition{}); got != "" {
		t.Errorf("gainsForgeAemberText on a plain card = %q, want empty", got)
	}
	def := NewCard("The Sting", Shadows, Artifact, Rare, WithGainsForgeAember())
	if !def.GainsForgeAember {
		t.Error("WithGainsForgeAember should set the flag")
	}
	if got := gainsForgeAemberText(
		&def,
	); got != "You gain all Æmber your opponent spends when forging a key." {
		t.Errorf("gainsForgeAemberText = %q", got)
	}
	if !strings.Contains(
		RenderCardRules(&def), "You gain all Æmber your opponent spends when forging a key.",
	) {
		t.Error("card rules should render the forge-aember-gain line")
	}

	g := NewGame("A", "B", 1)
	if _, ok := g.forgeAemberGainer(0); ok {
		t.Error("no gainer should be found with nothing in play")
	}
	sting := g.AddArtifact(def, 1)
	gainer, ok := g.forgeAemberGainer(0)
	if !ok || gainer != sting {
		t.Errorf("forgeAemberGainer = %d, %v, want %d, true", gainer, ok, sting)
	}
	if _, ok := g.forgeAemberGainer(1); ok {
		t.Error("a player should not gain their own forge spending")
	}

	g.State.Aember[0] = 6
	g.forgeKey(0)
	if g.Keys(0) != 1 {
		t.Fatal("the payer should still forge their key")
	}
	if g.State.Aember[0] != 0 {
		t.Errorf("payer pool = %d, want 0", g.State.Aember[0])
	}
	if g.State.Aember[1] != KeyCost {
		t.Errorf("gainer pool = %d, want %d", g.State.Aember[1], KeyCost)
	}
}
