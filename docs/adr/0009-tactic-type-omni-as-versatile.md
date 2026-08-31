# Card type "Tactic"; Omni authored as Versatile + Action

## Context

KeyForge overloads one word: it has an **action card** (a one-shot card type) and
an **"Action:"** ability (used from a card already in play). Naming the type
"Action" too makes "an Action" ambiguous in code, card text, and tests. KeyForge
also has an **"Omni:"** ability, usable on any turn rather than only when its house
is active.

## Decision

Name the one-shot card type **`Tactic`**, reserving "Action" for the ability
(`Trigger.Action`) exclusively. Author an **Omni** ability as the **`Versatile`**
keyword plus a `Trigger.Action` ability — the engine has no Omni trigger.

## Consequences

- "Tactic" (the type) and "Action" (the ability) are unambiguous everywhere.
- One fewer trigger to model: Omni behavior falls out of `Versatile` ("use as
  though it were your active house") composed with an Action ability.
- A deliberate vernacular divergence from KeyForge; it is recorded in `CONTEXT.md`
  and enforced by the card-wording rules, so a reader who expects "Action" the type
  or an "Omni" trigger finds the reason here rather than re-litigating it.
