module gitlab.com/gomidi/midi/v2/drivers/midicatdrv

go 1.16

require (
    gitlab.com/gomidi/midi/v2 v2.0.0-alpha.9
    gitlab.com/gomidi/midicat v0.3.6
    gitlab.com/metakeule/config v1.21.0
)

replace (
	gitlab.com/gomidi/midi/v2 => ../../
)