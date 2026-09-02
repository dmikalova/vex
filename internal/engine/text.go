package engine

import (
	"fmt"
	"strings"
	"unicode"
)

// RenderAbility renders a single triggered ability to its printed card line,
// e.g. "After you forge a key, deal 2 damage to each enemy creature."
func RenderAbility(a Ability) string {
	if a.Trigger == TriggerAfterCardPlayed {
		if s, ok := afterYouPlayText(a.Effect); ok {
			return punctuate(capitalizeFirst(s))
		}
	}
	if a.Trigger == TriggerAfterChooseHouse {
		if s, ok := afterChooseHouseText(a.Effect); ok {
			return punctuate(capitalizeFirst(s))
		}
	}
	if a.Trigger == TriggerEntersPlay {
		return SelfName + " enters play " + enterStateWord(a.Effect) + "."
	}
	prefix, capitalize := a.Trigger.prefix()
	body := a.Effect.Text()
	if capitalize {
		body = capitalizeFirst(body)
	}
	return prefix + punctuate(body)
}

// afterYouPlayText folds an "after you play a card" ability whose effect only
// gates on the played card's shape — a Conditional{ItIs} — into the natural
// "after you play a <shape>, <then>" wording (Carlo Phantom's "after you play an
// artifact, steal 1 Æmber"), rather than the literal "after you play a card, if
// it is a <shape>, ...". Any other effect renders with the broad prefix, so an
// unconditional or state-gated reaction still reads "After you play a card, ...".
func afterYouPlayText(e Effect) (string, bool) {
	cond, ok := e.(Conditional)
	if !ok {
		return "", false
	}
	it, ok := cond.Cond.(ItIs)
	if !ok {
		return "", false
	}
	return "after you play " + indefinite(
		houseTypeNoun(it.House, it.Type),
	) + ", " + cond.Then.Text(), true
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

// enterStateWord renders the state an "enters play" ability leaves its creature in,
// e.g. Stun -> "stunned", so the ability reads "<name> enters play stunned." An
// effect without a dedicated enter word falls back to its ordinary text.
func enterStateWord(e Effect) string {
	switch e.(type) {
	case Stun:
		return "stunned"
	case Ready:
		return "ready"
	default:
		return e.Text()
	}
}

// punctuate ends an ability body with a period. A body that already ends in a
// period is left alone; one that ends in a closing quote (an embedded ability
// such as Charge!'s granted "Play: ...") takes its period inside the quote, so
// the line reads `... an enemy creature."` rather than doubling or misplacing it.
func punctuate(body string) string {
	if strings.HasSuffix(body, `"`) {
		if inner := body[:len(body)-1]; !strings.HasSuffix(inner, ".") {
			return inner + `."`
		}
		return body
	}
	if strings.HasSuffix(body, ".") {
		return body
	}
	return body + "."
}

// renderAbilityLine renders an ability with its source card's self-references
// (the SelfName placeholder) resolved to the card's name.
func renderAbilityLine(def *CardDefinition, a Ability) string {
	return abilityTextWithNames(RenderAbility(a), def.Name, def.Name)
}

// abilityTextWithNames resolves the two placeholders an ability line may use: the
// card/creature named by "this", and, for an Upgrade resolving on its host, the
// Upgrade's own name.
func abilityTextWithNames(line, self, upgrade string) string {
	line = strings.ReplaceAll(line, SelfName, self)
	return strings.ReplaceAll(line, UpgradeName, upgrade)
}

// abilityLines renders a card's triggered abilities, one printed line each,
// resolving self-references to the card's name (or "this creature" for an
// upgrade's own abilities). Adjacent abilities that share one effect and fire on
// distinct action triggers — the pairs/triples WithFightOrReap, WithPlayReap, and
// WithPlayFightReap add — merge into one "Fight/Reap:" / "Play/Reap:" /
// "Play/Fight/Reap:" line, the KeyForge shorthand.
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
// (in either order) that share one effect — the pair FightOrReap and
// WithFightOrReap add. It prints as a single "Fight/Reap:" line regardless of
// which of the two is listed first.
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
	if def.AemberBonus > 0 {
		fields = append(fields, field{"Æmber", fmt.Sprintf("%d", def.AemberBonus)})
	}
	if len(def.Traits) > 0 {
		traits := make([]string, len(def.Traits))
		for i, t := range def.Traits {
			traits[i] = string(t)
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
// if it were printed on the creature: `Reap: Steal 1 Æmber.`
func RenderUpgradeOnCreature(def *CardDefinition) string {
	return strings.Join(cardRules(def, true), "\n")
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
	if s := attackDamageText(def); s != "" {
		rules = append(rules, s)
	}
	if fr := def.FightRestriction; fr != (Target{}) {
		rules = append(rules, def.Name+" can only fight "+singularNoun(fr.Text())+"s.")
	}
	if s := attackIgnoresText(def); s != "" {
		rules = append(rules, s)
	}
	rules = append(rules, restrictionText(def.Restricts)...)
	if def.PreventSteal {
		rules = append(rules, "Your Æmber cannot be stolen.")
	}
	if def.SpendableAember {
		rules = append(rules, "You may spend Æmber on "+def.Name+" when forging keys.")
	}
	if s := def.PlayRequirement.text(); s != "" {
		rules = append(rules, strings.ReplaceAll(s, SelfName, def.Name))
	}
	if s := keyCostText(def.KeyCostChange); s != "" {
		rules = append(rules, strings.ReplaceAll(s, SelfName, def.Name))
	}
	if s := def.HouseLock.text(); s != "" {
		rules = append(rules, strings.ReplaceAll(s, SelfName, def.Name))
	}
	if s := drawModifierText(def.DrawModifier); s != "" {
		rules = append(rules, s)
	}
	if s := playPermissionText(def.PlayPermission); s != "" {
		rules = append(rules, s)
	}
	if s := captureOpponentAemberText(def); s != "" {
		rules = append(rules, s)
	}
	rules = append(rules, upgradeStaticLines(def, hosted)...)
	if s := constantText(def); s != "" {
		rules = append(rules, s)
	}
	rules = append(rules, constantGrantedText(def)...)
	rules = append(rules, grantedText(def.Static, hosted)...)
	rules = append(rules, abilityLines(def)...)
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
// hand refill, e.g. `During your "draw cards" step, refill your hand to 1
// additional card.` Returns "" when the modifier is zero.
func drawModifierText(m DrawModifier) string {
	if m.Amount == 0 {
		return ""
	}
	word, n := "additional", m.Amount
	if n < 0 {
		word, n = "less", -n
	}
	noun := "card"
	if n != 1 {
		noun = "cards"
	}
	switch m.Player {
	case Controller:
		return fmt.Sprintf(
			"During your %q step, refill your hand to %d %s %s.",
			"draw cards",
			n,
			word,
			noun,
		)
	case Opponent:
		return fmt.Sprintf(
			"During their %q step, your opponent refills their hand to %d %s %s.",
			"draw cards",
			n,
			word,
			noun,
		)
	default: // EachPlayer
		return fmt.Sprintf(
			"During their %q step, each player refills their hand to %d %s %s.",
			"draw cards",
			n,
			word,
			noun,
		)
	}
}

// staticText renders an Upgrade's continuous modifier, e.g.
// "This creature gains +5 power."
func staticText(m StaticModifier) string {
	s := staticBonuses(m)
	if s == "" {
		return ""
	}
	if m.WhileOnFlank {
		return "While this creature is on a flank, it gains " + s + "."
	}
	return "This creature gains " + s + "."
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
	for _, kw := range m.Keywords {
		parts = append(parts, strings.ToLower(string(kw)))
	}
	return strings.Join(parts, " and ")
}

// upgradeStaticLines renders an Upgrade's continuous modifier and replacement
// text, combining them when both are printed on the same Upgrade. On a host's
// face (hosted) each stands on its own line, unframed.
func upgradeStaticLines(def *CardDefinition, hosted bool) []string {
	static := staticText(def.Static)
	replacement := destructionReplacementText(def)
	if hosted {
		var lines []string
		if s := staticBonuses(def.Static); s != "" {
			if def.Static.WhileOnFlank {
				s = "while on a flank, " + s
			}
			lines = append(lines, capitalizeFirst(s)+".")
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

// destructionReplacementText renders an Upgrade-granted replacement for its host
// being destroyed, naming the Upgrade that is destroyed instead.
func destructionReplacementText(def *CardDefinition) string {
	r := def.Static.Replaces
	if !r.valid() || r.When != EventCreatureDestroyed {
		return ""
	}
	return "If this creature would be destroyed, instead " + strings.ReplaceAll(
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
func grantedText(m StaticModifier, hosted bool) []string {
	frame := func(body string) string {
		if hosted {
			return body
		}
		return `This creature gains, "` + body + `"`
	}
	lines := make([]string, 0, len(m.Granted))
	for i := 0; i < len(m.Granted); i++ {
		ab := m.Granted[i]
		if i+1 < len(m.Granted) && isFightReapPair(ab, m.Granted[i+1]) {
			body := capitalizeFirst(
				abilityTextWithNames(ab.Effect.Text(), "this creature", "this upgrade"),
			)
			lines = append(lines, frame(`Fight/Reap: `+body+`.`))
			i++ // the partner prints as part of this line
			continue
		}
		body := abilityTextWithNames(RenderAbility(ab), "this creature", "this upgrade")
		lines = append(lines, frame(body))
	}
	if s := keyCostText(m.KeyCostChange); s != "" {
		lines = append(lines, frame(s))
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
		var parts []string
		if c.PowerBonus != 0 {
			parts = append(parts, fmt.Sprintf("%+d power", c.PowerBonus))
		}
		if c.ArmorBonus != 0 {
			parts = append(parts, fmt.Sprintf("%+d armor", c.ArmorBonus))
		}
		for _, k := range c.Keywords {
			parts = append(parts, strings.ToLower(string(k)))
		}
		if len(parts) == 0 {
			continue
		}
		who := capitalizeFirst(c.target().Text())
		line := who + " gains " + strings.Join(parts, " and ")
		if c.Per != nil {
			line += " for each " + c.Per.CountText()
		}
		if tgt := c.target(); tgt.Kind == TargetThisCreature && tgt.onFlank {
			line += " while it is on a flank"
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
		subject := capitalizeFirst(c.target().Text())
		for _, ab := range c.Granted {
			body := abilityTextWithNames(RenderAbility(ab), "this creature", def.Name)
			lines = append(lines, subject+` gains, "`+body+`"`)
		}
	}
	return lines
}

// restrictionText renders a card's constant "cannot" rules, one line each, e.g.
// "You cannot play creatures." Returns nil when the card imposes none.
func restrictionText(r Restrictions) []string {
	var lines []string
	if c := r.UseCondition; c != nil {
		lines = append(
			lines,
			"You cannot use this card unless "+strings.TrimPrefix(c.CondText(), "if ")+".",
		)
	}
	if r.Fighting {
		lines = append(lines, "You cannot use creatures to fight.")
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
				"Your opponent must give you %d Æmber in order to %s.",
				t.Amount,
				t.Action.phrase(),
			),
		)
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
		sentence += " for each " + kc.per.CountText()
	}
	if kc.whileOnFlank {
		return "While " + SelfName + " is on a flank, " + sentence + "."
	}
	return capitalizeFirst(sentence) + "."
}

// playPermissionText renders a continuous permission to play cards of a house
// while that house is not active, e.g. Witch of the Wilds.
func playPermissionText(p PlayPermission) string {
	if !p.granted() {
		return ""
	}
	noun := p.House.String() + " card"
	if p.count() != 1 {
		noun += "s"
	}
	return fmt.Sprintf(
		"During each turn in which %s is not your active house, you may play %s %s.",
		p.House,
		countWord(p.count()),
		noun,
	)
}

// countWord renders a small count as an English word ("one") for the common
// single-card grant, falling back to the numeral for larger counts.
func countWord(n int) string {
	if n == 1 {
		return "one"
	}
	return fmt.Sprintf("%d", n)
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

// keywordText renders a card's keywords as a single leading line, e.g.
// "Skirmish. Poison." or "Assault 2.". Returns "" when the card has none.
func keywordText(def *CardDefinition) string {
	var parts []string
	for _, k := range def.Keywords {
		parts = append(parts, string(k)+".")
	}
	if def.Assault > 0 {
		parts = append(parts, fmt.Sprintf("Assault %d.", def.Assault))
	}
	if def.Hazardous > 0 {
		parts = append(parts, fmt.Sprintf("Hazardous %d.", def.Hazardous))
	}
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, " ")
}

// attackDamageText renders a creature's custom fight damage, e.g. "Valdr deals +2
// Damage while attacking an enemy creature on the flank." or "Ether Spider deals
// no damage when fighting." Returns "" when the creature deals its plain power.
func attackDamageText(def *CardDefinition) string {
	ad := def.AttackDamage
	switch {
	case ad.Fixed && ad.Amount == 0:
		return def.Name + " deals no damage when fighting."
	case ad.Fixed:
		return fmt.Sprintf("%s deals %d damage when fighting.", def.Name, ad.Amount)
	case ad.Amount != 0 && ad.FlankOnly:
		return fmt.Sprintf(
			"%s deals +%d Damage while attacking an enemy creature on the flank.",
			def.Name,
			ad.Amount,
		)
	case ad.Amount != 0:
		return fmt.Sprintf("%s deals +%d Damage when fighting.", def.Name, ad.Amount)
	default:
		return ""
	}
}

// attackIgnoresText renders the defensive keywords a creature ignores while it is
// attacking, e.g. "While Niffle Ape is attacking, ignore taunt and elusive."
func attackIgnoresText(def *CardDefinition) string {
	if len(def.AttackIgnores) == 0 {
		return ""
	}
	words := make([]string, len(def.AttackIgnores))
	for i, kw := range def.AttackIgnores {
		words[i] = strings.ToLower(string(kw))
	}
	return fmt.Sprintf(
		"While %s is attacking, ignore %s.",
		def.Name,
		strings.Join(words, " and "),
	)
}

// capitalizeFirst upper-cases the first rune of s.
func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// indefinite prefixes a noun with the indefinite article "a" or "an", choosing
// "an" before a word that starts with a vowel — e.g. "an Urchin", "a Knight".
func indefinite(noun string) string {
	if noun == "" {
		return noun
	}
	switch unicode.ToLower([]rune(noun)[0]) {
	case 'a', 'e', 'i', 'o', 'u':
		return "an " + noun
	default:
		return "a " + noun
	}
}
