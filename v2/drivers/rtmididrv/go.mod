module gitlab.com/gomidi/midi/v2/drivers/rtmididrv

go 1.14

require (
	gitlab.com/gomidi/midi/v2 v2.0.0-alpha.8
)

replace (
	gitlab.com/gomidi/midi/v2 => ../../
	gitlab.com/gomidi/midi/v2/drivers/rtmididrv/imported/rtmidi => ./imported/rtmidi
)