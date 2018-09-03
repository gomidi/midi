package midireader

// Option is a configuration option for a reader
type Option func(rd *reader)

// NoteOffVelocity is an option for the reader that lets it differentiate between "fake" noteoff messages
// (which are in fact noteon messages (typ 9) with velocity of 0) and "real" noteoff messages (typ 8)
// having their own velocity.
// The later are returned as NoteOffVelocity messages and keep the given velocity, the former
// are returned as NoteOff messages without velocity. That means in order to get all noteoff messages,
// there must be checks for NoteOff and NoteOffVelocity (if this option is set).
// If this option is not set, both kinds are returned as NoteOff (default).
func NoteOffVelocity() Option {
	return func(rd *reader) {
		rd.readNoteOffPedantic = true
	}
}
