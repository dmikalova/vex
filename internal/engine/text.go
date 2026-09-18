package engine

import (
	"fmt"
	"strings"
)

// RenderAbility renders a single triggered ability to its printed card line,
// e.g. "After you forge a key, deal 2 damage to each enemy creature."
func RenderAbility(a Ability) string {
	if a.Trigger == TriggerAfterCardPlayed {
		if s, ok := afterYouActOnText("play", a.Effect); ok {
			return punctuate(capitalizeFirst(s))
		}
	}
	if a.Trigger == TriggerAfterUse {
		if s, ok := afterYouActOnText("use", a.Effect); ok {
			return punctuate(capitalizeFirst(s))
		}
	}
	if a.Trigger == TriggerAfterDiscardFromHand {
		if s, ok := afterYouActOnText("discard", a.Effect); ok {
			return punctuate(capitalizeFirst(s))
		}
	}
	if a.Trigger == TriggerAfterChooseHouse {
		if s, ok := afterChooseHouseText(a.Effect); ok {
			return punctuate(capitalizeFirst(s))
		}
	}
	if a.Trigger == TriggerAfterAnyPlayerChoosesHouse {
		if s, ok := afterAnyPlayerChooseHouseText(a.Effect); ok {
			return punctuate(capitalizeFirst(s))
		}
	}
	if a.Trigger == TriggerAfterCreaturePlayedAdjacent {
		if s, ok := afterCreaturePlayedAdjacentText(a.Effect); ok {
			return punctuate(capitalizeFirst(s))
		}
	}
	if a.Trigger == TriggerAfterCreatureReaps {
		if s, ok := afterCreatureReapsText(a.Effect); ok {
			return punctuate(capitalizeFirst(s))
		}
	}
	if a.Trigger == TriggerAfterCreatureDestroyed {
		if s, ok := afterCreatureDestroyedText(a.Effect); ok {
			return punctuate(capitalizeFirst(s))
		}
	}
	if a.Trigger == TriggerAfterCreatureFights {
		if s, ok := afterCreatureFightsText(a.Effect); ok {
			return punctuate(capitalizeFirst(s))
		}
	}
	if a.Trigger == TriggerAfterEnemyCardPlayed {
		if s, ok := afterEnemyPlaysCreatureOnFlankText(a.Effect); ok {
			return punctuate(capitalizeFirst(s))
		}
	}
	if a.Trigger == TriggerEntersPlay {
		if s, ok := entersPlayConditionalText(a.Effect); ok {
			return punctuate(capitalizeFirst(s))
		}
		return SelfName + " enters play " + enterStateWord(a.Effect) + "."
	}
	prefix, capitalize := a.Trigger.prefix()
	body := a.Effect.Text()
	if capitalize {
		body = capitalizeFirst(body)
	}
	return prefix + punctuate(body)
}

// afterYouActOnText folds an "after you <verb> a card" reaction gated only on the
// subject card's shape — a Conditional{ItIs} — into the natural "after you <verb>
// a <shape>, <then>" wording: Carlo Phantom's "after you play an artifact, steal
// 1 Æmber", Veylan Analyst's "after you use an artifact, gain 1 Æmber", and Baron
// Mengevin's "after you discard a Sanctum card, ...", rather than the literal
// "after you <verb> a card, if it is a <shape>, ...". A Conditional{ItIsNamed}
// folds the same way to a named card — Chain Gang's "after you play Subtle Chain,
// ready Chain Gang" — and a Conditional{ItIsOfTrait} to a trait creature — Dark
// Æmber Vault's "after you play a Mutant creature, draw a card". Any other effect
// renders with the broad prefix, so an unconditional or state-gated reaction still
// reads "After you <verb> a card, ...".
func afterYouActOnText(verb string, e Effect) (string, bool) {
	cond, ok := e.(Conditional)
	if !ok {
		return "", false
	}
	switch it := cond.Cond.(type) {
	case ItIs:
		return "after you " + verb + " " + indefinite(
			it.shapeNoun(),
		) + ", " + cond.Then.Text(), true
	case ItIsNamed:
		return "after you " + verb + " " + it.Name + ", " + cond.Then.Text(), true
	case ItIsOfTrait:
		return "after you " + verb + " " + indefinite(
			it.Trait.String()+" creature",
		) + ", " + cond.Then.Text(), true
	case ItHasBonusIcon:
		return "after you " + verb + " a card with a bonus icon, " + cond.Then.Text(), true
	default:
		return "", false
	}
}

// afterCreatureDestroyedText folds an AfterCreatureDestroyed reaction gated only on
// whose creature was destroyed — a Conditional{ItIsFriendly} or {And{ItIsEnemy,
// ItIsYourTurn}} — into "after a friendly creature is destroyed, <then>" or "after
// an enemy creature is destroyed during your turn, <then>", rather than the literal
// "after a creature is destroyed, if …". Any other effect shape reports false and
// renders with the broad prefix.
func afterCreatureDestroyedText(e Effect) (string, bool) {
	cond, ok := e.(Conditional)
	if !ok || cond.Else != nil {
		return "", false
	}
	switch c := cond.Cond.(type) {
	case ItIsFriendly:
		return "after a friendly creature is destroyed, " + cond.Then.Text(), true
	case And:
		if isEnemyDuringYourTurn(c) {
			return "after an enemy creature is destroyed during your turn, " +
				cond.Then.Text(), true
		}
	}
	return "", false
}

// isEnemyDuringYourTurn reports whether an And is exactly the enemy-creature,
// your-turn pair that "after an enemy creature is destroyed during your turn" folds.
func isEnemyDuringYourTurn(a And) bool {
	if len(a.Conditions) != 2 {
		return false
	}
	_, enemy := a.Conditions[0].(ItIsEnemy)
	_, yourTurn := a.Conditions[1].(ItIsYourTurn)
	return enemy && yourTurn
}

// afterCreatureFightsText folds an AfterCreatureFights reaction gated only on whose
// creature fought — a Conditional{ItIsFriendly} or {ItIsEnemy} — into "after a
// friendly creature is used to fight, <then>" (Lieutenant Gorvenal) rather than the
// literal "after a creature is used to fight, if it is a friendly creature, …". Any
// other effect shape reports false and renders with the broad prefix.
func afterCreatureFightsText(e Effect) (string, bool) {
	cond, ok := e.(Conditional)
	if !ok || cond.Else != nil {
		return "", false
	}
	adj, ok := scopedCreatureAdjective(cond.Cond)
	if !ok {
		return "", false
	}
	return "after " + indefinite(adj+" creature") + " is used to fight, " +
		cond.Then.Text(), true
}

// afterCreatureReapsText folds an AfterCreatureReaps reaction gated only on whose
// creature reaped — a Conditional{ItIsEnemy} or {ItIsFriendly} — into "after an
// enemy creature reaps, <then>" or "after a friendly creature reaps, <then>",
// rather than the literal "after a creature reaps, if it is an enemy creature, …".
// Aember Conduction Unit's inner Conditional (FirstReapOfTurn) rides along in
// <then>. Any other effect shape reports false and renders with the broad prefix.
func afterCreatureReapsText(e Effect) (string, bool) {
	cond, ok := e.(Conditional)
	if !ok || cond.Else != nil {
		return "", false
	}
	adj, ok := scopedCreatureAdjective(cond.Cond)
	if !ok {
		return "", false
	}
	return "after " + indefinite(adj+" creature") + " reaps, " + cond.Then.Text(), true
}

