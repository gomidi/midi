module gitlab.com/gomidi/midi/v2/drivers/portmididrv

go 1.14

require (
	github.com/rakyll/portmidi v0.0.0-20170620004031-e434d7284291
	gitlab.com/gomidi/midi/v2 v2.0.0-alpha.8
)

replace (
	gitlab.com/gomidi/midi/v2 => ../../
)