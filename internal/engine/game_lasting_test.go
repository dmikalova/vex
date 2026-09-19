package engine

import "testing"

// TestAddLastingRejectsFilteredReplacement pins that a lasting record on a
// replacement event may not narrow which subject it applies to. lastingReplacement
// matches on controller and event alone, so a filtered record would quietly replace
// outcomes it was never meant to reach; AddLasting refuses it at the door instead.
// A replacement that genuinely needs to discriminate needs lastingReplacement
// taught the filters first.
func TestAddLastingRejectsFilteredReplacement(t *testing.T) {
	for name, le := range map[string]LastingEffect{
		"subject": {On: EventReapAember, Do: actSteal, HasSubject: true, Subject: 3},
		"except":  {On: EventReapAember, Do: actSteal, HasExcept: true, Except: 3},
		"type":    {On: EventReapAember, Do: actSteal, Type: Creature},
		"trait":   {On: EventReapAember, Do: actSteal, Trait: Beast},
		"house": {
			On: EventReapAember,
			Do: actSteal,
			House: HouseMatcher{
				Kind:  MatchNamedHouse,
				House: Mars,
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Errorf("AddLasting(%s-filtered replacement) did not panic", name)
				}
			}()
			g := NewGame("A", "B", 1)
			g.AddLasting(le)
		})
	}

	// The unfiltered record is what Dimension Door installs, and it still registers.
	g := NewGame("A", "B", 1)
	g.AddLasting(LastingEffect{
		On: EventReapAember,
		Do: actSteal,
	})
	if _, ok := g.lastingReplacement(0, EventReapAember); !ok {
		t.Error("unfiltered replacement was not registered")
	}
}

// TestLastingActionDescriptions pins the short labels the ordering prompt shows
// when several reactions fire at once (ADR 0013).
func TestLastingActionDescriptions(t *testing.T) {
	for act, want := range map[lastingAction]string{
		actGainAember:  "gain Æmber",
		actDealDamage:  "deal damage",
		actReadyPlayed: "ready the creature",
		actDraw:        "draw a card",
		actLoseAember:  "opponent loses Æmber",
		actSteal:       "steal Æmber",
		actExalt:       "exalt the creature",
		actStun:        "stun the creature",
	} {
		if got := act.describe(); got != want {
			t.Errorf("describe(%d) = %q, want %q", act, got, want)
		}
	}
}

func TestLastingReactionsResolveOnPlay(t *testing.T) {
	g := started(t) // player 0 active, Brobnar
	foe := g.AddToBattleline(testCreature("foe", 5), 1)
	// Two reactions fire on the same event; both resolve when a creature is played.
	g.AddLasting(
		LastingEffect{
			On:         EventCreaturePlayed,
			Do:         actGainAember,
			Controller: 0,
			Amount:     1,
		},
	)
	g.AddLasting(
		LastingEffect{
			On:         EventCreaturePlayed,
			Do:         actDealDamage,
			Controller: 0,
			Amount:     2,
		},
	)

	g.AddToHand(testCreature("minion", 4), 0)
	before := g.Aember(0)
	if _, err := g.PlayCreature(0, handIdx(g, 0, "minion"), false); err != nil {
		t.Fatalf("PlayCreature: %v", err)
	}
	if got := g.Aember(0) - before; got != 1 {
		t.Errorf("Æmber gained = %d, want 1", got)
	}
	if g.Damage(foe) != 2 {
		t.Errorf("foe damage = %d, want 2", g.Damage(foe))
	}
}

func TestLastingOnceExpiresAtEndOfTurn(t *testing.T) {
	g := started(t)
	g.AddLasting(
		LastingEffect{
			On:         EventCreaturePlayed,
			Do:         actReadyPlayed,
			Controller: 0,
			Amount:     0,
			House:      namedHouse(Mars),
			Type:       Creature,
			Once:       true,
		},
	)
	if g.State.LastingCount != 1 {
		t.Fatalf("setup: lasting count = %d, want 1", g.State.LastingCount)
	}

	g.EndPlayPhase(0)

	if g.State.LastingCount != 0 {
		t.Errorf(
			"a one-shot reaction that never fired should be cleared at end of turn, count = %d",
			g.State.LastingCount,
		)
	}
}

