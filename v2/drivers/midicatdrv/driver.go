package midicatdrv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"

	"gitlab.com/gomidi/midi/v2"
)

/*
func init() {
	drv, err := New()
	if err != nil {
		panic(fmt.Sprintf("could not register midicatdrv: %s", err.Error()))
	}
	midi.RegisterDriver(drv)
}
*/

type Driver struct {
	opened []midi.Port
	sync.RWMutex
}

func (d *Driver) String() string {
	return "midicatdrv"
}

// Close closes all open ports. It must be called at the end of a session.
func (d *Driver) Close() (err error) {
	d.Lock()
	var e CloseErrors

	for _, p := range d.opened {
		err = p.Close()
		if err != nil {
			e = append(e, err)
		}
	}

	d.Unlock()

	if len(e) == 0 {
		return nil
	}

	return e
}

// New returns a driver based on the default rtmidi in and out
func New() (*Driver, error) {
	return &Driver{}, nil
}

// Ins returns the available MIDI input ports
func (d *Driver) Ins() (ins []midi.In, err error) {

	c := midiCatCmd("ins --json")
	res, err := c.Output()

	if err != nil {
		return nil, fmt.Errorf("can't get in ports from midicat: %s", err.Error())
	}

	dec := json.NewDecoder(bytes.NewReader(res))

	var insjs = map[string]string{}

	err = dec.Decode(&insjs)

	if err != nil {
		return nil, fmt.Errorf("got invalid json from midicat: %s", string(res))
	}

	for idx, name := range insjs {
		i, err := strconv.Atoi(idx)
		if err != nil {
			return nil, fmt.Errorf("got invalid index from midicat: %s", string(res))
		}
		ins = append(ins, newIn(d, i, name))
	}

	inss := inPorts(ins)
	sort.Sort(inss)

	return inss, nil
}

// Outs returns the available MIDI output ports
func (d *Driver) Outs() (outs []midi.Out, err error) {
	c := midiCatCmd("outs --json")
	res, err := c.Output()

	if err != nil {
		return nil, fmt.Errorf("can't get out ports from midicat: %s", err.Error())
	}

	dec := json.NewDecoder(bytes.NewReader(res))

	var outsjs = map[string]string{}

	err = dec.Decode(&outsjs)

	if err != nil {
		return nil, fmt.Errorf("got invalid json from midicat: %s", string(res))
	}

	for idx, name := range outsjs {
		i, err := strconv.Atoi(idx)
		if err != nil {
			return nil, fmt.Errorf("got invalid index from midicat: %s", string(res))
		}
		outs = append(outs, newOut(d, i, name))
	}

	outss := outPorts(outs)
	sort.Sort(outss)

	return outss, nil

}

// CloseErrors collects error from closing multiple MIDI ports
type CloseErrors []error

func (c CloseErrors) Error() string {
	if len(c) == 0 {
		return "no errors"
	}

	var bd strings.Builder

	bd.WriteString("the following closing errors occured:\n")

	for _, e := range c {
		bd.WriteString(e.Error() + "\n")
	}

	return bd.String()
}
