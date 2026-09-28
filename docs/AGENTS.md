# AGENTS.md — docs

Rules for the Markdown under `docs/` and every other Markdown file in the repo,
including `AGENTS.md` files and skill docs.

## Keep Markdown markdownlint-clean

Markdown is linted by [goldmark-lint](https://github.com/mrueg/goldmark-lint), a
Go port of markdownlint. `mage ci:markdown` reports what the gate would flag
(`mage ci:fix` applies the autofixable part); fix what it flags. Which rules
are off is in the generated `.markdownlint-cli2.yaml`; rely on no other rule
being off. Change a rule under `tools.markdownlint` in `mklv.config.json`,
never in the generated file. Wrap prose at about 80 columns. The human's
[todo.md](todo.md) is never touched, so its lint state is not your concern.

## Markdownlint pitfalls (running log)

Every time a markdownlint error is hit (outside `todo.md`), record it here with
the fix: the rule, the cause, the fix.

- **MD029 / ol-prefix — "Ordered list item prefix [Expected: N; Actual: M]".**
  A line at **column 0** inside an ordered-list item ends the list, so later
  items restart numbering. The trigger is a **long inline code span**: a reflow
  splits a span wider than the wrap column and the tail lands at column 0;
  indenting the continuation by hand does not survive the next reflow. Keep
  every inline code span short enough to fit on one line (split one long
  `` `card.New("Name", …, With*)` `` into `` `card.New(…)` ``,
  `` `card.Type` ``, `` `card.Provenance(card.<Set>, n)` ``).