func TestClearLastingKeepsOtherPlayer(t *testing.T) {
	g := started(t)
	g.AddLasting(
		LastingEffect{
			On:         EventCreaturePlayed,
			Do:         actGainAember,
			Controller: 0,
			Amount:     1,
		},
	)
	g.AddLasting(
		LastingEffect{
			On:         EventCreaturePlayed,
			Do:         actGainAember,
			Controller: 1,
			Amount:     1,
		},
	) // the opponent's reaction

	g.clearLasting(0)

	if g.State.LastingCount != 1 {
		t.Fatalf("lasting count = %d, want 1 (the opponent's kept)", g.State.LastingCount)
	}
	if g.State.Lasting[0].Controller != 1 {
		t.Errorf("kept controller = %d, want 1", g.State.Lasting[0].Controller)
	}
}

func TestAddLastingCap(t *testing.T) {
	g := started(t)
	for range maxLasting + 3 {
		g.AddLasting(
			LastingEffect{
				On:         EventCreaturePlayed,
				Do:         actGainAember,
				Controller: 0,
				Amount:     1,
			},
		)
	}
	if int(g.State.LastingCount) != maxLasting {
		t.Errorf("lasting count = %d, want %d (capped)", g.State.LastingCount, maxLasting)
	}
}

func TestLastingOnceReadiesMatchingHouseAndSelfRemoves(t *testing.T) {
	g := NewGame("A", "B", 1)
	mars := g.AddToBattleline(NewCard("m", Mars, Creature, Common, WithPower(2)), 0)
	sanc := g.AddToBattleline(NewCard("s", Sanctum, Creature, Common, WithPower(2)), 0)
	g.State.Cards[mars].Exhausted = true
	g.State.Cards[sanc].Exhausted = true

	g.AddLasting(
		LastingEffect{
			On:         EventCreaturePlayed,
			Do:         actReadyPlayed,
			Controller: 0,
			Amount:     0,
			House:      namedHouse(Mars),
			Type:       Creature,
			Once:       true,
		},
	)

	// A non-Mars subject is filtered out: not readied, and the entry stays armed.
	g.resolveLastingWindow(EventCreaturePlayed, 0, sanc)
	if !g.State.Cards[sanc].Exhausted {
		t.Error("a non-Mars creature should not be readied")
	}
	if g.State.LastingCount != 1 {
		t.Errorf("entry should remain after a filtered subject, count = %d", g.State.LastingCount)
	}

	// A Mars subject is readied, and the one-shot entry removes itself.
	g.resolveLastingWindow(EventCreaturePlayed, 0, mars)
	if g.State.Cards[mars].Exhausted {
		t.Error("the next Mars creature should enter ready")
	}
	if g.State.LastingCount != 0 {
		t.Errorf("the one-shot entry should be removed, count = %d", g.State.LastingCount)
	}
}

// TestEnteringPlayIsOneEventWhicheverWay pins that a card put into play fires the
// EventCardEntersPlay lasting reactions, exactly as a card played from hand does.
// The two arrivals once disagreed: putIntoPlay fired the AfterCreatureEnters
// trigger but never the lasting event of the same name, because each play path
// appended that event by hand and putIntoPlay was not a play path. Both now gather
// it through abilityWindow.addEntersPlay, so a new arrival cannot miss it.
func TestEnteringPlayIsOneEventWhicheverWay(t *testing.T) {
	for _, tc := range []struct {
		name string
		typ  CardType
	}{{"creature", Creature}, {"artifact", Artifact}} {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGame("A", "B", 1)
			id := g.Register(NewCard("arrival", Mars, tc.typ, Common, WithPower(3)), 0)
			g.State.Deck[0].add(id)

			g.AddLasting(LastingEffect{
				On:         EventCardEntersPlay,
				Do:         actReadyPlayed,
				Controller: 0,
				House:      anyHouse,
				Type:       AnyType,
				Once:       true,
			})

			g.putIntoPlay(id, 0)

			if g.State.LastingCount != 0 {
				t.Errorf(
					"a card put into play should consume the entry, count = %d",
					g.State.LastingCount,
				)
			}
			if g.State.Cards[id].Exhausted {
				t.Error("the lasting reaction should have readied the arrival")
			}
		})
	}
}

