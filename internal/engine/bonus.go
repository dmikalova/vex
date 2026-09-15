package engine

import "strings"

// BonusIcon is a printed bonus icon in a card's upper-left corner. After a card is
// played its icons resolve one at a time, top to bottom, before its "Play:"
// abilities (the game's ordinary bonus-icon step). The four modelled here are the
// icons Mass Mutation distributes with Enhance; Discard, House, and the +1-power
// icon are not modelled yet.
type BonusIcon uint8

const (
	// bonusUnset is the invalid zero value: a bonus icon must name its kind.
	bonusUnset BonusIcon = iota
	// BonusAember gains 1 Æmber from the common supply.
	BonusAember
	// BonusCapture has a friendly creature capture 1 Æmber from the opponent.
	BonusCapture
	// BonusDamage deals 1 damage to a creature in play.
	BonusDamage
	// BonusDraw draws 1 card.
	BonusDraw
	// bonusIconCount bounds the enum for validation.
	bonusIconCount
)

// String names the icon, e.g. "Æmber".
func (b BonusIcon) String() string {
	switch b {
	case BonusAember:
		return "Æmber"
	case BonusCapture:
		return "Capture"
	case BonusDamage:
		return "Damage"
	case BonusDraw:
		return "Draw"
	default:
		return "Unset"
	}
}

// valid reports whether the icon names a real kind.
func (b BonusIcon) valid() bool { return b > bonusUnset && b < bonusIconCount }

// allBonusIcons lists every bonus-icon kind, the closed catalog the rulebook
// completeness test ranges over (ADR 0018).
func allBonusIcons() []BonusIcon {
	out := make([]BonusIcon, 0, int(bonusIconCount)-1)
	for b := bonusUnset + 1; b < bonusIconCount; b++ {
		out = append(out, b)
	}
	return out
}

// countBonus returns how many icons of one kind sit in a list.
func countBonus(icons []BonusIcon, kind BonusIcon) int {
	n := 0
	for _, ic := range icons {
		if ic == kind {
			n++
		}
	}
	return n
}

// AemberBonus reports how many Æmber bonus icons the card prints, the count the
// client shows as the card's Æmber pips.
func (d *CardDefinition) AemberBonus() int { return countBonus(d.Bonuses, BonusAember) }

// bonusIconsText renders an ordered icon list as space-separated names, e.g.
// "Æmber Æmber Draw", for a card's generated text.
func bonusIconsText(icons []BonusIcon) string {
	names := make([]string, len(icons))
	for i, ic := range icons {
		names[i] = ic.String()
	}
	return strings.Join(names, " ")
}
