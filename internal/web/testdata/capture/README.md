# Replay captures

A JSON file here is one **open finding**: a match the browser client saved and
then could not replay, written out by the dev server's `POST /debug/capture`
endpoint (`mage web` turns it on; a deployed build has no such endpoint).

`TestCaptures` replays every file here and **fails** with the panic it recorded,
so a capture in a commit is a lapse rather than an archive entry — fix the fault
and run `mage capturePrune`, which deletes the findings that no longer reproduce
along with every one recorded against a different command-log version or card
pool.

A capture whose version or card-pool digest does not match the current tree is
**skipped with a note**, never failed: the same seed deals different cards from a
different pool, so it is not evidence about this tree.

See [docs/testing.md](../../../../docs/testing.md) and `internal/web/capture.go`.
