package smftrack

import (
	"bytes"
	"fmt"
	"github.com/gomidi/midi"
	"github.com/gomidi/midi/live/midireader"
	"github.com/gomidi/midi/messages/meta"
	"io"
	"time"

	"github.com/gomidi/midi/smf"
)

// NewRecorder records into a single track on target, it needs the tempo to calculate the ticks from seconds
func NewRecorder(target smf.Writer, tempo uint) (io.Writer, error) {
	_, err := target.WriteHeader()

	if err != nil {
		return nil, err
	}

	if target.Header().Format != smf.SMF0 {
		return nil, fmt.Errorf("only SMF0 files supported, sorry, please convert your file")
	}

	ti, isMetric := target.Header().TimeFormat.(smf.MetricTicks)
	if !isMetric {
		return nil, fmt.Errorf("only metric timeformat supported, sorry")
	}

	_, err = target.Write(meta.Tempo(tempo))
	if err != nil {
		return nil, err
	}

	p := &recorder{}
	p.metricTicks = ti
	p.to = target
	p.from = midireader.New(&p.bf, nil)
	p.start = time.Now()
	return p, nil
}

// each write does in fact write to the midi.Writer passed to new
func (p *recorder) Write(data []byte) (n int, err error) {

	_, err = p.bf.Write(data)
	if err != nil {
		return
	}

	var msg midi.Message
	msg, err = p.from.Read()

	if err != nil {
		return
	}

	p.bf.Reset()

	now := time.Now()
	d := now.Sub(p.start)
	p.start = now
	p.to.SetDelta(p.metricTicks.TempoTicks(p.tempo, d))

	return p.to.Write(msg)
}

type recorder struct {
	metricTicks smf.MetricTicks
	bf          bytes.Buffer
	from        midi.Reader
	tempo       uint32
	to          smf.Writer
	start       time.Time
}