// scopedCreatureAdjective returns the house-relative adjective a subject-scope
// condition names — "enemy" for ItIsEnemy, "friendly" for ItIsFriendly — so a
// board-wide creature trigger gated on one folds into "an enemy creature" or "a
// friendly creature" phrasing.
func scopedCreatureAdjective(c Condition) (string, bool) {
	switch c.(type) {
	case ItIsEnemy:
		return "enemy", true
	case ItIsFriendly:
		return "friendly", true
	}
	return "", false
}

// afterCreaturePlayedAdjacentText folds an "after a creature is played adjacent
// to <self>" reaction gated only on the played creature's trait — a
// Conditional{ItIsOfTrait} — into the natural "after a <Trait> creature is played
// adjacent to <self>, <then>" wording (Stilt-Kin's "after a Giant creature is
// played adjacent to Stilt-Kin, ready and fight with Stilt-Kin"), rather than the
// literal "after a creature is played adjacent to <self>, if it is a <Trait>
// creature, ...". Any other effect shape reports false and renders with the broad
// prefix.
func afterCreaturePlayedAdjacentText(e Effect) (string, bool) {
	cond, ok := e.(Conditional)
	if !ok {
		return "", false
	}
	trait, ok := cond.Cond.(ItIsOfTrait)
	if !ok {
		return "", false
	}
	return "after a " + trait.Trait.String() + " creature is played adjacent to " +
		SelfName + ", " + cond.Then.Text(), true
}

// afterEnemyPlaysCreatureOnFlankText folds an AfterEnemyCardPlayed reaction gated
// only on which flank the played creature landed on — a Conditional{OnFlank{OfIt}} —
// into "after your opponent plays a creature on their <side> flank, <then>" (Dexus
// right, Sinestra left), rather than the literal "after your opponent plays a card,
// if it is on the right flank, ...". Only a creature holds a flank, so the fused
// wording names the creature the position already requires. Any other effect shape
// reports false and renders with the broad prefix.
func afterEnemyPlaysCreatureOnFlankText(e Effect) (string, bool) {
	cond, ok := e.(Conditional)
	if !ok || cond.Else != nil {
		return "", false
	}
	fl, ok := cond.Cond.(OnFlank)
	if !ok || !fl.OfIt || fl.Where == AnyFlank {
		return "", false
	}
	return "after your opponent plays a creature on their " +
		fl.sideName() + " flank, " + cond.Then.Text(), true
}

// afterChooseHouseText folds an AfterChooseHouse ability, whose effect is a
// Conditional gated on the chosen house (a ChoseHouse condition), into the
// natural "after you choose <House> as your active house, <then>" wording (Jehu
// the Bureaucrat's "after you choose Sanctum as your active house, gain 2
// Æmber"). Any other effect shape reports false and renders with the ordinary
// prefix.
func afterChooseHouseText(e Effect) (string, bool) {
	cond, ok := e.(Conditional)
	if !ok {
		return "", false
	}
	ch, ok := cond.Cond.(ChoseHouse)
	if !ok {
		return "", false
	}
	return "after you choose " + ch.House.String() + " as your active house, " + cond.Then.Text(), true
}

// afterAnyPlayerChooseHouseText folds an AfterAnyPlayerChoosesHouse ability whose
// effect is a Conditional gated on the chosen house — a ChoseHouse condition, into
// "after a player chooses <House> as their active house, <then>" (the house
// plants' "after a player chooses Brobnar as their active house, gain 1 Æmber"),
// or an ActiveHouseMatchesNoCardsInPlay condition, into "after a player chooses an
// active house which matches no cards in play, <then>" (Sci. Officer Qincan). Any
// other effect shape reports false and renders with the ordinary prefix.
func afterAnyPlayerChooseHouseText(e Effect) (string, bool) {
	cond, ok := e.(Conditional)
	if !ok {
		return "", false
	}
	switch c := cond.Cond.(type) {
	case ChoseHouse:
		return "after a player chooses " + c.House.String() +
			" as their active house, " + cond.Then.Text(), true
	case ActiveHouseMatchesNoCardsInPlay:
		return "after a player chooses an active house " +
			c.CondText() + ", " + cond.Then.Text(), true
	}
	return "", false
}

// entersPlayConditionalText folds an "enters play" ability whose effect is gated
// on a board fact — a Conditional{Cond, Then: <state effect>} — into the natural
// "<cond>, <self> enters play ready" wording (Bramble Lynx's "if you have used a
// creature to reap this turn, Bramble Lynx enters play ready"), rather than the
// literal "if ..., ready this creature". Any other effect shape reports false and
// renders with the ordinary "<self> enters play <word>" form.
func entersPlayConditionalText(e Effect) (string, bool) {
	cond, ok := e.(Conditional)
	if !ok || cond.Else != nil {
		return "", false
	}
	return cond.Cond.CondText() + ", " + SelfName + " enters play " + enterStateWord(
		cond.Then,
	), true
}

// enterStateWord renders the state an "enters play" ability leaves its creature in,
// e.g. Stun -> "stunned", so the ability reads "<name> enters play stunned." A
// Sentences of state effects joins its words with "and" (Gizelhart's Zealot enters
// play ready and enraged). An effect without a dedicated enter word falls back to
// its ordinary text.
func enterStateWord(e Effect) string {
	switch ef := e.(type) {
	case Stun:
		return "stunned"
	case Ready:
		return "ready"
	case Enrage:
		return "enraged"
	case Sentences:
		words := make([]string, len(ef.Effects))
		for i, sub := range ef.Effects {
			words[i] = enterStateWord(sub)
		}
		return strings.Join(words, " and ")
	default:
		return e.Text()
	}
}

// abilityTextWithNames resolves the two placeholders an ability line may use: the
// host creature named by {self}, and the card's own name ({card}) for text that
// must name the card itself rather than the host it acts on.
func abilityTextWithNames(line, self, card string) string {
	line = strings.ReplaceAll(line, SelfName, self)
	return strings.ReplaceAll(line, CardName, card)
}

// abilityLines renders a card's triggered abilities, one printed line each,
// resolving self-references to the card's name (or "this creature" for an
// upgrade's own abilities). Adjacent abilities that share one effect and fire on
// distinct action triggers — the pairs/triples the composite triggers
// Trigger.FightReap, Trigger.PlayReap, and Trigger.PlayFightReap fan out into —
// merge into one "Fight/Reap:" / "Play/Reap:" / "Play/Fight/Reap:" line, the
// KeyForge shorthand.
func abilityLines(def *CardDefinition) []string {
	self := def.Name
	if def.Type == Upgrade {
		self = "this creature"
	}
	abs := def.Abilities
	lines := make([]string, 0, len(abs))
	for i := 0; i < len(abs); {
		if label, run := actionTriggerRun(abs, i); run > 1 {
			body := capitalizeFirst(abilityTextWithNames(abs[i].Effect.Text(), self, def.Name))
			lines = append(lines, label+": "+punctuate(body))
			i += run
			continue
		}
		lines = append(lines, abilityTextWithNames(RenderAbility(abs[i]), self, def.Name))
		i++
	}
	return lines
}

