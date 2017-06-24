package smfreader

type state int

const (
	// At the start of the MIDI file.
	// Expect SMF Header chunk.
	stateExpectHeader state = 0

	// Expect a chunk. Any kind of chunk. Except MThd.
	// But really, anything other than MTrk would be weird.
	stateExpectChunk state = 1

	// We're in a Track, expect a track midi.
	stateExpectTrackEvent state = 2

	// This has to happen sooner or later.
	stateDone state = 3
)
