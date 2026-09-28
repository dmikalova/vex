package engine

// This file is the census of the Effect family: one row per effect node the
// package declares, each naming the rulebook term it owes or the reason it owes
// none. See catalog.go for the classification rule, and effect.go for the AST
// the nodes belong to. The rows are grouped by mechanic, the way log_samples.go
// groups the log's variants and the way the effect_<mechanic>.go files group the
// nodes themselves.
//
// Each literal is filled only as far as the node needs to render and resolve
// sensibly — a Target where the node reads one, an Amount where the text prints
// one — rather than being an elaborate example. TestFamilyRowsWellFormed renders
// every one, so a literal too empty to print fails.

// effectFamily is the Effect family's census entry. An effect is discovered by
// its Resolve method, which no other family declares.
func effectFamily() Family {
	return newFamily(
		"Effect",
		"Resolve",
		[]string{"*EffectContext"},
		EffectCatalog(),
		Effect.Text,
	).gated()
}

// EffectCatalog returns one representative value of every effect node, with the
// rulebook term each owes.
func EffectCatalog() []Catalogued[Effect] {
	return []Catalogued[Effect]{
		// The Æmber economy: gaining, losing, and parking Æmber on a card.
		{
			Node:  GainAember{Player: Controller, Amount: 1},
			Rules: bears("Gain Æmber"),
		},
		{
			Node:  LoseAember{Player: Opponent, Amount: 1},
			Rules: bears("Lose Æmber"),
		},
		{
			Node: MoveAemberFromPool{
				Amount: 1,
				Target: Target{Kind: TargetThisCreature},
				Source: Controller,
			},
			Rules: bears("Move Æmber to a Card"),
		},
		{
			Node:  PlaceAemberOnThis{Amount: 1},
			Rules: bears("Move Æmber to a Card"),
		},

		// Capture: Æmber taken out of a pool and held on a creature.
		{
			Node: CaptureAember{
				Amount: 1,
				Target: Target{Kind: TargetThisCreature},
				Source: Opponent,
			},
			Rules: bears("Capture Æmber"),
		},
		{
			Node:  CaptureFromAnyPlayer{Amount: 1},
			Rules: bears("Capture Æmber"),
		},
		{
			Node:  DistributeCapture{All: true, Source: Opponent},
			Rules: bears("Capture Æmber"),
		},
		{
			Node:  RedistributeCapturedAember{Side: Controller},
			Rules: bears("Capture Æmber"),
		},

		// Æmber moving between pools, cards, and the common supply.
		{
			Node:  StealAember{Amount: 1},
			Rules: bears("Steal Æmber"),
		},
		{
			Node: MoveAember{
				Amount: 1,
				From:   Target{Kind: TargetChosenCreature},
				To:     Controller,
			},
			Rules: bears("Move Æmber"),
		},
		{
			Node: MoveAemberToSupply{
				Amount: 1,
				Target: Target{Kind: TargetChosenCreature},
			},
			Rules: bears("Move Æmber"),
		},
		{
			Node:  Exalt{Target: Target{Kind: TargetThisCreature}, Amount: 1},
			Rules: bears("Exalt"),
		},

		// Keys: forging one, undoing one, and what a key costs.
		{
			Node:  ForgeKey{Player: Controller},
			Rules: bears("Turn structure"),
		},
		{
			Node:  UnforgeKey{Player: Opponent},
			Rules: bears("Turn structure"),
		},
		{
			Node:  CancelForge{},
			Rules: bears("Turn structure"),
		},
		{
			Node:  SkipForgePhase{Player: Opponent},
			Rules: bears("Turn structure"),
		},
		{
			Node: RaiseKeyCost{
				Player:   Opponent,
				Amount:   1,
				Duration: OpponentNextTurn,
			},
			Rules: bears("Key Cost"),
		},
		{
			Node: LowerKeyCost{
				Player:   Controller,
				Amount:   1,
				Duration: RemainderOfPlayerTurn,
			},
			Rules: bears("Key Cost"),
		},
		{
			Node:  ScheduleOnLeave{Do: ForgeKey{Player: Controller}},
			Rules: plumbing("composition: defers its child until the source leaves play"),
		},

		// Chains, which tax the hand refill.
		{
			Node:  GainChains{Player: Controller, Amount: 1},
			Rules: bears("Gain Chains"),
		},

		// The turn: whose it is, which house is active, and ending it.
		{
			Node:  EndTurn{},
			Rules: bears("End the Turn"),
		},
		{
			Node:  ChangeActiveHouse{To: TheChosenHouse},
			Rules: bears("Active House"),
		},
		{
			Node:  ByActivePlayer{Do: Draw{Amount: 1}},
			Rules: plumbing("composition: re-aims its child at the active player"),
		},
		{
			Node:  ForEachHouse{Do: Draw{Amount: 1}},
			Rules: plumbing("composition: resolves its child once per house"),
		},

		// Houses: naming one, and a card belonging to one.
		{
			Node: BelongToHouse{
				Target:   Target{Kind: TargetThisCreature},
				House:    Mars,
				Duration: RemainderOfPlayerTurn,
			},
			Rules: bears("Belong to House"),
		},
		{
			Node:  NameHouse{Player: Opponent},
			Rules: bears("Name a House"),
		},

		// Playing and using cards, including out of house.
		{
			Node:  PlayFrom{From: Discard, Player: Controller},
			Rules: bears("Play a Card from Hand or Discard Pile"),
		},
		{
			Node:  PlayFromOpponent{From: Archives},
			Rules: bears("Play a Card from Hand or Discard Pile"),
		},
		{
			Node:  PlayItFromOpponentDiscard{},
			Rules: bears("Play a Card from Hand or Discard Pile"),
		},
		{
			Node: PlayOrUse{
				House: HouseMatcher{Kind: MatchChosenHouse},
				Grant: GrantPlay | GrantUse,
			},
			Rules: bears("Play or Use a Card"),
		},
		{
			Node: MayPlayOrUse{
				Houses: HouseSelector{
					Match: HouseMatcher{Kind: MatchNamedHouse, House: Mars},
				},
				Grant: GrantPlay,
			},
			Rules: bears("Play or Use a Card"),
		},
		{
			Node:  Use{Max: 1, Target: Target{Kind: TargetEachFriendlyCreature}},
			Rules: bears("Use a Card"),
		},

		// Readying and exhausting.
		{
			Node:  Ready{Target: Target{Kind: TargetThisCreature}},
			Rules: bears("Ready"),
		},
		{
			Node: ReadyCreatures{
				Max:    ForgedKeys{Player: Controller},
				Target: Target{Kind: TargetEachFriendlyCreature},
			},
			Rules: bears("Ready"),
		},
		{
			Node:  Exhaust{Target: Target{Kind: TargetChosenEnemyCreature}},
			Rules: bears("Exhaust"),
		},
		{
			Node: ExhaustCreatures{
				Max:    2,
				Target: Target{Kind: TargetEachEnemyCreature},
			},
			Rules: bears("Exhaust"),
		},

		// The hand: drawing, discarding it, and refilling it.
		{
			Node:  Draw{Amount: 1},
			Rules: bears("Draw"),
		},
		{
			Node:  DiscardHand{Player: Controller},
			Rules: bears("Discard"),
		},
		{
			Node:  RefillHand{Player: Controller},
			Rules: bears("Refill Hand"),
		},

		// Damage: dealing it, moving it, and refusing it.
		{
			Node:  DealDamage{Amount: 2, Target: Target{Kind: TargetChosenCreature}},
			Rules: bears("Deal Damage"),
		},
		{
			Node:  RedistributeDamage{},
			Rules: bears("Redistribute Damage"),
		},
		{
			Node:  RedirectFightDamage{Target: Target{Kind: TargetChosenCreature}},
			Rules: bears("Redirect Fight Damage"),
		},
		{
			Node: TakesExtraDamage{
				Target: Target{Kind: TargetEachEnemyCreature},
				Amount: 1,
			},
			Rules: bears("Deal Damage"),
		},
		{
			Node:  DamageOthersAfterUsingTrait{Trait: Dinosaur, Amount: 1},
			Rules: bears("Deal Damage"),
		},
		{
			Node: CannotBeDealtDamage{
				Target:   Target{Kind: TargetThisCreature},
				Duration: RemainderOfPlayerTurn,
			},
			Rules: bears("Cannot Be Dealt Damage"),
		},
		{
			Node:  Heal{Amount: 2, Target: Target{Kind: TargetChosenCreature}},
			Rules: bears("Heal"),
		},

		// Destruction.
		{
			Node:  Destroy{Target: Target{Kind: TargetChosenCreature}},
			Rules: bears("Destroy"),
		},
		{
			Node: DestroyChosen{
				Target: Target{Kind: TargetEachEnemyCreature},
				Amount: 2,
			},
			Rules: bears("Destroy"),
		},
		{
			Node: BatchDestroy{Gather: ChosenFromEach{
				Target{Kind: TargetEachFriendlyCreature},
				Target{Kind: TargetEachEnemyCreature},
			}},
			Rules: bears("Destroy"),
		},
		{
			Node:  DestroyEachCreatureAtEndOfTurn{},
			Rules: bears("Destroy"),
		},

		// Ward and armor, the two ways damage and destruction are absorbed.
		{
			Node:  Ward{Target: Target{Kind: TargetThisCreature}, Amount: 1},
			Rules: bears("Ward"),
		},
		{
			Node:  RemoveWard{Target: Target{Kind: TargetChosenCreature}},
			Rules: bears("Ward"),
		},
		{
			Node:  LoseArmor{Target: Target{Kind: TargetChosenCreature}},
			Rules: bears("Lose Armor"),
		},

		// Stun and enrage, the two states that spend a creature's next use.
		{
			Node:  Stun{Target: Target{Kind: TargetChosenEnemyCreature}},
			Rules: bears("Stun"),
		},
		{
			Node:  Unstun{Target: Target{Kind: TargetThisCreature}},
			Rules: bears("Unstun"),
		},
		{
			Node:  Enrage{Target: Target{Kind: TargetChosenCreature}},
			Rules: bears("Enrage"),
		},

		// Power and armor, whether granted, overridden, copied, or countered.
		{
			Node:  GainStats{Target: Target{Kind: TargetThisCreature}, Power: 2},
			Rules: bears("Power and Armor"),
		},
		{
			Node: OverrideStats{
				Power:    1,
				HasPower: true,
				Duration: RemainderOfPlayerTurn,
			},
			Rules: bears("Power and Armor"),
		},
		{
			Node: CopyPrintedStats{
				Target: Target{Kind: TargetThisCreature},
				Source: Target{Kind: TargetChosenCreature},
			},
			Rules: bears("Power and Armor"),
		},
		{
			Node: AddPowerCounter{
				Target: Target{Kind: TargetChosenCreature},
				Amount: 1,
			},
			Rules: bears("Power Counter"),
		},
		{
			Node: PlaceCounter{
				Kind:   CounterDoom,
				Target: Target{Kind: TargetChosenCreature},
				Amount: 1,
			},
			Rules: bears("Generic Counters"),
		},
		{
			Node: RemoveCounters{
				Kind:   CounterDoom,
				Target: Target{Kind: TargetChosenCreature},
				Amount: 1,
			},
			Rules: bears("Generic Counters"),
		},

		// Keywords, traits, and text boxes: what a card says it is.
		{
			Node: GainKeywords{
				Target:   Target{Kind: TargetThisCreature},
				Keywords: []Keyword{Taunt},
				Duration: RemainderOfPlayerTurn,
			},
			Rules: bears("Gain and Lose Keywords"),
		},
		{
			Node: LoseKeywords{
				Target:   Target{Kind: TargetEachEnemyCreature},
				Keywords: []Keyword{Taunt},
				Duration: RemainderOfPlayerTurn,
			},
			Rules: bears("Gain and Lose Keywords"),
		},
		{
			Node:  LoseKeyword{Keyword: Taunt},
			Rules: bears("Gain and Lose Keywords"),
		},
		{
			Node: GainTrait{
				Target: Target{Kind: TargetThisCreature},
				Trait:  Dinosaur,
			},
			Rules: bears("Gain a Trait"),
		},
		{
			Node: GainTextBox{
				Target: Target{Kind: TargetChosenCreature},
				Source: Target{Kind: TargetThisCreature},
			},
			Rules: bears("Text Box"),
		},
		{
			Node:  LendTextBoxFromHand{},
			Rules: bears("Text Box"),
		},
		{
			Node:  BlankEnemyText{},
			Rules: bears("Text Box"),
		},

		// Control and attachment.
		{
			Node: TakeControl{
				Target:   Target{Kind: TargetChosenEnemyCreature},
				Duration: RemainderOfPlayerTurn,
			},
			Rules: bears("Take Control"),
		},
		{
			Node:  AttachSelfTo{Target: Target{Kind: TargetChosenCreature}},
			Rules: bears("Attach an Upgrade"),
		},

		// Assault, granted as a keyword by an effect.
		{
			Node: GainAssault{
				Target: Target{Kind: TargetThisCreature},
				Amount: Fixed(2),
			},
			Rules: bears("Assault"),
		},
		{
			Node: GainAssaultUntilNextTurn{
				Target: Target{Kind: TargetThisCreature},
				Amount: Fixed(2),
			},
			Rules: bears("Assault"),
		},

		// The battleline: where a creature stands, and what counts as a flank.
		{
			Node:  ConsiderFlank{Target: Target{Kind: TargetThisCreature}},
			Rules: bears("Flank"),
		},
		{
			Node:  MoveToFlank{Target: Target{Kind: TargetChosenCreature}},
			Rules: bears("Battleline Position"),
		},
		{
			Node:  MoveWithinBattleline{Target: Target{Kind: TargetThisCreature}},
			Rules: bears("Battleline Position"),
		},
		{
			Node:  RearrangeBattleline{},
			Rules: bears("Battleline Position"),
		},
		{
			Node:  Swap{With: Target{Kind: TargetChosenCreature}},
			Rules: bears("Battleline Position"),
		},
		{
			Node:  SwapChosen{},
			Rules: bears("Battleline Position"),
		},
		{
			Node: TurnIntoCreature{
				Target:   Target{Kind: TargetThisCreature},
				Duration: UntilThisLeavesPlay,
			},
			Rules: bears("Become a Creature"),
		},

		// Choosing a creature and acting on it, which reads from its children.
		{
			Node: ChooseCreatureThen{
				Target: Target{Kind: TargetChosenCreature},
				Then:   Destroy{Target: Target{Kind: TargetTheChosenCreature}},
			},
			Rules: plumbing("composition: puts a chosen creature in context for its child"),
		},
		{
			Node: OnChooseCreature{
				Target: Target{Kind: TargetChosenFriendlyCreature},
				Verbs:  []CreatureVerb{ReadyVerb{}},
			},
			Rules: plumbing("composition: applies its verbs to a chosen creature"),
		},
		{
			Node: OneAtATime{
				Times:  Fixed(2),
				Target: Target{Kind: TargetEachFriendlyCreature},
				Verbs:  []CreatureVerb{UseVerb{}},
			},
			Rules: plumbing("composition: applies its verbs to one creature at a time"),
		},
		{
			Node: RepeatedFight{
				Times:  Fixed(2),
				Target: Target{Kind: TargetChosenFriendlyCreature},
			},
			Rules: plumbing(
				"composition: repeats a ready-and-fight against a different enemy each time",
			),
		},

		// Restrictions on what a player or their creatures may do.
		{
			Node: Restrict{
				Player:   Opponent,
				Action:   RestrictFighting,
				Duration: OpponentNextTurn,
			},
			Rules: bears("Restriction"),
		},
		{
			Node: CreaturesCannot{
				Action:   ReapUse,
				Duration: RemainderOfPlayerTurn,
			},
			Rules: bears("Restriction"),
		},
		{
			Node: CannotPlay{
				Player:   Opponent,
				Type:     Creature,
				Duration: OpponentNextTurn,
			},
			Rules: bears("Restriction"),
		},
		{
			Node: PlayersCannotPlay{
				Type:     Creature,
				Duration: RemainderOfPlayerTurn,
			},
			Rules: bears("Restriction"),
		},
		{
			Node: MustChooseHouse{
				Player:    Opponent,
				Reference: ChosenActiveHouse,
			},
			Rules: bears("Name a House"),
		},
		{
			Node: CannotChooseHouse{
				Player:    Opponent,
				Reference: JustChosenActiveHouse,
			},
			Rules: bears("Name a House"),
		},
		{
			Node:  WagerOpponentChoosesChosenHouse{Amount: 1},
			Rules: bears("Name a House"),
		},

		// Lasting effects: a child attached to a later event, or an outcome replaced.
		{
			Node: ForRemainderOfTurn{
				On: EventCreaturePlayed,
				Do: GainAember{Player: Controller, Amount: 1},
			},
			Rules: plumbing("composition: attaches its child to a later event this turn"),
		},
		{
			Node: ForOpponentNextTurn{
				On: EventCreaturePlayed,
				Do: GainAember{Player: Opponent, Amount: 1},
			},
			Rules: plumbing(
				"composition: attaches its child to a later event on the opponent's next turn",
			),
		},
		{
			Node: NextPlayed{
				Type:       Creature,
				EntersPlay: Ready{Target: Target{Kind: TargetThisCreature}},
			},
			Rules: plumbing("composition: attaches its child to the next card played"),
		},
		{
			Node: GainAbility{
				Target: Target{Kind: TargetEachFriendlyCreature},
				Ability: Ability{
					Trigger: TriggerAfterReap,
					Effect:  GainAember{Player: Controller, Amount: 1},
				},
				Duration: RemainderOfPlayerTurn,
			},
			Rules: bears("Gain an Ability"),
		},
		{
			Node: Instead{
				Of:     EventAemberAddedToPool,
				With:   Steal,
				Player: Controller,
			},
			Rules: bears("Replacement"),
		},
		{
			Node:  PutNextTacticIntoHand{},
			Rules: bears("Replacement"),
		},
		{
			Node:  FuseTriggersForTurn{A: TriggerAfterReap, B: TriggerAfterFight},
			Rules: bears("Fused Triggers"),
		},
		{
			Node: TriggerAbility{
				Trigger: TriggerAfterPlay,
				Target:  Target{Kind: TargetChosenFriendlyCreature},
			},
			Rules: bears("Trigger Another Card's Ability"),
		},

		// Bonus icons resolved by an effect rather than by playing the card.
		{
			Node:  ResolveBonusIcons{Target: Target{Kind: TargetChosenFriendlyCreature}},
			Rules: bears("Resolve Bonus Icons"),
		},
		{
			Node:  ExtraBonusIconResolution{},
			Rules: bears("Resolve Bonus Icons"),
		},

		// One pool handing Æmber to the other.
		{
			Node:  GiveAember{Amount: 1},
			Rules: bears("Move Æmber"),
		},

		// Archiving.
		{
			Node: ArchiveCard{
				Zone:      Hand,
				From:      Controller,
				Selection: Chosen{},
				Quantity:  Takes{N: Fixed(1)},
			},
			Rules: bears("Archive"),
		},
		{
			Node:  ArchiveFromPlay{Target: Target{Kind: TargetChosenFriendlyCreature}},
			Rules: bears("Archive"),
		},
		{
			Node:  ArchiveSource{},
			Rules: bears("Archive"),
		},
		{
			Node:  ArchiveGrantingUpgrade{},
			Rules: bears("Archive"),
		},
		{
			Node:  ArchiveCardUnder{},
			Rules: bears("Archive"),
		},
		{
			Node:  ArchiveDiscardedThisWay{},
			Rules: bears("Archive"),
		},

		// Discarding.
		{
			Node: DiscardCard{
				Player:    Controller,
				Zones:     []Zone{Hand},
				Selection: Chosen{},
				Quantity:  Takes{N: Fixed(1)},
			},
			Rules: bears("Discard"),
		},
		{
			Node:  DiscardFromOpponent{Sources: []Zone{Deck}},
			Rules: bears("Discard"),
		},
		{
			Node:  DiscardTop{Player: Opponent, Amount: 1},
			Rules: bears("Discard"),
		},
		{
			Node:  DiscardUntil{Player: Opponent, MayStop: true},
			Rules: bears("Discard"),
		},
		{
			Node:  DiscardArchives{Player: Opponent},
			Rules: bears("Discard"),
		},

		// Purging.
		{
			Node: PurgeCard{
				Zones:     []Zone{Discard},
				Player:    Controller,
				Selection: Chosen{},
				Quantity:  Takes{N: Fixed(1)},
			},
			Rules: bears("Purge"),
		},
		{
			Node:  PurgeCreature{Target: Target{Kind: TargetChosenCreature}},
			Rules: bears("Purge"),
		},
		{
			Node:  PurgeSource{},
			Rules: bears("Purge"),
		},
		{
			Node:  PurgeArchivedCardThen{Then: Draw{Amount: 1}},
			Rules: bears("Purge"),
		},

		// Moving a card between zones, which ADR 0031 treats as one mechanism.
		{
			Node: PutCard{
				Zones:       []Zone{Discard},
				Selection:   Chosen{},
				Destination: ToHand,
			},
			Rules: bears("Put a Card into Another Zone"),
		},
		{
			Node: PutChosen{
				Quantity:    Takes{N: Fixed(1)},
				Target:      Target{Kind: TargetEachEnemyCreature},
				Destination: ToHand,
			},
			Rules: bears("Put a Card into Another Zone"),
		},
		{
			Node: PutFromPlay{
				Target:      Target{Kind: TargetChosenCreature},
				Destination: ToHand,
			},
			Rules: bears("Put a Card into Another Zone"),
		},
		{
			Node:  PutItIntoHand{},
			Rules: bears("Put a Card into Another Zone"),
		},
		{
			Node:  PutFromHand{Filter: Filter{Type: Creature}},
			Rules: bears("Put a Card into Another Zone"),
		},
		{
			Node:  PutDiscardedIntoHand{Noun: Creature},
			Rules: bears("Put a Card into Another Zone"),
		},
		{
			Node:  PutRevealedCard{To: IntoArchives},
			Rules: bears("Put a Card into Another Zone"),
		},

		// Putting a card into play, the movement that skips being played.
		{
			Node: PutIntoPlay{
				Target:  Target{Kind: TargetChosenCreature},
				Control: ControlYours,
			},
			Rules: bears("Put a Card into Play"),
		},
		{
			Node:  PutDiscardedIntoPlay{Noun: Creature},
			Rules: bears("Put a Card into Play"),
		},
		{
			Node:  EachPlayerPutsHandCreaturesIntoPlay{Ready: true},
			Rules: bears("Put a Card into Play"),
		},

		// Cards set under another card.
		{
			Node:  PutUnderFromHand{Filter: Filter{Type: Creature}},
			Rules: bears("Put a Card From Hand Under a Card"),
		},
		{
			Node:  PutUnderIntoPlay{},
			Rules: bears("Put the Cards Under a Card Into Play"),
		},
		{
			Node:  PlayCardUnder{},
			Rules: bears("Play the Card Under a Card"),
		},
		{
			Node:  Graft{Target: Target{Kind: TargetChosenFriendlyCreature}},
			Rules: bears("Graft a Card"),
		},
		{
			Node:  TriggerGraftedPlayEffect{},
			Rules: bears("Graft a Card"),
		},

		// Shuffling.
		{
			Node:  Shuffle{Zones: []Zone{Discard}},
			Rules: bears("Shuffle"),
		},
		{
			Node: ShuffleIntoDeck{
				Player:    Controller,
				From:      []Zone{Discard},
				Selection: Chosen{},
				Quantity:  Takes{N: Fixed(1)},
			},
			Rules: bears("Shuffle"),
		},
		{
			Node:  SwapDeckAndDiscard{},
			Rules: bears("Swap Deck And Discard"),
		},

		// Reading and revealing cards.
		{
			Node:  RevealHand{Player: Opponent},
			Rules: bears("Reveal a Card"),
		},
		{
			Node:  RevealRandomFromHand{},
			Rules: bears("Reveal a Card"),
		},
		{
			Node:  RevealChosenFromHand{},
			Rules: bears("Reveal a Card"),
		},
		{
			Node: RevealTopOfDeck{
				Amount: 1,
				Player: Controller,
				Then:   []TopAct{ReorderRest{}},
			},
			Rules: bears("Reveal Top of Deck"),
		},
		{
			Node: LookAtTopOfDeck{
				Amount: 1,
				Then:   []TopAct{MayDiscardLookedAt{}},
			},
			Rules: bears("Reveal Top of Deck"),
		},
		{
			Node:  Search{Sources: []Zone{Deck}, Max: 1, Dest: ToHand},
			Rules: bears("Search"),
		},
		{
			Node:  NameCard{},
			Rules: bears("Name a Card"),
		},

		// Playing a card an effect turned up.
		{
			Node:  PlayTopOfDeck{},
			Rules: bears("Play a Card from Hand or Discard Pile"),
		},
		{
			Node:  PlayRevealedCard{},
			Rules: bears("Play a Card from Hand or Discard Pile"),
		},
		{
			Node:  CancelFight{},
			Rules: bears("Cancel a Fight"),
		},

		// Composition: the nodes that join, gate, repeat, or time their children.
		{
			Node: Sequence{Effects: []Effect{
				Draw{Amount: 1},
				GainAember{Player: Controller, Amount: 1},
			}},
			Rules: plumbing("composition: resolves its children in order"),
		},
		{
			Node: Then{
				First:  Destroy{Target: Target{Kind: TargetChosenCreature}},
				Result: GainAember{Player: Controller, Amount: 1},
			},
			Rules: bears("Result Gate"),
		},
		{
			Node:  May{Do: Draw{Amount: 1}},
			Rules: bears("May"),
		},
		{
			Node: ChooseOne{Options: []Effect{
				Draw{Amount: 1},
				GainAember{Player: Controller, Amount: 1},
			}},
			Rules: bears("Choose One"),
		},
		{
			Node:  ChooseHouseThen{Then: ChangeActiveHouse{To: TheChosenHouse}},
			Rules: plumbing("composition: puts a chosen house in context for its child"),
		},
		{
			Node:  OpponentNamesHouse{},
			Rules: bears("Name a House"),
		},
		{
			Node: Conditional{
				Cond: SourceReady{},
				Then: Draw{Amount: 1},
			},
			Rules: bears("Conditional"),
		},
		{
			Node: ForEach{
				Times: ForgedKeys{Player: Controller},
				Do:    Draw{Amount: 1},
			},
			Rules: plumbing("composition: resolves its child once per count"),
		},
		{
			Node: ForEachDiscarded{
				Filter: Filter{Type: Creature},
				Do:     GainAember{Player: Controller, Amount: 1},
			},
			Rules: plumbing("composition: resolves its child once per card discarded this way"),
		},
		{
			Node: Repeat{
				Do:   Draw{Amount: 1},
				Gate: ByExalting{Creature: Target{Kind: TargetThisCreature}},
			},
			Rules: bears("Repeat"),
		},
		{
			Node: ForDuration{
				Duration: RemainderOfPlayerTurn,
				Effects: []Effect{
					BelongToHouse{
						Target:   Target{Kind: TargetThisCreature},
						House:    Mars,
						Duration: RemainderOfPlayerTurn,
					},
					CannotBeDealtDamage{
						Target:   Target{Kind: TargetThisCreature},
						Duration: RemainderOfPlayerTurn,
					},
				},
			},
			Rules: plumbing("composition: gives its children a duration"),
		},
		{
			Node: GainUntilNextTurn{Effects: []Effect{
				GainKeywords{
					Target:   Target{Kind: TargetThisCreature},
					Keywords: []Keyword{Taunt},
					Duration: StartOfPlayerNextTurn,
				},
				GainAssaultUntilNextTurn{
					Target: Target{Kind: TargetThisCreature},
					Amount: Fixed(2),
				},
			}},
			Rules: plumbing("composition: gives its children a duration"),
		},
	}
}
