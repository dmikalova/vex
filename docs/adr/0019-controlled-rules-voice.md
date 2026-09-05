# A controlled Rules voice for every player-facing surface

## Context

Vactrol writes rules text in several places — printed card text, the game log,
the generated rulebook, and the prose fragments that frame it — and each was
written to its own taste. The result reads unevenly: some lines are methodical
("Deal 2 damage to a creature"), others are flourishy or vague, and the same idea
is phrased differently depending on which surface it lands on. Flourish is not
free. A rule the reader has to interpret is a rule two readers interpret
differently, and a generator that has to render varied phrasings grows special
cases. `docs/card-wording-rules.md` already pins down 21 card-text conventions
this way, but it is scoped to card text, is framed as a diff against one
KeyForge set, and does not govern the log, the rulebook, or its own prose.

The text also has an authority problem. KeyForge is the source game and its
Master Rulebook (`docs/keyforge-master-rulebook.md`) is the canonical wording for
anything Vactrol has not already decided — but Vactrol deliberately diverges in
places (it renames the action card type to Tactic, retires `pay` for `give`,
re-expresses Omni as Versatile). Without a stated precedence rule, an author can't
tell whether to follow KeyForge or the local divergence, and "match KeyForge"
silently overwrites a deliberate choice.

## Decision

Vactrol adopts a single **Rules voice**: a controlled, methodical style that every
player-facing surface follows. The voice is a small curated stack, one register
per surface, so each surface reads consistently and renders predictably.

- **Base voice — STE-lite.** All rules text follows a lightweight subset of
  Simplified Technical English (ASD-STE100): short declarative sentences, one
  instruction per sentence, a controlled vocabulary (one word per meaning — a
  creature is _destroyed_, never also _killed_ or _removed_), present-tense active
  voice for what a card does, and no flourish. This is the default register unless
  a surface names a more specific one below.

- **Card text — MTG "Oracle" templating.** Printed card text uses the reference
  templating discipline of Magic's Oracle text: fixed phrasings for recurring
  shapes, keyword actions in place of longhand, and the ordering rules already
  captured in `docs/card-wording-rules.md` (front-loaded `for each`, front-loaded
  `instead`, `Choose one:` as a bulleted list, and so on). Card text is unbound,
  present-tense imperative, consistent with ADR 0006's `Text()`.

- **Rulebook rules — plain declarative statements.** A rule in the rulebook is a
  short declarative sentence in the present tense: it states what happens or what
  a player does, and a permission is phrased as a permission ("a player may …").
  A hard rule and a permission read differently because they say different things,
  not because they carry ceremonial keywords.

- **Example scenarios — Gherkin.** A bound example — a concrete situation
  demonstrating a rule — is written **Given / When / Then**. This is the shape the
  accuracy ratchet in ADR 0018 cites, and the shape an engine-backed rulebook page
  can execute directly.

- **Comments and docs — minimalism / plain language.** Doc comments and Markdown
  docs follow the same STE-lite plainness. Extending the controlled style all the
  way into code comments is a larger, separate migration and is its own decision
  (ADR 0020); this ADR governs player-facing and doc surfaces.

**Authority and precedence.** The KeyForge Master Rulebook is the wording
authority Vactrol consults **only when Vactrol has not already decided** a term or
phrasing. Where Vactrol deliberately diverges from KeyForge, **Vactrol wins**, and
the divergence is recorded in the Vactrol⇄KeyForge divergence register (the
divergence half of the split `card-wording-rules.md`). "Match KeyForge" never
silently overrides a recorded divergence.

`docs/card-wording-rules.md` is split to serve this. The wording conventions
(house style, surface-independent) stay in
[card-wording-rules.md](../card-wording-rules.md); the divergence register (where
and why Vactrol departs from KeyForge, plus the precedence rule) is
[keyforge-divergences.md](../keyforge-divergences.md). The numbered rules keep
their numbers so that the many `rule N` cross-references in code and docs stay
valid; the register indexes the divergent rules rather than renumbering them.

## Consequences

- Every rules surface reads in one voice, so a reader learns the dialect once and
  a generator renders fewer special cases.
- Æmber is spelled **Æmber** everywhere — the engine, tests, rulebook, and UI
  already do; the stray `Aember` in the wording rules and the rulebook prose
  fragments is corrected. One word per meaning applies to the game's own name for
  its resource.
- An author has a precedence rule: decide in the Vactrol voice; fall back to the
  KeyForge Master Rulebook only for what Vactrol has not settled; record any
  deliberate divergence rather than letting a later "match KeyForge" erase it.
- A pre-pass conforms existing card text, log lines, and docs to the voice. It is
  a text pass with no behavior change; where a phrasing maps to a structural rule
  (a keyword action, a front-loaded clause) it is already enforced by ADR 0006's
  renderer.
- The voice is documentation and review discipline, not a compiler. An optional
  controlled-vocabulary linter is possible later but is not required by this
  decision.
