module gitlab.com/gomidi/midi/v2/drivers/rtmididrv

go 1.16

require (
	gitlab.com/gomidi/midi/v2 v2.0.3
	gitlab.com/gomidi/midi/v2/drivers/rtmididrv/imported/rtmidi v0.9.0 // indirect
)

replace (
	gitlab.com/gomidi/midi/v2 => ../../
	gitlab.com/gomidi/midi/v2/drivers/rtmididrv/imported/rtmidi => ./imported/rtmidi
)
