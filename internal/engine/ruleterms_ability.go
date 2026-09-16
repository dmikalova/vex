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
			Title:      "After You Resolve a Damage Bonus Icon",
			Definition: "An ability that resolves after its controller resolves a Damage bonus icon, acting on the creature it hit.",
			Body: `This ability resolves after its controller resolves a Damage bonus icon, with
the creature that damage hit referred to as "it". It fires only for the resolving
player's own icons.`,
		},
		{
			Section:    SectionAbility,
			Title:      "After You Resolve a Draw Bonus Icon",
			Definition: "An ability that resolves after its controller resolves a Draw bonus icon.",
			Body: `This ability resolves after its controller resolves a Draw bonus icon. It fires
only for the resolving player's own icons, and only when the draw actually drew a
card.`,
		},
		{
			Section:    SectionAbility,
			Title:      "After You Forge a Key",
			Definition: "An ability that resolves after its controller forges a key.",
			Body:       `This ability resolves after its controller forges a key.`,
		},
		{
			Section:    SectionAbility,
			Title:      "After a Player Forges a Key",
			Definition: "An ability that resolves after any player forges a key, acting on the player who forged.",
			Body: `This ability resolves after any player forges a key, its own controller or
the opponent. The player who forged is the one it acts on, so "they" always
means whoever forged.`,
		},
		{
			Section:    SectionAbility,
			Title:      "After Your Opponent Forges a Key",
			Definition: "An ability that resolves after the opponent forges a key.",
			Body:       `This ability resolves after the opponent forges a key. It fires only on the opponent's forge, on the non-forging player's cards.`,
		},
		{
			Section:    SectionAbility,
			Title:      "When Your Opponent Would Forge a Key",
			Definition: "An ability that resolves when the opponent would forge a key, before the forge, and can prevent it.",
			Body: `This ability resolves when the opponent would forge a key, before that forge
happens, so it can prevent the forge. It fires only on the opponent's forge, and
the forging opponent is called "they". A prevented forge leaves the opponent's
Æmber unspent.`,
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
			Title:      "After an Upgrade Enters Play",
			Definition: "An ability that resolves after any upgrade enters play, including the opponent's.",
			Body: `This ability resolves after any upgrade enters play, whether you or your
opponent attached it.`,
		},
		{
			Section:    SectionAbility,
			Title:      "After a Creature Is Played Adjacent",
			Definition: "An ability that resolves after a creature is played into a battleline position next to this card.",
			Body: `This ability resolves after a creature is played into a battleline position
adjacent to the card holding it. A creature is only ever played onto its own
controller's battleline, so only the controller's own plays reach it. It does
not fire for a creature merely put into play by another effect, only for one
that is played.`,
		},
		{
			Section:    SectionAbility,
			Title:      "After a Creature Is Played",
			Definition: "An ability that resolves after any creature is played from hand, including the opponent's.",
			Body: `This ability resolves after any creature is played from hand, friendly or
enemy, with the played creature referred to as "it". Unlike After a Creature
Enters Play it fires only on an actual play, not on a creature put into play by
another effect, and it reaches every card in play whatever its battleline
position, so an artifact can watch the whole board.`,
		},
		{
			Section:    SectionAbility,
			Title:      "After a Neighbor Is Used to Fight",
			Definition: "An ability that resolves after a battleline neighbor of this card is used to fight.",
			Body: `This ability resolves after a battleline neighbor of the card holding it is
used to fight — the neighbor that fought is referred to as "it". It fires
whether or not the neighbor survives the fight.`,
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
			Title:      "After a Creature Is Destroyed in a Fight With",
			Definition: "An ability that resolves on a creature that survives a fight in which the other combatant was destroyed.",
			Body: `"In a fight with" names the timing window of a single fight: the moment power
damage is exchanged between the two combatants. This ability resolves on a
creature that survives a fight with the creature it is used against when that
other combatant is destroyed in the exchange; the destroyed creature is the one
referred to as "it".`,
		},
		{
			Section:    SectionAbility,
			Title:      "After an Enemy Creature Is Destroyed While Fighting",
			Definition: "An ability that resolves on a bystander when an enemy creature is destroyed in a fight.",
			Body: `This ability resolves whenever a creature is destroyed while fighting — either
combatant killed in the exchange of power damage — on a card whose controller is
the enemy of that creature. It fires on a bystander, not on a combatant, so an
artifact off to the side reacts to every enemy death in combat; the destroyed
creature is the one referred to as "it".`,
		},
		{
			Section:    SectionAbility,
			Title:      "After a Creature Is Destroyed by This Creature's Assault Damage",
			Definition: "An ability that resolves on a creature whose own Assault damage destroys the creature it attacks.",
			Body: `Assault deals its damage before the fight itself resolves. When that Assault
damage destroys the creature it attacks, the fight does not occur, and this
ability resolves on the attacker; the destroyed creature is the one referred to
as "it". It fires only when the Assault is the kill, not when the fight damage
that follows is.`,
		},
		{
			Section:    SectionAbility,
			Title:      "After This Creature Prevents Damage With Its Armor",
			Definition: "An ability that resolves after the card prevents damage with its own armor, scaled by the amount prevented.",
			Body: `This ability resolves after the card holding it prevents damage with its own
armor. The amount just prevented is the total armor it spent absorbing the damage
instance. Armor spent by a shield taking the damage in its place does not fire it,
only armor the card itself spends.`,
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
			Title:      "After a Tactic Is Played but Before It Resolves",
			Definition: "An ability that resolves after a Tactic is played, before that Tactic's own effect resolves.",
			Body: `This ability resolves after a Tactic is played — by either player —
before that Tactic's own effect resolves, so the reaction acts on the board the
Tactic is about to affect (Encounter Suit wards its host before the Tactic can
reach it). It fires on every card in play whoever played the Tactic.`,
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
			Title:      "At the Start of Each Player's Turn",
			Definition: "An ability that resolves at the start of every player's turn, resolving as the player whose turn it is.",
			Body: `This ability resolves at the start of every player's turn, its own
controller's and the opponent's, resolving as the player whose turn is starting so
"they"/"that player" is that active player rather than the card's controller
(Gambling Den, General Order 24). It is the whole-board companion to Start of Turn,
which fires only on its own controller's turn.`,
		},
		{
			Section:    SectionAbility,
			Title:      "At the End of Each Player's Turn",
			Definition: "An ability that resolves at the end of every player's turn, resolving as the player whose turn it is.",
			Body: `This ability resolves at the end of every player's turn, its own
controller's and the opponent's, resolving as the player whose turn is ending so
"they"/"that player" is that active player rather than the card's controller
(Pincerator). It is the whole-board companion to End of Turn, which fires only on
its own controller's turn.`,
		},
		{
			Section:    SectionAbility,
			Title:      "End of Ready Step",
			Definition: `An ability that resolves at the end of its controller's "ready cards" step, after every card has readied.`,
			Body: `An End of Ready Step ability resolves at the end of its controller's
"ready cards" step, once every card has readied (Greater Oxtet purges a card
from hand to grow).`,
		},
		{
			Section:    SectionAbility,
			Title:      "After a Player Chooses a House",
			Definition: "An ability that resolves after any player chooses their active house, whoever's turn it is.",
			Body: `This ability resolves after any player chooses their active house,
whether the choice was made by its controller or their opponent (Snag's Mirror
bars the chooser's opponent from that same house on their next turn).`,
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
			Title:      "After a Creature Is Destroyed",
			Definition: "An ability that resolves after any creature is destroyed, with that creature as \"it\".",
			Body: `This ability resolves after any creature is destroyed, friendly or enemy,
with the destroyed creature as "it" (Neffru gains that creature's owner 1 Æmber).
It resolves only after the whole destruction is complete and the destroyed cards
have reached their discard piles, so a card destroyed alongside them does not
react.`,
		},
		{
			Section:    SectionAbility,
			Title:      "After a Friendly Creature Is Destroyed",
			Definition: "An ability that resolves after a friendly creature is destroyed, with that creature as \"it\".",
			Body: `This ability resolves after a friendly creature is destroyed, with the
destroyed creature as "it" (Spartasaur destroys each non-Dinosaur creature). It
resolves only after the whole destruction is complete and the destroyed cards
have reached their discard piles, so a card destroyed alongside them does not
react.`,
		},
		{
			Section:    SectionAbility,
			Title:      "After Æmber Is Stolen From You",
			Definition: "An ability that resolves after Æmber is stolen from its controller, scaled by the amount stolen.",
			Body: `This ability resolves after Æmber is stolen from its controller. The number of
Æmber taken in that single theft is available to the effect, so it can scale by
it (Molephin deals 1 damage to each enemy creature for each Æmber stolen). A
theft from the other player does not fire it.`,
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
			Title:      "After a Creature Is Used to Fight",
			Definition: "An ability that resolves after any creature is used to fight.",
			Body: `This ability resolves after any creature is used to fight — friendly or enemy —
with the fighting creature as "it" (Shattered Throne makes it capture 1 Æmber). It
fires on every in-play card, including the fighting creature itself. Fighting
happens only on the attacker's own turn, so the fighter is always the active
player's creature.`,
		},
		{
			Section:    SectionAbility,
			Title:      "After a Friendly Creature Is Used to Fight",
			Definition: "An ability that resolves after a creature on the controller's own side is used to fight.",
			Body: `This ability resolves after a friendly creature — the controller's own, itself
or another — is used to fight, with the fighting creature as "it" (Lieutenant
Gorvenal captures 1 Æmber). An enemy creature fighting does not fire it.`,
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
