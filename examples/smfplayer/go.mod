module gitlab.com/gomidi/midi/examples/smfplayer

go 1.22.2

require gitlab.com/gomidi/midi/v2 v2.2.0

replace (
	gitlab.com/gomidi/midi/portmididrv => ../../portmididrv
	gitlab.com/gomidi/midi/v2 => ../../v2
)
