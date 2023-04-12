
# Drivers for gomidi

All drivers must obey the following rules. If not, it is considered a bug.

1. Implement the midi.Driver, midi.In and midi.Out interfaces
2. Autoregister via midi.RegisterDriver(drv) within init function
3. The String() method must return a unique and indentifying string
4. The driver expects multiple goroutines to use its writing and reading methods and locks accordingly (to be threadsafe).
5. The midi.In ports respects any running status bytes when converting to midi.Message(s)
6. The midi.Out ports accept running status bytes.
7. The New function may take optional driver specific arguments. These are all organized as functional optional arguments.
8. The midi.In port is responsible for buffering of sysex data. Only complete sysex data is passed to the listener.
9. The midi In port must check the ListenConfig and act accordingly.
10. incomplete sysex data must be cached inside the sender and flushed, if the data is complete.

## Reason for the timestamp choice of the listener callback

1. we don't want floating point, but integers of small fractions of time (easier to calculate with).
2. not every driver tracks a delta timing. the ones that don't should indicate with a -1. also some underlying drivers
   use floats for delta timing, so we can't be sure to don't get negative values.
3. for the size, we find that int32 is large enough, if we take a reasonable resolution 
   of 1 millisecond. Then we get
     max int32  = 2147483647 / 10 (ms) / 1000 (sec) / 60 (min) / 60 (hours) / 24 = 2,48 days 
   of maximal duration when converting to absolute timing (starting from the first message), which should be 
   long enough for a midi recording.
    int64 would double the needed ressources for no real benefit.
