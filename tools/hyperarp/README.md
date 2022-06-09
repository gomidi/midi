# hyperarp

A special breed of MIDI arpeggiator.

Note: If you are reading this on Github, please note that the repo is located at Gitlab (gitlab.com/gomidi/hyperarp) and this is only a mirror.

- Go version: >= 1.14

## Installation

    go install gitlab.com/gomidi/hyperarp
    
or download the binaries at https://github.com/gomidi/hyperarp/releases/tag/v0.0.18

## Usage

    hyperarp help
    
returns
	
	usage: 
	  hyperarp [command] OPTION... 
	
	options:
	  [--ctrlch=<integer>]          channel for control messages (only needed if 
	                                not the same as the input channel 
	
	  -i, --in=<integer>            number of the input MIDI port (use hyperarp 
	                                list to see the available MIDI ports) 
	
	  [--ccdir=80]                  controller number for the direction switch 
	
	  [--ccstyle=17]                controller number to select the playing style 
	                                (staccato, non-legato, legato) 
	
	  [--notetiming=<integer>]      note (key) for the timing interval 
	
	  [--notestyle=<integer>]       note (key) for the playing style (staccato, 
	                                non-legato, legato) 
	
	  -o, --out=<integer>           number of the output MIDI port (use hyperarp 
	                                list to see the available MIDI ports) 
	
	  [-t, --transpose=0]           transpose (number of semitones) 
	
	  [-b, --tempo=120]             tempo (BPM) 
	
	  [--cctiming=16]               controller number to set the timing interval 
	
	  [--notedir=<integer>]         note (key) for the direction switch 
	
	  [--config-locations]          prints the locations of current configuration 
	
	  [--config-files]              prints the locations of the config files 
	
	  [--version]                   prints the current version of the program 
	
	  [--help]                      prints the help 
	
	  [--config-spec]               prints the specification of the configurable 
	                                options 
	
	  [--config-env]                prints the environmental variables of the 
	                                configurable options 
	
	
	commands:
	  list                          show the available MIDI ports 
	
	
	for help about a specific command, run 
	  hyperarp help <command>

