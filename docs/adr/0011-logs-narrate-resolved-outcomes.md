# The game log narrates resolved outcomes, not card text

## Context

ADR 0006 makes every `Effect` render its own printed text and carry itself out, so
a card's rules text can never drift from its behavior. The obvious next step is to
let the same tree render the game log too — one generator for card text, card
images, and logs alike.

It does not work. `Text()` renders an **unbound, present-tense imperative**:
"Deal 2 damage to a creature." A log line has to state a **bound, past-tense
outcome**: which creature, how much armor absorbed, whether a ward was spent,
whether the steal was capped by an empty pool, whether a replacement effect turned
a destruction into a purge. None of that is knowable from the AST — it is only
knowable after resolution. Deriving logs from the tree would force the tree to
depend on runtime outcomes, which is exactly the purity ADR 0006 buys.

The log we had did derive from the tree in two places, printing a card's rendered
ability line as its log entry. That is the failure this ADR names: the log
described what the card said instead of what happened.

## Decision

The game log is narrated from **resolved outcomes**. The sites that change state
emit typed **log entries**; the effect AST supplies only attribution — which card,
which ability category, which granting card.

- A `LogEntry` is an interface with one small variant per verb family, each
  rendering itself via `Text()` — the same Interpreter shape as `Effect`, one file
  per family. Extrinsic passes over entries (the client's icon and card-link
  rendering) are standalone type-switch functions, per ADR 0006's rule.
- `Text()` on an effect and `Text()` on a log entry are **two renderers over
  different inputs**. They share a vocabulary of phrasings — Æmber, flanks,
  "prevented by armor", card names, list joining — and neither derives from the
  other. Procedural card art is a third reader of the AST and is unaffected.
- Attribution is not passed to each emission site. Resolution opens a **frame**
  (actor, source card, ability category, granting card); entries inherit it.
  Frames nest, and the client groups a top-level frame's entries into one bubble.
- The log is the **public** view of the match: identical for both players, naming
  no card in a hidden zone. `Zone` carries whether it is public; a card in a hidden
  zone is named only after a `Reveal`.
- The log lives on `Game`, not `GameState`. Entries are interface values and so
  cannot satisfy ADR 0005's flat, pointerless, comparable state. Recording is
  switchable off, so `GameState.FastCopy` and the MCTS bot allocate nothing.

## Consequences

- Every observable state change is on the log and every line matches the state it
  describes. A missing line means the effect did not resolve; a zero means it
  resolved to zero.
- Cards author no log strings, and there is no free-form escape hatch. `Logf` is
  strangled during the migration and removed from the `Resolver` port's `Logger`
  role, which becomes a typed recorder. A card that appears to need free-form text
  is a missing entry variant, the same way it would be a missing effect node.
- The client stops parsing prose. Card links and the Æmber icon come from a
  `LocalID` on an entry, not from substring matching, and turn and phase
  boundaries come from typed entries rather than a `"--- "` prefix.
- Log behavior can never be load-bearing for game rules, since a simulated state
  has no log.
