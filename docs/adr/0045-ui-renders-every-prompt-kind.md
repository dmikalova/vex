# 45. The UI is proven to render every prompt kind

This decision records how the client is validated to handle every prompt the
engine can raise. It guides a staged change: steps 1 and 2 exist in code; step 3
is staged. It is the prompt-surface counterpart to ADR 0022 (iconography is a
Visitor guarded by a totality test) and ADR 0018 (the rulebook registry is
complete by construction).

## Context

The engine raises a prompt by calling an **optional-capability interface** on the
player's `Chooser` — `OptionChooser`, `PositionChooser`, `DeclinableChooser`,
`Orderer`, `ReactionChooser`, `BadgeChooser` — each discovered at runtime with
`if c, ok := chooser.(OptionChooser); ok { … } else { fallback }`
(`internal/engine/game.go`). The base `Chooser` requires only `ChooseCreature`;
every richer capability is structural and unasserted.

This design **fails open**, in three places:

- **The engine falls back.** A chooser that does not satisfy a capability takes the
  engine's default path. No panic.
- **The web dispatch is flag guards, not a switch.** `controls()`
  (`internal/web/view_controls.go`) is a chain of `if g.choosingX` guards, and
  `optionChooser()` sniffs label _shape_, ending in a benign catch-all that renders
  generic buttons. Nothing fails.
- **Signature drift is silent.** If a chooser method's signature changes, the
  chooser quietly stops satisfying the capability and downgrades to fallback, with
  no compile error.

So adding a new prompt route does not force the UI to handle it — the opposite
pressure from the two proven totality patterns already in the tree: the
iconography Visitor with `TestIconTotality` (ADR 0022, data-driven, no allowlist)
and the rulebook term registry (ADR 0018, enumerate a closed catalog and fail on
any undescribed member). Both make a new variant fail the build until it is
handled.

One constraint shapes the fix: **fallback is load-bearing.** `FirstChooser`, the
bot, and MCTS rollouts deliberately implement only the base `Chooser` and answer
every richer prompt with the engine default, rendering no UI. A totality rule must
target the **human-facing** chooser, not every `Chooser` implementation, or it
would force the bot to render prompts it never shows.

## Decision

Validate prompt rendering by **totality, not fallback**, in three staged steps.

**Step 1 — compile-time capability assertions (done).** Each chooser binds, with
`var _ engine.Capability = (*chooser)(nil)`, to the capability interfaces it is
intended to satisfy: `suspendChooser` — the one the web client's session answers
through — to all seven, `bridgeChooser` to the ones card tests reach,
`FirstChooser` to only the base `Chooser`. This catches
signature drift and accidental capability loss at compile time. It does **not**
catch a brand-new capability interface, because nothing forces a new assertion to
be written — that is step 2's job.

**Step 2 — a closed `PromptKind` registry with a web-scoped totality test
(done).** `internal/engine/prompt.go` holds `PromptKind`, an enumerable catalog
with one member per Chooser capability (`PromptCreature`, `PromptCardOrDecline`,
`PromptOption`, `PromptPosition`, `PromptReaction`, `PromptOrder`, `PromptBadge`),
exposed as the closed list `PromptKinds()`. In `internal/web`, `promptControls`
is the single dispatch seam `controls()` routes its live prompt branches through,
and `TestPromptTotality` asserts it returns a non-fallback rendering for every
`PromptKind` — so a new kind fails the test until it is rendered. The many card
kinds share one cluster; `PromptBadge` is display-only (drawn on the candidate,
no dock cluster). The silent catch-all in `optionChooser()` is gone: option
rendering now dispatches on a named, enumerable `optionKind` through
`optionControls`, guarded by `TestOptionKindTotality`, so a new option shape is a
named case rather than a fall-through. Label shape still derives the option kind —
the seam that survives until ADR 0040's `Request` carries the kind from the
engine.

This **converges with ADR 0040.** That ADR inverts `Chooser` into a
`Request → Command` answerer where `Request` is "a thin decision-point marker —
which player owes a decision, and in what context." `PromptKind` is the
enumerable form of that context: the same closed catalog serves the totality test
now and `Request`'s context tag after the step-function refactor. Build
`PromptKind` so it can become `Request`'s context, not a throwaway.

**Step 3 — scenario snapshot exhaustiveness (future).** Once the "load a test
situation and run it in the engine for Playwright" harness lands
([../todo.md](../todo.md)), drive the engine through one scenario per `PromptKind`
and assert each yields a rendered control cluster. This validates _presentation_,
not just dispatch — it catches a kind that dispatches but renders wrong. It reuses
the scenario harness rather than mocking prompts, and is deferred until that
harness exists.

## A prompt names its source by identity, not by name

Every asking capability takes a leading `PromptSource` (`internal/engine/prompt.go`)
that identifies the card whose ability raised the prompt by `LocalID`. A prompt
never identifies its source by name. A name cannot distinguish two copies of one
card, so a client that matches on it picks whichever copy it finds first and
renders that copy's live state — the wrong house, the wrong Maverick mark, the
wrong damage — for the copy that actually asked. `TestPromptSourceReadsTheAskingCopysHouse`
(`internal/web/view_controls_test.go`) pins this with two copies in play.

`LocalID` 0 is a real card, so presence rides on the paired `HasCard` boolean
rather than a sentinel id (ADR 0010). A prompt with no card source — a
turn-structure question — carries the zero `PromptSource`.

The source deliberately does **not** carry the card's name alongside its id.
`renderPrompt` substitutes `SelfName` inside the engine before the chooser is
called, so the prompt text is already rendered: a client uses the source to point
at the card, never to rebuild the sentence. New prompt context is a new field on
`PromptSource`, not another parameter on every chooser implementation.

## Consequences

- Step 1 hardens the silent-fallback risk immediately, at zero runtime cost, and
  documents each chooser's intended capability set at its definition.
- Step 2 reverses the fallback-versus-totality tradeoff for prompts, matching
  ADR 0022's stance for icons: a new mechanic cannot ship a prompt without a
  rendering decision. The `PromptKind` enum stays scoped to human-facing rendering;
  `FirstChooser` and the bot stay on fallback by design.
- Step 2's enum is forward-built as ADR 0040's `Request` context, so the two do
  not diverge when the engine becomes a suspendable step function.
- Step 3 depends on the scenario harness and waits for it.
