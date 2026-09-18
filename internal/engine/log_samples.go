package engine

// This file holds the log's sample catalog (ADR 0046): one representative value
// of every LogEntry variant the game can narrate. It is the log counterpart to
// Keywords() — a closed, enumerable set co-located with the entries it names — so
// a consumer that wants to show or exercise "every kind of log line" has a single
// source to draw from rather than reconstructing the set by playing games.
//
// The catalog is kept complete by a totality test (TestLogEntrySamplesTotality):
// a variant with no sample fails the build. Its wording is pinned, index for
// index, by TestLogEntryText, which renders every sample and asserts its text.
// So this list and that test move together — adding a variant means adding a
// sample here and its wording there, and a mismatch is loud.

// LogEntrySamples returns one representative value of every LogEntry variant, in
// a stable order grouped the way the log families are. The gallery uses it as
// the catalog it covers real games against, falling back to a sample's own
// rendering for a variant a game run never produced (ADR 0046).
func LogEntrySamples() []LogEntry {
	return []LogEntry{
		// Turn shape.
		TurnBegan{Player: 0, Turn: 3},
		GameStarted{FirstPlayer: 0},
		Mulliganed{Player: 1, Hand: 5},
		PhaseBegan{Player: 1, Phase: PhaseReady},
		CardsReadied{Player: 0, Cards: []LocalID{4, 7}},
		CardsDrawn{Player: 1, Count: 3, Hand: 6},
		CardsDrawn{Player: 1, Count: 1, Hand: 6},
		CardsDrawn{Player: 0, Count: 0, Hand: 2},
		CardsDrawnBy{Player: 0, Count: 1},
		HouseChosen{Player: 1, House: Brobnar},
		ForgeSkipped{Player: 0},
		AemberGainedFromForging{Card: 5, From: 1, Amount: 6},
		KeyForged{Player: 0, Color: KeyColorRed, Keys: 1, Needed: 3},
		KeyForged{Player: 1, Color: KeyColorColorless, Keys: 4, Needed: 3},
		KeyUnforged{Player: 0, Keys: 1, Needed: 3},
		ChainShed{Player: 1, Remaining: 4},
		GameWon{Player: 0},
		PlayerConceded{Player: 1},
		PlayerStanding{Player: 0, Aember: 4, KeyColors: []KeyColor{KeyColorRed}},

		// Æmber. The source-card subject is exercised in
		// TestRecordTextSubjectsToSourceCard; here the entries render unframed, so
		// they name the player.
		AemberGained{Player: 0, Amount: 2},
		AemberLost{Player: 1, Amount: 1},
		AemberStolen{Player: 0, From: 1, Amount: 2},
		AemberStolen{Player: 0, From: 1, Amount: 2, FromSupply: true, Cause: 9},
		AemberCaptured{Creature: 7, Amount: 3, Source: 7},
		AemberCaptured{Creature: 7, Amount: 1, Source: 3},
		AemberCaptured{Creature: 7, Amount: 3, Source: 7, FromSupply: true, Cause: 9, From: 0},
		AemberCaptured{Creature: 7, Amount: 1, Source: 3, FromSupply: true, Cause: 9, From: 0},
		AemberMovedToCommonSupply{Creature: 7, Amount: 1},
		AemberCapturedInsteadOfGain{Creature: 7, Player: 1, Amount: 1},
		AemberCapturedInsteadOfSteal{Cause: 9, Creature: 7, Player: 0, Amount: 2},
		AemberCapturedInsteadOfSteal{
			Cause:      9,
			Creature:   7,
			Player:     0,
			Amount:     2,
			FromSupply: true,
		},
		AemberExalted{Player: 0, Creature: 4, Amount: 2},
		AemberMovedToPool{Player: 0, From: 4, To: 1, Amount: 2},
		AemberMovedToCard{Player: 0, From: 4, To: 5, Amount: 1},
		AemberLostToMaximum{Amount: 2},

		// Creatures and cards in play.
		CreatureReadied{Creature: 2},
		CreatureGainedKeyword{Creature: 2, Keyword: Skirmish},
		CreatureLostKeyword{Creature: 2, Keyword: Elusive},
		CreatureGainedStats{Creature: 2, Armor: 1},
		CreatureGainedStats{Creature: 2, Power: 2, Armor: 2},
		CreatureGainedAssault{Creature: 2, Amount: 3},
		CreatureGainedTrait{Creature: 2, Trait: Mutant},
		CreatureConsideredFlank{Creature: 2},
		CreatureExhausted{Creature: 2},
		CreatureStunned{Creature: 2, By: 2},
		CreatureStunned{Creature: 2, By: 5},
		CreatureStunned{Creature: 2, By: 5, AlreadyStunned: true},
		CreaturesUnstunned{Player: 0, Creatures: []LocalID{2, 5}},
		CreatureEnraged{Creature: 2, By: 2},
		CreatureEnraged{Creature: 2, By: 5},
		CreatureEnraged{Creature: 2, By: 5, AlreadyEnraged: true},
		CreatureWarded{Creature: 2, By: 2},
		CreatureWarded{Creature: 2, By: 5},
		CreatureWarded{Creature: 2, By: 5, AlreadyWarded: true},
		WardRemoved{Creature: 2, By: 5},
		WardRemoved{Creature: 2, By: 5, AlreadyUnwarded: true},
		WardAbsorbed{Creature: 2},
		WardAbsorbed{Creature: 2, Prevented: wardDestruction},
		WardAbsorbed{Creature: 2, Prevented: wardDamage, Amount: 5},
		NoCreatureToFight{Creature: 2},
		CardsRevealedToAll{Player: 0, Cards: []LocalID{1, 2}},
		KeyForgePrevented{Player: 1, By: 3},
		PositionsSwapped{A: 1, B: 2},
		CardsSwapped{A: 1, B: 2, FromPlayer: 1, FromZone: Discard},
		MovedToFlank{Creature: 2, Right: true},
		MovedToFlank{Creature: 2},
		MovedWithinBattleline{Creature: 2},
		TurnedIntoCreature{Card: 2, Right: true},
		TurnedIntoCreature{Card: 2},
		RevertedToArtifact{Card: 2},
		ControlTaken{Player: 1, Card: 3},
		ControlReturned{Card: 3, Owner: 0},
		CardDestroyed{Card: 3},
		CardsDestroyedBy{Source: 5, Cards: []LocalID{3}},
		CardsDestroyedBy{Source: 5, Cards: []LocalID{3, 8}},
		CardsDestroyedBy{Source: 5, Cards: []LocalID{3, 8, 9}},
		DestructionReplaced{Card: 3, By: 8},
		AemberOnCardReleased{Card: 3, Amount: 2, To: 1},
		StunRecovered{Player: 1, Creature: 3},
		CardCannotBeUsed{Card: 3},

		// Copied stats and text box. These entries are recorded beside the effects
		// that produce them (effect_copy.go, effect_textbox.go), so they are easy to
		// miss when enumerating the log by family; the catalog keeps them in view.
		CreatureCopiedStats{Creature: 2, Source: 5},
		CreatureGainedTextBox{Creature: 2, Source: 5},

		// Combat.
		FightCancelled{Attacker: 1},
		Fought{Attacker: 1, AttackerPower: 4, Defender: 2, DefenderPower: 3},
		Fought{
			Attacker:         1,
			AttackerPower:    6,
			AttackerKeywords: FightSkirmish,
			Defender:         2,
			DefenderPower:    3,
			DefenderKeywords: FightElusive | FightPoison,
		},
		PoisonKills{Source: 1, Victim: 2},
		DamageRefused{Creature: 2},
		ArmorAbsorbed{Creature: 2, Amount: 1},
		DamageTaken{Creature: 2, Amount: 3, Total: 4},
		AssaultDealt{Source: 1, Amount: 2, Target: 2},
		HazardousDealt{Source: 2, Amount: 5, Target: 1},
		AbilityDamageDealt{Amount: 4, Target: 2},

		// Zones.
		ArchivesTakenIntoHand{Player: 0, Count: 2},
		CardMoved{Player: 0, Card: 6, From: Hand, To: Archives},
		CardMoved{Player: 1, Card: 6, From: Discard, To: Archives},
		TopOfDeckArchived{Player: 0, Card: 6},
		ArchivesDiscarded{Player: 1, Count: 3},
		TopOfDeckDiscarded{Player: 0, Card: 6},
		CardMoved{Player: 0, Card: 6, From: Deck, To: Discard},
		DeckAndDiscardSwapped{Player: 1},
		ShuffledIntoDeck{Player: 0, DiscardCards: []LocalID{6}, HandCount: 2},
		ShuffledIntoDeck{Player: 1, DiscardCards: []LocalID{4, 6}},
		ShuffledIntoDeck{Player: 0, ArchivesCount: 1},
		ShuffledIntoDeck{Player: 1},
		CardDiscarded{Player: 0, Card: 6},
		CardMoved{Player: 0, Card: 6, From: Archives, To: Discard},
		CardMoved{Player: 0, Card: 6, From: Discard, To: purged},
		CardMoved{Player: 1, Card: 6, From: Hand, To: purged},
		CardPurgedFromHand{Card: 6, Owner: 0},
		CardMoved{Player: 1, Card: 6, From: Archives, To: purged},
		CardMoved{Player: 1, Card: 6, From: Deck, To: purged},
		CardPurged{Card: 6},
		CardPutOnTopOfDeck{Card: 6, Owner: 0},
		CardPutIntoHand{Card: 6, Owner: 1},
		CardPutIntoArchives{Card: 6, Owner: 0},
		CardArchivedFromPurge{Player: 0, Card: 6},
		CardShuffledIntoDeck{Card: 6, Owner: 1},
		DeckShuffled{Player: 1},
		DiscardRecycledIntoDeck{Player: 1},
		CardsShuffledIntoDeckBy{Owner: 1, Cards: []LocalID{3, 8}},
		CardAbducted{Player: 0, Card: 6, Owner: 1},
		CardPutFromDiscardIntoHand{Player: 0, Card: 6},
		CardPutFromDeckIntoHand{Player: 1, Card: 6},
		CardPutFromDiscardOnTopOfDeck{Player: 0, Card: 6},
		CardPutUnder{Player: 0, Card: 6, Host: 9, FaceDown: false},
		CardPutUnder{Player: 0, Card: 6, Host: 9, FaceDown: true},
		CardGrafted{Card: 6, Host: 9},

		// Playing and using.
		CardPlayedToBattleline{Player: 0, Card: 9},
		CardPlayedToBattleline{Player: 0, Card: 9, FlankLeft: true},
		CardPlayedToBattleline{Player: 0, Card: 9, Interior: true},
		ArtifactPlayed{Player: 1, Card: 9},
		TacticPlayed{Player: 0, Card: 9},
		UpgradeAttached{Player: 0, Upgrade: 9, Host: 2},
		CardPutIntoPlay{Player: 1, Card: 9},
		PlayedFromTopOfDeck{Card: 9, Player: 0},
		BonusAemberGained{Player: 0, Card: 9, Amount: 2},
		BonusAemberCaptured{Creature: 7, Card: 9, Player: 0, Amount: 2},
		BonusCaptured{Creature: 7, Card: 9, Amount: 1},
		BonusDamageDealt{Source: 9, Amount: 1, Target: 2},
		BonusCardDrawn{Player: 0, Card: 9, Amount: 1},
		AemberSpentToPlay{Player: 0, Card: 9, Amount: 1},
		Reaped{Player: 0, Card: 2},
		ReapedStealing{Player: 0, Card: 2, Amount: 1, Cause: 9},
		ReapedStealing{Player: 0, Card: 2, Cause: 9},
		ReapedCaptured{Player: 0, Card: 2, Creature: 7},
		ActionAbilityUsed{Player: 1, Card: 2},

		// Lasting effects. A lasting Æmber gain records the plain AemberGained/
		// AemberCapturedInsteadOfGain entry under its source-card frame, so only the
		// draw keeps its own lasting entry (and its event context).
		LastingDraw{Player: 1, Amount: 2, On: EventFight},
		AemberGivenAfterForging{Player: 0, To: 1, Amount: 3},
		AemberGiven{Giver: 0, Receiver: 1, Amount: 1},
		AemberGiven{Giver: 0, Receiver: 1, Amount: 1, Reason: TollUseArtifact},

		// Grants, chains, and manual mode.
		MayPlayOrUseGranted{
			Player: 0,
			Houses: HouseSelector{Match: HouseMatcher{Kind: MatchNamedHouse, House: Brobnar}},
			Grant:  GrantFight,
		},
		MayPlayOrUseGranted{
			Player: 1,
			Houses: HouseSelector{Match: HouseMatcher{Kind: MatchAnyHouse}},
			Grant:  GrantFight,
		},
		MayPlayOrUseGranted{
			Player: 0,
			Houses: HouseSelector{Match: HouseMatcher{Kind: MatchNamedHouse, House: Dis}},
			Grant:  GrantUse,
		},
		MayPlayOrUseGranted{
			Player: 0,
			Houses: HouseSelector{Match: HouseMatcher{Kind: MatchAnyHouse}},
			Grant:  GrantUse,
			Types:  CardTypesOf(Artifact),
		},
		MayPlayOrUseGranted{
			Player: 0,
			Houses: HouseSelector{Match: HouseMatcher{Kind: MatchNamedHouse, House: Mars}},
			Grant:  GrantPlay,
		},
		MayPlayOrUseGranted{
			Player: 0,
			Houses: HouseSelector{Match: HouseMatcher{Kind: MatchNamedHouse, House: Mars}},
			Grant:  GrantPlay | GrantUse,
		},
		MayPlayOrUseGranted{
			Player: 0,
			Houses: HouseSelector{Match: HouseMatcher{Kind: MatchExceptHouse, House: StarAlliance}},
			Grant:  GrantPlay,
			Count:  1,
		},
		MayPlayOrUseGranted{
			Player: 0,
			Houses: HouseSelector{Controlled: true},
			Grant:  GrantPlay,
		},
		MayUseTraitGranted{Player: 0, Trait: Mutant},
		HouseForcedNextTurn{Player: 1, House: Logos},
		HouseForbiddenNextTurn{Player: 1, House: Logos},
		HouseWagerArmed{Predictor: 0, Player: 1, House: Logos, Amount: 2},
		KeywordLostByAll{Keyword: Elusive},
		ChainsGained{Player: 0, Amount: 2, Total: 5},
		ManualCardMoved{Player: 0, Card: 3, To: ManualPurge},
		ManualExhaustSet{Card: 3, Exhausted: true},
		ManualExhaustSet{Card: 3},
		ManualPlacedInPlay{Player: 0, Card: 3},
		ManualMatchFull{Player: 1},
		ManualCardAdded{Player: 0, Card: 3},
		ManualAemberSet{Player: 0, Amount: 7},
		ManualChainsSet{Player: 1, Amount: 1},
		ManualHouseChosen{Player: 0, House: Mars},
		ManualKeyForged{Player: 0, Color: KeyColorBlue, Keys: 2, Needed: 3},
		ManualKeyUnforged{Player: 1, Keys: 0, Needed: 3},

		// A restored entry reads back exactly as it was narrated.
		RestoredEntry{Line: "P0 gains 1 Æmber"},
	}
}
