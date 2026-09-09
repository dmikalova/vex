# Agent scratchpad

Working notes the agent writes for itself, to translate a request into concrete
work and show what is left. This is **not** [todo.md](todo.md) — that is the
human's personal list, which agents never write into. Rules:

- When an item is **done, delete it** — do not mark it done. This file only ever
  shows outstanding work, so it reads as a live "here is what I still mean to do"
  surface for coordinating with the human.
- Keep items concrete and **grouped by area or mechanic**, so related work is
  built together: add the shared primitive once, then knock out the group.
- Cite the ADR or doc that decided an item where one exists.

## Card wording / authoring

- **Festering Touch** reword, e.g. "Play: Choose 2 creatures. Deal 1 damage to the
  chosen creatures with no damage. Deal 3 damage to the chosen creatures with
  damage."
- **Orator Hissaro** could read: "Play: Exalt and ready each neighboring creature.
  For the remainder of the turn, those creatures belong to house Saurian."
- **Borr-Nit** and similar could be atomized and recomposed further (decompose
  fused effects into shared nodes).
- **Memory Chip**: after choosing self house, archive a card from hand (check the
  original printed card text and match it).

## Card catalog / provenance

_No outstanding items._

## Web — mobile, previews, layout

- **Log hover preview on mobile = single tap.** Currently needs a double tap on the
  card name to raise the preview; make it one tap.
- **Facedown cards face down for BOTH players.** A facedown card shows its back to
  everyone; only the controller can hover (desktop) / tap (mobile) to peek at its
  face. Must be visually distinguishable from a faceup card.
- **Card title dynamic shrink (QUESTION).** Yshi's title shrinks by a fixed step
  rather than to the exact needed size. Explain whether we still use predetermined
  shrink amounts, whether dynamic fit is possible, and why smooth fit is/ isn't
  feasible.
- **Glyph bar margin** — remove the margin around the glyph bar (or the surrounding
  boxes) to save space.
- **Glyph strip — residual effect-level unknowns.** Card-level features, all
  triggers (incl. phase triggers), Static/Restrictions/Toll/Replaces/KeyCost, and
  Splash-attack now transcribe. Still falling back to the abstract `glyph-unknown`
  (`*`): the inner `.Then` effect of `ChooseHouseThen`/`ChooseCreatureThen`
  (Restringuntus, Deep Probe, Niffle Grounds) and the granted-ability effects on
  the Blasters, Evasion Sigil, and Rocket Boots. Map each inner effect to a glyph
  (ADR 0022). The remaining empty strips (Dust Pixie, Toad, Mega, etc.) are
  Æmber-bonus-pip-only creatures — the pips render on the card face, not the strip,
  so those are intentional, not gaps.

## Tooling / tests

- **Unused-asset test.** Add a test that fails when an asset in web/assets is not
  referenced by the code (no dead assets).
