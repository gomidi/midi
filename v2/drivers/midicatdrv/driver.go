package midicatdrv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	"gitlab.com/gomidi/midi/v2/drivers"
	//	"gitlab.com/metakeule/config"
)

func init() {
	drv, err := New()
	if err != nil {
		panic(fmt.Sprintf("could not register midicatdrv: %s", err.Error()))
	}
	drivers.Register(drv)
}

type Driver struct {
	opened []drivers.Port
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

	d.opened = nil

	d.Unlock()

	if len(e) == 0 {
		return nil
	}

	return e
}

// const midicatVersion = "0.3.6"
// const midicatVersion = "0.4.0"
const midicatVersion = "0.6.3"
const midicatDownloadURL = "https://gitlab.com/gomidi/midi/-/releases/v2.1.1 midicat binaries (v0.6.3)"

func barkTo(wr io.Writer) {
	fmt.Fprintf(wr, "can't find midicat binary version >= %s in your PATH, please download from: %s\n", midicatVersion, midicatDownloadURL)
}

/*
func isVersionLess(a, b *config.Version) bool {
	if a.Major != b.Major {
		return a.Major < b.Major
	}

	if a.Minor != b.Minor {
		return a.Minor < b.Minor
	}

	return a.Patch < b.Patch
}

func CheckMIDICatBinary(barkTarget io.Writer) error {
	//
	b, err := midiCatVersionCmd().Output()

	if err != nil {
		if barkTarget != nil {
			barkTo(barkTarget)
		}
		return fmt.Errorf("missing binary 'midicat'")
	}

	s := string(b)

	idx := strings.LastIndex(s, " ")

	if idx < 4 {
		return fmt.Errorf("wrong version string of 'midicat'")
	}

	ver := strings.TrimSpace(s[idx:])

	vsGot, err := config.ParseVersion(ver)

	if err != nil {
		if barkTarget != nil {
			barkTo(barkTarget)
		}
		return fmt.Errorf("wrong version string of 'midicat'")
	}

	vsWant, err2 := config.ParseVersion(midicatVersion)

	if err2 != nil {
		panic(err2.Error())
	}

	if isVersionLess(vsGot, vsWant) {
		if barkTarget != nil {
			barkTo(barkTarget)
		}
		return fmt.Errorf("wrong version of 'midicat' (got %s required >= %s)", ver, midicatVersion)
	}

	return nil
}
*/

func checkMIDICAT() bool {
	b, err := midiCatVersionCmd().Output()

	if err != nil {
		barkTo(os.Stdout)
		panic("missing binary 'midicat'")
	}

	s := string(b)

	if s != midicatVersion {
		barkTo(os.Stdout)
		panic(fmt.Sprintf("%q", s))
	}
	return true
}

// New returns a driver based on the default rtmidi in and out
func New() (*Driver, error) {
	/*
		err := CheckMIDICatBinary(nil)
		if err != nil {

			return nil, err
		}
	*/
	checkMIDICAT()
	return &Driver{}, nil
}

// Ins returns the available MIDI input ports
func (d *Driver) Ins() (ins []drivers.In, err error) {

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
func (d *Driver) Outs() (outs []drivers.Out, err error) {
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