// actionTriggerRun finds the run of consecutive abilities starting at i that
// share one effect and each fire on a distinct action trigger (Play, Fight, or
// Reap). It returns the combined label ("Play/Fight/Reap") in canonical order and
// the run length; a run of one is left for the normal per-trigger rendering.
func actionTriggerRun(abs []Ability, i int) (string, int) {
	text := abs[i].Effect.Text()
	var play, fight, reap bool
	n := 0
	for j := i; j < len(abs) && abs[j].Effect.Text() == text; j++ {
		switch t := abs[j].Trigger; {
		case t == TriggerAfterPlay && !play:
			play = true
		case t == TriggerAfterFight && !fight:
			fight = true
		case t == TriggerAfterReap && !reap:
			reap = true
		default:
			// Stop: not an action trigger, or a repeated one.
			return canonicalTriggerLabel(play, fight, reap), n
		}
		n++
	}
	return canonicalTriggerLabel(play, fight, reap), n
}

// canonicalTriggerLabel joins the present action-trigger labels in the fixed
// order Play, Fight, Reap.
func canonicalTriggerLabel(play, fight, reap bool) string {
	var parts []string
	if play {
		parts = append(parts, "Play")
	}
	if fight {
		parts = append(parts, "Fight")
	}
	if reap {
		parts = append(parts, "Reap")
	}
	return strings.Join(parts, "/")
}

// isFightReapPair reports whether two adjacent abilities are a Fight and a Reap
// (in either order) that share one effect — the pair the FightReap granted
// helper and the Trigger.FightReap composite add. It prints as a single
// "Fight/Reap:" line regardless of which of the two is listed first.
func isFightReapPair(a, b Ability) bool {
	if a.Effect.Text() != b.Effect.Text() {
		return false
	}
	return (a.Trigger == TriggerAfterReap && b.Trigger == TriggerAfterFight) ||
		(a.Trigger == TriggerAfterFight && b.Trigger == TriggerAfterReap)
}

// RenderCardText renders a card's details as labeled, colon-aligned lines
// (House, Type, Rarity, stats, Æmber, Traits), followed by the card's rules text
// (keywords, upgrade modifier, and ability lines). Labels are padded by rune
// width so the multi-byte "Æmber" label still aligns.
func RenderCardText(def *CardDefinition) string {
	return renderCardText(def, false)
}

// RenderCardDetail is RenderCardText with the card's name as an initial
// "Name:" line, for a detail pane that shows a card on its own.
func RenderCardDetail(def *CardDefinition) string {
	return renderCardText(def, true)
}

func renderCardText(def *CardDefinition, withName bool) string {
	type field struct{ label, value string }
	var fields []field
	if withName {
		fields = append(fields, field{"Name", def.Name})
	}
	fields = append(fields,
		field{"House", def.House.String()},
		field{"Type", def.Type.String()},
		field{"Rarity", string(def.Rarity)},
	)
	if def.Type == Creature {
		fields = append(fields, field{"Power", fmt.Sprintf("%d", def.Power)})
		if def.Armor > 0 {
			fields = append(fields, field{"Armor", fmt.Sprintf("%d", def.Armor)})
		}
	}
	if len(def.Bonuses) > 0 {
		fields = append(fields, field{"Bonus", bonusIconsText(def.Bonuses)})
	}
	if len(def.Traits) > 0 {
		traits := make([]string, len(def.Traits))
		for i, t := range def.Traits {
			traits[i] = t.String()
		}
		fields = append(fields, field{"Traits", strings.Join(traits, " • ")})
	}

	// Widest label, measured in runes so "Æmber" (multi-byte Æ) aligns visually.
	width := 0
	for _, f := range fields {
		if n := len([]rune(f.label)); n > width {
			width = n
		}
	}

	var lines []string
	for _, f := range fields {
		pad := strings.Repeat(" ", width-len([]rune(f.label))+1)
		lines = append(lines, f.label+":"+pad+f.value)
	}

	// Rules text (keywords, upgrade modifier, abilities) follows the labeled
	// header, separated by a blank line.
	if rules := cardRules(def, false); len(rules) > 0 {
		lines = append(lines, "")
		lines = append(lines, rules...)
	}

	return strings.Join(lines, "\n")
}

// RenderCardRules renders just a card's rules text — the keyword, restriction,
// key-cost, static-modifier, constant-ability, granted-ability, and triggered-
// ability lines that RenderCardText shows below the labeled header — joined one
// per line with no header. It is the text drawn on the compact card face, so an
// upgrade shows its granted keywords and abilities there too.
func RenderCardRules(def *CardDefinition) string {
	return strings.Join(cardRules(def, false), "\n")
}

// RenderUpgradeOnCreature renders an Upgrade's rules as they read once it is
// attached. Printed on its own an Upgrade has to say who it is talking about —
// `This creature gains, "Reap: Steal 1 Æmber."` — but drawn on the host's face
// that creature is right there, so the framing is dropped and the line reads as
// if it were printed on the creature: `Reap: Steal 1 Æmber.` A creature played as
// an upgrade shows only what its Static grants the host, not its own creature
// keywords and abilities, which do not apply while it is an upgrade.
func RenderUpgradeOnCreature(def *CardDefinition) string {
	if def.PlayableAsUpgrade {
		return strings.Join(upgradeGrantLines(def, true), "\n")
	}
	return strings.Join(cardRules(def, true), "\n")
}

// upgradeGrantLines gathers what an Upgrade grants its host: its static modifier,
// its non-flank fight protection, and its granted abilities, in printed order.
// hosted drops the "this creature" framing for a face already showing the host.
func upgradeGrantLines(def *CardDefinition, hosted bool) []string {
	var lines []string
	if s := houseOverrideLine(def); s != "" {
		lines = append(lines, s)
	}
	lines = append(lines, upgradeStaticLines(def, hosted)...)
	if def.Static.ProtectsFromNonFlank {
		lines = append(lines,
			"Creatures not on a flank cannot fight this creature.")
	}
	// A house-override line already folds in the granted abilities.
	if def.Static.HouseOverride == HouseNone {
		lines = append(lines, grantedText(def.Static, def.Name, hosted)...)
	}
	return lines
}

// houseOverrideLine renders an Upgrade's house-override clause, folding in the
// abilities the Upgrade grants so they read as one sentence, or "" when the
// Upgrade overrides no house — `This creature belongs to Logos and this creature
// gains "Reap: Draw a card."` (Academy Training). Because it carries the grants,
// callers must not also print them through grantedText.
func houseOverrideLine(def *CardDefinition) string {
	m := def.Static
	if m.HouseOverride == HouseNone {
		return ""
	}
	line := "This creature belongs to " + m.HouseOverride.String()
	frame := func(body string) string { return `this creature gains "` + body + `"` }
	for _, g := range grantedLines(m, def.Name, frame) {
		line += " and " + g
	}
	return line
}

