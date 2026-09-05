# Under: a per-host, out-of-play, facedown-capable attachment

## Context

Masterplan puts a card from hand facedown under itself, then later plays it and
destroys itself; Jargogle does the same shape of thing; Graft puts a card
faceup under itself and, per its own printed rule, sends that card to its
owner's discard pile if Graft is destroyed. All three need a per-host
collection of cards that is unbounded (any number of cards can end up under one
host), order-preserving, and — unlike an attached Upgrade — **out of play**: a
card placed under a host is not on the battleline or in the artifact row, it
does not fight or reap, and its owner cannot look at it unless the host's
controller allows it (a facedown under-card is exactly as hidden as a card in
hand).

`GameState` is flat and pointerless (ADR 0005), so an unbounded per-card
collection cannot be a slice; ADR 0001 already solved this shape once for
Upgrades with an intrusive singly-linked list threaded through `CardCore`. The
question this ADR answers is whether Under is a second instance of that same
pattern or needs something new, and how to model the one genuinely new rule it
introduces: controller-gated visibility of a facedown card sitting outside
every other zone.

## Decision

Under reuses ADR 0001's intrusive-list idiom exactly, with its own three
`CardCore` fields — `FirstUnderPlus`, `NextUnderPlus`, `UnderHostPlus` — plus one
field the Upgrade chain has no need of: `UnderFaceDown`, a bool recording
whether this particular card is hidden while buried. The bookkeeping
(`AttachUnder`, `detachUnder`, `firstUnder`/`nextUnder`, `underHostOf`) lives in
`game_under.go`, mirroring `game_upgrades.go` function-for-function.

Two differences from Upgrades follow from Under's own rules rather than from
the storage mechanism:

- **No type restriction.** An Upgrade's chain only ever holds cards of type
  Upgrade; Masterplan and Jargogle bury "a card from your hand" unrestricted by
  type, so `invariants.go`'s conservation walk checks Under's back-links but
  never a card's `Type`.
- **`Resolver` exposes Under, unlike Upgrades.** No card ability ever reads or
  moves an attached Upgrade directly — only the internal play/attach path
  touches that chain. Under-cards, by contrast, are themselves the target of
  later card text ("play the card under {self}"), so `ZoneReader`/
  `ZoneResolver` (ADR 0008) gained `Under`, `PlayFromUnder`, and `PutCardUnder`
  for effects to call through `EffectContext`.

Visibility is a single query, `Peekable(viewer, host) bool`, defined as "the
viewer controls the host" — the master rulebook's FACEDOWN CARDS rule applied
to this zone. It is a pure read with no separate state to keep in sync: a
facedown under-card's owner may still not look, only the host's controller may,
matching a card in a hand belonging to its owner but visible to no one else by
default.

When a host leaves play, every card placed under it goes to **its own owner's**
discard pile (`discardUnder`, called from every leave-play path alongside
`discardUpgrades`). Graft's printed text states this explicitly; Masterplan and
Jargogle do not repeat it, so this generalizes Graft's rule as the shared
default for the mechanic rather than restating it once per card.

The engine names this field `UnderFaceDown` rather than a bare `FaceDown` on
purpose: KeyForge has a second, unrelated facedown mechanic — facedown token
creatures left in play (Winds of Exchange) — which this ADR does not implement
and which will need its own name (and likely its own bool) when it is built, so
the two are not free to collide on one flag.

## Consequences

- A third and a fourth intrusive chain (after Upgrades) confirm ADR 0001
  generalizes without a new storage ADR; only the new domain rule (Under itself,
  and `Peekable`) needed deciding here.
- `discardUnder`'s "always to owner's discard" is an inferred default beyond the
  letter of the only card that states it (Graft) — if a future card buries a
  card and specifies a different fate for it on destruction, that card's
  ability, not this default, must say so explicitly.
- Facedown token creatures remain unimplemented; `UnderFaceDown`'s name leaves
  that mechanic a clean name of its own rather than forcing a later rename.
