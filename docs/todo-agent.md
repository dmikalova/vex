# Agent scratchpad

Working notes the agent writes for itself, to translate a request into concrete
work and show what is left. This is **not** [todo.md](todo.md) — that is the
human's personal list, which agents never write into. Rules:

- When an item is **done, delete it** — do not mark it done. This file only ever
  shows outstanding work, so it reads as a live "here is what I still mean to do"
  surface for coordinating with the human.
- Keep items concrete and **grouped by area or mechanic**, so related work is
  built together: add the shared primitive once, then knock out the group.
- Cite the ADR or doc that decided an item where one exists.

## Flank / placement prompts (engine)

- **Exhume must ask for a flank, not silently pick right.** It still auto-places
  on the right flank. Introduce an unset/invalid flank zero value so an
  unspecified flank errors by default; a card that dictates a flank (e.g. Amasser)
  sets it explicitly.

## Card behavior fixes

_(All of the first batch are done: Triumph "6 or more", City-State Interest
capture-from-opponent, Medic Ingram double-prompt, dying-creature Æmber
recipient, Poltergeist constant-ability artifacts. Remaining card-behavior work
lives under the other headings.)_

## Card text / glyphs

## Logs (engine)

- Poison and Skirmish are not logged.
- Æmber gains should name their source (e.g. "Harmonia …").

## Web — zone views, prompts, action bar

- Upgrade-attach prompt wording: change "Choose a creature to attach Stunner to"
  to "Choose a creature to attach Stunner **onto**". BLOCKED web-side: the string is
  produced by the engine (`AttachSelfTo.Text` / `game_play.go`'s
  "Choose a creature to attach {self} to" prompt), and engine AGENTS notes the
  prompt is deliberately kept identical to the printed card text. The mobile
  single-line clickable-green-name presentation is done (the prompt source card's
  face is dropped on mobile, leaving its green `prompt-source-name` token).
