# The rulebook is a typed term registry, complete by construction

> **Update.** The typed registry decided here still stands, but the rulebook is no
> longer generated to `docs/rulebook.md`. The web client renders the registry live
> at `/rulebook` and `/glossary` (via `engine.RuleBook()` / `engine.Glossary()`),
> so the `magefiles/genrules` renderer and the freshness check have been removed.
> References below to the generated file describe the original decision.

## Context

The rulebook (`docs/rulebook.md`) is generated. Foundational prose lives in
`docs/rulebook/*.md`; every keyword, trigger, card type, and rule-bearing effect
is meant to describe itself in a doc comment next to the code that enforces it, so
the two can never drift. `magefiles/genrules` walks `internal/engine`, harvests
each `//rulebook:<section> <Title>` directive it finds, and splices the harvested
bodies into the fixed section spine. This is the same "text lives next to code"
idea as ADR 0006's `Text()` and `gencomments`.

Three gaps make the rulebook untrustworthy in a way per-card text is not:

- **It is not complete.** A directive is opt-in. `Elusive` and `Taunt` are valid
  `Keyword` enum values with no `//rulebook:` comment, so they are silently absent
  from the rulebook. Nothing counts the keywords the game actually has against the
  keywords the rulebook describes. A missing term is invisible.
- **It is not fresh.** `mage gen` regenerates the file, but nothing fails if a
  human forgets to run it. The committed rulebook can lag the code that generates
  it, and `mage check` stays green while it does.
- **It is not bound to behavior.** A directive body is free prose. It can say
  anything, including something the engine no longer does, and no test notices.
  `Text()` on an effect can't lie because it renders the very node that resolves;
  a `//rulebook:` comment has no such tether.

The harvest is also indirect: it re-parses Go source as text (an AST walk over
comments) to recover facts the engine already holds as typed values. The set of
keywords is `Keywords()`; the set of triggers and card types are enums. The
generator reconstructs by string-scraping what the package could hand it directly.

## Decision

The rulebook's terms are a **typed registry** in `internal/engine`, co-located
with the enums they describe, and the registry is **complete by construction**,
**fresh by gate**, and **bound to behavior by example**.

- **Registry.** Each term is a `RuleTerm` value: `Section` (which part of the
  spine), `Title`, a one-line `Definition` written in the Rules voice (ADR 0019)
  that doubles as the glossary entry, and a Markdown `Body`. Terms register next
  to their enum — a keyword's term sits beside the `Keyword` it describes — and
  the package exposes them as an enumerable table (`RuleTerms()`), the same way it
  already exposes `Keywords()`. The registry is the single source shared by the
  rulebook and the eventual in-client glossary; a glossary is the `Definition`
  column of the same table.

- **Complete by construction.** A test enumerates the **closed catalogs** the
  game defines — every `Keyword`, every `Trigger`, every card type, every turn and
  combat step — and fails if any member has no registered term. Adding a keyword
  without describing it breaks the build. `Elusive` and `Taunt` are the first two
  the test lights up. Effects are **not** a closed catalog in v1 and are exempt
  until ADR 0019's classification lands; see Consequences.

- **Fresh by gate.** `mage gen` regenerates the rulebook from the registry and the
  prose fragments, and `mage check` fails on any diff between the committed file
  and a fresh regeneration. The rulebook cannot lag the registry through a
  forgotten `mage gen`.

- **Bound by example (accuracy ratchet).** A term's body starts bound only by
  co-location — it sits on the symbol it describes, so a reviewer changing the
  symbol sees it. For terms whose correctness is subtle, a term may cite an engine
  test or scenario that demonstrates the rule, and the completeness test requires
  the citation to name a real test. The end state is an engine-backed rulebook
  page where a reader runs the cited scenario and watches the rule happen, so the
  rulebook doubles as regression coverage. The ratchet is opt-in per term now and
  tightened over time; it is not required of every term at once.

- **`genrules` becomes a renderer.** It imports the engine registry instead of
  re-parsing source, groups terms by section, and splices the `docs/rulebook/*.md`
  fragments as it does today. The AST harvest and the `//rulebook:` directive are
  removed. The section spine (turn, combat, cardtype, keyword, ability, effect)
  and alphabetical ordering are unchanged.

## Consequences

- A term can no longer be missing without the build knowing. Completeness is a
  test over the same enum the game plays with, not a hope that every author
  remembered a comment.
- The committed rulebook is always the rulebook the current code generates. A
  stale file is a red gate, not a silent lie.
- The generator stops string-scraping Go source and reads typed values, so a
  renamed term or a new keyword flows through as data, not as a re-parsed comment.
- **Effects are deferred.** They are not a closed catalog and many are pure
  composition or plumbing that bear no player-facing rule. ADR 0019 classifies
  each effect node as RuleBearing, Composition, or Implementation; only
  RuleBearing effects owe a term, and a non-RuleBearing effect that nonetheless
  carries a rule must cite an example. Until that classification exists, effects
  keep contributing bodies but are exempt from the completeness test, so this ADR
  does not force a rule onto every composition node.
- The registry lives in `internal/engine`, which is under the 100% coverage gate,
  so the registry, the completeness test, and the renderer's engine-facing surface
  are fully tested. The freshness gate runs the existing `mage gen`, so it adds no
  new untested engine code.
- The six framing fragments in `docs/rulebook/` stayed hand-authored and outside
  the registry when this ADR landed. They have since moved into the registry too
  (see Update below).

## Update: framing fragments moved into the registry

The six framing fragments — the document overview and each section's intro — now
live in the engine registry alongside the terms, not as loose Markdown under
`docs/rulebook/`. The engine exposes `RuleOverview()` and `RuleSectionIntro()`
next to `RuleTerms()`; the overview registers from `ruleterms_overview.go` and
each section intro registers from its `ruleterms_<section>.go` init, beside that
section's terms. `genrules` reads all three from the engine and no longer touches
the filesystem for prose, so the whole rulebook — terms and framing alike — flows
through one typed contract. The `docs/rulebook/` directory is removed.
