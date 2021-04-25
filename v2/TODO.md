# make helper functions (methods on *Message) for the data of the syscommon messages

# make transparent running status to explicit status reader; make it the default in listener, let it start listening at the first explicit status

# make transparent running status writer

# improve sysex

# Pipelines / Builders


```go
package main

import (
	_ gitlab.com/gomidi/rtmididrv  // autoregisters driver in central midi.DRIVERS hash, like database/sql
	gitlab.com/gomidi/midi
)

func main() {
	midi.Ins()  // returns in ports
	in, err := midi.Listen("port-description").
	    Only(midi.Channel1Msg & midi.NoteMsg).
		Do(func (msg midi.Message, deltamicroSec int64) {
		  fmt.Printf("[%v] %s\n", deltamicroSec, msg)
	    })
	in.Close()	
}

```


```go
package main

import (
	_ gitlab.com/gomidi/rtmididrv  // autoregisters driver in central midi.DRIVERS hash, like database/sql
	gitlab.com/gomidi/midi
	gitlab.com/gomidi/midi/smf
)

func main() {
	sm, err := smf.ReadTracks("midifile.mid",3).
	  Only(midi.Channel1Msg & midi.NoteMsg).
	  Do(func (trackNo int, msg midi.Message, delta int64, deltamicroSec int64) {
		fmt.Printf("T%v [%v] %s\n", trackNo, delta, msg)
	  })
}

```


```go
package main

import (
	_ gitlab.com/gomidi/rtmididrv  // autoregisters driver in central midi.DRIVERS hash, like database/sql
	gitlab.com/gomidi/midi
	gitlab.com/gomidi/midi/smf
)

func main() {
	
	file := smf.New("record.mid")
	defer file.Close()
	
	// single track recording, for multitrack we would have to collect the messages first (separated by port / midi channel)
	// and the write them after the recording on the different tracks
	in, err := midi.Listen("port-description").
	    Only(midi.Channel1Msg & midi.NoteMsg).
		Do(func (msg midi.Message, deltamicroSec int64) {
		  delta := file.DeltaFromMicroSec(deltamicroSec)
		  file.Delta(delta)
		  file.Write(msg)
	    })
	in.Close()	
}

```

```go
package main

import (
	_ gitlab.com/gomidi/rtmididrv  // autoregisters driver in central midi.DRIVERS hash, like database/sql
	gitlab.com/gomidi/midi
	gitlab.com/gomidi/midi/smf
)

func main() {
	out := midi.OutByName("port-description")
	defer out.Close()	

    // single track playing, for multitrack we would have to collect the tracks events first and properly synchronize playback	
	sm, err := smf.ReadTracks("midifile.mid", 3).
	   Only(midi.Channel1Msg & midi.NoteMsg).
	   Do(func (trackNo int, msg midi.Message, delta int64, deltamicroSec int64) {
		time.Sleep(time.Microseconds(deltamicroSec))
		out.Write(msg)
	})
	// you may do something with the sm
}

```
