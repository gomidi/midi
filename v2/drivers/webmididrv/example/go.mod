module gitlab.com/gomidi/midi/v2/drivers/webmididrv/example

go 1.16

require gitlab.com/gomidi/midi/v2 v2.0.0-alpha.9
require gitlab.com/gomidi/midi/v2/drivers/webmididrv v0.0.0



replace (
	gitlab.com/gomidi/midi/v2 => ../../../
	gitlab.com/gomidi/midi/v2/drivers/webmididrv => ../
	
)
