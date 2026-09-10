# Gaining a text box copies traits, keywords, and triggered abilities

## Context

Two cards let a creature take on another card's text box. Mimic Gel is a Logos
creature whose `Play:` chooses another creature in play, gives itself +1 power
counters equal to that creature's power, and **gains the chosen creature's text
box** until Mimic Gel leaves play. Creed of Nurture is an Untamed artifact whose
`Action:` destroys itself, reveals a creature from your hand, and lends that
revealed card's text box to a chosen creature in play for the remainder of the
turn.

"Gaining a text box" is a deliberately narrow reimplementation of KeyForge's
copy/gain wording: the recipient takes on the **printed** traits, keywords, and
triggered abilities of the source card — as if they were printed on the recipient
— and nothing else. It does not copy the source's name, power, armor, type,
house, rarity, or any other stat. Power is handled separately (Mimic Gel copies it
as power counters); the text box is only the text.

State is flat, pointerless, and comparable (ADR 0005), so a gain cannot be a
closure or a pointer to another card — it has to be a plain, comparable reference
stored on the recipient.

## Decision

`CardCore` carries two `+1`-encoded source references, following the same
encoding as `HostPlus`/`ControlPlus` (`0` means none, otherwise `LocalID+1`):

- `TextBoxSourcePlus` — the source whose text box the creature gained **until it
  leaves play** (Mimic Gel). `resetCore` zeroes the whole core on exit, so the
  gain never outlives the recipient.
- `TextBoxTurnSourcePlus` — the source whose text box the creature gained for the
  **remainder of the turn** (Creed of Nurture). The ready phase clears it in the
  same both-players loop that clears the other turn-scoped grants
  (`TempPowerBonus`, `GrantedKeywords`, ...).

`grantedTextBoxSources(id)` folds both fields into the list of sources a creature
has gained. The text box is read **additively** at each read seam — a gained text
box adds to the recipient's own printed text, it does not replace it — and folds
in at exactly three seams:

1. **Triggered abilities** — `triggeredBy` gathers each gained source's printed
   `Abilities`, keeping `Source` set to the recipient so any self-referential text
   ("destroy this creature", "{self} gains…") rebinds to the gaining creature.
2. **Keywords** — `hasKeyword` folds in each gained source's printed keywords.
3. **Traits** — `HasTrait` folds in each gained source's printed traits.

All three seams gate on `!textBlanked(id)`: a blanked text box ignores a gained
text box just as it ignores the creature's own.

`GrantTextBox(recipient, source, remainderOfTurn)` is the one `CreatureResolver`
port method both effects share; it writes the appropriate field and logs
`CreatureGainedTextBox`. `GainTextBox` is the effect node Mimic Gel drives
(`Source` is `TheChosenCreature`); `LendTextBoxFromHand` is the composed effect
Creed of Nurture drives (reveal a hand creature, choose a creature in play, lend
for the turn).

## Read the source's printed text box, never a gained one

The source's text box is read from the immutable card catalog
(`g.cat.def(source)`), so it stays available even after the source card leaves
play — Mimic Gel keeps a copy of a creature that has since been destroyed. It is
always the source's **printed** text box, never a text box the source itself
gained, so gaining never chains or recurses. This matches KeyForge's rule that a
copy takes the printed base card, and it keeps the read a single non-recursive
lookup.

## Scope: what a gained text box does not include (yet)

A gained text box is deliberately limited to traits, keywords, and triggered
abilities. It does **not** include:

- **Constant abilities** (`ConstantAbility`) — a card whose text is "each creature
  gains…" is not folded into a gained text box.
- **Numeric combat keywords** (Assault N, Hazardous N, Splash-attack N) — these
  read their value from the printed definition at the combat seam, which the
  gained sources are not yet folded into.

Neither Mimic Gel nor Creed of Nurture needs these, so they are left out rather
than built speculatively. When a card does need one, fold the gained sources in at
that read seam (the constant-ability gather in `triggeredBy` / the numeric-keyword
readers in `game_combat.go`), the same way the three current seams do.

## Consequences

- A creature can never have a base power of 0 and rely on a later ability to raise
  it: a 0-power creature is destroyed at the entry settle boundary
  (`Power <= 0` is lethal, ADR 0029) before its `Play:` ability resolves. Mimic Gel
  therefore has a base power of 1 — a body that survives entering play — and copies
  the chosen creature's power as +1 counters on top of it.
- Because the gain is two plain enum-tagged fields, undo stays a snapshot and the
  state stays comparable; there is nothing to unwind but the fields clearing when
  the recipient leaves play (permanent gain) or at the ready phase (turn loan).
