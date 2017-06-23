package smfreader

type Option func(*reader)

func Debug(l Logger) Option {
	return func(r *reader) {
		r.logger = l
	}
}

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
		r.mthd.numTracks = remainingtracks
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
		r.mthd.numTracks = remainingtracks
		r.state = stateExpectTrackEvent
		r.headerIsRead = true
	}
}