// playableAsUpgradeText renders the clause a creature played as an upgrade prints
// below its own text — `<name> may be played as an upgrade instead of a creature,
// with the text: "…"` — quoting what it grants a host. The granted text's own
// double quotes become single quotes so they nest inside the clause. Returns ""
// for a card that cannot be played as an upgrade.
func playableAsUpgradeText(def *CardDefinition) string {
	if !def.PlayableAsUpgrade {
		return ""
	}
	body := strings.ReplaceAll(strings.Join(upgradeGrantLines(def, false), " "), `"`, `'`)
	return def.Name +
		` may be played as an upgrade instead of a creature, with the text: "` +
		body + `"`
}

// cardRules assembles a card's rules lines in printed order: keywords, "cannot"
// restrictions, key-cost, an upgrade's static modifier, a constant ability, the
// abilities an upgrade grants its host, and finally the card's own triggered
// abilities. An upgrade's own abilities name their host "this creature" since it
// is unknown at print time; hosted drops that framing for a face already showing
// the creature (see RenderUpgradeOnCreature).
func cardRules(def *CardDefinition, hosted bool) []string {
	var rules []string
	if s := keywordText(def); s != "" {
		rules = append(rules, s)
	}
	if def.TauntReachesNeighborsNeighbors {
		rules = append(rules,
			def.Name+"'s taunt also applies to its neighbors' neighbors.")
	}
	if s := attackDamageText(def); s != "" {
		rules = append(rules, s)
	}
	if def.DealsNoDamageWhenAttacked {
		rules = append(rules, def.Name+" deals no damage when attacked.")
	}
	if n := def.StealsInsteadOfDamageWhenAttacked; n > 0 {
		rules = append(rules, fmt.Sprintf(
			"When %s would deal damage, steal %d Æmber instead.", def.Name, n))
	}
	if s := entersReadyText(def.EntersReadyGrant); s != "" {
		rules = append(rules, s)
	}
	if fr := def.FightRestriction; fr != (Target{}) {
		rules = append(rules, def.Name+" can only fight "+singularNoun(fr.Text())+"s.")
	}
	for _, k := range def.CannotBeUsedTo {
		rules = append(rules, def.Name+" cannot "+k.verb()+".")
	}
	if c := def.CannotBeUsedWhile; c != nil {
		rules = append(rules,
			"While "+trimCondPrefix(c.CondText())+", "+def.Name+" cannot be used.")
	}
	if dw := def.DestroyedWhen; dw != nil {
		cond := strings.ReplaceAll(dw.CondText(), SelfName, def.Name)
		rules = append(rules, capitalizeFirst(cond)+", destroy "+def.Name+".")
	}
	if t := def.TakesDamageFor; t.valid() {
		rules = append(rules,
			"Damage dealt to "+t.Text()+" is dealt to "+def.Name+" instead.")
	}
	if def.AlsoTakesNeighborFightDamage {
		rules = append(rules,
			"Damage dealt to "+def.Name+"'s neighbors during fights is also dealt to "+
				def.Name+".")
	}
	if m := def.CannotBeDealtDamageBy; m.Active() {
		rules = append(rules, def.Name+" cannot be dealt damage by "+m.clause()+".")
	}
	if px := def.PowerX; px != nil {
		rules = append(rules,
			strings.ReplaceAll("X is "+cardinalCountText(px)+".", SelfName, def.Name))
	}
	if s := attackIgnoresText(def); s != "" {
		rules = append(rules, s)
	}
	if s := attackKeywordsText(def); s != "" {
		rules = append(rules, s)
	}
	for _, r := range restrictionText(def.Restricts, def.Type == Upgrade) {
		rules = append(rules, strings.ReplaceAll(r, SelfName, def.Name))
	}
	if b := def.CannotPlayWhile; b.When != nil {
		rules = append(rules, conditionalPlayBarText(b))
	}
	if c := def.AemberCannotBeStolen; c != nil {
		if _, ok := c.(AlwaysMet); ok {
			rules = append(rules, "Your Æmber cannot be stolen.")
		} else {
			cond := strings.ReplaceAll(c.CondText(), SelfName, def.Name)
			rules = append(rules, "While "+trimCondPrefix(cond)+", your Æmber cannot be stolen.")
		}
	}
	if def.SpendableAember {
		rules = append(rules, "You may spend Æmber on "+def.Name+" when forging keys.")
	}
	if s := def.PlayRequirement.text(); s != "" {
		rules = append(rules, strings.ReplaceAll(s, SelfName, def.Name))
	}
	for _, kc := range def.KeyCostChanges {
		if s := keyCostText(kc); s != "" {
			rules = append(rules, strings.ReplaceAll(s, SelfName, def.Name))
		}
	}
	if s := def.HouseLock.text(); s != "" {
		rules = append(rules, strings.ReplaceAll(s, SelfName, def.Name))
	}
	if s := drawModifierText(def.DrawModifier); s != "" {
		rules = append(rules, strings.ReplaceAll(s, SelfName, def.Name))
	}
	if s := playPermissionText(def.PlayPermission); s != "" {
		rules = append(rules, strings.ReplaceAll(s, SelfName, def.Name))
	}
	if s := captureOpponentAemberText(def); s != "" {
		rules = append(rules, s)
	}
	if s := takeFromSupplyText(def); s != "" {
		rules = append(rules, s)
	}
	if s := captureStolenAemberText(def); s != "" {
		rules = append(rules, s)
	}
	if s := gainsForgeAemberText(def); s != "" {
		rules = append(rules, s)
	}
	if s := bonusInsteadText(def); s != "" {
		rules = append(rules, s)
	}
	// A creature played as an upgrade folds its Static grant into the "may be played
	// as an upgrade" clause below, so it does not also print as standalone lines.
	if !def.PlayableAsUpgrade {
		if s := houseOverrideLine(def); s != "" {
			rules = append(rules, s)
		}
		rules = append(rules, upgradeStaticLines(def, hosted)...)
		if def.Static.ProtectsFromNonFlank {
			rules = append(rules,
				"Creatures not on a flank cannot fight this creature.")
		}
	}
	if s := constantText(def); s != "" {
		rules = append(rules, s)
	}
	rules = append(rules, constantGrantedText(def)...)
	rules = append(rules, spendAsPoolLines(def, hosted)...)
	if !def.PlayableAsUpgrade && def.Static.HouseOverride == HouseNone {
		rules = append(rules, grantedText(def.Static, def.Name, hosted)...)
	}
	rules = append(rules, abilityLines(def)...)
	if s := playableAsUpgradeText(def); s != "" {
		rules = append(rules, s)
	}
	// Enhance is a deck-building note, so it prints on the last line of the text box.
	if len(def.Enhances) > 0 {
		rules = append(rules, "Enhance "+bonusIconsText(def.Enhances)+".")
	}
	return rules
}

