# Vactrol

Vactrol is a rules engine for a [KeyForge](https://keyforge.fandom.com/)-style
card game, written in Go. It models a full two-player match — houses, Æmber,
keys, combat, and card abilities — behind a small, pointerless state type that is
cheap to copy (for AI search), and ships with a WebAssembly client for playing
in the browser.

## Quick start

```sh
mage web      # build the wasm client and serve it at http://localhost:8000
mage test     # run the test suite
mage cover    # engine test coverage (kept at 100%)
```

Requires Go 1.26+.

## Layout

| Path                        | Package                 | Responsibility                                                                                    |
| --------------------------- | ----------------------- | ------------------------------------------------------------------------------------------------- |
| `internal/engine`           | `engine`                | Core rules engine: game state, combat, and the card-effect AST. Pointerless and clone-friendly.   |
| `internal/card`             | `card`                  | Authoring facade over the engine (grouped namespaces like `card.House.X`) plus the card registry. |
| `internal/cards`            | `cards`                 | Card-database aggregator; blank-imports every set so its cards self-register.                     |
| `internal/cards/sets/<set>` | e.g. `callofthearchons` | One self-registering file per card.                                                               |
| `internal/cards/cardtest`   | `cardtest`              | Shared test harness for the set packages.                                                         |
| `internal/match`            | `match`                 | Shared match setup (random decks, house list) used by every frontend.                             |
| `internal/web`              | `web`                   | [go-app](https://github.com/maxence-charriere/go-app) WebAssembly client, Monokai-themed.         |
| `cmd/web`                   | `main`                  | Serves the web client and, compiled to wasm, runs it in the browser.                              |

The web UI ships today; further frontends (e.g. an MCTS bot) and a lobby server
are planned as their own `cmd/…` binaries and `internal/…` packages on the same
engine.

For how the pieces fit together — the pointerless state, the card-effect AST, and
how a turn and an ability flow through the code — see
[`docs/architecture.md`](docs/architecture.md). For the testing options and what
to test at each layer, see [`docs/testing.md`](docs/testing.md). A full index of
every doc is in [`docs/README.md`](docs/README.md).

## Development

`mage -l` lists the available targets: `web`, `webWasm`, `build`, `test`,
`cover`, `vet`, `fmt`, `lint`, `check`, `tidy`, `gen`. `mage check` is the full
green gate (fmt-check, build, vet, lint, markdown lint, test, coverage) and is
what CI runs.

The card-management commands live under the `tools` namespace, for researching
and implementing cards:

```sh
mage tool:lookup "ether spider"   # find source cards by name, with a ready-made card.Provenance(...)
mage tool:missing         # pick a set (↑/↓), then list its cards still to implement
mage tool:coverage                 # per-set count of implemented cards
mage tool:stub callofthearchons    # scaffold build-excluded stubs for a set's unimplemented cards
```

Card-authoring conventions live in
[`internal/cards/AGENTS.md`](internal/cards/AGENTS.md).

Commit hooks are managed by [lefthook](https://github.com/evilmartians/lefthook),
extending the shared base config in
[`dmikalova/github-meta`](https://github.com/dmikalova/github-meta). Install the
tooling once (`lefthook`, `commitlint`, `gitleaks`, `typos`, `quickmark`, `mage`),
then run `lefthook install`. Commits follow
[Conventional Commits](https://www.conventionalcommits.org/) (enforced by
commitlint), which also drives semantic-release versioning on deploy.

## Deployment

Vactrol runs on Google Cloud Run at
[vactrol.mklv.tech](https://vactrol.mklv.tech), served by the `cmd/web` binary
(the native build serves the WebAssembly client). The container listens on
`$PORT`, which Cloud Run injects.

CI/CD is a thin caller in
[`.github/workflows/cicd.yaml`](.github/workflows/cicd.yaml) that invokes the
reusable `go-cloudrun.yaml` workflow in
[`dmikalova/github-meta`](https://github.com/dmikalova/github-meta). On a push to
`main` it runs `mage check`, then semantic-release cuts a version, buildah builds
[`Dockerfile`](Dockerfile) into an image, pushes it to GHCR (mirrored to Artifact
Registry), and `gcloud run deploy` rolls it out. Pull requests run `mage check`
only. The infrastructure itself — the Cloud Run service, the `vactrol.mklv.tech`
domain mapping, DNS, and CI deploy permissions — is defined as a Terramate stack
in the [`infrastructure`](https://github.com/dmikalova/infrastructure) repo under
`gcp/apps/vactrol`.

## License

See [LICENSE](LICENSE).
