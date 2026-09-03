# Vactrol style guide

This is the one place the whole repo's coding style lives. The `AGENTS.md` files
defer here for _how to write the code_; they keep only the structural facts about
where things go (build targets, file organization, engine constraints).

The guide is adapted from [TigerBeetle's TigerStyle](https://github.com/tigerbeetle/tigerbeetle/blob/main/docs/TIGER_STYLE.md). TigerStyle is
written for a Zig database; vactrol is a Go game engine, so the Zig-specific
mechanics (manual allocation, `zig fmt`, `snake_case` files, 100-column Zig) are
translated into their Go and vactrol equivalents, and a few rules are
deliberately relaxed or dropped where they don't fit a card-game engine. Where
vactrol departs from TigerStyle, the departure is called out.

## Why have style

> Another word for style is design. — TigerStyle

Style is not decoration on top of working code; it is how the code comes to work
well. Every rule below exists to advance one of vactrol's design goals, in this
priority order:

1. **Correctness and safety.** The engine is the referee. A rules bug is worse
   than a slow engine or an awkward API.
2. **Composability and maintainability.** Every card is an instance of a general
   mechanic, and the set keeps growing. Code that composes cleanly is worth more
   than code that is merely short.
3. **Performance.** The flat, pointerless state exists so a position can be
   cloned for free (MCTS). Performance is designed in, not bolted on.
4. **Developer experience.** Readability is table stakes — a means to the three
   goals above, not an end in itself.

When two rules conflict, the higher goal wins.

## Simplicity and zero technical debt

- **Simplicity is the hardest revision, not the first attempt.** "Let's do
  something simple" is easy to say; achieving it takes multiple passes and the
  willingness to throw a first draft away. Spend the thinking up front, in
  design, where it is cheapest.
- **Do it right the first time.** Prefer fixing a problem in design over
  implementation, and in implementation over production. Don't let a known
  latency spike, an unbounded loop, or an exponential algorithm slip through with
  a "we'll fix it later" — the later may not come.
- **Refactoring is welcome, and preferred over working around code.** Do not
  contort a new feature to fit today's shape. If a mechanic lands more cleanly
  after renaming, splitting, or generalizing an existing type, do the refactor.
  A good fit for the new feature (and the ones that follow it) matters more than
  preserving the current implementation. Keep such refactors focused, keep
  everything green (including 100% `internal/engine` coverage), and lean on the
  test suite — it pins each card's behavior and rendered text, so a regression
  surfaces immediately.

## Safety and correctness

These are [NASA's Power of Ten][powerbten] rules, translated to Go and to
vactrol's flat-state engine. Go has no `assert`, no manual memory, and a garbage
collector, so the letter of several rules changes even though the spirit does
not.

[powerbten]: https://spinroot.com/gerard/pdf/P10.pdf

### Control flow

- **Use only simple, explicit control flow.** Prefer straight-line code and
  shallow branching. Reserve recursion for structurally-bounded trees — the
  effect AST is walked recursively by `validateEffect`, and that is fine because
  the tree is finite and built at card-init — but never recurse over
  runtime-unbounded data.
- **Put a limit on everything.** Everything real has a bound. Loops and queues
  should have a fixed upper bound; where a loop genuinely cannot terminate (an
  event loop), make that explicit. This is why engine state is fixed-size arrays
  (`[maxCards]CardCore`, `[maxLasting]LastingEffect`) rather than unbounded
  slices — the limit is part of the design.
- **Push `if`s up and `for`s down.** When splitting a large function, keep the
  branching (`switch`/`if`) in the parent and move non-branchy work into helpers.
  Centralize control flow and state manipulation in one function; keep leaf
  functions pure — they compute what should change and return it, rather than
  reaching out and mutating.
- **Split compound conditions.** A condition that ANDs or ORs several booleans is
  hard to verify. Prefer nested `if`/`else` (or early returns) so each case is
  visible, and ask whether a positive branch also needs its negative counterpart
  handled.
- **State invariants positively.** Compare in the form the invariant naturally
  takes — `if index < count { /* holds */ }` — rather than negating it. Negations
  are error-prone.

### Errors and the invalid boundary

- **Handle every error.** Most catastrophic failures come from mishandled
  non-fatal errors. Never discard an `error` you can act on; never `_ =` an error
  return without a reason.
- **Validate at the boundary; then trust within it.** vactrol's "assertions" are
  its `validate()` pass and its invalid-zero sentinels. A card runs `validate()`
  at registration (`card.New` → `init`), so a malformed definition fails at
  program start, not mid-game. Sentinels like `playerUnset` / `targetUnset` /
  `durationUnset` / `eventUnset` turn an omitted required field into an
  init-time error rather than a silent zero-value default. Add these guards when
  you add a required field. A required _numeric_ field follows the same rule: do
  not silently coerce a meaningless zero into a working default (e.g.
  `CardsDiscarded.Amount` rejects `< 1` in `validate()` rather than treating `0`
  as `1`, because "discarded 0 or more cards" is always true). Make the author
  state the value; fail at registration if they don't.
- **Never silently or gracefully absorb bad data or an internal error.** A
  quietly dropped part is a bug that hides until someone notices it is missing —
  usually much later, and never through a test. Fail loudly at the boundary
  (`validate()`, a sentinel, a `panic` at registration) so the gap is impossible
  to miss. For example, a connected card (one whose behavior references another
  card that must exist) must **not** paper over a missing piece: if a connected
  part is not yet implemented, do not silently omit it from the connecting card,
  and do not have the connecting card gracefully tolerate the absent part.
  Instead leave the whole connecting card unimplemented until every part is in
  place, so its `TODO` stub marks it as outstanding and it is trivial to find.
  Half-wiring a card so it "works" with pieces missing trades a loud, findable
  gap for a silent, forgettable one.
- **`panic` only for the truly impossible.** A `panic` (e.g. `PlayerFor` on
  `playerUnset`) is acceptable only because `validate()` guarantees a real card
  can never reach it — it fires at authoring time, never in a live match. Never
  let a _computed_ value reach such a panic; reject bad input at the boundary and
  reserve the panic for the case that cannot occur.
- **Assert the positive space _and_ the negative space.** Tests must exercise not
  only valid data but invalid data, and data crossing from valid to invalid —
  that boundary is where the interesting bugs live. The 100% coverage gate on
  `internal/engine` exists to force every branch (including the error and
  impossible ones) to be exercised.

### Explicitness

- **Pass options explicitly at the call site instead of relying on defaults.**
  Being explicit avoids latent bugs if a default ever changes, and documents
  intent where it is read. In practice this is why cards are authored with each
  effect field named (`card.DealDamage{Amount: 3, Target: …}`) rather than
  leaning on zero values.
- **Always say why.** Never forget to motivate a non-obvious decision. Explaining
  the rationale increases understanding, invites correction, and shares the
  criteria by which the decision can be re-evaluated later.
- **No sentinel values or magic numbers — name a distinct state with `iota`.**
  When a value stands for a distinct state ("no rarity mark", "Connected", "unset
  player"), do not overload a number with a magic meaning (e.g. `-1` for
  Connected, `0` for none). Model the axis as an `iota` enum whose members name
  each state, and derive any count from the enum rather than hard-coding it. This
  matches the invalid-zero sentinels above: the zero value is the "unset/none"
  member, and the remaining members are ordered so ordinal arithmetic replaces
  literals where it reads naturally (e.g. the rarity marks run `rarityCommon…
  raritySpecial`, so a mark's ordinal *is* its diamond count).

## Composition and design

vactrol's defining rule. Read every change through the lens of _idiomatic,
maintainable, composable Go that will keep being extended_ as more of the game is
implemented. When a request is ambiguous, choose the option a senior Go engineer
would find easiest to build on — not the shortest path to a passing build.

- **Implement the mechanic, not the card.** Every card is an instance of a more
  general mechanic. Before adding a type, ask: what is the smallest orthogonal
  piece this card needs, and how would the _next_ card reuse it? Prefer many
  small pieces that snap together over one bespoke effect that does everything a
  single card happens to want.
- **Decompose fused effects.** A card that "does A, then B if C" is a `Sequence`
  of an effect, a `Conditional`, and a condition — not one `DoAThenBIfC` node.
  Bonkers Killing Machine is `DiscardTopOfEachDeck` →
  `DestroyOfEachDiscardedHouse` → `Conditional{CardsDestroyedFewerThan, Destroy}`,
  not a single `DiscardAndDestroyByHouse`.
- **Thread state through the context; don't fuse producers and consumers.** When
  one step produces a value a later step consumes (cards revealed, houses
  discarded, creatures healed, a card put in focus as `ctx.It`), record it on the
  `EffectContext` and let any following effect or condition read it. This is how
  `Heal` + `CreaturesHealed`, or `RevealTopOfDeck` + `ItIsOfHouse` +
  `PlayRevealedCard`, compose without a combined type.
- **Parameterize over enums, not new types.** "At least N of a house" is a
  reusable count with a `House` and `Amount` field, not a new `…OfHouseAtLeast`
  type. A new lifetime is a `Duration` field (`EndOfTurn`,
  `UntilThisLeavesPlay`), not a `…ForRemainderOfTurn` variant. "Which house" is a
  house choice/reference, not a `…OfChosenHouse` / `…OfActiveHouse` pair.
- **Reuse the shared vocabularies.** `Target` already filters by house, trait,
  type, chosen/active house, and set-relative selectors — reach for it (or extend
  it) before inventing a parallel filter. Events (`EventCreaturePlayed`,
  `EventReap`, `EventCreatureDestroyed`, …) with a subject already drive triggers,
  lasting reactions (`ForRemainderOfTurn`), and replacements (`Instead`,
  `Replace`) — extend that spine rather than adding a one-off flag on
  `CardDefinition`.
- **Self-reference through existing seams.** "Destroy this Upgrade" is
  `Destroy{Target: This}` (the destroy path detaches an attached upgrade), not a
  `DestroyThisUpgrade`. "This creature captures the Æmber" is a replacement of the
  add-to-pool event, not a `CapturesOpponentAember` bool.
- **Vary behavior with a strategy that renders its own text.** When behavior
  changes along an axis, model the axis as a small strategy — a `Chooser`,
  `Selector`, `Count`, or `Condition` — each carrying both its behavior and its
  text fragment, so it plugs into the AST without the printed text ever drifting
  from the behavior. Reach for this before adding another `Target` field or a
  `bool`.
- **A one-off name is a smell.** If a name encodes a whole card sentence
  (`ReadyAndBelongToHouseAfterYouPlayCreature`), it is hiding reusable pieces —
  split it. The cost of a slightly larger card definition is worth an engine the
  next card can build on. When a card genuinely needs something new, add the
  smallest general primitive and one card that uses it — never speculative, but
  always shaped so the next card can reach for it.
- **A phrasing helper belongs in the shared vocabulary, not beside its first
  caller.** `indefinite`, `plural`, `countNoun`, `singularNoun` and friends live
  in `internal/engine/text.go`; the log's equivalents (`namedCards`, `because`,
  `nameMoved`) live in `log.go`. A helper left in the `effect_*.go` that first
  needed it is invisible to the next author, who then writes a second copy — that
  is how `"a " + noun`, `noun + "s"` and `card(s)` all ended up hand-rolled in
  five places each. When you find yourself formatting English, look in the shared
  file first, and put anything new there.

The engine's structural realization of these principles — the effect AST, the
`Resolver` port, the strategy families, the lasting registry — is documented in
[internal/engine/AGENTS.md](../internal/engine/AGENTS.md). This section is the
_why_; that file is the _what goes where_.

## Performance

> The lack of back-of-the-envelope performance sketches is the root of all evil.

- **Design for performance from the outset.** The biggest wins come in the design
  phase, before anything can be profiled. vactrol's flat, pointerless,
  comparable `GameState` is exactly this: it makes cloning a position a plain
  value copy with no allocation, which is what makes AI search viable. Do not
  introduce a pointer, slice, or map into `GameState` (or into a value compared
  against it, like `Target`) — that one change would forfeit the win.
- **Have mechanical sympathy.** Work with the grain of the machine. Give the CPU
  large, predictable chunks of work; batch instead of context-switching per
  event. Don't react directly to every external event — run at your own pace and
  batch, which is both safer (control flow stays yours) and faster.
- **Be explicit; don't lean on the compiler to save you.** Extract hot loops into
  standalone functions with primitive arguments where it makes the redundant work
  visible to a human reader, not just the optimizer.

Vactrol's departure from TigerStyle: Go is garbage-collected, so the "all memory
statically allocated at startup, none allocated after init" rule does not apply
literally. Its spirit survives as the flat fixed-size state and a preference for
avoiding per-turn allocation on hot paths.

## Naming and vocabulary

> There are only two hard things in Computer Science: cache invalidation, naming
> things, and off-by-one errors.

- **Get the nouns and verbs just right.** Great names capture what a thing is or
  does and show you understand the domain. Take the time to find them.
- **Follow Go's naming conventions**, not Zig's. Use `MixedCaps` /
  `mixedCaps`, not `snake_case`, for identifiers. Keep acronyms in one case
  (`LocalID`, `HTTPServer`, not `LocalId` / `HttpServer`). File names are
  `snake_case.go`.
- **Do not abbreviate.** Spell names out; the reader should not have to expand
  them. (Single-letter loop indices in a tight numeric loop are the only
  exception.)
- **Put units and qualifiers last, most-significant word first.** `latencyMsMax`,
  not `maxLatencyMs`, so related names line up (`latencyMsMin` sits right beside
  it) and group by subject.
- **Prefer nouns to participles for things that get talked about.**
  `replica.pipeline` reads better than `replica.preparing` because a noun can be
  used directly in prose and composes into derived names (`pipelineMax`).
- **Don't overload a name with context-dependent meanings.** Reusing a term for
  two concepts causes exactly the confusion you'd expect.
- **Speak KeyForge, not generic game-speak.** Names — types, methods, fields,
  effects, targets — must stay within KeyForge's own vocabulary. Say
  `ExceptMostPowerfulCreature`, not `ExceptStrongest`; `CannotFight`, not
  `PreventFight`. When you need the right word, source it in this order:
  1. the provenance files (`internal/cards/provenance`) and the original card
     text they point at — the canonical wording;
  2. existing implementations in this repo — reuse an established term rather than
     coining a synonym.

  If neither has a term, pick the phrasing closest to how KeyForge cards are
  actually written, and flag it for the user.

## Comments and commit messages

- **Say why, not what.** Code shows what it does; a comment exists to explain what
  the code cannot — the rationale, the surprising constraint, the reason this
  shape was chosen. Don't restate the next line. Keep it to one short line where
  one will do; don't write a doc-comment paragraph for a one-line point.
- **Say how, for tests.** A short description at the top of a non-obvious test
  explaining its goal and method helps the reader get up to speed or skip past it.
  In vactrol, the behavior description goes in the `t.Run("…")` subtest name
  (describe/it shape), not in the generated doc comment.
- **Comments are prose.** Full comments are sentences: a space after `//`, a
  capital letter, and a full stop (or a colon when they introduce what follows).
  Trailing end-of-line comments may be terse phrases with no punctuation.
- **Write commit messages that inform.** The commit message is stored in git and
  read in `git blame`; a pull-request description is not, so it is no substitute.
  Explain the change and why it was made.

## Off-by-one errors

- **Treat `index`, `count`, and `size` as distinct types** even though they are
  all integers, with clear rules to convert: index→count adds one (indexes are
  0-based, counts 1-based); count→size multiplies by the unit. Naming with the
  qualifier last makes these conversions legible.
- **Show your intent with division.** Make it obvious whether a division is
  exact, floored, or ceilinged, rather than leaving the reader to guess whether
  you considered the rounding case.

## Formatting and tooling by the numbers

- **Run everything through `mage`, never raw `go`.** The mage targets wrap the
  project's conventions (coverage gate, comment/rulebook generation, golines). The
  gate is `mage check` (fmt-check, build, vet, lint, test, coverage); it must
  print `ALL GREEN` before work is considered done. `mage fmt` formats;
  `mage cover` keeps `internal/engine` at 100%. See the root
  [AGENTS.md](../AGENTS.md) for the full target list.
- **Let `golines` own formatting.** Indentation, alignment, and wrapping are not
  matters of taste here — `golines` (via `mage fmt`) decides them: it applies
  `gofmt` and additionally shortens code lines over 100 columns. Do not fight it.
- **Keep functions short enough to see at once.** There is a real
  discontinuity between a function that fits on screen and one you must scroll to
  read; aim to keep functions within roughly 70 lines, splitting by pushing `if`s
  up and `for`s down. This is a guideline, not a gate.
- **Keep lines to a readable measure.** `golines` shortens over-long Go code
  lines automatically (100 columns); wrap markdown and comment prose at roughly 80
  columns to match the existing docs.
- **Brace `if` statements** unless the whole statement fits on one line, as
  defense against "goto fail"-style bugs (Go's `gofmt` and `go vet` largely
  enforce this already).

## Dependencies and tooling

- **Keep dependencies minimal and deliberate.** This is vactrol's most explicit
  departure from TigerStyle's _zero-dependencies_ policy: the frontends do depend
  on libraries (a TUI toolkit, a WebAssembly framework). But the core `engine`
  imports nothing upward and stays dependency-light on purpose, and a new
  dependency — especially anywhere near the engine — must earn its place against
  the supply-chain, safety, and maintenance cost it adds.
- **Prefer the tools already in the box.** A small, standardized toolbox is
  simpler to operate than an array of specialized instruments. vactrol's is Go
  plus `mage`; when you need a script or a task, add a mage target rather than a
  one-off shell script, so it stays cross-platform and typed.
