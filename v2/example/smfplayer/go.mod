module gitlab.com/gomidi/midi/v2/example/smfplayer

go 1.18

require (
	gitlab.com/gomidi/midi/v2 v2.0.0-beta.5
	gitlab.com/gomidi/midi/v2/drivers/rtmididrv v2.0.0-beta.5 // indirect
	gitlab.com/gomidi/midi/v2/drivers/portmididrv v2.0.0-beta.5
)

replace (
	gitlab.com/gomidi/midi/v2 => ../../
	gitlab.com/gomidi/midi/v2/drivers/rtmididrv => ../../drivers/rtmididrv
	gitlab.com/gomidi/midi/v2/drivers/portmididrv => ../../drivers/portmididrv
)
