# Vactrol Rulebook

Vactrol is a two-player card game — a clone of **KeyForge** — in which rival
players race to unlock a shared vault. On each of your turns you channel a single
**house**, marshal its creatures and artifacts, and gather **Aember**: the raw
resource that forges keys. The first player to forge three keys wins.

This rulebook is generated directly from the game engine's source. Every keyword,
ability, and effect below is described by the doc comment on the code that
implements it, so what you read here is always what the game actually does.

## Core concepts

- **Aember** is the resource players collect. It sits in a player's pool until it
  is spent forging a key, stolen by the opponent, or captured onto a creature.
- **Keys** are the win condition. At the start of your turn, if you have at least
  six Aember you forge a key, spending that Aember. Forge three keys to win.
- **Houses** are the factions a card belongs to. Your deck draws on three houses,
  but each turn you may act with the cards of only one **active house** you choose
  that turn.
- **Cards** come in four types — creatures, tactics, artifacts, and upgrades (see
  _Card Types_).
- **Exhaustion** gates a creature's activity. Using a creature to reap, fight, or
  take an action exhausts it; creatures and artifacts also enter play exhausted.
  Everything you control readies at the end of your turn.
