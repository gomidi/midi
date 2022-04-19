module gitlab.com/gomidi/midi/v2/example/smfplayer

go 1.18

require (
	gitlab.com/gomidi/midi/v2 v2.0.3
	gitlab.com/gomidi/midi/v2/drivers/rtmididrv v0.0.0-20210412062545-442b1d8545e9 // indirect
	gitlab.com/gomidi/midi/v2/drivers/portmididrv v0.0.0
	gitlab.com/gomidi/midi/v2/drivers/midicatdrv v0.0.0
)

replace (
	gitlab.com/gomidi/midi/v2 => ../../
	gitlab.com/gomidi/midi/v2/drivers/rtmididrv => ../../drivers/rtmididrv
	gitlab.com/gomidi/midi/v2/drivers/portmididrv => ../../drivers/portmididrv
	gitlab.com/gomidi/midi/v2/drivers/midicatdrv => ../../drivers/midicatdrv
)
