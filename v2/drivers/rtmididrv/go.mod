module gitlab.com/gomidi/midi/v2/drivers/rtmididrv

go 1.16

require (
	gitlab.com/gomidi/midi/v2 v2.0.0-alpha.15
	gitlab.com/gomidi/midi/v2/drivers/rtmididrv/imported/rtmidi v0.0.0-20210425073027-dcb5d7eb9e83 // indirect
)

replace (
	gitlab.com/gomidi/midi/v2 => ../../
	gitlab.com/gomidi/midi/v2/drivers/rtmididrv/imported/rtmidi => ./imported/rtmidi
)
