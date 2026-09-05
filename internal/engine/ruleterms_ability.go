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
			Section: SectionAbility,
			Title:   "Constant Ability",
			Body: `A constant ability is a continuous rule a card applies while it stays in play,
with no trigger of its own — "Each friendly creature gains +1 power", or a card
that grants every creature a keyword or a "Destroyed:" ability. Its effect
applies for as long as the source card remains in play and stops the moment it
leaves; applying it is not "using" the card and never exhausts it.`,
		},
		{
			Section: SectionAbility,
			Title:   "Play",
			Body: `A Play ability resolves right after you play the card from your hand. On a
creature or artifact it fires as the card enters play; on an action it is the
card's one-shot effect.`,
		},
		{
			Section: SectionAbility,
			Title:   "Reap",
			Body: `A Reap ability resolves after you use a ready creature to reap. Reaping gains
you 1 Æmber and exhausts the creature; the ability resolves in addition.`,
		},
		{
			Section: SectionAbility,
			Title:   "Fight",
			Body: `A Fight ability resolves after a creature you used to fight has dealt and
taken its combat damage and any resulting destruction has been carried out.`,
		},
		{
			Section: SectionAbility,
			Title:   "Action",
			Body: `An Action ability is one you resolve by using the card directly, without
reaping or fighting; using it this way exhausts the card.`,
		},
		{
			Section: SectionAbility,
			Title:   "After You Forge a Key",
			Body:    `This ability resolves after its controller forges a key.`,
		},
		{
			Section: SectionAbility,
			Title:   "After a Creature Enters Play",
			Body: `This ability resolves after any creature enters play, including creatures
your opponent plays.`,
		},
		{
			Section: SectionAbility,
			Title:   "Destroyed",
			Body: `A Destroyed ability resolves as the card is destroyed, before it reaches the
discard pile, so it can still act on the board it is leaving.`,
		},
		{
			Section: SectionAbility,
			Title:   "Before Fight",
			Body: `A Before Fight ability resolves when a creature is used to fight, before any
combat damage is dealt.`,
		},
		{
			Section: SectionAbility,
			Title:   "After a Creature Is Destroyed Fighting",
			Body: `This ability resolves on a creature that survives a fight in which the other
combatant was destroyed; the destroyed creature is the one referred to as
"it".`,
		},
		{
			Section: SectionAbility,
			Title:   "After You Play a Card",
			Body: `This ability resolves after its controller plays a card — a creature,
artifact, or action — from hand. Putting a card into play by another effect is
not "playing" it and does not fire this (that is TriggerAfterCreatureEnters).
AfterCardPlayed narrows it to a house and/or type.`,
		},
		{
			Section: SectionAbility,
			Title:   "End of Turn",
			Body: `An End of Turn ability resolves during the end of its controller's turn,
after cards ready and the controller draws (Shaffles drains the opponent at
each turn's end).`,
		},
		{
			Section: SectionAbility,
			Title:   "Start of Turn",
			Body: `A Start of Turn ability resolves at the start of its controller's turn,
before they forge, so an ability that changes what a key costs still has time
to.`,
		},
		{
			Section: SectionAbility,
			Title:   "After an Enemy Creature Is Destroyed",
			Body: `This ability resolves after an enemy creature is destroyed during its
controller's turn (Pile of Skulls captures Æmber onto a friendly creature
whenever an enemy creature is destroyed on your turn).`,
		},
		{
			Section: SectionAbility,
			Title:   "After Your Opponent Plays a Card",
			Body: `This ability resolves after your opponent plays a card (Teliga gains its
controller Æmber whenever the opponent plays a card).`,
		},
		{
			Section: SectionAbility,
			Title:   "After You Use a Card",
			Body: `This ability resolves after its controller uses a card — reaps or fights with a
creature, or fires an "Action:" — with the used card as "it" (Veylan Analyst
gains Æmber whenever you use an artifact).`,
		},
		{
			Section: SectionAbility,
			Title:   "After You Discard a Card From Your Hand",
			Body: `This ability resolves after its controller discards a card from their hand,
with the discarded card as "it" (Rock-Hurling Giant). Discarding from anywhere
else — the top of a deck, the archives — is not this.`,
		},
		{
			Section: SectionAbility,
			Title:   "Leaves Play",
			Body: `A Leaves Play ability resolves as the card leaves play by any route — destroyed,
purged, returned to hand, archived, or shuffled away. It fires before the card's
teardown, so the card is still on the board when it resolves. TriggerDestroyed is
the narrower "only when destroyed" version.`,
		},
	})
}
