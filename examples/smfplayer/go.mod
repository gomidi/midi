module gitlab.com/gomidi/midi/examples/smfplayer

go 1.18

require (
	gitlab.com/gomidi/midi/v2 v2.0.20
	gitlab.com/gomidi/midi/portmididrv v0.0.0
)

replace (
	gitlab.com/gomidi/midi/v2 => ../../v2
	gitlab.com/gomidi/midi/portmididrv => ../../portmididrv
)