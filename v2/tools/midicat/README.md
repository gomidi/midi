# midicat

Download the binaries for Linux and Windows [here](https://github.com/gomidi/midicat/releases/download/v0.3.6/midicat-binaries.zip).

When using windows, run the commands inside `cmd.exe`.

## Usage / Examples

get the list of MIDI in ports

    midicat ins
    
get the list of MIDI out ports

    midicat outs
    
pass the MIDI data from in port 11 to out port 12 

    midicat in -i=11 | midicat out -i=12
    
pass the MIDI data from in port 11 to out port 12 while logging it to stderr

    midicat in -i=11 | midicat log | midicat out -i=12
    
log the input from MIDI in port 11 without printing the raw MIDI bytes to stdout

    midicat in -i=11 | midicat log --nopass

get help

    midicat help