// TestLastingOnceFiltersByCardType covers the type filter a "next creature or
// artifact" entry carries (Soft Landing): an upgrade play is passed over, while a
// creature and an artifact both satisfy it.
func TestLastingOnceFiltersByCardType(t *testing.T) {
	g := NewGame("A", "B", 1)
	up := g.Register(NewCard("u", Mars, Upgrade, Common), 0)
	artifact := g.Register(NewCard("a", Mars, Artifact, Common), 0)
	g.State.Artifacts[0].add(artifact)
	g.State.Cards[artifact].Exhausted = true

	g.AddLasting(
		LastingEffect{
			On:         EventCardEntersPlay,
			Do:         actReadyPlayed,
			Controller: 0,
			Amount:     0,
			House:      anyHouse,
			Type:       AnyType,
			Once:       true,
		},
	)

	g.resolveLastingWindow(EventCardEntersPlay, 0, up)
	if g.State.LastingCount != 1 {
		t.Errorf("an upgrade should not satisfy the entry, count = %d", g.State.LastingCount)
	}

	g.resolveLastingWindow(EventCardEntersPlay, 0, artifact)
	if g.State.Cards[artifact].Exhausted {
		t.Error("the artifact should enter ready")
	}
	if g.State.LastingCount != 0 {
		t.Errorf("the one-shot entry should be removed, count = %d", g.State.LastingCount)
	}
}

func TestLastingOnceOrdersWithPersistentReaction(t *testing.T) {
	g := NewGame("A", "B", 1)
	mars := g.AddToBattleline(NewCard("m", Mars, Creature, Common, WithPower(2)), 0)
	g.State.Cards[mars].Exhausted = true

	// The one-shot sits between two persistent reactions, so removing it scans past
	// the first and shifts the last down.
	g.AddLasting(
		LastingEffect{
			On:         EventCreaturePlayed,
			Do:         actGainAember,
			Controller: 0,
			Amount:     1,
		},
	)
	g.AddLasting(
		LastingEffect{
			On:         EventCreaturePlayed,
			Do:         actReadyPlayed,
			Controller: 0,
			Amount:     0,
			House:      namedHouse(Mars),
			Type:       Creature,
			Once:       true,
		},
	)
	g.AddLasting(
		LastingEffect{
			On:         EventCreaturePlayed,
			Do:         actGainAember,
			Controller: 0,
			Amount:     1,
		},
	)

	g.resolveLastingWindow(EventCreaturePlayed, 0, mars) // three reactions fire; ordering path runs

	if g.Aember(0) != 2 {
		t.Errorf("aember = %d, want 2", g.Aember(0))
	}
	if g.State.Cards[mars].Exhausted {
		t.Error("the Mars creature should be readied")
	}
	if g.State.LastingCount != 2 {
		t.Errorf("count = %d, want 2 (both persistent reactions remain)", g.State.LastingCount)
	}
}

// recordAember is a test effect that records its controller's Æmber pool at the
// moment it resolves, so a test can tell whether another reaction in the same
// window resolved before or after it.
type recordAember struct{ got *int }

func (recordAember) Text() string { return "record Æmber" }
func (e recordAember) Resolve(ctx *EffectContext) {
	*e.got = ctx.Resolver.Aember(ctx.Controller)
}