// CardDocComment renders a card's details as a Go doc comment block, the form
// `mage generateComments` writes above each card's declaration. It is the card
// name, a blank separator, then RenderCardText's lines, each turned into a
// comment line: the title as "// <Name>", blanks as "//", and detail lines as
// tab-indented "//\t..." so godoc renders the labeled block preformatted. The
// result has no trailing newline.
func CardDocComment(def *CardDefinition) string {
	var b strings.Builder
	b.WriteString("// " + def.Name + "\n//\n")
	for _, line := range strings.Split(RenderCardText(def), "\n") {
		if line == "" {
			b.WriteString("//\n")
		} else {
			b.WriteString("//\t" + line + "\n")
		}
	}
	return strings.TrimRight(b.String(), "\n")
}

// drawModifierText renders a card's continuous change to a player's end-of-turn
// hand size, e.g. `Your hand size is 1 more.` or, scaled by a running count,
// `For each friendly Sin creature your hand size is 1 more.` Returns "" when the
// modifier is zero.
func drawModifierText(m DrawModifier) string {
	if m.Amount == 0 {
		return ""
	}
	word, n := "more", m.Amount
	if n < 0 {
		word, n = "less", -n
	}
	var subject string
	switch m.Player {
	case Controller:
		subject = "your hand size"
	case Opponent:
		subject = "your opponent's hand size"
	default: // EachPlayer
		subject = "each player's hand size"
	}
	body := fmt.Sprintf("%s is %d %s", subject, n, word)
	var s string
	if m.Per != nil {
		s = "For each " + m.Per.CountText() + " " + body + "."
	} else {
		s = capitalizeFirst(body) + "."
	}
	if m.OnlyWhileOffFlank {
		s = "While " + SelfName + " is not on a flank, " + strings.ToLower(s[:1]) + s[1:]
	}
	if m.OnlyWhileInCenter {
		s = "While " + SelfName + " is in the center of the battleline, " + strings.ToLower(
			s[:1],
		) + s[1:]
	}
	return s
}

// staticText renders an Upgrade's continuous modifier, e.g.
// "This creature gains +5 power."
func staticText(m StaticModifier) string {
	var lines []string
	if s := staticBonuses(m); s != "" {
		if m.Per != nil {
			s += " for each " + m.Per.perTargetText()
		}
		if m.WhileOnFlank {
			lines = append(lines, "While this creature is on a flank, it gains "+s+".")
		} else {
			lines = append(lines, "This creature gains "+s+".")
		}
	}
	if s := keywordGrantsText(m); s != "" {
		lines = append(lines, s)
	}
	for _, k := range m.CannotBeUsedTo {
		lines = append(lines, "This creature cannot "+k.verb()+".")
	}
	return strings.Join(lines, " ")
}

// keywordGrantsText renders the keywords an Upgrade grants to creatures around its
// host, one sentence per grant with the reach spelled out — "This creature and
// each of its neighbors gains elusive." for a grant reaching both. It is empty
// when the modifier grants no keywords by reach.
func keywordGrantsText(m StaticModifier) string {
	var lines []string
	for _, grant := range m.KeywordGrants {
		who := keywordGrantReach(grant)
		if who == "" || len(grant.Keywords) == 0 {
			continue
		}
		words := make([]string, len(grant.Keywords))
		for i, kw := range grant.Keywords {
			words[i] = strings.ToLower(kw.String())
		}
		lines = append(lines, who+" gains "+oxfordAnd(words)+".")
	}
	return strings.Join(lines, " ")
}

// keywordGrantReach names the creatures a KeywordGrant reaches, as the subject of
// its sentence. It is empty when the grant reaches nobody.
func keywordGrantReach(grant KeywordGrant) string {
	switch {
	case grant.Host && grant.Neighbors:
		return "This creature and each of its neighbors"
	case grant.Host:
		return "This creature"
	case grant.Neighbors:
		return "Each of this creature's neighbors"
	default:
		return ""
	}
}

// staticBonuses lists what an Upgrade's continuous modifier adds, without the
// "This creature gains" framing — e.g. "+5 power and elusive". Empty when the
// modifier adds nothing.
func staticBonuses(m StaticModifier) string {
	var parts []string
	if m.PowerBonus != 0 {
		parts = append(parts, fmt.Sprintf("%+d power", m.PowerBonus))
	}
	if m.ArmorBonus != 0 {
		parts = append(parts, fmt.Sprintf("%+d armor", m.ArmorBonus))
	}
	if m.AssaultBonus != 0 {
		parts = append(parts, fmt.Sprintf("%+d assault", m.AssaultBonus))
	}
	if m.HazardousBonus != 0 {
		parts = append(parts, fmt.Sprintf("%+d hazardous", m.HazardousBonus))
	}
	if m.SplashAttackBonus != 0 {
		parts = append(parts, fmt.Sprintf("%+d splash-attack", m.SplashAttackBonus))
	}
	for _, kw := range m.Keywords {
		parts = append(parts, strings.ToLower(kw.String()))
	}
	return oxfordAnd(parts)
}

// oxfordAnd joins parts into one clause: "a" alone, "a and b" for two, and an
// Oxford-comma list "a, b, and c" for three or more.
// upgradeStaticLines renders an Upgrade's continuous modifier and replacement
// text, combining them when both are printed on the same Upgrade. On a host's
// face (hosted) each stands on its own line, unframed.
func upgradeStaticLines(def *CardDefinition, hosted bool) []string {
	static := staticText(def.Static)
	replacement := destructionReplacementText(def)
	// A creature carrying its own destruction replacement (Reassembling Automaton)
	// states it plainly in its own voice — "If this creature would be destroyed,
	// instead …" — rather than through the "This creature gains, …" framing an
	// Upgrade uses to grant the replacement to its host.
	if replacement != "" && def.Type != Upgrade {
		return []string{capitalizeFirst(replacement) + "."}
	}
	if hosted {
		var lines []string
		if s := staticBonuses(def.Static); s != "" {
			if def.Static.WhileOnFlank {
				s = "while on a flank, " + s
			}
			lines = append(lines, capitalizeFirst(s)+".")
		}
		if s := keywordGrantsText(def.Static); s != "" {
			lines = append(lines, s)
		}
		for _, k := range def.Static.CannotBeUsedTo {
			lines = append(lines, "This creature cannot "+k.verb()+".")
		}
		if replacement != "" {
			lines = append(lines, capitalizeFirst(replacement)+".")
		}
		return lines
	}
	switch {
	case static != "" && replacement != "":
		return []string{strings.TrimSuffix(static, ".") + ` and, "` + replacement + `."`}
	case static != "":
		return []string{static}
	case replacement != "":
		return []string{`This creature gains, "` + replacement + `."`}
	default:
		return nil
	}
}

// destructionReplacementText renders a replacement for a creature being destroyed
// — an Upgrade granting it to its host, or a creature carrying its own
// (Reassembling Automaton) — naming the card that resolves the replacement. A
// conditional replacement folds its condition into the "would be destroyed" clause.
func destructionReplacementText(def *CardDefinition) string {
	r := def.Static.Replaces
	if !r.valid() || r.When != EventCreatureDestroyed {
		return ""
	}
	cond := ""
	if r.Cond != nil {
		cond = " and " + strings.TrimPrefix(r.Cond.CondText(), "if ")
	}
	return "If this creature would be destroyed" + cond + ", instead " + strings.ReplaceAll(
		r.With.Text(),
		SelfName,
		def.Name,
	)
}

