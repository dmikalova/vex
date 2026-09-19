package engine

import (
	"fmt"
	"strings"
)

// This file holds the log entries that narrate combat and damage (ADR 0011):
// the fight itself, what stopped it, and what each point of damage actually did
// once armor and immunity had their say.

// FightCancelled narrates a fight that an effect stopped before it happened.
type FightCancelled struct{ Attacker LocalID }

// Text renders a fight an effect stopped before it happened.
func (e FightCancelled) Text(n Namer) string {
	return fmt.Sprintf("%s's fight does not occur", n.Name(e.Attacker))
}

// FightKeywords is the set of a combatant's keywords that change what a fight
// does, annotated onto it in the fight line.
type FightKeywords uint8

const (
	// FightSkirmish spares its attacker the fight's return damage.
	FightSkirmish FightKeywords = 1 << iota
	// FightElusive turns the fight's damage aside entirely.
	FightElusive
	// FightPoison destroys whatever this combatant's fight damage lands on.
	FightPoison
)

// fightKeywordWords lists every annotatable keyword with the word it reads as, in
// the order a fight line names them.
var fightKeywordWords = []struct {
	bit  FightKeywords
	word string
}{
	{FightSkirmish, "skirmish"},
	{FightElusive, "elusive"},
	{FightPoison, "poison"},
}

// annotate renders a combatant as "<name> (<n> power, <keyword>, …)". These are
// the log's only parentheses, and they are structure rather than an aside: the
// fight line carries both combatants' stats and keywords, so spelling them out
// would turn one line into three.
func (k FightKeywords) annotate(name string, power int) string {
	parts := []string{fmt.Sprintf("%d power", power)}
	for _, kw := range fightKeywordWords {
		if k&kw.bit != 0 {
			parts = append(parts, kw.word)
		}
	}
	return fmt.Sprintf("%s (%s)", name, strings.Join(parts, ", "))
}

// Fought narrates two creatures fighting: the power each brought to it, and the
// keywords each brings that change what the fight does. The keywords annotate the
// combatant that has them instead of getting lines of their own, because that
// attributes each one correctly — skirmish is the attacker's, elusive the
// defender's, poison either's — and every annotation is known when the line is
// written. Poison's lethality is not, since it depends on damage surviving armor,
// so it stays a separate outcome line (PoisonKills).
type Fought struct {
	Attacker         LocalID
	AttackerPower    int
	AttackerKeywords FightKeywords
	Defender         LocalID
	DefenderPower    int
	DefenderKeywords FightKeywords
}

// Text renders the fight, and what each creature brought to it.
func (e Fought) Text(n Namer) string {
	return fmt.Sprintf("%s fights %s",
		e.AttackerKeywords.annotate(n.Name(e.Attacker), e.AttackerPower),
		e.DefenderKeywords.annotate(n.Name(e.Defender), e.DefenderPower))
}

// PoisonKills narrates a poison creature's fight damage proving lethal: the
// creature it damaged is destroyed, however much power it had left.
type PoisonKills struct {
	Source LocalID
	Victim LocalID
}

// Text renders poison destroying a creature the poison combatant damaged.
func (e PoisonKills) Text(n Namer) string {
	return fmt.Sprintf("%s's poison is lethal to %s", n.Name(e.Source), n.Name(e.Victim))
}

// DamageRefused narrates damage that a creature could not be dealt at all.
type DamageRefused struct{ Creature LocalID }

// Text renders damage a creature could not be dealt at all.
func (e DamageRefused) Text(n Namer) string {
	return fmt.Sprintf("%s cannot be dealt damage", n.Name(e.Creature))
}

// ArmorAbsorbed narrates the part of an incoming hit that armor soaked up.
type ArmorAbsorbed struct {
	Creature LocalID
	Amount   int
}

// Text renders the part of an incoming hit that armor soaked up.
func (e ArmorAbsorbed) Text(n Namer) string {
	return fmt.Sprintf("%s's armor absorbs %d damage", n.Name(e.Creature), e.Amount)
}

