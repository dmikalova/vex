package engine

// Abilities rulebook terms (ADR 0018): each describes itself next to the code it
// governs; the completeness test fails the build if a member of the matching
// closed catalog has no term here.
func init() {
	registerRuleSectionIntro(
		SectionAbility,
		`An ability is a line of rules text on a card whose **trigger** says when it
resolves. When the trigger's condition occurs, the ability's effect (see
_Effects_) resolves. The triggers below cover every way an ability can fire.`,
	)
	registerRuleTerms([]RuleTerm{
		{
			Section:    SectionAbility,
			Title:      "Constant Ability",
			Definition: "A continuous rule a card applies while it stays in play, with no trigger of its own.",
			Body: `A constant ability is a continuous rule a card applies while it stays in play,
with no trigger of its own — "Each friendly creature gains +1 power", or a card
that grants every creature a keyword or a "Destroyed:" ability. Its effect
applies for as long as the source card remains in play and stops the moment it
leaves; applying it is not "using" the card and never exhausts it.`,
		},
		{
			Section:    SectionAbility,
			Title:      "Play",
			Definition: "An ability that resolves right after you play the card from your hand.",
			Body: `A Play ability resolves right after you play the card from your hand. On a
creature or artifact it fires as the card enters play; on an action it is the
card's one-shot effect.`,
		},
		{
			Section:    SectionAbility,
			Title:      "Reap",
			Definition: "An ability that resolves after you use a ready creature to reap.",
			Body: `A Reap ability resolves after you use a ready creature to reap. Reaping gains
you 1 Æmber and exhausts the creature; the ability resolves in addition.`,
		},
		{
			Section:    SectionAbility,
			Title:      "Fight",
			Definition: "An ability that resolves after a creature you used to fight deals and takes its combat damage.",
			Body: `A Fight ability resolves after a creature you used to fight has dealt and
taken its combat damage and any resulting destruction has been carried out.`,
		},
		{
			Section:    SectionAbility,
			Title:      "Action",
			Definition: "An ability you resolve by using the card directly, without reaping or fighting, which exhausts it.",
			Body: `An Action ability is one you resolve by using the card directly, without
reaping or fighting; using it this way exhausts the card.`,
		},
		{
			Section:    SectionAbility,
			Title:      "After You Forge a Key",
			Definition: "An ability that resolves after its controller forges a key.",
			Body:       `This ability resolves after its controller forges a key.`,
		},
		{
			Section:    SectionAbility,
			Title:      "After a Creature Enters Play",
			Definition: "An ability that resolves after any creature enters play, including the opponent's.",
			Body: `This ability resolves after any creature enters play, including creatures
your opponent plays.`,
		},
		{
			Section:    SectionAbility,
			Title:      "Destroyed",
			Definition: "An ability that resolves as the card is destroyed, before it reaches the discard pile.",
			Body: `A Destroyed ability resolves as the card is destroyed, before it reaches the
discard pile, so it can still act on the board it is leaving.`,
		},
		{
			Section:    SectionAbility,
			Title:      "Before Fight",
			Definition: "An ability that resolves when a creature is used to fight, before any combat damage.",
			Body: `A Before Fight ability resolves when a creature is used to fight, before any
combat damage is dealt.`,
		},
		{
			Section:    SectionAbility,
			Title:      "After a Creature Is Destroyed Fighting",
			Definition: "An ability that resolves on a creature that survives a fight in which the other combatant was destroyed.",
			Body: `This ability resolves on a creature that survives a fight in which the other
combatant was destroyed; the destroyed creature is the one referred to as
"it".`,
		},
		{
			Section:    SectionAbility,
			Title:      "After You Play a Card",
			Definition: "An ability that resolves after its controller plays a card from hand.",
			Body: `This ability resolves after its controller plays a card — a creature,
artifact, or action — from hand. Putting a card into play by another effect is
not "playing" it and does not fire this (that is TriggerAfterCreatureEnters).
AfterCardPlayed narrows it to a house and/or type.`,
		},
		{
			Section:    SectionAbility,
			Title:      "End of Turn",
			Definition: "An ability that resolves at the end of its controller's turn, after cards ready and they draw.",
			Body: `An End of Turn ability resolves during the end of its controller's turn,
after cards ready and the controller draws (Shaffles drains the opponent at
each turn's end).`,
		},
		{
			Section:    SectionAbility,
			Title:      "Start of Turn",
			Definition: "An ability that resolves at the start of its controller's turn, before they forge.",
			Body: `A Start of Turn ability resolves at the start of its controller's turn,
before they forge, so an ability that changes what a key costs still has time
to.`,
		},
		{
			Section:    SectionAbility,
			Title:      "After an Enemy Creature Is Destroyed",
			Definition: "An ability that resolves after an enemy creature is destroyed during its controller's turn.",
			Body: `This ability resolves after an enemy creature is destroyed during its
controller's turn (Pile of Skulls captures Æmber onto a friendly creature
whenever an enemy creature is destroyed on your turn).`,
		},
		{
			Section:    SectionAbility,
			Title:      "After Your Opponent Plays a Card",
			Definition: "An ability that resolves after your opponent plays a card.",
			Body: `This ability resolves after your opponent plays a card (Teliga gains its
controller Æmber whenever the opponent plays a card).`,
		},
		{
			Section:    SectionAbility,
			Title:      "After You Use a Card",
			Definition: "An ability that resolves after its controller reaps, fights, or fires an Action.",
			Body: `This ability resolves after its controller uses a card — reaps or fights with a
creature, or fires an "Action:" — with the used card as "it" (Veylan Analyst
gains Æmber whenever you use an artifact).`,
		},
		{
			Section:    SectionAbility,
			Title:      "After You Discard a Card From Your Hand",
			Definition: "An ability that resolves after its controller discards a card from their hand.",
			Body: `This ability resolves after its controller discards a card from their hand,
with the discarded card as "it" (Rock-Hurling Giant). Discarding from anywhere
else — the top of a deck, the archives — is not this.`,
		},
		{
			Section:    SectionAbility,
			Title:      "After This Creature Is Used",
			Definition: "An ability that resolves after the card carrying it is itself used.",
			Body: `This ability resolves after the card that carries it is itself used — reaped or
fought with, or fired as an "Action:". Unlike After You Use a Card, which fires on
your other cards with the used card as "it", this fires on the used card itself,
so an upgrade can punish its own host (Containment Field destroys its host after
it is used).`,
		},
		{
			Section:    SectionAbility,
			Title:      "After a Creature Reaps",
			Definition: "An ability that resolves after any creature reaps, friendly or enemy.",
			Body: `This ability resolves after any creature reaps — friendly or enemy — with the
reaping creature as "it" (Orb of Invidius stuns whatever just reaped). It fires
on every in-play card, including the reaper itself.`,
		},
		{
			Section:    SectionAbility,
			Title:      "After an Enemy Creature Reaps",
			Definition: "An ability that resolves after an enemy creature reaps.",
			Body: `This ability resolves after an enemy creature reaps, with the reaping creature
as "it" (Pip Pip stuns the enemy that just reaped). Reaping happens only on the
reaper's own turn, so this naturally fires only for the reaper's opponent.`,
		},
		{
			Section:    SectionAbility,
			Title:      "Leaves Play",
			Definition: "An ability that resolves as the card leaves play by any route.",
			Body: `A Leaves Play ability resolves as the card leaves play by any route — destroyed,
purged, returned to hand, archived, or shuffled away. It fires before the card's
teardown, so the card is still on the board when it resolves. TriggerDestroyed is
the narrower "only when destroyed" version.`,
		},
	})
}