// grantedText renders the triggered abilities an Upgrade grants its host,
// combining matching Reap/Fight pairs into the printed "Fight/Reap:" shorthand,
// e.g. `This creature gains, "Reap: Steal 1 Æmber."`. Self-references resolve to
// "this creature" since the host is unknown when the Upgrade prints; hosted drops
// the framing so the line reads as if printed on the creature.
func grantedText(m StaticModifier, upgrade string, hosted bool) []string {
	frame := func(body string) string {
		if hosted {
			return body
		}
		return `This creature gains, "` + body + `"`
	}
	return grantedLines(m, upgrade, frame)
}

// grantedLines renders the triggered abilities and static grants an Upgrade gives
// its host, applying frame to each so the caller controls the surrounding phrase
// (grantedText's "This creature gains, …" or the house-override line's lowercase
// "this creature gains …").
func grantedLines(m StaticModifier, upgrade string, frame func(string) string) []string {
	lines := make([]string, 0, len(m.Granted))
	for i := 0; i < len(m.Granted); i++ {
		ab := m.Granted[i]
		if i+1 < len(m.Granted) && isFightReapPair(ab, m.Granted[i+1]) {
			body := capitalizeFirst(
				abilityTextWithNames(ab.Effect.Text(), "this creature", upgrade),
			)
			lines = append(lines, frame(`Fight/Reap: `+body+`.`))
			i++ // the partner prints as part of this line
			continue
		}
		body := abilityTextWithNames(RenderAbility(ab), "this creature", upgrade)
		lines = append(lines, frame(body))
	}
	if s := keyCostText(m.KeyCostChange); s != "" {
		lines = append(lines, frame(s))
	}
	if c := m.AemberCannotBeStolen; c != nil {
		if _, ok := c.(AlwaysMet); ok {
			lines = append(lines, frame("Your Æmber cannot be stolen."))
		} else {
			lines = append(lines, frame(
				"While "+trimCondPrefix(c.CondText())+", your Æmber cannot be stolen."))
		}
	}
	return lines
}

// constantText renders a card's constant ability, e.g. "Each friendly creature
// gains +1 power." or "Each neighboring creature gains +2 armor." The subject is
// the constant ability's Target ("each creature" when unset). Returns "" when the
// card has no constant ability.
func constantText(def *CardDefinition) string {
	if len(def.ConstantAbilities) == 0 {
		return ""
	}
	var lines []string
	for _, c := range def.ConstantAbilities {
		who := capitalizeFirst(c.target().Text())
		for _, t := range c.DisableTriggers {
			lines = append(lines, t.String()+" effects cannot trigger.")
		}
		if c.BlankText {
			lines = append(lines, who+"'s text box is considered blank (except for traits).")
		}
		if c.RemovesTraits {
			lines = append(lines, who+" loses each of its traits.")
		}
		if c.SelectiveArchivePickup {
			lines = append(
				lines,
				"Instead of picking up all of your archives, you may pick up any number of cards in your archives.",
			)
		}
		for _, k := range c.CannotBeUsedTo {
			line := who + " cannot " + k.verb() + "."
			lines = append(lines, strings.ReplaceAll(line, SelfName, def.Name))
		}
		for _, m := range c.AlsoTriggers {
			from, onto := triggerEffectNoun(m.From), triggerEffectNoun(m.Onto)
			line := who + "'s " + from + " effect is a " + from + "/" + onto + " effect."
			lines = append(lines, strings.ReplaceAll(line, SelfName, def.Name))
		}
		var parts []string
		if c.PowerBonus != 0 {
			parts = append(parts, fmt.Sprintf("%+d power", c.PowerBonus))
		}
		if c.ArmorBonus != 0 {
			parts = append(parts, fmt.Sprintf("%+d armor", c.ArmorBonus))
		}
		if c.HazardousBonus != 0 {
			parts = append(parts, fmt.Sprintf("hazardous %d", c.HazardousBonus))
		}
		if c.AssaultBonus != 0 {
			parts = append(parts, fmt.Sprintf("assault %d", c.AssaultBonus))
		}
		for _, k := range c.Keywords {
			parts = append(parts, strings.ToLower(k.String()))
		}
		if len(parts) == 0 {
			continue
		}
		line := who + " gains " + oxfordAnd(parts)
		if c.WhileInCenter {
			line = "While " + SelfName + " is in the center of your battleline, " +
				c.target().Text() + " gains " + oxfordAnd(parts)
		}
		if c.WhileCondition != nil {
			line = "While " + trimCondPrefix(c.WhileCondition.CondText()) + ", " +
				SelfName + " gains " + oxfordAnd(parts)
		}
		if c.Per != nil {
			// A count read from the source names it, unless the buffed creature is the
			// source itself, where a trailing "it" reads unambiguously (Centurion
			// Stenopius "for each Æmber on it").
			noun := c.Per.CountText()
			if c.target().Kind != TargetThisCreature {
				noun = countLeadText(c.Per)
			}
			line += " for each " + noun
		}
		if c.PerTarget != nil {
			line += " for each " + c.PerTarget.perTargetText()
		}
		if tgt := c.target(); tgt.Kind == TargetThisCreature && tgt.onFlank {
			line += " while it is on a flank"
		}
		if tgt := c.target(); tgt.Kind == TargetThisCreature && tgt.damaged {
			line += " while it is damaged"
		}
		if c.WhileOffFlank {
			line += " while it is not on a flank"
		}
		line += "."
		lines = append(lines, strings.ReplaceAll(line, SelfName, def.Name))
	}
	return strings.Join(lines, "\n")
}

// constantGrantedText renders the triggered abilities a card's constant ability
// grants the creatures it reaches, one line each, e.g.
// `Each creature gains, "Destroyed: purge this creature."`. Self-references in the
// granted ability resolve to "this creature", the creature that gains it.
func constantGrantedText(def *CardDefinition) []string {
	if len(def.ConstantAbilities) == 0 {
		return nil
	}
	var lines []string
	for _, c := range def.ConstantAbilities {
		if len(c.Granted) == 0 {
			continue
		}
		subject := capitalizeFirst(strings.ReplaceAll(c.target().Text(), SelfName, def.Name))
		for _, ab := range c.Granted {
			body := abilityTextWithNames(RenderAbility(ab), "this creature", def.Name)
			if c.WhileInCenter {
				lines = append(lines, "While "+def.Name+
					" is in the center of your battleline, it gains, \""+body+`"`)
				continue
			}
			if c.WhileCondition != nil {
				lines = append(lines, "While "+trimCondPrefix(c.WhileCondition.CondText())+
					", "+def.Name+` gains, "`+body+`"`)
				continue
			}
			lines = append(lines, subject+` gains, "`+body+`"`)
		}
	}
	return lines
}

