# smfimage

Converts a Standard MIDI File (SMF) file to an image (PNG), based on keys and timing; separated by track.

## Status

usable (experimental)

## Installation of the CLI tool

     go install gitlab.com/gomidi/midi/tools/smfimage/cmd/smfimage

## Usage 

The following command creates the file `my-midi-file.png` reflecting the notes of `midi-file.mid`

     smfimage my-midi-file.mid

Colors are set based on the interval to the basenote (defaults to "C"). 
You can pass the base note to the command line via `-b`.

For more options

     smfimage help

The following color mapping is used

    prime/octave  -> yellow
    minor second  -> mint
    major second  -> orange
    minor third   -> sky blue
    major third   -> spring green
    fourth        -> cyan
    tritone       -> lime/chartreuse
    fifth         -> royalblue
    minor sixth   -> pink
    major sixth   -> violett/purple
    minor seventh -> magenta
    major seventh -> red

If you want to define own mappings, use the library and create a ColorMapper.


## Examples

color note mapping

![example image](https://gitlab.com/gomidi/midi/-/raw/master/tools/smfimage/example.png)

monochrome mode: velocity mapping: quiet (red) <-> loud (yellow) 

![example image](https://gitlab.com/gomidi/midi/-/raw/master/tools/smfimage/example2.png)

## Documentation

smfimage can also be used as library, see https://pkg.go.dev/gitlab.com/gomidi/midi/tools/smfimage