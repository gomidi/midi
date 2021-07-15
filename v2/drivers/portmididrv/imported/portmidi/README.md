# portmidi
Want to output to an MIDI device or listen your MIDI device as an input? This
package contains Go bindings for PortMidi. `libportmidi` (v. 217) is required as a dependency, it's available via apt-get and brew.

~~~ sh
apt-get install libportmidi-dev
# or
brew install portmidi
~~~

Or, alternatively you can download the source and build it by yourself. See
the instructions on [PortMidi homepage](http://portmedia.sourceforge.net/portmidi/).

In order to start, go get this repository:
~~~ sh
go get github.com/rakyll/portmidi
~~~

## Usage

### Initialize
~~~ go
portmidi.Initialize()
~~~

### About MIDI Devices

~~~ go
portmidi.CountDevices() // returns the number of MIDI devices
portmidi.Info(deviceID) // returns info about a MIDI device
portmidi.DefaultInputDeviceID() // returns the ID of the system default input
portmidi.DefaultOutputDeviceID() // returns the ID of the system default output
~~~

### Write to a MIDI Device

~~~ go
out, err := portmidi.NewOutputStream(deviceID, 1024, 0)
if err != nil {
    log.Fatal(err)
}

// note on events to play C major chord
out.WriteShort(0x90, 60, 100)
out.WriteShort(0x90, 64, 100)
out.WriteShort(0x90, 67, 100)

// notes will be sustained for 2 seconds
time.Sleep(2 * time.Second)

// note off events
out.WriteShort(0x80, 60, 100)
out.WriteShort(0x80, 64, 100)
out.WriteShort(0x80, 67, 100)

out.Close()
~~~

### Read from a MIDI Device
~~~ go
in, err := portmidi.NewInputStream(deviceID, 1024)
if err != nil {
    log.Fatal(err)
}
defer in.Close()

events, err := in.Read(1024)
if err != nil {
    log.Fatal(err)
}

// alternatively you can filter the input to listen
// only a particular set of channels
in.SetChannelMask(portmidi.Channel(1) | portmidi.Channel.(2))
in.Read(1024) // will retrieve events from channel 1 and 2

// or alternatively listen events
ch := in.Listen()
event := <-ch
~~~

### Cleanup
Cleanup your input and output streams once you're done. Likely to be called on graceful termination.
~~~ go
portmidi.Terminate()
~~~
