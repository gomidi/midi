# mid
Porcelain library for reading and writing MIDI and SMF (Standard MIDI File) 

## Description

Package mid provides an easy abstraction for reading and writing of "live" `MIDI` and `SMF` 
(Standard MIDI File) data.

`MIDI` data could be written the following ways:

- `NewWriter` is used to write "live" MIDI to an `io.Writer`.
- `NewSMF` is used to write SMF MIDI to an `io.Writer`.
- `NewSMFFile` is used to write a complete SMF file.
- `WriteTo` writes "live" MIDI to an `connect.Out`, aka MIDI out port

To read, create a `Reader` and attach callbacks to it.
Then MIDI data could be read the following ways:

- `Reader.Read` reads "live" MIDI from an `io.Reader`.
- `Reader.ReadSMF` reads SMF MIDI from an `io.Reader`.
- `Reader.ReadSMFFile` reads a complete SMF file.
- `Reader.ReadFrom` reads "live" MIDI from an `connect.In`, aka MIDI in port

For a simple example with "live" MIDI and `io.Reader` and `io.Writer` see the example below.

## Example

We use an `io.Writer` to write to and `io.Reader` to read from. They are connected by the same `io.Pipe`.

```go
package main

import (
    "fmt"
    "gitlab.com/gomidi/midi/mid"
    "io"
    "time"
)

// callback for note on messages
func noteOn(p *mid.Position, channel, key, vel uint8) {
    fmt.Printf("NoteOn (ch %v: key %v vel: %v)\n", channel, key, vel)
}

// callback for note off messages
func noteOff(p *mid.Position, channel, key, vel uint8) {
    fmt.Printf("NoteOff (ch %v: key %v)\n", channel, key)
}

func main() {
    fmt.Println()

    // to disable logging, pass mid.NoLogger() as option
    rd := mid.NewReader()

    // set the functions for the messages you are interested in
    rd.Msg.Channel.NoteOn = noteOn
    rd.Msg.Channel.NoteOff = noteOff

    // to allow reading and writing concurrently in this example
    // we need a pipe
    piperd, pipewr := io.Pipe()

    go func() {
        wr := mid.NewWriter(pipewr)
        wr.SetChannel(11) // sets the channel for the next messages
        wr.NoteOn(120, 50)
        time.Sleep(time.Second) // let the note ring for 1 sec
        wr.NoteOff(120)
        pipewr.Close() // finishes the writing
    }()

    for {
        if rd.Read(piperd) == io.EOF {
            piperd.Close() // finishes the reading
            break
        }
    }

    // Output:
    // channel.NoteOn channel 11 key 120 velocity 50
    // NoteOn (ch 11: key 120 vel: 50)
    // channel.NoteOff channel 11 key 120
    // NoteOff (ch 11: key 120)
}
```