# Vex

Vex is a rules engine for a [KeyForge](https://keyforge.fandom.com/)-style
card game, written in Go. It models a full two-player match — houses, Æmber,
keys, combat, and card abilities — behind a small, pointerless state type that is
cheap to copy (for AI search), and ships with a WebAssembly client for playing
in the browser.

## Quick start

```sh
mage web        # build the wasm client and serve it at http://localhost:8000
mage ci:test    # run the test suite
mage ci:cover   # coverage of the gated areas (kept at 100%)
```

Requires Go 1.27+.

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

`mage -l` lists the available targets. The generic gate comes from
[`dmikalova/project-standards`](https://github.com/dmikalova/project-standards)'
shared `ci` targets (`ci:fix`, `ci:check`, `ci:build`, `ci:test`, `ci:cover`,
`ci:lint` and the rest); vex adds its own (`web`, `webWasm`, `webAssets`, `gen`,
`fuzz`, `soak`, `debug`, `trace`, `profile` and the `tool:` namespace). The
local validator is `mage ci:fix && mage ci:check`: `ci:fix` applies every
autofix, and `ci:check` verifies without writing (format, tidy, build including
js/wasm, vet, lint, markdown, spelling, secrets, commit messages, generated
config drift, tests and the 100% coverage gates). `ci:check` is what CI runs.

Tool configs (`.golangci.yaml`, `.markdownlint-cli2.yaml`, `.commitlint.yaml`,
`.gitleaks.toml`, `.ruleguard.go`, `.gitignore`, `.dockerignore`) are generated
from project-standards' base configs and the overrides in
[`mklv.config.json`](mklv.config.json). Edit the overrides there, never the
generated files.

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
[`dmikalova/project-standards`](https://github.com/dmikalova/project-standards);
pre-commit runs `mage ci:check`. Every checker runs through `go run`, so the
only tools to install are `lefthook` and `mage`; then run `lefthook install`.
Commits follow [Conventional Commits](https://www.conventionalcommits.org/)
(enforced by `ci:commits`), which also drive versioning on deploy.

## Deployment

Vex runs on Google Cloud Run at
[vex.mklv.tech](https://vex.mklv.tech), served by the `cmd/web` binary
(the native build serves the WebAssembly client). The container listens on
`$PORT`, which Cloud Run injects.

CI/CD is a thin caller in
[`.github/workflows/cicd.yaml`](.github/workflows/cicd.yaml) that invokes the
reusable `cicd.yaml` workflow in
[`dmikalova/project-standards`](https://github.com/dmikalova/project-standards). On a push to
`main` it runs `mage ci:check`, then cuts a version from the Conventional
Commits, builds [`Dockerfile`](Dockerfile) into an image, and deploys it to
Cloud Run. Pull requests run `mage ci:check` only, and a weekly scheduled run
checks the project without deploying. The infrastructure itself — the Cloud Run service, the `vex.mklv.tech`
domain mapping, DNS, and CI deploy permissions — is defined as a Terramate stack
in the [`infrastructure`](https://github.com/dmikalova/infrastructure) repo under
`gcp/apps/vex`.

## License

See [LICENSE](LICENSE).
