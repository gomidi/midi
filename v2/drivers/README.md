
# Drivers for gomidi

All drivers must obey the following rules. If not, it is considered a bug.

1. Implement the midi.Driver, midi.In and midi.Out interfaces
2. Autoregister via midi.RegisterDriver(drv) within init function
3. The String() method must return a unique and indentifying string
4. The driver expects multiple goroutines to use its writing and reading methods and locks accordingly (to be threadsafe).
5. The midi.In ports respects any running status bytes when converting to midi.Message(s)
6. The midi.Out ports may convert explicit status bytes to running status bytes.
7. The New function may take optional driver specific arguments. These are all organized as functional optional arguments.
8. The midi.In port is responsible for buffering of sysex data. Only complete sysex data is passed to the listener.
9. The midi In port must check, if the receiver Implements any of the optional interfaces for handling sysex data, syscommon
   data or realtime data and serve accordingly. Such data is only passed, if the Receiver implements the according interface.
   This is the way, the Receiver tells the driver, which data it is interested in.
10. incomplete sysex data must be cached inside the sender and flushed, if the data is complete.
11. The driver must have a pass through method which passes the data as is (for debugging).