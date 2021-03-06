package midi

import (
	"bytes"
	"fmt"

	"gitlab.com/gomidi/midi/v2/internal/utils"
)

func (m Message) BPM() float64 {
	if m.Type.IsNot(MetaTempoMsg) {
		fmt.Println("not tempo message")
		return -1
	}

	rd := bytes.NewReader(m.metaDataWithoutVarlength())
	microsecondsPerCrotchet, err := utils.ReadUint24(rd)
	if err != nil {
		fmt.Println("cant read")
		return -1
	}

	return float64(60000000) / float64(microsecondsPerCrotchet)
}
