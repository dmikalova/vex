# Code comments follow the same plain, controlled style

## Context

ADR 0019 defines a controlled Rules voice for player-facing surfaces — card text,
the game log, the rulebook, and the docs around them. It deliberately stops short
of code comments, because bringing comments into a controlled style is a different
kind of change with a different blast radius.

Comments are worth the same discipline for the same reason player-facing text is:
a comment the reader has to interpret is a comment two readers interpret
differently, and this repo leans hard on comments (the doc comment on every card
_is_ its rendered text under ADR 0006; `AGENTS.md` requires a comment and its
implementation to agree). But comments differ from rules text in ways that make a
single sweeping rewrite risky:

- There are far more of them, spread across every package, not four surfaces.
- A comment's audience is the next engineer, not a player, so the register is
  plain technical English, not card templating or rulebook phrasing.
- Rewriting a comment can subtly change what it claims the code does, which is
  exactly the kind of drift `AGENTS.md` warns against. A bulk pass touching every
  comment at once would be unreviewable.

## Decision

Code comments follow the same **plain, controlled style** as the Rules voice —
short declarative sentences, one idea per sentence, controlled vocabulary, no
flourish — adapted to a technical audience (plain English explaining _why_, not
card templating or rulebook phrasing).

This is a **decision now, migration later**. The style is adopted as the standard
for new and edited comments immediately, so nothing new is written in the old
flourishy register. Converting the existing body of comments is a **separate,
staged migration** run on its own schedule, package by package, small enough that
each step is reviewable and each rewrite can be checked against the code it
describes. The registry and voice work in ADR 0018 and ADR 0019 does **not** wait
on this migration.

The invariant that a comment and its implementation must agree (ADR 0006,
`AGENTS.md`) is the guard rail for every migration step: a comment is rewritten to
say more plainly what the code does, never to say something the code does not do.
Where a comment describes non-obvious behavior, the migration is the moment to
bind it to an example, as the accuracy ratchet does for rulebook terms.

## Consequences

- New and edited comments read in the plain controlled style from now on, so the
  comment base improves as the code is touched even before any dedicated pass.
- The large mechanical rewrite is decoupled from the rulebook and voice work, so
  neither blocks the other and neither turns into an unreviewable diff.
- Each migration step is bounded and checked against its code, so the "comment and
  implementation agree" invariant is strengthened by the pass, not endangered by
  it.
- This ADR records the intent; the migration's scope, order, and any tooling are
  worked out when it is scheduled, not here.
