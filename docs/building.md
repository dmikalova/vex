# Building, testing, and the mage targets

Every build, test, format, and coverage task runs through `mage`. The generic
gate comes from project-standards' shared `ci` targets (imported in
`magefiles/build.go`); the rest are vex's own targets, which wrap the project's
conventions (comment generation, the simulator, the card tools). Run `mage -l`
to see every target.

## The gate

- `mage ci:fix && mage ci:check` — **the local validator**; run it before
  considering work done. `ci:fix` applies every autofix (go fix, golangci-lint
  `--fix`, goimports, golines, gci, goldmark-lint `--fix`, misspell `-w`) and
  regenerates the generated configs. `ci:check` only verifies — it writes no
  file and fails on anything `ci:fix` would change — then runs format, tidy,
  build, vet, lint, markdown, spell, secrets, commits and drift, then the tests
  and the coverage gates. It must print `ALL GREEN`. Agent sessions under
  diatom do not run it themselves: diatom runs `ci:fix` and then the gate when
  a session finishes.
- `mage ci:build` — build all packages for the host, then the web client for
  js/wasm (into a temporary directory; `mage webWasm` writes `web/app.wasm`).
- `mage ci:test` — run all tests.
- `mage testRun <pattern>` — run only the tests matching a name pattern.
- `mage ci:vet` — run `go vet`.
- `mage ci:lint` — run golangci-lint, fixing nothing.
- `mage ci:format` — check formatting without writing (`mage ci:fix` formats).
- `mage ci:markdown` / `mage ci:spell` — lint Markdown, check spelling.
- `mage ci:cover` — report coverage for each gated area, which must stay at
  100%. The areas are `internal/engine`, the card definitions under
  `internal/cards/sets`, `internal/cards/cardtest`, and `internal/deckgen`;
  they are listed as `ci.CoverGates` in `magefiles/build.go`. `internal/web` is
  deliberately ungated. See [testing.md](testing.md).

## Generation

- `mage generateComments` — rewrite each card's doc comment from its
  definition.
- `mage gen` — regenerate the card comments. It is deterministic and
  regenerates from the tree, so run it freely whenever a definition changes.

## Simulation, profiling, and open findings

- `mage debug` — replay a simulated game with the game log on and print the log
  tail next to the invariant violation that ended it. With no `-script` it finds
  the first failing game in the fixed-seed property batch `mage ci:test` plays;
  pass `-script` the hex a failure printed to replay that one, and `-tail` to
  widen the log (`mage debug -script=<hex> -tail=200`). How to read what it
  prints is in [sim-debugging.md](sim-debugging.md).
- `mage trace` — play the fixed-seed property games once with the game log on
  and write every line to `tmp/sim/trace.log` (gitignored), so a whole game
  reads end to end. Where `mage debug` shows the tail of the game that broke, a
  trace is the full log of games that pass. `-count` sets how many games
  (default 1), `-out` the destination
  (`mage trace -count=25 -out=tmp/sim/mine.log`).
- `mage profile` — profile the engine under the whole-game simulator (the
  closest proxy for MCTS load), write CPU and allocation profiles to `tmp/sim`,
  and print the per-game counts next to a committed baseline so a regression
  shows as a delta. Default runs the 1000 seeded games; `-random` runs a fresh
  random outlier batch; `-save` re-blesses the baseline. It never blocks and
  never gates — the baseline is advisory and lives outside `mage ci:check` and
  CI. See [testing.md](testing.md).
- `mage profileServer` — open a profile from the last `mage profile` run in the
  interactive pprof web UI (flame graph, call graph, source view); blocks until
  stopped. `-mem` serves the allocation profile, `-random` the outlier run's.
- `mage uiTest` — drive the `/ui-test` browser scenarios headlessly: build the
  wasm client, serve it, and run each scenario registered in
  `internal/web/uitest_scenarios.go` in headless Chrome (go-rod). The driver
  package `internal/web/uitest` is behind the `uitest` build tag, so it is
  invisible to `./...` and **not** part of `mage ci:check` or `ci:test` — a
  wasm build plus a browser boot is not a per-save gate. See
  [testing.md](testing.md).
- `mage corpusPrune` — replay every entry in `FuzzPlay`'s seed corpus and
  rewrite it as one minimized entry per bug that still reproduces, dropping the
  entries whose bug is fixed. The corpus is the list of open findings, not an
  archive of every script a soak ever saw; run this after fixing a soak or fuzz
  find.
- `mage capturePrune` — the same, for the browser client's replay captures in
  `internal/web/testdata/capture`: replay every entry, delete the ones whose
  fault no longer reproduces and every one recorded against a different
  command-log version or card pool, keep the rest. A capture is written by the
  dev server when the client fails to replay a saved match; `TestCaptures`
  replays them all and fails, so the directory is the list of **open findings**
  and a capture in a commit is a lapse. See [testing.md](testing.md).

## Tools and docs

- `mage tool:*` — card research and authoring tools (`tool:lookup`,
  `tool:missing`, `tool:nextCard`, `tool:coverage`, `tool:stub`); described in
  [internal/cards/AGENTS.md](../internal/cards/AGENTS.md). `tool:nodeUsage`
  and `tool:census` audit the engine and are described in
  [internal/engine/AGENTS.md](../internal/engine/AGENTS.md).
- `mage tool:gameSize` — report the in-memory `GameState` size (the cost of one
  undo snapshot), the number of implemented cards, and — after building a fresh
  wasm — the shipped web bundle's size raw and compressed (brotli and gzip, the
  levels `WebAssets` ships), split into the WASM bundle, the other assets, and
  the total.
- `mage docs` — serve this module's Go documentation at
  `http://localhost:6060` (pkgsite, the pkg.go.dev renderer); read-only, no
  gate depends on it.

## Generated configs and ignore files

Tool configs are generated, never hand-edited: `.golangci.yaml`,
`.markdownlint-cli2.yaml`, `.commitlint.yaml`, `.gitleaks.toml` and
`.ruleguard.go` are project-standards' base configs merged with the overrides
in `mklv.config.json`. Change an exclusion or a lint setting there, then run
`mage ci:fix` to regenerate; `ci:check` fails on a hand edit (drift). The same
goes for `.gitignore` and `.dockerignore`: add entries under `ignore.git` /
`ignore.docker` in `mklv.config.json`, because the weekly conformance run
rewrites both files from the universal templates plus those additions.

## Writing a mage target

The rules for adding or changing a target live in
[magefiles/AGENTS.md](../magefiles/AGENTS.md).