// reactionOrderRecorder is a ReactionChooser test chooser: it records the first
// window it is handed and picks the last reaction each call, reversing the window,
// so a test can prove a duration reaction was ordered against a card ability in one
// window. It returns an out-of-range index when invalid is set, to exercise the
// gathered-order fallback.
type reactionOrderRecorder struct {
	got     []OrderableReaction
	invalid bool
	asked   bool
}

func (reactionOrderRecorder) ChooseCreature(_, _ string, cands []LocalID) (LocalID, bool) {
	if len(cands) == 0 {
		return 0, false
	}
	return cands[0], true
}

func (r *reactionOrderRecorder) ChooseReaction(_ string, rs []OrderableReaction) int {
	if !r.asked {
		r.got = rs
		r.asked = true
	}
	if r.invalid {
		return -1 // out of range: forces the gathered-order fallback
	}
	return len(rs) - 1 // pick the last, reversing the window
}

// TestReapWindowInterleavesLastingReaction proves a duration reaction (Crystal
// Hive's "after a creature reaps: gain Æmber") is ordered in the same window as a
// creature's own "Reap:" ability, not a trailing one: a ReactionChooser that
// reverses the window resolves the duration reaction before the card ability, so
// the card ability sees the duration Æmber already gained.
func TestReapWindowInterleavesLastingReaction(t *testing.T) {
	g := started(t)
	var recorded int
	reaper := g.AddToBattleline(
		testCreature("reaper", 3, WithAbility(TriggerAfterReap, recordAember{&recorded})),
		0,
	)
	g.AddLasting(
		LastingEffect{
			On:         EventReap,
			Do:         actGainAember,
			Controller: 0,
			Amount:     1,
			Once:       true,
		},
	)
	rec := &reactionOrderRecorder{}
	g.SetChooser(0, rec)

	before := g.Aember(0)
	g.reapWith(reaper)

	// The reaper gains 1 for reaping, then the reversed window resolves the duration
	// gain (before + 2) before the card ability records the pool.
	if recorded != before+2 {
		t.Errorf(
			"recorded pool = %d, want %d (duration reaction resolved first)",
			recorded,
			before+2,
		)
	}
	if len(rec.got) != 2 {
		t.Fatalf(
			"ordered %d reactions, want 2 (the card ability and the duration reaction)",
			len(rec.got),
		)
	}
	if !rec.got[0].HasCard || rec.got[0].Card != reaper {
		t.Errorf("first reaction = %+v, want the reaper's card ability", rec.got[0])
	}
	if rec.got[1].HasCard || rec.got[1].Label != actGainAember.describe() {
		t.Errorf("second reaction = %+v, want the duration reaction", rec.got[1])
	}
	if g.State.LastingCount != 0 {
		t.Errorf("lasting count = %d, want 0 (the one-shot removed itself)", g.State.LastingCount)
	}
}

// TestReapWindowFallsBackOnBadOrder checks that a ReactionChooser returning an
// out-of-range index is ignored: the window resolves in its gathered order (card
// ability before the duration reaction), so the card ability records the pool
// before the duration gain lands.
func TestReapWindowFallsBackOnBadOrder(t *testing.T) {
	g := started(t)
	var recorded int
	reaper := g.AddToBattleline(
		testCreature("reaper", 3, WithAbility(TriggerAfterReap, recordAember{&recorded})),
		0,
	)
	g.AddLasting(
		LastingEffect{
			On:         EventReap,
			Do:         actGainAember,
			Controller: 0,
			Amount:     1,
		},
	)
	g.SetChooser(0, &reactionOrderRecorder{invalid: true})

	before := g.Aember(0)
	g.reapWith(reaper)

	if recorded != before+1 {
		t.Errorf(
			"recorded pool = %d, want %d (default order: card ability first)",
			recorded,
			before+1,
		)
	}
}

