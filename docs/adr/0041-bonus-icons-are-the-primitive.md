# 41. Bonus icons are the primitive; Enhance distributes them at generation

## Context

KeyForge cards carry **bonus icons** in the upper-left corner — Æmber, Capture,
Damage, Draw (and Discard, House, +1-power, not modelled yet). After a card is
played its icons resolve one at a time, top to bottom, before its "Play:" ability.
The **Enhance** keyword is a deck-building rule: an Enhance card adds the listed
icons to random _other_ cards in the deck at generation time and gains nothing
itself during play.

The engine modelled only one of these: a scalar `AemberBonus int` that gained the
controller that much Æmber on play. There was no representation of the other icon
kinds, no ordered icon list, and no Enhance pass — the deck-generation seam for it
was documented (ADR 0004) but inert. Roughly a third of Mass Mutation is Enhance
cards, so the set could not be finished without modelling this.

## Decision

Make the **bonus icon** the primitive, not Æmber-specifically.

- A `BonusIcon` enum names the kinds (`BonusAember`, `BonusCapture`, `BonusDamage`,
  `BonusDraw` to start). `CardDefinition.AemberBonus int` becomes
  `Bonuses []BonusIcon`, an **ordered** list, authored with
  `WithBonus(card.Bonus.Aember, …)`. `WithAemberBonus` is removed; every call site
  migrated. Æmber is just one icon kind among the four.
- Playing a card resolves its `Bonuses` as a first-class step (`resolveBonusIcons`),
  one icon at a time in printed order, before its "Play:" abilities — the same
  window KeyForge uses. Each icon is resolved through the ordinary game seams
  (`gainAember`, `draw`, `dealDamage`, capture), so continuous replacements (Ether
  Spider) still intercept a bonus Æmber gain, and the step is a clean place to hang
  the icon-manipulating cards (Wild Bounty, Amphora Captura, Master of the Grey)
  when they are built.
- An Enhance source declares `Enhances []BonusIcon` (authored with `WithEnhance`),
  which renders its "Enhance …" line and drives a **deck-wide finishing pass** in
  `deckgen` (the pass ADR 0004 reserved). After the deck is valid, each source's
  icons land on uniformly-random eligible cards, capped at five icons per card
  (printed icons count), deterministic from the deck seed. `Bonuses` (icons on me,
  resolve on my play) and `Enhances` (icons I contribute to the deck) are separate
  fields with opposite direction.

Two deliberate divergences from KeyForge, both recorded in the divergence register:
the **source** of a bonus icon's effect is the card carrying it (KeyForge treats
the game as the source), and a **creature that leaves play before or while its icons
resolve does not resolve the rest** (KeyForge resolves them all once the card is
played). A third, non-KeyForge addition: a card may bar specific bonus-icon kinds
from _receiving_ Enhance icons (`WithoutEnhancement(kinds…)`), for a bonus that
would only weaken it.

## Consequences

- Bonus icons render in the Rules voice and the log says "bonus" ("Splinter bonus
  deals 1 damage to …"), so a bonus icon is never confused with the same effect from
  an ability. They have a rulebook section, complete by construction (ADR 0018): a
  new icon kind without a term fails the build.
- The engine's flat-state invariant (ADR 0005) is untouched — `Bonuses`/`Enhances`
  are immutable catalog data on `CardDefinition`, like `Keywords`. The deckgen pass
  clones a card's `Bonuses` before landing an icon, because materialize copies defs
  shallowly and two slots of the same card would otherwise share one backing array.
- Upgrades attach _before_ their icons resolve, so a Damage bonus that kills the
  intended host finds the upgrade already attached and it sheds cleanly instead of
  attaching to a destroyed host.
- The icon-manipulating cards are deliberately deferred: the resolution step is
  first-class and interceptable now, but Wild Bounty, Amphora Captura, Scrivener
  Favian, Master of the Grey, and Ensign El-Samra are built later against it.