// DamageTaken narrates damage that actually landed, and the total now on the
// creature.
type DamageTaken struct {
	Creature LocalID
	Amount   int
	Total    int
}

// Text renders the damage that landed, and the creature's new total. The total is
// stated only when it differs from the hit, since the first hit is its own total.
func (e DamageTaken) Text(n Namer) string {
	if e.Total == e.Amount {
		return fmt.Sprintf("%s takes %d damage", n.Name(e.Creature), e.Amount)
	}
	return fmt.Sprintf("%s takes %d damage and now has %d damage",
		n.Name(e.Creature), e.Amount, e.Total)
}

// AssaultDealt narrates the pre-fight Assault an attacker deals the creature it
// attacks. It reports the damage that landed, not the keyword's value: where the
// two differ, the ArmorAbsorbed line already narrates the difference.
type AssaultDealt struct {
	Source LocalID
	Amount int
	Target LocalID
}

// Text renders the Assault damage, naming the keyword as the kind of damage it
// is rather than inventing a verb, so it reads like every other damage line.
func (e AssaultDealt) Text(n Namer) string {
	return fmt.Sprintf("%s deals %d assault damage to %s",
		n.Name(e.Source), e.Amount, n.Name(e.Target))
}

// HazardousDealt narrates the pre-fight Hazardous a defender deals its attacker.
type HazardousDealt struct {
	Source LocalID
	Amount int
	Target LocalID
}

// Text renders the Hazardous damage, naming the keyword as the kind of damage it
// is, so the dealing creature stays the subject as in every other damage line.
func (e HazardousDealt) Text(n Namer) string {
	return fmt.Sprintf("%s deals %d hazardous damage to %s",
		n.Name(e.Source), e.Amount, n.Name(e.Target))
}

// AbilityDamageDealt narrates damage a card's ability dealt another creature,
// naming the card that dealt it — the ability-damage counterpart to AssaultDealt
// and HazardousDealt. The dealing card is the record's frame source, not a field
// on the entry. Damage a card deals itself is left to the passive DamageTaken
// line, which reads better than naming the card twice.
type AbilityDamageDealt struct {
	Amount int
	Target LocalID
}

// Text renders the ability damage, naming the card that dealt it as the frame's
// source; with no frame it degrades to the passive line naming only the target.
func (e AbilityDamageDealt) Text(n Namer) string {
	if s, ok := framedSource(n); ok {
		return fmt.Sprintf("%s deals %d damage to %s", s, e.Amount, n.Name(e.Target))
	}
	return fmt.Sprintf("%s takes %d damage", n.Name(e.Target), e.Amount)
}

// damageEntry chooses how a landed hit narrates: pre-fight Assault or Hazardous
// names its striking creature and keyword value, ability damage names the card
// that dealt it, and any other damage is the bare DamageTaken line with the
// creature's new total. dealt is the damage left after armor; total is the
// creature's damage after taking it.
func (t DamageTarget) damageEntry(target LocalID, dealt, total int) LogEntry {
	switch t.SourceKeyword {
	case assaultDamage:
		return AssaultDealt{
			Source: t.Source,
			Amount: dealt,
			Target: target,
		}
	case hazardousDamage:
		return HazardousDealt{
			Source: t.Source,
			Amount: dealt,
			Target: target,
		}
	case abilityDamage:
		if t.Source != target {
			return AbilityDamageDealt{
				Amount: dealt,
				Target: target,
			}
		}
		return DamageTaken{
			Creature: target,
			Amount:   dealt,
			Total:    total,
		}
	case bonusDamage:
		return BonusDamageDealt{
			Source: t.Source,
			Amount: dealt,
			Target: target,
		}
	default:
		return DamageTaken{
			Creature: target,
			Amount:   dealt,
			Total:    total,
		}
	}
}