// TestCaptureChosenReactionOnCardPlayed covers Commandeer: after you play another
// card, a friendly creature the active player chooses captures 1 Æmber.
func TestCaptureChosenReactionOnCardPlayed(t *testing.T) {
	if actCaptureChosen.describe() != "a friendly creature captures Æmber" {
		t.Errorf("describe = %q", actCaptureChosen.describe())
	}

	g := started(t)
	g.State.Aember[1] = 3
	src := g.AddToDiscard(NewCard("Commandeer", Sanctum, Tactic, Common), 0)

	ForRemainderOfTurn{
		On: EventCardPlayed,
		Do: CaptureAember{
			Amount: 1,
			Target: Target{Kind: TargetChosenFriendlyCreature},
			Source: Opponent,
		},
	}.Resolve(&EffectContext{
		Resolver:   g,
		Controller: 0,
		Source:     src,
	})

	capper := g.AddToBattleline(testCreature("capper", 3), 0)
	played := g.AddToBattleline(testCreature("played", 2), 0)
	g.SetChooser(0, idChooser{id: capper})

	g.resolveLastingWindow(EventCardPlayed, 0, played)

	if g.AmberOn(capper) != 1 {
		t.Errorf("captured Æmber on chosen creature = %d, want 1", g.AmberOn(capper))
	}
	if g.AmberOn(played) != 0 {
		t.Errorf("the played card should not capture; got %d", g.AmberOn(played))
	}
	if g.State.Aember[1] != 2 {
		t.Errorf("opponent pool = %d, want 2 after a capture", g.State.Aember[1])
	}
}

// TestCaptureChosenReactionFizzlesWithoutCreature checks the reaction does nothing
// when the active player controls no creature to capture with.
func TestCaptureChosenReactionFizzlesWithoutCreature(t *testing.T) {
	g := started(t)
	g.State.Aember[1] = 3
	src := g.AddToDiscard(NewCard("Commandeer", Sanctum, Tactic, Common), 0)

	ForRemainderOfTurn{
		On: EventCardPlayed,
		Do: CaptureAember{
			Amount: 1,
			Target: Target{Kind: TargetChosenFriendlyCreature},
			Source: Opponent,
		},
	}.Resolve(&EffectContext{
		Resolver:   g,
		Controller: 0,
		Source:     src,
	})

	played := g.AddToDiscard(NewCard("Another", Sanctum, Tactic, Common), 0)
	g.resolveLastingWindow(EventCardPlayed, 0, played)

	if g.State.Aember[1] != 3 {
		t.Errorf("opponent pool = %d, want 3 (no capture)", g.State.Aember[1])
	}
}

func TestCaptureReactionOnFight(t *testing.T) {
	if actCapture.describe() != "capture Æmber" {
		t.Errorf("describe = %q", actCapture.describe())
	}

	g := started(t)
	g.State.Aember[1] = 3

	// Register the lasting the way Take Hostages authors it, to exercise the
	// CaptureAember -> actCapture mapping and validation.
	ForRemainderOfTurn{
		On: EventFight,
		Do: CaptureAember{
			Amount: 1,
			Target: Target{Kind: TargetTriggeringCreature},
			Source: Opponent,
		},
	}.
		Resolve(
			&EffectContext{
				Resolver:   g,
				Controller: 0,
			},
		)

	att := g.AddToBattleline(NewCard("att", Brobnar, Creature, Common, WithPower(4)), 0)
	def := g.AddToBattleline(testCreature("def", 2), 1)
	if err := g.Fight(0, att, def); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	if g.AmberOn(att) != 1 {
		t.Errorf("captured Æmber on fighter = %d, want 1", g.AmberOn(att))
	}
	if g.State.Aember[1] != 2 {
		t.Errorf("opponent pool = %d, want 2 after a capture", g.State.Aember[1])
	}
}