// conditionalPlayBarText renders a symmetric, board-wide play bar as one line,
// e.g. "If a player has more creatures in play than their opponent, they cannot
// play creatures." (Quixxle Stone). It phrases the condition in the third person
// because the bar names whichever player it applies to, not the controller.
func conditionalPlayBarText(b ConditionalPlayBar) string {
	return "If a player " + symmetricCondText(b.When) +
		", they cannot play " + typeWord(b.Type) + "s."
}

// symmetricCondText renders a condition in the third-person, board-wide voice a
// ConditionalPlayBar needs ("a player ... they"), rather than the controller voice
// of CondText ("if you ..."). A condition whose two voices differ by more than the
// lead word renders its own, so the wording stays with the condition instead of
// being type-switched on here.
func symmetricCondText(c Condition) string {
	if s, ok := c.(symmetricCondTexter); ok {
		return s.symmetricCondText()
	}
	return strings.TrimPrefix(c.CondText(), "if ")
}

// symmetricCondTexter is the optional capability a Condition implements when its
// board-wide third-person wording is not its CondText with "if " trimmed.
type symmetricCondTexter interface {
	symmetricCondText() string
}

// trimCondPrefix drops the leading "if " a CondText opens with, so a condition can
// be reused after a different lead word — "While " + trimCondPrefix(...) reads
// "While your red key is forged".
func trimCondPrefix(s string) string {
	return strings.TrimPrefix(s, "if ")
}

// restrictionText renders a card's constant "cannot" rules, one line each, e.g.
// "You cannot play creatures." Returns nil when the card imposes none. isUpgrade
// phrases a use-condition against the host creature (Earthbind) rather than the
// card itself (Giant Sloth).
func restrictionText(r Restrictions, isUpgrade bool) []string {
	var lines []string
	if c := r.UseCondition; c != nil {
		cond := strings.TrimPrefix(c.CondText(), "if ")
		if isUpgrade {
			lines = append(lines, "This creature cannot be used unless "+cond+".")
		} else {
			lines = append(lines, "You cannot use this card unless "+cond+".")
		}
	}
	if r.Fighting {
		lines = append(lines, "You cannot use creatures to fight.")
	}
	switch r.Reaping {
	case Controller:
		lines = append(lines, "Your creatures cannot reap.")
	case Opponent:
		lines = append(lines, "Enemy creatures cannot reap.")
	case EachPlayer:
		lines = append(lines, "Creatures cannot reap.")
	}
	switch r.BonusIcons {
	case Controller:
		lines = append(lines, "You cannot resolve bonus icons on cards you play.")
	case Opponent:
		lines = append(lines, "Your opponent cannot resolve bonus icons on cards they play.")
	case EachPlayer:
		lines = append(lines, "Players cannot resolve bonus icons on cards they play.")
	}
	if r.CannotPlay != TypeUnset {
		lines = append(lines, "You cannot play "+strings.ToLower(r.CannotPlay.String())+"s.")
	}
	if l := r.PlayCardLimit; l.Amount > 0 {
		who := "You"
		switch l.Player {
		case Opponent:
			who = "Your opponent"
		case EachPlayer:
			who = "Each player"
		}
		lines = append(
			lines,
			fmt.Sprintf("%s cannot play more than %d cards each turn.", who, l.Amount),
		)
	}
	if t := r.Toll; t.Amount > 0 {
		lines = append(
			lines,
			fmt.Sprintf(
				"In order to %s, your opponent must give you %d Æmber.",
				t.Action.phrase(),
				t.Amount,
			),
		)
	}
	if r.SkipForge {
		lines = append(lines, `You skip your "forge a key" phase.`)
	}
	if n := r.NoForgeKeyNumber; n > 0 {
		lines = append(lines, fmt.Sprintf("Players cannot forge their %s key.", ordinalWord(n)))
	}
	if r.NoForgeWhileAheadOnKeys {
		lines = append(lines,
			"Each player cannot forge keys while they have more forged keys than their opponent.")
	}
	if r.MustFightIfAble {
		lines = append(lines, "Creatures must fight when used, if able.")
	}
	return lines
}

// keyCostText renders a card's key-cost change, e.g. "Your opponent's keys cost +1
// Æmber." (or "Your keys cost…" / "Each player's keys cost…"). Returns "" when the
// change is zero. Both a card that prints the rule and an Upgrade that grants it
// use this same sentence.
func keyCostText(kc KeyCostChange) string {
	if kc.amount == 0 {
		return ""
	}
	whose := "your"
	switch kc.player {
	case Opponent:
		whose = "your opponent's"
	case EachPlayer:
		whose = "each player's"
	}
	sentence := fmt.Sprintf("%s keys cost %+d Æmber", whose, kc.amount)
	if kc.per != nil {
		// A self-referential count reads forward as a leading, named clause
		// (Nyzyk Resonator); an aggregate count trails the sentence.
		if _, ok := kc.per.(leadingCounter); ok {
			return capitalizeFirst(forEach(kc.per, sentence)) + "."
		}
		sentence += " for each " + kc.per.CountText()
	}
	if kc.whileOnFlank {
		return "While " + SelfName + " is on a flank, " + sentence + "."
	}
	if kc.whileOffFlank {
		return "While " + SelfName + " is not on a flank, " + sentence + "."
	}
	if kc.whileCondition != nil {
		// A .While gate reads as a standing "While ..." clause, so the condition's
		// own "if " lead is swapped for it (Proclamation 346E).
		return "While " + trimCondPrefix(kc.whileCondition.CondText()) + ", " + sentence + "."
	}
	return capitalizeFirst(sentence) + "."
}

// playPermissionText renders a continuous permission to play cards of a house
// while that house is not active, e.g. Witch of the Wilds. The rendered text
// drops the "not your active house" qualifier — a player can already play cards
// of their active house — and reads as the plain standing permission it is. The
// NonActive form (Captain Val Jericho) instead grants cards of any non-active
// house, gated by an optional condition that names the source via SelfName.
func playPermissionText(p PlayPermission) string {
	if !p.granted() {
		return ""
	}
	if p.Types != 0 {
		return fmt.Sprintf(
			"You may play %s as if they were in the active house.",
			p.Types.playablePlural(),
		)
	}
	if p.NonActive {
		noun := "card"
		if p.count() != 1 {
			noun += "s"
		}
		if p.Condition != nil {
			return fmt.Sprintf(
				"During your turn, %s, you may play %s %s that is not of the active house.",
				p.Condition.CondText(),
				countWord(p.count()),
				noun,
			)
		}
		return fmt.Sprintf(
			"During your turn you may play %s %s that is not of the active house.",
			countWord(p.count()),
			noun,
		)
	}
	noun := p.House.String() + " card"
	if p.count() != 1 {
		noun += "s"
	}
	return fmt.Sprintf(
		"Each turn you may play %s %s.",
		countWord(p.count()),
		noun,
	)
}

