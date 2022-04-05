# midispy
Spy on MIDI data transmitted between a sending and a receiving device

## Usage (CLI)

    go install gitlab.com/gomidi/v2/tools/midispy/cmd/midispy

To get a list of available MIDI devices, run

    midispy list

Example of output:

    ---MIDI input ports---
    0 "Midi Through:Midi Through Port-0 14:0"
    1 "Virtual Raw MIDI 2-0:VirMIDI 2-0 24:0"
    2 "MPKmini2:MPKmini2 MIDI 1 28:0"

    ---MIDI output ports---
    0 "Midi Through:Midi Through Port-0 14:0"
    1 "Virtual Raw MIDI 2-0:VirMIDI 2-0 24:0"
    2 "MPKmini2:MPKmini2 MIDI 1 28:0"
    3 "FLUID Synth (qsynth):Synth input port (qsynth:0) 128:0"

Then you can use the given ids to tell midispy to listen:

    midispy in=10 out=11
    [10] "VMPK Output:VMPK Output 128:0"
    ->
    [11] "FLUID Synth (qsynth):Synth input port (qsynth:0) 130:0"
    -----------------------
    [10->11] 13:38:00.152725 channel.NoteOn channel 0 key 58 velocity 100
    [10->11] 13:38:00.265286 channel.NoteOff channel 0 key 58
    [10->11] 13:38:01.180286 channel.NoteOn channel 0 key 58 velocity 100
    [10->11] 13:38:01.276850 channel.NoteOff channel 0 key 58
    [10->11] 13:38:01.701086 channel.NoteOn channel 0 key 71 velocity 100
    [10->11] 13:38:01.786206 channel.NoteOff channel 0 key 71
    [10->11] 13:38:01.990798 channel.NoteOn channel 0 key 69 velocity 100
    [10->11] 13:38:02.081213 channel.NoteOff channel 0 key 69
    [10->11] 13:38:02.255353 channel.NoteOn channel 0 key 58 velocity 100
    [10->11] 13:38:02.345848 channel.NoteOff channel 0 key 58
    [10->11] 13:38:32.193058 channel.NoteOn channel 0 key 58 velocity 100
    [10->11] 13:38:33.215804 channel.NoteOn channel 0 key 74 velocity 100
    [10->11] 13:38:34.087821 channel.NoteOff channel 0 key 74
    [10->11] 13:38:34.580737 channel.NoteOff channel 0 key 58


To get help:

    midispy help   