func TestEnemyCreatureDestroyedReaction(t *testing.T) {
	if !EventEnemyCreatureDestroyed.isReaction() {
		t.Error("EventEnemyCreatureDestroyed should be a reaction point")
	}
	if EventEnemyCreatureDestroyed.clause() != "each time an enemy creature is destroyed" {
		t.Errorf("clause = %q", EventEnemyCreatureDestroyed.clause())
	}

	g := NewGame("A", "B", 1)
	g.AddLasting(
		LastingEffect{
			On:         EventEnemyCreatureDestroyed,
			Do:         actGainAember,
			Controller: 0,
			Amount:     1,
		},
	)
	foe := g.AddToBattleline(testCreature("foe", 3), 1)
	g.destroyEach(0, []LocalID{foe})
	if g.State.Aember[0] != 1 {
		t.Errorf(
			"controller Æmber = %d, want 1 after an enemy creature was destroyed",
			g.State.Aember[0],
		)
	}
}

func TestFightFiresLasting(t *testing.T) {
	g := started(t)
	g.AddLasting(LastingEffect{
		On:         EventFight,
		Do:         actGainAember,
		Controller: 0,
		Amount:     1,
	})
	att := g.AddToBattleline(NewCard("att", Brobnar, Creature, Common, WithPower(4)), 0)
	def := g.AddToBattleline(testCreature("def", 2), 1)

	if err := g.Fight(0, att, def); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	if g.State.Aember[0] != 1 {
		t.Errorf(
			"controller Æmber = %d, want 1 after a Warsong-style fight reaction",
			g.State.Aember[0],
		)
	}
}

func TestFightFiresLastingLoseAember(t *testing.T) {
	g := started(t)
	g.AddLasting(LastingEffect{
		On:         EventFight,
		Do:         actLoseAember,
		Controller: 0,
		Amount:     1,
	})
	att := g.AddToBattleline(NewCard("att", Brobnar, Creature, Common, WithPower(4)), 0)
	def := g.AddToBattleline(testCreature("def", 2), 1)
	g.State.Aember[1] = 3

	if err := g.Fight(0, att, def); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	if g.State.Aember[1] != 2 {
		t.Errorf(
			"opponent Æmber = %d, want 2 after a Barn Razing-style fight reaction",
			g.State.Aember[1],
		)
	}
}

func TestEventFightIsReaction(t *testing.T) {
	if !EventFight.isReaction() {
		t.Error("EventFight should be a reaction point")
	}
	if EventFight.clause() != "each time a friendly creature fights" {
		t.Errorf("clause = %q", EventFight.clause())
	}
}

// TestEventTableIsTotal is the guard the event descriptor table exists for: it
// walks the whole Event enum, from the first real event to the last, and requires
// every value to be described. An event added to the enum but not to the table
// fails here instead of silently reading as a replacement point that renders with
// the generic played-a-creature clause.
func TestEventTableIsTotal(t *testing.T) {
	for e := eventUnset + 1; e <= EventNextTacticIntoHand; e++ {
		if _, ok := events[e]; !ok {
			t.Errorf("event %d has no table entry", e)
		}
	}
	// The unset zero value is not an event, so it is described by nothing.
	if _, ok := events[eventUnset]; ok {
		t.Error("the unset zero value should have no table entry")
	}
	if got := len(events); got != int(EventNextTacticIntoHand) {
		t.Errorf("table describes %d events, the enum holds %d", got, EventNextTacticIntoHand)
	}
}

// Only a reaction point carries a rendered clause of its own: a replacement or
// modifier point is queried at an outcome and never printed as a "when" phrase,
// so giving one a clause would mean it is rendered somewhere it should not be.
func TestOnlyReactionsCarryAClause(t *testing.T) {
	for e := eventUnset + 1; e <= EventNextTacticIntoHand; e++ {
		info := events[e]
		if info.isReaction {
			continue
		}
		if info.clause != "" || info.onOpponentTurn != "" {
			t.Errorf("event %d is not a reaction but carries a clause", e)
		}
	}
}
