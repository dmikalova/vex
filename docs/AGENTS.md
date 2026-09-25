# AGENTS.md — docs

Guidance for the Markdown under `docs/` (and, by extension, every Markdown file
in the repo).

## Keep Markdown markdownlint-clean

Run `mage ci:fix && mage ci:markdown` after editing any Markdown and fix what
it flags before considering the work done — `ci:markdown` is part of
`mage ci:check`. Both run [goldmark-lint](https://github.com/mrueg/goldmark-lint)
(a Go port of markdownlint, `go run`-pinned so it is always available):
`ci:fix` applies its `--fix` corrections in place, and `ci:markdown` only
reports, failing on whatever cannot be autofixed. The config is generated into
`.markdownlint-cli2.yaml` from project-standards' base and vex's overrides under
`tools.markdownlint` in `mklv.config.json` (line-length, inline-HTML, both
emphasis rules, no-space-in-emphasis, fenced-code-language, first-line-heading,
blanks-around-fences/lists, and table-column-style are deliberately off; do not
rely on other rules being off). Change a rule there, never in the generated
file. The human's [todo.md](todo.md) is the one file agents never touch, so
its lint state is not your concern.

## Markdownlint pitfalls (running log)

Every time a markdownlint error is hit (outside `todo.md`), record it here with
the fix, so the same mistake is not made again. Keep each entry to the rule, the
cause, and the fix.

- **MD029 / ol-prefix — "Ordered list item prefix [Expected: N; Actual: M]".**
  A line at **column 0** in the middle of an ordered-list item terminates the
  list, so the items after it restart as a new list and their numbers no longer
  match what MD029 expects. The trigger is a **long inline code span** in a list
  item: when a reflow (editor format-on-save, or a sibling agent) rewraps the
  paragraph, a span wider than the wrap column is split, and the tail lands at
  column 0. Manually indenting the continuation does **not** hold — the next
  reflow undoes it. Fix durably by keeping each inline code span short enough to
  fit on one line (break one long `` `card.New("Name", …, With*)` `` span into
  several short ones like `` `card.New(…)` ``, `` `card.Type` ``,
  `` `card.Provenance(card.<Set>, n)` ``), so no span ever wraps.
