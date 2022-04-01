module gitlab.com/gomidi/midi/v2/tools/midispy

go 1.16

require (
	gitlab.com/gomidi/midi/v2 v2.0.0-alpha.15
	gitlab.com/gomidi/midi/v2/drivers/rtmididrv v0.0.0-20210412062545-442b1d8545e9 // indirect
	gitlab.com/metakeule/config v1.21.0
)

replace (
	gitlab.com/gomidi/midi/v2 => ../../
	gitlab.com/gomidi/midi/v2/drivers/rtmididrv => ../../drivers/rtmididrv
)


