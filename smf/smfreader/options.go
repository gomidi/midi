package smfreader

// Option is an option for the Reader
type Option func(*reader)

func debug(l logger) Option {
	return func(r *reader) {
		r.logger = l
	}
}

// NoteOffVelocity lets the reader differentiate between "fake" noteoff messages
// (which are in fact noteon messages (typ 9) with velocity of 0) and "real" noteoff messages (typ 8)
// that have a velocity.
// The former are returned as NoteOffVelocity messages and keep the given velocity, the later
// are returned as NoteOff messages without velocity. That means in order to get all noteoff messages,
// there must be checks for NoteOff and NoteOffVelocity (if this option is set).
// If this option is not set, both kinds are returned as NoteOff (default).
func NoteOffVelocity() Option {
	return func(rd *reader) {
		rd.readNoteOffPedantic = true
	}
}

type logger interface {
	Printf(format string, vals ...interface{})
}

/*
func FailOnUnknownChunks() Option {
	return func(r *reader) {
		r.failOnUnknownChunks = true
	}
}

// PostHeader tells the reader that next read is after the smf header
// remainingtracks are the number of tracks that are going to be parsed (must be > 0)
func PostHeader(remainingtracks uint16) Option {
	if remainingtracks == 0 {
		panic("remainingtracks must be at least 1")
	}
	return func(r *reader) {
		r.NumTracks = remainingtracks
		r.state = stateExpectChunk
		r.headerIsRead = true
	}
}

// InsideTrack tells the reader that next read is inside a track (after the track header)
// remainingtracks are the number of tracks that are going to be parsed (must be > 0)
func InsideTrack(remainingtracks uint16) Option {
	if remainingtracks == 0 {
		panic("remainingtracks must be at least 1")
	}
	return func(r *reader) {
		r.NumTracks = remainingtracks
		r.state = stateExpectTrackEvent
		r.headerIsRead = true
	}
}
*/
