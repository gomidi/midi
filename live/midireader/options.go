package midireader

// Option is a configuration option for a reader
type Option func(rd *reader)

// NoteOffPedantic is an option for the reader that lets it differenciate between "fake" noteoff messages
// (which are in fact noteon messages (typ 9) with velocity of 0) and "real" noteoff messages (typ 8)
// The former are returned as NoteOffPedantic messages and keep the given velocity, the later
// are returned as NoteOff messages without velocity. That means in order to get all noteoff messages,
// there must be checks for NoteOff and NoteOffPedantic (if this option is set).
// If this option is not set, both kinds are returned as NoteOff (default).
func NoteOffPedantic() Option {
	return func(rd *reader) {
		rd.readNoteOffPedantic = true
	}
}
