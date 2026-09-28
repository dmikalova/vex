# AGENTS.md — `magefiles`

The target reference is [docs/building.md](../docs/building.md).

- `mage` lists each target with the **first sentence** of its doc comment: keep
  that sentence to one 80-column line (roughly 55 characters after mage's
  indent) and put every detail in the sentences after it.
- A target that writes a binary writes it to an ignored location (for example
  `./bin`), or its output goes under `ignore.git` in `mklv.config.json`. Never
  commit a generated binary.
