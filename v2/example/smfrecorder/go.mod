module gitlab.com/gomidi/midi/v2/example/smfrecorder

go 1.14

require (
	gitlab.com/gomidi/midi/v2 v2.0.0-alpha.15
	gitlab.com/gomidi/midi/v2/drivers/rtmididrv v0.0.0-20210412062545-442b1d8545e9 // indirect
)

replace (
	gitlab.com/gomidi/midi/v2 => ../../
	gitlab.com/gomidi/midi/v2/drivers/rtmididrv => ../../drivers/rtmididrv
)
