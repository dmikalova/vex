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
  it, and `mage ci:check` stays green while it does.
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
  until ADR 0019's classification lands; see Consequences. (That exemption is
  withdrawn — see "Update: the node census closes the effect catalog".)

- **Fresh by gate.** `mage gen` regenerates the rulebook from the registry and the
  prose fragments, and `mage ci:check` fails on any diff between the committed file
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
- **Effects are deferred** (withdrawn — see "Update: the node census closes the
  effect catalog"). They are not a closed catalog and many are pure
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

## Update: the node census closes the effect catalog

Effects are no longer exempt, and the deferral to ADR 0019 is **withdrawn**. ADR
0019 is the controlled Rules voice; it contains no RuleBearing / Composition /
Implementation classification and never did, so the exemption above pointed at a
promise nothing kept. The hole it stood in for is now closed by a **node census**
in `internal/engine`, built on the mechanism ADR 0046 introduced for
`LogEntrySamples`: a hand-written catalog of constructed values, proved complete
against the package's own source by an AST scan.

**Closed catalogs, one per node family.** The AST node families — `Effect` and
the strategies beside it (`Refinement`, `Condition`, `Count`, `Selection`,
`Spread`, `Gather`, `Quantity`, and the rest) — are struct types rather than an
enum, so nothing enumerated them. Each now has a `catalog_<family>.go` holding one
`Catalogued` row per member, carrying a **constructed node** rather than a type
name, so a consumer can render and walk what the row names. `Families()`
(`catalog.go`) returns them all, and `TestFamilyTotality` binds each gated family
to the types the package's non-test source declares, **in both directions**: a
node with no row fails, and a row naming no declared type fails. `Target` is the
odd member — a flag struct kept comparable by ADR 0005 — so it is catalogued as
both an enum of its kinds and a family of its filter builders discovered by
shape; `Destination` is covered the same way.

**Interfaces and constants are guarded too.** A new `type X interface` in
`internal/engine` is a red build until `Interfaces()` (`catalog_interface.go`)
classifies it as a catalogued node family or records why it is not one
(`TestInterfaceTotality`) — so adding a family cannot skip the census. One level
down, `Enums()` (`catalog_enum.go`) scans the constants declared with each
text-bearing enum's type, so a `Keyword`, `Duration` or `CounterKind` that its
enumerating function forgot fails the build rather than going silently
undescribed (`TestEnumTotality`). That closes the original hole from underneath:
"complete by construction" no longer rests on an author remembering to extend
`Keywords()`.

**Two classes, not three.** Every row carries a mandatory `RulesBearing` with
exactly one of two columns set: `Term`, the rulebook title the node owes, or
`NoTerm`, a required one-line reason it owes none. Both empty and both set fail
(`TestFamilyRowsWellFormed`, `TestEnumRowsWellFormed`). The rule for `NoTerm` is
stated plainly: **a node is plumbing when its text contributes no vocabulary of
its own** — it only joins, repeats, gates or re-aims its children, and everything
a player reads comes from those children (`Sequence`, `Then`, `ForEach`,
`Repeat`, the duration wrappers). The three classes the Consequences promised
collapsed to two because nothing downstream reads the Composition /
Implementation distinction: every consumer asks only whether a term is owed. A
row also renders its node's text, so a zero-valued literal cannot stand in for a
node that prints nothing.

**Binding rules.** Terms bind **many-to-one**: several nodes may name one title
(all the Archive nodes are "Archive"; every `Duration` window is "Duration"),
because a term per node would turn the rulebook into an API listing. Binding is
**cross-section**: a row may name a title in any section, matched across the
whole registry, so the ward nodes point at the Keyword section's "Ward" and a
forge node at the Turn section's step rather than restating the rule
(`TestCatalogTermsAreRegistered`). And the binding is checked in **reverse**:
every term in the Effects and Card text sections must be claimed by at least one
row (`TestCensusSectionTermsAreClaimed`), so a term left behind by a deleted or
renamed node fails the build instead of rotting in the rulebook. A rule with no
census member — a standing restriction a card definition carries rather than an
effect node, such as "Must Fight When Used" — belongs in another section, not in
an exemption.

**The spine gained a "Card text" section.** The Decision names the spine as
"(turn, combat, cardtype, keyword, ability, effect)". It is now turn, combat,
cardtype, keyword, bonus, ability, **cardtext**, effect. Card text holds the
sentence parts a node leans on — Target, For Each, Duration, the conditional
frame — which the Effects section's verbs would otherwise have to restate one
verb at a time. It and Effects are the two sections the reverse check claims.

`mage tool:census` reports the whole census — uncatalogued nodes grouped by
declaring file, rows naming no node, unclassified interfaces, titles no term
carries, and unclaimed terms — so the state is readable without reading the
tests. The tests are what fail the build.
