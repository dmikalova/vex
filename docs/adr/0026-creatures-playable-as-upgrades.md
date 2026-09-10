# Creatures playable as upgrades

## Context

Worlds Collide introduces creatures that read "may be played as an upgrade
instead of a creature." Explo-rover is a skirmish creature that, as an upgrade,
grants its host skirmish; CALV-1N draws on Fight/Reap as a creature and, as an
upgrade, grants its host that same Fight/Reap draw. The player chooses which way
to play the card each time.

Such a card is two things at once. In the battleline it is an ordinary creature
with its own power, keywords, and abilities. Attached to a host it is an upgrade
that grants a continuous modifier and nothing else — its own creature keywords
and abilities do not apply, because it is a card in the upgrade chain, not a
creature on the battleline. The engine already models the second thing:
`StaticModifier` (`WithStatic`) is exactly "what an upgrade grants its host," and
the upgrade chain (ADR 0001) is type-agnostic, so a Creature-typed card can sit
in it without any storage change.

The open questions were where the creature-or-upgrade branch belongs, how it
interacts with a creature ban (Grommid), and how the two faces of the card
render without desyncing text from behavior (ADR 0006).

## Decision

A creature is marked `PlayableAsUpgrade` (option `WithPlayableAsUpgrade`) and
carries in its `Static` field what it grants a host. The card is authored twice
over on purpose — its creature keywords/abilities are its battleline identity,
and `Static` is its upgrade grant — because they are genuinely two different
things. `NewCard` fails loud (ADR 0010) unless the card is a Creature and its
`Static` actually grants something (`StaticModifier.grants()`).

**The branch is chosen at play time, in `playCardFromZone`'s `case Creature`.**
When the card is `PlayableAsUpgrade` and a host exists on either battleline, the
player is asked, through a plain `ChooseOption`, whether to play it as a creature
or an upgrade. Choosing "upgrade" picks a host (`pickCreature`) and routes to the
existing `playUpgradeCard` path — the same attach, Æmber bonus, and rendering an
ordinary upgrade uses. Choosing "creature", or having no host, falls through to
the normal creature play. Playing as an upgrade is not "playing a creature", so
it never emits the creature-entered reactions and never enters the battleline.

**A creature ban allows upgrade mode.** Playing a card as an upgrade is not
playing a creature, so Grommid's "you cannot play creatures" does not stop it —
as long as a host exists, the ban simply forces upgrade mode (the "creature"
option is not offered). With no host, a banned `PlayableAsUpgrade` creature
cannot be played at all (`ErrCannotPlayCreature`), and `CanPlay` mirrors this so
a UI can dim it correctly.

**Rendering folds the grant into a clause, and shows only the grant when
attached** (ADR 0006). In the hand/detail view the card prints its own creature
text plus one clause — `<name> may be played as an upgrade instead of a creature,
with the text: "…"` — quoting what it grants a host (the granted text's own
double quotes become single quotes so they nest). The standalone `Static` lines
are suppressed there, since they belong inside that clause. When the card is
attached to a host, `RenderUpgradeOnCreature` shows only the grant — its creature
keywords and abilities do not read on the host, because they do not apply while
it is an upgrade.

## Consequences

- The first place the engine branches on how a card is played based on the
  card's own text, rather than which `PlayX` method was called. The branch lives
  in one spot (`playCardFromZone`), so every zone that can play a creature
  (hand, deck, discard) gets the choice for free.
- Identity is by attachment, not printed type. `TypeOf(id)` returns `Upgrade`
  for a card that is currently attached in an upgrade chain, whatever its printed
  type, so `IsCreature(id)` is false for a creature-as-upgrade while it is
  attached. This is what a creature-reaching effect reads: `Destroy an artifact,
a creature, and an upgrade` (Destroy Them All) targets an attached
  creature-as-upgrade as the upgrade, and the view hides its power/armor, because
  an attached card has no creature stats to show ("null power"). The type is read
  from the attachment (`HostPlus != 0`), not a stored field, because the discard
  path for a shed upgrade (`discardUpgrades`) does not reset the core — a stored
  type would leak into the discard pile. Its `Static` still buffs its host like
  any upgrade. The conservation invariant that guards the chain therefore checks
  the **printed** type (a chain may hold an Upgrade or a Creature, ADR 0026;
  anything else is corruption), not `TypeOf`, which would be tautological.
- The choice is one-way per play. A card played as an upgrade is an upgrade until
  it leaves play; there is no flip back to a creature. When its host leaves play
  it is discarded like any upgrade (`discardUpgrades`, already type-agnostic).
- `PlayableAsUpgrade` is a capability flag, not a keyword, trigger, or card type,
  so it needs no rulebook term (ADR 0018): the printed clause carries its whole
  meaning.