// captureOpponentAemberText renders a continuous replacement that captures Æmber
// added to a pool, e.g. "If Æmber would be added to your opponent's pool, instead
// Ether Spider captures it."
func captureOpponentAemberText(def *CardDefinition) string {
	r := def.Replaces
	if r.Of != EventAemberAddedToPool || r.With != Capture {
		return ""
	}
	whose := "your"
	if r.Player == Opponent {
		whose = "your opponent's"
	}
	return "If Æmber would be added to " + whose + " pool, instead " + def.Name + " captures it."
}

// takeFromSupplyText renders a continuous replacement that draws a steal or capture
// from the common supply instead of the target's pool, e.g. "Æmber stolen or
// captured from your pool is taken from the common supply instead." (Po's Pixies).
func takeFromSupplyText(def *CardDefinition) string {
	r := def.Replaces
	if r.Of != EventAemberTakenFromPool || r.With != FromCommonSupply {
		return ""
	}
	return "Æmber stolen or captured from your pool is taken from the common supply instead."
}

// captureStolenAemberText renders a continuous replacement that captures every
// stolen Æmber onto a creature the active player controls, e.g. "Each Æmber that
// would be stolen is captured by a creature controlled by the active player
// instead." (Gargantodon).
func captureStolenAemberText(def *CardDefinition) string {
	r := def.Replaces
	if r.Of != EventAemberStolen || r.With != Capture {
		return ""
	}
	return "Each Æmber that would be stolen is captured by a creature controlled by the active player instead."
}

// bonusInsteadText renders a card's continuous bonus-icon substitution, e.g. "When
// resolving a bonus icon, you may resolve it as a Capture bonus icon instead."
// (Amphora Captura) or "When you resolve a Capture bonus icon, steal 1 Æmber
// instead." (Scrivener Favian). A May substitution reads "you may"; a mandatory one
// omits it.
func bonusInsteadText(def *CardDefinition) string {
	rm := def.BonusInstead
	if !rm.set() {
		return ""
	}
	trigger := "When resolving a bonus icon, "
	if rm.From != bonusUnset {
		trigger = "When you resolve a " + rm.From.String() + " bonus icon, "
	}
	clause := bonusInsteadClause(rm)
	if rm.May {
		clause = "you may " + clause
	}
	return trigger + clause + " instead."
}

// bonusInsteadClause renders the verb phrase a bonus-icon substitution applies:
// "resolve it as a Capture bonus icon" for an icon swap, or the effect's own text
// ("steal 1 Æmber") for a replacement effect.
func bonusInsteadClause(rm BonusInstead) string {
	if rm.Instead != nil {
		return lowerFirst(rm.Instead.Text())
	}
	return "resolve it as a " + rm.As.String() + " bonus icon"
}

// gainsForgeAemberText renders a card that gains all the Æmber its controller's
// opponent spends forging a key, e.g. "You gain all Æmber your opponent spends
// when forging a key."
func gainsForgeAemberText(def *CardDefinition) string {
	if !def.GainsForgeAember {
		return ""
	}
	return "You gain all Æmber your opponent spends when forging a key."
}

// keywordText renders a card's keywords as a single leading sentence, e.g.
// "Skirmish, Poison." or "Assault 2.". Returns "" when the card has none.
func keywordText(def *CardDefinition) string {
	var parts []string
	for _, k := range def.Keywords {
		parts = append(parts, k.String())
	}
	if def.Assault > 0 {
		parts = append(parts, fmt.Sprintf("Assault %d", def.Assault))
	}
	if def.Hazardous > 0 {
		parts = append(parts, fmt.Sprintf("Hazardous %d", def.Hazardous))
	}
	if def.SplashAttack > 0 {
		parts = append(parts, fmt.Sprintf("Splash-attack %d", def.SplashAttack))
	}
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, ", ") + "."
}

// attackDamageText renders a creature's custom fight damage, e.g. "Valdr deals +2
// damage while attacking an enemy creature on the flank." or "Ether Spider deals
// no damage when fighting." Returns "" when the creature deals its plain power.
func attackDamageText(def *CardDefinition) string {
	ad := def.AttackDamage
	switch {
	case ad.Fixed && ad.Amount == 0:
		return def.Name + " deals no damage when fighting."
	case ad.Fixed:
		return fmt.Sprintf("%s deals %s when fighting.", def.Name, damageAmount(ad.Amount))
	case ad.Amount != 0 && ad.FlankOnly:
		return fmt.Sprintf(
			"%s deals +%d damage while attacking an enemy creature on the flank.",
			def.Name,
			ad.Amount,
		)
	case ad.Amount != 0:
		return fmt.Sprintf("%s deals +%d damage when fighting.", def.Name, ad.Amount)
	default:
		return ""
	}
}

// entersReadyText renders the "enter play ready" grant line, matching the printed
// wording: an ungated creature grant reads "Your creatures enter play ready.", an
// ungated artifact grant "Friendly artifacts enter play ready.", and a gated or
// house-filtered grant leads with its "While you have N or more Æmber" clause and
// a "non-<House>" qualifier (Fandangle).
func entersReadyText(g EntersReadyGrant) string {
	var noun string
	switch g.Type {
	case Creature:
		noun = "creatures"
	case Artifact:
		noun = "artifacts"
	default:
		return ""
	}
	// Preserve the established simple wording when there is no gate or filter.
	if g.MinAember == 0 && g.ExceptHouse == HouseNone {
		if g.Type == Artifact {
			return "Friendly artifacts enter play ready."
		}
		return "Your creatures enter play ready."
	}
	qualifier := ""
	if g.ExceptHouse != HouseNone {
		qualifier = "non-" + g.ExceptHouse.String() + " "
	}
	if g.MinAember > 0 {
		return fmt.Sprintf("While you have %d or more Æmber, your %s%s enter play ready.",
			g.MinAember, qualifier, noun)
	}
	return "Your " + qualifier + noun + " enter play ready."
}

// attackIgnoresText renders the defensive keywords a creature ignores while it is
// attacking, e.g. "While Niffle Ape is attacking, ignore taunt and elusive."
func attackIgnoresText(def *CardDefinition) string {
	if len(def.AttackIgnores) == 0 {
		return ""
	}
	words := make([]string, len(def.AttackIgnores))
	for i, kw := range def.AttackIgnores {
		words[i] = strings.ToLower(kw.String())
	}
	return fmt.Sprintf(
		"While %s is attacking, ignore %s.",
		def.Name,
		strings.Join(words, " and "),
	)
}

// attackKeywordsText renders the keywords a creature gains while attacking, e.g.
// "Spyyyder gains poison while attacking an enemy flank creature." Returns "" when
// the creature gains none.
func attackKeywordsText(def *CardDefinition) string {
	ak := def.AttackKeywords
	if len(ak.Keywords) == 0 {
		return ""
	}
	words := make([]string, len(ak.Keywords))
	for i, kw := range ak.Keywords {
		words[i] = strings.ToLower(kw.String())
	}
	if ak.FlankOnly {
		return fmt.Sprintf(
			"%s gains %s while attacking an enemy flank creature.",
			def.Name,
			strings.Join(words, " and "),
		)
	}
	return fmt.Sprintf(
		"%s gains %s while attacking.",
		def.Name,
		strings.Join(words, " and "),
	)
}
